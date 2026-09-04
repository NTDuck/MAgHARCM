package agents

import (
	"context"
	"sort"
)

// Backlink: [[1.0.0 PRIM-18]] Jaccard-Coupling Architecture Recovery (MSR4SA-2017).
// Computes Jaccard similarity of commit neighborhoods between modules to detect
// hidden coupling invisible to static import edges.
// J(A, B) = |Commits(A) ∩ Commits(B)| / |Commits(A) ∪ Commits(B)|.

// JaccardCouplingEdge represents an inferred temporal coupling between two files.
type JaccardCouplingEdge struct {
	FileA      string  `json:"file_a"`
	FileB      string  `json:"file_b"`
	Similarity float64 `json:"similarity"`
	CoCommits  int     `json:"co_commits"`
}

// JaccardCouplingAnalyzer implements PRIM-18 temporal co-change analysis.
type JaccardCouplingAnalyzer struct{}

// NewJaccardCouplingAnalyzer constructs a new PRIM-18 analyzer.
func NewJaccardCouplingAnalyzer() *JaccardCouplingAnalyzer {
	return &JaccardCouplingAnalyzer{}
}

// ComputeCoupling calculates pairwise Jaccard similarity from a map of file -> set of commit hashes.
func (j *JaccardCouplingAnalyzer) ComputeCoupling(ctx context.Context, fileCommits map[string]map[string]struct{}, minThreshold float64) []JaccardCouplingEdge {
	files := make([]string, 0, len(fileCommits))
	for f := range fileCommits {
		files = append(files, f)
	}
	sort.Strings(files)

	var edges []JaccardCouplingEdge

	for i := 0; i < len(files); i++ {
		for k := i + 1; k < len(files); k++ {
			fA := files[i]
			fB := files[k]

			commitsA := fileCommits[fA]
			commitsB := fileCommits[fB]

			if len(commitsA) == 0 || len(commitsB) == 0 {
				continue
			}

			intersectionCount := 0
			for c := range commitsA {
				if _, ok := commitsB[c]; ok {
					intersectionCount++
				}
			}

			if intersectionCount == 0 {
				continue
			}

			// Union = |A| + |B| - |A ∩ B|
			unionCount := len(commitsA) + len(commitsB) - intersectionCount
			similarity := float64(intersectionCount) / float64(unionCount)

			if similarity >= minThreshold {
				edges = append(edges, JaccardCouplingEdge{
					FileA:      fA,
					FileB:      fB,
					Similarity: similarity,
					CoCommits:  intersectionCount,
				})
			}
		}
	}

	sort.Slice(edges, func(i, k int) bool {
		return edges[i].Similarity > edges[k].Similarity
	})

	return edges
}
