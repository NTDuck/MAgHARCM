package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// AllTools returns the full suite of tools across all groups (fs, lsp, pa, execution, git, validation).
func AllTools(backend filesystem.Backend, shell filesystem.Shell) []tool.BaseTool {
	var list []tool.BaseTool
	list = append(list, NewFSTools(backend)...)
	list = append(list, NewLSPTools()...)
	list = append(list, NewPATools()...)
	list = append(list, NewExecutionTools(shell)...)
	list = append(list, NewGitTools()...)
	list = append(list, NewValidationTools()...)
	return list
}

// NewToolsConfig creates a centralized ToolsConfig for all tool groups,
// including comprehensive aliases, argument normalization, and error recovery middleware.
func NewToolsConfig(tools ...tool.BaseTool) adk.ToolsConfig {
	return adk.ToolsConfig{
		ToolsNodeConfig: compose.ToolsNodeConfig{
			Tools:               tools,
			ExecuteSequentially: true,
			ToolAliases: map[string]compose.ToolAliasConfig{
				"read_file": {
					NameAliases: []string{"read", "cat", "view_file", "readFile"},
					ArgumentsAliases: map[string][]string{
						"file_path": {"path", "filepath", "file", "filename", "target"},
					},
				},
				"write_file": {
					NameAliases: []string{"write", "create_file", "save_file", "writeFile"},
					ArgumentsAliases: map[string][]string{
						"file_path": {"path", "filepath", "file", "filename", "target"},
						"content":   {"text", "data", "contents", "code", "body"},
					},
				},
				"edit_file": {
					NameAliases: []string{"edit", "replace_in_file", "patch_file", "editFile"},
					ArgumentsAliases: map[string][]string{
						"file_path": {"path", "filepath", "file", "filename", "target"},
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
				"definition": {
					NameAliases: []string{"goto_definition", "find_definition", "get_definition"},
					ArgumentsAliases: map[string][]string{
						"symbol":    {"name", "identifier"},
						"file_path": {"path", "file"},
					},
				},
				"diagnostics": {
					NameAliases: []string{"check_diagnostics", "get_diagnostics", "lint"},
					ArgumentsAliases: map[string][]string{
						"file_path":   {"path", "file"},
						"project_dir": {"dir", "project"},
					},
				},
				"get_directory_tree": {
					NameAliases: []string{"tree", "directory_tree", "project_tree"},
					ArgumentsAliases: map[string][]string{
						"path": {"dir", "root"},
					},
				},
				"get_file_structure": {
					NameAliases: []string{"file_structure", "outline", "symbols"},
					ArgumentsAliases: map[string][]string{
						"file_path": {"path", "file"},
					},
				},
			},
			ToolArgumentsHandler: func(ctx context.Context, name, arguments string) (string, error) {
				var m map[string]any
				if err := json.Unmarshal([]byte(arguments), &m); err != nil {
					return arguments, nil
				}
				changed := false

				// Generic file path normalization
				for _, key := range []string{"path", "filepath", "file", "filename", "target"} {
					if val, ok := m[key].(string); ok && val != "" {
						if _, exists := m["file_path"]; !exists {
							m["file_path"] = val
							changed = true
						}
					}
				}

				// Generic content normalization
				for _, key := range []string{"text", "data", "contents", "code", "body"} {
					if val, ok := m[key].(string); ok && val != "" {
						if _, exists := m["content"]; !exists {
							m["content"] = val
							changed = true
						}
					}
				}

				// Path defaulting for directory / glob / tree tools
				if name == "glob" || name == "grep" || name == "list_dir" || name == "get_directory_tree" {
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

				// Execute command normalization
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
				return fmt.Sprintf("Tool %s not found. Please use available tools (fs, lsp, pa, execution, git, validation)", name), nil
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
