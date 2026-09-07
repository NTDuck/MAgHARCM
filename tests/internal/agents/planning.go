package agents_test

import (
	"MAgHARCM/internal/agents"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDefaultProjectSkeleton(t *testing.T) {
	// 1. Rust skeleton
	rustSkel := agents.DefaultProjectSkeleton(".artifacts/sample_project/rust", "Rust", []string{"foo.c:bar"})
	if _, ok := rustSkel["Cargo.toml"]; !ok {
		t.Errorf("expected Cargo.toml in rust skeleton")
	}

	// 2. Go skeleton
	goSkel := agents.DefaultProjectSkeleton(".artifacts/sample_project/go", "Go", []string{"foo.c:bar"})
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

func TestComputeReverseTopoOrder(t *testing.T) {
	// Test acyclic graph
	items := []string{"main.go", "util.go", "db.go"}
	deps := map[string][]string{
		"main.go": {"util.go", "db.go"},
		"util.go": {"db.go"},
		"db.go":   {},
	}
	order := agents.ComputeReverseTopoOrder(items, deps)
	// db.go should come before util.go, and util.go before main.go
	indexOf := func(slice []string, val string) int {
		for i, v := range slice {
			if v == val {
				return i
			}
		}
		return -1
	}
	if indexOf(order, "db.go") > indexOf(order, "util.go") {
		t.Errorf("expected db.go before util.go, got order %v", order)
	}
	if indexOf(order, "util.go") > indexOf(order, "main.go") {
		t.Errorf("expected util.go before main.go, got order %v", order)
	}

	// Test cyclic graph with back-edge removal
	cycleItems := []string{"A", "B", "C"}
	cycleDeps := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {"A"}, // cycle!
	}
	cycleOrder := agents.ComputeReverseTopoOrder(cycleItems, cycleDeps)
	if len(cycleOrder) != 3 {
		t.Fatalf("expected 3 items in cycleOrder, got %d: %v", len(cycleOrder), cycleOrder)
	}
}
