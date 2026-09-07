package agents

import (
	"context"
	"testing"

	"MAgHARCM/internal/compiletime"
)

func TestEvidenceFirstAdaptor(t *testing.T) {
	adaptor := NewEvidenceFirstAdaptor()
	report, err := adaptor.Adapt(context.Background(), "LegacyAlgorithm", "void run() { unsafe { free(p); } }", "Rust")
	if err != nil {
		t.Fatalf("Adapt failed: %v", err)
	}
	if report.Spec.SourceUnit != "LegacyAlgorithm" {
		t.Errorf("expected source unit LegacyAlgorithm, got %s", report.Spec.SourceUnit)
	}
	if !report.SafetyModel.ContainsUnsafeBlocks {
		t.Errorf("expected safety model to detect unsafe blocks")
	}
	if !report.SafetyModel.PotentialMemoryLeaks {
		t.Errorf("expected safety model to detect potential memory leaks")
	}
}

func TestSpecLifecycleManager(t *testing.T) {
	manager := NewSpecLifecycleManager()
	if manager.CurrentPhase != PhaseConstitution {
		t.Fatalf("expected initial phase Constitution, got %s", manager.CurrentPhase)
	}

	// Fail empty artifact
	if _, err := manager.Advance("empty", ""); err == nil {
		t.Errorf("expected error on empty artifact")
	}

	// Advance Constitution -> Specify
	next, err := manager.Advance("spec.md", "Specification contents")
	if err != nil {
		t.Fatalf("advance to specify failed: %v", err)
	}
	if next != PhaseSpecify {
		t.Errorf("expected phase Specify, got %s", next)
	}

	// Advance Specify -> Plan
	next, err = manager.Advance("plan.md", "Plan contents")
	if err != nil {
		t.Fatalf("advance to plan failed: %v", err)
	}
	if next != PhasePlan {
		t.Errorf("expected phase Plan, got %s", next)
	}
}

func TestJaccardCouplingAnalyzer(t *testing.T) {
	analyzer := NewJaccardCouplingAnalyzer()
	fileCommits := map[string]map[string]struct{}{
		"a.go": {"c1": {}, "c2": {}, "c3": {}},
		"b.go": {"c2": {}, "c3": {}, "c4": {}},
		"c.go": {"c5": {}},
	}

	edges := analyzer.ComputeCoupling(context.Background(), fileCommits, 0.4)
	if len(edges) != 1 {
		t.Fatalf("expected 1 coupling edge above threshold 0.4, got %d", len(edges))
	}
	// Intersection(a, b) = 2 (c2, c3), Union(a, b) = 4 (c1, c2, c3, c4), J = 2/4 = 0.5
	if edges[0].Similarity != 0.5 {
		t.Errorf("expected similarity 0.5, got %f", edges[0].Similarity)
	}
}

func TestDRHierarchyPartitioner(t *testing.T) {
	partitioner := NewDRHierarchyPartitioner()
	elements := []string{"interface.go", "service.go", "leaf.go"}
	deps := map[string][]string{
		"service.go":   {"interface.go", "leaf.go"},
		"interface.go": {"leaf.go"}, // Intentionally illegal dependency for testing violations
		"leaf.go":      {},
	}

	hierarchy := partitioner.Partition(elements, deps)
	if len(hierarchy.L1Interfaces) == 0 && len(hierarchy.L2Subsystems) == 0 {
		t.Fatalf("partition produced empty hierarchy")
	}
}

func TestConceptAssigner(t *testing.T) {
	assigner := NewConceptAssigner()
	files := map[string]string{
		"calc.go": "func calculateMean(data []float64) float64 { return stats.Mean(data) }",
		"auth.go": "func validateToken(token string) bool { return checkSignature(token) }",
	}

	report := assigner.AssignConcepts(context.Background(), files)
	if len(report.Concepts) == 0 {
		t.Fatalf("expected discovered concepts, got none")
	}
}

func TestComprehensionPipeline(t *testing.T) {
	comp := NewComprehensionPipeline()
	files := []string{"cmd/main.go", "pkg/math/calc.go", "pkg/math/calc_test.go"}
	contents := map[string]string{
		"cmd/main.go":          "import fmt",
		"pkg/math/calc.go":     "import math",
		"pkg/math/calc_test.go": "import testing",
	}

	phases := comp.Comprehend(context.Background(), files, contents)
	if len(phases.Decomposition) != 3 {
		t.Errorf("expected 3 decomposed units, got %d", len(phases.Decomposition))
	}
	if len(phases.Recognition) == 0 {
		t.Errorf("expected recognized patterns, got none")
	}
}

func TestIncrementalStrategySwitch(t *testing.T) {
	reg := NewDefaultRegistry()
	// Start with BigBang
	p := Profile{FileCount: 2, LoC: 200, HasTests: true, HasBuild: true}
	next, err := reg.NextStrategy(context.Background(), compiletime.StrategyBigBang, p)
	if err != nil {
		t.Fatalf("NextStrategy failed: %v", err)
	}
	if next != compiletime.StrategyIncremental {
		t.Errorf("expected next strategy Incremental, got %s", next)
	}

	// Test SwitchToNextStrategy on state
	state := &compiletime.State{
		AnalyzerOutput: compiletime.AnalyzerOutput{
			Research: compiletime.DocumentWrapper[compiletime.SourceProjectResearch]{
				Data: compiletime.SourceProjectResearch{
					MigrationStrategy: string(compiletime.StrategyBigBang),
				},
			},
		},
	}
	newStrat, switched := SwitchToNextStrategy(context.Background(), state)
	if !switched {
		t.Fatalf("expected strategy switch to succeed")
	}
	if newStrat != compiletime.StrategyFrozenLegacy {
		t.Errorf("expected FROZEN_LEGACY when HasTests is false, got %s", newStrat)
	}
}
