package artifacts

// SourceProjectResearch represents the research document produced in phase 3.2.1.
type SourceProjectResearch struct {
	Overview           string   `json:"overview"`
	DirectoryStructure string   `json:"directory_structure"`
	StructsInterfaces  string   `json:"structs_and_interfaces"`
	DataModels         string   `json:"data_models"`
	ErrorHandling      string   `json:"error_handling"`
	Dependencies       []string `json:"dependencies"`
	MigrationStrategy  string   `json:"migration_strategy,omitempty"`
	StrategyRationale  string   `json:"strategy_rationale,omitempty"`
	RawDocument        string   `json:"raw_document"`
}

// ThirdPartyLibraryAnalysis represents the library analysis document produced in phase 3.2.2.
type ThirdPartyLibraryAnalysis struct {
	Libraries   []LibraryMapping `json:"libraries"`
	RawDocument string           `json:"raw_document"`
}

// LibraryMapping details how a source library maps to a target library.
type LibraryMapping struct {
	SourceLibrary   string `json:"source_library"`
	TargetLibrary   string `json:"target_library"`
	Overview        string `json:"overview"`
	Usage           string `json:"usage"`
	Recommendations string `json:"recommendations"`
}

// TargetProjectDesign represents the design document produced in phase 3.2.3.
type TargetProjectDesign struct {
	Overview                string   `json:"overview"`
	TranslationRequirements string   `json:"translation_requirements"`
	SourceFilesToTranslate  []string `json:"source_files_to_translate"`
	ModuleStructure         string   `json:"module_structure"`
	ErrorHandling           string   `json:"error_handling"`
	ThirdPartyLibraries     []string `json:"third_party_libraries"`
	RawDocument             string   `json:"raw_document"`
}

// AnalyzerOutput aggregates research, library mapping, and architectural design documents.
type AnalyzerOutput struct {
	Research DocumentWrapper[SourceProjectResearch]     `json:"research"`
	Library  DocumentWrapper[ThirdPartyLibraryAnalysis] `json:"library"`
	Design   DocumentWrapper[TargetProjectDesign]       `json:"design"`
}
