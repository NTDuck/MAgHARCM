---
title: "P-112 — Yang et al. 2024 — SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering"
backlink: "[[1.0.0 P-112]]"
aliases:
  - "1.0.0 P-112"
  - "P-112"
  - "P-112-Yang-SWE-Agent-2024"
  - "P-112-Yang-SWE-Agent-2024"
  - "P-112-Yang-SWE-Agent-2024"
  - "P-112-Yang-SWE-Agent-2024"
  - "Yang-SWE-agent-2024"
tags: [paper, agent, code, swebench, tool-use, [[1.0.0 PRIM-5]], [[1.0.0 PRIM-6]], [[1.0.0 PRIM-25]], [[1.0.0 PRIM-29]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-112 — Yang et al. 2024 — SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering]]

- **Authors**: John Yang, Carlos E. Jimenez, Alexander Wettig, Hila Weintraub, Ofir Press, Karthik Narasimhan (Princeton Language and Intelligence / Princeton NLP).
- **Venue / Year**: NeurIPS 2024 (Main track, Large Language Models for Code workshop also presented earlier); arXiv:2405.15793 v4.
- **URL**: https://arxiv.org/abs/2405.15793 ; code at https://github.com/princeton-nlp/SWE-agent
- **Anchors**: PRIM-5 (Test Suite Co-Translation), PRIM-6 (Multi-Stage Build/Test Repair), PRIM-25 (Role-Flip De-Hallucination), PRIM-29 (Recruitment-Adaptive Planning); the canonical tool-calling agent scaffold that defines the SWE-bench evaluation protocol OpenHands, AutoCodeRover, and Aider all build on.

## 1. Core Contribution

SWE-agent is the **first autonomous agent scaffold purpose-built for repository-level software-engineering tasks**, and the canonical reference point for what an "agent" looks like on [[1.0.0 P-109]] (SWE-bench Verified). The paper makes two complementary contributions:

1. **Agent-Computer Interfaces (ACIs) as a first-class design object.** Where prior agent work treated the terminal as a black box (the LM writes `bash` commands and hopes), SWE-agent elevates the **interface layer between the agent and the computer** to a research artefact. The ACI is a curated, task-tailored set of commands — `open`, `scroll_up`, `scroll_down`, `find_file`, `search_file`, `search_dir`, `edit`, `insert`, `submit` — that re-shape what the agent can *perceive* of a repository and what it can *do* with that perception. Concretely:
   - **File viewing as a first-class operation.** Instead of cat-ing a 10k-line file (which blows the context window and hides structure), the ACI exposes `open`/`scroll_up`/`scroll_down` that load a windowed view plus line numbers. The agent reads a file the way a human does: top-down, with the cursor advancing.
   - **Search as a structured operation.** `search_file` and `search_dir` return regex-matched lines with surrounding context, not raw `grep` blobs. The ACI normalises search output so the model can parse it without bespoke parsing logic per repository.
   - **Edit as a structured operation.** `edit` takes a `start_line:end_line` block plus a replacement; `insert` inserts above a target line. The model never invents a `sed` script — it submits a structured edit primitive that the host validates and applies.
   - **Submit and reset.** The `submit` command finalises the candidate patch; a reset hook rolls the environment back to the pre-edit state for retries.
2. **Headline result: 12.5% on SWE-bench Verified with GPT-4.** SWE-agent was the **first** autonomous agent to clear the 10% bar on Verified and set the early baseline that every subsequent paper compares against. It was re-evaluated on [[1.0.0 P-109]] (SWE-bench Verified, OpenAI 2024) and remains the standard scaffold for SLM-grade agents even when the underlying model is swapped to Qwen2.5-Coder 7B, Llama-3.1 8B, or DeepSeek-Coder-V2-Lite.

The paper's central empirical claim is that **the interface matters more than the model for many repository-level tasks**. When the same GPT-4 backbone is given a vanilla `bash`-shell interface, its SWE-bench Verified score drops sharply (single-digit territory in their ablation); with the ACI, it lands at 12.5%. The lift is not from a smarter model — it is from a *better-shaped perception*.

Beyond the headline number, the paper establishes the **SWE-agent evaluation protocol** that [[1.0.0 P-110]]-adjacent scaffolds (OpenHands/CodeAct, AutoCodeRover, Aider) inherit:
- **Single Docker environment per instance**, mounted from the SWE-bench task definition.
- **Read-only initial state** with a `submit`-then-exit run loop; no in-place retries on the candidate patch (the host resets the environment for any retry attempt).
- **Per-instance test suite** as the oracle; the agent's job is to make the failing tests in the suite pass without breaking the passing ones.
- **Token budget** measured in agent steps, not wall-clock; a 75-step default is the SWE-agent norm, with up to 250 for hard instances.

**Why this matters for the SLM regime**: SWE-agent's ACI pattern is exactly the **structural de-hallucination mechanism** that a 4B–7B model needs. A vanilla `bash` interface lets the SLM invent shell pipelines that fail silently or return garbled output; a curated ACI returns **parseable, bounded, structured responses** that an SLM can read reliably. This is the structural analogue of the cloze reformulation in [[1.0.0 P-55]] — the ACI converts free-form tool use into slot-filling over a typed interface, which is exactly what an SLM does best.

## 2. Application in MAgHARCM

SWE-agent is one of the four pillars of MAgHARCM's agent-scaffold design (alongside [[1.0.0 P-110]]'s graph retrieval, [[1.0.0 P-111]]'s SWE-bench task substrate, and the [[1.0.0 P-109]] Verified oracle). It anchors four primitives.

### 2.1 PRIM-25 (Role-Flip De-Hallucination)

The ACI pattern is **the structural instantiation of PRIM-25 at the tool-call layer**. A 4B–7B SLM asked to produce free-form `bash` will hallucinate pipelines (chained `awk | sed | grep | cut` constructs that the model has no way to verify mid-stream). SWE-agent's ACI constrains every observation to a **structured tool schema**: `open <path>:<line>` returns a windowed file view with line numbers; `search_file <regex> <path>` returns a list of `(file, line, match)` triples; `edit <start>:<end>` requires a literal block replacement. The schema forces the SLM into **fill-in-the-slots** generation: instead of composing shell logic, the model selects a tool and fills in the schema fields. The hallucination surface shrinks from "any string of bash" to "typed arguments to a known tool". This is the tool-call analogue of [[1.0.0 P-55]]'s cloze reformulation — same idea (structured output slots), different layer (tools instead of natural language). SWE-agent's empirical lift (from sub-5% to 12.5% on Verified by changing only the interface) is the **headline evidence** that PRIM-25's structural pattern works at the agent layer.

### 2.2 PRIM-5 (Test Suite Co-Translation & Synthesis)

SWE-agent's `submit` primitive is paired with a **per-instance test runner hook** that reports the failing tests back into the agent's context window. This is the canonical SWE-bench-style oracle feedback loop, and it is the exact pattern that MAgHARCM's Translator's `Validator.generateAdditionalTests` pathway needs to consume: a **bounded, parseable report** of which tests fail and how, fed back into the next repair iteration. SWE-agent's `search_dir` / `search_file` primitives also feed PRIM-5's test-synthesis prompt: the test-discovery prompt `find all tests related to <symbol>` is implemented as `search_dir <symbol_pattern> tests/` followed by `open` on the candidate test files. The agent then extracts the existing test idiom (setup, fixtures, assertion style) before synthesising new tests, and the structured search output is what makes the extraction reliable for an SLM.

### 2.3 PRIM-6 (Multi-Stage Build/Test Feedback Repair)

SWE-agent's loop is the **canonical validator cascade** that PRIM-6 generalises:
1. **Stage 1 (build/import)**: the agent's `submit` triggers a host-side `pytest --collect-only`. If collection fails (import error, syntax error, missing dependency), the error is reported back and the agent iterates without ever seeing test outcomes.
2. **Stage 2 (first failing test)**: once collection succeeds, the first failing assertion in the run is reported with its traceback. The agent iterates on **that one assertion**, not the full test outcome — a critical insight that PRIM-6 inherits. SWE-agent showed that reporting only the first failing test, not the full suite, gives the agent a **focal target** for repair; this is what enables the SLM's effective reasoning to converge.
3. **Stage 3 (regression check)**: on a candidate patch that fixes the failing test, the host re-runs the **entire** suite and reports any *previously passing* tests that now fail. The agent iterates until either the full suite passes or the step budget is exhausted.

The **kind-of-failure signal** (collection error vs assertion error vs regression error) is exactly the discriminator that `Registry.TryInOrder` in PRIM-21 consumes to pick the next repair strategy. SWE-agent's evaluation log on SWE-bench Verified provides a calibration dataset for the per-strategy hit rate.

### 2.4 PRIM-29 (Recruitment-Adaptive Planning)

The **dynamic tool selection** pattern — where the agent decides at each step which ACI primitive to invoke based on the current observation — is the **template for PRIM-29's recruitment loop**. PRIM-29's recruitment-adaptive planning reframes "which tool should I call?" as "which agent role should I recruit next?", and the structural pattern is identical:
- **Observation-driven selection**: the next tool/agent is chosen by inspecting the current context, not a pre-baked plan.
- **Bounded step budget**: SWE-agent's 75–250 step cap is the analogue of PRIM-29's max-recruit-rounds parameter; both protect against runaway planning.
- **Submit/retry boundary**: SWE-agent's `submit` is PRIM-29's "promote-to-final" boundary; on either side of the boundary, the recovery model is different (in-place repair vs full re-rollout).
- **Per-step verbosity budget**: SWE-agent truncates long file views to a window; PRIM-29 truncates long agent transcripts to a summary. Both are the **bounded-perception** discipline that lets an SLM stay coherent over a long planning horizon.

PRIM-29's recruitment registry (`registry/recruitment.go` in the MAgHARCM codebase) borrows the ACI-style interface schema: each recruit role has a name, an input contract, an output contract, and a `submit`-equivalent final-action. SWE-agent is the proof point that this kind of typed registry, rather than a free-form `run bash` interface, is what makes agent-level planning tractable for SLMs.

### 2.5 SLM-era headline framing

Together with [[1.0.0 P-109]] (SWE-bench Verified) and [[1.0.0 P-102]] (SmallCode 4B at 87% HumanEval), SWE-agent closes the SLM-era loop: a 4B–7B SLM scaffolded with a SWE-agent-style ACI can clear **12.5% on SWE-bench Verified**, which is enough to be measured as a real signal and to leave meaningful headroom for the structural lift that MAgHARCM's 8-agent pipeline is designed to add. The gap between single-function synthesis (87% on HumanEval) and repository-level issue resolution (12.5% on Verified, with a frontier ceiling near 49%) is exactly the gap MAgHARCM is designed to close.

## 3. Hop-1 References (papers cited by SWE-agent)

- **Jimenez et al. 2024** — [[1.0.0 P-111]] — *SWE-bench: Can Language Models Resolve Real-World GitHub Issues?* (ICLR 2024). The 2,294-instance benchmark SWE-agent is evaluated on; the SWE-agent paper is one of the three canonical agents evaluated in the SWE-bench technical report.
- **OpenAI 2024** — [[1.0.0 P-109]] — *SWE-bench Verified*. The 500-instance human-validated subset; SWE-agent's 12.5% is reported on this subset and is the canonical SWE-agent number cited by every follow-on paper.
- **Yao et al. 2022** — *ReAct: Synergizing Reasoning and Acting in Language Models* (ICLR 2023). The reasoning+acting framework SWE-agent inherits and adapts to repository-level work; the "Thought / Action / Observation" trace is the SWE-agent prompt template's backbone.
- **Schick et al. 2023** — *Toolformer: Language Models Can Teach Themselves to Use Tools* (NeurIPS 2023). Foundational tool-use learning paper; SWE-agent's ACI is the structured, hand-engineered alternative to Toolformer's learned tool use.
- **Shinn et al. 2023** — *Reflexion: Language Agents with Verbal Reinforcement Learning* (NeurIPS 2023). The verbal-self-reflection loop SWE-agent partially inherits; SWE-agent's `submit`-and-reset boundary is the analogue of Reflexion's episode boundary.
- **Wang et al. 2024** — *OpenHands (CodeAct): An Open Platform for AI Software Developers as Generalist Agents*. The multi-turn CodeAct agent; SWE-agent's ACI is the single-turn, single-edit predecessor, and OpenHands' CodeAct is the multi-turn successor that subsumes SWE-agent's interface philosophy.
- **Zhang et al. 2024** — [[1.0.0 P-113]] — *AutoCodeRover: Autonomous Program Improvement*. The autonomous software-engineering agent that combines retrieval and program synthesis on SWE-bench; co-evaluated with SWE-agent and the canonical open-source baseline for Verified.
- **Aider (2024)** — [[1.0.0 P-115]] (candidate slot) — *Aider polyglot / chat-with-your-codebase*. The per-file-diff code-edit agent; cross-evaluates with SWE-agent on the Aider benchmark, and shares the "structured edit primitive" interface philosophy (the Aider `edit` block is structurally identical to SWE-agent's `edit` ACI).
- **Peng et al. 2023** — *RepoCoder: Repository-Level Code Completion Through Iterative Retrieval and Generation*. Iterative-retrieval-augmented code completion; one of the agent scaffolds evaluated in the SWE-bench technical report alongside SWE-agent.
- **Chen et al. 2021** — *Evaluating Large Language Models Trained on Code* (HumanEval / Codex). Foundational single-function synthesis benchmark; SWE-agent is the repository-level descendant, and HumanEval remains the unit-test reference for SWE-agent's per-function behaviour.

## 4. Hop-2 References (papers-cited-by-hop-1)

- **Vaswani et al. 2017** — *Attention Is All You Need* (NeurIPS 2017). Transformer backbone for every model SWE-agent evaluates (GPT-4, Claude, open-source SLMs).
- **OpenAI 2023** — *GPT-4 Technical Report*. The backbone SWE-agent's 12.5% number is reported on; the default model in the SWE-agent Docker harness.
- **Wei et al. 2022** — *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. CoT is the substrate ReAct extends; SWE-agent inherits CoT inside each Thought/Action/Observation step.
- **Zhou et al. 2023** — [[1.0.0 P-101]] — *Least-to-Most Prompting Enables Complex Reasoning in Large Language Models*. Problem-decomposition pattern; SWE-agent's natural per-issue decomposition (read issue → locate relevant code → edit → verify) is the operational form of least-to-most.
- **Gao et al. 2022** — *PAL: Program-aided Language Models*. The pattern of offloading computation to a deterministic executor; SWE-agent's `submit`-and-host-validates loop is the PAL pattern applied at the repository-edit boundary.
- **Touvron et al. 2023** — *LLaMA 2: Open Foundation and Fine-Tuned Chat Models*. The open-model backbone for many of SWE-agent's ablation runs; Qwen2.5-Coder and Llama-3.1 inherit from this lineage.
- **Chen et al. 2021** — *Evaluating Large Language Models Trained on Code* (code-davinci-002 / Codex). The Codex backbone for the original SWE-agent-style prototypes; SWE-agent's first prototype was Codex, not GPT-4.
- **Hendrycks et al. 2021** — *Measuring Massive Multitask Language Understanding* (MMLU). General-capability baseline reported alongside SWE-agent's SWE-bench numbers; SWE-agent's authors report MMLU as a sanity check on the underlying models.
- **Zhou et al. 2023** — [[1.0.0 P-94]] — *LIMA: Less Is More for Alignment*. Quality-over-quantity alignment result that SWE-agent's prompt-design discipline inherits; the SWE-agent ACI prompt is hand-curated in the LIMA style (few, high-quality exemplars, no RLHF).

## 5. Backlinks

- **PRIM-5** (Test Suite Co-Translation & Synthesis): SWE-agent's `submit`-driven test report is the oracle feedback loop that PRIM-5's `Validator.generateAdditionalTests` consumes; the structured search primitives feed PRIM-5's test-discovery prompt.
- **PRIM-6** (Multi-Stage Build/Test Feedback Repair): SWE-agent's three-stage loop (collection → first failing test → regression) is the canonical validator cascade that PRIM-6 generalises; the kind-of-failure signal feeds `Registry.TryInOrder` in PRIM-21.
- **PRIM-25** (Role-Flip De-Hallucination): SWE-agent's ACI is the **structural instantiation of PRIM-25 at the tool-call layer**; structured tool schemas convert free-form bash hallucination into typed slot-filling, with the 12.5% number as the headline evidence.
- **PRIM-29** (Recruitment-Adaptive Planning): SWE-agent's dynamic tool selection is the model for PRIM-29's recruitment loop; the typed registry (name + input + output + final-action) is the ACI pattern generalised to agent-role recruitment.
- **PRIM-21** (Migration Strategy Selection): SWE-agent's per-failure-kind branching (collection error vs assertion error vs regression) is the discriminator that `Registry.TryInOrder` consumes to pick the next repair strategy; cross-link with [[1.0.0 P-109]] §2 for the calibration dataset.
- **PRIM-22** (Four Phases Comprehension): SWE-agent's `open`/`scroll` interface is the canonical artefact for the **articulate phase**'s comprehension slot-filling — the agent reads code in windows the same way MAgHARCM's comprehension agent formats its `Returns:`/`Calls:`/`Side effects:` slots.
- **Cross-ref**: add P-112 to PRIM-5, PRIM-6, PRIM-25, PRIM-29, PRIM-21, PRIM-22 rows in `primitives/Primitives-Index.md` and `Software-Archaeology-Lineage.md` (handled by the separate "Update primitives INDEX + lineage matrix" orchestration task). P-112 is the **canonical tool-calling agent scaffold anchor** for the SLM era; together with [[1.0.0 P-109]] (Verified oracle) and [[1.0.0 P-110]] (graph retrieval substrate), it forms the evaluation-protocol pillar of MAgHARCM's agent design.

## Summary

SWE-agent (Yang et al., NeurIPS 2024) is the **canonical tool-calling agent scaffold** for repository-level software engineering, and the canonical reference for what an "agent" looks like on SWE-bench Verified. Its central contribution — **Agent-Computer Interfaces (ACIs) as a first-class design object** — establishes that the interface between the LM and the computer matters more than the model for many repository-level tasks; the same GPT-4 backbone jumps from sub-5% to 12.5% on Verified by changing only the interface. SWE-agent is hop-1 of [[1.0.0 P-109]] (SWE-bench Verified) and [[1.0.0 P-111]] (SWE-bench original), and is the scaffold that OpenHands/CodeAct, AutoCodeRover, and Aider all build on. For MAgHARCM, it anchors PRIM-5 (test-co-translation oracle feedback), PRIM-6 (multi-stage validator cascade), PRIM-25 (structural de-hallucination via typed tool schemas), and PRIM-29 (recruitment-adaptive planning via observation-driven tool selection). The ACI pattern — typed tool schemas that convert free-form shell hallucination into structured slot-filling — is the tool-call analogue of the cloze reformulation in [[1.0.0 P-55]], and is what makes a 4B–7B SLM tractable as an autonomous software-engineering agent.
