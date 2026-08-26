package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"MAgHARCM/pkg/config"
	"MAgHARCM/pkg/graph"
	"MAgHARCM/pkg/llm"
	"MAgHARCM/pkg/types"
)

const defaultRequirementsPrompt = `Input codebase: assets/samples/GildedRose-Refactoring-Kata/C
Output directory: .artifacts/GildedRose-Refactoring-Kata/rust
Input language: C
Output language: Rust`

func main() {
	var promptFlag string
	var sourceFlag string
	var targetFlag string
	var srcLangFlag string
	var tgtLangFlag string

	flag.StringVar(&promptFlag, "prompt", "", "Requirements prompt for code translation")
	flag.StringVar(&sourceFlag, "source", "", "Source codebase path")
	flag.StringVar(&targetFlag, "target", "", "Target output directory")
	flag.StringVar(&srcLangFlag, "source-lang", "", "Source language")
	flag.StringVar(&tgtLangFlag, "target-lang", "", "Target language")
	flag.Parse()

	// 1. Determine Prompt
	prompt := promptFlag
	if prompt == "" && len(flag.Args()) > 0 {
		prompt = strings.Join(flag.Args(), " ")
	}
	if prompt == "" {
		// Check stdin if piped
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			if stdinBytes, err := io.ReadAll(os.Stdin); err == nil && len(stdinBytes) > 0 {
				prompt = strings.TrimSpace(string(stdinBytes))
			}
		}
	}
	if prompt == "" {
		prompt = defaultRequirementsPrompt
	}

	fmt.Println("==================================================")
	fmt.Println("ReCodeAgent (arXiv:2604.07341) Translation Pipeline")
	fmt.Println("==================================================")
	fmt.Printf("Received Prompt:\n%s\n\n", prompt)

	// 2. Parse Requirements from Prompt
	task := parsePromptRequirements(prompt)
	if sourceFlag != "" {
		task.SourceDir = sourceFlag
	}
	if targetFlag != "" {
		task.TargetDir = targetFlag
	}
	if srcLangFlag != "" {
		task.SourceLang = srcLangFlag
	}
	if tgtLangFlag != "" {
		task.TargetLang = tgtLangFlag
	}

	fmt.Printf("Parsed Task Configuration:\n")
	fmt.Printf("  • Source Codebase: %s\n", task.SourceDir)
	fmt.Printf("  • Output Dir:      %s\n", task.TargetDir)
	fmt.Printf("  • Source Language: %s\n", task.SourceLang)
	fmt.Printf("  • Target Language: %s\n\n", task.TargetLang)

	// 3. Centralized Configuration (Must pattern)
	cfg := config.MustLoad(
		config.WithSourceTarget(task.SourceDir, task.TargetDir, task.SourceLang, task.TargetLang),
	)

	ctx := context.Background()

	// 4. Initialize Ollama Models (Must pattern)
	fmt.Printf("Initializing Ollama models from %s:\n", cfg.OllamaBaseURL)
	fmt.Printf("  • Reasoning Model: %s\n", cfg.ReasoningModel)
	fmt.Printf("  • Coding Model:    %s\n\n", cfg.CodingModel)
	models := llm.MustNewModels(ctx, cfg)

	// 5. Build Eino Graph with 4 ReCodeAgent agents and repair loop (Must pattern)
	fmt.Println("Constructing 4-agent Eino Graph (Analyzer -> Planning -> Translator <-> Validator)...")
	recodeGraph := graph.MustNewReCodeGraph(ctx, models)

	// 6. Ensure target directory exists
	_ = os.MkdirAll(cfg.TargetDir, 0755)

	// 7. Initial State
	initialState := &types.State{
		Task:          task,
		MaxIterations: cfg.MaxIterations,
		TranslatedProject: types.TranslatedProject{
			Files: make(map[string]string),
		},
	}

	// 8. Execute Translation Pipeline
	fmt.Println("Executing multi-agent translation pipeline...")
	finalState, err := recodeGraph.Execute(ctx, initialState)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing ReCodeAgent pipeline: %v\n", err)
		os.Exit(1)
	}

	// 9. Output Results
	fmt.Println("\n==================================================")
	fmt.Println("Translation Pipeline Summary")
	fmt.Println("==================================================")
	for _, l := range finalState.Logs {
		fmt.Println(l)
	}

	fmt.Println("\nValidation Status:")
	fmt.Println(finalState.ValidationReport.String())

	if finalState.ValidationReport.IsAllSuccess() {
		fmt.Println("\n✓ SUCCESS: Translation and validation completed successfully!")
		fmt.Printf("Output available in: %s\n", filepath.Clean(cfg.TargetDir))
	} else {
		fmt.Println("\n⚠ COMPLETED WITH WARNINGS: Some tests or checks failed after maximum iterations.")
		os.Exit(1)
	}
}

func parsePromptRequirements(prompt string) types.TranslationTask {
	task := types.TranslationTask{
		SourceDir:  "assets/samples/GildedRose-Refactoring-Kata/C",
		TargetDir:  ".artifacts/GildedRose-Refactoring-Kata/rust",
		SourceLang: "C",
		TargetLang: "Rust",
		Prompt:     prompt,
	}

	lines := strings.Split(prompt, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)

		if strings.HasPrefix(lower, "input codebase:") {
			raw := strings.TrimSpace(strings.TrimPrefix(trimmed[len("input codebase:"):], " "))
			task.SourceDir = strings.TrimPrefix(raw, "@")
		} else if strings.HasPrefix(lower, "output directory:") {
			raw := strings.TrimSpace(strings.TrimPrefix(trimmed[len("output directory:"):], " "))
			task.TargetDir = strings.TrimPrefix(raw, "@")
		} else if strings.HasPrefix(lower, "input language:") {
			task.SourceLang = strings.TrimSpace(strings.TrimPrefix(trimmed[len("input language:"):], " "))
		} else if strings.HasPrefix(lower, "output language:") {
			task.TargetLang = strings.TrimSpace(strings.TrimPrefix(trimmed[len("output language:"):], " "))
		}
	}

	// If single line or natural language sentence
	if task.SourceDir == "" || task.SourceDir == "assets/samples/GildedRose-Refactoring-Kata/C" {
		if m := regexp.MustCompile(`@?([a-zA-Z0-9_\-\./]+/GildedRose[a-zA-Z0-9_\-\./]*)`).FindStringSubmatch(prompt); len(m) > 1 {
			clean := strings.TrimPrefix(m[1], "@")
			if !strings.Contains(clean, ".artifacts") {
				task.SourceDir = clean
			}
		}
	}

	if m := regexp.MustCompile(`(?:output to|into|to|Output directory:?)\s+([a-zA-Z0-9_\-\./]*\.artifacts[a-zA-Z0-9_\-\./]*)`).FindStringSubmatch(prompt); len(m) > 1 {
		task.TargetDir = m[1]
	}

	if strings.Contains(prompt, "from C") || strings.Contains(prompt, "Input language: C") || strings.Contains(prompt, "translates @assets/samples/GildedRose-Refactoring-Kata/C") {
		task.SourceLang = "C"
	}
	if strings.Contains(prompt, "Rust") || strings.Contains(prompt, "rust") {
		task.TargetLang = "Rust"
	}

	return task
}
