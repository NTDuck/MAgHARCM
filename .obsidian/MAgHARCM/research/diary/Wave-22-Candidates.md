---
title: Wave-22 Candidate Evaluation (2026-09-07)
backlink: "[[1.0.0 Wave-22 Candidates]]"
tags: [candidates, wave-22, icse-2026, fse-2026, iterative-translation, hallucination-benchmark, "[[1.0.0 P-151]]", "[[1.0.0 P-152]]", slm-era, status: FIRED]
status: FIRED
date: 2026-09-07
last_updated: 2026-09-07 (iter-6, wave-22)
aliases:
 - "Wave-22-Candidates"
 - "Wave 22 Candidates"
 - "Wave-22"
---

# [[1.0.0 Wave-22 Candidates]]

## Trigger-gate evaluation

§7 trigger gate (Wave-22 firing criterion): Q1 venue confirmation (peer-reviewed at NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE; workshop-track permitted ONLY if method-level single-paper) + Q2 mechanism-vs-benchmark (concrete method, not benchmark / eval / prompt tweak alone) + Q3 anchoring-to-existing-primitive (defends or refutes an existing `[[1.0.0 PRIM-NN]]`).

## Wave-22 focus areas

1. **C-to-Rust iterative translation substrate** — ICSE 2026 / FSE 2026 / ASE 2026 scout for LLM-driven iterative translation feedback mechanisms. P-151 SmartC2Rust closes this loop (the C-to-Rust analogue of P-124 Syzygy Go-to-Rust).
2. **Code-summarization hallucination benchmark + mitigation** — FSE 2026 scout for systematic hallucination benchmarks that defend PRIM-22 (comprehension) and PRIM-25 (de-hallucination). P-152 Hallu-Eval + Hallu-Det + Hallu-Shield close this slot.
3. **NeurIPS 2026 / ICML 2027 / ICLR 2027 strict program-comprehension-mechanism query** — venue sweep confirms all three venues are pre-notification as of 2026-09-07 (NeurIPS 2026 notifications 2026-09-24, ICML 2027 submissions Jan 22 2027, ICLR 2027 submissions 2026-09-25). Strict slot remains open; Wave-23 must re-query after NeurIPS 2026 author notifications.
4. **Wave-21 W1 SWE-TRACE re-verification** — fourth carry; arXiv:2604.14820 (April 2026) still preprint-only; NeurIPS 2026 notifications 2026-09-24 still 17 days out. Carry to Wave-23.

## Candidates (5 triaged, 2 ACCEPT + 3 REJECT + 1 UNVERIFIED)

### Candidate 1 — SmartC2Rust: Iterative, Feedback-Driven C-to-Rust Translation via Large Language Models (Shiraishi, Cao & Shinagawa, ICSE 2026, DOI 10.1145/3744916.3773259, arXiv:2409.10506)

**Status: ACCEPT.** ICSE 2026 (April 2026, Rio de Janeiro) — flagship SE venue, in-list per §7. Mechanism is concrete: (a) **context-aware code segmentation** divides large-scale C source into context-bounded units to overcome LLM context-window limitations; (b) **iterative feedback loop** feeds three orthogonal signals back into the LLM — Rust compiler errors (lexical/syntactic feedback), semantic-equivalence diffs (semantic feedback), and `unsafe`-block counts (security feedback); (c) the LLM rewrites, the compiler validates, the loop repeats until `unsafe` is minimised and semantic equivalence holds. Reduces unsafe-block residuals and achieves higher semantic equivalence than prior C-to-Rust translation baselines (Syzygy P-124 was Go-to-Rust; SmartC2Rust closes the analogous C-to-Rust slot). Anchors `[[1.0.0 PRIM-23]]` Chunked Translation (context-aware C segmentation = chunking substrate), `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter (iterative feedback loop = recruiter-driven repair cycle), and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (compiler-feedback + unsafe-residual feedback = retrieval-refinement substrate). Citation hops: Hop 1 (mechanism) — SmartC2Rust itself, ICSE 2026; Hop 2 (foundational) — Syzygy P-124 (Go-to-Rust dual-oracle) + EAGLE-3 P-78 (speculative-decoding feedback loop) for the iterative-repair pattern. See [[1.0.0 P-151]].

### Candidate 2 — Hallucinations in LLM-Based Code Summarization: Unveiling, Detection, and Mitigation (FSE 2026, DOI 10.1145/3808189, PACMSE Vol. 3 Issue FSE)

**Status: ACCEPT.** FSE 2026 (CCF-A, in-list per §7). Introduces three artifacts: (a) **Hallu-Eval** — 800-pair benchmark of original + semantically-perturbed code/summary pairs (natural + induced logical hallucinations); (b) **Hallu-Det** — detection approach for code-summarization hallucinations; (c) **Hallu-Shield** — inference-time mitigation method. Mechanism concrete (benchmark + detector + inference-time mitigation triplet). Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Hallu-Eval as the systematic hallucination-evaluation substrate for the Comprehension phase) and `[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination (Hallu-Shield as inference-time mitigation alongside the existing cloze-reformulation anchor P-55/P-81). Citation hops: Hop 1 (mechanism) — Hallu-Eval paper itself, FSE 2026; Hop 2 (foundational) — P-143 HalluShield (FSE 2026, Wave-20) for the SLM-grounded speculative-decoding hallucination defence and P-122 ReasoningBank (Wave-16) for the test-time-scaling defence against hallucinated plans. See [[1.0.0 P-152]].

## Rejected (3)

### R1 — Code vs. Serialized AST Inputs for LLM-Based Code Summarization: An Empirical Study (Dong, Zhao & Harvey, LLM4Code 2026 workshop, co-located with ICSE 2026, arXiv:2602.06671)

**Status: REJECT Q1 (off-list workshop venue).** LLM4Code is a workshop co-located with ICSE 2026, **not** the ICSE 2026 main track. §7 trigger-list permits workshop-track ACCEPT only for method-level single-paper mechanism work; this paper is an empirical comparison study, not a method-level mechanism. Off-list per §7. arXiv:2602.06671; venue confirmed at conf.researchr.org LLM4Code 2026 program page. Mechanism (empirical comparison of raw code vs serialized AST inputs for LLM summarization) is informative but not method-level; empirical-only papers do not qualify for the workshop-track exception. Watch for ICSE 2027 main-track version.

### R2 — SmartComment: Detecting Code-Comment Inconsistencies in Smart Contracts by Combining LLM and Program Analysis (Zhang et al., FSE 2026, no DOI surfaced)

**Status: REJECT Q3 (off-axis target language).** FSE 2026 venue confirmed (in-list), mechanism concrete (LLM-driven review/developer-simulator workflow + comment-propagation/code-context extraction + program-variant differential analysis; 79.9% precision, 81.3% F1 on 1,000 real-world smart contracts), but the **target domain is Solidity smart contracts** — explicitly off-axis from the MAgHARCM modernization substrate (C → Rust, Go → Rust, Java → Rust). The four-benchmark target list (GildedRose C → Rust, Gohistogram Go → Rust, Stats Go → Rust, Commons-Validator Java → Rust) does not include Solidity. No MAgHARCM primitive anchor: PRIM-22 (comprehension) is anchored on C/Go/Java/legacy substrates; PRIM-25 (de-hallucination) is anchored on P-152 Hallu-Eval + P-143 HalluShield. Off-axis mechanism cannot defend an existing primitive. Watch for generalisation to C/Go/Java.

### R3 — Beyond Accuracy: Characterizing Code Comprehension Capabilities in (Large) Language Models (Mächtle, Serr, Loose & Eisenbarth, DeepTest 2026 workshop, co-located with ICSE 2026, arXiv:2601.12951)

**Status: REJECT Q1 (off-list workshop venue).** DeepTest 2026 is a workshop co-located with ICSE 2026, **not** the ICSE 2026 main track. §7 trigger-list permits workshop-track ACCEPT only for method-level single-paper mechanism work; this paper proposes a **diagnostic framework** (binary I/O consistency task + shadow models) rather than a method-level mechanism for code comprehension itself. Diagnostic-frameworks are evaluation tools, not method-level mechanism. AUROC 0.63 for human-metric predictors vs AUROC 0.86 for shadow models is empirically interesting but does not constitute a method-level mechanism in the §7 sense. Watch for ICSE 2027 / FSE 2027 main-track method-level extension.

## Watchlist — UNVERIFIED (1)

### W1 — SWE-TRACE: Optimizing Long-Horizon SWE Agents Through Rubric Process Reward Models and Heuristic Test-Time Scaling (Han, Xie, Ma, Zhu, Zhang, Long, Chen & Ye, arXiv:2604.14820)

**Status: UNVERIFIED (fourth carry, Wave-19 W2 → Wave-20 W1 → Wave-21 W1 → Wave-22 W1).** Re-verification 2026-09-07 (iter-6): arXiv:2604.14820 (April 2026) remains a preprint only; no peer-reviewed venue acceptance at NeurIPS 2026 / ICML 2026 / ICLR 2026 / FSE 2026 / ICSE 2026 / ASE 2026 / TOSEM 2026 / TSE 2026. NeurIPS 2026 author notifications scheduled 2026-09-24 (post this wave, 17 days out). Mechanism concrete (60K SFT corpus distillation + rubric-PRM + heuristic TTS); 8B scale matches Wave-19 SLM focus. Watchlist for Wave-23 venue confirmation after NeurIPS 2026 decisions.

> **Partition summary.** ACCEPT = 2 (P-151..P-152). REJECT Q1 = 2 (R1 Code vs. Serialized AST off-list workshop, R3 Beyond Accuracy off-list workshop). REJECT Q3 = 1 (R2 SmartComment off-axis Solidity target). UNVERIFIED = 1 (W1 SWE-TRACE fourth carry). Total triaged this wave = 5.

## Wave-22 outcome

- **2 new P-NN anchors persisted**: `[[1.0.0 P-151]]` SmartC2Rust (ICSE 2026), `[[1.0.0 P-152]]` Hallu-Eval / Hallucinations in LLM-Based Code Summarization (FSE 2026).
- **3 rejections** logged with per-candidate rationale (R1 Code vs. Serialized AST off-list workshop, R2 SmartComment off-axis Solidity target, R3 Beyond Accuracy off-list workshop diagnostic-framework).
- **1 UNVERIFIED item** carried (W1 SWE-TRACE; awaiting NeurIPS 2026 notifications on 2026-09-24, fourth carry).
- **Vault coverage gains**:
  - `[[1.0.0 PRIM-23]]` Chunked Translation: P-151 (SmartC2Rust context-aware C segmentation = chunking substrate for the C-to-Rust target). One new chunked-translation anchor at C-to-Rust.
  - `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter: P-151 (SmartC2Rust iterative Rust-compiler + semantic-diff + unsafe-residual feedback loop = recruiter-driven repair cycle for C-to-Rust). One new iterative-repair anchor at C-to-Rust; complements P-124 Syzygy Go-to-Rust anchor from Wave-17.
  - `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement: P-151 (SmartC2Rust compiler-feedback + unsafe-residual feedback = retrieval-refinement substrate). One new retrieval-refinement anchor.
  - `[[1.0.0 PRIM-22]]` Four Phases of Comprehension: P-152 (Hallu-Eval as systematic hallucination-evaluation substrate for the Comprehension phase). One new comprehension-evaluation anchor.
  - `[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination: P-152 (Hallu-Shield as inference-time mitigation alongside P-55/P-81 cloze-reformulation and P-143 HalluShield SLM-grounded speculative-decoding defence). One new inference-time mitigation anchor.

- **Vault paper count: 150 → 152.**

## Strict program-comprehension-mechanism slot — status update

- Wave-17 first opened this slot; Wave-18 partial closure via architecture-recovery trio + ADI; Wave-19 closed LLM-augmented static-analysis slots (P-129 LλMDA, P-139 TypePro, P-140 Panta); Wave-20 partial closure via P-142 NESA Datalog policy; Wave-21 partial closure via P-150 TestPrune coverage-driven context pruning.
- **Wave-22 status**: P-152 Hallu-Eval is a *systematic hallucination-evaluation substrate* for code summarization; this is an adjacent slot (evaluation, not mechanism). The strict-mechanism slot remains open for a *function-level → partition-aligned summary pass* at NeurIPS 2026 / ICML 2027 / ICLR 2027.
- Venue sweep 2026-09-07 confirms: NeurIPS 2026 author notifications 2026-09-24 (17 days out); ICML 2027 submissions Jan 22 2027 (4.5 months out); ICLR 2027 submissions 2026-09-25 (18 days out, embargoed). Wave-23 must re-query after NeurIPS 2026 author notifications.

## Watchlist for Wave-23+

- **W1 SWE-TRACE** (UNVERIFIED fourth carry) — re-verify after NeurIPS 2026 author notifications (2026-09-24).
- **R1 Code vs. Serialized AST** (REJECT Q1, off-list workshop) — watch for ICSE 2027 main-track version.
- **R2 SmartComment** (REJECT Q3, off-axis Solidity target) — watch for generalisation to C/Go/Java.
- **R3 Beyond Accuracy** (REJECT Q1, off-list workshop diagnostic-framework) — watch for ICSE 2027 / FSE 2027 main-track method-level extension.
- **R1-W21 ABC** (REJECT Q1, arXiv-only) — watch for ICML 2026 / NeurIPS 2026 / ICLR 2027 companion paper at a §7 whitelist venue.
- **Strict program-comprehension-mechanism slot** — concept-assignment / dynamic-invariants / temporal-coupling mechanism at NeurIPS 2026 / ICML 2027 / ICLR 2027 (function-level → partition-aligned summary pass).
- **Cross-agent KV cache policies** beyond P-147 SpecKV / P-148 LookaheadKV / P-149 SSD/Saguaro / P-141 KVFlow at NeurIPS 2026 / ICLR 2026 workshops.
- **Source-to-source modernization** beyond P-145 TerraMod + P-151 SmartC2Rust — ICSE 2027 / FSE 2027.

**TOTAL FIRED: 2 ACCEPT + 3 REJECT + 1 UNVERIFIED (Watchlist, carried).** Headline ACCEPT count: 2. Headline REJECT count: 3 (2 Q1 + 1 Q3). UNVERIFIED: 1 (carried).
