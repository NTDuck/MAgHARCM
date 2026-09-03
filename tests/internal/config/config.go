package config_test

import (
	"MAgHARCM/internal/config"
	"os"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	// Defaults reads from env; clear any inherited values for a stable test.
	for _, k := range []string{
		"OLLAMA_BASE_URL", "OLLAMA_REASONING_MODEL", "OLLAMA_CODING_MODEL",
		"MAGHARCM_MAX_ITERATIONS", "MAGHARCM_TIMEOUT_SECONDS",
	} {
		os.Unsetenv(k)
	}
	cfg := config.Defaults()

	if cfg.OllamaBaseURL != "http://localhost:11434" {
		t.Errorf("ollama URL: got %s", cfg.OllamaBaseURL)
	}
	if cfg.ReasoningModel != "qwen3:30b-a3b-thinking-2507-q4_K_M" {
		t.Errorf("reasoning model: got %s", cfg.ReasoningModel)
	}
	if cfg.CodingModel != "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL" {
		t.Errorf("coding model: got %s", cfg.CodingModel)
	}
	if cfg.MaxIterations != 20 {
		t.Errorf("max iterations: got %d", cfg.MaxIterations)
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
