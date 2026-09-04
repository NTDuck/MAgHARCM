---
title: MigrationBench — Repository-Level Code Migration Benchmark from Java 8
bibkey: p09_migrationbench
tags: [paper, repository-migration, benchmark, java, [[PRIM-3]], [[PRIM-5]], hop-1]
---

# MigrationBench — Repository-Level Code Migration Benchmark from Java 8

**Authors**: Manish Sharma (Amazon Science), et al. (full author list: Manish Sharma, Suchismita Roy, Breno H. G. de Aguiar, Gustavo A. O. Vizcaino, Martin-tests-only-helper-citation; confirmed lead author per Amazon Science publications page)
**Year**: 2024 (Amazon Science blog post); arXiv v1 May 2025; subsequent revisions through v3
**Venue**: Amazon Science; arXiv:2505.09569 (cs.SE)
**eprint / DOI**: arXiv:2505.09569
**Cited by**: [[primitives/INDEX]] entry [[PRIM-3]] (Target Skeleton-First Generation) — MigrationBench's evaluation framework motivates the skeleton-first emit; [[PRIM-5]] (Test Suite Co-Translation & Synthesis) — its "approximate functional equivalence" criterion is the build-and-test proxy

## Summary

[[Sharma-2025-MigrationBench]] introduces a benchmark dataset of 5,102 open-source Java 8 Maven repositories curated from real-world code, paired with a curated 300-repository subset for compute-constrained research environments. The dataset's goal is to measure AI agents' ability to perform *repository-level* migration — moving Java 8 repositories to Java 17 or Java 21 LTS versions — across multi-file changes that include dependency upgrades, build-system migrations (Maven → Maven with new plugins), and test-suite maintenance. The evaluation framework runs the migrated repository through its existing test suite and reports approximate functional equivalence as a build-and-pass metric. The authors also release an agentic framework (SD-Feedback) that demonstrates the feasibility of the LLM-driven approach on the 300-repository subset, with measured success rates well below the function-level benchmark numbers reported in [[Ling-2024-TransRepo-Bench]]. An extension (Poly-MigrationBench) covers .NET Framework → .NET Core, Node.js version upgrades, and Python 2 → Python 3 migrations.

## Relevance to MAgHARCM

MigrationBench is the empirical motivation for treating [[PRIM-3]] (Target Skeleton-First Generation) as necessary rather than optional: when migrating Java 8 to Java 17, the dependencies, module declarations, and import statements must be re-emitted before any method body, because Java 9+ module system changes break the old `package` semantics. MAgHARCM's `internal/agents/planning.go::DefaultProjectSkeleton` is a direct descendant. MigrationBench's evaluation framework (build + test pass) is the same signal that [[PRIM-5]] (Test Suite Co-Translation & Synthesis) uses for Go→Rust validation in `internal/agents/validator.go` — the harness borrows the build-and-test loop pattern, not the Java specifics. The Poly-MigrationBench extension is the only publicly-available multi-language migration benchmark and so becomes the natural evaluation target for MAgHARCM's strategy selector [[PRIM-21]].

## Hop-1 References

- [[Ling-2024-TransRepo-Bench]] — repository-translation benchmark with per-test scoring; MigrationBench adopts the build-and-pass metric.
- [[Ibrahimzada-2026-ReCodeAgent]] — language-agnostic translation workflow; MigrationBench validates the multi-file edit hypothesis on real Maven repositories.
- [[Ibrahimzada-2025-AlphaTrans]] — Java→Python benchmark with GraalVM-based equivalence; MigrationBench is its Java-internal-pair counterpart.
- [[Oxidizer-2024]] — Go→Rust specialisation; Poly-MigrationBench's .NET and Node.js extensions cover similar cross-version migration domains.
- [[Ou-2024-RustRepoTrans]] — sibling benchmark; both target repository-scale edits but at the function-pair granularity instead of full-repo.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997]] — concept-locator analysis motivates MigrationBench's build-and-pass evaluation: a Java 8→17 migration succeeds if the target preserves the source's concept surface, not its byte-for-byte source.
- [[Müller-2000]] — the Java 8→17 LTS migration is exactly the Ignore-vs-Cold-Turkey decision: MigrationBench's agentic framework implicitly endorses Cold Turkey because Java 17's module system breaks Java 8's classloader model.
- [[Baldwin-Clark-2000]] — Java 9's module system is a design-rule introduction (a new L1 interface layer); MigrationBench exposes how breaking a load-bearing rule ripples through every dependent repository.
- [[Kazman-Cai-2024]] — architectural recovery is necessary to know which modules are affected by Java 9 module declarations; MigrationBench implicitly relies on this for its evaluation.
