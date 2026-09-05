package agents

// Backlink: [[Methodology]] §2 The 4+1 Agents and [[Primitives]] §NEW-PRIM-26.

import (
	"context"
	"path/filepath"

	"github.com/cloudwego/eino/components/model"

	"MAgHARCM/internal/compiletime"
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
	// CodeGraph is the PRIM-9 HybridCodeGraph (AST ∪ CPG ∪ SDG) consulted
	// when LSP returns no definition. Optional; nil disables the fallback.
	CodeGraph *HybridCodeGraph
}

// NewNavigator returns a Navigator bound to the given LSPProvider.
// Pass `nil` to disable LSP-backed lookups; LookupSymbol will return
// Error=ErrNoLSPProvider in that case. Callers should fall back to
// LLM-only symbol resolution in this scenario.
func NewNavigator(provider tools.LSPProvider) *Navigator {
	return &Navigator{Provider: provider}
}

// MustNavigator returns a Navigator bound to the given LSPProvider and panics
// if the provider is nil. Use this at startup where a missing LSP provider is
// a fatal configuration error.
func MustNavigator(provider tools.LSPProvider) *Navigator {
	if provider == nil {
		panic("compiletime.MustNavigator: LSPProvider must not be nil")
	}
	return NewNavigator(provider)
}

// ErrNoLSPProvider is the back-compat alias for compiletime.ErrNavigatorNoProvider.
// New code should reference compiletime.ErrNavigatorNoProvider directly.
var ErrNoLSPProvider = compiletime.ErrNavigatorNoProvider

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

	// PRIM-9 HybridCodeGraph fallback: when LSP has no definition, consult
	// the AST ∪ CPG ∪ SDG graph and synthesize one. Useful in multi-language
	// projects where the LSP only covers the host language.
	if res.Definition == nil && n.CodeGraph != nil {
		if gr, err := n.CodeGraph.ResolveSymbol(ctx, symbol); err == nil && gr.Found {
			res.Definition = &tools.DefinitionOutput{
				Symbol: symbol,
				Definitions: []tools.DefinitionLocation{
					{FilePath: firstLocation(gr.Locations), Snippet: gr.Snippet},
				},
			}
		}
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

	logger.LogTool("navigator", "SymbolResolution name=%q defined=%v refs=%d hover=%v",
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
func ProjectDirOrDot(filePath string) string {
	if filePath == "" {
		return compiletime.ProjectDirPlaceholder
	}
	dir := filepath.Dir(filePath)
	if dir == "" {
	}
	return dir
}

// Backlink: [[Methodology]] §1 Stage 2.5 (Symbol Navigation) and [[Primitives]] §NEW-PRIM-26.
//
// NavigatorAgent is the 5th node in the MAgHARCM graph: it sits between the
// Analyzer and the Planner and resolves any unresolved symbols in the
// analyzer's output via the LSP-backed Navigator. For iteration 1 it is a
// best-effort no-op when no LSP provider is configured (e.g., a translate
// run without LSP) — the existing AnalyzerAgent still owns the embedded
// Navigator and will continue to issue lookups directly. Once the
// analyzer's symbol-aware path is retired, NavigatorAgent becomes the
// sole owner of symbol resolution.
type NavigatorAgent struct {
	Model     model.BaseChatModel
	Navigator *Navigator // LSP-backed symbol resolver; nil disables lookups
}

// NewNavigatorAgent builds a NavigatorAgent. Pass provider=nil to keep
// the agent in disabled no-op mode (Run returns state unchanged).
func NewNavigatorAgent(m model.BaseChatModel, provider tools.LSPProvider) *NavigatorAgent {
	return &NavigatorAgent{
		Model:     m,
		Navigator: NewNavigator(provider),
	}
}

// maxNavigatorLookups caps the number of symbols resolved per Run. The
// NameMapping is bounded by source-file count in practice, but a hostile
// or pathological mapping should not stall the pipeline on N^2 LSP calls.
const maxNavigatorLookups = 50

// Run resolves every symbol in state.PlanningOutput.NameMapping by
// delegating to the underlying Navigator. Returned state is the same
// pointer as the input; the agent does not mutate state in iteration 1
// (resolved information is logged but not yet folded back into the
// planning output). When no Navigator or no Provider is configured the
// state is forwarded unchanged so the rest of the pipeline is unaffected.
func (n *NavigatorAgent) Run(ctx context.Context, state *compiletime.State) (*compiletime.State, error) {
	if state == nil {
		return nil, nil
	}
	if n == nil || n.Navigator == nil || n.Navigator.Provider == nil {
		logger.LogStep("navigator: disabled, no provider; forwarding state unchanged")
		return state, nil
	}

	symbols := make([]string, 0, len(state.PlanningOutput.NameMapping))
	for sym := range state.PlanningOutput.NameMapping {
		if len(symbols) >= maxNavigatorLookups {
			break
		}
		symbols = append(symbols, sym)
	}
	if len(symbols) == 0 {
		return state, nil
	}

	logger.LogStep("navigator: resolving %d symbols, capped at %d",
		len(state.PlanningOutput.NameMapping), maxNavigatorLookups)

	// Default file path: derived from the task source dir. The LSP provider
	// may ignore this when the symbol is project-global; passing "" is
	// safer than fabricating a non-existent path.
	filePath := state.Task.SourceDir
	for _, sym := range symbols {
		res := n.Navigator.LookupSymbol(ctx, sym, filePath)
		if res.Error != nil {
			logger.LogTool("navigator", "lookup %q failed: %v", sym, res.Error)
			continue
		}
		logger.LogTool("navigator", "resolved %q: defined=%v refs=%d hover=%v",
			sym, res.Definition != nil, RefCount(res.References), res.Hover != nil)
	}
	return state, nil
}

// firstLocation returns the first entry of locations, or "" when empty.
// Used by the PRIM-9 HybridCodeGraph fallback to extract a FilePath for the
// synthesized DefinitionLocation.
func firstLocation(locations []string) string {
	if len(locations) == 0 {
		return ""
	}
	return locations[0]
}
