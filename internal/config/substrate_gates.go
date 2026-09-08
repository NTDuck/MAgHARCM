package config

import "gopkg.in/yaml.v3"

// SubstrateGates is the typed schema for the §11.x substrate opt-in keys
// referenced across .obsidian/MAgHARCM/research/Methodology.md §11 and
// Architecture.md §9. Each boolean opt-in corresponds to one Wave-NN SLM-era
// substrate anchor:
//
//   - MemoryDistilled                      → Wave-17 (P-122 ReasoningBank ICLR 2026)
//   - MaTTSEnabled                         → Wave-17 (MaTTS compute-memory loop)
//   - OracleCrossLingual                   → Wave-17 (P-123 CodeChemist ICML 2026)
//   - TranslationDynamicSpecs              → Wave-17 (P-124 Syzygy ICLR 2025 Wksp)
//   - ComprehensionGraphSelfEvolving       → Wave-21/23 (P-153 CoReX ICSE 2026)
//   - ComprehensionHallucinationEvaluation → Wave-22 (P-152 Hallu-Eval FSE 2026)
//   - TranslationFeedbackDriven            → Wave-22/23 (P-151 SmartC2Rust + P-154 TransAgent)
//
// STATUS: Wave-NN-OPEN. No code path in cmd/ or internal/ currently reads
// these flags — they exist as a typed schema only. Adding a `substrates:`
// block to configs/agents.yml and wiring `LoadSubstrateGates()` into a
// consumer is a future-sprint task that MUST ship alongside the substrate
// implementation (no silent wiring of unimplemented features). Per ADR-C-014
// the struct lives alongside its future consumer; per ADR-C-005 the on-disk
// defaults MUST be `false` (safe-blind).
type SubstrateGates struct {
	MemoryDistilled                      bool `yaml:"memory.distilled"`
	MaTTSEnabled                         bool `yaml:"mattts.enabled"`
	OracleCrossLingual                   bool `yaml:"oracle.cross_lingual"`
	TranslationDynamicSpecs              bool `yaml:"translation.dynamic_specs"`
	ComprehensionGraphSelfEvolving       bool `yaml:"comprehension.graph_self_evolving"`
	ComprehensionHallucinationEvaluation bool `yaml:"comprehension.hallucination_evaluation"`
	TranslationFeedbackDriven            bool `yaml:"translation.feedback_driven"`
}

// allSubstrateKeys is the canonical ordered list of yaml tags this schema
// recognises. The TestSubstrateGatesKeys test guards against accidental
// drift between the struct yaml tags and the docs-cited key set.
var allSubstrateKeys = []string{
	"memory.distilled",
	"mattts.enabled",
	"oracle.cross_lingual",
	"translation.dynamic_specs",
	"comprehension.graph_self_evolving",
	"comprehension.hallucination_evaluation",
	"translation.feedback_driven",
}

// MarshalSubstrateGates renders the schema as a YAML document with every
// key explicitly set to false. Used by future config-generators that emit
// configs/agents.yml scaffolds; today it is the canonical example used in
// unit tests.
func MarshalSubstrateGates() ([]byte, error) {
	return yaml.Marshal(SubstrateGates{})
}
