package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds centralized configuration for the MAgHARCM translation pipeline.
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

// Defaults returns a Config with default values populated from the
// environment, falling back to the documented shipped defaults.
func Defaults() *Config {
	return &Config{
		OllamaBaseURL:  envOrDefault("OLLAMA_BASE_URL", "http://localhost:11434"),
		ReasoningModel: envOrDefault("OLLAMA_REASONING_MODEL", "qwen3:30b-a3b-thinking-2507-q4_K_M"),
		CodingModel:    envOrDefault("OLLAMA_CODING_MODEL", "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL"),
		MaxIterations:  envIntOrDefault("MAGHARCM_MAX_ITERATIONS", 10),
		Timeout:        time.Duration(envIntOrDefault("MAGHARCM_TIMEOUT_SECONDS", 5000)) * time.Second,
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envIntOrDefault(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			return i
		}
	}
	return def
}
