package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudwego/eino-ext/adk/backend/local"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"

	. "MAgHARCM/internal/patterns"
)

func TestToolGroups(t *testing.T) {
	ctx := context.Background()
	tmpDir := t.TempDir()

	backend := Must(local.NewBackend(ctx, &local.Config{}))

	// 1. Test FS tools
	fsTools := NewFSTools(backend)
	if len(fsTools) != 6 {
		t.Fatalf("expected 6 FS tools, got %d", len(fsTools))
	}

	testFile := filepath.Join(tmpDir, "test.rs")
	testContent := "pub fn hello() -> &'static str {\n    \"world\"\n}\n"
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// 2. Test LSP tools
	lspTools := NewLSPTools()
	if len(lspTools) != 6 {
		t.Fatalf("expected 6 LSP tools, got %d", len(lspTools))
	}

	// 3. Test PA tools (tree-sitter integration)
	paTools := NewPATools()
	if len(paTools) != 2 {
		t.Fatalf("expected 2 PA tools, got %d", len(paTools))
	}

	// 4. Test Execution tools
	execTools := NewExecutionTools(backend)
	if len(execTools) != 1 {
		t.Fatalf("expected 1 execution tool, got %d", len(execTools))
	}

	// 5. Test Git tools
	gitTools := NewGitTools()
	if len(gitTools) != 3 {
		t.Fatalf("expected 3 git tools, got %d", len(gitTools))
	}

	// 6. Test Validation tools
	valTools := NewValidationTools()
	if len(valTools) != 2 {
		t.Fatalf("expected 2 validation tools, got %d", len(valTools))
	}

	// 7. Test Centralized Registry
	all := AllTools(backend, backend)
	if len(all) != 20 {
		t.Fatalf("expected 20 tools total in registry, got %d", len(all))
	}

	cfg := NewToolsConfig(all...)
	if len(cfg.Tools) != 20 {
		t.Fatalf("expected 20 tools in ToolsConfig, got %d", len(cfg.Tools))
	}

	// 8. Test tree-sitter file structure parser directly
	rustCode := []byte("pub struct GildedRose { pub items: Vec<Item> }\npub fn update_quality() {}\n")
	lang, langName := getTreeSitterLanguage(".rs")
	if lang == nil || langName != "rust" {
		t.Fatalf("expected rust language from getTreeSitterLanguage")
	}
	parser := tree_sitter.NewParser()
	defer parser.Close()
	_ = parser.SetLanguage(lang)
	tree := parser.Parse(rustCode, nil)
	defer tree.Close()
	elements, _ := extractElementsFromAST(tree.RootNode(), rustCode, "rust")
	if len(elements) != 2 {
		t.Fatalf("expected 2 elements from tree-sitter parser, got %d", len(elements))
	}
}
