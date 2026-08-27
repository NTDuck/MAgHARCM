package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"MAgHARCM/internal/config"
	"MAgHARCM/internal/graph"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/types"
)

func main() {
	var configFlag string
	var sourceFlag string
	var targetFlag string
	var srcLangFlag string
	var tgtLangFlag string
	var toolchainFlag string

	flag.StringVar(&configFlag, "config", "config.yml", "Path to YAML configuration file")
	flag.StringVar(&sourceFlag, "source", "", "Source codebase path (overrides config)")
	flag.StringVar(&targetFlag, "target", "", "Target output directory (overrides config)")
	flag.StringVar(&srcLangFlag, "source-lang", "", "Source language (overrides config)")
	flag.StringVar(&tgtLangFlag, "target-lang", "", "Target language (overrides config)")
	flag.StringVar(&toolchainFlag, "toolchain", "", "Target build/test toolchain (e.g. cargo, go, make)")
	flag.Parse()

	// Resolve configuration from YAML
	configFile := configFlag
	if _, err := os.Stat(configFile); err != nil && configFile == "config.yml" {
		if _, err := os.Stat("configs/config.yml"); err == nil {
			configFile = "configs/config.yml"
		}
	}

	cfg, err := config.LoadYAML(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load YAML config from `%s`: %v\n", configFile, err)
		os.Exit(1)
	}

	// Apply CLI overrides if specified
	if sourceFlag != "" {
		cfg.SourceDir = sourceFlag
	}
	if targetFlag != "" {
		cfg.TargetDir = targetFlag
	}
	if srcLangFlag != "" {
		cfg.SourceLang = srcLangFlag
	}
	if tgtLangFlag != "" {
		cfg.TargetLang = tgtLangFlag
	}
	if toolchainFlag != "" {
		cfg.Toolchain = toolchainFlag
	}

	task := types.TranslationTask{
		SourceDir:   cfg.SourceDir,
		TargetDir:   cfg.TargetDir,
		SourceLang:  cfg.SourceLang,
		TargetLang:  cfg.TargetLang,
		Toolchain:   cfg.Toolchain,
		LSPProvider: cfg.LSPProvider,
	}

	if task.SourceDir == "" || task.TargetDir == "" {
		fmt.Fprintln(os.Stderr, "Error: source_dir and target_dir must be specified in config.yml or via flags.")
		os.Exit(1)
	}


	fmt.Println("================================================================")
	fmt.Println("       MAgHARCM - Multi-Agent Code Translation Engine           ")
	fmt.Println("================================================================")
	fmt.Printf("  Configuration:   %s\n", configFile)
	fmt.Printf("  Source Codebase: %s (%s)\n", task.SourceDir, task.SourceLang)
	fmt.Printf("  Target Directory:%s (%s)\n", task.TargetDir, task.TargetLang)
	if task.Toolchain != "" {
		fmt.Printf("  Toolchain:       %s\n", task.Toolchain)
	}
	fmt.Println()

	ctx := context.Background()

	// Instantiate reasoning and coding model endpoints
	logger.LogStep("Connecting to Ollama models at `%s`", cfg.OllamaBaseURL)
	logger.LogStep("Reasoning Model: `%s`", cfg.ReasoningModel)
	logger.LogStep("Coding Model:    `%s`", cfg.CodingModel)
	models := llm.MustNewModels(ctx, cfg)

	// Construct the Eino multi-agent graph with automated repair feedback
	logger.LogStep("Constructing 4-agent Eino Graph (Analyzer, Planning, Translator, Validator)")
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
	logger.LogStep("Starting multi-agent translation execution")
	finalState, err := magharcmGraph.Execute(ctx, initialState)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing MAgHARCM pipeline: %v\n", err)
		os.Exit(1)
	}

	// Completion Status
	if finalState.ValidationReport.IsAllSuccess() {
		logger.LogAgent("MAgHARCM", "Translation and validation completed successfully: %s", finalState.ValidationReport.String())
		logger.LogStep("Target project ready in `%s`", filepath.Clean(cfg.TargetDir))
	} else {
		logger.LogWarning("Execution finished: %s", finalState.ValidationReport.String())
		os.Exit(1)
	}
}
