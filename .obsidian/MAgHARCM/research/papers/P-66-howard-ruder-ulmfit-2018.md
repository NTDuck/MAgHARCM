---
title: "P-66 — Howard & Ruder 2018 — Universal Language Model Fine-tuning for Text Classification (ULMFiT)"
backlink: "[[2.0.0 P-66]]"
tags: [paper, transfer-learning, slm, language-model, fine-tuning, discriminative-learning-rate, [[1.0.0 P-55]], hop-2]
---

# [[2.0.0 P-66 — ULMFiT]]

## Citation

Howard, J., & Ruder, S. (2018). *Universal Language Model Fine-tuning for Text Classification*. ACL 2018. arXiv:1801.06146. URL: https://arxiv.org/abs/1801.06146.

## Summary

ULMFiT introduced the three-stage transfer-learning recipe that became the SLM playbook: pre-train a language model on a large general corpus; fine-tune it on the target task's data with **discriminative learning rates** (lower layers learn slower, top layers learn faster, mirroring the fact that low-level features transfer across tasks); apply **gradual unfreezing** (unfreeze one layer at a time during fine-tuning, top-down, to avoid catastrophic forgetting). The paper reports **18-24% error reduction** on six text-classification benchmarks vs. training-from-scratch baselines, while fine-tuning the same model on as few as **100 labeled examples**. The recipe is the SLM-friendly precedent for MAgHARCM's per-domain LoRA adapters (cf. [[1.0.0 P-58]] Qwen2.5-Coder synthetic-data recipe).

## Method

Three stages:
1. **LM pre-training** on a large general corpus (Wikitext-103, 103M tokens).
2. **LM fine-tuning** on target task data with discriminative learning rates (each layer $l$ gets $\eta^{l}$ where $\eta$ decays geometrically from top to bottom).
3. **Classifier fine-tuning** with gradual unfreezing (unfreeze the topmost layer first, add a new classifier on top, fine-tune for one epoch; then unfreeze the next layer down; repeat).

## Findings Relevant to MAgHARCM

- **Discriminative learning rates** transfer directly to per-domain LoRA adapter training on Qwen2.5-Coder (cf. [[1.0.0 P-58]]): when adapting the SLM to a specific legacy codebase's domain patterns, lower transformer layers (token embeddings, low-level syntax) need less perturbation than upper layers (code-specific patterns, domain vocabulary).
- **Gradual unfreezing** maps to staged adapter training: train the LoRA on the comprehension task first, then the planning task, then the translation task — each task's adapter is trained without disturbing the others.
- **Few-shot transferability** (100 labeled examples) is the empirical justification for MAgHARCM's per-fragment training set in [[1.0.0 PRIM-22]] comprehension phases: the SLM can absorb a new codebase's conventions from a small labeled sample.

## How MAgHARCM Uses It

[[1.0.0 P-55]] (Schick & Schütze 2021) cites ULMFiT as the foundational transfer-learning recipe for SLMs. ULMFiT's discriminative-learning-rate discipline is the empirical template for MAgHARCM's per-domain LoRA adapter training: when adapting Qwen2.5-Coder-7B-Instruct to a specific legacy codebase's domain, freeze the bottom 16 transformer layers and train a small LoRA on the top 8 layers plus a 2-layer domain classifier head.

## References

### Hop-1 (Howard & Ruder 2018 cites)
- Mikolov, T. et al. (2013). *Distributed Representations of Words and Phrases and their Compositionality*. (word2vec)
- Peters, M. et al. (2018). *Deep contextualized word representations* (ELMo).
- McCann, B. et al. (2017). *Learned in Translation: Contextualized Word Vectors* (CoVe).
- Kingma, D. & Ba, J. (2015). *Adam: A Method for Stochastic Optimization*. ICLR 2015. arXiv:1412.6980. See [[2.0.0 P-67]].
- Ruder, S. (2016). *An Overview of Gradient Descent Optimization Algorithms*. arXiv:1609.04747.

### Hop-2
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762.
- Devlin, J. et al. (2019). *BERT: Pre-training of Deep Bidirectional Transformers for Language Understanding*. NAACL 2019. arXiv:1810.04805.

## Backlinks

[[1.0.0 P-55]], [[1.0.0 P-58]], [[1.0.0 P-67]], [[1.0.0 PRIM-22]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-66 is the **transfer-learning precursor** for MAgHARCM's per-domain LoRA recipe.
