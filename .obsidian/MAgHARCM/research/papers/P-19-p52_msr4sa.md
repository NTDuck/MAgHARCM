---
title: "MSR4SA — Mining Software Repositories for Software Architecture: A Systematic Mapping Study"
bibkey: p52_msr4sa
tags: [paper, systematic-mapping, mining-software-repositories, architecture-recovery, jaccard-coupling, [[PRIM-18]], hop-1]
---

# [[1.0.0 P-19]] MSR4SA — Mining Software Repositories for Software Architecture: A Systematic Mapping Study

**Authors**: Mohamed Soliman, Michel Albonico, Ivano Malavolta (and collaborators)
**Year**: 2025
**Venue**: Information and Software Technology (Elsevier), Volume 181, Article 107677
**eprint / DOI**: https://doi.org/10.1016/j.infsof.2025.107677 ; open-access at ivanomalavolta.com/files/papers/IST_MSR4SA_2025.pdf ; replication package at zenodo.org/records/14614349
**Cited by**: [[primitives/INDEX]] entry [[PRIM-18]] (Jaccard-Coupling Architecture Recovery); referenced commit-neighbourhood Jaccard similarity for hidden coupling detection

## Summary

[[Soliman-2025-MSR4SA]] (IST 2025) is the first systematic mapping study of *Mining Software Repositories* (MSR) techniques applied to software-architecture activities. The authors survey the literature through a structured protocol and classify every primary study along four axes: which architecture activity it supports (recovery, tactics identification, smell detection, evolution analysis, conformance checking), which MSR data source it mines (commit logs, issue trackers, mailing lists, code-review history, release notes), which mining technique it uses (co-change mining, topic models, information retrieval, network analysis, process metrics), and at which architectural level it operates (component, connector, configuration, rationale). The mapping surfaces a dominant pattern: **co-change mining** (often cast as Jaccard similarity over commit neighbourhoods) is by far the most popular technique for hidden-coupling recovery, because static import edges systematically under-report the coupling that emerges through shared change histories. The authors package a replication artefact on Zenodo (S2-group/msr4sa-systematic-mapping-study-rep-pkg) so future MSR4SA work can be reproduced.

## Relevance to MAgHARCM

[[PRIM-18]] cites MSR4SA for the specific technique of *Jaccard similarity over commit neighbourhoods* to detect coupling invisible to static imports. MAgHARCM's `internal/agents/analyzer.go` ([[PRIM-21]] migration-strategy selection) and the planned `internal/agents/cpg.go` ([[PRIM-9]] tri-representation code graph) both need a "hidden coupling" signal that pure static analysis cannot supply — particularly for the commons-validator and Apache HTTP client samples, where modules that never `import` each other still get modified together. The MSR4SA protocol's axes (activity × data source × technique × level) give MAgHARCM a vocabulary for documenting which MSR signals it consumes and at which architectural level each signal is fed. The replication package on Zenodo is the right starting point for a per-repository hidden-coupling analyser that mines `git log --name-only` over the source project and emits a Jaccard-similarity matrix the planning agent can read.

## Hop-1 References

- [[Kazman-Cai-2015-ICSE-SEIP]] — same architectural-recovery lineage; MSR4SA frames Kazman & Cai's design-rule-space recovery as one of the activities MSR techniques can support.
- [[Rajlich-1997-ICSE]] — concept-locator analysis; MSR4SA's "co-change mining" is the empirical operationalisation of Rajlich's claim that concepts live in change histories, not in identifier names.
- [[Müller-2000-IWPC]] — MSR4SA classifies legacy-system evolution analysis as one of the four architecture activities; Müller's LIS decomposability assessment is the target decision that MSR signals feed into.
- [[Khajeh-Hosseini-2019-AutomatedTransformation]] — modernising legacy systems requires MSR signals to plan the migration order; MSR4SA's mapping is the literature backbone that justifies mining commit histories before touching code.
- [[Feathers-2004-WorkingEffectivelyWithLegacyCode]] — Feathers' seam-finding is a static analogue; MSR4SA's co-change mining is the dynamic (history-based) counterpart.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000-DesignRules]] — MSR4SA's Jaccard co-change signal recovers *implicit* design rules that the source code's static structure hides; per Baldwin & Clark, those implicit rules are exactly the load-bearing knowledge that migration must preserve.
- [[Kazman-Cai-2024-ArchitecturalRecovery]] — the design-rule-space decomposition provides the architectural skeleton that MSR4SA's co-change matrix populates with empirical coupling weights; together they yield a layered (static-skeleton + dynamic-coupling) recovery.
- [[Foltz-2023-DR-JONES]] — DR.JONES' "breadth before depth, recency bias" heuristics predict that a Jaccard-coupling matrix should be *first* read at a coarse module granularity (breadth), then refined — MSR4SA does not currently enforce this and MAgHARCM should.
