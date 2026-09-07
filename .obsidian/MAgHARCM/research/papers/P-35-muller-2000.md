---
title: Reverse Engineering - A Roadmap
bibkey: p35_muller
tags: [paper, roadmap, reverse-engineering, migration-strategies, [[1.0.0 PRIM-21]], hop-1]
---

# [[1.0.0 P-35]] Reverse Engineering: A Roadmap

**Authors**: Hausi A. Müller, Jens H. Jahnke, Dennis B. Smith, Margaret-Anne Storey, Scott R. Tilley, Kenny Wong  
**Year**: 2000  
**Venue**: The Future of Software Engineering (ICSE 2000), ACM Press, pp. 47–60  
**eprint / DOI**: 10.1145/336512.336526  
**Cited by**: `[[2.0.0 Software-Archaeology-Lineage]]`, `[[1.0.0 PRIM-21]]` (Migration Strategy Selection)

## Summary

Müller et al. provide a comprehensive roadmap synthesizing the state of reverse engineering research at the turn of the millennium. The paper decomposes program comprehension into three interacting dimensions:
1. **Cognitive Aspects**: Mental models (top-down hypothesis-driven models by Soloway, bottom-up abstraction models by Pennington, and integrated meta-models by Mayrhauser & Vans).
2. **Analysis Techniques**: Static parsing, AST extraction, program slicing, control-flow analysis, and dynamic tracing.
3. **Tool Interoperability**: Common exchange formats (GXL, RSF) and visualization environments.

The paper formalizes the canonical **Five Legacy Migration Strategies**:
1. *Big Bang (Cold Turkey)*: Monolithic replacement in one stroke; appropriate only for small, bounded systems.
2. *Incremental (Chicken Little)*: Evolutionary component-by-component migration in topological order.
3. *Pilot Subsystem*: Isolating and migrating a representative low-risk subsystem before full cutover.
4. *Frozen Legacy*: Freezing the legacy codebase and developing characterization tests before attempting any transformation.
5. *Parallel Cutover*: Running legacy and migrated systems concurrently in production with differential comparison before cutover.

## Relevance to MAgHARCM

Direct intellectual source of `[[1.0.0 PRIM-21]]` (Migration Strategy Selection). MAgHARCM's `internal/agents/strategy.go` implements Müller's five strategies as a dynamic try-and-fail registry:
- Tiny projects match Big Bang.
- Untested systems match Frozen Legacy (requiring `[[1.0.0 PRIM-5]]` test synthesis).
- Large codebases trigger Pilot translation.
- Comprehensive modular suites trigger Parallel Cutover.
- All systems fall back safely to Incremental translation.

## Hop-1 References

- [[Chikofsky-Cross-1990]] — Foundation for the reverse engineering roadmap.
- [[Baldwin-Clark-2000]] — Provides the modular substitution mechanics for incremental migration.
- [[Foltz-2023]] — Deepens Müller's cognitive models into the modern DR. JONES traversal model.

## Hop-2 Anchors (Software Archaeology Lean)

- [[pp-besm-Software-Archaeology]] — Translates Müller's cognitive models into practical developer investigation passes.
- [[Feathers-2004-WELC]] — Operationalizes Frozen Legacy by putting characterization tests in place first.
