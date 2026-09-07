---
title: Seacord, Comella-Dorda, Lewis, Place & Plakosh — Modernizing Legacy Systems (SEI / Addison-Wesley 2003)
bibkey: p46_seacord_modernizing_2003
tags: [paper, legacy-modernization, risk-driven, architecture-recovery, [[1.0.0 PRIM-14]], hop-1]
---

# [[1.0.0 P-46]] Seacord, Comella-Dorda, Lewis, Place & Plakosh — Modernizing Legacy Systems

**Authors**: Robert C. Seacord, Santiago Comella-Dorda, Grace A. Lewis, Pat R. Place, Daniel Plakosh (Carnegie Mellon, Software Engineering Institute)
**Year**: 2003 (book); 2001 (SEI technical report CMU/SEI-2001-TR-025)
**Venue**: *Modernizing Legacy Systems: Software Technologies, Engineering Processes, and Business Practices*, Addison-Wesley Professional; SEI Technical Report CMU/SEI-2001-TR-025
**eprint / DOI**: ISBN 0-201-77045-9 (Addison-Wesley); DTIC ADA396063 (SEI TR)
**Cited by**: [[primitives/INDEX]] entry [[PRIM-14]] (Software-Archaeology Stage) — provides the industry-grade, risk-driven modernization playbook that the archaeological excavation pass executes.

## Summary

[[Seacord-2003-Modernizing]] is the canonical SEI reference for modernizing long-lived legacy systems. Drawing on dozens of SEI modernization engagements (air-traffic control, manufacturing execution, insurance claims, payroll), the authors argue that the central failure mode of legacy modernization is not technical, but *risk-management*: teams take a Big-Bang rip-and-replace stance, lose familiarity with the system's hidden invariants, and ship a modernized system that loses the original's tacit guarantees (audit trails, accounting rounding, regulatory reporting, peak-load behaviour).

The book's central framework is a four-stage risk-driven modernization process:

1. *Inventory and characterization* — every legacy component is inventoried and characterized by age, volatility, defect density, business criticality, and modernization cost. This is the analogue of the five archaeological strata in [[pp-besm]] (`[[1.0.0 P-36]]`).
2. *Architecture recovery* — the legacy architecture is reconstructed using a hybrid of source-code analysis, execution tracing, and developer interviews. The authors explicitly describe architecture recovery as "software archaeology" — borrowing Chikofsky & Cross's design-recovery taxonomy (`[[1.0.0 P-33]]`).
3. *Modernization strategy selection* — choose from a portfolio of incremental strategies (Seacord et al. list *encapsulate*, *rehost*, *replatform*, *refactor*, *rebuild*, *replace*) based on the risk/cost profile from stage 1. The list extends Müller's earlier five-strategy taxonomy (`[[1.0.0 P-35]]`).
4. *Incremental migration with continuous validation* — modernized slices are deployed alongside the legacy system (parallel cutover), with characterization tests (Feathers 2004, `[[1.0.0 P-25]]`) acting as the regression oracle. The legacy system is the *specification*; the modernized system must match its behaviour on the captured characterization tests.

The book's strongest practical contribution is its insistence that characterization tests be built *before* any code is touched. Modernization without characterization tests is blind translation; modernization with them is constrained re-implementation. This is the canonical industry framing behind MAgHARCM's combination of [[PRIM-5]] (Test Suite Co-Translation), [[PRIM-8]] (Mock-Based Validation), and [[PRIM-13]] (Adversarial Test-Weakening Guard).

The book also introduces the *legacy system screener* — a structured questionnaire that scores legacy components on twelve dimensions (size, complexity, volatility, maintainability, business value, dependency density, test coverage, etc.) and uses the score to triage which components to wrap, refactor, or rewrite first. This screener is the practical ancestor of MAgHARCM's [[PRIM-14]] triage layer.

## Relevance to MAgHARCM

1. [[PRIM-14]] (Software-Archaeology Stage, `internal/agents/archaeology.go`): the archaeological excavation pass is a direct realization of Seacord's stage 1 + stage 2. The screener scores correspond to the per-module archaeological metadata that the archaeologist agent emits; the architecture-recovery step corresponds to the design-rule and concept-assignment extraction.
2. [[PRIM-21]] (Migration Strategy Selection, `internal/agents/strategy.go::Registry.TryInOrder`): Seacord et al.'s six-strategy portfolio (encapsulate / rehost / replatform / refactor / rebuild / replace) extends Müller's five strategies (`[[1.0.0 P-35]]`); [[PRIM-21]]'s try-and-fail registry is the policy layer that walks through this portfolio.
3. [[PRIM-5]] (Test Suite Co-Translation & Synthesis) and [[PRIM-13]] (Adversarial Test-Weakening Guard): the book's "characterize before you change" mandate is the canonical industrial framing behind both primitives. Characterization tests are the bridge between the legacy system (specification) and the modernized system (implementation).
4. [[PRIM-25]] (Role-Flip De-Hallucination Gate): the book's insistence on developer-interview triangulation is the human-validated analogue of the role-flip gate; the book's "skeptical review" practice maps onto the role-flip pattern.

## Hop-1 References

- [[Chikofsky-Cross-1990]] — the reverse-engineering taxonomy that the book explicitly imports and extends in its stage-2 architecture-recovery section.
- [[Muller-2000-Strategies]] — Müller's five migration strategies; the book extends the list with encapsulate/rehost/replatform/refactor/rebuild/replace.
- [[Feathers-2004-Legacy-Code]] — "legacy code is code without tests"; the characterization-test mandate the book advocates for is the central Feathers discipline.
- [[Bennett-2000-Staged-Life-Cycle]] — the staged lifecycle that motivates Seacord's stage-1 inventory step.
- [[Rajlich-1997-Concept-Assignment]] — concept assignment; one of the three techniques Seacord lists under stage-2 architecture recovery (alongside DSM clustering and scenario reconstruction).
- [[Kazman-Cai-2017-DRSpaces]] — design-rule-space extraction; the modern realization of the architecture-recovery stage Seacord describes.

## Hop-2 Anchors (software-archaeology lean)

- [[Lehman-1980-Laws-Evolution]] — Seacord's risk-driven modernization exists precisely because of Lehman's "increasing complexity" law: without continuous work, legacy systems decay into unmaintainable form.
- [[Baldwin-Clark-2000]] — the encapsulate / refactor / rebuild / replace strategy is essentially a Baldwin-Clark modular-operator sequence: substituting → augmenting → inverting → porting, applied to legacy subsystems.
- [[MacCormack-2006-DSM]] (`[[1.0.0 P-42]]`) — the DSM cycle-density scoring that motivates Seacord's volatility score on the legacy-system screener.
- [[Curtis-1992-Process-Modeling]] (`[[1.0.0 P-45]]`) — Seacord's four-stage process is itself a software-process model; Curtis's multi-aspect framework explains why the four stages must be enacted (not merely documented) to work.
- [[Chikofsky-Cross-1990]] (`[[1.0.0 P-33]]`) — design recovery; Seacord's stage-2 architecture-recovery pass is one of the canonical applications of Chikofsky & Cross's taxonomy.

## Backlinks

- [[METHODOLOGY]] §PRIM-14 — Seacord et al. (2003) is the industry-grade academic anchor for the Software-Archaeology Stage; complements [[pp-besm]] (`[[1.0.0 P-36]]`) and [[Chikofsky-Cross-1990]] (`[[1.0.0 P-33]]`) in the excavation-pass lineage.
- [[primitives/INDEX]] — [[1.0.0 PRIM-14]] status entry now references [[P-46]] alongside pp-besm and Chikofsky & Cross (1990).
- [[Software-Archaeology-Lineage]] §4 — the risk-driven modernization playbook is the overarching methodology that PRIM-14's archaeological excavation pass executes.
