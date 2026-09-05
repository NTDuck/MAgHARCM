---
title: "P-62 — Khot et al. 2022 — Decomposed Prompting: A Modular Approach for Solving Complex Tasks"
backlink: "[[2.0.0 P-62]]"
tags: [paper, prompting, decomposition, modular, multi-step, slm, chain-of-thought, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-24]], [[1.0.0 PRIM-31]], [[2.0.0 MAgHARCM]]]
---

# [[2.0.0 P-62 — Khot et al. — Decomposed Prompting]]

## Citation

Khot, T., Richardson, P., Khashabi, D., Sabharwal, A. (2022). *Decomposed Prompting: A Modular Approach for Solving Complex Tasks*. ICLR 2023. arXiv:2210.02406 (October 2022). AI2 & University of Washington. URL: https://arxiv.org/abs/2210.02406 ; https://openreview.net/forum?id=_nGgzakjvf.

## Summary

Decomposed Prompting (DP) attacks a single empirical observation: **one monolithic prompt asking an LM to "solve the whole task" scales worse than a small library of role-specialised prompts glued together by a lightweight planner**. The recipe:

1. **A small, fixed library of sub-task modules**, each implemented as its own short prompt — e.g. `string-splitter`, `object-finder`, `question-answerer`, `verifier`. Each module's prompt is *short enough to fit comfortably in the SLM attention window* and exposes a typed I/O signature (input string + label, output string).
2. **A few-shot "decomposer" prompt** that consumes the original question and emits a *directed acyclic graph* (DAG) of module invocations, where each edge is an I/O type match. The DAG is textual (`QA -> STRING_SPLIT -> OBJECT_FINDER -> QA`) and parsed by a thin orchestrator, not by the LM.
3. **Per-module black-box calls** — the orchestrator invokes each module's prompt exactly once and threads outputs through the DAG. There is no recursive LM call inside a module; each module's prompt is a single self-contained prompt.
4. **Empirical sweep across 23 task domains** (date understanding, sports understanding, reasoning about objects, penguins-in-a-table, logical-args, etc.) at the Codex / Codex-Davinci scale, showing DP recovers most of the gains from in-prompt chain-of-thought (Wei et al. 2022) *while making the intermediate steps inspectable and replayable*.

**SLM-aware relevance for MAgHARCM**: the paper is the *theoretical twin* of MAgHARCM's prompt architecture. Where Wei et al. 2022 (CoT) put the reasoning trace *inside* one prompt, DP puts the reasoning trace *outside* the LM, in a deterministic orchestrator that hands each step its own short prompt. That is exactly the decomposition `internal/agents/prompts.go` embodies: one template per agent role, a thin orchestrator (the Eino graph in `internal/graph/graph.go`) connecting them, and shared typed artifacts (`internal/compiletime/state.go`) threaded as the I/O signature.

## Method

DP's two-stage protocol:

1. **Decomposer prompt.** A few-shot prompt that, given the natural-language task, emits a textual DAG specifying which module calls to make in which order. Modules are typed by their I/O: `QA` (free text → free text), `STRING_SPLIT` (string → [string]), `OBJECT_FINDER` ([string] → [object]), `BOOL` (any → bool), etc.
2. **Module calls.** For each module invocation, the orchestrator substitutes the I/O from upstream modules into the module's prompt and queries the LM. The module's output is fed downstream.

The empirical contribution is the catalogue of 23 tasks and the proof that DP *replaces* in-prompt CoT for many tasks — the reasoning trace lives in the DAG, not in the prompt, which is the SLM-friendly choice.

## Findings Relevant to MAgHARCM

- **Short prompts scale better than long prompts on SLMs.** The 4B-30B SLM fleet MAgHARCM is tuned for (cf. [[1.0.0 P-58]] Qwen2.5-Coder) loses accuracy as prompt length grows beyond ~2-4K tokens (cf. [[1.0.0 P-53]] Liu et al. *Lost in the Middle*). DP's short-module discipline fits the SLM constraint; CoT-style long-prompt reasoning does not.
- **Typed I/O signatures map to typed state.** MAgHARCM's `compiletime.State` is a typed container; the prompts in `internal/agents/prompts.go` declare which typed fields they consume (e.g., `SourceProjectResearch`, `ImplementationPlan`). DP's typed module signature is the theoretical anchor for this design.
- **Orchestrator outside the LM is inspectable.** DP's orchestrator is a thin Python (or Go) loop that emits intermediate artifacts. MAgHARCM's Eino graph (`internal/graph/graph.go`) is exactly this: each node's typed input/output is inspectable, replayable, and serializable to disk (cf. [[1.0.0 PRIM-28]] Conversable State Checkpoints).
- **DAG vs. chain.** DP's DAG generalizes CoT's chain. MAgHARCM's Eino graph supports both DAG (multiple parallel paths within an iteration) and cycles (repair loop), which is a strict superset of DP's protocol.

## How MAgHARCM Uses It

The architecture of `internal/agents/prompts.go` (six role-specialised templates — Analyzer, Planning, Translator-Translate, Translator-Chunked, Translator-Repair, Validator-Coverage) plus the orchestrator in `internal/graph/graph.go` is the textbook Decomposed Prompting layout. The shared preamble `compiletime.SLMPromptContractPreamble` prepended to every module is DP's "library of short typed prompts" pattern. The `compiletime.DocumentWrapper[T]` typed intermediate artifact is DP's typed I/O signature. The TranslatorChunkedPromptTemplate's "Previously Emitted Modules" state block is DP's per-module re-invocation discipline applied within one role.

## References

### Hop-1
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. NeurIPS 2022. arXiv:2201.11903.
- Kojima, T. et al. (2022). *Large Language Models are Zero-Shot Reasoners*. NeurIPS 2022. arXiv:2205.11916.
- Brown, T. B. et al. (2020). *Language Models are Few-Shot Learners*. arXiv:2005.14165. (GPT-3)
- Chen, M. et al. (2021). *Evaluating Large Language Models Trained on Code*. arXiv:2107.03374. (Codex / HumanEval)
- Schick, T. & Schütze, H. (2021). *It's Not Just Size That Matters*. NAACL 2021. arXiv:2009.07118. See [[1.0.0 P-55]].

### Hop-2
- Suzgun, M. et al. (2022). *Challenging BIG-Bench Tasks and Whether Chain-of-Thought Can Solve Them*. arXiv:2210.09261.
- Wang, X. et al. (2023). *Self-Consistency*. ICLR 2023. arXiv:2203.11171. See [[1.0.0 P-52]].
- Wei, J. et al. (2022). *Emergent Abilities of Large Language Models*. TMLR 2022. arXiv:2206.07682.

## Backlinks

[[1.0.0 PRIM-22]], [[1.0.0 PRIM-24]], [[1.0.0 PRIM-31]], [[1.0.0 PRIM-28]], [[1.0.0 P-53]], [[1.0.0 P-55]], [[1.0.0 P-58]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-62 is the **theoretical anchor** for MAgHARCM's per-agent prompt-template architecture: short typed prompts + deterministic orchestrator + typed intermediate artifacts. This is the SLM-friendly alternative to long-prompt CoT.
