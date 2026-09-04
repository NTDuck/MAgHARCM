package agents_test

import (
	"context"
	"testing"

	"MAgHARCM/internal/agents"
	"MAgHARCM/internal/compiletime"
)

// TestTryInOrderEmptyProfile exercises the catch-all default path: a
// profile that matches none of the gated strategies (mid-sized, with tests)
// falls through to INCREMENTAL as the first match in the registry.
func TestTryInOrderEmptyProfile(t *testing.T) {
	reg := agents.NewDefaultRegistry()
	kind, err := reg.TryInOrder(context.Background(), agents.Profile{
		FileCount: 5,
		LoC:       800,
		HasTests:  true,
		HasBuild:  true,
	})
	if err != nil {
		t.Fatalf("TryInOrder returned error: %v", err)
	}
	if kind != compiletime.StrategyIncremental {
		t.Errorf("expected INCREMENTAL as the first match for default profile, got %s", kind)
	}
}

// TestTryInOrderPilotProfile exercises the large-codebase path: >50 files
// or >10k LoC routes to PILOT.
func TestTryInOrderPilotProfile(t *testing.T) {
	reg := agents.NewDefaultRegistry()
	kind, err := reg.TryInOrder(context.Background(), agents.Profile{
		FileCount: 60,
		LoC:       15000,
		HasTests:  true,
		HasBuild:  true,
	})
	if err != nil {
		t.Fatalf("TryInOrder returned error: %v", err)
	}
	if kind != compiletime.StrategyPilot {
		t.Errorf("expected PILOT for large repo, got %s", kind)
	}
}

// TestTryInOrderFrozenLegacyProfile exercises the no-test path: a profile
// with HasTests=false routes to FROZEN_LEGACY.
func TestTryInOrderFrozenLegacyProfile(t *testing.T) {
	reg := agents.NewDefaultRegistry()
	kind, err := reg.TryInOrder(context.Background(), agents.Profile{
		FileCount: 10,
		LoC:       2000,
		HasTests:  false,
		HasBuild:  true,
	})
	if err != nil {
		t.Fatalf("TryInOrder returned error: %v", err)
	}
	if kind != compiletime.StrategyFrozenLegacy {
		t.Errorf("expected FROZEN_LEGACY for HasTests=false repo, got %s", kind)
	}
}

// TestTryInOrderBigBangProfile exercises the small-codebase path: <=3
// files AND <500 LoC routes to BIG_BANG.
func TestTryInOrderBigBangProfile(t *testing.T) {
	reg := agents.NewDefaultRegistry()
	kind, err := reg.TryInOrder(context.Background(), agents.Profile{
		FileCount: 2,
		LoC:       200,
		HasTests:  true,
		HasBuild:  true,
	})
	if err != nil {
		t.Fatalf("TryInOrder returned error: %v", err)
	}
	if kind != compiletime.StrategyBigBang {
		t.Errorf("expected BIG_BANG for small repo, got %s", kind)
	}
}

// TestTryInOrderParallelCutoverProfile exercises the modular-with-tests
// path: >10 files AND HasTests routes to PARALLEL_CUTOVER.
func TestTryInOrderParallelCutoverProfile(t *testing.T) {
	reg := agents.NewDefaultRegistry()
	kind, err := reg.TryInOrder(context.Background(), agents.Profile{
		FileCount: 15,
		LoC:       3000,
		HasTests:  true,
		HasBuild:  true,
	})
	if err != nil {
		t.Fatalf("TryInOrder returned error: %v", err)
	}
	if kind != compiletime.StrategyParallelCutover {
		t.Errorf("expected PARALLEL_CUTOVER for modular repo with tests, got %s", kind)
	}
}

// TestSelectAndTryStrategies verifies the convenience runner returns the
// strategy kind and a non-empty rationale for a small-repo profile that
// matches BIG_BANG.
func TestSelectAndTryStrategies(t *testing.T) {
	kind, rationale, err := agents.SelectAndTryStrategies(context.Background(), agents.Profile{
		FileCount: 2,
		LoC:       200,
		HasTests:  true,
		HasBuild:  true,
	})
	if err != nil {
		t.Fatalf("SelectAndTryStrategies returned error: %v", err)
	}
	if kind != compiletime.StrategyBigBang {
		t.Errorf("expected BIG_BANG, got %s", kind)
	}
	if rationale == "" {
		t.Errorf("expected non-empty rationale for %s", kind)
	}
}
