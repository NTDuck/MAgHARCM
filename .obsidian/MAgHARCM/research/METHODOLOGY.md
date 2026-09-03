---
title: MAgHARCM Methodology
backlink: [[Methodology]]
tags: [methodology, architecture, pipeline]
---

# MAgHARCM Methodology

## Pipeline

Four agents, wired as Eino graph nodes. Each node reads from the shared
`types.State` and writes back into it. Repair is a back-edge from
Validator to Translator rather than a separate agent.

```
START → analyzer → planner → translator → validator → branch
                                                ^         |
                                                └─repair──┘
```

| Stage | Agent    | What it does |
| ----- | -------- | ------------ |
| 1     | analyzer   | Walks the source tree, runs the 8B reasoning model, and emits three documents: Source Project Research, Third-Party Library Analysis, Target Project Design. Picks a migration strategy via `internal/agents/strategy.go::Registry.TryInOrder`. |
| 2     | planner    | Reads analyzer output, fragments the AST into translation units, resolves ambiguous symbol names through the configured LSP provider (Tree-sitter or ABCoder MCP), and emits an `ImplementationPlan` with two ordered lists: source translation steps (Part A) and test generation + validation steps (Part B). Builds a reverse-topological DAG with back-edge cuts. |
| 3     | translator | Runs the 4B coding model per fragment in topological order. Each fragment sees bounded context (Symbol Navigator + Prior-Modules Memory, each capped at 4KB) and the previously emitted modules as a context prefix. Persists a checkpoint after each fragment. |
| 4     | validator  | Compiles the target, runs the test suite, and emits a `ValidationReport`. The cascade runs AST pre-check → cargo check → cargo test → adversarial weakening guard → CodaMOSA plateau detector. When the cascade detects failure, the report triggers a repair iteration back on the translator. |

## Migration Strategy Registry

The analyzer picks a strategy via `Registry.TryInOrder`. Each strategy has a
`Matches(Profile) bool` predicate and an `Attempt(*State) error` method. The
registry walks the slice and tries the first strategy that accepts the
profile. On `error`, the runner skips to the next strategy; only an exhausted
registry terminates the pipeline.

| Strategy         | Matches when |
| ---------------- | ------------ |
| INCREMENTAL      | always true (fallback) |
| PILOT            | source is large (`FileCount > 50` or `LoC > 10000`) |
| FROZEN_LEGACY    | source has no test suite |
| PARALLEL_CUTOVER | source has tests and `FileCount > 10` |
| BIG_BANG         | source is tiny (`FileCount <= 3` and `LoC < 500`) |

The cascade used to be a hardcoded 4-branch `if`/`else if` chain in
`SelectMigrationStrategy`. That function is gone; the registry replaces it
with try-and-fail iteration so adding a strategy is a one-line append.

## Chunked Translator

The translator runs in two modes sharing `TranslatorAgent`:

- **One-shot** (`generateTranslation`): single prompt with the full source
  tree. Used when the source fits in context.
- **Chunked** (`RunChunked`): fragments the source by topological order and
  prompts the coding model per fragment. Each fragment sees the previously
  emitted modules as a small context prefix.

The chunked mode emits a `Cargo.toml` skeleton and falls back to a minimal
manifest when the analyzer did not produce one.

## Validator Cascade

The validator runs five checks in sequence. Each check can short-circuit the
later checks:

1. **AST pre-check** — parse the generated target with the configured
   Tree-sitter grammar; emit `ASTSyntaxErrors` before invoking the compiler.
2. **Compiler build** — run `cargo check`; capture rustc error codes for
   borrow-checker failures (`E0382`, `E0502`, `E0597`).
3. **Test suite** — run `cargo test`; emit per-test pass/fail matrix.
4. **Weakening guard** — diff the current test AST against the source test
   AST; if any assertion was dropped, halt the repair loop and emit
   `AdversarialWeakeningDetected`.
5. **Plateau detector** — CodaMOSA-style coverage-plateau check; if coverage
   has not improved across the last N iterations, exit the loop with
   `PlateauDetected` to prevent unbounded spinning.

The runner stops when `IsAllSuccess()` returns true, when
`Iteration == MaxIterations`, or when the weakening guard fires.

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
