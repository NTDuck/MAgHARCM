package config

import (
	"testing"
	"time"
)

func TestParseNestedYAML(t *testing.T) {
	yamlData := `
translation:
  source:
    dir: "assets/samples/my_project/C"
    language: "C"
  target:
    dir: ".artifacts/my_project/rust"
    language: "Rust"
    toolchain: "cargo"
  models:
    reasoning: "custom-reasoning:latest"
    coding: "custom-coding:latest"
    ollama_url: "http://localhost:11434"
  execution:
    max_iterations: 15
    timeout_seconds: 300
  lsp:
    provider: "abcoder"
`

	cfg, err := ParseYAML([]byte(yamlData))
	if err != nil {
		t.Fatalf("failed to parse yaml: %v", err)
	}

	if cfg.SourceDir != "assets/samples/my_project/C" {
		t.Errorf("expected source dir assets/samples/my_project/C, got %s", cfg.SourceDir)
	}
	if cfg.TargetDir != ".artifacts/my_project/rust" {
		t.Errorf("expected target dir .artifacts/my_project/rust, got %s", cfg.TargetDir)
	}
	if cfg.SourceLang != "C" {
		t.Errorf("expected source lang C, got %s", cfg.SourceLang)
	}
	if cfg.TargetLang != "Rust" {
		t.Errorf("expected target lang Rust, got %s", cfg.TargetLang)
	}
	if cfg.Toolchain != "cargo" {
		t.Errorf("expected toolchain cargo, got %s", cfg.Toolchain)
	}
	if cfg.ReasoningModel != "custom-reasoning:latest" {
		t.Errorf("expected reasoning model custom-reasoning:latest, got %s", cfg.ReasoningModel)
	}
	if cfg.CodingModel != "custom-coding:latest" {
		t.Errorf("expected coding model custom-coding:latest, got %s", cfg.CodingModel)
	}
	if cfg.MaxIterations != 15 {
		t.Errorf("expected max iterations 15, got %d", cfg.MaxIterations)
	}
	if cfg.Timeout != 300*time.Second {
		t.Errorf("expected timeout 300s, got %v", cfg.Timeout)
	}
	if cfg.LSPProvider != "abcoder" {
		t.Errorf("expected lsp provider abcoder, got %s", cfg.LSPProvider)
	}
}

func TestParseFlatYAML(t *testing.T) {
	yamlData := `
source_dir: "src/c"
target_dir: "dst/rust"
source_lang: "C"
target_lang: "Rust"
toolchain: "cargo"
max_iterations: 5
`

	cfg, err := ParseYAML([]byte(yamlData))
	if err != nil {
		t.Fatalf("failed to parse flat yaml: %v", err)
	}

	if cfg.SourceDir != "src/c" || cfg.TargetDir != "dst/rust" {
		t.Errorf("unexpected dirs: %s -> %s", cfg.SourceDir, cfg.TargetDir)
	}
	if cfg.MaxIterations != 5 {
		t.Errorf("expected 5 iterations, got %d", cfg.MaxIterations)
	}
}
