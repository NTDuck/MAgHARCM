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
	"MAgHARCM/internal/checkpoint"
	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/config"
	"MAgHARCM/internal/graph"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/logger"
)

// ErrMissingFields is returned when cfg is nil.
var ErrMissingFields = errors.New("cfg is required")

// Run executes the full analyzer -> planning -> translator -> validator
// pipeline and returns the final state. cfg MUST have every required field
func Run(ctx context.Context, cfg *config.Config) (*compiletime.State, error) {
	if cfg == nil {
		return nil, ErrMissingFields
	}
	if err := config.Require(cfg); err != nil {
		return nil, err
	}

	task := compiletime.Task{
		SourceDir:   cfg.SourceDir,
		TargetDir:   cfg.TargetDir,
		SourceLang:  cfg.SourceLang,
		TargetLang:  cfg.TargetLang,
		Toolchain:   cfg.Toolchain,
		LSPProvider: cfg.LSPProvider,
	}

	logger.LogStep("Run: source=%s srcLang=%s target=%s tgtLang=%s toolchain=%s",
		task.SourceDir, task.SourceLang, task.TargetDir, task.TargetLang, task.Toolchain)

	runID := checkpoint.RunIDForSourceDir(task.SourceDir)
	logger.LogStep("Run ID: %s", runID)

	resumed, err := checkpoint.LoadLatest(runID)
	if err != nil {
		return nil, fmt.Errorf("load checkpoint for %s: %w", runID, err)
	}
	var initialState *compiletime.State
	if resumed != nil && resumed.State != nil {
		logger.LogStep("Resume from checkpoint iter-%d", resumed.Iteration)
		// cfg wins on resume: the checkpoint may predate a config change
		// (target_dir, lsp.provider, model set, iteration budget). Task is
		// re-bound and MaxIterations refreshed; all other artifact state
		// (analysis, plan, translated files, validation report) is carried
		// over untouched.
		resumed.State.Task = task
		resumed.State.MaxIterations = cfg.MaxIterations
		initialState = resumed.State
	} else {
		if resumed != nil {
			logger.LogWarning("Checkpoint iter-%d has nil state; starting fresh run instead of resuming", resumed.Iteration)
		}
		initialState = &compiletime.State{
			Task:          task,
			MaxIterations: cfg.MaxIterations,
			TranslatedProject: compiletime.TranslatedProject{
				ArtifactSchemaVersion: compiletime.CurrentSchemaVersion,
				Files:                 make(map[string]string),
			},
		}
	}

	logger.LogStep("Connect to Ollama models at `%s`", cfg.OllamaBaseURL)
	logger.LogStep("Reasoning Model: `%s`", cfg.ReasoningModel)
	logger.LogStep("Coding Model:    `%s`", cfg.CodingModel)

	models, err := llm.NewModels(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize Ollama models: %w", err)
	}

	logger.LogStep("Build 8-agent Eino Graph: Archaeologist, Analyzer, Planner, Translator, Reviewer, Validator, VerdictPanel, Recruiter")
	magharcmGraph, err := graph.NewMAgHARCMGraph(ctx, models, runID)
	if err != nil {
		return nil, fmt.Errorf("construct graph: %w", err)
	}

	logger.LogStep("Start multi-agent translation execution")
	finalState, err := magharcmGraph.Execute(ctx, initialState)
	if err != nil {
		return nil, fmt.Errorf("execute pipeline: %w", err)
	}

	if finalState.ValidationReport.IsAllSuccess() {
		if err := checkpoint.Cleanup(runID); err != nil {
			logger.LogWarning("Cannot remove checkpoints for run `%s`: %v", runID, err)
		}
		logger.LogAgent("MAgHARCM", "Translation and validation completed: %s", finalState.ValidationReport.String())
		logger.LogStep("Target project ready in `%s`", filepath.Clean(cfg.TargetDir))
	} else {
		logger.LogWarning("Execution stopped: %s", finalState.ValidationReport.String())
		// PRIM-29 Recruitment-Adaptive Planning (AgentVerse style):
		// Evaluate validation metrics through Recruiter to derive next-iteration plan.
		recruiter := agents.NewRecruiter()
		summary := agents.ValidationSummary{
			CompilationSuccess:           finalState.ValidationReport.CompilationSuccess,
			PassRate:                     finalState.ValidationReport.TestPassRate,
			PlateauDetected:              finalState.ValidationReport.PlateauDetected,
			AdversarialWeakeningDetected: false,
		}
		if plan, err := recruiter.Recruit(ctx, agents.Profile{}, summary); err == nil {
			logger.LogAgent("Recruiter", "Plan for next iteration: %s tools=%v agents=%v",
				plan.Rationale, plan.Tools, plan.Agents)
		}
	}
	return finalState, nil
}

// Success reports whether the final state cleared validation.
func Success(s *compiletime.State) bool {
	return s != nil && s.ValidationReport.IsAllSuccess()
}
