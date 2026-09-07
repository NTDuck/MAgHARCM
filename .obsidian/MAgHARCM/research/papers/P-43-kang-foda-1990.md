---
title: Kang et al. — Feature-Oriented Domain Analysis (FODA) Feasibility Study (SEI 1990)
bibkey: p43_kang_foda_1990
tags: [paper, feature-modeling, domain-analysis, foda, [[1.0.0 PRIM-10]], hop-1]
---

# [[1.0.0 P-43]] Kang et al. — Feature-Oriented Domain Analysis (FODA) Feasibility Study (SEI 1990)

**Authors**: Kyo C. Kang, Sholom G. Cohen, James A. Hess, William E. Novak, A. Spencer Peterson (Software Engineering Institute, Carnegie Mellon)
**Year**: 1990 (November)
**Venue**: SEI Technical Report CMU/SEI-90-TR-21, ESC-TR-90-21
**eprint / DOI**: DTIC ADA235785; SEI institutional repository
**Cited by**: [[primitives/INDEX]] entry [[PRIM-10]] (Feature-Mapping & Type-Compatibility Validation) — the canonical feature-modeling hop-1 reference.

## Summary

[[Kang-1990-FODA]] introduces the Feature-Oriented Domain Analysis (FODA) methodology, which models a software product family as a tree of features — distinguishable end-user-visible characteristics of the system — and analyses commonality and variability across the family. The methodology has three phases:

1. *Context analysis* — identify the domain boundary and the actors that interact with the family.
2. *Domain modelling* — extract the feature model: a hierarchical tree of mandatory, optional, alternative, and mutually-exclusive features.
3. *Architecture modelling* — derive a reference architecture that realises the feature model; commonality is realised by shared components, variability by parameterised or alternative ones.

The FODA feature model is the antecedent of every modern feature-model formalism (formally analysed by Batory 2005, van der Storm 2007, Apel 2009). Its primitives — mandatory / optional / alternative / mutually-exclusive features — are still the vocabulary of feature modelling in software product-line engineering (SPLE) and in modern model-driven software engineering tools.

The paper also introduces the **capability matrix**, an early artefact for mapping feature-model capabilities onto architectural components. The capability matrix is the direct ancestor of MAgHARCM's [[PRIM-10]] feature-mapping step.

## Relevance to MAgHARCM

1. [[PRIM-10]] (Feature-Mapping & Type-Compatibility Validation, `internal/agents/feature_mapping.go`): the feature-mapping algorithm uses FODA-style capability matrices to record which source-language library maps to which target-language crate. The mandatory/optional/alternative feature classification is the basis for distinguishing *must-preserve* features from *nice-to-have* ones during translation.
2. [[PRIM-3]] (Target Skeleton-First Generation): the FODA reference architecture is the conceptual ancestor of MAgHARCM's Rust-trait skeleton. Both emit a stable interface first, then parameterise the implementation per feature.
3. [[PRIM-26]] (Symbol-Aware Navigator, ABCoder-based): the navigator's per-feature symbol resolution is the modern realisation of FODA's capability matrix lookup.
4. [[P-05]] (Oxidizer, Siddiq et al. 2024): Oxidizer's feature-mapping step is a contemporary descendant of FODA, applied to the C→Rust translation case.

## Hop-1 References

- [[Batory-2005-Feature-Models]] — feature algebras; the mathematical formalisation of FODA feature-model primitives.
- [[Apel-2009-Feature-Oriented-Programming]] — feature modules and feature composition; FODA's "commonality + variability" split as a programming-language concept.
- [[Czarnecki-Generative-Programming]] — generative programming; the broader umbrella under which FODA falls.
- [[Czarnecki-2000-Generative-Programming-Map]] — the canonical taxonomy that places FODA in the wider generative-programming landscape.
- [[Pohl-2005-Software-Product-Line-Engineering]] — SPLE textbook; FODA is the first chapter.
- [[Donohoe-1990-Software-Reuse]] — earlier reuse methodology; FODA is its generalisation to families.

## Hop-2 Anchors (software-archaeology lean)

- [[Muller-2000-Strategies]] — Müller migration strategies map onto FODA's reference architecture: a successful migration identifies the feature model and reuses the reference architecture; a stalled migration fails to identify commonality.
- [[Baldwin-Clark-2000]] — design rules correspond to mandatory features in FODA: the L1 interface is the set of features the family cannot vary.
- [[Rajlich-1997-Staged-Life-Cycle]] — concept assignment is the offline extraction of feature-model nodes from source code.
- [[Parnas-1972]] — information hiding as the design rule; FODA's mandatory features are Parnas's hidden modules.
- [[Kazman-Cai-2017-DRSpaces]] — design-rule extraction is the architecture-level analogue of feature-model extraction; both operate on the same L1/L3 partition.

## Backlinks

- [[METHODOLOGY]] §PRIM-10 — Kang (1990) FODA is the canonical feature-model hop-1 reference.
- [[primitives/INDEX]] — [[1.0.0 PRIM-10]] status entry now references [[P-43]] alongside [[P-05]] (Oxidizer) and [[Czarnecki]].
- [[Software-Archaeology-Lineage]] §4 — FODA feature modelling is the offline design-recovery technique that PRIM-10's per-feature mapping consumes.
