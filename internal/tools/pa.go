package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	"MAgHARCM/internal/languages"
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

// CodeElement represents a structural item in canonical Code IR.
type CodeElement = languages.CodeElement

// FileStructureInput parameters for get_file_structure.
type FileStructureInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path of the source file to analyze"`
}

// FileStructureOutput result of get_file_structure.
type FileStructureOutput = languages.FileStructureResult

// ParseFileStructure extracts canonical Code IR from any source file.
func ParseFileStructure(filePath string) (*FileStructureOutput, error) {
	return languages.ExtractFileStructure(filePath)
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

	cleanRoot := filepath.Clean(root)
	sb.WriteString(fmt.Sprintf("%s/\n", filepath.Base(cleanRoot)))

	var walk func(dir string, prefix string, depth int) error
	walk = func(dir string, prefix string, depth int) error {
		if depth > maxDepth {
			return nil
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		// Filter out ignored / build directories
		var filtered []os.DirEntry
		for _, e := range entries {
			name := e.Name()
			if strings.HasPrefix(name, ".") || name == "target" || name == "bin" || name == "node_modules" || name == "vendor" {
				continue
			}
			filtered = append(filtered, e)
		}

		for i, entry := range filtered {
			isLast := i == len(filtered)-1
			connector := "├── "
			childPrefix := prefix + "│   "
			if isLast {
				connector = "└── "
				childPrefix = prefix + "    "
			}

			path := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				sb.WriteString(fmt.Sprintf("%s%s%s/\n", prefix, connector, entry.Name()))
				if err := walk(path, childPrefix, depth+1); err != nil {
					return err
				}
			} else {
				sb.WriteString(fmt.Sprintf("%s%s%s\n", prefix, connector, entry.Name()))
				filesList = append(filesList, path)
			}
		}
		return nil
	}

	if err := walk(cleanRoot, "", 1); err != nil {
		return "", nil, err
	}

	return sb.String(), filesList, nil
}

// NewPATools creates repository analysis and structural code extraction tools.
func NewPATools() []tool.BaseTool {
	treeTool, _ := utils.InferTool("get_directory_tree", "Returns a structured representation of the repository directory hierarchy across all languages",
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

	structureTool, _ := utils.InferTool("get_file_structure", "Parses and extracts structured Code IR (functions, methods, types, structs, classes, interfaces, and imports) from any source file",
		func(ctx context.Context, input *FileStructureInput) (*FileStructureOutput, error) {
			return ParseFileStructure(input.FilePath)
		})

	return []tool.BaseTool{treeTool, structureTool}
}
