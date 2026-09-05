package agents

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"github.com/cloudwego/eino/components/model"
	"MAgHARCM/internal/compiletime"
	"github.com/cloudwego/eino/schema"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
)

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
	Model                model.BaseChatModel
	// the source project produces a runnable binary, the analyzer invokes
	// SpecMiner.Recover before synthesizeAnalysis so allocation sizes,
	// pointer nullability, aliasing, lifetime ranges, and branch coverage
	// feed the LLM's migration-strategy rationale. Source-Binary path is
	// taken from AnalyzerSpecMinerConfig.SourceBinary; nil disables.
	SpecMiner            *SpecMiner
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
func NewAnalyzerAgent(m model.BaseChatModel) *AnalyzerAgent {
	return &AnalyzerAgent{Model: m}
}

// Run executes the 3-phase analysis workflow and returns updated state.
func (a *AnalyzerAgent) Run(ctx context.Context, state *compiletime.State) (*compiletime.State, error) {
	logger.LogAgent("Analyzer", "Starting source project analysis: `%s` (`%s` -> `%s`)",
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
	a.populateAnalyzerOutput(state, rawDoc, string(kind), rationale)
	logger.LogAgent("Analyzer", "Analysis complete: strategy=%s (%s), Research, Library Analysis, and Target Design generated",
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
func (a *AnalyzerAgent) synthesizeAnalysis(ctx context.Context, state *compiletime.State, treeStr, structureSummary, allCode string) (string, error) {
	prompt, err := renderPromptTemplate("analyzer", analyzerPromptTemplate, map[string]any{
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
func (a *AnalyzerAgent) populateAnalyzerOutput(state *compiletime.State, rawDoc, strategy, rationale string) {
	state.AnalyzerOutput.Research = compiletime.DocumentWrapper[compiletime.SourceProjectResearch]{
		ArtifactSchemaVersion: compiletime.CurrentSchemaVersion,
		Data: compiletime.SourceProjectResearch{
			Overview:           extractSection(rawDoc, "## 1. Overview", "## 2. Directory Structure"),
			DirectoryStructure: extractSection(rawDoc, "## 2. Directory Structure", "## 3. Data Structures"),
			MigrationStrategy:  strategy,
			StrategyRationale:  rationale,
		},
		RawMarkdown: rawDoc,
	}
	state.AnalyzerOutput.Library = compiletime.DocumentWrapper[compiletime.ThirdPartyLibraryAnalysis]{
		Data: compiletime.ThirdPartyLibraryAnalysis{
			Libraries: []compiletime.LibraryMapping{},
		},
		RawMarkdown: extractSection(rawDoc, "=== SECTION: LIBRARY_ANALYSIS ===", "=== SECTION: TARGET_DESIGN ==="),
	}

	state.AnalyzerOutput.Design = compiletime.DocumentWrapper[compiletime.TargetProjectDesign]{
		Data: compiletime.TargetProjectDesign{
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
