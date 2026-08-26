package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildDirectoryTree(t *testing.T) {
	sampleDir := "../../assets/samples/GildedRose-Refactoring-Kata/C"
	if _, err := os.Stat(sampleDir); err != nil {
		sampleDir = "assets/samples/GildedRose-Refactoring-Kata/C"
	}

	treeStr, files, err := BuildDirectoryTree(sampleDir, 3)
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
	cFile := "../../assets/samples/GildedRose-Refactoring-Kata/C/GildedRose.c"
	if _, err := os.Stat(cFile); err != nil {
		cFile = "assets/samples/GildedRose-Refactoring-Kata/C/GildedRose.c"
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
