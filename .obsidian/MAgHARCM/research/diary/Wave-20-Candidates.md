---
title: Wave-20 Candidate Evaluation (2026-09-07)
backlink: "[[1.0.0 Wave-20 Candidates]]"
tags: [candidates, wave-20, mechanism-vs-benchmark, "[[1.0.0 P-142]]", "[[1.0.0 P-143]]", "[[1.0.0 P-144]]", "[[1.0.0 P-145]]", "[[1.0.0 P-146]]", software-archaeology, neuro-symbolic, dynamic-analysis, code-summarization, kv-cache, slm-as-judge, status: FIRED]
status: FIRED
date: 2026-09-07
last_updated: 2026-09-07
---

# [[1.0.0 Wave-20 Candidates]]

## Trigger-gate evaluation

§7 trigger gate (rewritten 2026-09-25): Q1 venue confirmation (peer-reviewed at NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE; workshop-track permitted if method-level) + Q2 mechanism-vs-benchmark (concrete method, not benchmark / eval / prompt tweak) + Q3 anchoring-to-existing-primitive (defends or refutes an existing `[[1.0.0 PRIM-NN]]`).

## Wave-20 focus areas

1. **Wave-19 watchlist re-verification** — SliceMate (W1) and SWE-TRACE (W2). SliceMate was confirmed at ISSTA 2025 (PACMSE Vol 2) — ISSTA is **off-list** per §7 whitelist, so W1 promotes to REJECT Q1 (off-list venue). SWE-TRACE remains UNVERIFIED (arXiv-only, no peer-reviewed acceptance).
2. **Strict program-comprehension-mechanism slot** — NeurIPS 2026 / ICML 2027 / ICLR 2027 / FSE 2026 / ICSE 2026 scout for a concrete mechanism (not a benchmark) introducing a new program-comprehension substrate at SLM scale.
3. **Wave-19 SLM-anchor cluster reinforcement** — KV-cache / speculative-decoding / TTS-mechanism papers at 2025+ venues that strengthen the SLM-scale PRIM-7 / PRIM-21 / PRIM-31 substrate.

## Candidates (8 triaged, 5 ACCEPT + 2 REJECT + 1 UNVERIFIED)

### Candidate 1 — NESA: Relational Neuro-Symbolic Static Program Analysis (Wang, Gao, Zhang, Liu, Guo, Zheng, Shi & Zhang, FSE 2026)

**Status: ACCEPT.** Restricted Datalog analysis-policy language decomposes complex program-analysis problems into syntactic (handled by parsing-based analysis) and semantic (handled by LLM) sub-problems. Lazy and incremental prompting reduces hallucinations. Compilation-free and customizable. F1 0.72 on TaintBench (surpasses industrial baseline by +0.20); 13 real-world memory-leak bugs detected and fixed by developers. Anchors `[[1.0.0 PRIM-22]]` (Four Phases of Comprehension: Reorganization + Insight via Datalog policy) and `[[1.0.0 PRIM-9]]` (Tri-Representation Hybrid Code Graph: a *fourth representation* — neuro-symbolic policy — on top of syntactic / lexical / semantic). See [[1.0.0 P-142]].

### Candidate 2 — Hallucinations in LLM-based Code Summarization: Unveiling, Detection, and Mitigation (Hallu-Eval + Hallu-Det + Hallu-Shield, FSE 2026)

**Status: ACCEPT.** Two mechanisms bundled in one paper:
- **Hallu-Det** (detection): synergistic entity-level detection + synonymous-mutation-based refinement (F1 0.95 on Qwen2.5-Coder-7B summaries).
- **Hallu-Shield** (inference-time mitigation): external value model guides LLM toward more faithful summaries. 10.6% relative reduction in hallucination rate on DeepSeek-Coder-6.7B (66% → 59%); 74.0% win rate under LLM-as-a-judge majority vote.

Anchors `[[1.0.0 PRIM-7]]` (Verdict Validation: external value model is a SLM-scale verifier replacing frontier-PRM-as-judge) and `[[1.0.0 PRIM-22]]` (Four Phases of Comprehension: Insight-phase code-summarization fidelity). PACMSE Vol 3 FSE Issue Article 66, DOI: 10.1145/3808139. See [[1.0.0 P-143]].

### Candidate 3 — TraceCoder: A Trace-Driven Multi-Agent Framework for Automated Debugging of LLM-Generated Code (Huang, Ye, Sun, Zhang, Zhang & Liu, ICSE 2026 Research Track)

**Status: ACCEPT.** Four-component mechanism:
- **Runtime instrumentation** captures fine-grained execution traces (beyond binary pass/fail).
- **Causal analysis** uses traces for accurate root-cause localization.
- **Historical Lesson Learning Mechanism (HLLM)** distills insights from previous failed repair attempts, preventing repetitive cycles.
- **Rollback Mechanism (RM)** enforces that every iterative repair step is a strict improvement toward a correct solution.

Up to 34.43% relative Pass@1 improvement; iterative repair alone contributes 65.61% relative gain. Anchors `[[1.0.0 PRIM-22]]` (Four Phases of Comprehension: Observation phase via trace instrumentation) and `[[1.0.0 PRIM-25]]` (Role-Flip Reviewer: HLLM is the multi-agent substrate for cross-iteration review). See [[1.0.0 P-144]].

### Candidate 4 — TerraMod: Automating Terraform Code Migration through Provider Evolution Knowledge (Gupta, Aggarwal, Paulovicks, Mohapatra, Lee & Sheinin, IBM Research, ICSE 2026 NIER Track)

**Status: ACCEPT (borderline, NIER track).** §7 workshop-track rationale analog: NIER is **not** a workshop — it is a peer-reviewed, double-anonymous sub-track of ICSE main conference (treated as in-list per §7 ICSE entry). Method-level threshold met: concrete mechanism (changelogs + API schemas + deprecation links → structured migration context → LLM-guided migration) is method-level, not benchmark. Anchors `[[1.0.0 PRIM-21]]` (Migration Strategy Selection: knowledge-augmented migration context as a runtime signal for strategy choice) and `[[1.0.0 PRIM-22]]` (Four Phases of Comprehension: Reorganization phase via external knowledge context). See [[1.0.0 P-145]].

### Candidate 5 — ContextPRM: Leveraging Contextual Coherence for Multi-Domain Test-Time Scaling (Zhang, Liu, Yu, Qiu, Xiao, Ren, Chen & Liu, ICLR 2026)

**Status: ACCEPT.** Domain-agnostic logical-flow PRM trained on contextual coherence rather than domain-specific knowledge. 6.5% average accuracy improvement on MMLU-Pro across nine non-mathematical domains over majority voting (WMV); outperforms VersaPRM (2.2%) and math-focused PRMs (0.5%). Strong cross-domain generalization even when fine-tuned on a single domain. Anchors `[[1.0.0 PRIM-7]]` (Verdict Validation: domain-agnostic PRM is a SLM-scale verifier for cross-domain judgement; complements P-125 T1, P-126 ARC-Decode, P-127 SLM-as-a-Judge, P-135 SPECS, P-136 CaTS, P-137 SuffixDecoding, P-143 Hallu-Shield as the *eighth* PRIM-7 cluster anchor). See [[1.0.0 P-146]].

## Rejected (2)

### R1 — SliceMate: Accurate and Scalable Static Program Slicing via LLM-Powered Agents (Chang, Shi, Lyu, Zhou, Wang, Yang, Li & Lo, ISSTA 2025 / PACMSE Vol 2)

**Status: REJECT Q1 (off-list venue).** Re-verification of Wave-19 W1: Yunbo Lyu's homepage lists PTMware at ISSTA 2026, but SliceMate itself is published at **ISSTA 2025** (PACMSE Vol 2, ISSTA:2362–2383). §7 whitelist is explicitly: NeurIPS / ICML / ICLR / FSE / ICSE / ASE / TOSEM / TSE (workshop-track permitted if method-level). ISSTA is **not in §7 whitelist** — although it is closely adjacent to FSE (both ACM SIGSOFT software-engineering venues), it is not literally named. Methodological note: §7's workshop-track ACCEPT exception covers workshops on a WHITELISTED venue (e.g., NeurIPS 2025 LAW Workshop, ICLR 2026 workshop tracks) — not a separate non-whitelisted venue even if method-level. SliceMate's ISSTA publication does not satisfy the §7 workshop-track rationale. Mechanism is concrete (three LLM agents replace explicit PDG/SDG construction; +22% acc / +28% F1 on SliceBench 2,200-program benchmark) and high-relevance for Wave-20 focus area (b) program-comprehension-mechanism. Watchlist for FSE 2027 / TOSEM 2026 companion paper at an in-list venue.

### R2 — Environment-in-the-Loop: Rethinking Code Migration with LLM-based Agents (Li, Fei, Ma, Zhang, Sarro & Ye, ReCode 2026 Workshop, co-located with ICSE 2026)

**Status: REJECT Q1 (off-list workshop).** ReCode 2026 is the *1st International Workshop on Code Translation, Transformation, and Modernization* co-located with ICSE 2026 — workshop-track. §7 workshop-track ACCEPT rationale requires that the main-track alternative be unavailable in the focus area; MAgHARCM already anchors the environment-in-the-loop substrate via P-145 TerraMod (ICSE 2026 NIER, main-track subvenue) and the existing PRIM-1/22 substrate, so this paper is redundant. Mechanism (semantic code-error feedback guiding LLM agents during migration) is concrete but redundant.

## Watchlist — UNVERIFIED (1)

### W1 — SWE-TRACE: Optimizing Long-Horizon SWE Agents Through Rubric Process Reward Models and Heuristic Test-Time Scaling (Han, Xie, Ma, Zhu, Zhang, Long, Chen & Ye, arXiv:2604.14820)

**Status: UNVERIFIED (carried from Wave-19 W2).** Re-verification 2026-09-07: arXiv:2604.14820 (April 2026) remains a preprint only; no peer-reviewed venue acceptance at NeurIPS 2026 / ICML 2026 / ICLR 2026 / FSE 2026 / ICSE 2026 / ASE 2026 / TOSEM 2026 / TSE 2026. NeurIPS 2026 author notifications are scheduled for 2026-09-24 (post this wave). Mechanism concrete (60K SFT corpus distillation + rubric-PRM + heuristic TTS); 8B scale matches Wave-19 SLM focus. Watchlist for Wave-21 venue confirmation after NeurIPS 2026 decisions.

> **Partition summary.** ACCEPT = 5 (P-142..P-146). REJECT Q1 = 2 (R1 SliceMate, R2 Environment-in-the-Loop). UNVERIFIED = 1 (W1 SWE-TRACE carried from Wave-19). Total triaged this wave = 8.

## Wave-20 outcome

- **5 new P-NN anchors persisted**: `[[1.0.0 P-142]]` NESA, `[[1.0.0 P-143]]` Hallu-Shield, `[[1.0.0 P-144]]` TraceCoder, `[[1.0.0 P-145]]` TerraMod, `[[1.0.0 P-146]]` ContextPRM.
- **2 rejections** logged with per-candidate rationale (R1 SliceMate ISSTA off-list, R2 Environment-in-the-Loop workshop redundancy).
- **1 UNVERIFIED item** carried from Wave-19 W2 (SWE-TRACE; awaiting NeurIPS 2026 notifications).
- **Vault coverage gains**:
  - `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation: P-143 (Hallu-Shield external value model), P-146 (ContextPRM domain-agnostic PRM). Two new SLM-scale verifier anchors, bringing the cluster to 8 papers (P-125, P-126, P-127, P-135, P-136, P-137, P-143, P-146).
  - `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph: P-142 (NESA neuro-symbolic policy as a fourth representation on top of syntactic/lexical/semantic).
  - `[[1.0.0 PRIM-21]]` Migration Strategy Selection: P-145 (TerraMod knowledge-augmented migration context as a runtime strategy signal).
  - `[[1.0.0 PRIM-22]]` Four Phases of Comprehension: P-142 (NESA Reorganization via Datalog policy), P-143 (Hallu-Shield Insight-phase summary fidelity), P-144 (TraceCoder Observation phase via trace instrumentation), P-145 (TerraMod Reorganization via external knowledge context). Four Wave-20 comprehension-substrate anchors.
  - `[[1.0.0 PRIM-25]]` Role-Flip Reviewer: P-144 (TraceCoder HLLM as cross-iteration review substrate).

- **Vault paper count: 141 → 146.**

## Strict program-comprehension-mechanism slot — status update

- Wave-18 R4 deferred this slot; Wave-19 closed LLM-augmented static-analysis slots but the *strict program-comprehension-mechanism* slot remained open.
- **Wave-20 partial closure**: P-142 NESA (Datalog-policy + LLM hybrid program-analysis mechanism), P-143 Hallu-Shield (code-summarization fidelity mechanism), P-144 TraceCoder (trace-driven debugging comprehension), and P-145 TerraMod (knowledge-context migration comprehension) collectively fill the gap at FSE 2026 / ICSE 2026.
- The strict-mechanism slot remains open for a *concept-assignment / dynamic-invariants / temporal-coupling* mechanism specifically — see Wave-21 watchlist.

## Watchlist for Wave-21+

- **W1 SWE-TRACE** (UNVERIFIED carried from Wave-19 W2) — re-verify after NeurIPS 2026 author notifications (2026-09-24).
- **R1 SliceMate** (REJECT Q1, off-list ISSTA venue) — watch for FSE 2027 / TOSEM 2026 companion paper at an in-list venue.
- **R2 Environment-in-the-Loop** (REJECT Q1, workshop redundancy vs P-145) — watch for ICSE 2027 main-track version.
- **Strict program-comprehension-mechanism slot** — concept-assignment / dynamic-invariants / temporal-coupling mechanism at NeurIPS 2026 / ICML 2027 / ICLR 2027.
- **Cross-agent KV cache policies** beyond P-134 RelayCaching / P-141 KVFlow at NeurIPS 2026.
- **Source-to-source modernization** beyond P-145 TerraMod — ICSE 2027 / FSE 2027.

**TOTAL FIRED: 5 ACCEPT + 2 REJECT + 1 UNVERIFIED (Watchlist, carried).** Headline ACCEPT count: 5. Headline REJECT count: 2. UNVERIFIED: 1 (carried).
