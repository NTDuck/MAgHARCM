---
title: MAgHARCM Domain Acronyms & Terminology Glossary
date: 2026-09-08
last_updated: 2026-09-08
aliases:
  - "Glossary"
  - "Domain Glossary"
  - "Terminology"
tags: [adhoc, glossary, definitions, terminology, "[[2.0.0 MAgHARCM]]"]
---

# [[2.0.0 MAgHARCM Domain Acronyms & Terminology Glossary]]

> **Executive Overview**: Reference dictionary defining domain-specific terminology, acronyms, and architectural abstractions used across the MAgHARCM modernization reports, research waves, and codebase specifications.

---

## 1. Core Machine Learning & SLM Terminology

| Term / Acronym | Full Form | Definition & MAgHARCM Context |
| :--- | :--- | :--- |
| **SLM** | Small Language Model | Parameter-efficient models (typically 4B–30B parameters, e.g., Qwen2.5-Coder:7B/32B, Phi-3, StarCoder2) deployed locally on workstation GPUs (RTX 4090 / Apple Silicon) to eliminate inference costs and enterprise IP leakage. |
| **KV-Cache** | Key-Value Cache | Memory buffers storing attention keys and values for previously generated tokens in Transformers. In long-running multi-agent pipelines, uncompressed KV-caches exhaust GPU VRAM; MAgHARCM incorporates KV eviction and compression (`SpecKV`, `ChunkKV`, `KVFlow`). |
| **MaTTS** | Memory-augmented Test-Time Scaling | Algorithmic loop (`[[1.0.0 P-122]]` ReasoningBank) coupling persistent memory triples `(Title, Description, Content)` with allocated test-time compute budgets to iteratively improve migration candidates without human guidance. |
| **PRM** | Process Reward Model | An evaluation model trained to score intermediate execution or reasoning steps rather than assigning a binary score only to final code (`ContextPRM` `[[1.0.0 P-146]]`). |
| **Speculative Decoding** | Speculative Decoding | Acceleration technique using an efficient draft model or suffix tree (`SuffixDecoding` `[[1.0.0 P-137]]`) to propose candidate tokens verified in parallel by the target model. |

---

## 2. Program Analysis & Software Engineering Terminology

| Term / Acronym | Full Form | Definition & MAgHARCM Context |
| :--- | :--- | :--- |
| **ACI** | Agent-Computer Interface | The structured interaction layer (tools, file I/O, compilation feedback, diagnostics) through which autonomous agents inspect and mutate repository codebases (`SWE-agent` `[[1.0.0 P-112]]`). |
| **AST** | Abstract Syntax Tree | Hierarchical tree representation of program syntax used by compilers and code analyzers for syntactic verification and chunked translation (`[[1.0.0 PRIM-23]]`). |
| **CFG** | Control Flow Graph | Directed graph representation of all paths that might be traversed through a program during its execution. |
| **PDG** | Program Dependence Graph | Graph representation making data and control dependencies between program statements explicit (`LλMDA` `[[1.0.0 P-129]]`). |
| **CPG** | Code Property Graph | Unified graph data structure merging AST, CFG, and PDG into a queryable property graph for whole-repository static reasoning (`GraphCoder` `[[1.0.0 P-110]]`). |
| **SDG** | System Dependence Graph | Inter-procedural extension of the PDG representing multiple functions, callsites, and global data flows across an entire repository (`TypePro` `[[1.0.0 P-139]]`). |
| **DRSpaces** | Design Rule Spaces | Architectural analysis method (Baldwin & Clark; Xiao et al.) decomposing software into design-rule hierarchies and identifying coupling and churn hotspots (`[[1.0.0 PRIM-14]]`, `[[1.0.0 PRIM-18]]`). |
| **DA** | Dynamic Analysis | Software analysis performed by executing code with concrete inputs to observe runtime invariants, aliasing, and memory safety (`Syzygy` `[[1.0.0 P-124]]`). |
| **SA** | Static Analysis | Code analysis conducted without execution, analyzing ASTs, types, and control/data flows (`POLA-Tester` `[[1.0.0 P-155]]`). |

---

## 3. MAgHARCM System Primitives & Concepts

| Concept | Description |
| :--- | :--- |
| **Reverse Topological Planning** | Planning translation order starting from leaves (zero-dependency utility files) up to entrypoints (`[[1.0.0 PRIM-1]]`), preventing cyclic compilation deadlocks. |
| **Target Skeleton-First** | Emitting empty structs, traits, and interface signatures in the target language before generating function bodies (`[[1.0.0 PRIM-3]]`) so dependent files compile cleanly. |
| **Binary Compilation Status** | Strict invariant: per-project compilation is exclusively **Pass** or **Fail** (no partial compilation percentages permitted). |
| **Locality of Behaviour (LoB)** | ADR-C-014 architectural rule: artifact struct definitions reside in the same Go file as the producer function, centralizing shared state in `internal/compiletime/state.go`. |
| **Try-and-Fail Strategy Registry** | Dynamic fallback sequence (`[[1.0.0 PRIM-21]]`) that sequentially tries alternative translation strategies when compiler or test validation fails. |
