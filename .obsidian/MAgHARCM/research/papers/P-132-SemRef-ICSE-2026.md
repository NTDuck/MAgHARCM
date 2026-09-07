---
title: "P-132 SemRef: Semantic-Enhanced Automatic Refinement of Architecture Recovery Results Using LLMs"
backlink: "[[1.0.0 P-132]]"
aliases:
  - "1.0.0 P-132"
  - "P-132"
  - "P-132-SemRef-ICSE-2026"
  - "SemRef-ICSE-2026"
  - "zhang2026semref"
tags: [paper, architecture-recovery, refinement, llm, "[[1.0.0 PRIM-9]]", "[[1.0.0 PRIM-31]]", wave-18]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICSE 2026 (Research Track)
---

# [[1.0.0 P-132]] SemRef

## TL;DR

SemRef is a **post-hoc refinement loop**: take the imprecise output of an existing architecture-recovery tool, feed it to an LLM together with explicit + implicit structural dependencies, and let the LLM rebalance the partitioning. Achieves **+17.72-43.35% normalized accuracy** over five metrics, including **+118.57% MoJoFM and +100.41% a2a_adj**. Anchors `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement.

## Mechanism (Q2)

The pipeline is a refinement loop with four stages:

1. **Baseline recovery** — run any existing architecture-recovery tool (structural, semantic, or hybrid) to produce an initial partition. The partition may be imprecise but serves as a starting point.
2. **Evidence assembly** — collect explicit dependencies (call/import graph) and implicit dependencies (naming, co-change, semantic co-occurrence) for every file. The LLM is shown both the initial partition and the assembled evidence.
3. **LLM rebalance** — the LLM proposes re-assignments (move file F from module A to module B) when the evidence supports the move. The LLM operates on a structured representation, not raw code.
4. **Iterate** — re-run the baseline recovery on the rebalanced partition and repeat until a convergence criterion (no further moves, or accuracy ceiling) is met.

The key insight is **composability**: SemRef is a drop-in refinement layer that does not require re-implementing the underlying recovery tool. Any off-the-shelf tool can be lifted to LLM-refinement quality by wrapping it in the loop.

## Anchoring (Q3)

| Primitive | Pre-wave-18 behaviour | SemRef substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Single-pass recovery; outputs are final | LLM-refinement loop on top of any recovery tool; explicit + implicit dependency evidence drives rebalance |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Single-pass retrieval of architectural boundaries | Iterative rebalance loop with LLM-as-judge; convergence-based termination |

## Hop-1 Citations

- Baseline recovery tools (Maqbool et al. 2014; various 2018-2024 architectures).
- MoJoFM / a2a / a2a_adj / c2cavg / MoJoFM2 — recovery quality metrics.
- LLM-as-judge lineage (Zheng et al. 2023).

## Hop-2 Citations

- Iterative refinement (Shinn et al. 2023 *Reflexion*; Madaan et al. 2023 *Self-Refine*) — ancestor self-refinement pattern.
- Evidence-grounded LLM reasoning (Bohnet et al. 2024; Wei et al. 2022 CoT).

## MAgHARCM integration

- **YAML config key**: `architecture_recovery.semref.enabled: true` (opt-in, default `false`); `architecture_recovery.semref.max_iters: <int>`.
- **Implementation file**: `internal/code_graph/semref.go::RefinePartition` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-9]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **LLM hallucination risk**: the LLM may propose moves that improve a metric while degrading another; multi-metric guardrails are needed.
- **Iteration cost**: each iteration is an LLM call plus a re-run of the baseline tool; convergence may be slow on large repos.
- **Composability caveat**: SemRef's gains depend on the quality of the input partition. A wildly wrong baseline may require multiple loops to recover, and may never fully converge.

## Source

- URL: conf.researchr.org/details/icse-2026/icse-2026-research-track/72.
