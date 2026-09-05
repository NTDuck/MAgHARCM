---
title: "P-84 — Muennighoff et al. 2025 — s1: Simple Test-Time Scaling"
backlink: "[[1.0.0 P-84]]"
tags: [paper, test-time-scaling, reasoning, slm, budget-forcing, wait-token, [[1.0.0 PRIM-3]], [[1.0.0 PRIM-22]], [[1.0.0 PRIM-25]], hop-1]
---

# [[1.0.0 P-84 — Muennighoff et al. — s1 Simple Test-Time Scaling]]

## Citation

Muennighoff, N., Yang, Z., Shi, W., Li, X., Zhang, L.-F., Wei, F., Cheng, H., Liu, H., Gao, L., Chen, Y., et al. (2025). *s1: Simple Test-Time Scaling*. arXiv:2501.19393 (January 2025). Project page: https://github.com/simplescaling/s1.

## Summary

The s1 paper asks the narrow but consequential question: *can test-time compute scaling for reasoning be achieved without reinforcement learning, by fine-tuning a small set of high-quality reasoning traces on an already-strong base model?* The authors' answer is yes. Fine-tuning Qwen2.5-32B-Instruct on **s1K** — a curated dataset of 1,000 questions paired with reasoning traces selected on three criteria (quality, difficulty, diversity) — produces **s1-32B**, which outperforms OpenAI o1-preview on competition-math benchmarks (MATH, AIME24) by up to 27%. The technique, **budget forcing**, controls test-time compute at inference time by either terminating generation early (injecting the end-of-thinking token) or extending it (suppressing the stop token and appending the literal word *"Wait"* to the trace, prompting the model to keep reasoning).

Three findings matter most for MAgHARCM:

1. **1,000 high-quality samples beat 100,000 mediocre ones.** The s1K dataset's three-criterion curation (quality + difficulty + diversity) is the existence proof that targeted data curation — not data scale — unlocks reasoning. For MAgHARCM, this validates the principle that carefully-architected comprehension prompts (PRIM-22) and SOP-anchored artifact schemas (PRIM-24) extract more value from a 32B SLM than naive few-shotting.

2. **Budget forcing is the operational primitive.** The literal "Wait" injection — appending a single token to force continued deliberation — is the cheapest possible mechanism for allocating extra test-time compute. MAgHARCM can adopt this directly in the optional-checks agent (PRIM-25): when confidence is below a threshold, inject "Wait — re-check the boundary condition" and let the model re-deliberate on the same prompt.

3. **Test-time scaling has a clear knee.** s1-32B shows monotone improvement on math reasoning as budget grows up to ~4,000 "thinking" tokens, then plateaus. This gives MAgHARCM a calibration point: rather than spending compute indiscriminately, the verdict panel (PRIM-7) and the role-flip gate (PRIM-25) can budget deliberation per-query rather than per-batch.

## Method

- **Base model**: Qwen2.5-32B-Instruct (chosen because pre-existing reasoning is already strong; s1 simply amplifies latent capability).
- **Data**: s1K = 1,000 carefully selected reasoning traces. Selection criteria: (a) quality (verified correct), (b) difficulty (fails greedy decoding from Qwen2.5-32B-Instruct baseline), (c) diversity (broad task taxonomy). Sample budgets studied from 100 to 100,000 examples; performance saturates near 1,000.
- **Fine-tuning**: standard SFT on s1K, ~26 minutes on 16×H100. No RL, no reward model.
- **Inference (budget forcing)**: a wrapper monitors generated tokens. When the model produces an end-of-thinking token earlier than the configured budget, the wrapper *appends it* to force termination; when the model produces it later than the budget, the wrapper *suppresses it and appends "Wait"* to force continuation. This is a single-token surgical intervention — far cheaper than RL rollouts or MCTS.
- **Evaluation**: AIME24, MATH (500-problem subset), GPQA Diamond, code benchmarks. s1-32B beats o1-preview by 27 percentage points on AIME24 under matched token budgets.

## Findings Relevant to MAgHARCM

- **Curated 1K > raw 100K** is the operational justification for MAgHARCM's prompt-engineering investment. The 31 primitives and the SOP-anchored role schema (PRIM-24) are the "curation" layer for SLM reasoning; s1 demonstrates the leverage is real.
- **"Wait" injection is the cheapest test-time compute allocator.** MAgHARCM's optional-checks agent (PRIM-25) and verdict panel (PRIM-7) can adopt a literal "Wait — re-check…" continuation token when confidence is below threshold. Cost: one extra token of input + one short re-decoding pass. Benefit: demonstrable accuracy gain on hard reasoning sub-tasks.
- **Reasoning budget is a knob, not a property.** MAgHARCM can configure reasoning depth per-task: light for comprehension (PRIM-22), heavy for adversarial-review verdicts (PRIM-13, PRIM-25). s1's monotone curve through ~4K tokens gives a calibration target.
- **s1-32B is in the SLM range MAgHARCM targets.** 32B is at the upper edge of the 4B-30B+ fleet (Qwen2.5-Coder [[P-21]], Phi-3 [[P-54]], StarCoder2 [[P-22]]). s1 demonstrates that a 32B model with the right training data can match a much larger closed model — reinforcing P-81 (Gorilla)'s SLM-tool-use thesis and the P-82 SLM-migration thesis.

## How MAgHARCM Uses It

- **PRIM-22 (Four Phases Comprehension)**: during the *explanation* and *search* phases, inject a "Wait — re-derive from first principles" continuation token when the comprehension trace is shallow. This is the budget-forcing pattern applied at the comprehension layer.
- **PRIM-25 (Communicative De-Hallucination Gate)**: the role-flip reviewer should issue a second-pass prompt with literal "Wait" prepended to the system message when the first-pass verdict disagrees with the verdict panel. Cost: negligible (single token). Benefit: catches the cases where the first-pass reviewer capitulated.
- **PRIM-3 (Target Skeleton-First Generation)**: s1's "quality + difficulty + diversity" curation matches MAgHARCM's skeleton-first approach — provide the hardest-pinned-down interface (skeleton) and let the SLM fill in only the body, then budget-force the body for clarity.
- **Future work (not implemented)**: per-query budget allocation in `verdict_panel.go`. Track disagreement rate over a rolling window; increase budget for high-disagreement queries. The DSC (difficulty-adaptive self-consistency) finding from P-83 composes naturally with s1's budget-forcing.

## References

### Hop-1 (papers cited by s1 or directly related)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. Foundational CoT; s1K's reasoning traces are CoT-style.
- Snell, C., Lee, J., Xu, K., & Kumar, A. (2024). *Scaling LLM Test-Time Compute Optimally Can be More Effective than Scaling Model Parameters*. arXiv:2408.03314. The test-time-compute-scaling precursor s1 builds on.
- Zhou, L., Schellaert, W., Martínez-Plumed, F., et al. (2024). *Larger and More Instructable Language Models Become Less Reliable*. Nature 634. Cited by s1 for the inverse-scaling concerns motivating small-data curation.
- Lightman, H. et al. (2023). *Let's Verify Step by Step* (PRM800K). The RL-reward-model baseline s1 explicitly avoids.
- OpenAI (2024). *Learning to Reason with LLMs* (o1 system card). The closed-source baseline s1-32B beats.
- Liu, A. et al. (2024). *Lost in the Middle* [[P-53]]. Cited for evidence that structured prompts outperform free-form.

### Hop-2 (foundational anchors referenced)
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners* (GPT-3). arXiv:2005.14165. The in-context-learning foundation.
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762. Transformer architecture.
- Rafailov, R. et al. (2023). *Direct Preference Optimization (DPO)*. arXiv:2305.18290. The alternative alignment signal to RLHF/RM.
- Zhou, C. et al. (2023). *LIMA: Less Is More for Alignment*. arXiv:2305.11206. The "1K high-quality beats 100K random" hypothesis s1 inherits.
- Team Qwen (2024). *Qwen2.5 Technical Report*. arXiv:2412.15115. The base model for s1-32B.

### MAgHARCM lineage cross-refs
- [[1.0.0 P-21]] — Qwen2.5-Coder 2024 (s1-32B's base model is from the same Qwen family).
- [[1.0.0 P-52]] — Self-Consistency (Wang 2023). The sampling-ensemble ancestor to s1's test-time allocation.
- [[1.0.0 P-53]] — Lost in the Middle (Liu 2023). Justifies s1's structured trace prompt.
- [[1.0.0 P-81]] — Gorilla (Patil 2023). The SLM-tool-use thesis s1 extends to reasoning.
- [[1.0.0 P-83]] — Code-specialised self-consistency. Composes with s1 via DSC (difficulty-adaptive sampling).
- [[1.0.0 P-55]] — Schick & Schütze SLM Few-Shot. The earlier "1K samples can be enough" anchor.

## Backlinks

[[1.0.0 PRIM-3]], [[1.0.0 PRIM-7]], [[1.0.0 PRIM-22]], [[1.0.0 PRIM-25]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-84 is the **test-time-scaling SLM anchor** for MAgHARCM: it demonstrates that budget forcing + curated traces unlocks o1-class reasoning on a 32B open model, with direct cost (one "Wait" token) that the optional-checks agent and verdict panel can adopt.
