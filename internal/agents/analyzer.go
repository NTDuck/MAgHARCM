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

	prompt := fmt.Sprintf(`You are the Analyzer Agent in the ReCodeAgent multi-agent repository-level translation workflow.
Your goal is to perform analysis of the source project and generate 3 documents:
1. Source Project Research
2. Third-Party Library Analysis
3. Target Project Design

Source Language: %s
Target Language: %s
Source Codebase Directory: %s

Directory Tree:
%s

AST File Structure Summary:
%s

Source Files Content:
%s

Generate the 3 documents in structured markdown format with exact headers:

# Source Project Research
## Overview
## Directory Structure
## Structs and Interfaces
## Data Models
## Error Handling
## Dependencies

# Third-Party Library Analysis
## Standard and Third-Party Libraries
Identify all libraries/headers used and their idiomatic equivalents in %s.

# Target Project Design
## Overview
## Translation Requirements
## Source Files to Translate
## Module Structure
## Error Handling
## Third-Party Libraries
`, state.Task.SourceLang, state.Task.TargetLang, state.Task.SourceDir, treeStr, structureSummary, allCode, state.Task.TargetLang)

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
