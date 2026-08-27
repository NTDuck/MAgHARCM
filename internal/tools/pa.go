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
)

// DirectoryTreeInput parameters for get_directory_tree.
type DirectoryTreeInput struct {
	DirectoryPath string `json:"directory_path" jsonschema_description:"Path of the directory to analyze (default: current directory)"`
	MaxDepth      int    `json:"max_depth,omitempty" jsonschema_description:"Maximum depth to traverse (default: 5)"`
}

// DirectoryTreeOutput result of get_directory_tree.
type DirectoryTreeOutput struct {
	DirectoryPath string   `json:"directory_path"`
	TreeString    string   `json:"tree_string"`
	Files         []string `json:"files"`
}

// CodeElement represents a structural item extracted by tree-sitter AST.
type CodeElement struct {
	Kind      string `json:"kind"` // "function", "struct", "class", "typedef", "global", "macro"
	Name      string `json:"name"`
	Signature string `json:"signature,omitempty"`
	Line      int    `json:"line"`
	EndLine   int    `json:"end_line"`
}

// FileStructureInput parameters for get_file_structure.
type FileStructureInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path of the source file to analyze"`
}

// FileStructureOutput result of get_file_structure.
type FileStructureOutput struct {
	FilePath string        `json:"file_path"`
	Language string        `json:"language"`
	Elements []CodeElement `json:"elements"`
	Imports  []string      `json:"imports"`
	RawCode  string        `json:"raw_code,omitempty"`
}

func getTreeSitterLanguage(ext string) (*tree_sitter.Language, string) {
	switch strings.ToLower(ext) {
	case ".c", ".h":
		return tree_sitter.NewLanguage(tree_sitter_c.Language()), "c"
	case ".cc", ".cpp", ".cxx", ".hpp":
		return tree_sitter.NewLanguage(tree_sitter_cpp.Language()), "cpp"
	case ".go":
		return tree_sitter.NewLanguage(tree_sitter_go.Language()), "go"
	case ".rs":
		return tree_sitter.NewLanguage(tree_sitter_rust.Language()), "rust"
	default:
		return nil, ""
	}
}

// BuildDirectoryTree recursively walks a directory and formats a tree representation.
func BuildDirectoryTree(root string, maxDepth int) (string, []string, error) {
	if maxDepth <= 0 {
		maxDepth = 10
	}
	var filesList []string
	var sb strings.Builder

	info, err := os.Stat(root)
	if err != nil {
		return "", nil, err
	}
	if !info.IsDir() {
		return info.Name(), []string{root}, nil
	}

	sb.WriteString(filepath.Clean(root) + "/\n")

	var walk func(path string, prefix string, depth int) error
	walk = func(path string, prefix string, depth int) error {
		if depth > maxDepth {
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for i, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, ".") && name != ".artifacts" {
				continue
			}
			isLast := i == len(entries)-1
			connector := "├── "
			nextPrefix := prefix + "│   "
			if isLast {
				connector = "└── "
				nextPrefix = prefix + "    "
			}

			fullPath := filepath.Join(path, name)
			if entry.IsDir() {
				sb.WriteString(fmt.Sprintf("%s%s%s/\n", prefix, connector, name))
				if err := walk(fullPath, nextPrefix, depth+1); err != nil {
					return err
				}
			} else {
				sb.WriteString(fmt.Sprintf("%s%s%s\n", prefix, connector, name))
				filesList = append(filesList, fullPath)
			}
		}
		return nil
	}

	if err := walk(root, "", 1); err != nil {
		return "", nil, err
	}

	return sb.String(), filesList, nil
}

// ParseFileStructure uses Tree-Sitter to extract structured code elements from a source file.
func ParseFileStructure(filePath string) (*FileStructureOutput, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	ext := filepath.Ext(filePath)
	tsLang, langName := getTreeSitterLanguage(ext)
	if tsLang == nil {
		// Fallback for languages without tree-sitter grammar: simple line scan
		return fallbackParseStructure(filePath, string(content)), nil
	}

	parser := tree_sitter.NewParser()
	defer parser.Close()
	if err := parser.SetLanguage(tsLang); err != nil {
		return nil, fmt.Errorf("failed to set tree-sitter language: %w", err)
	}

	tree := parser.Parse(content, nil)
	if tree == nil {
		return fallbackParseStructure(filePath, string(content)), nil
	}
	defer tree.Close()

	elements, imports := extractElementsFromTree(tree.RootNode(), content, langName)

	return &FileStructureOutput{
		FilePath: filePath,
		Language: langName,
		Elements: elements,
		Imports:  imports,
		RawCode:  string(content),
	}, nil
}

func extractElementsFromTree(root *tree_sitter.Node, src []byte, lang string) ([]CodeElement, []string) {
	var elements []CodeElement
	var imports []string

	if root == nil {
		return elements, imports
	}

	count := root.NamedChildCount()
	for i := uint(0); i < count; i++ {
		node := root.NamedChild(i)
		if node == nil {
			continue
		}
		nodeType := node.Kind()
		startRow := int(node.StartPosition().Row) + 1
		endRow := int(node.EndPosition().Row) + 1

		switch lang {
		case "c", "cpp":
			switch nodeType {
			case "preproc_include":
				imports = append(imports, strings.TrimSpace(node.Utf8Text(src)))
			case "function_definition":
				name := extractChildByKind(node, src, "function_declarator")
				if name == "" {
					name = node.Utf8Text(src)
					if idx := strings.Index(name, "{"); idx > 0 {
						name = strings.TrimSpace(name[:idx])
					}
				}
				sig := firstLine(node.Utf8Text(src))
				elements = append(elements, CodeElement{
					Kind:      "function",
					Name:      name,
					Signature: sig,
					Line:      startRow,
					EndLine:   endRow,
				})
			case "declaration":
				txt := node.Utf8Text(src)
				kind := "declaration"
				if strings.HasPrefix(txt, "typedef") {
					kind = "typedef"
				}
				elements = append(elements, CodeElement{
					Kind:      kind,
					Name:      firstLine(txt),
					Signature: firstLine(txt),
					Line:      startRow,
					EndLine:   endRow,
				})
			case "struct_specifier":
				elements = append(elements, CodeElement{
					Kind:      "struct",
					Name:      firstLine(node.Utf8Text(src)),
					Signature: firstLine(node.Utf8Text(src)),
					Line:      startRow,
					EndLine:   endRow,
				})
			}

		case "rust":
			switch nodeType {
			case "use_declaration":
				imports = append(imports, strings.TrimSpace(node.Utf8Text(src)))
			case "function_item":
				nameNode := node.ChildByFieldName("name")
				name := ""
				if nameNode != nil {
					name = nameNode.Utf8Text(src)
				}
				elements = append(elements, CodeElement{
					Kind:      "function",
					Name:      name,
					Signature: firstLine(node.Utf8Text(src)),
					Line:      startRow,
					EndLine:   endRow,
				})
			case "struct_item":
				nameNode := node.ChildByFieldName("name")
				name := ""
				if nameNode != nil {
					name = nameNode.Utf8Text(src)
				}
				elements = append(elements, CodeElement{
					Kind:      "struct",
					Name:      name,
					Signature: firstLine(node.Utf8Text(src)),
					Line:      startRow,
					EndLine:   endRow,
				})
			case "impl_item":
				elements = append(elements, CodeElement{
					Kind:      "impl",
					Name:      firstLine(node.Utf8Text(src)),
					Signature: firstLine(node.Utf8Text(src)),
					Line:      startRow,
					EndLine:   endRow,
				})
			}
		}
	}

	return elements, imports
}

func extractChildByKind(n *tree_sitter.Node, src []byte, kind string) string {
	for i := uint(0); i < n.ChildCount(); i++ {
		c := n.Child(i)
		if c.Kind() == kind {
			return c.Utf8Text(src)
		}
	}
	return ""
}

func firstLine(s string) string {
	if idx := strings.Index(s, "\n"); idx >= 0 {
		return strings.TrimSpace(s[:idx])
	}
	return strings.TrimSpace(s)
}

func fallbackParseStructure(filePath, content string) *FileStructureOutput {
	lines := strings.Split(content, "\n")
	var elements []CodeElement
	var imports []string

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#include") || strings.HasPrefix(trimmed, "import") || strings.HasPrefix(trimmed, "use ") {
			imports = append(imports, trimmed)
		} else if strings.Contains(trimmed, "struct ") {
			elements = append(elements, CodeElement{
				Kind:      "struct",
				Name:      trimmed,
				Signature: trimmed,
				Line:      i + 1,
				EndLine:   i + 1,
			})
		} else if strings.Contains(trimmed, "(") && strings.Contains(trimmed, ")") && (strings.HasSuffix(trimmed, "{") || strings.HasSuffix(trimmed, ";")) {
			elements = append(elements, CodeElement{
				Kind:      "function",
				Name:      trimmed,
				Signature: trimmed,
				Line:      i + 1,
				EndLine:   i + 1,
			})
		}
	}

	return &FileStructureOutput{
		FilePath: filePath,
		Language: filepath.Ext(filePath),
		Elements: elements,
		Imports:  imports,
		RawCode:  content,
	}
}

// NewPATools creates Tree-Sitter AST inspection and directory hierarchy tools for repository analysis.
func NewPATools() []tool.BaseTool {
	treeTool, _ := utils.InferTool("get_directory_tree", "Returns a structured representation of the project directory hierarchy",
		func(ctx context.Context, input *DirectoryTreeInput) (*DirectoryTreeOutput, error) {
			dir := input.DirectoryPath
			if dir == "" {
				dir = "."
			}
			treeStr, files, err := BuildDirectoryTree(dir, input.MaxDepth)
			if err != nil {
				return nil, err
			}
			return &DirectoryTreeOutput{
				DirectoryPath: dir,
				TreeString:    treeStr,
				Files:         files,
			}, nil
		})

	structureTool, _ := utils.InferTool("get_file_structure", "Generates a structured representation of a given source file via Tree-Sitter AST (classes, functions, structs, globals, imports)",
		func(ctx context.Context, input *FileStructureInput) (*FileStructureOutput, error) {
			return ParseFileStructure(input.FilePath)
		})

	return []tool.BaseTool{treeTool, structureTool}
}
