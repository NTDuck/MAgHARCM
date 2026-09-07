---
title: "P-126 ARC-Decode: Accelerated Decoding with Risk-Bounded Acceptance"
backlink: "[[1.0.0 P-126]]"
aliases:
  - "1.0.0 P-126"
  - "P-126"
  - "P-126-ARC-Decode-Speculative-Decoding-2026"
  - "ARC-Decode-Speculative-Decoding-2026"
  - "li2026arcdecode"
tags: [paper, speculative-decoding, sampling, inference-acceleration, training-free, "[[1.0.0 PRIM-7]]", "[[1.0.0 PRIM-21]]", "[[1.0.0 P-57]]", "[[1.0.0 P-64]]", "[[1.0.0 P-78]]", "[[1.0.0 P-108]]", "[[1.0.0 P-114]]", wave-18]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICML 2026
---

# [[1.0.0 P-126]] ARC-Decode

## TL;DR

ARC-Decode is a **training-free, risk-bounded acceptance rule** for speculative decoding under sampling (T > 0). It uses a Jensen-Shannon divergence bound plus confidence-based pre-verification to accept statistically-safe draft tokens, achieving **1.6× end-to-end speedup over EAGLE-3** in sampling conditions with negligible quality impact. Anchors `[[1.0.0 PRIM-7]]` Verdict Validation and `[[1.0.0 PRIM-21]]` Migration Strategy Selection.

## Mechanism (Q2)

Two coupled mechanisms compose the acceptance decision:

1. **Jensen-Shannon divergence bound** — the algorithm computes a distribution-level distance between the draft model and the target model at each candidate token position. Accept only when the JS divergence falls under a probabilistic threshold that bounds the worst-case distributional drift.
2. **Confidence-based pre-verification** — before the JS test, a cheap confidence filter rejects draft tokens whose top-1 probability is below an empirical floor, eliminating obvious mismatches without invoking the more expensive divergence computation.

The framework is **training-free**: no fine-tuning of either draft or target model is required, and the risk bound is a statistical guarantee derived from the divergence measure. Under sampling (T > 0), speculative decoding typically degrades because exact distribution matching breaks down; ARC-Decode's JS-bound approach recovers most of the lost speedup by accepting tokens that are statistically safe rather than identical.

## Anchoring (Q3)

| Primitive | Pre-wave-18 behaviour | ARC-Decode substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Greedy or exact-match speculative acceptance; sampling-mode speculative decoding loses most of its speedup | Risk-bounded acceptance via JS divergence + confidence pre-filter; statistical-safety test instead of exact-match |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Greedy-vs-sampling is a binary toggle; no middle ground under speculative decoding | Continuous risk-budget knob: the JS threshold exposes a speed-quality trade-off curve |

## Hop-1 Citations

- `[[1.0.0 P-57]]` Leviathan et al. 2023 *Fast Inference from Transformers via Speculative Decoding* (ICML 2023) — original speculative decoding.
- `[[1.0.0 P-64]]` Li et al. 2024 *EAGLE: Speculative Sampling Requires Rethinking Feature Uncertainty* (EAGLE v1).
- `[[1.0.0 P-78]]` / `[[1.0.0 P-108]]` Li et al. 2025 *EAGLE-3: Scaling up Inference Acceleration of Large Language Models via Training-Time Test* (NeurIPS 2025) — EAGLE-3 baseline that ARC-Decode outperforms.
- `[[1.0.0 P-114]]` Cai et al. 2024 *Medusa: Simple LLM Inference Acceleration Framework with Multiple Decoding Heads* (heads-based alternative).

## Hop-2 Citations

- Chen et al. 2023 *Accelerating Large Language Model Decoding with Speculative Sampling* (speculative sampling under T=0).
- Sun et al. 2024 *Spectr: Fast Speculative Decoding via Optimal Transport* (OT-based acceptance, ancestor of risk-bounded acceptance).

## MAgHARCM integration

- **YAML config key**: `inference.arc_decode.risk_budget: <float>` (default `0.05`); `inference.arc_decode.enabled: true`.
- **Implementation file**: `internal/inference/arc_decode.go::AcceptToken` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-21]]`.

## Caveats

- **Risk-budget calibration**: the JS threshold is dataset- and temperature-dependent; production deployment needs an offline calibration pass to set the bound without quality regression.
- **Draft-model pairing**: ARC-Decode assumes a high-quality draft model is available. MAgHARCM's existing `internal/inference/speculative.go` provides the pairing layer but must be extended to expose the JS-bound hooks.
- **Numerical precision**: JS divergence under fp16/bf16 can introduce drift; the implementation should compute the bound in fp32 to avoid spurious rejections near the threshold.

## Source

- arXiv: forthcoming.
- OpenReview: `jhJjW2DFKD` / `0K57Wtg15V`.
