package agent

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

func TestNewRetryConfig(t *testing.T) {
	ctx := context.Background()
	cfg := NewRetryConfig(
		WithMaxRetries(3),
		WithValidator(func(ctx context.Context) (bool, string) {
			return false, ""
		}),
	)

	if cfg.MaxRetries != 3 {
		t.Fatalf("expected maxRetries 3, got %d", cfg.MaxRetries)
	}

	// 1. Test error retry
	d1 := cfg.ShouldRetry(ctx, &adk.RetryContext{
		Err: errors.New("inference timeout"),
	})
	if !d1.Retry {
		t.Fatalf("expected retry on error")
	}

	// 2. Test tool calls allowed to proceed without retry
	d2 := cfg.ShouldRetry(ctx, &adk.RetryContext{
		OutputMessage: &schema.Message{
			ToolCalls: []schema.ToolCall{{}},
		},
	})
	if d2.Retry {
		t.Fatalf("expected tool calls to execute without retry")
	}

	// 3. Test empty output retry
	d3 := cfg.ShouldRetry(ctx, &adk.RetryContext{
		OutputMessage: &schema.Message{
			Content: "",
		},
	})
	if !d3.Retry {
		t.Fatalf("expected empty content to retry")
	}
}
