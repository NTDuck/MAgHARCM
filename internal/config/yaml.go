package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// YAMLConfig represents the declarative YAML configuration format for MAgHARCM.
type YAMLConfig struct {
	Translation struct {
		Source struct {
			Dir      string `yaml:"dir"`
			Language string `yaml:"language"`
		} `yaml:"source"`
		Target struct {
			Dir       string `yaml:"dir"`
			Language  string `yaml:"language"`
			Toolchain string `yaml:"toolchain"`
		} `yaml:"target"`
		Models struct {
			Reasoning string `yaml:"reasoning"`
			Coding    string `yaml:"coding"`
			OllamaURL string `yaml:"ollama_url"`
		} `yaml:"models"`
		Execution struct {
			MaxIterations  int `yaml:"max_iterations"`
			TimeoutSeconds int `yaml:"timeout_seconds"`
		} `yaml:"execution"`
		LSP struct {
			Provider string `yaml:"provider"`
		} `yaml:"lsp"`
	} `yaml:"translation"`

	// Flat format support for convenience
	SourceDir      string `yaml:"source_dir"`
	TargetDir      string `yaml:"target_dir"`
	SourceLang     string `yaml:"source_lang"`
	TargetLang     string `yaml:"target_lang"`
	Toolchain      string `yaml:"toolchain"`
	ReasoningModel string `yaml:"reasoning_model"`
	CodingModel    string `yaml:"coding_model"`
	OllamaURL      string `yaml:"ollama_url"`
	MaxIterations  int    `yaml:"max_iterations"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
	LSPProvider    string `yaml:"lsp_provider"`
}

// ParseYAML parses raw YAML bytes into a Config struct.
func ParseYAML(data []byte) (*Config, error) {
	var y YAMLConfig
	if err := yaml.Unmarshal(data, &y); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml config: %w", err)
	}

	cfg, err := New()
	if err != nil {
		return nil, err
	}

	// Structured nested format mapping
	if y.Translation.Source.Dir != "" {
		cfg.SourceDir = y.Translation.Source.Dir
	}
	if y.Translation.Source.Language != "" {
		cfg.SourceLang = y.Translation.Source.Language
	}
	if y.Translation.Target.Dir != "" {
		cfg.TargetDir = y.Translation.Target.Dir
	}
	if y.Translation.Target.Language != "" {
		cfg.TargetLang = y.Translation.Target.Language
	}
	if y.Translation.Target.Toolchain != "" {
		cfg.Toolchain = y.Translation.Target.Toolchain
	}
	if y.Translation.Models.Reasoning != "" {
		cfg.ReasoningModel = y.Translation.Models.Reasoning
	}
	if y.Translation.Models.Coding != "" {
		cfg.CodingModel = y.Translation.Models.Coding
	}
	if y.Translation.Models.OllamaURL != "" {
		cfg.OllamaBaseURL = y.Translation.Models.OllamaURL
	}
	if y.Translation.Execution.MaxIterations > 0 {
		cfg.MaxIterations = y.Translation.Execution.MaxIterations
	}
	if y.Translation.Execution.TimeoutSeconds > 0 {
		cfg.Timeout = time.Duration(y.Translation.Execution.TimeoutSeconds) * time.Second
	}
	if y.Translation.LSP.Provider != "" {
		cfg.LSPProvider = y.Translation.LSP.Provider
	}

	// Flat format overrides if present
	if y.SourceDir != "" {
		cfg.SourceDir = y.SourceDir
	}
	if y.TargetDir != "" {
		cfg.TargetDir = y.TargetDir
	}
	if y.SourceLang != "" {
		cfg.SourceLang = y.SourceLang
	}
	if y.TargetLang != "" {
		cfg.TargetLang = y.TargetLang
	}
	if y.Toolchain != "" {
		cfg.Toolchain = y.Toolchain
	}
	if y.ReasoningModel != "" {
		cfg.ReasoningModel = y.ReasoningModel
	}
	if y.CodingModel != "" {
		cfg.CodingModel = y.CodingModel
	}
	if y.OllamaURL != "" {
		cfg.OllamaBaseURL = y.OllamaURL
	}
	if y.MaxIterations > 0 {
		cfg.MaxIterations = y.MaxIterations
	}
	if y.TimeoutSeconds > 0 {
		cfg.Timeout = time.Duration(y.TimeoutSeconds) * time.Second
	}
	if y.LSPProvider != "" {
		cfg.LSPProvider = y.LSPProvider
	}

	return cfg, nil
}

// LoadYAML reads a YAML file from disk and parses it into a Config.
func LoadYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read yaml config file %s: %w", path, err)
	}
	return ParseYAML(data)
}

// MustLoadYAML reads and parses a YAML config or panics using the Must pattern.
func MustLoadYAML(path string) *Config {
	cfg, err := LoadYAML(path)
	if err != nil {
		panic(err)
	}
	return cfg
}
