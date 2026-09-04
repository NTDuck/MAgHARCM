package agents

// Backlink: [[Methodology]] §1 Stage 2 and [[Primitives]] §NEW-PRIM-20, §NEW-PRIM-21, §NEW-PRIM-26.

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
)

// AnalyzerAgent maps the source codebase hierarchy, identifies third-party library dependencies, and drafts the target architecture.
type AnalyzerAgent struct {
	Model     model.BaseChatModel
	Navigator *Navigator // optional symbol-aware lookup helper (NEW-PRIM-26)
	// SpecMiner is the PRIM-4 dynamic-invariant recovery hook. When set and
	// the source project produces a runnable binary, the analyzer invokes
	// SpecMiner.Recover before synthesizeAnalysis so allocation sizes,
	// pointer nullability, aliasing, lifetime ranges, and branch coverage
	// feed the LLM's migration-strategy rationale. Source-Binary path is
	// taken from AnalyzerSpecMinerConfig.SourceBinary; nil disables.
	SpecMiner            *SpecMiner
	AnalyzerSpecMinerConfig AnalyzerSpecMinerConfig
}

// DocumentWrapper keeps both structured data and markdown representation.
type DocumentWrapper[T any] struct {
	ArtifactSchemaVersion string `json:"schema_version"`
	Data                  T      `json:"data"`
	RawMarkdown           string `json:"raw_markdown"`
}

func (d DocumentWrapper[T]) SchemaVersion() string { return d.ArtifactSchemaVersion }

// SourceProjectResearch represents the research document produced by AnalyzerAgent.
type SourceProjectResearch struct {
	Overview           string   `json:"overview"`
	DirectoryStructure string   `json:"directory_structure"`
	StructsInterfaces  string   `json:"structs_and_interfaces"`
	DataModels         string   `json:"data_models"`
	ErrorHandling      string   `json:"error_handling"`
	Dependencies       []string `json:"dependencies"`
	MigrationStrategy  string   `json:"migration_strategy,omitempty"`
	StrategyRationale  string   `json:"strategy_rationale,omitempty"`
	RawDocument        string   `json:"raw_document"`
}

// ThirdPartyLibraryAnalysis represents the library analysis document produced by AnalyzerAgent.
type ThirdPartyLibraryAnalysis struct {
	Libraries   []LibraryMapping `json:"libraries"`
	RawDocument string           `json:"raw_document"`
}

// LibraryMapping details how a source library maps to a target library.
type LibraryMapping struct {
	SourceLibrary   string `json:"source_library"`
	TargetLibrary   string `json:"target_library"`
	Overview        string `json:"overview"`
	Usage           string `json:"usage"`
	Recommendations string `json:"recommendations"`
}

// TargetProjectDesign represents the design document produced by AnalyzerAgent.
type TargetProjectDesign struct {
	Overview                string   `json:"overview"`
	TranslationRequirements string   `json:"translation_requirements"`
	SourceFilesToTranslate  []string `json:"source_files_to_translate"`
	ModuleStructure         string   `json:"module_structure"`
	ErrorHandling           string   `json:"error_handling"`
	ThirdPartyLibraries     []string `json:"third_party_libraries"`
	RawDocument             string   `json:"raw_document"`
}

// AnalyzerOutput aggregates research, library mapping, and architectural design documents.
type AnalyzerOutput struct {
	ArtifactSchemaVersion string                                     `json:"schema_version"`
	Research              DocumentWrapper[SourceProjectResearch]     `json:"research"`
	Library               DocumentWrapper[ThirdPartyLibraryAnalysis] `json:"library"`
	Design                DocumentWrapper[TargetProjectDesign]       `json:"design"`
}

func (a AnalyzerOutput) SchemaVersion() string { return a.ArtifactSchemaVersion }

// AnalyzerSpecMinerConfig captures the inputs SpecMiner needs at analyzer
// time. SourceBinary is the compiled source artifact path; Inputs are the
// representative inputs the analyzer passes to SpecMiner.Recover to
// exercise the runtime. Zero value disables SpecMiner entirely.
type AnalyzerSpecMinerConfig struct {
	SourceBinary string
	Inputs       []string
}
// NewAnalyzerAgent creates an AnalyzerAgent instance without a Navigator.
// Use NewAnalyzerAgentWithNavigator to attach a Symbol-Aware Navigator.
func NewAnalyzerAgent(m model.BaseChatModel) *AnalyzerAgent {
	return &AnalyzerAgent{Model: m}
}

// NewAnalyzerAgentWithNavigator creates an AnalyzerAgent with an attached
// Symbol-Aware Navigator. Pass provider=nil to keep Navigator in disabled
// fallback mode (callers may then fall back to LLM-only symbol resolution).
func NewAnalyzerAgentWithNavigator(m model.BaseChatModel, provider tools.LSPProvider) *AnalyzerAgent {
	return &AnalyzerAgent{Model: m, Navigator: NewNavigator(provider)}
}

// Run executes the 3-phase analysis workflow and returns updated state.
func (a *AnalyzerAgent) Run(ctx context.Context, state *State) (*State, error) {
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
			state.SpecMinerInvariants = SpecMinerInvariants{
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
	treeStr, files, err := tools.BuildDirectoryTree(sourceDir, 15)
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
func (a *AnalyzerAgent) synthesizeAnalysis(ctx context.Context, state *State, treeStr, structureSummary, allCode string) (string, error) {
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
func (a *AnalyzerAgent) populateAnalyzerOutput(state *State, rawDoc, strategy, rationale string) {
	state.AnalyzerOutput.Research = DocumentWrapper[SourceProjectResearch]{
		ArtifactSchemaVersion: CurrentSchemaVersion,
		Data: SourceProjectResearch{
			Overview:           extractSection(rawDoc, "## 1. Overview", "## 2. Directory Structure"),
			DirectoryStructure: extractSection(rawDoc, "## 2. Directory Structure", "## 3. Data Structures"),
			MigrationStrategy:  strategy,
			StrategyRationale:  rationale,
		},
		RawMarkdown: rawDoc,
	}
	state.AnalyzerOutput.Library = DocumentWrapper[ThirdPartyLibraryAnalysis]{
		ArtifactSchemaVersion: CurrentSchemaVersion,
		Data: ThirdPartyLibraryAnalysis{
			Libraries: []LibraryMapping{},
		},
		RawMarkdown: extractSection(rawDoc, "=== SECTION: LIBRARY_ANALYSIS ===", "=== SECTION: TARGET_DESIGN ==="),
	}

	state.AnalyzerOutput.Design = DocumentWrapper[TargetProjectDesign]{
		ArtifactSchemaVersion: CurrentSchemaVersion,
		Data: TargetProjectDesign{
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
