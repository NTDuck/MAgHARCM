---
title: "P-100 — Khot, Trivedi, Finlayson, Fu, Richardson, Clark & Sabharwal 2022 — Decomposed Prompting: A Modular Approach for Solving Complex Tasks"
backlink: "[[1.0.0 P-100]]"
aliases:
  - "1.0.0 P-100"
  - "P-100"
  - "P-100-Khot-Decomposed-Prompting-2022"
  - "P-100-Khot-Decomposed-Prompting-2022"
  - "Khot-Decomposed-Prompting-2022"
tags: [paper, decomposed-prompting, modular, slm, multi-agent, slm-multi-agent, dag, typed-io, sop, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-24]], hop-2]
---

# [[1.0.0 P-100 — Khot et al. — Decomposed Prompting (SLM-era re-anchor)]]

## Citation

Khot, T., Trivedi, H., Finlayson, M., Fu, Y., Richardson, K., Clark, P., & Sabharwal, A. (2022). *Decomposed Prompting: A Modular Approach for Solving Complex Tasks*. ICLR 2023. arXiv:2210.02406 (5 Oct 2022 v1). Allen Institute for AI (AI2) & University of Washington. URL: https://arxiv.org/abs/2210.02406 ; https://openreview.net/forum?id=_nGgzakjvf.

This is the **SLM-era re-anchor** of the same paper as `[[1.0.0 P-62]]`; P-62 framed the paper in Codex/PaLM-class LLM terms, P-100 reframes it in the **SLM multi-agent** terms MAgHARCM is built on (see `[[1.0.0 P-100]]` ↔ `[[1.0.0 P-62]]` cross-link in Backlinks).

## Summary

Decomposed Prompting (DP) observes a single empirical fact: **a monolithic prompt that asks one LM to "solve the whole task" scales worse — in cost, accuracy, and inspectability — than a small library of role-specialised short prompts connected by a deterministic orchestrator**. The recipe is four moves:

1. **A small, fixed library of sub-task modules**, each implemented as its own short prompt with a *typed I/O signature* (input typed slot → output typed slot). Modules are short enough to fit comfortably inside an SLM attention window: `string-splitter(string) → [string]`, `object-finder([string]) → [object]`, `qa(question) → answer`, `verifier(claim, evidence) → bool`, etc.
2. **A few-shot "decomposer" prompt** that consumes the natural-language task and emits a *directed acyclic graph (DAG) of module invocations*, where each edge is an I/O type match. The DAG is textual (`QA → STRING_SPLIT → OBJECT_FINDER → QA`) and parsed by a thin orchestrator, *not by the LM*.
3. **Per-module black-box calls** — the orchestrator invokes each module's prompt exactly once and threads outputs through the DAG. No recursive LM call inside a module; each module's prompt is a single self-contained prompt.
4. **An empirical sweep across 23 task domains** (date understanding, sports understanding, reasoning about objects, penguins-in-a-table, logical-args, etc.) on Codex / Codex-Davinci, showing DP recovers most of the accuracy gains of in-prompt chain-of-thought `[[1.0.0 P-90]]` Wei et al. 2022 *while making the intermediate steps inspectable, replayable, and replaceable*.

**SLM multi-agent relevance for MAgHARCM**: the central claim of P-100 is the **equivalence** *SLM multi-agent = decomposed prompting*. Each agent in a multi-agent SLM system is exactly one DP module: a short typed prompt, a typed I/O signature, an external orchestrator routing typed artifacts between agents. P-100 is the *theoretical proof* that the MAgHARCM layout — six role-specialised agent prompts (`internal/agents/prompts.go`) connected by an Eino orchestrator (`internal/graph/graph.go`) over a typed state container (`internal/compiletime/state.go`) — is not an ad-hoc engineering choice but the *canonical* way to do complex-task prompting under SLM budget constraints.

## Method

DP's two-stage protocol:

1. **Decomposer prompt.** A few-shot prompt that, given the natural-language task, emits a textual DAG specifying which module calls to make in which order. Modules are typed by their I/O (`QA`, `STRING_SPLIT`, `OBJECT_FINDER`, `BOOL`, `INTEGER`, etc.); the decomposer's job is to emit a *type-respecting* call plan.
2. **Module calls.** For each module invocation, the orchestrator substitutes the I/O from upstream modules into the module's prompt and queries the LM. The module's output is fed downstream; intermediate artifacts are first-class, inspectable objects on the orchestrator's data path.

The empirical contribution is the catalogue of 23 tasks and the proof that DP **replaces** in-prompt CoT for many tasks: the reasoning trace lives in the DAG, not in the prompt, which is the SLM-friendly choice (short prompts, no chained few-shot exemplars, no lost-in-the-middle risk).

**Specialist-module design rules** distilled from the paper and adopted as MAgHARCM's per-agent prompt-engineering convention:

- **One role per module.** Each module does *one* sub-task. A `qa` module does not also extract entities; a `string-split` module does not also score sentiment. The composability of DP falls apart the moment modules become multi-purpose.
- **Typed I/O on every module.** Inputs and outputs are typed slots (`string`, `[string]`, `bool`, `object`), not free text. Type mismatches are caught at DAG-build time by the orchestrator, not at LM-decode time by a hallucination.
- **Short prompts.** Each module's prompt fits in <1K tokens of system prompt + typed I/O. SLM context budgets (cf. `[[1.0.0 P-80]]` StreamingLLM) reward shortness; the prompt's value is the *typed signature*, not the prose.
- **No in-module reasoning chains.** A module may emit a *short* typed answer; if reasoning is required, it is the *orchestrator's* job to route through a `reasoner` module that itself has a typed I/O. Reasoning never lives inside a `qa` module's prose.
- **Deterministic orchestrator.** The DAG parser and routing are code, not LM. The LM is only ever called once per module invocation. No recursive LM calls, no LM-driven tool selection inside a module.

## Findings Relevant to MAgHARCM

- **SLM multi-agent = decomposed prompting (the central claim).** A multi-agent SLM system is *literally* a DP deployment: one module per agent, typed I/O matching between agents, an external orchestrator routing typed artifacts. MAgHARCM's six-agent layout (Analyzer, Planning, Translator-Translate, Translator-Chunked, Translator-Repair, Validator-Coverage) is the DP pattern instantiated for SLM-era software-archaeology migration tasks. P-100 is the theoretical anchor for *why this layout is canonical*.
- **Specialist modules are the SLM-friendly answer to long-prompt reasoning.** The 4B-30B SLM fleet MAgHARCM is tuned for (Qwen2.5-Coder `[[1.0.0 P-58]]`, StarCoder2 `[[1.0.0 P-22]]`, Phi-3-mini `[[1.0.0 P-54]]`, Code Llama `[[1.0.0 P-95]]`) loses accuracy as prompt length grows past ~2-4K tokens (cf. `[[1.0.0 P-53]]` Lost in the Middle). DP's specialist-module discipline is the SLM-era answer: do not ask one prompt to do comprehension + planning + translation + verification; ask *four short typed prompts*, each on its own agent.
- **Typed I/O signatures map directly to typed state.** MAgHARCM's `compiletime.DocumentWrapper[T]` and `compiletime.State` (`internal/compiletime/state.go`) are the typed I/O slots DP prescribes. The compile-time schema check (cf. `[[1.0.0 PRIM-24]]` SOP-Anchored Role-Artifact Schema) is the *static* version of DP's runtime type-matching: MAgHARCM catches type errors at compile time, DP catches them at DAG-build time.
- **Orchestrator outside the LM is inspectable and replayable.** DP's orchestrator is a thin loop that emits intermediate artifacts. MAgHARCM's Eino graph (`internal/graph/graph.go`) emits per-node typed state to disk (cf. `[[1.0.0 PRIM-28]]` Conversable State Checkpoints). Every intermediate artifact in MAgHARCM — the comprehension summary, the plan, each translated chunk — is inspectable, replayable, and serialisable. This is the DP property in production.
- **DAG generalises chain; cycles generalise DAG.** DP's static DAG is the planning substrate; MAgHARCM adds the repair cycle (Translator-Repair → Validator-Coverage → Translator-Repair) on top of the static DAG. Cycles are not in DP but are a natural extension; P-100 justifies the static DAG, and MAgHARCM extends it to a DAG+cycles graph.
- **Module library is small and stable.** DP's 23 modules do not grow unboundedly; the library is *closed* once it covers the task taxonomy. MAgHARCM's six-agent library is similarly closed — adding a seventh agent is a deliberate API change, not a per-task prompt edit. This is the *SLM budget discipline*: a fixed library can be pre-compiled, pre-validated, and pre-tuned.
- **Decomposed Prompting *replaces* in-prompt CoT for many tasks.** Where `[[1.0.0 P-90]]` puts the reasoning trace inside the prompt (an SLM-unfriendly choice on sub-30B models), DP puts the trace in the orchestrator (SLM-friendly). MAgHARCM inherits this: the `TranslatorChunkedPromptTemplate`'s "Previously Emitted Modules" state block is DP's per-module re-invocation discipline, applied inside one agent role.

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-22]]` Four Phases of Comprehension** (`internal/agents/comprehension.go`): the four comprehension phases — *explain*, *search*, *plan*, *summarise* — are P-100's specialist-module pattern applied to source-code understanding. Each phase is a *separate short typed prompt* (Explain: source → natural-language explanation; Search: source + explanation → list-of-locations; Plan: explanation + locations → step-list; Summarise: step-list + locations → comprehension-summary). The typed intermediate artifacts (explanation, locations, plan, summary) are DP's typed I/O slots, stored in `compiletime.State` and threaded between phases by the orchestrator. P-100 is the *theoretical justification* for splitting comprehension into four typed phases instead of one big "understand this code" prompt.
- **`[[1.0.0 PRIM-24]]` SOP-Anchored Role-Artifact Schema** (`internal/compiletime/compiletime.go`, `internal/agents/state.go`): the SOP-style role-artifact schema (each agent role declares its typed input slots, typed output slots, and the contract between them) is DP's typed I/O discipline enforced *at compile time*. Where DP catches type mismatches at DAG-build time, MAgHARCM catches them at Go compile time via the generic `DocumentWrapper[T]` wrapper. P-100 is the *theoretical justification* for the SOP schema: typed I/O on every module is what makes decomposition composable.
- **Six-agent layout (`internal/agents/prompts.go`)**: Analyzer, Planning, Translator-Translate, Translator-Chunked, Translator-Repair, Validator-Coverage — the MAgHARCM agent library. This is the *concrete instance* of P-100's "small fixed library of specialist modules" pattern, applied to SLM-era code migration. Cite P-100 as the theoretical anchor for *why the library is small, fixed, and specialist* (rather than a single generalist prompt, or a sprawling per-task prompt zoo).
- **`internal/graph/graph.go` orchestrator**: the Eino graph that connects the six agents is P-100's "thin orchestrator outside the LM" in production code. It parses the typed call plan, threads typed artifacts between agents, emits per-node state to disk, and never calls an LM recursively. P-100 is the *theoretical justification* for keeping the orchestrator LM-free.
- **`compiletime.SLMPromptContractPreamble`** (`internal/agents/prompts.go`): the shared preamble prepended to every agent's prompt is the DP "library of short typed prompts" convention realised as a Go constant. Each prompt template is *short* (specialist, single-role) and *typed* (declares its typed I/O slots).
- **Repair cycle (Translator-Repair ↔ Validator-Coverage)**: MAgHARCM's repair loop is a *natural extension* of P-100's static DAG to cycles. DP does not include cycles; MAgHARCM adds them for the validator→repair feedback. Cite P-100 for the static DAG justification, then note the extension.
- **Future work (not implemented)**: a per-migration-task DP "library selector" — pick a 4-7 module subset of MAgHARCM's six-agent library for the task at hand (e.g. a chunked-translation task uses Translator-Chunked + Validator-Coverage + Translator-Repair; a comprehension-only task uses Analyzer + Planning). The selection itself is a small DP decomposer prompt over the closed library.

## References

### Hop-1 (papers that build directly on decomposed prompting or are direct descendants)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. NeurIPS 2022. arXiv:2201.11903. See `[[1.0.0 P-90]]` — the in-prompt CoT baseline DP replaces. DP and CoT are the *two routes to multi-step reasoning*: CoT puts the trace in the prompt (SLM-unfriendly), DP puts it in the orchestrator (SLM-friendly).
- Khot, T. et al. (2022). *Decomposed Prompting: A Modular Approach for Solving Complex Tasks*. arXiv:2210.02406. See `[[1.0.0 P-62]]` — the original MAgHARCM anchor for this paper. P-62 frames DP in Codex/PaLM-class terms; P-100 re-frames it in SLM multi-agent terms. P-62 and P-100 are *two anchors on the same paper*; cite both when the SLM multi-agent framing matters.

### Hop-2 (foundational anchors referenced)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. NeurIPS 2022. arXiv:2201.11903. See `[[1.0.0 P-90]]` — the in-prompt CoT baseline DP generalises. DP's static DAG is the *external* counterpart to CoT's in-prompt chain; both are routes to multi-step reasoning, but only DP is SLM-friendly on sub-30B models.
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners* (GPT-3). arXiv:2005.14165. The in-context-learning foundation both CoT and DP build on; the decomposer prompt is itself a few-shot prompt.
- Chen, M. et al. (2021). *Evaluating Large Language Models Trained on Code*. arXiv:2107.03374. (Codex / HumanEval). The model class and benchmark family DP's 23-task sweep targets.
- Suzgun, M. et al. (2022). *Challenging BIG-Bench Tasks and Whether Chain-of-Thought Can Solve Them*. arXiv:2210.09261. The 23-task evaluation harness DP shares; DP and the Suzgun et al. CoT study are companion papers for the BIG-Bench/Held-Out suite.
- Wei, J. et al. (2022). *Emergent Abilities of Large Language Models*. TMLR 2022. arXiv:2206.07682. The scale-dependence framework that explains *why* DP is the SLM-friendly answer (specialist modules avoid the emergent-CoT scale threshold).

## Backlinks

`[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-24]]`, `[[1.0.0 PRIM-28]]`, `[[1.0.0 P-53]]`, `[[1.0.0 P-54]]`, `[[1.0.0 P-55]]`, `[[1.0.0 P-58]]`, `[[1.0.0 P-62]]`, `[[1.0.0 P-80]]`, `[[1.0.0 P-90]]`, `[[1.0.0 P-95]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-100 is the **SLM multi-agent re-anchor** of Decomposed Prompting for MAgHARCM. Where P-62 frames the paper as the *prompt-architecture* anchor (short typed prompts + deterministic orchestrator + typed intermediate artifacts), P-100 frames the paper as the *theoretical proof* that **SLM multi-agent = decomposed prompting** — i.e. that the MAgHARCM six-agent layout is not an ad-hoc engineering choice but the canonical SLM-era instantiation of the DP pattern. Cite P-100 whenever a primitive's design is justified in terms of specialist modules, typed I/O, or the equivalence between multi-agent systems and decomposed prompting; cite P-62 whenever the prompt-architecture framing (Codex/PaLM-era) suffices.
