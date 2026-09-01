package graph

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"

	"MAgHARCM/internal/agents"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/types"
)

// MAgHARCMGraph wraps the compiled Eino runnable for the multi-agent pipeline.
type MAgHARCMGraph struct {
	Runnable compose.Runnable[*types.State, *types.State]
	// RunID identifies the current translation run and is forwarded to the
	// translator and validator agents so they can persist per-iteration
	// checkpoints under .artifacts/<RunID>/checkpoints/. Empty disables
	// checkpointing.
	RunID string
}

// checkpointLambda is the lambda body for the save_translator_ckpt and
// save_validator_ckpt graph nodes. The actual save is performed inside the
// translator/validator Run via a defer that runs on every return (success or
// error). These lambda nodes exist as explicit graph checkpoints so the
// repair-cycle branch can re-enter the translator only after the validator
// checkpoint is on disk.
func checkpointLambda(_ context.Context, state *types.State) (*types.State, error) {
	return state, nil
}

// NewMAgHARCMGraph constructs and compiles the multi-agent cyclic graph with automated translation repair loop.
// runID is forwarded to the translator and validator so each iteration is
// checkpointed under .artifacts/<runID>/checkpoints/. Pass "" to disable
// checkpointing entirely.
func NewMAgHARCMGraph(ctx context.Context, models *llm.Models, runID string) (*MAgHARCMGraph, error) {
	g := compose.NewGraph[*types.State, *types.State]()

	// Initialize reasoning and coding agent instances
	analyzerAgent := agents.NewAnalyzerAgent(models.Reasoning)
	planningAgent := agents.NewPlanningAgent(models.Reasoning)
	translatorAgent := agents.NewTranslatorAgent(models.Coding, runID)
	validatorAgent := agents.NewValidatorAgent(models.Reasoning, runID)

	// Register agent execution units as graph nodes
	if err := g.AddLambdaNode("analyzer", compose.InvokableLambda(analyzerAgent.Run)); err != nil {
		return nil, err
	}
	if err := g.AddLambdaNode("planning", compose.InvokableLambda(planningAgent.Run)); err != nil {
		return nil, err
	}
	if err := g.AddLambdaNode("translator", compose.InvokableLambda(translatorAgent.Run)); err != nil {
		return nil, err
	}
	if err := g.AddLambdaNode("validator", compose.InvokableLambda(validatorAgent.Run)); err != nil {
		return nil, err
	}
	// Checkpoint lambda nodes: each is a no-op forwarder placed after the
	// agent that has already persisted a snapshot via its own defer. They
	// give the repair branch a stable anchor node so the cycle resumes only
	// after the previous iteration's checkpoint is durable on disk.
	if err := g.AddLambdaNode("save_translator_ckpt", compose.InvokableLambda(checkpointLambda)); err != nil {
		return nil, err
	}
	if err := g.AddLambdaNode("save_validator_ckpt", compose.InvokableLambda(checkpointLambda)); err != nil {
		return nil, err
	}

	// Wire the primary forward pipeline from analysis through validation
	if err := g.AddEdge(compose.START, "analyzer"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("analyzer", "planning"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("planning", "translator"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("translator", "save_translator_ckpt"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("save_translator_ckpt", "validator"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("validator", "save_validator_ckpt"); err != nil {
		return nil, err
	}
	// Implement dynamic repair feedback loop returning to translator on validation failure.
	// Branch off save_validator_ckpt so the cycle resumes only after the
	// validator checkpoint is durable.
	repairBranch := compose.NewGraphBranch(
		func(ctx context.Context, state *types.State) (string, error) {
			if state.IsComplete || state.ValidationReport.IsAllSuccess() || state.Iteration >= state.MaxIterations {
				logger.LogStep("Graph pipeline terminating: complete=%v, all_success=%v, iteration=%d/%d",
					state.IsComplete, state.ValidationReport.IsAllSuccess(), state.Iteration, state.MaxIterations)
				return compose.END, nil
			}
			logger.LogStep("Validation incomplete, cycling back to translator for repair (iteration %d/%d)",
				state.Iteration, state.MaxIterations)
			return "translator", nil
		},
		map[string]bool{
			compose.END:  true,
			"translator": true,
		},
	)
	if err := g.AddBranch("save_validator_ckpt", repairBranch); err != nil {
		return nil, err
	}

	// Compile graph with step budget accommodating multiple repair iterations
	runnable, err := g.Compile(ctx,
		compose.WithGraphName("MAgHARCM"),
		compose.WithMaxRunSteps(50),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to compile MAgHARCM graph: %w", err)
	}

	return &MAgHARCMGraph{Runnable: runnable, RunID: runID}, nil
}


// Execute runs the translation graph with initial state and returns final state.
func (rg *MAgHARCMGraph) Execute(ctx context.Context, initialState *types.State) (*types.State, error) {
	return rg.Runnable.Invoke(ctx, initialState)
}
