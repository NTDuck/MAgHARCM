---
title: "Kazman et al. — A Case Study in Locating the Architectural Roots of Technical Debt"
bibkey: p51_kazman_drspaces
tags: [paper, design-rules, architectural-recovery, technical-debt, load-bearing, [[PRIM-19]], hop-1]
---

# [[1.0.0 P-20]] Kazman et al. — A Case Study in Locating the Architectural Roots of Technical Debt

**Authors**: Rick Kazman, Yuanfang Cai, Ran Mo (and collaborators, including SoftServe engineering)
**Year**: 2015
**Venue**: Proceedings of the 37th IEEE/ACM International Conference on Software Engineering (ICSE 2015), Software Engineering in Practice (SEIP) track, pp. 179–188
**eprint / DOI**: https://doi.org/10.1109/ICSE.2015.139 ; SEI summary at sei.cmu.edu/library/locating-the-architectural-roots-of-technical-debt/ ; PDF at cs.drexel.edu/~yfcai/papers/2015/icse2015-seip.pdf
**Cited by**: [[primitives/INDEX]] entry [[PRIM-19]] (Design Rule Hierarchy Partitioning); references L1/L2/L3 stability layers and the load-bearing concentration check

## Summary

[[Kazman-2015-CaseStudy]] (ICSE 2015 SEIP) is the canonical industrial case study of the *design-rule-space* approach to architectural recovery. Working with SoftServe Inc. on a large project with mounting technical debt, the authors reconstruct the system's architecture as a layered design-rule space: **L1** rules (the small, stable set of load-bearing decisions that constrain downstream work — module interfaces, key abstractions, build invariants), **L2** rules (subsystem-level structural decisions, derived from L1), and **L3** rules (leaf-level implementation choices that are trivially substitutable). They then introduce the **load-bearing concentration check**: a small set of files/symbols that account for a disproportionate fraction of the project's defects is a marker that those files/symbols are doing implicit L1-rule work without being recognised as such. The case study demonstrates that this single analysis pinpoints where re-architecting effort will pay off, and that the same framework predicts future defect clusters by combining change-history with design-rule assignment. The paper is the empirical anchor for the "stability layering + concentration" diagnostic that MAgHARCM uses under [[PRIM-19]].

## Relevance to MAgHARCM

[[PRIM-19]] cites Kazman et al. for the L1/L2/L3 classification and the load-bearing concentration check. MAgHARCM's planned `internal/agents/cpg.go` ([[PRIM-9]] tri-representation code graph) should annotate every recovered architectural element with a stability layer and a concentration score; the planning agent ([[PRIM-1]]/[[PRIM-2]]) should refuse to translate an L1-rule element until the operator explicitly signs off, because translation errors at L1 propagate to the entire system. The `internal/agents/analyzer.go` ([[PRIM-21]] migration-strategy selection) should make Kazman's concentration check one of its inputs — a project where a tiny set of files dominates the change history is a strong candidate for the Integrate-in-Place strategy ([[Müller-2000-IWPC]]) rather than Cold Turkey, because re-architecting those L1 nodes is too risky for a clean cutover.

## Hop-1 References

- [[Baldwin-Clark-2000-DesignRules]] — the conceptual ancestor; Kazman's L1/L2/L3 layering is the operationalisation of Baldwin & Clark's "stable intermediate forms" applied to an existing legacy system.
- [[Müller-2000-IWPC]] — Kazman's concentration check feeds Müller's LIS decomposability assessment: a project with concentrated load-bearing nodes is structurally *less* decomposable and demands a different migration strategy.
- [[Rajlich-1997-ICSE]] — Kazman's design-rule hierarchy is one possible *output* of Rajlich's concept-locator analysis: the concepts Rajlich finds become L1 rules once their downstream impact is quantified.
- [[Feathers-2004-WorkingEffectivelyWithLegacyCode]] — Feathers' seams and characterisation tests provide the practical mechanism for safely refactoring L1 nodes once Kazman has identified them.
- [[Soliman-2025-MSR4SA]] — MSR4SA's co-change mining is the empirical signal that Kazman-style concentration checks consume; together they form a static-skeleton + dynamic-coupling recovery stack.

## Hop-2 Anchors (software-archaeology lean)

- [[Kazman-Cai-2024-ArchitecturalRecovery]] — Kazman & Cai's later work formalises the design-rule-space decomposition into a tool-supported recovery procedure; MAgHARCM should consume that procedure rather than re-implement the 2015 checklist.
- [[Rajlich-1997-ICSE]] — Kazman's L1 nodes are best discovered *after* Rajlich's concept map is built, because the concept map explains *why* a given node is load-bearing rather than just *that* it is.
- [[Müller-2000-LegacyIntegration]] — the load-bearing concentration check directly drives Müller's strategy selector: high concentration → Integrate-in-Place, low concentration → Cold Turkey, mixed → Chicken Little.
