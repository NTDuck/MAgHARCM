package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/languages"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
	"MAgHARCM/internal/types"
)

// PlanningAgent extracts translation units, maps symbols to target conventions, generates project skeletons, and devises execution plans.
type PlanningAgent struct {
	Model model.BaseChatModel
}

// NewPlanningAgent creates a PlanningAgent instance.
func NewPlanningAgent(m model.BaseChatModel) *PlanningAgent {
	return &PlanningAgent{Model: m}
}

// Run executes the planning phase and populates PlanningOutput in state.
func (p *PlanningAgent) Run(ctx context.Context, state *types.State) (*types.State, error) {
	logger.LogAgent("Planning", "Decomposing translation into granular translation units and constructing plan")

	fragments, sourceSummaries, err := p.extractFragments(state.Task.SourceDir)
	if err != nil {
		return nil, err
	}
	state.PlanningOutput.Fragments = fragments

	rawContent, err := p.generatePlanningArtifacts(ctx, state, sourceSummaries)
	if err != nil {
		return nil, err
	}

	state.PlanningOutput.NameMapping = p.parseNameMapping(rawContent)
	skeletonFiles := p.resolveSkeletonFiles(rawContent, state, fragments)
	state.PlanningOutput.SkeletonFiles = skeletonFiles

	if err := p.writeSkeletonFiles(state.Task.TargetDir, skeletonFiles); err != nil {
		return nil, err
	}

	state.PlanningOutput.Plan = p.parseImplementationPlan(rawContent, fragments)
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
	_, files, err := tools.BuildDirectoryTree(sourceDir, 15)
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

// generatePlanningArtifacts queries the reasoning model for name mapping, skeleton, and implementation plan.
func (p *PlanningAgent) generatePlanningArtifacts(ctx context.Context, state *types.State, sourceSummaries []string) (string, error) {
	prompt, err := renderPromptTemplate("planning", planningPromptTemplate, map[string]any{
		"SourceLang":   state.Task.SourceLang,
		"TargetLang":   state.Task.TargetLang,
		"SourceFiles":  strings.Join(sourceSummaries, "\n"),
		"TargetDesign": state.AnalyzerOutput.Design.RawMarkdown,
	})
	if err != nil {
		return "", fmt.Errorf("failed to render planning prompt: %w", err)
	}

	resp, err := p.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert software engineer and project planner specializing in language-agnostic code translation."),
		schema.UserMessage(prompt),
	})
	if err != nil {
		return "", fmt.Errorf("planning model call failed: %w", err)
	}
	return resp.Content, nil
}

// parseNameMapping extracts and decodes the JSON symbol name mapping from model output.
func (p *PlanningAgent) parseNameMapping(rawContent string) map[string]string {
	nameMapping := make(map[string]string)
	if jsonStr := extractBlock(rawContent, "=== NAME_MAPPING_JSON ===", "==="); jsonStr != "" {
		_ = json.Unmarshal([]byte(jsonStr), &nameMapping)
	}
	logger.LogTool("name_mapping", "Created %d symbol mappings", len(nameMapping))
	return nameMapping
}

// resolveSkeletonFiles extracts skeleton files or synthesizes default boilerplate for the target language.
func (p *PlanningAgent) resolveSkeletonFiles(rawContent string, state *types.State, fragments []string) map[string]string {
	skeletonFiles := parseFileBlocks(rawContent, "=== SKELETON_FILES ===")
	if len(skeletonFiles) == 0 {
		logger.LogWarning("Planning LLM did not emit explicit skeleton files; generating fallback skeleton for `%s`", state.Task.TargetLang)
		skeletonFiles = defaultProjectSkeleton(state.Task.TargetDir, state.Task.TargetLang, fragments)
	} else if strings.EqualFold(state.Task.TargetLang, "Rust") {
		if _, hasCargo := skeletonFiles["Cargo.toml"]; !hasCargo {
			projectName := sanitizeProjectName(filepath.Base(state.Task.TargetDir))
			if projectName == "" || projectName == "." || strings.EqualFold(projectName, state.Task.TargetLang) {
				parent := filepath.Base(filepath.Dir(state.Task.TargetDir))
				if parent != "" && parent != "." && parent != "/" {
					projectName = sanitizeProjectName(parent)
				}
			}
			if projectName == "" {
				projectName = "translated_project"
			}
			skeletonFiles["Cargo.toml"] = fmt.Sprintf("[package]\nname = \"%s\"\nversion = \"0.1.0\"\nedition = \"2021\"\n\n[dependencies]\n", projectName)
		}
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
		logger.LogTool("write_file", "Wrote skeleton to `%s` (%d bytes)", relPath, len(content))
	}
	return nil
}

// parseImplementationPlan parses the implementation plan sections into structured steps
// and schedules them in reverse topological order (NEW-PRIM-1, NEW-PRIM-2 / GAP-02).
func (p *PlanningAgent) parseImplementationPlan(rawContent string, fragments []string) types.ImplementationPlan {
	planStr := extractBlock(rawContent, "=== IMPLEMENTATION_PLAN ===", "")
	if planStr == "" {
		planStr = rawContent
	}

	orderedFrags := ComputeReverseTopoOrder(fragments, nil)
	var partASteps []types.PlanStep
	if len(orderedFrags) > 0 {
		for i, frag := range orderedFrags {
			partASteps = append(partASteps, types.PlanStep{
				ID:              fmt.Sprintf("A%d", i+1),
				Description:     fmt.Sprintf("Translate module fragment: %s", frag),
				Type:            "source",
				StepName:        frag,
				ReverseTopoRank: i + 1,
			})
		}
	} else {
		partASteps = []types.PlanStep{
			{ID: "A1", Description: "Translate all source modules to target language", Type: "source", ReverseTopoRank: 1},
		}
	}

	return types.ImplementationPlan{
		Overview: extractSection(planStr, "## Overview", "## Part A"),
		PartA:    partASteps,
		PartB: []types.PlanStep{
			{ID: "B1", Description: "Translate and execute test suite", Type: "test", ReverseTopoRank: len(partASteps) + 1},
		},
		RawPlan: planStr,
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
func defaultProjectSkeleton(targetDir string, targetLang string, fragments []string) map[string]string {
	projectName := sanitizeProjectName(filepath.Base(targetDir))
	if projectName == "" || projectName == "." || strings.EqualFold(projectName, targetLang) {
		parent := filepath.Base(filepath.Dir(targetDir))
		if parent != "" && parent != "." && parent != "/" {
			projectName = sanitizeProjectName(parent)
		}
	}
	if projectName == "" {
		projectName = "translated_project"
	}

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

func extractBlock(doc, startTag, endTag string) string {
	startIdx := strings.Index(doc, startTag)
	if startIdx == -1 {
		return ""
	}
	content := doc[startIdx+len(startTag):]
	if endTag != "" {
		endIdx := strings.Index(content, endTag)
		if endIdx != -1 {
			content = content[:endIdx]
		}
	}
	return strings.TrimSpace(content)
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

