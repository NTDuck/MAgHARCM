package languages_test

import (
	"MAgHARCM/internal/languages"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/ebitengine/purego"
)

func TestLanguageRegistryLookup(t *testing.T) {
	reg := languages.GetRegistry()

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
	loader := languages.GetLoader()
	loader.AddSearchPath("/tmp/custom-grammars")

	reg := languages.GetRegistry()
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

func TestRuntimePuregoDlopenDlsym(t *testing.T) {
	gccPath, err := exec.LookPath("gcc")
	if err != nil {
		t.Skip("gcc not available for compiling dynamic shared library test")
	}

	tmpDir := t.TempDir()
	cSource := filepath.Join(tmpDir, "test_lib.c")
	soPath := filepath.Join(tmpDir, "libtest_plugin.so")

	// C source exporting a valid safe function
	cCode := `
int calculate_runtime(int a, int b) {
    return a + b;
}
`
	if err := os.WriteFile(cSource, []byte(cCode), 0644); err != nil {
		t.Fatalf("failed to write C source: %v", err)
	}

	cmd := exec.CommandContext(context.Background(), gccPath, "-shared", "-fPIC", "-o", soPath, cSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to compile shared library: %v, output: %s", err, string(out))
	}

	// Test purego dynamic loading and function execution at runtime
	handle, err := purego.Dlopen(soPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		t.Fatalf("purego.Dlopen failed on %s: %v", soPath, err)
	}

	sym, err := purego.Dlsym(handle, "calculate_runtime")
	if err != nil || sym == 0 {
		t.Fatalf("purego.Dlsym failed for symbol calculate_runtime: %v", err)
	}

	var calcFn func(int, int) int
	purego.RegisterFunc(&calcFn, sym)

	res := calcFn(15, 27)
	if res != 42 {
		t.Errorf("expected calculate_runtime(15, 27) == 42, got %d", res)
	}

	t.Logf("Successfully verified runtime purego Dlopen, Dlsym, and C function execution (15 + 27 = %d)", res)
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

	cRes, err := languages.ExtractFileStructure(cFile)
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

	pyRes, err := languages.ExtractFileStructure(pyFile)
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
