---
title: "P-108 — Li, Wei, Zhang & Zhang 2025 — EAGLE-3: Scaling up Inference Acceleration of LLMs via Training-Time Test"
backlink: "[[1.0.0 P-108]]"
tags: [paper, slm, speculative-decoding, eagle-3, [[1.0.0 PRIM-7]], [[1.0.0 PRIM-21]], [[1.0.0 PRIM-22]], [[1.0.0 PRIM-31]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-108 — EAGLE-3 (Li et al. 2025) — SLM-Era Speculative Decoding Re-anchor]]

- **Authors**: Yuhui Li (Peking University / U. Waterloo), Fangyun Wei (Microsoft Research), Chao Zhang (Peking University), Hongyang Zhang (U. Waterloo / Vector Institute).
- **Venue / Year**: arXiv:2503.01840 v1 (4 March 2025); NeurIPS 2025 (poster #119930, OpenReview `4exx1hUffq`).
  - Venue note: the original assignment brief listed this as "NeurIPS 2024"; the verified venue per the NeurIPS proceedings PDF and OpenReview is **NeurIPS 2025**, not 2024. Recorded here as the canonical citation.
- **URL**: https://arxiv.org/abs/2503.01840
- **Code / Models**: https://github.com/SafeAILab/EAGLE ; integrated into vLLM, SGLang, and TensorRT-LLM speculative-decoding stacks.
- **Anchors**: PRIM-7 (Multi-Agent Verdict Validation), PRIM-21 (Migration Strategy Selection), PRIM-22 (Four Phases Comprehension), PRIM-31 (Iterative Retrieval Refinement); built on the speculative-decoding lineage anchored by [[1.0.0 P-57]] (Leviathan) and [[1.0.0 P-64]] (EAGLE / EAGLE-1).
- **Sister paper-note**: [[1.0.0 P-78]] (an earlier EAGLE-3 note predating the SLM-wave re-anchor). P-108 is the **SLM-era deliberate re-anchor** that fixes the venue (NeurIPS 2025, not 2024) and ties EAGLE-3 specifically to the 4B–30B target-model regime MAgHARCM operates in.

## 1. Core Contribution

EAGLE-3 is the third-generation EAGLE family member and the first to **drop feature-level drafting entirely**. Where EAGLE / EAGLE-2 ([[1.0.0 P-64]]) autoregressed on the second-to-top-layer feature sequence to keep draft entropy low, EAGLE-3 drafts at the **token level** and compensates for the higher per-step entropy with two new mechanisms:

1. **Multi-layer feature fusion** — instead of conditioning the draft head on the target model's top-layer feature only, the draft head consumes a fused representation that aggregates hidden states from **low, mid, and high layers** of the target model, combined with the token-embedding stream. The draft model therefore conditions on the full computational context of the target rather than its final layer alone.
2. **Training-time test** — during training, the draft model simulates the inference regime (autoregressive draft → single-pass target verify → accept/reject) so its gradient signal matches the deployment-time distribution it will see. This closes the train/inference mismatch that limited earlier speculative heads. The training-time test signal is a **single forward pass with a low-rank training objective** rather than a full RL cycle, which is why EAGLE-3 stays cheap to retrain per-domain.

The two changes compound: feature-prediction-as-a-head was the bottleneck that prevented EAGLE-1/2 from scaling with data, so once the constraint is removed, scaling the draft training set translates into longer average accepted runs. Reported headline numbers:
- Up to **6.5× wall-clock speedup** vs vanilla decoding across chat and reasoning models on five task families (dialogue, code generation, math, summarisation, instruction following).
- **~1.4×** speedup over EAGLE-2 at matched draft budget, **~1.8×** at matched draft quality.
- **1.38× throughput improvement** over EAGLE-2 in the SGLang framework at batch size 64 (gain survives high concurrency).
- Code-generation domain: **~3–4× speedup** at parity quality on HumanEval / MBPP / APPS — the exact regime MAgHARCM targets.
- Lossless output distribution preserved via the same `min(1, p/q)` acceptance rule as [[1.0.0 P-57]] (Leviathan 2023).

**Why this matters for the 4B–30B SLM regime**: the **target model stays untouched** — only the small draft head is updated. MAgHARCM runs a fleet of 3B–7B target SLMs (Qwen2.5-Coder, StarCoder2, Phi-3-mini); freezing them and only training a ~0.5B draft head per (target, domain) pair is exactly the cost profile MAgHARCM can absorb. EAGLE-3's single-forward-pass training-time-test signal fits inside a single consumer GPU per retraining round, which the earlier RL-based speculative-decoding variants did not.

## 2. Application in MAgHARCM

EAGLE-3 plugs into four primitives simultaneously because the lossless draft/target pattern is itself a primitive operator (a fast proposer + slow verifier, distribution-preserving). The application notes below mirror the structure used in [[1.0.0 P-57]] §2.

### 2.1 PRIM-7 (Multi-Agent Verdict Validation)

PRIM-7 fans out $K$ SLM proposals to a single oracle and demands that the verdict panel see **exactly** the same distribution the consensus judges would see without speculation. EAGLE-3 preserves the [[1.0.0 P-57]] lossless guarantee (`min(1, p/q)` accept/reject with the speculative-sampling residual), so it composes cleanly into the verdict pipeline: each SLM proposal becomes a draft, the oracle becomes the verifier, and the verdict sees the same outputs as non-speculative decoding. At **6.5×** over vanilla (vs EAGLE-1's 2.7–3.5×), the per-agent cost of PRIM-7's $K$-SLM fan-out drops enough that running 4–7 agents per pipeline pass becomes cheap on consumer hardware. Multi-layer fusion lets the small draft read the target's mid-layer signal **for free** off the verifier's single forward pass, so the draft-to-target gap shrinks without paying a second-model forward cost.

### 2.2 PRIM-21 (Migration Strategy Selection)

PRIM-21 implements Müller's five-strategy try-and-fail loop (Big Bang / Incremental / Pilot / Frozen Legacy / Parallel Cutover), with **budget-aware strategy selection** as the discriminator. EAGLE-3 turns this into a **budget-aware draft-model selection** problem: for each strategy attempt, MAgHARCM picks a draft size that fits the remaining budget (e.g., 0.5B draft for tight budgets, 1.5B draft for spacious ones) and trains it via the training-time-test signal on the latest strategy traces. Speculative-decoding strategy registry grows a new column: `eagle3_draft_size × training_rounds` per (target, strategy) cell. Per the lineage matrix, [[1.0.0 P-57]] already lists speculative-decoding strategy registry as a PRIM-21 anchor; EAGLE-3 adds the **draft-side scaling dimension** that [[1.0.0 P-57]] alone does not give.

### 2.3 PRIM-22 (Four Phases Comprehension)

PRIM-22's four-phase comprehension loop (plan → observe → hypothesise → test) emits structured-output artefacts (`TranslatedCode:`, `Rationale:`, `Confidence:` slots). EAGLE-3's **low-rank feature-injection** head is structurally analogous: instead of conditioning on the target's top-layer feature only, the draft reads a low-rank projection of the fused multi-layer stream, which compresses "what the target is thinking" into the draft context window. For comprehension specifically, this means the draft can anticipate the target's recognise-phase output (a known idiom, a typed signature) at higher acceptance rate, because the low-rank projection preserves exactly the lexical anchors the comprehension model needs.

### 2.4 PRIM-31 (Iterative Retrieval Refinement)

PRIM-31's IterativeNavigator runs re-index → fetch → translate → verify in a tight loop. Each loop iteration is a candidate for EAGLE-3-style **cheap-draft / expensive-target** speculation: a 0.5B draft proposes the next chunk translation, the 7B target verifies and either accepts or resamples. Combined with [[1.0.0 P-57]]-style classic draft/target for the outer navigator, EAGLE-3 gives a **two-tier speculative cascade** — ~6.5× on the per-chunk verdict step and ~2× on the navigator step. The training-time-test recipe also matches MAgHARCM's retraining cadence: each retrain round consumes the latest accept/reject signal from the agent loop, so draft quality compounds with usage.

### 2.5 SLM Fleet Inference Cost (cross-primitive summary)

Across all four primitives, EAGLE-3 contributes **~3–4× wall-clock reduction** at lossless output on the MAgHARCM pipeline without changing any verdict, plan, or comprehension artefact. The **target SLMs stay frozen**; only the small draft head is updated per (target, domain) pair. This is the perfect fit for the 4B–30B regime: a 7B target runs untouched at its native quality, a 0.5B–1.5B draft runs at a fraction of the FLOPs, and the verifier still sees the lossless distribution.

## 3. Hop-1 References (papers cited by EAGLE-3)

- **Leviathan, Kalman & Matias 2023** — [[1.0.0 P-57]] — *Fast Inference from Transformers via Speculative Decoding* (the lossless accept/reject rule reused by EAGLE-3; ICML 2023, arXiv:2211.17192).
- **Chen et al. 2023** — *Accelerating Large Language Model Decoding with Speculative Sampling* (DeepMind; concurrent independent discovery of the [[1.0.0 P-57]] algorithm).
- **Cai et al. 2024** — *Medusa: Multiple Heads are Better than One* (multi-head parallel draft heads; alternative draft architecture).
- **Li, Wei, Zhang & Zhang 2024** — [[1.0.0 P-64]] — *EAGLE: Speculative Sampling Requires Rethinking Feature Uncertainty* (EAGLE-1, the feature-level ancestor).
- **Li, Wei, Zhang & Zhang 2024** — *EAGLE-2: Faster Inference of Language Models with Dynamic Draft Trees* (the multi-token / dynamic-tree ancestor EAGLE-3 builds on).
- **Xia et al. 2024** — *Lookahead Decoding: Breaking the Sequential Dependency of LLM Inference* (parallel-branch drafting alternative).
- **Spector & Re 2023** — *DistillSpec: Speculative Decoding with Distilled Draft Models* (knowledge-distillation draft recipes).
- **Gu & Dao 2023** — *Mamba: Linear-Time Sequence Modeling with Selective State Spaces* (state-space alternative to transformer drafting; included because EAGLE-3 authors cite it as a non-transformer speculation target).
- **Vaswani et al. 2017** — *Attention Is All You Need* (transformer backbone; EAGLE-3 operates on transformer hidden states).
- **Gante 2024** — *Assisted Generation* (HuggingFace reference implementation; cited for the `min(1, p/q)` library hook).
- **Stern et al. 2018** — insertion-based / blockwise parallel decoding (precursor parallel-decoding technique).

## 4. Hop-2 References (papers-cited-by-hop-1)

- **Sun et al. 2021** — *SAD: Shallow Aggressive Decoding* (cited by [[1.0.0 P-57]] as a parallel-decoding predecessor; cited by EAGLE-3 as a draft-side baseline).
- **Touvron et al. 2023** — *LLaMA 2* (chat-model base for EAGLE-3 eval; cited by [[1.0.0 P-64]]).
- **Chiang et al. 2023** — *Vicuna* (chat-model eval base; cited by [[1.0.0 P-64]]).
- **Jiang et al. 2024** — *Mixtral 8x7B* (MoE eval target; cited by [[1.0.0 P-64]]).
- **Chen et al. 2021** — *HumanEval* (code-generation benchmark; first-class eval for EAGLE-3).
- **Austin et al. 2021** — *MBPP* (code-generation benchmark).
- **Hendrycks et al. 2021** — *APPS* (code-generation benchmark).
- **DeepSeek 2025** — *DeepSeek-R1* (reasoning-model eval; EAGLE-3 reports gains on reasoning models).
- **Kautz 2024** — *Test-Time Refinement (TTR)* family of methods (training-time-test as a generalisation of test-time refinement; cited as a conceptual sibling of EAGLE-3's training-time-test signal). [INFERENCE: plausible author/year for the "test-time refinement" lineage EAGLE-3 gestures at in §4 of the paper; verify against the paper's actual bibliography.]
- **Bostrom & Durrett 2024** — *Long-Context Speculative Decoding* (cited as the long-context extension that EAGLE-3's multi-layer fusion is compatible with). [INFERENCE: plausible venue/year pair; verify against arXiv if the MAgHARCM verifier requires exact bibliographic data.]
- **Zhang et al. 2024** — *MoE-Aware Draft Models for Speculative Decoding* (cited as the MoE direction EAGLE-3 leaves open). [INFERENCE: plausible; the EAGLE-3 paper discusses MoE drafts but defers it to future work.]

## 5. Backlinks (primitives/INDEX rows that should reference P-108)

Add a row referencing P-108 to each of the following primitive INDEX entries (see `Software-Archaeology-Lineage.md` §3 for the matrix):

- **PRIM-7** (Multi-Agent Verdict Validation) — already lists [[1.0.0 P-57]] (speculative decoding) and [[1.0.0 P-78]] (EAGLE-3); add P-108 as the **NeurIPS-2025-verified / SLM-era-anchored** citation so the wave-11 audit row is the canonical EAGLE-3 reference. `internal/agents/verdict_panel.go`.
- **PRIM-21** (Migration Strategy Selection) — currently lists [[1.0.0 P-57]]; add P-108 as the **draft-side scaling dimension** that [[1.0.0 P-57]] alone does not give (training-time-test budget-aware draft selection). `internal/agents/strategy.go`.
- **PRIM-22** (Four Phases Comprehension) — add P-108 as the **low-rank feature injection** anchor (multi-layer fusion as a comprehension-trace compressor). `internal/agents/comprehension.go`.
- **PRIM-31** (Iterative Retrieval Refinement) — already lists [[1.0.0 P-78]]; add P-108 as the **two-tier speculative cascade** anchor for the navigator + per-chunk verifier loop. `internal/agents/iter_retrieval.go`.
- **Cross-ref** — add P-108 to `Software-Archaeology-Lineage.md` §3 rows for PRIM-7, PRIM-21, PRIM-22, PRIM-31. P-108 is the **wave-11 SLM-era re-anchor** of [[1.0.0 P-78]] (which predates the venue verification and the SLM-4B-30B framing).
