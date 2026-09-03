package config

import (
	"fmt"
	"os"
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
	applyTranslationOverrides(cfg, &y)
	return cfg, nil
}

// applyTranslationOverrides applies parsed YAML translation attributes to Config.
// Backlink: [[Design Space]] §Configuration.
func applyTranslationOverrides(cfg *Config, y *YAMLConfig) {
	t := y.Translation
	applySourceConfig(cfg, t.Source.Dir, t.Source.Language)
	applyTargetConfig(cfg, t.Target.Dir, t.Target.Language, t.Target.Toolchain)
	applyModelsConfig(cfg, t.Models.Reasoning, t.Models.Coding, t.Models.OllamaURL)
	applyExecutionConfig(cfg, t.Execution.MaxIterations, t.Execution.TimeoutSeconds)
	applyLSPConfig(cfg, t.LSP.Provider)
}

func applySourceConfig(cfg *Config, dir, lang string) {
	if dir != "" {
		cfg.SourceDir = dir
	}
	if lang != "" {
		cfg.SourceLang = lang
	}
}

func applyTargetConfig(cfg *Config, dir, lang, toolchain string) {
	if dir != "" {
		cfg.TargetDir = dir
	}
	if lang != "" {
		cfg.TargetLang = lang
	}
	if toolchain != "" {
		cfg.Toolchain = toolchain
	}
}

func applyModelsConfig(cfg *Config, reasoning, coding, ollamaURL string) {
	if reasoning != "" {
		cfg.ReasoningModel = reasoning
	}
	if coding != "" {
		cfg.CodingModel = coding
	}
	if ollamaURL != "" {
		cfg.OllamaBaseURL = ollamaURL
	}
}

func applyExecutionConfig(cfg *Config, maxIterations, timeoutSeconds int) {
	if maxIterations > 0 {
		cfg.MaxIterations = maxIterations
	}
	if timeoutSeconds > 0 {
		cfg.Timeout = time.Duration(timeoutSeconds) * time.Second
	}
}

func applyLSPConfig(cfg *Config, provider string) {
	if provider != "" {
		cfg.LSPProvider = provider
	}
}
