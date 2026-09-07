---
title: AlphaTrans — A Neuro-Symbolic Compositional Approach for Repository-Level Code Translation and Validation
bibkey: p02_alphatrans
tags: [paper, repository-translation, neuro-symbolic, [[PRIM-1]], [[PRIM-2]], [[PRIM-3]], [[PRIM-5]], [[PRIM-6]], hop-1]
---

# [[1.0.0 P-02]] AlphaTrans — A Neuro-Symbolic Compositional Approach for Repository-Level Code Translation and Validation

**Authors**: Ali Reza Ibrahimzada, Kaiyao Ke, Mrigank Pawagi, Muhammad Salman Abid, Rangeet Pan, Saurabh Sinha, Reyhaneh Jabbarvand
**Year**: 2025
**Venue**: Proceedings of the ACM on Software Engineering (PACMSE) Vol. 2, FSE, Article FSE1097
**eprint / DOI**: arXiv:2410.24117 (also indexed as `ibrahimzada2024alphatras` in [[refs.bib]])
**Cited by**: [[primitives/INDEX]] entries [[PRIM-1]] (Reverse Topological Translation Ordering), [[PRIM-2]] (Back-Edge-Conditioned Reverse-Topological Scheduling), [[PRIM-3]] (Target Skeleton-First Generation), [[PRIM-5]] (Test Suite Co-Translation & Synthesis), [[PRIM-6]] (Multi-Stage Build/Test Feedback Repair)

## Summary

[[Ibrahimzada-2025-AlphaTrans]] combines static type analysis with a deterministic LLM translation pipeline and GraalVM-based test translation for repository-scale Java-to-Python migration. The method translates a call graph's DAG after cycle-breaking back-edge removal ([P02 §6 ¶1]); emits a target-language type skeleton before any method body ([§5.2]); and uses a fixed-budget repair loop where each iteration runs an LLM-driven `testCheck` pass, collects GraalVM-style I/O traces from the target, and feeds diffs back to the translator ([§6 Alg 2]). Evaluation on 116 projects reports 2.4 % false positives and 7.6 % false negatives under GraalVM-based validation; LLMs hallucinate types, method names, and assertion APIs when the callee is not yet emitted, motivating a fragment-by-fragment order.

## Relevance to MAgHARCM

AlphaTrans supplies the algorithmic underpinnings of [[PRIM-1]] and [[PRIM-2]] — the back-edge-conditioned reverse-topological scheduler that `internal/agents/planning.go` implements. Its skeleton-first emit ([§5.2]) is the formal origin of [[PRIM-3]]'s `DefaultProjectSkeleton`. The testCheck sub-routine ([§6]) is the canonical instance of [[PRIM-5]] and the most-evaluated form of [[PRIM-6]] in the literature. MAgHARCM replaces GraalVM with `cargo check` + `cargo test` because it targets Rust (not Java→Python) and must run on local SLMs, but the cascade topology is borrowed unchanged from Algorithm 2.

## Hop-1 References

- [[ReCodeAgent-2026]] — sister paper by the same group; extends AlphaTrans's cascade into a four-agent explicit loop.
- [[Oxidizer-2024]] — same neuro-symbolic family (LLM + type-driven analysis); Go→Rust instead of Java→Python.
- [[Roziere-2021-Tests]] — Leveraging Automated Unit Tests for Unsupervised Code Translation; the test-as-supervision signal AlphaTrans adopts.
- [[Skel-2025]] — Skeleton-guided translation; converges with AlphaTrans's skeleton emit step on a different PL pair.
- [[Ling-2024-TransRepo-Bench]] — skeleton-first benchmark; AlphaTrans's skeleton step is one of its conceptual ancestors.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000]] — design-rule decomposition mirrors AlphaTrans's stable-type skeleton as the L1 interface layer; method bodies are L3 substitutable leaves.
- [[Rajlich-1997]] — concept-locator analysis gives an offline mechanism for verifying that AlphaTrans's `testCheck` oracle actually covers the source's conceptual surface.
- [[Foltz-2023]] — DR.JONES cognitive model motivates the fragment-by-fragment dispatch: small linear traversals empirically beat deep ones, which is why AlphaTrans schedules leaves first.

## Hop-2 Deep Archaeology & Upstream Lineage

- **Software Modernization Archaeology (Rajlich & Müller)**: Connects modern LLM decompilation/migration back to early software reverse engineering (program slicing, concept assignment, redocumentation).
- **Cognitive Traversal (Foltz & Landauer / Latent Semantic Analysis)**: How human engineers comprehend legacy systems across semantic hops versus how LLM context windows navigate fragmented symbols.
- **Architectural Coupling & Decomposition (Baldwin & Clark / Kazman)**: Modularity theory and Design Structure Matrices (DSM) underpinning why reverse-topological scheduling ([[PRIM-1]]) avoids cyclic cascade failures.
