---
title: MAgHARCM Vault
tags: [vault, index, [[v0.2.0]]]
---

# MAgHARCM Vault

This is the Obsidian vault for the MAgHARCM project. It documents the
architecture, methodology, primitives, and research notes.

[[v0.2.0]] is the current vault versioning (was [[v0.1.0]] before this
sprint). The [[Sprint-Recon-2026-09-04]] entry is the parity audit that
drove this sprint's edits.

## Quick links

- [[research/METHODOLOGY]] — the 4-agent pipeline, strategy registry
  ([[PRIM-21]]), chunked translator ([[PRIM-23]]), validator cascade
  ([[PRIM-6]] [[PRIM-13]] [[PRIM-27]]), state container ([[PRIM-24]]),
  checkpointing ([[PRIM-28]]).
- [[research/Architecture]] — package layout and dependency DAG.
- [[primitives/INDEX]] — catalog of named capabilities [[PRIM-1]]..[[PRIM-30]].

## Structure

```
.obsidian/MAgHARCM/
├── research/
│   ├── METHODOLOGY.md            — pipeline, strategies, validator cascade,
│   │                              state, checkpointing
│   ├── Architecture.md           — package DAG, 4-agent decomposition
│   └── papers/                   — bibliography research (recursive)
├── primitives/
│   └── INDEX.md                  — catalog of [[PRIM-1]]..[[PRIM-30]] entries
└── diary/                        — session notes (sprint recon, etc.)
```

## Editing

Use Obsidian's link syntax `[[note name]]` to connect related pages. The
vault is stored in the project root so it stays versioned with the code.

## Versioning Convention

Use `[[x.y.z ...]]` for every ID/version marker. Spontaneous parentheses
(`(CodaMOSA)`, `(Syzygy)`, `(ReCodeAgent)`) are NOT allowed; replace with
`[[Author-YEAR]]` or `[[P-NN]]` style backlinks.
