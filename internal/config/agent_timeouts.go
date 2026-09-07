package config

import "time"

// AgentTimeouts is the typed shape of the `agent_timeouts:` block in
// configs/agents.yml. It is the canonical example of a per-agent typed
// config loaded via MustLoadConfig[AgentTimeouts]. Per ADR-C-014, this
// struct lives alongside the consumers that need it; new typed configs
// MUST follow the same pattern (file in package consumers OR a dedicated
// configs/<name>_types.go next to the YAML file).
type AgentTimeouts struct {
	// ReasoningTimeout caps the wall-clock budget of a single
	// reasoning-model invocation (Analyzer, Planner, Validator, etc.).
	ReasoningTimeout time.Duration `yaml:"reasoning_timeout"`
	// CodingTimeout caps the wall-clock budget of a single coding-model
	// invocation (Translator, repair loop, test synthesis).
	CodingTimeout time.Duration `yaml:"coding_timeout"`
	// MaxIterations caps the coverage-remediation loop. The Validator
	// overrides this at runtime when a tighter bound is appropriate;
	// the YAML value is the upper ceiling.
	MaxIterations int `yaml:"max_iterations"`
}
