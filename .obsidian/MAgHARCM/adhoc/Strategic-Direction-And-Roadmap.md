---
title: MAgHARCM Strategic Direction & Future Roadmap
date: 2026-09-07
last_updated: 2026-09-07 (iter-6, wave-22)
aliases:
  - "Strategic-Direction-And-Roadmap"
  - "Strategic Direction and Roadmap"
  - "Roadmap"
  - "Direction"
tags: [adhoc, strategy, roadmap, future, slm, "[[2.0.0 MAgHARCM]]"]
---

# [[2.0.0 MAgHARCM Strategic Direction & Roadmap]]

> **Executive Overview**: Outlines the long-term vision, technological direction, and phased roadmap for MAgHARCM as an autonomous, edge-capable repository modernization platform.

---

## 1. Core Strategic Thesis

### Edge-Native, Small-Language-Model (SLM) First (4B–30B)
- **Constraint**: Cloud frontier models (GPT-4o, Claude 3.5 Sonnet) are cost-prohibitive for large enterprise migrations ($100k+ LoC) and raise severe IP/data privacy concerns for proprietary legacy systems.
- **Solution**: MAgHARCM is engineered to achieve frontier-grade repository translation using **local 4B–30B open models** (e.g. Qwen2.5-Coder:7B/32B `[[1.0.0 P-21]]`, StarCoder2:15B `[[1.0.0 P-22]]`, Phi-3-mini:3.8B `[[1.0.0 P-27]]`).
- **Mechanism**: Offset small-model capacity limits through **structured multi-agent decomposition**, **reverse-topological skeleton planning**, **formal CPG graph retrieval**, and **iterative compiler feedback cascades**.

### Test-Time Scaling ($S^*$, s1, ReasoningBank)
- Transitioning from simple greedy sampling to test-time search and reasoning over migration candidate solutions.
- Incorporating distilled self-reflection strategies (`[[1.0.0 P-122]]` ReasoningBank) and test-time synthesis scaling (`[[1.0.0 P-123]]` CodeChemist).

---

## 2. Phased Roadmap

```
Phase 1: Foundations & Parity (COMPLETED)
  ├── 31 Primitives implemented & wired
  ├── Eino 8-agent execution graph
  └── Automated vault linting (scripts/lint_vault.sh)
         │
         ▼
Phase 2: Local SLM Fine-Tuning & Prompt Specialization (CURRENT)
  ├── Structured output cloze slots (PRIM-22, PRIM-24)
  ├── KV Cache chunk compression (KVPress, ChunkKV P-105; KVzip P-128 + RepairKV P-138 + KVFlow P-141 for agent-loop substrate)
  ├── Memory-augmented Test-Time Scaling (MaTTS in internal/memorystore)
  ├── SLM-as-Judge pipeline (T1 P-125 + SLM-as-a-Judge P-127; defended by HalluShield P-143 + validated by SPECS P-135 + CaTS P-136 + ContextPRM P-146)
  └── Program-comprehension mechanism substrate (NESA P-142 + SSAR/SemArc/SemRef P-130..P-132 + ADI P-133; §11.6)
         │
         ▼
Phase 3: Deep Enterprise Modernization (UPCOMING)
  ├── C/C++ to Safe Rust memory safety oracle (Syzygy P-124)
  ├── Java Spring / Jakarta to Go microservices decomposition
  └── Automated dependency tree migration via UniAST (ABCoder P-14)
         │
         ▼
Phase 4: Autonomous Verified SWE-Bench Zero-Shot Modernizer (FUTURE)
  ├── Multi-SWE-bench cross-language benchmark evaluation
  └── Human-out-of-the-loop verified self-repair loop
```

---

## 3. Near-Term Initiatives

1. **Resolve MemoryStore MaTTS Integration**:
   - Establish the canonical `MaTTSDefaultBudget` configuration in `compiletime` to restore `memorystore` buildability.
   - Wire vector memory retrieval for cross-sprint strategy reuse.
2. **Expand Language Pairs Beyond C $	o$ Rust**:
   - Primary: Java $	o$ Rust (addressing the `Commons-Validator` plateau).
   - Secondary: Python 2 $	o$ Python 3 / Go (addressing scientific computing stacks).
3. **Formalize Benchmark Telemetry**:
   - Integrate automated wall-clock, memory, and token telemetry into the Charm TUI dashboard.
4. **Wave-20 program-comprehension-mechanism substrate**:
   - The strict program-comprehension-mechanism slot (open since Wave-17) is partially closed by `[[1.0.0 P-142]]` NESA self-evolving graph pre-analysis. The substrate is opt-in via `configs/agents.yml:comprehension.graph_self_evolving: true`.
   - Wave-21 priority: function-level → partition-aligned summary pass (residual gap); re-scout NeurIPS 2026 / ICML 2027 / ICLR 2027.
5. **Wave-20 SLM-as-Judge closed-loop**:
   - PRIM-7 substrate is now closed-loop: judgement (`[[1.0.0 P-125]]` T1 / `[[1.0.0 P-127]]` SLM-as-a-Judge) → defence (`[[1.0.0 P-143]]` HalluShield) → validation (`[[1.0.0 P-135]]` SPECS / `[[1.0.0 P-136]]` CaTS / `[[1.0.0 P-146]]` ContextPRM).
6. **Wave-20 PRIM-21 strategy-selection closed-loop**:
   - PRIM-21 substrate is now closed-loop: try-and-fail registry (`[[1.0.0 P-122]]` ReasoningBank) → strategy selection (`[[1.0.0 P-145]]` TerraMod / `[[1.0.0 P-123]]` CodeChemist) → execution (`[[1.0.0 P-137]]` SuffixDecoding) → judgement (`[[1.0.0 P-143]]` HalluShield).
7. **Wave-20 prompt-trace training data**:
   - `[[1.0.0 P-144]]` TraceCoder extends the PRIM-29 + PRIM-31 substrate with prompt-trace training-data construction.
8. **Wave-21 KV-cache × speculative-decoding substrate**:
   - PRIM-21 / PRIM-31 substrate is now closed-loop for KV-cache eviction + speculative-decoding: draft-model-driven lookahead (`[[1.0.0 P-147]]` SpecKV) → parameter-efficient LoRA-modules on target model (`[[1.0.0 P-148]]` LookaheadKV) → asynchronous draft-verify pipeline (`[[1.0.0 P-149]]` SSD/Saguaro). All three are opt-in via `configs/agents.yml:agents.kv_cache.eviction.strategy` / `agents.speculative.async_pipeline`.
9. **Wave-21 coverage-driven context pruning**:
   - `[[1.0.0 P-150]]` TestPrune extends PRIM-22 Observation-phase substrate with coverage-driven test minimisation. The substrate is opt-in via `configs/agents.yml:agents.comprehension.observation.test_prune: true`.
10. **Wave-21 forward to Wave-22**:
   - The residual strict-mechanism program-comprehension-mechanism slot (function-level → partition-aligned summary pass) carries to Wave-22; explicit query against NeurIPS 2026 / ICML 2027 / ICLR 2027 listings.
   - W1 SWE-TRACE re-verification after NeurIPS 2026 author notifications (2026-09-24); fourth carry if no peer-reviewed venue.
   - Wave-22 to integrate SpecKV's draft-model KV-eviction path into the existing benchmark harness and re-run the four repositories if BLK-04 (GPU/LLM endpoint) resolves.

11. **Wave-22 feedback-driven C-to-Rust translation substrate**:
   - `[[1.0.0 P-151]]` SmartC2Rust (ICSE 2026, DOI 10.1145/3744916.3773259) extends the Phase-3 C/C++ → Safe Rust oracle with iterative context-aware segmentation + three-signal feedback loop (target compiler errors, semantic-equivalence diffs, residual unsafe-block counts). The C-to-Rust analogue of `[[1.0.0 P-124]]` Syzygy's Go-to-Rust three-signal loop. Anchors `PRIM-23` Chunked Translation, `PRIM-29` Dynamic Iteration Recruiter, `PRIM-31` Iterative Retrieval. The substrate is opt-in via `configs/agents.yml:translation.feedback_driven: true`.
12. **Wave-22 hallucination-evaluation triplet substrate**:
   - `[[1.0.0 P-152]]` Hallu-Eval / Hallu-Shield (FSE 2026, DOI 10.1145/3808189) closes the systematic hallucination-evaluation slot for SLM-scale comprehension + translation. Triplet = (Hallu-Eval 800-pair benchmark with natural + induced logical hallucinations) + (Hallu-Det detection approach) + (Hallu-Shield inference-time mitigation). Anchors `PRIM-22` Four Phases Comprehension, `PRIM-25` Role-Flip De-Hallucination. The substrate is opt-in via `configs/agents.yml:comprehension.hallucination_evaluation: true`.
13. **Wave-22 forward to Wave-23**:
   - The residual strict-mechanism program-comprehension-mechanism slot (function-level → partition-aligned summary pass) carries to Wave-23; explicit query against NeurIPS 2026 / ICML 2027 / ICLR 2027 listings.
   - W1 SWE-TRACE re-verification after NeurIPS 2026 author notifications (2026-09-24); fifth carry if no peer-reviewed venue, then retire.
   - Wave-23 to integrate SmartC2Rust's three-signal feedback loop + Hallu-Eval's hallucination-evaluation triplet into the benchmark harness and re-run the four repositories if BLK-04 (GPU/LLM endpoint) resolves.
