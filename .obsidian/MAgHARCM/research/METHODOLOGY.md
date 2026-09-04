---
title: MAgHARCM Methodology
backlink: [[2.0.0 Methodology]]
tags: [methodology, architecture, pipeline, [[2.0.0 MAgHARCM]]]
---

# [[2.0.0 MAgHARCM Methodology]]

## 1. Multi-Agent Pipeline Overview

MAgHARCM structures software modernization as a collaborative multi-agent lifecycle wired via an Eino execution graph.
Agents communicate across explicit boundaries through the shared typed pipeline state.
Every primitive from `[[1.0.0 PRIM-1]]` through `[[1.0.0 PRIM-31]]` represents a concrete capability implemented in the codebase.

```
START ──► archaeologist ──► analyzer ──► planner ──► translator ──► reviewer ──► validator ──► verdict_panel ──► branch
                                            ▲                                        │                              │
                                            │                                        ▼                              ▼
                                            └──────────────── repair_loop ───────────┴────────── recruiter ─────────┘
```

---

## 2. Agent Roles and Boundaries

| Stage | Node Name | Agent Role | Primitive Grounding | Output Contract |
| :--- | :--- | :--- | :--- | :--- |
| 1 | `archaeologist` | Software Archaeologist | `[[1.0.0 PRIM-14]]`, `[[1.0.0 PRIM-18]]`, `[[1.0.0 PRIM-19]]`, `[[1.0.0 PRIM-20]]`, `[[1.0.0 PRIM-22]]` | `ArchaeologyReport`: boundaries, churn coupling, DRSpaces, concept maps |
| 2 | `analyzer` | System & Library Analyzer | `[[1.0.0 PRIM-4]]`, `[[1.0.0 PRIM-10]]`, `[[1.0.0 PRIM-15]]`, `[[1.0.0 PRIM-21]]` | `AnalyzerOutput`: project research, third-party library mapping, target design |
| 3 | `planner` | Topological Planner | `[[1.0.0 PRIM-1]]`, `[[1.0.0 PRIM-2]]`, `[[1.0.0 PRIM-3]]`, `[[1.0.0 PRIM-16]]`, `[[1.0.0 PRIM-30]]` | `PlanningOutput`: reverse-topo implementation plan, skeleton files |
| 4 | `translator` | Code & Test Synthesizer | `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-26]]`, `[[1.0.0 PRIM-31]]` | `TranslatedProject`: translated source and unit test files |
| 5 | `reviewer` | Role-Flip De-Hallucination Gate | `[[1.0.0 PRIM-25]]` | `ReviewerReport`: adversarial sanity check, sycophancy rejection |
| 6 | `validator` | Build & Test Cascade | `[[1.0.0 PRIM-5]]`, `[[1.0.0 PRIM-6]]`, `[[1.0.0 PRIM-13]]`, `[[1.0.0 PRIM-27]]` | `ValidationReport`: binary pass/fail compilation, test pass rate, error codes |
| 7 | `verdict_panel` | Multi-Agent Consensus | `[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-8]]`, `[[1.0.0 PRIM-11]]`, `[[1.0.0 PRIM-12]]` | `VerdictReport`: multi-judge equivalence agreement |
| 8 | `recruiter` | Dynamic Iteration Recruiter | `[[1.0.0 PRIM-29]]` | `RecruitmentPlan`: targeted repair tools and focus areas |

---

## 3. Incremental Try-and-Fail Strategy Registry (`[[1.0.0 PRIM-21]]`)

Rather than committing irrevocably to a static heuristic choice, the analyzer and repair loop evaluate migration strategies dynamically.
Each strategy implements `Matches(Profile) bool` and `Attempt(context.Context, Profile) error`.
If a strategy encounters validation plateau or structural repair failure, the system falls back incrementally to the next viable strategy:

1. `BIG_BANG`: Evaluated first for small codebases ($\le 3$ files, $< 500$ LoC).
2. `PILOT`: Evaluated for large systems ($> 50$ files or $> 10000$ LoC), translating an isolated subsystem first.
3. `PARALLEL_CUTOVER`: Modular systems with comprehensive tests ($> 10$ files).
4. `FROZEN_LEGACY`: Systems with zero existing tests, requiring characterization synthesis.
5. `INCREMENTAL`: Canonical universal baseline strategy.

---

## 4. Software Archaeology Suite

Pre-translation comprehension executes through five integrated archaeological primitives:
1. `[[1.0.0 PRIM-14]]` **Archaeology Stage**: Discovers module boundaries, historical build time capsules, and git churn hotspots.
2. `[[1.0.0 PRIM-18]]` **Jaccard-Coupling Recovery**: Computes temporal co-change similarity across commits to expose hidden coupling.
3. `[[1.0.0 PRIM-19]]` **Design Rule Hierarchy**: Partitions the codebase into L1 interfaces, L2 subsystems, and L3 leaves per Baldwin & Clark.
4. `[[1.0.0 PRIM-20]]` **Concept Assignment**: Locates domain concepts across lexical clusters per Rajlich.
5. `[[1.0.0 PRIM-22]]` **Four Phases of Comprehension**: Applies Foltz's DR. JONES cognitive traversal model.

---

## 5. Binary Compilation & Validation Cascade

Per-project compilation status is strictly binary: **Pass** or **Fail**. There is no partial compilation percentage.
The validation cascade enforces:
1. **AST Syntax Pre-check**: Fast parsing via configured LSP provider (`abcoder-mcp` by default).
2. **Native Toolchain Compilation**: Strict type-checking and borrow-checker inspection (`cargo check`).
3. **Automated Test Suite**: Execution of translated and synthesized tests (`cargo test`).
4. **Adversarial Weakening Guard (`[[1.0.0 PRIM-13]]`)**: Halts if test assertions are removed or widened.
5. **Coverage-Guided Plateau Detector (`[[1.0.0 PRIM-27]]`)**: Exits loop when test improvements stagnate.

---

## 6. Centralized Compile-time Config & Locality of Behaviour

Compile-time invariants, enums, sentinels, and initialization helpers reside in `internal/compiletime`.
Every agent declares its intermediate artifacts within its own module file, upholding Locality of Behaviour.
The pipeline state coordinates data flow across agents through explicit typed contracts.
