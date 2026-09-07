---
title: Qwen2.5-Coder Technical Report
bibkey: p21_qwen2_5coder
tags: [paper, code-llm, [[PRIM-23]], [[PRIM-25]], [[PRIM-26]], hop-1]
---

# [[1.0.0 P-21]] Qwen2.5-Coder Technical Report

**Authors**: Qwen Team (An Yang, Baosong Yang, Beichen Zhang, Binyuan Hui, Bo Zheng, Bowen Yu, Chengpeng Li, Chengen Huang, Dayiheng Liu, Fei Huang, Haoran Wei, Huan Lin, Jian Yang, Jianhong Tu, Jianwei Zhang, Jianxin Yang, Jiaxi Yang, Jing Zhou, Junyang Lin, Kai Dang, Keming Lu, Keqin Bao, Kexin Yang, Lei Yu, Lianghao Deng, Mei Li, Mingfeng Xue, Mingze Li, Pei Zhang, Peng Wang, Qin Zhu, Rui Men, Ruize Gao, Shixuan Liu, Shuang Luo, Tianhao Li, Tianyi Tang, Wenbiao Yin, Xingzhang Ren, Xinyu Wang, Xinyu Zhang, Xuancheng Ren, Yang Fan, Yang Su, Yichang Zhang, Yu Qiao, Yuxuan Cai, Zhenyu Gu, Zhiyuan Liu, Zonghong Dai — Alibaba Qwen Team)
**Year**: 2024
**Venue**: arXiv preprint (cs.CL / cs.SE / cs.AI)
**eprint / DOI**: [[arXiv-2409.12186]]
**Cited by**: [[primitives/INDEX]] entry [[PRIM-23]] (Chunked Translation), [[PRIM-25]] (Communicative-De-hallucination Role-Flip Gate), [[PRIM-26]] (Symbol-Aware Navigator)

## Summary

[[Qwen-2024-Qwen2.5-Coder]] is the technical report behind the Qwen2.5-Coder family of code-specialised LLMs (0.5B / 1.5B / 3B / 7B / 14B / 32B parameters). The base model continues pre-training of Qwen2.5 on 5.5 trillion tokens of curated code, math, and general text; instruction tuning then injects synthetic code-instruction data generated via a multi-stage pipeline (code execution filtering, hard-example mining, and self-distilled solver traces). The report documents three claims of empirical interest to MAgHARCM: (i) the 7B-Plus models match GPT-4o on HumanEval, MBPP, and LiveCodeBench at the time of release; (ii) a long-context window of 128K tokens — relevant to repository-scale translation where one fragment may reference dozens of source modules; and (iii) the family includes Qwen2.5-Coder-32B-Instruct, which the MAgHARCM eval harness can run via Ollama at high quality. The report also reports code-completion, code-repair, and code-reasoning benchmarks separately, so practitioners can pick the smallest variant that meets a target capability.

## Relevance to MAgHARCM

Qwen2.5-Coder is the actual model family that MAgHARCM evaluates against in [[PRIM-23]] (Chunked Translation) — `internal/agents/chunked_translator.go` and the translator prompt templates assume a Qwen2.5-Coder-4B-Instruct (GGUF Q4_K_XL) or larger local model on the coding-model slot of `state.Config`. The 128K context window informs the upper bound on the bounded-context prefix the [[PRIM-26]] Symbol-Aware Navigator hands to the translator (current cap is 4 KB per excerpt; Qwen2.5-Coder's window means we could enlarge this if the profile demands it). The role-flip [[PRIM-25]] gate ([[PRIM-25 Communicative-De-hallucination Role-Flip Gate]]) uses Qwen2.5-Coder as its reviewer-model with a system prompt that inverts the role to "critical reviewer who must find at least one bug"; the report's coverage of code-reasoning and code-debug benchmarks supports the assumption that the same model is competent at the review role. The chat template and tokenizer output are pinned via `internal/consts/consts.go` constants; Qwen2.5-Coder's permissive Apache-2.0 license is the licensing reason MAgHARCM can ship pre-quantised GGUF weights inside the eval artefact set.

## Hop-1 References

- [[DeepSeek-Coder-V2-2024]] — competing code-LLM family with comparable HumanEval numbers; MAgHARCM evaluates both in [[sec_eval]] ablation.
- [[StarCoder2-2024]] — predecessor open code-LLM family with the same permissive licensing posture; data-pipeline contrast informs the Qwen2.5-Coder synthesis filter design.
- [[Code-Llama-2023]] — Meta's code-LLM release; cited for benchmark comparison in the Qwen2.5-Coder report.
- [[RepoCoder-2024]] — retrieval-augmented code completion; pairs naturally with Qwen2.5-Coder's long context window for repository-scale completions.
- [[HyperAgent-2026]] — multi-agent SE framework that uses Qwen2.5-Coder-class backbones for its translator agent.

## Hop-2 Anchors (software-archaeology lean)

- [[Feathers-2004-WELC]] — characterisation tests and dependency-breaking seams are the offline analogue of Qwen2.5-Coder's hard-example mining: both isolate the cases the model (or the developer) cannot yet handle.
- [[Rajlich-1997]] — concept-locator analysis motivates why long-context code-LLMs are useful: a single prompt can carry the concept → location map the offline analysis produces.
- [[Kazman-Cai-2024]] — architectural recovery predicts that Qwen2.5-Coder's recall on 128K-token inputs is best when the context is organised by L1/L2/L3 stability layers, mirroring the [[PRIM-21]] strategy-selection logic.
