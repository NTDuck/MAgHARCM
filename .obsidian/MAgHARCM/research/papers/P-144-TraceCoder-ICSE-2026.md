---
title: "P-144 TraceCoder: A Trace-Driven Multi-Agent Framework for Automated Debugging of LLM-Generated Code"
backlink: "[[1.0.0 P-144]]"
aliases:
  - "1.0.0 P-144"
  - "P-144"
  - "P-144-TraceCoder-ICSE-2026"
  - "TraceCoder-ICSE-2026"
  - "tracecoder2026icse"
tags: [paper, multi-agent, debugging, trace-instrumentation, role-flip-reviewer, "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-25]]", wave-20]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICSE 2026 (Research Track)
---

# [[1.0.0 P-144]] TraceCoder

## TL;DR

TraceCoder is a **trace-driven multi-agent debugging framework** that emulates the observe-analyze-repair loop used by human software engineers. Four mechanisms: (1) **runtime instrumentation** captures fine-grained execution traces beyond binary pass/fail; (2) **causal analysis** localises root causes from traces; (3) **Historical Lesson Learning Mechanism (HLLM)** distills insights from previous failed repair attempts to prevent repetitive cycles; (4) **Rollback Mechanism (RM)** enforces that every iterative repair step is a strict improvement. Up to **34.43% relative Pass@1 improvement** over advanced baselines; iterative repair alone contributes **65.61% relative gain**. Anchors `[[1.0.0 PRIM-22]]` (Observation phase via trace instrumentation) and `[[1.0.0 PRIM-25]]` (Role-Flip Reviewer: HLLM is the cross-iteration review substrate).

## Mechanism (Q2)

1. **Runtime instrumentation** — diagnostic probes injected into the candidate code capture fine-grained execution traces (function calls, control-flow paths, intermediate variable values). Provides deep visibility into *why* a failure occurs, not just *whether* it occurs.
2. **Causal analysis** — a multi-agent system consumes the traces and performs root-cause localisation.
3. **Historical Lesson Learning Mechanism (HLLM)** — distills insights from previous failed repair attempts; informs future correction strategies with past lessons; prevents repetitive, inefficient cycles.
4. **Rollback Mechanism (RM)** — enforces that every iterative repair step is a strict improvement toward a correct solution; rejects non-monotonic steps to ensure convergence.

Together these give MAgHARCM a *trace-driven* debugging substrate that complements P-133 ADI's Frame Lifetime Trace (line-level DA) by adding HLLM (cross-iteration lesson learning) and RM (convergence enforcement).

## Anchoring (Q3)

| Primitive | Pre-wave-20 behaviour | TraceCoder substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Observation phase) | Observation relied on compiler errors + test pass/fail signals; no fine-grained execution trace | TraceCoder's runtime instrumentation captures the full execution path, giving Observation phase the trace substrate that P-133 ADI frames for line-level DA |
| `[[1.0.0 PRIM-25]]` Role-Flip Reviewer | Reviewer used a fresh-context pass; no memory of previous failed iterations | HLLM distills lessons across iterations; reviewer now has cross-iteration context for more grounded feedback |

## Hop-1 Citations

- P-133 ADI (FSE 2026, Frame Lifetime Trace — sister mechanism for function-level DA).
- Execution-trace analysis in debuggers (gdb, lldb — classical substrate).
- SWE-bench Verified (OpenAI 2024, P-109) — benchmark lineage for the Pass@1 evaluation.
- OpenHands CodeAct (Wang et al. 2024, P-115) — multi-agent debugging scaffold baseline.

## Hop-2 Citations

- Multi-agent debate / role-flip pattern (Du et al. 2023, Liang et al. 2023).
- Rollback / monotonic-improvement search (classic CS, e.g., Levesque 1984).
- Classical execution tracing (Knuth 1971, Sedgewick 1986).

## MAgHARCM integration

- **YAML config key**: `debugging.tracecoder.enabled: true`; `debugging.tracecoder.hllm.max_lessons: <int>`; `debugging.tracecoder.rollback.strict_monotonic: true`.
- **Implementation file**: `internal/debugging/tracecoder.go::NewTraceCoderAgent` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-25]]`.

## Caveats

- **Instrumentation cost** — runtime probes add execution overhead; production runs may need to disable instrumentation in the final verification step.
- **Lesson staleness** — HLLM lessons accumulate; periodic pruning is required to prevent stale lessons from dominating.
- **Rollback scope** — RM operates at the patch level; finer-grained rollback (e.g., per-statement) requires a richer diff substrate.

## Source

- arXiv: 2602.06875.
- Venue: ICSE 2026 (Research Track), verified via researchr.org/details/icse-2026/icse-2026-research-track/145; CSMA-Research-Group/TraceCoder GitHub repository.
