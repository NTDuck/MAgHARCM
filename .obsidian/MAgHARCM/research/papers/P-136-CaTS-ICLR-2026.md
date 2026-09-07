---
title: "P-136 CaTS: Calibrated Test-Time Scaling for Efficient LLM Reasoning"
backlink: "[[1.0.0 P-136]]"
aliases:
  - "1.0.0 P-136"
  - "P-136"
  - "P-136-CaTS-ICLR-2026"
  - "CaTS-ICLR-2026"
  - "huang2026cats"
tags: [paper, test-time-scaling, calibration, slm-as-judge, "[[1.0.0 PRIM-7]]", wave-19]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICLR 2026 (Poster)
---

# [[1.0.0 P-136]] CaTS

## TL;DR

CaTS distills **self-consistency-derived confidence into the SLM itself** via Self-Calibration, enabling reliable single-forward-pass confidence estimation that drives adaptive early-stopping TTS (CaTS-ES) and adaptive best-of-N (CaTS-BoN). Improves MathQA accuracy from **73.7 to 83.6** at sample budget 16. Anchors `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation as a frontier-PRM-replacement primitive.

## Mechanism (Q2)

1. **Self-Calibration** — a one-time distillation pass takes self-consistency-derived confidence scores (from a small calibration set) and trains the SLM to output calibrated confidence in a single forward pass.
2. **CaTS-ES (early-stopping TTS)** — when the calibrated confidence exceeds a threshold, the model commits the answer without sampling further; saves compute at no accuracy cost.
3. **CaTS-BoN (adaptive best-of-N)** — confidence drives dynamic sample count: low-confidence instances get more samples; high-confidence instances stop early.

Together these give MAgHARCM a **two-tier verifier architecture** when combined with P-127 SLM-as-a-Judge: cheap confidence pre-filter (CaTS) + multi-judge panel for hard cases.

## Anchoring (Q3)

| Primitive | Pre-wave-19 behaviour | CaTS substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation | Each candidate required a full multi-judge evaluation; no cheap pre-filter | Self-Calibrated confidence gives a single-forward-pass pre-filter; the multi-judge panel fires only on low-confidence cases |

## Hop-1 Citations

- P-52 Wang 2023 Self-Consistency.
- P-91 Snell 2024 Test-Time Compute.
- P-98 Brown 2024 LLM Monkeys.

## Hop-2 Citations

- P-92 Lightman 2023 PRM (process reward modeling).
- P-96 Kojima 2022 Zero-Shot CoT (CoT prompting foundation).
- Guo 2017 Calibration (temperature scaling).

## MAgHARCM integration

- **YAML config key**: `tts.cats.enabled: true`; `tts.cats.confidence_threshold: <float in [0,1]>`.
- **Implementation file**: `internal/agents/cats.go::CalibratedConfidence` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-7]]`.

## Caveats

- **Calibration set cost** — Self-Calibration requires a small calibration set; cost is amortised across deployment but must be re-run on domain shift.
- **Calibration drift** — long-lived deployments must recalibrate; calibrated confidence can degrade under data drift.
- **Threshold tuning** — confidence_threshold must be tuned per-task; aggressive thresholds risk false-confidence on hard cases.

## Source

- arXiv: openreview jrSc4RJXy1.
- Venue: ICLR 2026 Poster (verified at iclr.cc/virtual/2026/poster/10007848; ICLR proceedings hash 9f9ecbf4062842df17ec3f4ea3ad7f54).
