package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"MAgHARCM/internal/agents"
	"MAgHARCM/internal/config"
	"MAgHARCM/internal/graph"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/types"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "config.yml", "Path to YAML configuration file")
	flag.Parse()

	cfg, err := config.LoadYAML(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load YAML configuration from `%s`: %v\n", configFile, err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Error: source_dir and target_dir must be configured in `%s`.\n", configFile)
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

	// Resume from a prior checkpoint if one exists for this run. The run ID is
	// derived from the source directory so re-running MAgHARCM against the same
	// source repo (e.g. after a VRAM-OOM crash mid-translation) automatically
	// picks up where the last run left off. A fresh source dir produces a fresh
	// run ID and skips checkpoint loading.
	runID := agents.RunIDForTask(task)
	logger.LogStep("Run ID: %s", runID)
	resumed, err := agents.LoadLatest(runID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load checkpoint for run `%s`: %v\n", runID, err)
		os.Exit(1)
	}
	var initialState *types.State
	if resumed != nil {
		logger.LogStep("Resuming from checkpoint iter-%d", resumed.Iteration)
		initialState = resumed.State
	} else {
		initialState = &types.State{
			Task:          task,
			MaxIterations: cfg.MaxIterations,
			TranslatedProject: types.TranslatedProject{
				Files: make(map[string]string),
			},
		}
	}

	// Instantiate reasoning and coding model endpoints
	logger.LogStep("Connecting to Ollama models at `%s`", cfg.OllamaBaseURL)
	logger.LogStep("Reasoning Model: `%s`", cfg.ReasoningModel)
	logger.LogStep("Coding Model:    `%s`", cfg.CodingModel)
	models, err := llm.NewModels(ctx, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Ollama models: %v\n", err)
		os.Exit(1)
	}

	logger.LogStep("Constructing 4-agent Eino Graph (Analyzer, Planning, Translator, Validator)")
	magharcmGraph, err := graph.NewMAgHARCMGraph(ctx, models, runID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error constructing MAgHARCM graph: %v\n", err)
		os.Exit(1)
	}

	// Execute translation graph until validation criteria are met or iteration limit reached
	logger.LogStep("Starting multi-agent translation execution")
	finalState, err := magharcmGraph.Execute(ctx, initialState)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing MAgHARCM pipeline: %v\n", err)
		os.Exit(1)
	}

	// On successful validation, drop the checkpoint files so subsequent runs
	// of the same source dir start clean. Failed runs keep their checkpoints
	// so the next invocation can resume.
	if finalState.ValidationReport.IsAllSuccess() {
		if err := agents.Cleanup(runID); err != nil {
			logger.LogWarning("Failed to clean up checkpoints for run `%s`: %v", runID, err)
		}
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
