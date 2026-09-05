---
title: "P-56 — Zan et al. 2025 — Multi-SWE-bench: A Multilingual Benchmark for Issue Resolving"
backlink: "[[1.0.0 P-56]]"
tags: [paper, benchmark, multilingual, code-translation, issue-resolving, evaluation, [[1.0.0 PRIM-5]], [[1.0.0 PRIM-11]], [[1.0.0 PRIM-12]], [[1.0.0 PRIM-13]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-56 — Zan et al. 2025 — Multi-SWE-bench]]

- **Authors**: Daoguang Zan, Zhirong Huang, Wei Liu, Hanwu Chen, Linhao Zhang, Shulin Xin, Lu Chen, Qi Liu, Xiaojian Zhong, Aoyan Li, Siyao Liu, Yongsheng Xiao, Liangqiang Chen, Yuyu Zhang, Jing Su, Tianyu Liu, Rui Long, Kai Shen, Liang Xiang (ByteDance).
- **Venue / Year**: NeurIPS 2025 Datasets & Benchmarks Track; arXiv:2504.02605 (April 2025).
- **URL**: https://arxiv.org/abs/2504.02605
- **Anchors**: PRIM-5 (Test Suite Co-Translation & Synthesis), PRIM-11 (Implementation-Agnostic Test), PRIM-12 (Wasm Reference Oracle), PRIM-13 (Adversarial Test-Weakening Guard); the canonical multilingual benchmark for code-translation evaluation MAgHARCM needs.

## 1. Core Contribution

Multi-SWE-bench extends the original SWE-bench (Jimenez et al. 2024, Python-only) to **seven production-grade languages**: Java, TypeScript, JavaScript, Go, Rust, C, and C++ (plus Python). Total dataset size: **1,632 high-quality, expert-annotated instances** selected from 2,456 candidates by 68 expert annotators. Methodology:

1. **Per-language issue harvesting** — pull real-world GitHub issues from high-quality repositories in each target language, with verified fix commits and reproducible Docker evaluation environments.
2. **Three evaluation modes** — the benchmark evaluates models under three canonical agent scaffolds: **Agentless** (procedural, single-pass patch generation), **SWE-agent** (tool-calling autonomous agent), and **OpenHands** (the CodeAct-style multi-turn agent).
3. **Multi-SWE-RL companion release** — alongside the benchmark, the authors release a 4,723-instance RL training dataset and an open data-production pipeline for community expansion.
4. **Headline empirical findings** — GPT-4-class models score 20-30% on Java/Rust issue resolution; smaller open models (Qwen2.5-Coder 7B-Instruct, DeepSeek-Coder-V2-Lite) score 5-12%; the gap is largest for Rust and C++ where the standard-library surface area and macro system are most unfamiliar to LLMs.

**SLM-aware relevance for MAgHARCM**: this benchmark is the multilingual ground-truth MAgHARCM needs to evaluate the Translator + Validator pipeline. Prior work (Cassano et al. [[P-51]]; Wang et al. [[P-52]]; Liu et al. [[P-53]]) covers single-language benchmarks (HumanEval, MBPP) or single-edit evaluation (CanItEdit), but no existing benchmark gives per-language, real-issue-resolution scores that MAgHARCM can map directly onto its pipeline stages.

## 2. Application in MAgHARCM

- **PRIM-5 (Test Suite Co-Translation & Synthesis)** — Multi-SWE-bench's per-language test infrastructure is a direct production-quality reference for what MAgHARCM's `Validator` should generate. Specifically, the `Agentless` evaluation mode (procedural, no agentic tool-calls) matches MAgHARCM's compile-time-validated pipeline (no runtime agent loop). The benchmark's *test patching acceptance criterion* — does the model's patch make the existing test suite pass — is exactly what `ValidationReport.AllSuccess` tracks.
- **PRIM-11 (Implementation-Agnostic Test)** — the benchmark's evaluation methodology is implementation-agnostic by design: any agent that produces a patch that passes the verified test-suite is scored as correct, with no constraint on the agent's internal strategy. This validates MAgHARCM's implementation-agnostic test design.
- **PRIM-12 (Wasm Reference Oracle)** — Multi-SWE-bench does not include cross-language reference oracles (e.g., a Python reference implementation being checked against a Rust translation), but the methodology is extensible: the Docker evaluation environment can host a *legacy reference binary* and compare output against the new build. This is the conceptual template for adding a cross-language oracle to MAgHARCM's validator pipeline.
- **PRIM-13 (Adversarial Test-Weakening Guard)** — the benchmark's expert-annotated test patches are the reference standard for *test quality*. MAgHARCM's adversarial guard (which detects when a translator weakens tests to make them pass) can use Multi-SWE-bench's human-validated test diffs as a "test-quality training set" to learn what legitimate test changes look like vs. illegitimate weakening.
- **Per-language baseline scores** — Multi-SWE-bench's headline numbers (Qwen2.5-Coder 7B at ~12% on Rust, Phi-3-mini-128K at ~5%) tell us exactly where MAgHARCM's current SLM fleet will land without pipeline scaffolding. This is the calibration target: PRIM-3 (skeleton-first) + PRIM-22 (comprehension) + PRIM-31 (iterative navigator) exist precisely to lift SLMs from this floor.

## 3. Hop-1 References (papers cited by Zan et al.)

- Jimenez et al. (2024) — SWE-bench (the Python-only predecessor; Multi-SWE-bench's direct parent).
- Chen et al. (2021) — Codex / HumanEval (foundational code-generation benchmark referenced for evaluation methodology).
- Austin et al. (2021) — Program Synthesis with Large Language Models (APPS benchmark; referenced for problem-difficulty classification).
- Li et al. (2023) — StarCoder [[P-22]] in MAgHARCM lineage (referenced as a baseline SLM evaluated on Multi-SWE-bench).
- Roziere et al. (2023) — Code Llama (referenced as an instruction-tuned code baseline).
- Cassano et al. (2024) — Can It Edit? [[P-51]] in MAgHARCM lineage (referenced for multi-step code-edit evaluation).
- Yang et al. (2024) — SWE-agent (cited as one of the three evaluation scaffolds; canonical agentic code-edit agent).
- Wang et al. (2024) — OpenHands (CodeAct-style multi-turn agent; cited as the third evaluation scaffold).
- OpenAI et al. (2024) — GPT-4 technical report (referenced for GPT-4-class baseline scores).
- Liu et al. (2024) — DeepSeek-V2 (referenced for the DeepSeek-Coder-V2-Lite baseline).
- Anthropic (2024) — Claude 3.5 Sonnet system card (referenced for closed-model baseline scores); cross-link [[1.0.0 P-38]] (Anthropic sycophancy) for the related alignment-evaluation framework.

## 4. Hop-2 References (papers-cited-by-hop-1)

- Wang et al. (2023) — Self-Consistency [[P-52]] in MAgHARCM lineage (cited by SWE-agent as a reasoning technique).
- Liu et al. (2024) — Lost in the Middle [[P-53]] in MAgHARCM lineage (cited by SWE-bench for context-length trade-offs).
- Wei et al. (2022) — Chain-of-Thought Prompting (foundational; cited by OpenHands ReAct loop).
- Yao et al. (2022) — ReAct (the canonical reasoning+acting framework that SWE-agent and OpenHands both inherit).
- Hendrycks et al. (2021) — Measuring Massive Multitask Language Understanding (MMLU; cited as general benchmark).
- Touvron et al. (2023) — LLaMA 2 (foundational architecture for many of the open-model baselines).
- Schick et al. (2023) — Toolformer (referenced by SWE-agent for tool-calling methodology).
- Khattab et al. (2024) — DSPy (referenced by OpenHands as a prompt-optimization layer).
- Austin et al. (2021) — Program Synthesis with Large Language Models (foundational APPS benchmark; cited by HumanEval lineage).
- Peng et al. (2023) — RepoCoder (referenced by SWE-agent for iterative retrieval-augmented code completion).

## 5. Backlinks

- **PRIM-5** (Test Suite Co-Translation & Synthesis): the benchmark's test-patching acceptance criterion is `ValidationReport.AllSuccess`.
- **PRIM-11** (Implementation-Agnostic Test): the benchmark scores any agent that makes tests pass; matches MAgHARCM's design.
- **PRIM-12** (Wasm Reference Oracle): extensible via Docker-evaluated cross-language reference binaries.
- **PRIM-13** (Adversarial Test-Weakening Guard): expert-validated test diffs are the "legitimate-test-change" training set.
- **PRIM-22** (Four Phases Comprehension): per-language comprehension baselines come from Multi-SWE-bench scores.
- **PRIM-31** (Iterative Retrieval Refine): the Agentless vs. SWE-agent gap quantifies the value of iterative refinement.
- **Cross-ref**: add P-56 to PRIM-5, PRIM-11, PRIM-12, PRIM-13 rows in `Software-Archaeology-Lineage.md`. This is the **canonical multilingual evaluation anchor** for MAgHARCM's validator pipeline.
