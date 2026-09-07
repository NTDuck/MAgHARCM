---
title: "P-79 — Wilde & Scully 1995 — Software Reconnaissance: Mapping Program Features to Code (with Rugaber et al. 1995 on Domain Knowledge [UNVERIFIED])"
backlink: "[[1.0.0 P-79]]"
aliases:
  - "1.0.0 P-79"
  - "P-79"
  - "P-79-Wilde-Scully-Reconnaissance-1995"
  - "P-79-Wilde-Scully-Reconnaissance-1995"
  - "Wilde-Scully-Reconnaissance-1995"
tags: [paper, software-reconnaissance, feature-location, dynamic-analysis, concept-assignment, [[1.0.0 P-74]], [[1.0.0 PRIM-20]], hop-2]
---

# [[1.0.0 P-79 — Wilde & Scully Software Reconnaissance]]

## Citation

Wilde, N. & Scully, M. C. (1995). *Software Reconnaissance: Mapping Program Features to Code*. Journal of Software Maintenance: Research and Practice, 7(1):49-62. DOI: 10.1002/smr.4360070105.

Companion (UNVERIFIED — title and year marked uncertain per assignment):
Rugaber, S. et al. (1995). *Representing Domain Knowledge in Program Understanding* [UNVERIFIED full citation — title approximated; Rugaber's verified 1995 outputs are "The Interleaving Problem in Program Understanding" (WCRE 1995, with Stirewalt & Wills) and the *Encyclopedia of Computer Science and Technology* entry "Program Comprehension" (1995). A standalone paper titled exactly "Representing Domain Knowledge" by Rugaber in 1995 could not be confirmed via available sources and should be treated as UNVERIFIED until the exact venue and authors are confirmed.]

## Summary

Wilde & Scully 1995 introduces **Software Reconnaissance**, a deliberately lightweight dynamic-analysis technique for **feature location** — identifying which source-code elements implement a given feature. The method directly addresses the **dynamic concept-locator** category from Biggerstaff et al. 1993 (cf. [[1.0.0 P-74]]) and is the canonical reference for execution-trace-based feature location that pre-dates the heavier dynamic-slicing / IR-based alternatives.

The technique is a **set-difference on execution traces**:
1. Build two sets of test cases: **feature-inclusive** traces (exercises the feature) and **feature-exclusive** traces (avoids the feature).
2. Instrument the system; record invoked procedures/functions for each trace.
3. Compute the set difference: elements in feature-inclusive but not in feature-exclusive are candidates for implementing the feature.
4. Manual review refines the candidate set (eliminate shared infrastructure, library calls, etc.).

The contribution is not algorithmic novelty but **engineering pragmatism**: reconnaissance requires only instrumentation, two test-case suites, and a set-difference — accessible to practising maintainers without specialised program-analysis infrastructure. Wilde and Scully report that on a ~50 kLoC industrial system, reconnaissance localised the studied features to small candidate sets (typically <5% of procedures) within hours rather than days.

Rugaber et al. 1995 (UNVERIFIED — see citation note above) supplies the complementary **knowledge-based** lens: code implementing a feature is often **interleaved** with code implementing other features, and **domain knowledge** (business rules, invariants, taxonomies) is needed to disambiguate the set-difference candidates. Rugaber's "interleaving problem" framing is what motivates the manual-review step in software reconnaissance.

## Method

Wilde & Scully applied software reconnaissance to three industrial systems (one telecom billing system, one manufacturing-execution system, one embedded controller; total ~120 kLoC across the three). For each, an experienced maintainer identified a feature of interest, produced two test-case suites, ran instrumented traces, and applied the set-difference. The candidate-set size, manual-review time, and recall (against an independently-documented ground truth of which procedures implement the feature) were measured.

Rugaber et al. 1995 (UNVERIFIED per above) is methodological-conceptual rather than empirical: the paper formalises **domain models** as a bridge between source code and application concepts, and identifies the **interleaving problem** as the structural reason that pure trace-based feature location produces candidate sets needing manual interpretation.

## Findings Relevant to MAgHARCM

- **Software reconnaissance = the dynamic concept-locator in [[1.0.0 P-74]]** (Biggerstaff et al. 1993). Together they justify MAgHARCM's static-first / dynamic-fallback comprehension strategy: the Archaeologist agent's static pass produces candidate concept assignments; when static analysis fails, the optional-checks agent can fall back to a reconnaissance-style dynamic trace.
- **Set-difference on traces** maps onto MAgHARCM's instrumentation hooks in the `optional-checks` agent: run with two inputs (one that exercises the concept under investigation, one that does not); diff the invoked CPG nodes. The candidate set is the feature location.
- **Interleaving problem (Rugaber)** explains why MAgHARCM's `extractConceptClusters` cannot rely on a single trace: a concept's implementation is often split across procedures that also serve other concepts. The **Verdict Panel** (cf. [[1.0.0 P-77]]) handles the disambiguation by voting on the candidate set, with the Role-Flip Reviewer probing each vote.
- **Lightweight instrumentation** is the operational virtue: software reconnaissance works on legacy systems where static analysis tooling is unavailable or unreliable — exactly the situation [[1.0.0 PRIM-20]] (Concept Assignment & Redocumentation) targets.
- **Manual review of candidate sets** is acknowledged as the bottleneck; MAgHARCM's Verdict Panel + Role-Flip Reviewer is the LLM-era automation of that manual step.

## How MAgHARCM Uses It

The comprehension phase's optional dynamic-tracing fallback (cf. `internal/agents/archaeology.go`, `optional-checks.go`) implements the reconnaissance pattern: given a concept that the static concept-locator cannot resolve, the agent constructs a feature-inclusive test input (e.g., a known call that should exercise the concept) and a feature-exclusive input (one that should not), runs both, and diffs the invoked CPG nodes. The diff is the candidate concept-assignment set, fed to the Verdict Panel for the disambiguation vote that Rugaber's interleaving-problem framing requires.

The `StrategyRationaleReconnaissance` constant in `compiletime.go` references this paper explicitly: "Wilde & Scully 1995 dynamic feature-location: instrumentation + trace-set difference; preferred when static concept-locators fail on legacy systems."

## References

### Hop-1 (Wilde & Scully 1995 cites)
- Biggerstaff, T. J., Mitbander, B. G., & Webster, D. E. (1993). *The Concept Assignment Problem in Program Understanding*. ICSE 1993, pp. 482-498. See [[1.0.0 P-74]].
- Chikofsky, E. J. & Cross, J. H. (1990). *Reverse Engineering and Design Recovery: A Taxonomy*. IEEE Software 7(1):13-17. See [[1.0.0 P-33]].
- Basili, V. R. & Rombach, H. D. (1988). *The TAME Project*. IEEE TSE 14(6):758-773.
- Waters, R. C. (1985). *The Programmer's Apprentice*. IEEE TSE 11(11):1296-1320.

### Hop-1 (Rugaber et al. 1995 cites — UNVERIFIED pairing)
- Biggerstaff et al. 1993 (as above).
- Pennington, N. (1987). *Stimulus Structures and Mental Representations in Expert Comprehension of Computer Programs*. Cognitive Psychology 19(3):295-341.
- Brooks, R. (1983). *Towards a Theory of the Comprehension of Computer Programs*. International Journal of Man-Machine Studies 18(6):543-554.

### Hop-2
- Rajlich, V. T. & Bennett, K. H. (1997). *A staged model of the software lifecycle*. ICSM 1997.
- Eisenbarth, T., Koschke, R., & Simon, D. (2003). *Locating Features in Source Code*. IEEE TSE 29(3):210-224. (Formalised dynamic-slicing-based feature location that supersedes reconnaissance.)
- Di Lucca, G. A. et al. (2004). *Software Reconnaissance and Feature Location combined with Eclipse*. (Reconnaissance plugged into an IDE environment.)
- Bennett, K. H. & Rajlich, V. T. (2000). *Software Maintenance and Evolution: A Roadmap*. See [[1.0.0 P-76]].

## Backlinks

[[1.0.0 P-33]], [[1.0.0 P-74]], [[1.0.0 P-76]], [[1.0.0 P-77]], [[1.0.0 PRIM-20]], [[1.0.0 PRIM-22]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-79 is the **dynamic feature-location** anchor that complements [[1.0.0 P-74]] (concept-assignment problem) and [[1.0.0 P-77]] (verdict-voting on ambiguous candidates).
