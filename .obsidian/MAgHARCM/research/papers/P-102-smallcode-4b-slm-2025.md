---
title: "P-102 — SmallCode 4B SLM Coding (fp8.co 2025)"
backlink: "[[1.0.0 P-102]]"
tags: [paper, slm, code-generation, 4b-parameters, humaneval, [[1.0.0 PRIM-25]], [[1.0.0 PRIM-26]], hop-2]
---

# [[1.0.0 P-102 — SmallCode: 4B-Parameter SLM at 87% HumanEval]]

## Citation

*SmallCode: How a 4B-Parameter Model Achieves 87% on AI Coding Benchmarks*. fp8.co Articles, 2025.
URL: https://fp8.co/articles/SmallCode-AI-Coding-Agent-Small-LLM-Deep-Dive

## Summary

SmallCode demonstrates that aggressive specialisation compresses frontier-model coding capability into 4 billion parameters, achieving **87% on HumanEval**. The architectural pattern: a 4B SLM trained on a curated subset of code corpora (filtered for style consistency, version-pinned APIs, and unit-test co-occurrence) outperforms 70B-class general-purpose models on code-specific benchmarks while running at single-digit token/sec on commodity edge hardware.

Key tactics:
1. **Curated corpus distillation**: training data filtered by commit density, test coverage, and lint-clean ratio.
2. **Function-level tokenisation**: custom BPE merges favouring camelCase / snake_case boundaries.
3. **Chain-of-thought emulation via scaffolding**: external planner scaffolds CoT prompts; the SLM only generates final code, not reasoning.
4. **Tool-calling prefix-conditioning**: tool-call schema injected into context window via deterministic prefix template (BFCL v4 compatible).

## Relevance to MAgHARCM

- **Direct empirical anchor for SLM-era PRIM-25 Role-Flip Gate**: a 4B model with specialised training passes adversarial inspection when scaffolding is correct. This validates the gate's "reviewer + translator" split where the reviewer is a smaller, specialised model.
- **Token-budget evidence for PRIM-21 Strategy Selection**: strategy decisions should consider 4B-model token costs (1/15th of a 70B model), making sampling-heavy strategies (P-98 LLM Monkeys) viable at scale.
- **Curated-corpus lesson for PRIM-22 Comprehension**: the corpus-distillation mechanism suggests MAgHARCM's archaeology stage should weight recent, well-tested commits over raw coverage.
- **Tool-calling prefix conditioning (P-85 UNVERIFIED → P-102 verified)**: concrete evidence that 4B models can do function calling with the right scaffolding.

## References (hop-1)

- [[1.0.0 P-58]] Qwen2.5-Coder (4B variant) — same family, verified
- [[1.0.0 P-81]] Gorilla function-calling
- [[1.0.0 P-85]] Yue et al. 2025 function-calling (UNVERIFIED predecessor)
- [[1.0.0 P-95]] Code Llama 7B (larger cousin, instruction-tuned)

## References (hop-2)

- BFCL v4 leaderboard: https://gorilla.cs.berkeley.edu/blogs/8_berkeley_function_calling_leaderboard.html
- ToolBench: https://github.com/openbmb/toolbench
