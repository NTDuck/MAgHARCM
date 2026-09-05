---
title: MacCormack, Rusnak & Glassman — Exploring the Structure of Complex Software Designs (Design Structure Matrix)
bibkey: p42_maccormack_dsm_2006
tags: [paper, design-structure-matrix, technical-debt, modularity, [[1.0.0 PRIM-3]], hop-1]
---

# [[1.0.0 P-42]] MacCormack, Rusnak & Glassman — Exploring the Structure of Complex Software Designs

**Authors**: Alan MacCormack, John Rusnak, Carliss Y. Baldwin (HBS)
**Year**: 2006 (HBS Working Paper 06-016); later published in *Management Science* 2007
**Venue**: Harvard Business School Working Paper 06-016; Management Science 53(5), 2007
**eprint / DOI**: DOI 10.1287/mnsc.1060.0699 (Management Science)
**Cited by**: [[primitives/INDEX]] entry [[PRIM-3]] (Target Skeleton-First Generation) — provides the empirical Design Structure Matrix (DSM) analysis that motivates the skeleton-first synthesis order.

## Summary

[[MacCormack-2006-DSM]] analyses the dependency structure of two large open-source systems (Linux kernel and Apache HTTP server) using the Design Structure Matrix (DSM) methodology. A DSM is a square matrix whose rows and columns correspond to system components (files, modules, subsystems) and whose cells encode the directed dependencies between them. Clusters of mutually-dependent components are flagged as "cycles"; the size and number of cycles correlate with the system's maintenance burden.

Three empirical findings anchor the paper:

1. *Cycle size and number predict release-defect density*: larger and more numerous dependency cycles correlate with higher post-release bug counts in the Linux kernel. Cycles are an empirically validated technical-debt signal.
2. *Interface stability precedes mass refactoring*: stable interface layers (L1 in Baldwin-Clark terms) emerge before the bulk of internal restructuring, mirroring the design-rule stability claim of [[Baldwin-Clark-2000]] but with measurable empirical support.
3. *Refactoring reduces cycle density over time*: the Linux kernel's DSM cycle count decreases monotonically across major releases, indicating that explicit refactoring pass — not organic evolution — is required to keep dependency structure healthy.

The paper closes with a DSM-based refactoring strategy: identify high-cycle components, freeze their interface (becoming a new L1 rule), then extract sub-modules to break the cycle. This three-step pattern is the operational form of [[PRIM-3]]'s skeleton-first synthesis.

## Relevance to MAgHARCM

1. [[PRIM-3]] (Target Skeleton-First Generation, `internal/agents/planning.go::DefaultProjectSkeleton`): the skeleton-first synthesis order is the empirical refactoring strategy from [[MacCormack-2006-DSM]] translated into a translation pipeline. The Rust traits emitted before method bodies are the L1 interface freeze step; method bodies come after, when their dependencies are stable.
2. [[PRIM-19]] (Design Rule Hierarchy Partitioning, `internal/agents/design_rule_hierarchy.go`): the per-file stability classification (L1/L2/L3) uses the empirical finding that L1 stability precedes refactoring as its scoring heuristic.
3. [[PRIM-1]] (Reverse Topological Ordering, `internal/agents/planning.go`): the topological synthesis order reflects the DSM insight that downstream components (those with high in-degree) should be generated last, after their dependencies are stable.
4. [[P-19]] (MSR4SA, Soliman 2025): the MSR4SA systematic mapping classifies commit-co-change mining as one MSR technique; the DSM analysis in this P-42 entry is the static-analysis complement.

## Hop-1 References

- [[Baldwin-Clark-2000]] — design-rule decomposition; the conceptual ancestor of the L1/L2/L3 distinction this paper measures empirically.
- [[Steward-1981-DSM]] — the original Design Structure Matrix methodology for system decomposition, transferred from mechanical engineering to software.
- [[Browning-2001-Applying-DSM]] — applied DSM tutorial; the methodological backbone this paper draws on.
- [[Yoder-1997-Architecture-Recovery]] — architecture recovery via DSM in object-oriented systems; an earlier empirical application.
- [[Sangal-2005-Using-DSM-Dependencies]] — using DSM to model software package dependencies; the direct methodological ancestor of the Linux/Apache DSM analysis.
- [[Parnas-1972]] — information hiding as the design-rule kernel; the L1 stability claim is a measurable operationalisation of Parnas's interface-stability criterion.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000]] — design-rule decomposition predicts that L1 stability precedes internal restructuring; [[MacCormack-2006-DSM]] supplies the empirical measurement.
- [[Lehman-1980-Laws-Evolution]] — Lehman's "increasing complexity" law is the dynamic counterpart of the static DSM cycle-density finding.
- [[Rajlich-1997-Staged-Life-Cycle]] — concept assignment is the offline analogue of DSM cycle identification: a cycle in the DSM signals a missing concept assignment.
- [[Muller-2000-Strategies]] — Müller's five migration strategies map onto DSM cycle-resolution patterns: Big Bang ignores cycles, Incremental resolves them incrementally, Frozen freezes the cycle boundary first.
- [[Kazman-Cai-2017-DRSpaces]] — design-rule-space extraction uses DSM clustering as one of its inputs.

## Backlinks

- [[METHODOLOGY]] §PRIM-3 — [[1.0.0 P-42]] MacCormack DSM is the empirical anchor for the skeleton-first synthesis order.
- [[primitives/INDEX]] — [[1.0.0 PRIM-3]] status entry now references [[P-42]] alongside [[Baldwin-Clark-2000]] and [[ReCodeAgent]].
- [[Software-Archaeology-Lineage]] §4 — DSM cycle detection is one of the offline recovery techniques that PRIM-19's partition consumes.
