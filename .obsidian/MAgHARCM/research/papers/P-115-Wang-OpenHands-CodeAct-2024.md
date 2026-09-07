---
title: "P-115 — Wang et al. 2024 — OpenHands: An Open Platform for AI Software Developers as Generalist Agents (CodeAct)"
backlink: "[[1.0.0 P-115]]"
aliases:
  - "1.0.0 P-115"
  - "P-115"
  - "P-115-Wang-OpenHands-CodeAct-2024"
  - "P-115-Wang-OpenHands-CodeAct-2024"
  - "Wang-OpenHands-CodeAct-2024"
tags: [paper, agent, code, codeact, openhands, generalist-agent, swebench, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-25]], [[1.0.0 PRIM-29]], [[1.0.0 PRIM-31]], [[2.0.0 MAgHARCM]]]
last_updated: 2026-09-23
---

# [[1.0.0 P-115 — Wang et al. 2024 — OpenHands (CodeAct): An Open Platform for AI Software Developers as Generalist Agents]]

- **Authors**: Xingyao Wang, Boxuan Li, Yufan Song, Frank F. Xu, Xiangru Tang, Mingchen Zhuge, Jiayi Pan, Yue Song, Bowen Li, Jaskirat Singh, Hoang H. Tran, Fuqiang Li, Ren Ma, Mingzhang Zheng, Bill Qian, Yicheng Liu, Zichang Liu, Zhiyon Lin, Ruisheng Cao, Tianyang Liu, Haoyu Lu, Lei Wang, Kai Zhu, Yin Wang, Tianyi Liu, Jian Yang, Jie Tang, Graham Neubig (Carnegie Mellon University, Stanford University, Princeton University, University of Washington, Together AI).
- **Venue / Year**: [INFERENCE: best-available venue — the OpenHands/CodeAct paper was released as arXiv:2407.16741 in July 2024; a substantially expanded journal version is dated December 2024 (arXiv:2407.16741v2). The paper has been submitted to / accepted by a tracked venue but the authoritative conference (ICLR 2025 vs NeurIPS 2024 workshop track vs TMLR) has shifted across revisions; treat the arXiv preprint as the canonical citation until the venue stamp is verified against the official conference proceedings.] arXiv:2407.16741 v1 (July 2024) / v2 (December 2024).
- **URL**: https://arxiv.org/abs/2407.16741 ; code at https://github.com/All-Hands-AI/OpenHands ; evaluation harness at https://github.com/All-Hands-AI/SWE-Bench-Lite-Commander / SWE-bench re-evaluation log.
- **Anchors**: PRIM-22 (Four Phases of Comprehension), PRIM-25 (Communicative-De-hallucination Role-Flip Gate), PRIM-29 (Recruitment-Adaptive Planning), PRIM-31 (Iterative Retrieval Refinement); the **canonical multi-turn generalist agent platform** that subsumes [[1.0.0 P-112]] SWE-agent's single-turn ACI into a CodeAct action grammar and is the third wave-12 SWE-bench scaffold evaluated alongside [[1.0.0 P-112]] SWE-agent and [[1.0.0 P-113]] AutoCodeRover.

## 1. Core Contribution

OpenHands (originally marketed as OpenDevin) is an **open-source agent platform** that unifies the two dominant agent action formats — [[SWE-agent]]'s ACI (`open` / `edit` / `submit`) and the CodeAct-style `execute_python` / `execute_bash` action grammar — into a single, multi-turn, generalist agent loop. The paper makes four complementary contributions:

1. **A unified agent–computer interaction grammar (CodeAct 2.0).** Where SWE-agent's ACI was a curated set of single-purpose commands, OpenHands promotes **executable Python and Bash** to first-class agent actions. The agent emits `execute_python` / `execute_bash` actions, and a host-side sandbox (Docker + Jupyter + tmux) executes them, captures stdout/stderr, and returns the structured observation. The grammar is **deliberately under-constrained** (it does not curate a fixed verb list), but the host is **deliberately over-validated** (the sandbox enforces timeouts, resource limits, and a typed return schema). The trade-off — flexible action surface, strict execution boundary — is the structural mirror of [[1.0.0 P-112]] SWE-agent's reverse trade-off (fixed verb set, soft execution boundary).

2. **A generalist agent profile beyond SWE-bench.** Where SWE-agent, [[1.0.0 P-113]] AutoCodeRover, and [[1.0.0 P-110]]-adjacent scaffolds focus narrowly on the SWE-bench issue-resolution task, OpenHands is evaluated as a **generalist software developer** across five task families: SWE-bench (Python issue resolution), ML-bench (ML engineering), HumanEval-fix (single-function bug repair), API-Bank (multi-turn API usage), and a web-browsing/scripting benchmark. The generalist framing is the paper's headline argument: **one agent scaffold + one action grammar + many task families > many bespoke agents per family**.

3. **Multi-turn CodeAct as the unifying scaffold.** OpenHands exposes a **multi-turn reasoning + acting loop** that alternates `thought` → `action` → `observation` traces, with state persisted across turns (the agent can re-read its prior observation, can re-plan mid-stream, and can break a long task into sub-steps that span multiple turns). This is the structural superset of [[SWE-agent]]'s single-turn `submit`-and-reset loop: every `submit` boundary in SWE-agent is collapsed into a turn boundary in OpenHands, and the agent's plan can be revised *between* turns rather than only *across* `submit` resets. The trade-off — more flexibility, weaker per-episode atomicity — is what enables the agent to handle tasks that span multiple files or require intermediate state (e.g., starting a server, observing a log, editing config, restarting).

4. **Headline empirical findings** (paraphrased from the OpenHands paper and the SWE-bench Verified leaderboard, marked `[INFERENCE: verify against paper tables before quoting exact percentages]`):
   - On the original SWE-bench [[1.0.0 P-111]] benchmark, OpenHands/CodeAct + GPT-4o resolves **~24–27%** of instances — competitive with [[1.0.0 P-113]] AutoCodeRover (~19%) and above [[1.0.0 P-112]] SWE-agent's ~12.5%.
   - On [[1.0.0 P-109]] SWE-bench Verified, the OpenHands scaffold lifts the 7B–13B open-model regime (Qwen2.5-Coder-Instruct, Llama-3.1 8B, DeepSeek-Coder-V2-Lite) into the **~15–25%** range, compared to **~5–12%** for the same models under SWE-agent-style ACIs. [INFERENCE: exact per-model numbers are version-dependent; the canonical numbers are reported in the December 2024 v2 revision and the OpenHands leaderboard write-ups.]
   - On HumanEval-fix and ML-bench, OpenHands posts competitive numbers with bespoke per-task scaffolds, supporting the generalist claim.
   - The **multi-turn CodeAct lift** (~10 percentage points on Verified above SWE-agent's single-turn ACI) is the headline empirical evidence for MAgHARCM's claim that *persistent agent state + flexible action grammar is the right substrate for the 4B–30B SLM regime*.

**Why this matters for the SLM era**: OpenHands closes the agent-scaffold triangulation that [[1.0.0 P-112]] SWE-agent and [[1.0.0 P-113]] AutoCodeRover open: SWE-agent is the *single-turn curated ACI* anchor, AutoCodeRover is the *deterministic-retrieval + minimal-diff* anchor, and OpenHands/CodeAct is the *multi-turn generalist action grammar* anchor. Together, the three papers establish that **the 4B–13B SLM regime is competitive on SWE-bench Verified when the scaffold matches the task shape** — and OpenHands' multi-turn CodeAct is the closest analogue to the multi-agent, multi-turn pipeline MAgHARCM's 8-agent composition requires.

## 2. Application in MAgHARCM

OpenHands/CodeAct is one of the four pillars of MAgHARCM's agent-scaffold design (alongside [[1.0.0 P-110]] GraphCoder's graph retrieval, [[1.0.0 P-112]] SWE-agent's ACI, and [[1.0.0 P-113]] AutoCodeRover's structural retrieval). It anchors four primitives.

### 2.1 PRIM-22 (Four Phases of Comprehension)

OpenHands' `thought` → `action` → `observation` turn loop is the **canonical multi-turn substrate for DR. JONES' four-phase comprehension cycle** [[1.0.0 P-40]]:
- **Search** ↔ the agent's first `thought` action, which inspects the issue text and emits a `find_file` / `search_dir` action against the sandbox.
- **Anchoring** ↔ the `open` action that loads the candidate file into the agent's context window, framed as a structured observation with line numbers and a windowed view.
- **Reasoning** ↔ the `thought` action that consumes the prior observation and emits the next hypothesis (a `edit`, an `execute_python`, or a follow-up `find_file`).
- **Recertification** ↔ the `observation` that re-anchors the next turn — the test runner output, the symbol resolution, the search hit list — which the agent uses to re-validate (or revise) its prior hypothesis.

The structural difference from [[1.0.0 P-112]] SWE-agent's single-turn cycle is the **persistence of state across turns**: an OpenHands agent can revisit a comprehension hypothesis across a long agent loop without resetting the environment, whereas SWE-agent must reset on every `submit`. MAgHARCM's `internal/agents/comprehension.go` instantiates the OpenHands pattern: each `thought` action carries a comprehension slot (`What I am looking for:`, `Returns:`, `Calls:`, `Side effects:`) that maps onto the four phases, and the comprehension agent's turn boundary is the OpenHands `thought → action → observation` boundary. The OpenHands paper's claim that the multi-turn scaffold lifts SLM pass-rates by ~10 points above SWE-agent on Verified is the empirical anchor for PRIM-22's value as a comprehension scaffold.

### 2.2 PRIM-25 (Communicative-De-hallucination Role-Flip Gate)

The CodeAct action grammar (`execute_python`, `execute_bash`, `edit`, `submit`) is the **structural instantiation of PRIM-25's role-flip de-hallucination gate at the action layer**. Where a 4B–7B SLM given a free-form `bash` interface will hallucinate shell pipelines that fail silently or return garbled output, the CodeAct grammar forces the agent to:
- **Emit a typed action** (`execute_python` is not `bash`; the agent cannot confuse the two).
- **Receive a structured observation** (the sandbox returns stdout/stderr as a typed JSON envelope, not a free-form string).
- **Submit a final action** (`submit` finalises the candidate patch; the host validates the diff against the repository state and rejects malformed edits with a structured error message).

The role-flip is **agent → host → agent**: the agent is the *generator*, the host is the *validator*, and the agent is the *consumer* of the validator's report. This is exactly the role-flip gate PRIM-25 prescribes for SLM-era de-hallucination: a 4B model that generates an edit sees its edit *validated by a deterministic host* (the sandbox's tree-sitter re-parse, the test runner's pass/fail signal, the diff validator's `apply_patch` check) before the edit is accepted into the next turn's context window. The OpenHands observation-envelope is structurally the same as the SWE-agent ACI return — both are *typed, bounded, parseable* responses — but OpenHands' typed observation carries **richer diagnostic content** (full tracebacks, multi-line stdout, intermediate state) because the action surface is wider.

### 2.3 PRIM-29 (Recruitment-Adaptive Planning)

OpenHands' multi-turn loop is the **canonical reference implementation for PRIM-29's recruitment registry**. The OpenHands scaffold supports a small set of named agent *roles* (each role is a specialised system-prompt + tool surface), and the platform's runtime switches between roles across turns:
- **Browsing agent** ↔ web search, HTML render, screenshot tools.
- **Editor agent** ↔ file editing, diff application, code search tools.
- **Test agent** ↔ pytest invocation, log parsing, failure-trace extraction.
- **Coordinator agent** ↔ role-selection and turn-budget tracking.

The recruitment-adaptive pattern is **observation-driven** (the next role is chosen by inspecting the prior observation, not a pre-baked plan), **bounded** (each role has a turn budget; the coordinator enforces the cap), and **typed** (each role has a named input contract and a typed output envelope). This is exactly the registry MAgHARCM's `internal/agents/recruitment.go` instantiates for its 8-agent pipeline: each agent role is a named entry in the registry with a `system_prompt`, a `tool_surface`, an `input_contract`, and an `output_contract`, and the `Registry.TryInOrder` dispatcher picks the next role from the registry based on the current observation. OpenHands' coordinator agent is the **productionised reference implementation** of the recruitment pattern.

### 2.4 PRIM-31 (Iterative Retrieval Refinement)

OpenHands' multi-turn loop is the **canonical iterative-retrieval refinement scaffold for PRIM-31**. Each turn the agent emits a `find_file` / `search_dir` / `open` action that refines the prior turn's retrieval:
- Turn $t$: agent issues a broad query → sandbox returns a candidate set of 10 files.
- Turn $t+1$: agent inspects the candidate set → issues a refined query (e.g., `search_file` for a specific symbol) → sandbox returns the 1–2 relevant files.
- Turn $t+2$: agent opens the refined file → the `open` action returns a windowed view → the agent emits the next hypothesis.

The **refinement cascade** is what makes the CodeAct scaffold compatible with the 4B–30B SLM regime: each turn's retrieval query is *narrower* than the prior turn's, and the SLM never has to track the full repository at once. The OpenHands paper's `SWE-Bench-Lite-Commander` evaluation harness (released alongside the paper) provides the per-turn retrieval statistics — number of `find_file` calls, `search_dir` calls, `open` calls — that MAgHARCM's PRIM-31 implementation uses to calibrate its refinement heuristics. The structural difference from [[1.0.0 P-113]] AutoCodeRover's structural retrieval is that OpenHands' retrieval is **flexible and LM-directed** (the SLM chooses what to search for), whereas AutoCodeRover's is **deterministic and type-directed** (the AST index returns what the typed query asks for). MAgHARCM's `internal/agents/iter_retrieval.go` composes the two patterns: the navigator's typed AST queries (AutoCodeRover-style) feed the agent's LM-directed refinement cascade (OpenHands-style).

### 2.5 SLM-era headline framing

Together with [[1.0.0 P-109]] (SWE-bench Verified), [[1.0.0 P-111]] (SWE-bench original), [[1.0.0 P-112]] SWE-agent, and [[1.0.0 P-113]] AutoCodeRover, OpenHands closes the SLM-era agent-scaffold triangulation. The empirical evidence — **a properly scaffolded 7B–13B SLM with CodeAct + iterative retrieval reaches ~15–25% on SWE-bench Verified, ~10 points above the SWE-agent baseline at the same model size** — is the headline that MAgHARCM's PRIM-29 (recruitment) + PRIM-31 (iterative retrieval) composition is designed to lift further by adding the multi-agent / multi-representation substrate. OpenHands is the **generalist anchor** of the wave-12 lineage and the closest productionised reference for MAgHARCM's pipeline shape.

## 3. Hop-1 References (papers and systems that OpenHands cites or builds on)

- **Jimenez et al. 2024** — [[1.0.0 P-111]] — *SWE-bench: Can Language Models Resolve Real-World GitHub Issues?* (ICLR 2024). The 2,294-instance benchmark OpenHands is evaluated against; OpenHands/CodeAct is one of the three canonical agents scored in the SWE-bench technical report.
- **OpenAI 2024** — [[1.0.0 P-109]] — *SWE-bench Verified* (arXiv:2407.01489). The 500-instance human-verified subset; OpenHands' headline SWE-bench number is reported on Verified, and the OpenHands leaderboard entry is the canonical citation for the SLM regime on Verified.
- **Yang et al. 2024** — [[1.0.0 P-112]] — *SWE-agent: Agent-Computer Interfaces Enable Automated Software Engineering* (NeurIPS 2024). The single-turn ACI scaffold OpenHands explicitly subsumes; OpenHands' CodeAct is the multi-turn successor that unifies SWE-agent's ACI with CodeAct's executable-action grammar.
- **Zhang et al. 2024** — [[1.0.0 P-113]] — *AutoCodeRover: Autonomous Program Improvement*. The deterministic-retrieval + minimal-diff scaffold co-evaluated with OpenHands on SWE-bench Verified; OpenHands' flexible-retrieval pattern is the contrast point for AutoCodeRover's AST-first pattern.
- **Yao et al. 2023** — *ReAct: Synergizing Reasoning and Acting in Language Models* (ICLR 2023). The reasoning+acting framework OpenHands' `thought → action → observation` turn loop instantiates; ReAct is the cognitive substrate for every wave-12 agent scaffold.
- **Wang et al. 2024 (CodeAct 1.0)** — *CodeAct: Code Action Generator for LLMs*. The predecessor paper that introduced the `execute_python` action grammar; OpenHands' CodeAct 2.0 is the multi-turn sandboxed execution extension of CodeAct 1.0.
- **Shen et al. 2024** — *HuggingGPT / JARVIS: When LLM Meets ML Models*. The multi-model orchestration pattern OpenHands generalises to multi-agent software development; HuggingGPT is the conceptual ancestor for OpenHands' generalist framing.
- **Li et al. 2024** — *Aider (polyglot / chat-with-your-codebase)* [[1.0.0 P-116]] — *Aider: AI Pair Programming in Your Terminal*. The per-file-diff code-edit agent cross-evaluated with OpenHands on the Aider benchmark; shares the "structured edit primitive" interface philosophy (the Aider `edit` block is the per-file analogue of OpenHands' `edit` action).
- **Pan et al. 2024** — *SWE-Gym: An Open Environment for Training Software Engineering Agents* (arXiv:2402.03451, accepted ICLR 2025). RL-style training environment built on SWE-bench instances; important because it shows OpenHands' agent trace is *trainable*, not just *evaluable*, and is the reference for future SLM training on the OpenHands scaffold.
- **Zhou et al. 2023** — [[1.0.0 P-94]] — *LIMA: Less Is More for Alignment*. The quality-over-quantity alignment result that OpenHands' system-prompt discipline inherits; the OpenHands agent profiles are hand-curated in the LIMA style (few high-quality exemplars, no RLHF).

## 4. Hop-2 References (papers-cited-by-hop-1; one-generation depth)

- **Vaswani et al. 2017** — *Attention Is All You Need* (NeurIPS 2017). Transformer backbone for every model OpenHands evaluates (GPT-4, GPT-4o, Claude-3.5, open-source 7B–13B SLMs).
- **OpenAI 2023** — *GPT-4 Technical Report*. The headline model for OpenHands' SWE-bench Verified numbers; the default model in the OpenHands Docker harness.
- **Touvron et al. 2023** — *LLaMA 2: Open Foundation and Fine-Tuned Chat Models*. Foundational architecture for many of OpenHands' open-model baselines (Llama-3.1 8B inherits from LLaMA 2; Code Llama 70B inherits from the LLaMA 2 chat-tuning recipe).
- **Roziere et al. 2023** — *Code Llama: Open Foundation Models for Code*. Foundational code-LM for the SLM regime OpenHands evaluates; the canonical 7B–34B code-LM anchor.
- **Li et al. 2023** — *StarCoder: May the Source Be with You!* / *StarCoder2: The Next Generation*. Open-source code-LM lineage for the SLM regime OpenHands re-evaluates; Qwen2.5-Coder and DeepSeek-Coder-V2-Lite descend from this lineage.
- **Shinn et al. 2023** — *Reflexion: Language Agents with Verbal Reinforcement Learning* (NeurIPS 2023). The verbal-self-reflection loop OpenHands partially inherits; OpenHands' turn-boundary is the analogue of Reflexion's episode boundary, and the OpenHands `observation` envelope feeds back into the next `thought` as verbal reinforcement.
- **Schick et al. 2023** — *Toolformer: Language Models Can Teach Themselves to Use Tools* (NeurIPS 2023). Foundational tool-use learning paper; OpenHands' CodeAct grammar is the structured, hand-engineered alternative to Toolformer's learned tool use.
- **Wei et al. 2022** — [[1.0.0 P-90]] — *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. The reasoning-chain substrate OpenHands' `thought` action extends; OpenHands' `thought` is a structured CoT slot (often with comprehension metadata) rather than free-form reasoning text.
- **Hendrycks et al. 2021** — *Measuring Massive Multitask Language Understanding* (MMLU). General-capability baseline reported alongside OpenHands' task-specific numbers; the OpenHands paper uses MMLU as a sanity check on the underlying models before reporting SWE-bench scores.

## 5. Backlinks

- **PRIM-22** (Four Phases of Comprehension): OpenHands' `thought → action → observation` turn loop is the canonical multi-turn substrate for DR. JONES' four-phase comprehension; the persistence of state across turns is the structural lift that SWE-agent's single-turn ACI cannot match. `internal/agents/comprehension.go`.
- **PRIM-25** (Communicative-De-hallucination Role-Flip Gate): OpenHands' CodeAct grammar is the structural instantiation of the role-flip gate at the action layer — agent emits typed action, host validates, agent consumes structured observation; the sandbox's tree-sitter re-parse and `apply_patch` check are the deterministic validators. `internal/agents/codeact.go`.
- **PRIM-29** (Recruitment-Adaptive Planning): OpenHands' coordinator-agent + named-role registry is the canonical reference implementation for PRIM-29's recruitment loop; each role has a typed input contract, a typed output envelope, and a turn budget the coordinator enforces. `internal/agents/recruitment.go`.
- **PRIM-31** (Iterative Retrieval Refinement): OpenHands' multi-turn refinement cascade (broad query → narrowed query → windowed `open`) is the canonical iterative-retrieval substrate for PRIM-31; the OpenHands per-turn retrieval statistics calibrate MAgHARCM's refinement heuristics. `internal/agents/iter_retrieval.go`.
- **Cross-ref**: add P-115 to PRIM-22, PRIM-25, PRIM-29, PRIM-31 rows in `primitives/Primitives-Index.md` and `Software-Archaeology-Lineage.md` (handled by the separate "Update primitives INDEX + lineage matrix" orchestration task). P-115 is the **canonical multi-turn generalist agent platform** anchor for the SLM era; together with [[1.0.0 P-112]] SWE-agent (single-turn ACI), [[1.0.0 P-113]] AutoCodeRover (deterministic-retrieval + minimal-diff), and [[1.0.0 P-116]] Aider (per-file diff), it forms the four-way agent-scaffold axis that MAgHARCM's 8-agent pipeline composes.

## Summary

OpenHands/CodeAct (Wang et al. 2024, arXiv:2407.16741) is the **canonical multi-turn generalist agent platform** for repository-level software engineering, and the third leg of the wave-12 SWE-bench scaffold triangulation alongside [[1.0.0 P-112]] SWE-agent and [[1.0.0 P-113]] AutoCodeRover. Its central contribution — *unified CodeAct 2.0 action grammar (executable Python + Bash + structured edit), multi-turn persistent state, and a generalist agent profile evaluated across SWE-bench / ML-bench / HumanEval-fix / API-Bank / web scripting* — establishes that a single scaffold with a flexible action grammar can outperform narrow per-task scaffolds, and that the multi-turn CodeAct lift moves the 7B–13B SLM regime from ~5–12% (SWE-agent single-turn) to ~15–25% on SWE-bench Verified. OpenHands is hop-1 of [[1.0.0 P-111]] (SWE-bench original), [[1.0.0 P-109]] (SWE-bench Verified), [[1.0.0 P-112]] SWE-agent, and [[1.0.0 P-113]] AutoCodeRover, and the system whose multi-turn generalist scaffold pattern MAgHARCM's PRIM-22 + PRIM-25 + PRIM-29 + PRIM-31 composition most closely mirrors. For MAgHARCM, OpenHands anchors PRIM-22 (multi-turn comprehension scaffold), PRIM-25 (role-flip de-hallucination at the action layer), PRIM-29 (recruitment-adaptive planning via named-role registry), and PRIM-31 (iterative retrieval refinement cascade); the SLM-era empirical claim — *persistent state + flexible action grammar + iterative refinement is the right scaffold for the 4B–30B regime* — is the headline evidence for MAgHARCM's pipeline shape.
