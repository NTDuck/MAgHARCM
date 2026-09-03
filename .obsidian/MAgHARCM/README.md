---
title: MAgHARCM Vault
tags: [vault, index]
---

# MAgHARCM Vault

This is the Obsidian vault for the MAgHARCM project. It documents the
architecture, methodology, primitives, and research notes.

## Quick links

- [[research/METHODOLOGY]] — the 4-agent pipeline, strategy registry, chunked translator, validator cascade, state container, checkpointing.
- [[research/Architecture]] — package layout and dependency DAG.
- [[primitives/INDEX]] — catalog of named capabilities.

## Structure

```
.obsidian/MAgHARCM/
├── research/
│   ├── METHODOLOGY.md   — pipeline, strategies, validator cascade, state, checkpointing
│   └── Architecture.md  — package DAG, 4-agent decomposition
├── primitives/
│   └── INDEX.md         — catalog of PRIM-NN entries
└── diary/               — session notes (optional)
```

## Editing

Use Obsidian's link syntax `[[note name]]` to connect related pages. The
vault is stored in the project root so it stays versioned with the code.
