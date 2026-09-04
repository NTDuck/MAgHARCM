---
title: Software Archaeology & Legacy System Modernization Lineage
tags: [research, software-archaeology, modernization, lineage, synthesis]
---

# Software Archaeology & Legacy System Modernization Lineage

## Overview

Software modernization via Large Language Models (LLMs) and Small Language Models (SLMs) does not emerge from a vacuum. Modern multi-agent pipelines like MAgHARCM, ReCodeAgent, AlphaTrans, and CAID re-discover, formalize, and automate fundamental software modernization principles established across 30+ years of software engineering research.

## Core Intellectual Lineages

### 1. Program Comprehension & Concept Assignment
- **Rajlich (1997) & Müller (2000)**: Defined concept assignment and redocumentation. Before code can be safely migrated, legacy code concepts must be mapped to target domain concepts. In MAgHARCM, this maps directly to [[PRIM-14]] (Software-Archaeology Stage) and [[PRIM-20]].
- **Foltz et al. (1998, 2023)**: Explored mental models and cognitive traversal of source text using semantic spaces (LSA/embeddings). In LLM pipelines, this manifests as bounded context retrieval ([[PRIM-26]], [[PRIM-31]]).

### 2. Modularity Theory & Architectural Decoupling
- **Baldwin & Clark (2000) *Design Rules***: Established modularity theory and the economic/structural value of decoupling. Breaking dependencies via interfaces is the theoretical basis for skeleton-first code generation ([[PRIM-3]]).
- **Kazman et al. (2000, 2017) DRSpaces**: Formalized architectural degradation analysis and design structure matrices (DSM). In MAgHARCM, DAG linearization and cycle elimination ([[PRIM-1]], [[PRIM-2]]) ensure that translation leaves are translated prior to dependent callers, preserving architectural boundaries.

### 3. Verification & Dynamic Recovery
- **Syzygy (2024)** / SpecMiner lineage: Runtime tracing to discover hidden invariants (nullability, allocation ranges) so target language types (e.g. Rust Option/Result) can accurately reflect runtime reality rather than raw pointer syntax ([[PRIM-4]]).
- **MSR & Legacy Mining (pp-besm, dev.to)**: Treats code repositories as geological strata. Git churn hotspots, author co-commit metrics, and temporal coupling indicate structural risk areas for translation agents.
