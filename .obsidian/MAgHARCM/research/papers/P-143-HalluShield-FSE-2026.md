---
title: "P-143 Hallucinations in LLM-based Code Summarization: Unveiling, Detection, and Mitigation (Hallu-Eval / Hallu-Det / Hallu-Shield)"
backlink: "[[1.0.0 P-143]]"
aliases:
  - "1.0.0 P-143"
  - "P-143"
  - "P-143-HalluShield-FSE-2026"
  - "HalluShield-FSE-2026"
  - "hallushield2026fse"
tags: [paper, code-summarization, hallucination-detection, value-model, slm-as-judge, "[[1.0.0 PRIM-7]]", "[[1.0.0 PRIM-22]]", wave-20]
date: 2026-09-07
last_updated: 2026-09-07
venue: FSE 2026 (Research Track)
---

# [[1.0.0 P-143]] Hallu-Shield

## TL;DR

Hallu-Shield is an **inference-time mitigation mechanism** for LLM-based code-summarization hallucinations, layered on top of **Hallu-Det** (synergistic entity-level detection + synonymous-mutation-based refinement; F1 **0.95** on Qwen2.5-Coder-7B summaries) and **Hallu-Eval** (an 800-pair benchmark of natural + semantically perturbed hallucinations). Hallu-Shield uses an **external value model** to guide the LLM toward more faithful summaries during inference (no retraining required). On DeepSeek-Coder-6.7B, achieves **10.6% relative reduction** in hallucination rate (66% → 59%) and **74.0% win rate** under an LLM-as-a-judge majority vote. Anchors `[[1.0.0 PRIM-7]]` (external value model as a SLM-scale verifier replacing frontier-PRM-as-judge) and `[[1.0.0 PRIM-22]]` (Insight-phase code-summarization fidelity).

## Mechanism (Q2)

1. **Hallu-Det (detection)** — two-prong strategy:
   - **Entity-level detection** — explicitly identifies hallucinations tied to code entities (function names, types, side-effects).
   - **Synonymous-mutation-based refinement** — generates synonymous code mutations; hallucinations that survive the mutation confirm a factual mismatch.
2. **Hallu-Shield (inference-time mitigation)** — an external value model (separate from the LLM) scores candidate summaries during decoding and steers the LLM toward high-fidelity outputs. No model retraining required.
3. **Hallu-Eval (benchmark)** — 800 code-summary pairs combining naturally occurring hallucinations with semantically perturbed code designed to induce complex logical hallucinations.

Together these give MAgHARCM's verdict-validation pattern a **summarization-specific verifier** that operates at inference time, complementing the Wave-18 SLM-as-judge family (P-125, P-126, P-127).

## Anchoring (Q3)

| Primitive | Pre-wave-20 behaviour | Hallu-Shield substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation | Verdict panel relied on either a frontier-PRM (P-92 Lightman PRM800K, expensive) or a SLM-as-a-judge panel (P-127) | Hallu-Shield's external value model is a *summarization-specific verifier* — a third axis beyond PRM and SLM-as-judge, operating at inference time |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Insight phase) | Insight-phase code summaries were vulnerable to entity-level hallucinations (function names, types) — downstream Reverse-Topo-Plan and Migration Strategy Selection acted on hallucinated names | Hallu-Det's entity-level detection + Hallu-Shield's value-model steering suppress hallucinations before they enter downstream primitives |

## Hop-1 Citations

- Hallu-Eval 2024 (Wang et al., code-summarization hallucination benchmark — substrate for the 800-pair dataset).
- P-127 SLM-as-a-Judge (ICSE 2026, Crupi et al.) — sister SLM-scale verifier in the Wave-18 PRIM-7 cluster.
- P-92 Lightman PRM800K 2023 — frontier-PRM-as-judge baseline that Hallu-Shield supersedes for summarization.
- LLM-as-a-judge majority-vote (Zheng et al. NeurIPS 2023).

## Hop-2 Citations

- Value-model guidance in RLHF (Christiano et al. 2017, Stiennon et al. 2020).
- Code-summarization benchmark lineage (CodeXGLUE 2021, TLCodeSum 2022).
- P-87 Hou et al. TOSEM 2024 (LLM4SE SLR).

## MAgHARCM integration

- **YAML config key**: `verdict.hallu_shield.enabled: true`; `verdict.hallu_shield.value_model: <model_id>`; `verdict.hallu_shield.detect_entity_level: true`.
- **Implementation file**: `internal/verdict/hallu_shield.go::NewHalluShieldVerifier` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-22]]`.

## Caveats

- **Task-specific verifier** — Hallu-Shield is tuned for code summarization; transferring to other PRIM-7 tasks (e.g., type inference, function generation) requires a new value-model training set.
- **External value model cost** — the value model adds an inference step; effective cost = LLM decode × value-model scoring.
- **Synonymous-mutation refinement cost** — Hallu-Det's refinement step invokes N mutations per candidate; N is a configuration knob.

## Source

- Venue: FSE 2026 (Research Track), verified via PACMSE Vol 3 Issue FSE Article 66, DOI: 10.1145/3808139; researchr.org/details/fse-2026/fse-2026-research-papers/66; conference-publishing.com/authors/FSE26.
