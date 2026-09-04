package artifacts

// ArchaeologyReport is the side-channel payload produced by the PRIM-14
// software-archaeology pre-planning pass. It captures the structural
// decomposition (BoundaryMap), time-capsule reproduction commands,
// legacy-churn hotspots, naming forensics, and a co-occurrence concept
// map. Defined in artifacts (not agents) to avoid a types<->agents import
// cycle: artifacts sits at the bottom of the dependency DAG.
type ArchaeologyReport struct {
	BoundaryMap         []string               `json:"boundary_map"`
	TimeCapsuleCommands []string               `json:"time_capsule_commands"`
	ChurnHotspots       []string               `json:"churn_hotspots"`
	NamingForensics     []NamingFinding        `json:"naming_forensics"`
	ConceptMap          map[string][]string    `json:"concept_map"`
}

// NamingFinding is a single legacy-naming observation produced by the
// archaeologist's forensic-naming pass.
type NamingFinding struct {
	OldName    string `json:"old_name"`
	NewName    string `json:"new_name,omitempty"`
	Location   string `json:"location"`
	Rationale  string `json:"rationale,omitempty"`
}
