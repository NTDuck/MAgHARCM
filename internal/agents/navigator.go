package agents

// Backlink: [[Methodology]] §2 The 4+1 Agents and [[Primitives]] §NEW-PRIM-26.

import (
	"context"
	"fmt"
	"path/filepath"

	"MAgHARCM/internal/logger"
	"MAgHARCM/internal/tools"
)

// SymbolResolution is the planner-callable result of a Navigator.LookupSymbol call.
// It bundles the definition, all references, and the hover info into a single
// structured response that the translator can consume without making 3 separate
// LSP calls. Per P13 HyperAgent §6.2 Table 6, this is the "SymbolAwareRetriever"
// output — the substrate of "navigator-guided LLM translation".
type SymbolResolution struct {
	Symbol     string                  // the symbol name queried
	Definition *tools.DefinitionOutput // where the symbol is defined
	References *tools.ReferencesOutput // all usages across the project
	Hover      *tools.HoverOutput      // documentation/signature
	Error      error                   // first non-nil error encountered (partial results above)
}

// Navigator wraps the LSP tools in a planner-callable interface.
// The Navigator is intentionally stateful (it holds an LSPProvider reference)
// so the planner can reuse the same Navigator across many Lookups without
// re-establishing the LSP connection.
type Navigator struct {
	Provider tools.LSPProvider
}

// NewNavigator returns a Navigator bound to the given LSPProvider.
// Pass `nil` to disable LSP-backed lookups; LookupSymbol will return
// Error=ErrNoLSPProvider in that case. Callers should fall back to
// LLM-only symbol resolution in that scenario.
func NewNavigator(provider tools.LSPProvider) *Navigator {
	return &Navigator{Provider: provider}
}

// ErrNoLSPProvider is returned by LookupSymbol when the Navigator has no
// LSP provider configured. Callers should fall back to LLM-only symbol
// resolution in this case (the translator can still ask the LLM directly).
var ErrNoLSPProvider = fmt.Errorf("navigator: no LSP provider configured")

// LookupSymbol performs a single combined symbol-resolution call.
// Returns a SymbolResolution with whatever fields could be resolved;
// Error is non-nil only if a fatal error occurred (e.g., no provider).
// Partial results are returned when individual sub-calls fail: a failed
// sub-call leaves its output as nil but does not erase the others.
func (n *Navigator) LookupSymbol(ctx context.Context, symbol, filePath string) SymbolResolution {
	if n.Provider == nil {
		return SymbolResolution{Symbol: symbol, Error: ErrNoLSPProvider}
	}
	res := SymbolResolution{Symbol: symbol}

	def, err := n.Provider.GetDefinition(ctx, &tools.DefinitionInput{
		Symbol:   symbol,
		FilePath: filePath,
	})
	if err != nil {
		logger.LogTool("navigator", "definition lookup failed for %q: %v", symbol, err)
		res.Error = err
	} else {
		res.Definition = def
	}

	refs, err := n.Provider.GetReferences(ctx, &tools.ReferencesInput{
		Symbol:     symbol,
		ProjectDir: ProjectDirOrDot(filePath),
	})
	if err != nil {
		logger.LogTool("navigator", "references lookup failed for %q: %v", symbol, err)
		if res.Error == nil {
			res.Error = err
		}
	} else {
		res.References = refs
	}

	hover, err := n.Provider.GetHover(ctx, &tools.HoverInput{
		Symbol:   symbol,
		FilePath: filePath,
	})
	if err != nil {
		logger.LogTool("navigator", "hover lookup failed for %q: %v", symbol, err)
		if res.Error == nil {
			res.Error = err
		}
	} else {
		res.Hover = hover
	}

	logger.LogTool("navigator", "SymbolResolution(%q): defined=%v refs=%d hover=%v",
		symbol, res.Definition != nil, RefCount(res.References), res.Hover != nil)
	return res
}

// LookupSymbols performs multiple symbol resolutions in one call.
// Returns a map keyed by symbol name. Convenience method for the planner
// that wants to resolve a batch of symbols at once.
func (n *Navigator) LookupSymbols(ctx context.Context, symbols []string, filePath string) map[string]SymbolResolution {
	out := make(map[string]SymbolResolution, len(symbols))
	for _, sym := range symbols {
		out[sym] = n.LookupSymbol(ctx, sym, filePath)
	}
	return out
}

// refCount returns the number of references in the output, or 0 if nil.
func RefCount(refs *tools.ReferencesOutput) int {
	if refs == nil {
		return 0
	}
	return len(refs.References)
}

// projectDirOrDot returns the directory containing filePath, or "." if filePath is empty.
// Used as a default ProjectDir for tools that take a project root rather than a file.
func ProjectDirOrDot(filePath string) string {
	if filePath == "" {
		return "."
	}
	dir := filepath.Dir(filePath)
	if dir == "" {
		return "."
	}
	return dir
}
