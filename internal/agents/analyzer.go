package agents

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"MAgHARCM/internal/artifacts"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
	"MAgHARCM/internal/types"
)

// AnalyzerAgent performs source project research, library analysis, and target project design (§3.2).
type AnalyzerAgent struct {
	Model model.BaseChatModel
}

// NewAnalyzerAgent creates an AnalyzerAgent instance.
func NewAnalyzerAgent(m model.BaseChatModel) *AnalyzerAgent {
	return &AnalyzerAgent{Model: m}
}

// Run executes the 3-phase analysis workflow and returns updated state.
func (a *AnalyzerAgent) Run(ctx context.Context, state *types.State) (*types.State, error) {
	logger.LogAgent("Analyzer", "Starting source project analysis: %s (%s -> %s)",
		state.Task.SourceDir, state.Task.SourceLang, state.Task.TargetLang)
	state.Log("[Analyzer] Starting Source Project Research for %s (%s -> %s)", state.Task.SourceDir, state.Task.SourceLang, state.Task.TargetLang)

	// 1. Source Project Research (§3.2.1)
	logger.LogStep("Scanning directory hierarchy via get_directory_tree...")
	treeStr, files, err := tools.BuildDirectoryTree(state.Task.SourceDir, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to build directory tree: %w", err)
	}
	logger.LogTool("get_directory_tree", "Found %d source/header/test files in %s", len(files), state.Task.SourceDir)

	var fileStructures []string
	var fileContents []string
	for _, f := range files {
		logger.LogStep("Parsing AST structure with Tree-Sitter: %s", filepath.Base(f))
		structOut, err := tools.ParseFileStructure(f)
		if err == nil {
			elemDesc := fmt.Sprintf("File: %s\nLanguage: %s\nElements:\n", filepath.Base(f), structOut.Language)
			for _, el := range structOut.Elements {
				elemDesc += fmt.Sprintf("  - [%s] %s (lines %d-%d): %s\n", el.Kind, el.Name, el.Line, el.EndLine, el.Signature)
			}
			fileStructures = append(fileStructures, elemDesc)
			logger.LogTool("get_file_structure", "%s -> %d AST elements, %d imports",
				filepath.Base(f), len(structOut.Elements), len(structOut.Imports))
			if structOut.RawCode != "" {
				fileContents = append(fileContents, fmt.Sprintf("=== File: %s ===\n%s\n", filepath.Base(f), structOut.RawCode))
			}
		}
	}

	allCode := strings.Join(fileContents, "\n")
	structureSummary := strings.Join(fileStructures, "\n")

	logger.LogStep("Synthesizing Source Research, 3rd-Party Library Analysis, and Target Design with Reasoning Model...")

	prompt, err := renderPrompt("analyzer.md", map[string]any{
		"SourceLang":          state.Task.SourceLang,
		"TargetLang":          state.Task.TargetLang,
		"SourceDir":           state.Task.SourceDir,
		"DirectoryTree":       treeStr,
		"StructureSummary":    structureSummary,
		"SourceFilesContent":  allCode,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render analyzer prompt: %w", err)
	}
	resp, err := a.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert software architect and compiler researcher specializing in repository-level code translation."),
		schema.UserMessage(prompt),
	})
	if err != nil {
		return nil, fmt.Errorf("analyzer model call failed: %w", err)
	}

	rawDoc := resp.Content

	state.AnalyzerOutput.Research = types.DocumentWrapper[types.SourceProjectResearch]{
		Data: types.SourceProjectResearch{
			Overview:           extractSection(rawDoc, "Source Project Research", "Third-Party Library Analysis"),
			DirectoryStructure: treeStr,
			RawDocument:        rawDoc,
		},
		RawMarkdown: rawDoc,
	}

	state.AnalyzerOutput.Library = types.DocumentWrapper[types.ThirdPartyLibraryAnalysis]{
		Data: types.ThirdPartyLibraryAnalysis{
			RawDocument: extractSection(rawDoc, "Third-Party Library Analysis", "Target Project Design"),
		},
		RawMarkdown: rawDoc,
	}

	state.AnalyzerOutput.Design = types.DocumentWrapper[types.TargetProjectDesign]{
		Data: types.TargetProjectDesign{
			Overview:    extractSection(rawDoc, "Target Project Design", ""),
			RawDocument: rawDoc,
		},
		RawMarkdown: rawDoc,
	}

	_ = artifacts.SaveAnalyzerOutput(state.Task.TargetDir, state.AnalyzerOutput)
	logger.LogAgent("Analyzer", "Analysis complete: Research, Library Analysis, and Target Design generated")
	state.Log("[Analyzer] Research and Target Project Design generated successfully")
	return state, nil
}

func extractSection(doc, startHeader, endHeader string) string {
	startIdx := strings.Index(doc, startHeader)
	if startIdx == -1 {
		return doc
	}
	sub := doc[startIdx:]
	if endHeader != "" {
		endIdx := strings.Index(sub, endHeader)
		if endIdx != -1 {
			return strings.TrimSpace(sub[:endIdx])
		}
	}
	return strings.TrimSpace(sub)
}
