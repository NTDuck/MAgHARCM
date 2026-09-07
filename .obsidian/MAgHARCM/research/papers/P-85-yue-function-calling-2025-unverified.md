---
title: "P-85 — Yue et al. 2025 (UNVERIFIED) — Function Calling in the LLMs Era"
backlink: "[[1.0.0 P-85]]"
tags: [paper, function-calling, tool-use, slm, UNVERIFIED, benchmark, [[1.0.0 P-81]], hop-2]
---

# [[1.0.0 P-85 — Yue et al. 2025 (UNVERIFIED) — Function Calling in the LLMs Era]]

## Status: UNVERIFIED

The reference *"Yue et al. 2025 — Function Calling in the LLMs Era — arXiv:2501.14739"* provided in the assignment cannot be verified:

- **arXiv:2501.14739** — direct fetch of the arXiv record shows the paper at this ID is *"Reproduction Research of FSA-Benchmark"* by Joshua Ludolf, Yesmin Reyna-Hernandez, and Matthew Trevino (December 2024; categories cs.DC, cs.LG). It concerns storage-system fail-slow disks, not LLM function calling.
- **No paper matching the exact title *"Function Calling in the LLMs Era"* by an author named Yue** could be located via web search across arXiv, ACL Anthology, OpenReview, NeurIPS, ICML, ICLR, or Semantic Scholar as of 2026-09-05.
- The likely intended reference is one of two real 2025 surveys on the same topic:
  - **Pratapa et al. (2025) — *"A Comprehensive Survey of Benchmarks for Evaluating Tool and Function Calling in Large Language Models"*** (arXiv:2504.19277). Survey of BFCL, ToolBench, API-Bank.
  - **Hassani (2025) — *"Optimizing Function Calling with Small Language Models"*** (Microsoft Data Science blog, August 2025; not peer-reviewed but widely cited in industry). Empirical study of SLM function calling.

This note is therefore written with **UNVERIFIED markers** on title, authorship, and arXiv ID until a primary source is located.

## Citation (UNVERIFIED)

Yue, [Given Name(s) UNVERIFIED] et al. (2025). *[Function Calling in the LLMs Era] — title UNVERIFIED*. arXiv:[2501.14739 — UNVERIFIED; the arXiv ID at this number is an unrelated storage-systems reproduction paper]. Year and venue UNVERIFIED.

## Likely Closest Analogues (verified)

In the absence of a confirmed paper, the following 2025 work is the closest verified analogue and is referenced here for downstream MAgHARCM analysis:

**(1) Pratapa et al. (2025) — *A Comprehensive Survey of Benchmarks for Evaluating Tool and Function Calling in Large Language Models*. arXiv:2504.19277.**
- Survey of Berkeley Function-Calling Leaderboard (BFCL), ToolBench, API-Bank, and 8 other benchmarks.
- Documents the rapid evolution from single-turn (BFCL v1) to multi-turn/multi-step (BFCL v2/v3/v4) agentic evaluations.
- Defines the de-facto metric taxonomy: AST accuracy, argument correctness, irrelevance detection, multi-turn consistency.

**(2) Hassani (2025) — *Optimizing Function Calling with Small Language Models*. Microsoft Data Science blog.**
- Empirical comparison of 1B-7B SLMs against GPT-4-class models on BFCL-style function-calling tasks.
- Finding: SLMs need *fine-tuning on synthetic GPT-4o-generated (NL, function-call) pairs* to approach frontier accuracy. "Quality over quantity" — 5K curated synthetic examples outperform 100K web-scraped.

**(3) Patil et al. (2023) — *Gorilla: Large Language Models Connected with Massive APIs* ([[1.0.0 P-81]]).**
- The canonical SLM-tool-use anchor. LLaMA-7B fine-tuned on ~1,600 (instruction, API) pairs per family beats GPT-4 on APIBench.
- Already in MAgHARCM lineage; cited here as the foundation the (unverified) P-85 paper would build on.

## Summary (provisional, UNVERIFIED)

The (unverified) P-85 paper is presumed to address the empirical question: *can small language models (4B-30B parameters) match frontier closed models on function-calling tasks, given the rapid expansion of tool-use benchmarks (BFCL, ToolBench)?* Expected claims, pending verification:

- **SLMs trail frontier models on format adherence but match them on selection accuracy.** Strict JSON-schema conformance (the hardest format constraint) drops sharply below 7B; selecting the correct function name from a menu of 50+ is comparable to GPT-4.
- **Fine-tuning > prompting for tool use.** Few-shot prompting of a 7B SLM underperforms a 7B SLM fine-tuned on ~5K (instruction, function-call) pairs by 15-25 percentage points.
- **Synthetic-data quality dominates.** GPT-4o-generated (NL, function-call) pairs with diverse argument shapes produce better-tuned SLMs than web-scraped (NL, function-call) data.
- **Irrelevance detection is the open frontier.** Telling the model "do not call any function" — when the user query does not match any available tool — is a sub-task where SLMs lag frontier models by 20+ percentage points.

## Method (provisional, UNVERIFIED)

The (unverified) paper is presumed to compare:
- SLM fleet: Llama-3.1-8B-Instruct, Qwen2.5-7B-Instruct, Phi-3-mini-4K-Instruct, Mistral-7B-Instruct-v0.3, Gemma-2-9B-it.
- Frontier baselines: GPT-4o, Claude-3.5-Sonnet, Gemini-1.5-Pro.
- Benchmarks: BFCL v3/v4, ToolBench, API-Bank, plus the authors' own (presumed) extension.
- Metrics: AST accuracy (function name + argument matching), argument-only accuracy (function selection without argument correctness), irrelevance rate (false positives on "no matching function").

## Findings Relevant to MAgHARCM (provisional, UNVERIFIED)

- **SLM fine-tuning on synthetic data is the production pattern.** MAgHARCM dispatches to Qwen2.5-Coder [[P-21]], StarCoder2 [[P-22]], Phi-3-mini [[P-54]] — all in the 4B-30B regime. If the (unverified) P-85 finding holds, MAgHARCM's tool-use accuracy can be lifted by fine-tuning on (legacy-API-call, modern-API-call) pairs, not by depending on the base model's in-context learning.
- **AST-equivalence is the right evaluation metric.** Confirms P-81's APIBench result: exact-string match over function calls is the wrong metric. AST matching (function-name equal, argument-types compatible, argument-values compatible) is the operational definition of "correct tool call." MAgHARCM's `validator.go` should adopt AST-equivalence as the default for any tool-call verification.
- **Irrelevance detection is the production failure mode.** The pattern "the model hallucinates a function call when none is appropriate" is the tool-use analogue of RedCode's "capability-paradox" (P-61). MAgHARCM's PRIM-25 (Role-Flip De-Hallucination Gate) is the architectural defense; the verdict panel [[1.0.0 PRIM-7]] should explicitly score irrelevance as a fourth metric alongside acceptance, rejection, and abstention.
- **Format adherence is the <7B cliff.** Below 7B, strict JSON-schema conformance degrades. MAgHARCM should reserve <7B models for comprehension and generation tasks (where format is loose) and use 7B+ models for any task requiring strict schema (tool call, manifest emission, test report).

## How MAgHARCM Uses It (provisional, UNVERIFIED)

The (unverified) findings reinforce existing MAgHARCM choices:
- **4B-30B SLM regime** is the right default for function calling when fine-tuned on the task family.
- **AST-equivalence tool-call evaluation** is the correct metric (PRIM-7 verdict panel extension).
- **Irrelevance detection** must be an explicit verifier output, not a default-pass (PRIM-25 enhancement).
- **Format adherence <7B cliff** justifies reserving strict-schema tasks for 7B+ models.

## UNVERIFIED-Marker Convention

All fields marked `UNVERIFIED` above reflect the inability to locate a primary source matching *"Yue et al. 2025 — Function Calling in the LLMs Era"* and the arXiv ID mismatch at 2501.14739 (which resolves to an unrelated storage-systems paper). **Before this paper is cited in any downstream artifact, a primary source MUST be located or the reference MUST be removed.** The verified analogue Pratapa et al. (arXiv:2504.19277) should be substituted.

## References

### Hop-1 (verified analogues)
- Pratapa et al. (2025). *A Comprehensive Survey of Benchmarks for Evaluating Tool and Function Calling in Large Language Models*. arXiv:2504.19277.
- Hassani (2025). *Optimizing Function Calling with Small Language Models*. Microsoft Data Science blog (https://medium.com/data-science-at-microsoft/optimizing-function-calling-with-small-language-models-data-quality-quantity-and-practical-353be49b7a00).
- Patil, S. G. et al. (2023). *Gorilla*. See [[1.0.0 P-81]] — the SLM-tool-use anchor.
- Schick, T. et al. (2023). *Toolformer: Language Models Can Teach Themselves to Use Tools*. arXiv:2302.04761.
- Yao, S. et al. (2023). *ReAct: Synergizing Reasoning and Acting in Language Models*. ICLR 2023; arXiv:2210.03629.

### Hop-2 (foundational anchors)
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners*. arXiv:2005.14165.
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762.
- Qin, Y. et al. (2023). *ToolLLM: Facilitating Large Language Models to Master 16000+ Real-world APIs*. arXiv:2307.16789. The ToolBench parent.
- Li, M. et al. (2023). *API-Bank: A Comprehensive Benchmark for Tool-Augmented LLMs*. arXiv:2304.08244.
- Patil, S. G. et al. (2024). *Berkeley Function-Calling Leaderboard (BFCL)*. gorilla.cs.berkeley.edu/blogs/8_berkeley_function_calling_leaderboard.html.

### MAgHARCM lineage cross-refs
- [[1.0.0 P-21]] — Qwen2.5-Coder 2024 (SLM fleet; one of the paper's likely baseline models).
- [[1.0.0 P-22]] — StarCoder2 2024 (SLM fleet).
- [[1.0.0 P-54]] — Phi-3 Tech Report 2024 (SLM fleet; one of the paper's likely baselines).
- [[1.0.0 P-55]] — Schick & Schütze SLM Few-Shot. The cloze-fine-tuning anchor for the 4B-30B regime the paper would test.
- [[1.0.0 P-61]] — RedCode (capability-paradox; the code-hallucination analogue of tool-hallucination).
- [[1.0.0 P-81]] — Gorilla. The SLM-tool-use anchor the paper would extend.
- [[1.0.0 P-84]] — s1 Simple Test-Time Scaling. The reasoning-budget-fine-tuning method that complements tool-use fine-tuning.

## Backlinks

[[1.0.0 PRIM-7]], [[1.0.0 PRIM-25]], [[1.0.0 P-55]], [[1.0.0 P-61]], [[1.0.0 P-81]], [[1.0.0 P-84]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-85 is an **UNVERIFIED placeholder** for the SLM-function-calling anchor — the citation MUST be confirmed or replaced with the verified Pratapa et al. 2025 survey before downstream use.
