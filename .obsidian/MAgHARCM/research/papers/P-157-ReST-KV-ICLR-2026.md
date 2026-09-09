---
title: "P-157 ReST-KV: Robust KV Cache Eviction with Layer-wise Output Reconstruction and Spatial-Temporal Smoothing"
backlink: "[[1.0.0 P-157]]"
aliases:
 - "1.0.0 P-157"
 - "P-157"
 - "P-157-ReST-KV-ICLR-2026"
 - "ReST-KV-ICLR-2026"
 - "restkv2026iclr"
tags: [paper, kv-cache-eviction, output-reconstruction, spatial-temporal-smoothing, slm-era, "[[1.0.0 PRIM-21]]", "[[1.0.0 PRIM-31]]", wave-24]
date: 2026-09-09
last_updated: 2026-09-09 (iter-1, wave-24)
venue: ICLR 2026 (Poster)
---

# [[1.0.0 P-157]] ReST-KV

## TL;DR

ReST-KV is a **training-free KV-cache eviction** method that replaces raw-attention importance scoring with **layer-wise output reconstruction**: model how each candidate token's removal perturbs the layer output, evict the minimum-discrepancy KV pairs, and smooth importance spatial-temporally (EMA + adaptive windows). Anchors `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (eviction fidelity governs what retained context the retrieval loop sees) and `[[1.0.0 PRIM-21]]` Migration Strategy Selection (reconstruction-based vs attention-heuristic = eviction-strategy knob alongside P-147 SpecKV / P-148 LookaheadKV).

## Mechanism (Q2)

1. **Layer-wise output reconstruction** — eviction is formulated as an optimization problem that minimizes output discrepancy. Instead of ranking by attention weight magnitude, ReST-KV directly models how each token's removal affects the model output, naturally capturing attention-redistribution effects (removing a token forces the model to re-allocate attention to remaining tokens, which raw scoring misses).
2. **Spatial-temporal smoothing** — exponential-moving-average smoothing handles temporal variation in token importance across decoding steps; an adaptive window mechanism captures spatial patterns across the token sequence.
3. **Efficiency** — +2.58% LongBench, +15.2% RULER over state-of-the-art eviction baselines; 10.61x decoding-latency reduction at 128k context; consistent wins on Needle-in-a-Haystack and InfiniteBench.

## Anchoring (Q3)

| Primitive | Pre-wave-24 behaviour | ReST-KV substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Retrieval refinement anchored on P-117 RepoCoder, P-128 KVzip, P-134 RelayCaching, P-138 RepairKV, P-141 KVFlow, P-150 TestPrune, P-153 CoReX; eviction quality treated as a black box | ReST-KV output-reconstruction-governed eviction determines what retained context the retrieval loop sees; eviction fidelity is a retrieval-refinement precondition, not a downstream detail |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Strategy selection anchored on P-122 ReasoningBank, P-123 CodeChemist, P-126 ARC-Decode, P-135 SPECS, P-137 SuffixDecoding, P-145 TerraMod, P-147 SpecKV, P-148 LookaheadKV, P-149 SSD/Saguaro | ReST-KV adds the reconstruction-objective eviction strategy to the KV-cache strategy family: choose output-reconstruction scoring when attention redistribution dominates, attention-heuristic when eviction budget is tight |

## Hop-1 Citations

- P-147 SpecKV (ICLR 2026, Wave-21) — draft-model-driven eviction. ReST-KV shares the venue family; replaces lookahead scoring with reconstruction optimization.
- P-148 LookaheadKV (ICLR 2026, Wave-21) — LoRA-augmented lookahead. ReST-KV is the training-free alternative.
- P-134 RelayCaching (ICML 2026, Wave-19) — cross-agent KV reuse. ReST-KV governs per-layer eviction quality that determines what the reuse substrate retains.
- P-141 KVFlow (Wave-19) — workflow-aware eviction; ReST-KV supplies the output-fidelity objective.

## Hop-2 Citations

- H2O / heavy-hitter oracle eviction (Zhang 2023) — raw-attention-importance baseline family shown to be blind to attention redistribution.
- Optimal Brain Surgeon lineage (LeCun 1990; SparseGPT, Frantar & Alistarh 2023) — remove-and-reconstruct objective transferred from weight pruning to KV eviction.

## MAgHARCM integration

- **YAML config key**: `agents.kv_cache.eviction.strategy: restkv` (joins `speckv` / `lookaheadkv` from Wave-21 gates).
- **Implementation file**: none yet (KV-cache substrate gates remain disabled pending BLK-04 re-run decision).
- **Affected primitives**: `[[1.0.0 PRIM-31]]`, `[[1.0.0 PRIM-21]]`.

## Caveats

- Reconstruction scoring adds per-layer compute during prefill; the paper reports net latency wins at 128k context but does not profile SLM-scale (4B-30B) prefill overhead.
- LongBench/RULER are natural-language benchmarks; code-translation context distributions differ.
- No multi-agent workflow evaluation; KVFlow-style workflow signals are untested in combination.

## Source

- Venue: ICLR 2026 (Poster), verified via iclr.cc/virtual/2026/poster/10009650 and OpenReview forum PhEHuo7oMm.
- arXiv: 2605.08840.
- Authors' affiliation: Institute of Automation, Chinese Academy of Sciences (CASIA) et al.
- Code: github.com/an-yongqi/rest-kv.

## BibTeX

```
@inproceedings{restkv2026iclr,
  title     = {ReST-KV: Robust KV Cache Eviction with Layer-wise Output Reconstruction and Spatial-Temporal Smoothing},
  author    = {An, Yongqi and Lu, Chang and Zhu, Kuan and Yu, Tao and Zhao, Chaoyang and Wu, Hong and Tang, Ming and Wang, Jinqiao},
  booktitle = {Proceedings of the Fourteenth International Conference on Learning Representations (ICLR)},
  year      = {2026}
}
```
