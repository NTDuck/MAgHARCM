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

func TestLSPProviderProviders(t *testing.T) {
	ctx := context.Background()

	// Native Provider
	native := GetLSPProvider("native")
	if native.Name() != "native" {
		t.Errorf("expected native provider name, got %s", native.Name())
	}

	// ABCoder MCP Provider with fallback
	abcoder := GetLSPProvider("abcoder")
	if abcoder.Name() != "abcoder-mcp" {
		t.Errorf("expected abcoder-mcp provider name, got %s", abcoder.Name())
	}

	// Test Definition on Native
	def, err := native.GetDefinition(ctx, &DefinitionInput{
		Symbol:  "calculate",
		Project: t.TempDir(),
	})
	if err != nil {
		t.Fatalf("unexpected error on definition: %v", err)
	}
	if def.Symbol != "calculate" {
		t.Errorf("expected symbol calculate, got %s", def.Symbol)
	}

	// Test ABCoder fallback execution
	hover, err := abcoder.GetHover(ctx, &HoverInput{
		Symbol: "Item",
	})
	if err != nil {
		t.Fatalf("unexpected error on hover fallback: %v", err)
	}
	if !hover.Found {
		t.Errorf("expected hover found")
	}
}

func TestMultiLanguageExecutionTools(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	// No build file
	out, err := ValidateProjectBuild(ctx, tmpDir, "Rust", "cargo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Success {
		t.Errorf("expected failure when no build config exists")
	}

	// Go build check on actual directory
	goOut, err := ValidateProjectBuild(ctx, ".", "Go", "go")
	if err != nil {
		t.Fatalf("unexpected error on go build check: %v", err)
	}
	if !goOut.Success {
		t.Errorf("expected current repo go build to succeed: %v", goOut.Errors)
	}
}
