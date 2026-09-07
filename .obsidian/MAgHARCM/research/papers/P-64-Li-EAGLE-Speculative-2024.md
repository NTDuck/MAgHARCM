---
title: "P-64 — Li, Wei, Zhang & Zhang 2024 — EAGLE: Speculative Sampling Requires Rethinking Feature Uncertainty"
backlink: "[[1.0.0 P-64]]"
aliases:
  - "1.0.0 P-64"
  - "P-64"
  - "P-64-Li-EAGLE-Speculative-2024"
  - "P-64-Li-EAGLE-Speculative-2024"
  - "Li-EAGLE-Speculative-2024"
tags: [paper, inference, decoding, speculative, feature-level, slm, lossless, code-generation, [[1.0.0 PRIM-7]], [[1.0.0 P-57]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-64 — EAGLE: Speculative Sampling Requires Rethinking Feature Uncertainty]]

## Citation

Li, Y., Wei, F., Zhang, C., & Zhang, H. (2024). *EAGLE: Speculative Sampling Requires Rethinking Feature Uncertainty*. ICML 2024. arXiv:2401.15077. URL: https://arxiv.org/abs/2401.15077

## Abstract Summary

EAGLE reformulates speculative decoding at the **feature level** rather than the token level. The draft model autoregresses on the second-to-top-layer feature sequence shifted by one time step, then the target model verifies the proposed tokens in one parallel forward pass using the feature-conditioned probability. The feature-shift trick makes the draft distribution easy to predict (low entropy) while the target verification reuses a single forward pass, yielding lossless output distribution and **2.7×–3.5× wall-clock speedup on LLaMA2-Chat 70B** with doubled throughput. Evaluated on dialogue, **code generation**, math, and instruction following across Vicuna, LLaMA2-Chat, and Mixtral 8x7B.

## Method

EAGLE sits inside the speculative-decoding lineage of [[1.0.0 P-57]] (Leviathan, Kalman & Matias 2023) but removes the requirement for a separate draft model. The single trained model runs two heads: the feature-draft head predicts the next feature vector from the shifted second-to-top-layer sequence, and the LM head predicts the token distribution conditioned on the verified feature. Acceptance uses the same `min(1, p/q)` rejection rule as classic spec decoding. The shift-by-one construction is the key insight — raw second-to-top features are unstable, but shifted features form a smooth autoregressive sequence.

## Findings Relevant to MAgHARCM

- Code generation is a first-class evaluation target in EAGLE (HumanEval / MBPP), so the speedup claim holds on the task domain MAgHARCM cares about.
- Lossless output distribution matters for [[1.0.0 PRIM-7]] (Multi-Agent Verdict Validation): the verdict panel cannot tolerate a draft that shifts the distribution the consensus judges.
- Single-model design (no separate draft) is a clean fit for MAgHARCM's SLM fleet — adding a heavyweight second model would defeat the cost thesis.
- EAGLE composes with [[1.0.0 P-57]]'s classic draft/target pattern for the larger [[1.0.0 PRIM-31]] (Iterative Retrieval Refinement) loop.

## How MAgHARCM Uses It

[[1.0.0 PRIM-7]] currently fans out $K$ SLM proposals to a single oracle. EAGLE turns each SLM **into its own draft** at the feature level, so the same model both proposes tokens and verifies them via the top-layer head in one pass — collapsing PRIM-7's two-stage "SLM proposes, oracle verifies" into one distributional step. Combined with [[1.0.0 P-57]]-style classic draft/target for the bigger PRIM-31 navigator loop, expected cost drop is ~2.5× on the verdict step without changing the verdict.

## References

### Hop-1 (EAGLE's direct citations)
- Vaswani et al. 2017 — *Attention Is All You Need* — transformer backbone.
- Leviathan, Kalman & Matias 2023 — [[1.0.0 P-57]] — speculative decoding baseline.
- Chen et al. 2023 — DeepMind — improved draft model architectures.
- Cai et al. 2024 — Medusa — multi-head parallel draft heads.
- Gante 2024 — *Assisted Generation* — HF reference implementation.
- Sun et al. 2021 — *SAD* — sequence-level draft.
- Stern et al. 2018 — insertion-based decoding.

### Hop-2 (transformer / scaling-law / code-LLM eval lineage)
- Touvron et al. 2023 — *LLaMA 2* — base model for EAGLE eval.
- Chiang et al. 2023 — *Vicuna* — dialogue eval.
- Jiang et al. 2024 — *Mixtral 8x7B* — MoE eval.
- Chen et al. 2021 — *HumanEval* — code-generation benchmark.
- Austin et al. 2021 — *MBPP* — code-generation benchmark.
- Hendrycks et al. 2021 — *APPS* — code-generation benchmark.

## Backlinks

[[1.0.0 PRIM-7]], [[1.0.0 P-57]], [[1.0.0 PRIM-3]], [[1.0.0 PRIM-31]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].
