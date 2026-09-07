---
title: "P-128 KVzip: Query-Agnostic KV Cache Compression with Context Reconstruction"
backlink: "[[1.0.0 P-128]]"
aliases:
  - "1.0.0 P-128"
  - "P-128"
  - "P-128-KVzip-KV-Cache-Compression-2025"
  - "KVzip-KV-Cache-Compression-2025"
  - "kim2025kvzip"
tags: [paper, kv-cache, inference-acceleration, compression, query-agnostic, "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-31]]", "[[1.0.0 P-80]]", "[[1.0.0 P-105]]", "[[1.0.0 P-68]]", wave-18]
date: 2026-09-07
last_updated: 2026-09-07
venue: NeurIPS 2025 (Oral)
---

# [[1.0.0 P-128]] KVzip

## TL;DR

KVzip is a **query-agnostic KV cache eviction method** that lets the LLM itself quantify per-pair importance by reconstructing the original context, producing a compressed cache that supports diverse future queries. Achieves **3-4× cache reduction and 2× decoding-latency decrease** in context-dependent mode, generic across LLaMA3, Qwen2.5/3, and Gemma3. Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement.

## Mechanism (Q2)

Two coupled mechanisms:

1. **Self-reconstruction importance scoring** — to score a KV pair, the LLM is asked to **reconstruct the original context** using only the remaining pairs. If reconstruction succeeds without a specific pair, that pair is low-importance and evictable; if reconstruction fails, the pair is high-importance.
2. **Query-agnostic eviction** — because importance is computed against the context itself (not a future query), the resulting compressed cache supports **diverse future queries** without re-compression.

This contrasts with query-aware eviction (StreamingLLM, ChunkKV) which optimizes the cache for a known query and must re-evict when the query changes. KVzip trades a one-time reconstruction pass for amortized savings across many queries on the same context. In context-dependent mode (long contexts with many follow-up queries) the amortization dominates.

## Anchoring (Q3)

| Primitive | Pre-wave-18 behaviour | KVzip substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Linear full-context attention; no cache eviction; query-aware eviction (StreamingLLM, ChunkKV) ties cache to a known query | LLM-self-reconstructed importance scoring; query-agnostic eviction; compressed cache is reusable across diverse queries |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Re-rank or re-fetch on each iteration; no cache-aware design | Iterative reuse of a single compressed cache across retrieval iterations; cost amortized |

## Hop-1 Citations

- `[[1.0.0 P-80]]` Xiao et al. 2024 *Efficient Streaming Language Models with Attention Sinks* (StreamingLLM, attention sinks).
- `[[1.0.0 P-105]]` ChunkKV (NeurIPS 2025) — chunk-wise KV compression baseline.
- `[[1.0.0 P-68]]` Shazeer 2019 *Fast Transformer Decoding: Multi-Query Attention* (MQA ancestor).
- H2O (Zhang et al. 2023) and Scissorhands (Tang et al. 2024) — KV-eviction lineage.

## Hop-2 Citations

- DuoAttention (NVIDIA, 2024) — separate full/sparse attention heads.
- FastGen (Ge et al. 2023) — dynamic KV compression policy.
- GQA (Ainslie et al. 2023) — grouped-query attention KV sharing.

## MAgHARCM integration

- **YAML config key**: `kv_cache.kvzip.enabled: true` (opt-in, default `false`); `kv_cache.kvzip.reconstruction_passes: <int>`.
- **Implementation file**: `internal/inference/kvzip.go::CompressContext` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **Reconstruction overhead**: the self-reconstruction pass is non-trivial; for single-query workloads it may not amortize. KVzip wins in long-context multi-query scenarios (agent loops, RAG re-ranking).
- **Reconstruction fidelity**: the quality of importance scores depends on the LLM's ability to reconstruct context from a partial cache. Very aggressive eviction can produce a self-fulfilling degradation.
- **Model coverage**: validated on LLaMA3, Qwen2.5/3, Gemma3; portability to other architectures (especially MoE) needs separate validation.

## Source

- arXiv:2505.23416 — KVzip (NeurIPS 2025 Oral).
