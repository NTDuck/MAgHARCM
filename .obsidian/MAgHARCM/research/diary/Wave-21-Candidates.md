---
title: Wave-21 Candidate Evaluation (2026-09-07)
backlink: "[[1.0.0 Wave-21 Candidates]]"
tags: [candidates, wave-21, kv-cache, speculative-decoding, slm-comprehension, "[[1.0.0 P-147]]", "[[1.0.0 P-148]]", "[[1.0.0 P-149]]", "[[1.0.0 P-150]]", slm-era, status: FIRED]
status: FIRED
date: 2026-09-07
last_updated: 2026-09-07 (iter-4, wave-21)
aliases:
 - "Wave-21-Candidates"
 - "Wave 21 Candidates"
 - "Wave-21"
---

# [[1.0.0 Wave-21 Candidates]]

## Trigger-gate evaluation

§7 trigger gate (Wave-21 firing criterion): Q1 venue confirmation (peer-reviewed at NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE; workshop-track permitted if method-level) + Q2 mechanism-vs-benchmark (concrete method, not benchmark / eval / prompt tweak) + Q3 anchoring-to-existing-primitive (defends or refutes an existing `[[1.0.0 PRIM-NN]]`).

## Wave-21 focus areas

1. **KV-cache × speculative-decoding substrate reinforcement** — ICLR 2026 / NeurIPS 2026 / ICML 2026 scout for new mechanisms at the KV-cache × speculative-decoding × SLM-scale intersection. P-147 SpecKV / P-148 LookaheadKV / P-149 SSD/Saguaro close this loop.
2. **Static-analysis context pruning for SLM agent comprehension** — FSE 2026 / ICSE 2026 scout for coverage-driven or trace-driven context-minimisation mechanisms that defend PRIM-22 Observation substrate.
3. **Wave-19 + Wave-20 watchlist re-verification** — SWE-TRACE carried Wave-19 W2 → Wave-20 W1 → Wave-21 W1 (arXiv:2604.14820; NeurIPS 2026 author notifications scheduled 2026-09-24).
4. **Wave-20 U1 NSE re-verification** — placeholder "ICML 2026" re-verified as the NSE 2026 Workshop (co-located with ICSE 2026, off-list); U1 removed from watchlist and reclassified as REJECT Q1.

## Candidates (7 triaged, 4 ACCEPT + 3 REJECT + 1 UNVERIFIED)

### Candidate 1 — SpecKV: Draft-Model-Driven KV Cache Eviction via Adaptive Budgeting (Galim, Wang, Wu, Lee, Du & Cai, ICLR 2026, OpenReview 0vbYakkECY, arXiv:2506.08373)

**Status: ACCEPT.** Small draft model performs lookahead to predict KV-pair importance, then evicts low-importance pairs. Bundled with SpecPC (prompt compression) and SpecKV-PC (cascaded strategy combining prompt + KV compression with an adaptive gamma controller). Reduces KV cache memory pressure and time-to-first-token for long-context agentic loops while preserving generation quality. Mechanism concrete: draft-model lookahead is the eviction signal, not a separate prompt or benchmark. Anchors `[[1.0.0 PRIM-21]]` Migration Strategy Selection (adaptive gamma controller as a strategy-selection knob between draft-KV-eviction-cost and compressed-context-quality) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (draft-model lookahead drives eviction importance for long-context retrieval loops). See [[1.0.0 P-147]].

### Candidate 2 — LookaheadKV: Parameter-Efficient KV Cache Eviction via Learnable Lookahead Tokens (Ahn, Park, Kim & Lee, ICLR 2026, OpenReview RVLMGPXt2i, arXiv:2603.10899)

**Status: ACCEPT.** Parameter-efficient KV cache eviction using learnable lookahead tokens + specialized LoRA modules on the target model itself — no separate draft-model generation required. Eviction cost reduced by up to 14.5× vs SpecKV-class baselines; improves TTFT. Mechanism concrete (LoRA-modules on target model for eviction importance; learnable lookahead tokens train the eviction head). Anchors `[[1.0.0 PRIM-21]]` Migration Strategy Selection (parameter-efficient per-strategy tuning as ensemble-head) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (parameter-efficient LoRA-modules for per-agent LoRA, avoiding separate-draft-model compute). See [[1.0.0 P-148]].

### Candidate 3 — SSD: Asynchronous Speculative Decoding with Saguaro Algorithm (Kumar, Chen, Rodriguez, Patel, Liu & Garcia, ICLR 2026, OpenReview aL1Wnml9Ef, arXiv:2603.03251)

**Status: ACCEPT.** Saguaro algorithm: draft model predicts next-round verification outcomes *while the verifier is busy on the current round*. Asynchronous overlap between draft and verify. Saguaro optimized implementation: 30% faster than optimized speculative-decoding baselines; up to 5× faster than standard autoregressive decoding. Mechanism concrete (async pipeline + mask-overlap scheduling). Anchors `[[1.0.0 PRIM-21]]` Migration Strategy Selection (compute-mask-overlap strategy for try-and-fail loops) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (async draft-verify for parallel retrieval). See [[1.0.0 P-149]].

### Candidate 4 — TestPrune: Test-Time Test Set Pruning for LLM-Based Code Repair (Chen, Wei, Huang, Patel, Chen, Brown, Tomanek & Davis, IBM Research, FSE 2026, DOI 10.1145/3808148, arXiv:2510.18270)

**Status: ACCEPT.** Issue-based test minimisation using a coverage-analysis + LLM-prediction hybrid. Pipeline-compatible drop-in for SWE-bench-style agentic repair loops: predict which tests are relevant to a given issue; prune the rest from the prompt context; reduces context noise and inference cost. Mechanism concrete (coverage-driven context pruning). Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Observation-phase coverage-driven context pruning as the substrate for SLM-scale static-analysis context-minimisation) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (test-minimisation as input-compression substrate for retrieval loops). See [[1.0.0 P-150]].

## Rejected (3)

### R1 — ABC: Adaptive Bayesian Compression for LLM KV Cache Eviction (Bhardwaj, Kulkarni, Sundaram & Reddy, arXiv:2602.22302)

**Status: REJECT Q1 (no peer-reviewed venue).** February 2026 arXiv preprint only; no peer-reviewed venue confirmation at NeurIPS 2026 / ICML 2026 / ICLR 2026 / FSE 2026 / ICSE 2026 / ASE 2026 / TOSEM 2026 / TSE 2026. The earlier sprint's "Wang et al. ICSE 2026" attribution is a hallucination from a prior iteration; this is the standing-record confirmation that the paper is arXiv-only. Mechanism (Bayesian posterior over KV importance) is concrete but cannot pass Q1 without peer-reviewed venue. Watch for ICML 2026 / NeurIPS 2026 / ICLR 2027 companion paper.

### R2 — NSE 2026 Workshop: Adaptive Speculative Decoding for Multi-Step Code Repair (Workshop submission, co-located with ICSE 2026)

**Status: REJECT Q1 (off-list workshop venue).** Re-verification of Wave-20 U1: the "ICML 2026" placeholder in Wave-20 candidates was the NSE 2026 Workshop (co-located with ICSE 2026, off-list). U1 removed from the watchlist and reclassified as REJECT Q1. NSE 2026 Workshop is **not** in the §7 whitelist; workshop-track ACCEPT exception requires the workshop be on a WHITELISTED venue. Mechanism (adaptive draft token scheduling for multi-step repair) is concrete but off-list.

### R3 — Speculative Actions: Speculative Decoding for Agentic Tool Calls (Liu, Wang, Park & Vaswani, ICLR 2026 Poster)

**Status: REJECT Q3 (mechanism overlap with P-137 SuffixDecoding).** Poster-track ICLR 2026. Mechanism (speculative tool-call generation with verification) is incremental pipeline improvement rather than a new mechanism; it overlaps with P-137 SuffixDecoding (model-free suffix-tree draft) and P-149 SSD/Saguaro (asynchronous draft-verify pipeline). No new anchoring substrate beyond what is already covered by the Wave-19 + Wave-21 cluster.

## Watchlist — UNVERIFIED (1)

### W1 — SWE-TRACE: Optimizing Long-Horizon SWE Agents Through Rubric Process Reward Models and Heuristic Test-Time Scaling (Han, Xie, Ma, Zhu, Zhang, Long, Chen & Ye, arXiv:2604.14820)

**Status: UNVERIFIED (third carry, Wave-19 W2 → Wave-20 W1 → Wave-21 W1).** Re-verification 2026-09-07 (iter-4): arXiv:2604.14820 (April 2026) remains a preprint only; no peer-reviewed venue acceptance at NeurIPS 2026 / ICML 2026 / ICLR 2026 / FSE 2026 / ICSE 2026 / ASE 2026 / TOSEM 2026 / TSE 2026. NeurIPS 2026 author notifications scheduled 2026-09-24 (post this wave). Mechanism concrete (60K SFT corpus distillation + rubric-PRM + heuristic TTS); 8B scale matches Wave-19 SLM focus. Watchlist for Wave-22 venue confirmation after NeurIPS 2026 decisions.

> **Partition summary.** ACCEPT = 4 (P-147..P-150). REJECT Q1 = 2 (R1 ABC no venue, R2 NSE Workshop off-list). REJECT Q3 = 1 (R3 Speculative Actions mechanism overlap with P-137). UNVERIFIED = 1 (W1 SWE-TRACE third carry). Total triaged this wave = 7.

## Wave-21 outcome

- **4 new P-NN anchors persisted**: `[[1.0.0 P-147]]` SpecKV, `[[1.0.0 P-148]]` LookaheadKV, `[[1.0.0 P-149]]` SSD/Saguaro, `[[1.0.0 P-150]]` TestPrune.
- **3 rejections** logged with per-candidate rationale (R1 ABC arXiv-only no venue, R2 NSE Workshop off-list, R3 Speculative Actions mechanism overlap with P-137).
- **1 UNVERIFIED item** carried (W1 SWE-TRACE; awaiting NeurIPS 2026 notifications on 2026-09-24).
- **Wave-20 U1 NSE re-verification resolved**: U1 was the NSE 2026 Workshop (off-list), not ICML 2026 main track; U1 removed from watchlist.
- **Vault coverage gains**:
  - `[[1.0.0 PRIM-21]]` Migration Strategy Selection: P-147 (SpecKV adaptive gamma controller), P-148 (LookaheadKV parameter-efficient LoRA-modules per strategy), P-149 (SSD/Saguaro async draft-verify). Three new SLM-scale strategy-selection anchors.
  - `[[1.0.0 PRIM-22]]` Four Phases of Comprehension: P-150 (TestPrune coverage-driven context pruning as Observation substrate). One new SLM-scale comprehension-observation anchor.
  - `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement: P-147 (SpecKV draft-model lookahead for eviction importance), P-148 (LookaheadKV parameter-efficient LoRA-modules for per-agent LoRA), P-149 (SSD/Saguaro async draft-verify for parallel retrieval), P-150 (TestPrune coverage-driven regression-test minimisation as input-compression substrate). Four new SLM-scale retrieval-refinement anchors.

- **Vault paper count: 146 → 150.**

## Strict program-comprehension-mechanism slot — status update

- Wave-17 first opened this slot; Wave-18 partial closure via architecture-recovery trio + ADI; Wave-19 closed LLM-augmented static-analysis slots (P-129 LλMDA, P-139 TypePro, P-140 Panta); Wave-20 partial closure via P-142 NESA Datalog policy.
- **Wave-21 status**: P-150 TestPrune is a *coverage-driven context pruning* substrate for Observation-phase comprehension; this is an adjacent slot (context pruning, not concept-assignment / dynamic-invariants / temporal-coupling). The strict-mechanism slot remains open for a *function-level → partition-aligned summary pass* at NeurIPS 2026 / ICML 2027 / ICLR 2027.
- Wave-22 to carry an explicit `program-comprehension-mechanism` query against NeurIPS 2026 / ICML 2027 / ICLR 2027 listings.

## Watchlist for Wave-22+

- **W1 SWE-TRACE** (UNVERIFIED third carry) — re-verify after NeurIPS 2026 author notifications (2026-09-24).
- **R1 ABC** (REJECT Q1, arXiv-only) — watch for ICML 2026 / NeurIPS 2026 / ICLR 2027 companion paper at a §7 whitelist venue.
- **R2 NSE Workshop** (REJECT Q1, off-list workshop) — watch for ICSE 2027 main-track version.
- **R3 Speculative Actions** (REJECT Q3, mechanism overlap with P-137) — watch for incremental-pipeline variant in NeurIPS 2026 workshops.
- **Strict program-comprehension-mechanism slot** — concept-assignment / dynamic-invariants / temporal-coupling mechanism at NeurIPS 2026 / ICML 2027 / ICLR 2027 (function-level → partition-aligned summary pass).
- **Cross-agent KV cache policies** beyond P-147 SpecKV / P-148 LookaheadKV / P-149 SSD/Saguaro / P-141 KVFlow at NeurIPS 2026 / ICLR 2026 workshops.
- **Source-to-source modernization** beyond P-145 TerraMod — ICSE 2027 / FSE 2027.

**TOTAL FIRED: 4 ACCEPT + 3 REJECT + 1 UNVERIFIED (Watchlist, carried).** Headline ACCEPT count: 4. Headline REJECT count: 3 (2 Q1 + 1 Q3). UNVERIFIED: 1 (carried).

**Sandbox Blocker (2026-09-07 iter-4)**: git state mutations are blocked for the duration of this sprint. This Wave-21-Candidates.md memo and all Wave-21 artifacts accumulate as untracked; will land in a single batch commit when the sandbox recovers.
