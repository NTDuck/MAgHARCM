# MAgHARCM

**Multi-Agent Harness for Archaeology-guided Repository-level Code Migration**

Go-based translation harness that pairs a reasoning model with a coding model through a 6-stage pipeline (Intake → Analyze → Plan → Translate → Validate → Repair-loop) to migrate whole repositories between languages.

## Requirements

- [mise](https://mise.jdx.dev) (manages Go + Rust toolchains via `mise.toml`)
- Ollama running on `http://localhost:11434` with `lfm2.5:8b-a1b-q4_K_M` (reasoning) and `hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL` (coding) pulled.

## Quick start

```bash
mise install               # install pinned Go + Rust
mise run start             # translate sample defined in config.yml
mise run clean             # remove bin/ and .artifacts/
```

Switch samples with `--config`:

```bash
go run ./cmd/MAgHARCM --config config.sample1.yml   # GildedRose    C    -> Rust
go run ./cmd/MAgHARCM --config config.sample2.yml   # stats         Go   -> Rust
go run ./cmd/MAgHARCM --config config.sample3.yml   # gohistogram   Go   -> Rust
go run ./cmd/MAgHARCM --config config.sample4.yml   # commons-valid. Java -> Rust
```

## Tasks

| Task              | What it does                              |
| ----------------- | ----------------------------------------- |
| `mise run start`  | Run MAgHARCM with `config.yml`            |
| `mise run dev`    | Alias of `start`                          |
| `mise run build`  | Compile `bin/MAgHARCM`                    |
| `mise run test`   | `go test ./...`                           |
| `mise run lint`   | `go vet ./...`                            |
| `mise run fmt`    | `go fmt ./...`                            |
| `mise run clean`  | `rm -rf bin/ .artifacts/`                 |

## Method

The translator chunks large repositories by topological fragment order and reuses prior emitted modules as context. The validator compiles, runs the test suite, and triggers repair iterations when tests fail or coverage is low. Per-run checkpoints under `.artifacts/<run-id>/checkpoints/` enable resume after a crash.

See `paper/` for the empirical evaluation and `.obsidian/MAgHARCM/research/METHODOLOGY.md` for the methodology.
