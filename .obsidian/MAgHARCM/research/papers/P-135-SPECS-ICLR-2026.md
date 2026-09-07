---
title: "P-135 SPECS: Faster Test-Time Scaling through Speculative Drafts"
backlink: "[[1.0.0 P-135]]"
aliases:
  - "1.0.0 P-135"
  - "P-135"
  - "P-135-SPECS-ICLR-2026"
  - "SPECS-ICLR-2026"
  - "cemri2026specs"
tags: [paper, test-time-scaling, speculative-decoding, slm-as-judge, "[[1.0.0 PRIM-7]]", "[[1.0.0 PRIM-21]]", wave-19]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICLR 2026 (Conference)
---

# [[1.0.0 P-135]] SPECS

## TL;DR

SPECS combines **speculative drafts** (small SLM proposes reasoning steps) with **reward-guided soft verification** and a **dynamic draft↔target switch** based on a continuous budget knob. Achieves **18-19% latency reduction at matched accuracy** vs. beam-search TTS. Closes Wave-18 R3 watchlist. Anchors `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation and `[[1.0.0 PRIM-21]]` Migration Strategy Selection.

## Mechanism (Q2)

1. **Speculative drafts** — a smaller SLM proposes k reasoning steps ahead of the target model's verification pass; the target evaluates each step under a soft verifier (reward model head) rather than a binary accept/reject.
2. **Reward-conditioned dynamic switch** — a continuous budget knob (analogous to ARC-Decode's risk-budget, P-126) decides at every step whether to use the draft or fall back to the full target-model pass. Low-confidence steps trigger immediate fallback.
3. **Soft verification** — replaces the rigid accept/reject of classical speculative decoding with a graded confidence that feeds back into the dynamic switch.

Together these give MAgHARCM's `internal/agents/strategy.go::TryInOrder` registry a concrete mechanism for inserting a speculative-decoding layer ahead of the multi-judge verdict panel.

## Anchoring (Q3)

| Primitive | Pre-wave-19 behaviour | SPECS substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation | Fixed multi-judge panel; no early-exit on high-confidence drafts | Draft + soft-verifier pre-filter + dynamic switch replaces the panel's rigid evaluation; hard cases still flow to the full panel |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | TryInOrder registry had no continuous budget knob — strategies were discrete (try, fallback, escalate) | Continuous draft↔target budget knob analogous to ARC-Decode's risk-budget; integrates with the registry's existing escalation rules |

## Hop-1 Citations

- P-52 Wang 2023 Self-Consistency.
- P-91 Snell 2024 Test-Time Scaling.
- P-92 Lightman 2023 PRM.

## Hop-2 Citations

- P-57 Leviathan 2023 Speculative Decoding.
- P-78 / P-108 EAGLE-3 2025 — draft-model lineage.
- P-98 Brown 2024 LLM Monkeys (best-of-N).

## MAgHARCM integration

- **YAML config key**: `tts.specs.enabled: true`; `tts.specs.budget_knob: <float in [0,1]>`.
- **Implementation file**: `internal/agents/specs.go::SpeculativeDraft` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-21]]`.

## Caveats

- **Draft model selection** — the small SLM must share tokenizer and a reasonable embedding space with the target; mismatched tokenizers break the soft verifier.
- **Soft verifier calibration** — out-of-distribution reasoning domains degrade the reward signal; periodic recalibration is needed for long-lived deployments.
- **Latency-accuracy frontier** — 18-19% latency reduction is at matched accuracy; aggressive budget-knob settings risk accuracy regression.

## Source

- arXiv: 2506.15733.
- Venue: ICLR 2026 (OpenReview 5XixaecZ8W; verified via aimodeling.com acceptance news).
