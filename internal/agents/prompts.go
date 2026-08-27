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
