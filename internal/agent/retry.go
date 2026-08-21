package agent

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
)

// ValidatorFunc checks whether the current state meets the agent goal or requires feedback.
// Returns shouldRetry (bool) and feedback message (string).
type ValidatorFunc func(ctx context.Context) (shouldRetry bool, feedback string)

// RetryOption configures retry behavior.
type RetryOption func(*retryOptions)

type retryOptions struct {
	maxRetries int
	validators []ValidatorFunc
}

// WithMaxRetries sets custom max retries.
func WithMaxRetries(n int) RetryOption {
	return func(o *retryOptions) {
		if n > 0 {
			o.maxRetries = n
		}
	}
}

// WithValidator attaches a domain validator to the retry cycle.
func WithValidator(v ValidatorFunc) RetryOption {
	return func(o *retryOptions) {
		if v != nil {
			o.validators = append(o.validators, v)
		}
	}
}

// NewRetryConfig constructs an experiment-agnostic, prompt-agnostic ModelRetryConfig.
func NewRetryConfig(opts ...RetryOption) *adk.ModelRetryConfig {
	opt := &retryOptions{
		maxRetries: 5,
	}
	for _, o := range opts {
		o(opt)
	}

	return &adk.ModelRetryConfig{
		MaxRetries: opt.maxRetries,
		ShouldRetry: func(ctx context.Context, retryCtx *adk.RetryContext) *adk.RetryDecision {
			// 1. Retry on model transport/inference errors
			if retryCtx.Err != nil {
				return &adk.RetryDecision{Retry: true}
			}
			if retryCtx.OutputMessage == nil {
				return &adk.RetryDecision{Retry: false}
			}

			// 2. If the model emitted tool calls, let the tool pipeline execute
			if len(retryCtx.OutputMessage.ToolCalls) > 0 {
				return &adk.RetryDecision{Retry: false}
			}

			// 3. If the model produced empty or whitespace content without tool calls
			if strings.TrimSpace(retryCtx.OutputMessage.Content) == "" {
				return &adk.RetryDecision{
					Retry: true,
					ModifiedInputMessages: append(retryCtx.InputMessages, schema.UserMessage(
						"Please proceed with the requested task using the appropriate tools.",
					)),
					PersistModifiedInputMessages: true,
				}
			}

			// 4. Run any attached domain validators (if configured by caller)
			for _, v := range opt.validators {
				if shouldRetry, feedback := v(ctx); shouldRetry && feedback != "" {
					return &adk.RetryDecision{
						Retry:                        true,
						ModifiedInputMessages:        append(retryCtx.InputMessages, schema.UserMessage(feedback)),
						PersistModifiedInputMessages: true,
					}
				}
			}

			return &adk.RetryDecision{Retry: false}
		},
	}
}
