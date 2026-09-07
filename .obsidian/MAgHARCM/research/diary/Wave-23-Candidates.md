---
title: Wave-23 Candidate Evaluation (2026-09-07)
backlink: "[[1.0.0 Wave-23 Candidates]]"
tags: [candidates, wave-23, icse-2026, fse-2026, slicing, code-translation, static-analysis, test-generation, "[[1.0.0 P-153]]", "[[1.0.0 P-154]]", "[[1.0.0 P-155]]", "[[1.0.0 P-156]]", slm-era, status: FIRED]
status: FIRED
date: 2026-09-07
last_updated: 2026-09-07 (iter-7, wave-23)
aliases:
 - "Wave-23-Candidates"
 - "Wave 23 Candidates"
 - "Wave-23"
---

# [[1.0.0 Wave-23 Candidates]]

## Trigger-gate evaluation

§7 trigger gate (Wave-23 firing criterion): Q1 venue confirmation (peer-reviewed at NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE; workshop-track permitted ONLY if method-level single-paper) + Q2 mechanism-vs-benchmark (concrete method, not benchmark / eval / prompt tweak alone) + Q3 anchoring-to-existing-primitive (defends or refutes an existing `[[1.0.0 PRIM-NN]]`).

## Wave-23 focus areas

1. **Post-FSE-2026 / post-ICSE-2026 concluded-program scout** — both flagship SE venues concluded in 2026-Q2 (ICSE 2026 April Rio de Janeiro; FSE 2026 July Montréal, published as PACMSE Vol 3 Issue FSE). Scout for strict-mechanism papers in slicing, code translation, static-analysis warning repair, and LLM-augmented test generation.
2. **CoReX strict-mechanism closure** — P-153 CoReX (ICSE 2026) closes the strict program-comprehension-mechanism slot that Wave-22 §"Strict program-comprehension-mechanism slot" carried forward (concept-assignment / dynamic-invariants / slicing residual).
3. **W1 SWE-TRACE re-verification** — fifth carry. arXiv:2604.14820 (April 2026) still preprint-only per dblp "Informal or Other Publication"; NeurIPS 2026 author notifications 2026-09-24 (16 days out, post this wave). Re-verify in Wave-24.
4. **CodeCureAgent / TransAgent / TestWeaver triplet** — three additional mechanism papers surfaced from FSE 2026 + ICSE 2026 concluded programs. Two FSE 2026 anchors (TransAgent translation; CodeCureAgent static-analysis warning repair) plus one ICSE 2026 anchor (TestWeaver LLM-augmented backward-slicing test generation).

## Candidates (6 triaged, 4 ACCEPT + 1 REJECT + 1 UNVERIFIED)

### Candidate 1 — CoReX: Context-Aware Refinement-Based Slicing for Debugging Regression Failures (Badihi & Rubin, ICSE 2026 Research Track, University of British Columbia)

**Status: ACCEPT.** ICSE 2026 (April 12-18 2026, Rio de Janeiro) — flagship SE venue, in-list per §7. Mechanism is concrete: (a) **context-aware refinement-based slicing** — instead of static or observed-dependency-only slicing (e.g., SliceMate P-rejected-Wave-19), CoReX conditions the slice on the regression failure context by iteratively refining the candidate slice using observed execution behaviour and recent change history; (b) **regression-failure localisation** — the refined slice localises the regression-causing code with higher precision than vanilla slicing; (c) **empirical validation** on Java regression-failure benchmarks. Closes the strict program-comprehension-mechanism residual slot carried forward from Wave-17 first opening. **Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Structure phase, function-level DA pattern extended from P-133 ADI) and `[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination (refinement-based slicing = context-conditioned summary pass = partition-aligned summary-pass substrate).** Citation hops: Hop 1 (mechanism-support) → P-133 ADI (Frame Lifetime Trace, function-level DA pattern); Hop 2 (foundational) → Weiser program slicing (1981 IEEE TSE). One new strict-mechanism anchor.

### Candidate 2 — TransAgent: Enhancing LLM-Based Code Translation via Fine-Grained Execution Alignment (Yuan, W.Chen, H.Wang, Peng, Z.Chen & Lou, FSE 2026 / PACMSE Vol 3 Issue FSE, DOI 10.1145/3797099, Fudan University)

**Status: ACCEPT.** FSE 2026 (July 5-9 2026, Montréal; published in PACMSE Vol 3 Issue FSE) — flagship SE venue, in-list per §7. Mechanism is concrete: (a) **multi-agent translation pipeline** — multiple LLM agents collaborate on source-to-target translation with role decomposition (translator + critic + repairer); (b) **fine-grained execution alignment** — error-prone code blocks are localised by comparing the execution behaviour of the source code against the translated (target) code rather than relying on test-output-only diffs; (c) **iterative repair loop** drives the agents to fix the localised errors. Outperforms UniTrans by up to 33.3% in translation accuracy and improves program repair performance by 56.7% on average over Agentless. **Anchors `[[1.0.0 PRIM-2]]` Multi-Agent Translation Pipeline (TransAgent as multi-agent collaborator + execution-aligned critic, complementing P-151 SmartC2Rust single-LLM iterative feedback) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (execution-aligned critic feedback = retrieval refinement substrate).** Citation hops: Hop 1 (mechanism-support) → P-151 SmartC2Rust (iterative compiler-feedback repair, Wave-22); Hop 2 (foundational) → UniTrans (Yan et al., ASE 2023, code-translation survey baseline). One new multi-agent-translation anchor at PACMSE/FSE.

### Candidate 3 — CodeCureAgent: Automatic Classification and Repair of Static Analysis Warnings (Joos, Bouzenia & Pradel, FSE 2026 / PACMSE Vol 3 Issue FSE, DOI 10.1145/3808139, arXiv:2509.11787, CISPA Helmholtz Center for Information Security)

**Status: ACCEPT.** FSE 2026 (PACMSE Vol 3 Issue FSE) — flagship SE venue, in-list per §7. Mechanism is concrete: (a) **agentic LLM framework** — iterates code-search + multi-file-edit tool invocations rather than a fixed algorithmic pipeline; (b) **classification + repair twin objective** — distinguishes false positives from true positives before repair, suppressing false positives and patching true positives; (c) **three-step validation heuristic** for patch approval — (i) project builds clean, (ii) target warning removed without new warnings, (iii) test suite green. Evaluated on 1,000 SonarQube warnings across 106 Java projects: 96.8% plausible-fix rate, 86.3% manual-inspection correct-fix rate, ~$0.029 + ~4 min per warning. **Anchors `[[1.0.0 PRIM-12]]` Static Analysis Co-Evolution (agentic warning-classification+repair twin-objective extends the static-analysis substrate) and `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter (three-step build+test validation = iteration-recruiter pattern with concrete acceptance heuristic).** Citation hops: Hop 1 (mechanism-support) → P-124 Syzygy (iterative compilation-feedback loop, Wave-17); Hop 2 (foundational) → SonarQube static-analysis warning taxonomy + prior LLM4Code warning-repair work. One new static-analysis-agent anchor.

### Candidate 4 — TestWeaver: Execution-aware, Feedback-driven Regression Testing Generation with Large Language Models (Le, Van, Vu, M.V.T.Pham, H.N.Phan, H.N.Phan & T.N.Nguyen, ICSE 2026, DOI 10.1145/3744916.3787805, arXiv:2508.01255, FPT Software AI Center / VinUniversity / Nanyang Technological University / University of Texas at Dallas)

**Status: ACCEPT.** ICSE 2026 (April 2026, Rio de Janeiro) — flagship SE venue, in-list per §7. Mechanism is concrete: (a) **backward slicing for context reduction** — instead of full-program context, the LLM receives a backward slice from the target line to reduce hallucinations and keep focus; (b) **close-test retrieval** — identifies prior tests that share control-flow similarity with the path to the target line and includes them as execution context; (c) **execution in-line annotations** — variable-state comments along the executed path are injected into the prompt. Addresses the coverage-plateau problem in LLM-based regression test generation. **Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Structure phase — backward slicing + execution annotations = LLM-augmented static-analysis pattern complementing P-129 LλMDA's partial-PDG pattern and P-156's context-aware refinement pattern) and `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter (close-test retrieval + execution annotation = retrieval-and-recruitment substrate for test-generation agents).** Citation hops: Hop 1 (mechanism-support) → P-129 LλMDA (LLM-augmented partial program dependence analysis, Wave-18) + P-133 ADI (Frame Lifetime Trace, Wave-18); Hop 2 (foundational) → traditional backward slicing (Weiser 1981 IEEE TSE) + execution-aware test generation lineage (e.g., Evosuite, Panko). One new LLM-augmented-test-generation anchor.

## Rejected (1)

### R1 — AutoCodeSherpa: Symbolic Explanations for Agentic Code (Yunbo Lyu et al., ISSTA 2026)

**Status: REJECT Q1 (off-list venue).** ISSTA 2026 is **not** on the §7 trigger list (NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE). §7 trigger list permits workshop-track ACCEPT only for method-level single-paper mechanism work, and ISSTA is a main conference (not workshop). Off-list per §7 strict-venue interpretation. Watch for a §7-list venue companion paper (ICSE 2027 / FSE 2027).

## Watchlist — UNVERIFIED (1)

### W1 — SWE-TRACE: Optimizing Long-Horizon SWE Agents Through Rubric Process Reward Models and Heuristic Test-Time Scaling (Han, Xie, Ma, Zhu, Zhang, Long, Chen & Ye, arXiv:2604.14820)

**Status: UNVERIFIED (fifth carry, Wave-19 W2 → Wave-20 W1 → Wave-21 W1 → Wave-22 W1 → Wave-23 W1).** Re-verification 2026-09-07 (iter-7): arXiv:2604.14820 (April 2026) remains a preprint only per dblp "Informal or Other Publication" entry (Aug 2026); no peer-reviewed venue acceptance at NeurIPS 2026 / ICML 2026 / ICLR 2026 / FSE 2026 / ICSE 2026 / ASE 2026 / TOSEM 2026 / TSE 2026. NeurIPS 2026 author notifications scheduled 2026-09-24 (16 days out, post this wave). Mechanism concrete (60K SFT corpus distillation + rubric-PRM + heuristic TTS); 8B scale matches Wave-19 SLM focus. Watchlist for Wave-24 venue confirmation after NeurIPS 2026 decisions.

> **Partition summary.** ACCEPT = 4 (P-153..P-156). REJECT Q1 = 1 (R1 AutoCodeSherpa off-list ISSTA). UNVERIFIED = 1 (W1 SWE-TRACE fifth carry). Total triaged this wave = 6.

## Wave-23 outcome

- **4 new P-NN anchors persisted**: `[[1.0.0 P-153]]` CoReX (ICSE 2026), `[[1.0.0 P-154]]` TransAgent (FSE 2026 / PACMSE), `[[1.0.0 P-155]]` CodeCureAgent (FSE 2026 / PACMSE), `[[1.0.0 P-156]]` TestWeaver (ICSE 2026).
- **1 rejection** logged with per-candidate rationale (R1 AutoCodeSherpa off-list ISSTA venue).
- **1 UNVERIFIED item** carried (W1 SWE-TRACE; awaiting NeurIPS 2026 notifications on 2026-09-24, fifth carry).
- **Vault coverage gains**:
  - `[[1.0.0 PRIM-22]]` Four Phases of Comprehension: P-153 (CoReX refinement-based slicing = Structure-phase mechanism complementing P-133 ADI Frame Lifetime Trace + P-129 LλMDA partial-PDG + P-156 TestWeaver backward slicing). Three new comprehension-mechanism anchors (one strict-mechanism; two LLM-augmented static-analysis).
  - `[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination: P-153 (CoReX refinement-based slicing = partition-aligned summary-pass substrate). One new de-hallucination-substrate anchor.
  - `[[1.0.0 PRIM-2]]` Multi-Agent Translation Pipeline: P-154 (TransAgent multi-agent LLM translation + execution-aligned critic). One new multi-agent-translation anchor.
  - `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement: P-154 (TransAgent execution-aligned critic feedback = retrieval refinement substrate). One new retrieval-refinement anchor.
  - `[[1.0.0 PRIM-12]]` Static Analysis Co-Evolution: P-155 (CodeCureAgent agentic warning-classification + repair twin-objective). One new static-analysis-agent anchor.
  - `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter: P-155 (three-step build+test validation = iteration-recruiter with concrete acceptance heuristic) + P-156 (close-test retrieval + execution annotation). Two new iteration-recruiter anchors.

- **Vault paper count: 152 → 156.**

## Strict program-comprehension-mechanism slot — closure

- Wave-17 first opened this slot; Wave-18 partial closure via architecture-recovery trio (P-130 SSAR + P-131 SemArc + P-132 SemRef) + ADI (P-133); Wave-19 closed LLM-augmented static-analysis slots (P-129 LλMDA, P-139 TypePro, P-140 Panta); Wave-20 partial closure via P-142 NESA Datalog policy; Wave-21 partial closure via P-150 TestPrune coverage-driven context pruning; Wave-22 adjacent (P-152 Hallu-Eval is evaluation-substrate, not mechanism).
- **Wave-23 status**: P-153 CoReX closes the **strict mechanism** slot — context-aware refinement-based slicing is a concrete mechanism (not a benchmark/evaluation/prompt-tweak). The slot that Wave-22 carried forward is now CLOSED with P-153 as the canonical anchor. The residual residual: any Wave-24+ advances to concept-assignment / temporal-coupling / dependency-structure-matrix partitions remain open at NeurIPS 2026 / ICML 2027 / ICLR 2027.

## Watchlist for Wave-24+

- **W1 SWE-TRACE** (UNVERIFIED fifth carry) — re-verify after NeurIPS 2026 author notifications (2026-09-24, 16 days out).
- **R1 AutoCodeSherpa** (REJECT Q1, off-list ISSTA) — watch for ICSE 2027 / FSE 2027 / TOSEM companion paper.
- **Strict concept-assignment / temporal-coupling / DSM-partition mechanism** — remaining residual at NeurIPS 2026 / ICML 2027 / ICLR 2027.
- **Cross-agent KV cache policies** beyond P-147 SpecKV / P-148 LookaheadKV / P-149 SSD/Saguaro / P-141 KVFlow at NeurIPS 2026 / ICLR 2026 workshops.
- **Source-to-source modernization** beyond P-145 TerraMod + P-151 SmartC2Rust + P-154 TransAgent — ICSE 2027 / FSE 2027.

**TOTAL FIRED: 4 ACCEPT + 1 REJECT + 1 UNVERIFIED (Watchlist, carried).** Headline ACCEPT count: 4. Headline REJECT count: 1 (1 Q1). UNVERIFIED: 1 (carried).
