# MAgHARCM

**Multi-Agent Harness for Archaeology-guided Repository-level Code Migration**

Go-based translation harness that pairs a reasoning model with a coding model through a 6-stage pipeline (Intake → Analyze → Plan → Translate → Validate → Repair-loop) to migrate whole repositories between languages. Targets local SLMs via Ollama — no cloud API required.

## Requirements

- [mise](https://mise.jdx.dev) (manages Go + Rust toolchains via `mise.toml`)
- Ollama running on `http://localhost:11434` with the reasoning + coding models pulled.
  Defaults: `lfm2.5:8b-a1b-q4_K_M` (reasoning) and
  `hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL` (coding).
  Override per-sample via `.config/*.yml`.

## Quick start

```bash
mise install            # install pinned Go + Rust
mise run start          # run with .config/default.yml (if present)
mise run bench          # run every .config/<cat>-<proj>.yml once
mise run bench-k        # run every config K times (K=3) for stability
mise run clean          # remove bin/ and .artifacts/
```

Run a single sample by name:

```bash
go run ./cmd/MAgHARCM --config .config/oxidizer-gohistogram.yml
go run ./cmd/MAgHARCM-tui --config .config/gildedrose.yml
```

Configs are generated from the corpus under `assets/samples/`. Regenerate them
with `python3 scripts/generate-configs.py` whenever the corpus changes.

## Tasks

| Task                | What it does                                              |
| ------------------- | --------------------------------------------------------- |
| `mise run start`    | Run MAgHARCM with `.config/default.yml` if it exists      |
| `mise run dev`      | Alias for `start`                                         |
| `mise run bench`    | `bash scripts/run-all-samples.sh` — one run per config    |
| `mise run bench-k`  | `K=3 bash scripts/run-samples-k.sh` — K runs per config   |
| `mise run retry`    | `bash scripts/retry-stable.sh <config>`                  |
| `mise run build`    | Compile `bin/MAgHARCM` and `bin/MAgHARCM-tui`             |
| `mise run test`     | `go test ./...`                                           |
| `mise run lint`     | `go vet ./...`                                            |
| `mise run fmt`      | `go fmt ./...`                                            |
| `mise run clean`    | `rm -rf bin/ .artifacts/`                                 |

## Method

The translator chunks large repositories by topological fragment order and
reuses prior emitted modules as context. The validator compiles, runs the test
suite, and triggers repair iterations when tests fail or coverage is low.
Per-run checkpoints under `.artifacts/<run-id>/checkpoints/` enable resume
after a crash.

See `docs/.paper/` for the empirical evaluation and
`.obsidian/MAgHARCM/research/METHODOLOGY.md` for the methodology.
