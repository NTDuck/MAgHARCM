package agents

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/types"
)

// chunkedLoCThreshold and chunkedFragmentsThreshold together gate when the
// chunked translator takes over from the single-shot translator. Both must be
// exceeded. The thresholds are deliberately conservative so small repos still
// take the cheap single-shot path.
const (
	chunkedFragmentsThreshold = 5
	chunkedLoCThreshold       = 1000
	priorModulesBudgetBytes   = 4096
	priorModulesSnippetBytes  = 200
)

// shouldUseChunkedTranslation returns true when the planning output indicates
// a translation large enough to justify per-file chunked model calls.
//
// We gate on the count of *distinct source-file basenames* rather than the raw
// fragment count: planning.go's extractFragments emits one fragment per AST
// element, so a single source file with N functions contributes N fragments.
// Counting distinct basenames gives a stable signal that scales with file
// count (the actual unit of work for the chunked translator) instead of with
// AST-element density, which is what downstream consumers (validator
// MinRealTests, the per-file emit loop) really care about.
func shouldUseChunkedTranslation(state *types.State) bool {
	if state == nil {
		return false
	}
	distinctSources := len(GroupFragmentsBySourceFile(state.PlanningOutput.Fragments))
	if distinctSources <= chunkedFragmentsThreshold {
		return false
	}
	if countSourceLoC(state.Task.SourceDir) <= chunkedLoCThreshold {
		return false
	}
	return true
}

// countSourceLoC counts newline-terminated lines across every regular file
// under sourceDir. Used purely as a heuristic gate for chunked translation.
func countSourceLoC(sourceDir string) int {
	total := 0
	_ = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		if len(data) == 0 {
			return nil
		}
		total += strings.Count(string(data), "\n") + 1
		return nil
	})
	return total
}

// sourceFileIndex walks sourceDir and groups every regular file by basename so
// fragment IDs of the form "file.py:func_a" can be resolved to the full path
// and raw content of that file. When multiple source files share a basename,
// the first match (deterministic via filepath.Walk's lexical order) wins.
type sourceFileIndex struct {
	byBase  map[string]string // basename -> absolute path
	content map[string]string // basename -> raw content
}

func newSourceFileIndex(sourceDir string) (*sourceFileIndex, error) {
	idx := &sourceFileIndex{
		byBase:  make(map[string]string),
		content: make(map[string]string),
	}
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if _, seen := idx.byBase[base]; seen {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		idx.byBase[base] = path
		idx.content[base] = string(data)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return idx, nil
}

// fragmentBase extracts the source-file basename from a fragment ID of the
// form "base.py:func_a". Returns "" if the fragment cannot be parsed.
func fragmentBase(fragment string) string {
	idx := strings.Index(fragment, ":")
	if idx <= 0 {
		return ""
	}
	return fragment[:idx]
}
// groupFragmentsBySourceFile buckets fragment IDs by their source-file basename
// using fragmentBase. Fragments with an unparsable basename ("") are dropped.
// Iteration order of the returned map is unspecified; callers must sort keys
func GroupFragmentsBySourceFile(fragments []string) map[string][]string {
	grouped := make(map[string][]string)
	for _, frag := range fragments {
		base := fragmentBase(frag)
		if base == "" {
			continue
		}
		grouped[base] = append(grouped[base], frag)
	}
	return grouped
}

// RunChunked translates the source code in fragment-sized chunks, feeding the
// running "previously emitted modules" state into each subsequent model call.
// It populates state.TranslatedProject.Files (merged with the existing skeleton
// entries), persists everything to disk via syncFilesToDisk, and returns the
// merged TranslatedProject.
func (t *TranslatorAgent) RunChunked(ctx context.Context, state *types.State) (*types.TranslatedProject, error) {
	if state == nil {
		return nil, fmt.Errorf("nil state passed to RunChunked")
	}
	if state.TranslatedProject.Files == nil {
		state.TranslatedProject.Files = make(map[string]string)
	}

	fragments := append([]string(nil), state.PlanningOutput.Fragments...)
	if len(fragments) == 0 {
		logger.LogStep("RunChunked: no fragments in plan; preserving existing translated files")
		return &state.TranslatedProject, nil
	}

	sort.Strings(fragments)
	idx, err := newSourceFileIndex(state.Task.SourceDir)
	if err != nil {
		return nil, fmt.Errorf("failed to index source dir %q: %w", state.Task.SourceDir, err)
	}

	packageName := t.resolvePackageName(state.Task.TargetDir, state.Task.TargetLang)
	logger.LogStep("Chunked translation: %d fragments, package=%s", len(fragments), packageName)
	logger.LogStep("[CHUNKED] selected: fragments=%d loc=%d", len(fragments), countSourceLoC(state.Task.SourceDir))

	grouped := GroupFragmentsBySourceFile(fragments)
	// Sort basenames for deterministic dispatch order. GroupFragmentsBySourceFile
	// drops unparsable fragment IDs, so each surviving key corresponds to at
	// least one fragment.
	bases := make([]string, 0, len(grouped))
	for base := range grouped {
		if _, ok := idx.content[base]; !ok {
			logger.LogWarning("Skipping group for unknown source file %q (%d fragments)", base, len(grouped[base]))
			continue
		}
		bases = append(bases, base)
	}
	sort.Strings(bases)

	logger.LogStep("Chunked translation: %d fragments across %d source file(s), package=%s", len(fragments), len(bases), packageName)

	for i, base := range bases {
		fileFrags := grouped[base]
		// Render one fenced block per fragment. Each block carries the fragment
		// ID and the shared source-file content so the model can locate the
		// symbol without re-reading the file for every fragment.
		var sb strings.Builder
		fmt.Fprintf(&sb, "=== Source File: %s ===\n", base)
		for _, frag := range fileFrags {
			fmt.Fprintf(&sb, "```\n# Fragment: %s\n%s\n```\n", frag, idx.content[base])
		}
		srcBlock := sb.String()

		priorSummary := buildPriorModulesSummary(state.TranslatedProject.Files)

		logger.LogStep("Chunked translation: source %d/%d (%s, %d fragment(s))", i+1, len(bases), base, len(fileFrags))
		emitted, err := t.translateFragment(ctx, state, packageName, srcBlock, priorSummary)
		if err != nil {
			return nil, fmt.Errorf("chunked translation failed for source file %q: %w", base, err)
		}
		if len(emitted) == 0 {
			logger.LogWarning("Source file %q produced no files", base)
			continue
		}
		if err := t.syncFilesToDisk(state.Task.TargetDir, emitted, state); err != nil {
			return nil, fmt.Errorf("failed to persist chunked output for source file %q: %w", base, err)
		}
	}

	t.ensureRustCargoManifest(state)

	logger.LogAgent("Translator", "Chunked translation complete: %d total files written", len(state.TranslatedProject.Files))
	return &state.TranslatedProject, nil
}

// ensureRustCargoManifest writes a minimal Cargo.toml to state.Task.TargetDir
// when the target language is Rust and no Cargo.toml was emitted by the
// chunked loop. This unblocks the validator when the initial Translator
// skeleton wrote no build manifest (the chunked translator only emits
// per-source-file outputs, never the workspace manifest).
func (t *TranslatorAgent) ensureRustCargoManifest(state *types.State) {
	if state == nil || state.Task.TargetLang != "Rust" {
		return
	}
	if _, ok := state.TranslatedProject.Files["Cargo.toml"]; ok {
		return
	}
	packageName := t.resolvePackageName(state.Task.TargetDir, state.Task.TargetLang)
	manifest := fmt.Sprintf(
		"[package]\nname = %q\nversion = \"0.1.0\"\nedition = \"2021\"\n\n[lib]\npath = \"src/lib.rs\"\n",
		packageName,
	)
	state.TranslatedProject.Files["Cargo.toml"] = manifest
	if err := t.syncFilesToDisk(state.Task.TargetDir, map[string]string{"Cargo.toml": manifest}, state); err != nil {
		logger.LogWarning("Failed to emit fallback Cargo.toml: %v", err)
		return
	}
	logger.LogStep("Emitted fallback Cargo.toml for target %q", state.Task.TargetDir)
}


// translateFragment calls the Coding Model once for a single source fragment,
// supplying a compact summary of previously emitted modules so the model can
// re-use symbols already declared. Returns just the new fragment's target
// files extracted from the model's response.
func (t *TranslatorAgent) translateFragment(
	ctx context.Context,
	state *types.State,
	packageName, sourceBlock, priorSummary string,
) (map[string]string, error) {
	prompt, err := renderPromptTemplate("translator_chunked", translatorChunkedPromptTemplate, map[string]any{
		"PackageName":     packageName,
		"TargetLang":      state.Task.TargetLang,
		"TargetLangLower": strings.ToLower(state.Task.TargetLang),
		"SourceLang":      state.Task.SourceLang,
		"SourceFiles":     sourceBlock,
		"TargetDesign":    state.AnalyzerOutput.Design.RawMarkdown,
		"PriorModules":    priorSummary,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to render chunked translator prompt: %w", err)
	}

	resp, err := t.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage("You are an expert systems programmer translating a single source file into idiomatic, safe target code. Output only the requested code file inside a fenced code block."),
		schema.UserMessage(prompt),
	})
	if err != nil {
		return nil, fmt.Errorf("chunked translator model call failed: %w", err)
	}

	return parseAllFileMarkers(resp.Content), nil
}

// buildPriorModulesSummary renders a compact "what's already been emitted"
// block for the chunked prompt. Each already-written target file contributes
// one header line plus the first priorModulesSnippetBytes of its content.
// The whole summary is capped at priorModulesBudgetBytes; further modules are
// elided with a trailing marker once the budget is exceeded.
func buildPriorModulesSummary(files map[string]string) string {
	if len(files) == 0 {
		return "(no modules emitted yet)"
	}
	keys := make([]string, 0, len(files))
	for k := range files {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	budget := priorModulesBudgetBytes
	for _, path := range keys {
		content := files[path]
		snippet := content
		if len(snippet) > priorModulesSnippetBytes {
			snippet = snippet[:priorModulesSnippetBytes] + "\n... (truncated)\n"
		}
		entry := fmt.Sprintf("--- %s ---\n%s\n", path, snippet)
		if sb.Len()+len(entry) > budget {
			sb.WriteString("--- (further modules elided to fit summary budget) ---\n")
			break
		}
		sb.WriteString(entry)
	}
	return sb.String()
}
