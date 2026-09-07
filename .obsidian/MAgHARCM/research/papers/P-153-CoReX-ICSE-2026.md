---
title: "P-153 CoReX: Context-Aware Refinement-Based Slicing for Debugging Regression Failures"
backlink: "[[1.0.0 P-153]]"
aliases:
 - "1.0.0 P-153"
 - "P-153"
 - "P-153-CoReX-ICSE-2026"
 - "CoReX-ICSE-2026"
 - "corex2026icse"
tags: [paper, program-slicing, regression-debugging, context-refinement, strict-mechanism, "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-31]]", wave-23]
date: 2026-09-07
last_updated: 2026-09-07 (iter-7, wave-23)
venue: ICSE 2026 (Research Track)
---

# [[1.0.0 P-153]] CoReX

## TL;DR

CoReX is a **context-aware refinement-based slicing** mechanism for **regression-failure debugging**. Instead of static or observed-dependency-only slicing (e.g., SliceMate P-rejected-Wave-19), CoReX iteratively refines the candidate slice using observed execution behaviour and recent change history, conditioned on the regression failure context. Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Structure phase — function-level DA pattern extended from P-133 ADI) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (refinement-based slicing = iterative retrieval refinement substrate). **Closes the strict program-comprehension-mechanism slot carried forward from Wave-17 first opening.**

## Mechanism (Q2)

1. **Context-conditioned refinement** — instead of producing a single static slice or an observed-dependency slice, CoReX conditions the slice on the regression failure context (the failing test + the recent change set). The initial candidate slice is refined iteratively using (a) observed execution behaviour at the failing test, (b) recent change history, and (c) regression-relevant control-flow landmarks.
2. **Iterative refinement loop** — each refinement step narrows the slice to the subset of code whose behaviour actually changed under the failing test, removing false-positive slicing edges that appear in static or observed-only slices.
3. **Regression-failure localisation** — the refined slice localises the regression-causing code with higher precision than vanilla slicing (paper reports F1 / MAP improvements over baselines on Java regression-failure benchmarks).
4. **Function-level granularity** — the refinement step operates at function-level granularity, matching the function-level DA pattern from P-133 ADI rather than line-level gdb-style DA.

Together these give MAgHARCM a **strict mechanism** for program-comprehension substrate in the regression-debugging setting — the first context-aware refinement-based slicing paper in the SLM-era comprehension substrate.

## Anchoring (Q3)

| Primitive | Pre-wave-23 behaviour | CoReX substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Comprehension anchored on P-128 KVzip, P-129 LλMDA (LLM-aided partial PDG), P-133 ADI (Frame Lifetime Trace), P-139 TypePro (inter-procedural SDG slicing), P-140 Panta (iterative test generation), P-142 NESA (Datalog policy), P-150 TestPrune (coverage-driven context pruning), P-151 SmartC2Rust (chunked translation), P-152 Hallu-Eval (hallucination benchmark); **strict-mechanism residual slot still open** | CoReX as the **strict-mechanism closure** for the residual slot. Closes Wave-17 first-opened program-comprehension-mechanism gap. |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Retrieval refinement anchored on P-117 RepoCoder, P-122 ReasoningBank, P-128 KVzip, P-134 RelayCaching, P-137 SuffixDecoding, P-141 KVFlow, P-147 SpecKV, P-148 LookaheadKV, P-150 TestPrune; no refinement-based slicing substrate | CoReX refinement-based slicing as iterative retrieval refinement substrate = partition-aligned summary-pass pattern. |

## Hop-1 Citations

- P-133 ADI (ICSE/FSE 2026, Wave-18) — Frame Lifetime Trace + function-level navigational commands. CoReX extends the function-level DA pattern to the regression-debugging setting with refinement-based slicing.
- P-129 LλMDA (ICSE 2026, Wave-18) — LLM-aided partial program dependence analysis. CoReX shares the "augment static with execution + refinement" pattern.
- P-139 TypePro (Wave-19) — inter-procedural slicing for type inference. CoReX extends the inter-procedural slicing lineage to regression failure debugging.
- P-140 Panta (Wave-19) — iterative hybrid static+dynamic test generation. CoReX shares the iterative-refinement loop pattern.

## Hop-2 Citations

- Program slicing (Weiser 1981 IEEE TSE) — foundational program-slicing substrate. CoReX is a context-aware refinement-based extension of the Weiser slicing substrate.
- Context-sensitive interprocedural slicing (Hind 2001) — context-sensitive slicing lineage. CoReX inherits the context-sensitivity machinery.
- Regression fault localization (Zeller 2009) — foundational regression-debugging substrate. CoReX is the slicing-specific instance of the regression-debugging lineage.
- Slicing for debugging (Agrawal 1993) — early slicing-for-debugging work. CoReX extends the slicing-for-debugging lineage with refinement-based context conditioning.

## MAgHARCM integration

- **YAML config key**: `agents.comprehension.slicing.refinement_context: "regression_failure"`; `agents.comprehension.slicing.refinement_iterations: <int>`; `agents.comprehension.slicing.coarseness: "function_level"`.
- **Implementation file**: `internal/comprehension/slicing.go::NewContextAwareRefiner` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **Java-only benchmark** — CoReX is evaluated on Java regression-failure benchmarks; generalisation to Go / Rust regression failures is future work.
- **Refinement cost** — the iterative refinement loop adds compile + execute cost per iteration; must be amortised over the localisation-quality gain.
- **Context-conditioning risk** — the regression failure context biases the slice toward the failing test; cross-test generalisation is not in scope.

## Source

- Venue: ICSE 2026 (Research Track), verified via conf.researchr.org/details/icse-2026/icse-2026-research-track/60. Held April 12-18 2026 in Rio de Janeiro.
- DOI: pending ACM assignment (researchr slot #60; check ACM DL after camera-ready publication).
- arXiv: not surfaced (no arXiv preprint confirmed); ICSE 2026 is the canonical venue.

## BibTeX

```
@inproceedings{corex2026icse,
  title     = {CoReX: Context-Aware Refinement-Based Slicing for Debugging Regression Failures},
  author    = {Badihi, Sahar and Rubin, Julia},
  booktitle = {Proceedings of the 48th IEEE/ACM International Conference on Software Engineering (ICSE)},
  year      = {2026},
  address   = {Rio de Janeiro, Brazil}
}
```
