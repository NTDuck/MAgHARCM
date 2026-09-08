package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/checkpoint"
	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
)

// Type aliases (cycle-free Locality of Behaviour).
// Producer files declare the algorithms; canonical artifact types
// live in internal/compiletime/state.go. Aliases let method receivers
// reference the type name without package qualification.
type TranslatedProject = compiletime.TranslatedProject

// TranslatorAgent generates target source and test implementations, and iteratively resolves compiler errors in repair mode.
type TranslatorAgent struct {
	Model model.BaseChatModel
	// RunID identifies the current translation run; when set, a checkpoint of
	// *types.State is persisted under .artifacts/<RunID>/checkpoints/ at the
	// end of every Run (both success and error paths). Empty RunID disables
	// checkpointing — used by callers that don't want disk persistence.
	RunID string
	// IterativeNavigator is the PRIM-31 dynamic fragment index.
	IterativeNavigator *IterativeNavigator
}

// compiletime.TranslatedProject contains the files written or edited in the target
// repository. Producer: TranslatorAgent (this file).
// NewTranslatorAgent creates a TranslatorAgent instance. runID enables
// per-run checkpoint persistence; pass "" to disable checkpointing.
func NewTranslatorAgent(m model.BaseChatModel, runID string) *TranslatorAgent {
	return &TranslatorAgent{Model: m, RunID: runID}
}

// Run translates the code and tests, or executes repairs if validation report has failures.
func (t *TranslatorAgent) Run(ctx context.Context, state *compiletime.State) (*compiletime.State, error) {
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
		logger.LogAgent("Translator", "Repair Mode iteration %d of %d: Diagnosing validation errors and fixing code",
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

	logger.LogAgent("Translator", "Initial Translation Mode: Implementing Part A Source and Part B Tests")
	return t.translate(ctx, state)
}

// checkpoint persists a snapshot of state if RunID is set. Errors are logged
// but never propagated — checkpointing is best-effort and must not abort the
// pipeline when the disk is full or the runID is empty.
func (t *TranslatorAgent) checkpoint(state *State) {
	if t.RunID == "" || state == nil {
		return
	}
	if path, err := checkpoint.Save(t.RunID, state); err != nil {
		logger.LogWarning("Translator checkpoint save failed: %v", err)
	} else {
		logger.LogStep("Translator checkpoint saved: %s", path)
	}
}

// translate generates the initial translation from source modules, design, and implementation plan.
func (t *TranslatorAgent) translate(ctx context.Context, state *compiletime.State) (*compiletime.State, error) {
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
func (t *TranslatorAgent) repair(ctx context.Context, state *compiletime.State) (*compiletime.State, error) {
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
func (t *TranslatorAgent) collectCurrentTargetFiles(state *State) []string {
	var targetFilesData []string
	for relPath, content := range state.TranslatedProject.Files {
		targetFilesData = append(targetFilesData, fmt.Sprintf("=== Current File: %s ===\n%s\n", relPath, content))
	}
	return targetFilesData
}

// resolvePackageName determines the canonical package/crate name for import statements and project metadata.
func resolvePackageName(targetDir, targetLang string) string {
	packageName := sanitizeProjectName(filepath.Base(targetDir))
	if packageName == "" || packageName == "." || strings.EqualFold(packageName, targetLang) {
		parent := filepath.Base(filepath.Dir(targetDir))
		if parent != "" && parent != "." && parent != "/" {
			packageName = sanitizeProjectName(parent)
		}
	}
	if packageName == "" {
		packageName = compiletime.TranslatedPackagePlaceholder
	}
	return packageName
}

func (t *TranslatorAgent) resolvePackageName(targetDir, targetLang string) string {
	return resolvePackageName(targetDir, targetLang)
}

// generateTranslation renders the translation prompt and queries the coding model.
func (t *TranslatorAgent) generateTranslation(ctx context.Context, state *compiletime.State, sourceFiles []string, packageName string) (map[string]string, error) {
	logger.LogStep("Prompting Coding Model for complete `%s` translation", state.Task.TargetLang)

	prompt, err := renderPromptTemplate("translator_translate", translatorTranslatePromptTemplate, map[string]any{
		"PackageName":        packageName,
		"SourceLang":         state.Task.SourceLang,
		"TargetLang":         state.Task.TargetLang,
		"TargetLangLower":    strings.ToLower(state.Task.TargetLang),
		"SourceFiles":        strings.Join(sourceFiles, "\n"),
		"TargetDesign":       state.AnalyzerOutput.Design.RawMarkdown,
		"compiletime.ImplementationPlan": state.PlanningOutput.Plan.RawPlan,
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
func (t *TranslatorAgent) generateRepair(ctx context.Context, state *compiletime.State, targetFiles []string, packageName string) (map[string]string, error) {
	logger.LogStep("Feeding compiler diagnostics, current files, and test failures to Coding Model for targeted repair")
	prompt, err := renderPromptTemplate("translator_repair", translatorRepairPromptTemplate, map[string]any{
		"PackageName":         packageName,
		"TargetLang":          state.Task.TargetLang,
		"TargetLangLower":     strings.ToLower(state.Task.TargetLang),
		"Diagnostics":         state.ValidationReport.Diagnostics,
		"CurrentFiles":        strings.Join(targetFiles, "\n"),
		"CrateCanonicalHints": CrateCanonicalHints(state.Task.TargetLang),
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
func (t *TranslatorAgent) syncFilesToDisk(targetDir string, files map[string]string, state *compiletime.State) error {
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
			clean = CanonicalizeCargoToml(clean)
		}
		if strings.EqualFold(state.Task.TargetLang, "Rust") && strings.HasPrefix(relPath, "tests/") {
			clean = NormalizeRustTestImports(clean, state)
		}
		fullPath := filepath.Join(targetDir, relPath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte(clean), 0644); err != nil {
			return fmt.Errorf("failed to write translated file %s: %w", fullPath, err)
		}
		state.TranslatedProject.Files[relPath] = clean
		logger.LogTool("write_file", "Wrote `%s` to `%s`, %d bytes", relPath, targetDir, len(content))
	}
	return nil
}

// NormalizeRustTestImports injects the package-root glob import into every
// top-level Rust `mod tests { ... }` block that lacks an import of the
// package (or of crate::). Rust child modules do not inherit a file-top
// `use`, so an SLM that emits `mod tests { ... }` with no import inside the
// block fails with E0425 no matter how often the repair prompt repeats the
// guideline. Deterministic write-time normalization removes that reliance
// entirely. An unused glob is only a warning, so the import is injected
// unconditionally — no project-specific symbol list. The import lands
// immediately after the mod's opening line, at the block's own indentation,
// so it is at mod scope and precedes every use site (including any
// #[test]-attributed fn).
func NormalizeRustTestImports(content string, state *compiletime.State) string {
	if !strings.EqualFold(state.Task.TargetLang, "Rust") {
		return content
	}
	pkg := resolvePackageName(state.Task.TargetDir, state.Task.TargetLang)
	importLine := "use " + pkg + "::*;"

	// Pass 1: locate each top-level `mod tests {` block, its indentation,
	// and whether it already imports the package or crate::.
	lines := strings.Split(content, "\n")
	type modBlock struct {
		openLine int    // index of the `mod tests {` line
		indent   string // the mod block's own leading whitespace
		needs    bool   // lacks an in-block package import
	}
	var blocks []modBlock
	depth := 0
	current := -1 // index into blocks of the mod currently being scanned
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if current >= 0 {
			if strings.Contains(trimmed, "use "+pkg) || strings.Contains(trimmed, "use crate::") {
				blocks[current].needs = false
			}
			depth += strings.Count(line, "{") - strings.Count(line, "}")
			if depth <= 0 {
				current = -1
			}
			continue
		}
		if isRustTestsModOpen(trimmed) {
			indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			blocks = append(blocks, modBlock{openLine: i, indent: indent, needs: true})
			current = len(blocks) - 1
			depth = 1
		}
	}
	if len(blocks) == 0 {
		return content
	}

	// Pass 2: rebuild with the import right after each qualifying mod-open.
	inject := map[int]string{} // line index -> import line to insert after
	for _, b := range blocks {
		if b.needs {
			inject[b.openLine] = b.indent + "\t" + importLine
		}
	}
	var out []string
	for i, line := range lines {
		out = append(out, line)
		if ins, ok := inject[i]; ok {
			out = append(out, ins)
		}
	}
	return strings.Join(out, "\n")
}

// isRustTestsModOpen reports whether the line opens a test module block.
func isRustTestsModOpen(trimmed string) bool {
	return strings.HasPrefix(trimmed, "mod ") && strings.HasSuffix(trimmed, "{") &&
		!strings.Contains(trimmed, ";")
}
