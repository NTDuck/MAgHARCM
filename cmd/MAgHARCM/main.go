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

	"MAgHARCM/internal/config"
	"MAgHARCM/internal/graph"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/types"
)

func main() {
	var configFlag string
	var promptFileFlag string
	var promptFlag string
	var sourceFlag string
	var targetFlag string
	var srcLangFlag string
	var tgtLangFlag string
	var toolchainFlag string

	flag.StringVar(&configFlag, "config", "", "Path to YAML configuration file (default: config.yml)")
	flag.StringVar(&promptFileFlag, "prompt-file", "", "Path to file containing requirements prompt")
	flag.StringVar(&promptFlag, "prompt", "", "Direct string requirements prompt for code translation")
	flag.StringVar(&sourceFlag, "source", "", "Source codebase path (overrides config/prompt)")
	flag.StringVar(&targetFlag, "target", "", "Target output directory (overrides config/prompt)")
	flag.StringVar(&srcLangFlag, "source-lang", "", "Source language (overrides config/prompt)")
	flag.StringVar(&tgtLangFlag, "target-lang", "", "Target language (overrides config/prompt)")
	flag.StringVar(&toolchainFlag, "toolchain", "", "Target build/test toolchain (e.g. cargo, go, make)")
	flag.Parse()

	var cfg *config.Config
	var task types.TranslationTask

	// Prioritize YAML Configuration File
	if configFlag != "" {
		var err error
		cfg, err = config.LoadYAML(configFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading YAML config %s: %v\n", configFlag, err)
			os.Exit(1)
		}
		task = types.TranslationTask{
			SourceDir:   cfg.SourceDir,
			TargetDir:   cfg.TargetDir,
			SourceLang:  cfg.SourceLang,
			TargetLang:  cfg.TargetLang,
			Toolchain:   cfg.Toolchain,
			LSPProvider: cfg.LSPProvider,
		}
	} else if _, err := os.Stat("config.yml"); err == nil && promptFlag == "" && promptFileFlag == "" && len(flag.Args()) == 0 {
		cfg, err = config.LoadYAML("config.yml")
		if err == nil {
			task = types.TranslationTask{
				SourceDir:   cfg.SourceDir,
				TargetDir:   cfg.TargetDir,
				SourceLang:  cfg.SourceLang,
				TargetLang:  cfg.TargetLang,
				Toolchain:   cfg.Toolchain,
				LSPProvider: cfg.LSPProvider,
			}
		}
	} else if _, err := os.Stat("configs/config.yml"); err == nil && promptFlag == "" && promptFileFlag == "" && len(flag.Args()) == 0 {
		cfg, err = config.LoadYAML("configs/config.yml")
		if err == nil {
			task = types.TranslationTask{
				SourceDir:   cfg.SourceDir,
				TargetDir:   cfg.TargetDir,
				SourceLang:  cfg.SourceLang,
				TargetLang:  cfg.TargetLang,
				Toolchain:   cfg.Toolchain,
				LSPProvider: cfg.LSPProvider,
			}
		}
	}

	// If not loaded from YAML, process prompts
	if cfg == nil {
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
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				if stdinBytes, err := io.ReadAll(os.Stdin); err == nil && len(stdinBytes) > 0 {
					prompt = strings.TrimSpace(string(stdinBytes))
				}
			}
		}
		if prompt == "" {
			fmt.Fprintln(os.Stderr, "Error: No configuration or prompt provided. Specify --config <config.yml> or create config.yml in workspace.")
			flag.Usage()
			os.Exit(1)
		}

		task = parsePromptRequirements(prompt)
		cfg = config.MustLoad(
			config.WithSourceTarget(task.SourceDir, task.TargetDir, task.SourceLang, task.TargetLang),
		)
	}
	if sourceFlag != "" {
		task.SourceDir = sourceFlag
		cfg.SourceDir = sourceFlag
	}
	if targetFlag != "" {
		task.TargetDir = targetFlag
		cfg.TargetDir = targetFlag
	}
	if srcLangFlag != "" {
		task.SourceLang = srcLangFlag
		cfg.SourceLang = srcLangFlag
	}
	if tgtLangFlag != "" {
		task.TargetLang = tgtLangFlag
		cfg.TargetLang = tgtLangFlag
	}
	if toolchainFlag != "" {
		task.Toolchain = toolchainFlag
		cfg.Toolchain = toolchainFlag
	}

	if task.SourceDir == "" || task.TargetDir == "" {
		fmt.Fprintln(os.Stderr, "Error: Could not extract source codebase and output directory.")
		os.Exit(1)
	}

	fmt.Printf("Task Specification:\n")
	fmt.Printf("  Source Codebase: %s\n", task.SourceDir)
	fmt.Printf("  Output Dir:      %s\n", task.TargetDir)
	fmt.Printf("  Source Language: %s\n", task.SourceLang)
	fmt.Printf("  Target Language: %s\n", task.TargetLang)
	if task.Toolchain != "" {
		fmt.Printf("  Toolchain:       %s\n", task.Toolchain)
	}
	fmt.Println()
	ctx := context.Background()

	// Instantiate reasoning and coding model endpoints
	logger.LogStep("Connecting to Ollama models at %s...", cfg.OllamaBaseURL)
	logger.LogStep("Reasoning Model: %s", cfg.ReasoningModel)
	logger.LogStep("Coding Model:    %s", cfg.CodingModel)
	models := llm.MustNewModels(ctx, cfg)

	// Construct the Eino multi-agent graph with automated repair feedback
	logger.LogStep("Constructing 4-agent Eino Graph (Analyzer, Planning, Translator, Validator)...")
	magharcmGraph := graph.MustNewMAgHARCMGraph(ctx, models)

	// Initialize pipeline state with parsed task specification
	initialState := &types.State{
		Task:          task,
		MaxIterations: cfg.MaxIterations,
		TranslatedProject: types.TranslatedProject{
			Files: make(map[string]string),
		},
	}

	// Execute translation graph until validation criteria are met or iteration limit reached
	logger.LogStep("Starting multi-agent translation execution...")
	finalState, err := magharcmGraph.Execute(ctx, initialState)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing MAgHARCM pipeline: %v\n", err)
		os.Exit(1)
	}

	// Completion Status
	if finalState.ValidationReport.IsAllSuccess() {
		logger.LogAgent("MAgHARCM", "Translation and validation completed successfully: %s", finalState.ValidationReport.String())
		logger.LogStep("Target project ready in %s", filepath.Clean(cfg.TargetDir))
	} else {
		logger.LogWarning("Execution finished: %s", finalState.ValidationReport.String())
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
