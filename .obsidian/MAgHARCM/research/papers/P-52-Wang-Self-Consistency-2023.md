---
title: "P-52 — Wang et al. 2023 — Self-Consistency Improves Chain of Thought Reasoning in Language Models"
backlink: "[[1.0.0 P-52]]"
aliases:
  - "1.0.0 P-52"
  - "P-52"
  - "P-52-Wang-Self-Consistency-2023"
  - "P-52-Wang-Self-Consistency-2023"
  - "Wang-Self-Consistency-2023"
tags: [paper, reasoning, decoding, ensemble, verdict, [[1.0.0 PRIM-7]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-52 — Wang et al. 2023 — Self-Consistency Improves CoT]]

- **Authors**: Xuezhi Wang, Jason Wei, Dale Schuurmans, Quoc Le, Ed Chi, Sharan Narang, Aakanksha Chowdhery, Denny Zhou (Google Brain).
- **Venue / Year**: ICLR 2023; arXiv:2203.11171.
- **URL**: https://arxiv.org/abs/2203.11171
- **Anchors**: PRIM-7 (Verdict Panel), PRIM-21 (Migration Strategy Selection via try-and-fail registry).

## 1. Core Contribution

**Self-consistency**: replace greedy decoding with a sampling-based ensemble.
1. Sample N diverse reasoning paths from the model (temperature > 0).
2. Aggregate final answers via majority vote (or weighted vote).
3. Empirically improves CoT accuracy on arithmetic, commonsense, symbolic reasoning by 5-20 percentage points.

Key insight: the *answer distribution* matters more than the *single best sample*. Models that produce a wide answer distribution under diverse sampling also produce higher-accuracy majority answers.

## 2. Application in MAgHARCM

- The **PRIM-7 Verdict Panel** runs N optional checks (verdict_panel, mock_validator, implementation_agnostic, wasm_oracle) — this is structurally a self-consistency ensemble across DIFFERENT models/strategies, not N samples of one model. P-52's weighted-vote aggregation pattern applies.
- **PRIM-21 Registry.TryInOrder** try-and-fail pattern: each strategy's `Attempt` returns `AttemptResult` with success/partial/miss; the registry treats "vote" as "did this strategy succeed?" — analogous to a degenerate 1-of-N self-consistency over STRATEGIES rather than SAMPLES.
- Validates the role-flip gate's "ask the same question to two roles, vote on agreement" pattern.

## 3. Hop-1 References (papers cited by Wang et al.)

- Wei et al. (2022) — Chain-of-Thought Prompting (foundational).
- Cobbe et al. (2021) — Training Verifiers to Solve Math Word Problems (precursor verifier-based ensemble).
- Lewkowycz et al. (2022) — Minerva (math-focused CoT scaling).
- Brown et al. (2020) — GPT-3 / few-shot prompting (foundational).
- Chowdhery et al. (2022) — PaLM (large-model CoT reasoning baseline).

## 4. Hop-2 References (papers-cited-by-hop-1)

- Vaswani et al. (2017) — Attention Is All You Need (foundational transformer architecture referenced by GPT-3 and PaLM).
- Raffel et al. (2020) — T5 (referenced as a decoder-only baseline for CoT).
- Sanh et al. (2022) — T0 / Multitask Prompted Training (referenced as few-shot baseline).
- Cobbe et al. (2021) — GSM8K (math word problem benchmark).
- Hendrycks et al. (2021) — Measuring Massive Multitask Language Understanding (MMLU; cited as a benchmark).

## 5. Backlinks

- PRIM-7 (Verdict Panel): use self-consistency's majority-vote aggregation across N judges.
- PRIM-21 (Migration Strategy Selection): the try-and-fail registry IS a strategy-level self-consistency.
- Cross-ref: add P-52 to PRIM-7 + PRIM-21 rows in `Software-Archaeology-Lineage.md`.
