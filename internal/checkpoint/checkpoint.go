// Package checkpoint persists snapshots of *compiletime.State to disk
// after every pipeline stage. It is a leaf package that depends only on
// internal/compiletime — no agent-internal imports — so the wider
// dependency graph (compiletime ↔ agents) stays acyclic.
//
// Backlink: [[1.0.0 PRIM-28]] Conversable State Checkpoints & Interrupts.
package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"MAgHARCM/internal/compiletime"
)

// pathSepReplacer replaces path separators with '-' when deriving a Run
// ID from a source directory. Hoisted to package-level for clarity.
var pathSepReplacer = strings.NewReplacer(string(filepath.Separator), "-")

// iterRe matches the per-iteration checkpoint filename pattern. Compiled
// once at package init.
var iterRe = regexp.MustCompile(`iter-(\d{4})` + regexp.QuoteMeta(compiletime.CheckpointExt))

// Dir returns the per-run checkpoint directory under
// .artifacts/<run-id>/checkpoints/.
func Dir(runID string) string {
	return filepath.Join(compiletime.DefaultArtifactDir, runID, "checkpoints")
}

// RunIDForSourceDir derives a stable, human-readable run identifier from
// a source directory path. The run ID is the cleaned source-dir path
// with path separators replaced by '-', so re-running against the same
// source directory always lands on the same checkpoint directory.
func RunIDForSourceDir(sourceDir string) string {
	src := filepath.Clean(sourceDir)
	if src == "" || src == compiletime.ProjectDirPlaceholder {
		return compiletime.DefaultRunID
	}
	return pathSepReplacer.Replace(src)
}

// Checkpoint is a snapshot of *compiletime.State plus a version +
// timestamp. Versioned so future schema changes don't break old
// checkpoints.
type Checkpoint struct {
	Version   int                  `json:"version"`
	CreatedAt time.Time            `json:"created_at"`
	Iteration int                   `json:"iteration"`
	State     *compiletime.State   `json:"state"`
}

// CurrentCheckpointVersion is the schema version for new checkpoints.
const CurrentCheckpointVersion = 1

// Save writes a checkpoint for the given state under
// .artifacts/<runID>/checkpoints/iter-N.json. runID must be non-empty
// (use RunIDForSourceDir if the caller has no explicit run ID).
// Returns the path written so the caller can log it.
func Save(runID string, state *compiletime.State) (string, error) {
	if runID == "" {
		return "", fmt.Errorf("checkpoint: runID must not be empty")
	}
	dir := Dir(runID)
	if err := os.MkdirAll(dir, compiletime.CheckpointDirMode); err != nil {
		return "", fmt.Errorf("checkpoint: mkdir %s: %w", dir, err)
	}
	ckpt := Checkpoint{
		Version:   CurrentCheckpointVersion,
		CreatedAt: time.Now().UTC(),
		Iteration: state.Iteration,
		State:     state,
	}
	name := fmt.Sprintf(compiletime.CheckpointFilePattern, state.Iteration)
	path := filepath.Join(dir, name)
	data, err := json.MarshalIndent(ckpt, "", "  ")
	if err != nil {
		return "", fmt.Errorf("checkpoint: marshal: %w", err)
	}
	// Atomic write: temp file then rename.
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, compiletime.CheckpointFileMode); err != nil {
		return "", fmt.Errorf("checkpoint: write %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return "", fmt.Errorf("checkpoint: rename %s -> %s: %w", tmp, path, err)
	}
	return path, nil
}

// LoadLatest returns the most recent checkpoint for runID, or
// (nil, nil) if no checkpoint exists. Errors are returned only for I/O
// or JSON-decode failures.
func LoadLatest(runID string) (*Checkpoint, error) {
	if runID == "" {
		return nil, fmt.Errorf("checkpoint: runID must not be empty")
	}
	dir := Dir(runID)
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
		if e.IsDir() || filepath.Ext(e.Name()) != compiletime.CheckpointExt {
			continue
		}
		match := iterRe.FindStringSubmatch(e.Name())
		if match == nil {
			continue
		}
		iter, err := strconv.Atoi(match[1])
		if err != nil {
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
	if ckpt.Version != CurrentCheckpointVersion {
		return nil, fmt.Errorf("checkpoint: version mismatch in %s: got %d, want %d",
			latestPath, ckpt.Version, CurrentCheckpointVersion)
	}
	return &ckpt, nil
}

// Cleanup removes all checkpoints for runID (typically called on
// successful completion). Returns nil if the directory doesn't exist.
func Cleanup(runID string) error {
	dir := Dir(runID)
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("checkpoint: cleanup %s: %w", dir, err)
	}
	return nil
}
