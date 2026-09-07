---
title: "P-125 T1: Tool-Integrated Verification for Test-time Compute Scaling in Small Language Models"
backlink: "[[1.0.0 P-125]]"
aliases:
  - "1.0.0 P-125"
  - "P-125"
  - "P-125-T1-Tool-Integrated-Verification-2026"
  - "T1-Tool-Integrated-Verification-2026"
  - "kang2026t1"
tags: [paper, slm, test-time-scaling, tool-integration, verification, "[[1.0.0 PRIM-7]]", "[[1.0.0 P-91]]", "[[1.0.0 P-92]]", "[[1.0.0 P-97]]", "[[1.0.0 P-98]]", "[[1.0.0 P-84]]", "[[1.0.0 P-104]]", wave-18]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICLR 2026
---

# [[1.0.0 P-125]] T1

## TL;DR

T1 is a **two-stage verification framework** that lets a 1B-parameter SLM match or exceed an 8B model on math reasoning by offloading memorization-heavy steps to an external code interpreter and using the SLM only for final verification on a filtered candidate subset. Llama-3.2 1B with T1 outperforms Llama-3.1 8B on MATH and improves both PRM and critic-model verification accuracy. Anchors `[[1.0.0 PRIM-7]]` Verdict Validation.

## Mechanism (Q2)

The framework decomposes verification into two stages:

1. **Stage 1 — external-tool filtering**: a code interpreter (Python REPL) executes candidate solutions and filters those whose numerical or symbolic computation can be confirmed by code execution. This offloads memorization-heavy steps (arithmetic, algebraic simplification) away from the SLM.
2. **Stage 2 — SLM final verification**: the SLM performs verification only on the filtered subset, focusing its limited capacity on reasoning-heavy steps that cannot be mechanized.

The architecture exploits an asymmetry: code execution is cheap and reliable for procedural sub-tasks, while learned language reasoning is needed for the higher-order semantic checks. By routing the procedural work to the tool, the SLM becomes a **specialized verifier** rather than a general solver.

## Anchoring (Q3)

| Primitive | Pre-wave-18 behaviour | T1 substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Verdict Validation | PRM/critic scores all candidates with no task-routing; small PRMs underperform large PRMs | Tool-filtered then SLM-verified: code interpreter handles procedural steps, SLM handles semantic verification on filtered subset |

## Hop-1 Citations

- `[[1.0.0 P-92]]` Lightman et al. 2023 *Let's Verify Step by Step* (PRM training, process reward models).
- `[[1.0.0 P-91]]` Snell et al. 2024 *Scaling LLM Test-Time Compute Optimally* (compute allocation framework).
- `[[1.0.0 P-97]]` Welleck et al. 2024 *Self-Correcting Language Models* (critic/verifier-as-generator).
- `[[1.0.0 P-98]]` Brown et al. 2024 *Large Language Monkeys* (best-of-N sampling baseline).

## Hop-2 Citations

- `[[1.0.0 P-84]]` Muennighoff et al. 2025 *s1: Simple Test-Time Scaling* (budget forcing).
- `[[1.0.0 P-104]]` Li et al. 2025 *S\*: Test Time Scaling for Code Generation* (hybrid TTS for code).

## MAgHARCM integration

- **YAML config key**: `verifier.t1_enabled: true` (opt-in, default `false`).
- **Implementation file**: `internal/verifier/t1_tool_filter.go::FilterAndVerify` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-7]]`.

## Caveats

- **Code interpreter dependency**: Stage 1 requires an available sandboxed Python interpreter. MAgHARCM's existing `internal/languages/oracle.go` can be reused, but the verifier pipeline must preserve tool output between filter and SLM-verification stages.
- **Domain transferability**: T1 is demonstrated on MATH; generalization to non-numeric domains (e.g. natural-language reasoning) depends on whether a comparable external tool exists to mechanize procedural sub-tasks.
- **Filter-pass cost**: the code interpreter execution is the new bottleneck when most candidates pass Stage 1; budget allocation must account for tool-call latency in the overall TTS compute budget.

## Source

- arXiv:2504.04718 — T1 (ICLR 2026).
- OpenReview: `tBkLWfmugI`.
