package graph

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"

	"MAgHARCM/pkg/agents"
	"MAgHARCM/pkg/llm"
	"MAgHARCM/pkg/pattern"
	"MAgHARCM/pkg/types"
)

// MAgHARCMGraph wraps the compiled Eino runnable for the multi-agent pipeline.
type MAgHARCMGraph struct {
	Runnable compose.Runnable[*types.State, *types.State]
}

// NewMAgHARCMGraph builds and compiles the 4-agent Eino Graph with the repair loop (§3 & Algorithm 1).
func NewMAgHARCMGraph(ctx context.Context, models *llm.Models) (*MAgHARCMGraph, error) {
	g := compose.NewGraph[*types.State, *types.State]()

	// 1. Initialize the 4 agents
	analyzerAgent := agents.NewAnalyzerAgent(models.Reasoning)
	planningAgent := agents.NewPlanningAgent(models.Reasoning)
	translatorAgent := agents.NewTranslatorAgent(models.Coding)
	validatorAgent := agents.NewValidatorAgent(models.Reasoning)

	// 2. Add agent nodes to the graph
	pattern.Must0(g.AddLambdaNode("analyzer", compose.InvokableLambda(analyzerAgent.Run)))
	pattern.Must0(g.AddLambdaNode("planning", compose.InvokableLambda(planningAgent.Run)))
	pattern.Must0(g.AddLambdaNode("translator", compose.InvokableLambda(translatorAgent.Run)))
	pattern.Must0(g.AddLambdaNode("validator", compose.InvokableLambda(validatorAgent.Run)))

	// 3. Connect sequential flow: START -> analyzer -> planning -> translator -> validator
	pattern.Must0(g.AddEdge(compose.START, "analyzer"))
	pattern.Must0(g.AddEdge("analyzer", "planning"))
	pattern.Must0(g.AddEdge("planning", "translator"))
	pattern.Must0(g.AddEdge("translator", "validator"))

	// 4. Add conditional branch on validator: if success -> END, else -> loop back to translator
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

	// 5. Compile runnable graph with higher max steps for repair cycles
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
