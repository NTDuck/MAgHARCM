---
title: "P-78 — Li, Wei, Zhang & Zhang 2025 — EAGLE-3: Scaling up Inference Acceleration of LLMs via Training-Time Test"
backlink: "[[1.0.0 P-78]]"
tags: [paper, inference, decoding, speculative, feature-level, multi-token, slm, lossless, code-generation, [[1.0.0 PRIM-7]], [[1.0.0 P-57]], [[1.0.0 P-64]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-78 — EAGLE-3 (Li et al. 2025)]]

- **Authors**: Yuhui Li, Fangyun Wei, Chao Zhang, Hongyang Zhang.
- **Venue / Year**: arXiv:2503.01840 (March 2025); later integrated into vLLM, SGLang, and TensorRT-LLM speculative-decoding stacks.
- **URL**: https://arxiv.org/abs/2503.01840
- **Code**: https://github.com/SafeAILab/EAGLE
- **Anchors**: [[1.0.0 P-57]] (Leviathan 2023 — speculative-decoding anchor), [[1.0.0 P-64]] (EAGLE / EAGLE-1 / EAGLE-2 lineage); relevant to [[1.0.0 PRIM-7]] (Multi-Agent Verdict Validation) and [[1.0.0 PRIM-31]] (Iterative Retrieval Refinement).

## 1. Core Contribution

EAGLE-3 is the third-generation EAGLE family member and the first to **drop feature-level drafting entirely**. Where EAGLE / EAGLE-2 (cf. [[1.0.0 P-64]]) autoregressed on the second-to-top-layer feature sequence to keep draft entropy low, EAGLE-3 **drafts at the token level** and compensates for the higher per-step entropy with two new mechanisms:

1. **Multi-layer feature fusion** — instead of conditioning the draft head on the target model's top-layer feature only, the draft head consumes a fused representation that aggregates hidden states from **low, mid, and high layers** of the target model, combined with the token-embedding stream. The draft model therefore conditions on the full computational context of the target rather than its final layer alone.
2. **Training-time test** — during training, the draft model simulates the inference regime (autoregressive draft → single-pass target verify → accept/reject) so its gradient signal matches the deployment-time distribution it will see. This closes the train/inference mismatch that limited earlier speculative heads.

The two changes compound: feature-prediction-as-a-head was the bottleneck that prevented EAGLE-1/2 from scaling with data, so once the constraint is removed, scaling the draft training set translates into longer average accepted runs. Reported results:
- Up to **6.5× wall-clock speedup** vs vanilla decoding across chat and reasoning models on five tasks (dialogue, code generation, math, summarisation, instruction following).
- **~1.4×** speedup over EAGLE-2 at matched draft budget, **~1.8×** at matched draft quality.
- In the SGLang framework at batch size 64: **1.38× throughput improvement** over the EAGLE-2 baseline, indicating the gain survives high concurrency.
- Lossless output distribution preserved (same `min(1, p/q)` acceptance rule as [[1.0.0 P-57]]).

## 2. Method (in the speculative-decoding lineage)

EAGLE-3 sits inside the [[1.0.0 P-57]] (Leviathan, Kalman & Matias 2023) speculative-decoding lineage and inherits its verify step:

1. Draft model $\mathcal{D}_\theta$ proposes $\gamma$ candidate tokens autoregressively, conditioned on the multi-layer-fused feature stream plus the token embedding of the verified prefix.
2. Target model $M_p$ runs a **single parallel forward pass** over the $\gamma+1$ extended prefixes, reusing its own hidden states.
3. Each draft token is accepted with probability $\min(1,\, p(x)/q(x))$; on rejection the next token is resampled from the standard speculative-sampling residual $\mathrm{norm}\bigl(\max(0,\, p(x) - q(x))\bigr)$.

The novelty is exclusively on the **draft side**:
- *Why drop feature prediction?* EAGLE authors observed that EAGLE-1/2 saturated quickly as draft training data scaled: the second-to-top-layer feature sequence has bounded entropy and cannot encode more information than the top-layer LM head already does, so extra data has nothing to fit. Token-level drafting, conditioned on multi-layer fusion, has a strictly larger hypothesis space and so benefits from scale.
- *Why multi-layer fusion?* Top-layer features are the most "decision-shaped" and lose lower-level syntactic / lexical signal. Low/mid layers preserve lexical anchors that token-level drafting can exploit directly. Concatenating or summing low+mid+high gives the draft head access to both ends of the abstraction hierarchy.
- *Why training-time test?* Speculative-decoding drafts see a different input distribution at inference (only accepted prefixes survive) than at training (the draft sees its own outputs). Simulating the accept/reject loop during training makes the draft robust to the deployment distribution.

## 3. Findings Relevant to MAgHARCM

- **Lossless output matters for [[1.0.0 PRIM-7]] (Multi-Agent Verdict Validation)** — the verdict panel cannot tolerate a draft that shifts the distribution the consensus judges. EAGLE-3 preserves the [[1.0.0 P-57]] lossless guarantee, so it composes cleanly into the verdict pipeline: each SLM proposal becomes a draft, the oracle becomes the verifier, and the verdict sees exactly the same outputs as non-speculative decoding.
- **Higher speedup directly attacks the verdict cost** — at **6.5×** over vanilla (vs EAGLE-1's 2.7–3.5×), the per-agent cost of PRIM-7's $K$-SLM fan-out drops enough that running 4–7 agents per pipeline pass becomes cheap on consumer hardware. MAgHARCM's SLM-fleet cost thesis (cf. [[1.0.0 P-57]] §2) tightens further.
- **Multi-layer fusion is a clean fit for MAgHARCM's two-tier model mix** — MAgHARCM already pairs small draft SLMs with larger target SLMs (e.g., StarCoder2 3B → Qwen2.5-Coder 7B-Instruct, cf. [[1.0.0 P-58]]). Multi-layer fusion lets the small draft see the target's mid-layer signal cheaply (it reads off the target's forward pass for free), so the draft-to-target gap shrinks without paying a second-model forward cost.
- **Training-time test matches MAgHARCM's retraining cadence** — MAgHARCM periodically re-trains its small drafts on the latest agent traces. Training-time test is exactly the recipe needed: train the draft on the loop's accept/reject signal so it learns to anticipate verifier rejections, rather than training only on ground-truth next-token prediction.
- **Code generation as first-class eval** — EAGLE-3 evaluates on HumanEval / MBPP / APPS, so the speedup claim holds on the task domain MAgHARCM cares about.
- **Production-framework integrations** (vLLM, SGLang, TensorRT-LLM) mean MAgHARCM can adopt EAGLE-3 without rebuilding the inference stack.

## 4. How MAgHARCM Uses It

[[1.0.0 PRIM-7]] currently fans out $K$ SLM proposals to a single oracle. The path EAGLE-3 opens:

1. **Each SLM proposal is its own EAGLE-3 draft**: a 3B-parameter Qwen2.5-Coder/StarCoder2-class draft is fused against the multi-layer hidden states of a 7B target during the verifier's single forward pass. The proposal cost drops by ~5–6× while the verifier still sees the same lossless distribution.
2. **Per-token rejection becomes per-translation verdict**: the `min(1, p/q)` accept/reject discipline generalises from single tokens to full translations; the verdict panel runs the same speculative-sampling residual on each candidate proposal.
3. **Training-time test on the agent loop**: re-train each draft head on (prompt, oracle-verdict-prefix) pairs collected from the latest pipeline runs, simulating the accept/reject distribution the draft will see at inference. The draft's average accepted-run length grows with data, compounding the wall-clock gain.
4. **Combined with [[1.0.0 P-57]] classic draft/target** for the larger [[1.0.0 PRIM-31]] (Iterative Retrieval Refinement) navigator loop: the outer navigator uses EAGLE-3-style multi-layer fusion drafts, the inner per-chunk translation uses Leviathan-style draft/target pairs. Stacking gives ~5× on the verdict step and ~2.5× on the navigator step.

Net effect: the SLM-fleet inference cost drops ~3–4× without changing any output, making the four-to-seven-agent MAgHARCM pipeline feasible on a single consumer GPU.

## 5. References

### Hop-1 (EAGLE-3's direct citations)
- Leviathan, Kalman & Matias 2023 — [[1.0.0 P-57]] — speculative decoding baseline (the lossless accept/reject rule reused here).
- Chen et al. 2023 — *Accelerating Large Language Model Decoding with Speculative Sampling* (DeepMind) — concurrent independent discovery of speculative decoding.
- Cai et al. 2024 — *Medusa: Multiple Heads are Better than One* — multi-head parallel draft heads (alternative draft architecture).
- Li, Wei, Zhang & Zhang 2024 — [[1.0.0 P-64]] — *EAGLE: Speculative Sampling Requires Rethinking Feature Uncertainty* (EAGLE-1, the feature-level ancestor).
- Li, Wei, Zhang & Zhang 2024 — *EAGLE-2: Faster Inference of Language Models with Dynamic Draft Trees* (the multi-token / dynamic-tree ancestor EAGLE-3 builds on).
- Gante 2024 — *Assisted Generation* — HF reference implementation.
- Sun et al. 2021 — *SAD* — sequence-level draft.
- Stern et al. 2018 — insertion-based decoding.
- Vaswani et al. 2017 — *Attention Is All You Need* — transformer backbone.

### Hop-2 (transformer / scaling-law / code-LLM eval lineage)
- Touvron et al. 2023 — *LLaMA 2* — chat-model base for EAGLE-3 eval.
- Chiang et al. 2023 — *Vicuna* — chat-model base.
- Jiang et al. 2024 — *Mixtral 8x7B* — MoE eval.
- Chen et al. 2021 — *HumanEval* — code-generation benchmark.
- Austin et al. 2021 — *MBPP* — code-generation benchmark.
- Hendrycks et al. 2021 — *APPS* — code-generation benchmark.
- DeepSeek 2025 — *DeepSeek-R1* — reasoning-model eval (EAGLE-3 reports gains on reasoning models).

## 6. Backlinks

[[1.0.0 P-57]], [[1.0.0 P-64]], [[1.0.0 PRIM-7]], [[1.0.0 PRIM-31]], [[1.0.0 PRIM-3]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-78 is the **current inference-time acceleration frontier** for the MAgHARCM SLM fleet: extends the [[1.0.0 P-57]] anchor and the [[1.0.0 P-64]] EAGLE-1 lineage with multi-layer fusion + training-time test for ~6.5× lossless speedup.
