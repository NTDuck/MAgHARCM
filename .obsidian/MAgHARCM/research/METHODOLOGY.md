---
title: MAgHARCM Methodology
backlink: [[Methodology]]
tags: [methodology, architecture, pipeline]
---

# MAgHARCM Methodology

## Pipeline

Six stages, wired as Eino graph nodes. Each node reads from the shared
`types.State` and writes back into it.

```
START → analyzer → navigator → planning → translator → save_translator_ckpt
      → validator → save_validator_ckpt → branch
```

| Stage | Agent | What it does |
| ----- | ----- | ------------ |
| 1     | analyzer     | Walks the source tree, runs the analyzer reasoning model, and emits three documents: project research, library mapping, and target design. Picks a migration strategy via `internal/agents/strategy.go::Registry.TryInOrder`. |
| 2     | navigator    | Resolves ambiguous symbols in `PlanningOutput.NameMapping` by querying the configured LSP provider (native tree-sitter or ABCoder MCP). No-op when no provider is wired. |
| 3     | planner      | Reads analyzer output, fragments the AST into translation units, and emits an `ImplementationPlan` with two ordered lists: source translation steps (Part A) and test generation + validation steps (Part B). |
| 4     | translator   | Runs the coding model per fragment in topological order. Persists a checkpoint after each fragment. Re-uses the previously emitted modules as context for later fragments. |
| 5     | validator    | Compiles the target, runs the test suite, and emits a `ValidationReport`. When tests fail or coverage is thin, the report triggers a repair iteration on the translator. |

## Migration Strategy Registry

The analyzer picks one of five strategies via `Registry.TryInOrder`. Each
strategy has a `Matches(Profile) bool` predicate and an `Attempt(*State) error`
method. The first strategy that accepts the profile is the one that runs.

| Strategy         | Matches when |
| ---------------- | ------------ |
| INCREMENTAL      | always true (fallback) |
| PILOT            | source is large (`FileCount > 50` or `LoC > 10000`) |
| FROZEN_LEGACY    | source has no test suite |
| PARALLEL_CUTOVER | source has tests and `FileCount > 10` |
| BIG_BANG         | source is tiny (`FileCount <= 3` and `LoC < 500`) |

The cascade used to be a 4-branch `if`/`else if` chain in
`SelectMigrationStrategy`. That function is gone; the registry replaces it.

## Translator Modes

Two modes share the `TranslatorAgent`:

- **One-shot** (`generateTranslation`): single prompt with the full source
  tree. Used when the source fits in context.
- **Chunked** (`RunChunked`): fragments the source by topological order and
  prompts the coding model per fragment. Each fragment sees the previously
  emitted modules as a small context prefix.

The chunked mode emits a `Cargo.toml` skeleton and falls back to a minimal
manifest when the analyzer did not produce one.

## Validator Loop

The validator drives the repair loop. After each iteration it emits a
`ValidationReport` with:

- `AllSuccess` — compilation + tests pass with the real-test gate satisfied.
- `PlateauDetected` — coverage has stopped improving.
- `AdversarialWeakeningDetected` — tests have been gutted across iterations.
- `ASTSyntaxErrors` — pre-compilation parse errors.

The runner stops when `IsAllSuccess()` returns true, or when
`Iteration == MaxIterations`. Plateau remedy iterations (NEW-PRIM-27 /
CodaMOSA) run inside the validator before the loop exits.

## State Container

`internal/types.State` is the shared bag. The artifact types it embeds
(`AnalyzerOutput`, `PlanningOutput`, `TranslatedProject`, `ValidationReport`)
live in `internal/artifacts`, which has no upstream dependencies. This
three-package layout (`artifacts` → `types` → producers) avoids the
import cycle that a naive split would create.

See [[Architecture]] for the package DAG.

## Checkpointing

After each translator and validator run, the state is serialised to
`.artifacts/<run-id>/checkpoints/iter-NNNN.json`. A crashed run resumes from
the most recent checkpoint on the next invocation.

## Local-First Models

No cloud API. Both models run via Ollama at the URL the YAML config gives
(`OllamaBaseURL`). Defaults are documented in the README but the YAML config
is authoritative — the runner refuses to start without it.
