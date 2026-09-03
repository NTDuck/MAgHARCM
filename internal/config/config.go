package config

import (
	"time"
)

// Config holds centralized configuration for the MAgHARCM translation pipeline.
// Backlink: [[Design Space]] §Configuration and [[Methodology]] §4.
//
// Every field is required and must be supplied via the YAML configuration
// file. There are no default fallbacks; missing fields are reported by
// Require or ParseYAML as a structured error.
type Config struct {
	OllamaBaseURL  string
	ReasoningModel string
	CodingModel    string
	MaxIterations  int
	Timeout        time.Duration
	SourceDir      string
	TargetDir      string
	SourceLang     string
	TargetLang     string
	Toolchain      string
	LSPProvider    string
}

// Defaults returns a zero-valued *Config for use in tests only. Production
// code MUST obtain *Config via LoadYAML / ParseYAML and then call Require.
func Defaults() *Config {
	return &Config{}
}
