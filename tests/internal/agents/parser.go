package agents_test

import (
	"MAgHARCM/internal/agents"
	"testing"
)

func TestParseAllFileMarkersVariations(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name: "Standard FILE: header with fenced code",
			input: "FILE: src/lib.rs\n```rust\npub fn hello() {}\n```\n\nFILE: Cargo.toml\n```toml\n[package]\nname = \"test\"\n```",
			expected: map[string]string{
				"src/lib.rs": "pub fn hello() {}",
				"Cargo.toml": "[package]\nname = \"test\"",
			},
		},
		{
			name: "Markdown bold and backtick header: ### FILE: `src/item.rs`",
			input: "### FILE: `src/item.rs`\n```rust\npub struct Item;\n```\n\n**Path:** `tests/main.rs`\n```rust\n#[test]\nfn test_it() {}\n```",
			expected: map[string]string{
				"src/item.rs":  "pub struct Item;",
				"tests/main.rs": "#[test]\nfn test_it() {}",
			},
		},
		{
			name: "Fence annotation: ```rust src/lib.rs",
			input: "```rust src/lib.rs\npub fn add(a: i32, b: i32) -> i32 { a + b }\n```\n\n```toml Cargo.toml\n[package]\nname = \"test\"\n```",
			expected: map[string]string{
				"src/lib.rs": "pub fn add(a: i32, b: i32) -> i32 { a + b }",
				"Cargo.toml": "[package]\nname = \"test\"",
			},
		},
		{
			name: "Inline comment marker: // FILE: src/main.rs",
			input: "// FILE: src/main.rs\n```rust\nfn main() {}\n```",
			expected: map[string]string{
				"src/main.rs": "fn main() {}",
			},
		},
		{
			name: "Unfenced file block with descriptions",
			input: "FILE: src/main.rs (optional example usage)\n```rust\nfn main() {}\n```",
			expected: map[string]string{
				"src/main.rs": "fn main() {}",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agents.ParseAllFileMarkers(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d files, got %d: %v", len(tt.expected), len(result), result)
			}
			for path, expectedContent := range tt.expected {
				content, found := result[path]
				if !found {
					t.Errorf("file %s not found in result: %v", path, result)
					continue
				}
				if content != expectedContent {
					t.Errorf("file %s content mismatch.\nExpected:\n%q\nGot:\n%q", path, expectedContent, content)
				}
			}
		})
	}
}
