---
title: "P-134 RelayCaching: Accelerating LLM Collaboration via Decoding KV Cache Reuse"
backlink: "[[1.0.0 P-134]]"
aliases:
  - "1.0.0 P-134"
  - "P-134"
  - "P-134-RelayCaching-ICML-2026"
  - "RelayCaching-ICML-2026"
  - "geng2026relaycaching"
tags: [paper, kv-cache, multi-agent, speculative-decoding, "[[1.0.0 PRIM-31]]", "[[1.0.0 PRIM-9]]", wave-19]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICML 2026 (Poster #1915)
---

# [[1.0.0 P-134]] RelayCaching

## TL;DR

RelayCaching extends single-query KV cache reuse to **multi-agent LLM pipelines** by detecting sparse prefix-induced deviations across agent outputs and selectively recomputing only the deviation positions in shared KV tensors. Achieves **>80% reuse rate and 4.7× TTFT reduction** in agent collaborations. Anchors `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (cross-agent prefill) and `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph (shared context as structural graph).

## Mechanism (Q2)

1. **Sparse deviation recompute** — when two agents share a long prefix (e.g. system prompt + prior tool outputs), RelayCaching detects the small set of positions where the second agent's KV differs from the first agent's; only those positions are recomputed; the rest of the KV is reused as-is.
2. **Training-free** — no draft model, no fine-tuning. Compatible with any backbone that uses KV cache (transformer-family decoders).
3. **Multi-agent workflow-aware** — exposes a hookable API for the workflow scheduler to declare "shared prefix boundary" at agent-handoff time.

Together these give the 8-agent MAgHARCM graph (research→planner→architect→developer→tester→reviewer→verifier→historian) a concrete cross-turn KV reuse substrate that Wave-18 KVzip (P-128, single-query) does not address.

## Anchoring (Q3)

| Primitive | Pre-wave-19 behaviour | RelayCaching substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Per-agent prefill cost dominates end-to-end latency on the 8-agent graph; TTFT is linear in agent count | Sparse-deviation recompute shares prefill across adjacent agents; TTFT drops from O(N·P) to O(P + N·Δ) where Δ ≪ P |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Cross-agent token overlap was implicit (prefix-sharing heuristics) | Shared-context graph with explicit agent↔agent overlap edges drives the recompute policy |

## Hop-1 Citations

- P-57 Leviathan 2023 Speculative Decoding — cross-token reuse lineage.
- P-141 KVFlow 2025 (this wave) — workflow-aware eviction policy inspiration.
- P-128 KVzip 2025 — single-query KV compression baseline.

## Hop-2 Citations

- Shazeer 2019 Multi-Query Attention (P-68) — KV sharing foundations.
- Ainslie 2023 GQA — grouped-query attention KV sharing.
- Prefix caching vLLM / SGLang radix tree (2024).

## MAgHARCM integration

- **YAML config key**: `kv_cache.relay.enabled: true` (opt-in, default `false`); `kv_cache.relay.deviation_threshold: <float>`.
- **Implementation file**: `internal/kvrelay/relay_cache.go::NewRelayCache` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-31]]`, `[[1.0.0 PRIM-9]]`.

## Caveats

- **Deviation detection cost** — sparse deviation requires a fast KV-comparison pass; cost grows with shared-prefix length. Acceptable when Δ ≪ P but degrades on divergent agent outputs.
- **Shared-prefix semantics** — workflow scheduler must declare prefix boundaries correctly; misdeclaration under-reuses and risks cross-agent contamination.
- **Backbone constraint** — only transformer-family decoders with explicit KV cache.

## Source

- arXiv: 2603.13289.
- Venue: ICML 2026 Poster #1915 (OpenReview id 1tbhBSXcyX; verified at icml.cc/virtual/2026/poster/66638).
