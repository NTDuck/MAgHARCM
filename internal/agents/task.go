package agents

import (
	"fmt"

	"MAgHARCM/internal/compiletime"
)

// Backlink: [[1.0.0 ADR-2026-09-07-Sprint-Conventions]] rule ADR-C-014
// (Locality of Behaviour for producer artifact structs).
//
// Task is the user-supplied translation specification consumed by the
// pipeline entry-points (cmd/...) and the configuration loader. The
// canonical declaration lives here in the producer agent file per
// ADR-C-014; callers across the codebase keep accessing the type via
// the compiletime.Task re-export alias declared in
// internal/compiletime/state.go.

// Task defines the specification for a translation task. Every required
// field is populated by the configuration loader before pipeline execution.
type Task struct {
	SourceDir   string
	TargetDir   string
	SourceLang  string
	TargetLang  string
	Toolchain   string
	LSPProvider string
	RequestFile string
}

// Validate verifies that every required field of Task is populated.
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

// MustTask returns the task or panics if invalid. Use at startup where a
// malformed task is a fatal configuration error.
func MustTask(t Task) Task {
	compiletime.Must(struct{}{}, t.Validate())
	return t
}
