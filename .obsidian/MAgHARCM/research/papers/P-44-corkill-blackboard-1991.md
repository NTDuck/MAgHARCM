---
title: Corkill — Blackboard Systems (AI Expert 1991)
bibkey: p44_corkill_blackboard_1991
tags: [paper, blackboard-architecture, multi-agent, knowledge-source, [[1.0.0 PRIM-17]], hop-1]
---

# [[1.0.0 P-44]] Corkill — Blackboard Systems (AI Expert 1991)

**Authors**: Daniel D. Corkill (University of Massachusetts Amherst)
**Year**: 1991 (September / October)
**Venue**: *AI Expert*, Vol. 6, No. 9, pp. 40-47; reprinted in *Blackboard Systems* (Engelmore & Morgan, eds.) 1988 / 1997
**eprint / DOI**: AAAI archive copy; published by Academic Press
**Cited by**: [[primitives/INDEX]] entry [[PRIM-17]] (Asynchronous SE Agent Blackboard) — the canonical practical hop-1 reference for blackboard systems.

## Summary

[[Corkill-1991-Blackboard-Systems]] is the practitioner-level survey of blackboard architectures, complementary to the canonical Nii 1986 "Blackboard Systems" two-part paper in *AI Magazine*. Corkill organises the survey around five recurring components:

1. *Knowledge sources (KSs)* — independent, specialised modules that produce partial solutions when triggered by blackboard state.
2. *Blackboard* — a shared data structure holding the partial solution in a hierarchy of abstraction levels.
3. *Control shell* — the scheduler that selects which KS to invoke next based on blackboard state and KS preconditions.
4. *Knowledge-base interaction protocol* — the structured way KSs read and write blackboard entries.
5. *Triggering / event mechanism* — the precondition evaluation that activates a KS when its input conditions are met.

Corkill identifies four practical lessons that recur across Hearsay-II, BB1, GBB, and every blackboard system that followed:

1. *Domain-driven partitioning*: knowledge sources should mirror the domain decomposition, not the runtime call structure.
2. *Abstraction-level hierarchy*: the blackboard should be organised by abstraction level (signal → segment → phrase → word) so KSs trigger only when data at the right level appears.
3. *Scheduler is policy, not mechanism*: the control shell's selection logic is policy that can be replaced without touching KSs.
4. *Granularity trade-off*: too many small KSs maximise parallelism but explode the scheduler search space; too few large KSs minimise scheduling overhead but lose the modularity that justified the blackboard in the first place.

## Relevance to MAgHARCM

1. [[PRIM-17]] (Asynchronous SE Agent Blackboard, `internal/agents/blackboard.go`): the blackboard implementation is a direct realisation of the five-component architecture. Knowledge sources map to the eight agent units (Archaeologist, Analyzer, Planner, Translator, RoleFlipReviewer, Validator, VerdictPanel, Recruiter); the blackboard is the shared `State` struct; the control shell is the Eino graph scheduler; the trigger mechanism is the per-iteration precondition evaluation.
2. [[P-18]] (CAID-2024, Multi-Agent SE CAID): CAID uses Nii 1986 as its primary citation; Corkill 1991 is the more accessible practitioner-level companion, the one most CS practitioners actually read.
3. [[P-11]] (MetaGPT-2024 SOP contracts): MetaGPT's "publish-subscribe over shared state" pattern is a blackboard variant with a restricted set of allowed KSs (the SOP-defined roles).
4. [[P-26]] (AgentVerse-2023, dynamic recruitment): AgentVerse's "post-iteration evaluation" is a blackboard trigger mechanism (lesson 4 above: granularity trade-off).

## Hop-1 References

- [[Nii-1986-Blackboard-Part-1]] — the canonical blackboard systems survey in *AI Magazine*; Corkill's article is the practitioner summary.
- [[Nii-1986-Blackboard-Part-2]] — blackboard application case studies; the Hearsay-II and BB1 analyses.
- [[Engelmore-Morgan-1988-Blackboard-Systems]] — the canonical blackboard-systems anthology; includes Hearsay-II, BB1, and the original O-Plan chapters.
- [[Hayes-Roth-1985-Blackboard-Control]] — BB1 control shell; the policy-vs-mechanism lesson (3) is from this paper.
- [[Lesser-1975-Hearsay-II]] — the canonical blackboard system for speech understanding; the first system to demonstrate the abstraction-level hierarchy (lesson 2).
- [[Erman-1980-Hearsay-II-Speech]] — the published Hearsay-II account; the original demonstration of knowledge-source partitioning.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997-Staged-Life-Cycle]] — staged lifecycle maps to abstraction levels on the blackboard: each stage promotes the artefact from one abstraction level to the next.
- [[Baldwin-Clark-2000]] — knowledge sources map to L3 leaves (substitutable modules); the blackboard itself is the L1 design rule (shared state schema).
- [[Muller-2000-Strategies]] — migration strategies are knowledge sources that the blackboard's control shell selects between; MAgHARCM's [[PRIM-21]] Registry is the policy implementation.
- [[Lehman-1980-Laws-Evolution]] — Lehman's "increasing complexity" is the dynamic counterpart of blackboard bloat: as knowledge sources multiply, the scheduler search space grows super-linearly.
- [[Kazman-Cai-2017-DRSpaces]] — design-rule extraction operates on the blackboard's shared-state schema, which is itself a design-rule artefact.

## Backlinks

- [[METHODOLOGY]] §PRIM-17 — Corkill (1991) is the practitioner-level blackboard-systems reference; complements Nii (1986) as the canonical pair.
- [[primitives/INDEX]] — [[1.0.0 PRIM-17]] status entry now references [[P-44]] alongside [[P-18]] (CAID) and Nii (1986).
- [[Software-Archaeology-Lineage]] §4 — blackboard architecture is the orchestration pattern that PRIM-17's eight-agent pipeline consumes.
