// Package runner executes the MAgHARCM translation pipeline given a
//	populated *config.Config. Both cmd/MAgHARCM (one-shot) and cmd/MAgHARCM-tui
// (interactive REPL) call into Run so the wiring stays in one place.
package runner

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"MAgHARCM/internal/agents"
	"MAgHARCM/internal/config"
	"MAgHARCM/internal/graph"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/types"
)

// ErrMissingFields is returned when the config lacks the source/target dirs
// that the pipeline cannot derive from defaults.
var ErrMissingFields = errors.New("source_dir and target_dir are required")

// Run executes the full analyzer -> planning -> translator -> validator
// pipeline and returns the final state. cfg must have SourceDir and
// TargetDir set; other fields fall back to config.Defaults().
func Run(ctx context.Context, cfg *config.Config) (*types.State, error) {
	if cfg == nil {
		return nil, ErrMissingFields
	}
	if cfg.SourceDir == "" || cfg.TargetDir == "" {
		return nil, ErrMissingFields
	}

	def := config.Defaults()
	if cfg.OllamaBaseURL == "" {
		cfg.OllamaBaseURL = def.OllamaBaseURL
	}
	if cfg.ReasoningModel == "" {
		cfg.ReasoningModel = def.ReasoningModel
	}
	if cfg.CodingModel == "" {
		cfg.CodingModel = def.CodingModel
	}
	if cfg.MaxIterations == 0 {
		cfg.MaxIterations = def.MaxIterations
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = def.Timeout
	}

	task := types.TranslationTask{
		SourceDir:   cfg.SourceDir,
		TargetDir:   cfg.TargetDir,
		SourceLang:  cfg.SourceLang,
		TargetLang:  cfg.TargetLang,
		Toolchain:   cfg.Toolchain,
		LSPProvider: cfg.LSPProvider,
	}

	fmt.Println("================================================================")
	fmt.Println("       MAgHARCM - Multi-Agent Code Translation Engine           ")
	fmt.Println("================================================================")
	fmt.Printf("  Source Codebase: %s (%s)\n", task.SourceDir, task.SourceLang)
	fmt.Printf("  Target Directory:%s (%s)\n", task.TargetDir, task.TargetLang)
	if task.Toolchain != "" {
		fmt.Printf("  Toolchain:       %s\n", task.Toolchain)
	}
	fmt.Println()

	runID := agents.RunIDForTask(task)
	logger.LogStep("Run ID: %s", runID)

	resumed, err := agents.LoadLatest(runID)
	if err != nil {
		return nil, fmt.Errorf("load checkpoint for %s: %w", runID, err)
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

	logger.LogStep("Connecting to Ollama models at `%s`", cfg.OllamaBaseURL)
	logger.LogStep("Reasoning Model: `%s`", cfg.ReasoningModel)
	logger.LogStep("Coding Model:    `%s`", cfg.CodingModel)

	models, err := llm.NewModels(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize Ollama models: %w", err)
	}

	logger.LogStep("Constructing 4-agent Eino Graph (Analyzer, Planning, Translator, Validator)")
	magharcmGraph, err := graph.NewMAgHARCMGraph(ctx, models, runID)
	if err != nil {
		return nil, fmt.Errorf("construct graph: %w", err)
	}

	logger.LogStep("Starting multi-agent translation execution")
	finalState, err := magharcmGraph.Execute(ctx, initialState)
	if err != nil {
		return nil, fmt.Errorf("execute pipeline: %w", err)
	}

	if finalState.ValidationReport.IsAllSuccess() {
		if err := agents.Cleanup(runID); err != nil {
			logger.LogWarning("Failed to clean up checkpoints for run `%s`: %v", runID, err)
		}
		logger.LogAgent("MAgHARCM", "Translation and validation completed successfully: %s", finalState.ValidationReport.String())
		logger.LogStep("Target project ready in `%s`", filepath.Clean(cfg.TargetDir))
	} else {
		logger.LogWarning("Execution finished: %s", finalState.ValidationReport.String())
	}

	return finalState, nil
}

// Success reports whether the final state cleared validation.
func Success(s *types.State) bool {
	return s != nil && s.ValidationReport.IsAllSuccess()
}
