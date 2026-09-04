---
title: TRAM — Advancing Automated In-Isolation Validation in Repository-Level Code Translation
bibkey: p10_tram
tags: [paper, repository-translation, mock-validation, in-isolation, [[PRIM-8]], [[PRIM-10]], hop-1]
---

# TRAM — Advancing Automated In-Isolation Validation in Repository-Level Code Translation

**Authors**: Kaiyao Ke, Ali Reza Ibrahimzada, Rangeet Pan, Saurabh Sinha, Reyhaneh Jabbarvand
**Year**: 2025 (arXiv preprint 2025-11; no conference acceptance surfaced as of 2026-08-30)
**Venue**: arXiv preprint cs.SE (arXiv:2511.21878, Nov 2025)
**eprint / DOI**: arXiv:2511.21878
**Cited by**: [[primitives/INDEX]] entry [[PRIM-8]] (State-Grounded Mock-Based In-Isolation Validation); influences [[PRIM-10]] (Feature-Mapping & Type-Compatibility Validation) through its context-aware type-resolution step

## Summary

[[Ke-2025-TRAM]] addresses the cost of validating repository-level code translations by introducing mock-based in-isolation validation. Rather than running the entire translated repository against an agent-based verdict ([[Ibrahimzada-2025-MatchFixAgent]]) or against a full compiler + test cascade ([[Ibrahimzada-2025-AlphaTrans]] §6), TRAM validates each translated method *in isolation* by constructing a mock of its dependencies in the target language. The construction proceeds in two stages. First, a context-aware type-resolution step retrieves API documentation and contextual code snippets for the source-language types used by the method, then asks an LLM to produce a target-language type mapping that respects field names, mutability, and nullability. Second, a custom serialisation/deserialisation workflow emits mock objects in the target language that capture the state-equivalent behaviour of the source dependency. The translated method is then compiled and exercised against the mock; only when each method passes its in-isolation check does the system assemble and validate the full system. Evaluation on Java→Python reports state-of-the-art translation accuracy with substantially lower LLM cost than [[Ibrahimzada-2025-MatchFixAgent]] and lower human engineering effort than [[Abid-2024-GlueTest]].

## Relevance to MAgHARCM

TRAM is the direct citation for [[PRIM-8]] (State-Grounded Mock-Based In-Isolation Validation). MAgHARCM's `internal/agents/mock_validator.go` (planned this sprint) realises the per-fragment validation pipeline: for each Go method being translated, generate a Rust trait-implementation mock of its callee types using the [[PRIM-10]] idiom-mapping table, compile the translated method against the mock, and record pass/fail before proceeding to the next fragment. The cost argument is crucial for MAgHARCM's local-SLM-only operating mode — TRAM's per-method validation cost is roughly one-fifth of MatchFixAgent's per-pair verdict cost, putting it within budget for the 30K–50K LoC commons-validator sample. TRAM's context-aware type-resolution step extends [[PRIM-10]] with field-level fidelity (not just identifier mapping), and its assembly-then-validate-second-pass design matches the chunked-translation strategy [[PRIM-23]].

## Hop-1 References

- [[Ibrahimzada-2025-MatchFixAgent]] — same group; TRAM is the per-method mock-based alternative to MatchFixAgent's whole-translation verdict.
- [[Ibrahimzada-2025-AlphaTrans]] — same group; TRAM replaces AlphaTrans's GraalVM-based test translation with mock-based isolation.
- [[Oxidizer-2024]] — same neuro-symbolic family; TRAM's context-aware type resolution extends Oxidizer's feature-mapping table to field-level semantics.
- [[Abid-2024-GlueTest]] — language-interoperability validation; TRAM is the mock-based, non-interop alternative.
- [[Ibrahimzada-2026-ReCodeAgent]] — language-agnostic translation; TRAM is the validation subsystem that could be slotted into ReCodeAgent's pipeline.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000]] — TRAM's per-method isolation is a design-rule decomposition: each method's contract (its type signature + behaviour) is a stable load-bearing rule, and the mock exists to make that rule independently testable.
- [[Rajlich-1997]] — concept-locator analysis explains why TRAM needs context-aware type resolution: the type mapping is concept-level (what does this `Map<K,V>` *mean*?), not identifier-level (which struct replaces `java.util.HashMap`?).
- [[Kazman-Cai-2024]] — architectural recovery underwrites TRAM's mock construction: knowing which dependencies a method actually touches (vs. which ones it imports) requires call-graph + data-flow recovery.
- [[Müller-2000]] — TRAM's "isolate-then-assemble" pattern is the Integrate-in-Place migration strategy applied at the method granularity; [[PRIM-21]]'s strategy selector should consider this option for large repositories.
- [[Foltz-2023]] — DR.JONES explains why per-method isolation improves validation accuracy: linear traversal of a small mock beats whole-program execution in terms of cognitive traceability.

## Hop-2 Deep Archaeology & Upstream Lineage

- **Software Modernization Archaeology (Rajlich & Müller)**: Connects modern LLM decompilation/migration back to early software reverse engineering (program slicing, concept assignment, redocumentation).
- **Cognitive Traversal (Foltz & Landauer / Latent Semantic Analysis)**: How human engineers comprehend legacy systems across semantic hops versus how LLM context windows navigate fragmented symbols.
- **Architectural Coupling & Decomposition (Baldwin & Clark / Kazman)**: Modularity theory and Design Structure Matrices (DSM) underpinning why reverse-topological scheduling ([[PRIM-1]]) avoids cyclic cascade failures.
