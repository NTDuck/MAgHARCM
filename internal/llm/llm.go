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
// later.
//
// NumCtx is intentionally small (1024) to keep the KV cache footprint
// down on the 6GB RTX 3060 we're running on. Larger contexts will
// trigger VRAM pressure when the foreign llama-server (PID 274826 on
// port 8080) is also holding ~3GB of GPU memory; with NumCtx=1024 the
// Ollama model fits in ~1.4GB and survives that contention.
// Stop tokens match the delimiters our prompts ask the model to emit
// at the end of a structured response. `<|im_end|>` is the chatml
// token used by Qwen-style models. We deliberately do NOT stop on
// `\n### ` because our own prompts emit `### FILE:` / `### Source File:`
// / `### Fragment:` headers inside legitimate structured output, and
// premature termination on those headers would truncate every response
// mid-generation.
var defaultChatOptions = &ollamaapi.Options{
	Runner:     ollamaapi.Runner{NumCtx: 1024},
	NumPredict: 4096,
	Stop: []string{
		"<|im_end|>",
		"<|endoftext|>",
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
