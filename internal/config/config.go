package config

import (
	"os"
	"strconv"
	"time"

	"MAgHARCM/internal/pattern"
)

// Config holds centralized configuration for the MAgHARCM translation pipeline.
type Config struct {
	OllamaBaseURL  string        `json:"ollama_base_url"`
	ReasoningModel string        `json:"reasoning_model"`
	CodingModel    string        `json:"coding_model"`
	MaxIterations  int           `json:"max_iterations"`
	Timeout        time.Duration `json:"timeout"`
	SourceDir      string        `json:"source_dir"`
	TargetDir      string        `json:"target_dir"`
	SourceLang     string        `json:"source_lang"`
	TargetLang     string        `json:"target_lang"`
	Toolchain      string        `json:"toolchain"`
	LSPProvider    string        `json:"lsp_provider"`
}

// Option configures a Config instance.
type Option func(*Config)

// WithOllamaBaseURL sets the base URL for Ollama.
func WithOllamaBaseURL(url string) Option {
	return func(c *Config) {
		if url != "" {
			c.OllamaBaseURL = url
		}
	}
}

// WithReasoningModel sets the reasoning model name.
func WithReasoningModel(model string) Option {
	return func(c *Config) {
		if model != "" {
			c.ReasoningModel = model
		}
	}
}

// WithCodingModel sets the coding model name.
func WithCodingModel(model string) Option {
	return func(c *Config) {
		if model != "" {
			c.CodingModel = model
		}
	}
}

// WithMaxIterations sets the maximum repair iterations.
func WithMaxIterations(n int) Option {
	return func(c *Config) {
		if n > 0 {
			c.MaxIterations = n
		}
	}
}

// WithTimeout sets the overall execution timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Config) {
		if d > 0 {
			c.Timeout = d
		}
	}
}

// WithToolchain sets the build/test toolchain.
func WithToolchain(toolchain string) Option {
	return func(c *Config) {
		if toolchain != "" {
			c.Toolchain = toolchain
		}
	}
}

// WithLSPProvider sets the LSP provider ("native" or "abcoder").
func WithLSPProvider(provider string) Option {
	return func(c *Config) {
		if provider != "" {
			c.LSPProvider = provider
		}
	}
}

// WithSourceTarget sets the source and target project directories and languages.
func WithSourceTarget(sourceDir, targetDir, sourceLang, targetLang string) Option {
	return func(c *Config) {
		if sourceDir != "" {
			c.SourceDir = sourceDir
		}
		if targetDir != "" {
			c.TargetDir = targetDir
		}
		if sourceLang != "" {
			c.SourceLang = sourceLang
		}
		if targetLang != "" {
			c.TargetLang = targetLang
		}
	}
}

// New creates a Config with defaults and applied options.
func New(opts ...Option) (*Config, error) {
	cfg := &Config{
		OllamaBaseURL:  getEnvOrDefault("OLLAMA_BASE_URL", "http://localhost:11434"),
		ReasoningModel: getEnvOrDefault("OLLAMA_REASONING_MODEL", "gpt-oss:20b"),
		CodingModel:    getEnvOrDefault("OLLAMA_CODING_MODEL", "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL"),
		MaxIterations:  getEnvIntOrDefault("MAGHARCM_MAX_ITERATIONS", 10),
		Timeout:        time.Duration(getEnvIntOrDefault("MAGHARCM_TIMEOUT_SECONDS", 5000)) * time.Second,
		SourceDir:      getEnvOrDefault("MAGHARCM_SOURCE_DIR", ""),
		TargetDir:      getEnvOrDefault("MAGHARCM_TARGET_DIR", ""),
		SourceLang:     getEnvOrDefault("MAGHARCM_SOURCE_LANG", ""),
		TargetLang:     getEnvOrDefault("MAGHARCM_TARGET_LANG", ""),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg, nil
}

// MustLoad returns a loaded Config using the Must pattern.
func MustLoad(opts ...Option) *Config {
	return pattern.Must(New(opts...))
}

func getEnvOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func getEnvIntOrDefault(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil && i > 0 {
			return i
		}
	}
	return def
}
