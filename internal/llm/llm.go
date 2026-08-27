package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"

	"MAgHARCM/internal/config"
)

// Models contains initialized ChatModel instances for reasoning and coding.
type Models struct {
	Reasoning model.BaseChatModel
	Coding    model.BaseChatModel
}

// NewModels initializes the Ollama ChatModels for reasoning and coding.
func NewModels(ctx context.Context, cfg *config.Config) (*Models, error) {
	reasoningModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: cfg.OllamaBaseURL,
		Model:   cfg.ReasoningModel,
		Timeout: cfg.Timeout,
	})
	if err != nil {
		return nil, err
	}

	codingModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: cfg.OllamaBaseURL,
		Model:   cfg.CodingModel,
		Timeout: cfg.Timeout,
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
