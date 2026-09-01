package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
	"MAgHARCM/internal/types"
)

// TranslatorAgent generates target source and test implementations, and iteratively resolves compiler errors in repair mode.
type TranslatorAgent struct {
	Model model.BaseChatModel
	// RunID identifies the current translation run; when set, a checkpoint of
	// *types.State is persisted under .artifacts/<RunID>/checkpoints/ at the
	// end of every Run (both success and error paths). Empty RunID disables
	// checkpointing — used by callers that don't want disk persistence.
	RunID string
}

// NewTranslatorAgent creates a TranslatorAgent instance. runID enables
// per-run checkpoint persistence; pass "" to disable checkpointing.
func NewTranslatorAgent(m model.BaseChatModel, runID string) *TranslatorAgent {
	return &TranslatorAgent{Model: m, RunID: runID}
}

// Run translates the code and tests, or executes repairs if validation report has failures.
func (t *TranslatorAgent) Run(ctx context.Context, state *types.State) (*types.State, error) {
	defer t.checkpoint(state)
	if state.TranslatedProject.Files == nil {
		state.TranslatedProject.Files = make(map[string]string)
	}
	if len(state.TranslatedProject.Files) == 0 && len(state.PlanningOutput.SkeletonFiles) > 0 {
		for k, v := range state.PlanningOutput.SkeletonFiles {
			state.TranslatedProject.Files[k] = v
		}
	}

	if state.Iteration > 0 && !state.ValidationReport.IsAllSuccess() {
		logger.LogAgent("Translator", "Repair Mode (Iteration %d/%d): Diagnosing validation errors and fixing code",
			state.Iteration, state.MaxIterations)
		return t.repair(ctx, state)
	}

	if shouldUseChunkedTranslation(state) {
		logger.LogAgent("Translator", "Chunked translation mode: %d fragments, source LoC > %d", len(state.PlanningOutput.Fragments), chunkedLoCThreshold)
		if _, err := t.RunChunked(ctx, state); err != nil {
			return nil, err
		}
		return state, nil
	}

	logger.LogAgent("Translator", "Initial Translation Mode: Implementing Part A (Source) and Part B (Tests)")
	return t.translate(ctx, state)
}

// checkpoint persists a snapshot of state if RunID is set. Errors are logged
// but never propagated — checkpointing is best-effort and must not abort the
// pipeline when the disk is full or the runID is empty.
func (t *TranslatorAgent) checkpoint(state *types.State) {
	if t.RunID == "" || state == nil {
		return
	}
	if path, err := Save(t.RunID, state); err != nil {
		logger.LogWarning("Translator checkpoint save failed: %v", err)
	} else {
		logger.LogStep("Translator checkpoint saved: %s", path)
	}
}

// translate generates the initial translation from source modules, design, and implementation plan.
func (t *TranslatorAgent) translate(ctx context.Context, state *types.State) (*types.State, error) {
	sourceFiles := t.collectSourceFiles(state.Task.SourceDir)
	packageName := t.resolvePackageName(state.Task.TargetDir, state.Task.TargetLang)

	files, err := t.generateTranslation(ctx, state, sourceFiles, packageName)
	if err != nil {
		return nil, err
	}

	if err := t.syncFilesToDisk(state.Task.TargetDir, files, state); err != nil {
		return nil, err
	}
	logger.LogAgent("Translator", "Successfully wrote %d translated files to `%s`", len(files), state.Task.TargetDir)
	return state, nil
}

// repair prompts the coding model with compiler diagnostics and test failure output to fix code.
func (t *TranslatorAgent) repair(ctx context.Context, state *types.State) (*types.State, error) {
	targetFiles := t.collectCurrentTargetFiles(state)
	packageName := t.resolvePackageName(state.Task.TargetDir, state.Task.TargetLang)

	repairedFiles, err := t.generateRepair(ctx, state, targetFiles, packageName)
	if err != nil {
		return nil, err
	}

	if err := t.syncFilesToDisk(state.Task.TargetDir, repairedFiles, state); err != nil {
		return nil, err
	}
	logger.LogAgent("Translator", "Repairs applied across %d files", len(repairedFiles))
	return state, nil
}

// collectSourceFiles walks the source directory and reads raw content for each file.
func (t *TranslatorAgent) collectSourceFiles(sourceDir string) []string {
	var sourceFilesData []string
	_ = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err == nil {
			sourceFilesData = append(sourceFilesData, fmt.Sprintf("=== Source File: %s ===\n%s\n", filepath.Base(path), string(data)))
		}
		return nil
	})
	return sourceFilesData
}

// collectCurrentTargetFiles gathers the in-memory translated file contents.
func (t *TranslatorAgent) collectCurrentTargetFiles(state *types.State) []string {
	var targetFilesData []string
	for relPath, content := range state.TranslatedProject.Files {
		targetFilesData = append(targetFilesData, fmt.Sprintf("=== Current File: %s ===\n%s\n", relPath, content))
	}
	return targetFilesData
}

// resolvePackageName determines the canonical package/crate name for import statements.
func (t *TranslatorAgent) resolvePackageName(targetDir, targetLang string) string {
	packageName := sanitizeProjectName(filepath.Base(targetDir))
	if packageName == "" || packageName == "." || strings.EqualFold(packageName, targetLang) {
		parent := filepath.Base(filepath.Dir(targetDir))
		if parent != "" && parent != "." && parent != "/" {
			packageName = sanitizeProjectName(parent)
		}
	}
	if packageName == "" {
		packageName = "translated_project"
	}
	return packageName
}

// generateTranslation renders the translation prompt and queries the coding model.
func (t *TranslatorAgent) generateTranslation(ctx context.Context, state *types.State, sourceFiles []string, packageName string) (map[string]string, error) {
	logger.LogStep("Prompting Coding Model for complete `%s` translation", state.Task.TargetLang)

	prompt, err := renderPromptTemplate("translator_translate", translatorTranslatePromptTemplate, map[string]any{
		"PackageName":        packageName,
		"SourceLang":         state.Task.SourceLang,
		"TargetLang":         state.Task.TargetLang,
		"TargetLangLower":    strings.ToLower(state.Task.TargetLang),
		"SourceFiles":        strings.Join(sourceFiles, "\n"),
		"TargetDesign":       state.AnalyzerOutput.Design.RawMarkdown,
		"ImplementationPlan": state.PlanningOutput.Plan.RawPlan,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render translator prompt: %w", err)
	}

	resp, err := t.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert systems programmer translating source code into idiomatic, safe target code."),
		schema.UserMessage(prompt),
	})
	if err != nil {
		return nil, fmt.Errorf("translator model call failed: %w", err)
	}

	return parseAllFileMarkers(resp.Content), nil
}

// generateRepair renders the repair prompt and queries the coding model for targeted fixes.
func (t *TranslatorAgent) generateRepair(ctx context.Context, state *types.State, targetFiles []string, packageName string) (map[string]string, error) {
	logger.LogStep("Feeding compiler diagnostics and test failures to Coding Model for targeted repair")
	prompt, err := renderPromptTemplate("translator_repair", translatorRepairPromptTemplate, map[string]any{
		"PackageName":     packageName,
		"TargetLang":      state.Task.TargetLang,
		"TargetLangLower": strings.ToLower(state.Task.TargetLang),
		"Diagnostics":     state.ValidationReport.Diagnostics,
		"CrateCanonicalHints": crateCanonicalHints(state.Task.TargetLang),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render repair prompt: %w", err)
	}

	resp, err := t.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert systems programmer debugging compiler errors and test failures. Output only the requested code files inside fenced code blocks."),
		schema.UserMessage(prompt),
	})
	if err != nil {
		return nil, fmt.Errorf("translator repair call failed: %w", err)
	}

	return parseAllFileMarkers(resp.Content), nil
}

// syncFilesToDisk cleans, writes files to disk, and updates the in-memory state.
func (t *TranslatorAgent) syncFilesToDisk(targetDir string, files map[string]string, state *types.State) error {
	hasNewTest := false
	for relPath := range files {
		if strings.HasPrefix(relPath, "tests/") {
			hasNewTest = true
			break
		}
	}
	if hasNewTest {
		testsDir := filepath.Join(targetDir, "tests")
		_ = os.RemoveAll(testsDir)
		for k := range state.TranslatedProject.Files {
			if strings.HasPrefix(k, "tests/") {
				delete(state.TranslatedProject.Files, k)
			}
		}
	}

	for relPath, content := range files {
		clean := tools.CleanCodeContent(content)
		if relPath == "Cargo.toml" && strings.EqualFold(state.Task.TargetLang, "Rust") {
			clean = canonicalizeCargoToml(clean)
		}
		fullPath := filepath.Join(targetDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(clean), 0644); err != nil {
			return fmt.Errorf("failed to write translated file %s: %w", fullPath, err)
		}
		state.TranslatedProject.Files[relPath] = clean
		logger.LogTool("write_file", "Wrote `%s` to `%s` (%d bytes)", relPath, targetDir, len(content))
	}
	return nil
}
