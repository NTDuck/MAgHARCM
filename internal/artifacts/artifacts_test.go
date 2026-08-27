package artifacts

import (
	"os"
	"path/filepath"
	"testing"

	"MAgHARCM/internal/types"
)

func TestArtifactPersistence(t *testing.T) {
	tmpDir := t.TempDir()

	// Validate analyzer artifact persistence
	analyzerOut := types.AnalyzerOutput{
		Research: types.DocumentWrapper[types.SourceProjectResearch]{
			RawMarkdown: "# Research\nOverview",
		},
		Library: types.DocumentWrapper[types.ThirdPartyLibraryAnalysis]{
			RawMarkdown: "# Libraries\nNone",
		},
		Design: types.DocumentWrapper[types.TargetProjectDesign]{
			RawMarkdown: "# Target Design\nRust architecture",
		},
	}
	if err := SaveAnalyzerOutput(tmpDir, analyzerOut); err != nil {
		t.Fatalf("failed to save analyzer output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".MAgHARCM", "01_analyzer", "research.md")); err != nil {
		t.Errorf("expected research.md to exist: %v", err)
	}

	// Validate planning artifact persistence
	planningOut := types.PlanningOutput{
		NameMapping: map[string]string{"foo": "bar"},
		SkeletonFiles: map[string]string{
			"Cargo.toml": "[package]\nname = \"foo\"",
		},
		Plan: types.ImplementationPlan{
			RawPlan: "# Plan",
		},
	}
	if err := SavePlanningOutput(tmpDir, planningOut); err != nil {
		t.Fatalf("failed to save planning output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".MAgHARCM", "02_planning", "name_mapping.json")); err != nil {
		t.Errorf("expected name_mapping.json to exist: %v", err)
	}

	// Validate translation snapshot persistence
	transProj := types.TranslatedProject{
		Files: map[string]string{
			"src/lib.rs": "pub fn hello() {}",
		},
	}
	if err := SaveTranslationIteration(tmpDir, 1, transProj); err != nil {
		t.Fatalf("failed to save translation iteration: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".MAgHARCM", "03_translation", "iteration_01", "files_manifest.json")); err != nil {
		t.Errorf("expected files_manifest.json to exist: %v", err)
	}

	// Validate validation report persistence
	valRep := types.ValidationReport{
		AllSuccess:         true,
		CompilationSuccess: true,
		PassedTests:        3,
		TotalTests:         3,
		Diagnostics:        "All passed",
	}
	if err := SaveValidationIteration(tmpDir, 1, valRep); err != nil {
		t.Fatalf("failed to save validation iteration: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".MAgHARCM", "04_validation", "iteration_01", "validation_report.json")); err != nil {
		t.Errorf("expected validation_report.json to exist: %v", err)
	}
}
