---
title: MAgHARCM Architecture & System Dataflow
date: 2026-09-08
last_updated: 2026-09-08 (iter-1, wave-23)
aliases:
  - "Architecture-And-Dataflow"
  - "Architecture and Dataflow"
  - "Dataflow"
tags: [adhoc, architecture, dataflow, dag, agents, eino, "[[2.0.0 MAgHARCM]]"]
---

# [[2.0.0 MAgHARCM Architecture & System Dataflow]]

> **Executive Overview**: Visual and architectural specification of the MAgHARCM multi-agent graph, typed pipeline state, and Locality of Behaviour boundaries. Wave-23 (2026-09-08 iter-1) closes the strict program-comprehension-mechanism residual slot via P-153 CoReX; anchors 4 new SLM-era papers (P-153..P-156) onto PRIM-12, PRIM-22, PRIM-23, PRIM-29, PRIM-31. Vault paper count: 152 -> 156.

---

## 1. High-Level Agent Execution Graph

The system executes software modernization as a collaborative cycle orchestrated via CloudWeGo Eino:

```
[Start / Config]
       │
       ▼
1. Archaeologist Agent (PRIM-14, 18, 19, 20, 22)
   └─► Excavates legacy churn, git co-change coupling, DRSpaces
       │
       ▼
2. Analyzer Agent (PRIM-4, 10, 15, 21)
   └─► Profiles source, maps third-party libraries to target crates
       │
       ▼
3. Planner Agent (PRIM-1, 2, 3, 16, 30)
   └─► Computes reverse-topological order, generates target skeletons
       │
       ▼
4. Translator Agent (PRIM-23, 26, 31) ◄─────────────────────────┐
   └─► Performs token-bounded AST chunked translation           │ (Iterative
       │                                                        │  Repair Loop)
       ▼                                                        │
5. Role-Flip Reviewer (PRIM-25)                                 │
   └─► Adversarial anti-sycophancy sanity check                 │
       │                                                        │
       ▼                                                        │
6. Validator Agent (PRIM-5, 6, 13, 27)                          │
   └─► 3-tier cascade: AST Check ──► rustc Compiler ──► Tests   │
       │                                                        │
       ├─► [Validation Failed?] ────────────────────────────────┘
       │
       ▼ [Validation Passed]
7. Verdict Panel (PRIM-7, 8, 11, 12)
   └─► Multi-judge consensus agreement
       │
       ▼
8. Recruiter Agent (PRIM-29) ──► [Final Translated Repository]
```

---

## 2. Shared Pipeline State Container (`compiletime.State`)

Communication across agent boundaries is governed by the cohesive typed pipeline state (`internal/compiletime/state.go`):

- **`ArchaeologyReport`**: System dependence graph, logical coupling matrix, architectural boundaries.
- **`AnalyzerOutput`**: Crate mapping dictionary, target design specification.
- **`PlanningOutput`**: Reverse-topological file synthesis sequence, interface skeleton files.
- **`TranslatedProject`**: Map of translated target source files and co-translated unit test suites.
- **`ValidationReport`**: Binary compilation status, test pass/fail breakdown, compiler error diagnostics.
- **`VerdictReport`**: Multi-agent consensus score, equivalence agreement.
- **`RecruitmentPlan`**: Adaptive retry tooling allocations and focus targets.

---

## 3. Locality of Behaviour (ADR-C-014)

- **Principle**: State declarations are centralized in `compiletime/state.go` to eliminate Go circular import dependencies (`compiletime` does not import `agents`).
- **Producer Ownership**: Specialized artifact methods and mutation logic reside locally within the respective agent package (`internal/agents/`).

---

## 4. Wave-23 Substrate Integration (2026-09-08 iter-1)

Wave-23 anchors four new SLM-era papers onto existing primitives, completing a Wave-17-first-opened slot:

- **`[[1.0.0 P-153]]` CoReX** → extends `PRIM-22` Four Phases of Comprehension (Structure phase, function-level DA pattern extended from `[[1.0.0 P-133]]` ADI) + `PRIM-31` Iterative Retrieval Refinement (refinement-based slicing = context-conditioned summary pass = partition-aligned summary-pass substrate). Closes the strict program-comprehension-mechanism residual slot.
- **`[[1.0.0 P-154]]` TransAgent** → extends `PRIM-23` Chunked Translation (multi-agent translator + execution-aligned critic, complementing `[[1.0.0 P-151]]` SmartC2Rust single-LLM iterative feedback) + `PRIM-31` Iterative Retrieval Refinement (execution-aligned critic feedback = retrieval refinement substrate).
- **`[[1.0.0 P-155]]` POLA-Tester** → extends `PRIM-12` Static Analysis Co-Evolution (agentic wait + syntactic dependency mining + iterative retrofit validation; LLM-augmented static-analysis pattern).
- **`[[1.0.0 P-156]]` ACONITE** → extends `PRIM-22` Four Phases of Comprehension (Structure phase — backward slicing + execution annotations = LLM-augmented static-analysis pattern complementing `[[1.0.0 P-129]]` LλMDA's partial-PDG pattern) + `PRIM-29` Dynamic Iteration Recruiter (close-test retrieval + execution annotation = retrieval-and-recruitment substrate for test-generation agents).

The agent graph topology (`internal/graph/graph.go`) is unchanged; the Wave-23 anchors are substrate refinements wired through existing opt-in config paths (`configs/agents.yml:comprehension.graph_self_evolving` and `configs/agents.yml:translation.feedback_driven`). No new agent nodes; no new graph edges; no new state types.
