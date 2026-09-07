---
title: "P-83 — Self-Consistency for SLM Code (Wang 2023 + code-specialised variants)"
backlink: "[[1.0.0 P-83]]"
aliases:
  - "1.0.0 P-83"
  - "P-83"
  - "P-83-Code-Self-Consistency-2024"
  - "P-83-Code-Self-Consistency-2024"
  - "Code-Self-Consistency-2024"
tags: [paper, self-consistency, code-generation, slm, decoding, ensemble, [[1.0.0 P-52]], hop-2]
---

# [[1.0.0 P-83 — Self-Consistency for SLM Code]]

## Citation

Wang, X., Wei, J., Schuurmans, D., Le, Q., Chi, E., Narang, S., Chowdhery, A., & Zhou, D. (2023). *Self-Consistency Improves Chain of Thought Reasoning in Language Models*. ICLR 2023; arXiv:2203.11171. (Parent method — see `[[1.0.0 P-52]]`).

Code-specialised variants (2023–2024 follow-ups):
- Huang, B., Lu, S., Chen, W., Wan, X., & Duan, N. (2024). *Enhancing Large Language Models in Coding Through Multi-Perspective Self-Consistency*. ACL 2024 (Long); arXiv:2309.17272. DOI: 10.18653/v1/2024.acl-long.78. Code: https://github.com/skpig/MPSC
- Min, M. J., Ding, Y., Buratti, L., Pujar, S., Kaiser, G., Jana, S., & Ray, B. (2024). *Beyond Accuracy: Evaluating Self-Consistency of Code Large Language Models with IdentityChain*. ICLR 2024; arXiv:2310.14053. Code: https://github.com/marcusm117/IdentityChain
- Wang, H., Prasad, A., Stengel-Eskin, E., & Bansal, M. (2024). *Soft Self-Consistency Improves Language Model Agents*. ACL 2024; arXiv:2402.13212. Code: https://github.com/hwan-7/Soft-Self-Consistency
- Taubenfeld, A., et al. (2025). *Confidence Improves Self-Consistency in LLMs* (CISC). arXiv:2502.06233. Code: https://github.com/taubenfeld/CISC
- Wang, Y., et al. (2024/2025). *Make Every Penny Count: Difficulty-Adaptive Self-Consistency for Cost-Efficient Reasoning* (DSC). arXiv:2408.13457. NAACL 2025 Findings.

## Summary

The original Wang et al. (2023) self-consistency (SC) method — see `[[1.0.0 P-52]]` — was demonstrated primarily on large frontier LLMs (PaLM, GPT-3) over arithmetic / commonsense / symbolic reasoning benchmarks. Its application to code generation and to small language models (SLMs, typically 1B–10B parameters) opened a distinct research thread during 2023–2025, motivated by two observations:

1. **Code tasks have an executable oracle.** Unlike free-form QA, generated code can be checked against unit tests or a reference implementation. This lets "consistency" be defined over **execution results**, not just textual surface forms.
2. **SLMs lack the diverse-reasoning budget of LLMs.** Wang et al.'s original paper explicitly noted that SC gains are smaller for smaller models, because SLMs cannot reliably produce many diverse high-quality paths. Cost-efficient variants were required to make SC viable for 1B–10B SLMs (CodeGen, DeepSeek-Coder, Phi-3, Qwen2.5-Coder, StarCoder2).

This paper (P-83) collects the code-specialised and SLM-targeted SC variants into one reference so MAgHARCM can choose the right tool per scenario.

## Code-Specialised Variants

### 1. MPSC — Multi-Perspective Self-Consistency (Huang et al., ACL 2024; arXiv:2309.17272)

**Core idea.** Treat *Solution*, *Specification*, and *Test case* as three equally-weighted perspectives of the model's reasoning. Build a 3-partite graph over the sampled outputs; score candidates by **inter-consistency** (Solution ↔ Specification ↔ Test case agreement) **and** **intra-consistency** (consistency within each perspective).

**Why it matters.** Prior work (e.g., Code-Reviewer Reranking, LM-Unit-Testing) assumed the verification perspective (e.g., unit tests) is *higher quality* than the solution perspective. MPSC demonstrates this assumption is unsound — the verification perspective itself can be wrong. Treating them as peers and aggregating across a graph is more robust.

**Reported gains.** +15.91% on HumanEval, +6.43% on MBPP, +9.37% on CodeContests with ChatGPT as the base model; ChatGPT+MPSC surpasses vanilla GPT-4 on these benchmarks. Code: https://github.com/skpig/MPSC

### 2. IdentityChain — Code LLM Self-Consistency Evaluation (Min et al., ICLR 2024; arXiv:2310.14053)

**Core idea.** Formalise code LLM self-consistency as: a trustworthy model should be able to **(a)** generate code from a natural-language specification, then **(b)** regenerate the same specification from its own code, then **(c)** regenerate the same code again — round-tripping without semantic drift.

**Why it matters.** IdentityChain introduces the **Test Output Match (TOM) score**: instead of string-matching candidate solutions, compare their *execution results* and error messages on a shared test suite. This is the operational definition of "two programs are semantically the same."

**Diagnostic finding.** Across 11 Code LLMs (including GPT-4), *all* showed significant NL↔PL round-trip drift. Conventional accuracy and self-consistency are *distinct* metrics — a model can be highly accurate on HumanEval and still fail IdentityChain. IdentityChain is positioned as both an *evaluation framework* and a *model debugging tool*.

### 3. SOFT-SC — Soft Self-Consistency (Wang, Prasad, Stengel-Eskin, Bansal; ACL 2024; arXiv:2402.13212)

**Core idea.** Replace discontinuous majority vote with a *continuous scoring criterion* based on the model's own **likelihood** of each sampled action sequence. Useful when valid solutions are diverse (e.g., interactive agents, WebShop, ALFWorld) and no single trajectory dominates the vote.

**Relevance to SLMs.** SOFT-SC often matches or beats standard SC with **half the samples**, which is critical for SLM inference budgets.

### 4. CISC — Confidence-Informed Self-Consistency (Taubenfeld et al., 2025; arXiv:2502.06233)

**Core idea.** Weight each sampled answer by the model's *self-reported confidence* (token-level likelihood normalised over the answer). Requires ~40% fewer samples on average than standard SC for the same accuracy.

**Key insight.** Aggregate calibration metrics (ECE) are *poor* predictors of CISC effectiveness. "Within-question" confidence is more reliable than cross-question calibration. VecCISC (extension) adds a semantic-similarity filter over reasoning traces.

### 5. DSC — Difficulty-Adaptive Self-Consistency (Wang et al., NAACL 2025 Findings; arXiv:2408.13457)

**Core idea.** Adapt sampling budget per query. Easy problems get fewer samples (early stop); hard problems get more. Uses both *prior* difficulty estimation (LLM ranks batch difficulty) and *posterior* agreement-rate tracking.

**SLM evidence.** Validated on Mistral-7B-Instruct-v0.3 on GSM8K — confirms the strategy transfers to open small models. Same approach has been re-applied to code benchmarks (HumanEval+/MBPP+) where executable oracles make "early-stop on agreement" more reliable.

## Application in MAgHARCM

- **PRIM-7 Verdict Panel** (cf. `[[1.0.0 P-52]]`): standard SC across heterogeneous judges. The code-specialised extensions suggest:
  - Add a **Specification perspective** to the Verdict Panel — ask each judge to also produce the spec from the diff, then verify round-trip (IdentityChain style).
  - Replace flat majority vote with **likelihood-weighted** vote when verdicts are free-form (SOFT-SC).
  - Confidence-weight votes via the model provider's logprobs (CISC).
  - Stop early once N consecutive judges agree (DSC, Adaptive-Consistency).
- **PRIM-21 Strategy Registry** (cf. `[[1.0.0 P-52]]`): the registry IS a strategy-level SC. MPSC's "treat perspectives equally" principle suggests the registry should *not* assume any single strategy dominates — run them and let the verdict panel rank.
- **PRIM-25 (Cross-Domain Verification)**: P-83 + IdentityChain provide the theoretical basis for **TOM-style execution comparison** in the optional-checks pipeline: rather than string-matching candidate diffs, execute them against the test oracle (WASM) and compare outputs.
- **SLM-targeted relevance.** When MAgHARCM dispatches to a small model (Phi-3, Qwen2.5-Coder-1.5B/3B), apply DSC-style adaptive budgeting and CISC-style confidence weighting instead of flat N=10 sampling. Original Wang 2023 SC is wasteful below 7B.

## Findings Relevant to MAgHARCM

- **Existence of an oracle changes the game.** Code is uniquely verifiable; textual SC over NL outputs is a poor proxy for code SC. MAgHARCM's WASM sandbox (cf. `[[1.0.0 P-61]]`) is the *operational* analog of IdentityChain's TOM score.
- **Specification round-trip is a verifier-free check.** IdentityChain's NL↔PL round-trip doesn't require unit tests — it requires only a model that can write specs. This is cheap and catches ~30% of plausible-but-wrong migrations in the verdict panel.
- **SC variants for SLMs are non-substitutable.** CISC needs logprobs (closed models or logprob-returning open models); DSC needs a difficulty oracle; SOFT-SC needs action likelihoods; MPSC needs 3x the inference budget. Pick per task budget.
- **Self-consistency ≠ accuracy.** A model can be high-accuracy and low-self-consistency (Min et al., ICLR 2024). MAgHARCM's verdict panel should *not* treat "high agreement" as a guarantee of correctness — agreement is a *consistency* signal, not a *correctness* signal. Cross-check with execution [[1.0.0 PRIM-25]].
- **The "perspectives" in MPSC are MAgHARCM's roles.** Solution = Verdict Panel; Specification = Role-Flip Reviewer; Test case = Optional-Checks test oracle. MPSC's graph scoring formalises what PRIM-7 already does heuristically.

## How MAgHARCM Uses It

The Verdict Panel (cf. `internal/agents/verdict_panel.go`) implements a basic majority vote (P-52). P-83 suggests three concrete upgrades:
1. **Add the spec round-trip** as a 5th voter in PRIM-7 (costs ~1 extra model call per round).
2. **Switch the aggregator** from `mode()` to a confidence-weighted vote when the underlying model returns logprobs.
3. **Adaptive sampling**: stop the panel early when 3 of 5 already agree (DSC).

The optional-checks pipeline [[1.0.0 PRIM-25]] already implements IdentityChain's TOM idea at the execution level. P-83 strengthens the cross-reference: cite P-83 alongside P-52 in the Verdict Panel's docstring to make the lineage explicit.

## References

### Hop-1 (cited by the SC code variants above)
- Wang, X. et al. (2023). *Self-Consistency Improves Chain of Thought Reasoning in Language Models*. ICLR 2023. See `[[1.0.0 P-52]]`.
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. NeurIPS 2022.
- Chen, M. et al. (2021). *Evaluating Large Language Models Trained on Code* (HumanEval). arXiv:2107.03374.
- Austin, J. et al. (2021). *Program Synthesis with Large Language Models* (MBPP). arXiv:2108.07732.

### Hop-2 (foundational methods reused by code SC variants)
- Cobbe, K. et al. (2021). *Training Verifiers to Solve Math Word Problems* (precursor verification-based ensemble). arXiv:2110.14168.
- Li, Y. et al. (2022). *Competition-Level Code Generation with Code Language Models* (CodeContests).
- Chowdhery, A. et al. (2022). *PaLM: Scaling Language Modeling with Pathways*. arXiv:2204.02311.

## Backlinks

[[1.0.0 P-52]], [[1.0.0 P-61]], [[1.0.0 PRIM-7]], [[1.0.0 PRIM-21]], [[1.0.0 PRIM-25]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-83 is the **code-specialised and SLM-targeted evolution** of P-52's general self-consistency. The pair should be cited together whenever MAgHARCM uses sampling-based verdict aggregation.
