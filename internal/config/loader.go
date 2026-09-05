package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// MustLoadConfig is the generic, panic-on-failure companion to LoadYAML /
// ParseYAML. It reads the YAML file at path, decodes it into a fresh value
// of type T, and returns it. Any error — file missing, unreadable, or
// malformed YAML — produces a panic. There are no default fallbacks: a
// caller who supplies a typed config MUST also supply the corresponding
// file or the program is misconfigured.
//
// Generic instantiation relies on the caller's type parameter; callers
// typically declare a typed wrapper alongside the corresponding YAML file:
//
//	var cfg AgentTimeouts = config.MustLoadConfig[AgentTimeouts]("configs/agents.yml")
//
// Per ADR-C-014, the typed config struct itself lives alongside the
// consumer that needs it (typically its producer agent). This loader is
// purely the generic plumbing.
func MustLoadConfig[T any](path string) T {
	var cfg T
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("MustLoadConfig %s failed: read: %v", path, err))
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(fmt.Sprintf("MustLoadConfig %s failed: parse: %v", path, err))
	}
	return cfg
}
