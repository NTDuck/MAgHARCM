package agents

import (
	"strings"
	"testing"
)

func TestCanonicalizeCargoToml(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "renames underscore serde_xml_rs",
			input: `[dependencies]
serde_xml_rs = "0.4"
serde_json = "1.0"
`,
			expected: `[dependencies]
serde-xml-rs = "0.4"
serde_json = "1.0"
`,
		},
		{
			name: "renames hyphenated quick_xml and yaml_rust in dev-dependencies",
			input: `[dev-dependencies]
quick_xml = "0.31"
yaml_rust = "0.4"
`,
			expected: `[dev-dependencies]
quick-xml = "0.31"
yaml-rust = "0.4"
`,
		},
		{
			name: "preserves unrelated keys and lines",
			input: `[package]
name = "demo"
version = "0.1.0"

[dependencies]
serde_xml_rs = "0.4"
`,
			expected: `[package]
name = "demo"
version = "0.1.0"

[dependencies]
serde-xml-rs = "0.4"
`,
		},
		{
			name:     "no-op when key absent",
			input:    `[dependencies]\nserde_json = "1.0"\n`,
			expected: `[dependencies]\nserde_json = "1.0"\n`,
		},
		{
			name:     "empty string is empty",
			input:    "",
			expected: "",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := canonicalizeCargoToml(tc.input)
			if got != tc.expected {
				t.Fatalf("canonicalizeCargoToml(%q)\n  got:  %q\n  want: %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestCrateCanonicalHints(t *testing.T) {
	hints := crateCanonicalHints("Rust")
	if hints == "" {
		t.Fatal("expected non-empty hints for Rust")
	}
	if !strings.Contains(hints, "serde_xml_rs -> serde-xml-rs") {
		t.Fatalf("expected serde_xml_rs rename in hints, got: %s", hints)
	}
	if !strings.Contains(hints, "quick_xml -> quick-xml") {
		t.Fatalf("expected quick_xml rename in hints, got: %s", hints)
	}
	// Same-name entries (anyhow, tokio, ...) should NOT appear.
	if strings.Contains(hints, "anyhow -> anyhow") {
		t.Fatalf("identical rename should not appear in hints: %s", hints)
	}
	if crateCanonicalHints("Go") != "" {
		t.Fatal("expected empty hints for non-Rust target")
	}
}

func TestRenameCrateKey(t *testing.T) {
	got := renameCrateKey("serde_xml_rs = \"0.4\"\n", "serde_xml_rs", "serde-xml-rs")
	if got != "serde-xml-rs = \"0.4\"\n" {
		t.Fatalf("got %q", got)
	}
	// non-matching line passes through
	got = renameCrateKey("serde_json = \"1.0\"\n", "serde_xml_rs", "serde-xml-rs")
	if got != "serde_json = \"1.0\"\n" {
		t.Fatalf("non-matching should pass through, got %q", got)
	}
	// partial prefix should not match
	got = renameCrateKey("serde_xml_rs_extra = \"0.4\"\n", "serde_xml_rs", "serde-xml-rs")
	if got != "serde_xml_rs_extra = \"0.4\"\n" {
		t.Fatalf("partial prefix should not match, got %q", got)
	}
}
