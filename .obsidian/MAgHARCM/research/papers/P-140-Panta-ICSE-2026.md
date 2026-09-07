---
title: "P-140 Panta: LLM Test Generation via Iterative Hybrid Program Analysis"
backlink: "[[1.0.0 P-140]]"
aliases:
  - "1.0.0 P-140"
  - "P-140"
  - "P-140-Panta-ICSE-2026"
  - "Panta-ICSE-2026"
  - "gu2026panta"
tags: [paper, test-generation, hybrid-analysis, coverage-driven, "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-21]]", wave-19]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICSE 2026 (Research Track)
---

# [[1.0.0 P-140]] Panta

## TL;DR

Panta drives LLM test generation via **iterative hybrid program analysis** — static control-flow analysis (cyclomatic-complexity-targeted path ranking) combined with dynamic code coverage analysis (instrumented run feedback) steers the LLM to generate test cases targeting uncovered branches. Achieves **26% higher line coverage** and **23% higher branch coverage** than SOTA. Anchors `[[1.0.0 PRIM-22]]` (Search + Explanation phases get a static-flow ∪ dynamic-coverage signal) and `[[1.0.0 PRIM-21]]` (coverage-driven test prioritisation knob).

## Mechanism (Q2)

1. **Static control-flow ranking** — cyclomatic-complexity-targeted path ranking identifies hard-to-cover branches; the LLM is steered toward those paths first.
2. **Dynamic coverage feedback** — instrumented run feedback closes the loop: tests that miss targeted branches trigger a re-prompt with the uncovered branches as context.
3. **Iterative hybrid loop** — static + dynamic signals alternate, driving the LLM toward coverage convergence.

Together these give MAgHARCM's `internal/agents/validator.go` a **targeted test generation** substrate that complements its existing compile+run loop.

## Anchoring (Q3)

| Primitive | Pre-wave-19 behaviour | Panta substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Comprehension phase identified untested code paths only via static analysis or only via coverage feedback — never both in a loop | Static-flow ∪ dynamic-coverage signal drives the comprehension phase's search for untested paths |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Validator picked samples by order; no coverage-driven prioritisation | Coverage-driven knob lets the validator prioritise samples whose untested branches are most informative |

## Hop-1 Citations

- P-113 AutoCodeRover (AST-aware structural retrieval + in-loop verification).
- P-115 HyperAgent (multi-agent SE agent baseline).
- AlphaTrans (test-driven translation baseline).

## Hop-2 Citations

- McCabe 1976 Cyclomatic Complexity (control-flow complexity metric).
- Floyd 1967 Assigning Meanings to Programs (verification foundation).
- P-48 Jia-Harman 2011 Mutation Testing (test quality metric).

## MAgHARCM integration

- **YAML config key**: `test_gen.panta.enabled: true`; `test_gen.panta.coverage_target: <float in [0,1]>`.
- **Implementation file**: `internal/testgen/panta.go::HybridLoop` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-21]]`.

## Caveats

- **Instrumentation cost** — dynamic coverage feedback requires instrumented runs; cost grows with codebase size.
- **Static ranking** — cyclomatic complexity is a proxy; high-complexity ≠ always high-value branches.
- **LLM context budget** — feeding uncovered-branch context to the LLM consumes tokens; budget for long iterations.

## Source

- arXiv: 2503.13580.
- Venue: ICSE 2026 Research Track (verified via conf.researchr.org/details/icse-2026/icse-2026-research-track/36/; Wed 15 Apr 2026 14:45 talk, Asia I, "AI for Software Engineering 4" session).
