---
title: Architecture
backlink: [[Architecture]]
tags: [architecture, package-graph]
---

# Architecture

Package layout and dependency direction. Arrows point from importer to
importee.

```
internal/artifacts   (no internal imports; owns the output data shapes)
       ^
       |
internal/types       (State + TranslationTask; embeds artifacts.*)
       ^
       |
internal/agents      (5 agents: analyzer, navigator, planning, translator, validator)
       |
internal/graph       (Eino graph wiring the 5 agents + 2 checkpoint lambdas)
       |
internal/runner      (orchestrates Run -> graph.Compile -> Run; CLI/UI entry point)
```

Cross-cutting:

- `internal/llm` — Ollama client; imported by `runner` and `graph`.
- `internal/tools` — exec, fs, LSP providers, prompt archive; imported by `agents`.
- `internal/languages` — language registry (C, Rust, Go, Java, Python); imported by `tools`.
- `internal/logger` — structured logger + ring buffer; imported everywhere.
- `internal/config` — YAML config loader with hard-fail Require check; imported by `runner`.

## Why three packages for state

A naive split would put `AnalyzerOutput` in `internal/agents/analyzer`,
`PlanningOutput` in `internal/agents/planning`, etc. — but `internal/types.State`
embeds all of them. Moving them out would force `types` to import `agents`,
which already imports `types` (cycle).

The `internal/artifacts` package owns the output shapes, has zero upstream
imports, and is imported by both `types` and `agents`. No cycle forms.

## Why a strategy registry

The first cut of `SelectMigrationStrategy` was a 4-branch if-cascade with
hardcoded strategy names. Adding a strategy meant editing the function and
risking a regression on the existing branches.

The registry (`internal/agents/strategy.go`) is a slice of named strategies
with `Matches(Profile) bool` and `Attempt(*State) error`. Adding a strategy
is a one-line append to `NewDefaultRegistry`. The analyzer iterates the
slice via `Registry.TryInOrder` and picks the first match.

## Why a fifth agent

The original graph had four agents. The navigator was added as a parallel
lookup node between analyzer and planning because some repositories have
name collisions across modules that the analyzer alone cannot resolve. The
navigator calls `LSPProvider.GetDefinition` and `GetReferences` to confirm
the target of each ambiguous symbol before the planner fragments the AST.

When no LSP provider is wired, the navigator is a no-op. It still appears
in the graph so the wiring is consistent.

See [[METHODOLOGY]] for the runtime behaviour.
