---
title: MAgHARCM Vault
tags: [vault, index]
---

# MAgHARCM Vault

This is the Obsidian vault for the MAgHARCM project. It documents the
architecture, methodology, primitives, and research notes.

## Quick links

- [[research/METHODOLOGY]] — the 6-stage pipeline, strategy registry, translator modes, validator loop, state container, checkpointing.
- [[research/Architecture]] — package layout and dependency DAG.
- [[primitives/INDEX]] — catalog of named capabilities (PRIM-1 through PRIM-27).

## Structure

```
.obsidian/MAgHARCM/
├── research/
│   ├── METHODOLOGY.md   — pipeline, strategies, state, checkpointing
│   └── Architecture.md  — package DAG, why three packages for state
├── primitives/
│   └── INDEX.md         — catalog of PRIM-NN entries
└── diary/               — session notes (optional)
```

## Editing

Use Obsidian's link syntax `[[note name]]` to connect related pages. The
vault is stored in the project root so it stays versioned with the code.
