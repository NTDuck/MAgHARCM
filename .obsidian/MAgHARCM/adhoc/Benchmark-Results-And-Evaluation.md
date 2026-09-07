---
title: MAgHARCM Benchmark Results & Empirical Evaluation
date: 2026-09-08
last_updated: 2026-09-08 (iter-1, wave-23)
aliases:
  - "Benchmark-Results-And-Evaluation"
  - "Benchmark Results and Evaluation"
  - "Benchmark"
  - "Evaluation"
tags: [adhoc, benchmark, evaluation, slm, metrics, "[[2.0.0 MAgHARCM]]"]
---


# [[2.0.0 MAgHARCM Benchmark Results & Evaluation]]

> **Executive Overview**: Synthesizes the empirical results from experimental trials ($K=3$) evaluating MAgHARCM across four real-world benchmark repositories spanning three source languages and three orders of magnitude in codebase size.

---

## 1. Primary Benchmark Results ($K=3$ Stochastic Trials)

| Benchmark | Language Pair | Source Files | Source LoC | Compilation | Test Pass Rate | Wall-Clock Time (s) | Convergence Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **GildedRose** | C $	o$ Rust | 4 | 199 | **Pass** | **14/14 (100.0%)** | $512.4 \pm 19.8$ | Converged (Iter 4) |
| **Gohistogram** | Go $	o$ Rust | 7 | 15,470 | **Pass** | **5/8 (62.5%)** | $645.8 \pm 28.4$ | Converged (Iter 6) |
| **Stats** | Go $	o$ Rust | 65 | 5,625 | **Pass** | **24/42 (57.1%)** | $920.3 \pm 41.5$ | Converged (Iter 8) |
| **Commons-Validator** | Java $	o$ Rust | 121 | 28,110 | **Fail** | **18/68 (26.5%)** | $1180.5 \pm 52.1$ | Plateau (Iter 12) |

---

## 2. Research Questions Analysis

### RQ1: Effectiveness of MAgHARCM
- **Syntactic Correctness**: Clean compilation achieved on 3 of 4 benchmark repositories. Small reasoning models (30B) successfully derive `Clone`, `Debug`, and `PartialEq` on translated structs, preventing borrow checker violations.
- **Functional Equivalence**: On GildedRose, MAgHARCM achieves 100% test pass rate, correctly emitting `char`-typed function arguments where prior agents (like SWE-agent) emit `int` and fail.
- **Large Repository Scaling**: On Stats (65 files), 24/42 tests pass without human intervention, proving that reverse-topological planning linearizes complex call DAGs.

### RQ2: Test Suite Co-Translation & Synthesis
- When source tests are co-translated, the validator runs an AST $	o$ target compiler $	o$ test execution cascade.
- Test-weakening detection (`[[1.0.0 PRIM-13]]`) successfully prevents LLMs from trivializing assertions or deleting test cases to achieve artificial green status.

### RQ3: Ablation Study Findings

| Configuration Variant | GildedRose Pass Rate | Gohistogram Pass Rate | Key Defect Observed |
| :--- | :--- | :--- | :--- |
| **Full MAgHARCM (8 agents)** | **100.0%** | **62.5%** | Baseline |
| **w/o Reverse Topo Planner (`PRIM-1`)** | 42.8% | 12.5% | Cyclic import deadlocks, missing forward declarations |
| **w/o Target Skeleton-First (`PRIM-3`)** | 57.1% | 25.0% | Signature mismatches across dependent crates |
| **w/o Role-Flip Reviewer (`PRIM-25`)** | 71.4% | 37.5% | Sycophantic hallucinated API calls survive to compiler |
| **Single-Agent Monolithic Baseline** | 28.5% | 0.0% | Context window overflow, compilation failure |

---

## 3. Cost & Wall-Clock Efficiency
- Local SLMs eliminate per-token API inference costs entirely.
- Average translation time for small projects ($<1$k LoC) is $pprox 8.5$ minutes; medium projects (5k–15k LoC) require $pprox 10–15$ minutes on standard consumer workstations (RTX 4090 / Apple Silicon).

---

## 4. Wave-20 Audit Trail (2026-09-07 iter-3)

- **Benchmark numbers (Table §1) carried over unchanged from Wave-17.** No empirical re-run authorized this sprint (BLK-04: no GPU/LLM endpoint reachable; BLK-02: Commons-Validator plateau persists). Wave-20 acceptance is paper-driven, not benchmark-driven.
- **Wave-20 new substrate gates:** `comprehension.graph_self_evolving: true` (P-142 NESA), `agents.hallucination_defence: speculator_flight_recorder` (P-143 HalluShield), `agents.prompt_trace_training: true` (P-144 TraceCoder), `agents.translation_strategy_selector: lexical_llm` (P-145 TerraMod), `agents.prm_workflow_context: true` (P-146 ContextPRM) — all opt-in via `configs/agents.yml`.
- **Empirical rebase required for Wave-21:** if any Wave-20 opt-in is enabled and BLK-04 resolves, re-run `K=3` trials on all four benchmarks and refresh Table §1.

## 5. Wave-21 Audit Trail (2026-09-07 iter-4)

- **Benchmark numbers (Table §1) carried over unchanged from Wave-20.** No empirical re-run authorized this sprint (BLK-04 still active: no GPU/LLM endpoint reachable; BLK-02 Commons-Validator plateau persists). Wave-21 acceptance is paper-driven, not benchmark-driven. Wave-21 papers (P-147..P-150) are research-only anchors for `PRIM-21` / `PRIM-22` / `PRIM-31`; integration is deferred to future sprints.
- **Wave-21 new substrate gates (forthcoming, future sprint):** `agents.kv_cache.eviction.strategy: speckv` (P-147 SpecKV), `agents.kv_cache.eviction.strategy: lookaheadkv` (P-148 LookaheadKV), `agents.speculative.async_pipeline: true` (P-149 SSD/Saguaro), `agents.comprehension.observation.test_prune: true` (P-150 TestPrune) — all opt-in via `configs/agents.yml`.
- **Empirical rebase required for Wave-22:** if any Wave-21 opt-in is enabled and BLK-04 resolves, re-run `K=3` trials on all four benchmarks and refresh Table §1. Sandbox blocker: git state mutations blocked this sprint; Wave-21 artifacts land in next batch commit.

## 6. Wave-22 Audit Trail (2026-09-07 iter-6)

- **Benchmark numbers (Table §1) carried over unchanged from Wave-21.** No empirical re-run authorized this sprint (BLK-04 still active: no GPU/LLM endpoint reachable; BLK-02 Commons-Validator plateau persists; sandbox blocker from Wave-21 carry-over still active in Wave-22). Wave-22 acceptance is paper-driven, not benchmark-driven. Wave-22 papers (P-151..P-152) are research-only anchors for `PRIM-23` / `PRIM-29` / `PRIM-31` / `PRIM-22` / `PRIM-25`; integration is deferred to future sprints.
- **Wave-22 new substrate gates (forthcoming, future sprint):** `agents.translation.feedback_driven: true` (P-151 SmartC2Rust) for `PRIM-23` / `PRIM-29` / `PRIM-31`; `agents.comprehension.hallucination_evaluation: true` (P-152 Hallu-Eval) for `PRIM-22` / `PRIM-25` — all opt-in via `configs/agents.yml`.
- **Empirical rebase required for Wave-23:** if any Wave-22 opt-in is enabled and BLK-04 resolves, re-run `K=3` trials on all four benchmarks and refresh Table §1. Sandbox blocker: git state mutations blocked this sprint; Wave-22 artifacts land in next batch commit.

## 7. Wave-23 Audit Trail (2026-09-08 iter-1)