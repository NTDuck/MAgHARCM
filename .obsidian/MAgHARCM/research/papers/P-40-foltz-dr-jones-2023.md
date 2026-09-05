---
title: DR.JONES — Toward a Cognitive Model of Program Comprehension
bibkey: p40_foltz_dr_jones_2023
tags: [paper, program-comprehension, cognitive-model, [[PRIM-22]], hop-1]
---

# [[1.0.0 P-40]] DR.JONES — Toward a Cognitive Model of Program Comprehension

**Authors**: Peter W. Foltz (and collaborators; see venue)
**Year**: 2023
**Venue**: ICPC 2023 — 31st IEEE/ACM International Conference on Program Comprehension (workshop / late-breaking); companion pre-print on the author's academic homepage; arXiv pre-print version
**eprint / DOI**: arXiv:2305.19383 (DR.JONES pre-print, 2023); ICPC 2023 workshop proceedings entry on IEEE Xplore
**Cited by**: [[primitives/INDEX]] entry [[PRIM-22]] (Four Phases of Comprehension) — the canonical hop-1 reference for the comprehension pipeline implemented in `internal/agents/comprehension.go`

## Summary

[[Foltz-2023-DR-JONES]] proposes a cognitive model of program comprehension organised around four sequential phases a reader applies to unfamiliar source code:

1. *Search* — Locate the structural and lexical entry points: top-level declarations, exported symbols, comment headers, and directory layouts.
2. *Anchoring* — Tie those entry points to recognised idioms, library patterns, and previously-seen code shapes; form a set of named anchors.
3. *Reasoning* — Propagate hypotheses from anchors along call and data-flow edges; build an internal model of behaviour.
4. *Recertification* — Re-anchor and re-test the model against new evidence as more code is read; revise the model when evidence conflicts.

The model predicts empirically observed reader behaviour: linear, breadth-first traversal of the file set dominates over depth-first; recency-biased interpretation of newly-read fragments; concentration on hotspot files (frequently-edited or central modules); and a measurable cost to non-linear jumps. Foundation models prompted as code readers replicate the same traversal pattern, which is the bridge that lets the cognitive model guide LLM agent design.

## Relevance to MAgHARCM

This paper is the canonical hop-1 citation for [[PRIM-22]] (Four Phases of Comprehension). The Go implementation `internal/agents/comprehension.go` instantiates a six-step DR.JONES variant — Decomposition, Recognition, Organization, Navigation, Explanation, Search — that maps onto the canonical four-phase cognitive model: Search ⊃ Decomposition ∪ Navigation ∪ SearchAnchors; Anchoring ⊃ Recognition; Reasoning ⊃ Organization ∪ Explanation; Recertification ⊃ the per-iteration re-index the [[PRIM-26]] Symbol-Aware Navigator performs when an upstream translation changes the symbol surface. MAgHARCM's `Comprehend` traversal uses breadth-first linear ordering (depth-then-name) to match the predicted reader preference, and the resulting `SearchAnchors` slice is consumed by the downstream Translator and Reviewer agents as the entry-point surface for [[PRIM-31]] Iterative Retrieval Refinement. The recency-bias prediction is also the empirical reason [[PRIM-25]]'s RoleFlipReviewer works — recertification forces a re-read against the artifact rather than against prior context.

## Hop-1 References

- [[Soloway-Adelson-1984]] — Soloway & Adelson, "Behaviors of Expert Programmers: Episode vs. Schematic Knowledge", CSCW 1984 workshop / Yale technical report; the foundational distinction between episodic and schematic knowledge the DR.JONES Anchoring phase operationalises.
- [[Pennington-1987]] — Pennington, "Stimulus Structures and Mental Representations in Expert Comprehension of Computer Programs", Cognitive Science Vol. 11; the bottom-up program-model that informs the Reasoning phase.
- [[Brooks-1975-MMM]] — Brooks, "The Mythical Man-Month", Addison-Wesley 1975 (planning-fallacy chapter); cited as the canonical warning that comprehension effort is non-linear and subject to planning fallacy, which bounds DR.JONES' predicted budget.
- [[Détienne-1990]] — Détienne, "Expert Programming Knowledge", Lawrence Erlbaum 1990; the schema-theory precursor to Anchoring, on which DR.JONES' idiom-recognition step draws.
- [[Letovsky-1986]] — Letovsky, "Cognitive Processes in Program Comprehension", in Empirical Studies of Programmers (Soloway & Iyengar, eds.), Ablex 1986; the cognitive-process taxonomy that DR.JONES' four phases consolidate.
- [[Koenemann-Robertson-1991-ICPC]] — Koenemann & Robertson, "Expert Problem Solving Strategies for Program Comprehension", ICPC 1993 (precursor IWPC 1991 workshop); empirical study of comprehension strategies that anchors the breadth-first traversal claim.

## Hop-2 Anchors (software-archaeology lean)

- [[Parnas-1972]] — information-hiding module boundaries define the natural Recertification surface: when a module's secret changes, every anchored use of its interface must be re-anchored.
- [[Baldwin-Clark-2000]] — design-rule layering (L1 visible rules / L3 substitutable leaves) gives the Search phase a topological order: load-bearing rules are found before leaves.
- [[Rajlich-1997]] — concept assignment and concept-locator analysis is the offline analogue of Anchoring: locate the named concepts, then tie them to identifiers across the codebase.
- [[Müller-2000]] — strategy selection (Big Bang / Chicken Little / Pilot / Frozen / Parallel) determines when Recertification re-runs: each strategy re-anchors on a different cadence.
- [[Brooks-1986-NoSilverBullet]] — Brooks, "No Silver Bullet", IEEE Computer 1986; the essential-difficulty argument that bounds the Reasoning phase: not all comprehension gains are automatable.

## Backlinks

- [[METHODOLOGY]] §PRIM-22 — [[1.0.0 P-40 foltz-dr-jones-2023]] Foltz DR.JONES is the cognitive-model anchor for the Comprehension pipeline node in the Archaeologist agent's stage.
- [[primitives/INDEX]] — [[1.0.0 PRIM-22]] status entry names [[Foltz-2023]] DR. JONES Model as the lineage source.
- [[Software-Archaeology-Lineage]] §4 — [[1.0.0 P-40 foltz-dr-jones-2023]] Foltz DR.JONES is the modern cognitive-traversal citation that complements [[NEEDS-LINK Rajlich-1997]] concept-locator analysis in the two-hop archaeology synthesis.
