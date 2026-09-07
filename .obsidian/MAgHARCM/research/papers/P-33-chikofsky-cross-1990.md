---
title: Reverse Engineering and Design Recovery: A Taxonomy
bibkey: p33_chikofsky
tags: [paper, reverse-engineering, design-recovery, redocumentation, [[1.0.0 PRIM-14]], [[1.0.0 PRIM-20]], hop-1]
---

# [[1.0.0 P-33]] Reverse Engineering and Design Recovery: A Taxonomy

**Authors**: Elliot J. Chikofsky, James H. Cross II  
**Year**: 1990  
**Venue**: IEEE Software, Vol. 7, No. 1, pp. 13–17  
**eprint / DOI**: 10.1109/52.43044  
**Cited by**: `[[2.0.0 Software-Archaeology-Lineage]]`, `[[1.0.0 PRIM-14]]` (Software Archaeology Stage), `[[1.0.0 PRIM-20]]` (Concept Assignment and Redocumentation)

## Summary

Chikofsky and Cross formulated the standard taxonomy distinguishing the six fundamental processes in software maintenance and evolution:
1. **Forward Engineering**: The traditional development movement from high-level abstractions and logical designs to physical implementation.
2. **Reverse Engineering**: Analyzing a subject system to identify components and interrelationships, and to create representations of the system at a higher level of abstraction. Reverse engineering does *not* alter the system.
3. **Redocumentation**: Creating or revising semantically equivalent representations at the same level of abstraction (e.g. reformatted call trees, data-flow diagrams).
4. **Design Recovery**: Recreating design abstractions from code, domain knowledge, and external observation.
5. **Restructuring**: Transforming code from one representation form to another at the same abstraction level while preserving external behavior (refactoring).
6. **Reengineering (Modernization)**: The examination and alteration of a legacy system to reconstitute it into a new target architecture (Reverse Engineering + Restructuring + Forward Engineering).

## Relevance to MAgHARCM

MAgHARCM implements the complete Chikofsky & Cross modernization lifecycle:
- Stage 1 (`[[1.0.0 PRIM-14]]` Archaeology & `[[1.0.0 PRIM-20]]` Concept Assignment) is pure **Reverse Engineering** and **Design Recovery**.
- Stage 2 & 3 (`[[1.0.0 PRIM-21]]` Strategy Selection & `[[1.0.0 PRIM-1]]` Planning) represent **Restructuring**.
- Stage 4 (`[[1.0.0 PRIM-23]]` Chunked Translation) executes **Forward Engineering** into the target language.

## Hop-1 References

- [[Parnas-1972]] — Provides the modular boundaries that design recovery aims to reconstruct.
- [[Müller-2000]] — Expands Chikofsky's taxonomy into an actionable reverse engineering roadmap.
- [[Rajlich-1997]] — Uses design recovery as the basis for concept-locator analysis.

## Hop-2 Anchors (Software Archaeology Lean)

- [[pp-besm-Software-Archaeology]] — Translates Chikofsky's Design Recovery into terminal-based developer workflows.
- [[Feathers-2004-WELC]] — Restructuring preserves behavioral invariants through characterization tests.
