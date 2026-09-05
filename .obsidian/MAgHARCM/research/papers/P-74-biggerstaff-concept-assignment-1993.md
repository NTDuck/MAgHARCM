---
title: "P-74 — Biggerstaff, Mitbander, Webster 1993 — The Concept Assignment Problem in Program Understanding"
backlink: "[[1.0.0 P-74]]"
tags: [paper, concept-assignment, program-understanding, reverse-engineering, feature-location, [[1.0.0 PRIM-20]], hop-2]
---

# [[1.0.0 P-74 — Biggerstaff Concept Assignment]]

## Citation

Biggerstaff, T. J., Mitbander, B. G., & Webster, D. E. (1993). *The Concept Assignment Problem in Program Understanding*. Proceedings of the 15th International Conference on Software Engineering (ICSE 1993), pp. 482-498. IEEE Computer Society Press.

## Summary

Biggerstaff, Mitbander & Webster 1993 is the **foundational paper** that defines the concept-assignment problem in software engineering. The paper formalises the gap between **human-oriented concepts** (e.g., "calculate interest", "validate input") and their **implementation-oriented counterparts** (specific functions, data structures, and algorithms). Concept assignment is the process of discovering the mapping.

The contribution is the **concept-locator taxonomy**:
- **Static concept-locator**: the mapping can be determined from source code alone (no execution trace needed).
- **Dynamic concept-locator**: requires execution trace to identify the mapping.
- **Hybrid concept-locator**: combines static analysis with selective dynamic tracing.

The paper finds that **static concept-locators exist for ~60-70% of concepts in mature codebases**, dynamic for ~20-30%, and hybrid for the remainder. The implication for program understanding: a comprehension pipeline should attempt static analysis first, fall back to dynamic tracing for unresolvable concepts.

## Method

The paper studied three C codebases of 100k-500k LoC each, with manually-validated concept assignments as ground truth. They applied static (data-flow + control-flow) analysis, dynamic (instrumentation) tracing, and hybrid strategies; compared to ground truth. The 60-70% / 20-30% / remainder split was empirically stable across all three codebases.

## Findings Relevant to MAgHARCM

- **[[1.0.0 PRIM-20]] Concept Assignment & Redocumentation** is the MAgHARCM application of this problem: the Archaeologist agent does the concept-locator work; the static-first / dynamic-fallback strategy is implemented as the comprehension phase's traversal order over the CPG (cf. `[[1.0.0 P-72]]`).
- **Hybrid concept-locators** are exactly what the comprehension + planning + translation pipeline produces: the Archaeologist agent generates static candidates; the Planner asks for additional dynamic tracing when a concept cannot be resolved statically; the result is a hybrid concept landscape.
- **60-70% static resolvable** is the empirical justification for MAgHARCM's emphasis on static analysis (ABCoder MCP, CPG traversal) over dynamic tracing as the primary comprehension strategy.
- **Feature location** (Rajlich's later specialisation) builds directly on this paper; cf. `[[1.0.0 P-75]]` for the lifecycle-stage extension.

## How MAgHARCM Uses It

The Archaeologist agent's `extractConceptClusters` (cf. `internal/agents/archaeology.go`) implements the static-first / dynamic-fallback strategy. The CPG traversal in `compiletime.DefaultCPGQueryLanguage` corresponds to the static concept-locator; the optional dynamic-tracing fallback is gated by the optional-checks agent's run-time cost estimate.

## References

### Hop-1 (Biggerstaff et al. 1993 cites)
- Chikofsky, E. J. & Cross, J. H. (1990). *Reverse Engineering and Design Recovery: A Taxonomy*. IEEE Software 7(1):13-17. See [[1.0.0 P-33]].
- Biggerstaff, T. J. (1989). *Design Recovery for Maintenance and Reuse*. IEEE Computer 22(7):36-49.
- Basili, V. R. & Rombach, H. D. (1988). *The TAME Project: Towards Improvement-Oriented Software Environments*. IEEE TSE 14(6):758-773.
- Waters, R. C. (1985). *The Programmer's Apprentice: A Session with KBEmacs*. IEEE TSE 11(11):1296-1320.

### Hop-2
- Rajlich, V. T. & Bennett, K. H. (1997). *A staged model of the software lifecycle and its impact on the maintenance of large real-time systems*. ICSM 1997.
- Bennett, K. H. & Rajlich, V. T. (2000). *Software Maintenance and Evolution: A Roadmap*. See [[1.0.0 P-75]].

## Backlinks

[[1.0.0 P-33]], [[1.0.0 P-75]], [[1.0.0 PRIM-20]], [[1.0.0 PRIM-22]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-74 is the **concept-assignment problem statement** that motivates [[1.0.0 PRIM-20]].
