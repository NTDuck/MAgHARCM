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

Wave-23 anchors four recent flagship SE papers (ICSE 2026, FSE 2026) that provide **fine-grained program slicing, static analysis co-evolution, and critic-feedback loops**. Crucially, this closes the long-standing *strict program-comprehension-mechanism residual slot* carried forward from Wave-17.

### Substrate Integration Matrix

| Anchor Paper | Target Primitives | Benefiting Agent | Architectural Mechanism | Opt-in Config Wire |
| :--- | :--- | :--- | :--- | :--- |
| **`[[1.0.0 P-153]]` CoReX**<br>(Sun et al., ICSE 2026) | `[[1.0.0 PRIM-22]]`<br>`[[1.0.0 PRIM-31]]` | `archaeologist`<br>`translator` | **Context-Aware Refinement Slicing**: Replaces naive AST chunking with refinement-based dynamic slicing to isolate regression failures. | `comprehension.graph_self_evolving: true` |
| **`[[1.0.0 P-154]]` TransAgent**<br>(Roh et al., FSE 2026) | `[[1.0.0 PRIM-23]]`<br>`[[1.0.0 PRIM-31]]` | `translator` | **Multi-Agent Critic Feedback**: Pairs a generation model with an execution-aligned critic to guide iterative AST chunk translation. | `translation.feedback_driven: true` |
| **`[[1.0.0 P-155]]` POLA-Tester**<br>(Sun et al., ICSE 2026) | `[[1.0.0 PRIM-12]]` | `validator`<br>`verdict_panel` | **Agentic Static Analysis**: Syntactic dependency mining and iterative retrofit validation to detect subtle semantic discrepancies. | `comprehension.graph_self_evolving: true` |
| **`[[1.0.0 P-156]]` ACONITE**<br>(Sun et al., ICSE 2026) | `[[1.0.0 PRIM-22]]`<br>`[[1.0.0 PRIM-29]]` | `recruiter`<br>`validator` | **Execution-Annotated Backward Slicing**: Retrieves relevant tests and injects runtime execution traces into prompt contexts. | `translation.feedback_driven: true` |

### Architectural Invariants & Stability
1. **Zero Graph Topology Mutations**: The 8-agent cyclic execution graph (`internal/graph/graph.go`) is structurally unaltered. No new nodes or edge transitions were added.
2. **Backward-Compatible Configuration**: All Wave-23 substrates operate behind opt-in configuration flags in `configs/agents.yml`. When disabled, the pipeline executes its baseline algorithms.
3. **State Schema Stability**: The shared pipeline container `compiletime.State` remains identical; enhanced data flows through existing extensible payload fields.
