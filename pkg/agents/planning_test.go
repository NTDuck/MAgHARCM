package agents

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDefaultRustSkeletonCompiles(t *testing.T) {
	if _, err := exec.LookPath("cargo"); err != nil {
		t.Skip("cargo not found in PATH")
	}

	tmpDir := t.TempDir()
	skeleton := defaultRustSkeleton()

	for relPath, content := range skeleton {
		fullPath := filepath.Join(tmpDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write skeleton file %s: %v", relPath, err)
		}
	}

	cmd := exec.CommandContext(context.Background(), "cargo", "check")
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("cargo check failed on defaultRustSkeleton: %v\nOutput:\n%s", err, string(out))
	}
}
