package main

import (
	"testing"

	"MAgHARCM/internal/config"
)

func TestConfigLoadAndOverrides(t *testing.T) {
	yamlContent := `
translation:
  source:
    dir: "assets/samples/test_c"
    language: "C"
  target:
    dir: ".artifacts/test_rust"
    language: "Rust"
    toolchain: "cargo"
  models:
    reasoning: "gpt-oss:20b"
    coding: "qwen:coder"
    ollama_url: "http://localhost:11434"
  execution:
    max_iterations: 5
    timeout_seconds: 1800
  lsp:
    provider: "native"
`
	cfg, err := config.ParseYAML([]byte(yamlContent))
	if err != nil {
		t.Fatalf("failed to parse yaml config: %v", err)
	}

	if cfg.SourceDir != "assets/samples/test_c" {
		t.Errorf("expected source_dir assets/samples/test_c, got %s", cfg.SourceDir)
	}
	if cfg.TargetDir != ".artifacts/test_rust" {
		t.Errorf("expected target_dir .artifacts/test_rust, got %s", cfg.TargetDir)
	}
	if cfg.SourceLang != "C" {
		t.Errorf("expected source_lang C, got %s", cfg.SourceLang)
	}
	if cfg.TargetLang != "Rust" {
		t.Errorf("expected target_lang Rust, got %s", cfg.TargetLang)
	}
	if cfg.Toolchain != "cargo" {
		t.Errorf("expected toolchain cargo, got %s", cfg.Toolchain)
	}
	if cfg.MaxIterations != 5 {
		t.Errorf("expected max_iterations 5, got %d", cfg.MaxIterations)
	}
}
