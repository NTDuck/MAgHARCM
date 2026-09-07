---
title: On the Criteria To Be Used in Decomposing Systems into Modules
backlink: "[[1.0.0 P-31]]"
bibkey: p31_parnas
aliases:
  - "1.0.0 P-31"
  - "P-31"
  - "P-31-Parnas-1972"
  - "P-31-Parnas-1972"
  - "Parnas-1972"
tags: [paper, modularity, information-hiding, [[1.0.0 PRIM-3]], [[1.0.0 PRIM-19]], hop-1]
---

# [[1.0.0 P-31]] On the Criteria To Be Used in Decomposing Systems into Modules

**Authors**: David L. Parnas  
**Year**: 1972  
**Venue**: Communications of the ACM (CACM), Vol. 15, No. 12, pp. 1053–1058  
**eprint / DOI**: 10.1145/361598.361623  
**Cited by**: `[[2.0.0 Software-Archaeology-Lineage]]`, `[[1.0.0 PRIM-3]]` (Target Skeleton-First Generation), `[[1.0.0 PRIM-19]]` (Design Rule Hierarchy Partitioning)

## Summary

Parnas's seminal 1972 paper introduces the concept of **Information Hiding** as the foundational principle for modular software decomposition. Parnas compares two decomposition strategies for a KWIC (Key Word in Context) index system:
1. *Decomposition based on flowchart / processing steps*: modules represent chronological phases of execution (input, circular shift, alphabetize, output).
2. *Decomposition based on information hiding*: each module hides a design decision or "secret" (data representation, sorting algorithm, file storage format) behind an abstract, stable interface.

Parnas demonstrates that decomposition by information hiding yields superior maintainability, independent development viability, and comprehensive comprehensibility because changes to implementation secrets do not ripple across module boundaries.

## Relevance to MAgHARCM

Parnas (1972) is the direct theoretical justification for:
1. `[[1.0.0 PRIM-3]]` (Target Skeleton-First Generation): Generating target language trait and interface declarations before any method bodies forces the LLM to commit to stable module boundaries, enclosing implementation details behind Parnas contracts.
2. `[[1.0.0 PRIM-19]]` (Design Rule Hierarchy Partitioning): L1 design rule interfaces represent Parnas boundaries that must remain invariant when volatile L3 leaf implementations are refactored.

## Hop-1 References

- [[Baldwin-Clark-2000]] — Formalizes Parnas's information hiding into economic modularity theory and the six modular operators.
- [[Chikofsky-Cross-1990]] — Relies on Parnas modularity to define reverse engineering boundaries and design recovery.
- [[Rajlich-1997]] — Uses information hiding modules as the physical boundaries for concept locator analysis.

## Hop-2 Anchors (Software Archaeology Lean)

- [[Feathers-2004-WELC]] — Parnas modules define the "seams" where characterization tests can be injected without editing client code.
- [[Kazman-Cai-2017]] — DRSpaces architecture debt patterns measure violations of Parnas information hiding where secrets leak across interfaces.
- [[pp-besm-Software-Archaeology]] — Legacy excavation begins by reconstructing Parnas's original module secrets through header inspection.
