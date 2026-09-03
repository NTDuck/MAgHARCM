// Package runner executes the MAgHARCM translation pipeline given a
//
//	populated *config.Config. Both cmd/MAgHARCM (one-shot) and cmd/MAgHARCM-tui
//
// (interactive REPL) call into Run so the wiring stays in one place.
package runner

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"MAgHARCM/internal/agents"
	"MAgHARCM/internal/artifacts"
	"MAgHARCM/internal/config"
	"MAgHARCM/internal/graph"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/types"
)

// ErrMissingFields is returned when cfg is nil.
var ErrMissingFields = errors.New("cfg is required")

// Run executes the full analyzer -> planning -> translator -> validator
// pipeline and returns the final state. cfg MUST have every required field
// populated; missing fields produce a structured error from config.Require.
func Run(ctx context.Context, cfg *config.Config) (*types.State, error) {
	if cfg == nil {
		return nil, ErrMissingFields
	}
	if err := config.Require(cfg); err != nil {
		return nil, err
	}

	task := types.TranslationTask{
		SourceDir:   cfg.SourceDir,
		TargetDir:   cfg.TargetDir,
		SourceLang:  cfg.SourceLang,
		TargetLang:  cfg.TargetLang,
		Toolchain:   cfg.Toolchain,
		LSPProvider: cfg.LSPProvider,
	}

	logger.LogStep("Run: source=%s (%s) -> target=%s (%s) toolchain=%s",
		task.SourceDir, task.SourceLang, task.TargetDir, task.TargetLang, task.Toolchain)

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
			TranslatedProject: artifacts.TranslatedProject{
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

	logger.LogStep("Constructing 5-agent Eino Graph (Analyzer, Navigator, Planning, Translator, Validator)")
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
