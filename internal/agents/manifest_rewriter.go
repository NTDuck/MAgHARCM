package agents

// Backlink: [[Primitives]] §NEW-PRIM-31 (Source-to-Target Manifest Rewriter).

import (
	"fmt"
	"strings"
)

// CrateMapping maps a source-language dependency identifier to a target
// Cargo crate name + version + feature flags. The mappings are seeded
// from P60's empirical analysis (Java commons-* libraries → Rust equivalents).
// This is a curated, hand-maintained list — automatic mapping would require
// a corpus + LLM, which is out of scope for this primitive.
type CrateMapping struct {
	Crate    string   // canonical crate name (e.g. "serde")
	Version  string   // semver range (e.g. "1.0")
	Features []string // optional feature flags (e.g. ["derive"])
	Optional bool     // true = include only when explicitly requested
}

// SourceDep is a minimal representation of a source-language dependency,
// independent of the source manifest format. The planner translates the
// raw POM/package.json/etc. into a list of SourceDeps before calling the
// manifest-rewriter.
type SourceDep struct {
	Name     string // e.g. "commons-validator"
	Version  string // e.g. "1.11.0" or ">=1.0"
	Scope    string // "compile", "test", "runtime", etc.
	Optional bool
}

// ManifestRewriter produces minimal Cargo.toml content from a list of
// SourceDeps using the CrateMapping table. The output is intentionally
// minimal — full transitive dependency resolution is delegated to the
// post-translation `cargo` invocation + cargo-deny/cargo-geiger verification.
type ManifestRewriter struct {
	Mappings map[string]CrateMapping
}

// NewManifestRewriter returns a ManifestRewriter seeded with the P60
// empirical mappings (commons-validator → validator, commons-codec → ...,
// etc.). The seed list is hand-curated and lives in seedMappings() below.
func NewManifestRewriter() *ManifestRewriter {
	return &ManifestRewriter{Mappings: seedMappings()}
}

// GenerateCargoToml is the public entry point for the manifest-rewriter.
// It constructs a ManifestRewriter with the default P60 seed mappings and
// invokes Rewrite. The planner (or future round) calls this directly when
// it needs a Cargo.toml from a source-language dep set without retaining
// a long-lived ManifestRewriter handle.
//
// Returns the minimal Cargo.toml string and the list of source dep names
// that did not have a known mapping (caller decides how to handle these).
func GenerateCargoToml(pkgName, sourceLang string, deps []SourceDep) (string, []string) {
	return NewManifestRewriter().Rewrite(pkgName, sourceLang, deps)
}

// Rewrite produces a minimal Cargo.toml string from the given source deps.
// Only deps that have a known mapping are included; unmapped deps are
// flagged in the returned Unmapped slice for the planner to handle.
// The returned Cargo.toml uses edition 2021 (safest default for MAgHARCM's
// toolchain; edition 2024 is newer but not yet broadly adopted in the
// translation model's training data).
func (m *ManifestRewriter) Rewrite(pkgName string, sourceLang string, deps []SourceDep) (cargoToml string, unmapped []string) {
	var b strings.Builder
	b.WriteString("[package]\n")
	b.WriteString(fmt.Sprintf("name = %q\n", sanitizePkgName(pkgName)))
	b.WriteString("version = \"0.1.0\"\n")
	b.WriteString("edition = \"2021\"\n\n")

	var depsLines, devDepsLines []string
	for _, dep := range deps {
		mapping, ok := m.Mappings[dep.Name]
		if !ok {
			unmapped = append(unmapped, dep.Name)
			continue
		}
		if dep.Optional && !mapping.Optional {
			// Optional deps that map to mandatory crates are still included
			// (the planner can decide whether to actually use them).
		}
		line := FormatDepLine(mapping.Crate, mapping.Version, mapping.Features)
		if dep.Scope == "test" {
			devDepsLines = append(devDepsLines, line)
		} else {
			depsLines = append(depsLines, line)
		}
	}
	if len(depsLines) > 0 {
		b.WriteString("[dependencies]\n")
		b.WriteString(strings.Join(depsLines, "\n"))
		b.WriteString("\n")
	}
	if len(devDepsLines) > 0 {
		b.WriteString("\n[dev-dependencies]\n")
		b.WriteString(strings.Join(devDepsLines, "\n"))
		b.WriteString("\n")
	}
	return b.String(), unmapped
}

// seedMappings returns the curated crate-mapping table seeded from P60.
// Keyed by the source-language dep name; value is the Rust equivalent.
// Add new mappings as patterns emerge.
func seedMappings() map[string]CrateMapping {
	return map[string]CrateMapping{
		// Apache Commons → Rust equivalents (P60 §3 verified)
		"commons-validator":   {Crate: "validator", Version: "0.18", Features: []string{"derive"}},
		"commons-lang3":       {Crate: "chrono", Version: "0.4"}, // date/time slice
		"commons-codec":       {Crate: "base64", Version: "0.22"},
		"commons-csv":         {Crate: "csv", Version: "1.3"},
		"commons-io":          {Crate: "tempfile", Version: "3"}, // closest fit
		"commons-collections": {Crate: "indexmap", Version: "2"},
		"commons-beanutils":   {Crate: "serde", Version: "1", Features: []string{"derive"}},
		"commons-digester":    {Crate: "quick-xml", Version: "0.36", Features: []string{"serialize"}},
		"commons-logging":     {Crate: "log", Version: "0.4"},
		// JUnit → Rust test framework
		"junit": {Crate: "rstest", Version: "0.21", Optional: true},
		// Common utilities
		"guava":            {Crate: "itertools", Version: "0.12"},
		"slf4j-api":        {Crate: "log", Version: "0.4"},
		"jackson-databind": {Crate: "serde_json", Version: "1", Features: []string{"derive"}},
		// Python equivalents (for completeness)
		"requests": {Crate: "reqwest", Version: "0.12", Features: []string{"json"}},
		"numpy":    {Crate: "ndarray", Version: "0.15"},
		"pytest":   {Crate: "rstest", Version: "0.21", Optional: true},
	}
}

// formatDepLine formats a single Cargo.toml dependency line.
func FormatDepLine(crate, version string, features []string) string {
	if len(features) == 0 {
		return fmt.Sprintf("%s = %q", crate, version)
	}
	return fmt.Sprintf("%s = { version = %q, features = [%s] }", crate, version, QuoteFeatures(features))
}

// quoteFeatures formats a slice of feature names as a comma-separated list of
// quoted strings: ["derive", "serde1"].
func QuoteFeatures(features []string) string {
	quoted := make([]string, len(features))
	for i, f := range features {
		quoted[i] = fmt.Sprintf("%q", f)
	}
	return strings.Join(quoted, ", ")
}

// sanitizePkgName converts an arbitrary package name into a Cargo-compatible
// identifier (lowercase ASCII letters, digits, underscores, hyphens).
// Returns "magharcm_translated" if the input is empty after sanitization.
func sanitizePkgName(name string) string {
	if name == "" {
		return "magharcm_translated"
	}
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}

// DetectSourceLang is a best-effort heuristic for the source language from
// a manifest file extension or path fragment. Returns the canonical
// language name (e.g. "Java", "Python", "JavaScript") or "Unknown".
func DetectSourceLang(manifestPath string) string {
	lower := strings.ToLower(manifestPath)
	switch {
	case strings.HasSuffix(lower, "pom.xml"):
		return "Java"
	case strings.HasSuffix(lower, "package.json"):
		return "JavaScript"
	case strings.HasSuffix(lower, "requirements.txt") || strings.HasSuffix(lower, "pyproject.toml"):
		return "Python"
	case strings.HasSuffix(lower, "configure.ac") || strings.HasSuffix(lower, "makefile.am"):
		return "C"
	case strings.HasSuffix(lower, "cmakelists.txt") || strings.HasSuffix(lower, "meson.build"):
		return "C"
	case strings.HasSuffix(lower, "cargo.toml"):
		return "Rust" // already in target — pass through
	default:
		return "Unknown"
	}
}
