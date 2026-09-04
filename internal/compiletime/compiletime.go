package compiletime

import (
	"fmt"
)

// Backlink: [[1.0.0 Architecture]] and [[1.0.0 Methodology]].
// Package compiletime provides centralized compile-time configuration,
// constants, sentinels, enums, and initialization helpers using the Must pattern.
// These are strictly distinguished from runtime YAML configurations (package config).

// -------------------------------------------------------------------------
// Compile-time Sentinels & Constants
// -------------------------------------------------------------------------


// CurrentSchemaVersion is the canonical schema version string for role artifacts.
const CurrentSchemaVersion = "1.0.0"
// ErrorUnknown is the sentinel compiler/toolchain string when no language matches.
const ErrorUnknown = "unknown"

// ToolchainCargo is the canonical Cargo toolchain name.
const ToolchainCargo = "cargo"

// ToolchainGo is the canonical Go toolchain name.
const ToolchainGo = "go"

// LangRust is the lower-case Rust language key.
const LangRust = "rust"

// LangRustCanonical is the display-canonical Rust language name.
const LangRustCanonical = "Rust"

// LangGo is the Go language key.
const LangGo = "go"

// LangC is the C language key.
const LangC = "c"

// LangJava is the Java language key.
const LangJava = "java"

// LangPython is the Python language key.
const LangPython = "python"

// LSPProviderNative is the tree-sitter based LSP provider.
const LSPProviderNative = "native"

// LSPProviderABCoder is the canonical abcoder MCP LSP provider name (default).
const LSPProviderABCoder = "abcoder-mcp"

// SpecMinerToolName is the canonical PRIM-4 SpecMiner tool name.
const SpecMinerToolName = "specminer"

// DefaultArtifactDir is the canonical artifact/checkpoint directory.
const DefaultArtifactDir = ".artifacts"

// SourceReindexed marks a symbol served from the emitted-fragment index (PRIM-31).
const SourceReindexed = "reindexed"

// SourceFresh marks a symbol requiring a fresh Navigator round-trip.
const SourceFresh = "fresh"

// DefaultRequestFile is the canonical YAML request file name.
const DefaultRequestFile = "magharcm-request.yml"

// -------------------------------------------------------------------------
// Strict Enums
// -------------------------------------------------------------------------

// CompilationStatus represents binary per-project compilation outcomes.
// There is no partial compilation percentage: status is strictly Pass or Fail.
type CompilationStatus string

const (
	CompilationStatusPass CompilationStatus = "PASS"
	CompilationStatusFail CompilationStatus = "FAIL"
)

// StrategyKind names one of Mueller's five migration strategies (PRIM-21).
type StrategyKind string

const (
	StrategyBigBang         StrategyKind = "BIG_BANG"
	StrategyIncremental     StrategyKind = "INCREMENTAL"
	StrategyPilot           StrategyKind = "PILOT"
	StrategyFrozenLegacy    StrategyKind = "FROZEN_LEGACY"
	StrategyParallelCutover StrategyKind = "PARALLEL_CUTOVER"
)

// WorkUnitStatus defines P17 Blackboard scheduler lifecycle states (CAID).
type WorkUnitStatus string

const (
	WorkUnitPending   WorkUnitStatus = "PENDING"
	WorkUnitClaimed   WorkUnitStatus = "CLAIMED"
	WorkUnitCompleted WorkUnitStatus = "COMPLETED"
	WorkUnitFailed    WorkUnitStatus = "FAILED"
	WorkUnitTerminal  WorkUnitStatus = "TERMINAL"
)

// -------------------------------------------------------------------------
// Must Pattern Helpers (No Fallbacks)
// -------------------------------------------------------------------------

// Must panics if err is non-nil; otherwise it returns v.
// Used for compile-time and startup initializations where failure is fatal.
func Must[T any](v T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("compiletime initialization failed: %v", err))
	}
	return v
}

// MustNotNil panics if ptr is nil.
func MustNotNil[T any](ptr *T, fieldName string) *T {
	if ptr == nil {
		panic(fmt.Sprintf("compiletime invariant violated: %s must not be nil", fieldName))
	}
	return ptr
}

// MustNotEmpty panics if s is empty.
func MustNotEmpty(s string, fieldName string) string {
	if s == "" {
		panic(fmt.Sprintf("compiletime invariant violated: %s must not be empty", fieldName))
	}
	return s
}

// -------------------------------------------------------------------------
// Compile-time Task Specification
// -------------------------------------------------------------------------

// Task defines the specification for a translation task.
type Task struct {
	SourceDir   string `json:"source_dir"`
	TargetDir   string `json:"target_dir"`
	SourceLang  string `json:"source_lang"`
	TargetLang  string `json:"target_lang"`
	Toolchain   string `json:"toolchain,omitempty"`
	LSPProvider string `json:"lsp_provider,omitempty"`
}

// Validate verifies that all required fields of Task are populated.
func (t Task) Validate() error {
	if t.SourceDir == "" {
		return fmt.Errorf("task: source_dir is required")
	}
	if t.TargetDir == "" {
		return fmt.Errorf("task: target_dir is required")
	}
	if t.SourceLang == "" {
		return fmt.Errorf("task: source_lang is required")
	}
	if t.TargetLang == "" {
		return fmt.Errorf("task: target_lang is required")
	}
	if t.Toolchain == "" {
		return fmt.Errorf("task: toolchain is required")
	}
	if t.LSPProvider == "" {
		return fmt.Errorf("task: lsp_provider is required")
	}
	return nil
}

// MustTask returns the task or panics if invalid.
func MustTask(t Task) Task {
	Must(struct{}{}, t.Validate())
	return t
}
