---
title: RustRepoTrans — Repository-Level Code Translation Benchmark Targeting Rust
bibkey: p07_rustrepotrans
tags: [paper, repository-translation, benchmark, rust, [[PRIM-10]], [[PRIM-3]], hop-1]
---

# [[1.0.0 P-07]] RustRepoTrans — Repository-Level Code Translation Benchmark Targeting Rust

**Authors**: Mengyang Ou, Yang Liu, Mingyue Jiang, Qingyang Zhang, Yanjun Zhang, Dechuang Zhao, Zhe Hou, Yafei Wu, Shikun Zhang
**Year**: 2024
**Venue**: arXiv preprint cs.SE (arXiv:2411.13990, multiple revisions through v4)
**eprint / DOI**: arXiv:2411.13990
**Cited by**: [[primitives/INDEX]] entry [[PRIM-10]] (Feature-Mapping & Type-Compatibility Validation); informs [[PRIM-3]] skeleton choices for non-Rust→Rust translation tasks

## Summary

[[Ou-2024-RustRepoTrans]] introduces the first repository-level context code-translation benchmark targeting Rust, comprising 375 tasks spanning C→Rust (145), Java→Rust (122), and Python→Rust (108) function pairs. The dataset is built from real-world open-source repositories and is paired with full test suites so that translation correctness can be measured at the function-pair level. The authors evaluate four state-of-the-art LLMs (GPT-4o, Claude-3.5-Sonnet, DeepSeek-Coder-V2, and CodeLlama-70B) and report that all four perform substantially worse than on function-level benchmarks (e.g. [[TransRepo-Bench-2025]]), with success rates below 30 % on Java→Rust and below 15 % on Python→Rust. Failure analysis surfaces three dominant error classes: (1) feature-mapping errors (Go-style `interface` not lifted to a `trait`; Java checked exceptions not mapped to `Result<T,E>`); (2) lifetime/ownership errors that the LLM hallucinates around mutable borrows; and (3) cross-crate dependency errors where the LLM imports a crate that does not exist on crates.io. RustRepoTrans is hosted at SYSUSELab/RustRepoTrans on GitHub with the evaluation harness, prompts, and per-task result cache.

## Relevance to MAgHARCM

RustRepoTrans is the primary empirical motivation for [[PRIM-10]] (Feature-Mapping & Type-Compatibility Validation). Its failure-analysis tables directly inform the idiom-mapping table that MAgHARCM's `internal/agents/feature_mapping.go` (planned this sprint) maintains — specifically the Go-interface ↔ Rust-trait, Go-error-return ↔ `Result<T,E>`, and cross-crate import resolution rules. Its Java→Rust subset is MAgHARCM's evaluation target for `internal/agents/manifest_rewriter.go` [[PRIM-10 partial]] because it surfaces the dependency-list discrepancy problem (source pom.xml imports do not have a 1:1 Cargo.toml equivalent). The benchmark also gives MAgHARCM a known-function-pair evaluation set on which to compare translator variants; the per-task result cache is reusable for regression testing of `internal/agents/chunked_translator.go` [[PRIM-23]].

## Hop-1 References

- [[Ibrahimzada-2025-AlphaTrans]] — sibling Java-to-Python repository translation paper; RustRepoTrans extends its evaluation design to a harder target language (Rust).
- [[Ibrahimzada-2026-ReCodeAgent]] — language-agnostic translation workflow; RustRepoTrans validates its claims about Rust as a low-resource target.
- [[Oxidizer-2024]] — Go→Rust specialisation; RustRepoTrans is the multi-source-language companion benchmark.
- [[Ling-2024-TransRepo-Bench]] — Java→C# benchmark; RustRepoTrans adopts the per-function-pair test scoring approach.
- [[Roziere-2021-Tests]] — test-as-supervision signal that RustRepoTrans uses for pass/fail evaluation.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000]] — RustRepoTrans's error class 1 (feature-mapping) is exactly the design-rule decomposition problem: source-language idioms belong at different stability layers than target-language idioms.
- [[Rajlich-1997]] — concept-locator analysis explains error class 2 (lifetime/ownership hallucination): the LLM is generating from surface syntax, but the concept of ownership has no source-language equivalent and must be discovered offline.
- [[Kazman-Cai-2024]] — architectural recovery underwrites error class 3 (cross-crate dependency errors): recovering the source module-dependency graph is a prerequisite to mapping it to crates.io.
- [[Foltz-2023]] — DR.JONES cognitive model explains the LLM's bias toward plausible-looking but uncompilable Rust: recency-biased generation of syntactically familiar constructs without deeper ownership reasoning.
