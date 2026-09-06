---
title: "P-53 — Liu et al. 2023 — Lost in the Middle: How Language Models Use Long Contexts"
backlink: "[[1.0.0 P-53]]"
tags: [paper, long-context, retrieval, chunking, [[1.0.0 PRIM-31]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-53 — Liu et al. 2023 — Lost in the Middle]]

- **Authors**: Nelson F. Liu, Kevin Lin, John Hewitt, Ashwin Paranjape, Michele Bevilacqua, Fabio Petroni, Percy Liang (Stanford / Princeton / Samaya AI).
- **Venue / Year**: TACL 2024 (arXiv:2307.03172, July 2023).
- **URL**: https://arxiv.org/abs/2307.03172
- **Anchors**: PRIM-31 (IterativeNavigator), PRIM-22 (Comprehension), chunked translation path.

## 1. Core Contribution

Empirically demonstrates **U-shaped attention degradation** in transformer LMs:
- Models perform well when relevant information is at the BEGINNING or END of a long context.
- Performance drops 20-40% when relevant information is in the MIDDLE.
- The effect is consistent across open (MPT-30B, LongChat-13B) and closed (Claude 1.3, GPT-3.5-16k) models.
- Practical implication: **re-order context** so critical information is at the start/end, OR use retrieval-augmented context (RAG) to inject only relevant snippets.

## 2. Application in MAgHARCM

- Validates the **chunked translation** strategy: when a source file exceeds `IterativeContextBudgetBytes` (4 KiB, see `compiletime.IterativeContextBudgetBytes`), we split into navigable fragments rather than feeding the entire file.
- Validates the **IterativeNavigator** [[1.0.0 PRIM-31]] re-indexed/fresh split: re-indexed symbols go to the START of the prompt context (high-attention zone), fresh navigations go to the END.
- Justifies the `compiletime.DefaultProjectDir` plus per-symbol navigation — we never dump the whole repo into context.

## 3. Hop-1 References (papers cited by Liu et al.)

- Vaswani et al. (2017) — Attention Is All You Need (foundational).
- Press et al. (2022) — ALiBi positional encoding (referenced as a long-context enabler).
- Chen et al. (2023) — LongLoRA / extending context windows via fine-tuning.
- Tworkowski et al. (2023) — Focused Transformer (contrastive retrieval-augmented attention).
- Beltagy et al. (2020) — Longformer (sparse-attention baseline).
- Dao et al. (2022) — FlashAttention (memory-efficient attention).

## 4. Hop-2 References (papers-cited-by-hop-1)

- Devlin et al. (2019) — BERT (foundational transformer encoder).
- Su et al. (2021) — RoPE (rotary position embeddings, cited by LongLoRA).
- Su et al. (2024) — RoPE scaling extensions.
- Zaheer et al. (2020) — Big Bird (sparse-attention model for long documents).
- Choromanski et al. (2021) — Performers (linear-attention approximation).

## 5. Backlinks

- PRIM-31 (IterativeNavigator): the `reindexed` / `fresh` source-status split IS the "U-shaped attention" mitigation strategy.
- PRIM-22 (Comprehension): comprehension runs on selected fragments, never whole-repo context.
- PRIM-3 (Target Skeleton-First): skeleton-first ensures the target structure is at the END of the prompt (high-attention zone) when refinement prompts run.
- Cross-ref: add P-53 to PRIM-31 + PRIM-22 + PRIM-3 rows in `Software-Archaeology-Lineage.md`.
