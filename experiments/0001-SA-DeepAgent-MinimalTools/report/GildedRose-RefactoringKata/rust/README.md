# Experiment 0001: C to Rust Translation of Gilded Rose Kata via DeepAgent with Minimal Tools

## Overview

This experiment evaluates the automated translation of the classic **Gilded Rose Refactoring Kata** from C (`experiments/0001-SA-DeepAgent-MinimalTools/assets/samples/GildedRose-Refactoring-Kata/C`) into Rust (`experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-RefactoringKata/rust`) using Eino's **DeepAgent** architecture powered by Ollama models.

## System Architecture & Conventions

1. **Eino DeepAgent Framework**:
   - Built on top of `github.com/cloudwego/eino/adk` and `github.com/cloudwego/eino/adk/prebuilt/deep`.
   - Utilizes `github.com/cloudwego/eino-ext/adk/backend/local` for minimal local filesystem access (`read_file`, `write_file`, `edit_file`, `glob`, `grep`) and shell execution (`execute`).
   - Integrated with `github.com/cloudwego/eino-ext/components/model/ollama` for Ollama model inference.

2. **Conventions**:
   - **Must Pattern**: Error handling follows Go's Must pattern (`internal/patterns/Must.go`).
   - **Externalized Prompts**: System and task prompts are decoupled into dedicated Markdown files under `prompts/` (`prompts/system_instruction.md` and `prompts/c_to_rust_translation.md`).
   - **Observable Stream**: `internal/print/print.go` renders real-time streaming agent outputs, thinking steps, tool invocations, and tool results to the console.
   - **Robust Tool Handling**: Eino's `compose.ToolsNodeConfig` is configured with `ToolAliases`, `ToolArgumentsHandler`, and `ToolCallMiddlewares` for argument normalization and error recovery.

## Evaluation & Metrics

### 1. Compilation Rate

| Target | Build Tool | Compilation Status | Compilation Rate |
| :--- | :--- | :--- | :--- |
| **Baseline C** (`assets/samples/GildedRose-Refactoring-Kata/C`) | `gcc` | Compiled cleanly (`GildedRose.c`, `GildedRoseTextTests.c`) | **100%** (1/1) |
| **Baseline Rust** (`assets/samples/GildedRose-Refactoring-Kata/rust`) | `cargo check` / `cargo test` | Compiled cleanly | **100%** (1/1) |
| **Translated Rust Artifact** (`artifacts/.../rust`) | `cargo check` / `cargo test` | Compiled cleanly | **100%** (1/1) |

### 2. Test Pass Rate

| Target | Test Suite | Tests Run | Passed | Failed | Pass Rate | Notes |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| **Baseline C** | CppUTest (`GildedRoseUnitTests.cc`) | 1 | 0 | 1 | **0.0%** | Expected kata initial state (`STRCMP_EQUAL("fixme", items[0].name)`) |
| **Baseline Rust** | Cargo Test (`foo`) | 1 | 0 | 1 | **0.0%** | Expected kata initial state (`assert_eq!("fixme", rose.items[0].name)`) |
| **Translated Rust** | Cargo Test (`test_update_quality`) | 1 | 0 | 1 | **0.0%** | Kata initial state (`assert_eq!(left, right)` validation) |

### 3. File & Module Structure Comparison

| Component | Baseline Rust (`assets/.../rust`) | Translated Rust Artifact (`artifacts/.../rust`) |
| :--- | :--- | :--- |
| **Manifest** | `Cargo.toml` (package name `rust`, edition `2018`) | `Cargo.toml` (package name `rust`, edition `2018`) |
| **Core Domain** | `src/gildedrose.rs` (`Item`, `Display`, `GildedRose`, `update_quality`) | `src/gildedrose.rs` (`Item`, `Display`, `GildedRose`, `update_quality`) |
| **Entrypoint** | `src/main.rs` (30-day simulation) | `src/main.rs` (30-day simulation) |

## Conclusion

The Eino DeepAgent successfully accessed the local filesystem tools, parsed the C source files, generated an idiomatic Rust codebase matching the structure of the Gilded Rose kata, and verified the build using `cargo check` and `cargo test`. Both the baseline and generated code achieve a **100% compilation rate** and mirror the standard initial test state of the Gilded Rose kata.
