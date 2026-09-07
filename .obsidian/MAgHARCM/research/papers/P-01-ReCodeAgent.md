---
title: ReCodeAgent — A Multi-Agent Workflow for Language-agnostic Translation and Validation of Large-scale Repositories
backlink: "[[1.0.0 P-01]]"
bibkey: p01_recodeagent
aliases:
  - "1.0.0 P-01"
  - "P-01"
  - "P-01-ReCodeAgent"
  - "P-01-ReCodeAgent"
  - "ReCodeAgent"
tags: [paper, repository-translation, multi-agent, [[PRIM-1]], [[PRIM-3]], [[PRIM-5]], [[PRIM-6]], hop-1]
---

# [[1.0.0 P-01]] ReCodeAgent — A Multi-Agent Workflow for Language-agnostic Translation and Validation of Large-scale Repositories

**Authors**: Ali Reza Ibrahimzada, Brandon Paulsen, Daniel Kroening, Reyhaneh Jabbarvand
**Year**: 2026
**Venue**: Proceedings of the ACM International Conference on the Foundations of Software Engineering (FSE)
**eprint / DOI**: arXiv:2604.07341
**Cited by**: [[primitives/Primitives-Index]] entries [[PRIM-1]] (Reverse Topological Translation Ordering), [[PRIM-3]] (Target Skeleton-First Generation), [[PRIM-5]] (Test Suite Co-Translation & Synthesis), [[PRIM-6]] (Multi-Stage Build/Test Feedback Repair)

## Summary

[[Ibrahimzada-2026-ReCodeAgent]] decomposes repository-scale translation into four specialised LLM agents — Analyzer, Planning, Translator, Validator — coordinated through an explicit `Analyze → Plan → Translate → Validate` loop. The Planner emits a bottom-up module order from a constructed call graph ([ReCodeAgent Sec 3.3, p. 5]); the Translator starts from a target-language skeleton of interfaces, imports, and signature-only functions before any method body is produced ([Sec 3.3.3]); source-language unit tests are co-translated and, when source tests are missing, additional LLM-generated tests are appended to the suite ([Sec 3.5.2; Alg 1 lines 14-21]); and the Validator runs a three-tier cascade — AST/syntax check, target compiler check, translated test suite execution — and feeds diagnostics back to the Translator under a fixed budget ([Alg 1]). Evaluated on 118 projects across four language pairs and six languages, reporting 99.4 % compile success and 86.5 % test pass on Anthropic Claude Sonnet via Claude Code; the cloud LLM is the principal reason MAgHARCM's local-SLM topology must re-derive the same primitives.

## Relevance to MAgHARCM

ReCodeAgent is the closest prior work and the spine of [[PRIM-1]], [[PRIM-3]], [[PRIM-5]], [[PRIM-6]]. MAgHARCM's `internal/agents/planning.go::ComputeReverseTopoOrder` ([[PRIM-1]]) is the local-SLM re-implementation of ReCodeAgent's bottom-up module plan; `internal/agents/planning.go::DefaultProjectSkeleton` ([[PRIM-3]]) mirrors the skeleton-first emit step; `internal/agents/validator.go::generateAdditionalTests` ([[PRIM-5]]) realises the algorithm-1 lines 14-21 test-synthesis path; and `internal/agents/validator.go` ([[PRIM-6]]) runs the same AST → compiler → test cascade with rustc diagnostics fed back to the chunked translator ([[PRIM-23]]). The four-agent decomposition is adopted wholesale; only the cloud-LLM assumption is replaced by the local 30 B reasoning + 4 B coding topology.

## Hop-1 References

- [[AlphaTrans-2024]] — reverse-topological method translation; same back-edge-conditioned DAG view of the call graph.
- [[Oxidizer-2024]] — Go-to-Rust validated translation with mock-test synthesis; co-translation of tests.
- [[Roziere-2020-Transcoder]] — unsupervised neural code translation; baseline family ReCodeAgent contrasts against.
- [[Ibrahimzada-2024-AlphaTrans-Validation]] — sibling paper by overlapping authors on GraalVM-based validation; informs ReCodeAgent's 3-tier cascade.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997]] — concept-locator analysis is the offline complement to ReCodeAgent's Analyzer agent; together they map concepts→locations and locations→translations.
- [[Baldwin-Clark-2000]] — design-rule decomposition underlies ReCodeAgent's skeleton-first decision: identify stable interfaces (L1), then treat modules as substitutable leaves (L3).
- [[Kazman-Cai-2024]] — architectural recovery techniques motivate ReCodeAgent's pre-Planning call-graph construction; an archaeological pass can recover the call graph when static analysis overapproximates cycles.

## Hop-2 Deep Archaeology & Upstream Lineage

- **Software Modernization Archaeology (Rajlich & Müller)**: Connects modern LLM decompilation/migration back to early software reverse engineering (program slicing, concept assignment, redocumentation).
- **Cognitive Traversal (Foltz & Landauer / Latent Semantic Analysis)**: How human engineers comprehend legacy systems across semantic hops versus how LLM context windows navigate fragmented symbols.
- **Architectural Coupling & Decomposition (Baldwin & Clark / Kazman)**: Modularity theory and Design Structure Matrices (DSM) underpinning why reverse-topological scheduling ([[PRIM-1]]) avoids cyclic cascade failures.
