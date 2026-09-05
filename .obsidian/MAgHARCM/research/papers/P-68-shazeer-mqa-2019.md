---
title: "P-68 — Shazeer 2019 — Fast Transformer Decoding: One Write-Head is All You Need (MQA)"
backlink: "[[2.0.0 P-68]]"
tags: [paper, transformer, inference-speed, multi-query-attention, kv-cache, [[1.0.0 P-57]], hop-2]
---

# [[2.0.0 P-68 — Multi-Query Attention]]

## Citation

Shazeer, N. (2019). *Fast Transformer Decoding: One Write-Head is All You Need*. arXiv:1911.02150. URL: https://arxiv.org/abs/1911.02150.

## Summary

Multi-Query Attention (MQA) replaces the per-head query/key/value projections with a **single shared key and value projection** across all heads, while keeping separate query projections per head. The change:

- **Memory**: KV-cache size shrinks by a factor of $H$ (the number of attention heads), reducing from $2 \cdot H \cdot d_{\text{head}}$ per token to $2 \cdot d_{\text{head}}$ per token.
- **Compute**: autoregressive attention still requires $H \cdot d_{\text{head}}$ per layer, but the KV-projection bandwidth saved is significant at long context lengths.
- **Quality**: empirically **<1% quality drop** on machine-translation benchmarks, often with **no measurable drop** when fine-tuned.

This is the foundational inference-speedup technique for serving LLMs at scale. Every modern open-weights model with fast inference (PaLM, Gemma, Qwen2.5-Coder, LLaMA-3) uses either MQA or its generalisation GQA (Grouped-Query Attention, Ainslie et al. 2023).

## Method

Standard multi-head attention projects $x \to (Q, K, V)$ with $H$ separate head projections. MQA shares $K$ and $V$ across all $H$ heads:

$$\text{MQA}(x) = \text{Concat}(\text{head}_1, \ldots, \text{head}_H)$$
$$\text{head}_i = \text{softmax}\!\left(\frac{Q_i K^\top}{\sqrt{d}}\right) V$$
where $K$ and $V$ are computed once and reused across all $H$ heads.

The paper shows that the quality loss from sharing K/V is small because the attention pattern is dominated by the query projection's content; the per-head value projection mostly averages to a similar per-position weight.

## Findings Relevant to MAgHARCM

- **KV-cache size** directly determines **inference latency** at long contexts. MAgHARCM's comprehension phase reads entire legacy source files (often >10k tokens). MQA / GQA makes this tractable on consumer GPUs.
- **Qwen2.5-Coder (cf. [[1.0.0 P-58]])** uses **GQA** (Grouped-Query Attention) — 8 KV-heads shared across 14 query heads, giving roughly 2× KV-cache reduction over full MHA. This is why Qwen2.5-Coder fits in 24GB VRAM at 32k context.
- **Speculative decoding** (cf. [[1.0.0 P-57]], [[1.0.0 P-61]], [[1.0.0 P-64]]) composes with MQA: the draft model runs MQA, the target model runs MQA, KV-cache bandwidth is the bottleneck in both. The two techniques are orthogonal and stack.

## How MAgHARCM Uses It

When configuring the inference backend for the comprehension + planning + translation agents, MAgHARCM prefers models with MQA or GQA. The `compiletime.DefaultInferenceKVHeads` constant (a small integer like 8) controls the expected KV-head count; the runner asserts the loaded model matches.

## References

### Hop-1 (Shazeer 2019 cites)
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762.
- Devlin, J. et al. (2019). *BERT*. NAACL 2019. arXiv:1810.04805.
- Radford, A. et al. (2019). *Language Models are Unsupervised Multitask Learners*. (GPT-2)
- Ho, J., et al. (2019). *Axial Attention in Multidimensional Transformers*. arXiv:1912.12180.

### Hop-2
- Ainslie, J. et al. (2023). *GQA: Training Generalized Multi-Query Transformer Models from Multi-Head Checkpoints*. arXiv:2305.13245.
- Leviathan, Y., Kalman, M., Matias, Y. (2023). *Fast Inference from Transformers via Speculative Decoding*. See [[1.0.0 P-57]].
- Touvron, H. et al. (2023). *LLaMA: Open and Efficient Foundation Language Models*. arXiv:2302.13971.

## Backlinks

[[1.0.0 P-57]], [[1.0.0 P-58]], [[1.0.0 P-61]], [[1.0.0 P-64]], [[1.0.0 PRIM-22]], [[2.0.0 MAgHARCM]].

P-68 is the **KV-cache compression** technique that makes long-context SLM inference tractable; it composes with speculative decoding (P-57, P-61, P-64) for end-to-end latency reduction.
