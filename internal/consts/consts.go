// Package consts centralizes compile-time constants previously scattered as
// inline literals across agent and tool code. None of these act as config
// fallbacks — every config field must come from the external YAML file.
package consts

// ErrorUnknown is the compiler / toolchain-name sentinel used when no
// language registry entry matches (internal/tools/exec.go).
const ErrorUnknown = "unknown"

// ToolchainCargo is the canonical Cargo toolchain name.
const ToolchainCargo = "cargo"

// LangRust is the lower-case Rust language key used by the language registry.
const LangRust = "rust"

// LangRustCanonical is the display-canonical Rust language name.
const LangRustCanonical = "Rust"

// LSPProviderNative is the tree-sitter based LSP provider name.
const LSPProviderNative = "native"

// LSPProviderABCoder is the new canonical abcoder MCP LSP provider name.
const LSPProviderABCoder = "abcoder-mcp"

// The five Mueller migration strategies produced by
// agents.SelectMigrationStrategy (internal/agents/analyzer.go).
const (
	StrategyBigBang         = "BIG_BANG"
	StrategyIncremental     = "INCREMENTAL"
	StrategyPilot           = "PILOT"
	StrategyFrozenLegacy    = "FROZEN_LEGACY"
	StrategyParallelCutover = "PARALLEL_CUTOVER"
)
// SpecMinerToolName is the tool key used by the dynamic invariant recovery agent.
// Backlink: [[P04]] SpecMiner-Style Dynamic Invariant Recovery ([[Syzygy-2024]]).
const SpecMinerToolName = "specminer"

// DefaultArtifactDir is the default artifact/checkpoint directory.
const DefaultArtifactDir = ".artifacts"

// SourceReindexed marks a Lookup result that was served from the iterative
// navigator's already-emitted-fragment index (P30, RepoCoder-style).
const SourceReindexed = "reindexed"

// SourceFresh marks a Lookup result that required a fresh Navigator call
// because the symbol was not present in the emitted-fragment index.
const SourceFresh = "fresh"

// DefaultRequestFile is the default YAML request file name.
const DefaultRequestFile = "magharcm-request.yml"
