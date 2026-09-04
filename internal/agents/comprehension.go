package agents

import (
	"context"
	"path/filepath"
	"sort"
	"strings"
)

// Backlink: [[1.0.0 PRIM-22]] Four Phases of Comprehension (Foltz-2023).
// Applies the DR. JONES cognitive model:
// Decomposition, Recognition, Organization, Navigation, Explanation, Search.
// Formulates empirical traversal: linear traversal preferred, recency bias,
// breadth before depth, and hotspot concentration.

// DRJonesPhases represents the cognitive steps during code comprehension.
type DRJonesPhases struct {
	Decomposition []string            `json:"decomposition"` // structural units (files/modules)
	Recognition   []string            `json:"recognition"`   // recognized libraries, idioms, patterns
	Organization  map[string][]string `json:"organization"`  // hierarchy groupings
	Navigation    []string            `json:"navigation"`    // ordered reading path
	Explanation   string              `json:"explanation"`   // mental model summary
	SearchAnchors []string            `json:"search_anchors"`// key lookup symbols
}

// ComprehensionPipeline implements PRIM-22.
type ComprehensionPipeline struct{}

// NewComprehensionPipeline constructs a new PRIM-22 pipeline.
func NewComprehensionPipeline() *ComprehensionPipeline {
	return &ComprehensionPipeline{}
}

// Comprehend executes the DR. JONES cognitive traversal over the source file set.
func (c *ComprehensionPipeline) Comprehend(ctx context.Context, files []string, fileContents map[string]string) *DRJonesPhases {
	phases := &DRJonesPhases{
		Decomposition: make([]string, len(files)),
		Recognition:   make([]string, 0),
		Organization:  make(map[string][]string),
		Navigation:    make([]string, len(files)),
		SearchAnchors: make([]string, 0),
	}

	copy(phases.Decomposition, files)
	sort.Strings(phases.Decomposition)

	// Recognition: scan imports and common signatures
	seenPatterns := make(map[string]bool)
	for _, content := range fileContents {
		if strings.Contains(content, "import") || strings.Contains(content, "#include") {
			if strings.Contains(content, "math") && !seenPatterns["Math/Numerics Library"] {
				seenPatterns["Math/Numerics Library"] = true
				phases.Recognition = append(phases.Recognition, "Math/Numerics Library")
			}
			if strings.Contains(content, "test") && !seenPatterns["Testing Framework"] {
				seenPatterns["Testing Framework"] = true
				phases.Recognition = append(phases.Recognition, "Testing Framework")
			}
		}
	}

	// Organization: group files by directory
	for _, f := range phases.Decomposition {
		dir := filepath.Dir(f)
		phases.Organization[dir] = append(phases.Organization[dir], filepath.Base(f))
	}

	// Navigation: breadth before depth (linear traversal preferred)
	copy(phases.Navigation, phases.Decomposition)
	sort.Slice(phases.Navigation, func(i, j int) bool {
		depthI := strings.Count(phases.Navigation[i], "/")
		depthJ := strings.Count(phases.Navigation[j], "/")
		if depthI != depthJ {
			return depthI < depthJ
		}
		return phases.Navigation[i] < phases.Navigation[j]
	})

	// Search anchors: top-level files
	for _, f := range phases.Navigation {
		if strings.Count(f, "/") <= 1 {
			phases.SearchAnchors = append(phases.SearchAnchors, filepath.Base(f))
		}
	}

	phases.Explanation = "Decomposed into structural units with breadth-first linear traversal ordering."

	return phases
}
