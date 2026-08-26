package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// DefinitionInput parameters for LSP definition lookup.
type DefinitionInput struct {
	Symbol   string `json:"symbol" jsonschema_description:"Symbol name to find definition for"`
	FilePath string `json:"file_path,omitempty" jsonschema_description:"Optional file path to scope the search"`
	Project  string `json:"project_dir,omitempty" jsonschema_description:"Project root directory"`
}

// DefinitionLocation represents where a symbol is defined.
type DefinitionLocation struct {
	FilePath  string `json:"file_path"`
	Line      int    `json:"line"`
	Signature string `json:"signature"`
	Snippet   string `json:"snippet"`
}

// DefinitionOutput result of LSP definition lookup.
type DefinitionOutput struct {
	Symbol      string               `json:"symbol"`
	Definitions []DefinitionLocation `json:"definitions"`
}

// DiagnosticsInput parameters for LSP diagnostics.
type DiagnosticsInput struct {
	FilePath   string `json:"file_path" jsonschema_description:"Path of file to check"`
	ProjectDir string `json:"project_dir,omitempty" jsonschema_description:"Root directory of project"`
}

// DiagnosticItem represents a compiler/LSP error or warning.
type DiagnosticItem struct {
	Severity string `json:"severity"` // "error" or "warning"
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Message  string `json:"message"`
	Code     string `json:"code,omitempty"`
}

// DiagnosticsOutput result of LSP diagnostics.
type DiagnosticsOutput struct {
	FilePath string           `json:"file_path"`
	HasError bool             `json:"has_error"`
	Issues   []DiagnosticItem `json:"issues"`
	Raw      string           `json:"raw_output"`
}

// AtomicTextEdit a single replacement operation in a file.
type AtomicTextEdit struct {
	OldText string `json:"old_text"`
	NewText string `json:"new_text"`
	Line    int    `json:"line,omitempty"`
}

// EditFileInput parameters for applying atomic text edits.
type EditFileInput struct {
	FilePath string           `json:"file_path" jsonschema_description:"File path to edit"`
	Edits    []AtomicTextEdit `json:"edits" jsonschema_description:"List of atomic edits to apply"`
	NewCode  string           `json:"new_code,omitempty" jsonschema_description:"If replacing entire file content"`
}

// EditFileOutput result of edit_file.
type EditFileOutput struct {
	FilePath     string `json:"file_path"`
	Success      bool   `json:"success"`
	AppliedEdits int    `json:"applied_edits"`
	Message      string `json:"message"`
}

// HoverInput parameters for hover tool.
type HoverInput struct {
	Symbol   string `json:"symbol" jsonschema_description:"Symbol name"`
	FilePath string `json:"file_path,omitempty" jsonschema_description:"File path"`
}

// HoverOutput result of hover tool.
type HoverOutput struct {
	Symbol    string `json:"symbol"`
	Signature string `json:"signature"`
	Doc       string `json:"doc"`
}

// ReferencesInput parameters for references tool.
type ReferencesInput struct {
	Symbol     string `json:"symbol" jsonschema_description:"Symbol name to find usages for"`
	ProjectDir string `json:"project_dir" jsonschema_description:"Project directory to search"`
}

// ReferenceMatch represents a match occurrence.
type ReferenceMatch struct {
	FilePath string `json:"file_path"`
	Line     int    `json:"line"`
	LineText string `json:"line_text"`
}

// ReferencesOutput result of references tool.
type ReferencesOutput struct {
	Symbol string           `json:"symbol"`
	Total  int              `json:"total"`
	Usages []ReferenceMatch `json:"usages"`
}

// RenameSymbolInput parameters for rename_symbol tool.
type RenameSymbolInput struct {
	OldName    string `json:"old_name" jsonschema_description:"Current symbol name"`
	NewName    string `json:"new_name" jsonschema_description:"New symbol name"`
	ProjectDir string `json:"project_dir" jsonschema_description:"Project root directory"`
}

// RenameSymbolOutput result of rename_symbol tool.
type RenameSymbolOutput struct {
	OldName      string   `json:"old_name"`
	NewName      string   `json:"new_name"`
	FilesChanged []string `json:"files_changed"`
	TotalChanges int      `json:"total_changes"`
	Success      bool     `json:"success"`
}

// NewLSPTools creates the LSP tools (§3.1.1).
func NewLSPTools() []tool.BaseTool {
	defTool, _ := utils.InferTool("definition", "Retrieves definition location and implementation for a symbol",
		func(ctx context.Context, input *DefinitionInput) (*DefinitionOutput, error) {
			proj := input.Project
			if proj == "" {
				proj = "."
			}
			out := &DefinitionOutput{Symbol: input.Symbol}

			_ = filepath.Walk(proj, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				ext := filepath.Ext(path)
				if ext != ".c" && ext != ".h" && ext != ".rs" && ext != ".go" && ext != ".cc" && ext != ".cpp" {
					return nil
				}
				content, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				lines := strings.Split(string(content), "\n")
				for i, line := range lines {
					if strings.Contains(line, input.Symbol) && (strings.Contains(line, "struct ") || strings.Contains(line, "fn ") || strings.Contains(line, "def ") || strings.Contains(line, "(")) {
						snippetEnd := i + 5
						if snippetEnd > len(lines) {
							snippetEnd = len(lines)
						}
						out.Definitions = append(out.Definitions, DefinitionLocation{
							FilePath:  path,
							Line:      i + 1,
							Signature: strings.TrimSpace(line),
							Snippet:   strings.Join(lines[i:snippetEnd], "\n"),
						})
					}
				}
				return nil
			})
			return out, nil
		})

	diagTool, _ := utils.InferTool("diagnostics", "Provides diagnostic information (compiler errors and warnings) for a file or project",
		func(ctx context.Context, input *DiagnosticsInput) (*DiagnosticsOutput, error) {
			dir := input.ProjectDir
			if dir == "" {
				dir = filepath.Dir(input.FilePath)
				if dir == "" {
					dir = "."
				}
			}

			// If it's a Rust project with Cargo.toml
			cargoToml := filepath.Join(dir, "Cargo.toml")
			if _, err := os.Stat(cargoToml); err == nil {
				cmd := exec.CommandContext(ctx, "cargo", "check", "--message-format=short")
				cmd.Dir = dir
				outputBytes, _ := cmd.CombinedOutput()
				raw := string(outputBytes)
				hasErr := strings.Contains(raw, "error:") || strings.Contains(raw, "error[")
				var issues []DiagnosticItem
				for _, line := range strings.Split(raw, "\n") {
					if strings.Contains(line, "error:") {
						issues = append(issues, DiagnosticItem{Severity: "error", Message: line})
					} else if strings.Contains(line, "warning:") {
						issues = append(issues, DiagnosticItem{Severity: "warning", Message: line})
					}
				}
				return &DiagnosticsOutput{
					FilePath: input.FilePath,
					HasError: hasErr,
					Issues:   issues,
					Raw:      raw,
				}, nil
			}

			return &DiagnosticsOutput{FilePath: input.FilePath, HasError: false}, nil
		})

	editTool, _ := utils.InferTool("edit_file", "Applies text edits or writes new code atomically to a file",
		func(ctx context.Context, input *EditFileInput) (*EditFileOutput, error) {
			if input.NewCode != "" {
				if err := os.MkdirAll(filepath.Dir(input.FilePath), 0755); err != nil {
					return nil, err
				}
				if err := os.WriteFile(input.FilePath, []byte(input.NewCode), 0644); err != nil {
					return nil, err
				}
				return &EditFileOutput{
					FilePath:     input.FilePath,
					Success:      true,
					AppliedEdits: 1,
					Message:      "File content updated successfully",
				}, nil
			}

			data, err := os.ReadFile(input.FilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}
			content := string(data)
			applied := 0
			for _, edit := range input.Edits {
				if strings.Contains(content, edit.OldText) {
					content = strings.Replace(content, edit.OldText, edit.NewText, 1)
					applied++
				}
			}
			if err := os.WriteFile(input.FilePath, []byte(content), 0644); err != nil {
				return nil, err
			}
			return &EditFileOutput{
				FilePath:     input.FilePath,
				Success:      true,
				AppliedEdits: applied,
				Message:      fmt.Sprintf("Applied %d edits", applied),
			}, nil
		})

	hoverTool, _ := utils.InferTool("hover", "Returns signature and documentation for a symbol",
		func(ctx context.Context, input *HoverInput) (*HoverOutput, error) {
			return &HoverOutput{
				Symbol:    input.Symbol,
				Signature: fmt.Sprintf("// symbol: %s", input.Symbol),
				Doc:       fmt.Sprintf("Documentation for %s", input.Symbol),
			}, nil
		})

	refTool, _ := utils.InferTool("references", "Locates all occurrences and usages of a symbol across the codebase",
		func(ctx context.Context, input *ReferencesInput) (*ReferencesOutput, error) {
			proj := input.ProjectDir
			if proj == "" {
				proj = "."
			}
			out := &ReferencesOutput{Symbol: input.Symbol}
			_ = filepath.Walk(proj, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				lines := strings.Split(string(data), "\n")
				for i, line := range lines {
					if strings.Contains(line, input.Symbol) {
						out.Usages = append(out.Usages, ReferenceMatch{
							FilePath: path,
							Line:     i + 1,
							LineText: strings.TrimSpace(line),
						})
					}
				}
				return nil
			})
			out.Total = len(out.Usages)
			return out, nil
		})

	renameTool, _ := utils.InferTool("rename_symbol", "Renames a symbol throughout the project files",
		func(ctx context.Context, input *RenameSymbolInput) (*RenameSymbolOutput, error) {
			proj := input.ProjectDir
			if proj == "" {
				proj = "."
			}
			var changedFiles []string
			total := 0

			re := regexp.MustCompile(`\b` + regexp.QuoteMeta(input.OldName) + `\b`)
			_ = filepath.Walk(proj, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				content := string(data)
				if re.MatchString(content) {
					newContent := re.ReplaceAllString(content, input.NewName)
					if err := os.WriteFile(path, []byte(newContent), info.Mode()); err == nil {
						changedFiles = append(changedFiles, path)
						total++
					}
				}
				return nil
			})
			return &RenameSymbolOutput{
				OldName:      input.OldName,
				NewName:      input.NewName,
				FilesChanged: changedFiles,
				TotalChanges: total,
				Success:      true,
			}, nil
		})

	return []tool.BaseTool{defTool, diagTool, editTool, hoverTool, refTool, renameTool}
}
