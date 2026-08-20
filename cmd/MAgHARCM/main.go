package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/adk/backend/local"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/adk"
	fsmw "github.com/cloudwego/eino/adk/middlewares/filesystem"
	"github.com/cloudwego/eino/adk/prebuilt/deep"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	. "MAgHARCM/internal/patterns"
	"MAgHARCM/internal/print"
)

func main() {
	ctx := context.Background()

	systemPrompt := string(Must(os.ReadFile("prompts/system_instruction.md")))
	taskPrompt := string(Must(os.ReadFile("prompts/c_to_rust_translation.md")))

	backend := Must(local.NewBackend(ctx, &local.Config{}))
	reasoningModel := Must(ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL",
	}))

	codingModel := Must(ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL",
	}))

	fsMiddleware := Must(fsmw.New(ctx, &fsmw.MiddlewareConfig{
		Backend: backend,
		Shell:   backend,
	}))

	toolsConfig := adk.ToolsConfig{
		ToolsNodeConfig: compose.ToolsNodeConfig{
			ExecuteSequentially: true,
			ToolAliases: map[string]compose.ToolAliasConfig{
				"read_file": {
					NameAliases: []string{"read", "cat", "view_file", "readFile"},
					ArgumentsAliases: map[string][]string{
						"file_path": {"path", "filepath", "file", "filename"},
					},
				},
				"write_file": {
					NameAliases: []string{"write", "create_file", "save_file", "writeFile"},
					ArgumentsAliases: map[string][]string{
						"file_path": {"path", "filepath", "file", "filename"},
						"content":   {"text", "data", "contents", "code"},
					},
				},
				"edit_file": {
					NameAliases: []string{"edit", "replace_in_file", "patch_file", "editFile"},
					ArgumentsAliases: map[string][]string{
						"file_path": {"path", "filepath", "file", "filename"},
					},
				},
				"glob": {
					NameAliases: []string{"find", "list_files"},
					ArgumentsAliases: map[string][]string{
						"path":    {"dir", "directory"},
						"pattern": {"glob", "expr"},
					},
				},
				"grep": {
					NameAliases: []string{"search"},
					ArgumentsAliases: map[string][]string{
						"path":    {"dir", "directory", "file_path"},
						"pattern": {"query", "regex"},
					},
				},
				"execute": {
					NameAliases: []string{"bash", "sh", "shell", "cmd", "run", "terminal", "assistant", "exec"},
					ArgumentsAliases: map[string][]string{
						"command": {"cmd", "script", "code", "input"},
					},
				},
			},
			ToolArgumentsHandler: func(ctx context.Context, name, arguments string) (string, error) {
				var m map[string]any
				if err := json.Unmarshal([]byte(arguments), &m); err != nil {
					return arguments, nil
				}
				changed := false
				if name == "read_file" || name == "write_file" || name == "edit_file" {
					for _, key := range []string{"path", "filepath", "file", "filename", "target"} {
						if val, ok := m[key].(string); ok && val != "" {
							if _, exists := m["file_path"]; !exists {
								m["file_path"] = val
								changed = true
							}
						}
					}
				}
				if name == "write_file" {
					for _, key := range []string{"text", "data", "contents", "code", "body"} {
						if val, ok := m[key].(string); ok && val != "" {
							if _, exists := m["content"]; !exists {
								m["content"] = val
								changed = true
							}
						}
					}
					if _, hasPath := m["file_path"]; !hasPath || m["file_path"] == "" {
						if content, ok := m["content"].(string); ok {
							if strings.Contains(content, "[package]") {
								m["file_path"] = "experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-RefactoringKata/rust/Cargo.toml"
								changed = true
							} else if strings.Contains(content, "struct Item") || strings.Contains(content, "GildedRose") {
								m["file_path"] = "experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-RefactoringKata/rust/src/gildedrose.rs"
								changed = true
							} else if strings.Contains(content, "fn main") {
								m["file_path"] = "experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-RefactoringKata/rust/src/main.rs"
								changed = true
							}
						}
					}
				}
				if name == "glob" || name == "grep" || name == "ls" {
					if p, ok := m["path"].(string); !ok || p == "" || p == "/" {
						m["path"] = "."
						changed = true
					}
					if name == "glob" {
						if patterns, ok := m["patterns"].([]any); ok && len(patterns) > 0 {
							if firstPat, ok := patterns[0].(string); ok {
								m["pattern"] = firstPat
								changed = true
							}
						}
					}
				}
				if name == "execute" {
					for _, key := range []string{"cmd", "script", "code", "input"} {
						if val, ok := m[key].(string); ok && val != "" {
							if _, exists := m["command"]; !exists {
								m["command"] = val
								changed = true
							}
						}
					}
				}
				if changed {
					b, _ := json.Marshal(m)
					return string(b), nil
				}
				return arguments, nil
			},
			UnknownToolsHandler: func(ctx context.Context, name, input string) (string, error) {
				return fmt.Sprintf("Tool %s not found. Available tools: read_file, write_file, edit_file, glob, grep, execute", name), nil
			},
			ToolCallMiddlewares: []compose.ToolMiddleware{
				{
					Invokable: func(next compose.InvokableToolEndpoint) compose.InvokableToolEndpoint {
						return func(ctx context.Context, input *compose.ToolInput) (*compose.ToolOutput, error) {
							out, err := next(ctx, input)
							if err != nil {
								return &compose.ToolOutput{
									Result: fmt.Sprintf("Error executing tool %s: %v", input.Name, err),
								}, nil
							}
							return out, nil
						}
					},
				},
			},
		},
	}

	retryConfig := &adk.ModelRetryConfig{
		MaxRetries: 5,
		ShouldRetry: func(ctx context.Context, retryCtx *adk.RetryContext) *adk.RetryDecision {
			if retryCtx.Err != nil {
				return &adk.RetryDecision{Retry: true}
			}
			if retryCtx.OutputMessage == nil {
				return &adk.RetryDecision{Retry: false}
			}
			// If the model produced tool calls, let them execute!
			if len(retryCtx.OutputMessage.ToolCalls) > 0 {
				return &adk.RetryDecision{Retry: false}
			}

			// The model produced a text response without tool calls.
			// Check if all files exist
			reqFiles := []string{
				"experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-RefactoringKata/rust/Cargo.toml",
				"experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-RefactoringKata/rust/src/gildedrose.rs",
				"experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-RefactoringKata/rust/src/main.rs",
			}
			var missing []string
			for _, f := range reqFiles {
				if _, err := os.Stat(f); os.IsNotExist(err) {
					missing = append(missing, f)
				}
			}
			if len(missing) > 0 {
				return &adk.RetryDecision{
					Retry:                        true,
					ModifiedInputMessages:        append(retryCtx.InputMessages, schema.UserMessage(fmt.Sprintf("Do not stop or ask questions. Call write_file immediately to create the missing files: %v, then run cargo test.", missing))),
					PersistModifiedInputMessages: true,
				}
			}

			return &adk.RetryDecision{Retry: false}
		},
	}
	codingAgent := Must(adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:             "coding-agent",
		Description:      "Expert coding agent that writes code, creates files, and runs builds/tests using filesystem tools.",
		Instruction:      systemPrompt,
		Model:            codingModel,
		ToolsConfig:      toolsConfig,
		Handlers:         []adk.ChatModelAgentMiddleware{fsMiddleware},
		ModelRetryConfig: retryConfig,
	}))
	agent := Must(deep.New(ctx, &deep.Config{
		Name:             "reasoning-deep-agent",
		ChatModel:        reasoningModel,
		Instruction:      systemPrompt,
		SubAgents:        []adk.Agent{codingAgent},
		Backend:          backend,
		Shell:            backend,
		ToolsConfig:      toolsConfig,
		ModelRetryConfig: retryConfig,
	}))

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	events := runner.Query(ctx, taskPrompt)
	print.Events(events)
}
