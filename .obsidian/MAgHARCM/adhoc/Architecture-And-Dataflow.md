---
title: MAgHARCM Architecture & System Dataflow
date: 2026-09-28
last_updated: 2026-09-28
aliases:
  - "Architecture-And-Dataflow"
  - "Architecture and Dataflow"
  - "Dataflow"
tags: [adhoc, architecture, dataflow, dag, agents, eino, "[[2.0.0 MAgHARCM]]"]
---

# [[2.0.0 MAgHARCM Architecture & System Dataflow]]

> **Executive Overview**: Visual and architectural specification of the MAgHARCM multi-agent graph, typed pipeline state, and Locality of Behaviour boundaries.

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
