package agents_test

import (
	"MAgHARCM/internal/checkpoint"
	"MAgHARCM/internal/compiletime"
	"os"
	"path/filepath"
	"strings"
	"testing"

)

// TestCheckpointSaveLoadRoundTrip verifies that a saved checkpoint can be loaded
// back with identical iteration and state contents.
func TestCheckpointSaveLoadRoundTrip(t *testing.T) {
	runID := "test-roundtrip"
	t.Cleanup(func() { _ = checkpoint.Cleanup(runID) })

	state := &compiletime.State{
		Iteration:     3,
		MaxIterations: 50,
		IsComplete:    false,
		TranslatedProject: compiletime.TranslatedProject{
			Files: map[string]string{"src/lib.rs": "fn hi() {}"},
		},
	}

	path, err := checkpoint.Save(runID, state)
	if err != nil {
		t.Fatalf("checkpoint.Save: %v", err)
	}
	if path == "" {
		t.Fatalf("checkpoint.Save returned empty path")
	}
	if !strings.Contains(path, "iter-0003.json") {
		t.Errorf("expected iter-0003.json in path, got %q", path)
	}

	ckpt, err := checkpoint.LoadLatest(runID)
	if err != nil {
		t.Fatalf("checkpoint.LoadLatest: %v", err)
	}
	if ckpt == nil {
		t.Fatal("checkpoint.LoadLatest returned nil checkpoint")
	}
	if ckpt.Version != checkpoint.CurrentCheckpointVersion {
		t.Errorf("Version mismatch: got %d, want %d", ckpt.Version, checkpoint.CurrentCheckpointVersion)
	}
	if ckpt.Iteration != 3 {
		t.Errorf("Iteration mismatch: got %d, want 3", ckpt.Iteration)
	}
	if ckpt.State == nil {
		t.Fatal("checkpoint.State is nil")
	}
	if ckpt.State.Iteration != 3 {
		t.Errorf("State.Iteration mismatch: got %d, want 3", ckpt.State.Iteration)
	}
	if ckpt.State.TranslatedProject.Files["src/lib.rs"] != "fn hi() {}" {
		t.Errorf("State.TranslatedProject content mismatch: got %q", ckpt.State.TranslatedProject.Files["src/lib.rs"])
	}
}

// TestCheckpointLoadLatestEmpty verifies that checkpoint.LoadLatest returns (nil, nil) when
// no checkpoint exists for the runID.
func TestCheckpointLoadLatestEmpty(t *testing.T) {
	runID := "test-empty"
	t.Cleanup(func() { _ = checkpoint.Cleanup(runID) })

	ckpt, err := checkpoint.LoadLatest(runID)
	if err != nil {
		t.Fatalf("checkpoint.LoadLatest: expected nil error, got %v", err)
	}
	if ckpt != nil {
		t.Errorf("checkpoint.LoadLatest: expected nil checkpoint, got %+v", ckpt)
	}
}

// TestCheckpointLoadLatestPicksHighestIteration verifies that checkpoint.LoadLatest returns
// the checkpoint with the highest iteration number when multiple are present.
func TestCheckpointLoadLatestPicksHighestIteration(t *testing.T) {
	runID := "test-multiple"
	t.Cleanup(func() { _ = checkpoint.Cleanup(runID) })

	// checkpoint.Save in non-monotonic order to make sure sorting is on disk content, not save order.
	for _, iter := range []int{1, 7, 3} {
		if _, err := checkpoint.Save(runID, &compiletime.State{Iteration: iter}); err != nil {
			t.Fatalf("checkpoint.Save iter=%d: %v", iter, err)
		}
	}

	ckpt, err := checkpoint.LoadLatest(runID)
	if err != nil {
		t.Fatalf("checkpoint.LoadLatest: %v", err)
	}
	if ckpt == nil {
		t.Fatal("checkpoint.LoadLatest returned nil checkpoint")
	}
	if ckpt.Iteration != 7 {
		t.Errorf("expected highest iteration 7, got %d", ckpt.Iteration)
	}
}

// TestCheckpointSaveEmptyRunIDReturnsError verifies that checkpoint.Save rejects an empty runID.
func TestCheckpointSaveEmptyRunIDReturnsError(t *testing.T) {
	if _, err := checkpoint.Save("", &compiletime.State{Iteration: 1}); err == nil {
		t.Fatal("checkpoint.Save with empty runID: expected error, got nil")
	}
}

// TestCheckpointCleanupRemovesDirectory verifies that checkpoint.Cleanup removes the
// checkpoint directory entirely.
func TestCheckpointCleanupRemovesDirectory(t *testing.T) {
	runID := "test-cleanup"

	if _, err := checkpoint.Save(runID, &compiletime.State{Iteration: 2}); err != nil {
		t.Fatalf("checkpoint.Save: %v", err)
	}
	dir := checkpoint.Dir(runID)
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("expected checkpoint dir to exist after checkpoint.Save: %v", err)
	}

	if err := checkpoint.Cleanup(runID); err != nil {
		t.Fatalf("checkpoint.Cleanup: %v", err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("expected checkpoint dir gone after checkpoint.Cleanup, got err=%v", err)
	}
	// checkpoint.Cleanup on a missing directory must be a no-op (no error).
	if err := checkpoint.Cleanup(runID); err != nil {
		t.Errorf("checkpoint.Cleanup on missing dir: expected nil error, got %v", err)
	}
}

// TestCheckpointVersionMismatchReturnsError verifies that a checkpoint with an
// unsupported version is rejected by checkpoint.LoadLatest.
func TestCheckpointVersionMismatchReturnsError(t *testing.T) {
	runID := "test-version"
	t.Cleanup(func() { _ = checkpoint.Cleanup(runID) })

	dir := checkpoint.Dir(runID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	bad := `{"version":999,"created_at":"2026-01-01T00:00:00Z","iteration":1,"state":{"iteration":1}}`
	path := filepath.Join(dir, "iter-0001.json")
	if err := os.WriteFile(path, []byte(bad), 0o644); err != nil {
		t.Fatalf("write bad checkpoint: %v", err)
	}

	_, err := checkpoint.LoadLatest(runID)
	if err == nil {
		t.Fatal("expected version-mismatch error, got nil")
	}
	if !strings.Contains(err.Error(), "version mismatch") {
		t.Errorf("expected version-mismatch error, got: %v", err)
	}
}
