package agents_test

import (
	"MAgHARCM/internal/agents"
	"strings"
	"testing"
)

func TestParseAllFileMarkersFenced(t *testing.T) {
	input := `Here is the translated code:

FILE: Cargo.toml
` + "```toml" + `
[package]
name = "my_crate"
version = "0.1.0"
` + "```" + `

Conversational prose in between.

FILE: src/lib.rs
` + "```rust" + `
pub fn hello() -> &'static str {
    "world"
}
` + "```" + `

More conversational prose at the end.`

	files := agents.ParseAllFileMarkers(input)
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(files), files)
	}

	if !strings.Contains(files["Cargo.toml"], "my_crate") {
		t.Errorf("Cargo.toml content mismatch: %s", files["Cargo.toml"])
	}
	if !strings.Contains(files["src/lib.rs"], "pub fn hello") {
		t.Errorf("src/lib.rs content mismatch: %s", files["src/lib.rs"])
	}
	if strings.Contains(files["src/lib.rs"], "Conversational prose") {
		t.Errorf("Prose leaked into source file: %s", files["src/lib.rs"])
	}
}

func TestParseAllFileMarkersUnfenced(t *testing.T) {
	input := `FILE: src/lib.rs
pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

FILE: tests/integration_tests.rs
use my_crate::*;

#[test]
fn test_add() {
    assert_eq!(add(2, 3), 5);
}
`

	files := agents.ParseAllFileMarkers(input)
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d: %v", len(files), files)
	}
	if !strings.Contains(files["src/lib.rs"], "pub fn add") {
		t.Errorf("src/lib.rs content mismatch: %s", files["src/lib.rs"])
	}
	if !strings.Contains(files["tests/integration_tests.rs"], "test_add") {
		t.Errorf("tests mismatch: %s", files["tests/integration_tests.rs"])
	}
}
