---
title: MAgHARCM Methodology
backlink: [[Methodology]]
tags: [methodology, architecture, pipeline]
---

# MAgHARCM Methodology

## Pipeline ([[PRIM-1]], [[PRIM-2]], [[PRIM-3]], [[PRIM-23]])

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
| 1     | analyzer   | Walks the source tree, runs the 8B reasoning model, and emits three documents: Source Project Research, Third-Party Library Analysis, Target Project Design. Picks a migration strategy via `internal/agents/strategy.go::Registry.TryInOrder` ([[PRIM-21]]). |
| 2     | planner    | Reads analyzer output, fragments the AST ([[PRIM-2]]) into translation units, resolves ambiguous symbol names ([[PRIM-26]]), and emits an `ImplementationPlan` with two ordered lists: source translation steps (Part A) and test generation + validation steps (Part B). Builds reverse-topological DAG ([[PRIM-1]]). |
| 3     | translator | Runs the 4B coding model per fragment in topological order ([[PRIM-23]]). Each fragment sees bounded context ([[PRIM-26]] + [[PRIM-31]], each capped at 4KB) and the previously emitted modules as a context prefix. Optional role-flip gate ([[PRIM-25]]) inspects output for hallucinated bug-free claims. Persists a checkpoint ([[PRIM-28]]) after each fragment. |
| 4     | validator  | Compiles the target, runs the test suite, and emits a `ValidationReport`. The cascade runs AST pre-check ([[PRIM-6]]) → cargo check → cargo test → adversarial weakening guard ([[PRIM-13]]) → plateau detector ([[PRIM-27]]). Optional auxiliary checks: multi-agent verdict ([[PRIM-7]]), mock-based in-isolation ([[PRIM-8]]), implementation-agnostic I/O ([[PRIM-11]]), Wasm oracle ([[PRIM-12]]). When the cascade detects failure, the report triggers a repair iteration back on the translator. |

## Migration Strategy Registry ([[PRIM-21]])

The analyzer picks a strategy via `Registry.TryInOrder`. Each strategy
has a `Matches(Profile) bool` predicate and an `Attempt(*State) error`
method. The registry walks the slice and tries the first strategy that
accepts the profile. On `error`, the runner skips to the next strategy;
only an exhausted registry terminates the pipeline.

| Strategy         | Matches when |
| ---------------- | ------------ |
| INCREMENTAL      | always true (fallback) |
| PILOT            | source is large (`FileCount > 50` or `LoC > 10000`) |
| FROZEN_LEGACY    | source has no test suite |
| PARALLEL_CUTOVER | source has tests and `FileCount > 10` |
| BIG_BANG         | source is tiny (`FileCount <= 3` and `LoC < 500`) |

The cascade used to be a hardcoded 4-branch `if`/`else if` chain in
`SelectMigrationStrategy`. That function is gone; the registry replaces
it with try-and-fail iteration so adding a strategy is a one-line
append.

## Chunked Translator ([[PRIM-23]])

The translator runs in two modes sharing `TranslatorAgent`:

- **One-shot** (`generateTranslation`): single prompt with the full
  source tree. Used when the source fits in context.
- **Chunked** (`RunChunked`): fragments the source by topological order
  ([[PRIM-1]]) and prompts the coding model per fragment. Each fragment
  sees the previously emitted modules as a small context prefix
  ([[PRIM-31]] Iterative Retrieval Refinement).

The chunked mode emits a `Cargo.toml` skeleton ([[PRIM-3]]) and falls
back to a minimal manifest when the analyzer did not produce one
([[PRIM-10]]).

## Validator Cascade ([[PRIM-6]], [[PRIM-13]], [[PRIM-27]])

The validator runs five core checks in sequence. Each check can
short-circuit the later checks:

1. **AST pre-check** ([[PRIM-6]]) — parse the generated target with the
   configured Tree-sitter grammar; emit `ASTSyntaxErrors` before
   invoking the compiler.
2. **Compiler build** — run `cargo check`; capture rustc error codes for
   borrow-checker failures (`E0382`, `E0502`, `E0597`).
3. **Test suite** — run `cargo test`; emit per-test pass/fail matrix.
4. **Weakening guard** ([[PRIM-13]]) — diff the current test AST against
   the source test AST; if any assertion was dropped, halt the repair
   loop and emit `AdversarialWeakeningDetected`.
5. **Plateau detector** ([[PRIM-27]]) — CodaMOSA-style
   coverage-plateau check; if coverage has not improved across the
   last N iterations, exit the loop with `PlateauDetected` to prevent
   unbounded spinning.

Optional auxiliary validators that the runner can compose in:

- [[PRIM-7]] Multi-Agent Verdict Validation — three independent LLM
  judges prompt the translator to repair on majority disagreement.
- [[PRIM-8]] State-Grounded Mock-Based In-Isolation Validation — for each
  module, validate against state-equivalent stub mocks.
- [[PRIM-11]] Implementation-Agnostic Testing — black-box I/O vectors,
  no AST inspection.
- [[PRIM-12]] Wasm-Based Reference Execution Oracle — Wasm-compiled
  source as ground truth.

The runner stops when `IsAllSuccess()` returns true, when `Iteration ==
MaxIterations`, or when the weakening guard fires.

## Software-Archaeology Stage ([[PRIM-14]])

Optional pre-translation comprehension pass that emits an
`ArchaeologyReport` consumed by the planner. Sub-steps:

1. `ExtractBoundaries` — find module/file/function boundaries.
2. `BuildTimeCapsule` — emit shell commands to reproduce the original
   build/test environment.
3. `FindChurnHotspots` — git-churn-driven hotspots (empty if no git).
4. `ForensicNaming` — scan identifiers for legacy encoding markers.
5. `MapConcepts` — co-occurrence of identifiers across files; cluster
   into concept groups.

## State Container ([[PRIM-24]])

`internal/types.State` is the shared bag. The artifact types it embeds
(`AnalyzerOutput`, `PlanningOutput`, `TranslatedProject`,
`ValidationReport`) live in `internal/artifacts`, which has no upstream
dependencies. Every artifact struct carries a `SchemaVersion` string
field (`[[PRIM-24]]` SOP-Anchored Role-Artifact Schema) for
forward-compatible migrations. This three-package layout (`artifacts` →
`types` → producers) avoids the import cycle that a naive split would
create.

See [[Architecture]] for the package DAG.

## Checkpointing ([[PRIM-28]])

After each translator and validator run, the state is serialised to
`.artifacts/<run-id>/checkpoints/iter-NNNN.json` ([[PRIM-28]]). A crashed
run resumes from the most recent checkpoint on the next invocation.

## Recruitment-Adaptive Planning ([[PRIM-29]])

Before each iteration the runner calls `Recruiter.Recruit(profile,
lastReport)` and selects which tools and agents to invoke based on the
last validator report. Default policy: re-run validator on
compilation failure; re-run translator with `recruit_translator_v2`
flavour on plateau; otherwise standard pipeline.

## Local-First Models

No cloud API. Both models run via Ollama at the URL the YAML config
gives (`OllamaBaseURL`). Defaults are documented in the README but the
YAML config is authoritative — the runner refuses to start without it.

## Versioning Convention

All version/ID markers use the `[[x.y.z ...]]` convention. Spontaneous
parentheses (`(CodaMOSA)`, `(Syzygy)`, `(ReCodeAgent)`) are NOT allowed
in vault files; replace with `[[Author-YEAR]]` or `[[P-NN]]` style
backlinks.
