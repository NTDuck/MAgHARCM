package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ollama"
)

// stubServer returns an httptest.Server that mimics Ollama's /api/chat
// response shape and writes the given content (with no tool_calls).
func stubServer(content string) *httptest.Server {
	c, err := json.Marshal(content)
	if err != nil {
		panic(err) // test stub; marshal of a Go string never fails.
	}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body := `{"model": "stub", "created_at": "2026-09-09T00:00:00Z", "done": true, "done_reason": "stop", "message": {"role": "assistant", "content": ` + string(c) + `, "tool_calls": null}}`
		_, _ = w.Write([]byte(body))
	}))
}

func newStubChatModel(t *testing.T, baseURL string) *ollama.ChatModel {
	t.Helper()
	m, err := ollama.NewChatModel(context.Background(), &ollama.ChatModelConfig{
		BaseURL: baseURL,
		Model:   "stub",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		t.Fatalf("stub chat model build: %v", err)
	}
	return m
}

// TestExtractFallback_Path2: model emits bare schema JSON as content.
func TestExtractFallback_Path2(t *testing.T) {
	srv := stubServer(`{"overview": "the project adds two integers"}`)
	defer srv.Close()
	m := newStubChatModel(t, srv.URL)
	type schema struct {
		Overview string `json:"overview"`
	}
	ext, err := NewStructuredExtractor[schema](m, "emit_analysis",
		"Emit the structured source-code analysis")
	if err != nil {
		t.Fatalf("extractor build: %v", err)
	}
	got, err := ext.Extract(context.Background(),
		"You are an analyst. Always respond by calling the emit_analysis tool.",
		"Analyse int add(int a, int b) { return a+b; }")
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if !strings.Contains(got.Overview, "add") {
		t.Fatalf("overview wrong: %q", got.Overview)
	}
}

// TestExtractFallback_Envelope: model wraps the args in {name, arguments}.
func TestExtractFallback_Envelope(t *testing.T) {
	srv := stubServer(`{"name":"emit_analysis","arguments":{"overview":"envelope form"}}`)
	defer srv.Close()
	m := newStubChatModel(t, srv.URL)
	type schema struct {
		Overview string `json:"overview"`
	}
	ext, err := NewStructuredExtractor[schema](m, "emit_analysis", "Emit analysis")
	if err != nil {
		t.Fatalf("extractor build: %v", err)
	}
	got, err := ext.Extract(context.Background(), "system", "user")
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if !strings.Contains(got.Overview, "envelope") {
		t.Fatalf("overview wrong: %q", got.Overview)
	}
}

// TestExtractFallback_EnvelopeStringArgs: arguments is a JSON-encoded string.
func TestExtractFallback_EnvelopeStringArgs(t *testing.T) {
	srv := stubServer(`{"name":"emit_analysis","arguments":"{\"overview\":\"string args form\"}"}`)
	defer srv.Close()
	m := newStubChatModel(t, srv.URL)
	type schema struct {
		Overview string `json:"overview"`
	}
	ext, err := NewStructuredExtractor[schema](m, "emit_analysis", "Emit analysis")
	if err != nil {
		t.Fatalf("extractor build: %v", err)
	}
	got, err := ext.Extract(context.Background(), "system", "user")
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if !strings.Contains(got.Overview, "string args") {
		t.Fatalf("overview wrong: %q", got.Overview)
	}
}

// TestExtractFallback_RefusesFreeText: prose-only content fails cleanly.
func TestExtractFallback_RefusesFreeText(t *testing.T) {
	srv := stubServer("# Source Project Research\n## Overview\nThis project...")
	defer srv.Close()
	m := newStubChatModel(t, srv.URL)
	type schema struct {
		Overview string `json:"overview"`
	}
	ext, err := NewStructuredExtractor[schema](m, "emit_analysis", "Emit analysis")
	if err != nil {
		t.Fatalf("extractor build: %v", err)
	}
	_, err = ext.Extract(context.Background(), "system", "user")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "no tool call arguments") {
		t.Fatalf("expected 'no tool call arguments' error, got: %v", err)
	}
}
