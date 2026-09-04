---
title: MAgHARCM Vault
backlink: [[2.0.0 MAgHARCM Vault]]
tags: [vault, index, [[2.0.0 MAgHARCM]]]
---

# [[2.0.0 MAgHARCM Vault]]

This is the canonical Obsidian vault for the MAgHARCM project. It documents the
architecture, methodology, primitives, and software-archaeology research notes.

`[[2.0.0 MAgHARCM]]` is the current vault versioning.
The `[[2.0.0 Sprint Modernization & Architectural Convergence]]` entry captures the
full parity convergence where all 31 primitives are active implementations in the codebase.

## Quick Links

- `[[2.0.0 Methodology]]` — the 8-agent pipeline, dynamic try-and-fail strategy registry
  (`[[1.0.0 PRIM-21]]`), chunked translator (`[[1.0.0 PRIM-23]]`), validator cascade
  (`[[1.0.0 PRIM-6]]`, `[[1.0.0 PRIM-13]]`, `[[1.0.0 PRIM-27]]`), state container (`[[1.0.0 PRIM-24]]`),
  and durable checkpointing (`[[1.0.0 PRIM-28]]`).
- `[[2.0.0 Architecture]]` — package layout, dependency DAG, and Locality of Behaviour.
- `[[2.0.0 Primitives Index]]` — complete catalog of all 31 primitives `[[1.0.0 PRIM-1]]` through `[[1.0.0 PRIM-31]]`.
- `[[2.0.0 Software-Archaeology-Lineage]]` — 2-hop foundational literature synthesis (Parnas, Lehman, Chikofsky, Rajlich, Müller, Baldwin & Clark, Feathers, Kazman, Foltz).

## Structure

```
.obsidian/MAgHARCM/
├── research/
│   ├── METHODOLOGY.md            — pipeline, strategies, validator cascade, state
│   ├── Architecture.md           — package DAG, 8-agent decomposition, Locality of Behaviour
│   ├── Software-Archaeology-Lineage.md — 2-hop literature lineage & synthesis
│   └── papers/                   — bibliography research P-01 through P-30
├── primitives/
│   └── INDEX.md                  — complete catalog of all 31 primitives
└── diary/                        — session and sprint tracking notes
```

## Versioning Convention

All version markers use the `[[x.y.z ...]]` convention. Spontaneous parentheses
are forbidden; replace them with `[[Author-YEAR]]` or `[[P-NN]]` style backlinks.
Every primitive is implemented in the codebase, establishing a single source of truth.
