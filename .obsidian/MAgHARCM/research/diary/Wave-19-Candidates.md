---
title: Wave-19 Candidate Evaluation (2026-09-07)
backlink: "[[1.0.0 Wave-19 Candidates]]"
tags: [candidates, wave-19, mechanism-vs-benchmark, "[[1.0.0 P-134]]", "[[1.0.0 P-135]]", "[[1.0.0 P-136]]", "[[1.0.0 P-137]]", "[[1.0.0 P-138]]", "[[1.0.0 P-139]]", "[[1.0.0 P-140]]", "[[1.0.0 P-141]]", speculative-decoding, kv-cache, test-time-scaling, slm-as-judge, software-archaeology, dynamic-analysis, multi-agent-workflow, system-dependence-graph, status: FIRED]
status: FIRED
date: 2026-09-07
last_updated: 2026-09-07
---

# [[1.0.0 Wave-19 Candidates]]

## Trigger-gate evaluation

§7 trigger gate (rewritten 2026-09-25): Q1 venue confirmation (peer-reviewed at NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE; workshop-track permitted if method-level) + Q2 mechanism-vs-benchmark (concrete method, not benchmark / eval / prompt tweak) + Q3 anchoring-to-existing-primitive (defends or refutes an existing `[[1.0.0 PRIM-NN]]`).

## Wave-19 focus areas

1. **Speculative-decoding × KV-cache compression hybrids** (combining EAGLE-3 / ARC-Decode with KVzip / StreamingLLM / ChunkKV / suffix-tree drafts).
2. **Test-time-scaling search at SLM 4B-30B scale** (TTS + verification, TTS + tool-use, TTS + retrieval, calibrated single-pass confidence).
3. **LLM-augmented static analysis beyond dependence graphs** (inter-procedural slicing, dynamic invariants, control-flow recovery, agentic hybrid program analysis).

## Candidates (13 triaged, 8 ACCEPT + 3 REJECT + 2 UNVERIFIED)

### Candidate 1 — RelayCaching: Accelerating LLM Collaboration via Decoding KV Cache Reuse (Geng, Gao, Wu, Liu & Liu, ICML 2026 Poster #1915)

**Status: ACCEPT.** Cross-agent KV cache reuse via sparse deviation recompute (>80% reuse, 4.7× TTFT reduction). Anchors PRIM-31 + PRIM-9 by extending single-query compression to multi-agent prefill reuse on the 8-agent graph. See [[1.0.0 P-134]].

### Candidate 2 — SPECS: Faster Test-Time Scaling through Speculative Drafts (Cemri et al., ICLR 2026)

**Status: ACCEPT.** Speculative drafts + reward-guided soft verification + dynamic draft↔target switching. Closes Wave-18 R3 watchlist. Anchors PRIM-7 + PRIM-21 with continuous budget knob. See [[1.0.0 P-135]].

### Candidate 3 — CaTS: Calibrated Test-Time Scaling for Efficient LLM Reasoning (Huang, Huang, Leng, Liu & Huang, ICLR 2026 Poster)

**Status: ACCEPT.** Self-Calibration distills self-consistency confidence into the SLM itself; CaTS-ES (early-stopping TTS) + CaTS-BoN (adaptive best-of-N). 4B-30B scale (LLaMA3, Qwen2.5, Mistral) matches MAgHARCM substrate. Anchors PRIM-7 as a frontier-PRM-replacement primitive. See [[1.0.0 P-136]].

### Candidate 4 — SuffixDecoding: Extreme Speculative Decoding for Emerging AI Applications (Oliaro, Jia, Campos & Qiao, NeurIPS 2025 Spotlight)

**Status: ACCEPT.** Model-free speculative decoding using suffix trees over prompt + previous output tokens; deployed in Snowflake ArcticInference and vLLM. Up to 5.3× over vanilla, 2.8× over EAGLE-2/3 on AgenticSQL. Anchors PRIM-7 + PRIM-21 by exploiting template repetitiveness in agent loops. See [[1.0.0 P-137]].

### Candidate 5 — RepairKV (Cache You Later): Post-Compression KV Repair for Long-Context Agentic LLM Inference (Rusli, Paliwal, Zhang & Jiao, ICML 2026 AdaptFM Workshop)

**Status: ACCEPT (borderline, workshop-track).** §7 workshop-track ACCEPT rationale: (a) §7 method-level threshold met — concrete post-compression KV repair mechanism, not a benchmark; (b) no main-track alternative covers the same axis within Wave-19 focus area (a) — KVzip (P-128) targets single-query compression, KVFlow (P-141) targets workflow eviction, RelayCaching (P-134) targets cross-agent prefill reuse; P-138 alone covers post-compression KV repair for long-context agentic retrieval (91.0% vs 24.5% no-repair baseline at 32K on Qwen2.5-7B). Anchors PRIM-22 + PRIM-31 by giving multi-turn retrieval a runtime correction substrate. See [[1.0.0 P-138]].

### Candidate 6 — TypePro: Boosting LLM-Based Type Inference via Inter-Procedural Slicing (Lin, Fan, Huang, Shen & Wu, FSE 2026 Research Papers)

**Status: ACCEPT.** System Dependence Graph + inter-procedural backward/forward slicing + structural-similarity candidate-type selector. 88.9% Top-1 EM on ManyTypes4Py (+7.1 pp); 86.6% on ManyTypes4TypeScript (+10.3 pp). Anchors PRIM-9 (extends CPG with slicing as fourth representation) + PRIM-22 (type-inference substrate for Recognition/Explanation). See [[1.0.0 P-139]].

### Candidate 7 — Panta: LLM Test Generation via Iterative Hybrid Program Analysis (Gu, Nashid & Mesbah, ICSE 2026 Research Track)

**Status: ACCEPT.** Iterative hybrid static (cyclomatic-complexity path ranking) + dynamic (instrumented coverage feedback) program analysis drives LLM test generation. 26% higher line coverage, 23% higher branch coverage than SOTA. Anchors PRIM-22 + PRIM-21 by giving validator a coverage-driven test-prioritisation knob. See [[1.0.0 P-140]].

### Candidate 8 — KVFlow: Efficient Prefix Caching for Multi-Agent LLM Serving via Workflow-Aware Eviction (NeurIPS 2025 Poster)

**Status: ACCEPT.** Workflow-aware eviction using Agent Step Graph + "steps-to-execution" metric + KV prefetching overlapping CPU→GPU transfer. Up to 1.83× single-workflow, 2.19× concurrent vs SGLang. Anchors PRIM-31 + PRIM-21 by extending Wave-18 KVzip to multi-agent workflow retention. See [[1.0.0 P-141]].

## Rejected (3)

### R1 — TTA*: Test-Time A* Search for Multistep Reasoning in Small Language Models (Braverman, Zhang & Gu, NeurIPS 2025 LAW Workshop)

**Status: REJECT Q1 (workshop redundancy).** NeurIPS 2025 LAW workshop at NeurIPS 2025 is workshop-track; per §7 workshop-track is permitted only if main-track alternative is unavailable. P-135 SPECS and P-136 CaTS cover the same focus area at ICLR 2026 main track, so TTA*'s workshop venue makes it redundant. Mechanism (A* search wrapper) concrete but adds little beyond s1 (P-84) beam-search lineage. Watchlist for NeurIPS 2026 main-track companion.

### R2 — HELIOS: Hierarchical Graph Abstraction for Structure-Aware LLM Decompilation (Achamyeleh, Thomare & Al Faruque, NDSS 2026 LAST-X Workshop)

**Status: REJECT Q1 (off-list venue + off-axis target).** NDSS LAST-X workshop is off-list (§7 whitelist excludes NDSS workshop). Mechanism concrete (CFG+FCG → textual prompt; compilability 45.0%→85.2% on Gemini 2.0) but binary-decompilation target is off-axis for MAgHARCM's source→source translation. Watchlist only.

### R3 — LongSpec: Long-Context Lossless Speculative Decoding (ACL 2026 Main)

**Status: REJECT Q1 (off-list venue).** ACL is not in the Wave-19 whitelist (NeurIPS/ICML/ICLR/ICSE/FSE/ESEC/ASE/TOSEM/TSE). Same rationale as Wave-18 R1 ReflexiCoder. Mechanism fits Wave-19 (a) precisely (memory-efficient constant-KV-cache draft + novel position indices for 32K+ contexts); venue is the only blocker. Watchlist for ICLR/NeurIPS companion.

## Watchlist — UNVERIFIED (2)

### W1 — SliceMate: Accurate and Scalable Static Program Slicing via LLM-Powered Agents (Chang, Shi, Lyu, Zhou, Wang, Yang, Li & Lo, arXiv:2507.18957)

**Status: UNVERIFIED.** Yunbo Lyu's homepage claims ISSTA 2026 acceptance, but no ISSTA 2026 program slot for SliceMate is listed on conf.researchr.org or conference-publishing.com/authors/ISSTA26. Mechanism is concrete (three LLM agents replace explicit PDG/SDG construction; 22% acc / 28% F1 improvement). Method-level relevance is high (Wave-19 focus area c), but venue cannot be confirmed. Watchlist for Wave-20 venue re-verification.

### W2 — SWE-TRACE: Optimizing Long-Horizon SWE Agents Through Rubric PRMs and Heuristic Test-Time Scaling (Han, Xie, Ma, Zhu, Zhang, Long, Chen & Ye, arXiv:2604.14820)

**Status: UNVERIFIED.** April 2026 arXiv preprint only; no peer-reviewed venue confirmation at NeurIPS/ICML/ICLR/FSE/ICSE/ASE/TOSEM as of 2026-09-07. Mechanism concrete (cascaded trajectory optimisation + rubric-PRM + heuristic TTS) and 8B scale matches Wave-19 SLM focus. Watchlist for Wave-20 venue confirmation.
> **Partition summary.** ACCEPT = 8 (P-134..P-141). REJECT Q1 = 3 (R1 TTA*, R2 HELIOS, R3 LongSpec). UNVERIFIED = 2 (W1 SliceMate, W2 SWE-TRACE). Total triaged = 13. Watchlist carries UNVERIFIED items only; off-list-venue items live in REJECT.

## Wave-19 outcome

- **8 new P-NN anchors persisted**: `[[1.0.0 P-134]]` RelayCaching, `[[1.0.0 P-135]]` SPECS, `[[1.0.0 P-136]]` CaTS, `[[1.0.0 P-137]]` SuffixDecoding, `[[1.0.0 P-138]]` RepairKV, `[[1.0.0 P-139]]` TypePro, `[[1.0.0 P-140]]` Panta, `[[1.0.0 P-141]]` KVFlow.
- **3 rejections** logged with per-candidate rationale (R1 TTA*, R2 HELIOS, R3 LongSpec).
- **2 UNVERIFIED items** logged for Wave-20 venue re-verification (W1 SliceMate, W2 SWE-TRACE).
- **Vault coverage gains**:
  - PRIM-7 Multi-Agent Verdict Validation: P-135 (SPECS), P-136 (CaTS), P-137 (SuffixDecoding). Three 2025-2026 frontier-PRM-replacement anchors. (P-138 RepairKV is borderline workshop-track; credited to PRIM-22 + PRIM-31 only.)
  - PRIM-9 Tri-Representation Hybrid Code Graph: P-134 (RelayCaching), P-139 (TypePro). Two inter-agent + inter-procedural extensions.
  - PRIM-21 Migration Strategy Selection: P-135 (SPECS), P-137 (SuffixDecoding), P-140 (Panta), P-141 (KVFlow). Four draft / coverage / scheduling budget knobs.
  - PRIM-22 Four Phases of Comprehension: P-138 (RepairKV), P-139 (TypePro), P-140 (Panta). Three comprehension-substrate anchors at long-context / type-inference / coverage scale.
  - PRIM-31 Iterative Retrieval Refinement: P-134 (RelayCaching), P-138 (RepairKV), P-141 (KVFlow). Three multi-turn / multi-agent correction anchors.
- **Vault paper count: 133 → 141.**

## Watchlist for Wave-20+

- **W1 SliceMate** (UNVERIFIED) — re-verify ISSTA 2026 venue; mechanism is high-relevance for Wave-19 (c) if accepted.
- **W2 SWE-TRACE** (UNVERIFIED) — re-verify NeurIPS/ICML 2026 acceptance; rubric-PRM is a strong PRIM-7 anchor candidate.
- **R3 LongSpec** (REJECT Q1, off-list venue) — re-verify ICLR/NeurIPS companion; constant-KV-cache draft fits PRIM-22 long-context substrate.
- **R1 TTA\*** (workshop) — watch for NeurIPS 2026 main-track version; A* TTS extends P-84 s1 lineage.
- **R2 HELIOS** (NDSS workshop) — off-axis target; revisit only if a source-to-source variant appears.
- Software-archaeology *strict-mechanism* gap at NeurIPS 2026 / ICML 2027 / ICLR 2027 still open (carried from Wave-18 R4).
- Cross-agent KV cache policies beyond RelayCaching/KVFlow at NeurIPS 2026 / ICLR 2027.

**TOTAL FIRED: 8 ACCEPT + 3 REJECT + 2 UNVERIFIED (Watchlist).** Headline ACCEPT count: 8. Headline REJECT count: 3. UNVERIFIED: 2.
