---
title: "P-57 — Leviathan, Kalman & Matias 2023 — Fast Inference from Transformers via Speculative Decoding"
backlink: "[[1.0.0 P-57]]"
tags: [paper, inference, decoding, speculative, draft-model, slm, lossless, [[1.0.0 PRIM-3]], [[1.0.0 PRIM-7]], [[1.0.0 PRIM-31]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-57 — Leviathan, Kalman & Matias 2023 — Speculative Decoding]]

- **Authors**: Yaniv Leviathan, Matan Kalman, Yossi Matias (Google Research).
- **Venue / Year**: ICML 2023; arXiv:2211.17192 (November 2022, revised May 2023).
- **URL**: https://arxiv.org/abs/2211.17192
- **Anchors**: PRIM-3 (Target Skeleton-First), PRIM-7 (Multi-Agent Verdict Validation), PRIM-31 (Iterative Retrieval Refine); the canonical inference-time acceleration anchor for SLM-flavored agent loops.

## 1. Core Contribution

**Speculative decoding** accelerates inference from a large "target" autoregressive model $M_p$ by pairing it with a smaller, faster "draft" model $M_q$ that proposes the next $\gamma$ tokens autoregressively, which the target then verifies in a single parallel forward pass. The key result: **the output distribution is mathematically identical** to sampling from $M_p$ alone (lossless), yet wall-clock latency drops **2×–3×** on T5-XXL and 137B-parameter LaMDA.

Algorithm (SpeculativeDecodingStep):
1. Draft $\gamma$ candidate tokens from $M_q$ autoregressively.
2. Run $M_p$ on all $\gamma + 1$ prefixes in parallel.
3. Accept each draft token with probability $\min\!\bigl(1,\, p(x) / q(x)\bigr)$; reject and resample from the adjusted distribution $p'(x) = \mathrm{norm}\bigl(\max(0,\, p(x) - q(x))\bigr)$.

This is **strictly better** than rejection sampling: the expected acceptance rate is higher because speculative sampling is not constrained to resample from the unmodified $p(x)$. Empirical results:
- T5-XXL English→German translation: **2.5×–3.4×** speedup using T5-small as the draft model.
- 137B-parameter LaMDA dialog: **2×–2.3×** speedup using a smaller LaMDA variant as the draft.
- Acceptance rate $\alpha$ is the key knob: the higher the agreement between $M_q$ and $M_p$, the more tokens per parallel pass.

**Concurrent independent discovery**: Chen et al. (2023, DeepMind) — *Accelerating Large Language Model Decoding with Speculative Sampling* — derived the same algorithm in the same timeframe. Both papers are cited as the foundational references.

**SLM-aware relevance for MAgHARCM**: in the SLM regime MAgHARCM targets (Qwen2.5-Coder 3B-7B, StarCoder2 3B-7B, Phi-3-mini 3.8B), speculative decoding enables pairing a **small draft model** (e.g., StarCoder2 3B) with a **larger target model** (e.g., Qwen2.5-Coder 7B-Instruct) without changing the output distribution. This is the natural inference architecture for MAgHARCM's planner-executor split: a fast draft model proposes candidate translations, a slower verifier model validates them.

## 2. Application in MAgHARCM

- **PRIM-3 (Target Skeleton-First)** — speculative decoding pairs naturally with skeleton-first translation: the **draft model** generates the cheap skeleton (trait signatures, function signatures, type declarations), and the **target model** verifies the skeleton AND fills in the method bodies in a single parallel pass. The skeleton's high acceptance rate (structurally constrained output) maximizes the speedup factor.
- **PRIM-7 (Multi-Agent Verdict Validation)** — the speculative-decoding acceptance step is structurally a **single-token verdict**: "does the draft's token match the target's distribution within tolerance?" PRIM-7's verdict panel extends this from N tokens to N full translations, with the same accept/reject discipline. The pattern of "fast proposer + slow verifier, distribution-preserving" is exactly the MAgHARCM verdict architecture.
- **PRIM-31 (Iterative Retrieval Refine)** — each iteration of the IterativeNavigator (re-index → fetch → translate → verify) is a candidate for speculative decoding: the **draft model** proposes the chunk translation, the **target model** verifies and either accepts or resamples. Lossless output guarantees that the validation oracle sees exactly the same code as it would have without speculation.
- **SLM Fleet Inference Cost** — MAgHARCM runs 4-7 agents per pipeline execution; if each agent uses speculative decoding with a draft/target pair from the same family (e.g., StarCoder2 3B draft → Qwen2.5-Coder 7B target), the total inference cost drops ~2× without changing any output, making the pipeline feasible on consumer hardware.

## 3. Hop-1 References (papers cited by Leviathan, Kalman & Matias)

- Vaswani et al. (2017) — Attention Is All You Need (the transformer architecture that speculative decoding operates on).
- Brown et al. (2020) — GPT-3 / Language Models are Few-Shot Learners (cited as the canonical "large autoregressive model" that motivates speedup).
- Chowdhery et al. (2022) — PaLM (referenced as a large-model baseline in the same regime as GPT-3).
- Thoppilan et al. (2022) — LaMDA (referenced as the 137B dialog model used in experiments).
- Raffel et al. (2020) — T5 (the family of models used in translation experiments; T5-XXL is the target, T5-small/base/large are the drafts).
- Roberts et al. (2022) — T5X / SeqIO (the production T5 implementation that the wallclock baseline uses).
- Hinton et al. (2015) — Distilling the Knowledge in a Neural Network (foundational draft-from-large-model recipe; referenced as an alternative to speculative decoding).
- Shazeer (2019) — Fast Transformer Decoding: One Write-Head is All You Need (multi-query attention; related decoding-speedup technique).
- Stern et al. (2018) — Blockwise Parallel Decoding for Deep Autoregressive Models (the precursor parallel-decoding technique).
- Sun et al. (2021) — Instantaneous Grammatical Error Correction with Shallow Aggressive Decoding (parallel-decoding predecessor).

## 4. Hop-2 References (papers-cited-by-hop-1)

- Sanh et al. (2019) — DistilBERT (knowledge distillation; cited by Hinton et al. distillation lineage).
- Jiao et al. (2020) — TinyBERT (further distillation of BERT; cited by knowledge-distillation lineage).
- Devlin et al. (2019) — BERT (foundational transformer encoder; cited by T5 and DistilBERT).
- Howard & Ruder (2018) — ULMFiT (transfer-learning recipe; cited by T5).
- Kaplan et al. (2020) — Scaling Laws for Neural Language Models (the scaling hypothesis that motivates speculative decoding as an alternative to scale-up).
- Hoffmann et al. (2022) — Chinchilla compute-optimal scaling (referenced by PaLM as the compute-optimal precedent).
- Liu et al. (2024) — Lost in the Middle [[P-53]] in MAgHARCM lineage (cited as a long-context challenge that speculative decoding helps with).
- Schick et al. (2023) — Toolformer (cited by speculative-decoding extensions as a downstream use case).
- Beltagy et al. (2020) — Longformer (sparse-attention baseline referenced by Blockwise Parallel Decoding).
- Child et al. (2019) — Sparse Transformers (foundational sparse-attention work referenced by Stern et al. 2018).

## 5. Backlinks

- **PRIM-3** (Target Skeleton-First): skeleton generation is the high-acceptance-rate regime for speculative decoding.
- **PRIM-7** (Multi-Agent Verdict Validation): single-token accept/reject discipline generalizes to N-translation verdict.
- **PRIM-31** (Iterative Retrieval Refine): per-chunk draft/target pair within the navigator loop.
- **PRIM-21** (Migration Strategy Selection): strategy-level ensemble can use draft/target pairs per strategy.
- **Inference cost optimization**: ~2× cost reduction on the entire MAgHARCM pipeline without output change.
- **Cross-ref**: add P-57 to PRIM-3, PRIM-7, PRIM-31 rows in `Software-Archaeology-Lineage.md`. This is the **canonical inference-time acceleration anchor** for the SLM fleet.
