package agents

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/llm"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
)
// Canonical artifact types live in internal/compiletime/state.go (cycle-free
// leaf package; agents → compiletime). This file owns the analyzer algorithm
// and re-exports the artifact types as type aliases for ergonomic use.
// Backlink: [[1.0.0 PRIM-4]] SpecMiner Dyn-Inferred Requirements (see body).

// Type aliases (cycle-free Locality of Behaviour).
// Producer files declare the algorithms; canonical artifact types
// live in internal/compiletime/state.go. Aliases let method receivers
// reference the type name without package qualification.
type AnalyzerOutput = compiletime.AnalyzerOutput
type DocumentWrapper[T any] = compiletime.DocumentWrapper[T]
type LibraryMapping = compiletime.LibraryMapping
type SourceProjectResearch = compiletime.SourceProjectResearch
type TargetProjectDesign = compiletime.TargetProjectDesign
type ThirdPartyLibraryAnalysis = compiletime.ThirdPartyLibraryAnalysis

// AnalyzerAgent maps the source codebase hierarchy, identifies third-party library dependencies, and drafts the target architecture.
type AnalyzerAgent struct {
	Model model.BaseChatModel
	// structured is the typed-output extractor bound to AnalyzerSchema. It is
	// built once in NewAnalyzerAgent and reused across every Run call so the
	// underlying ChatModel is never mutated (WithTools returns a new instance).
	structured *llm.StructuredExtractor[AnalyzerSchema]
	// the source project produces a runnable binary, the analyzer invokes
	// SpecMiner.Recover before synthesizeAnalysis so allocation sizes,
	// pointer nullability, aliasing, lifetime ranges, and branch coverage
	// feed the LLM's migration-strategy rationale. Source-Binary path is
	// taken from AnalyzerSpecMinerConfig.SourceBinary; nil disables.
	SpecMiner               *SpecMiner
	AnalyzerSpecMinerConfig AnalyzerSpecMinerConfig
}

// compiletime.DocumentWrapper keeps both structured data and markdown representation.
// time. SourceBinary is the compiled source artifact path; Inputs are the
// representative inputs the analyzer passes to SpecMiner.Recover to
// exercise the runtime. Zero value disables SpecMiner entirely.
type AnalyzerSpecMinerConfig struct {
	SourceBinary string
	Inputs       []string
}

// NewAnalyzerAgent wires the analyzer to a chat model and prepares its typed-output
// schema (AnalyzerSchema -> emit_analysis tool). Returns an error when the model
// cannot be used as a ToolCallingChatModel — older Ollama builds or test doubles
// that lack WithTools will surface the failure here rather than at Run time.
func NewAnalyzerAgent(m model.BaseChatModel) (*AnalyzerAgent, error) {
	extractor, err := llm.NewStructuredExtractor[AnalyzerSchema](m, "emit_analysis",
		"Emit the structured source-code analysis: overview, directory structure, structs/interfaces, data models, error handling, dependencies, third-party library mappings, and target design.")
	if err != nil {
		return nil, fmt.Errorf("analyzer: build structured extractor: %w", err)
	}
	return &AnalyzerAgent{Model: m, structured: extractor}, nil
}
// Run executes the 3-phase analysis workflow and returns updated state.
func (a *AnalyzerAgent) Run(ctx context.Context, state *compiletime.State) (*compiletime.State, error) {
	logger.LogAgent("Analyzer", "Starting source project analysis: source=%s srcLang=%s tgtLang=%s",
		state.Task.SourceDir, state.Task.SourceLang, state.Task.TargetLang)

	treeStr, files, err := a.discoverSourceFiles(state.Task.SourceDir)
	if err != nil {
		return nil, err
	}

	fileStructures, fileContents := a.extractFileStructures(files)
	// PRIM-4 SpecMiner pre-analysis pass: when a runnable source binary is
	// available, recover dynamic invariants and stash them into the analyzer
	// output so downstream agents see allocation/nullability/lifetime
	// evidence alongside the static AST. Failures are logged but non-fatal —
	// a failed SpecMiner run must not abort the pipeline.
	if a.SpecMiner != nil && a.AnalyzerSpecMinerConfig.SourceBinary != "" {
		invars, err := a.SpecMiner.Recover(ctx, a.AnalyzerSpecMinerConfig.SourceBinary, a.AnalyzerSpecMinerConfig.Inputs)
		if err != nil {
			logger.LogWarning("PRIM-4 SpecMiner pre-analysis failed: %v", err)
		} else {
			logger.LogStep("PRIM-4 SpecMiner recovered %d allocation sizes across %d inputs", len(invars.AllocSizes), len(a.AnalyzerSpecMinerConfig.Inputs))
			state.SpecMinerInvariants = compiletime.SpecMinerInvariants{
				AllocSizes:         invars.AllocSizes,
				PointerNullability: invars.PointerNullability,
				AliasingPairs:      invars.AliasingPairs,
				LifetimeRanges:     invars.LifetimeRanges,
				BranchCoverage:     invars.BranchCoverage,
			}
		}
	}
	rawDoc, err := a.synthesizeAnalysis(ctx, state, treeStr, strings.Join(fileStructures, "\n"), strings.Join(fileContents, "\n"))
	if err != nil {
		return nil, err
	}
	kind, rationale, err := SelectAndTryStrategies(ctx, Profile{
		FileCount: len(files),
		LoC:       countSourceLoC(state.Task.SourceDir),
		HasTests:  true,
		HasBuild:  true,
	})
	if err != nil {
		return nil, err
	}
	a.populateAnalyzerOutput(state, rawDoc, kind, rationale)
	logger.LogAgent("Analyzer", "Analysis complete: strategy=%s rationale=%s; Research, Library Analysis, and Target Design generated",
		kind, rationale)
	return state, nil
}

// discoverSourceFiles scans the source directory hierarchy.
func (a *AnalyzerAgent) discoverSourceFiles(sourceDir string) (string, []string, error) {
	logger.LogStep("Scanning directory hierarchy via `get_directory_tree`")
	treeStr, files, err := tools.BuildDirectoryTree(sourceDir, compiletime.DefaultSourceTreeDepth)
	if err != nil {
		return "", nil, fmt.Errorf("failed to build directory tree: %w", err)
	}
	logger.LogTool("get_directory_tree", "Found %d source/header/test files in `%s`", len(files), sourceDir)
	return treeStr, files, nil
}

// extractFileStructures parses the AST structure and imports for discovered files.
// maxAnalyzerCodeBudgetBytes caps the total raw source code fed to the analyzer reasoning prompt.
const maxAnalyzerCodeBudgetBytes = 32 * 1024

// extractFileStructures parses the AST structure and imports for discovered files.
func (a *AnalyzerAgent) extractFileStructures(files []string) ([]string, []string) {
	var fileStructures []string
	var fileContents []string
	var totalBytes int

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
			if structOut.RawCode != "" && totalBytes < maxAnalyzerCodeBudgetBytes {
				chunk := structOut.RawCode
				if len(chunk) > 2048 {
					chunk = chunk[:2048] + "\n// ... (truncated)"
				}
				fileContents = append(fileContents, fmt.Sprintf("=== File: %s ===\n%s\n", filepath.Base(f), chunk))
				totalBytes += len(chunk)
			}
		}
	}
	return fileStructures, fileContents
}

// synthesizeAnalysis queries the reasoning model with directory structure, AST elements, and file contents.
// The model is forced to emit a single tool call matching AnalyzerSchema so the artifact is
// produced as typed JSON instead of a free-form markdown blob the agent has to regex apart.
func (a *AnalyzerAgent) synthesizeAnalysis(ctx context.Context, state *compiletime.State, treeStr, structureSummary, allCode string) (AnalyzerSchema, error) {
	if a.structured == nil {
		return AnalyzerSchema{}, fmt.Errorf("analyzer: structured extractor not initialised (call NewAnalyzerAgent with a ToolCallingChatModel)")
	}
	prompt, err := renderPromptTemplate("analyzer", analyzerPromptTemplate, map[string]any{
		"SourceLang":         state.Task.SourceLang,
		"TargetLang":         state.Task.TargetLang,
		"SourceDir":          state.Task.SourceDir,
		"DirectoryTree":      treeStr,
		"StructureSummary":   structureSummary,
		"SourceFilesContent": allCode,
	})
	if err != nil {
		return AnalyzerSchema{}, fmt.Errorf("failed to render analyzer prompt: %w", err)
	}

	system := "You are an expert software architect and compiler researcher specializing in repository-level code translation. Always respond by calling the emit_analysis tool with the structured fields; do not emit any other text."
	return a.structured.Extract(ctx, system, prompt)
}

// populateAnalyzerOutput writes the structured AnalyzerSchema into the
// canonical compiletime.AnalyzerOutput on state. The RawMarkdown fields are
// populated from a deterministic render of the structured data so downstream
// agents (planner, etc.) can still consume prose when needed.
func (a *AnalyzerAgent) populateAnalyzerOutput(state *compiletime.State, schema AnalyzerSchema, strategy compiletime.StrategyKind, rationale string) {
	rendered := renderAnalyzerMarkdown(schema, strategy, rationale)
	state.AnalyzerOutput = schema.toAnalyzerOutput(strategy, rationale, rendered)
}

// renderAnalyzerMarkdown produces a deterministic prose view of the typed
// schema. The planner and downstream prompts occasionally request a
// human-readable version of the research/design; this keeps that contract
// without dragging in the brittle "scan for ## headers" parsing path.
func renderAnalyzerMarkdown(s AnalyzerSchema, strategy compiletime.StrategyKind, rationale string) string {
	var b strings.Builder
	b.WriteString("# Source Project Research\n\n")
	b.WriteString("## 1. Overview\n\n")
	b.WriteString(s.Overview)
	b.WriteString("\n\n## 2. Directory Structure\n\n")
	b.WriteString(s.DirectoryStructure)
	b.WriteString("\n\n## 3. Structs and Interfaces\n\n")
	b.WriteString(s.StructsInterfaces)
	b.WriteString("\n\n## 4. Data Models\n\n")
	b.WriteString(s.DataModels)
	b.WriteString("\n\n## 5. Error Handling\n\n")
	b.WriteString(s.ErrorHandling)
	b.WriteString("\n\n## 6. Dependencies\n\n")
	for _, d := range s.Dependencies {
		b.WriteString("- ")
		b.WriteString(d)
		b.WriteString("\n")
	}
	b.WriteString("\n# Third-Party Library Analysis\n\n")
	for _, lib := range s.Libraries {
		fmt.Fprintf(&b, "## %s -> %s\n\n%s\n\nUsage: %s\n\nRecommendations: %s\n\n",
			lib.SourceLibrary, lib.TargetLibrary, lib.Overview, lib.Usage, lib.Recommendations)
	}
	b.WriteString("# Target Project Design\n\n")
	b.WriteString("## Overview\n\n")
	b.WriteString(s.Design.Overview)
	b.WriteString("\n\n## Module Decomposition\n\n")
	b.WriteString(s.Design.ModuleStructure)
	b.WriteString("\n\n## Migration Strategy: ")
	b.WriteString(string(strategy))
	b.WriteString("\n\n")
	b.WriteString(rationale)
	b.WriteString("\n")
	return b.String()
}

// extractSection is retained as a stub for any caller that still passes a
// raw markdown blob. New code MUST use AnalyzerSchema directly.
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
