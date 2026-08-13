package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/deep"

	"MAgHARCM/internal/MAgHARCM"
)

func printMessageOutput(w io.Writer, event *adk.AgentEvent) {
	if event.Err != nil || event.Output == nil || event.Output.MessageOutput == nil {
		return
	}
	mo := event.Output.MessageOutput
	if mo.IsStreaming && mo.MessageStream != nil {
		defer mo.MessageStream.Close()
		for {
			chunk, err := mo.MessageStream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				break
			}
			if chunk != nil {
				if chunk.ReasoningContent != "" {
					io.WriteString(w, chunk.ReasoningContent)
				}
				if chunk.Content != "" {
					io.WriteString(w, chunk.Content)
				}
			}
		}
	} else if mo.Message != nil {
		if mo.Message.ReasoningContent != "" {
			io.WriteString(w, mo.Message.ReasoningContent)
		}
		if mo.Message.Content != "" {
			io.WriteString(w, mo.Message.Content)
		}
	}
}

func main() {
	ctx := context.Background()

	reasoningModel := MAgHARCM.Must(ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "lfm2.5:8b-a1b-q4_K_M",
		// Model:   "gpt-oss:20b",
	}))

	codingModel := MAgHARCM.Must(ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL",
	}))

	codingAgent := MAgHARCM.Must(adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:  "coding-agent",
		Model: codingModel,
	}))

	agent := MAgHARCM.Must(deep.New(ctx, &deep.Config{
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
	iter := runner.Query(ctx, "9.11 or 9.9 which is bigger")
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			fmt.Fprintf(os.Stderr, "event error: %v\n", event.Err)
			continue
		}
		printMessageOutput(os.Stdout, event)
	}
	fmt.Println()
}
