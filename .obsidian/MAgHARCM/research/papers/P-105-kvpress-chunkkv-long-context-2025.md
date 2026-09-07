---
title: "P-105 — KV Cache Compression for Long-Context LLM Code (NeurIPS 2025)"
backlink: "[[1.0.0 P-105]]"
tags: [paper, kv-cache, long-context, compression, slm, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-31]], hop-1]
---

# [[1.0.0 P-105 — ChunkKV: Semantic-Preserving KV Cache Compression]]

## Citation

Liu, A. et al. (2025). *ChunkKV: Semantic-Preserving KV Cache Compression for Efficient Long-Context Inference*. NeurIPS 2025.
URL: https://proceedings.neurips.cc/paper_files/paper/2025/file/2987f911151b39cd3a1761e212319e8e-Paper-Conference.pdf

## Summary

ChunkKV compresses KV cache by **chunking context into semantically coherent segments** and evicting low-importance chunks. For long-context LLM inference (1M+ tokens), it achieves **4-8× memory reduction** with negligible quality loss. Companion tool: **NVIDIA/kvpress** provides a unified interface to 12+ KV-cache compression methods (StreamingLLM, H2O, Scissorhands, etc.).

The SLM angle: 7B-13B models at 100K context can run on a single 24GB GPU with ChunkKV, whereas uncompressed they require 2-3 GPUs. Critical for edge-deployed code-comprehension SLMs.

## Relevance to MAgHARCM

- **Direct evolution of P-80 StreamingLLM** for SLM long-context. MAgHARCM's PRIM-22 Comprehension Stage + PRIM-31 Iterative Retrieval both consume large context windows; ChunkKV enables both on SLM hardware.
- **Chunking strategy mirrors PRIM-23 Chunked Translation**: ChunkKV's semantic-chunk eviction aligns with chunked-translator's chunk boundary detection. Same chunk granularity → better cross-chunk state preservation.
- **KVPress library**: a drop-in candidate for the LLM provider layer (`internal/llm/`). The user's directive: "Adopt more externalities" — KVPress is a direct opportunity.
- **Edge-deployment feasibility**: enables MAgHARCM on 4B-13B models at 100K+ token context without GPU cluster.

## References (hop-1)

- [[1.0.0 P-80]] StreamingLLM attention-sink sliding window
- [[1.0.0 P-68]] Shazeer Multi-Query Attention
- [[1.0.0 P-57]] Leviathan Speculative Decoding (composes with KV compression)

## References (hop-2)

- SCOPE (ACL 2025 oral): stage-level KV cache optimisation
- H2O (Zhang et al.): heavy-hitter oracle
- Scissorhands: importance-score-based eviction
