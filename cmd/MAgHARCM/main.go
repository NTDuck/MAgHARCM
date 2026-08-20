package main

import (
	"context"

	"github.com/cloudwego/eino-ext/adk/backend/local"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/adk"
	fsmw "github.com/cloudwego/eino/adk/middlewares/filesystem"
	"github.com/cloudwego/eino/adk/prebuilt/deep"

	"MAgHARCM/internal/agent"
	. "MAgHARCM/internal/patterns"
	"MAgHARCM/internal/print"
	"MAgHARCM/internal/prompts"
	"MAgHARCM/internal/tools"
)

func main() {
	ctx := context.Background()

	systemPrompt := prompts.MustLoad("assets/prompts/system_instruction.md")
	taskPrompt := prompts.MustLoad("assets/prompts/c_to_rust_translation.md")

	backend := Must(local.NewBackend(ctx, &local.Config{}))
	fsMiddleware := Must(fsmw.New(ctx, &fsmw.MiddlewareConfig{
		Backend: backend,
		Shell:   backend,
	}))

	reasoningModel := Must(ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL",
	}))

	codingModel := Must(ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL",
	}))

	toolsConfig := tools.NewConfig()
	retryConfig := agent.NewRetryConfig()

	codingAgent := Must(adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:             "coding-agent",
		Description:      "Expert coding agent that writes code, creates files, and runs builds/tests using filesystem tools.",
		Instruction:      systemPrompt,
		Model:            codingModel,
		ToolsConfig:      toolsConfig,
		Handlers:         []adk.ChatModelAgentMiddleware{fsMiddleware},
		ModelRetryConfig: retryConfig,
	}))

	deepAgent := Must(deep.New(ctx, &deep.Config{
		Name:             "default-agent",
		ChatModel:        reasoningModel,
		Instruction:      systemPrompt,
		SubAgents:        []adk.Agent{codingAgent},
		Backend:          backend,
		Shell:            backend,
		ToolsConfig:      toolsConfig,
		ModelRetryConfig: retryConfig,
	}))

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           deepAgent,
		EnableStreaming: true,
	})

	events := runner.Query(ctx, taskPrompt)
	print.Events(events)
}
