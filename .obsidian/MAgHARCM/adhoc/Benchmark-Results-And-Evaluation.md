---
title: MAgHARCM Benchmark Results & Empirical Evaluation
date: 2026-09-08
last_updated: 2026-09-09 (iter-1, wave-24)
aliases:
  - "Benchmark-Results-And-Evaluation"
  - "Benchmark Results and Evaluation"
  - "Benchmark"
  - "Evaluation"
tags: [adhoc, benchmark, evaluation, slm, metrics, "[[2.0.0 MAgHARCM]]"]
---

# [[2.0.0 MAgHARCM Benchmark Results & Evaluation]]

> **Executive Overview**: Synthesizes empirical results from experimental trials ($K=3$) evaluating MAgHARCM across four real-world benchmark repositories spanning three source languages and three orders of magnitude in codebase size. Acronyms and domain concepts are defined in the [[Glossary|Domain Acronyms & Terminology Glossary]].

---

## 1. Primary Benchmark Results ($K=3$ Stochastic Trials)

| Benchmark | Language Pair | Source Files | Source LoC | Compilation | Test Pass Rate | Wall-Clock Time (s) | Convergence Status |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **GildedRose** | C to Rust | 4 | 199 | **Pass** | **14/14 (100.0%)** | $512.4 \pm 19.8$ | Converged (Iter 4) |
| **Gohistogram** | Go to Rust | 7 | 15,470 | **Pass** | **5/8 (62.5%)** | $645.8 \pm 28.4$ | Converged (Iter 6) |
| **Stats** | Go to Rust | 65 | 5,625 | **Pass** | **24/42 (57.1%)** | $920.3 \pm 41.5$ | Converged (Iter 8) |
| **Commons-Validator** | Java to Rust | 121 | 28,110 | **Fail** | **18/68 (26.5%)** | $1180.5 \pm 52.1$ | Plateau (Iter 12) |

---

## 2. Research Questions Analysis

### RQ1: Effectiveness of MAgHARCM
- **Syntactic Correctness**: Clean compilation achieved on 3 of 4 benchmark repositories. Small reasoning models (30B) successfully derive `Clone`, `Debug`, and `PartialEq` on translated structs, preventing borrow checker violations.
- **Functional Equivalence**: On GildedRose, MAgHARCM achieves a 100% test pass rate, correctly emitting `char`-typed function arguments where prior agents (e.g. SWE-agent) emit `int` and fail.
- **Large Repository Scaling**: On Stats (65 files), 24/42 tests pass without human intervention, demonstrating that reverse-topological planning linearizes complex call DAGs.

### RQ2: Test Suite Co-Translation & Synthesis
- When source tests are co-translated, the validator executes an AST to target compiler to test execution cascade.
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
- Average translation time for small projects ($<1$k LoC) is $\approx 8.5$ minutes; medium projects (5k–15k LoC) require $\approx 10–15$ minutes on standard consumer workstations (RTX 4090 / Apple Silicon).

---

## 4. Historical Audit Trail & Substrate Evolution (Waves 20–24)

### Empirical Baseline Disclaimer (Standing Rule)
Benchmark numbers reported in Section 1 reflect the verified Wave-17 empirical run ($K=3$). In subsequent sprints (Waves 20–23), empirical re-runs were deferred because:
1. **BLK-04**: The offline execution sandbox lacks a reachable GPU daemon (Ollama/vLLM) to perform fresh stochastic trials.
2. **BLK-02**: The Commons-Validator translation plateau (18/68 tests passing) requires dedicated regex/inheritance profile tuning prior to re-measurement.
3. Waves 20–23 primarily landed algorithmic substrates and verified literature anchors rather than empirical regressions.

An empirical rebase will be scheduled once BLK-04 resolves and the respective opt-in substrate gates are toggled.


### Wave-24 Empirical Probe (2026-09-09)

**BLK-04 partially resolved on this workstation**: the local Ollama daemon is reachable (`http://localhost:11434`, models `qwen3:30b-a3b-thinking-2507-q4_K_M` + `hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL` per the crust configs), so a bounded empirical probe ran instead of a deferral.

| Project | Run | Status | Artifact |
| :--- | :--- | :--- | :--- |
| 2dpartint (C to Rust) | crust runner, single trial | **Fail** — pipeline aborted at the `analyzer` node: `structured: unmarshal tool args: unexpected end of JSON input (raw=)` (model emitted neither a tool call nor content; no FINALSUM emitted) | `benchmarks/crust/results/768d4bd9f5120500a6909fb9acdc0d988dbb1422/2dpartint/.log` |

- **Runner fix landed first**: `benchmarks/crust/scripts/run.sh` passed the bare project name where the binary expected a config path; fixed in commit `aa75ccb` (`fix(benchmark): pass config path through crust runner`). The two pre-fix runs (`839c52ea…`, `6df14943…`) panicked at config load and were discarded; their result dirs were removed.
- **No Section-1 row changed**: the run produced no compilation or test-pass numbers, so the verified Wave-17 baseline in Section 1 stands unchanged.
- **Known bug surface for the next code phase** (located, not fixed this sprint): `internal/llm/structured.go::Extract` is single-shot — no retry on empty/unparseable model output, and the `AnalyzerSchema.libraries` nested field can arrive as a JSON string instead of an object (same qwen3-thinking JSON-as-content family fixed for `PlanningSchema` in commit `3f4249e`). A corrective-prompt retry loop on unmarshal failure is the recommended fix.

### Consolidated Substrate Evolution Matrix

| Wave | Sprint Date | Newly Anchored Substrates & Literature | Specified Target Config Gates (configs/agents.yml) | Re-run Condition |
| :--- | :--- | :--- | :--- | :--- |
| **Wave 20** | 2026-09-07 (iter-3) | - Graph Pre-Analysis (`[[1.0.0 P-142]]` NESA)<br>- Hallucination Defense (`[[1.0.0 P-143]]` HalluShield)<br>- Prompt-Trace Training (`[[1.0.0 P-144]]` TraceCoder)<br>- Strategy Selection (`[[1.0.0 P-145]]` TerraMod)<br>- Workflow PRM (`[[1.0.0 P-146]]` ContextPRM) | `comprehension.graph_self_evolving: true`<br>`agents.hallucination_defence: speculator_flight_recorder`<br>`agents.prompt_trace_training: true`<br>`agents.translation_strategy_selector: lexical_llm`<br>`agents.prm_workflow_context: true` | Re-run $K=3$ once BLK-04 resolves and any Wave-20 gate is toggled. |
| **Wave 21** | 2026-09-07 (iter-4) | - Draft KV Eviction (`[[1.0.0 P-147]]` SpecKV)<br>- LoRA-augmented Lookahead (`[[1.0.0 P-148]]` LookaheadKV)<br>- Async Speculative Pipeline (`[[1.0.0 P-149]]` SSD/Saguaro)<br>- Coverage Test Minimization (`[[1.0.0 P-150]]` TestPrune) | `agents.kv_cache.eviction.strategy: speckv`<br>`agents.kv_cache.eviction.strategy: lookaheadkv`<br>`agents.speculative.async_pipeline: true`<br>`agents.comprehension.observation.test_prune: true` | Re-run $K=3$ once BLK-04 resolves and any Wave-21 gate is toggled. |
| **Wave 22** | 2026-09-07 (iter-6) | - Feedback-Driven C-to-Rust (`[[1.0.0 P-151]]` SmartC2Rust)<br>- Hallucination Benchmark Triplet (`[[1.0.0 P-152]]` Hallu-Eval) | `agents.translation.feedback_driven: true`<br>`agents.comprehension.hallucination_evaluation: true` | Re-run $K=3$ once BLK-04 resolves and any Wave-22 gate is toggled. |
| **Wave 23** | 2026-09-08 (iter-1) | - Context-Aware Refinement Slicing (`[[1.0.0 P-153]]` CoReX)<br>- Critic-Feedback Multi-Agent (`[[1.0.0 P-154]]` TransAgent)<br>- Agentic Warning Classification + Repair (`[[1.0.0 P-155]]` CodeCureAgent)<br>- LLM Test-Generation (`[[1.0.0 P-156]]` TestWeaver) | `comprehension.graph_self_evolving: true` (shared with §11.6)<br>`translation.feedback_driven: true` (shared with SmartC2Rust) | Re-run $K=3$ once BLK-04 resolves and Wave-23 substrates are active. |
| **Wave 24** | 2026-09-09 (iter-1) | - Output-Reconstruction KV Eviction (`[[1.0.0 P-157]]` ReST-KV) | `agents.kv_cache.eviction.strategy: restkv` (hypothetical gate; not wired) | Re-run $K=3$ once the structured-output retry fix lands and any Wave-24 gate is toggled. |
