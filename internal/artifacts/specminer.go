package artifacts

// SpecMinerInvariants is a structural alias for the SpecMiner-style dynamic
// invariants (allocation sizes, pointer nullability, aliasing, lifetime
// ranges, branch coverage). Defined here — rather than imported from
// internal/agents — to avoid an import cycle: artifacts is the bottom of
// the dependency DAG and types/agents depend on it, never the reverse.
// Producers in internal/agents populate a value of this shape; consumers
// read it from state.SpecMinerInvariants.
type SpecMinerInvariants struct {
	AllocSizes         []int            `json:"alloc_sizes"`
	PointerNullability map[string]bool  `json:"pointer_nullability"`
	AliasingPairs      [][2]string      `json:"aliasing_pairs"`
	LifetimeRanges     map[string][2]int `json:"lifetime_ranges"`
	BranchCoverage     map[string]float64 `json:"branch_coverage"`
}
