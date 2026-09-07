---
title: "P-147 SpecKV: Draft-based Approximate Inference for LLMs"
backlink: "[[1.0.0 P-147]]"
aliases:
 - "1.0.0 P-147"
 - "P-147"
 - "P-147-SpecKV-ICLR-2026"
 - "SpecKV-ICLR-2026"
 - "speckv2026iclr"
tags: [paper, kv-cache, speculative-decoding, slm, eviction, "[[1.0.0 PRIM-21]]", "[[1.0.0 PRIM-31]]", wave-21]
date: 2026-09-07
last_updated: 2026-09-07 (iter-4, wave-21)
venue: ICLR 2026 (Conference)
---

# [[1.0.0 P-147]] SpecKV

## TL;DR

SpecKV is a **draft-model-driven KV cache eviction** strategy in which a small draft model performs lookahead to predict which KV pairs will be most important, then evicts low-importance pairs before the target model attends. Bundled with SpecPC (prompt compression) and SpecKV-PC (cascaded strategy combining prompt + KV compression with an adaptive gamma controller). Reduces KV cache memory pressure and time-to-first-token for long-context agentic loops while preserving generation quality. Anchors `[[1.0.0 PRIM-21]]` Migration Strategy Selection (adaptive gamma controller as a strategy-selection knob between draft-KV-eviction-cost and compressed-context-quality) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (draft-model lookahead drives eviction importance for long-context retrieval loops).

## Mechanism (Q2)

1. **Draft-model lookahead** — a small draft model performs a cheap lookahead pass to predict which KV pairs will be most important for upcoming target-model attention.
2. **Importance-weighted eviction** — KV pairs ranked lowest by the draft-model lookahead are evicted; the target model attends over the remaining (importance-weighted) KV cache.
3. **SpecPC prompt compression** — companion mechanism that compresses the prompt using a draft-model-driven signal.
4. **SpecKV-PC cascaded strategy** — combines prompt + KV compression with an **adaptive gamma controller** that tunes the budget allocation between the two based on observed downstream quality.
5. **Adaptive gamma** — the controller trades off draft-KV-eviction-cost (gamma high → evict more) against compressed-context-quality (gamma low → preserve more KV).

Together these give MAgHARCM a **draft-model-driven importance signal** for KV eviction that complements P-128 KVzip (query-agnostic context reconstruction) and P-141 KVFlow (workflow-aware eviction) for multi-step agentic loops.

## Anchoring (Q3)

| Primitive | Pre-wave-21 behaviour | SpecKV substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Strategy registry used prompt size + verifier cost as signals; no eviction-vs-quality budget knob | SpecKV's adaptive gamma controller = budget knob between draft-KV-eviction-cost and compressed-context-quality for SLM-scale long-context strategies |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Retrieval loop used StreamingLLM attention sinks / KVzip query-agnostic reconstruction; no draft-model-driven importance signal | SpecKV's draft-model lookahead = eviction importance signal that adapts per-context rather than per-query |

## Hop-1 Citations

- P-128 KVzip (NeurIPS 2025, Kim et al.) — query-agnostic KV cache compression with context reconstruction; complementary approach (compression rather than eviction) for the same long-context SLM substrate.
- P-141 KVFlow (NeurIPS 2025 Poster) — workflow-aware KV cache eviction for multi-agent serving; orthogonal axis (workflow signal vs. draft-model signal).
- P-57 Leviathan et al. (NeurIPS 2023) — original speculative decoding; SpecKV's draft-model lookahead is a KV-cache analogue.
- P-108 EAGLE-3 (NeurIPS 2025) — training-time-test draft model for SLM speculative decoding; SpecKV's draft-model lookahead is KV-cache-focused but borrows the "use a small model to predict the big model's behaviour" pattern.

## Hop-2 Citations

- Speculative decoding (Leviathan 2023, Chen 2023) — original drafter-verifier pattern.
- KV cache compression (Pope et al. 2023 Efficient Memory Management) — foundational KV cache management.
- Long-context attention sinks (Xiao et al. 2024 StreamingLLM, P-80) — alternative KV-cache substrate for SLM-scale long-context.
- Adaptive budget allocation (Snell et al. P-91 test-time compute allocation) — pattern of a controller tuning between two cost-quality axes.

## MAgHARCM integration

- **YAML config key**: `agents.kv_cache.eviction.strategy: speckv`; `agents.kv_cache.eviction.lookahead_tokens: <int>`; `agents.kv_cache.compression.gamma: <float>`.
- **Implementation file**: `internal/iter_retrieval/kv_eviction.go::NewSpecKVEviction` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **Draft-model cost** — the draft-model lookahead pass adds compute before the target model's first token; cost must be amortized over long enough generations to pay off.
- **Importance signal quality** — draft-model-importance is an approximation; evicted pairs may be more important than predicted (no second-chance correction like P-138 RepairKV).
- **Adaptive gamma tuning** — the gamma controller is itself a learning signal; cold-start needs a default heuristic.
- **Composition with prompt compression** — SpecKV-PC combines prompt + KV compression, but the cascaded composition can compound errors if both signals are wrong.

## Source

- arXiv: 2506.08373.
- Venue: ICLR 2026 (Conference), verified via OpenReview 0vbYakkECY.

## BibTeX

```
@inproceedings{speckv2026iclr,
  title  = {SpecKV: Draft-based Approximate Inference for LLMs},
  author = {Galim, Kevin and Wang, Yuxin and Wu, Fei and Lee, Sungho and Du, Nan and Cai, Yi},
  booktitle = {International Conference on Learning Representations (ICLR)},
  year   = {2026},
  eprint = {arXiv:2506.08373},
  url    = {https://openreview.net/forum?id=0vbYakkECY}
}
```
