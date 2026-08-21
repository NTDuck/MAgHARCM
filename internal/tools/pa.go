package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	. "MAgHARCM/internal/patterns"
)

// DirectoryTreeInput parameters for project directory tree analysis.
type DirectoryTreeInput struct {
	Path          string `json:"path,omitempty" jsonschema_description:"Root directory path to analyze (default: .)"`
	MaxDepth      int    `json:"max_depth,omitempty" jsonschema_description:"Maximum directory traversal depth (default: 4)"`
	IncludeHidden bool   `json:"include_hidden,omitempty" jsonschema_description:"Whether to include dotfiles (default: false)"`
}

// DirectoryTreeOutput result of directory tree analysis.
type DirectoryTreeOutput struct {
	Path      string `json:"path"`
	Tree      string `json:"tree"`
	FileCount int    `json:"file_count"`
	DirCount  int    `json:"dir_count"`
}

// CodeElement key structural element extracted from a source file.
type CodeElement struct {
	Type      string `json:"type"` // "function", "struct", "class", "enum", "trait", "impl", "variable", "macro"
	Name      string `json:"name"`
	Line      int    `json:"line"`
	Signature string `json:"signature"`
}

// FileStructureInput parameters for file structure analysis.
type FileStructureInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path of the source file to analyze"`
}

// FileStructureOutput result of file structure analysis.
type FileStructureOutput struct {
	FilePath  string        `json:"file_path"`
	Language  string        `json:"language"`
	Elements  []CodeElement `json:"elements"`
	Imports   []string      `json:"imports"`
	LineCount int           `json:"line_count"`
}

// NewPATools constructs the Project Analysis (PA) tool group (ReCodeAgent 3.1.2).
func NewPATools() []tool.BaseTool {
	directoryTreeTool := Must(utils.InferTool(
		"get_directory_tree",
		"Project Analysis: Returns a structured ASCII directory tree representation of the project",
		func(ctx context.Context, input *DirectoryTreeInput) (*DirectoryTreeOutput, error) {
			basePath := input.Path
			if basePath == "" {
				basePath = "."
			}
			basePath = filepath.Clean(basePath)
			maxDepth := input.MaxDepth
			if maxDepth <= 0 {
				maxDepth = 4
			}

			var sb strings.Builder
			sb.WriteString(basePath + "\n")
			fileCount := 0
			dirCount := 0

			var walk func(dir string, prefix string, depth int) error
			walk = func(dir string, prefix string, depth int) error {
				if depth > maxDepth {
					return nil
				}
				entries, err := os.ReadDir(dir)
				if err != nil {
					return nil
				}

				// Filter out ignored / heavy directories
				var validEntries []os.DirEntry
				for _, e := range entries {
					name := e.Name()
					if !input.IncludeHidden && strings.HasPrefix(name, ".") {
						continue
					}
					if name == "target" || name == "node_modules" || name == "vendor" || name == ".git" {
						continue
					}
					validEntries = append(validEntries, e)
				}

				for i, e := range validEntries {
					isLast := i == len(validEntries)-1
					pointer := "├── "
					extension := "│   "
					if isLast {
						pointer = "└── "
						extension = "    "
					}

					name := e.Name()
					if e.IsDir() {
						dirCount++
						sb.WriteString(fmt.Sprintf("%s%s%s/\n", prefix, pointer, name))
						if err := walk(filepath.Join(dir, name), prefix+extension, depth+1); err != nil {
							return err
						}
					} else {
						fileCount++
						sb.WriteString(fmt.Sprintf("%s%s%s\n", prefix, pointer, name))
					}
				}
				return nil
			}

			if err := walk(basePath, "", 1); err != nil {
				return nil, fmt.Errorf("error generating directory tree: %w", err)
			}

			return &DirectoryTreeOutput{
				Path:      basePath,
				Tree:      sb.String(),
				FileCount: fileCount,
				DirCount:  dirCount,
			}, nil
		},
	))

	fileStructureTool := Must(utils.InferTool(
		"get_file_structure",
		"Project Analysis: Generates a structured representation of a source file, identifying functions, structs, classes, globals, and imports",
		func(ctx context.Context, input *FileStructureInput) (*FileStructureOutput, error) {
			path := filepath.Clean(input.FilePath)
			if path == "" {
				return nil, fmt.Errorf("file_path is required")
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", path, err)
			}

			ext := strings.ToLower(filepath.Ext(path))
			language := "unknown"
			switch ext {
			case ".rs":
				language = "rust"
			case ".c", ".h":
				language = "c"
			case ".cc", ".cpp", ".hpp":
				language = "cpp"
			case ".go":
				language = "go"
			case ".py":
				language = "python"
			}

			var elements []CodeElement
			var imports []string

			scanner := bufio.NewScanner(strings.NewReader(string(data)))
			lineNo := 1

			// Regex matchers for code elements across languages
			reRustFn := regexp.MustCompile(`^\s*(?:pub(?:\([^)]+\))?\s+)?(?:async\s+)?fn\s+([a-zA-Z0-9_]+)`)
			reRustStruct := regexp.MustCompile(`^\s*(?:pub(?:\([^)]+\))?\s+)?struct\s+([a-zA-Z0-9_]+)`)
			reRustEnum := regexp.MustCompile(`^\s*(?:pub(?:\([^)]+\))?\s+)?enum\s+([a-zA-Z0-9_]+)`)
			reRustTrait := regexp.MustCompile(`^\s*(?:pub(?:\([^)]+\))?\s+)?trait\s+([a-zA-Z0-9_]+)`)
			reRustImpl := regexp.MustCompile(`^\s*impl(?:\s*<[^>]+>)?\s+(?:([a-zA-Z0-9_]+)\s+for\s+)?([a-zA-Z0-9_]+)`)
			reRustUse := regexp.MustCompile(`^\s*use\s+([^;]+);`)

			reCFn := regexp.MustCompile(`^(?:[a-zA-Z0-9_*]+\s+)+([a-zA-Z0-9_]+)\s*\([^)]*\)\s*[{;]`)
			reCStruct := regexp.MustCompile(`^(?:typedef\s+)?struct\s*(?:([a-zA-Z0-9_]+))?\s*[{]?`)
			reCInclude := regexp.MustCompile(`^\s*#include\s+([<"][^>"]+[>"])`)

			reGoFunc := regexp.MustCompile(`^\s*func\s+(?:\([^)]+\)\s+)?([a-zA-Z0-9_]+)`)
			reGoType := regexp.MustCompile(`^\s*type\s+([a-zA-Z0-9_]+)\s+(struct|interface)`)

			for scanner.Scan() {
				line := scanner.Text()
				trimmed := strings.TrimSpace(line)

				switch language {
				case "rust":
					if m := reRustUse.FindStringSubmatch(trimmed); len(m) > 1 {
						imports = append(imports, m[1])
					} else if m := reRustFn.FindStringSubmatch(trimmed); len(m) > 1 {
						elements = append(elements, CodeElement{Type: "function", Name: m[1], Line: lineNo, Signature: trimmed})
					} else if m := reRustStruct.FindStringSubmatch(trimmed); len(m) > 1 {
						elements = append(elements, CodeElement{Type: "struct", Name: m[1], Line: lineNo, Signature: trimmed})
					} else if m := reRustEnum.FindStringSubmatch(trimmed); len(m) > 1 {
						elements = append(elements, CodeElement{Type: "enum", Name: m[1], Line: lineNo, Signature: trimmed})
					} else if m := reRustTrait.FindStringSubmatch(trimmed); len(m) > 1 {
						elements = append(elements, CodeElement{Type: "trait", Name: m[1], Line: lineNo, Signature: trimmed})
					} else if m := reRustImpl.FindStringSubmatch(trimmed); len(m) > 2 {
						target := m[2]
						if m[1] != "" {
							target = m[1] + " for " + m[2]
						}
						elements = append(elements, CodeElement{Type: "impl", Name: target, Line: lineNo, Signature: trimmed})
					}
				case "c", "cpp":
					if m := reCInclude.FindStringSubmatch(trimmed); len(m) > 1 {
						imports = append(imports, m[1])
					} else if m := reCStruct.FindStringSubmatch(trimmed); len(m) > 1 && m[1] != "" {
						elements = append(elements, CodeElement{Type: "struct", Name: m[1], Line: lineNo, Signature: trimmed})
					} else if m := reCFn.FindStringSubmatch(trimmed); len(m) > 1 && !strings.HasPrefix(trimmed, "if") && !strings.HasPrefix(trimmed, "for") && !strings.HasPrefix(trimmed, "while") {
						elements = append(elements, CodeElement{Type: "function", Name: m[1], Line: lineNo, Signature: trimmed})
					}
				case "go":
					if strings.HasPrefix(trimmed, "import") {
						imports = append(imports, trimmed)
					} else if m := reGoFunc.FindStringSubmatch(trimmed); len(m) > 1 {
						elements = append(elements, CodeElement{Type: "function", Name: m[1], Line: lineNo, Signature: trimmed})
					} else if m := reGoType.FindStringSubmatch(trimmed); len(m) > 2 {
						elements = append(elements, CodeElement{Type: m[2], Name: m[1], Line: lineNo, Signature: trimmed})
					}
				}

				lineNo++
			}

			return &FileStructureOutput{
				FilePath:  path,
				Language:  language,
				Elements:  elements,
				Imports:   imports,
				LineCount: lineNo - 1,
			}, nil
		},
	))

	return []tool.BaseTool{directoryTreeTool, fileStructureTool}
}
