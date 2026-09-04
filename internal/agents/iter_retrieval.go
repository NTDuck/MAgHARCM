package agents

// Backlink: [[Primitives]] §NEW-PRIM-30 (Iterative Retrieval Refinement, RepoCoder-style
// [P61]). Sits between the Navigator and the Translator: every translated fragment
// emitted by upstream stages is reindexed so subsequent Lookups for symbols those
// fragments already touched can be served from the index instead of another LSP
// round-trip. The "4KB budget" is RepoCoder's local context window — it caps how
// much stale context follows a reindexed symbol so the prompt stays bounded.

import (
	"context"
	"strings"
	"sync"

	"MAgHARCM/internal/consts"
	"MAgHARCM/internal/logger"
)

// TranslatedFragment is one unit of upstream work the iterative navigator has
// already produced. Path is the source file the fragment came from; Symbols are
// the names the fragment references or defines; Body is the translated text.
type TranslatedFragment struct {
	Path    string
	Symbols []string
	Body    string
}

// IterativeResolution is the planner-callable result of IterativeNavigator.Lookup.
// Source is one of consts.SourceReindexed (served from the local fragment
// index) or consts.SourceFresh (served by a fresh Navigator round-trip).
type IterativeResolution struct {
	Symbol    string
	Body      string
	SizeBytes int
	Source    string
}

// contextBudgetBytes is the maximum body size returned by Lookup. RepoCoder
// found 4 KiB is enough to anchor a symbol without dragging in the rest of the
// file. Symbols whose Body exceeds this are truncated before being returned.
const contextBudgetBytes = 4 * 1024

// IterativeNavigator wraps a Navigator with a per-symbol index of the
// fragments already emitted by the translation pipeline. Lookups first try
// the index; on a miss they fall back to the wrapped Navigator and tag the
// result as SourceFresh. Safe for concurrent use — the index is guarded by
// a mutex and the wrapped Navigator's own concurrency model applies to the
// fallback path.
type IterativeNavigator struct {
	*Navigator
	mu      sync.RWMutex
	index   map[string][]TranslatedFragment
}

// NewIterativeNavigator returns an IterativeNavigator wrapping the given
// Navigator. A nil Navigator is permitted; Lookup will then return an empty
// IterativeResolution tagged SourceReindexed with the symbol name only, mirroring
// the Navigator's own no-provider semantics.
func NewIterativeNavigator(n *Navigator) *IterativeNavigator {
	return &IterativeNavigator{
		Navigator: n,
		index:     make(map[string][]TranslatedFragment),
	}
}

// Reindex records the fragments produced by an upstream translation pass so
// later Lookups can answer from the index. Fragments with an empty Path are
// ignored — there is no way to attribute them later. Reindex is the only
// writer to the index; concurrent calls are safe.
func (in *IterativeNavigator) Reindex(_ context.Context, previouslyTranslated []TranslatedFragment) error {
	in.mu.Lock()
	defer in.mu.Unlock()
	for _, frag := range previouslyTranslated {
		if frag.Path == "" {
			continue
		}
		for _, sym := range frag.Symbols {
			sym = strings.TrimSpace(sym)
			if sym == "" {
				continue
			}
			in.index[sym] = append(in.index[sym], frag)
		}
	}
	logger.LogStep("iter_retrieval: indexed %d fragments across %d symbols", len(previouslyTranslated), len(in.index))
	return nil
}

// Lookup resolves a symbol. If at least one indexed fragment mentions the
// symbol and its Body fits the 4 KiB budget, the trimmed Body is returned
// with Source=consts.SourceReindexed. Otherwise the wrapped Navigator is
// consulted; its SymbolResolution's Definition content (or an empty string
// if no provider is configured) is returned with Source=consts.SourceFresh.
// Stale fragments whose Body exceeds the budget are dropped from the
// returned Body but stay in the index for callers that want the full text.
func (in *IterativeNavigator) Lookup(ctx context.Context, symbol string) (IterativeResolution, error) {
	if symbol == "" {
		return IterativeResolution{}, nil
	}

	if body, ok := in.indexedBody(symbol); ok {
		logger.LogStep("iter_retrieval: hit %q from index (%d bytes)", symbol, len(body))
		return IterativeResolution{
			Symbol:    symbol,
			Body:      body,
			SizeBytes: len(body),
			Source:    consts.SourceReindexed,
		}, nil
	}

	res := IterativeResolution{Symbol: symbol, Source: consts.SourceFresh}
	if in.Navigator != nil {
		sr := in.LookupSymbol(ctx, symbol, "")
		if sr.Definition != nil && len(sr.Definition.Definitions) > 0 {
			res.Body = sr.Definition.Definitions[0].Snippet
		}
	}
	if len(res.Body) > contextBudgetBytes {
		res.Body = res.Body[:contextBudgetBytes]
	}
	res.SizeBytes = len(res.Body)
	logger.LogStep("iter_retrieval: miss %q; fresh=%d bytes", symbol, res.SizeBytes)
	return res, nil
}

// indexedBody returns the trimmed body of the first indexed fragment for
// symbol that fits the 4 KiB budget, or ("", false) if no indexed fragment
// qualifies. It does not mutate the index — fragments that exceed the
// budget remain available for callers that ask for them by other paths.
func (in *IterativeNavigator) indexedBody(symbol string) (string, bool) {
	in.mu.RLock()
	frags, ok := in.index[symbol]
	in.mu.RUnlock()
	if !ok {
		return "", false
	}
	for _, frag := range frags {
		if len(frag.Body) <= contextBudgetBytes {
			return frag.Body, true
		}
	}
	return "", false
}
