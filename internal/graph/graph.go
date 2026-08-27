package graph

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"

	"MAgHARCM/internal/agents"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/pattern"
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
	pattern.Must0(g.AddLambdaNode("analyzer", compose.InvokableLambda(analyzerAgent.Run)))
	pattern.Must0(g.AddLambdaNode("planning", compose.InvokableLambda(planningAgent.Run)))
	pattern.Must0(g.AddLambdaNode("translator", compose.InvokableLambda(translatorAgent.Run)))
	pattern.Must0(g.AddLambdaNode("validator", compose.InvokableLambda(validatorAgent.Run)))

	// Wire the primary forward pipeline from analysis through validation
	pattern.Must0(g.AddEdge(compose.START, "analyzer"))
	pattern.Must0(g.AddEdge("analyzer", "planning"))
	pattern.Must0(g.AddEdge("planning", "translator"))
	pattern.Must0(g.AddEdge("translator", "validator"))

	// Implement dynamic repair feedback loop returning to translator on validation failure
	repairBranch := compose.NewGraphBranch(
		func(ctx context.Context, state *types.State) (string, error) {
			if state.IsComplete || state.ValidationReport.IsAllSuccess() || state.Iteration >= state.MaxIterations {
				state.Log("[Graph] Pipeline terminating: complete=%v, all_success=%v, iteration=%d/%d",
					state.IsComplete, state.ValidationReport.IsAllSuccess(), state.Iteration, state.MaxIterations)
				return compose.END, nil
			}
			state.Log("[Graph] Validation incomplete, cycling back to translator for repair (iteration %d/%d)",
				state.Iteration, state.MaxIterations)
			return "translator", nil
		},
		map[string]bool{
			compose.END:  true,
			"translator": true,
		},
	)
	pattern.Must0(g.AddBranch("validator", repairBranch))

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
	return pattern.Must(NewMAgHARCMGraph(ctx, models))
}

// Execute runs the translation graph with initial state and returns final state.
func (rg *MAgHARCMGraph) Execute(ctx context.Context, initialState *types.State) (*types.State, error) {
	return rg.Runnable.Invoke(ctx, initialState)
}
