---
title: Primitives Index
backlink: [[Primitives]]
tags: [primitives, catalog, status]
---

# Primitives Index

A primitive is a small, named capability the harness reuses across the
pipeline. Each entry links to the file that implements it and notes its
status (active, deprecated, experimental).

## Active

| ID | Name | Where | Status |
| -- | ---- | ----- | ------ |
| PRIM-1 | Strategy registry (try-and-fail) | `internal/agents/strategy.go` | active |
| PRIM-3 | Chunked translator | `internal/agents/chunked_translator.go` | active |
| PRIM-4 | Validator repair loop | `internal/agents/validator.go` | active |
| PRIM-5 | Bubble Tea TUI | `internal/tui/tui.go` | active |
| PRIM-6 | AST pre-compile check | `internal/agents/validator.go::scanASTSyntax` | active |
| PRIM-7 | Checkpoint resume | `internal/agents/checkpoint.go` | active |
| PRIM-13 | Adversarial weakening detector | `internal/agents/validator.go` | active |
| PRIM-27 | Plateau-bounded remedy loop (CodaMOSA) | `internal/agents/plateau.go` | active |

## Newly Added (this refactor)

- **PRIM-1 (Strategy registry)** — replaced the 4-branch if-cascade in
  `SelectMigrationStrategy`. Adding a strategy is one line; on a strategy
  failure the registry tries the next entry.
- **PRIM-5 (Bubble Tea TUI)** — replaces the ad-hoc stdin scanner with a
  Bubble Tea Model. The exported command surface (`Phase1Step`, `HandleSlash`,
  `Phase1Help`, `slashHelp`) is preserved so the test mirror stays in sync.
- **PRIM-27 (Plateau detector)** — folded into the validator cascade as the
  final check before the repair back-edge. Halts the repair loop when the
  coverage delta across the last N iterations is below threshold.

## Removed

- **PRIM-2 (Navigator agent)** — was a 5th parallel graph node. Removed:
  symbol-name resolution is now a sub-mechanism of the Planner agent
  invoking the LSP provider directly. The figure `magh-pipeline.workflow.json`
  reflects this.

## Package Cohesion

- **`internal/artifacts`** — output data shapes; no upstream imports.
- **`internal/types`** — `State` + `TranslationTask`.
- **`internal/consts`** — compile-time constants (error values, toolchain
  names).

See [[Architecture]] for the dependency DAG and [[METHODOLOGY]] for the
runtime pipeline.
