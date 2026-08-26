package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildDirectoryTree(t *testing.T) {
	tmpDir := t.TempDir()
	sampleC := filepath.Join(tmpDir, "main.c")
	_ = os.WriteFile(sampleC, []byte("#include <stdio.h>\nint main() { return 0; }\n"), 0644)

	treeStr, files, err := BuildDirectoryTree(tmpDir, 3)
	if err != nil {
		t.Fatalf("failed to build directory tree: %v", err)
	}

	if treeStr == "" {
		t.Errorf("expected non-empty tree string")
	}
	if len(files) == 0 {
		t.Errorf("expected files in directory tree, got 0")
	}
}

func TestParseFileStructure(t *testing.T) {
	tmpDir := t.TempDir()
	cFile := filepath.Join(tmpDir, "sample.c")
	sampleCode := `#include <stdio.h>
struct Item {
    char *name;
    int value;
};
int calculate(struct Item *item) {
    return item->value * 2;
}
`
	if err := os.WriteFile(cFile, []byte(sampleCode), 0644); err != nil {
		t.Fatalf("failed to write test fixture: %v", err)
	}

	structOut, err := ParseFileStructure(cFile)
	if err != nil {
		t.Fatalf("failed to parse file structure: %v", err)
	}

	if structOut.Language != "c" {
		t.Errorf("expected language c, got %s", structOut.Language)
	}
	if len(structOut.Elements) == 0 {
		t.Errorf("expected elements extracted, got 0")
	}
}

func TestFSTools(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	fsTools := NewFSTools()
	if len(fsTools) != 3 {
		t.Errorf("expected 3 fs tools, got %d", len(fsTools))
	}

	// Write
	err := os.WriteFile(testFile, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Read
	data, err := os.ReadFile(testFile)
	if err != nil || string(data) != "hello world" {
		t.Errorf("read mismatch: %s", string(data))
	}
}

func TestPAToolsCreation(t *testing.T) {
	paTools := NewPATools()
	if len(paTools) != 2 {
		t.Errorf("expected 2 PA tools, got %d", len(paTools))
	}
	info, err := paTools[0].Info(context.Background())
	if err != nil || info.Name != "get_directory_tree" {
		t.Errorf("unexpected tool info: %v, %v", info, err)
	}
}

func TestLSPToolsCreation(t *testing.T) {
	lspTools := NewLSPTools()
	if len(lspTools) != 6 {
		t.Errorf("expected 6 LSP tools, got %d", len(lspTools))
	}
}
