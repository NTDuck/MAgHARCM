package runner_test

import (
	"MAgHARCM/internal/runner"
	"context"
	"errors"
	"strings"
	"testing"

	"MAgHARCM/internal/config"
)

func TestRunRejectsNilConfig(t *testing.T) {
	if _, err := runner.Run(context.Background(), nil); !errors.Is(err, runner.ErrMissingFields) {
		t.Errorf("nil cfg: got %v want runner.ErrMissingFields", err)
	}
}

func TestRunRejectsEmptyDirs(t *testing.T) {
	cfg := &config.Config{}
	_, err := runner.Run(context.Background(), cfg)
	if err == nil {
		t.Fatalf("empty cfg: expected error, got nil")
	}
	if !strings.Contains(err.Error(), "config: missing required fields:") {
		t.Errorf("empty cfg: got %v want substring %q", err, "config: missing required fields:")
	}

	cfg.SourceDir = "src/c"
	_, err = runner.Run(context.Background(), cfg)
	if err == nil {
		t.Fatalf("only SourceDir set: expected error, got nil")
	}
	if !strings.Contains(err.Error(), "config: missing required fields:") {
		t.Errorf("only SourceDir set: got %v want substring %q", err, "config: missing required fields:")
	}
}

// TestSuccessNilState verifies runner.Success() does not panic on a nil final
// state — useful for callers that early-return before completion.
func TestSuccessNilState(t *testing.T) {
	if runner.Success(nil) {
		t.Errorf("runner.Success(nil) should be false")
	}
}
