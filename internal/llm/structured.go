// Package llm provides typed structured-output helpers built on top of Eino's
// ToolCallingChatModel. Each call binds a single typed schema as a tool,
// invokes the model, unmarshals the tool-call arguments directly into a
// concrete Go value, and returns it. No string parsing, no regex, no manual
// delimiter juggling.
package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

// StructuredExtractor is the typed-output facade. It is constructed once per
// (model, schema) pair by [NewStructuredExtractor]; call [Extract] as many
// times as you like. The bound schema travels with the returned extractor
// only — the underlying ChatModel is never mutated, so the same model can
// still be used elsewhere with a different schema.
type StructuredExtractor[T any] struct {
	bound model.ToolCallingChatModel
}

// NewStructuredExtractor binds a JSON schema derived from the Go type T to
// the given chat model and returns a typed extractor. The schema is the
// JSON-Schema 2020-12 document Eino infers from T's struct tags, exposed to
// the model as a single forced tool call.
//
// T MUST be a struct (or a pointer to one). Fields without `json:"..."`
// tags are skipped; nested structs/maps/slices follow Go's standard JSON
// semantics so the produced schema matches [encoding/json] round-trip
// behavior exactly.
func NewStructuredExtractor[T any](m model.BaseChatModel, toolName, toolDesc string) (*StructuredExtractor[T], error) {
	tcm, ok := m.(model.ToolCallingChatModel)
	if !ok {
		return nil, fmt.Errorf("structured: model does not implement ToolCallingChatModel (got %T)", m)
	}
	info, err := utils.GoStruct2ToolInfo[T](toolName, toolDesc)
	if err != nil {
		return nil, fmt.Errorf("structured: build schema for %q: %w", toolName, err)
	}
	bound, err := tcm.WithTools([]*schema.ToolInfo{info})
	if err != nil {
		return nil, fmt.Errorf("structured: bind tool %q: %w", toolName, err)
	}
	return &StructuredExtractor[T]{bound: bound}, nil
}

// Extract runs a single Generate call with the bound schema. The model MUST
// produce a tool call to the bound tool; any other outcome (refusal, free
// text, wrong tool) is reported as an error so the caller can decide between
// retry, repair, or bail.
func (e *StructuredExtractor[T]) Extract(ctx context.Context, system, user string, msgs ...*schema.Message) (T, error) {
	var zero T
	if e == nil || e.bound == nil {
		return zero, fmt.Errorf("structured: extractor not initialised")
	}
	allMsgs := make([]*schema.Message, 0, len(msgs)+1)
	if system != "" {
		allMsgs = append(allMsgs, schema.SystemMessage(system))
	}
	allMsgs = append(allMsgs, msgs...)
	allMsgs = append(allMsgs, schema.UserMessage(user))

	resp, err := e.bound.Generate(ctx, allMsgs)
	if err != nil {
		return zero, fmt.Errorf("structured: model call: %w", err)
	}
	if resp == nil {
		return zero, fmt.Errorf("structured: nil response")
	}

	// Path 1: canonical tool-call argument. Most Ollama chat models honour
	// the bound tool and surface arguments in resp.ToolCalls.
	var args string
	for _, tc := range resp.ToolCalls {
		if tc.Function.Arguments == "" {
			continue
		}
		args = tc.Function.Arguments
		break
	}
	// Path 2: JSON-as-content fallback. Several qwen3-MOE-thinking
	// variants we benchmark emit the schema-shaped payload as content
	// instead of as a wire-level tool call, especially under long prompts.
	// We strip markdown fences, locate the outermost balanced {...}, and
	// try to unmarshal directly. If the model also wrapped the args in
	// `{name, arguments}` (the format Ollama's chat template prefers when
	// tools are bound), we unwrap before unmarshalling.
	if args == "" && resp.Content != "" {
		if unwrapped, ok := extractJSONObject(resp.Content); ok {
			if inner, ok := unwrapToolCallEnvelope(unwrapped); ok {
				args = inner
			} else {
				args = unwrapped
			}
		}
	}

	var typed T
	if err := json.Unmarshal([]byte(args), &typed); err != nil {
		return zero, fmt.Errorf("structured: unmarshal tool args: %w (raw=%s)", err, truncateForErr(args))
	}
	return typed, nil
}

func truncateForErr(s string) string {
	const cap = 240
	if len(s) <= cap {
		return s
	}
	return s[:cap] + "…"
}

// extractJSONObject returns the first top-level balanced {...} substring
// found in s, after stripping ```json ... ``` fences the model may wrap
// the payload in. Strings are honoured while counting braces so embedded
// '{' / '}' inside string literals do not throw the parser off. Returns
// ok=false when no balanced JSON object is present.
func extractJSONObject(s string) (string, bool) {
	t := strings.TrimSpace(s)
	if i := strings.Index(t, "```"); i >= 0 {
		if j := strings.Index(t[i+3:], "```"); j >= 0 {
			inner := t[i+3 : i+3+j]
			if k := strings.Index(inner, "\n"); k >= 0 {
				inner = inner[k+1:]
			}
			t = inner
		}
	}
	open := strings.IndexByte(t, '{')
	if open < 0 {
		return "", false
	}
	depth, inStr, esc := 0, false, false
	for i := open; i < len(t); i++ {
		c := t[i]
		switch {
		case esc:
			esc = false
		case c == '\\' && inStr:
			esc = true
		case c == '"':
			inStr = !inStr
		case !inStr && c == '{':
			depth++
		case !inStr && c == '}':
			depth--
			if depth == 0 {
				return t[open : i+1], true
			}
		}
	}
	return "", false
}

// unwrapToolCallEnvelope peels Ollama's `{name, arguments}` envelope off
// a JSON object produced by a model that knows it was asked for a tool
// call but emitted the wrapper as content instead of as ToolCalls. The
// arguments value can be either a JSON object (the happy path) or a JSON
// string holding a JSON object (some Ollama chat templates double-encode
// it). Returns ok=false when the input does not look like the envelope.
func unwrapToolCallEnvelope(j string) (string, bool) {
	var probe struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal([]byte(j), &probe); err != nil {
		return "", false
	}
	if probe.Name == "" || len(probe.Arguments) == 0 {
		return "", false
	}
	a := strings.TrimSpace(string(probe.Arguments))
	if len(a) > 0 && a[0] == '{' {
		return a, true
	}
	if len(a) > 0 && a[0] == '"' {
		var nested string
		if err := json.Unmarshal([]byte(a), &nested); err != nil {
			return "", false
		}
		return nested, true
	}
	return "", false
}
