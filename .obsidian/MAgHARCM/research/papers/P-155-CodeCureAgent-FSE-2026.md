---
title: "P-155 CodeCureAgent: Automatic Classification and Repair of Static Analysis Warnings"
backlink: "[[1.0.0 P-155]]"
aliases:
 - "1.0.0 P-155"
 - "P-155"
 - "P-155-CodeCureAgent-FSE-2026"
 - "CodeCureAgent-FSE-2026"
 - "codecureagent2026fse"
tags: [paper, static-analysis, warning-classification, automated-repair, agentic-llm, validation-heuristic, "[[1.0.0 PRIM-13]]", "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-29]]", "[[1.0.0 P-42]]", "[[1.0.0 P-124]]", "[[1.0.0 P-142]]", "[[1.0.0 P-151]]", wave-23]
date: 2026-09-07
last_updated: 2026-09-07 (iter-7, wave-23)
venue: FSE 2026 (Research Track, PACMSE Vol. 3 Issue FSE)
---

# [[1.0.0 P-155]] CodeCureAgent

## TL;DR

CodeCureAgent is an **agentic LLM framework** for **static-analysis warning classification + repair**. Iteratively invokes code-search + multi-file-edit tools; classifies false vs true positives before repair; employs a **three-step validation heuristic** (project builds, target warning removed without new warnings, test suite green). 96.8% plausible-fix / 86.3% correct-fix on 1,000 SonarQube warnings across 106 Java projects, ~$0.029 + ~4 min per warning. Anchors `[[1.0.0 PRIM-13]]` Adversarial Test-Weakening Guard (agentic warning-classification + repair twin-objective extends static-analysis validation lineage) and `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Observation phase) and `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter (three-step build+test validation = concrete acceptance heuristic).

## Mechanism (Q2)

1. **Agentic LLM framework** — instead of a fixed algorithmic pipeline, CodeCureAgent uses an LLM-driven agent that iteratively invokes tools: code-search (gather context from the codebase), file-reader (read multi-file context), edit (perform multi-file edits). The agent decides which tool to invoke and when.
2. **Classification + repair twin objective** — the agent classifies warnings before repair: false positives are suppressed, true positives are passed to the repair step. This avoids wasting repair effort on warnings that are not real bugs.
3. **Three-step validation heuristic** — every patch must pass (a) project builds clean, (b) target warning is removed without new warnings, (c) test suite passes. Patches that fail any step are rejected and the agent retries with the failure signal as context.
4. **Multi-file edit support** — unlike single-file LLM repair, CodeCureAgent can perform multi-file edits when the fix requires it (e.g., introducing a helper method in a separate file).

Together these give MAgHARCM an **agentic static-analysis warning-classification + repair framework** — the first concrete acceptance heuristic (build + warning-removed + test-green) for agent-loop validation in the static-analysis substrate.

## Anchoring (Q3)

| Primitive | Pre-wave-23 behaviour | CodeCureAgent substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-13]]` Adversarial Test-Weakening Guard | Static-analysis validation lineage anchored on P-42 (basic static-analysis substrate), P-48 Jia & Harman Mutation Survey, P-124 Syzygy (compilation feedback), P-142 NESA (Datalog policy); no agentic LLM warning-classification + repair framework | CodeCureAgent as the **agentic LLM warning-classification + repair twin-objective framework**. Closes the agent-loop static-analysis gap; extends the static-analysis validation lineage to LLM-driven classification + repair. |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Comprehension Observation phase anchored on P-128 KVzip, P-133 ADI, P-142 NESA; no agentic warning-classification substrate | CodeCureAgent as the agentic warning-classification substrate for the Observation phase |
| `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter | Iteration recruiter anchored on P-151 SmartC2Rust (Rust compiler feedback loop), P-124 Syzygy (compilation feedback), abstract pattern only; no concrete three-step acceptance heuristic | CodeCureAgent three-step build+test validation = **concrete acceptance heuristic** for the iteration recruiter |

## Hop-1 Citations

- P-124 Syzygy (ICSE 2025, Wave-17) — iterative compilation-feedback loop. CodeCureAgent extends the iterative-compilation-feedback substrate with the three-step acceptance heuristic.
- P-142 NESA (FSE 2026, Wave-20) — Datalog policy for static analysis. CodeCureAgent shares the static-analysis substrate but specialises for warning-classification + repair.
- P-151 SmartC2Rust (ICSE 2026, Wave-22) — iterative compiler + semantic + unsafe feedback loop. CodeCureAgent inherits the iterative-feedback-loop pattern.
- LLM4Code warning-repair lineage — broader LLM-based warning-repair substrate. CodeCureAgent is the agentic LLM instance.

## Hop-2 Citations

- SonarQube static-analysis warning taxonomy — foundational static-analysis warning substrate. CodeCureAgent is evaluated on SonarQube warnings.
- Build-and-test patch validation lineage — patch-validation substrate. CodeCureAgent's three-step validation extends the build-and-test lineage.
- Agentic framework tooling (ReAct, Toolformer) — agentic framework substrate. CodeCureAgent uses an agentic LLM with tool invocations.
- Multi-file code repair lineage — multi-file repair substrate. CodeCureAgent supports multi-file edits.

## MAgHARCM integration

- **YAML config key**: `agents.static_analysis.classification: true`; `agents.static_analysis.repair: true`; `agents.static_analysis.validation_heuristic: ["build", "warning_removed", "tests_pass"]`; `agents.static_analysis.max_iterations: <int>`.
- **Implementation file**: `internal/staticanalysis/codecureagent.go::NewCodeCureAgent` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-13]]`, `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-29]]`.

## Caveats

- **Java + SonarQube only** — evaluated on Java projects with SonarQube warnings; generalisation to other languages / static analysers is future work.
- **Cost per warning** — ~$0.029 + ~4 min per warning; the cost is small but non-trivial at scale (1,000 warnings = ~$29 + ~67 hours).
- **Test-suite dependency** — the three-step validation requires a runnable test suite; projects without tests will fail the validation heuristic.
- **Build-system dependency** — the build step requires a configured build system; non-buildable projects are out of scope.

## Source

- Venue: FSE 2026 (Research Track), verified via dl.acm.org/doi/10.1145/3808139; PACMSE Vol. 3, Issue FSE, DOI: 10.1145/3808139. Held July 5-9 2026 in Montréal, Canada.
- arXiv: 2509.11787 (preprint).

## BibTeX

```
@inproceedings{codecureagent2026fse,
  title     = {CodeCureAgent: Automatic Classification and Repair of Static Analysis Warnings},
  author    = {Joos, Pascal and Bouzenia, Islem and Pradel, Michael},
  booktitle = {Proceedings of the ACM International Conference on the Foundations of Software Engineering (FSE)},
  year      = {2026},
  volume    = {3},
  number    = {FSE},
  series    = {Proceedings of the ACM on Software Engineering (PACMSE)},
  doi       = {10.1145/3808139}
}
```
