---
title: CodeS-bench (placeholder) — Repository-Scale Code Translation Benchmark
bibkey: p06_codes_bench
tags: [paper, repository-translation, benchmark, [[PRIM-6]], [[PRIM-13]], hop-1, phantom-bibkey]
---

# [[1.0.0 P-06]] CodeS-bench (placeholder) — Repository-Scale Code Translation Benchmark

**Authors**: Tang et al. (cited in methodology ledger; no confirmed peer-reviewed publication located)
**Year**: 2024 (per ledger)
**Venue**: ACM ISSTA 2024 (per ledger — UNVERIFIED, the ISSTA 2024 proceedings do not surface a benchmark of this exact name)
**eprint / DOI**: ledger shows malformed `2404.x]`; no arXiv ID found
**Cited by**: [[primitives/INDEX]] entry [[PRIM-6]] (Multi-Stage Build/Test Feedback Repair); referenced indirectly by [[PRIM-13]] (Adversarial Test-Weakening Guard) as a same-class benchmark

## Summary

The bibliography slot `p06_codes_bench` is reserved for a repository-scale code-translation benchmark purportedly authored by Tang et al. and published at ACM ISSTA 2024. Targeted web search 2026-09-04 against the ISSTA 2024 proceedings (Vienna, 16–20 Sep 2024), Researchr, and arXiv cs.SE does not surface any benchmark titled "CodeS-bench" by an author named Tang; the closest neighbours are class-level work ([[ClassEval-T-2025]]) and translation-benchmark studies (e.g. [[TransRepo-Bench-2025]]). Because P06's appendix primitive (Multi-Stage Build/Test Feedback Repair) is algorithmically instantiated in [[Ibrahimzada-2026-ReCodeAgent]] Algorithm 1 and [[Ibrahimzada-2025-AlphaTrans]] [[PRIM-6]] §6 Alg 2, the *actual* empirical evidence underpinning [[PRIM-6]] is borrowed from those two papers rather than from this phantom benchmark. [[AlphaTrans-2025-P02]] and [[ReCodeAgent-2026]] together describe a 3-tier cascade (AST/syntax check → target-compiler check → translated-test execution) with a fixed-budget feedback loop whose diagnostics are returned to the translator up to a cap of repair iterations.

## Relevance to MAgHARCM

P06's primitive is the architectural backbone of `internal/agents/validator.go::generateAdditionalTests` and the surrounding repair-loop logic. Even though `p06_codes_bench` is a phantom reference, the primitive is *measured* in [[AlphaTrans-2025]] §6 and [[ReCodeAgent-2026]] Alg 1, and MAgHARCM borrows both the cascade topology (parser → compiler → test) and the fixed-budget feedback semantics. The placeholder entry is retained in the ledger so that future empirical work on a real CodeS-bench-class benchmark (e.g. a Java→Rust or Go→Rust repository-scale evaluation) can be slotted in without renaming the primitive. Until then, MAgHARCM treats [[PRIM-6]] as having no independent measurement — its claims are inherited from [[PRIM-6]]'s two real citations.

## Hop-1 References

- [[Ibrahimzada-2026-ReCodeAgent]] — co-author of the ledger's P06 placeholder; Algorithm 1 is the canonical 3-tier cascade instance.
- [[Ibrahimzada-2025-AlphaTrans]] — same group; §6 Alg 2 is the second canonical instance with GraalVM test translation.
- [[Oxidizer-2024]] — Go→Rust cascade with type-driven analysis and mock-test synthesis; corroborates the multi-stage pattern.
- [[Ling-2024-TransRepo-Bench]] — repository-scale translation benchmark; same class as the placeholder CodeS-bench.
- [[Roziere-2021-Tests]] — the test-as-supervision signal that motivates the third tier of the cascade.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997]] — concept-locator analysis underwrites the *what to compile and test* decision in tier 2/3: tests should cover the source's concept surface, not its call-graph surface.
- [[Foltz-2023]] — DR.JONES explains why the cascade stops at three tiers: empirically, repair budgets beyond three feedback dimensions do not improve LLM translation outcomes.
- [[Baldwin-Clark-2000]] — the cascade's tier 2 (target compiler) is the L1 load-bearing check; tier 3 (test execution) is the L3 leaf check. A design-rule decomposition clarifies why tier 2 must precede tier 3.
