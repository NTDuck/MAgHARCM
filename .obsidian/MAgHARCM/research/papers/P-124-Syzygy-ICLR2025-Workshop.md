---
title: "P-124 Syzygy: Dual Code-Test C-to-Safe-Rust Translation using LLMs and Dynamic Analysis"
backlink: "[[1.0.0 P-124]]"
aliases:
  - "1.0.0 P-124"
  - "P-124"
  - "P-124-Syzygy-ICLR2025-Workshop"
  - "P-124-Syzygy-ICLR2025-Workshop"
  - "Syzygy-ICLR2025-Workshop"
tags: [paper, type-aware-translation, dynamic-analysis, code-translation, rust, "[[1.0.0 PRIM-9]]", "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-30]]", "[[1.0.0 P-88]]", "[[1.0.0 P-85]]", wave-17, iclr-2025-workshop]
date: 2026-09-28
last_updated: 2026-09-28
venue: ICLR 2025 VerifAI Workshop (peer-reviewed)
---

# [[1.0.0 P-124]] Syzygy

## TL;DR

Syzygy is a **dual code-test C-to-safe-Rust translation framework** that combines LLM-driven code/test generation with dynamic-analysis-mined specifications (type, bounds, nullability, aliasing) collected at runtime. The dynamic-analysis step is the **runtime counterpart** to HiTyper's static Type Dependency Graph (TDG) — together they close the type/safety substrate for `[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter and `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph. Demonstrated C-to-safe-Rust translation validated by I/O equivalence tests.

## Mechanism (Q2)

Three stages:

1. **SpecMiner (dynamic analysis)** — Clang/LLVM IR instrumentation passes collect runtime properties: pointer aliasing, allocation sizes, void-pointer types, nullability. These are **syntactically hidden** properties that no static analysis can recover reliably.
2. **Dual code+test generation (LLM-driven)** — for each translation unit (function, struct, typedef), the LLM produces both Rust code AND a Rust equivalence test, conditioned on the SpecMiner properties. Generation proceeds incrementally in dependency order.
3. **Equivalence validation** — the Rust translation must satisfy the mined I/O examples; failures trigger multi-round repair (the analogue of Le Goues automated program repair for the C→Rust direction).

The dual code+test pattern is what makes Syzygy **safe-Rust enforceable** rather than just functionally-equivalent.

## Anchoring (Q3)

| Primitive | Pre-wave-17 behaviour | Syzygy substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Static AST + CFG + DFG; missing runtime semantics | Runtime-mined properties (aliasing, bounds, nullability) injected as a fourth representation |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Static-only concept location | Static TDG (HiTyper P-88) + dynamic property mining (Syzygy) as complementary comprehension dimensions |
| `[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter | Manifests static only | Type/bounds/nullability-enriched manifests for safe-Rust generation |

## Hop-1 Citations

- `[[1.0.0 P-88]]` Peng et al. 2022 HiTyper (ICSE 2022) — Python hybrid static+neural type inference; Syzygy is the C→Rust analogue with the **dynamic** half added.
- C2Rust (Emery et al.) — rule-based C-to-Rust transpiler baseline; Syzygy strictly improves by enforcing safe-Rust + I/O test equivalence.
- VERT (Anonymous 2024) — test-validated C-to-Rust translation baseline; Syzygy improves via dual code+test + dynamic analysis.

## Hop-2 Citations

- Le Goues et al. 2019 automated program repair — Syzygy's multi-round repair loop is the C→Rust analogue.
- LLVM instrumentation framework — SpecMiner's Clang/LLVM IR instrumentation passes.

## MAgHARCM integration

- **YAML config key**: `translation.dynamic_specs: true` (opt-in, default `false`).
- **Implementation file**: `internal/agents/manifest_rewriter.go::EnrichWithDynamicSpecs` (forthcoming — sprint 2026-09-30+).
- **Affected primitives**: `[[1.0.0 PRIM-9]]`, `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-30]]`.

## Caveats

- **ICLR 2025 workshop venue** — peer-reviewed but workshop-track; the §7 trigger gate rewritten 2026-09-25 permits workshop-track papers only when method-level (which Syzygy is). Q1 = PASS under the workshop exemption.
- **SpecMiner requires executable C** — legacy C codebases that cannot be compiled (no headers, missing libraries) cannot use the dynamic-analysis step. The static-only path remains as fallback.
- **SLM prompt budget** — bounded per translation unit (signature + SpecMiner properties + tests), making 4B-30B SLMs feasible per the published evaluation.

## Source

- arXiv:2412.14234 — Syzygy: Dual Code-Test C to (safe) Rust Translation using LLMs and Dynamic Analysis.
- OpenReview: `LYVyioTwvF` (ICLR 2025 VerifAI Workshop).
