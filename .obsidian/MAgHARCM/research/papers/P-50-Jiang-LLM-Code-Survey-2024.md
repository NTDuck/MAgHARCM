---
title: "P-50 — Jiang et al. 2024 — A Survey on Large Language Models for Code Generation"
backlink: "[[1.0.0 P-50]]"
aliases:
  - "1.0.0 P-50"
  - "P-50"
  - "P-50-Jiang-LLM-Code-Survey-2024"
  - "P-50-Jiang-LLM-Code-Survey-2024"
  - "Jiang-LLM-Code-Survey-2024"
tags: [paper, code-generation, survey, slm, [[1.0.0 PRIM]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-50 — Jiang et al. 2024 — Survey on LLMs for Code Generation]]

- **Authors**: Juyong Jiang, Fan Wang, Jiasi Shen, Sungju Kim, Sunghun Kim (HKUST).
- **Venue / Year**: arXiv:2406.00515 (June 2024).
- **URL**: https://arxiv.org/abs/2406.00515
- **Anchors**: general PRIM for prompt engineering + SLM-aware generation; cross-cuts PRIM-1 (Analyzer), PRIM-3 (Planning), PRIM-6 (Translator).

## 1. Core Contribution

Comprehensive survey of LLMs for code generation across the full stack:
1. **Model families** — encoder-only (CodeBERT), decoder-only (CodeGen, Code Llama, StarCoder, Qwen-Coder), encoder-decoder (PLBART, Codet5).
2. **Prompting strategies** — zero-shot, few-shot, chain-of-thought (CoT), structured CoT (SCoT), retrieval-augmented (RAG), agentic (Toolformer, ChatDev).
3. **Benchmarks** — HumanEval, MBPP, APPS, CodeContests, DS-1000, RepoBench.
4. **SLM-relevant insight** — smaller decoder-only models (≤7B) reach 80–90% of large-model performance on HumanEval-style tasks when paired with few-shot prompts and structural output constraints.

## 2. Application in MAgHARCM

- The survey's SLM findings inform the SLM-aware prompt preamble (`compiletime.SLMPromptContractPreamble`).
- The structured-output recommendation informs the JSON-section parsing used by every agent's `RenderPromptTemplate` consumer.
- The benchmark table clarifies what `ValidationReport.AllSuccess` should track against (HumanEval, MBPP-equivalent for translated repos).

## 3. Hop-1 References (papers cited by Jiang et al.)

- Chen et al. (2021) — CodeX / HumanEval ([[P-22 starcoder2]], [[P-21 qwen2_5coder]] adopt the HumanEval task format).
- Nijkamp et al. (2023) — CodeGen family (cited as canonical decoder-only code LLM).
- Roziere et al. (2023) — Code Llama (cited as the open-source instruction-tuned reference).
- Li et al. (2023) — StarCoder ([[P-22]] in MAgHARCM lineage).
- Wang et al. (2023) — Self-Consistency [[P-52]] in MAgHARCM lineage.
- Roziere et al. (2023) — LLaMA2 backbone family.
- Austin et al. (2021) — Program Synthesis with Large Language Models (foundational few-shot benchmark).
- Cheng et al. (2023) — Codet5+ foundation model survey.
- Liu et al. (2024) — Lost in the Middle [[P-53]] in MAgHARCM lineage.

## 4. Hop-2 References (papers-cited-by-hop-1)

- Wei et al. (2022) — Chain-of-Thought Prompting (foundational CoT, cited by Wang et al. → [[P-52]]).
- Chowdhery et al. (2022) — PaLM (cited by CoT paper).
- Hendrycks et al. (2021) — Measuring Coding Challenge Competence With APPS (cited by HumanEval lineage).
- Touvron et al. (2023) — LLaMA 2 (cited by Code Llama paper).
- Bommasani et al. (2022) — On the Opportunities and Risks of Foundation Models (cited by survey on SLM scaling laws).
- Khattab et al. (2024) — DSPy (cited by RAG section — supports `optimise()` loops over prompts).
- Sanh et al. (2022) — Multitask Prompted Training Enables Zero-Shot Task Generalization (T0 cited as few-shot baseline).

## 5. Backlinks

- Software archaeology / chunked translation: anchors PRIM-31 (IterativeNavigator) hybrid strategy.
- Prompt engineering: anchors `compiletime.SLMPromptContractPreamble` directly.
- Lineage cross-ref: add P-50 to general PRIM row in `Software-Archaeology-Lineage.md`.
