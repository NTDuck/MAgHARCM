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
internal/agents      (4 agents: analyzer, planner, translator, validator)
       |
internal/graph       (Eino graph wiring the 4 agents + 2 checkpoint lambdas)
       |
internal/runner      (orchestrates Run -> graph.Compile -> Run; CLI/UI entry point)
```

Cross-cutting:

- `internal/llm` — Ollama client; imported by `runner` and `graph`.
- `internal/tools` — exec, fs, LSP providers, prompt archive; imported by `agents`.
- `internal/languages` — language registry (C, Rust, Go, Java, Python); imported by `tools`.
- `internal/logger` — structured logger + ring buffer; imported everywhere.
- `internal/config` — YAML config loader with hard-fail `Require` check; imported by `runner`.

## Why a strategy registry

The first cut of `SelectMigrationStrategy` was a 4-branch if-cascade with
hardcoded strategy names. Adding a strategy meant editing the function and
risking a regression on the existing branches.

The registry (`internal/agents/strategy.go`) is a slice of named strategies
with `Matches(Profile) bool` and `Attempt(*State) error`. Adding a strategy
is a one-line append to `NewDefaultRegistry`. The analyzer iterates the
slice via `Registry.TryInOrder` and picks the first match; on a strategy
failure it tries the next one rather than failing the whole run.

## Why four agents

The agent decomposition mirrors ReCodeAgent's Reasoning-Coding split: a
30B reasoning model handles analysis and planning where longer context and
tool use pay off, and a small 4B coding model handles the per-fragment
translation where latency and determinism dominate.

- **Analyzer** (reasoning) — emits Source Project Research, Library Mapping,
  Target Project Design.
- **Planner** (reasoning) — fragments AST, resolves ambiguous symbol names,
  builds reverse-topological DAG, emits Part A + Part B implementation plan.
- **Translator** (coding) — executes the plan in topological order, replacing
  skeleton stubs with concrete target code. Each fragment is synthesised with
  bounded context (Symbol Navigator + Prior-Modules Memory).
- **Validator** (coding) — runs AST pre-check, compiler build, test suite,
  weakening guard, plateau detector; feeds diagnostics back to the Translator.

See [[METHODOLOGY]] for the runtime behaviour.
