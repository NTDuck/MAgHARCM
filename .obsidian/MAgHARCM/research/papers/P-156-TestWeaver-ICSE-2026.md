---
title: "P-156 TestWeaver: Execution-aware, Feedback-driven Regression Testing Generation with Large Language Models"
backlink: "[[1.0.0 P-156]]"
aliases:
 - "1.0.0 P-156"
 - "P-156"
 - "P-156-TestWeaver-ICSE-2026"
 - "TestWeaver-ICSE-2026"
 - "testweaver2026icse"
tags: [paper, test-generation, llm-augmented, backward-slicing, execution-annotation, coverage-plateau, "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-29]]", "[[1.0.0 PRIM-31]]", wave-23]
date: 2026-09-07
last_updated: 2026-09-07 (iter-7, wave-23)
venue: ICSE 2026 (Research Track)
---

# [[1.0.0 P-156]] TestWeaver

## TL;DR

TestWeaver is an **LLM-augmented regression-test generation** framework combining three mechanisms: (a) **backward slicing** for context reduction, (b) **close-test retrieval** for execution context, (c) **execution in-line annotations** for variable-state context. Addresses the coverage-plateau problem in LLM-based regression test generation. Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (backward slicing + execution annotations extend the LLM-augmented static-analysis pattern alongside P-129 LλMDA + P-153 CoReX + P-156 itself) and `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter (close-test retrieval + execution annotation = retrieval-and-recruitment substrate for test-generation agents) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (close-test retrieval = iterative retrieval refinement substrate for test-generation agents).

## Mechanism (Q2)

1. **Backward slicing for context reduction** — instead of providing the full program context to the LLM, TestWeaver provides a backward slice from the target line. This reduces hallucinations, keeps the LLM focused, and reduces token cost.
2. **Close-test retrieval** — TestWeaver identifies prior tests that share control-flow similarity with the path to the target line and includes them as execution context. The retrieval substrate is the LLM-augmented analogue of P-150 TestPrune's coverage-driven context pruning.
3. **Execution in-line annotations** — variable-state comments along the executed path are injected into the LLM prompt. The annotations encode the runtime state at each statement, giving the LLM ground-truth execution context.
4. **Coverage-plateau mitigation** — the three mechanisms together mitigate the coverage plateau (the point at which additional test-generation iterations stop improving coverage) by providing targeted, high-signal context to the LLM.

Together these give MAgHARCM an **LLM-augmented test-generation framework** with three concrete context-conditioning mechanisms — the first tri-component (slicing + retrieval + annotation) LLM-test-gen substrate.

## Anchoring (Q3)

| Primitive | Pre-wave-23 behaviour | TestWeaver substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Comprehension anchored on P-128 KVzip, P-129 LλMDA (LLM-aided partial PDG), P-133 ADI (Frame Lifetime Trace), P-139 TypePro (inter-procedural SDG slicing), P-140 Panta (iterative test generation), P-142 NESA (Datalog policy), P-150 TestPrune (coverage-driven context pruning), P-151 SmartC2Rust (chunked translation), P-152 Hallu-Eval (hallucination benchmark), P-153 CoReX (refinement-based slicing); test-generation substrate is anchored on P-140 Panta only | TestWeaver extends the LLM-augmented static-analysis pattern with backward slicing + close-test retrieval + execution annotation. Three new comprehension-mechanism anchors. |
| `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter | Iteration recruiter anchored on P-124 Syzygy, P-151 SmartC2Rust, P-155 CodeCureAgent (three-step validation); no close-test retrieval + execution annotation substrate for test-generation agents | TestWeaver close-test retrieval + execution annotation = retrieval-and-recruitment substrate for test-generation agents |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Retrieval refinement anchored on P-117 RepoCoder, P-128 KVzip, P-137 SuffixDecoding, P-150 TestPrune; no close-test retrieval + backward-slicing substrate for test-generation agents | TestWeaver close-test retrieval + backward slicing = iterative retrieval refinement substrate for test-generation agents |

## Hop-1 Citations

- P-129 LλMDA (ICSE 2026, Wave-18) — LLM-aided partial program dependence analysis. TestWeaver shares the LLM-augmented dependence-analysis substrate.
- P-133 ADI (ICSE/FSE 2026, Wave-18) — Frame Lifetime Trace + function-level navigational commands. TestWeaver inherits the function-level slicing substrate.
- P-140 Panta (Wave-19) — iterative hybrid static+dynamic test generation. TestWeaver extends the iterative test-generation substrate with LLM-augmented slicing + retrieval + annotation.
- P-150 TestPrune (FSE 2026, Wave-21) — coverage-driven context pruning. TestWeaver shares the coverage-driven pruning substrate but adds LLM augmentation.

## Hop-2 Citations

- Backward slicing (Weiser 1981 IEEE TSE) — foundational slicing substrate. TestWeaver uses backward slicing for context reduction.
- Evosuite execution-aware test generation (Fraser 2011) — execution-aware test-generation substrate. TestWeaver inherits the execution-aware lineage.
- LLM test-generation lineage (Chen 2024) — LLM-based test-generation substrate. TestWeaver is the first tri-component (slicing + retrieval + annotation) instance.
- Coverage-plateaus (Rojas 2023) — coverage-plateau analysis. TestWeaver's three mechanisms directly target the plateau.

## MAgHARCM integration

- **YAML config key**: `agents.testgen.backward_slicing: true`; `agents.testgen.close_test_retrieval: true`; `agents.testgen.execution_annotations: true`; `agents.testgen.coverage_plateau_threshold: <int>`.
- **Implementation file**: `internal/testgen/testweaver.go::NewTestWeaver` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-29]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **Backward-slicing dependency** — the backward-slicing step requires a slicing substrate (Weiser-style); projects without compiler-grade dependency info are out of scope.
- **Execution-trace availability** — execution in-line annotations require the program to be executable with tracing enabled; non-executable code is out of scope.
- **Close-test retrieval cost** — the close-test retrieval step adds a control-flow-similarity computation per candidate; the cost grows with the test-suite size.
- **Multi-language generalisation** — TestWeaver is evaluated on Java projects; generalisation to Go / Rust / Python is future work.

## Source

- Venue: ICSE 2026 (Research Track), verified via dl.acm.org/doi/10.1145/3744916.3787805. Held April 12-18 2026 in Rio de Janeiro, Brazil.
- arXiv: 2508.01255 (preprint).
- Authors' affiliation: FPT Software AI Center / VinUniversity / Nanyang Technological University / University of Texas at Dallas.

## BibTeX

```
@inproceedings{testweaver2026icse,
  title     = {TestWeaver: Execution-aware, Feedback-driven Regression Testing Generation with Large Language Models},
  author    = {Le, Cuong Chi and Van, Cuong Duc and Vu, Tung Duy and Pham, Minh V. T. and Phan, Hoang N. and Phan, Huy N. and Nguyen, Tien N.},
  booktitle = {Proceedings of the 48th IEEE/ACM International Conference on Software Engineering (ICSE)},
  year      = {2026},
  doi       = {10.1145/3744916.3787805}
}
```
