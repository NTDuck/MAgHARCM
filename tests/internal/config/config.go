package config_test

import (
	"MAgHARCM/internal/config"
	"os"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	// Defaults returns a zero-valued *Config for tests only. Unset any
	// environment values that older releases used to honor so this test
	// stays deterministic across environments.
	for _, k := range []string{
		"OLLAMA_BASE_URL", "OLLAMA_REASONING_MODEL", "OLLAMA_CODING_MODEL",
		"MAGHARCM_MAX_ITERATIONS", "MAGHARCM_TIMEOUT_SECONDS",
	} {
		os.Unsetenv(k)
	}
	cfg := config.Zero()

	if cfg.OllamaBaseURL != "" {
		t.Errorf("ollama URL: got %q want empty", cfg.OllamaBaseURL)
	}
	if cfg.ReasoningModel != "" {
		t.Errorf("reasoning model: got %q want empty", cfg.ReasoningModel)
	}
	if cfg.CodingModel != "" {
		t.Errorf("coding model: got %q want empty", cfg.CodingModel)
	}
	if cfg.MaxIterations != 0 {
		t.Errorf("max iterations: got %d want 0", cfg.MaxIterations)
	}
}

func TestParseYAMLOverrides(t *testing.T) {
	yamlData := `
translation:
  source:
    dir: "src/c"
    language: "C"
  target:
    dir: "dst/rust"
    language: "Rust"
    toolchain: "cargo"
  models:
    reasoning: "custom-reasoning:latest"
    coding: "custom-coding:latest"
    ollama_url: "http://127.0.0.1:11434"
  execution:
    max_iterations: 15
    timeout_seconds: 300
  lsp:
    provider: "native"
`
	cfg, err := config.ParseYAML([]byte(yamlData))
	if err != nil {
		t.Fatalf("ParseYAML: %v", err)
	}
	if cfg.SourceDir != "src/c" || cfg.SourceLang != "C" {
		t.Errorf("source: %s/%s", cfg.SourceDir, cfg.SourceLang)
	}
	if cfg.TargetDir != "dst/rust" || cfg.TargetLang != "Rust" || cfg.Toolchain != "cargo" {
		t.Errorf("target: %s/%s (%s)", cfg.TargetDir, cfg.TargetLang, cfg.Toolchain)
	}
	if cfg.ReasoningModel != "custom-reasoning:latest" {
		t.Errorf("reasoning: %s", cfg.ReasoningModel)
	}
	if cfg.CodingModel != "custom-coding:latest" {
		t.Errorf("coding: %s", cfg.CodingModel)
	}
	if cfg.OllamaBaseURL != "http://127.0.0.1:11434" {
		t.Errorf("ollama URL: %s", cfg.OllamaBaseURL)
	}
	if cfg.MaxIterations != 15 {
		t.Errorf("max iterations: %d", cfg.MaxIterations)
	}
	if cfg.LSPProvider != "native" {
		t.Errorf("lsp provider: %s", cfg.LSPProvider)
	}
}
