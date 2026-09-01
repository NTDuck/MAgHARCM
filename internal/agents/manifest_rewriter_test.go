package agents

import (
	"strings"
	"testing"
)

func TestManifestRewriterBasicRewrite(t *testing.T) {
	rw := NewManifestRewriter()
	deps := []SourceDep{
		{Name: "commons-validator", Version: "1.7.0", Scope: "compile"},
		{Name: "guava", Version: "32.0.0", Scope: "compile"},
		{Name: "commons-codec", Version: "1.15", Scope: "compile"},
		{Name: "unknown-thing", Version: "1.0", Scope: "compile"},
	}
	toml, unmapped := rw.Rewrite("commons-validator-demo", "Java", deps)
	if len(unmapped) != 1 || unmapped[0] != "unknown-thing" {
		t.Fatalf("expected unmapped=[unknown-thing], got %v", unmapped)
	}
	if !strings.Contains(toml, "[package]") {
		t.Fatalf("missing [package] section: %q", toml)
	}
	if !strings.Contains(toml, `name = "commons-validator-demo"`) {
		t.Fatalf("missing pkg name: %q", toml)
	}
	if !strings.Contains(toml, `version = "0.1.0"`) {
		t.Fatalf("missing version: %q", toml)
	}
	if !strings.Contains(toml, `edition = "2021"`) {
		t.Fatalf("missing edition: %q", toml)
	}
	if !strings.Contains(toml, "[dependencies]") {
		t.Fatalf("missing [dependencies]: %q", toml)
	}
	// commons-validator → validator with derive feature
	if !strings.Contains(toml, "validator") {
		t.Fatalf("missing validator mapping: %q", toml)
	}
	if !strings.Contains(toml, `"derive"`) {
		t.Fatalf("missing derive feature: %q", toml)
	}
	// guava → itertools
	if !strings.Contains(toml, "itertools") {
		t.Fatalf("missing itertools mapping: %q", toml)
	}
	// commons-codec → base64
	if !strings.Contains(toml, "base64") {
		t.Fatalf("missing base64 mapping: %q", toml)
	}
	// Unmapped dep should NOT appear in toml
	if strings.Contains(toml, "unknown-thing") {
		t.Fatalf("unmapped dep leaked into toml: %q", toml)
	}
}

func TestManifestRewriterTestScope(t *testing.T) {
	rw := NewManifestRewriter()
	deps := []SourceDep{
		{Name: "commons-validator", Version: "1.7.0", Scope: "compile"},
		{Name: "junit", Version: "4.13.2", Scope: "test"},
	}
	toml, unmapped := rw.Rewrite("my-pkg", "Java", deps)
	if len(unmapped) != 0 {
		t.Fatalf("expected 0 unmapped, got %v", unmapped)
	}
	if !strings.Contains(toml, "[dependencies]") {
		t.Fatalf("missing [dependencies]: %q", toml)
	}
	if !strings.Contains(toml, "[dev-dependencies]") {
		t.Fatalf("missing [dev-dependencies]: %q", toml)
	}
	// validator (from commons-validator) should be in [dependencies]
	depsIdx := strings.Index(toml, "[dependencies]")
	devIdx := strings.Index(toml, "[dev-dependencies]")
	if depsIdx < 0 || devIdx < 0 {
		t.Fatalf("section indices not found: deps=%d dev=%d", depsIdx, devIdx)
	}
	if depsIdx > devIdx {
		t.Fatalf("[dependencies] should come before [dev-dependencies]")
	}
	// validator must appear before [dev-dependencies]
	if !strings.Contains(toml[:devIdx], "validator") {
		t.Fatalf("validator missing from deps section: %q", toml[:devIdx])
	}
	// rstest (from junit) must appear in [dev-dependencies]
	if !strings.Contains(toml[devIdx:], "rstest") {
		t.Fatalf("rstest missing from dev-dependencies: %q", toml[devIdx:])
	}
	// And must NOT appear before [dev-dependencies]
	if strings.Contains(toml[:devIdx], "rstest") {
		t.Fatalf("rstest leaked into [dependencies]: %q", toml[:devIdx])
	}
}

func TestManifestRewriterEmptyDeps(t *testing.T) {
	rw := NewManifestRewriter()
	toml, unmapped := rw.Rewrite("minimal", "Java", []SourceDep{})
	if len(unmapped) != 0 {
		t.Fatalf("expected 0 unmapped, got %v", unmapped)
	}
	if !strings.Contains(toml, "[package]") {
		t.Fatalf("missing [package]: %q", toml)
	}
	if !strings.Contains(toml, `name = "minimal"`) {
		t.Fatalf("missing pkg name: %q", toml)
	}
	// Empty deps should produce no [dependencies] or [dev-dependencies] section.
	if strings.Contains(toml, "[dependencies]") {
		t.Fatalf("unexpected [dependencies] for empty deps: %q", toml)
	}
	if strings.Contains(toml, "[dev-dependencies]") {
		t.Fatalf("unexpected [dev-dependencies] for empty deps: %q", toml)
	}
}

func TestManifestRewriterAllUnmapped(t *testing.T) {
	rw := NewManifestRewriter()
	deps := []SourceDep{
		{Name: "obscure-lib-a", Version: "1.0", Scope: "compile"},
		{Name: "obscure-lib-b", Version: "2.0", Scope: "compile"},
		{Name: "obscure-lib-c", Version: "3.0", Scope: "compile"},
	}
	toml, unmapped := rw.Rewrite("pkg", "Java", deps)
	if len(unmapped) != 3 {
		t.Fatalf("expected 3 unmapped, got %v", unmapped)
	}
	if !strings.Contains(toml, "[package]") {
		t.Fatalf("missing [package]: %q", toml)
	}
	// When all deps are unmapped, the minimal Cargo.toml contains only
	// [package] (no [dependencies] / [dev-dependencies] sections, since
	// there are no mapped entries to put under them). The spec language
	// is loose here — "minimal" means we omit empty sections.
	if strings.Contains(toml, "[dependencies]") {
		t.Fatalf("unexpected [dependencies] for all-unmapped: %q", toml)
	}
	if strings.Contains(toml, "[dev-dependencies]") {
		t.Fatalf("unexpected [dev-dependencies] for all-unmapped: %q", toml)
	}
	// None of the obscure libs should appear in the toml
	for _, name := range []string{"obscure-lib-a", "obscure-lib-b", "obscure-lib-c"} {
		if strings.Contains(toml, name) {
			t.Fatalf("unmapped dep %q leaked into toml: %q", name, toml)
		}
	}
}
func TestDetectSourceLang(t *testing.T) {
	cases := []struct {
		path     string
		expected string
	}{
		{path: "/proj/pom.xml", expected: "Java"},
		{path: "/proj/package.json", expected: "JavaScript"},
		{path: "/proj/requirements.txt", expected: "Python"},
		{path: "/proj/pyproject.toml", expected: "Python"},
		{path: "/proj/configure.ac", expected: "C"},
		{path: "/proj/Makefile.am", expected: "C"},
		{path: "/proj/CMakeLists.txt", expected: "C"},
		{path: "/proj/meson.build", expected: "C"},
		{path: "/proj/Cargo.toml", expected: "Rust"},
		{path: "/proj/random.txt", expected: "Unknown"},
		{path: "", expected: "Unknown"},
	}
	for _, tc := range cases {
		got := DetectSourceLang(tc.path)
		if got != tc.expected {
			t.Fatalf("DetectSourceLang(%q) = %q, want %q", tc.path, got, tc.expected)
		}
	}
}

func TestFormatDepLine(t *testing.T) {
	// No features — flat string form
	got := formatDepLine("serde", "1.0", nil)
	if got != `serde = "1.0"` {
		t.Fatalf("no-features: got %q", got)
	}
	got = formatDepLine("serde", "1.0", []string{})
	if got != `serde = "1.0"` {
		t.Fatalf("empty features: got %q", got)
	}
	// Single feature — inline-table form
	got = formatDepLine("validator", "0.18", []string{"derive"})
	if got != `validator = { version = "0.18", features = ["derive"] }` {
		t.Fatalf("single feature: got %q", got)
	}
	// Multiple features — inline-table form
	got = formatDepLine("reqwest", "0.12", []string{"json", "rustls-tls"})
	if got != `reqwest = { version = "0.12", features = ["json", "rustls-tls"] }` {
		t.Fatalf("multiple features: got %q", got)
	}
}

func TestQuoteFeatures(t *testing.T) {
	if got := quoteFeatures(nil); got != "" {
		t.Fatalf("nil features: got %q", got)
	}
	if got := quoteFeatures([]string{}); got != "" {
		t.Fatalf("empty features: got %q", got)
	}
	if got := quoteFeatures([]string{"derive"}); got != `"derive"` {
		t.Fatalf("single: got %q", got)
	}
	if got := quoteFeatures([]string{"derive", "serde1"}); got != `"derive", "serde1"` {
		t.Fatalf("multiple: got %q", got)
	}
}

func TestManifestRewriterRoundTripWithCanonicalize(t *testing.T) {
	// Integration: take a known Java source-dep set, run Rewrite, then
	// canonicalizeCargoToml on the output, and verify the result is
	// structurally well-formed Cargo.toml (has [package], [dependencies]
	// sections; no leftover source-language dep aliases).
	deps := []SourceDep{
		{Name: "commons-validator", Version: "1.7.0", Scope: "compile"},
		{Name: "junit", Version: "4.13.2", Scope: "test"},
		{Name: "guava", Version: "32.0.0", Scope: "compile"},
	}
	// Use a package name that does NOT contain source-language substrings
	// so the leak check below is meaningful.
	toml, unmapped := GenerateCargoToml("validator-port", "Java", deps)
	if len(unmapped) != 0 {
		t.Fatalf("expected 0 unmapped for known deps, got %v", unmapped)
	}
	canonical := canonicalizeCargoToml(toml)
	if canonical == "" {
		t.Fatal("canonicalize returned empty string")
	}
	// Structural checks — must still parse as Cargo.toml sections.
	if !strings.Contains(canonical, "[package]") {
		t.Fatalf("canonicalized toml missing [package]: %q", canonical)
	}
	if !strings.Contains(canonical, "[dependencies]") {
		t.Fatalf("canonicalized toml missing [dependencies]: %q", canonical)
	}
	if !strings.Contains(canonical, "[dev-dependencies]") {
		t.Fatalf("canonicalized toml missing [dev-dependencies]: %q", canonical)
	}
	if !strings.Contains(canonical, "validator") {
		t.Fatalf("validator crate missing from canonical: %q", canonical)
	}
	if !strings.Contains(canonical, "rstest") {
		t.Fatalf("rstest crate missing from canonical: %q", canonical)
	}
	if !strings.Contains(canonical, "itertools") {
		t.Fatalf("itertools crate missing from canonical: %q", canonical)
	}
	// No raw source-language names should leak into the final manifest.
	for _, leaked := range []string{"commons-validator", "junit", "guava"} {
		if strings.Contains(canonical, leaked) {
			t.Fatalf("source-language name %q leaked into canonical manifest: %q", leaked, canonical)
		}
	}
}
