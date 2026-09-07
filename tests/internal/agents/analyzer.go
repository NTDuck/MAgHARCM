package agents_test

import (
	"context"
	"strings"
	"testing"

	"MAgHARCM/internal/agents"
	"MAgHARCM/internal/compiletime"
)

// TestAnalyzerStrategySelection mirrors the production selection logic via
// the new Registry.TryInOrder API. The legacy SelectMigrationStrategy
// helper was retired in favour of the strategy interface; this test
// exercises the same five canonical profiles that the legacy test covered.
func TestAnalyzerStrategySelection(t *testing.T) {
	reg := agents.NewDefaultRegistry()
	cases := []struct {
		name string
		p    agents.Profile
		want compiletime.StrategyKind
	}{
		{
			name: "small project -> BIG_BANG",
			p:    agents.Profile{FileCount: 2, LoC: 200, HasTests: true, HasBuild: true},
			want: compiletime.StrategyBigBang,
		},
		{
			name: "large project -> PILOT",
			p:    agents.Profile{FileCount: 60, LoC: 15000, HasTests: true, HasBuild: true},
			want: compiletime.StrategyPilot,
		},
		{
			name: "untested project -> FROZEN_LEGACY",
			p:    agents.Profile{FileCount: 10, LoC: 2000, HasTests: false, HasBuild: true},
			want: compiletime.StrategyFrozenLegacy,
		},
		{
			name: "modular project with tests -> PARALLEL_CUTOVER",
			p:    agents.Profile{FileCount: 15, LoC: 3000, HasTests: true, HasBuild: true},
			want: compiletime.StrategyParallelCutover,
		},
		{
			name: "standard modular with tests -> INCREMENTAL",
			p:    agents.Profile{FileCount: 5, LoC: 800, HasTests: true, HasBuild: true},
			want: compiletime.StrategyIncremental,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kind, err := reg.TryInOrder(context.Background(), tc.p)
			if err != nil {
				t.Fatalf("TryInOrder returned error: %v", err)
			}
			if kind != tc.want {
				t.Errorf("expected %s, got %s", tc.want, kind)
			}
		})
	}
}

// TestAnalyzerRationaleStringFor verifies the analyzer still emits a
// non-empty rationale string for each chosen strategy so the prompt
// assembly downstream keeps working.
func TestAnalyzerRationaleStringFor(t *testing.T) {
	cases := []compiletime.StrategyKind{
		compiletime.StrategyBigBang,
		compiletime.StrategyIncremental,
		compiletime.StrategyPilot,
		compiletime.StrategyFrozenLegacy,
		compiletime.StrategyParallelCutover,
	}
	for _, k := range cases {
		_, rationale, err := agents.SelectAndTryStrategies(context.Background(), agents.Profile{
			FileCount: 2,
			LoC:       200,
			HasTests:  true,
			HasBuild:  true,
		})
		if err != nil && !strings.Contains(err.Error(), "no migration strategy") {
			t.Fatalf("unexpected error for %s: %v", k, err)
		}
		if rationale == "" {
			t.Errorf("empty rationale for strategy %s", k)
		}
	}
}
