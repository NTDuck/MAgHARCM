package types

// TranslationTask defines the input specification for the translation pipeline.
type TranslationTask struct {
	SourceDir   string `json:"source_dir"`
	TargetDir   string `json:"target_dir"`
	SourceLang  string `json:"source_lang"`
	TargetLang  string `json:"target_lang"`
	Toolchain   string `json:"toolchain,omitempty"`
	LSPProvider string `json:"lsp_provider,omitempty"`
}
