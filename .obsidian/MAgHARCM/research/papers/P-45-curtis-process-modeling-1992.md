---
title: Curtis, Kellner & Over — Process Modeling (CACM 1992)
bibkey: p45_curtis_process_modeling_1992
tags: [paper, software-process, sop, workflow, [[1.0.0 PRIM-24]], hop-1]
---

# [[1.0.0 P-45]] Curtis, Kellner & Over — Process Modeling (CACM 1992)

**Authors**: Bill Curtis, Marc I. Kellner, James Over (Hewlett-Packard)
**Year**: 1992 (September)
**Venue**: *Communications of the ACM*, Vol. 35, No. 9, pp. 75-90
**eprint / DOI**: DOI 10.1145/135239.135249 (CACM September 1992)
**Cited by**: [[primitives/INDEX]] entry [[PRIM-24]] (SOP-Anchored Role-Artifact Schema) — the canonical process-modelling hop-1 reference for SOP contracts.

## Summary

[[Curtis-1992-Process-Modeling]] surveys the state of software-process modelling in 1992, organising the field around four orthogonal dimensions:

1. *Process model notation* — petri-nets, state-transition graphs, rule-based systems, data-flow diagrams; the trade-off between expressiveness and analysability.
2. *Enactment mechanism* — process models are either enacted by a process engine (active) or merely documented (passive).
3. *Process model content* — control-flow only (what runs next) versus multi-aspect (control-flow + data-flow + role-assignment + artefact structure).
4. *Use* — descriptive (recording what happened), prescriptive (driving what should happen), or retrospective (explaining why it happened).

The paper's central empirical claim is the "process maturity" finding: organisations with explicit, multi-aspect, enacted process models (the upper-right quadrant of the matrix) exhibit systematically lower defect rates than organisations with informal, single-aspect, descriptive processes. This claim prefigures the later CMM/CMMI evidence base by seven years.

The paper also introduces the **Role-Activity Diagram (RAD)** notation as an example of a multi-aspect process model, in which role assignment is a first-class modelling primitive alongside control flow and artefact structure. RAD is the direct ancestor of MetaGPT's SOP contracts and of MAgHARCM's [[PRIM-24]] role-artifact schema.

## Relevance to MAgHARCM

1. [[PRIM-24]] (SOP-Anchored Role-Artifact Schema, `internal/compiletime/compiletime.go`): the SOP definition language is a modern descendant of RAD. Each SOP node declares the role, the artefact type, the input/output schema, and the transition condition — exactly the four-dimensional decomposition Curtis, Kellner & Over prescribe.
2. [[P-11]] (MetaGPT-2024 SOP contracts): MetaGPT cites Curtis (1992) as one of its antecedents for SOP-style role decomposition. The [[P-45]] entry provides the academic anchor.
3. [[P-12]] (ChatDev-2024 chat chains): ChatDev's role-flipping is a process-modelling concern — the SOP defines when a role may be flipped, and Curtis's multi-aspect process modelling is the theoretical backdrop.
4. [[P-26]] (AgentVerse-2023, dynamic recruitment): AgentVerse's post-iteration evaluation is an enactment-mechanism concern; Curtis's prescriptive-vs-descriptive distinction is the conceptual axis.

## Hop-1 References

- [[Humphrey-1989-Managing-Software-Process]] — the CMM precursor; the process-maturity claim prefigured by the SEI Capability Maturity Model.
- [[Kellner-1990-Process-Modelling-Issues]] — the predecessor to the CACM survey; the RAD notation is introduced here.
- [[Osterweil-1987-Software-Processes-Are-Software-Too]] — the canonical claim that process models are themselves software artefacts and should be treated with the same engineering rigour.
- [[Zave-1989-Operational-Approach-Requirements]] — the "operational approach" to software specification; process modelling as specification.
- [[Jackson-1995-Software-Requirements-And-Specifications]] — requirements engineering as a process-modelling problem; the broader theoretical context.
- [[Wing-1990-Family-Of-Requirements]] — the family-of-requirements formalism; multi-aspect process modelling's requirements-side cousin.

## Hop-2 Anchors (software-archaeology lean)

- [[Muller-2000-Strategies]] — Müller's five migration strategies are themselves process models; [[Curtis-1992-Process-Modeling]] is the meta-theory they sit inside.
- [[Rajlich-1997-Staged-Life-Cycle]] — the staged lifecycle is a process model; Curtis's multi-aspect framework explains why staged lifecycles need both control-flow and role-assignment.
- [[Lehman-1980-Laws-Evolution]] — Lehman's "feedback-driven evolution" is a process-modelling insight: software processes are second-order systems that adapt to the artefacts they produce.
- [[Baldwin-Clark-2000]] — design-rule decomposition maps to process modelling: L1 design rules are the artefacts a process produces; the process itself is the L2 subsystem.
- [[Parnas-1972]] — information hiding is a process-modelling concern: the process must preserve module interfaces across stages, which is a control-flow invariant.

## Backlinks

- [[METHODOLOGY]] §PRIM-24 — Curtis, Kellner & Over (1992) is the academic anchor for the SOP-Anchored Role-Artifact Schema; complements [[P-11]] (MetaGPT) for the modern implementation lineage.
- [[primitives/INDEX]] — [[1.0.0 PRIM-24]] status entry now references [[P-45]] alongside [[P-11]] (MetaGPT) and Meyer (1988) DbC.
- [[Software-Archaeology-Lineage]] §4 — process modelling is the orchestration theory that PRIM-24's SOP contracts implement.
