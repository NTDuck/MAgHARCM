package graph

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"

	"MAgHARCM/internal/agents"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/logger"
)

// MAgHARCMGraph wraps the compiled Eino runnable for the multi-agent pipeline.
type MAgHARCMGraph struct {
	Runnable compose.Runnable[*agents.State, *agents.State]
	// RunID identifies the current translation run for disk checkpoints.
	RunID string
}

// checkpointLambda provides an explicit graph synchronization barrier
// so downstream branches execute only after previous state is durable.
func checkpointLambda(_ context.Context, state *agents.State) (*agents.State, error) {
	return state, nil
}

// NewMAgHARCMGraph constructs and compiles the 8-agent cyclic graph with automated repair.
func NewMAgHARCMGraph(ctx context.Context, models *llm.Models, runID string) (*MAgHARCMGraph, error) {
	g := compose.NewGraph[*agents.State, *agents.State]()

	// 1. Initialize independent agent execution units
	archaeologistAgent := agents.NewArchaeologist()
	analyzerAgent := agents.NewAnalyzerAgent(models.Reasoning)
	planningAgent := agents.NewPlanningAgent(models.Reasoning)
	translatorAgent := agents.NewTranslatorAgent(models.Coding, runID)
	reviewerAgent := agents.NewRoleFlipGate(models.Reasoning)
	validatorAgent := agents.NewValidatorAgent(models.Reasoning, runID)
	verdictPanelAgent := agents.NewVerdictPanel(models.Reasoning)
	recruiterAgent := agents.NewRecruiter()

	// 2. Register agent nodes in the Eino Graph

	// Node 1: Archaeologist (PRIM-14, 18, 19, 20, 22)
	if err := g.AddLambdaNode("archaeologist", compose.InvokableLambda(func(ctx context.Context, state *agents.State) (*agents.State, error) {
		logger.LogAgent("Archaeologist", "Starting pre-planning software archaeology")
		if state.Task.SourceDir != "" {
			rep, err := archaeologistAgent.Investigate(ctx, state.Task.SourceDir)
			if err != nil {
				logger.LogWarning("Archaeologist investigation had warnings: %v", err)
			} else {
				state.ArchaeologyReport = rep
				logger.LogAgent("Archaeologist", "Excavation complete: %d boundaries, %d churn hotspots, %d forensics",
					len(rep.BoundaryMap), len(rep.ChurnHotspots), len(rep.NamingForensics))
			}
		}
		return state, nil
	})); err != nil {
		return nil, err
	}

	// Node 2: Analyzer
	if err := g.AddLambdaNode("analyzer", compose.InvokableLambda(func(ctx context.Context, state *agents.State) (*agents.State, error) {
		models.PrepareReasoning()
		return analyzerAgent.Run(ctx, state)
	})); err != nil {
		return nil, err
	}

	// Node 3: Planner (PRIM-1, 2, 3)
	if err := g.AddLambdaNode("planning", compose.InvokableLambda(func(ctx context.Context, state *agents.State) (*agents.State, error) {
		models.PrepareReasoning()
		return planningAgent.Run(ctx, state)
	})); err != nil {
		return nil, err
	}

	// Node 4: Translator (PRIM-23, 31)
	if err := g.AddLambdaNode("translator", compose.InvokableLambda(func(ctx context.Context, state *agents.State) (*agents.State, error) {
		models.PrepareCoding()
		return translatorAgent.Run(ctx, state)
	})); err != nil {
		return nil, err
	}

	// Node 5: Checkpoint after translation
	if err := g.AddLambdaNode("save_translator_ckpt", compose.InvokableLambda(checkpointLambda)); err != nil {
		return nil, err
	}

	// Node 6: Reviewer (PRIM-25 Role-Flip Gate)
	if err := g.AddLambdaNode("reviewer", compose.InvokableLambda(func(ctx context.Context, state *agents.State) (*agents.State, error) {
		models.PrepareReasoning()
		if len(state.TranslatedProject.Files) > 0 {
			// Select a representative sample for role-flip inspection
			var sample string
			for _, content := range state.TranslatedProject.Files {
				sample = content
				break
			}
			verdict, err := reviewerAgent.Inspect(ctx, sample)
			if err != nil {
				logger.LogWarning("Reviewer role-flip gate error: %v", err)
			} else if verdict.DefectFound {
				logger.LogWarning("Reviewer detected potential defect: %s", verdict.Reason)
			} else {
				logger.LogAgent("Reviewer", "Role-flip sanity check passed")
			}
		}
		return state, nil
	})); err != nil {
		return nil, err
	}

	// Node 7: Validator (PRIM-5, 6, 13, 27)
	if err := g.AddLambdaNode("validator", compose.InvokableLambda(func(ctx context.Context, state *agents.State) (*agents.State, error) {
		models.PrepareReasoning()
		return validatorAgent.Run(ctx, state)
	})); err != nil {
		return nil, err
	}

	// Node 8: Checkpoint after validation
	if err := g.AddLambdaNode("save_validator_ckpt", compose.InvokableLambda(checkpointLambda)); err != nil {
		return nil, err
	}
	// Node 9: Verdict Panel (PRIM-7 Multi-Agent Consensus)
	if err := g.AddLambdaNode("verdict_panel", compose.InvokableLambda(func(ctx context.Context, state *agents.State) (*agents.State, error) {
		models.PrepareReasoning()
		if !state.ValidationReport.IsAllSuccess() && len(state.TranslatedProject.Files) > 0 {
			var sampleTarget string
			for _, c := range state.TranslatedProject.Files {
				sampleTarget = c
				break
			}
			res, err := verdictPanelAgent.Judge(ctx, "legacy source code", sampleTarget, agents.DefaultPanelSize)
			if err != nil {
				logger.LogWarning("Verdict panel error: %v", err)
			} else {
				logger.LogAgent("VerdictPanel", "Majority decision: agree=%v, disagreements=%d",
					res.Agree, len(res.Disagreements))
			}
		}
		return state, nil
	})); err != nil {
		return nil, err
	}

	// Node 10: Recruiter (PRIM-29 Dynamic Iteration Adaptation)
	if err := g.AddLambdaNode("recruiter", compose.InvokableLambda(func(ctx context.Context, state *agents.State) (*agents.State, error) {
		if !state.ValidationReport.IsAllSuccess() {
			// Trigger try-and-fail migration strategy switch if needed
			if state.ValidationReport.PlateauDetected || len(state.ValidationReport.CompilationErrors) > 0 {
				if next, switched := agents.SwitchToNextStrategy(ctx, state); switched {
					logger.LogAgent("Recruiter", "Activated next migration strategy: %s", next)
				}
			}
			summary := agents.ValidationSummary{
				CompilationSuccess:           state.ValidationReport.CompilationSuccess,
				PassRate:                     state.ValidationReport.TestPassRate,
				PlateauDetected:              state.ValidationReport.PlateauDetected,
				AdversarialWeakeningDetected: state.ValidationReport.AdversarialWeakeningDetected,
			}
			plan, err := recruiterAgent.Recruit(ctx, agents.Profile{}, summary)
			if err == nil {
				logger.LogAgent("Recruiter", "Adaptive plan for repair: %s (tools: %v)", plan.Rationale, plan.Tools)
			}
		}
		return state, nil
	})); err != nil {
		return nil, err
	}

	// 3. Connect Forward Pipeline Edges
	if err := g.AddEdge(compose.START, "archaeologist"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("archaeologist", "analyzer"); err != nil {
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
	if err := g.AddEdge("save_translator_ckpt", "reviewer"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("reviewer", "validator"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("validator", "save_validator_ckpt"); err != nil {
		return nil, err
	}

	// 4. Connect Validation Branching & Cyclic Repair Loop
	repairBranch := compose.NewGraphBranch(
		func(ctx context.Context, state *agents.State) (string, error) {
			if state.IsComplete || state.ValidationReport.IsAllSuccess() || state.Iteration >= state.MaxIterations {
				logger.LogStep("Pipeline termination condition met: complete=%v, all_success=%v, iteration=%d/%d",
					state.IsComplete, state.ValidationReport.IsAllSuccess(), state.Iteration, state.MaxIterations)
				return compose.END, nil
			}
			logger.LogStep("Validation incomplete; routing to verdict panel and recruiter for repair (iteration %d/%d)",
				state.Iteration, state.MaxIterations)
			return "verdict_panel", nil
		},
		map[string]bool{
			compose.END:     true,
			"verdict_panel": true,
		},
	)
	if err := g.AddBranch("save_validator_ckpt", repairBranch); err != nil {
		return nil, err
	}

	// Connect repair branch sequence: verdict_panel -> recruiter -> translator
	if err := g.AddEdge("verdict_panel", "recruiter"); err != nil {
		return nil, err
	}
	if err := g.AddEdge("recruiter", "translator"); err != nil {
		return nil, err
	}

	// 5. Compile Graph with safety run-step ceiling
	runnable, err := g.Compile(ctx,
		compose.WithGraphName("MAgHARCM-8Agent"),
		compose.WithMaxRunSteps(50),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to compile MAgHARCM graph: %w", err)
	}

	return &MAgHARCMGraph{Runnable: runnable, RunID: runID}, nil
}

// Execute runs the translation graph with initial state and returns final state.
func (rg *MAgHARCMGraph) Execute(ctx context.Context, initialState *agents.State) (*agents.State, error) {
	return rg.Runnable.Invoke(ctx, initialState)
}
