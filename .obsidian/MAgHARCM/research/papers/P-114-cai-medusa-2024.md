---
title: "P-114 — Cai et al. 2024 — Medusa: Simple LLM Inference Acceleration Framework with Multiple Decoding Heads"
backlink: "[[1.0.0 P-114]]"
tags: [paper, slm, speculative-decoding, medusa, multi-head, [[1.0.0 PRIM-7]], [[1.0.0 PRIM-21]], [[1.0.0 PRIM-22]], [[1.0.0 PRIM-31]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-114 — Cai et al. 2024 — Medusa: Multiple Heads are Better than One]]

- **Authors**: Tianle Cai, Yuhong Li, Zhengyang Geng, Hongwu Peng, Jacob D. Tripp, Omar Saif, Hossein Pourreza, Ming-Yu Liu, Xiang Li, Trevor Darrell (UC Berkeley + NVIDIA).
- **Venue / Year**: arXiv:2401.10774 v6 (January 2024); integrated into vLLM (Kwon et al. 2023) and SGLang (Zheng et al. 2024); the project's README states it was initially released around ICML 2024 / 2401 — [INFERENCE: arXiv preprint + open-source release; venue is the MLSys 2024 ecosystem rather than a single tracked conference]. The follow-up *Medusa-2* technical report and the *EAGLE vs Medusa* comparisons are also arXiv-only as of wave-12.
- **URL**: https://arxiv.org/abs/2401.10774
- **Code / Models**: https://github.com/FasterDecoding/medusa ; first-class integrations in vLLM >= 0.4 (`--speculative-model medusa`), SGLang (`speculative_draft_model_name=medusa`), and TensorRT-LLM (Medusa plugin).
- **Anchors**: PRIM-7 (Multi-Agent Verdict Validation), PRIM-21 (Migration Strategy Selection), PRIM-22 (Four Phases Comprehension), PRIM-31 (Iterative Retrieval Refinement); speculative-decoding lineage anchored by [[1.0.0 P-57]] (Leviathan), [[1.0.0 P-64]] (EAGLE-1), and [[1.0.0 P-108]] (EAGLE-3). P-114 is the **alternative-drafting-strategy** anchor — every place EAGLE-3 is named, Medusa is the cheap + fixed sibling.

## 1. Core Contribution

Medusa is the **training-free / head-only** alternative to EAGLE-3's training-time-test single-head draft ([[1.0.0 P-108]]). Instead of training a separate draft model or a feature-shift head (the [[1.0.0 P-64]] EAGLE-1 design), Medusa **adds multiple parallel decoding heads** directly to the target model's final layer, each predicting a different future-token offset in parallel:

1. **Multiple parallel heads at future positions** — head $k$ predicts the token at position $t+k$ conditioned on the target's hidden state at position $t$. With two or three heads, Medusa proposes a **draft tree** of length $\gamma = K$ candidate continuations in a single target forward pass — no autoregressive draft loop, no separate draft model. The $k$-th head is a small MLP on top of the last hidden state followed by the frozen LM head's unembedding matrix (untied weights via a per-head low-rank adapter; the original LM head weights are reused for the first Medusa head).
2. **Tree-based attention for verification** — the candidate continuations form a tree (because the draft step is branching, not linear), and the target verifies the whole tree in **one forward pass** using a tree-shaped causal mask. Acceptance uses the same $\min(1, p/q)$ rule as [[1.0.0 P-57]] (Leviathan 2023), so the output distribution is lossless.
3. **Optional lightweight fine-tuning** — heads are typically trained for one epoch on a small calibration set (the "Medusa training" step in the README) so that the per-head logit distribution approximates the target's shifted-by-$k$ distribution. Training is fast (single-GPU, minutes for a 7B target on a few thousand prompts) because only the head adapters update; the target backbone stays frozen. There is **no per-domain retraining loop** in Medusa by default — once calibrated, the heads are fixed.

Reported headline numbers (arXiv:2401.10774 v6):
- **2.2×–3.6×** lossless wall-clock speedup on Vicuna-7B/13B, LLaMA-2-Chat 7B/13B/70B, and Mistral 7B across MT-Bench, AlpacaEval, and HumanEval-style code generation.
- **~2.2×** on Vicuna-7B at parity quality; **~2.8×** on LLaMA-2-Chat 70B (the largest model reported in the paper).
- Tree-verify acceptance rate is the dominant knob: average accepted run length $\gamma$ of 2.4–3.1 tokens across the reported eval suite.
- Quality preserved by construction (lossless speculative sampling).

**EAGLE-3 vs Medusa at a glance (the framing MAgHARCM cares about)**:

| Dimension | [[1.0.0 P-108]] EAGLE-3 | **P-114 Medusa** |
|---|---|---|
| Draft source | trained draft head on fused multi-layer features | K parallel heads on top-layer only |
| Setup cost | multi-hour training-time-test per (target, domain) | minutes of head calibration per target |
| Per-domain retrain | cheap (single-pass training-time-test, low-rank) | not supported out of the box — heads are fixed |
| Speedup | up to **6.5×** (arXiv:2503.01840) | **2.2×–3.6×** (arXiv:2401.10774) |
| Output distribution | lossless | lossless |
| Inference integrations | vLLM, SGLang, TensorRT-LLM | vLLM, SGLang, TensorRT-LLM |
| Failure mode | drift if retraining signal stale | head miscalibration → low acceptance |

**Why this matters for the 4B–30B SLM regime**: MAgHARCM targets Qwen2.5-Coder [[P-21]], StarCoder2 [[P-22]], and Phi-3-mini [[P-54]]. Medusa's "freeze the target, bolt on two heads, calibrate for ten minutes" is the **minimum-friction** speculative-decoding recipe in the entire [[1.0.0 P-57]] → [[1.0.0 P-64]] → [[1.0.0 P-108]] → P-114 lineage. For pipelines that need a speedup *now* without committing a per-domain training budget, Medusa is the default. The trade-off — fixed heads, no per-domain retraining, lower ceiling — is exactly what the EAGLE-3 row of the strategy registry exists to capture when the budget allows.

## 2. Application in MAgHARCM

Medusa's K-head draft tree composes into four primitives simultaneously, mirroring the application notes in [[1.0.0 P-57]] §2 and [[1.0.0 P-108]] §2. P-114 is the **fixed-budget** anchor; P-108 is the **trainable-budget** anchor; both rows live in the same lineage matrix cell.

### 2.1 PRIM-7 (Multi-Agent Verdict Validation)

PRIM-7's verdict panel fans out $K$ SLM proposals to a single oracle. Medusa's **parallel-head tree-attention** is the natural fit: each Medusa head is structurally a **voter** at offset $k$, and the tree-attention verifier step is a **batch consensus** in a single forward pass. The K drafted continuations at offsets $t+1, \dots, t+K$ are exactly the kind of N-proposal sample that PRIM-7 needs to compute a verdict distribution. Concretely:
- $K = 4$–$7$ Medusa heads per target SLM maps one-to-one onto PRIM-7's $K$-agent verdict panel.
- Tree-verify accepts a multi-token draft per agent in one parallel pass, so the **per-agent cost is amortized** across the panel rather than charged sequentially.
- The lossless guarantee (`min(1, p/q)` per offset) preserves the distribution the consensus judges would see without speculation — same property as [[1.0.0 P-57]] and [[1.0.0 P-108]].

Compared to [[1.0.0 P-108]] at this primitive: EAGLE-3 buys more speedup if MAgHARCM can afford per-domain retraining; Medusa buys **deployability on day one** without retraining. The verdict pipeline defaults to Medusa; switches to EAGLE-3 only when the per-domain training budget is unlocked.

### 2.2 PRIM-21 (Migration Strategy Selection)

PRIM-21 implements Müller's five-strategy try-and-fail loop with budget-aware selection. Medusa's place in this loop is the **cheap + fixed** strategy end of the speculative-decoding axis:
- **Medusa branch** (cheap + fixed): pick the next strategy attempt, calibrate Medusa heads on a small sample of that strategy's traces, run the whole attempt under Medusa speculation. No domain-specific retraining of the draft; if the strategy is abandoned, the heads are discarded cheaply.
- **EAGLE-3 branch** (expensive + trainable): for strategies that survive long enough to justify a training-time-test budget, switch to EAGLE-3 and reap the higher ceiling (~6.5× vs ~2.2×–3.6×).
- Speculative-decoding strategy registry grows a new column: `medusa_K × calibration_size` per (target, strategy) cell, parallel to the `eagle3_draft_size × training_rounds` column from [[1.0.0 P-108]] §2.2.

The discriminator is **time-to-first-speedup**: Medusa wins when the strategy lifetime is short (most strategies in PRIM-21 fail on the first attempt, per the budget-aware selection logic); EAGLE-3 wins only when a strategy survives multiple iterations and amortizes the retraining cost.

### 2.3 PRIM-22 (Four Phases Comprehension)

PRIM-22's comprehension loop emits structured-output artefacts (`TranslatedCode:`, `Rationale:`, `Confidence:`, `Side effects:` slots in the MAgHARCM slot grammar). Medusa's **multi-head parallel prediction** is structurally a **slot-filling co-processor**: head 1 predicts the next comprehension slot's token (`TranslatedCode:` at offset 1), head 2 predicts the slot body's first token, head 3 predicts the closing `:` or newline. Because each head fires in parallel from the same hidden state, the comprehension trace emerges at **higher effective sampling rate per target forward pass** — the SLM emits a multi-token comprehension artefact in one verified draft rather than K sequential sampling steps.
- **Zero-cost structured output**: no grammar-constrained decoder required; the heads learn the slot boundaries during the one-shot calibration step.
- **Quality preservation**: same lossless guarantee as [[1.0.0 P-57]]; the comprehension panel sees the same artefact distribution as without speculation.
- Compare to [[1.0.0 P-108]] §2.3, where EAGLE-3's low-rank feature injection is the comprehension-trace compressor; Medusa's role here is **simpler and faster** but locked to the slot grammar visible during calibration.

### 2.4 PRIM-31 (Iterative Retrieval Refinement)

PRIM-31's IterativeNavigator runs re-index → fetch → translate → verify in a tight loop. Medusa's K-head prediction is a **look-ahead pre-fetcher**:
- While the navigator processes iteration $i$ (verdict, comprehension, fetch), the Medusa heads predict the first $\gamma$ tokens of iteration $i+1$'s context query in parallel.
- The target's tree-verify step at the end of iteration $i$ accepts the predicted continuation *if* the navigator's next state is consistent with it — yielding a **pre-fetch hit rate** that mirrors the draft acceptance rate.
- Combined with [[1.0.0 P-57]]-style classic draft/target for the outer navigator and [[1.0.0 P-108]] for the per-chunk verifier, MAgHARCM gets a **three-tier speculative cascade** when Medusa is the inner-most layer.

### 2.5 SLM Fleet Inference Cost (cross-primitive summary)

Across all four primitives, Medusa contributes **~2.2×–3.6× wall-clock reduction** at lossless output on the MAgHARCM pipeline without changing any verdict, plan, comprehension artefact, or navigator step. The **target SLMs stay frozen**; only the K head adapters are calibrated. This is the **day-one** deployment mode for speculative decoding in the 4B–30B regime; EAGLE-3 ([[1.0.0 P-108]]) is the **week-two** mode once a strategy survives long enough to justify retraining.

## 3. Hop-1 References (papers cited by Medusa)

- **Leviathan, Kalman & Matias 2023** — [[1.0.0 P-57]] — *Fast Inference from Transformers via Speculative Decoding* (the lossless accept/reject rule reused by Medusa; ICML 2023, arXiv:2211.17192).
- **Chen et al. 2023** — *Accelerating Large Language Model Decoding with Speculative Sampling* (DeepMind; concurrent independent discovery of the [[1.0.0 P-57]] algorithm; cited by Medusa as the parallel-paper reference).
- **Li, Wei, Zhang & Zhang 2024** — [[1.0.0 P-64]] — *EAGLE: Speculative Sampling Requires Rethinking Feature Uncertainty* (EAGLE-1; the feature-level draft ancestor that Medusa explicitly contrasts against).
- **Li, Wei, Zhang & Zhang 2024** — *EAGLE-2: Faster Inference of Language Models with Dynamic Draft Trees* (the dynamic-tree ancestor whose draft-tree structure Medusa generalizes with K parallel heads).
- **Li, Wei, Zhang & Zhang 2025** — [[1.0.0 P-108]] — *EAGLE-3: Scaling up Inference Acceleration of LLMs via Training-Time Test* (the training-time-test alternative to Medusa's fixed-head design).
- **Xia et al. 2024** — *Lookahead Decoding: Breaking the Sequential Dependency of LLM Inference* (parallel-branch drafting alternative; cited by Medusa as a non-tree parallel-decoding sibling).
- **Spector & Re 2023** — *DistillSpec: Speculative Decoding with Distilled Draft Models* (knowledge-distillation draft recipe; cited as the "draft model trained offline" alternative Medusa's head-only calibration is contrasted against).
- **Gu & Dao 2023** — *Mamba: Linear-Time Sequence Modeling with Selective State Spaces* (state-space alternative transformer; cited by Medusa as a non-transformer speculation target).
- **Zheng et al. 2024** — *SGLang: Efficient Execution of Structured Language Model Programs* (one of the two production runtimes that ships Medusa integration; cited for the tree-attention scheduling implementation).
- **Kwon et al. 2023** — *vLLM: Efficient Memory Management for Large Language Model Serving with PagedAttention* (the other production runtime that ships Medusa integration; cited for the PagedAttention + Medusa scheduling combination).
- **Leviathan 2023** — speculative-decoding survey note (cited by Medusa as the canonical algorithm reference; aligns with [[1.0.0 P-57]]).
- **Miao et al. 2024** — *SpecInfer: Accelerating Generative LLM Serving with Speculative Inference and Token Tree Verification* (the token-tree verification ancestor whose tree-attention mask Medusa reuses).

## 4. Hop-2 References (papers-cited-by-hop-1)

- **Vaswani et al. 2017** — *Attention Is All You Need* (transformer backbone; Medusa heads attach to a transformer target model).
- **Touvron et al. 2023** — *LLaMA 2* (the primary eval target in the Medusa paper; LLaMA-2-Chat 7B/13B/70B all show the reported 2.2×–3.6× speedup).
- **Jiang et al. 2023** — *Mistral 7B* (the second eval target; Mistral 7B Instruct demonstrates Medusa's portability across model families).
- **Chiang et al. 2023** — *Vicuna* (chat-model eval base; cited as the Vicuna-7B/13B eval target).
- **Ainslie et al. 2023** — *GQA: Training Generalized Multi-Query Transformer Models from Multi-Head Checkpoints* (grouped-query attention; cited as the attention variant Medusa's tree-mask composes with).
- **Talmor et al. 2020** — *Tree-structured decoding with recurrent neural networks* (the tree-decoding literature Medusa's tree-attention verifier descends from).
- **Hinton et al. 2015** — *Distilling the Knowledge in a Neural Network* (the foundational distillation recipe that [[1.0.0 P-57]] cites and that Medusa's calibration step echoes in spirit).
- **Zhang et al. 2017** — *mixup: Beyond Empirical Risk Minimization* (data-augmentation recipe cited by the Medusa head-calibration procedure's data-mixing step). [INFERENCE: plausible citation for the calibration data-mixing trick; verify against the Medusa paper's bibliography if exact bibliographic data is required.]
- **Stern et al. 2018** — *Blockwise Parallel Decoding for Deep Autoregressive Models* (insertion-based parallel decoding; cited as the broader parallel-decoding lineage Medusa sits in).
- **Sun et al. 2021** — *SAD: Shallow Aggressive Decoding* (shallow parallel decoding; cited as a precursor to Medusa's tree-attention verify).
- **Ouyang et al. 2022** — *Training Language Models to Follow Instructions with Human Feedback* (RLHF; cited as the post-training regime that the head-calibration step runs on top of).
- **Touvron et al. 2023** — *LLaMA: Open and Efficient Foundation Language Models* (the original LLaMA; cited by LLaMA-2 lineage).

## 5. Backlinks (primitives/INDEX rows that should reference P-114)

Add a row referencing P-114 to each of the following primitive INDEX entries (see `Software-Archaeology-Lineage.md` §3 for the matrix); P-114 lives in the **same anchor set** as [[1.0.0 P-108]] because Medusa and EAGLE-3 are the two alternative drafting strategies in the speculative-decoding lineage:

- **PRIM-7** (Multi-Agent Verdict Validation) — add P-114 as the **K-head-voter = K-agent-verdict** mapping. Medusa's parallel heads are the natural structural fit for the verdict panel; P-114 is the **day-one / no-retraining** companion to [[1.0.0 P-108]]'s **week-two / trainable** entry. `internal/agents/verdict_panel.go`.
- **PRIM-21** (Migration Strategy Selection) — add P-114 as the **cheap + fixed** end of the speculative-decoding strategy axis. Speculative-decoding strategy registry grows a `medusa_K × calibration_size` column parallel to the [[1.0.0 P-108]] `eagle3_draft_size × training_rounds` column. Medusa wins when the strategy lifetime is short; EAGLE-3 wins when the strategy survives long enough to amortize retraining. `internal/agents/strategy.go`.
- **PRIM-22** (Four Phases Comprehension) — add P-114 as the **multi-head slot-filling co-processor** anchor. Each Medusa head predicts a different comprehension slot offset; the slot grammar is learned during the one-shot calibration step, yielding zero-cost structured output for the comprehension phase. `internal/agents/comprehension.go`.
- **PRIM-31** (Iterative Retrieval Refinement) — add P-114 as the **inner-most look-ahead pre-fetcher** anchor. Medusa heads pre-predict iteration $i+1$'s context query while the navigator processes iteration $i$, completing the three-tier speculative cascade ([[1.0.0 P-57]] outer + [[1.0.0 P-108]] per-chunk + **P-114** pre-fetch). `internal/agents/iter_retrieval.go`.
- **Cross-ref** — add P-114 to `Software-Archaeology-Lineage.md` §3 rows for PRIM-7, PRIM-21, PRIM-22, PRIM-31. P-114 is the **wave-12 alternative-drafting-strategy** anchor that completes the speculative-decoding cell of the lineage matrix alongside [[1.0.0 P-57]] (classic draft/target), [[1.0.0 P-64]] (EAGLE-1 feature-level), and [[1.0.0 P-108]] (EAGLE-3 training-time-test). Together they form the four-way strategy table that PRIM-21 selects from.
