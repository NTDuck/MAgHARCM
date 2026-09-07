---
title: "P-127 SLM-as-a-Judge: Improving Code Generation via Small Language Model-as-a-Judge"
backlink: "[[1.0.0 P-127]]"
aliases:
  - "1.0.0 P-127"
  - "P-127"
  - "P-127-SLM-as-a-Judge-ICSE-2026"
  - "SLM-as-a-Judge-ICSE-2026"
  - "crupi2026slmjudge"
tags: [paper, slm, code-generation, verifier, sft, execution-independent, "[[1.0.0 PRIM-7]]", "[[1.0.0 P-92]]", "[[1.0.0 P-98]]", "[[1.0.0 P-83]]", "[[1.0.0 P-106]]", "[[1.0.0 P-121]]", wave-18]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICSE 2026
---

# [[1.0.0 P-127]] SLM-as-a-Judge

## TL;DR

SLM-as-a-Judge **fine-tunes 0.5B-4B small language models as binary correct/incorrect code classifiers** via SFT on a 722-task Java dataset assembled from CoderEval, HumanEval, and MBPP. The trained SLM ranks 10 candidate solutions from an SLM generator, providing an **execution-independent verifier** for code-generation pipelines. Anchors `[[1.0.0 PRIM-7]]` Verdict Validation.

## Mechanism (Q2)

Three coupled components:

1. **Training data construction** — 722 Java tasks combining CoderEval (industrial), HumanEval (synthetic), and MBPP (introductory). Each task is paired with multiple candidate solutions labeled correct/incorrect.
2. **SLM fine-tuning via SFT** — Qwen2.5-Coder 0.5B/3B, Gemma-3 4B, and Llama-3.2 3B are SFT-fine-tuned as binary classifiers. The classifier reads a (problem, candidate) pair and emits correct/incorrect.
3. **Best-of-N ranking** — at inference, the SLM generator produces N=10 candidates and the trained SLM judge selects the top-ranked solution.

The execution-independence is the key contribution: most existing code verifiers (execution-based oracles, AST matchers) require running the code or hand-written specs; SLM-as-a-Judge only needs the problem statement and the candidate text. This makes the verifier deployable in environments where code execution is unsafe or unavailable.

## Anchoring (Q3)

| Primitive | Pre-wave-18 behaviour | SLM-as-a-Judge substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Execution-based oracles (run code, compare output) or hand-crafted AST matchers; large-LLM-as-judge (~70B+) | Trained 0.5B-4B SLM as binary classifier; execution-independent; uses problem text only |

## Hop-1 Citations

- `[[1.0.0 P-92]]` Lightman et al. 2023 *Let's Verify Step by Step* (PRM-style training of verifiers).
- `[[1.0.0 P-98]]` Brown et al. 2024 *Large Language Monkeys* (best-of-N sampling that consumes the verifier output).
- `[[1.0.0 P-83]]` Huang et al. 2024 *Self-Consistency for Code* (code-specialised self-consistency baseline).
- `[[1.0.0 P-106]]` / `[[1.0.0 P-121]]` BFCL AST-verification work (2025) — AST-based verifier baseline that SLM-as-a-Judge complements.

## Hop-2 Citations

- Zheng et al. 2023 *Judging LLM-as-a-Judge with MT-Bench and Chatbot Arena* (large-LLM-as-judge ancestor).
- Chen et al. 2024 *CodeJudge* (code-specialised LLM-as-judge for larger models; SLM-as-a-Judge extends to the 0.5B-4B band).

## MAgHARCM integration

- **YAML config key**: `verifier.slm_judge.enabled: true` (opt-in, default `false`); `verifier.slm_judge.model: <hf-identifier>`.
- **Implementation file**: `internal/verifier/slm_judge.go::RankCandidates` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-7]]`.

## Caveats

- **Language coverage**: the training set is Java; non-Java code may degrade. Re-training or fine-tuning on a polyglot corpus is required for multi-language deployment.
- **Candidate-pool sensitivity**: ranking accuracy is bounded by generator diversity — if all 10 candidates share the same bug, the SLM judge can rank them but cannot promote correctness.
- **Verifier-vs-generator parity**: the SLM judge must be at least as capable as the generator in its target domain. Deploying a stronger generator than the judge risks the judge endorsing incorrect code.

## Source

- arXiv:2602.11911 — SLM-as-a-Judge (ICSE 2026).
