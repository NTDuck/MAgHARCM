package agents

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
	"MAgHARCM/internal/types"
)

// AnalyzerAgent maps the source codebase hierarchy, identifies third-party library dependencies, and drafts the target architecture.
type AnalyzerAgent struct {
	Model model.BaseChatModel
}

// NewAnalyzerAgent creates an AnalyzerAgent instance.
func NewAnalyzerAgent(m model.BaseChatModel) *AnalyzerAgent {
	return &AnalyzerAgent{Model: m}
}

// Run executes the 3-phase analysis workflow and returns updated state.
func (a *AnalyzerAgent) Run(ctx context.Context, state *types.State) (*types.State, error) {
	logger.LogAgent("Analyzer", "Starting source project analysis: `%s` (`%s` -> `%s`)",
		state.Task.SourceDir, state.Task.SourceLang, state.Task.TargetLang)

	treeStr, files, err := a.discoverSourceFiles(state.Task.SourceDir)
	if err != nil {
		return nil, err
	}

	fileStructures, fileContents := a.extractFileStructures(files)
	rawDoc, err := a.synthesizeAnalysis(ctx, state, treeStr, strings.Join(fileStructures, "\n"), strings.Join(fileContents, "\n"))
	if err != nil {
		return nil, err
	}

	a.populateAnalyzerOutput(state, rawDoc)
	logger.LogAgent("Analyzer", "Analysis complete: Research, Library Analysis, and Target Design generated")
	return state, nil
}

// discoverSourceFiles scans the source directory hierarchy.
func (a *AnalyzerAgent) discoverSourceFiles(sourceDir string) (string, []string, error) {
	logger.LogStep("Scanning directory hierarchy via `get_directory_tree`")
	treeStr, files, err := tools.BuildDirectoryTree(sourceDir, 5)
	if err != nil {
		return "", nil, fmt.Errorf("failed to build directory tree: %w", err)
	}
	logger.LogTool("get_directory_tree", "Found %d source/header/test files in `%s`", len(files), sourceDir)
	return treeStr, files, nil
}

// extractFileStructures parses the AST structure and imports for discovered files.
func (a *AnalyzerAgent) extractFileStructures(files []string) ([]string, []string) {
	var fileStructures []string
	var fileContents []string

	for _, f := range files {
		logger.LogStep("Parsing AST structure with Tree-Sitter: `%s`", filepath.Base(f))
		structOut, err := tools.ParseFileStructure(f)
		if err == nil {
			elemDesc := fmt.Sprintf("File: %s\nLanguage: %s\nElements:\n", filepath.Base(f), structOut.Language)
			for _, el := range structOut.Elements {
				elemDesc += fmt.Sprintf("  - %s: %s (lines %d-%d)\n", el.Kind, el.Name, el.Line, el.EndLine)
			}
			fileStructures = append(fileStructures, elemDesc)
			logger.LogTool("get_file_structure", "`%s` -> %d AST elements, %d imports",
				filepath.Base(f), len(structOut.Elements), len(structOut.Imports))
			if structOut.RawCode != "" {
				fileContents = append(fileContents, fmt.Sprintf("=== File: %s ===\n%s\n", filepath.Base(f), structOut.RawCode))
			}
		}
	}
	return fileStructures, fileContents
}

// synthesizeAnalysis queries the reasoning model with directory structure, AST elements, and file contents.
func (a *AnalyzerAgent) synthesizeAnalysis(ctx context.Context, state *types.State, treeStr, structureSummary, allCode string) (string, error) {
	prompt, err := renderPrompt("analyzer.md", map[string]any{
		"SourceLang":         state.Task.SourceLang,
		"TargetLang":         state.Task.TargetLang,
		"SourceDir":          state.Task.SourceDir,
		"DirectoryTree":      treeStr,
		"StructureSummary":   structureSummary,
		"SourceFilesContent": allCode,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render analyzer prompt: %w", err)
	}

	resp, err := a.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert software architect and compiler researcher specializing in repository-level code translation."),
		schema.UserMessage(prompt),
	})
	if err != nil {
		return "", fmt.Errorf("analyzer model call failed: %w", err)
	}
	return resp.Content, nil
}

// populateAnalyzerOutput unpacks markdown sections into structured documents on state.
func (a *AnalyzerAgent) populateAnalyzerOutput(state *types.State, rawDoc string) {
	state.AnalyzerOutput.Research = types.DocumentWrapper[types.SourceProjectResearch]{
		Data: types.SourceProjectResearch{
			Overview:           extractSection(rawDoc, "## 1. Overview", "## 2. Directory Structure"),
			DirectoryStructure: extractSection(rawDoc, "## 2. Directory Structure", "## 3. Data Structures"),
		},
		RawMarkdown: rawDoc,
	}

	state.AnalyzerOutput.Library = types.DocumentWrapper[types.ThirdPartyLibraryAnalysis]{
		Data: types.ThirdPartyLibraryAnalysis{
			Libraries: []types.LibraryMapping{},
		},
		RawMarkdown: extractSection(rawDoc, "=== SECTION: LIBRARY_ANALYSIS ===", "=== SECTION: TARGET_DESIGN ==="),
	}

	state.AnalyzerOutput.Design = types.DocumentWrapper[types.TargetProjectDesign]{
		Data: types.TargetProjectDesign{
			Overview: extractSection(rawDoc, "## Target Architecture", "## Module Decomposition"),
		},
		RawMarkdown: extractSection(rawDoc, "=== SECTION: TARGET_DESIGN ===", ""),
	}
}

func extractSection(doc, startHeader, endHeader string) string {
	startIdx := strings.Index(doc, startHeader)
	if startIdx == -1 {
		return ""
	}
	content := doc[startIdx+len(startHeader):]
	if endHeader != "" {
		endIdx := strings.Index(content, endHeader)
		if endIdx != -1 {
			content = content[:endIdx]
		}
	}
	return strings.TrimSpace(content)
}
