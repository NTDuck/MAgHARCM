package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// NewRetryConfig constructs the ModelRetryConfig for autonomous task execution.
func NewRetryConfig() *adk.ModelRetryConfig {
	return &adk.ModelRetryConfig{
		MaxRetries: 5,
		ShouldRetry: func(ctx context.Context, retryCtx *adk.RetryContext) *adk.RetryDecision {
			if retryCtx.Err != nil {
				return &adk.RetryDecision{Retry: true}
			}
			if retryCtx.OutputMessage == nil {
				return &adk.RetryDecision{Retry: false}
			}
			// If the model produced tool calls, allow them to execute
			if len(retryCtx.OutputMessage.ToolCalls) > 0 {
				return &adk.RetryDecision{Retry: false}
			}

			// If the model produced a text response without tool calls, verify deliverables
			reqFiles := []string{
				"experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/Cargo.toml",
				"experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/src/gildedrose.rs",
				"experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/src/main.rs",
			}
			var missing []string
			for _, f := range reqFiles {
				if _, err := os.Stat(f); os.IsNotExist(err) {
					missing = append(missing, f)
				}
			}
			if len(missing) > 0 {
				return &adk.RetryDecision{
					Retry: true,
					ModifiedInputMessages: append(retryCtx.InputMessages, schema.UserMessage(
						fmt.Sprintf("Do not stop or ask questions. Call write_file immediately to create the missing files: %v, then run cargo test.", missing),
					)),
					PersistModifiedInputMessages: true,
				}
			}

			// Verify compilation
			manifestPath := "experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/Cargo.toml"
			cmd := exec.CommandContext(ctx, "cargo", "check", "--manifest-path", manifestPath)
			out, err := cmd.CombinedOutput()
			if err != nil {
				outStr := string(out)
				if len(outStr) > 2000 {
					outStr = outStr[:2000]
				}
				return &adk.RetryDecision{
					Retry: true,
					ModifiedInputMessages: append(retryCtx.InputMessages, schema.UserMessage(
						fmt.Sprintf("Cargo check failed with error:\n```\n%s\n```\nPlease call write_file on `src/gildedrose.rs` / `src/main.rs` to fix this compiler error immediately.", outStr),
					)),
					PersistModifiedInputMessages: true,
				}
			}

			return &adk.RetryDecision{Retry: false}
		},
	}
}
