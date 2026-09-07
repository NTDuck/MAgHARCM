---
title: "P-146 ContextPRM: Leveraging Contextual Coherence for Multi-Domain Test-Time Scaling"
backlink: "[[1.0.0 P-146]]"
aliases:
  - "1.0.0 P-146"
  - "P-146"
  - "P-146-ContextPRM-ICLR-2026"
  - "ContextPRM-ICLR-2026"
  - "contextprm2026iclr"
tags: [paper, process-reward-model, domain-agnostic, contextual-coherence, tts, "[[1.0.0 PRIM-7]]", wave-20]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICLR 2026 (Conference Paper)
---

# [[1.0.0 P-146]] ContextPRM

## TL;DR

ContextPRM is a **Process Reward Model (PRM)** trained on **domain-agnostic contextual coherence** — logical transitions between consecutive chain-of-thought steps — rather than domain-specific knowledge. Achieves **6.5% average accuracy improvement** on MMLU-Pro across nine non-mathematical domains over majority voting (WMV); significantly outperforms VersaPRM (2.2%) and math-focused PRMs (0.5%). Strong cross-domain generalization even when fine-tuned on a single domain. Anchors `[[1.0.0 PRIM-7]]` as the *eighth* SLM-scale verifier in the Wave-18 + Wave-19 + Wave-20 cluster (P-125 T1, P-126 ARC-Decode, P-127 SLM-as-a-Judge, P-135 SPECS, P-136 CaTS, P-137 SuffixDecoding, P-143 Hallu-Shield, **P-146 ContextPRM**).

## Mechanism (Q2)

1. **Domain-agnostic training objective** — instead of training on domain-specific fact verification (mathematics, programming), the model is trained to evaluate **logical transitions** (coherence) between consecutive CoT steps.
2. **Contextual coherence scoring** — at inference, the PRM scores each candidate step's coherence with the prior context, regardless of subject matter (law, history, philosophy, code, etc.).
3. **Cross-domain fine-tuning** — fine-tuning on a single domain (e.g., Law or Psychology) generalizes to other domains because the underlying skill (logical coherence) is transferable.
4. **TTS integration** — used as the scoring signal for weighted majority voting (WMV) at test-time; the PRM weights each candidate answer by its contextual-coherence score.

Together these give MAgHARCM a **domain-agnostic verifier** for multi-domain modernisation tasks (legacy Java/C/Go modernisation spans law-domain contract code, history-domain genealogy code, philosophy-domain policy code).

## Anchoring (Q3)

| Primitive | Pre-wave-20 behaviour | ContextPRM substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation | Verdict cluster was task-specific: T1 (tool-integrated, P-125), ARC-Decode (speculative, P-126), SLM-as-a-Judge (code gen, P-127), SPECS (draft, P-135), CaTS (confidence, P-136), SuffixDecoding (template, P-137), Hallu-Shield (summarization, P-143). No *domain-agnostic* verifier | ContextPRM is the domain-agnostic coherence verifier — closes the gap for cross-domain modernisation tasks |

## Hop-1 Citations

- P-92 Lightman PRM800K 2023 — frontier-PRM-as-judge baseline (math-domain).
- VersaPRM (2025) — earlier PRM generalisation attempt.
- P-91 Snell 2024 — test-time compute-optimal allocation across strategies.

## Hop-2 Citations

- Process reward model lineage (Lightman 2023, Uesato 2022).
- Chain-of-thought prompting (Wei et al. 2022, P-90).
- Domain-agnostic evaluation (MMLU-Pro 2024).

## MAgHARCM integration

- **YAML config key**: `verdict.context_prm.model_id: <id>`; `verdict.context_prm.coherence_threshold: <float>`; `verdict.context_prm.cross_domain_tune: <source_domain>`.
- **Implementation file**: `internal/verdict/context_prm.go::NewContextPRMScorer` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-7]]`.

## Caveats

- **Coherence ≠ correctness** — contextual coherence is necessary but not sufficient; a CoT can be coherent but factually wrong.
- **Domain coverage** — training requires a curated multi-domain CoT dataset; MAgHARCM's runtime needs a domain-detection step before invoking ContextPRM.
- **Weighting-vs-selection** — WMV is one integration mode; alternative integration (best-of-N, beam search) is open future work.

## Source

- arXiv: 2509.24460.
- Venue: ICLR 2026 (Conference Paper), verified via iclr.cc/virtual/2026/poster/10011128; proceedings.iclr.cc paper ce186a37e63b37638ecd06dee6b9a355.
