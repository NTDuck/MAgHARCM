---
title: Oxidizer — Scalable, Validated Code Translation of Entire Projects using Large Language Models
bibkey: p05_oxidizer
tags: [paper, repository-translation, go-to-rust, type-mapping, [[PRIM-10]], hop-1]
---

# Oxidizer — Scalable, Validated Code Translation of Entire Projects using Large Language Models

**Authors**: Hanliang Zhang, Cristina David, Meng Wang, Brandon Paulsen, Daniel Kroening
**Year**: 2024
**Venue**: arXiv preprint (cs.PL, cs.SE) arXiv:2412.08035
**eprint / DOI**: arXiv:2412.08035
**Cited by**: [[primitives/INDEX]] entry [[PRIM-10]] (Feature-Mapping & Type-Compatibility Validation)

## Summary

[[Zhang-2024-Oxidizer]] targets Go-to-Rust translation by combining LLM synthesis with a type-driven program-analysis pipeline. It builds and maintains an explicit mapping table from Go language features to their Rust equivalents — interfaces ↔ traits, error returns ↔ `Result`, slices ↔ `Vec`, channels ↔ crossbeam — and validates every emitted Rust symbol against its Go counterpart ([P05 §3]). Test coverage is used as a proxy for translation correctness; when source tests are absent, Oxidizer synthesises mock tests in the target language whose I/O traces are compared against a sandboxed execution of the Go source. Evaluation reports that explicit feature mapping measurably reduces hallucinated type errors but does not solve the absent-test-suite problem (which MAgHARCM's [[PRIM-5]] test synthesis and [[PRIM-6]] cascade address).

## Relevance to MAgHARCM

Oxidizer is the direct citation for [[PRIM-10]] (Feature-Mapping & Type-Compatibility Validation). MAgHARCM's `internal/agents/feature_mapping.go` (planned this sprint per the methodology ledger) maintains the same idiom-mapping table — Go interfaces ↔ Rust traits, Go error returns ↔ `Result<T,E>` — and validates every emitted Rust symbol against its source counterpart via the navigator. `internal/agents/manifest_rewriter.go` ([PRIM-10 partial]) is the sibling file that applies a curated crate-mapping table to the dependency manifest, also from Oxidizer §3. The Go→Rust pair is MAgHARCM's primary evaluation domain, so Oxidizer's empirical claims translate directly into the harness's design constraints.

## Hop-1 References

- [[AlphaTrans-2024]] — sibling neuro-symbolic system; targets Java→Python with GraalVM validation rather than Go→Rust.
- [[ReCodeAgent-2026]] — same author overlap (Kroening, Paulsen) on a language-agnostic translation system; Oxidizer is the Go→Rust specialisation.
- [[RustMap-2025]] — C-to-Rust project-scale migration via program analysis and LLM; co-cited with Oxidizer in PL-translation surveys.
- [[Roziere-2021-Tests]] — test-as-supervision signal that Oxidizer's mock-test synthesis extends to a new PL pair.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000]] — Oxidizer's feature-mapping table is a domain-specific design-rule catalogue; stable inter-language mapping belongs at L1.
- [[Rajlich-1997]] — concept-locator analysis motivates *why* certain Go idioms map to specific Rust constructs: the underlying concept (error propagation, polymorphism, iteration) is what is being translated, not the surface syntax.
- [[Foltz-2023]] — DR.JONES cognitive model explains why an explicit mapping table beats leaving the LLM to infer type equivalents: recency-biased linear traversal of an explicit catalog outperforms deep inference.
- [[Kazman-Cai-2024]] — architectural recovery underwrites Oxidizer's manifest-mapping step: identifying which Go modules contain which idioms depends on recovery of the module-dependency graph.
