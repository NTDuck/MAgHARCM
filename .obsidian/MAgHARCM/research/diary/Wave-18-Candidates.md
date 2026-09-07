---
title: Wave-18 Candidate Evaluation (2026-09-07)
backlink: "[[1.0.0 Wave-18 Candidates]]"
tags: [candidates, wave-18, mechanism-vs-benchmark, "[[1.0.0 P-125]]", "[[1.0.0 P-126]]", "[[1.0.0 P-127]]", "[[1.0.0 P-128]]", "[[1.0.0 P-129]]", "[[1.0.0 P-130]]", "[[1.0.0 P-131]]", "[[1.0.0 P-132]]", "[[1.0.0 P-133]]", test-time-scaling, kv-cache, software-archaeology, speculative-decoding, slm-as-judge, dynamic-analysis, architecture-recovery]
status: FIRED
date: 2026-09-07
last_updated: 2026-09-07
---

# [[1.0.0 Wave-18 Candidates]]

## Trigger-gate evaluation

§7 trigger gate (rewritten 2026-09-25): Q1 venue confirmation (peer-reviewed at NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE; workshop-track permitted if method-level) + Q2 mechanism-vs-benchmark (concrete method, not benchmark / eval / prompt tweak) + Q3 anchoring-to-existing-primitive (defends or refutes an existing `[[1.0.0 PRIM-NN]]`).

## Candidates (13 triaged, 9 ACCEPT + 4 REJECT)

### Candidate 1 — T1: Tool-Integrated Verification for Test-time Compute Scaling in SLMs (Kang, Jeong & Cho, ICLR 2026)

**Status: ACCEPT.** (Per scout `ResearchVerdictValidation`.) See [[1.0.0 P-125]].

### Candidate 2 — ARC-Decode: Risk-Bounded Acceptance for Speculative Decoding (Li et al., ICML 2026)

**Status: ACCEPT.** (Per scout `ResearchSpeculativeDecoding`.) See [[1.0.0 P-126]].

### Candidate 3 — SLM-as-a-Judge for Code Generation (Crupi et al., ICSE 2026)

**Status: ACCEPT.** (Per scout `ResearchVerdictValidation`.) See [[1.0.0 P-127]].

### Candidate 4 — KVzip: Query-Agnostic KV Cache Compression with Context Reconstruction (Kim et al., NeurIPS 2025 Oral)

**Status: ACCEPT.** (Per scout `ResearchKVCache`.) See [[1.0.0 P-128]].

### Candidate 5 — L\*λ\*MDA: LLM-Aided Partial Program Dependence Analysis (Rong, Yadavally & Nguyen, ICSE 2026)

**Status: ACCEPT.** (Per scout `ResearchSoftwareArchaeology`.) See [[1.0.0 P-129]].

### Candidate 6 — SSAR: Software Architecture Recovery (Ding, Mo, Wu & Song, ICSE 2026)

**Status: ACCEPT.** (Per scout `ResearchSoftwareArchaeology`.) See [[1.0.0 P-130]].

### Candidate 7 — SemArc: Software Architecture Recovery Augmented with Semantics (Zhao et al., TSE 2026)

**Status: ACCEPT.** (Per scout `ResearchSoftwareArchaeology`.) See [[1.0.0 P-131]].

### Candidate 8 — SemRef: Semantic-Enhanced Refinement of Architecture Recovery (Zhang et al., ICSE 2026)

**Status: ACCEPT.** (Per scout `ResearchSoftwareArchaeology`.) See [[1.0.0 P-132]].

### Candidate 9 — ADI: Empowering Autonomous Debugging Agents with Efficient Dynamic Analysis (Xiang et al., FSE 2026 — SIGSOFT Distinguished Paper Award)

**Status: ACCEPT (borderline).** (Per scout `ResearchDynamicAnalysis`.) Borderline because the focus list pairs "dynamic invariants" — ADI is dynamic-analysis-for-debugging-agents which fits. See [[1.0.0 P-133]].

## Rejected (4)

### R1 — ReflexiCoder (Jiang et al., ACL 2026 Findings)

**Status: REJECT Q1.** ACL Findings is off-list; same rejection rationale as Wave-17 MemSearcher. Mechanism is real (RL-internalised self-reflection at 1.5B-14B) but venue fails Q1. Watchlist for NeurIPS 2026 / ICML 2027 method-level companion.

### R2 — Self-Distillation for Code Generation (Zhang et al., Apple, arXiv:2604.01193)

**Status: REJECT Q1.** arXiv-only (April 2026); no confirmed peer-reviewed acceptance. Strong mechanism (30B 42.4% → 55.3% pass@1 on LiveCodeBench v6). Watchlist for NeurIPS 2026 acceptance.

### R3 — SPECS: Faster Test-Time Scaling via Speculative Drafts and Dynamic Switching (Cemri et al., arXiv:2506.15733)

**Status: REJECT Q1.** ICLR 2026 submission; acceptance status as of 2026-09-07 unconfirmed in the venue whitelist. Watchlist for ICLR 2026 proceedings integration.

### R4 — Software-Archaeology / Program-Comprehension Mechanism Gap

**Status: REJECT (gap confirmed for Wave-18).** FSE 2026 (June 2026) and ICSE 2026 (April 2026) had multiple LLM4Code papers but none introduces a *new* program-comprehension or software-archaeology mechanism at the SLM scale. Closest candidates (FSE 2026 StackRepoQA, CoReX, Slicer4J successors) are benchmarks / empirical studies, not mechanisms. **However**, this gap is *partially closed* by Wave-18 ACCEPT #5-9 (L\*λ\*MDA, SSAR, SemArc, SemRef, ADI), which operationalize structural+semantic hybrid code graphs at the LLM-augmented program-analysis level. The strict program-comprehension-mechanism gap remains; the *LLM-augmented software-archaeology* gap is filled.

## Wave-18 outcome

- **9 new P-NN anchors persisted**: `[[1.0.0 P-125]]` T1, `[[1.0.0 P-126]]` ARC-Decode, `[[1.0.0 P-127]]` SLM-as-a-Judge, `[[1.0.0 P-128]]` KVzip, `[[1.0.0 P-129]]` L\*λ\*MDA, `[[1.0.0 P-130]]` SSAR, `[[1.0.0 P-131]]` SemArc, `[[1.0.0 P-132]]` SemRef, `[[1.0.0 P-133]]` ADI.
- **4 rejections** logged with per-candidate rationale.
- **Vault coverage gains**:
  - PRIM-7 Verdict Validation: P-125 (T1), P-126 (ARC-Decode), P-127 (SLM-as-a-Judge) — three SLM-scale replacement anchors for the frontier-PRM pattern (`[[1.0.0 P-92]]`).
  - PRIM-9 Tri-Representation Hybrid Code Graph: P-129 (L\*λ\*MDA), P-130 (SSAR), P-131 (SemArc), P-132 (SemRef) — four structural+semantic hybrid code graph substrate anchors.
  - PRIM-21 Migration Strategy Selection: P-126 (ARC-Decode risk-bounded acceptance).
  - PRIM-22 Four Phases of Comprehension: P-128 (KVzip), P-129, P-130, P-131, P-133.
  - PRIM-31 Iterative Retrieval Refinement: P-128 (KVzip), P-132 (SemRef), P-133 (ADI).
- **Vault paper count: 124 → 133.**

## Watchlist for Wave-19+

- Software-archaeology *strict-mechanism* gap at NeurIPS 2026 / ICML 2027 / ICLR 2027.
- KV-cache compression mechanisms (e.g., RepairKV, RelayCaching, STAR-KV) for SLM-agent multi-turn workloads.
- ReflexiCoder if a NeurIPS 2026 / ICML 2027 method-level companion lands.
- SSD (Apple) if NeurIPS 2026 acceptance lands.
- SPECS if ICLR 2026 acceptance lands (proceedings integration).
- Distillation-style self-reflection mechanisms beyond RL (trajectory-level distillation, error-aware SFT) at 2025+ venues.

**TOTAL FIRED: 9 ACCEPT + 4 REJECT.**
