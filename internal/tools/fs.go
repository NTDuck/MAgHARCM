package tools

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
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
	FilePath   string `json:"file_path"`
	Content    string `json:"content"`
	StartLine  int    `json:"start_line"`
	LinesCount int    `json:"lines_count"`
}

// WriteFileInput parameters for writing a file.
type WriteFileInput struct {
	FilePath string `json:"file_path" jsonschema_description:"The path of the file to write to"`
	Content  string `json:"content" jsonschema_description:"The full content to write to the file"`
}

// WriteFileOutput result of writing a file.
type WriteFileOutput struct {
	FilePath string `json:"file_path"`
	Bytes    int    `json:"bytes_written"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

// EditFileInput parameters for editing a file.
type EditFileInput struct {
	FilePath   string `json:"file_path" jsonschema_description:"The path of the file to edit"`
	OldString  string `json:"old_string" jsonschema_description:"The exact string to find and replace"`
	NewString  string `json:"new_string" jsonschema_description:"The new replacement string"`
	ReplaceAll bool   `json:"replace_all,omitempty" jsonschema_description:"If true, replaces all occurrences; if false, fails unless exactly one match is found"`
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

// GrepMatchItem single grep match result.
type GrepMatchItem struct {
	Path    string `json:"path"`
	Line    int    `json:"line"`
	Content string `json:"content"`
}

// GrepOutput result of grep search.
type GrepOutput struct {
	Pattern string          `json:"pattern"`
	Matches []GrepMatchItem `json:"matches"`
	Count   int             `json:"count"`
}

// ListDirInput parameters for listing directory entries.
type ListDirInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Directory path to list (default: .)"`
}

// ListDirOutput result of listing directory.
type ListDirOutput struct {
	Path    string   `json:"path"`
	Entries []string `json:"entries"`
	Count   int      `json:"count"`
}

// NewFSTools constructs the filesystem tool group.
func NewFSTools(backend filesystem.Backend) []tool.BaseTool {
	readFileTool := Must(utils.InferTool(
		"read_file",
		"Read the contents of a file at the specified path with optional offset and limit",
		func(ctx context.Context, input *ReadFileInput) (*ReadFileOutput, error) {
			path := filepath.Clean(input.FilePath)
			if path == "" || path == "." {
				return nil, fmt.Errorf("file_path is required")
			}
			file, err := os.Open(path)
			if err != nil {
				return nil, fmt.Errorf("failed to open file %s: %w", path, err)
			}
			defer file.Close()

			offset := input.Offset
			if offset <= 0 {
				offset = 1
			}
			limit := input.Limit
			if limit <= 0 {
				limit = 2000
			}

			scanner := bufio.NewScanner(file)
			var sb strings.Builder
			lineNum := 1
			linesRead := 0

			for scanner.Scan() {
				if lineNum >= offset {
					sb.WriteString(fmt.Sprintf("%6d\t%s\n", lineNum, scanner.Text()))
					linesRead++
					if linesRead >= limit {
						break
					}
				}
				lineNum++
			}
			if err := scanner.Err(); err != nil && err != io.EOF {
				return nil, fmt.Errorf("error reading file %s: %w", path, err)
			}

			return &ReadFileOutput{
				FilePath:   path,
				Content:    sb.String(),
				StartLine:  offset,
				LinesCount: linesRead,
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
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory for %s: %w", path, err)
			}
			if err := os.WriteFile(path, []byte(input.Content), 0644); err != nil {
				return nil, fmt.Errorf("failed to write file %s: %w", path, err)
			}
			return &WriteFileOutput{
				FilePath: path,
				Bytes:    len(input.Content),
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
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", path, err)
			}
			content := string(data)
			if !strings.Contains(content, input.OldString) {
				return nil, fmt.Errorf("old_string not found in %s", path)
			}
			var newContent string
			if input.ReplaceAll {
				newContent = strings.ReplaceAll(content, input.OldString, input.NewString)
			} else {
				if strings.Count(content, input.OldString) > 1 {
					return nil, fmt.Errorf("old_string matched multiple times in %s; set replace_all=true or provide more context", path)
				}
				newContent = strings.Replace(content, input.OldString, input.NewString, 1)
			}
			if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
				return nil, fmt.Errorf("failed to write edited file %s: %w", path, err)
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

			var matches []string
			err := filepath.WalkDir(basePath, func(p string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				rel, err := filepath.Rel(basePath, p)
				if err != nil || rel == "." {
					return nil
				}
				relSlash := filepath.ToSlash(rel)
				matched, _ := doublestar.Match(pattern, relSlash)
				if matched {
					matches = append(matches, filepath.ToSlash(filepath.Join(basePath, rel)))
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("glob walk error: %w", err)
			}
			sort.Strings(matches)
			return &GlobOutput{
				Path:    basePath,
				Pattern: pattern,
				Matches: matches,
				Count:   len(matches),
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

			flags := ""
			if input.CaseInsensitive {
				flags = "(?i)"
			}
			re, err := regexp.Compile(flags + input.Pattern)
			if err != nil {
				return nil, fmt.Errorf("invalid regex pattern %s: %w", input.Pattern, err)
			}

			var matches []GrepMatchItem
			err = filepath.WalkDir(basePath, func(p string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				if input.Glob != "" {
					matched, _ := doublestar.Match(input.Glob, d.Name())
					if !matched {
						return nil
					}
				}
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
						matches = append(matches, GrepMatchItem{
							Path:    filepath.ToSlash(p),
							Line:    lineNo,
							Content: text,
						})
					}
					lineNo++
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("grep error: %w", err)
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
			entries, err := os.ReadDir(path)
			if err != nil {
				return nil, fmt.Errorf("failed to list directory %s: %w", path, err)
			}
			var list []string
			for _, e := range entries {
				name := e.Name()
				if e.IsDir() {
					name += "/"
				}
				list = append(list, name)
			}
			return &ListDirOutput{
				Path:    path,
				Entries: list,
				Count:   len(list),
			}, nil
		},
	))

	return []tool.BaseTool{readFileTool, writeFileTool, editFileTool, globTool, grepTool, listDirTool}
}
