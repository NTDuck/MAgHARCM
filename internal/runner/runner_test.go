package runner

import (
	"context"
	"errors"
	"testing"

	"MAgHARCM/internal/config"
)

func TestRunRejectsNilConfig(t *testing.T) {
	if _, err := Run(context.Background(), nil); !errors.Is(err, ErrMissingFields) {
		t.Errorf("nil cfg: got %v want ErrMissingFields", err)
	}
}

func TestRunRejectsEmptyDirs(t *testing.T) {
	cfg := &config.Config{}
	if _, err := Run(context.Background(), cfg); !errors.Is(err, ErrMissingFields) {
		t.Errorf("empty dirs: got %v want ErrMissingFields", err)
	}
	cfg.SourceDir = "src/c"
	if _, err := Run(context.Background(), cfg); !errors.Is(err, ErrMissingFields) {
		t.Errorf("empty target_dir: got %v want ErrMissingFields", err)
	}
}

// TestSuccessNilState verifies Success() does not panic on a nil final
// state — useful for callers that early-return before completion.
func TestSuccessNilState(t *testing.T) {
	if Success(nil) {
		t.Errorf("Success(nil) should be false")
	}
}
