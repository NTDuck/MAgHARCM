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
	"MAgHARCM/pkg/logger"
	"MAgHARCM/pkg/types"
)

func main() {
	var promptFileFlag string
	var promptFlag string
	var sourceFlag string
	var targetFlag string
	var srcLangFlag string
	var tgtLangFlag string

	flag.StringVar(&promptFileFlag, "prompt-file", "", "Path to file containing requirements prompt")
	flag.StringVar(&promptFlag, "prompt", "", "Direct string requirements prompt for code translation")
	flag.StringVar(&sourceFlag, "source", "", "Source codebase path (overrides prompt)")
	flag.StringVar(&targetFlag, "target", "", "Target output directory (overrides prompt)")
	flag.StringVar(&srcLangFlag, "source-lang", "", "Source language (overrides prompt)")
	flag.StringVar(&tgtLangFlag, "target-lang", "", "Target language (overrides prompt)")
	flag.Parse()

	// 1. Determine Prompt Content
	prompt := promptFlag

	if prompt == "" && promptFileFlag != "" {
		data, err := os.ReadFile(promptFileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading prompt file %s: %v\n", promptFileFlag, err)
			os.Exit(1)
		}
		prompt = strings.TrimSpace(string(data))
	}

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
		// Fallback to prompts/c_to_rust.txt if it exists in current dir
		if data, err := os.ReadFile("prompts/c_to_rust.txt"); err == nil {
			prompt = strings.TrimSpace(string(data))
		}
	}

	if prompt == "" {
		fmt.Fprintln(os.Stderr, "Error: No prompt provided. Specify --prompt-file <path>, --prompt <text>, or pass via stdin.")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("==================================================")
	fmt.Println("MAgHARCM Multi-Agent Translation Pipeline")
	fmt.Println("==================================================")
	fmt.Printf("Input Prompt:\n%s\n\n", prompt)

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

	if task.SourceDir == "" || task.TargetDir == "" {
		fmt.Fprintln(os.Stderr, "Error: Could not extract source codebase and output directory from prompt.")
		os.Exit(1)
	}

	fmt.Printf("Parsed Task Specification:\n")
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
	logger.LogStep("Connecting to Ollama models at %s...", cfg.OllamaBaseURL)
	logger.LogStep("Reasoning Model: %s", cfg.ReasoningModel)
	logger.LogStep("Coding Model:    %s", cfg.CodingModel)
	models := llm.MustNewModels(ctx, cfg)

	// 5. Build Eino Graph with 4 ReCodeAgent agents and repair loop (Must pattern)
	logger.LogStep("Constructing 4-agent Eino Graph (Analyzer -> Planning -> Translator <-> Validator)...")
	magharcmGraph := graph.MustNewMAgHARCMGraph(ctx, models)
	// 6. Ensure target directory exists
	_ = os.MkdirAll(cfg.TargetDir, 0o755)

	// 7. Initial State
	initialState := &types.State{
		Task:          task,
		MaxIterations: cfg.MaxIterations,
		TranslatedProject: types.TranslatedProject{
			Files: make(map[string]string),
		},
	}

	// 8. Execute Translation Pipeline
	logger.LogStep("Starting multi-agent translation execution...")
	finalState, err := magharcmGraph.Execute(ctx, initialState)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing MAgHARCM pipeline: %v\n", err)
		os.Exit(1)
	}

	// 9. Output Results
	fmt.Println("\n==================================================")
	fmt.Println("MAgHARCM Execution Summary")
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
		Prompt: prompt,
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

	// If prompt is in natural language
	if task.SourceDir == "" {
		if m := regexp.MustCompile(`(?i)(?:translate[sd]?|codebase:)\s+@?([a-zA-Z0-9_\-\./]+)`).FindStringSubmatch(prompt); len(m) > 1 {
			task.SourceDir = strings.TrimPrefix(m[1], "@")
		}
	}

	if task.TargetDir == "" {
		if m := regexp.MustCompile(`(?i)(?:output to|at|directory:)\s+@?([a-zA-Z0-9_\-\./]+)`).FindStringSubmatch(prompt); len(m) > 1 {
			task.TargetDir = strings.TrimPrefix(m[1], "@")
		}
	}

	if task.SourceLang == "" {
		if m := regexp.MustCompile(`(?i)\bfrom\s+([a-zA-Z\+#]+)\b`).FindStringSubmatch(prompt); len(m) > 1 {
			task.SourceLang = m[1]
		}
	}

	if task.TargetLang == "" {
		if m := regexp.MustCompile(`(?i)\b(?:into|to)\s+(?:a\s+)?([a-zA-Z\+#]+)\b`).FindStringSubmatch(prompt); len(m) > 1 {
			task.TargetLang = m[1]
		}
	}

	return task
}
