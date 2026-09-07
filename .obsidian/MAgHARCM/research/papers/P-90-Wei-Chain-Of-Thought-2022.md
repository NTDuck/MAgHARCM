---
title: "P-90 — Wei et al. 2022 — Chain-of-Thought Prompting Elicits Reasoning in Large Language Models"
backlink: "[[1.0.0 P-90]]"
aliases:
  - "1.0.0 P-90"
  - "P-90"
  - "P-90-Wei-Chain-Of-Thought-2022"
  - "P-90-Wei-Chain-Of-Thought-2022"
  - "Wei-Chain-Of-Thought-2022"
tags: [paper, chain-of-thought, reasoning, prompting, few-shot, slm, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-25]], hop-1]
---

# [[1.0.0 P-90 — Wei et al. — Chain-of-Thought Prompting]]

## Citation

Wei, J., Wang, X., Schuurmans, D., Bosma, M., Ichter, B., Xia, F., Chi, E. H., Le, Q. V., & Zhou, D. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. Advances in Neural Information Processing Systems 35 (NeurIPS 2022), pp. 24824–24837. arXiv:2201.11903 (28 Jan 2022 v1; 10 Oct 2022 v6 with appendix).

## Summary

The Wei et al. (2022) paper is the foundational chain-of-thought (CoT) prompting paper. It asks the deceptively simple question: *can a few hand-written examples that include intermediate reasoning steps unlock problem-solving in sufficiently large language models?* The answer is yes — and the threshold is sharp.

The authors show that standard few-shot prompting ("input → answer") plateaus on multi-step reasoning tasks (arithmetic, commonsense, symbolic) for **PaLM 540B**, but a few-shot prompt where each exemplar is augmented with a chain of intermediate reasoning steps ("input → reasoning → answer") produces dramatic improvements on the same model. On GSM8K (grade-school math), PaLM 540B-CoT achieves **56.9% accuracy** versus 17.9% for standard few-shot — the largest single-prompting improvement demonstrated on that benchmark at the time. Crucially, the **gains are model-scale-dependent**: PaLM 540B shows large gains, while the 62B and 8B variants show only marginal or no improvement. The paper's most-cited claim is that *"chain-of-thought prompting is an emergent property of model scale."*

The paper introduces three task families:
1. **Arithmetic reasoning** — GSM8K, SVAMP, ASDiv, AQuA, MAWPS.
2. **Commonsense reasoning** — CSQA, StrategyQA, Date Understanding, Sports Understanding.
3. **Symbolic reasoning** — Last Letter Concatenation, Coin Flip.

The technique requires **no fine-tuning** — only an in-context exemplar with intermediate steps. The exemplars need not be task-matched (Robustness Experiment §3.3): even CoT exemplars from a different task class can elicit reasoning.

## Method

- **Prompt format**: a small set of `Q → [reasoning chain] → A` exemplars (typically 6–8), followed by a new `Q →` that the model completes.
- **No fine-tuning**: pure inference-time prompt composition.
- **Models studied**: LaMDA 137B, PaLM 8B / 62B / 540B. The 540B setting is the headline configuration.
- **Evaluation**: greedy decoding, exact-match accuracy, comparison to standard few-shot baselines.
- **Robustness experiments**: (a) vary the order of exemplars, (b) replace math exemplars with non-math CoT exemplars — both confirm CoT is robust to exemplar choice.

## Findings Relevant to MAgHARCM

- **CoT does NOT work uniformly on SLMs.** Wei et al.'s headline finding — "emerges with scale" — implies that CoT prompts may waste tokens on 1.5B-7B models. MAgHARCM's `[[1.0.0 PRIM-22]]` Four Phases Comprehension must therefore include an *adaptive CoT toggle*: apply CoT only on SLMs ≥ ~30B (Qwen2.5-Coder-32B [[1.0.0 P-21]], s1-32B [[1.0.0 P-84]]); use direct prompts on smaller SLMs.
- **CoT is the substrate for s1's reasoning traces.** The s1K dataset (1,000 curated reasoning traces, [[1.0.0 P-84]]) is structurally CoT — input → many-step reasoning → answer. P-90 is the upstream foundation paper that defines the format s1 inherits.
- **Robustness to exemplar choice matters for prompt engineering.** P-90 demonstrates that the *exemplar task* need not match the *target task*. This validates MAgHARCM's practice of using shared CoT exemplars across comprehension and translation prompts (cf. `internal/agents/comprehension.go`).
- **CoT + Self-Consistency is a known stack.** Wang et al. 2023 [[1.0.0 P-52]] builds on P-90 by sampling multiple CoT traces and majority-voting. MAgHARCM's verdict panel (`internal/agents/verdict_panel.go`) implements a related pattern.
- **No fine-tuning requirement is operationally important.** CoT is a *prompt-engineering* technique — it requires zero gradient updates and zero new training data. This is the cheapest possible reasoning improvement; it should be the default for MAgHARCM's `[[1.0.0 PRIM-25]]` Role-Flip reviewer when the underlying SLM is large enough.

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-22]]` Four Phases of Comprehension**: the *explanation* and *search* phases embed 4-shot CoT exemplars from prior comprehension tasks. Add a runtime gate that disables CoT on sub-7B SLMs (where it can hurt) and enables it on ≥30B models.
- **`[[1.0.0 PRIM-25]]` Communicative De-Hallucination Gate**: the role-flip reviewer's adversarial prompt includes a CoT exemplar where the reviewer is asked to *show its work* when finding a defect ("identify the missing semicolon → walk through the AST → cite the offending line"). This converts the reviewer from a one-shot answerer into a chain-reasoning adversary.
- **`[[1.0.0 PRIM-3]]` Target Skeleton-First Generation**: when generating long-form code skeletons, prepend 1-shot CoT exemplars that walk through "skeleton → fill body" — empirically improves skeleton coherence on 32B SLMs (consistent with the s1-32B finding [[1.0.0 P-84]] that inherits from P-90).
- **Future work (not implemented)**: per-model-size CoT tuning. Build a calibration table mapping model size + benchmark → optimal CoT exemplar count. Start with {1.5B: 0-shot, 7B: 2-shot, 14B: 4-shot, 32B+: 8-shot}.

## References

### Hop-1 (papers that build directly on P-90)
- Wang, X. et al. (2023). *Self-Consistency Improves Chain of Thought Reasoning in Language Models*. ICLR 2023. arXiv:2203.11171. See `[[1.0.0 P-52]]` for the sampling-based aggregation of CoT traces.
- Muennighoff, N. et al. (2025). *s1: Simple Test-Time Scaling*. arXiv:2501.19393. See `[[1.0.0 P-84]]` — s1K traces are CoT-augmented.
- Kojima, T. et al. (2022). *Large Language Models are Zero-Shot Reasoners*. NeurIPS 2022. arXiv:2205.11916. Shows that *"Let's think step by step"* unlocks CoT **without exemplars**.
- Zhou, D. et al. (2023). *Least-to-Most Prompting Enables Complex Reasoning in Large Language Models*. ICLR 2023. arXiv:2205.10625. Decomposes a hard problem into ordered subproblems; relevant to `[[1.0.0 PRIM-22]]` planning.
- Khot, T. et al (2022). *Decomposed Prompting: A Modular Approach for Solving Complex Tasks*. arXiv:2210.02406. See `[[1.0.0 P-62]]` — modular decomposition is a generalization of CoT.
- Suzgun, M. et al. (2022). *Challenging BIG-Bench Tasks and Whether Chain-of-Thought Can Solve Them*. arXiv:2210.09261. Systematic study of CoT across 23 BIG-Bench tasks.

### Hop-2 (foundational anchors referenced)
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners* (GPT-3). arXiv:2005.14165. The in-context-learning foundation CoT builds on.
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762. Transformer architecture.
- Cobbe, K. et al. (2021). *Training Verifiers to Solve Math Word Problems* (GSM8K). arXiv:2110.14168. The benchmark P-90 uses.
- Hendrycks, D. et al. (2021). *Measuring Mathematical Problem Solving With the MATH Dataset*. arXiv:2103.03874. The competition-math benchmark.
- Chowdhery, A. et al. (2022). *PaLM: Scaling Language Modeling with Pathways*. arXiv:2204.02311. The base model P-90 studies.

### MAgHARCM lineage cross-refs
- `[[1.0.0 P-21]]` — Qwen2.5-Coder (the 32B SLM where CoT unlocks reasoning per P-90's scale dependence).
- `[[1.0.0 P-52]]` — Self-Consistency (sampling-aggregation layer on top of CoT).
- `[[1.0.0 P-62]]` — Decomposed Prompting (modular generalization of CoT, used in PRIM-22/23/31).
- `[[1.0.0 P-84]]` — s1 (CoT-style reasoning traces + budget forcing).
- `[[1.0.0 P-83]]` — Code-specialised Self-Consistency (CoT-augmented code sampling).

## Backlinks

`[[1.0.0 PRIM-3]]`, `[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-25]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-90 is the **foundational CoT anchor** for MAgHARCM's reasoning layer. Every later reasoning paper — Self-Consistency (P-52), s1 (P-84), Decomposed Prompting (P-62) — inherits from P-90's chain-of-thought format. Cite P-90 whenever a prompt template embeds an explicit `reasoning → answer` exemplar.
