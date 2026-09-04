// Package consts re-exports centralized compile-time constants from internal/compiletime.
// None of these act as config fallbacks — all runtime configs must come from external YAML.
package consts

import "MAgHARCM/internal/compiletime"

const (
	ErrorUnknown       = compiletime.ErrorUnknown
	ToolchainCargo     = compiletime.ToolchainCargo
	LangRust           = compiletime.LangRust
	LangRustCanonical  = compiletime.LangRustCanonical
	LSPProviderNative  = compiletime.LSPProviderNative
	LSPProviderABCoder = compiletime.LSPProviderABCoder
	SpecMinerToolName  = compiletime.SpecMinerToolName
	DefaultArtifactDir = compiletime.DefaultArtifactDir
	SourceReindexed    = compiletime.SourceReindexed
	SourceFresh        = compiletime.SourceFresh
	DefaultRequestFile = compiletime.DefaultRequestFile

	StrategyBigBang         = compiletime.StrategyBigBang
	StrategyIncremental     = compiletime.StrategyIncremental
	StrategyPilot           = compiletime.StrategyPilot
	StrategyFrozenLegacy    = compiletime.StrategyFrozenLegacy
	StrategyParallelCutover = compiletime.StrategyParallelCutover

	WorkUnitPending   = compiletime.WorkUnitPending
	WorkUnitClaimed   = compiletime.WorkUnitClaimed
	WorkUnitCompleted = compiletime.WorkUnitCompleted
	WorkUnitFailed    = compiletime.WorkUnitFailed
	WorkUnitTerminal  = compiletime.WorkUnitTerminal
)
