package logger_test

import (
	"MAgHARCM/internal/logger"
	"bytes"
	"strings"
	"sync"
	"testing"
)

func TestEmitWritesExpectedLines(t *testing.T) {
	var buf bytes.Buffer
	logger.SetOutput(&buf)
	t.Cleanup(func() { logger.SetOutput(nil) })

	logger.LogAgent("Analyzer", "Starting analysis")
	logger.LogTool("get_file_structure", "Extracted elements")
	logger.LogStep("Processing step")
	logger.LogWarning("low coverage")
	logger.LogError("bad thing")
	logger.LogValidation("build succeeded")

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
	logger.SetOutput(&buf)
	t.Cleanup(func() { logger.SetOutput(nil) })

	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	for range n {
		go func() {
			defer wg.Done()
			logger.LogStep("hello")
		}()
	}
	wg.Wait()

	if got := strings.Count(buf.String(), "hello"); got != n {
		t.Errorf("expected %d lines, got %d", n, got)
	}
}
