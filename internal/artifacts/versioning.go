package artifacts

// CurrentSchemaVersion is the schema version stamped on every artifact produced
// by the MAgHARCM pipeline. It is a string-semver value (e.g. "1.0.0") so it
// can travel safely through JSON without numeric coercion. Consumers should
// compare against this constant to decide whether an artifact requires
// migration before use.
//
// This constant implements SOP-Anchored Role-Artifact Schema (paper appendix
// primitive P24, cited as MetaGPT [P42] §3.2.2): each role-produced artifact
// carries an explicit, comparable schema version so downstream roles can
// detect and adapt to schema drift across pipeline iterations.
const CurrentSchemaVersion = "1.0.0"

// SchemaVersioned is implemented by every artifact that carries a
// SchemaVersion field. Producers MUST stamp CurrentSchemaVersion on artifacts
// at construction time; consumers MUST call SchemaVersion() before using any
// artifact produced by an older or newer pipeline.
type SchemaVersioned interface {
	SchemaVersion() string
}

// Annotate is the future wrap helper for adding pipeline-wide metadata
// (timestamps, producer role, run id) to an artifact alongside the schema
// version. It is intentionally a no-op identity for now — the schema-version
// field is set directly on each artifact at construction. Once cross-cutting
// annotation fields stabilize, Annotate will be expanded to attach them
// without modifying every producer site.
//
// Passing nil returns nil.
func Annotate(v any) any {
	return v
}
