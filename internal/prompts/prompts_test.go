package prompts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMustLoad(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "prompt.md")
	content := "test prompt content"
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	loaded := MustLoad(tmp)
	if loaded != content {
		t.Fatalf("expected %q, got %q", content, loaded)
	}
}
