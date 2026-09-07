---
title: "P-130 SSAR: A Novel Software Architecture Recovery Approach Enhancing Accuracy and Scalability"
backlink: "[[1.0.0 P-130]]"
aliases:
  - "1.0.0 P-130"
  - "P-130"
  - "P-130-SSAR-Architecture-Recovery-ICSE-2026"
  - "SSAR-Architecture-Recovery-ICSE-2026"
  - "ding2026ssar"
tags: [paper, architecture-recovery, code-graph, community-detection, modularization, "[[1.0.0 PRIM-9]]", "[[1.0.0 PRIM-22]]", wave-18]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICSE 2026 (Research Track)
---

# [[1.0.0 P-130]] SSAR

## TL;DR

SSAR builds a **weighted file-graph** whose edges combine semantic similarity (LLM/SBERT-style embeddings of file content) with structural dependencies (call/import graph), then applies an optimized community-detection partition to recover software modules. Achieves **5.0-90.9% accuracy gains and 5-99% runtime reduction** over six SOTA baselines on nine ground-truth projects (MoJoFM, a2a, c2a metrics). Anchors `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph and `[[1.0.0 PRIM-22]]` Four Phases of Comprehension.

## Mechanism (Q2)

The pipeline has three stages:

1. **Dual-edge graph construction** — for every file pair, compute two edge weights: (a) **semantic similarity** from LLM or SBERT embeddings of the file content; (b) **structural dependency** from the call/import graph (number of shared callees, dependency depth, etc.). The final edge weight is a learned or weighted-sum combination of the two.
2. **Community detection** — apply an optimized community-detection algorithm (Louvain / Leiden / spectral variant) on the weighted file graph to partition files into modules. The combination of semantic + structural edges biases the partition toward architecturally meaningful boundaries.
3. **Scalability optimizations** — caching, edge-pruning, and incremental updates reduce wall-clock time, enabling recovery on large repos that scale-fail pure-semantic approaches.

The dual-edge design is the key insight: structural-only recovery misses hidden coupling (e.g. files that share a domain concept but no direct calls); semantic-only recovery is noisy and slow. SSAR's combination yields both better boundaries and faster runtime.

## Anchoring (Q3)

| Primitive | Pre-wave-18 behaviour | SSAR substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Single-edge graphs: either structural (call graph) or semantic (embeddings); structural misses hidden coupling, semantic is noisy and slow | Dual-edge file graph combining semantic similarity + structural dependencies; optimized community detection |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Architecture-recovery phase relies on single-modality signals | Multi-modal file-level representation (semantic + structural) enriches the comprehension phase |

## Hop-1 Citations

- MoJoFM (van Emden et al. 2003) — modularization quality metric.
- a2a / c2cavg / MoJoFM (Tzerpos et al. 1999) — architecture recovery benchmark metrics.
- Louvain (Blondel et al. 2008) and Leiden (Traag et al. 2019) — community detection algorithms.
- SBERT (Reimers et al. 2019) — embedding baseline.

## Hop-2 Citations

- Mancoridis et al. 1999 *Bunch* — clustering-based architecture recovery (early ancestor).
- LLM-based code understanding lineage (Ma et al. 2024, Wang et al. 2024).

## MAgHARCM integration

- **YAML config key**: `architecture_recovery.ssar.enabled: true` (opt-in, default `false`); `architecture_recovery.ssar.semantic_weight: <float>`.
- **Implementation file**: `internal/code_graph/ssar.go::RecoverArchitecture` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-9]]`, `[[1.0.0 PRIM-22]]`.

## Caveats

- **Embedding cost**: semantic similarity requires an embedding pass over all files; very large repos (10k+ files) need incremental embedding or sampling to keep wall-clock tractable.
- **Metric dependency**: MoJoFM / a2a / c2cavg are sensitive to ground-truth granularity; evaluation should report multiple metrics rather than relying on a single score.
- **Domain-language coverage**: validated on nine projects; languages with weak call-graph extraction (dynamic dispatch, reflection-heavy code) may degrade structural-edge quality.

## Source

- URL: conf.researchr.org/details/icse-2026/icse-2026-research-track/221.
