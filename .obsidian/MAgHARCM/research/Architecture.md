---
title: Architecture
backlink: [[Architecture]]
tags: [architecture, package-graph]
---

# Architecture

Package layout and dependency direction. Arrows point from importer to
importee. Pipeline spine per [[PRIM-1]], [[PRIM-2]], [[PRIM-3]], [[PRIM-23]],
[[PRIM-6]] and repair loop per [[PRIM-27]], [[PRIM-13]].

```
internal/artifacts   (no internal imports; owns the output data shapes,
                      versioned per [[PRIM-24]])
       ^
       |
internal/types       (State + TranslationTask; embeds artifacts.*)
       ^
       |
internal/agents      (4 spine agents + auxiliary primitives
                      [[PRIM-1]]..[[PRIM-30]] per [[primitives/INDEX]])
       |
internal/graph       (Eino graph wiring the 4 agents + 2 checkpoint lambdas)
       |
internal/runner      (orchestrates Run -> graph.Compile -> Run; CLI/UI entry point)
```

Cross-cutting:

- `internal/llm` — Ollama client; imported by `runner` and `graph`.
- `internal/tools` — exec, fs, LSP providers, prompt archive; imported
  by `agents`.
- `internal/languages` — language registry (C, Rust, Go, Java, Python);
  imported by `tools`.
- `internal/logger` — structured logger + ring buffer; imported
  everywhere.
- `internal/config` — YAML config loader with hard-fail `Require` check;
  imported by `runner`.

## Why a strategy registry ([[PRIM-21]])

The first cut of `SelectMigrationStrategy` was a 4-branch if-cascade with
hardcoded strategy names. Adding a strategy meant editing the function
and risking a regression on the existing branches.

The registry (`internal/agents/strategy.go`) is a slice of named
strategies with `Matches(Profile) bool` and `Attempt(*State) error`.
Adding a strategy is a one-line append to `NewDefaultRegistry`. The
analyzer iterates the slice via `Registry.TryInOrder` and picks the first
match; on a strategy failure it tries the next one rather than failing
the whole run.

## Why four agents

The agent decomposition mirrors [[ReCodeAgent-2024]]'s Reasoning-Coding
split: an 8B reasoning model handles analysis and planning where longer
context and tool use pay off, and a small 4B coding model handles the
per-fragment translation where latency and determinism dominate.

- **Analyzer** (reasoning) — emits Source Project Research, Library
  Mapping, Target Project Design. Selects strategy via [[PRIM-21]].
- **Planner** (reasoning) — fragments AST ([[PRIM-2]]), resolves
  ambiguous symbol names ([[PRIM-26]]), builds reverse-topological DAG
  ([[PRIM-1]]), emits skeleton ([[PRIM-3]]) + Part A + Part B
  implementation plan.
- **Translator** (coding) — executes the plan in topological order,
  replacing skeleton stubs with concrete target code ([[PRIM-23]]). Each
  fragment sees bounded context ([[PRIM-26]] + [[PRIM-30]]).
- **Validator** (coding) — runs AST pre-check ([[PRIM-6]]), compiler
  build, test suite, weakening guard ([[PRIM-13]]), plateau detector
  ([[PRIM-27]]); feeds diagnostics back to the Translator.

Auxiliary primitives plug into the spine:

- [[PRIM-7]] Multi-Agent Verdict Validation — optional repair-loop verifier
- [[PRIM-8]] State-Grounded Mock-Based In-Isolation Validation —
  per-module test
- [[PRIM-9]] Tri-Representation Hybrid Code Graph — feeds [[PRIM-26]]
- [[PRIM-11]] Implementation-Agnostic Testing — black-box validation
- [[PRIM-12]] Wasm-Based Reference Execution Oracle — reference oracle
- [[PRIM-25]] Communicative-De-hallucination Role-Flip Gate — translator
  output gate
- [[PRIM-29]] Recruitment-Adaptive Planning — per-iteration agent/tool
  selection

See [[METHODOLOGY]] for the runtime behaviour.
