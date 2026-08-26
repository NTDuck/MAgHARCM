package main

import (
	"testing"
)

func TestParseStructuredPrompt(t *testing.T) {
	prompt := `Input codebase: @assets/samples/sample_project/C
Output directory: .artifacts/sample_project/rust
Input language: C
Output language: Rust`

	task := parsePromptRequirements(prompt)

	if task.SourceDir != "assets/samples/sample_project/C" {
		t.Errorf("expected assets/samples/sample_project/C, got %s", task.SourceDir)
	}
	if task.TargetDir != ".artifacts/sample_project/rust" {
		t.Errorf("expected .artifacts/sample_project/rust, got %s", task.TargetDir)
	}
	if task.SourceLang != "C" {
		t.Errorf("expected C, got %s", task.SourceLang)
	}
	if task.TargetLang != "Rust" {
		t.Errorf("expected Rust, got %s", task.TargetLang)
	}
}

func TestParseNaturalLanguagePrompt(t *testing.T) {
	prompt := "Translate @sample_repo/src into Rust at target/out from C"
	task := parsePromptRequirements(prompt)

	if task.SourceDir != "sample_repo/src" {
		t.Errorf("expected sample_repo/src, got %s", task.SourceDir)
	}
	if task.TargetDir != "target/out" {
		t.Errorf("expected target/out, got %s", task.TargetDir)
	}
	if task.SourceLang != "C" {
		t.Errorf("expected C, got %s", task.SourceLang)
	}
	if task.TargetLang != "Rust" {
		t.Errorf("expected Rust, got %s", task.TargetLang)
	}
}
