// Backlink: [[Primitives]] §NEW-PRIM-09 (Tri-Representation Hybrid Code Graph, RepoGraph [[P62]]).
//
// HybridCodeGraph fuses three program-representation layers: AST (Tree-Sitter,
// consts.LSPProviderNative), CPG (Yamaguchi's AST ∪ CFG ∪ PDG), and SDG
// (Horwitz/Reps/Binkley's inter-procedural PDG). Skeleton + types; future work
// fills CFG/PDG/call-graph backends.

package agents

import (
	"context"
	"fmt"

	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/logger"
)

const (
	LayerAST = "ast"
	LayerCPG = "cpg"
	LayerSDG = "sdg"

	EdgeCPGAST  = "cpg_ast"
	EdgeSDGCall = "sdg_call"
)

// ASTNode is a single Tree-Sitter derived node fed into the graph.
type ASTNode struct {
	ID, Kind, Name, FilePath, Parent string
	StartLn, EndLn                   int
}

// CPGNode is a node in Yamaguchi's Code Property Graph.
type CPGNode struct {
	ID, ASTID, Function, FilePath string
	StartLn, EndLn                int
	IsAST, IsCFG, IsPDG           bool
}
type CPGEdge struct{ From, To, Kind string }

// CPG is Yamaguchi's unified intra-procedural graph (AST ∪ CFG ∪ PDG).
type CPG struct {
	Nodes []CPGNode
	Edges []CPGEdge
}

// SDGNode is a procedure-level node in the System Dependence Graph.
type SDGNode struct {
	ID, CPGID, Procedure, FilePath string
	IsEntry                        bool
}

// SDGEdge is a system-dependence edge between two SDGNode procedure nodes
// (Horwitz/Reps/Binkley). Kind tags the dependence class — call, parameter,
// return — and From / To are the source and target SDGNode IDs.
type SDGEdge struct {
	Kind string
	From string
	To   string
}

// SDG is Horwitz/Reps/Binkley's inter-procedural extension of the PDG.
type SDG struct {
	Nodes []SDGNode
	Edges []SDGEdge
}

// Resolution is the consolidated answer from ResolveSymbol.
type Resolution struct {
	Found     bool
	Layers    []string
	Locations []string
	Snippet   string
}

// LSPResolver is the slice of `tools.LSPProvider` the graph needs; kept
// local to avoid an import cycle and let tests pass an in-memory fake.
type LSPResolver interface{ Name() string }

// HybridCodeGraph fuses AST + CPG and SDG and exposes symbol resolution.
type HybridCodeGraph struct {
	ast  []ASTNode
	cpg  CPG
	sdg  SDG
	lsp  LSPResolver
	lang string
}

// NewHybridCodeGraph constructs a graph anchored on lang. lsp may be nil.
func NewHybridCodeGraph(lang string, lsp LSPResolver) *HybridCodeGraph {
	if lang == "" {
		lang = compiletime.ErrorUnknown
	}
	g := &HybridCodeGraph{lang: lang, lsp: lsp}
	logger.LogStep("hybrid code graph created lang=%s lsp=%v", lang, lsp != nil)
	return g
}

// BuildAST indexes the supplied AST nodes (future: delegate to Tree-Sitter).
func (g *HybridCodeGraph) BuildAST(nodes []ASTNode) {
	g.ast = append(g.ast[:0], nodes...)
	logger.LogStep("hybrid code graph indexed %d AST nodes", len(g.ast))
}

// BuildCPG computes Yamaguchi's CPG (AST ∪ CFG ∪ PDG) from AST nodes. AST
// sub-graph populated; CFG/PDG stay empty for future control-/data-flow.
func BuildCPG(astNodes []ASTNode) (CPG, error) {
	if len(astNodes) == 0 {
		return CPG{}, fmt.Errorf("BuildCPG: no AST nodes supplied")
	}
	cpg := CPG{Nodes: make([]CPGNode, 0, len(astNodes)), Edges: make([]CPGEdge, 0, len(astNodes))}
	for _, n := range astNodes {
		cpg.Nodes = append(cpg.Nodes, CPGNode{
			ID: "cpg:" + n.ID, ASTID: n.ID, Function: n.Name,
			FilePath: n.FilePath, StartLn: n.StartLn, EndLn: n.EndLn, IsAST: true,
		})
		if n.Parent != "" {
			cpg.Edges = append(cpg.Edges, CPGEdge{From: "cpg:" + n.Parent, To: "cpg:" + n.ID, Kind: EdgeCPGAST})
		}
	}
	return cpg, nil
}

// BuildSDG computes the System Dependence Graph on top of cpg. Procedure-
// level nodes derive from CPG function roots; call/param edges stay empty.
func BuildSDG(cpg CPG) (SDG, error) {
	if len(cpg.Nodes) == 0 {
		return SDG{}, fmt.Errorf("BuildSDG: empty CPG")
	}
	sdg := SDG{Nodes: make([]SDGNode, 0, len(cpg.Nodes))}
	seen := make(map[string]struct{}, len(cpg.Nodes))
	for _, n := range cpg.Nodes {
		if n.Function == "" {
			continue
		}
		key := n.FilePath + "::" + n.Function
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		sdg.Nodes = append(sdg.Nodes, SDGNode{
			ID:    "sdg:" + n.Function + "@" + n.FilePath,
			CPGID: n.ID, Procedure: n.Function, FilePath: n.FilePath, IsEntry: true,
		})
	}
	return sdg, nil
}

// ResolveSymbol consults AST, CPG, and SDG and returns a consolidated match.
func (g *HybridCodeGraph) ResolveSymbol(ctx context.Context, name string) (Resolution, error) {
	if err := ctx.Err(); err != nil {
		return Resolution{}, err
	}
	if name == "" {
		return Resolution{}, fmt.Errorf("ResolveSymbol: empty name")
	}
	var layers, locations []string
	for _, n := range g.ast {
		if n.Name != name {
			continue
		}
		layers = append(layers, LayerAST)
		locations = append(locations, fmt.Sprintf("%s:%d", n.FilePath, n.StartLn))
	}
	for _, n := range g.cpg.Nodes {
		if n.Function != name {
			continue
		}
		layers = append(layers, LayerCPG)
		locations = append(locations, fmt.Sprintf("%s:%d", n.FilePath, n.StartLn))
		break
	}
	for _, n := range g.sdg.Nodes {
		if n.Procedure != name {
			continue
		}
		layers = append(layers, LayerSDG)
		locations = append(locations, fmt.Sprintf("%s:0", n.FilePath))
		break
	}
	if len(layers) == 0 {
		logger.LogStep("ResolveSymbol miss name=%s lang=%s", name, g.lang)
		return Resolution{Found: false}, nil
	}
	logger.LogStep("ResolveSymbol hit name=%s layers=%d sites=%d", name, len(layers), len(locations))
	return Resolution{
		Found: true, Layers: layers, Locations: locations,
		Snippet: fmt.Sprintf("// skeleton snippet for %q", name),
	}, nil
}

// SetCPG / SetSDG let callers inject precomputed layers.
func (g *HybridCodeGraph) SetCPG(cpg CPG) { g.cpg = cpg }
func (g *HybridCodeGraph) SetSDG(sdg SDG) { g.sdg = sdg }

// ASTLen / CPGLen / SDGLen report current graph sizes.
func (g *HybridCodeGraph) ASTLen() int { return len(g.ast) }
func (g *HybridCodeGraph) CPGLen() int { return len(g.cpg.Nodes) }
func (g *HybridCodeGraph) SDGLen() int { return len(g.sdg.Nodes) }
