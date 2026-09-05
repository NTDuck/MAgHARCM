package agents

// Backlink: [[1.0.0 PRIM-22]] Four Phases of Comprehension (Foltz-2023).
// Applies the DR. JONES cognitive model:
// Decomposition, Recognition, Organization, Navigation, Explanation, Search.
// Formulates empirical traversal: linear traversal preferred, recency bias,
// breadth before depth, and hotspot concentration.

import (
	"context"
	"path/filepath"
	"slices"
	"strings"

	"MAgHARCM/internal/compiletime"
)

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
	slices.Sort(phases.Decomposition)

	// Recognition: scan imports and common signatures
	seenPatterns := make(map[string]bool)
	for _, content := range fileContents {
		if strings.Contains(content, "import") || strings.Contains(content, "#include") {
			if strings.Contains(content, "math") && !seenPatterns[compiletime.ComprehensionRecognitionMath] {
				seenPatterns[compiletime.ComprehensionRecognitionMath] = true
				phases.Recognition = append(phases.Recognition, compiletime.ComprehensionRecognitionMath)
			}
			if strings.Contains(content, "test") && !seenPatterns[compiletime.ComprehensionRecognitionTesting] {
				seenPatterns[compiletime.ComprehensionRecognitionTesting] = true
				phases.Recognition = append(phases.Recognition, compiletime.ComprehensionRecognitionTesting)
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
	slices.SortFunc(phases.Navigation, func(a, b string) int {
		depthA := strings.Count(a, "/")
		depthB := strings.Count(b, "/")
		if depthA != depthB {
			return depthA - depthB
		}
		return strings.Compare(a, b)
	})

	// Search anchors: top-level files
	for _, f := range phases.Navigation {
		if strings.Count(f, "/") <= 1 {
			phases.SearchAnchors = append(phases.SearchAnchors, filepath.Base(f))
		}
	}

	phases.Explanation = compiletime.ComprehensionExplanationDefault

	return phases
}
