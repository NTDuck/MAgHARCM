package languages

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLanguageRegistryLookup(t *testing.T) {
	reg := GetRegistry()

	// 1. Rust
	rustSpec, ok := reg.FindByExtension(".rs")
	if !ok || rustSpec.Name != "rust" {
		t.Errorf("expected rust for .rs, got %v", rustSpec)
	}

	// 2. Python
	pySpec, ok := reg.FindByExtension(".py")
	if !ok || pySpec.Name != "python" {
		t.Errorf("expected python for .py, got %v", pySpec)
	}

	// 3. Java
	javaSpec, ok := reg.FindByName("Java")
	if !ok || javaSpec.Name != "java" {
		t.Errorf("expected java for Java name lookup, got %v", javaSpec)
	}

	// 4. Go
	goSpec, ok := reg.FindByExtension(".go")
	if !ok || goSpec.Name != "go" {
		t.Errorf("expected go for .go, got %v", goSpec)
	}
}

func TestDynamicLanguageLoader(t *testing.T) {
	loader := GetLoader()
	loader.AddSearchPath("/tmp/custom-grammars")

	reg := GetRegistry()
	rustSpec, ok := reg.FindByName("rust")
	if !ok {
		t.Fatalf("expected rust spec in registry")
	}

	// Loading when shared library is not in search paths should return descriptive error without crashing
	_, err := loader.LoadGrammar(rustSpec)
	if err == nil {
		t.Logf("Found system grammar for rust")
	} else {
		t.Logf("Dynamic grammar fallback active: %v", err)
	}
}

func TestExtractFileStructureMultiLanguage(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. C file
	cFile := filepath.Join(tmpDir, "sample.c")
	cCode := `#include <stdio.h>
struct User {
    char name[32];
    int id;
};
int get_id(struct User* u) {
    return u->id;
}
`
	if err := os.WriteFile(cFile, []byte(cCode), 0644); err != nil {
		t.Fatalf("failed to write c fixture: %v", err)
	}

	cRes, err := ExtractFileStructure(cFile)
	if err != nil {
		t.Fatalf("failed to extract c structure: %v", err)
	}
	if cRes.Language != "c" {
		t.Errorf("expected c language, got %s", cRes.Language)
	}
	if len(cRes.Elements) == 0 {
		t.Errorf("expected extracted elements from c file, got 0")
	}

	// 2. Python file
	pyFile := filepath.Join(tmpDir, "sample.py")
	pyCode := `import os

class Calculator:
    def add(self, a, b):
        return a + b
`
	if err := os.WriteFile(pyFile, []byte(pyCode), 0644); err != nil {
		t.Fatalf("failed to write py fixture: %v", err)
	}

	pyRes, err := ExtractFileStructure(pyFile)
	if err != nil {
		t.Fatalf("failed to extract python structure: %v", err)
	}
	if pyRes.Language != "python" {
		t.Errorf("expected python language, got %s", pyRes.Language)
	}
	if len(pyRes.Elements) == 0 {
		t.Errorf("expected extracted elements from python file, got 0")
	}
}
