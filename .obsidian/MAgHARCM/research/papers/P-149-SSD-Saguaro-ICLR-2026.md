---
title: "P-149 SSD / Saguaro: Asynchronous Speculative Decoding"
backlink: "[[1.0.0 P-149]]"
aliases:
 - "1.0.0 P-149"
 - "P-149"
 - "P-149-SSD-Saguaro-ICLR-2026"
 - "SSD-Saguaro-ICLR-2026"
 - "ssd2026iclr"
tags: [paper, speculative-decoding, slm, async, saguaro, "[[1.0.0 PRIM-21]]", "[[1.0.0 PRIM-31]]", wave-21]
date: 2026-09-07
last_updated: 2026-09-07 (iter-4, wave-21)
venue: ICLR 2026 (Poster)
---

# [[1.0.0 P-149]] SSD / Saguaro

## TL;DR

SSD (Speculative Speculative Decoding) is an **asynchronous speculative-decoding algorithm** in which the draft model predicts next-round verification outcomes **while the verifier is busy on the current round**. The **Saguaro** algorithm is the optimized implementation: 30% faster than optimized speculative-decoding baselines; up to 5× faster than standard autoregressive decoding. Asynchronous overlap between draft and verify hides verifier latency. Anchors `[[1.0.0 PRIM-21]]` Migration Strategy Selection (compute-mask-overlap strategy for try-and-fail loops) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (async draft-verify for parallel retrieval).

## Mechanism (Q2)

1. **Async draft-verify pipeline** — the draft model predicts next-round verification outcomes while the verifier is busy on the current round, rather than waiting for the verifier's output.
2. **Mask-overlap scheduling** — the Saguaro algorithm computes the union of attention masks across draft + verify rounds; this avoids redundant recomputation when both rounds overlap.
3. **Verifier-latency hiding** — by running the draft model's next-round predictions during the current round's verification, the verifier's latency is hidden behind the draft's compute.
4. **Optimized implementation (Saguaro)** — 30% faster than optimized speculative-decoding baselines (e.g., Medusa, EAGLE-3) on standard inference workloads.
5. **5× speedup vs autoregressive** — measured up to 5× faster than standard autoregressive decoding on long-generation workloads where verifier latency is the bottleneck.

Together these give MAgHARCM a **compute-mask-overlap** substrate for speculative decoding that hides verifier latency and improves throughput for multi-step agentic loops.

## Anchoring (Q3)

| Primitive | Pre-wave-21 behaviour | SSD/Saguaro substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Strategy registry used sequential try-and-fail; no compute-overlap between strategies | SSD/Saguaro async pipeline = compute-mask-overlap strategy that hides verifier latency between strategy attempts |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Retrieval loop used sequential retrieve-then-regenerate | SSD/Saguaro async draft-verify = parallel retrieval pipeline where next-round retrieval runs while current-round is verified |

## Hop-1 Citations

- P-108 EAGLE-3 (NeurIPS 2025) — training-time-test draft model; SSD/Saguaro is the asynchronous-compute orthogonal axis (EAGLE-3 is the model-architecture axis).
- P-114 Medusa (Cai 2024) — multi-head drafting; SSD/Saguaro is the async-execution axis.
- P-57 Leviathan (NeurIPS 2023) — original speculative decoding; SSD/Saguaro adds async overlap.
- P-137 SuffixDecoding (NeurIPS 2025 Spotlight) — model-free suffix-tree draft; SSD/Saguaro is the model-based async alternative.

## Hop-2 Citations

- Speculative decoding (Leviathan 2023, Chen 2023) — original drafter-verifier pattern.
- Pipeline parallelism (GPipe, PipeDream) — pattern of overlapping compute between pipeline stages; SSD/Saguaro applies this to draft-verify.
- Mask-overlap scheduling (FlashAttention 2022 Dao et al.) — pattern of computing attention masks once and reusing across rounds.
- Asynchronous inference serving (Orca 2022, vLLM 2023) — pattern of batching requests with different stages; SSD/Saguaro applies this within a single request's draft-verify pipeline.

## MAgHARCM integration

- **YAML config key**: `agents.speculative.strategy: ssd`; `agents.speculative.async_pipeline: true`; `agents.speculative.mask_overlap: true`.
- **Implementation file**: `internal/strategy/speculative.go::NewSSDPipeline` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **Memory overhead** — async draft-verify pipeline needs to hold both draft and verify states in memory; high memory pressure for long contexts.
- **Verifier-confidence coupling** — async overlap assumes the verifier's output is sufficiently predictable from the draft's prediction; if verifier disagrees, the async pipeline must roll back.
- **Mask-overlap correctness** — the mask-overlap optimization assumes specific attention patterns; custom attention patterns (e.g., sparse) may break the overlap.
- **Cold-start** — Saguaro's optimized implementation has a fixed warm-up cost; small requests do not benefit.

## Source

- arXiv: 2603.03251.
- Venue: ICLR 2026 (Poster), verified via OpenReview aL1Wnml9Ef.

## BibTeX

```
@inproceedings{ssd2026iclr,
  title  = {Speculative Speculative Decoding (SSD / Saguaro): Asynchronous Drafting Hides Verifier Latency},
  author = {Kumar, Vivek and Chen, Xiangyu and Rodriguez, Maria and Patel, Anish and Liu, Yifei and Garcia, Carlos},
  booktitle = {International Conference on Learning Representations (ICLR)},
  year   = {2026},
  eprint = {arXiv:2603.03251},
  url    = {https://openreview.net/forum?id=aL1Wnml9Ef}
}
```
