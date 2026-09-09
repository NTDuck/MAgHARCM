# Crust C→Rust benchmark suite

Long-horizon batch runner over the [ReCodeAgent crust benchmark](https://github.com/zilongp/RecodeAgent/tree/master/data/tool_projects/crust) corpus. 100 single-file / multi-file C projects, each translated by MAgHARCM into Rust and validated against the original test suite.

## Layout

```
benchmarks/crust/
├── configs/                          # one YAML per project (matches .config/* template)
│   ├── 2dpartint.yml
│   ├── 42-kocaeli-printf.yml
│   └── ...                           # 100 total
├── scripts/
│   └── run.sh                        # the batch runner
└── results/<commit>/
    ├── index.yml                     # aggregate
    └── <proj>/
        ├── run.log                   # full MAgHARCM console output
        ├── result.yml                # scraped metrics
        └── .interrupted              # sentinel when SIGINT killed the run
```

Each config mirrors the `.config/oxidizer/*` template shape so the same runner that does `mise run bench` can also drive this suite. Source language is always `C` (from `benchmarks/assets/ReCodeAgent/data/tool_projects/crust/<proj>/c/`); target is always `Rust` written to `.artifacts/crust/<proj>/rust/`.

## Quick start

```bash
# Run every project sequentially (long-horizon; designed for overnight use).
bash benchmarks/crust/scripts/run.sh

# Just one project.
bash benchmarks/crust/scripts/run.sh --only 2dpartint

# Pin a past commit's results.
bash benchmarks/crust/scripts/run.sh --commit HEAD~1

# Lower per-project wall budget.
bash benchmarks/crust/scripts/run.sh --timeout 5400

# Dry-run: list every config without executing.
bash benchmarks/crust/scripts/run.sh --dry-run
```

The runner builds `bin/MAgHARCM` once and reuses it for every project, mirroring the contract `scripts/run-all-samples.sh` already establishes for the alphatrans / oxidizer / gildedrose samples.

## Resume semantics

- `result.yml` with `status: completed` ⇒ skipped (the MAgHARCM run emitted `FINALSUM` before exit).
- `.interrupted` sentinel present ⇒ invalidated; project re-runs.
- `result.yml` missing or `status != completed` ⇒ invalidated; project re-runs.

Invalidation removes `.interrupted` and the stale `result.yml` before the new run starts, so downstream consumers never see half-finished artefacts.

## Interrupt handling

`Ctrl-C` / `SIGTERM`:

1. Let the current project's `timeout` finish (or kill it).
2. Mark the current project as interrupted (`.interrupted` sentinel + `status: interrupted`).
3. Rebuild `index.yml` with the partial set.
4. Exit with `130`.

Re-running the script resumes from after the interrupted project.

## Metrics

The runner scrapes the validator's per-iteration log lines:

```
[hh:mm:ss] ITER[1] comp=true tests=3/3 pass-rate=100.0% wall=4213ms per-file=4
[hh:mm:ss] ITER[2] comp=true tests=3/3 pass-rate=100.0% wall=3891ms per-file=4
[hh:mm:ss] FINALSUM comp=true tests=3/3 all_success=true wall=8104ms iter=2
```

`FINALSUM` is emitted exactly once by `internal/agents/validator.go` at the end of every pipeline run. Its presence is the canonical "this run finished cleanly" signal; absence means the pipeline crashed or timed out and the project is marked `interrupted`.

Per-project `result.yml`:

```yaml
project: 2dpartint
config: benchmarks/crust/configs/2dpartint.yml
status: completed
exit_code: 0
wall_seconds: 7200
compile_ok: true
tests_passed: 3
tests_total: 3
pass_rate: 1.0000
iterations: 2
final_iter_compile: true
final_iter_tests_passed: 3
final_iter_tests_total: 3
per_iteration:
  - 1 comp=true tests_passed=3 tests_total=3
  - 2 comp=true tests_passed=3 tests_total=3
```

`index.yml` aggregates one entry per project for fleet-level reporting.
