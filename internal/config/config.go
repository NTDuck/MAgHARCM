package config

import (
	"time"
)

// Default constant values for MAgHARCM pipeline configuration.
// See Obsidian vault: [[Methodology]] §2 "The 4+1 Agents" and [[Primitives]] §NEW-PRIM-23.
const (
	DefaultOllamaBaseURL  = "http://localhost:11434"
	DefaultReasoningModel = "qwen3:30b-a3b-thinking-2507-q4_K_M"
	DefaultCodingModel    = "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL"
	DefaultMaxIterations  = 10
	DefaultTimeoutSeconds = 5000
)

// Config holds centralized configuration for the MAgHARCM translation pipeline.
// Backlink: [[Design Space]] §Configuration and [[Methodology]] §4.
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

// Defaults returns a Config populated directly with shipped baseline defaults.
// Default fallbacks via environment variables are omitted per repository configuration standards.
// Backlink: [[Methodology]] §4 (Balanced MAgHARCM CAND-08).
func Defaults() *Config {
	return &Config{
		OllamaBaseURL:  DefaultOllamaBaseURL,
		ReasoningModel: DefaultReasoningModel,
		CodingModel:    DefaultCodingModel,
		MaxIterations:  DefaultMaxIterations,
		Timeout:        time.Duration(DefaultTimeoutSeconds) * time.Second,
	}
}
