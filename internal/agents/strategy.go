package agents

import (
	"context"
	"fmt"

	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/logger"
)

// StrategyKind aliases compiletime.StrategyKind for centralized enums.
type StrategyKind = compiletime.StrategyKind

const (
	StrategyBigBang         = compiletime.StrategyBigBang
	StrategyIncremental     = compiletime.StrategyIncremental
	StrategyPilot           = compiletime.StrategyPilot
	StrategyFrozenLegacy    = compiletime.StrategyFrozenLegacy
	StrategyParallelCutover = compiletime.StrategyParallelCutover
)
// Profile is the static repository signal set used by each MigrationStrategy
// to decide whether its gating conditions apply. The fields are populated
// once at analyzer start and then handed to Registry.TryInOrder.
type Profile struct {
	FileCount int
	LoC       int
	HasTests  bool
	HasBuild  bool
}

// MigrationStrategy is one pluggable Mueller strategy. Kind() returns its
// stable identifier; Matches() returns true when the strategy's gating
// conditions apply to the given profile; Attempt() executes the strategy and
// returns nil on success or an error if the runtime gate short-circuits.
type MigrationStrategy interface {
	Kind() StrategyKind
	Matches(p Profile) bool
	Attempt(ctx context.Context, p Profile) error
}

// Registry holds the ordered list of strategies tried by TryInOrder.
type Registry struct {
	all []MigrationStrategy
}

// NewDefaultRegistry returns the five Mueller strategies in canonical
// priority order: BIG_BANG (tiny project), PILOT (very large), FROZEN_LEGACY
// (no tests), PARALLEL_CUTOVER (modular with tests), then INCREMENTAL as the
// catch-all default.
func NewDefaultRegistry() *Registry {
	return &Registry{
		all: []MigrationStrategy{
			bigBangStrategy{},
			pilotStrategy{},
			frozenLegacyStrategy{},
			parallelCutoverStrategy{},
			incrementalStrategy{},
		},
	}
}

// TryInOrder iterates the registry, returning the first strategy whose
// Matches() and Attempt() both succeed. A Miss (Matches true, Attempt err)
// is logged and the next strategy is tried. Exhaustion returns an error
// describing every tried strategy; this is unreachable with the default
// registry because INCREMENTAL is the universal default.
func (r *Registry) TryInOrder(ctx context.Context, p Profile) (StrategyKind, error) {
	for _, s := range r.all {
		if !s.Matches(p) {
			continue
		}
		if err := s.Attempt(ctx, p); err != nil {
			logger.LogAgent("Analyzer", "Strategy %s declined for profile (files=%d, loc=%d, tests=%v): %v",
				s.Kind(), p.FileCount, p.LoC, p.HasTests, err)
			continue
		}
		return s.Kind(), nil
	}
	return "", fmt.Errorf("no migration strategy matched profile (files=%d, loc=%d, tests=%v)",
		p.FileCount, p.LoC, p.HasTests)
}

// NextStrategy finds the next viable strategy in the registry after the current one failed.
func (r *Registry) NextStrategy(ctx context.Context, current StrategyKind, p Profile) (StrategyKind, error) {
	foundCurrent := false
	for _, s := range r.all {
		if !foundCurrent {
			if s.Kind() == current {
				foundCurrent = true
			}
			continue
		}
		if !s.Matches(p) {
			continue
		}
		if err := s.Attempt(ctx, p); err != nil {
			logger.LogAgent("StrategyRegistry", "Strategy %s declined during fallback: %v", s.Kind(), err)
			continue
		}
		return s.Kind(), nil
	}
	return "", fmt.Errorf("no subsequent fallback strategy available after %s", current)
}

// SwitchToNextStrategy advances state to the next viable strategy when current strategy fails or plateaus.
func SwitchToNextStrategy(ctx context.Context, state *State) (StrategyKind, bool) {
	if state == nil {
		return "", false
	}
	current := StrategyKind(state.AnalyzerOutput.Research.Data.MigrationStrategy)
	p := Profile{
		FileCount: len(state.PlanningOutput.Fragments),
		HasTests:  state.ValidationReport.TotalTests > 0,
		HasBuild:  state.ValidationReport.CompilationSuccess,
	}
	reg := NewDefaultRegistry()
	next, err := reg.NextStrategy(ctx, current, p)
	if err != nil {
		logger.LogWarning("No further migration strategies available after %s: %v", current, err)
		return current, false
	}
	logger.LogAgent("StrategyRegistry", "Incremental failover: switching strategy from %s to %s", current, next)
	state.AnalyzerOutput.Research.Data.MigrationStrategy = string(next)
	state.AnalyzerOutput.Research.Data.StrategyRationale = rationaleFor(next)
	return next, true
}

// ----- concrete strategies -----

// bigBangStrategy: tiny self-contained project — direct single-pass
// translation, no staged rollout.
type bigBangStrategy struct{}

func (bigBangStrategy) Kind() StrategyKind { return StrategyBigBang }
func (bigBangStrategy) Matches(p Profile) bool {
	return p.FileCount <= 3 && p.LoC < 500
}
func (bigBangStrategy) Attempt(ctx context.Context, p Profile) error {
	// Runtime gate: when wired, would abort translation if any module
	// exceeds the BIG_BANG budget mid-pass. For now the analyzer only
	// needs the strategy kind for prompt context.
	_ = ctx
	_ = p
	return nil
}

// pilotStrategy: very large codebase — chunked subsystem pilot before
// sweeping migration.
type pilotStrategy struct{}

func (pilotStrategy) Kind() StrategyKind { return StrategyPilot }
func (pilotStrategy) Matches(p Profile) bool {
	return p.FileCount > 50 || p.LoC > 10000
}
func (pilotStrategy) Attempt(ctx context.Context, p Profile) error {
	// Runtime gate: would refuse to start unless a pilot chunk is
	// selected and validated before sweeping migration.
	_ = ctx
	_ = p
	return nil
}

// frozenLegacyStrategy: no test harness — synthesise tests and freeze the
// legacy boundary before touching production code.
type frozenLegacyStrategy struct{}

func (frozenLegacyStrategy) Kind() StrategyKind { return StrategyFrozenLegacy }
func (frozenLegacyStrategy) Matches(p Profile) bool {
	return !p.HasTests
}
func (frozenLegacyStrategy) Attempt(ctx context.Context, p Profile) error {
	// Runtime gate: would require a synthetic-test pass to complete
	// before any source mutation begins.
	_ = ctx
	_ = p
	return nil
}

// parallelCutoverStrategy: modular project with comprehensive tests —
// multi-stage parallel module cutover with side-by-side validation.
type parallelCutoverStrategy struct{}

func (parallelCutoverStrategy) Kind() StrategyKind { return StrategyParallelCutover }
func (parallelCutoverStrategy) Matches(p Profile) bool {
	return p.HasTests && p.FileCount > 10
}
func (parallelCutoverStrategy) Attempt(ctx context.Context, p Profile) error {
	// Runtime gate: would require parallel harness bootstrapped and
	// per-module cutover checkpoints observed.
	_ = ctx
	_ = p
	return nil
}

// incrementalStrategy: default catch-all — standard reverse-topological
// incremental translation. Matches every profile, so it always wins as the
// last strategy in the registry.
type incrementalStrategy struct{}

func (incrementalStrategy) Kind() StrategyKind { return StrategyIncremental }
func (incrementalStrategy) Matches(p Profile) bool {
	_ = p
	return true
}
func (incrementalStrategy) Attempt(ctx context.Context, p Profile) error {
	// Runtime gate: would refuse to translate a module whose downstream
	// dependents have not yet been migrated.
	_ = ctx
	_ = p
	return nil
}

// rationaleFor maps a strategy kind to its short rationale text, used by
// the analyzer when populating state. Keeping the rationale outside the
// strategy itself preserves the strategy's focus on gating.
func rationaleFor(kind StrategyKind) string {
	switch kind {
	case StrategyBigBang:
		return "Small self-contained project (<500 LoC, <=3 files): single-pass direct translation."
	case StrategyPilot:
		return "Large-scale codebase (>50 files or >10k LoC): chunked subsystem pilot translation."
	case StrategyFrozenLegacy:
		return "Legacy codebase without test harness: requires test synthesis and boundary freezing."
	case StrategyParallelCutover:
		return "Modular project with comprehensive test suite: multi-stage parallel module cutover."
	case StrategyIncremental:
		return "Standard multi-module project: reverse-topological incremental translation."
	default:
		return "No matching strategy selected."
	}
}

// SelectAndTryStrategies is the convenience runner used by AnalyzerAgent.Run:
// it builds the default registry, runs it over the given profile, and
// returns the chosen kind plus its rationale. The error is returned for
// forward compatibility but the default registry never produces one.
func SelectAndTryStrategies(ctx context.Context, p Profile) (StrategyKind, string, error) {
	kind, err := NewDefaultRegistry().TryInOrder(ctx, p)
	if err != nil {
		return "", "", err
	}
	return kind, rationaleFor(kind), nil
}
