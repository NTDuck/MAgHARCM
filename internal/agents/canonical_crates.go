package agents

import "strings"

// canonicalCrateRenames is a hand-curated map of underscore-form crate names
// that the translation/repair model has been observed emitting, mapped to
// their canonical hyphenated form on crates.io. The repair model can't
// reliably guess crate names from underscore to hyphen; this list is the
// ground truth for the small set of crates MAgHARCM translations commonly
// reach for. Add to this list as new patterns are observed.
var canonicalCrateRenames = map[string]string{
	"serde_xml_rs":   "serde-xml-rs",
	"serde_json":     "serde_json",
	"quick_xml":      "quick-xml",
	"yaml_rust":      "yaml-rust",
	"toml_edit":      "toml_edit",
	"regex_set":      "regex-set",
	"nom_supreme":    "nom-supreme",
	"anyhow":         "anyhow",
	"thiserror":      "thiserror",
	"chrono":         "chrono",
	"tokio":          "tokio",
	"hyper":          "hyper",
	"reqwest":        "reqwest",
	"actix_web":      "actix-web",
	"warp":           "warp",
	"clap":           "clap",
	"structopt":      "structopt",
	"env_logger":     "env_logger",
	"tracing":        "tracing",
	"tracing_sub":    "tracing-subscriber",
	"pretty_assert":  "pretty_assertions",
	"assert_matches": "assert_matches",
}

// crateCanonicalHints returns a multiline hint string listing the
// crate-name corrections to pass to the repair prompt. Only includes
// renames where the key (with-underscore) actually differs from the
// value (canonical-hyphenated), so the list is short. The string is
// empty when the target language isn't Rust (no renames apply).
func crateCanonicalHints(targetLang string) string {
	if !strings.EqualFold(targetLang, "Rust") {
		return ""
	}
	var lines []string
	for from, to := range canonicalCrateRenames {
		if from != to {
			lines = append(lines, "- "+from+" -> "+to)
		}
	}
	return strings.Join(lines, "\n")
}

// canonicalizeCargoToml rewrites a small set of known crate-name misspellings
// in a Cargo.toml string. The translation and repair models emit
// underscore-form crate names for crates that have hyphenated forms on
// crates.io (e.g. serde_xml_rs, quick_xml, yaml_rust). Without this
// pass, the validator's `cargo` invocation fails with E0432 unresolved
// import errors. The pass rewrites only the LHS of `key = ...` lines
// to avoid touching multi-line value contexts.
func canonicalizeCargoToml(manifest string) string {
	if manifest == "" {
		return manifest
	}
	for from, to := range canonicalCrateRenames {
		if from == to {
			continue
		}
		manifest = renameCrateKey(manifest, from, to)
	}
	return manifest
}

// renameCrateKey rewrites a `from = ...` Cargo.toml key entry to `to = ...`.
// Operates line-by-line and only touches a line when the LHS token
// (whitespace-stripped) equals `from`. Returns the input unchanged when
// no rewrite occurs.
func renameCrateKey(manifest, from, to string) string {
	lines := strings.Split(manifest, "\n")
	changed := false
	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if !strings.HasPrefix(trimmed, from) {
			continue
		}
		rest := strings.TrimPrefix(trimmed, from)
		if rest == "" || rest[0] != '=' && rest[0] != ' ' && rest[0] != '\t' {
			continue
		}
		// from matches as LHS token — rewrite
		indent := line[:len(line)-len(trimmed)]
		lines[i] = indent + to + rest
		changed = true
	}
	if !changed {
		return manifest
	}
	return strings.Join(lines, "\n")
}
