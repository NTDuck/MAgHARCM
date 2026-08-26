package main

import (
	"testing"
)

func TestParsePromptRequirements(t *testing.T) {
	prompt := "Translate @assets/samples/GildedRose-Refactoring-Kata/C into a Rust codebase, output to .artifacts/GildedRose-Refactoring-Kata/rust"
	task := parsePromptRequirements(prompt)

	if task.SourceDir != "assets/samples/GildedRose-Refactoring-Kata/C" {
		t.Errorf("expected assets/samples/GildedRose-Refactoring-Kata/C, got %s", task.SourceDir)
	}
	if task.TargetDir != ".artifacts/GildedRose-Refactoring-Kata/rust" {
		t.Errorf("expected .artifacts/GildedRose-Refactoring-Kata/rust, got %s", task.TargetDir)
	}
	if task.SourceLang != "C" {
		t.Errorf("expected C, got %s", task.SourceLang)
	}
	if task.TargetLang != "Rust" {
		t.Errorf("expected Rust, got %s", task.TargetLang)
	}
}

func TestParseDefaultPrompt(t *testing.T) {
	task := parsePromptRequirements(defaultRequirementsPrompt)

	if task.SourceDir != "assets/samples/GildedRose-Refactoring-Kata/C" {
		t.Errorf("expected assets/samples/GildedRose-Refactoring-Kata/C, got %s", task.SourceDir)
	}
	if task.TargetDir != ".artifacts/GildedRose-Refactoring-Kata/rust" {
		t.Errorf("expected .artifacts/GildedRose-Refactoring-Kata/rust, got %s", task.TargetDir)
	}
	if task.SourceLang != "C" {
		t.Errorf("expected C, got %s", task.SourceLang)
	}
	if task.TargetLang != "Rust" {
		t.Errorf("expected Rust, got %s", task.TargetLang)
	}
}
