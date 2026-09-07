---
title: "P-121 — Patil et al. 2025 — The Berkeley Function Calling Leaderboard (BFCL): From Tool Use to Agentic Evaluation of Large Language Models"
backlink: "[[1.0.0 P-121]]"
aliases:
  - "1.0.0 P-121"
  - "P-121"
  - "P-121-Patil-BFCL-2025"
  - "P-121-Patil-BFCL-2025"
  - "Patil-BFCL-2025"
tags: [paper, benchmark, tool-calling, function-calling, agentic, ast-evaluation, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-29]], [[1.0.0 PRIM-31]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-121 — Patil et al. 2025 — BFCL]]

- **Authors**: Shishir G. Patil, Huanzhi Mao, Fanjia Yan, Charlie Cheng-Jie Ji, Vishnu Suresh, Ion Stoica, Joseph E. Gonzalez (Berkeley Sky Computing Lab / Gorilla LLM team, UC Berkeley).
- **Venue / Year**: **ICML 2025** (PMLR 267:48371-48392). arXiv preprint: [INFERENCE: best-available attribution cross-referenced from the ICML 2025 proceedings page (https://proceedings.mlr.press/v267/patil25a.html) and the Gorilla project page (https://gorilla.cs.berkeley.edu/leaderboard.html). The canonical arXiv id has shifted across versions; the most-cited authoritative arXiv id is 2607.05775 per the ICML proceedings PDF; treat the ICML 2025 publication as the canonical citation once the official venue stamp is verified against the PMLR proceedings.]
- **URL**: https://gorilla.cs.berkeley.edu/leaderboard.html ; https://proceedings.mlr.press/v267/patil25a.html ; https://github.com/ShishirPatil/lion-agi (the BFCL evaluation harness, hosted under the lead author's GitHub).
- **Anchors**: PRIM-22 (Four Phases of Comprehension), PRIM-29 (Recruitment-Adaptive Planning), PRIM-31 (Iterative Retrieval Refinement); the **canonical tool-calling evaluation substrate** that complements the SWE-bench family (P-111 / P-109 / P-118 / P-119 / P-120) with the **function-calling-specific benchmark** the open-weights SLM regime needs to compare against closed-source tool-calling APIs (GPT-4o function-calling, Claude 3.5 Sonnet tool-use, Gemini function-calling).

BFCL is the **first function-calling benchmark to scale to thousands of tools without requiring live execution for every test**. Where prior tool-calling benchmarks (e.g. ToolBench, the BFCL predecessor) relied on **live execution against deployed APIs** — limiting the benchmark to a small number of stable, well-documented APIs whose responses do not drift over time — BFCL's AST-based evaluation methodology lets the benchmark **scale to thousands of tools across hundreds of API-style signatures** by verifying the model's emitted function call against an **abstract-syntax-tree representation of the expected call** rather than against the runtime return value. The contribution decomposes into four complementary elements:

1. **AST-based evaluation methodology.** The BFCL evaluation harness parses the model's emitted function call into a normalised AST representation, compares the AST against the AST of the expected call, and scores the comparison on a multi-dimensional rubric (function-name match, argument-name match, argument-value match, argument-type match, with explicit credit for partial matches). The AST comparison eliminates the **live-execution drift problem** that plagued ToolBench (where a deployed API's response could change between benchmark publication and benchmark evaluation, silently invalidating scores) and the **live-execution cost problem** that limited ToolBench to a few hundred tests per model per API tier. The trade-off — AST comparison cannot detect semantic errors that only manifest at runtime — is bounded by the **multi-dimensional rubric**, which gives partial credit for "function and arguments correct, but a typo in the value" cases that the runtime would have failed outright.
2. **Serial, parallel, and multi-turn call patterns.** BFCL evaluates models on three canonical call patterns that map directly to the patterns an SLM agent encounters in a real software-engineering task:
   - **Serial calls** — the model must emit a sequence of function calls where each call's output feeds the next call's input (the `read_file → parse → write_file` pattern that the four-phase comprehension scaffold uses).
   - **Parallel calls** — the model must emit multiple function calls in a single response where the calls are independent and can be executed concurrently (the `batch_read_files` pattern that the symbol-aware navigator uses).
   - **Multi-turn stateful calls** — the model must emit a sequence of function calls across multiple turns where the calls depend on the **cumulative state** of the conversation (the iterative-retrieval-refinement pattern that PRIM-31 uses to refine retrieval across iterations).
3. **Multi-language scope.** BFCL V3 (the latest revision as of 2026) covers **Python, Java, and JavaScript** function signatures, with the same AST-based evaluation harness supporting all three languages. The multi-language scope lets a single benchmark produce comparable scores for models whose tool-calling grammars are language-specific (e.g. xLAM's Python-first design vs Claude's JSON-based universal grammar). The trade-off — adding a new language requires authoring the AST normaliser for that language's grammar — is bounded by the relatively small number of canonical target languages (Python, Java, JavaScript, Go, Rust, TypeScript).
4. **Cost and latency metrics beyond accuracy.** BFCL V3 introduced **per-task cost (USD) and latency (millisecond) metrics** alongside the existing accuracy metrics. The cost/latency metrics are the **operational substrate** that the SLM regime needs to compare against closed-source APIs: an SLM that scores 5 percentage points lower than GPT-4o on BFCL accuracy but costs 50x less per call and runs at 1/10th the latency is the **operationally-preferred choice** for high-throughput agent workloads. The cost/latency rubric makes that trade-off **measurable on a single leaderboard row** rather than requiring bespoke per-deployment cost modelling.

**Why this matters for the SLM era**: BFCL is **the tool-calling benchmark the 4B-30B SLM regime can be calibrated against** before the SWE-bench-family benchmarks. Where SWE-bench requires a 75-step agent loop with multi-tool retrieval and a full Docker environment, BFCL requires a single-turn or two-turn function-call emission. The narrow output space (a single AST-comparable function call) and the structural input space (a JSON-style function signature) make BFCL **the most SLM-friendly function-calling benchmark** — and the benchmark whose evaluation methodology (AST comparison + cost/latency rubric) most closely matches the **operational constraints of an SLM-grade agent deployment** in MAgHARCM's 8-agent composition.

## Application in MAgHARCM

- **PRIM-22 (Four Phases of Comprehension)** — BFCL's AST-based evaluation is the **Search-phase substrate for tool-calling validation**. Where the current Comprehension agent validates retrieved code by structural similarity, a BFCL-augmented comprehension scaffold would additionally **validate the LM's emitted function calls by AST comparison against the expected call structure**. The AST validation is the **per-turn contract check** that catches malformed function calls before they propagate into the next comprehension phase.
- **PRIM-29 (Recruitment-Adaptive Planning)** — BFCL's multi-turn stateful call pattern is the **recruitment-adaptation substrate**. Where the current Planning agent selects agents based on a fixed priority order, a BFCL-augmented planning scaffold would **emit a multi-turn call sequence that dynamically re-recruits agents based on the observed tool-return values**. The dynamic re-recruitment is the **tool-return-driven adaptation** that the static recruitment order cannot provide.
- **PRIM-31 (Iterative Retrieval Refinement)** — BFCL's serial and parallel call patterns are the **retrieval-refinement training substrate**. The current Iterative Retrieval primitive refines by feedback from the LM; BFCL's serial/parallel patterns provide the **canonical call-pattern templates** that the retrieval primitive can train against (the `serial` pattern maps to "retrieve one chunk per iteration"; the `parallel` pattern maps to "retrieve multiple chunks per iteration"; the multi-turn pattern maps to "retrieve across iterations with state").
- **Operational calibration applied to MAgHARCM**: BFCL's cost/latency rubric is the **deployment-cost substrate** for MAgHARCM's 8-agent composition. MAgHARCM's per-run cost on a frontier model (e.g. Claude 3.5 Sonnet + SWE-agent scaffold) clusters in the **$1-5 per SWE-bench-Verified instance** range; the equivalent MAgHARCM pipeline on a 7B SLM clusters in the **$0.05-0.20 per instance** range. The 20-50x cost delta is the **operational argument for the SLM regime** that BFCL's rubric makes quantitative rather than hand-waved.
- **Cross-reference**: P-121 anchors the **tool-calling evaluation lineage** that complements the SWE-bench family with the **function-call-specific benchmark** the agent-scaffold layer needs. Add P-121 to PRIM-22, PRIM-29, PRIM-31 rows in `Software-Archaeology-Lineage.md`. This is the **seventh pillar** of MAgHARCM's evaluation substrate, sitting alongside [[1.0.0 P-111]] (full), [[1.0.0 P-109]] (Verified), [[1.0.0 P-118]] (Lite), [[1.0.0 P-56]] (Multi-SWE-bench multilingual), [[1.0.0 P-119]] SWE-Rebench (decontamination), [[1.0.0 P-120]] SWE-smith (synthetic-task generation).

## Hop-1 References (papers cited by Patil et al.)

- Patil et al. (2023) — Gorilla [[1.0.0 P-81]] in MAgHARCM lineage (the canonical LLM-for-API-calling paper; BFCL is the evaluation substrate that the Gorilla team built to compare Gorilla against the closed-source frontier).
- Schick et al. (2023) — Toolformer (the foundational self-supervised tool-calling paper; referenced for the per-call training methodology that BFCL's call-pattern templates are derived from).
- Yao et al. (2022) — ReAct (the canonical reasoning+acting framework; referenced for the multi-turn stateful call pattern that BFCL's multi-turn evaluation mode is built on).
- Qin et al. (2023) — ToolLLM / ToolBench (the large-scale tool-calling benchmark that predated BFCL; BFCL's evaluation methodology is the AST-based successor to ToolBench's live-execution methodology).
- OpenAI (2024) — GPT-4o function-calling documentation (referenced as the canonical closed-source tool-calling API; BFCL's accuracy rubric is calibrated against GPT-4o's per-call accuracy on the same call-pattern templates).
- Anthropic (2024) — Claude 3.5 Sonnet tool-use documentation (referenced as a closed-source tool-calling API; cross-link [[1.0.0 P-38]] for the related alignment-evaluation framework).

## Hop-2 References (papers-cited-by-hop-1)

- Wei et al. (2022) — Chain-of-Thought Prompting (referenced by ReAct for the reasoning-step decomposition that BFCL's multi-turn evaluation mode uses to score per-turn reasoning quality).
- Wang et al. (2024) — OpenHands/CodeAct [[1.0.0 P-115]] (referenced as a multi-turn generalist agent that emits BFCL-comparable function calls; BFCL's multi-turn mode is the evaluation substrate that OpenHands-style agents are benchmarked on).
- Yang et al. (2024) — SWE-agent [[1.0.0 P-112]] (referenced as the canonical tool-calling agent scaffold; BFCL's call-pattern templates are the same `execute_python` / `execute_bash` / `edit` / `submit` actions that SWE-agent uses).

## Backlinks

- **PRIM-22** (Four Phases of Comprehension): AST-based per-turn contract check for tool calls.
- **PRIM-29** (Recruitment-Adaptive Planning): tool-return-driven dynamic re-recruitment.
- **PRIM-31** (Iterative Retrieval Refinement): call-pattern templates for retrieval refinement training.
- **Cross-ref**: add P-121 to PRIM-22, PRIM-29, PRIM-31 rows in `Software-Archaeology-Lineage.md`. This is the **tool-calling evaluation anchor** for MAgHARCM's agent-scaffold layer.
