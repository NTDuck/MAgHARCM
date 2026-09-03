package artifacts

// PlanStep represents a single step in Part A or Part B of the implementation plan.
type PlanStep struct {
	ID              string `json:"id"`
	Description     string `json:"description"`
	SourceFile      string `json:"source_file,omitempty"`
	TargetFile      string `json:"target_file,omitempty"`
	Type            string `json:"type,omitempty"` // "source" or "test"
	Completed       bool   `json:"completed,omitempty"`
	StepName        string `json:"step_name,omitempty"`
	Details         string `json:"details,omitempty"`
	ReverseTopoRank int    `json:"reverse_topo_rank,omitempty"`
}

// ImplementationPlan organizes code translation and test verification steps into ordered phases.
type ImplementationPlan struct {
	Overview string     `json:"overview"`
	PartA    []PlanStep `json:"part_a"` // Source code translation
	PartB    []PlanStep `json:"part_b"` // Test code translation & validation
	RawPlan  string     `json:"raw_plan"`
}

// PlanningOutput captures AST fragments, symbol mappings, generated skeletons, and translation steps.
type PlanningOutput struct {
	Fragments     []string           `json:"fragments"`      // file_name:fragment_name
	NameMapping   map[string]string  `json:"name_mapping"`   // source_name -> target_name
	SkeletonFiles map[string]string  `json:"skeleton_files"` // relative_path -> skeleton_content
	Plan          ImplementationPlan `json:"plan"`
}
