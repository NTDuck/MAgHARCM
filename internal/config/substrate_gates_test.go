package config

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestSubstrateGatesYAMLRoundTrip verifies the typed schema round-trips
// through yaml.v3 with zero-value defaults. This guards the Wave-NN-OPEN
// scaffold without requiring configs/agents.yml to actually carry the
// `substrates:` block (which is a future-sprint wiring decision, per the
// file header).
func TestSubstrateGatesYAMLRoundTrip(t *testing.T) {
	out, err := MarshalSubstrateGates()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got SubstrateGates
	if err := yaml.Unmarshal(out, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range allSubstrateKeys {
		if !strings.Contains(string(out), key+": false") {
			t.Errorf("expected %q to appear with value false in YAML output:\n%s", key, string(out))
		}
	}
}

// TestSubstrateGatesZeroIsSafeBlind asserts every field defaults to false.
// Per ADR-C-005 the Wave-NN-OPEN scaffold MUST be safe-blind; flipping a key
// to true is the deliberate activation surface, not this default.
func TestSubstrateGatesZeroIsSafeBlind(t *testing.T) {
	var g SubstrateGates
	if g.MemoryDistilled || g.MaTTSEnabled || g.OracleCrossLingual ||
		g.TranslationDynamicSpecs || g.ComprehensionGraphSelfEvolving ||
		g.ComprehensionHallucinationEvaluation || g.TranslationFeedbackDriven {
		t.Fatalf("SubstrateGates zero value must have all false; got %+v", g)
	}
}

// TestSubstrateGatesKeysMatchDocs asserts the 7 yaml tags exactly match the
// docs-cited §11.x key set (Methodology §11, Architecture §9.3). Drift here
// would silently break a future consumer that reads by string key.
func TestSubstrateGatesKeysMatchDocs(t *testing.T) {
	g := SubstrateGates{}
	// Re-encode through yaml.Node to harvest the field tags via a marshal/unmarshal
	// cycle, then assert the doc-cited key set matches the field-set order.
	data, err := yaml.Marshal(g)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, key := range allSubstrateKeys {
		if !strings.Contains(string(data), key) {
			t.Errorf("docs-cited key %q not present in SubstrateGates YAML output", key)
		}
	}
	if got := strings.Count(string(data), ": false"); got != len(allSubstrateKeys) {
		t.Errorf("expected %d false entries, got %d", len(allSubstrateKeys), got)
	}
}
