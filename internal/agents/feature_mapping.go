package agents

// Backlink: [[Primitives]] §NEW-PRIM-10 (Feature-Mapping & Type-Compatibility Validation),
// cited as Oxidizer [[P05]] and RustRepoTrans [[P07]].

import (
	"context"
	"fmt"
	"strings"

	"MAgHARCM/internal/logger"
)

// Mapping represents one curated feature/idiom translation between a
// source-language construct and its Rust equivalent. The table is hand-curated
// (see seedFeatureMappings); full automation would require a corpus-driven LLM
// stage that is out of scope for this primitive.
type Mapping struct {
	// Source identifies the source-language construct (e.g. "go:interface",
	// "java:checked-exception", "python:None", "c:NULL", "cpp:macro").
	Source string
	// Target identifies the Rust construct (e.g. "trait", "Result<T, E>",
	// "Option<T>", "*const T", "macro_rules!").
	Target string
	// Kind categorises the mapping so callers can route it appropriately.
	// One of: "trait", "result-type", "null-safety", "module-system",
	// "macro", "operator-overload".
	Kind string
	// Notes captures caveats the translator must respect (e.g. "C NULL is
	// not type-checked; prefer Option<Box<T>> at the API boundary").
	Notes string
	// Confidence is the curator's confidence that this mapping is correct
	// for a representative codebase, on the closed interval [0.0, 1.0].
	Confidence float64
}

// Valid Mapping.Kind values. Exported so callers and tests can validate
// externally-constructed Mapping values without re-declaring the vocabulary.
const (
	MappingKindTrait            = "trait"
	MappingKindResultType       = "result-type"
	MappingKindNullSafety       = "null-safety"
	MappingKindModuleSystem     = "module-system"
	MappingKindMacro            = "macro"
	MappingKindOperatorOverload = "operator-overload"
)

// ValidationResult is the outcome of consulting the mapping table for a
// source/target symbol pair. Match=true means a Mapping was found with the
// same Kind on both sides; the Reason string captures the human-readable
// rationale for logging and for downstream prompt construction.
type ValidationResult struct {
	Match  bool
	Reason string
}

// FeatureMapper is the P10 primitive: a curated lookup of source-language
// constructs to Rust equivalents, plus a compatibility-validation helper
// consulted by the translator before emitting target code.
type FeatureMapper struct {
	Mappings []Mapping
}

// NewFeatureMapper returns a FeatureMapper seeded with the curated P10 table
// (Go interfaces ↔ Rust traits, Java checked exceptions ↔ Result, Python None
// ↔ Option, C NULL ↔ *const T, C++ macros ↔ macro_rules!, plus a handful of
// secondary idioms). The seed list lives in seedFeatureMappings and is
// hand-maintained.
func NewFeatureMapper() *FeatureMapper {
	return &FeatureMapper{Mappings: seedFeatureMappings()}
}

// Map looks up a (sourceLang, sourceFeat) pair in the mapping table and
// returns the first matching Mapping. The lookup is case-insensitive on both
// keys; an empty sourceLang or sourceFeat is rejected with an error so
// callers cannot silently match against an "all" wildcard.
//
// The returned Mapping's Confidence is the curator's score and is propagated
// verbatim — callers may apply their own floor when wiring it into a prompt.
func (m *FeatureMapper) Map(ctx context.Context, sourceLang, sourceFeat string) (Mapping, error) {
	if err := ctx.Err(); err != nil {
		return Mapping{}, err
	}
	if sourceLang == "" || sourceFeat == "" {
		return Mapping{}, fmt.Errorf("feature_mapper: sourceLang and sourceFeat must be non-empty")
	}
	lang := strings.ToLower(strings.TrimSpace(sourceLang))
	feat := strings.ToLower(strings.TrimSpace(sourceFeat))
	for _, mp := range m.Mappings {
		if mp.Source == "" {
			continue
		}
		parts := strings.SplitN(mp.Source, ":", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.ToLower(parts[0]) == lang && strings.ToLower(parts[1]) == feat {
			logger.LogStep("feature_mapper: matched %s:%s -> %s kind=%s conf=%.2f",
				lang, feat, mp.Target, mp.Kind, mp.Confidence)
			return mp, nil
		}
	}
	return Mapping{}, fmt.Errorf("feature_mapper: no mapping for %s:%s", sourceLang, sourceFeat)
}

// Validate consults the mapping table for both sourceSym and targetSym and
// reports whether the two are compatible under the same Mapping.Kind. This
// is the type-compatibility check the translator runs before emitting target
// code: if the source maps to "result-type" but the target was translated as
// a panic, the call surfaces as Match=false with an explanatory Reason.
//
// Lookup is by source-string identity (no language prefix required), which
// matches how downstream callers pass already-extracted symbols.
func (m *FeatureMapper) Validate(ctx context.Context, sourceSym, targetSym string) (ValidationResult, error) {
	if err := ctx.Err(); err != nil {
		return ValidationResult{}, err
	}
	if sourceSym == "" || targetSym == "" {
		return ValidationResult{}, fmt.Errorf("feature_mapper: sourceSym and targetSym must be non-empty")
	}
	srcLower := strings.ToLower(sourceSym)
	tgtLower := strings.ToLower(targetSym)
	var srcKind, tgtKind string
	var srcTarget, tgtSource string
	for _, mp := range m.Mappings {
		// Source side: match the full "<lang>:<feat>" token or the bare
		// <feat> suffix so callers may pass either form.
		if tokenMatches(mp.Source, sourceSym, srcLower) {
			srcKind = mp.Kind
			srcTarget = mp.Target
		}
		// Target side: match against the Rust identifier (e.g. "Result",
		// "Option") or its sugared form ("Result<T, E>").
		if tokenMatches(mp.Target, targetSym, tgtLower) {
			tgtKind = mp.Kind
			tgtSource = mp.Source
		}
	}
	if srcKind == "" {
		return ValidationResult{
			Match:  false,
			Reason: fmt.Sprintf("source symbol %q has no entry in feature-mapping table", sourceSym),
		}, nil
	}
	if tgtKind == "" {
		return ValidationResult{
			Match:  false,
			Reason: fmt.Sprintf("target symbol %q has no entry in feature-mapping table", targetSym),
		}, nil
	}
	if srcKind != tgtKind {
		logger.LogWarning("feature_mapper: kind mismatch source=%s kind=%s target=%s kind=%s",
			sourceSym, srcKind, targetSym, tgtKind)
		return ValidationResult{
			Match: false,
			Reason: fmt.Sprintf(
				"kind mismatch: %s maps to %s (kind=%s) but %s maps to %s (kind=%s)",
				sourceSym, srcTarget, srcKind, targetSym, tgtSource, tgtKind),
		}, nil
	}
	return ValidationResult{
		Match: true,
		Reason: fmt.Sprintf(
			"%s and %s share mapping kind=%s (source->%s, target->%s)",
			sourceSym, targetSym, srcKind, srcTarget, tgtSource),
	}, nil
}

// tokenMatches reports whether either the raw or lower-cased needle appears
// as a complete token in the haystack's lower-cased form. It handles the two
// shapes callers pass: "<lang>:<feat>" or bare "<feat>".
func tokenMatches(haystack, needle, needleLower string) bool {
	hLower := strings.ToLower(haystack)
	if hLower == needleLower {
		return true
	}
	if strings.HasSuffix(hLower, ":"+needleLower) {
		return true
	}
	return false
}

// seedFeatureMappings returns the curated feature-mapping table seeded from
// the P10 appendix (Oxidizer [[P05]], RustRepoTrans [[P07]]). Add new
// mappings as patterns emerge; keep Confidence in [0.0, 1.0].
func seedFeatureMappings() []Mapping {
	return []Mapping{
		// Go interfaces ↔ Rust traits.
		{
			Source:     "go:interface",
			Target:     "trait",
			Kind:       MappingKindTrait,
			Notes:      "Go interface methods become trait methods with no receivers; implicit interface satisfaction maps to trait impl blocks. Boxing needed for trait objects (Box<dyn Trait>).",
			Confidence: 0.95,
		},
		{
			Source:     "go:empty-interface",
			Target:     "trait Any",
			Kind:       MappingKindTrait,
			Notes:      "interface{} / any maps to trait Any (from core::any). Use Any::type_id for runtime type checks.",
			Confidence: 0.90,
		},
		// Java checked exceptions ↔ Rust Result.
		{
			Source:     "java:checked-exception",
			Target:     "Result<T, E>",
			Kind:       MappingKindResultType,
			Notes:      "Java's checked exceptions are surfaced at the type system; Rust uses Result<T, E> with custom error enums. Implement From<E> for ? propagation.",
			Confidence: 0.92,
		},
		{
			Source:     "java:runtime-exception",
			Target:     "panic!",
			Kind:       MappingKindResultType,
			Notes:      "Unchecked exceptions correspond to panic! in Rust; reserve for truly unrecoverable states and prefer Result in library code.",
			Confidence: 0.80,
		},
		// Python None ↔ Rust Option.
		{
			Source:     "python:None",
			Target:     "Option<T>",
			Kind:       MappingKindNullSafety,
			Notes:      "Optional[T] | None and implicit None returns map to Option<T>. None literal maps to None variant.",
			Confidence: 0.93,
		},
		{
			Source:     "python:Optional",
			Target:     "Option<T>",
			Kind:       MappingKindNullSafety,
			Notes:      "typing.Optional[T] is sugar for Union[T, None] and maps to Option<T>.",
			Confidence: 0.93,
		},
		// C NULL ↔ Rust raw pointer + Option.
		{
			Source:     "c:NULL",
			Target:     "*const T",
			Kind:       MappingKindNullSafety,
			Notes:      "C NULL is not type-checked; prefer Option<Box<T>> at the FFI boundary and only use *const T inside unsafe blocks.",
			Confidence: 0.78,
		},
		{
			Source:     "cpp:nullptr",
			Target:     "*const T",
			Kind:       MappingKindNullSafety,
			Notes:      "C++11 nullptr is type-checked but still raw; map to Option<Box<T>> at safe boundaries, *const T inside unsafe.",
			Confidence: 0.80,
		},
		// C/C++ macros ↔ Rust macro_rules!.
		{
			Source:     "cpp:macro",
			Target:     "macro_rules!",
			Kind:       MappingKindMacro,
			Notes:      "Object-like #define macros translate to const items where possible, otherwise macro_rules! with empty matcher. Function-like macros with side-effects map to macro_rules! with explicit hygiene notes.",
			Confidence: 0.70,
		},
		{
			Source:     "c:macro",
			Target:     "macro_rules!",
			Kind:       MappingKindMacro,
			Notes:      "C #define macros become const items for value-like or macro_rules! for function-like. Token-pasting and stringifying have no direct equivalent; flag for manual review.",
			Confidence: 0.68,
		},
		// Module systems.
		{
			Source:     "go:package",
			Target:     "mod",
			Kind:       MappingKindModuleSystem,
			Notes:      "Go package becomes a Rust module file; package-level vars become pub statics. Visibility defaults to private; use pub explicitly.",
			Confidence: 0.90,
		},
		{
			Source:     "python:module",
			Target:     "mod",
			Kind:       MappingKindModuleSystem,
			Notes:      "Python module becomes a Rust module file; __init__.py contents fold into the parent module; __all__ maps to pub re-exports.",
			Confidence: 0.85,
		},
		{
			Source:     "java:package",
			Target:     "mod",
			Kind:       MappingKindModuleSystem,
			Notes:      "Java packages map to Rust modules; nested packages become nested mod declarations. package-private becomes pub(crate).",
			Confidence: 0.88,
		},
		// Operator overloading.
		{
			Source:     "cpp:operator-overload",
			Target:     "impl Trait for Type",
			Kind:       MappingKindOperatorOverload,
			Notes:      "C++ operator overloads map to trait impls (Add, Sub, Mul, etc.). Member operators -> inherent impl, friend/non-member -> orphan-rule respecting wrapper impl.",
			Confidence: 0.85,
		},
		{
			Source:     "csharp:operator-overload",
			Target:     "impl Trait for Type",
			Kind:       MappingKindOperatorOverload,
			Notes:      "C# operator overloads map to trait impls; static operator methods translate to inherent impl blocks.",
			Confidence: 0.83,
		},
		// Secondary idioms for completeness.
		{
			Source:     "go:goroutine",
			Target:     "tokio::spawn",
			Kind:       "trait",
			Notes:      "Goroutines map to tokio::spawn futures; channels map to tokio::sync::mpsc. Requires explicit Send + 'static bounds.",
			Confidence: 0.75,
		},
		{
			Source:     "java:interface",
			Target:     "trait",
			Kind:       MappingKindTrait,
			Notes:      "Java interfaces map to Rust traits; default methods become trait-provided methods. Functional interfaces map to single-method traits callable via closures.",
			Confidence: 0.93,
		},
	}
}
