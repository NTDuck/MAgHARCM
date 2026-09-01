package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"MAgHARCM/internal/config"
)

func TestSetField(t *testing.T) {
	cfg := &config.Config{}
	cases := []struct {
		key, val, want string
	}{
		{"source_dir", "src/foo", "src/foo"},
		{"source-language", "Go", "Go"},
		{"target_language", "Rust", "Rust"},
		{"toolchain", "cargo", "cargo"},
		{"reasoning", "qwen:7b", "qwen:7b"},
		{"coding", "qwen3:4b", "qwen3:4b"},
		{"ollama", "http://127.0.0.1:11434", "http://127.0.0.1:11434"},
		{"lsp", "native", "native"},
	}
	for _, c := range cases {
		if err := setField(cfg, c.key, c.val); err != nil {
			t.Errorf("setField(%s): %v", c.key, err)
			continue
		}
	}
	if cfg.SourceDir != "src/foo" {
		t.Errorf("source_dir: %q", cfg.SourceDir)
	}
	if cfg.SourceLang != "Go" {
		t.Errorf("source_language: %q", cfg.SourceLang)
	}
	if cfg.TargetLang != "Rust" {
		t.Errorf("target_language: %q", cfg.TargetLang)
	}
	if cfg.ReasoningModel != "qwen:7b" {
		t.Errorf("reasoning: %q", cfg.ReasoningModel)
	}
	if cfg.LSPProvider != "native" {
		t.Errorf("lsp: %q", cfg.LSPProvider)
	}
}

func TestSetFieldInt(t *testing.T) {
	cfg := &config.Config{}
	if err := setField(cfg, "max_iterations", "5"); err != nil {
		t.Fatalf("max_iterations: %v", err)
	}
	if cfg.MaxIterations != 5 {
		t.Errorf("max_iterations: %d", cfg.MaxIterations)
	}
	if err := setField(cfg, "timeout_seconds", "1800"); err != nil {
		t.Fatalf("timeout_seconds: %v", err)
	}
	if cfg.Timeout != 1800*time.Second {
		t.Errorf("timeout: %v", cfg.Timeout)
	}
	// negative rejected
	if err := setField(cfg, "max_iterations", "-1"); err == nil {
		t.Errorf("expected error for negative max_iterations")
	}
	if err := setField(cfg, "max_iterations", "abc"); err == nil {
		t.Errorf("expected error for non-int max_iterations")
	}
}

func TestSetFieldUnknown(t *testing.T) {
	if err := setField(&config.Config{}, "nope", "x"); err == nil {
		t.Errorf("expected error for unknown field")
	}
}

func TestPhase1FinalizeMissingFields(t *testing.T) {
	if err := phase1Finalize(&config.Config{}); err == nil {
		t.Errorf("expected error for empty source_dir/target_dir")
	}
	if err := phase1Finalize(&config.Config{SourceDir: "a", TargetDir: "b"}); err == nil {
		t.Errorf("expected error for empty source_language/target_language")
	}
}

func TestPhase1FinalizeWritesYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "req.yml")
	t.Setenv("MAGHARCM_REQUEST_PATH", path)

	cfg := &config.Config{
		SourceDir:      "src/foo",
		SourceLang:     "Go",
		TargetDir:      "out/foo",
		TargetLang:     "Rust",
		Toolchain:      "cargo",
		ReasoningModel: "qwen:7b",
		CodingModel:    "qwen3:4b",
		OllamaBaseURL:  "http://localhost:11434",
		MaxIterations:  3,
		Timeout:        1800 * time.Second,
		LSPProvider:    "native",
	}
	if err := phase1Finalize(cfg); err != nil {
		t.Fatalf("phase1Finalize: %v", err)
	}

	loaded, err := config.LoadYAML(path)
	if err != nil {
		t.Fatalf("load back: %v", err)
	}
	if loaded.SourceDir != "src/foo" || loaded.SourceLang != "Go" {
		t.Errorf("source not round-tripped: %s/%s", loaded.SourceDir, loaded.SourceLang)
	}
	if loaded.TargetDir != "out/foo" || loaded.TargetLang != "Rust" || loaded.Toolchain != "cargo" {
		t.Errorf("target not round-tripped: %s/%s/%s", loaded.TargetDir, loaded.TargetLang, loaded.Toolchain)
	}
	if loaded.MaxIterations != 3 {
		t.Errorf("max_iterations: %d", loaded.MaxIterations)
	}
	if loaded.Timeout != 1800*time.Second {
		t.Errorf("timeout: %v", loaded.Timeout)
	}
}

func TestHandleSlashShow(t *testing.T) {
	cfg := &config.Config{SourceDir: "x", SourceLang: "Go", TargetDir: "y", TargetLang: "Rust"}
	out := captureStdout(t, func() {
		if _, _, err := handleSlash("/show", cfg, phaseCollect); err != nil {
			t.Errorf("show: %v", err)
		}
	})
	if !strings.Contains(out, "source_dir        x") {
		t.Errorf("show output missing field: %q", out)
	}
}

func TestHandleSlashClearResets(t *testing.T) {
	def := config.Defaults()
	cfg := &config.Config{SourceDir: "x", SourceLang: "Go"}
	captureStdout(t, func() {
		next, _, err := handleSlash("/clear", cfg, phaseExecute)
		if err != nil {
			t.Errorf("clear: %v", err)
		}
		if next != phaseCollect {
			t.Errorf("clear should return phase 1, got %v", next)
		}
	})
	if cfg.SourceDir != "" {
		t.Errorf("clear should reset config, source_dir still %q", cfg.SourceDir)
	}
	if cfg.OllamaBaseURL != def.OllamaBaseURL {
		t.Errorf("clear should restore defaults, got ollama=%q", cfg.OllamaBaseURL)
	}
}

func TestHandleSlashLoadMalformed(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.yml")
	if err := os.WriteFile(bad, []byte("this: : : not yaml"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	cfg := &config.Config{}
	if _, _, err := handleSlash("/load "+bad, cfg, phaseCollect); err == nil {
		t.Errorf("expected error loading malformed YAML")
	}

	// valid YAML but missing required fields -> still rejected
	missing := filepath.Join(dir, "missing.yml")
	if err := os.WriteFile(missing, []byte("translation:\n  source:\n    dir: a\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, _, err := handleSlash("/load "+missing, cfg, phaseCollect); err == nil {
		t.Errorf("expected error: YAML missing required fields must be refused")
	}
}

func TestHandleSlashLoadValidJumpsToPhase2(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.yml")
	body := `translation:
  source:
    dir: a
    language: Go
  target:
    dir: b
    language: Rust
`
	if err := os.WriteFile(good, []byte(body), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	cfg := &config.Config{}
	next, _, err := handleSlash("/load "+good, cfg, phaseCollect)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if next != phaseExecute {
		t.Errorf("expected phase 2, got %v", next)
	}
	if cfg.SourceDir != "a" || cfg.SourceLang != "Go" {
		t.Errorf("cfg not loaded: %s/%s", cfg.SourceDir, cfg.SourceLang)
	}
}

func TestHandleSlashUnknownCommand(t *testing.T) {
	cfg := &config.Config{}
	if _, _, err := handleSlash("/nope", cfg, phaseCollect); err == nil {
		t.Errorf("expected error for unknown slash command")
	}
}

func TestHandleSlashQuit(t *testing.T) {
	cfg := &config.Config{}
	_, cont, err := handleSlash("/quit", cfg, phaseCollect)
	if err != nil {
		t.Fatalf("quit: %v", err)
	}
	if cont {
		t.Errorf("quit should set cont=false")
	}
}

// captureStdout runs fn and returns anything it printed.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	done := make(chan string, 1)
	go func() {
		var b strings.Builder
		buf := make([]byte, 4096)
		for {
			n, err := r.Read(buf)
			if n > 0 {
				b.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		done <- b.String()
	}()

	fn()
	w.Close()
	return <-done
}
