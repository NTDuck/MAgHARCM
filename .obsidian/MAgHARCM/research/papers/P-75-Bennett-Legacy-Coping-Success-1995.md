---
title: "P-75 — Bennett 1995 — Legacy Systems: Coping with Success"
backlink: "[[1.0.0 P-75]]"
aliases:
  - "1.0.0 P-75"
  - "P-75"
  - "P-75-Bennett-Legacy-Coping-Success-1995"
  - "P-75-Bennett-Legacy-Coping-Success-1995"
  - "Bennett-Legacy-Coping-Success-1995"
tags: [paper, legacy-systems, modernisation, strategy, [[1.0.0 P-74]], [[1.0.0 P-46]], hop-2]
---

# [[1.0.0 P-75 — Bennett Legacy Systems]]

## Citation

Bennett, K. H. (1995). *Legacy Systems: Coping with Success*. IEEE Software 12(1):19-23. DOI: 10.1109/52.368257.

## Summary

Bennett's "Legacy Systems: Coping with Success" is the canonical 2-page framing paper for legacy-system modernisation. The paper opens with the now-famous definition:

> "Large software systems that we don't know how to cope with but that are vital to our organisation."

The "Coping with Success" title highlights the paradox: a system becomes "legacy" **because** it has been successful — it has survived long enough, become deeply embedded in business processes, and grown to a size where no single person understands it fully. The paper's contribution is the **modernisation strategy taxonomy**:
- **Leave alone** (servicing-stage maintenance only).
- **Reengineer** (architectural redesign + phased migration).
- **Replace** (discard legacy and rebuild from scratch).
- **Wrap** (encapsulate behind a service interface, migrate callers incrementally).
- **Migrate** (lift-and-shift or platform-port).

The paper explicitly compares the strategies' costs and risks; **wrap** is the recommended default for large systems; **reengineer** is preferred when the legacy's structure is still sound but the technology is obsolete; **replace** is generally high-risk.

## Method

Bennett's method is **conceptual synthesis** of empirical observations from the software-engineering literature up to 1994 (Brooks, Belady-Lehman, Parnas, Ward, etc.) plus case-study reflections on real industrial maintenance projects. The paper is short (5 pages) and serves as the **common vocabulary** for legacy modernisation research that followed.

## Findings Relevant to MAgHARCM

- **Wrap strategy = Strangler Fig** (cf. `[[1.0.0 PRIM-21]]` Migration Strategy Selection): Bennett's "wrap" is the conceptual predecessor of the Strangler-Fig pattern that MAgHARCM recommends as the default for medium-large legacy systems.
- **Leave alone = Servicing-stage** (cf. `[[1.0.0 P-32]]` Lehman laws): for codebases that are still in maintenance, no translation is the right answer.
- **Replace = Big-Bang**: high-risk; reserved for small projects (<500 LoC, cf. `compiletime.BigBangLoCMax`).
- **Reengineer = In-Place Adaptation**: preferred when the architecture is sound but the technology is obsolete.
- **Migrate = the umbrella term**: any of the above depending on the strategy registry's `TryInOrder` outcome.

## How MAgHARCM Uses It

The `compiletime.StrategyRegistry` (cf. `internal/agents/strategy.go`) implements the five-strategy taxonomy from Bennett 1995 as the `TryInOrder` policy: each strategy is tried in sequence, the first to produce a verified translation wins. The `StrategyRationale*` constants in `compiletime.go` reference the Bennett 1995 framing (e.g., "Strangler Fig (Bennett 1995 wrap): incremental substitution of legacy modules behind a stable interface").

## References

### Hop-1 (Bennett 1995 cites)
- Brooks, F. P. (1975). *The Mythical Man-Month*. Addison-Wesley.
- Belady, L. A. & Lehman, M. M. (1976). *A Model of Large Program Development*. IBM Systems Journal 15(3):225-252. See [[1.0.0 P-32]].
- Parnas, D. L. (1972). *On the Criteria to Be Used in Decomposing Systems into Modules*. CACM 15(12):1053-1058. See [[1.0.0 P-31]].
- Ward, W. T. (1990). *A Quantitative Study of Software Maintenance*. PhD thesis, University of Liverpool.

### Hop-2
- Seacord, R. C. et al. (2003). *Modernizing Legacy Systems: Software Technologies, Engineering Processes, and Business Practices*. Addison-Wesley. See [[1.0.0 P-46]].
- Bennett, K. H. & Rajlich, V. T. (2000). *Software Maintenance and Evolution: A Roadmap*.

## Backlinks

[[1.0.0 P-32]], [[1.0.0 P-46]], [[1.0.0 P-74]], [[1.0.0 PRIM-21]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-75 is the **modernisation-strategy taxonomy** anchor for [[1.0.0 PRIM-21]].
