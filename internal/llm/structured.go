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
	out := [][]*schema.Message{}
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

	var args string
	for _, tc := range resp.ToolCalls {
		if tc.Function.Arguments == "" {
			continue
		}
		args = tc.Function.Arguments
		break
	}
	if args == "" {
		return zero, fmt.Errorf("structured: model returned no tool call arguments (finish_reason=%v, content=%q)",
			resp.ResponseMeta, truncateForErr(resp.Content))
	}

	var typed T
	if err := json.Unmarshal([]byte(args), &typed); err != nil {
		return zero, fmt.Errorf("structured: unmarshal tool args: %w (raw=%s)", err, truncateForErr(args))
	}
	_ = out // keep out reachable for future streaming variant
	return typed, nil
}

func truncateForErr(s string) string {
	const cap = 240
	if len(s) <= cap {
		return s
	}
	return s[:cap] + "…"
}
