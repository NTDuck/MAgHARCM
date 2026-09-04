---
title: Program Slicing
bibkey: p28_weiser
tags: [paper, program-analysis, slicing, dependency, [[PRIM-2]], [[PRIM-9]], hop-1]
---

# Program Slicing

**Authors**: Mark Weiser (Xerox PARC; later University of Maryland)
**Year**: 1984
**Venue**: IEEE Transactions on Software Engineering, Vol. SE-10, No. 4, pp. 352–357, July 1984
**eprint / DOI**: DOI: 10.1109/TSE.1984.5010248 (IEEE TSE)
**Cited by**: [[primitives/INDEX]] entries [[PRIM-2]] (Back-Edge-Conditioned Reverse-Topological Scheduling — data-flow preconditioning), [[PRIM-9]] (Tri-Representation Hybrid Code Graph — slicing layer for PDG/SDG)

## Summary

[[Weiser-1984-Slicing]] introduces *program slicing*: given a variable of interest at a program point, compute the minimal subset of statements that could affect its value. Weiser formalises slicing as a data-flow fixed-point over the control-flow graph, distinguishes static from dynamic slicing, and proves that the slice is computable, conservative, and often dramatically smaller than the original program. The paper also demonstrates the first practical applications: debugging, program comprehension, and (forward-looking) software maintenance. Slicing is the foundational primitive on which every later inter-procedural, object-oriented, and concurrent slicing technique is built (Horwitz-Reps-Binkley SDG, Ottenstein-PDG, Agrawal-Horgan dynamic slicing, etc.). It is also one of the few program-analysis papers whose algorithmic skeleton still ships in production tools 40 years later — `gdb`'s reverse-finish, IDE "find usages", and modern call-graph preprocessors are all direct descendants.

## Relevance to MAgHARCM

Weiser's slicing primitive is the data-flow precondition for two MAgHARCM primitives. [[PRIM-2]] (Back-Edge-Conditioned Reverse-Topological Scheduling) uses a forward slice from the call-graph entry points to enumerate the statements whose values can reach a callee; only those statements are part of the reverse-topo order's "context", so the [[PRIM-23]] chunked translator sees a tight fragment. [[PRIM-9]] (Tri-Representation Hybrid Code Graph) embeds slicing into the PDG / SDG layers so the [[PRIM-26]] Navigator can answer "what statements affect this variable?" by traversing the slice rather than the entire AST+CFG+PDG union. The Tri-Representation graph (AST + CPG + SDG) only becomes tractable at repository scale because slicing collapses most of the graph away from the locus of interest; without [[Weiser-1984-Slicing]] the [[PRIM-26]] Navigator's ≤ 4 KB context budget could not be honoured. Slicing also underlies the [[PRIM-14]] Software-Archaeology Stage's naming & type forensics: a slice from a renamed symbol's uses is the canonical recovery primitive when the source's original names are lost.

## Hop-1 References

- [[Ottenstein-1984-PDG]] — Program Dependence Graph; the immediate successor that fuses control-flow and data-flow into a single directed graph Weiser's slicing rides on.
- [[Horwitz-Reps-Binkley-1990-SDG]] — inter-procedural System Dependence Graph; the standard extension cited by MAgHARCM [[PRIM-9]] and made efficient only because slicing is built in.
- [[Agrawal-Horgan-1990-Dynamic-Slicing]] — dynamic slicing that conditions on a concrete execution trace; relevant to [[PRIM-4]] SpecMiner-style dynamic invariant recovery.
- [[Tip-1994-Survey]] — early survey of slicing algorithms and applications; situates Weiser's contributions.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997]] — concept-locator analysis is the human-readable wrapper around automatic slicing: the slice is what a concept "owns", and the concept is what the slice "is about".
- [[Muller-2022]] — legacy modernization case studies cite slicing as the canonical first-pass preprocessor before any other migration activity, validating MAgHARCM's choice to embed slicing in [[PRIM-2]] rather than in a separate stage.
- [[Kazman-Cai-2024]] — architectural recovery uses slicing to bound the impact set of a proposed refactor; the same technique that justifies the Navigator's tight context windows.
