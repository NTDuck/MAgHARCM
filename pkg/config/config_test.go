package config

import (
	"testing"
	"time"
)

func TestConfigDefaults(t *testing.T) {
	cfg := MustLoad()

	if cfg.OllamaBaseURL != "http://localhost:11434" {
		t.Errorf("expected default ollama URL http://localhost:11434, got %s", cfg.OllamaBaseURL)
	}
	if cfg.ReasoningModel != "gpt-oss:20b" {
		t.Errorf("expected default reasoning model gpt-oss:20b, got %s", cfg.ReasoningModel)
	}
	if cfg.CodingModel != "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL" {
		t.Errorf("expected default coding model hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL, got %s", cfg.CodingModel)
	}
	if cfg.MaxIterations != 5 {
		t.Errorf("expected max iterations 5, got %d", cfg.MaxIterations)
	}
}

func TestConfigOptions(t *testing.T) {
	cfg := MustLoad(
		WithOllamaBaseURL("http://127.0.0.1:11434"),
		WithReasoningModel("custom-reasoning"),
		WithCodingModel("custom-coding"),
		WithMaxIterations(3),
		WithTimeout(10*time.Second),
		WithSourceTarget("src/a", "dst/b", "c", "rust"),
	)

	if cfg.OllamaBaseURL != "http://127.0.0.1:11434" {
		t.Errorf("expected 127.0.0.1, got %s", cfg.OllamaBaseURL)
	}
	if cfg.ReasoningModel != "custom-reasoning" {
		t.Errorf("expected custom-reasoning, got %s", cfg.ReasoningModel)
	}
	if cfg.CodingModel != "custom-coding" {
		t.Errorf("expected custom-coding, got %s", cfg.CodingModel)
	}
	if cfg.MaxIterations != 3 {
		t.Errorf("expected 3, got %d", cfg.MaxIterations)
	}
	if cfg.SourceDir != "src/a" || cfg.TargetDir != "dst/b" {
		t.Errorf("source or target mismatch: %s -> %s", cfg.SourceDir, cfg.TargetDir)
	}
}
