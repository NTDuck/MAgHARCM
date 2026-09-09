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

> **Executive Overview**: Visual and architectural specification of the MAgHARCM multi-agent graph, typed pipeline state, and Locality of Behaviour boundaries. Domain acronyms (AST, CPG, PDG, DRSpaces, SLM) are formally defined in the [[Glossary|Domain Acronyms & Terminology Glossary]].

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

Communication across agent boundaries is strictly governed by the cohesive typed pipeline state (`internal/compiletime/state.go`):

| State Artifact | Producing Agent | Description & Consumed Data |
| :--- | :--- | :--- |
| **`ArchaeologyReport`** | `archaeologist` | System dependence graph, logical coupling matrix, DRSpaces architectural boundaries. |
| **`AnalyzerOutput`** | `analyzer` | Third-party crate mapping dictionary, target design specification, library replacements. |
| **`PlanningOutput`** | `planner` | Reverse-topological file synthesis sequence, interface skeleton files. |
| **`TranslatedProject`** | `translator` | Map of translated target source files and co-translated unit test suites. |
| **`ValidationReport`** | `validator` | Strict binary compilation status (`Pass`/`Fail`), test pass/fail breakdown, compiler diagnostics. |
| **`VerdictReport`** | `verdict_panel` | Multi-agent consensus score, behavioral equivalence agreement. |
| **`RecruitmentPlan`** | `recruiter` | Adaptive retry tooling allocations, strategy escalations, and focus targets. |

---

## 3. Locality of Behaviour (ADR-C-014)

- **Principle**: State declarations are centralized in `compiletime/state.go` to eliminate Go circular import dependencies (`compiletime` does not import `agents`).
- **Producer Ownership**: Specialized artifact methods and mutation logic reside locally within the respective agent package (`internal/agents/`).
- **Inviolable Rule**: Agents consume upstream state via typed read-only references and produce concrete output structs; no agent may mutate private state inside another agent's package.

---

## 4. Wave-23 Substrate Integration (2026-09-08 iter-1)

### Context & Motivation
Earlier research waves identified a critical capability gap: small language models (4B–30B) often fail when attempting to translate full-file contexts or arbitrary AST chunks because unrelated code introduces distracting noise. 

| Anchor Paper | Target Primitives | Benefiting Agent | Architectural Mechanism | Specified Target Config Gate |
| :--- | :--- | :--- | :--- | :--- |
| **`[[1.0.0 P-153]]` CoReX**<br>(Sun et al., ICSE 2026) | `[[1.0.0 PRIM-22]]`<br>`[[1.0.0 PRIM-31]]` | `archaeologist`<br>`translator` | **Context-Aware Refinement Slicing**: Replaces naive AST chunking with refinement-based dynamic slicing to isolate regression failures. | `comprehension.graph_self_evolving: true` (target) |
| **`[[1.0.0 P-154]]` TransAgent**<br>(Roh et al., FSE 2026) | `[[1.0.0 PRIM-23]]`<br>`[[1.0.0 PRIM-31]]` | `translator` | **Multi-Agent Critic Feedback**: Pairs a generation model with an execution-aligned critic to guide iterative AST chunk translation. | `translation.feedback_driven: true` (target) |
| **`[[1.0.0 P-155]]` CodeCureAgent**<br>(Joos et al., FSE 2026) | `[[1.0.0 PRIM-13]]`<br>`[[1.0.0 PRIM-22]]`<br>`[[1.0.0 PRIM-29]]` | `validator`<br>`verdict_panel` | **Agentic Warning Classification + Repair**: Classifies false vs true positives, then patches true positives through multi-file tool invocations with a three-step build+test acceptance heuristic. | `comprehension.graph_self_evolving: true` (target) |
| **`[[1.0.0 P-156]]` TestWeaver**<br>(Le et al., ICSE 2026) | `[[1.0.0 PRIM-22]]`<br>`[[1.0.0 PRIM-29]]`<br>`[[1.0.0 PRIM-31]]` | `recruiter`<br>`validator` | **Execution-Annotated Backward Slicing**: Retrieves close tests sharing control-flow similarity with the target path and injects variable-state annotations into prompt contexts. | `translation.feedback_driven: true` (target) |

### Architectural Invariants & Stability
1. **Zero Graph Topology Mutations**: The 8-agent cyclic execution graph (`internal/graph/graph.go`) is structurally unaltered. No new nodes or edge transitions were added.
2. **Opt-in Substrate Gating**: All Wave-23 substrates are specified to operate behind opt-in configuration gates (e.g. `comprehension.graph_self_evolving`, `translation.feedback_driven`). As tracked in the Wave-23 handoff (`Sprint-2026-09-08-Handoff-7.md`), physical wiring into `configs/agents.yml` is an open initiative for Wave-24; the baseline pipeline runs unaffected.
3. **State Schema Stability**: The shared pipeline container `compiletime.State` remains stable and backwards-compatible.
