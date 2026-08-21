package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	. "MAgHARCM/internal/patterns"
)

// DefinitionInput parameters for LSP definition lookup.
type DefinitionInput struct {
	Symbol   string `json:"symbol" jsonschema_description:"Name of the function, struct, class, type, or variable to find definition for"`
	FilePath string `json:"file_path,omitempty" jsonschema_description:"Optional file path to narrow the search scope"`
	ScopeDir string `json:"scope_dir,omitempty" jsonschema_description:"Base directory to search in (default: .)"`
}

// DefinitionOutput result of LSP definition lookup.
type DefinitionOutput struct {
	Symbol         string `json:"symbol"`
	FilePath       string `json:"file_path"`
	StartLine      int    `json:"start_line"`
	EndLine        int    `json:"end_line"`
	Implementation string `json:"implementation"`
	Found          bool   `json:"found"`
}

// DiagnosticsInput parameters for LSP diagnostics.
type DiagnosticsInput struct {
	FilePath   string `json:"file_path,omitempty" jsonschema_description:"Source file to check for diagnostics"`
	ProjectDir string `json:"project_dir,omitempty" jsonschema_description:"Project directory to run language diagnostics on"`
}

// DiagnosticIssue represents a single compiler/LSP error or warning.
type DiagnosticIssue struct {
	FilePath string `json:"file_path"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"` // "error", "warning", "info"
	Message  string `json:"message"`
}

// DiagnosticsOutput result of LSP diagnostics.
type DiagnosticsOutput struct {
	HasErrors bool              `json:"has_errors"`
	Issues    []DiagnosticIssue `json:"issues"`
	RawOutput string            `json:"raw_output,omitempty"`
}

// AtomicTextEdit a single replacement operation in a file.
type AtomicTextEdit struct {
	OldText string `json:"old_text" jsonschema_description:"Exact text snippet to replace"`
	NewText string `json:"new_text" jsonschema_description:"Replacement text"`
}

// AtomicEditInput parameters for applying multiple text edits to a file atomically.
type AtomicEditInput struct {
	FilePath string           `json:"file_path" jsonschema_description:"Path of the file to edit"`
	Edits    []AtomicTextEdit `json:"edits" jsonschema_description:"List of atomic replacements to apply"`
}

// AtomicEditOutput result of atomic edits.
type AtomicEditOutput struct {
	FilePath     string `json:"file_path"`
	Success      bool   `json:"success"`
	AppliedCount int    `json:"applied_count"`
	Message      string `json:"message"`
}

// HoverInput parameters for LSP hover documentation.
type HoverInput struct {
	Symbol   string `json:"symbol" jsonschema_description:"Symbol name to look up documentation for"`
	FilePath string `json:"file_path,omitempty" jsonschema_description:"File containing the symbol"`
	ScopeDir string `json:"scope_dir,omitempty" jsonschema_description:"Directory scope to search in"`
}

// HoverOutput result of LSP hover documentation.
type HoverOutput struct {
	Symbol    string `json:"symbol"`
	Docstring string `json:"docstring"`
	Signature string `json:"signature"`
	Found     bool   `json:"found"`
}

// ReferencesInput parameters for LSP symbol references lookup.
type ReferencesInput struct {
	Symbol   string `json:"symbol" jsonschema_description:"Symbol name to locate references for"`
	ScopeDir string `json:"scope_dir,omitempty" jsonschema_description:"Directory scope to search in (default: .)"`
}

// ReferenceItem single occurrence of a symbol.
type ReferenceItem struct {
	FilePath string `json:"file_path"`
	Line     int    `json:"line"`
	Content  string `json:"content"`
}

// ReferencesOutput result of LSP references lookup.
type ReferencesOutput struct {
	Symbol     string          `json:"symbol"`
	References []ReferenceItem `json:"references"`
	Count      int             `json:"count"`
}

// RenameSymbolInput parameters for LSP symbol rename.
type RenameSymbolInput struct {
	Symbol   string `json:"symbol" jsonschema_description:"Old symbol name to rename"`
	NewName  string `json:"new_name" jsonschema_description:"New name for the symbol"`
	ScopeDir string `json:"scope_dir,omitempty" jsonschema_description:"Directory scope for renaming (default: .)"`
}

// RenameSymbolOutput result of LSP symbol rename.
type RenameSymbolOutput struct {
	OldSymbol     string   `json:"old_symbol"`
	NewSymbol     string   `json:"new_symbol"`
	RenamedCount  int      `json:"renamed_count"`
	ModifiedFiles []string `json:"modified_files"`
	Success       bool     `json:"success"`
	Message       string   `json:"message"`
}

// NewLSPTools constructs the Language Server Protocol (LSP) tool group (ReCodeAgent 3.1.1).
func NewLSPTools() []tool.BaseTool {
	definitionTool := Must(utils.InferTool(
		"definition",
		"LSP Definition Tool: Retrieves the complete implementation, file path, and line numbers of a symbol (function, struct, class, type)",
		func(ctx context.Context, input *DefinitionInput) (*DefinitionOutput, error) {
			if input.Symbol == "" {
				return nil, fmt.Errorf("symbol is required for definition tool")
			}
			scope := input.ScopeDir
			if scope == "" {
				scope = "."
			}
			symbol := input.Symbol

			// Search pattern for definitions across C, C++, Rust, Go
			patterns := []*regexp.Regexp{
				regexp.MustCompile(fmt.Sprintf(`(?m)^(?:pub\s+)?(?:fn|struct|enum|trait|type|impl)\s+(?:<[^>]+>\s+)?%s\b`, regexp.QuoteMeta(symbol))),
				regexp.MustCompile(fmt.Sprintf(`(?m)^\s*(?:typedef\s+)?struct\s+%s\b`, regexp.QuoteMeta(symbol))),
				regexp.MustCompile(fmt.Sprintf(`(?m)^\s*(?:[a-zA-Z0-9_*]+\s+)+%s\s*\(`, regexp.QuoteMeta(symbol))),
				regexp.MustCompile(fmt.Sprintf(`(?m)\bfunc\s+(?:\([^)]+\)\s+)?%s\b`, regexp.QuoteMeta(symbol))),
			}

			var targetFiles []string
			if input.FilePath != "" {
				targetFiles = []string{input.FilePath}
			} else {
				_ = filepath.WalkDir(scope, func(p string, d os.DirEntry, err error) error {
					if err != nil || d.IsDir() {
						return nil
					}
					ext := strings.ToLower(filepath.Ext(p))
					if ext == ".rs" || ext == ".c" || ext == ".h" || ext == ".cc" || ext == ".cpp" || ext == ".go" {
						targetFiles = append(targetFiles, p)
					}
					return nil
				})
			}

			for _, file := range targetFiles {
				data, err := os.ReadFile(file)
				if err != nil {
					continue
				}
				content := string(data)
				lines := strings.Split(content, "\n")

				for _, pat := range patterns {
					loc := pat.FindStringIndex(content)
					if loc != nil {
						// Calculate start line
						startLine := strings.Count(content[:loc[0]], "\n") + 1
						endLine := startLine
						var blockLines []string
						braceCount := 0
						seenBrace := false

						for i := startLine - 1; i < len(lines); i++ {
							line := lines[i]
							blockLines = append(blockLines, line)
							braceCount += strings.Count(line, "{") - strings.Count(line, "}")
							if strings.Contains(line, "{") {
								seenBrace = true
							}
							if seenBrace && braceCount <= 0 {
								endLine = i + 1
								break
							}
							// For declarations ending in semicolon (like in .h files)
							if !seenBrace && strings.Contains(line, ";") {
								endLine = i + 1
								break
							}
						}

						return &DefinitionOutput{
							Symbol:         symbol,
							FilePath:       filepath.ToSlash(file),
							StartLine:      startLine,
							EndLine:        endLine,
							Implementation: strings.Join(blockLines, "\n"),
							Found:          true,
						}, nil
					}
				}
			}

			return &DefinitionOutput{
				Symbol: symbol,
				Found:  false,
			}, nil
		},
	))

	diagnosticsTool := Must(utils.InferTool(
		"diagnostics",
		"LSP Diagnostics Tool: Provides diagnostic information (errors and warnings) for a specified file or project",
		func(ctx context.Context, input *DiagnosticsInput) (*DiagnosticsOutput, error) {
			projDir := input.ProjectDir
			if projDir == "" && input.FilePath != "" {
				projDir = filepath.Dir(input.FilePath)
			}
			if projDir == "" {
				projDir = "."
			}

			// Check for Cargo.toml (Rust)
			cargoToml := filepath.Join(projDir, "Cargo.toml")
			if _, err := os.Stat(cargoToml); err == nil {
				cmd := exec.CommandContext(ctx, "cargo", "check", "--manifest-path", cargoToml, "--message-format=short")
				out, err := cmd.CombinedOutput()
				outStr := string(out)
				hasErrors := err != nil
				var issues []DiagnosticIssue

				scanner := bufio.NewScanner(strings.NewReader(outStr))
				re := regexp.MustCompile(`^(.*?):(\d+):(\d+):\s*(error|warning):\s*(.*)$`)
				for scanner.Scan() {
					line := scanner.Text()
					matches := re.FindStringSubmatch(line)
					if len(matches) == 6 {
						var lineNum, colNum int
						fmt.Sscanf(matches[2], "%d", &lineNum)
						fmt.Sscanf(matches[3], "%d", &colNum)
						issues = append(issues, DiagnosticIssue{
							FilePath: matches[1],
							Line:     lineNum,
							Column:   colNum,
							Severity: matches[4],
							Message:  matches[5],
						})
					}
				}
				return &DiagnosticsOutput{
					HasErrors: hasErrors,
					Issues:    issues,
					RawOutput: outStr,
				}, nil
			}

			// Check for C/C++ files
			if input.FilePath != "" && (strings.HasSuffix(input.FilePath, ".c") || strings.HasSuffix(input.FilePath, ".h")) {
				cmd := exec.CommandContext(ctx, "gcc", "-fsyntax-only", input.FilePath)
				out, err := cmd.CombinedOutput()
				outStr := string(out)
				return &DiagnosticsOutput{
					HasErrors: err != nil,
					RawOutput: outStr,
				}, nil
			}

			return &DiagnosticsOutput{
				HasErrors: false,
				RawOutput: "No language server diagnostic issues found.",
			}, nil
		},
	))

	atomicEditTool := Must(utils.InferTool(
		"edit_file_atomic",
		"LSP Edit File Tool: Applies a batch of text edits to a file atomically",
		func(ctx context.Context, input *AtomicEditInput) (*AtomicEditOutput, error) {
			path := filepath.Clean(input.FilePath)
			if path == "" {
				return nil, fmt.Errorf("file_path is required")
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", path, err)
			}
			content := string(data)
			applied := 0

			for i, edit := range input.Edits {
				if !strings.Contains(content, edit.OldText) {
					return &AtomicEditOutput{
						FilePath:     path,
						Success:      false,
						AppliedCount: applied,
						Message:      fmt.Sprintf("Edit %d failed: old_text not found in %s", i+1, path),
					}, nil
				}
				content = strings.Replace(content, edit.OldText, edit.NewText, 1)
				applied++
			}

			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				return nil, fmt.Errorf("failed to write updated file %s: %w", path, err)
			}

			return &AtomicEditOutput{
				FilePath:     path,
				Success:      true,
				AppliedCount: applied,
				Message:      fmt.Sprintf("Atomically applied %d edits to %s", applied, path),
			}, nil
		},
	))

	hoverTool := Must(utils.InferTool(
		"hover",
		"LSP Hover Tool: Returns documentation, docstrings, and signature for a symbol",
		func(ctx context.Context, input *HoverInput) (*HoverOutput, error) {
			if input.Symbol == "" {
				return nil, fmt.Errorf("symbol is required for hover")
			}
			scope := input.ScopeDir
			if scope == "" {
				scope = "."
			}
			symbol := input.Symbol

			var targetFiles []string
			if input.FilePath != "" {
				targetFiles = []string{input.FilePath}
			} else {
				_ = filepath.WalkDir(scope, func(p string, d os.DirEntry, err error) error {
					if err != nil || d.IsDir() {
						return nil
					}
					ext := strings.ToLower(filepath.Ext(p))
					if ext == ".rs" || ext == ".c" || ext == ".h" || ext == ".cc" || ext == ".go" {
						targetFiles = append(targetFiles, p)
					}
					return nil
				})
			}

			re := regexp.MustCompile(fmt.Sprintf(`(?m)((?:(?:///|//|/\*).*?\n)*)^\s*(?:pub\s+)?(?:fn|struct|enum|trait|type|impl|typedef)?.*?\b%s\b[^{;]*`, regexp.QuoteMeta(symbol)))
			for _, file := range targetFiles {
				data, err := os.ReadFile(file)
				if err != nil {
					continue
				}
				matches := re.FindStringSubmatch(string(data))
				if len(matches) > 0 {
					doc := strings.TrimSpace(matches[1])
					sig := strings.TrimSpace(strings.TrimPrefix(matches[0], matches[1]))
					return &HoverOutput{
						Symbol:    symbol,
						Docstring: doc,
						Signature: sig,
						Found:     true,
					}, nil
				}
			}

			return &HoverOutput{
				Symbol: symbol,
				Found:  false,
			}, nil
		},
	))

	referencesTool := Must(utils.InferTool(
		"references",
		"LSP References Tool: Locates all occurrences and usages of a symbol across the codebase",
		func(ctx context.Context, input *ReferencesInput) (*ReferencesOutput, error) {
			if input.Symbol == "" {
				return nil, fmt.Errorf("symbol is required for references")
			}
			scope := input.ScopeDir
			if scope == "" {
				scope = "."
			}
			symbol := input.Symbol
			re := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(symbol)))

			var refs []ReferenceItem
			_ = filepath.WalkDir(scope, func(p string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				ext := strings.ToLower(filepath.Ext(p))
				if ext == ".rs" || ext == ".c" || ext == ".h" || ext == ".cc" || ext == ".cpp" || ext == ".go" || ext == ".toml" {
					file, err := os.Open(p)
					if err != nil {
						return nil
					}
					defer file.Close()

					scanner := bufio.NewScanner(file)
					lineNo := 1
					for scanner.Scan() {
						text := scanner.Text()
						if re.MatchString(text) {
							refs = append(refs, ReferenceItem{
								FilePath: filepath.ToSlash(p),
								Line:     lineNo,
								Content:  text,
							})
						}
						lineNo++
					}
				}
				return nil
			})

			return &ReferencesOutput{
				Symbol:     symbol,
				References: refs,
				Count:      len(refs),
			}, nil
		},
	))

	renameSymbolTool := Must(utils.InferTool(
		"rename_symbol",
		"LSP Rename Symbol Tool: Renames a symbol and updates all corresponding references across the project",
		func(ctx context.Context, input *RenameSymbolInput) (*RenameSymbolOutput, error) {
			if input.Symbol == "" || input.NewName == "" {
				return nil, fmt.Errorf("symbol and new_name are required")
			}
			scope := input.ScopeDir
			if scope == "" {
				scope = "."
			}
			symbol := input.Symbol
			newName := input.NewName
			re := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(symbol)))

			totalRenamed := 0
			var modified []string

			_ = filepath.WalkDir(scope, func(p string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				ext := strings.ToLower(filepath.Ext(p))
				if ext == ".rs" || ext == ".c" || ext == ".h" || ext == ".cc" || ext == ".cpp" || ext == ".go" {
					data, err := os.ReadFile(p)
					if err != nil {
						return nil
					}
					content := string(data)
					if re.MatchString(content) {
						count := len(re.FindAllString(content, -1))
						newContent := re.ReplaceAllString(content, newName)
						if err := os.WriteFile(p, []byte(newContent), 0644); err == nil {
							totalRenamed += count
							modified = append(modified, filepath.ToSlash(p))
						}
					}
				}
				return nil
			})

			return &RenameSymbolOutput{
				OldSymbol:     symbol,
				NewSymbol:     newName,
				RenamedCount:  totalRenamed,
				ModifiedFiles: modified,
				Success:       true,
				Message:       fmt.Sprintf("Renamed %d occurrences across %d files", totalRenamed, len(modified)),
			}, nil
		},
	))

	return []tool.BaseTool{definitionTool, diagnosticsTool, atomicEditTool, hoverTool, referencesTool, renameSymbolTool}
}
