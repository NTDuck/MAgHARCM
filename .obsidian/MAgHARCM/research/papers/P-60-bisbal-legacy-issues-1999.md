---
title: "P-60 — Bisbal, Lawless, Wu & Grimson 1999 — Legacy Information Systems: Issues and Directions"
backlink: "[[1.0.0 P-60]]"
tags: [paper, legacy-information-systems, big-bang-vs-incremental, migration-strategies, cascading-fallback, [[1.0.0 PRIM-21]], hop-1]
---

# [[1.0.0 P-60 — Bisbal et al. — Legacy Information Systems]]

## Citation

Bisbal, J., Lawless, D., Wu, B., & Grimson, J. (1999). *Legacy Information Systems: Issues and Directions*. IEEE Software 16(5):103–111. DOI 10.1109/52.795108. Trinity College Dublin — Knowledge & Data Engineering Group, MILESTONE project.

## Summary

Bisbal et al. is the canonical 1990s survey of legacy information system (LIS) modernization from Trinity College Dublin's MILESTONE group. The paper defines a *Legacy Information System* not by age but by resistance to modification — "any information system that significantly resists modification and evolution" — and argues that the central failure mode of modernization projects is *risk blindness*: teams commit to a single migration strategy (most often "cold turkey") whose preconditions do not actually hold for the system in front of them, and ship a modernized system that fails precisely on the business invariants the legacy had tacitly preserved.

The paper's central contribution is a **risk-stratified taxonomy** of migration strategies organized along two axes — *extent of change* (conservative → radical) and *disruption tolerance* (high → low). Strategies covered:

1. **Big Bang (Cold Turkey)** — wholesale redevelopment. Flagged as the most common strategy in industry surveys and simultaneously the highest-risk; "typically too high for this method to be seriously contemplated for large-scale, mission-critical systems."
2. **Wrapping** — placing a modern interface (CORBA / EJB / SOA façade) in front of legacy components so they remain callable while internals are replaced piecewise.
3. **Incremental (Wrap-and-replace, pilot, strangler)** — small, well-defined subsets are migrated in priority order, with each migration validated against the live legacy before proceeding to the next.
4. **Migration + parallel operation** — keep the legacy running alongside the modernized system, route a percentage of traffic to the modernized version, watch for divergence, ramp up.

## Findings Relevant to MAgHARCM

- **Cascading-fallback as emergent property.** The paper documents that cascading-fallback (try-and-fail through a list of strategies) is not a clean theoretical pattern but an *emergent* property of real projects: teams that committed to a single strategy up-front routinely fell back to a less ambitious strategy when the original choice's preconditions turned out not to hold. The paper does not propose cascading-fallback as a design — it observes it as a property of failed projects.
- **Wrapping / incremental as the universal floor.** Across all the case studies the paper surveys, the *only* strategy that reliably succeeded was one whose preconditions were minimal — the equivalent of MAgHARCM's `incrementalStrategy` with `Matches() == true`. Every other strategy had a real-world preconditions-mismatch failure rate.
- **Risk stratification beats plan-purity.** Projects that explicitly stratified risk across strategies (e.g., wrap-and-replace for low-risk components, big-bang for low-coupling new modules) outperformed projects that committed to a single strategy across the whole codebase.

## How MAgHARCM Uses It

[[1.0.0 PRIM-21]] (Migration Strategy Selection) is implemented in `internal/agents/strategy.go` as `Registry.TryInOrder`: the registry iterates its five Mueller strategies in canonical priority order, returning the first whose `Matches(Profile)` accepts the repository's signal set. Bisbal et al. supplies the **empirical risk-driven justification** for the cascading execution order — Müller's P-35 supplies the conceptual taxonomy; Bisbal's P-60 supplies the empirical evidence that the order matters and that the cascade-with-universal-floor is the empirically best-performing design.

## References

### Hop-1
- Müller, R., Tzerpos, V., et al. (2000). *Five Strategies for Re-engineering*. See [[1.0.0 P-35]].
- Sneed, H. (1995). *Planning the Reengineering of Legacy Systems*. IEEE Software 12(1).
- Brodie, M. & Stonebraker, M. (1995). *Migrating Legacy Systems*. Morgan Kaufmann.
- Seacord, R., Plakosh, D., Lewis, G. (2003). *Modernizing Legacy Systems*. Addison-Wesley. See [[1.0.0 P-46]].
- Fowler, M. (2004). *StranglerFigApplication*. (cited inline in MAgHARCM as `[Strangler-2004]`)

### Hop-2
- Lehman, M. M. (1980). See [[1.0.0 P-32]].
- Parnas, D. L. (1972). See [[1.0.0 P-31]].
- Baldwin, C. Y. & Clark, K. B. (2000). See [[1.0.0 P-41]].

## Backlinks

[[1.0.0 PRIM-21]], [[1.0.0 P-35]], [[1.0.0 P-46]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].
