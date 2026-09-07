---
title: "P-81 — Patil et al. 2023 — Gorilla: Large Language Models Connected with Massive APIs"
backlink: "[[1.0.0 P-81]]"
aliases:
  - "1.0.0 P-81"
  - "P-81"
  - "P-81-Patil-Gorilla-2023"
  - "P-81-Patil-Gorilla-2023"
  - "Patil-Gorilla-2023"
tags: [paper, tool-use, api-calling, retrieval-augmented, slm, function-calling, hallucination, [[1.0.0 PRIM-7]], [[1.0.0 PRIM-25]], [[2.0.0 MAgHARCM]], hop-1]
---

# [[1.0.0 P-81 — Patil et al. — Gorilla]]

## Citation

Patil, S. G., Zhang, T., Wang, X., & Gonzalez, J. E. (2023). *Gorilla: Large Language Models Connected with Massive APIs*. arXiv:2305.15334 (May 2023); NeurIPS 2024. UC Berkeley Sky Computing Lab. URL: https://arxiv.org/abs/2305.15334 ; project page: https://gorilla.cs.berkeley.edu/

## Summary

Gorilla addresses a specific failure mode of LLM tool use at API granularity: when asked to call an API (HuggingFace model, TorchHub function, TensorHub module), frontier LLMs such as GPT-4 hallucinate plausible-looking but non-existent endpoints, wrong argument shapes, and outdated versions. Gorilla fixes this by fine-tuning a LLaMA-7B base model on a self-instruct dataset of (instruction, API-call) pairs drawn from a curated, deduplicated corpus of public API documentation, and shows that the resulting 7B model **outperforms GPT-4 on API-call accuracy while remaining open-weights**. This is the strongest 2023 evidence that a 4B-30B SLM, properly trained, can match or beat a 100x-larger frontier model on a domain-anchored tool-use task — the foundational SLM-tool-use anchor for MAgHARCM's tool-use design.

Three contributions matter most for MAgHARCM:

1. **Gorilla model.** A LLaMA-7B base fine-tuned on ~1,600 (instruction, API) pairs per API family (HuggingFace, TorchHub, TensorHub). Surpasses GPT-4 on APIBench functional-correctness accuracy, with comparable or better robustness to API documentation changes (Gorilla in zero-shot evals beats GPT-4 even when GPT-4 is given retrieved docs).

2. **Retriever-Aware Training (RAT).** During fine-tuning, Gorilla is conditioned on retrieved API documentation for ~70% of examples and on its parametric memory for the rest. This trains the model to (a) prefer retrieved context when available and (b) still produce correct calls when retrieval is noisy or absent. RAT is the **retrieval-conditioned prompt pattern** MAgHARCM's optional-checks agent should adopt when consulting external knowledge.

3. **APIBench.** A benchmark of 1,700+ prompts across HuggingFace, TorchHub, TensorHub APIs, with AST sub-tree matching as the functional-correctness metric (rather than exact-string match, which penalises correct-but-equivalent calls). APIBench is the methodological template for any "is this tool-call correct?" evaluation in MAgHARCM.

## Method

Patil et al. build a self-instruct pipeline:

- **Corpus construction.** Crawl public API documentation for HuggingFace, TorchHub, TensorHub. Deduplicate by API signature and parameter list (so the same function under multiple aliases is counted once). This produces a corpus of ~1,600 unique APIs.
- **Instruction synthesis.** Use GPT-4 to generate natural-language instructions that would call each API, conditioned on the API's documentation block. Filter for diversity.
- **Reference solution.** For each (instruction, API) pair, store the canonical, executable API call as the reference.
- **Gorilla fine-tuning.** Start from LLaMA-7B. Fine-tune on three regimes: zero-shot (no retrieved docs), retrieval-augmented (instruction + retrieved API docs → call), and retrieval-only fine-tuning (the ablation that ablates parametric memory). RAT mixes zero-shot and retrieval-augmented examples in a single fine-tune.
- **Evaluation.** APIBench prompts the model (Gorilla, GPT-4, GPT-3.5, Claude, LLaMA-7B base) for an API call, parses the call, and computes AST sub-tree match against the reference. Match is functional, not lexical — equivalent calls with reordered keyword args count as correct.

Key results:
- Gorilla with retrieval beats GPT-4 with retrieval on APIBench accuracy (84% vs 66% on HuggingFace subset).
- Gorilla zero-shot beats GPT-4 zero-shot by a smaller margin; the bulk of Gorilla's edge comes from learning to *use* retrieved context.
- Gorilla outperforms GPT-4 on robustness to API-version updates: when an API's signature changes between training and eval, Gorilla degrades less because it has learned to follow the retrieved signature rather than memorise the old one.

## Findings Relevant to MAgHARCM

- **SLM-grade tool use is achievable with domain fine-tuning.** Gorilla's headline result — a 7B open model beating GPT-4 on a domain-anchored tool-use benchmark — is the existence proof for MAgHARCM's choice to drive the comprehension-and-translation pipeline with 4B-30B SLMs (Qwen2.5-Coder [[P-21]], StarCoder2 [[P-22]], Phi-3-mini [[P-54]]) rather than frontier models. The caveat: Gorilla's win comes from *fine-tuning on a curated corpus*. MAgHARCM cannot fine-tune on every legacy codebase it encounters, so the lesson is "use SLMs where the corpus is curated and the prompt structure is constrained" — which is exactly the comprehension+translation regime.

- **Retriever-Aware Training is the SLM analogue of MAgHARCM's optional-checks agent.** Gorilla's RAT regime (70% retrieval-conditioned, 30% parametric) is structurally identical to the optional-checks pattern in `internal/agents/optional_checks.go`: when external knowledge is available, prefer it; when not, fall back to the model's trained priors. The 30/70 split is a calibration point — MAgHARCM's optional-checks agent should bias toward retrieval when confidence in parametric memory is low.

- **API hallucinations are the tool-use analogue of code hallucinations.** P-61's RedCode finding (capability-paradox: stronger models produce more risky code) generalises to tool-use: GPT-4 hallucinates API endpoints with high confidence. MAgHARCM's [[1.0.0 PRIM-7]] (Multi-Agent Verdict Panel) and [[1.0.0 PRIM-25]] (Role-Flip De-Hallucination Gate) are the architectural defenses against both code and tool-call hallucinations. Gorilla demonstrates that the alternative — fine-tune on a curated corpus — also works, but only when the corpus is available.

- **AST sub-tree match is the right tool-call evaluation metric.** Exact-string match penalises correct-but-reordered calls. MAgHARCM's validator (`internal/agents/validator.go`) should adopt AST-equivalence (for source-language calls) or JSON-schema-equivalence (for REST calls) as the correctness metric for any tool-call evaluation, not lexical equality.

- **APIBench is the benchmark template.** MAgHARCM does not have a Gorilla-equivalent benchmark for legacy-code modernisation, but the structure transfers: build a curated corpus of "what a correct call looks like" for each API family (in MAgHARCM's case, each language pair (legacy_target, modern_target)), generate (instruction, call) pairs, and evaluate AST-equivalence. This is the recommended pattern for future MAgHARCM benchmark work.

## How MAgHARCM Uses It

- The **4B-30B SLM regime** in MAgHARCM's pipeline (Qwen2.5-Coder [[P-21]], StarCoder2 [[P-22]], Phi-3-mini [[P-54]]) is justified by Gorilla: domain fine-tuning of a 7B model can beat a 100x-larger frontier model on a well-scoped tool-use task. The scope of legacy-code modernisation is well-defined (translate Java/C++/Python-2 to Rust/Go), so the regime is appropriate.

- The **optional-checks agent's external-knowledge lookup** mirrors Gorilla's retriever. When the model is uncertain about an API signature (a legacy crate's function, a deprecated Java method), the optional-checks agent retrieves the documentation and re-conditions the model on it. The 30/70 RAT split is the calibration: prefer retrieval when the model's prior probability on the API is below a threshold.

- The **verdict panel's tool-call validation** (`internal/agents/verdict_panel.go`) uses AST-equivalence rather than exact-string match for any source-language tool calls (test invocations, build commands, language-server queries). Gorilla's AST sub-tree metric is the template.

- **Future work (not implemented):** a MAgHARCM-7B fine-tune on a curated corpus of legacy-to-modern translation examples, modelled on Gorilla. This is a longer-term possibility; not currently scoped.

## Cite Chain

This paper sits on a hop-1 cite chain with three papers cited inline in the assignment:

- **Schick et al. 2023 — Toolformer** (arXiv:2302.04761). Predates Gorilla by three months. Toolformer is the *self-supervised* approach to tool use: a 7.7B model learns to insert API calls into its own training data based on a "does this API result reduce next-token loss?" filter. The fine-tune requires only a handful of human demonstrations per API. Gorilla is the *supervised-fine-tuning* approach: it requires a curated (instruction, API-call) corpus but is more controllable. Both papers reach the same conclusion (SLMs can be taught to use tools) from opposite starting points.

- **Yao et al. 2023 — ReAct** (arXiv:2210.03629, ICLR 2023). Predates both. ReAct is the *prompting-only* approach: a frozen LLM is given 1-2 in-context examples of interleaved `Thought:` / `Action:` / `Observation:` traces and asked to follow the pattern. No fine-tuning, no API-specific training. ReAct is the proof-of-concept that tool use does not require training; Gorilla and Toolformer are the demonstrations that fine-tuning improves it. Together, the three papers form the canonical 2023 tool-use trilogy: prompt → self-supervise → supervised-fine-tune.

- **Patil et al. 2023 — Gorilla** (this paper). The supervised-fine-tuning endpoint of the chain.

## References

### Hop-1 (Patil cites)
- Schick, T. et al. (2023). *Toolformer: Language Models Can Teach Themselves to Use Tools*. arXiv:2302.04761. The self-supervised prior art.
- Touvron, H. et al. (2023). *LLaMA: Open and Efficient Foundation Language Models*. arXiv:2302.13971. Gorilla's base model.
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners*. arXiv:2005.14165. The in-context-learning foundation that ReAct and Gorilla build on.
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762. Transformer architecture.
- Ouyang, L. et al. (2022). *Training Language Models to Follow Instructions with Human Feedback* (InstructGPT). arXiv:2203.02155. The instruction-tuning paradigm Gorilla inherits.
- Wang, B. et al. (2023). *Self-Instruct: Aligning Language Models with Self-Generated Instructions*. arXiv:2212.10560. Gorilla's instruction-synthesis pipeline.

### Hop-2 (papers-cited-by-hop-1, relevant)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. The reasoning-trace pattern that ReAct interleaves with actions.
- Khot, T. et al. (2023). *Decomposed Prompting: A Modular Approach for Solving Complex Tasks*. See [[1.0.0 P-62]] (Decomposed Prompting 2022). ReAct-style decomposition as a precursor.
- Karpukhin, V. et al. (2020). *Dense Passage Retrieval for Open-Domain Question Answering*. arXiv:2004.04906. The retrieval backbone for Gorilla's RAT.
- Schick, T. & Schütze, H. (2021). *It's Not Just Size That Matters: Small Language Models Are Also Few-Shot Learners*. See [[1.0.0 P-55]]. The SLM-prompting anchor for the 4B-30B regime MAgHARCM targets.
- Gao, L. et al. (2022). *PAL: Program-aided Language Models*. See [[1.0.0 P-65]]. Program-aided reasoning; a tool-use variant.
- Bubeck, S. et al. (2023). *Sparks of Artificial General Intelligence: Early Experiments with Large Language Models*. arXiv:2303.12712. The GPT-4 capability survey that Gorilla explicitly rebuts on the API-calling axis.

### MAgHARCM lineage cross-refs
- [[1.0.0 P-21]] — Qwen2.5-Coder 2024 (current SLM fleet anchor).
- [[1.0.0 P-22]] — StarCoder2 2024 (alternative SLM fleet anchor).
- [[1.0.0 P-54]] — Phi-3 Tech Report 2024 (alternative SLM fleet anchor).
- [[1.0.0 P-55]] — Schick & Schütze 2021 (SLM prompting anchor; foundational for the 4B-30B regime Gorilla validates).
- [[1.0.0 P-61]] — RedCode 2024 (the code-hallucination analogue; verifier-of-verifier architecture justification).
- [[1.0.0 P-62]] — Khot Decomposed Prompting 2022 (decomposition pattern that ReAct interleaves).
- [[1.0.0 P-65]] — Gao PAL 2022 (program-aided reasoning variant).

## Backlinks

[[1.0.0 PRIM-7]] (Multi-Agent Verdict Validation), [[1.0.0 PRIM-25]] (Role-Flip De-Hallucination Gate), [[1.0.0 P-55]], [[1.0.0 P-61]], [[1.0.0 P-62]], [[1.0.0 P-65]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-81 is the **tool-use SLM anchor** for MAgHARCM: it justifies the 4B-30B SLM regime and provides the retriever-aware template for the optional-checks agent.
