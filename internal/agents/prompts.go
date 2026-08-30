package agents

import (
	"bytes"
	"fmt"
	"text/template"
)

const analyzerPromptTemplate = `You are the Analyzer Agent in the multi-agent repository-level translation workflow.
Your goal is to perform analysis of the source project and generate 3 documents:
1. Source Project Research
2. Third-Party Library Analysis
3. Target Project Design

Source Language: {{.SourceLang}}
Target Language: {{.TargetLang}}
Source Codebase Directory: {{.SourceDir}}

Directory Tree:
{{.DirectoryTree}}

AST File Structure Summary:
{{.StructureSummary}}

Source Files Content:
{{.SourceFilesContent}}

Generate the 3 documents in structured markdown format with exact headers:

# Source Project Research
## Overview
## Directory Structure
## Structs and Interfaces
## Data Models
## Error Handling
## Dependencies

# Third-Party Library Analysis
## Standard and Third-Party Libraries
Identify all libraries/headers used and their idiomatic equivalents in {{.TargetLang}}.

# Target Project Design
## Overview
## Translation Requirements
## Source Files to Translate
## Module Structure
## Error Handling
## Third-Party Libraries`

const planningPromptTemplate = `You are the Planning Agent in the multi-agent repository-level translation workflow.
Using the Analyzer's Target Project Design and source fragments, you must output:
1. Name Mapping: A JSON map of source symbol -> target {{.TargetLang}} symbol name.
2. Skeleton Generation: Complete skeleton files for the target {{.TargetLang}} project containing all build definitions, module declarations, struct/type definitions, and function signatures.
3. Implementation Plan:
   - Part A: Source code translation steps (dependency ordered)
   - Part B: Test code translation and validation steps

Structure & Module Guidelines:
1. Rust Module Layout: Source files in ` + "`src/`" + ` must form a valid module tree (e.g. ` + "`src/lib.rs`" + `, and submodules like ` + "`src/item.rs`" + ` or ` + "`src/item/mod.rs`" + `). External integration tests live in ` + "`tests/`" + ` as separate test crates; do NOT declare ` + "`mod tests;`" + ` inside ` + "`src/lib.rs`" + `.
2. Go Package Layout: All source files in the same directory share the same package name.
3. Clean Paths: Skeleton file paths must only contain clean relative filepaths without descriptive comments (e.g. use ` + "`src/main.rs`" + `, never ` + "`src/main.rs (optional example usage)`" + `).

Source Language: {{.SourceLang}}
Target Language: {{.TargetLang}}

Source Files Content:
{{.SourceFiles}}

Target Design Document:
{{.TargetDesign}}

Output your response strictly using these delimiters:
=== NAME_MAPPING_JSON ===
{
  "source_symbol": "target_symbol"
}

=== SKELETON_FILES ===
FILE: path/to/file.ext
` + "```\n" + `// Declarations and signatures with todo!()/placeholder bodies
` + "```\n\n" + `=== IMPLEMENTATION_PLAN ===
## Overview
## Part A: Source Code Translation
A1: Translate ...
## Part B: Test Code Translation & Validation
B1: Translate and execute tests ...`

const translatorTranslatePromptTemplate = `You are the Translator Agent in the multi-agent repository-level translation workflow.
Your task is to translate the source codebase ({{.SourceLang}}) into fully working, idiomatic, safe {{.TargetLang}} code.
Translate ALL source data models, structs, classes, functions, and logic directly from the provided source files preserving exact semantic behavior.
Translate ALL unit tests and characterization tests into the target test framework.

Target Package / Module Name: {{.PackageName}}

Translation Guidelines:
1. Preserve exact functionality, algorithms, boundary conditions, and control flows from the source code.
2. In the main library root (e.g. ` + "`src/lib.rs` or `lib.go`" + `), declare all types, structs, constructors, and functions with public export visibility. If creating submodules, declare and re-export them publicly in the library root (e.g. ` + "`pub mod item; pub use item::*;`" + `).
3. External test files (e.g. in ` + "`tests/`" + `) and main application binaries (e.g. ` + "`src/main.rs`" + `) are separate compilation units. They must import all public library symbols from the target package (e.g. ` + "`use {{.PackageName}}::*;`" + ` in Rust at the top of each test file, ` + "`package {{.PackageName}}_test`" + ` in Go, ` + "`#include`" + ` in C). In Rust, NEVER use ` + "`use super::*;`" + ` in external test files and NEVER declare ` + "`mod tests;`" + ` in ` + "`src/lib.rs`" + ` for files in ` + "`tests/`" + `.
4. Ensure constructors (e.g. ` + "`Item::new`, `init_item`" + `) and function signatures match calling conventions and types. In Rust functions returning a reference from multiple reference arguments, include explicit lifetime annotations (e.g. ` + "`pub fn init_item<'a>(item: &'a mut Item, name: &str, sell_in: i32, quality: i32) -> &'a mut Item`" + `).
5. Rust Borrow-Checker: Avoid borrowing collections immutably while mutably borrowing them. Compute lengths, bounds, or indices into local variables BEFORE mutably slicing/indexing (e.g. ` + "`let len = items.len(); let slice = &mut items[..size.min(len)];`" + ` or use ` + "`items.iter_mut()`" + `).
6. Write clean, complete, fully closed code for all required files without placeholders, stubs, or truncation.

MANDATORY TEST REQUIREMENTS — failure to comply = non-deliverable:
- EVERY Rust test file you emit MUST contain at least 3 ` + "`#[test]`" + `-attributed functions. Translate every distinct scenario in source test files (one scenario per day in GildedRoseTextTests.c) into its own #[test].
- EVERY ` + "`#[test]`" + ` function MUST contain at least one ` + "`assert!`" + ` / ` + "`assert_eq!`" + ` / ` + "`assert_ne!`" + ` call.
- Test file paths MUST be under ` + "`tests/`" + ` (integration tests) — e.g. ` + "`tests/text_tests.rs`" + `, ` + "`tests/unit_tests.rs`" + `, ` + "`tests/characterization_tests.rs`" + `.
- At the TOP of every test file, ` + "`use {{.PackageName}}::*;`" + ` (replace ` + "`{{.PackageName}}`" + ` with the actual package name below).
- The translated test bodies MUST drive the same scenarios as the source tests (e.g. for GildedRose's ` + "`GildedRoseTextTests.c`" + `, emit day-by-day characterization tests asserting item name/sell_in/quality after N days of ` + "`update_quality`" + `).
- DO NOT emit empty test files or files with only comments. DO NOT emit ` + "`// This file is not needed`" + ` stubs.
- For C source with ` + "`GildedRoseTextTests.c`" + ` (a ` + "`main()`" + ` driver that calls ` + "`print_item`" + `), translate each scenario line into a Rust ` + "`#[test]`" + ` that asserts the printed string or the item's name/sell_in/quality fields after ` + "`update_quality`" + ` is called.

MANDATORY TEST REQUIREMENTS (other target languages) — failure to comply = non-deliverable:
- Go: every ` + "`_test.go`" + ` file MUST define at least 2 ` + "`func TestXxx(t *testing.T)`" + ` functions; each MUST call ` + "`t.Errorf`" + `, ` + "`t.Fatalf`" + `, ` + "`t.Errorf`" + ` via ` + "`if got != want { t.Errorf(...) }`" + `, or use ` + "`got := ...; if got != want { t.Fatalf(...) }`" + `. Package MUST be ` + "`{{.PackageName}}_test`" + `. DO NOT emit ` + "`_test.go`" + ` files containing only package declarations, comments, or empty bodies.
- C/C++ (Catch2 / GoogleTest): every test file MUST define at least 2 ` + "`TEST_CASE`" + ` (Catch2) or ` + "`TEST`" + ` (GoogleTest) cases; each MUST use ` + "`REQUIRE`" + ` / ` + "`CHECK`" + ` (Catch2) or ` + "`EXPECT_*`" + ` / ` + "`ASSERT_*`" + ` (GoogleTest). DO NOT emit empty or comment-only test files.
- Python (pytest / unittest): every test file MUST define at least 2 ` + "`def test_*`" + ` functions (pytest) or 2 ` + "`class TestX(unittest.TestCase)`" + ` methods named ` + "`test_*`" + `; each MUST contain at least one ` + "`assert`" + ` statement. DO NOT emit ` + "`pass`" + `-only test files or ` + "`# This file is not needed`" + ` stubs.
- JavaScript / TypeScript (Jest / Mocha): every test file MUST define at least 2 ` + "`test(`" + ` or ` + "`it(`" + ` blocks (Jest) or 2 ` + "`describe`" + ` cases with ` + "`it`" + ` (Mocha); each MUST contain at least one ` + "`expect(...).to*`" + ` (Jest) or ` + "`assert.*`" + ` (Chai/Node). DO NOT emit empty or stub test files.

=== Source Files ===
{{.SourceFiles}}

=== Target Project Design ===
{{.TargetDesign}}

=== Implementation Plan ===
{{.ImplementationPlan}}

Generate the complete, working, compilable, and tested {{.TargetLang}} codebase.
Format output strictly using file blocks:
FILE: path/to/target/file.ext
` + "```{{.TargetLangLower}}\n" + `// Full file implementation
` + "```"

const translatorChunkedPromptTemplate = `You are the Translator Agent translating ONE source file into {{.TargetLang}} as part of a larger chunked translation workflow.
The repository has already been partially translated. Use the "Previously Emitted Modules" block as state: it lists every target file already written and shows the first lines of its content, so you can re-use existing public symbols (types, constructors, helpers) without re-declaring them.

Target Package / Module Name: {{.PackageName}}

Single-File Scope:
- You are given EXACTLY ONE source file below. Translate ONLY that source file into the corresponding {{.TargetLang}} target file(s).
- Output ONE FILE block only — ` + "`### FILE: <target_path>`" + ` — for the source file you were given. Do not emit other files in this response.
- If the single source file logically requires declarations split across multiple target files (e.g. ` + "`src/lib.rs`" + ` re-exporting ` + "`src/item.rs`" + ` for Rust, or ` + "`lib.go`" + ` declaring package-level types for Go), emit ONE FILE block per resulting target file. Never re-emit files that already appear in the "Previously Emitted Modules" block — assume those already exist and import them.
- Preserve exact functionality, algorithms, boundary conditions, and control flows from the source file.
- Match calling conventions and types used in already-emitted modules: re-use struct names, enum names, function names, and module paths from the "Previously Emitted Modules" block instead of inventing new ones.
- Public visibility / export rules still apply: any symbol the rest of the project needs must be ` + "`pub`" + ` (Rust), capitalized (Go), or ` + "`export`" + `d (TS).
- Rust Borrow-Checker: avoid borrowing collections immutably while mutably borrowing them. Compute lengths, bounds, or indices into local variables BEFORE mutably slicing/indexing.
- If the source file is a test file, follow the same MANDATORY TEST REQUIREMENTS as the full-batch translator (minimum #[test]/Test*/def test_* count per language, use {{.PackageName}}::*, package {{.PackageName}}_test for Go, etc.).
- DO NOT emit empty files or stub comments. DO NOT re-emit already-emitted modules.

=== Previously Emitted Modules (state summary) ===
{{.PriorModules}}

=== Target Project Design (apply to this chunk) ===
{{.TargetDesign}}

=== Source File to Translate (single chunk) ===
{{.SourceFiles}}

Translate the single source file above into {{.TargetLang}} now. Format the output strictly using file blocks:
### FILE: path/to/target/file.ext
` + "```{{.TargetLangLower}}\n" + `// Full file implementation
` + "```"

const translatorRepairPromptTemplate = `You are the Translator Agent in repair mode.
The Validator Agent reported compilation or test failures for the translated {{.TargetLang}} codebase.

Target Package / Module Name: {{.PackageName}}

=== Validation Diagnostics and Errors ===
{{.Diagnostics}}

=== Current Codebase Files ===
{{.CurrentFiles}}

Please analyze the root cause of each error/failure and output the complete corrected files to fix all issues.
Guidelines:
1. Compiler Hints & Type Errors: Read the compiler error messages, line numbers, and suggestions carefully.
   - Rust Borrow-Checker (e.g. E0502: cannot borrow as immutable because also borrowed as mutable, such as calling .len() on a collection while indexing &mut into it): evaluate and save lengths/bounds/indices into local variables BEFORE mutably borrowing or slicing (e.g. ` + "`let len = items.len(); let slice = &mut items[..size.min(len)];`" + ` or iterate with ` + "`items.iter_mut()`" + `).
   - Fix all type mismatches, lifetime annotations (e.g. ` + "`pub fn init_item<'a>(item: &'a mut Item, name: &str, sell_in: i32, quality: i32) -> &'a mut Item`" + `), and ownership/clone issues.
2. Exported Symbols: In the library root (` + "`src/lib.rs` / `lib.go`" + `), ensure all types, structs, constructors (` + "`new`, `init_item`" + `), and functions (` + "`update_quality`, `print_item`" + `) are declared with public visibility or re-exported from submodules (` + "`pub mod ...; pub use ...::*;`" + `).
3. Package Imports in Tests & Binaries: External test files in ` + "`tests/`" + ` and binaries in ` + "`src/main.rs`" + ` are separate compilation units and MUST import public library symbols from the package root via ` + "`use {{.PackageName}}::*;`" + ` at the top of each test file (never ` + "`use super::*;`" + `, and never declare ` + "`mod tests;`" + ` in ` + "`src/lib.rs`" + ` for files in ` + "`tests/`" + `).
4. Emit All Failing Files: If a compiler error occurs in any source or test file (e.g. ` + "`src/gilded_rose/update.rs`" + ` or ` + "`tests/unit/mod.rs`" + `), you MUST emit the complete corrected file in your response.
5. Complete Working Files: Output complete, corrected file implementations with all braces and delimiters closed inside code blocks without placeholders or truncation.

Format output strictly using file blocks:
FILE: path/to/target/file.ext
` + "```{{.TargetLangLower}}\n" + `// Full corrected file implementation
` + "```"

const validatorCoveragePromptTemplate = `You are the Validator Agent performing Coverage-Guided Test Generation.
The following functions or modules are uncovered or need more test assertions in the translated {{.TargetLang}} project:
{{.UncoveredFunctions}}

Target Source Files:
{{.SourceFiles}}

Target minimum test count: {{.MinRealTests}}. You MUST emit at least this many #[test] functions across the test files.
Emit AT LEAST 2× the number of source functions (one happy path + one edge case per function). Always aim higher than the minimum — more tests is better.
Translate EACH scenario in the source test files (e.g. the GildedRoseTextTests.c main() driver scenarios, one per day) into a separate #[test] function.

Generate additional comprehensive unit test cases in the target test framework to thoroughly exercise all uncovered functions and edge cases.
Output the complete updated test file:

FILE: {{.TestFileRelPath}}
` + "```{{.TargetLangLower}}\n" + `// Complete updated test suite with additional test assertions
` + "```"

// renderPromptTemplate parses and executes a prompt template string with data.
func renderPromptTemplate(name, tmplStr string, data any) (string, error) {
	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse prompt template %s: %w", name, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute prompt template %s: %w", name, err)
	}

	return buf.String(), nil
}
