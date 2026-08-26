package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/pkg/logger"
	"MAgHARCM/pkg/tools"
	"MAgHARCM/pkg/types"
)

// TranslatorAgent executes Part A & Part B translation and applies repair fixes (§3.4).
type TranslatorAgent struct {
	Model model.BaseChatModel
}

// NewTranslatorAgent creates a TranslatorAgent instance.
func NewTranslatorAgent(m model.BaseChatModel) *TranslatorAgent {
	return &TranslatorAgent{Model: m}
}

// Run translates the code and tests, or executes repairs if validation report has failures.
func (t *TranslatorAgent) Run(ctx context.Context, state *types.State) (*types.State, error) {
	if state.TranslatedProject.Files == nil {
		state.TranslatedProject.Files = make(map[string]string)
	}

	isRepair := state.Iteration > 0 && !state.ValidationReport.IsAllSuccess()

	if isRepair {
		logger.LogAgent("Translator", "Repair Mode (Iteration %d/%d): Diagnosing validation errors and fixing code...",
			state.Iteration, state.MaxIterations)
		state.Log("[Translator] Repair Mode (Iteration %d/%d): Fixing reported validation errors", state.Iteration, state.MaxIterations)
		return t.repair(ctx, state)
	}

	logger.LogAgent("Translator", "Initial Translation Mode: Implementing Part A (Source) and Part B (Tests)...")
	state.Log("[Translator] Implementing Part A (Source) and Part B (Tests)")
	return t.translate(ctx, state)
}

func (t *TranslatorAgent) translate(ctx context.Context, state *types.State) (*types.State, error) {
	var sourceFilesData []string
	_ = filepath.Walk(state.Task.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err == nil {
			sourceFilesData = append(sourceFilesData, fmt.Sprintf("=== Source File: %s ===\n%s\n", filepath.Base(path), string(data)))
		}
		return nil
	})

	logger.LogStep("Prompting Coding Model for complete %s translation...", state.Task.TargetLang)

	prompt := fmt.Sprintf(`You are the Translator Agent in the ReCodeAgent multi-agent repository-level translation workflow.
Your task is to translate the source codebase (%s) into fully working, idiomatic, safe %s code.
Translate ALL source data models, structs, classes, functions, and logic directly from the provided source files preserving exact semantic behavior.
Translate ALL unit tests and characterization tests into the target test framework.

Language-Agnostic Translation Guidelines:
1. Preserve exact functionality, algorithms, boundary conditions, and control flows from the source code.
2. In %s, ensure exported functions and structs from the source headers are public and idiomatic.
3. In Rust, files in tests/ (e.g. tests/integration_tests.rs) are compiled as external crates linking to the library crate defined in Cargo.toml. Therefore, all test files in tests/ MUST import library symbols with 'use <crate_name>::*;' (where crate_name matches the package name in Cargo.toml, e.g. 'use gilded_rose::*;'), and MUST NOT use 'use super::*;'.
4. Write clean, complete code for all required files without placeholders.

=== Source Files ===
%s

=== Target Project Design ===
%s

=== Implementation Plan ===
%s

Generate the complete, working, compilable, and tested %s codebase.
Format output strictly using file blocks:
FILE: Cargo.toml
`+"```toml"+`
[package]
name = "..."
version = "0.1.0"
edition = "2021"

[dependencies]
`+"```"+`

FILE: src/lib.rs
`+"```rust"+`
// Full source translation
`+"```"+`

FILE: tests/integration_tests.rs
`+"```rust"+`
// Full test suite translation
`+"```"+`
`, state.Task.SourceLang, state.Task.TargetLang,
		state.Task.TargetLang,
		strings.Join(sourceFilesData, "\n"),
		state.AnalyzerOutput.Design.RawMarkdown,
		state.PlanningOutput.Plan.RawPlan,
		state.Task.TargetLang)

	resp, err := t.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert systems programmer translating source code into idiomatic, safe target code."),
		schema.UserMessage(prompt),
	})
	if err != nil {
		return nil, fmt.Errorf("translator model call failed: %w", err)
	}

	files := parseAllFileMarkers(resp.Content)

	if err := syncFilesToDisk(state.Task.TargetDir, files, state); err != nil {
		return nil, err
	}

	logger.LogAgent("Translator", "Successfully wrote %d translated files to %s", len(files), state.Task.TargetDir)
	state.Log("[Translator] Translated %d files to %s", len(files), state.Task.TargetDir)
	return state, nil
}

func (t *TranslatorAgent) repair(ctx context.Context, state *types.State) (*types.State, error) {
	var targetFilesData []string
	for relPath, content := range state.TranslatedProject.Files {
		targetFilesData = append(targetFilesData, fmt.Sprintf("=== Current File: %s ===\n%s\n", relPath, content))
	}

	logger.LogStep("Feeding compiler diagnostics and test failures to Coding Model for targeted repair...")

	prompt := fmt.Sprintf(`You are the Translator Agent in ReCodeAgent repair mode.
The Validator Agent reported compilation or test failures for the translated %s codebase.

=== Validation Diagnostics and Errors ===
%s

=== Current Codebase Files ===
%s

Please analyze the root cause of each error/failure and output the complete corrected files to fix all issues.
In Rust, remember that test files in tests/ are separate integration test crates that MUST import library symbols with 'use <crate_name>::*;' (e.g. 'use gilded_rose::*;'), and MUST NOT use 'use super::*;'.
Do NOT use placeholders. Output full working code.

Format output strictly using file blocks:
FILE: Cargo.toml
`+"```toml"+`
...
`+"```"+`

FILE: src/lib.rs
`+"```rust"+`
...
`+"```"+`

FILE: tests/integration_tests.rs
`+"```rust"+`
...
`+"```"+`
`, state.Task.TargetLang, state.ValidationReport.Diagnostics, strings.Join(targetFilesData, "\n"))

	resp, err := t.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert systems programmer debugging compiler errors and test failures."),
		schema.UserMessage(prompt),
	})
	if err != nil {
		return nil, fmt.Errorf("translator repair call failed: %w", err)
	}

	repairedFiles := parseAllFileMarkers(resp.Content)
	if err := syncFilesToDisk(state.Task.TargetDir, repairedFiles, state); err != nil {
		return nil, err
	}

	logger.LogAgent("Translator", "Repairs applied across %d files", len(repairedFiles))
	state.Log("[Translator] Applied repairs across %d files", len(repairedFiles))
	return state, nil
}

func syncFilesToDisk(targetDir string, files map[string]string, state *types.State) error {
	hasNewTest := false
	for relPath := range files {
		if strings.HasPrefix(relPath, "tests/") {
			hasNewTest = true
			break
		}
	}
	if hasNewTest {
		_ = filepath.Walk(filepath.Join(targetDir, "tests"), func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel(targetDir, path)
			if _, exists := files[rel]; !exists {
				_ = os.Remove(path)
				delete(state.TranslatedProject.Files, rel)
			}
			return nil
		})
	}

	for relPath, content := range files {
		state.TranslatedProject.Files[relPath] = content
		fullPath := filepath.Join(targetDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create dir for %s: %w", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", fullPath, err)
		}
		logger.LogTool("write_file", "Wrote %s to %s (%d bytes)", relPath, targetDir, len(content))
	}
	return nil
}

func parseAllFileMarkers(text string) map[string]string {
	files := make(map[string]string)
	lines := strings.Split(text, "\n")
	var currentFile string
	var currentLines []string

	flush := func() {
		if currentFile != "" && len(currentLines) > 0 {
			code := tools.CleanCodeContent(strings.Join(currentLines, "\n"))
			if code != "\n" && code != "" {
				files[currentFile] = code
			}
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "FILE:") || strings.HasPrefix(trimmed, "### File:") {
			flush()
			f := strings.TrimPrefix(trimmed, "FILE:")
			f = strings.TrimPrefix(f, "### File:")
			f = strings.TrimSpace(f)
			currentFile = f
			currentLines = nil
		} else if currentFile != "" {
			currentLines = append(currentLines, line)
		}
	}
	flush()
	return files
}
