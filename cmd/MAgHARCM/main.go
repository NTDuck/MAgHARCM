package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/MAgHARCM"
)

func main() {
	ctx := context.Background()

	model := MAgHARCM.Must(ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL",
	}))

	res := MAgHARCM.Must(model.Generate(ctx, []*schema.Message{
		schema.UserMessage("9.11 or 9.9 which is bigger"),
	}))

	fmt.Println(res.Content)
}
