package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestConsoleSink(t *testing.T) {
	var buf bytes.Buffer
	sink := NewConsoleSink(&buf)

	sink.WriteEvent(Event{
		Type:    EventAgentStart,
		Agent:   "Analyzer",
		Message: "Starting analysis",
	})
	sink.WriteEvent(Event{
		Type:    EventToolCall,
		Tool:    "get_file_structure",
		Message: "Extracted elements",
	})

	out := buf.String()
	if out == "" {
		t.Errorf("expected console sink output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("[Tool: `get_file_structure`]")) {
		t.Errorf("expected backticked tool format [Tool: `get_file_structure`], got: %s", out)
	}
}

func TestJSONFileSink(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "events.jsonl")

	sink, err := NewJSONFileSink(logPath)
	if err != nil {
		t.Fatalf("failed to create json file sink: %v", err)
	}

	sink.WriteEvent(Event{
		Type:    EventStep,
		Message: "Processing step",
	})
	_ = sink.Close()

	data, err := os.ReadFile(logPath)
	if err != nil || len(data) == 0 {
		t.Fatalf("expected non-empty jsonl file: %v", err)
	}
}
