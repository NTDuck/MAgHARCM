---
title: "P-70 — Baldwin & Clark 2006 — Design Rules: The Power of Modularity"
backlink: "[[1.0.0 P-70]]"
tags: [paper, modularity, design-rules, software-economics, [[1.0.0 P-41]], hop-2]
---

# [[1.0.0 P-70 — Design Rules Modularity]]

## Citation

Baldwin, C. Y., & Clark, K. B. (2006). *The Power of Modularity: Designing the Right Architecture at the Right Time / Design Rules, Vol. 1*. MIT Press. (Original "Design Rules" project 1997-2000; republished as a single volume 2006.) Working paper version: Harvard Business School Working Paper 96-018, 1996. URL: https://www.researchgate.net/publication/247906506.

## Summary

The book formalises the **economic theory of modular architectures**. The core argument: a complex system can be decomposed into **modules** that hide their internal design from each other and interact only through **design rules** (visible interface specifications, hidden implementation rules). The visible/hidden distinction is what allows **parallel exploration**: each module owner can innovate independently without coordinating with the others, as long as the design rules remain stable.

The book's quantitative contribution is the **option-value of modularity**: in a world with uncertainty about the best implementation, the value of being able to experiment in parallel is roughly the variance of the outcome distribution times the time saved. Modular systems are **option-rich** (the design rules create real options for substitution), and integrated systems are **option-poor** (any change requires re-coordination).

## Method

The book formalises the visible/hidden distinction using:

- **Six types of design rules**: architecture (overall structure), interfaces (module interconnection), standards (interoperability metrics), testing, process, and design rules for a "tested module" itself.
- **Operators of modular evolution**: split (a module is split when its internal complexity exceeds its option-value), substitute (swap one module for another), augment (add a new module), exclude (remove a module), invert (change which design rules are visible vs. hidden).

## Findings Relevant to MAgHARCM

- **[[1.0.0 PRIM-19]] Design Rule Hierarchy Layering** is the application of Baldwin & Clark's design-rule taxonomy to legacy codebases: the Archaeologist agent extracts the visible/hidden boundary, classifies each interface rule into one of the six types, and reports the module decomposition.
- **[[1.0.0 PRIM-20]] Concept Assignment & Redocumentation** maps to the "operators of modular evolution": each detected concept-cluster becomes a candidate for `split`, `substitute`, `augment`, `exclude`, or `invert`.
- **Option-value of modularity** is the empirical justification for MAgHARCM's iteration loop: the verification pass tests whether the new translation preserves the design rules; if it does, the modularity is preserved and parallel exploration can continue.
- **The 1997 working paper** is the primary citation in MAgHARCM's Design-Rule extraction algorithm.

## How MAgHARCM Uses It

The Archaeologist agent's `extractDesignRuleHierarchy` (cf. `internal/agents/archaeology.go`) uses the six-design-rule taxonomy as its categorisation scheme. The optional-checks agent runs a per-cluster option-value computation; clusters with high option-value receive the highest translate-priority score.

## References

### Hop-1 (Baldwin & Clark 2006 cites)
- Simon, H. A. (1962). *The Architecture of Complexity*. Proceedings of the American Philosophical Society 106(6):467-482. See [[1.0.0 P-31]].
- Parnas, D. L. (1972). *On the Criteria to Be Used in Decomposing Systems into Modules*. CACM 15(12):1053-1058. See [[1.0.0 P-31]].
- Langlois, R. N. (2002). *Modularity in Technology and Organization*. Journal of Economic Behavior & Organization 49:19-37.
- Schilling, M. A. (2000). *Toward a General Modular Systems Theory and Its Application to Interfirm Product Modularity*. Academy of Management Review 25(2):312-334.

### Hop-2
- Baldwin, C. Y. (2008). *Where do Transactions Come From? Modularity, Transactions, and the Boundaries of Firms*. Industrial and Corporate Change 17(1):155-195.
- Fleming, L. & Baldwin, C. Y. (2024). *The Architecture of Complexity Revisited*. (See [[1.0.0 P-71]].)
- MacCormack, A., Baldwin, C. Y., Rusnak, J. (2012). *Exploring the Duality between Product and Organizational Architectures*. HBS Working Paper 12-099. See [[1.0.0 P-42]].

## Backlinks

[[1.0.0 P-34]], [[1.0.0 P-41]], [[1.0.0 P-42]], [[1.0.0 PRIM-19]], [[1.0.0 PRIM-20]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-70 is the **modularity-economics anchor** for MAgHARCM's Design Rule Hierarchy Partitioning.
