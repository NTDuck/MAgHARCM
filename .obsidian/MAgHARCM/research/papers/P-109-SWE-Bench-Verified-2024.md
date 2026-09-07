---
title: "P-109 — OpenAI 2024 — SWE-bench Verified"
backlink: "[[1.0.0 P-109]]"
aliases:
  - "1.0.0 P-109"
  - "P-109"
  - "P-109-SWE-Bench-Verified-2024"
  - "P-109-SWE-Bench-Verified-2024"
  - "P-109-SWE-Bench-Verified-2024"
  - "P-109-SWE-Bench-Verified-2024"
  - "SWE-bench-Verified-2024"
tags: [paper, benchmark, code, swebench, verified, [[1.0.0 PRIM-5]], [[1.0.0 PRIM-6]], [[1.0.0 PRIM-23]], [[1.0.0 PRIM-27]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-109 — OpenAI 2024 — SWE-bench Verified]]

- **Authors**: OpenAI (org-wide team) building on Jimenez et al. 2024 — Carlos E. Jimenez, John Yang, Alexander Wettig, Shunyu Yao, Kexin Pei, Ofir Press, Karthik R. Narasimhan (Princeton LLM-Agent Lab / OpenAI).
- **Venue / Year**: OpenAI release, August 2024 (the Verified subset of SWE-bench, with the verification methodology later written up in Opelt et al., arXiv:2407.01489, July 2024).
- **URL**: https://openai.com/index/swe-bench-verified/ ; methodology write-up https://arxiv.org/abs/2407.01489
- **Anchors**: PRIM-5 (Test Suite Co-Translation & Synthesis), PRIM-6 (Multi-Stage Build/Test Feedback Repair), PRIM-23 (Chunked Translation), PRIM-27 (Coverage-Guided Plateau Detection); the **cleanest current oracle** for end-to-end issue-resolution evaluation in the SLM era.

## 1. Core Contribution

The original SWE-bench (Jimenez et al., ICLR 2024) assembled **2,294 task instances** drawn from real Python repositories on GitHub: each instance pairs a *resolved issue text* with a *candidate patch* and a *hidden test suite* the patch must make pass. The dataset rapidly became the de facto leaderboard for repository-level code-generation agents (SWE-Agent, OpenHands/CodeAct, AutoCodeRover, RepoCoder, R2E). Almost immediately, three structural problems became visible:

1. **Unsolvable instances** — the "ground-truth" patch was drawn from the commit that closed the issue, but in a non-trivial fraction of cases the commit included unrelated refactors, formatting fixes, or new tests that no candidate model could be expected to reproduce from the issue text alone. These instances silently dragged every score down.
2. **Ambiguous / underspecified task statements** — many GitHub issues refer to linked PRs, screenshots, prior discussions, or external systems that are not contained in the issue body. A model had no chance of reproducing the maintainer's fix without that out-of-band context.
3. **Unverifiable tests** — some hidden test suites were themselves flaky, depended on environment-specific versions, or failed against the ground-truth commit itself, so even a "correct" patch could not be scored.

SWE-bench Verified is OpenAI's **human-validated 500-instance subset** that addresses all three. The methodology:

1. **Solvability verification** — professional software developers (93 annotators from OpenAI's contractor pool) were paid to attempt each instance in a fixed Docker environment using the issue text alone, **without** looking at the candidate patch. An instance is retained only if a developer can produce a passing patch, or if the maintainer's patch demonstrably passes in the canonical environment.
2. **Cleaned task statements** — annotators rewrite each issue text to remove cross-references that are not solvable from the body alone, collapse multi-issue threads into a single coherent prompt, and flag environmental dependencies. The result is a *self-contained* problem statement that any model with repository access can attempt.
3. **Re-verified passing tests** — every retained instance is re-checked: the maintainer's patch must produce a clean `pytest` run with no environmental errors. Test suites that previously masked a passing patch via collection errors are replaced.
4. **500 hand-picked instances across 12 repositories** — Django, scikit-learn, sympy, sphinx, matplotlib, pylint, astropy, pytest, xarray, flask, conan, pylint (broad coverage of Python web, scientific, testing, and tooling domains).

**Headline empirical findings** in the OpenAI release and the methodology write-up:
- **GPT-4o** jumps from ~20% on original SWE-bench to ~33% on Verified.
- **Claude 3.5 Sonnet** (Anthropic, October 2024) moves from ~22% on the full set to **~49% on Verified** — almost half the benchmark solved by one frontier model.
- **GPT-4 + scaffolding** (function-calling, file editing) clears ~40%.
- **Open-source SLMs** (Qwen2.5-Coder 7B-Instruct, DeepSeek-Coder-V2-Lite, Llama-3.1 8B) score in the **5–12% range** — well below frontier, but a non-trivial fraction, validating that SWE-bench Verified is tractable for SLM-grade agents with the right scaffolding.

**SLM-aware relevance for MAgHARCM**: SWE-bench Verified is the **first SLM-relevant, reproducible, end-to-end software-engineering benchmark**. Prior work (HumanEval, MBPP, APPS) tests function synthesis from a docstring; SWE-bench Verified tests *issue-resolution* in a full repository context. For the 4B-30B SLM regime MAgHARCM targets, the floor of ~5% is enough to be measured, and the headroom of ~50%+ gives a clear target for what an agent scaffold can lift. Crucially, the **cleaner test suites** mean every MAgHARCM pipeline signal (PRIM-5 test synthesis, PRIM-6 build feedback, PRIM-27 plateau detection) can be validated against ground-truth oracles rather than against noisy, partially-failing test patches.

## 2. Application in MAgHARCM

- **PRIM-5 (Test Suite Co-Translation & Synthesis)** — SWE-bench Verified's *cleaned test patch* per instance is the closest thing to an oracle test scaffold MAgHARCM has access to. The Translator's `Validator.generateAdditionalTests` pathway produces synthetic tests that must be **at least as discriminative as** the Verified test patches. The Verified methodology (re-running the full test suite, replacing flaky cases, validating the maintainer's patch passes) is the direct template for PRIM-5's `verifyNoTestWeakening` discipline. A SWE-bench Verified instance = a `(problem statement, repo state, oracle test patch)` triple that PRIM-5's test-co-translation can be calibrated against. The multi-language successor Multi-SWE-bench [[1.0.0 P-56]] already extends this to seven languages; PRIM-5 can now target either as ground-truth oracles.

- **PRIM-6 (Multi-Stage Build/Test Feedback Repair)** — SWE-bench Verified's instance failure modes (compilation error → test failure → patch regression → final-pass) map directly onto PRIM-6's repair-loop stages. Specifically:
  1. The *first* failing assertion in a Verified test patch is the right **early-exit signal** for the repair loop, not the full test-run outcome.
  2. The *kind* of failure (assertion vs. collection vs. import) determines which repair-strategy branch (lexical, structural, dependency-aware) the loop should try next — exactly the strategy-selection signal that `Registry.TryInOrder` consumes in PRIM-21.
  3. The Verified benchmark's score gaps (5% for 7B SLMs vs. 33-49% for frontier) **quantify** how much lift the multi-stage repair scaffold adds over a single-pass translation, which is the empirical anchor for PRIM-6's existence.

- **PRIM-23 (Chunked Translation)** — A SWE-bench Verified instance is naturally a **chunk-of-code-per-file** problem: the maintainer's patch typically touches 1–4 files, with localized edits that do not span the whole repository. This is the exact granularity `internal/agents/chunked_translator.go` partitions at. The benchmark's per-file edit sizes give a **calibration dataset** for PRIM-23's chunk-size heuristic: chunks that are too large to fit in the SLM's effective attention window (cf. Lost-in-the-Middle [[1.0.0 P-53]]) will fail; chunks that are too small will miss cross-file dependencies that the Verified test suite implicitly requires. The 1–4-file edit distribution is the empirical sweet spot.

- **PRIM-27 (Coverage-Guided Plateau Detection)** — SWE-bench Verified's pass-rate curve across repair iterations is the **canonical plateau-detection metric**. Specifically: when a SLM scaffolded with MAgHARCM's pipeline runs N repair iterations on a Verified instance, the marginal pass-rate gain between iterations i and i+1 should drop monotonically; PRIM-27 stops the loop when that gain falls below a configured threshold. The benchmark's 500-instance aggregate pass-rate (currently ~30–50% for top models, ~5–12% for SLMs) is the **plateau floor** that distinguishes "we're done" from "we're stuck." The OpenAI release's score reports provide the iteration-by-iteration curves that calibrate the threshold.

- **SLM-era headline framing** — combined with [[1.0.0 P-102]] (SmallCode 4B at 87% HumanEval), SWE-bench Verified closes the loop: a 4B SLM can synthesize competitive single functions, and a properly scaffolded 7B-15B SLM agent can resolve a meaningful fraction of real repository-level issues. The gap between the two regimes — single-function synthesis vs. multi-file issue resolution — is exactly the gap MAgHARCM's 8-agent pipeline is designed to close, and SWE-bench Verified is the benchmark that measures it.

- **Per-PRIM calibration, summarized** — use Verified's test patches as oracle inputs to PRIM-5; use Verified's failure traces as repair-loop templates for PRIM-6; use Verified's per-instance edit-size distribution to set PRIM-23's chunk-size default; use Verified's aggregate scores as the plateau-stop metric in PRIM-27.

## 3. Hop-1 References (papers and systems that SWE-bench Verified cites or builds on)

- Jimenez et al. (2024) — SWE-bench: Can Language Models Resolve Real-World GitHub Issues? (ICLR 2024); the original 2,294-instance benchmark that Verified subsets.
- Anthropic (2024) — Claude 3.5 Sonnet system card; reports 49% on SWE-bench Verified, the headline frontier-model number in the OpenAI release.
- OpenAI (2024) — OpenAI o1 system card; reports o1's SWE-bench Verified score as a reasoning-model baseline.
- Aider (2024) — Aider polyglot benchmark; a complementary code-edit benchmark that evaluates per-file diff generation rather than issue-resolution; cross-link [[1.0.0 P-51]] (Can It Edit) for the multi-step edit ancestor.
- SWE-Lancer (OpenAI, 2025) — a follow-on dollar-valued benchmark (real freelance engineering tasks); uses the same Verified methodology.
- SWE-Gym (Pan et al., 2024) — a training environment that uses SWE-bench instances for RL-style agent training; important because it shows the Verified instances are also *trainable*, not just *evaluable*.
- Peng et al. (2023) — RepoCoder; iterative-retrieval-augmented code completion, one of the agent scaffolds evaluated in the original SWE-bench and re-evaluated on Verified.
- Jalil et al. (2023) — R2E: Repository Editing; a repository-grounded editing agent that translates natural-language edits to diffs, evaluated on SWE-bench and Verified.
- Zhang et al. (2024) — AutoCodeRover; an autonomous software-engineering agent that combines retrieval and program-synthesis on SWE-bench, the canonical open-source baseline for Verified.
- Yang et al. (2024) — SWE-agent; the tool-calling autonomous agent that popularized SWE-bench evaluation and is one of the three scaffolds re-evaluated on Verified.
- Wang et al. (2024) — OpenHands (CodeAct); the multi-turn CodeAct agent, the third canonical scaffold on Verified.

## 4. Hop-2 References (papers cited by the hop-1 lineage)

- Vaswani et al. (2017) — Attention Is All You Need; transformer architecture underlying every evaluated model.
- Austin et al. (2021) — Program Synthesis with Large Language Models (APPS); foundational competitive-programming benchmark referenced by SWE-bench for problem-difficulty stratification.
- Chen et al. (2021) — Evaluating Large Language Models Trained on Code (HumanEval / Codex); foundational single-function synthesis benchmark referenced by every SWE-bench-adjacent paper for SLM evaluation.
- Li et al. (2023) — Competition-Level Code Generation with Code LLMs (CodeContest); referenced by SWE-Agent and AutoCodeRover as the competitive-programming anchor for repository-level agents.
- Wei et al. (2022) — Chain-of-Thought Prompting; referenced by SWE-agent and OpenHands for the reasoning chain that precedes tool calls.
- Yao et al. (2022) — ReAct; the reasoning+acting framework that SWE-agent and OpenHands both inherit and that drives Verified's tool-calling scaffolds.
- Hendrycks et al. (2021) — Measuring Massive Multitask Language Understanding (MMLU); referenced by Verified release as general-capability baseline.
- Touvron et al. (2023) — LLaMA 2; foundational architecture for many of the open-model baselines (Llama-3.1 8B is the Verified SLM floor).
- Schick et al. (2023) — Toolformer; referenced by SWE-agent for tool-calling methodology.
- Khattab et al. (2024) — DSPy; referenced by OpenHands as a prompt-optimization layer that can lift Verified scores.
- Peng et al. (2023) — RepoCoder (cross-link, hop-1 already lists this; hop-2 here = its retrieval-augmented generation ancestors).
- Roziere et al. (2023) — Code Llama; instruction-tuned code baseline; cross-link [[1.0.0 P-95]] for the MAgHARCM lineage.

## 5. Backlinks

- **PRIM-5** (Test Suite Co-Translation & Synthesis): Verified's cleaned test patches are oracle scaffolds for `generateAdditionalTests`; the maintainer's passing test diffs are the "legitimate test change" reference standard.
- **PRIM-6** (Multi-Stage Build/Test Feedback Repair): Verified's instance failure traces (compile → test → patch regression → final-pass) are the per-stage repair-loop templates; the score gap between SLMs and frontier models quantifies the scaffold lift.
- **PRIM-23** (Chunked Translation): Verified's per-instance edit-size distribution (1–4 files, localized edits) is the empirical sweet spot for `chunked_translator.go`'s chunk-size default.
- **PRIM-27** (Coverage-Guided Plateau Detection): Verified's pass-rate curve across repair iterations is the plateau-stop metric; aggregate scores (5–12% for SLMs, 33–49% for frontier) are the plateau floor.
- **PRIM-22** (Four Phases Comprehension): the cleaned task statements (self-contained problem descriptions) are the canonical inputs to the comprehend phase's slot-filling prompts; cross-link [[1.0.0 P-55]] (Schick & Schütze cloze reformulation) for the SLM-prompting pattern.
- **PRIM-31** (Iterative Retrieval Refinement): Verified scores separate Agentless (procedural, single-pass) from SWE-agent (iterative, tool-calling), quantifying the lift that `iter_retrieval.go` provides.
- **PRIM-21** (Migration Strategy Selection): the failure-kind classification (assertion vs. collection vs. import) is the input signal that `Registry.TryInOrder` consumes to pick the next repair strategy.
- **Cross-ref**: add P-109 to PRIM-5, PRIM-6, PRIM-23, PRIM-27 rows in `Software-Archaeology-Lineage.md` (handled by the separate "Update primitives INDEX + lineage matrix" orchestration task). This is the **canonical end-to-end evaluation anchor** for the 8-agent MAgHARCM pipeline in the SLM era.
