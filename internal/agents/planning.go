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
	state.Log("[Planning] Extracting fragments and building implementation plan")

	// Extract fine-grained AST translation units from all discovered files
	logger.LogStep("Extracting translation units across source and test files")
	_, files, err := tools.BuildDirectoryTree(state.Task.SourceDir, 5)
	if err != nil {
		return nil, err
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

	state.PlanningOutput.Fragments = fragments
	logger.LogTool("fragment_extraction", "Extracted %d translation fragments from source files", len(fragments))
	for _, fr := range fragments {
		logger.LogStep("Fragment: `%s`", fr)
	}

	// Request name mapping, project skeleton, and implementation plan from the reasoning model

	prompt, err := renderPrompt("planning.md", map[string]any{
		"SourceLang":   state.Task.SourceLang,
		"TargetLang":   state.Task.TargetLang,
		"SourceFiles":  strings.Join(sourceSummaries, "\n"),
		"TargetDesign": state.AnalyzerOutput.Design.RawMarkdown,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render planning prompt: %w", err)
	}
	resp, err := p.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert software engineer and project planner specializing in language-agnostic code translation."),
		schema.UserMessage(prompt),
	})
	if err != nil {
		return nil, fmt.Errorf("planning model call failed: %w", err)
	}

	rawContent := resp.Content

	// Parse Name Mapping
	nameMapping := make(map[string]string)
	if jsonStr := extractBlock(rawContent, "=== NAME_MAPPING_JSON ===", "==="); jsonStr != "" {
		_ = json.Unmarshal([]byte(jsonStr), &nameMapping)
	}
	state.PlanningOutput.NameMapping = nameMapping
	logger.LogTool("name_mapping", "Created %d symbol mappings", len(nameMapping))
	skeletonFiles := parseFileBlocks(rawContent, "=== SKELETON_FILES ===")
	if len(skeletonFiles) == 0 {
		logger.LogWarning("Planning LLM did not emit explicit skeleton files; generating fallback skeleton for `%s`", state.Task.TargetLang)
		skeletonFiles = defaultProjectSkeleton(state.Task.TargetDir, state.Task.TargetLang, fragments)
	}
	state.PlanningOutput.SkeletonFiles = skeletonFiles

	// Write skeleton files directly to target directory on filesystem
	for relPath, content := range skeletonFiles {
		fullPath := filepath.Join(state.Task.TargetDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for %s: %w", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write skeleton file %s: %w", fullPath, err)
		}
		logger.LogTool("write_file", "Wrote skeleton to `%s` (%d bytes)", relPath, len(content))
	}

	// Parse Plan
	planStr := extractBlock(rawContent, "=== IMPLEMENTATION_PLAN ===", "")
	if planStr == "" {
		planStr = rawContent
	}
	state.PlanningOutput.Plan = types.ImplementationPlan{
		Overview: extractSection(planStr, "## Overview", "## Part A"),
		PartA: []types.PlanStep{
			{ID: "A1", Description: "Translate all source modules to target language", Type: "source"},
		},
		PartB: []types.PlanStep{
			{ID: "B1", Description: "Translate and execute test suite", Type: "test"},
		},
		RawPlan: planStr,
	}
	logger.LogAgent("Planning", "Planning complete: %d skeleton files written to filesystem, implementation plan ready", len(skeletonFiles))
	state.Log("[Planning] Skeleton generated with %d files, Plan ready with %d fragments", len(skeletonFiles), len(fragments))
	return state, nil
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
		return map[string]string{
			"src/main." + strings.ToLower(targetLang): "// Main entry point\n",
		}
	}

	files := make(map[string]string)
	for relPath, tmpl := range spec.DefaultSkeleton {
		content := strings.ReplaceAll(tmpl, "{{.ProjectName}}", projectName)
		files[relPath] = content
	}

	if len(files) == 0 {
		ext := ".txt"
		if len(spec.Extensions) > 0 {
			ext = spec.Extensions[0]
		}
		files["src/lib"+ext] = "// Generated project module\n"
	}

	return files
}

func sanitizeProjectName(name string) string {
	name = strings.ToLower(name)
	var sb strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			sb.WriteRune(r)
		} else if r == '-' {
			sb.WriteRune('_')
		}
	}
	res := strings.Trim(sb.String(), "_")
	if res == "" {
		return "translated_project"
	}
	return res
}

func extractBlock(doc, startTag, endTag string) string {
	start := strings.Index(doc, startTag)
	if start == -1 {
		return ""
	}
	sub := doc[start+len(startTag):]
	if endTag != "" {
		end := strings.Index(sub, endTag)
		if end != -1 {
			return strings.TrimSpace(sub[:end])
		}
	}
	return strings.TrimSpace(sub)
}

func parseFileBlocks(text, sectionHeader string) map[string]string {
	files := make(map[string]string)
	content := text
	if sectionHeader != "" {
		secIdx := strings.Index(text, sectionHeader)
		if secIdx == -1 {
			return files
		}
		content = text[secIdx+len(sectionHeader):]
		if nextSec := strings.Index(content, "==="); nextSec != -1 {
			content = content[:nextSec]
		}
	}

	lines := strings.Split(content, "\n")
	var currentFile string
	var currentLines []string

	flush := func() {
		if currentFile != "" && len(currentLines) > 0 {
			code := tools.CleanCodeContent(strings.Join(currentLines, "\n"))
			if code != "\n" && code != "" {
				files[currentFile] = code
			}
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "FILE:") || strings.HasPrefix(trimmed, "### File:") {
			flush()
			f := strings.TrimPrefix(trimmed, "FILE:")
			f = strings.TrimPrefix(f, "### File:")
			f = strings.TrimSpace(f)
			currentFile = f
			currentLines = nil
		} else if currentFile != "" {
			currentLines = append(currentLines, line)
		}
	}
	flush()
	return files
}
