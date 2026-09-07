package agents

import (
	"context"
	"slices"
	"strings"
	"unicode"

	"MAgHARCM/internal/compiletime"
)

// Backlink: [[1.0.0 PRIM-20]] Concept Assignment and Redocumentation (Rajlich-1997).
// Performs concept-locator analysis: identify what domain concepts the code implements,
// where they live, and how they relate. Output: concept -> location map.

// ConceptLocation represents a single location where a concept is instantiated.
type ConceptLocation struct {
	FilePath string `json:"file_path"`
	Symbol   string `json:"symbol,omitempty"`
	Line     int    `json:"line,omitempty"`
}

// ConceptBinding binds a domain concept to a set of physical source locations.
type ConceptBinding struct {
	Concept     string            `json:"concept"`
	Description string            `json:"description"`
	Locations   []ConceptLocation `json:"locations"`
}

// ConceptAssignmentReport holds the synthesized concept-to-location map.
type ConceptAssignmentReport struct {
	Concepts []ConceptBinding `json:"concepts"`
}

// ConceptAssigner implements PRIM-20.
type ConceptAssigner struct{}

// NewConceptAssigner constructs a new PRIM-20 assigner.
func NewConceptAssigner() *ConceptAssigner {
	return &ConceptAssigner{}
}

// AssignConcepts analyzes source files and clusters identifiers into domain concepts.
func (c *ConceptAssigner) AssignConcepts(ctx context.Context, fileContents map[string]string) *ConceptAssignmentReport {
	conceptLocations := make(map[string][]ConceptLocation)

	for file, content := range fileContents {
		tokens := tokenizeSource(content)
		for _, token := range tokens {
			if len(token) < compiletime.ConceptTokenMinLength {
				continue
			}
			lower := strings.ToLower(token)
			for _, cluster := range compiletime.DefaultConceptClusters {
				for _, kw := range cluster.Keywords {
					if strings.Contains(lower, kw) {
						conceptLocations[cluster.Label] = append(conceptLocations[cluster.Label], ConceptLocation{FilePath: file, Symbol: token})
						break
					}
				}
			}
		}
	}
	report := &ConceptAssignmentReport{
		Concepts: make([]ConceptBinding, 0, len(conceptLocations)),
	}

	for concept, locs := range conceptLocations {
		// Deduplicate locations by FilePath
		seen := make(map[string]bool)
		deduped := make([]ConceptLocation, 0, len(locs))
		for _, l := range locs {
			if !seen[l.FilePath] {
				seen[l.FilePath] = true
				deduped = append(deduped, ConceptLocation{FilePath: l.FilePath, Symbol: l.Symbol})
			}
		}
		report.Concepts = append(report.Concepts, ConceptBinding{
			Concept:     concept,
			Description: compiletime.ConceptDescriptionPlaceholder,
			Locations:   deduped,
		})
	}

	slices.SortFunc(report.Concepts, func(a, b ConceptBinding) int {
		return strings.Compare(a.Concept, b.Concept)
	})
	return report
}

func tokenizeSource(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}
