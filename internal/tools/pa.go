package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_c "github.com/tree-sitter/tree-sitter-c/bindings/go"
	tree_sitter_cpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
	tree_sitter_rust "github.com/tree-sitter/tree-sitter-rust/bindings/go"

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

// CodeElement key structural element extracted from a source file via tree-sitter AST.
type CodeElement struct {
	Type      string `json:"type"` // "function", "struct", "enum", "trait", "impl", "type", "variable", "include"
	Name      string `json:"name"`
	StartLine uint   `json:"start_line"`
	EndLine   uint   `json:"end_line"`
	Signature string `json:"signature,omitempty"`
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

// getTreeSitterLanguage returns the tree-sitter language for a given file extension.
func getTreeSitterLanguage(ext string) (*tree_sitter.Language, string) {
	switch strings.ToLower(ext) {
	case ".rs":
		return tree_sitter.NewLanguage(tree_sitter_rust.Language()), "rust"
	case ".c", ".h":
		return tree_sitter.NewLanguage(tree_sitter_c.Language()), "c"
	case ".cc", ".cpp", ".cxx", ".hpp":
		return tree_sitter.NewLanguage(tree_sitter_cpp.Language()), "cpp"
	case ".go":
		return tree_sitter.NewLanguage(tree_sitter_go.Language()), "go"
	default:
		return nil, "unknown"
	}
}

// extractElementsFromAST traverses tree-sitter AST to extract top-level code elements.
func extractElementsFromAST(root *tree_sitter.Node, src []byte, lang string) ([]CodeElement, []string) {
	var elements []CodeElement
	var imports []string

	childCount := root.NamedChildCount()
	for i := uint(0); i < childCount; i++ {
		child := root.NamedChild(i)
		if child == nil {
			continue
		}

		kind := child.Kind()
		startRow := child.StartPosition().Row + 1
		endRow := child.EndPosition().Row + 1

		switch lang {
		case "rust":
			switch kind {
			case "function_item":
				nameNode := child.ChildByFieldName("name")
				name := extractNodeText(nameNode, src)
				elements = append(elements, CodeElement{
					Type:      "function",
					Name:      name,
					StartLine: startRow,
					EndLine:   endRow,
					Signature: firstLineOf(child.Utf8Text(src)),
				})
			case "struct_item":
				nameNode := child.ChildByFieldName("name")
				name := extractNodeText(nameNode, src)
				elements = append(elements, CodeElement{
					Type:      "struct",
					Name:      name,
					StartLine: startRow,
					EndLine:   endRow,
					Signature: firstLineOf(child.Utf8Text(src)),
				})
			case "enum_item":
				nameNode := child.ChildByFieldName("name")
				name := extractNodeText(nameNode, src)
				elements = append(elements, CodeElement{
					Type:      "enum",
					Name:      name,
					StartLine: startRow,
					EndLine:   endRow,
					Signature: firstLineOf(child.Utf8Text(src)),
				})
			case "trait_item":
				nameNode := child.ChildByFieldName("name")
				name := extractNodeText(nameNode, src)
				elements = append(elements, CodeElement{
					Type:      "trait",
					Name:      name,
					StartLine: startRow,
					EndLine:   endRow,
					Signature: firstLineOf(child.Utf8Text(src)),
				})
			case "impl_item":
				typeNode := child.ChildByFieldName("type")
				traitNode := child.ChildByFieldName("trait")
				name := extractNodeText(typeNode, src)
				if traitNode != nil {
					name = extractNodeText(traitNode, src) + " for " + name
				}
				elements = append(elements, CodeElement{
					Type:      "impl",
					Name:      name,
					StartLine: startRow,
					EndLine:   endRow,
					Signature: firstLineOf(child.Utf8Text(src)),
				})
			case "use_declaration":
				imports = append(imports, strings.TrimSpace(child.Utf8Text(src)))
			}
		case "c", "cpp":
			switch kind {
			case "function_definition":
				decl := child.ChildByFieldName("declarator")
				name := extractDeclaratorName(decl, src)
				elements = append(elements, CodeElement{
					Type:      "function",
					Name:      name,
					StartLine: startRow,
					EndLine:   endRow,
					Signature: firstLineOf(child.Utf8Text(src)),
				})
			case "declaration":
				text := strings.TrimSpace(child.Utf8Text(src))
				if strings.Contains(text, "typedef") || strings.Contains(text, "struct") {
					elements = append(elements, CodeElement{
						Type:      "struct",
						Name:      firstLineOf(text),
						StartLine: startRow,
						EndLine:   endRow,
						Signature: firstLineOf(text),
					})
				} else {
					elements = append(elements, CodeElement{
						Type:      "declaration",
						Name:      firstLineOf(text),
						StartLine: startRow,
						EndLine:   endRow,
						Signature: firstLineOf(text),
					})
				}
			case "preproc_include":
				imports = append(imports, strings.TrimSpace(child.Utf8Text(src)))
			}
		case "go":
			switch kind {
			case "function_declaration", "method_declaration":
				nameNode := child.ChildByFieldName("name")
				name := extractNodeText(nameNode, src)
				elements = append(elements, CodeElement{
					Type:      "function",
					Name:      name,
					StartLine: startRow,
					EndLine:   endRow,
					Signature: firstLineOf(child.Utf8Text(src)),
				})
			case "type_declaration":
				text := strings.TrimSpace(child.Utf8Text(src))
				elements = append(elements, CodeElement{
					Type:      "type",
					Name:      firstLineOf(text),
					StartLine: startRow,
					EndLine:   endRow,
					Signature: firstLineOf(text),
				})
			case "import_declaration":
				imports = append(imports, strings.TrimSpace(child.Utf8Text(src)))
			}
		}
	}

	return elements, imports
}

func extractNodeText(n *tree_sitter.Node, src []byte) string {
	if n == nil {
		return ""
	}
	return strings.TrimSpace(n.Utf8Text(src))
}

func extractDeclaratorName(n *tree_sitter.Node, src []byte) string {
	if n == nil {
		return ""
	}
	if n.Kind() == "identifier" {
		return strings.TrimSpace(n.Utf8Text(src))
	}
	if n.Kind() == "function_declarator" {
		direct := n.ChildByFieldName("declarator")
		if direct != nil {
			return extractDeclaratorName(direct, src)
		}
	}
	for i := uint(0); i < n.NamedChildCount(); i++ {
		child := n.NamedChild(i)
		if child != nil && child.Kind() == "identifier" {
			return strings.TrimSpace(child.Utf8Text(src))
		}
	}
	return strings.TrimSpace(n.Utf8Text(src))
}

func firstLineOf(s string) string {
	idx := strings.Index(s, "\n")
	if idx >= 0 {
		return strings.TrimSpace(s[:idx])
	}
	return strings.TrimSpace(s)
}

// NewPATools constructs Project Analysis tools (ReCodeAgent 3.1.2) using tree-sitter AST parsing.
func NewPATools() []tool.BaseTool {
	directoryTreeTool := Must(utils.InferTool(
		"get_directory_tree",
		"Project Analysis: Returns a structured ASCII directory tree representation of the project hierarchy",
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
			fileCount, dirCount := 0, 0

			var walk func(dir string, prefix string, depth int) error
			walk = func(dir string, prefix string, depth int) error {
				if depth > maxDepth {
					return nil
				}
				entries, err := os.ReadDir(dir)
				if err != nil {
					return nil
				}

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
		"Project Analysis: Generates tree-sitter AST structured representation of source files (C, C++, Rust, Go), identifying functions, structs, enums, traits, and imports",
		func(ctx context.Context, input *FileStructureInput) (*FileStructureOutput, error) {
			path := filepath.Clean(input.FilePath)
			if path == "" {
				return nil, fmt.Errorf("file_path is required")
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", path, err)
			}

			ext := filepath.Ext(path)
			langObj, langName := getTreeSitterLanguage(ext)
			if langObj == nil {
				return &FileStructureOutput{
					FilePath:  path,
					Language:  langName,
					LineCount: len(strings.Split(string(data), "\n")),
				}, nil
			}

			parser := tree_sitter.NewParser()
			defer parser.Close()
			_ = parser.SetLanguage(langObj)

			tree := parser.Parse(data, nil)
			if tree == nil {
				return nil, fmt.Errorf("failed to parse %s with tree-sitter", path)
			}
			defer tree.Close()

			root := tree.RootNode()
			elements, imports := extractElementsFromAST(root, data, langName)

			return &FileStructureOutput{
				FilePath:  path,
				Language:  langName,
				Elements:  elements,
				Imports:   imports,
				LineCount: len(strings.Split(string(data), "\n")),
			}, nil
		},
	))

	return []tool.BaseTool{directoryTreeTool, fileStructureTool}
}
