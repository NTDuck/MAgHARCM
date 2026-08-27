package agents

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDefaultProjectSkeleton(t *testing.T) {
	// 1. Rust skeleton
	rustSkel := defaultProjectSkeleton(".artifacts/sample_project/rust", "Rust", []string{"foo.c:bar"})
	if _, ok := rustSkel["Cargo.toml"]; !ok {
		t.Errorf("expected Cargo.toml in rust skeleton")
	}

	// 2. Go skeleton
	goSkel := defaultProjectSkeleton(".artifacts/sample_project/go", "Go", []string{"foo.c:bar"})
	if _, ok := goSkel["go.mod"]; !ok {
		t.Errorf("expected go.mod in go skeleton")
	}

	// 3. Test compilation of Rust skeleton if cargo is available
	if _, err := exec.LookPath("cargo"); err == nil {
		tmpDir := t.TempDir()
		for relPath, content := range rustSkel {
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
			t.Fatalf("cargo check failed on defaultProjectSkeleton: %v\nOutput:\n%s", err, string(out))
		}
	}
}
