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

	state.PlanningOutput.Plan = p.parseImplementationPlan(rawContent)
	logger.LogAgent("Planning", "Planning complete: %d skeleton files written to filesystem, implementation plan ready", len(skeletonFiles))
	return state, nil
}

// extractFragments scans source files and extracts AST translation fragments and file summaries.
func (p *PlanningAgent) extractFragments(sourceDir string) ([]string, []string, error) {
	logger.LogStep("Extracting translation units across source and test files")
	_, files, err := tools.BuildDirectoryTree(sourceDir, 5)
	if err != nil {
		return nil, nil, err
	}

	var fragments []string
	var sourceSummaries []string
	for _, f := range files {
		structOut, err := tools.ParseFileStructure(f)
		if err == nil {
			base := filepath.Base(f)
			for _, el := range structOut.Elements {
				frag := fmt.Sprintf("%s:%s", base, el.Name)
				fragments = append(fragments, frag)
			}
			sourceSummaries = append(sourceSummaries, fmt.Sprintf("File %s:\n%s\n", base, structOut.RawCode))
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

// parseImplementationPlan parses the implementation plan sections into structured steps.
func (p *PlanningAgent) parseImplementationPlan(rawContent string) types.ImplementationPlan {
	planStr := extractBlock(rawContent, "=== IMPLEMENTATION_PLAN ===", "")
	if planStr == "" {
		planStr = rawContent
	}
	return types.ImplementationPlan{
		Overview: extractSection(planStr, "## Overview", "## Part A"),
		PartA: []types.PlanStep{
			{ID: "A1", Description: "Translate all source modules to target language", Type: "source"},
		},
		PartB: []types.PlanStep{
			{ID: "B1", Description: "Translate and execute test suite", Type: "test"},
		},
		RawPlan: planStr,
	}
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

func parseFileBlocks(text, sectionHeader string) map[string]string {
	files := make(map[string]string)
	sectionText := text
	if sectionHeader != "" {
		idx := strings.Index(text, sectionHeader)
		if idx != -1 {
			sectionText = text[idx+len(sectionHeader):]
		}
	}

	lines := strings.Split(sectionText, "\n")
	var currentFile string
	var currentContent strings.Builder
	inFence := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "File: ") || strings.HasPrefix(trimmed, "FILE: ") || strings.HasPrefix(trimmed, "Path: ") {
			if currentFile != "" {
				files[currentFile] = extractCodeFromSection(currentContent.String())
				currentContent.Reset()
			}
			currentFile = strings.TrimSpace(trimmed[6:])
			inFence = false
			continue
		}

		if strings.HasPrefix(trimmed, "```") {
			inFence = !inFence
			continue
		}

		if currentFile != "" {
			currentContent.WriteString(line + "\n")
		}
	}

	if currentFile != "" {
		files[currentFile] = extractCodeFromSection(currentContent.String())
	}

	return files
}
