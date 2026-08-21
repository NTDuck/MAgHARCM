package tools

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"

	. "MAgHARCM/internal/patterns"
)

// GitStatusInput parameters for git status.
type GitStatusInput struct {
	RepoDir string `json:"repo_dir,omitempty" jsonschema_description:"Repository directory (default: .)"`
}

// GitStatusOutput result of git status.
type GitStatusOutput struct {
	Clean     bool     `json:"clean"`
	Modified  []string `json:"modified"`
	Untracked []string `json:"untracked"`
	Staged    []string `json:"staged"`
	RawStatus string   `json:"raw_status"`
}

// GitDiffInput parameters for git diff.
type GitDiffInput struct {
	RepoDir  string `json:"repo_dir,omitempty" jsonschema_description:"Repository directory (default: .)"`
	FilePath string `json:"file_path,omitempty" jsonschema_description:"Optional file path to diff"`
	Staged   bool   `json:"staged,omitempty" jsonschema_description:"Show staged changes if true"`
}

// GitDiffOutput result of git diff.
type GitDiffOutput struct {
	Diff string `json:"diff"`
}

// GitLogInput parameters for git log.
type GitLogInput struct {
	RepoDir string `json:"repo_dir,omitempty" jsonschema_description:"Repository directory (default: .)"`
	Limit   int    `json:"limit,omitempty" jsonschema_description:"Number of commits to retrieve (default: 5)"`
}

// GitLogOutput result of git log.
type GitLogOutput struct {
	Commits []string `json:"commits"`
}

// NewGitTools constructs the git tool group.
func NewGitTools() []tool.BaseTool {
	gitStatusTool := Must(utils.InferTool(
		"git_status",
		"Git: Check the working tree status, listing modified, untracked, and staged files",
		func(ctx context.Context, input *GitStatusInput) (*GitStatusOutput, error) {
			dir := input.RepoDir
			if dir == "" {
				dir = "."
			}
			cmd := exec.CommandContext(ctx, "git", "status", "--porcelain")
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			if err != nil {
				return nil, fmt.Errorf("git status failed: %w (output: %s)", err, string(out))
			}

			var modified, untracked, staged []string
			scanner := bufio.NewScanner(strings.NewReader(string(out)))
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) < 3 {
					continue
				}
				code := line[:2]
				file := strings.TrimSpace(line[3:])

				if strings.HasPrefix(code, "??") {
					untracked = append(untracked, file)
				} else if code[0] != ' ' && code[0] != '?' {
					staged = append(staged, file)
				}
				if code[1] == 'M' || code[1] == 'D' {
					modified = append(modified, file)
				}
			}

			return &GitStatusOutput{
				Clean:     len(modified) == 0 && len(untracked) == 0 && len(staged) == 0,
				Modified:  modified,
				Untracked: untracked,
				Staged:    staged,
				RawStatus: strings.TrimSpace(string(out)),
			}, nil
		},
	))

	gitDiffTool := Must(utils.InferTool(
		"git_diff",
		"Git: Show changes in the repository or a specific file",
		func(ctx context.Context, input *GitDiffInput) (*GitDiffOutput, error) {
			dir := input.RepoDir
			if dir == "" {
				dir = "."
			}
			args := []string{"diff"}
			if input.Staged {
				args = append(args, "--staged")
			}
			if input.FilePath != "" {
				args = append(args, "--", input.FilePath)
			}
			cmd := exec.CommandContext(ctx, "git", args...)
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			if err != nil {
				return nil, fmt.Errorf("git diff failed: %w (output: %s)", err, string(out))
			}
			return &GitDiffOutput{
				Diff: string(out),
			}, nil
		},
	))

	gitLogTool := Must(utils.InferTool(
		"git_log",
		"Git: Show recent commit logs in oneline format",
		func(ctx context.Context, input *GitLogInput) (*GitLogOutput, error) {
			dir := input.RepoDir
			if dir == "" {
				dir = "."
			}
			limit := input.Limit
			if limit <= 0 {
				limit = 5
			}
			cmd := exec.CommandContext(ctx, "git", "log", fmt.Sprintf("-n%d", limit), "--oneline")
			cmd.Dir = dir
			out, err := cmd.CombinedOutput()
			if err != nil {
				return nil, fmt.Errorf("git log failed: %w (output: %s)", err, string(out))
			}
			var commits []string
			scanner := bufio.NewScanner(strings.NewReader(string(out)))
			for scanner.Scan() {
				if text := strings.TrimSpace(scanner.Text()); text != "" {
					commits = append(commits, text)
				}
			}
			return &GitLogOutput{
				Commits: commits,
			}, nil
		},
	))

	return []tool.BaseTool{gitStatusTool, gitDiffTool, gitLogTool}
}
