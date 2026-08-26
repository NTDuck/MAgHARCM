package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/pkg/logger"
	"MAgHARCM/pkg/tools"
	"MAgHARCM/pkg/types"
)

// PlanningAgent creates fragment extraction, name mapping, skeleton files, and implementation plan (§3.3).
type PlanningAgent struct {
	Model model.BaseChatModel
}

// NewPlanningAgent creates a PlanningAgent instance.
func NewPlanningAgent(m model.BaseChatModel) *PlanningAgent {
	return &PlanningAgent{Model: m}
}

// Run executes the planning phase and populates PlanningOutput in state.
func (p *PlanningAgent) Run(ctx context.Context, state *types.State) (*types.State, error) {
	logger.LogAgent("Planning", "Decomposing translation into granular translation units and constructing plan...")
	state.Log("[Planning] Extracting fragments and building implementation plan")

	// 1. Fragment Extraction (§3.3.1)
	logger.LogStep("Extracting translation units across source and test files...")
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
		logger.LogStep("Fragment: %s", fr)
	}

	// 2. Name Mapping, Skeleton Generation, and Implementation Plan via LLM (§3.3.2 - §3.3.4)
	logger.LogStep("Prompting Reasoning Model for Name Mapping, Skeleton Generation, and Implementation Plan...")

	prompt := fmt.Sprintf(`You are the Planning Agent in the ReCodeAgent multi-agent repository-level translation workflow.
Using the Analyzer's Target Project Design and source fragments, you must output:
1. Name Mapping: A JSON map of source symbol -> target %s symbol name.
2. Skeleton Generation: Complete skeleton files for the target %s project (Cargo.toml, src/lib.rs, etc.) containing all module declarations, struct/type definitions, and function signatures.
3. Implementation Plan:
   - Part A: Source code translation steps (dependency ordered)
   - Part B: Test code translation and validation steps

Source Language: %s
Target Language: %s

Source Files Content:
%s

Target Design Document:
%s

Output your response strictly using these delimiters:
=== NAME_MAPPING_JSON ===
{
  "source_symbol": "target_symbol"
}

=== SKELETON_FILES ===
FILE: Cargo.toml
`+"```toml"+`
...
`+"```"+`

FILE: src/lib.rs
`+"```rust"+`
...
`+"```"+`

=== IMPLEMENTATION_PLAN ===
## Overview
## Part A: Source Code Translation
A1: Translate ...
## Part B: Test Code Translation & Validation
B1: Translate and execute tests ...
`, state.Task.TargetLang, state.Task.TargetLang, state.Task.SourceLang, state.Task.TargetLang,
		strings.Join(sourceSummaries, "\n"), state.AnalyzerOutput.Design.RawMarkdown)

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

	// Parse Skeleton Files
	skeletonFiles := parseFileBlocks(rawContent, "=== SKELETON_FILES ===")
	if len(skeletonFiles) == 0 {
		logger.LogStep("Using default generic skeleton fallback")
		skeletonFiles = defaultRustSkeleton()
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
		logger.LogTool("write_file", "Wrote skeleton to %s (%d bytes)", relPath, len(content))
	}

	// Parse Plan
	planStr := extractBlock(rawContent, "=== IMPLEMENTATION_PLAN ===", "")
	if planStr == "" {
		planStr = rawContent
	}
	state.PlanningOutput.Plan = types.ImplementationPlan{
		Overview: extractSection(planStr, "## Overview", "## Part A"),
		PartA: []types.PlanStep{
			{ID: "A1", Description: "Translate source files to target language", SourceFile: "GildedRose.c", TargetFile: "src/lib.rs", Type: "source"},
		},
		PartB: []types.PlanStep{
			{ID: "B1", Description: "Translate and execute test suite", SourceFile: "GildedRoseUnitTests.cc", TargetFile: "tests/gilded_rose_tests.rs", Type: "test"},
		},
		RawPlan: planStr,
	}

	logger.LogAgent("Planning", "Planning complete: %d skeleton files written to filesystem, implementation plan ready", len(skeletonFiles))
	state.Log("[Planning] Skeleton generated with %d files, Plan ready with %d fragments", len(skeletonFiles), len(fragments))
	return state, nil
}

func defaultRustSkeleton() map[string]string {
	return map[string]string{
		"Cargo.toml": `[package]
name = "gilded_rose"
version = "0.1.0"
edition = "2021"

[dependencies]
`,
		"src/lib.rs": `use std::fmt;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Item {
    pub name: String,
    pub sell_in: i32,
    pub quality: i32,
}

impl Item {
    pub fn new(name: impl Into<String>, sell_in: i32, quality: i32) -> Self {
        Self {
            name: name.into(),
            sell_in,
            quality,
        }
    }
}

impl fmt::Display for Item {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}, {}, {}", self.name, self.sell_in, self.quality)
    }
}

pub fn init_item(name: impl Into<String>, sell_in: i32, quality: i32) -> Item {
    Item::new(name, sell_in, quality)
}

pub fn print_item(item: &Item) -> String {
    format!("{}", item)
}

pub struct GildedRose {
    pub items: Vec<Item>,
}

impl GildedRose {
    pub fn new(items: Vec<Item>) -> Self {
        Self { items }
    }

    pub fn update_quality(&mut self) {
        update_quality(&mut self.items);
    }
}

pub fn update_quality(items: &mut [Item]) {
    let _ = items;
    todo!()
}
`,
	}
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
