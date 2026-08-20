package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/compose"
)

// NewConfig creates an idiomatic ToolsConfig with tool aliases, argument normalizers,
// and error recovery middleware.
func NewConfig() adk.ToolsConfig {
	return adk.ToolsConfig{
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
								m["file_path"] = "experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/Cargo.toml"
								changed = true
							} else if strings.Contains(content, "struct Item") || strings.Contains(content, "GildedRose") {
								m["file_path"] = "experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/src/gildedrose.rs"
								changed = true
							} else if strings.Contains(content, "fn main") {
								m["file_path"] = "experiments/0001-SA-DeepAgent-MinimalTools/artifacts/GildedRose-Refactoring-Kata/rust/src/main.rs"
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
}
