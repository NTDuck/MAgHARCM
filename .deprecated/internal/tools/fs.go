package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	. "MAgHARCM/internal/patterns"
)

// ReadFileInput parameters for reading a file.
type ReadFileInput struct {
	FilePath string `json:"file_path" jsonschema_description:"The path of the file to read"`
	Offset   int    `json:"offset,omitempty" jsonschema_description:"1-based line number to start reading from (default: 1)"`
	Limit    int    `json:"limit,omitempty" jsonschema_description:"Maximum number of lines to read (default: 2000)"`
}

// ReadFileOutput result of reading a file.
type ReadFileOutput struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

// WriteFileInput parameters for writing a file.
type WriteFileInput struct {
	FilePath string `json:"file_path" jsonschema_description:"The path of the file to write to"`
	Content  string `json:"content" jsonschema_description:"The full content to write to the file"`
}

// WriteFileOutput result of writing a file.
type WriteFileOutput struct {
	FilePath string `json:"file_path"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

// EditFileInput parameters for editing a file.
type EditFileInput struct {
	FilePath   string `json:"file_path" jsonschema_description:"The path of the file to edit"`
	OldString  string `json:"old_string" jsonschema_description:"The exact string to find and replace"`
	NewString  string `json:"new_string" jsonschema_description:"The new replacement string"`
	ReplaceAll bool   `json:"replace_all,omitempty" jsonschema_description:"If true, replaces all occurrences"`
}

// EditFileOutput result of editing a file.
type EditFileOutput struct {
	FilePath string `json:"file_path"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

// GlobInput parameters for globbing files.
type GlobInput struct {
	Path    string `json:"path,omitempty" jsonschema_description:"Base directory path to search in (default: .)"`
	Pattern string `json:"pattern" jsonschema_description:"Glob pattern to match files (e.g. **/*.rs, *.c)"`
}

// GlobOutput result of globbing files.
type GlobOutput struct {
	Path    string   `json:"path"`
	Pattern string   `json:"pattern"`
	Matches []string `json:"matches"`
	Count   int      `json:"count"`
}

// GrepInput parameters for searching file contents.
type GrepInput struct {
	Path            string `json:"path,omitempty" jsonschema_description:"Directory or file path to search in (default: .)"`
	Pattern         string `json:"pattern" jsonschema_description:"Regex pattern to search for"`
	Glob            string `json:"glob,omitempty" jsonschema_description:"File filter pattern (e.g. *.rs)"`
	CaseInsensitive bool   `json:"case_insensitive,omitempty" jsonschema_description:"Ignore case in search"`
}

// GrepOutput result of grep search.
type GrepOutput struct {
	Pattern string                 `json:"pattern"`
	Matches []filesystem.GrepMatch `json:"matches"`
	Count   int                    `json:"count"`
}

// ListDirInput parameters for listing directory entries.
type ListDirInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Directory path to list (default: .)"`
}

// ListDirOutput result of listing directory.
type ListDirOutput struct {
	Path    string                `json:"path"`
	Entries []filesystem.FileInfo `json:"entries"`
	Count   int                   `json:"count"`
}

// NewFSTools constructs the filesystem tool group by delegating to Eino's filesystem.Backend.
func NewFSTools(backend filesystem.Backend) []tool.BaseTool {
	readFileTool := Must(utils.InferTool(
		"read_file",
		"Read the contents of a file at the specified path with optional offset and limit",
		func(ctx context.Context, input *ReadFileInput) (*ReadFileOutput, error) {
			path := filepath.Clean(input.FilePath)
			if path == "" || path == "." {
				return nil, fmt.Errorf("file_path is required")
			}
			res, err := backend.Read(ctx, &filesystem.ReadRequest{
				FilePath: path,
				Offset:   input.Offset,
				Limit:    input.Limit,
			})
			if err != nil {
				return nil, err
			}
			return &ReadFileOutput{
				FilePath: path,
				Content:  res.Content,
			}, nil
		},
	))

	writeFileTool := Must(utils.InferTool(
		"write_file",
		"Write complete text content to a file at the specified path (creates parent directories if needed)",
		func(ctx context.Context, input *WriteFileInput) (*WriteFileOutput, error) {
			path := filepath.Clean(input.FilePath)
			if path == "" || path == "." {
				return nil, fmt.Errorf("file_path is required for write_file")
			}
			err := backend.Write(ctx, &filesystem.WriteRequest{
				FilePath: path,
				Content:  input.Content,
			})
			if err != nil {
				return nil, err
			}
			return &WriteFileOutput{
				FilePath: path,
				Success:  true,
				Message:  fmt.Sprintf("Successfully wrote %d bytes to %s", len(input.Content), path),
			}, nil
		},
	))

	editFileTool := Must(utils.InferTool(
		"edit_file",
		"Edit an existing file by replacing old_string with new_string",
		func(ctx context.Context, input *EditFileInput) (*EditFileOutput, error) {
			path := filepath.Clean(input.FilePath)
			if path == "" || path == "." {
				return nil, fmt.Errorf("file_path is required for edit_file")
			}
			err := backend.Edit(ctx, &filesystem.EditRequest{
				FilePath:   path,
				OldString:  input.OldString,
				NewString:  input.NewString,
				ReplaceAll: input.ReplaceAll,
			})
			if err != nil {
				return nil, err
			}
			return &EditFileOutput{
				FilePath: path,
				Success:  true,
				Message:  fmt.Sprintf("Successfully edited %s", path),
			}, nil
		},
	))

	globTool := Must(utils.InferTool(
		"glob",
		"Find files matching a glob pattern (e.g. **/*.rs, *.c) in a directory",
		func(ctx context.Context, input *GlobInput) (*GlobOutput, error) {
			basePath := input.Path
			if basePath == "" {
				basePath = "."
			}
			basePath = filepath.Clean(basePath)
			pattern := input.Pattern
			if pattern == "" {
				pattern = "*"
			}

			matches, err := backend.GlobInfo(ctx, &filesystem.GlobInfoRequest{
				Path:    basePath,
				Pattern: pattern,
			})
			if err != nil {
				return nil, err
			}
			var list []string
			for _, m := range matches {
				list = append(list, m.Path)
			}
			return &GlobOutput{
				Path:    basePath,
				Pattern: pattern,
				Matches: list,
				Count:   len(list),
			}, nil
		},
	))

	grepTool := Must(utils.InferTool(
		"grep",
		"Search file contents for lines matching a regex pattern",
		func(ctx context.Context, input *GrepInput) (*GrepOutput, error) {
			basePath := input.Path
			if basePath == "" {
				basePath = "."
			}
			basePath = filepath.Clean(basePath)

			matches, err := backend.GrepRaw(ctx, &filesystem.GrepRequest{
				Path:            basePath,
				Pattern:         input.Pattern,
				Glob:            input.Glob,
				CaseInsensitive: input.CaseInsensitive,
			})
			if err != nil {
				return nil, err
			}

			return &GrepOutput{
				Pattern: input.Pattern,
				Matches: matches,
				Count:   len(matches),
			}, nil
		},
	))

	listDirTool := Must(utils.InferTool(
		"list_dir",
		"List entries in a directory",
		func(ctx context.Context, input *ListDirInput) (*ListDirOutput, error) {
			path := input.Path
			if path == "" {
				path = "."
			}
			path = filepath.Clean(path)
			entries, err := backend.LsInfo(ctx, &filesystem.LsInfoRequest{
				Path: path,
			})
			if err != nil {
				return nil, err
			}
			return &ListDirOutput{
				Path:    path,
				Entries: entries,
				Count:   len(entries),
			}, nil
		},
	))

	return []tool.BaseTool{readFileTool, writeFileTool, editFileTool, globTool, grepTool, listDirTool}
}
