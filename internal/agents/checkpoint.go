package agents

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"MAgHARCM/internal/types"
)

// CheckpointDir returns the per-run checkpoint directory under .artifacts/<run-id>/checkpoints/.
func CheckpointDir(runID string) string {
	return filepath.Join(".artifacts", runID, "checkpoints")
}

// RunIDForTask derives a stable, human-readable run identifier from the
// translation task's source directory. The run ID is the cleaned source-dir
// path with path separators replaced by '-', so re-running against the same
// source directory always lands on the same checkpoint directory under
// .artifacts/. Two source dirs that canonicalize to the same path share a
// run ID.
func RunIDForTask(task types.TranslationTask) string {
	src := filepath.Clean(task.SourceDir)
	if src == "" || src == "." {
		return "default"
	}
	return strings.NewReplacer(string(filepath.Separator), "-").Replace(src)
}

// Checkpoint is a snapshot of *types.State plus a version + timestamp.
// Fields are versioned via Version so future schema changes don't break old checkpoints.
type Checkpoint struct {
	Version   int         `json:"version"`
	CreatedAt time.Time   `json:"created_at"`
	Iteration int         `json:"iteration"`
	State     *types.State `json:"state"`
}

const currentCheckpointVersion = 1

// Save writes a checkpoint for the given state under .artifacts/<runID>/checkpoints/iter-N.json.
// If runID is empty, returns an error (caller must always provide a run ID — generate one if missing).
// Returns the path written so the caller can log it.
func Save(runID string, state *types.State) (string, error) {
	if runID == "" {
		return "", fmt.Errorf("checkpoint: runID must not be empty")
	}
	dir := CheckpointDir(runID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("checkpoint: mkdir %s: %w", dir, err)
	}
	ckpt := Checkpoint{
		Version:   currentCheckpointVersion,
		CreatedAt: time.Now().UTC(),
		Iteration: state.Iteration,
		State:     state,
	}
	name := fmt.Sprintf("iter-%04d.json", state.Iteration)
	path := filepath.Join(dir, name)
	data, err := json.MarshalIndent(ckpt, "", "  ")
	if err != nil {
		return "", fmt.Errorf("checkpoint: marshal: %w", err)
	}
	// Write to temp file then rename for atomicity
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return "", fmt.Errorf("checkpoint: write %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return "", fmt.Errorf("checkpoint: rename %s -> %s: %w", tmp, path, err)
	}
	return path, nil
}

// LoadLatest returns the most recent checkpoint for runID, or (nil, nil) if no checkpoint exists.
// Errors are returned only for I/O or JSON-decode failures.
func LoadLatest(runID string) (*Checkpoint, error) {
	if runID == "" {
		return nil, fmt.Errorf("checkpoint: runID must not be empty")
	}
	dir := CheckpointDir(runID)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("checkpoint: readdir %s: %w", dir, err)
	}
	var latestPath string
	var latestIter int = -1
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		var iter int
		if _, err := fmt.Sscanf(e.Name(), "iter-%04d.json", &iter); err != nil {
			continue
		}
		if iter > latestIter {
			latestIter = iter
			latestPath = filepath.Join(dir, e.Name())
		}
	}
	if latestPath == "" {
		return nil, nil
	}
	data, err := os.ReadFile(latestPath)
	if err != nil {
		return nil, fmt.Errorf("checkpoint: read %s: %w", latestPath, err)
	}
	var ckpt Checkpoint
	if err := json.Unmarshal(data, &ckpt); err != nil {
		return nil, fmt.Errorf("checkpoint: unmarshal %s: %w", latestPath, err)
	}
	if ckpt.Version != currentCheckpointVersion {
		return nil, fmt.Errorf("checkpoint: version mismatch in %s: got %d, want %d",
			latestPath, ckpt.Version, currentCheckpointVersion)
	}
	return &ckpt, nil
}

// Cleanup removes all checkpoints for runID (typically called on successful completion).
// Returns nil if the directory doesn't exist.
func Cleanup(runID string) error {
	dir := CheckpointDir(runID)
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("checkpoint: cleanup %s: %w", dir, err)
	}
	return nil
}
