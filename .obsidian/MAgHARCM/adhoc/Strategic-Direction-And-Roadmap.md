---
title: MAgHARCM Strategic Direction & Future Roadmap
date: 2026-09-08
last_updated: 2026-09-09 (iter-3, wave-25)
aliases:
  - "Strategic-Direction-And-Roadmap"
  - "Strategic Direction and Roadmap"
  - "Roadmap"
  - "Direction"
tags: [adhoc, strategy, roadmap, future, slm, "[[2.0.0 MAgHARCM]]"]
---

# [[2.0.0 MAgHARCM Strategic Direction & Roadmap]]

> **Executive Overview**: Outlines the long-term architectural vision, edge-native SLM strategy, and phased modernization roadmap for MAgHARCM. For definitions of system acronyms (SLM, KV-cache, MaTTS, ACI, CPG), consult the [[Glossary|Domain Acronyms & Terminology Glossary]].

---

## 1. Core Strategic Thesis

### Edge-Native, Small-Language-Model (SLM) First (4B–30B)
- **Constraint**: Cloud frontier models (GPT-4o, Claude 3.5 Sonnet) are cost-prohibitive for large enterprise migrations ($100k+ LoC) and raise severe intellectual property and data privacy concerns for proprietary legacy codebases.
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
  ├── KV Cache chunk compression (KVPress, ChunkKV P-105; KVzip P-128 + RepairKV P-138 + KVFlow P-141)
  ├── Memory-augmented Test-Time Scaling (MaTTS in internal/memorystore)
  ├── SLM-as-Judge pipeline (T1 P-125 + SLM-as-a-Judge P-127; defended by HalluShield P-143)
  └── Program-comprehension mechanism substrate (CoReX P-153 closes residual slot; NESA P-142, ADI P-133)
         │
         ▼
Phase 3: Deep Enterprise Modernization (UPCOMING)
  ├── C/C++ to Safe Rust memory safety oracle (Syzygy P-124; SmartC2Rust P-151; TransAgent P-154)
  ├── Java Spring / Jakarta to Go microservices decomposition (Commons-Validator resolution)
  └── Automated dependency tree migration via UniAST (ABCoder P-14)
         │
         ▼
Phase 4: Autonomous Verified SWE-Bench Zero-Shot Modernizer (FUTURE)
  ├── Multi-SWE-bench cross-language benchmark evaluation
  └── Human-out-of-the-loop verified self-repair loop
```

---

## 3. Near-Term Initiatives

### 3.1. Core Architecture & Benchmark Modernization

1. **Resolve MemoryStore MaTTS Integration**:
   - Canonical `MaTTSDefaultBudget = 10` established in `internal/compiletime/compiletime.go` (BLK-01 resolved).
   - Wire persistent vector memory retrieval for cross-sprint strategy reuse.
2. **Expand Language Pairs Beyond C to Rust**:
   - **Primary**: Java to Rust (addressing the `Commons-Validator` plateau by engineering regex/inheritance mappings).
   - **Secondary**: Python 2 to Python 3 / Go (addressing legacy enterprise and scientific computing stacks).
3. **Formalize Benchmark Telemetry**:
   - Integrate automated wall-clock, memory, and token telemetry streaming into the Charm TUI dashboard.

---

### 3.2. Closed-Loop Substrate Implementations (Waves 20–23)

1. **Program-Comprehension Mechanism**:
   - Residual program-comprehension slot (open since Wave-17) is **closed** via `[[1.0.0 P-153]]` CoReX context-aware refinement slicing, joining `[[1.0.0 P-142]]` NESA self-evolving graph pre-analysis and `[[1.0.0 P-133]]` ADI frame lifetime traces.
2. **SLM-as-Judge Closed-Loop (`PRIM-7`)**:
   - Judgement (`[[1.0.0 P-125]]` T1 / `[[1.0.0 P-127]]` SLM-as-a-Judge) $\to$ Defence (`[[1.0.0 P-143]]` HalluShield) $\to$ Validation (`[[1.0.0 P-135]]` SPECS / `[[1.0.0 P-136]]` CaTS / `[[1.0.0 P-146]]` ContextPRM).
3. **Strategy-Selection Closed-Loop (`PRIM-21`)**:
   - Registry (`[[1.0.0 P-122]]` ReasoningBank) $\to$ Selection (`[[1.0.0 P-145]]` TerraMod / `[[1.0.0 P-123]]` CodeChemist) $\to$ Execution (`[[1.0.0 P-137]]` SuffixDecoding) $\to$ Judgement (`[[1.0.0 P-143]]` HalluShield).
4. **KV-Cache Eviction & Speculative Decoding (`PRIM-21, 31`)**:
   - Draft-model lookahead (`[[1.0.0 P-147]]` SpecKV) $\to$ LoRA-augmented target retention (`[[1.0.0 P-148]]` LookaheadKV) $\to$ Asynchronous verification pipeline (`[[1.0.0 P-149]]` SSD/Saguaro).
5. **Feedback-Driven Multi-Agent Translation (`PRIM-23, 29, 31`)**:
   - Context-aware chunking and three-signal loop (`[[1.0.0 P-151]]` SmartC2Rust) paired with multi-agent critic feedback (`[[1.0.0 P-154]]` TransAgent).
6. **Context Pruning & Test Synthesis (`PRIM-12, 22, 29`)**:
   - Coverage-driven test minimization (`[[1.0.0 P-150]]` TestPrune), agentic warning classification + repair (`[[1.0.0 P-155]]` CodeCureAgent), and execution-annotated backward slicing (`[[1.0.0 P-156]]` TestWeaver).

---

### 3.3. Wave-24 Status & Forward Strategic Objectives

Wave-24 fired 2026-09-09 and anchored `[[1.0.0 P-157]]` ReST-KV (ICLR 2026) for cross-agent KV-cache eviction (`PRIM-21, 31`), with 4 watchlist carries (SWE-TRACE, MemArt, MemDecay, ReCache) pending NeurIPS 2026 notifications (2026-09-24). Wave-25 (session 3, 2026-09-09) anchored `[[1.0.0 P-158]]` Agentic Rubrics (ACL 2026) for execution-free verification (`PRIM-7`) and retired SWE-TRACE at the sixth carry; ReCache/MemArt/MemDecay carry to Wave-26.

1. **Watchlist Triage**:
   - Re-verify W24-W1 SWE-TRACE (`arXiv:2604.14820`) after NeurIPS 2026 author notifications (2026-09-24); retire-or-confirm at Wave-25 if still lacking a peer-reviewed venue.
2. **Substrate Configuration Wiring**:
   - Wire target opt-in configuration paths (`configs/agents.yml:comprehension.graph_self_evolving`, `configs/agents.yml:translation.feedback_driven`) and verify with integration tests.
3. **Cross-Pattern Literature Surveys**:
   - Formalize cross-pattern studies among Wave-17..23 anchors (e.g. `P-129 × P-150`, `P-154 × P-151`, `P-155 × P-153`, `P-156 × P-140`) as general-purpose SLM modernization patterns.
4. **Empirical Benchmark Re-Measurement**:
   - Re-run the four benchmark repositories ($K=3$ stochastic trials) when local GPU inference (BLK-04) is restored.
