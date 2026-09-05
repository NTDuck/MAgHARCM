package agents

import (
	"context"
	"fmt"

	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/logger"
)

// StrategyKind aliases compiletime.StrategyKind for centralized enums.
type StrategyKind = compiletime.StrategyKind

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
		logger.LogAgent(compiletime.LogScopeAnalyzer, "Strategy %s declined for profile: files=%d loc=%d tests=%v err=%v",
			s.Kind(), p.FileCount, p.LoC, p.HasTests, err)
			continue
		}
		return s.Kind(), nil
	}
	return "", fmt.Errorf("no migration strategy matched profile: files=%d loc=%d tests=%v",
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
			logger.LogAgent(compiletime.LogScopeStrategyRegistry, "Strategy %s declined during fallback: %v", s.Kind(), err)
		}
		return s.Kind(), nil
	}
	return "", fmt.Errorf("no subsequent fallback strategy available after %s", current)
}

// SwitchToNextStrategy advances state to the next viable strategy when current strategy fails or plateaus.
func SwitchToNextStrategy(ctx context.Context, state *compiletime.State) (StrategyKind, bool) {
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
	logger.LogAgent(compiletime.LogScopeStrategyRegistry, "Incremental failover: switching strategy from %s to %s", current, next)
	state.AnalyzerOutput.Research.Data.StrategyRationale = rationaleFor(next)
	return next, true
}

// ----- concrete strategies -----

// bigBangStrategy: tiny self-contained project — direct single-pass
// translation, no staged rollout.
type bigBangStrategy struct{}

func (bigBangStrategy) Kind() StrategyKind { return compiletime.StrategyBigBang }
func (bigBangStrategy) Matches(p Profile) bool {
	return p.FileCount <= compiletime.BigBangFileMax && p.LoC < compiletime.BigBangLoCMax
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

func (pilotStrategy) Kind() StrategyKind { return compiletime.StrategyPilot }
func (pilotStrategy) Matches(p Profile) bool {
	return p.FileCount > compiletime.PilotFileMin || p.LoC > compiletime.PilotLoCMin
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

func (frozenLegacyStrategy) Kind() StrategyKind { return compiletime.StrategyFrozenLegacy }
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

func (parallelCutoverStrategy) Kind() StrategyKind { return compiletime.StrategyParallelCutover }
func (parallelCutoverStrategy) Matches(p Profile) bool {
	return p.HasTests && p.FileCount > compiletime.ParallelCutoverFileMin
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

func (incrementalStrategy) Kind() StrategyKind { return compiletime.StrategyIncremental }
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

// strategyRationalates is the rationale lookup used by the analyzer when
// populating state. Keeping the table outside the strategy itself preserves
// each strategy's focus on gating.
var strategyRationalates = map[StrategyKind]string{
	compiletime.StrategyBigBang:         compiletime.StrategyRationaleBigBang,
	compiletime.StrategyPilot:           compiletime.StrategyRationalePilot,
	compiletime.StrategyFrozenLegacy:    compiletime.StrategyRationaleFrozenLegacy,
	compiletime.StrategyParallelCutover: compiletime.StrategyRationaleParallelCutover,
	compiletime.StrategyIncremental:     compiletime.StrategyRationaleIncremental,
}

func rationaleFor(kind StrategyKind) string {
	if r, ok := strategyRationalates[kind]; ok {
		return r
	}
	return compiletime.StrategyRationaleUnknown
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
