---
title: "P-138 RepairKV (Cache You Later): Post-Compression KV Repair for Long-Context Agentic LLM Inference"
backlink: "[[1.0.0 P-138]]"
aliases:
  - "1.0.0 P-138"
  - "P-138"
  - "P-138-RepairKV-ICML-2026"
  - "RepairKV-ICML-2026"
  - "rusli2026repairkv"
tags: [paper, kv-cache, repair, long-context, agentic, "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-31]]", wave-19]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICML 2026 AdaptFM Workshop (Poster) — borderline workshop-track ACCEPT per §7 method-level threshold
---

# [[1.0.0 P-138]] RepairKV

## TL;DR

RepairKV is a **post-compression KV cache repair** runtime operator that, between agent turns, re-evaluates evicted KV rows in a slower memory tier and promotes a small subset back to the active cache. Achieves **91.0% retrieval vs 24.5% no-repair baseline** at 32K context on Qwen2.5-7B. Two-way eviction+repair policy is a substantive architectural shift, not a prompt tweak. Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (comprehension phase "re-search" step) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (second-chance correction).

## Mechanism (Q2)

1. **Two-way eviction+repair policy** — between agent turns, RepairKV evaluates evicted KV rows in a slower memory tier (CPU or NVMe) and promotes a small subset back to the active GPU cache based on a learned or heuristic retrieval score.
2. **Runtime operator** — repair runs at agent-handoff time, not at initial prefill; the cost is amortised across the multi-turn loop.
3. **Long-context agentic target** — designed for 32K+ contexts where aggressive eviction (KVzip-style) loses precision; the repair pass recovers it.

Together these give MAgHARCM's multi-turn agent loop a **runtime correction substrate** when earlier evictions were wrong.

## Anchoring (Q3)

| Primitive | Pre-wave-19 behaviour | RepairKV substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Eviction was one-way; multi-turn retrieval could not recover evicted tokens | Two-way eviction+repair lets the comprehension phase's "re-search" step pull evicted tokens back into the active cache |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Iterative retrieval was bounded by the initial eviction's correctness | Repair gives the iterative loop a second-chance correction mechanism when an earlier eviction was wrong |

## Hop-1 Citations

- P-128 KVzip 2025 (query-agnostic eviction baseline).
- P-80 StreamingLLM 2024 (attention sinks eviction).
- P-105 ChunkKV NeurIPS 2025 (semantic-chunk eviction).

## Hop-2 Citations

- RestoreKV (2026 arXiv:2608.01247 — learned cache restore).
- KVzap (NVIDIA surrogate MLP, 2026).
- Anthropic prompt caching (2025, 90% input-cost reduction baseline).

## MAgHARCM integration

- **YAML config key**: `kv_cache.repair.enabled: true` (opt-in, default `false`); `kv_cache.repair.tier: cpu|nvme`.
- **Implementation file**: `internal/kvrepair/repair.go::NewRepairOperator` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **Tier cost** — slower-tier re-evaluation adds latency at agent-handoff time; cost is amortised but must be budgeted.
- **Repair policy** — heuristic repair may over-promote; learned repair requires labelled data. Threshold tuning is required per deployment.
- **Workshop venue** — ICML 2026 AdaptFM workshop-track ACCEPT relies on §7 method-level threshold; main-track companion would strengthen the anchor.

## Source

- OpenReview: LsrmZrp7tW.
- Venue: ICML 2026 AdaptFM Workshop (Poster). arXiv forthcoming.
