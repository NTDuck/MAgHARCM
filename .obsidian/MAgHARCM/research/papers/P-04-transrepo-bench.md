---
title: TransRepo-Bench — Skeleton-Guided Translation: A Benchmarking Framework for Code Repository Translation with Fine-Grained Quality Evaluation
bibkey: p04_transrepo_bench
tags: [paper, repository-translation, benchmark, [[PRIM-3]], [[PRIM-5]], [[PRIM-23]], hop-1]
---

# [[1.0.0 P-04]] TransRepo-Bench — Skeleton-Guided Translation: A Benchmarking Framework for Code Repository Translation with Fine-Grained Quality Evaluation

**Authors**: Xing Zhang, Jiaheng Wen, Fangkai Yang, Pu Zhao, Yu Kang, Junhao Wang, Maoquan Wang, Yufan Huang, Shengyu Fu, Elsie Nallipogu, Qingwei Lin, Yingnong Dang, Saravan Rajmohan, Dongmei Zhang
**Year**: 2025 (ACL 2024 per methodology ledger; ACL Anthology Findings of EMNLP 2025)
**Venue**: Findings of EMNLP 2025 (ACL Anthology 2025.findings-emnlp.986); arXiv:2501.16050 (the ledger's `2405.08900` is a placeholder eprint; the canonical arXiv ID is 2501.16050)
**eprint / DOI**: arXiv:2501.16050; ACL Anthology 2025.findings-emnlp.986
**Cited by**: [[primitives/INDEX]] entry [[PRIM-3]] (Target Skeleton-First Generation); influences [[PRIM-5]] and [[PRIM-23]] methodology choice in the validator cascade.

## Summary

[[Zhang-2025-TransRepo-Bench]] introduces a two-step translation framework ("Skeleton-Guided-Translation") and the TRANSREPO-BENCH dataset of high-quality Java repositories paired with unit tests and build configurations. Step 1 emits the target repository's structural skeleton — file structure, class/method signatures, imports — without bodies; Step 2 refines the full repository using those skeletons as a guide to preserve inter-module coherence and dependencies. Evaluation is fine-grained per test case rather than binary pass/fail, surfacing which specific translations regress. The benchmark and framework together demonstrate that emitting a skeleton first measurably improves test-level metrics on Java→C# translation across multi-module repositories.

## Relevance to MAgHARCM

TransRepo-Bench is the empirical motivation for [[PRIM-3]] (Target Skeleton-First Generation). `internal/agents/planning.go::DefaultProjectSkeleton` realises the same two-step pattern: skeleton emit first, body fill second, with the previously emitted skeleton as a context prefix. Its fine-grained test-level evaluation also motivates [[PRIM-13]] (Adversarial Test-Weakening Guard) — binary pass/fail can hide weakened assertions that per-test scoring exposes. The benchmark's Java→C# scope is a useful comparator for MAgHARCM's Go→Rust focus; both are static-typed targets with explicit manifests.

## Hop-1 References

- [[AlphaTrans-2024]] — skeleton-first emit ([§5.2]) is the immediate ancestor of TransRepo-Bench's two-step framework.
- [[ReCodeAgent-2026]] — also skeleton-first; differs in its four-agent loop rather than explicit two-step.
- [[Oxidizer-2024]] — provides target-language scaffolding for Go→Rust; TransRepo-Bench references its manifest-mapping work.
- [[Roziere-2021-Tests]] — test-coverage-as-supervision is the evaluation signal TransRepo-Bench refines into per-test scoring.
- [[CodeS-2024]] — companion repository-scale benchmark (P06 in the ledger) for cross-evaluator comparison.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000]] — the skeleton-first two-step is exactly the design-rule decomposition: fix the load-bearing interfaces (L1), then treat modules as leaves (L3) that can be freely regenerated.
- [[Rajlich-1997]] — concept-locator analysis provides the offline rationale for the skeleton: the skeleton encodes the *concepts*, and the leaves are *implementations*.
- [[Kazman-Cai-2024]] — architectural-recovery techniques inform TransRepo-Bench's inter-module coherence check, which is a static recovery of the call skeleton.
