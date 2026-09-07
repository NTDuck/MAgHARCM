---
title: "P-103 — AgentModernize (Ahmed & Galib 2026) Behavioral Specification Graphs"
backlink: "[[1.0.0 P-103]]"
tags: [paper, legacy-modernization, multi-agent, behavioral-preservation, arxiv-2026, [[1.0.0 PRIM-14]], [[1.0.0 PRIM-15]], hop-1]
---

# [[1.0.0 P-103 — AgentModernize: Behavioral Specification Graphs for Legacy Modernization]]

## Citation

Ahmed, S. N. & Galib, M. (2026). *AgentModernize: Preserving Business Logic in Legacy Modernization with Multi-Agent LLMs and Behavioral Specification Graphs*. arXiv:2605.17535, v1 May 17 2026; v2 Aug 4 2026.
URL: https://arxiv.org/abs/2605.17535

**Venue note**: This is an arXiv preprint, NOT an ICSE publication. Earlier web-search summaries that attributed ICSE 2025 to this paper were wrong; the canonical venue is the arXiv listing.

## Summary

AgentModernize reframes legacy modernization as a **behavioural preservation** problem rather than a syntax translation problem. Four specialised agents (Legacy Analyzer, Specification Generator, Modernization Transformer, Equivalence Validator) coordinate through an intermediate artifact called a **Behavioural Specification Graph (BSG)** — a typed intermediate representation featuring preconditions, postconditions, and invariants that captures business logic semantics across the modernisation boundary.

Key contributions:
1. **BSG intermediate**: structurally similar to MAgHARCM's `compiletime.State` but lighter (only behavioural invariants, not full AST). Acts as a "trust boundary" that makes business logic inspectable before code is generated.
2. **Multi-agent split**: validates MAgHARCM's 8-agent decomposition with empirical evidence that specialisation beats monolithic translation (v1 results; v2 nuances: for stronger model backbones, simpler approaches can outperform the multi-agent pipeline).
3. **Empirical COBOL→Java baseline**: demonstrates preserved business logic through differential testing rather than syntactic round-trip.

## Relevance to MAgHARCM

- **Behavioural Specification Graph concept validates PRIM-15 (Evidence-First Adaptation)**: BSG is the concrete instantiation of evidence-first translation at the legacy/target boundary. MAgHARCM's translator should expose a BSG-equivalent trace for verifier cross-check.
- **Multi-agent decomposition aligns with METHODOLOGY §1**: AgentModernize provides independent arXiv 2026 empirical evidence for MAgHARCM's 8-agent pipeline.
- **Differential testing as validation mode** complements PRIM-6 (Multi-Stage Build/Test Feedback Repair): the validator's differential check is the test that the BSG contract holds across iterations.
- **Software archaeology tie-in**: BSG extraction reads directly into PRIM-14 (Software-Archaeology Stage) and PRIM-19 (Design Rule Hierarchy Partitioning).

## References (hop-1)

- [[1.0.0 P-89]] Phan et al. ICSE-NIER 2024 legacy-modernization baseline (UNVERIFIED)
- [[1.0.0 P-86]] Xu et al. 2024 LLM-empowered modernization taxonomy (UNVERIFIED)
- [[1.0.0 P-59]] Mindell Digital Apollo 2008 (archaeology lineage)
- [[1.0.0 P-46]] Seacord Modernizing 2003

## References (hop-2)

- SEDCoT (FSE 2025) symbolic-execution-enhanced COBOL translation
- DocSearch (ICSE 2025) agent-oriented documentation generation
