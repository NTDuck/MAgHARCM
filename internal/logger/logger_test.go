package logger

import (
	"bytes"
	"strings"
	"sync"
	"testing"
)

func TestEmitWritesExpectedLines(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	t.Cleanup(func() { SetOutput(nil) })

	LogAgent("Analyzer", "Starting analysis")
	LogTool("get_file_structure", "Extracted elements")
	LogStep("Processing step")
	LogWarning("low coverage")
	LogError("bad thing")
	LogValidation("build succeeded")

	s := buf.String()
	for _, want := range []string{
		"[Analyzer] Starting analysis",
		"[Tool: `get_file_structure`] Extracted elements",
		"Processing step",
		"[WARNING] low coverage",
		"[ERROR] bad thing",
		"[Validation] build succeeded",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in:\n%s", want, s)
		}
	}
}

func TestEmitThreadsafe(t *testing.T) {
	var buf bytes.Buffer
	SetOutput(&buf)
	t.Cleanup(func() { SetOutput(nil) })

	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			LogStep("hello")
		}()
	}
	wg.Wait()

	if got := strings.Count(buf.String(), "hello"); got != n {
		t.Errorf("expected %d lines, got %d", n, got)
	}
}
