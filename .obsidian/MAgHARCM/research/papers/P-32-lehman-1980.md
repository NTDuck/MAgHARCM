---
title: Programs, Life Cycles, and Laws of Software Evolution
bibkey: p32_lehman
tags: [paper, software-evolution, lehman-laws, [[1.0.0 PRIM-14]], [[1.0.0 PRIM-18]], hop-1]
---

# [[1.0.0 P-32]] Programs, Life Cycles, and Laws of Software Evolution

**Authors**: Meir M. Lehman  
**Year**: 1980 (expanded 1996)  
**Venue**: Proceedings of the IEEE, Vol. 68, No. 9, pp. 1060–1076  
**eprint / DOI**: 10.1109/PROC.1980.11805  
**Cited by**: `[[2.0.0 Software-Archaeology-Lineage]]`, `[[1.0.0 PRIM-14]]` (Software Archaeology Stage), `[[1.0.0 PRIM-18]]` (Jaccard-Coupling Architecture Recovery)

## Summary

Lehman established the foundational empirical theory of software evolution by classifying software into three categories:
- **S-Type (Specified)**: Systems completely defined by a formal mathematical specification (e.g. matrix multiplication).
- **P-Type (Problem)**: Systems addressing practical real-world problems through approximations.
- **E-Type (Evolutionary)**: Systems embedded in real-world environments where the system becomes part of the world it models (e.g., enterprise ERPs, operating systems, financial transaction engines).

For E-Type systems, Lehman formulated the canonical **Laws of Software Evolution**:
1. *Continuing Change*: An E-type system must be continually adapted or it becomes progressively less satisfactory.
2. *Increasing Complexity*: As an E-type system evolves, its complexity increases unless work is done to maintain or reduce it.
3. *Self-Regulation*: Evolution processes are self-regulating with statistically determinable feedback distributions.
4. *Conservation of Familiarity*: Incremental system growth across releases remains bounded to prevent cognitive overload.

## Relevance to MAgHARCM

Lehman's laws explain why naive transpilation of legacy software fails:
1. Legacy systems are quintessential E-type systems whose accumulated complexity is emergent rather than strictly planned.
2. The Law of Increasing Complexity explains why static documentation diverges from real code behaviour over decades, necessitating empirical excavation (`[[1.0.0 PRIM-14]]`).
3. The Law of Continuing Change justifies `[[1.0.0 PRIM-18]]` (Jaccard-Coupling Recovery): temporal coupling across commit history reveals true architectural dependencies that outlive static syntax.

## Hop-1 References

- [[Chikofsky-Cross-1990]] — Reverse engineering is the methodology to counter Lehman's Law of Increasing Complexity.
- [[Rajlich-Bennett-2000]] — Expands Lehman's evolution phase into the five-stage software lifecycle model (Initial, Evolution, Servicing, Phase-out, Close-down).
- [[MSR4SA-2017]] — MSR leverages Lehman's feedback loops to discover hidden co-change coupling from commit logs.

## Hop-2 Anchors (Software Archaeology Lean)

- [[pp-besm-Software-Archaeology]] — Geological strata in code repositories are direct physical manifestations of Lehman's successive evolutionary releases.
- [[Kazman-Cai-2017]] — Architectural debt in DRSpaces measures the entropy accumulated through unmitigated Lehman evolution.
