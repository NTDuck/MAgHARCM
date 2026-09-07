---
title: "P-97 — Welleck, Lu & Raffel 2024 — Generating Sequences by Learning to Self-Correct"
backlink: "[[1.0.0 P-97]]"
tags: [paper, self-correction, slm, fine-tuning, inference-cost, verifier-search, code-generation, [[1.0.0 PRIM-7]], [[1.0.0 PRIM-25]], hop-2]
---

# [[1.0.0 P-97 — Welleck et al. — Self-Correct]]

## Citation

Welleck, S., Lu, X., & Raffel, M. (2024). *Generating Sequences by Learning to Self-Correct*. Proceedings of the 41st International Conference on Machine Learning (ICML 2024). arXiv:2211.00053 (31 Oct 2022 v1; 7 Nov 2023 v3 with appendix). Carnegie Mellon University.

## Summary

The Welleck et al. paper reframes self-correction: instead of *hoping* a base LLM will catch and fix its own errors via an inference-time prompt ("Let me re-check that..."), the paper **trains the model to self-correct** as a first-class capability. The headline contribution is a fine-tuning procedure that teaches the model to iteratively refine its own outputs, paired-edit by paired-edit, so that at inference time the base model itself — not an external verifier or PRM — is the corrector.

Two distinct settings are studied:

1. **Training-time self-correction** (the paper's contribution): construct paired training data `(y_t, y_{t+1})` of "incorrect output → corrected output", then fine-tune the base model with an iterative self-correction loss that conditions on the model's own previous prediction. The resulting model is *trained to refine*.
2. **Inference-time self-correction via prompting** (a baseline the paper ablates): prompt the base model at decode time to critique and revise its own answer. The paper shows this prompting-only approach *fails to actually improve outputs on most tasks* — the base model rarely produces a strictly better revision. The training-time intervention is what makes self-correction work.

Key empirical claims:

1. **Trained self-correction matches or beats a strong verifier-search baseline on code and translation.** On code generation (HumanEval / MBPP) and machine translation (WMT En-De, En-Fr), a self-correcting 7B-13B SLM reaches accuracy on par with — and on some tasks above — a *verifier-best-of-N* pipeline that runs an external PRM over many samples per query.
2. **Verifier-free inference.** Because the corrector lives inside the base model, no separate verifier model is run at decode time. The compute saving is substantial: a PRM best-of-N with N=64 samples is ~64× the per-step cost of one self-correcting decode; the paper shows comparable accuracy with a handful of self-correction rounds (typically 1–4).
3. **Per-step conditioning matters.** Naively fine-tuning on `(y_t, y_{t+1})` pairs and asking the model to output `y_{t+1}` given `y_t` does not work — the model still degrades. The paper's **iterative training** explicitly trains on multiple refinement rounds and uses an off-policy-correction objective to handle distribution shift from the base model.
4. **Task-agnostic recipe.** The same recipe applies to code and translation; the gains are largest on tasks where inference-time verifier-search is also known to be effective, suggesting self-correction and verifier-search tap into related but distinct improvements.

Operational contribution to MAgHARCM:
- Provides the **SLM-era self-correction anchor**: a fine-tuned 7B self-corrector can replace a verifier-best-of-N pipeline, removing the per-query PRM cost.
- Validates the *training-time* route to self-correction as a cheaper alternative to the *inference-time* verifier-search route that dominates `[[1.0.0 P-91]]` Snell et al. and `[[1.0.0 P-92]]` Lightman et al.

## Method

- **Base models**: 7B–13B decoder-only transformer LMs (the paper uses LLaMA-style architectures). The recipe is model-agnostic.
- **Paired-edit training data**: `(y_t, y_{t+1})` pairs of incorrect → corrected outputs, sourced from either (a) iterative prompting of a stronger teacher model with self-correction, (b) rejection-sampling + automated fixers, or (c) human-written edits. The paper shows the *source* of pairs matters less than the *structure* of the training objective.
- **Iterative self-correction loss**: at training time, the model is conditioned on its *own* previously sampled outputs `y_0, y_1, ..., y_{T-1}` and trained to produce `y_T`. The loss handles off-policy drift between the model's current sampling distribution and the teacher distribution via importance weights or a behaviour-cloning term.
- **Inference-time refinement loop**: at decode time, the trained self-corrector runs a small fixed number of rounds (1–4 in the headline experiments): `y_0 = sample(base)`; `y_{t+1} = corrector(y_0, ..., y_t)`; return the final `y_T`. No external verifier, no rejection sampling.
- **Baseline for comparison**: PRM/ORM best-of-N (cf. `[[1.0.0 P-92]]` Lightman et al.) — sample N candidates, score with the verifier, return the best. The paper shows the trained self-corrector reaches comparable accuracy with substantially lower total inference compute.

## Findings Relevant to MAgHARCM

- **SLM self-correction is cheaper than PRM search.** The headline economic claim: for 7B–13B SLMs, training a self-corrector once and running it for 1–4 rounds per query is cheaper than running a separate verifier over best-of-N samples every query. MAgHARCM's deployment budget is dominated by inference cost on SLMs; P-97 makes the case for an amortised training cost paid once, in exchange for per-query savings forever after.
- **Prompting-only self-correction does NOT work.** A critical ablation: telling the base SLM "let me re-check that" via inference-time prompting does not improve outputs on most tasks — the base model is no better at finding its own errors than at not making them. MAgHARCM's `[[1.0.0 PRIM-25]]` Role-Flip Gate uses prompting to elicit critique, but P-97 warns that without a verifier or fine-tuning, the gate's rejections may be unreliable. Cite P-97 when justifying the *training-time* route.
- **Trained self-correction complements `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation.** P-97 gives MAgHARCM a way to ship a single self-correcting 7B SLM that *internally* does what the verdict panel does externally (multiple voters + revision). For cost-sensitive migration tasks, the trained self-corrector can substitute for a 3-voter verdict panel; for high-stakes tasks, the verdict panel remains the gold standard.
- **Verifier-free inference removes a class of adversarial attacks.** PRM best-of-N pipelines can be reward-hacked by the generator (cf. `[[1.0.0 P-92]]` Lightman et al. note on Goodhart's law); a trained self-corrector does not have a separately optimisable verifier to attack. This is a robustness, not just a cost, argument.
- **Training-data construction is the bottleneck.** The paper shows the paired-edit training data is the dominant cost and the dominant failure mode. MAgHARCM already produces (incorrect, corrected) pairs naturally via its execution oracle `[[1.0.0 PRIM-12]]` — failed test runs and the subsequent passing fix are exactly the training signal P-97 wants. The migration task is a *natural source* of paired-edit data.
- **Iterative training is required, not optional.** Single-round training on `(y_t, y_{t+1})` pairs does not generalise. The paper's multi-round iterative training is what makes the corrector robust to its own off-policy outputs. MAgHARCM's training pipeline must implement the multi-round objective, not a simplified single-round proxy.

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation**: train a 7B–13B self-corrector on MAgHARCM-produced (incorrect, corrected) pairs derived from `[[1.0.0 PRIM-12]]` Wasm-Based Reference Execution Oracle failures. Deploy the self-corrector as a *cheaper alternative* to a 3-voter verdict panel for routine migration tasks; keep the full verdict panel for high-stakes tasks where compute budget permits. Cite P-97 as the training-time self-correction anchor that justifies the trade-off.
- **`[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination Gate**: the gate's prompting-only critique is unreliable per P-97's ablation. Replace the prompt-only gate with a *trained* self-corrector for code-review tasks where the gate is the primary defence. Keep the prompt-only gate as a fallback when a fine-tuned corrector is not available.
- **`[[1.0.0 PRIM-12]]` Wasm-Based Reference Execution Oracle → Training-Data Flywheel**: oracle failures are exactly the `(y_t, y_{t+1})` paired-edit data P-97 trains on. Wire the oracle's failure log into a self-correction training pipeline so every MAgHARCM run that catches a bug produces a future training example. Cite P-97 as the recipe that consumes the data.
- **`[[1.0.0 PRIM-21]]` Migration Strategy Selection**: when the strategy registry routes to BIG_BANG or INCREMENTAL on a high-difficulty migration, prefer a self-correcting SLM over a verifier-best-of-N pipeline to stay within budget. Cite P-97 for the cost comparison.
- **Future work (not implemented)**: train a 7B Qwen2.5-Coder self-corrector on MAgHARCM's accumulated oracle-failure pairs; benchmark against the existing 3-voter verdict panel on the comprehension/translation/review suite; report the compute vs accuracy trade-off.

## References

### Hop-1 (papers that build directly on P-97 or are direct descendants)
- Snell, C. et al. (2024). *Scaling LLM Test-Time Compute Optimally Can be More Effective than Scaling Model Parameters*. arXiv:2408.03314. See `[[1.0.0 P-91]]` — the compute-optimal framework that P-97's verifier-free inference competes against. P-91 and P-97 are the two routes to "more inference accuracy": P-91 spends inference compute on a verifier, P-97 spends training compute on a corrector.
- Lightman, H. et al. (2023). *Let's Verify Step by Step*. arXiv:2305.02188. See `[[1.0.0 P-92]]` — the PRM800K verifier training that P-97's trained self-corrector aims to *replace* at inference time. P-97 shows comparable accuracy with no verifier at decode.
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. See `[[1.0.0 P-52]]` Self-Consistency lineage (P-90 → P-52) — the inference-time-sampling family that P-97 contrasts against. P-97's training-time self-correction is the alternative to sampling-based aggregation.
- Huang, J. et al. (2023). *Self-Checker: LLMs Can Self-Checkout Their Reasoning Errors via Iterative Decoding*. arXiv:2402.02687. P-97-style iterative refinement with an LLM-as-judge; predates and parallels P-97 in the iterative self-correction family.
- Madaan, A. et al. (2023). *Self-Refine: Iterative Refinement with Self-Feedback*. arXiv:2303.08151. The prompting-only iterative refinement baseline that P-97's trained approach improves on.

### Hop-2 (foundational anchors referenced)
- Rafailov, R. et al. (2023). *Direct Preference Optimization: Your Language Model is Secretly a Reward Model*. NeurIPS 2023. arXiv:2305.18290. See `[[1.0.0 P-93]]` — DPO is the SLM-era alignment technique that P-97's fine-tuning procedure is most natural to combine with. The self-corrector is fine-tuned with DPO-style preference data on (incorrect, corrected) pairs.
- Schick, T. et al. (2022). *Peer Review in Large Language Models: A Case Study of Two Stack-Overflow-Inspired Self-Correction Methods*. arXiv:2302.01828. Early multi-agent self-correction framework that motivates P-97's trained approach.
- Welleck, S. et al. (2020). *Neural Text Generation with Unlikelihood Training*. ICLR 2020. arXiv:2009.06534. Earlier Welleck / Lu work on iterative sequence-level training; the technical lineage of the iterative self-correction objective.
- Stiennon, N. et al. (2020). *Learning to Summarize with Human Feedback*. NeurIPS 2020. arXiv:2009.01325. The RLHF-from-preferences training pipeline that P-97's paired-edit data connects to.
- Touvron, H. et al. (2023). *LLaMA: Open and Efficient Foundation Language Models*. arXiv:2302.13971. The base model architecture P-97's experiments use.

### MAgHARCM lineage cross-refs
- `[[1.0.0 P-91]]` — Test-Time Compute Scaling (the verifier-search route P-97 competes against).
- `[[1.0.0 P-92]]` — Lightman PRM (the verifier P-97's trained self-corrector replaces at decode time).
- `[[1.0.0 P-90]]` — Chain-of-Thought (the prompting substrate P-97's training builds on).
- `[[1.0.0 P-52]]` — Self-Consistency (the sampling-aggregation alternative to trained self-correction).
- `[[1.0.0 P-93]]` — DPO (the natural fine-tuning recipe to combine with P-97's paired-edit objective).
- `[[1.0.0 P-94]]` — LIMA (the data-curation discipline P-97's paired-edit data inherits from).
- `[[1.0.0 P-95]]` — Code Llama (the 7B–34B SLM fleet P-97's trained self-corrector runs on).

## Backlinks

`[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-12]]`, `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-25]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-97 is the **SLM-era cheap self-correction anchor** for MAgHARCM. It provides the training-time alternative to inference-time verifier search: amortise a one-time fine-tuning cost on paired-edit data (which MAgHARCM already produces via its execution oracle) and ship a self-correcting SLM that needs no PRM at decode. Cite P-97 whenever MAgHARCM chooses trained self-correction over verifier-best-of-N, or whenever the verdict panel's cost is questioned.
