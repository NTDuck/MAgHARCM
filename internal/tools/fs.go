package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// ReadFileInput parameters for read_file.
type ReadFileInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path of file to read"`
}

// ReadFileOutput result of read_file.
type ReadFileOutput struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
	Lines    int    `json:"lines"`
}

// WriteFileInput parameters for write_file.
type WriteFileInput struct {
	FilePath string `json:"file_path" jsonschema_description:"Path of file to write"`
	Content  string `json:"content" jsonschema_description:"Content to write to file"`
}

// WriteFileOutput result of write_file.
type WriteFileOutput struct {
	FilePath string `json:"file_path"`
	Success  bool   `json:"success"`
	Bytes    int    `json:"bytes"`
	Message  string `json:"message"`
}

// ListFilesInput parameters for list_files.
type ListFilesInput struct {
	Directory string `json:"directory" jsonschema_description:"Directory path to list files for"`
	Pattern   string `json:"pattern,omitempty" jsonschema_description:"Optional glob pattern, e.g. *.rs or *.c"`
}

// ListFilesOutput result of list_files.
type ListFilesOutput struct {
	Directory string   `json:"directory"`
	Files     []string `json:"files"`
	Count     int      `json:"count"`
}

// CleanCodeContent removes markdown code fences and placeholder lines from source code.
func CleanCodeContent(code string) string {
	lines := strings.Split(code, "\n")
	var cleaned []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "```") || trimmed == "[content]" || strings.HasPrefix(trimmed, "[full ") {
			continue
		}
		cleaned = append(cleaned, l)
	}
	return strings.TrimSpace(strings.Join(cleaned, "\n")) + "\n"
}

// NewFSTools creates file system tools.
func NewFSTools() []tool.BaseTool {
	readTool, _ := utils.InferTool("read_file", "Reads the entire content of a file",
		func(ctx context.Context, input *ReadFileInput) (*ReadFileOutput, error) {
			data, err := os.ReadFile(input.FilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", input.FilePath, err)
			}
			content := string(data)
			lines := len(strings.Split(content, "\n"))
			return &ReadFileOutput{
				FilePath: input.FilePath,
				Content:  content,
				Lines:    lines,
			}, nil
		})

	writeTool, _ := utils.InferTool("write_file", "Writes content to a file, creating parent directories if necessary",
		func(ctx context.Context, input *WriteFileInput) (*WriteFileOutput, error) {
			if err := os.MkdirAll(filepath.Dir(input.FilePath), 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory: %w", err)
			}
			if err := os.WriteFile(input.FilePath, []byte(input.Content), 0644); err != nil {
				return nil, fmt.Errorf("failed to write file %s: %w", input.FilePath, err)
			}
			return &WriteFileOutput{
				FilePath: input.FilePath,
				Success:  true,
				Bytes:    len(input.Content),
				Message:  "File written successfully",
			}, nil
		})

	listTool, _ := utils.InferTool("list_files", "Lists files in a directory matching an optional pattern",
		func(ctx context.Context, input *ListFilesInput) (*ListFilesOutput, error) {
			dir := input.Directory
			if dir == "" {
				dir = "."
			}
			pattern := input.Pattern
			if pattern == "" {
				pattern = "*"
			}

			var matched []string
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				rel, _ := filepath.Rel(dir, path)
				if match, _ := filepath.Match(pattern, info.Name()); match || pattern == "*" {
					matched = append(matched, rel)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			return &ListFilesOutput{
				Directory: dir,
				Files:     matched,
				Count:     len(matched),
			}, nil
		})

	return []tool.BaseTool{readTool, writeTool, listTool}
}
