package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// YAMLConfig is the on-disk schema; only the nested translation.* form
// is supported.
// Backlink: [[Methodology]] §1 "The 6-Stage Pipeline" and [[Primitives]] §NEW-PRIM-24.
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

// LoadYAML reads a YAML file and returns a fully-populated Config.
// The file MUST set every field; missing fields produce an error.
func LoadYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return ParseYAML(data)
}

// MustLoadYAML reads a YAML file and returns *Config or panics on any error.
func MustLoadYAML(path string) *Config {
	cfg, err := LoadYAML(path)
	if err != nil {
		panic(fmt.Sprintf("MustLoadYAML %s failed: %v", path, err))
	}
	return cfg
}

// MustParseYAML decodes YAML bytes and returns *Config or panics on any error.
func MustParseYAML(data []byte) *Config {
	cfg, err := ParseYAML(data)
	if err != nil {
		panic(fmt.Sprintf("MustParseYAML failed: %v", err))
	}
	return cfg
}

// ParseYAML decodes YAML bytes into a Config. Every field must be present in
// the file; missing fields are reported as a single structured error.
func ParseYAML(data []byte) (*Config, error) {
	var y YAMLConfig
	if err := yaml.Unmarshal(data, &y); err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	cfg := &Config{
		SourceDir:      y.Translation.Source.Dir,
		SourceLang:     y.Translation.Source.Language,
		TargetDir:      y.Translation.Target.Dir,
		TargetLang:     y.Translation.Target.Language,
		Toolchain:      y.Translation.Target.Toolchain,
		ReasoningModel: y.Translation.Models.Reasoning,
		CodingModel:    y.Translation.Models.Coding,
		OllamaBaseURL:  y.Translation.Models.OllamaURL,
		MaxIterations:  y.Translation.Execution.MaxIterations,
		Timeout:        time.Duration(y.Translation.Execution.TimeoutSeconds) * time.Second,
		LSPProvider:    y.Translation.LSP.Provider,
	}

	if err := Require(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Require checks that every config field is set and returns a single
// structured error listing the missing fields. Fields are required: there are
// no default fallbacks.
func Require(cfg *Config) error {
	var missing []string
	if cfg.OllamaBaseURL == "" {
		missing = append(missing, "models.ollama_url")
	}
	if cfg.ReasoningModel == "" {
		missing = append(missing, "models.reasoning")
	}
	if cfg.CodingModel == "" {
		missing = append(missing, "models.coding")
	}
	if cfg.MaxIterations == 0 {
		missing = append(missing, "execution.max_iterations")
	}
	if cfg.Timeout == 0 {
		missing = append(missing, "execution.timeout_seconds")
	}
	if cfg.SourceDir == "" {
		missing = append(missing, "source.dir")
	}
	if cfg.TargetDir == "" {
		missing = append(missing, "target.dir")
	}
	if cfg.SourceLang == "" {
		missing = append(missing, "source.language")
	}
	if cfg.TargetLang == "" {
		missing = append(missing, "target.language")
	}
	if cfg.Toolchain == "" {
		missing = append(missing, "target.toolchain")
	}
	if cfg.LSPProvider == "" {
		missing = append(missing, "lsp.provider")
	}
	if len(missing) > 0 {
		return fmt.Errorf("config: missing required fields: %s", strings.Join(missing, ", "))
	}
	return nil
}

