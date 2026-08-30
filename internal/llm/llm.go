package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	ollamaapi "github.com/eino-contrib/ollama/api"

	"MAgHARCM/internal/config"
)

// Models contains initialized ChatModel instances for reasoning and coding.
type Models struct {
	Reasoning model.BaseChatModel
	Coding    model.BaseChatModel
}

// defaultChatOptions sets hard caps and stop tokens that prevent base
// (non-instruction-tuned) models from rambling past the answer. Without
// these, qwen2.5-coder:7b-base-q5_0 in particular can emit thousands of
// tokens of unrelated text after the response is complete — hanging the
// generation until Ollama's default num_predict=-1 cut-off fires much
// later. Stop tokens here match the delimiters our prompts ask the model
// to emit at the end of a structured response.
//
// Per-model overrides can be set later via cfg if needed; for now the
// defaults are conservative enough that no observed 7B/4B model in this
// project requires a larger budget.
var defaultChatOptions = &ollamaapi.Options{
	NumPredict: 4096,
	Stop: []string{
		"\n### ",        // qwen-instruct style marker
		"\n\n### ",      // qwen-instruct alt
		"<|im_end|>",    // chatml stop
		"<|endoftext|>", // gpt2-style
	},
	RepeatLastN: 64,
}

// NewModels initializes the Ollama ChatModels for reasoning and coding.
func NewModels(ctx context.Context, cfg *config.Config) (*Models, error) {
	reasoningModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: cfg.OllamaBaseURL,
		Model:   cfg.ReasoningModel,
		Timeout: cfg.Timeout,
		Options: defaultChatOptions,
	})
	if err != nil {
		return nil, err
	}

	codingModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: cfg.OllamaBaseURL,
		Model:   cfg.CodingModel,
		Timeout: cfg.Timeout,
		Options: defaultChatOptions,
	})
	if err != nil {
		return nil, err
	}

	return &Models{
		Reasoning: reasoningModel,
		Coding:    codingModel,
	}, nil
}

// MustNewModels initializes models and panics on error.
func MustNewModels(ctx context.Context, cfg *config.Config) *Models {
	models, err := NewModels(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return models
}
