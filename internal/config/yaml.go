package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// YAMLConfig is the on-disk schema; only the nested translation.* form
// is supported. (The previous flat-yaml shim was unused outside tests and
// was a duplicate-source-of-truth hazard — kept just the live format.)
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
}

// LoadYAML reads a YAML file and returns a Config populated with defaults
// overridden by the values present in the file.
func LoadYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return ParseYAML(data)
}

// ParseYAML decodes YAML bytes into a Config.
func ParseYAML(data []byte) (*Config, error) {
	var y YAMLConfig
	if err := yaml.Unmarshal(data, &y); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	cfg := Defaults()
	t := y.Translation
	if t.Source.Dir != "" {
		cfg.SourceDir = t.Source.Dir
	}
	if t.Source.Language != "" {
		cfg.SourceLang = t.Source.Language
	}
	if t.Target.Dir != "" {
		cfg.TargetDir = t.Target.Dir
	}
	if t.Target.Language != "" {
		cfg.TargetLang = t.Target.Language
	}
	if t.Target.Toolchain != "" {
		cfg.Toolchain = t.Target.Toolchain
	}
	if t.Models.Reasoning != "" {
		cfg.ReasoningModel = t.Models.Reasoning
	}
	if t.Models.Coding != "" {
		cfg.CodingModel = t.Models.Coding
	}
	if t.Models.OllamaURL != "" {
		cfg.OllamaBaseURL = t.Models.OllamaURL
	}
	if t.Execution.MaxIterations > 0 {
		cfg.MaxIterations = t.Execution.MaxIterations
	}
	if t.Execution.TimeoutSeconds > 0 {
		cfg.Timeout = time.Duration(t.Execution.TimeoutSeconds) * time.Second
	}
	if t.LSP.Provider != "" {
		cfg.LSPProvider = t.LSP.Provider
	}
	return cfg, nil
}
