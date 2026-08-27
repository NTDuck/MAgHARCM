package agents

import (
	"strings"
	"testing"
)

func TestRenderAllPrompts(t *testing.T) {
	// 1. analyzer
	analyzerPrompt, err := renderPromptTemplate("analyzer", analyzerPromptTemplate, map[string]any{
		"SourceLang":         "C",
		"TargetLang":         "Rust",
		"SourceDir":          "assets/samples/c",
		"DirectoryTree":      "dir tree representation",
		"StructureSummary":   "ast summary",
		"SourceFilesContent": "source code content",
	})
	if err != nil {
		t.Fatalf("failed to render analyzer template: %v", err)
	}
	for _, expected := range []string{"Source Language: C", "Target Language: Rust", "dir tree representation", "ast summary", "source code content"} {
		if !strings.Contains(analyzerPrompt, expected) {
			t.Errorf("analyzer template missing expected text: %s", expected)
		}
	}

	// 2. planning
	planningPrompt, err := renderPromptTemplate("planning", planningPromptTemplate, map[string]any{
		"SourceLang":   "C",
		"TargetLang":   "Rust",
		"SourceFiles":  "file summary content",
		"TargetDesign": "design document content",
	})
	if err != nil {
		t.Fatalf("failed to render planning template: %v", err)
	}
	for _, expected := range []string{"Source Language: C", "Target Language: Rust", "file summary content", "design document content"} {
		if !strings.Contains(planningPrompt, expected) {
			t.Errorf("planning template missing expected text: %s", expected)
		}
	}

	// 3. translator_translate
	translatePrompt, err := renderPromptTemplate("translator_translate", translatorTranslatePromptTemplate, map[string]any{
		"PackageName":        "gilded_rose",
		"SourceLang":         "C",
		"TargetLang":         "Rust",
		"TargetLangLower":    "rust",
		"SourceFiles":        "source file data",
		"TargetDesign":       "target design data",
		"ImplementationPlan": "implementation plan data",
	})
	if err != nil {
		t.Fatalf("failed to render translator_translate template: %v", err)
	}
	for _, expected := range []string{"Target Package / Module Name: gilded_rose", "source file data", "target design data", "implementation plan data"} {
		if !strings.Contains(translatePrompt, expected) {
			t.Errorf("translator_translate template missing expected text: %s", expected)
		}
	}

	// 4. translator_repair
	repairPrompt, err := renderPromptTemplate("translator_repair", translatorRepairPromptTemplate, map[string]any{
		"PackageName":     "gilded_rose",
		"TargetLang":      "Rust",
		"TargetLangLower": "rust",
		"Diagnostics":     "error: type mismatch",
		"CurrentFiles":    "current file content",
	})
	if err != nil {
		t.Fatalf("failed to render translator_repair template: %v", err)
	}
	for _, expected := range []string{"Target Package / Module Name: gilded_rose", "error: type mismatch", "current file content"} {
		if !strings.Contains(repairPrompt, expected) {
			t.Errorf("translator_repair template missing expected text: %s", expected)
		}
	}

	// 5. validator_coverage
	coveragePrompt, err := renderPromptTemplate("validator_coverage", validatorCoveragePromptTemplate, map[string]any{
		"TargetLang":         "Rust",
		"TargetLangLower":    "rust",
		"UncoveredFunctions": "fn update_quality",
		"SourceFiles":        "source files data",
		"TestFileRelPath":    "tests/integration_tests.rs",
	})
	if err != nil {
		t.Fatalf("failed to render validator_coverage template: %v", err)
	}
	for _, expected := range []string{"fn update_quality", "source files data", "tests/integration_tests.rs"} {
		if !strings.Contains(coveragePrompt, expected) {
			t.Errorf("validator_coverage template missing expected text: %s", expected)
		}
	}
}
