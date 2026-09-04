package agents

import (
	"context"
	"sort"
	"strings"
	"unicode"
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
			if len(token) < 4 {
				continue
			}
			lower := strings.ToLower(token)
			// Cluster common domain concepts
			if strings.Contains(lower, "validat") || strings.Contains(lower, "check") {
				conceptLocations["Validation & Verification"] = append(conceptLocations["Validation & Verification"], ConceptLocation{FilePath: file, Symbol: token})
			} else if strings.Contains(lower, "parse") || strings.Contains(lower, "lex") || strings.Contains(lower, "token") {
				conceptLocations["Parsing & Lexical Analysis"] = append(conceptLocations["Parsing & Lexical Analysis"], ConceptLocation{FilePath: file, Symbol: token})
			} else if strings.Contains(lower, "stat") || strings.Contains(lower, "math") || strings.Contains(lower, "calc") {
				conceptLocations["Mathematical & Statistical Computation"] = append(conceptLocations["Mathematical & Statistical Computation"], ConceptLocation{FilePath: file, Symbol: token})
			} else if strings.Contains(lower, "item") || strings.Contains(lower, "store") || strings.Contains(lower, "repo") {
				conceptLocations["Entity & Storage Domain"] = append(conceptLocations["Entity & Storage Domain"], ConceptLocation{FilePath: file, Symbol: token})
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
			Description: "Domain concept located through lexical and structural identifier clustering",
			Locations:   deduped,
		})
	}

	sort.Slice(report.Concepts, func(i, j int) bool {
		return report.Concepts[i].Concept < report.Concepts[j].Concept
	})

	return report
}

func tokenizeSource(s string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}
