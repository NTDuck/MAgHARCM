package runner_test

import (
	"MAgHARCM/internal/runner"
	"context"
	"errors"
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
	if _, err := runner.Run(context.Background(), cfg); !errors.Is(err, runner.ErrMissingFields) {
		t.Errorf("empty dirs: got %v want runner.ErrMissingFields", err)
	}
	cfg.SourceDir = "src/c"
	if _, err := runner.Run(context.Background(), cfg); !errors.Is(err, runner.ErrMissingFields) {
		t.Errorf("empty target_dir: got %v want runner.ErrMissingFields", err)
	}
}

// TestSuccessNilState verifies runner.Success() does not panic on a nil final
// state — useful for callers that early-return before completion.
func TestSuccessNilState(t *testing.T) {
	if runner.Success(nil) {
		t.Errorf("runner.Success(nil) should be false")
	}
}
