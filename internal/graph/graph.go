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
}

// NewMAgHARCMGraph constructs and compiles the multi-agent cyclic graph with automated translation repair loop.
func NewMAgHARCMGraph(ctx context.Context, models *llm.Models) (*MAgHARCMGraph, error) {
	g := compose.NewGraph[*types.State, *types.State]()

	// Initialize reasoning and coding agent instances
	analyzerAgent := agents.NewAnalyzerAgent(models.Reasoning)
	planningAgent := agents.NewPlanningAgent(models.Reasoning)
	translatorAgent := agents.NewTranslatorAgent(models.Coding)
	validatorAgent := agents.NewValidatorAgent(models.Reasoning)

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
	if err := g.AddEdge("translator", "validator"); err != nil {
		return nil, err
	}
	// Implement dynamic repair feedback loop returning to translator on validation failure
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
	if err := g.AddBranch("validator", repairBranch); err != nil {
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

	return &MAgHARCMGraph{Runnable: runnable}, nil
}

// MustNewMAgHARCMGraph constructs the graph and panics on failure.
func MustNewMAgHARCMGraph(ctx context.Context, models *llm.Models) *MAgHARCMGraph {
	g, err := NewMAgHARCMGraph(ctx, models)
	if err != nil {
		panic(err)
	}
	return g
}

// Execute runs the translation graph with initial state and returns final state.
func (rg *MAgHARCMGraph) Execute(ctx context.Context, initialState *types.State) (*types.State, error) {
	return rg.Runnable.Invoke(ctx, initialState)
}
