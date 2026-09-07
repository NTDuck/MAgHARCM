---
title: "P-94 — Zhou et al. 2023 — LIMA: Less Is More for Alignment"
backlink: "[[1.0.0 P-94]]"
aliases:
  - "1.0.0 P-94"
  - "P-94"
  - "P-94-Zhou-LIMA-2023"
  - "P-94-Zhou-LIMA-2023"
  - "Zhou-LIMA-2023"
tags: [paper, lima, alignment, data-curation, slm, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-24]], [[1.0.0 PRIM-25]], hop-1]
---

# [[1.0.0 P-94 — Zhou et al. — LIMA: Less Is More for Alignment]]

## Citation

Zhou, C., Liu, P., Xu, P., Iyer, S., Sun, J., Mao, Y., Ma, X., Efrat, A., Yu, P., Zhang, S., Ghosh, G., Lewis, M., Zettlemoyer, L., & Levy, O. (2023). *LIMA: Less Is More for Alignment*. Advances in Neural Information Processing Systems 36 (NeurIPS 2023). arXiv:2305.11206 (18 May 2023 v1). Meta AI / FAIR.

## Summary

LIMA is the strongest empirical demonstration that **alignment data quality dominates alignment data quantity**. The authors fine-tune a 65B-parameter LLaMA model (LLaMA-7B ablations are also reported) on a deliberately tiny, carefully curated dataset of **1,000 training examples** — and achieve performance comparable to GPT-4 on the Alpaca benchmark.

The dataset construction is the contribution:
- **1,000 examples total**, split into 750 training + 250 validation.
- Sources: 1,000 prompt-response pairs sampled from community Q&A forums (Stack Exchange, WikiHow, Reddit), filtered and rewritten by AI assistants.
- Each prompt is paired with a *single*, high-quality response (as judged by human annotators).
- Responses are stylistically and factually polished; they avoid harmful, misleading, or low-quality content.

The model's behaviour after LIMA fine-tuning is characterised as **"superficial alignment"**: most of the LLM's knowledge and capabilities come from the pre-training, and the 1K fine-tuning examples merely steer the model toward the desired *style* and *format*. The implication is that high-quality examples with **broad coverage of the desired distribution** produce strong alignment with minimal data.

The headline results:
- LIMA-65B trained on 1K examples matches GPT-4 on Alpaca-style human preference (43.9% vs 43.7% win rate vs DaVinci003).
- Human preference for LIMA-65B vs Bard: 58% win rate despite 1/1000th the alignment data.
- LIMA-7B (7B variant) is notably worse than LIMA-65B, reinforcing the **scale dependence** of the technique.

The paper's most-quoted claim: *"Almost all knowledge in a large language model is learned during pretraining, and a relatively small amount of instruction tuning data is sufficient to elicit high-quality output."*

## Method

- **Base model**: LLaMA-65B (65B parameters, public release); ablations on LLaMA-7B/33B.
- **Training data**: 1,000 prompt-response pairs, hand-curated. ~750 train + ~250 val.
- **Data construction**:
  - Pull candidate prompts from Stack Exchange, WikiHow, Reddit (r/AskHistorians, r/AskScience).
  - Filter for diversity: aim for broad coverage of topics, styles, and formats.
  - For each prompt, AI assistants (Bard, Vicuna, then human rewriting) generate a single high-quality response.
  - Authors manually review and reject prompts/responses that are low-quality, harmful, or duplicative.
- **Fine-tuning**: standard SFT with the 1K examples, two epochs, learning rate 1e-5.
- **Evaluation**: Alpaca-style human preference (DaVinci003 as baseline), adversarial prompts, OOD prompts.

## Findings Relevant to MAgHARCM

- **1,000 high-quality examples beat 100,000 mediocre ones.** This is the operational justification for MAgHARCM's prompt-engineering investment. The 31 primitives and the SOP-anchored role schema (`[[1.0.0 PRIM-24]]`) are the "curation" layer for SLM behaviour; LIMA demonstrates the leverage is real.
- **LIMA directly motivates `[[1.0.0 P-84]]` s1.** s1's 1K-curated-reasoning-traces dataset (s1K) is the LIMA recipe applied to reasoning data. Cite P-94 whenever citing s1.
- **LIMA + DPO = a small-data alignment pipeline.** `[[1.0.0 P-93]]` DPO can run on as few as ~10K high-quality preference pairs when curated carefully (cf. Zephyr-7B). LIMA's 1K SFT stage + a ~10K DPO stage = ~11K total examples to align a 7B-65B SLM, far less than the 100K+ of naive RLHF recipes.
- **Scale matters: LIMA-7B is much weaker than LIMA-65B.** MAgHARCM's SLM tuning should target the upper edge of the SLM range (Qwen2.5-Coder-32B [[1.0.0 P-21]], CodeLlama-34B [[1.0.0 P-95]]) where LIMA-style alignment produces the strongest returns. Below 7B, fine-tuning on <1K examples is unlikely to produce a robustly-aligned model.
- **Curation discipline generalises.** The LIMA recipe — broad-coverage prompt selection, single high-quality response per prompt, human review for rejection — is the exact recipe MAgHARCM should apply when building the SLM-tuning corpus from its own migration logs.
- **"Superficial alignment" is operationally useful.** It means pre-training is the foundation, fine-tuning is the format. MAgHARCM should focus pre-training-corpus-quality efforts on the SLM base (Qwen2.5-Coder, StarCoder2, CodeLlama) rather than on bespoke fine-tuning data.

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-24]]` SOP-Anchored Role-Artifact Schema**: when building the SOP-tuning corpus for an MAgHARCM SLM, apply the LIMA discipline: aim for **~1,000 hand-curated** SOP examples, broadly covering the 31 primitives, each with a single canonical response. Reject duplicate or low-quality exemplars. Cite P-94 in the SOP corpus documentation.
- **`[[1.0.0 PRIM-22]]` Four Phases of Comprehension**: the comprehension exemplars should be drawn from the LIMA recipe — diverse, high-quality, single-response. Avoid mixing in noisy in-context examples.
- **`[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination Gate**: the role-flip reviewer's few-shot exemplars benefit from LIMA-style curation. A few hand-crafted "rejection" exemplars outperform a thousand random rejections.
- **`[[1.0.0 PRIM-3]]` Target Skeleton-First Generation**: skeletons should be tuned on ~1K curated examples (cf. MAgHARCM's own skeleton corpus in `internal/agents/planning.go`).
- **Future work (not implemented)**: build the MAgHARCM-SLM-tuning corpus using the LIMA recipe. Estimate ~1K hand-curated examples covering all 31 primitives, then run SFT (LIMA-style) followed by DPO (P-93) on ~10K verdict-panel-derived preference pairs.

## References

### Hop-1 (papers that build directly on P-94)
- Rafailov, R. et al. (2023). *Direct Preference Optimization*. arXiv:2305.18290. See `[[1.0.0 P-93]]` — DPO inherits LIMA's "less is more" data hypothesis.
- Muennighoff, N. et al. (2025). *s1: Simple Test-Time Scaling*. arXiv:2501.19393. See `[[1.0.0 P-84]]` — s1K is the LIMA recipe applied to reasoning traces.
- Tunstall, L. et al. (2023). *Zephyr: Direct Distillation of LM Alignment*. arXiv:2310.16944. Zephyr uses ~22K preference pairs (UltraFeedback) for DPO after SFT — a LIMA-inspired scale.
- Ivison, H. et al. (2023). *Camels in a Changing Climate: Enhancing LM Adaptation with Tulu 2*. arXiv:2311.10702. Tulu-2's data curation philosophy.
- Zhou, L. et al. (2024). *Larger and More Instructable Language Models Become Less Reliable*. Nature 634. Cited by s1 for the inverse-scaling concerns motivating small-data curation.

### Hop-2 (foundational anchors referenced)
- Ouyang, L. et al. (2022). *Training Language Models to Follow Instructions with Human Feedback* (InstructGPT). NeurIPS 2022. arXiv:2203.02155. The large-data RLHF baseline.
- Touvron, H. et al. (2023). *LLaMA: Open and Efficient Foundation Language Models*. arXiv:2302.13971. The base model LIMA uses.
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. See `[[1.0.0 P-90]]` — the prompt-format influence on response quality.
- Wang, Y. et al. (2023). *How Far Can Camels Go? Exploring the State of Instruction Tuning on Open Resources*. The Tulu-1 data-curation precedent.

### MAgHARCM lineage cross-refs
- `[[1.0.0 P-21]]` — Qwen2.5-Coder (the LIMA-curation target for MAgHARCM tuning).
- `[[1.0.0 P-84]]` — s1 (the LIMA-recipe reasoning descendant).
- `[[1.0.0 P-93]]` — DPO (the alignment technique that benefits from LIMA's data discipline).
- `[[1.0.0 P-22]]` — StarCoder2 (another LIMA-curation target).
- `[[1.0.0 P-95]]` — CodeLlama (Meta's code LLM family; the LIMA authors are from the same org).
- `[[1.0.0 P-55]]` — Schick & Schütze SLM Few-Shot (the earlier "1K samples can be enough" anchor).

## Backlinks

`[[1.0.0 PRIM-3]]`, `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-24]]`, `[[1.0.0 PRIM-25]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-94 is the **alignment-data-curation anchor** for MAgHARCM. It validates the discipline of small-data, high-quality fine-tuning over large-data, low-quality alignment. Cite P-94 whenever MAgHARCM builds a fine-tuning corpus for an SLM.
