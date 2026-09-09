package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"MAgHARCM/internal/compiletime"
	"github.com/cloudwego/eino/components/model"
	"MAgHARCM/internal/llm"

	"MAgHARCM/internal/languages"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
)

// Type aliases (cycle-free Locality of Behaviour).
// Producer files declare the algorithms; canonical artifact types
// live in internal/compiletime/state.go. Aliases let method receivers
// reference the type name without package qualification.
type ImplementationPlan = compiletime.ImplementationPlan
type PlanStep = compiletime.PlanStep
type PlanningOutput = compiletime.PlanningOutput

// PlanningAgent extracts translation units, maps symbols to target conventions, generates project skeletons, and devises execution plans.
type PlanningAgent struct {
	Model model.BaseChatModel
	// structured binds PlanningSchema -> emit_plan tool. The model is forced
	// to call this single tool, replacing the brittle "extract a JSON block
	// from a fenced code section" path the planner used to rely on.
	structured *llm.StructuredExtractor[PlanningSchema]
}

// NewPlanningAgent wires the planner to a chat model and prepares its typed-output
// schema (PlanningSchema -> emit_plan tool). Returns an error when the model
// cannot be used as a ToolCallingChatModel.
func NewPlanningAgent(m model.BaseChatModel) (*PlanningAgent, error) {
	extractor, err := llm.NewStructuredExtractor[PlanningSchema](m, "emit_plan",
		"Emit the implementation plan: source -> target symbol name map, skeleton files keyed by relative path, and the implementation plan overview / Part A (source) / Part B (test) steps.")
	if err != nil {
		return nil, fmt.Errorf("planning: build structured extractor: %w", err)
	}
	return &PlanningAgent{Model: m, structured: extractor}, nil
}
// compiletime.PlanStep represents a single step in Part A or Part B of the
// implementation plan. Producer: PlanningAgent (this file).
// Run executes the planning phase and populates compiletime.PlanningOutput in state.
func (p *PlanningAgent) Run(ctx context.Context, state *compiletime.State) (*compiletime.State, error) {
	logger.LogAgent("Planning", "Decomposing translation into granular translation units and constructing plan")
	state.PlanningOutput.ArtifactSchemaVersion = compiletime.CurrentSchemaVersion

	fragments, sourceSummaries, err := p.extractFragments(state.Task.SourceDir)
	if err != nil {
		return nil, err
	}
	state.PlanningOutput.Fragments = fragments

	plan, err := p.generatePlanningArtifacts(ctx, state, sourceSummaries, fragments)
	if err != nil {
		return nil, err
	}

	state.PlanningOutput.NameMapping = plan.NameMapping
	skeletonFiles := p.resolveSkeletonFiles(plan, state, fragments)
	state.PlanningOutput.SkeletonFiles = skeletonFiles

	if err := p.writeSkeletonFiles(state.Task.TargetDir, skeletonFiles); err != nil {
		return nil, err
	}

	state.PlanningOutput.Plan = p.buildImplementationPlan(plan, fragments)
	logger.LogAgent("Planning", "Planning complete: %d skeleton files written, %d steps in reverse-topological order",
		len(skeletonFiles), len(state.PlanningOutput.Plan.PartA))
	return state, nil
}

// extractFragments scans source files and extracts AST translation fragments and file summaries.
//
// When a source file fails to parse OR parses to zero AST elements (common for
// build/config/prose files: pom.xml, README.md, *.properties, etc.), we still
// emit a single per-file fragment "<base>:file" so the chunked translator sees
// the file in its dispatch loop. Otherwise the chunked translator silently
// drops 100+ files on Java corpora where most files are XML/Markdown/properties.
func (p *PlanningAgent) extractFragments(sourceDir string) ([]string, []string, error) {
	logger.LogStep("Extracting translation units across source and test files")
	_, files, err := tools.BuildDirectoryTree(sourceDir, compiletime.DefaultSourceTreeDepth)
	if err != nil {
		return nil, nil, err
	}

	var fragments []string
	var sourceSummaries []string
	var totalSummaryBytes int
	const maxPlannerSummaryBudget = 32 * 1024

	for _, f := range files {
		if !isTranslatableFile(f) {
			continue
		}
		base := filepath.Base(f)
		structOut, parseErr := tools.ParseFileStructure(f)
		if parseErr != nil {
			logger.LogWarning("ParseFileStructure failed for %q; emitting file-level fallback fragment", base)
			fragments = append(fragments, fmt.Sprintf("%s:file", base))
			if totalSummaryBytes < maxPlannerSummaryBudget {
				raw, readErr := os.ReadFile(f)
				if readErr == nil {
					chunk := string(raw)
					if len(chunk) > 2048 {
						chunk = chunk[:2048] + "\n// ... (truncated)"
					}
					sourceSummaries = append(sourceSummaries, fmt.Sprintf("File %s:\n%s\n", base, chunk))
					totalSummaryBytes += len(chunk)
				} else {
					sourceSummaries = append(sourceSummaries, fmt.Sprintf("File %s: <unreadable>\n", base))
				}
			}
			continue
		}
		if len(structOut.Elements) == 0 {
			fragments = append(fragments, fmt.Sprintf("%s:file", base))
			if totalSummaryBytes < maxPlannerSummaryBudget {
				chunk := structOut.RawCode
				if len(chunk) > 2048 {
					chunk = chunk[:2048] + "\n// ... (truncated)"
				}
				sourceSummaries = append(sourceSummaries, fmt.Sprintf("File %s:\n%s\n", base, chunk))
				totalSummaryBytes += len(chunk)
			}
			continue
		}
		for _, el := range structOut.Elements {
			frag := fmt.Sprintf("%s:%s", base, el.Name)
			fragments = append(fragments, frag)
		}
		if totalSummaryBytes < maxPlannerSummaryBudget {
			chunk := structOut.RawCode
			if len(chunk) > 2048 {
				chunk = chunk[:2048] + "\n// ... (truncated)"
			}
			sourceSummaries = append(sourceSummaries, fmt.Sprintf("File %s:\n%s\n", base, chunk))
			totalSummaryBytes += len(chunk)
		}
	}
	logger.LogTool("fragment_extraction", "Extracted %d translation fragments from source files", len(fragments))
	for _, fr := range fragments {
		logger.LogStep("Fragment: `%s`", fr)
	}
	return fragments, sourceSummaries, nil
}

// generatePlanningArtifacts queries the reasoning model via the structured
// emit_plan tool. The model returns the symbol name map, skeleton file map,
// and implementation plan as a single typed tool call — no string parsing
// required.
func (p *PlanningAgent) generatePlanningArtifacts(ctx context.Context, state *compiletime.State, sourceSummaries []string, fragments []string) (PlanningSchema, error) {
	if p.structured == nil {
		return PlanningSchema{}, fmt.Errorf("planning: structured extractor not initialised (call NewPlanningAgent with a ToolCallingChatModel)")
	}
	prompt, err := renderPromptTemplate("planning", planningPromptTemplate, map[string]any{
		"SourceLang":   state.Task.SourceLang,
		"TargetLang":   state.Task.TargetLang,
		"SourceFiles":  strings.Join(sourceSummaries, "\n"),
		"TargetDesign": state.AnalyzerOutput.Design.RawMarkdown,
		"Fragments":    strings.Join(fragments, "\n"),
	})
	if err != nil {
		return PlanningSchema{}, fmt.Errorf("failed to render planning prompt: %w", err)
	}
	system := "You are an expert software engineer and project planner specializing in language-agnostic code translation. Always respond by calling the emit_plan tool with the structured fields; do not emit any other text."
	return p.structured.Extract(ctx, system, prompt)
}

// resolveSkeletonFiles selects skeleton files from the typed plan, falling
// back to language-specific boilerplate when the LLM did not emit any.
func (p *PlanningAgent) resolveSkeletonFiles(plan PlanningSchema, state *compiletime.State, fragments []string) map[string]string {
	skeletonFiles := plan.SkeletonFiles
	if len(skeletonFiles) == 0 {
		logger.LogWarning("Planning LLM did not emit explicit skeleton files; generating fallback skeleton for `%s`", state.Task.TargetLang)
		skeletonFiles = DefaultProjectSkeleton(state.Task.TargetDir, state.Task.TargetLang, fragments)
	} else if strings.EqualFold(state.Task.TargetLang, "Rust") {
		if _, hasCargo := skeletonFiles["Cargo.toml"]; !hasCargo {
			projectName := resolvePackageName(state.Task.TargetDir, state.Task.TargetLang)
			skeletonFiles["Cargo.toml"] = fmt.Sprintf("[package]\nname = \"%s\"\nversion = \"0.1.0\"\nedition = \"2021\"\n\n[dependencies]\n", projectName)
		}
	}
	if plan.NameMapping != nil {
		logger.LogTool("name_mapping", "Created %d symbol mappings", len(plan.NameMapping))
	}
	return skeletonFiles
}

// writeSkeletonFiles persists all initial project scaffolding to disk.
func (p *PlanningAgent) writeSkeletonFiles(targetDir string, skeletonFiles map[string]string) error {
	for relPath, content := range skeletonFiles {
		fullPath := filepath.Join(targetDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write skeleton file %s: %w", fullPath, err)
		}
		logger.LogTool("write_file", "Wrote skeleton to `%s`, %d bytes", relPath, len(content))
	}
	return nil
}

// buildImplementationPlan merges the LLM-emitted Part A / Part B steps with
// the deterministic reverse-topological fragment order. When the LLM
// emitted no Part A steps, the planner falls back to the fragment list so
// downstream phases always have work to dispatch.
func (p *PlanningAgent) buildImplementationPlan(plan PlanningSchema, fragments []string) compiletime.ImplementationPlan {
	partASteps := plan.PartA
	if len(partASteps) == 0 {
		orderedFrags := ComputeReverseTopoOrder(fragments, nil)
		if len(orderedFrags) > 0 {
			for i, frag := range orderedFrags {
				partASteps = append(partASteps, compiletime.PlanStep{
					ID:              fmt.Sprintf("A%d", i+1),
					Description:     fmt.Sprintf("Translate module fragment: %s", frag),
					Type:            "source",
					StepName:        frag,
					ReverseTopoRank: i + 1,
				})
			}
		} else {
			partASteps = []compiletime.PlanStep{
				{ID: "A1", Description: "Translate all source modules to target language", Type: "source", ReverseTopoRank: 1},
			}
		}
	}
	partBSteps := plan.PartB
	if len(partBSteps) == 0 {
		partBSteps = []compiletime.PlanStep{
			{ID: "B1", Description: "Translate and execute test suite", Type: "test", ReverseTopoRank: len(partASteps) + 1},
		}
	}
	return compiletime.ImplementationPlan{
		ArtifactSchemaVersion: compiletime.CurrentSchemaVersion,
		Overview:              plan.Overview,
		PartA:                 partASteps,
		PartB:                 partBSteps,
		RawPlan:               plan.Overview,
	}
}

// ComputeReverseTopoOrder builds a reverse topological ordering (dependencies/leaves first)
// with back-edge removal (AlphaTrans NEW-PRIM-1 + NEW-PRIM-2 / GAP-02) to break cycles.
func ComputeReverseTopoOrder(items []string, dependencies map[string][]string) []string {
	if dependencies == nil {
		dependencies = make(map[string][]string)
	}
	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	var ordered []string

	var dfs func(node string)
	dfs = func(node string) {
		visited[node] = true
		inStack[node] = true

		for _, dep := range dependencies[node] {
			if inStack[dep] {
				// Back-edge detected (cycle): remove/skip back-edge
				continue
			}
			if !visited[dep] {
				dfs(dep)
			}
		}

		inStack[node] = false
		ordered = append(ordered, node)
	}

	for _, item := range items {
		if !visited[item] {
			dfs(item)
		}
	}

	return ordered
}

// defaultProjectSkeleton dynamically generates standard project boilerplate files for the target language.
func DefaultProjectSkeleton(targetDir string, targetLang string, fragments []string) map[string]string {
	projectName := resolvePackageName(targetDir, targetLang)

	registry := languages.GetRegistry()
	spec, found := registry.FindByName(targetLang)
	if !found {
		spec, found = registry.FindByExtension("." + strings.ToLower(targetLang))
	}

	files := make(map[string]string)
	if found && spec.DefaultSkeleton != nil {
		for relPath, tmpl := range spec.DefaultSkeleton {
			content := strings.ReplaceAll(tmpl, "{{.ProjectName}}", projectName)
			files[relPath] = content
		}
	}

	if len(files) == 0 {
		files["src/lib.rs"] = "// Auto-generated boilerplate\n"
		files["Cargo.toml"] = fmt.Sprintf("[package]\nname = \"%s\"\nversion = \"0.1.0\"\nedition = \"2021\"\n\n[dependencies]\n", projectName)
	}

	return files
}

func sanitizeProjectName(name string) string {
	name = strings.ToLower(name)
	var b strings.Builder
	for i, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			b.WriteRune(r)
		} else if r == '-' || r == ' ' || r == '.' {
			if i > 0 && b.Len() > 0 {
				b.WriteRune('_')
			}
		}
	}
	res := strings.Trim(b.String(), "_")
	if res == "" || unicode.IsDigit(rune(res[0])) {
		res = "p_" + res
	}
	return res
}



// isTranslatableFile determines if a file is a source, test, or build manifest
// eligible for AST fragment extraction.
func isTranslatableFile(path string) bool {
	base := filepath.Base(path)
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".go", ".c", ".h", ".cc", ".cpp", ".cxx", ".hpp", ".java", ".rs", ".py", ".ts", ".js", ".kt", ".scala", ".cs":
		return true
	}
	switch strings.ToLower(base) {
	case "makefile", "go.mod", "pom.xml", "build.gradle", "cargo.toml":
		return true
	}
	return false
}
