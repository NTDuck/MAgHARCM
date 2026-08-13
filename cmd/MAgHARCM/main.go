package main

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/deep"

	. "MAgHARCM/internal/patterns"
)

func main() {
	ctx := context.Background()

	reasoningModel := Must(ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "lfm2.5:8b-a1b-q4_K_M",
		// Model:   "gpt-oss:20b",
		// Model:   "qwen3:30b-a3b-thinking-2507-q4_K_M",
	}))

	codingModel := Must(ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL",
	}))

	codingAgent := Must(adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:  "coding-agent",
		Model: codingModel,
	}))

	agent := Must(deep.New(ctx, &deep.Config{
		Name:        "default-agent",
		ChatModel:   reasoningModel,
		SubAgents:   []adk.Agent{codingAgent},
		ToolsConfig: adk.ToolsConfig{},
	}))

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})

	// https://mindyourdecisions.com/blog/2024/08/07/9-11-is-larger-than-9-9-according-to-ai/
	events := runner.Query(ctx, "9.11 or 9.9 which is bigger")

}
