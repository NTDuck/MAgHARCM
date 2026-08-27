package artifacts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/types"
)

const ArtifactsSubdir = ".MAgHARCM"

// GetArtifactsDir returns the base path for interim artifacts in the target directory.
func GetArtifactsDir(targetDir string) string {
	return filepath.Join(targetDir, ArtifactsSubdir)
}

// SaveAnalyzerOutput writes Phase 1 interim artifacts.
func SaveAnalyzerOutput(targetDir string, out types.AnalyzerOutput) error {
	dir := filepath.Join(GetArtifactsDir(targetDir), "01_analyzer")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	_ = os.WriteFile(filepath.Join(dir, "research.md"), []byte(out.Research.RawMarkdown), 0644)
	_ = os.WriteFile(filepath.Join(dir, "libraries.md"), []byte(out.Library.RawMarkdown), 0644)
	_ = os.WriteFile(filepath.Join(dir, "target_design.md"), []byte(out.Design.RawMarkdown), 0644)

	rawJSON, _ := json.MarshalIndent(out, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "analyzer_output.json"), rawJSON, 0644)

	logger.LogArtifact("Saved Analyzer artifacts to %s", dir)
	return nil
}

// SavePlanningOutput writes Phase 2 interim artifacts.
func SavePlanningOutput(targetDir string, out types.PlanningOutput) error {
	dir := filepath.Join(GetArtifactsDir(targetDir), "02_planning")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	mapJSON, _ := json.MarshalIndent(out.NameMapping, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "name_mapping.json"), mapJSON, 0644)

	_ = os.WriteFile(filepath.Join(dir, "implementation_plan.md"), []byte(out.Plan.RawPlan), 0644)

	skelManifest, _ := json.MarshalIndent(out.SkeletonFiles, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "skeleton_manifest.json"), skelManifest, 0644)

	logger.LogArtifact("Saved Planning artifacts to %s", dir)
	return nil
}

// SaveTranslationIteration writes Phase 3 interim artifacts for a specific iteration.
func SaveTranslationIteration(targetDir string, iteration int, project types.TranslatedProject) error {
	dir := filepath.Join(GetArtifactsDir(targetDir), "03_translation", fmt.Sprintf("iteration_%02d", iteration))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	for relPath, content := range project.Files {
		dest := filepath.Join(dir, "files", relPath)
		_ = os.MkdirAll(filepath.Dir(dest), 0755)
		_ = os.WriteFile(dest, []byte(content), 0644)
	}

	manifest, _ := json.MarshalIndent(project.Files, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "files_manifest.json"), manifest, 0644)

	logger.LogArtifact("Saved Translation iteration %d snapshot to %s", iteration, dir)
	return nil
}

// SaveValidationIteration writes Phase 4 interim artifacts for a specific iteration.
func SaveValidationIteration(targetDir string, iteration int, report types.ValidationReport) error {
	dir := filepath.Join(GetArtifactsDir(targetDir), "04_validation", fmt.Sprintf("iteration_%02d", iteration))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	repJSON, _ := json.MarshalIndent(report, "", "  ")
	_ = os.WriteFile(filepath.Join(dir, "validation_report.json"), repJSON, 0644)

	if report.Diagnostics != "" {
		_ = os.WriteFile(filepath.Join(dir, "diagnostics.txt"), []byte(report.Diagnostics), 0644)
	}

	logger.LogArtifact("Saved Validation iteration %d report to %s", iteration, dir)
	return nil
}
