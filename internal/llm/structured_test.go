package llm

import (
	"strings"
	"testing"
)

func TestExtractJSONObject(t *testing.T) {
	cases := []struct {
		name    string
		in      string
		want    string
		wantOK  bool
	}{
		{
			name:   "plain JSON object",
			in:     `{"overview": "x"}`,
			want:   `{"overview": "x"}`,
			wantOK: true,
		},
		{
			name:   "JSON wrapped in markdown fence",
			in:     "```json\n{\"overview\": \"x\"}\n```",
			want:   `{"overview": "x"}`,
			wantOK: true,
		},
		{
			name:   "prose around JSON",
			in:     "Here you go:\n{\"a\": 1, \"b\": 2}\nEnjoy.",
			want:   `{"a": 1, "b": 2}`,
			wantOK: true,
		},
		{
			name:   "brace inside string is honoured",
			in:     `{"a": "}{", "b": 1}`,
			want:   `{"a": "}{", "b": 1}`,
			wantOK: true,
		},
		{
			name:   "nested objects",
			in:     `{"x": {"y": 1}, "z": [1,2]}`,
			want:   `{"x": {"y": 1}, "z": [1,2]}`,
			wantOK: true,
		},
		{
			name:   "no JSON",
			in:     "no braces here",
			want:   "",
			wantOK: false,
		},
		{
			name:   "unbalanced",
			in:     `{"a": 1`,
			want:   "",
			wantOK: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := extractJSONObject(tc.in)
			if ok != tc.wantOK {
				t.Fatalf("ok mismatch: got %v want %v", ok, tc.wantOK)
			}
			if got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestUnwrapToolCallEnvelope(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		want   string
		wantOK bool
	}{
		{
			name:   "object arguments",
			in:     `{"name": "emit_analysis", "arguments": {"overview": "x"}}`,
			want:   `{"overview": "x"}`,
			wantOK: true,
		},
		{
			name:   "double-encoded string arguments",
			in:     `{"name": "emit_analysis", "arguments": "{\"overview\": \"x\"}"}`,
			want:   `{"overview": "x"}`,
			wantOK: true,
		},
		{
			name:   "missing name",
			in:     `{"arguments": {"overview": "x"}}`,
			want:   "",
			wantOK: false,
		},
		{
			name:   "missing arguments",
			in:     `{"name": "emit_analysis"}`,
			want:   "",
			wantOK: false,
		},
		{
			name:   "unrelated JSON",
			in:     `{"overview": "x"}`,
			want:   "",
			wantOK: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := unwrapToolCallEnvelope(tc.in)
			if ok != tc.wantOK {
				t.Fatalf("ok mismatch: got %v want %v", ok, tc.wantOK)
			}
			if got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}

// TestExtract_RoundTrip exercises the Path-2 fallback end-to-end with the
// shape of content qwen3-30b-a3b-thinking emits in practice: prose, then
// JSON-as-content, both bare and tool-envelope wrapped.
func TestExtractFallbackShapes(t *testing.T) {
	type schema struct {
		Overview string `json:"overview"`
	}
	// simulate content with envelope
	envelope := `{"name":"emit_analysis","arguments":{"overview":"hello"}}`
	got, ok := unwrapToolCallEnvelope(envelope)
	if !ok || !strings.Contains(got, `"hello"`) {
		t.Fatalf("envelope unwrap failed: %q ok=%v", got, ok)
	}
	// simulate content bare
	bare := `{"overview":"hello"}`
	obj, ok := extractJSONObject(bare)
	if !ok {
		t.Fatalf("bare extract failed")
	}
	if !strings.Contains(obj, `"hello"`) {
		t.Fatalf("bare obj wrong: %s", obj)
	}
}
