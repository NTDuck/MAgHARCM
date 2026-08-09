package MAgHARCM

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/schema"
)

const ModelName = "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL"

func Ask(ctx context.Context, baseURL, question string) (string, error) {
	cm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: baseURL,
		Model:   ModelName,
	})
	if err != nil {
		return "", err
	}

	res, err := cm.Generate(ctx, []*schema.Message{
		schema.UserMessage(question),
	})
	if err != nil {
		return "", err
	}

	return res.Content, nil
}
