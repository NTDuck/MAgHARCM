---
title: "P-51 — Cassano et al. 2024 — Can It Edit? Evaluating LLMs to Follow Code Editing Instructions"
backlink: "[[1.0.0 P-51]]"
tags: [paper, code-editing, benchmark, repair, slm, [[1.0.0 PRIM-22]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-51 — Cassano et al. 2024 — Can It Edit?]]

- **Authors**: Federico Cassano, Luisa Li, Akul Sethi, Noah Shinn, Maxwell Brennan, et al. (Northeastern University / Roblox / Hugging Face).
- **Venue / Year**: COLM 2024 (Conference on Language Modeling); arXiv:2402.14304.
- **URL**: https://arxiv.org/abs/2402.14304
- **Anchors**: PRIM-22 (Comprehension Recognition), PRIM-25 (RoleFlip Gate), translator repair path.

## 1. Core Contribution

"Can It Edit?" (CanItEdit) benchmark evaluates LLMs on **instructional code editing**:
- 218 hand-crafted Python editing problems; each pairs a function with a natural-language edit directive.
- Evaluates multi-step edits, bug-fix edits, refactor edits, and feature-add edits.
- Finding: smaller instruction-tuned models (3B-7B) drop 20-40% vs. 70B-class models on multi-step edits; Qwen2.5-Coder 7B-Instruct trails GPT-4 by ~15 points on this benchmark but exceeds Llama-3-70B on simple single-edit cases.
- Recommendation: **decompose multi-step edits into single-step edits** before prompting — directly relevant to the Translator's chunked repair loop.

## 2. Application in MAgHARCM

- Validates the Translator's **repair-mode** prompt template (`TranslatorRepairPromptTemplate`) which decomposes the repair into "minimal-change diff" + "isolated edit" subtasks.
- Informs the PRIM-25 RoleFlip Gate's `Model is nil` sentinel — the gate must produce a structured "this edit is feasible / not feasible" response, not free-form prose.
- Recommends explicit pre-prompt with the original code structure, supporting the analyzer/planning pass that precedes translation.

## 3. Hop-1 References (papers cited by Cassano et al.)

- Chen et al. (2021) — Codex (HumanEval [[P-22]] lineage).
- Austin et al. (2021) — Program Synthesis (foundational benchmark).
- Jimenez et al. (2024) — SWE-bench (canonical agentic code-edit benchmark; multi-file context).
- Liu et al. (2024) — Lost in the Middle [[P-53]] in MAgHARCM lineage.
- Roziere et al. (2023) — Code Llama (referenced as instruction-tuned baseline).
- Li et al. (2023) — StarCoder2 [[P-22]].

## 4. Hop-2 References (papers-cited-by-hop-1)

- Wei et al. (2023) — Simple synthetic data reduces sycophancy in large language models (cited by SWE-bench for benchmark-design choices).
- Wang et al. (2023) — Self-Consistency [[P-52]].
- Liu et al. (2024) — Lost in the Middle [[P-53]].
- Touvron et al. (2023) — LLaMA 2 (foundational model referenced by Code Llama).
- Christakopoulou et al. (2024) — TRIE / RLHF for code (cited by SWE-bench evaluation methodology).
- Yao et al. (2022) — ReAct (cited by SWE-bench agents as canonical reasoning+acting framework).

## 5. Backlinks

- Translator repair mode: PRIM-22 (Comprehension Recognition) + PRIM-25 (RoleFlip Gate).
- Adversarial test weakening guard (PRIM-13): the benchmark's "remove this test" anti-pattern motivates the guard.
- Cross-ref: add P-51 to PRIM-22 + PRIM-25 rows in `Software-Archaeology-Lineage.md`.
