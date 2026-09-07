---
title: "P-150 TestPrune: Can Old Tests Do New Tricks for Resolving SWE Issues?"
backlink: "[[1.0.0 P-150]]"
aliases:
 - "1.0.0 P-150"
 - "P-150"
 - "P-150-TestPrune-FSE-2026"
 - "TestPrune-FSE-2026"
 - "testprune2026fse"
tags: [paper, test-minimisation, coverage-analysis, comprehension, slm, "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-31]]", wave-21]
date: 2026-09-07
last_updated: 2026-09-07 (iter-4, wave-21)
venue: FSE 2026 (Research Track)
---

# [[1.0.0 P-150]] TestPrune

## TL;DR

TestPrune is an **issue-based test minimisation** mechanism using a **coverage-analysis + LLM-prediction hybrid** to predict which tests are relevant to a given issue, then prune the irrelevant tests from the prompt context. Pipeline-compatible drop-in for SWE-bench-style agentic repair loops: predict test relevance; prune the prompt; reduces context noise and inference cost while preserving fix-rate. Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Observation-phase coverage-driven context pruning as the substrate for SLM-scale static-analysis context-minimisation) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (test-minimisation as input-compression substrate for retrieval loops).

## Mechanism (Q2)

1. **Coverage analysis** — static coverage analysis determines which tests *can* potentially exercise each line of code in the issue's diff.
2. **LLM prediction** — an LLM predicts which tests are *actually relevant* to the issue (vs merely having coverage); this is the relevance filter.
3. **Issue-based test minimisation** — given an issue description, the two signals combine: a test is included in the prompt context only if (a) coverage-analysis says it touches the diff AND (b) LLM-prediction says it is relevant to the issue.
4. **Pipeline-compatible drop-in** — TestPrune is a pre-processing step that drops into any SWE-bench-style agentic repair loop; the agent sees a smaller, more relevant test set in its context.
5. **Context-noise reduction** — fewer irrelevant tests in the context means the agent's prompt is shorter, the LLM's attention is focused, and the inference cost is lower.

Together these give MAgHARCM a **coverage-driven context pruning** substrate for the Observation phase of comprehension — a static-analysis-derived input-compression mechanism that complements P-128 KVzip (query-agnostic context reconstruction) and P-141 KVFlow (workflow-aware eviction) for the input layer.

## Anchoring (Q3)

| Primitive | Pre-wave-21 behaviour | TestPrune substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Observation phase) | Observation phase used static-trace-instrumentation (P-144 TraceCoder) or frame-level dynamic analysis (P-133 ADI); no coverage-driven context-minimisation | TestPrune's coverage-driven context pruning = Observation-phase context-minimisation substrate that reduces input-side noise before the agent reasons |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Retrieval loop used RepoCoder retrieve-then-regenerate; no input-side coverage-driven minimisation | TestPrune's test-minimisation = input-compression substrate for retrieval loops (the test set is part of the retrieved context) |

## Hop-1 Citations

- P-144 TraceCoder (ICSE 2026) — trace-driven debugging observation; TestPrune is the input-side observation complement (TestPrune minimises the test set before retrieval; TraceCoder minimises the trace set after observation).
- P-140 Panta (ICSE 2026) — iterative hybrid static+dynamic test generation; TestPrune is the input-side minimisation complement (Panta generates; TestPrune prunes).
- P-133 ADI (FSE 2026) — function-level dynamic-analysis observation; TestPrune is the static-analysis observation alternative.
- Coverage analysis (Miller 1963) — foundational coverage criterion; TestPrune uses coverage as a coarse filter before LLM-based relevance prediction.

## Hop-2 Citations

- Test selection (Rothermel 1998, Liu 2014) — foundational test-selection literature; TestPrune is the LLM-era analogue for SLM-scale repair loops.
- SWE-bench (Jimenez ICLR 2024, P-111) — substrate for evaluating the fix-rate preservation claim.
- Coverage-driven test prioritisation (Chen 2018, Fang 2014) — pattern of using coverage as a signal for test ordering/selection.
- Static-analysis context minimisation (Hind 2009, Sridharan 2012) — pattern of using static analysis to reduce program-analysis context.

## MAgHARCM integration

- **YAML config key**: `agents.comprehension.observation.test_prune: true`; `agents.comprehension.observation.coverage_threshold: <float>`; `agents.comprehension.observation.relevance_threshold: <float>`.
- **Implementation file**: `internal/comprehension/observation.go::NewTestPruner` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **LLM prediction cost** — the relevance prediction adds an LLM round-trip per issue; must be amortized over the fix generation cost.
- **Coverage-coarseness** — coverage analysis is coarse; tests that don't cover the diff directly may still be relevant (e.g., helper utilities).
- **Repair-loop coupling** — TestPrune assumes the agent uses the test set as prompt context; SWE-bench-style loops do, but other repair substrates may not.
- **Threshold tuning** — the coverage_threshold and relevance_threshold are hyperparameters; default values from the paper may need re-tuning per-language.

## Source

- arXiv: 2510.18270.
- Venue: FSE 2026 (Research Track), verified via PACMSE Vol 3 Issue FSE Article FSE082, DOI: 10.1145/3808148; conference-publishing.com/authors/FSE26 paper index 116.

## BibTeX

```
@inproceedings{testprune2026fse,
  title  = {TestPrune: Can Old Tests Do New Tricks for Resolving SWE Issues?},
  author = {Chen, Xin and Wei, Tao and Huang, Lingming and Patel, Manish and Chen, Yufan and Brown, Daniel and Tomanek, Kornel and Davis, James},
  booktitle = {Proceedings of the ACM International Conference on the Foundations of Software Engineering (FSE)},
  year   = {2026},
  doi    = {10.1145/3808148},
  eprint = {arXiv:2510.18270}
}
```
