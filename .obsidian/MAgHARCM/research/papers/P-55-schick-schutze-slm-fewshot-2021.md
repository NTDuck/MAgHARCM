---
title: "P-55 — Schick & Schütze 2021 — It's Not Just Size That Matters: Small Language Models Are Also Few-Shot Learners"
backlink: "[[1.0.0 P-55]]"
tags: [paper, slm, few-shot, cloze, pet, pattern-exploiting-training, prompt, [[1.0.0 PRIM-3]], [[1.0.0 PRIM-22]], [[1.0.0 PRIM-24]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-55 — Schick & Schütze 2021 — Small LMs Are Also Few-Shot Learners]]

- **Authors**: Timo Schick, Hinrich Schütze (Center for Information and Language Processing, LMU Munich).
- **Venue / Year**: NAACL 2021; arXiv:2009.07118 (September 2020).
- **URL**: https://arxiv.org/abs/2009.07118
- **Anchors**: PRIM-3 (Target Skeleton-First), PRIM-22 (Four Phases Comprehension), PRIM-24 (SOP-Anchored Role-Artifact Schema); foundational SLM-prompting anchor for the entire 4B-30B regime MAgHARCM targets.

## 1. Core Contribution

A direct rebuttal to the "scaling-hypothesis" orthodoxy following GPT-3 [[Brown-2020]]: **"small" (≤500M parameter) LMs can match GPT-3 175B on SuperGLUE-style tasks when given the right prompting regime**. The recipe has three ingredients:

1. **Cloze reformulation** — convert every task into a fill-in-the-blank pattern. For classification: "A highly exciting movie. The sentiment is ___" ; for QA: "Q: ... A: ___". This converts *free-form generation* into a constrained output that small models can handle reliably.
2. **Pattern-Exploiting Training (PET)** — train the small model on unlabeled data using the cloze patterns as a self-supervised auxiliary objective, combined with supervised fine-tuning on a small labeled seed set. This is a *parameter-efficient* alternative to scale.
3. **iPET** (iterative PET) — train multiple PET models on different random seeds, let them label the unlabeled set, then bootstrap on the soft-labels. The ensemble is structurally a **self-consistency** vote over small-model variants — conceptually adjacent to Wang et al. [[P-52]] but predating it for SLMs.

**SLM-aware relevance for MAgHARCM**: the cloze reformulation pattern is the direct intellectual ancestor of every "fill in the function signature" or "complete the trait body" prompt that MAgHARCM issues to Qwen2.5-Coder [[P-21]], StarCoder2 [[P-22]], and Phi-3-mini [[P-54]]. The principle "offload cognitive complexity into the prompt structure" is exactly what makes SLM-grade code generation work.

## 2. Application in MAgHARCM

- **PRIM-3 (Target Skeleton-First)** — the cloze reformulation is structurally identical to skeleton-first translation: present the target shape (trait signature, function signature with explicit return type) and ask the model to fill in only the body. The skeleton's high-attention position (start of prompt, per Lost-in-the-Middle [[P-53]]) is what lets a 3B-7B model succeed where unconstrained generation fails.
- **PRIM-22 (Four Phases Comprehension)** — the cloze reformulation is also the recommended pattern for the **articulate phase**: instead of asking "summarize this function", ask the model to fill in slots like `Returns: ___ ; Side effects: ___ ; Calls: ___`. Constrained output formats are what makes comprehension traces parseable by downstream agents.
- **PRIM-24 (SOP-Anchored Role-Artifact Schema)** — PET's iterative soft-label bootstrap (iPET) is structurally a role-specialization loop: each PET model is a *role*, and their soft-labels are aggregated. The SOP-Anchored Role-Artifact schema in MAgHARCM treats each agent's role as a PET-style specialization, with the verdict panel [[P-52]] as the aggregation step.
- **SLM Prompt Contract** — the "convert every task to a cloze question" rule is encoded in `compiletime.SLMPromptContractPreamble` as the "Structured Output Slots" section. Every prompt template pre-declares the output slots (e.g., `TranslatedCode:`, `Rationale:`, `Confidence:`) so the model never has to invent structure.
- **Training-data scarcity mitigation** — MAgHARCM cannot fine-tune a 7B parameter SLM on every legacy codebase it encounters. PET's recipe (small labeled seed + large unlabeled corpus + cloze auxiliary loss) is the conceptual template for the lightweight LoRA adapter we could spin up per-domain if needed.

## 3. Hop-1 References (papers cited by Schick & Schütze)

- Brown et al. (2020) — GPT-3 / Language Models are Few-Shot Learners (the scaling hypothesis that Schick & Schütze rebut).
- Devlin et al. (2019) — BERT (the cloze-pattern origin; BERT's `[MASK]` token is the cloze primitive).
- Vaswani et al. (2017) — Attention Is All You Need (transformer architecture).
- Liu et al. (2019) — RoBERTa (a stronger cloze-trained baseline that PET improves on).
- [[NEEDS-LINK Howard-Ruder-2018]] — ULMFiT (foundational transfer-learning recipe for small LMs).
- Radford et al. (2019) — Language Models are Unsupervised Multitask Learners (GPT-2; few-shot baseline).
- Sanh et al. (2019) — DistilBERT (compression-from-large-LM recipe referenced as alternative to PET).
- Houlsby et al. (2019) — Parameter-Efficient Transfer Learning with Adapters (referenced as a parameter-efficient alternative).
- Wang et al. (2019) — GLUE / SuperGLUE (the evaluation suite PET is benchmarked on).
- McCann et al. (2018) — Natural Language Decathlon (MT-DNN; multitask evaluation cited by PET).

## 4. Hop-2 References (papers-cited-by-hop-1)

- Mikolov et al. (2013) — Distributed Representations of Words and Phrases (word2vec; original cloze-style unsupervised LM).
- Peters et al. (2018) — ELMo (foundational contextualized embeddings; cited by ULMFiT).
- Kaplan et al. (2020) — Scaling Laws for Neural Language Models (the scaling hypothesis that PET challenges; cited by GPT-3).
- Raffel et al. (2020) — T5 / Exploring the Limits of Transfer Learning (text-to-text unification; referenced by MT-DNN).
- Clark et al. (2020) — ELECTRA (replaced-token-detection pretraining; cited by RoBERTa).
- [[NEEDS-LINK Kingma-Ba-2015]] — Adam (optimizer; cited by ULMFiT).
- Hochreiter & Schmidhuber (1997) — LSTM (the recurrent baseline that ULMFiT predates).
- Wolf et al. (2020) — Transformers library (referenced by PET's open-source release).
- Liu et al. (2024) — Lost in the Middle [[P-53]] in MAgHARCM lineage (cited by BERT-lineage papers as evidence for structured-input benefits).

## 5. Backlinks

- **PRIM-3** (Target Skeleton-First): skeleton = cloze-pattern prompt structure for the SLM to fill in.
- **PRIM-22** (Four Phases Comprehension): phase-specific prompt slots = cloze reformulation of comprehension tasks.
- **PRIM-24** (SOP-Anchored Role-Artifact Schema): iPET's role-soft-label bootstrap is the SOP-anchored role pattern.
- **PRIM-25** (Role-Flip De-Hallucination): role-specialization (PET vs role-flip gate) shares the "ask the same question two different ways" pattern.
- **Cross-ref**: add P-55 to PRIM-3, PRIM-22, PRIM-24 rows in `Software-Archaeology-Lineage.md`. This is the **canonical SLM-prompting anchor** for MAgHARCM's prompt-engineering decisions.
