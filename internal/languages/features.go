package languages

import (
	"strings"
	"sync"
)

// FeatureKind names a canonical, language-agnostic construct that the
// translator must reconcile across source and target languages.
// (Cross-language primitives used during planning, not the AST node
// categories in spec.go which describe tree-sitter node kinds.)
type FeatureKind string

const (
	FeatureErrorHandling FeatureKind = "error_handling" // checked exceptions, Result<T,E>, panic/Result
	FeatureGenerics      FeatureKind = "generics"       // parametric types (Java, C#, Rust, Go 1.18+, Python type hints)
	FeatureAsync         FeatureKind = "async"          // goroutines, async/await, futures, threads
	FeatureOwnership     FeatureKind = "ownership"      // Rust borrow-checker, ARC, manual malloc/free
	FeatureGC            FeatureKind = "gc"             // tracing garbage collector (Go, Java, Python, C#)
	FeatureTraits        FeatureKind = "traits"         // Rust traits, Java interfaces, Go interfaces, C++ concepts
	FeatureModules       FeatureKind = "modules"        // cargo crates, go packages, java packages, python modules
	FeatureMacros        FeatureKind = "macros"         // Rust macro_rules!, C macros, Java annotations
	FeatureNullSafety    FeatureKind = "null_safety"    // Option<T> (Rust), @Nullable (Java), *T vs Optional (Go pre-1.18), None (Python)
	FeaturePatternMatch  FeatureKind = "pattern_match"  // match arms, switch, instanceof + cast
	FeatureClosures      FeatureKind = "closures"       // anonymous functions / lambdas
	FeatureReflection    FeatureKind = "reflection"     // runtime type inspection
)

// FeatureTemplate is the target-language snippet for a given feature kind.
// Templates use Go text/template syntax with a small surface:
//   {{.PackageName}} — resolved target package
//   {{.Name}}        — symbol name being constructed
//   {{.Lifetime}}    — 'a / 'static — only relevant for Rust
type FeatureTemplate struct {
	Kind        FeatureKind `json:"kind" yaml:"kind"`
	TargetLang  string      `json:"target_lang" yaml:"target_lang"`
	Template    string      `json:"template" yaml:"template"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
}

// FeatureMap holds (kind x targetLang) -> template lookup.
type FeatureMap struct {
	entries map[FeatureKind]map[string]FeatureTemplate // kind -> lang -> template
}

// NewFeatureMap returns an empty feature map.
func NewFeatureMap() *FeatureMap {
	return &FeatureMap{entries: map[FeatureKind]map[string]FeatureTemplate{}}
}

// Add inserts a template.
func (m *FeatureMap) Add(t FeatureTemplate) {
	if m.entries[t.Kind] == nil {
		m.entries[t.Kind] = map[string]FeatureTemplate{}
	}
	m.entries[t.Kind][strings.ToLower(t.TargetLang)] = t
}

// Lookup returns the template for (kind, targetLang) and whether it was found.
func (m *FeatureMap) Lookup(kind FeatureKind, targetLang string) (FeatureTemplate, bool) {
	if m.entries[kind] == nil {
		return FeatureTemplate{}, false
	}
	t, ok := m.entries[kind][strings.ToLower(targetLang)]
	return t, ok
}

// Kinds returns the set of feature kinds registered.
func (m *FeatureMap) Kinds() []FeatureKind {
	out := make([]FeatureKind, 0, len(m.entries))
	for k := range m.entries {
		out = append(out, k)
	}
	return out
}

// builtInFeatureMap is the static, embedded cross-language feature table.
// Initial coverage: Rust, Go, Java, Python, C#, JavaScript/TypeScript.
var builtInFeatureMap *FeatureMap
var featureMapOnce sync.Once

// GetFeatureMap returns the singleton built-in feature map.
func GetFeatureMap() *FeatureMap {
	featureMapOnce.Do(func() { builtInFeatureMap = newBuiltInFeatureMap() })
	return builtInFeatureMap
}

func newBuiltInFeatureMap() *FeatureMap {
	m := NewFeatureMap()
	// --- Error handling ---
	m.Add(FeatureTemplate{FeatureErrorHandling, "rust", `// Idiomatic Rust: Result<T, E> with the '?' operator for propagation.
{{if .ReturnErr}}pub fn {{.Name}}(...) -> Result<T, E> {
    ...
    other_fn(...)?;
    Ok(result)
}{{else}}pub fn {{.Name}}(...) -> T {{
    // infallible path
}}{{end}}`, "Result<T, E> + ? operator"})
	m.Add(FeatureTemplate{FeatureErrorHandling, "go", `// Idiomatic Go: explicit (T, error) returns with sentinel wrapping.
func {{.Name}}(...) (T, error) {
    ...
    if err != nil {
        return zero, fmt.Errorf("...: %w", err)
    }
    return result, nil
}`, "Multi-value return with error"})
	m.Add(FeatureTemplate{FeatureErrorHandling, "java", `// Idiomatic Java: checked exceptions + try/catch.
public T {{.Name}}(...) throws E {
    try {
        ...
    } catch (E e) {
        throw new RuntimeException("...", e);
    }
}`, "Checked exceptions"})
	m.Add(FeatureTemplate{FeatureErrorHandling, "python", `// Idiomatic Python: raise + custom exception classes.
def {{.Name}}(...) -> T:
    ...
    raise ValueError("...")
    return result`, "EAFP with raise"})

	// --- Generics ---
	m.Add(FeatureTemplate{FeatureGenerics, "rust", `pub struct {{.Name}}<T> { /* ... */ }`, "Native parametric types"})
	m.Add(FeatureTemplate{FeatureGenerics, "go", `type {{.Name}}[T any] struct{ /* ... */ }`, "Go 1.18+ type parameters"})
	m.Add(FeatureTemplate{FeatureGenerics, "java", `public class {{.Name}}<T> { /* ... */ }`, "Type erasure generics"})
	m.Add(FeatureTemplate{FeatureGenerics, "python", `T = TypeVar('T')\nclass {{.Name}}(Generic[T]): ...`, "TypeVar + Generic[T]"})
	m.Add(FeatureTemplate{FeatureGenerics, "typescript", `class {{.Name}}<T> { /* ... */ }`, "Structural generics"})

	// --- Async ---
	m.Add(FeatureTemplate{FeatureAsync, "rust", `pub async fn {{.Name}}(...) -> T { ... }`, "async/await + tokio runtime"})
	m.Add(FeatureTemplate{FeatureAsync, "go", `func {{.Name}}(...) (T, error) { /* goroutine-friendly */ }`, "Goroutines + channels"})
	m.Add(FeatureTemplate{FeatureAsync, "java", `CompletableFuture<T> {{.Name}}(...) { ... }`, "CompletableFuture / virtual threads"})
	m.Add(FeatureTemplate{FeatureAsync, "python", `async def {{.Name}}(...) -> T: ...`, "asyncio coroutines"})
	m.Add(FeatureTemplate{FeatureAsync, "typescript", `async {{.Name}}(...): Promise<T> { ... }`, "Promise + async/await"})

	// --- Ownership ---
	m.Add(FeatureTemplate{FeatureOwnership, "rust", `// Default: ownership + borrow checker.
fn {{.Name}}(s: &str) -> &{{.Lifetime}}str { ... }`, "&T, &mut T, Box<T>, Arc<T>"})

	// --- GC ---
	m.Add(FeatureTemplate{FeatureGC, "go", "// Go has tracing GC; no explicit free.\n", "runtime.GC"})
	m.Add(FeatureTemplate{FeatureGC, "java", "// Java has tracing GC.\n", "System.gc()"})
	m.Add(FeatureTemplate{FeatureGC, "python", "// CPython refcounting + cycle GC.\n", "gc.collect()"})
	m.Add(FeatureTemplate{FeatureGC, "csharp", "// .NET has tracing GC.\n", "GC.Collect()"})

	// --- Traits / interfaces ---
	m.Add(FeatureTemplate{FeatureTraits, "rust", `pub trait {{.Name}} {
    fn method(&self) -> T;
}`, "trait + impl"})
	m.Add(FeatureTemplate{FeatureTraits, "go", `type {{.Name}} interface {
    Method() T
}`, "implicit interface satisfaction"})
	m.Add(FeatureTemplate{FeatureTraits, "java", `public interface {{.Name}} {
    T method();
}`, "interface + implements"})
	m.Add(FeatureTemplate{FeatureTraits, "python", `class {{.Name}}(Protocol):
    def method(self) -> T: ...`, "Protocol / ABC"})
	m.Add(FeatureTemplate{FeatureTraits, "typescript", `interface {{.Name}} {
    method(): T;
}`, "structural interface"})

	// --- Null safety ---
	m.Add(FeatureTemplate{FeatureNullSafety, "rust", `Option<T>`, "None | Some(t)"})
	m.Add(FeatureTemplate{FeatureNullSafety, "java", `Optional<T>`, "@Nullable / Optional<T>"})
	m.Add(FeatureTemplate{FeatureNullSafety, "go", `// pre-1.18: *T with nil check
// 1.18+: generic *T or custom Option`, "nil pointers, no syntactic sugar"})
	m.Add(FeatureTemplate{FeatureNullSafety, "python", `Optional[T]`, "None | obj"})
	m.Add(FeatureTemplate{FeatureNullSafety, "typescript", `T | null`, "T | null | undefined"})

	// --- Pattern matching ---
	m.Add(FeatureTemplate{FeaturePatternMatch, "rust", `match expr {
    Pattern1 => result1,
    Pattern2 => result2,
    _ => default,
}`, "Exhaustive match"})
	m.Add(FeatureTemplate{FeaturePatternMatch, "go", `switch v := x.(type) {
case T1: ...
case T2: ...
default: ...
}`, "Type switch"})
	m.Add(FeatureTemplate{FeaturePatternMatch, "java", `switch (x) {
    case A -> ...
    case B -> ...
    default -> ...
}`, "Switch expressions (Java 21+)"})
	m.Add(FeatureTemplate{FeaturePatternMatch, "python", `match x:
    case pattern1: ...
    case pattern2: ...`, "match/case (3.10+)"})

	// --- Closures ---
	m.Add(FeatureTemplate{FeatureClosures, "rust", `|x: i32| -> i32 { x + 1 }`, "Closures with capture modes"})
	m.Add(FeatureTemplate{FeatureClosures, "go", `func(x int) int { return x + 1 }`, "Function values"})
	m.Add(FeatureTemplate{FeatureClosures, "java", `(x) -> x + 1`, "Lambda expressions"})
	m.Add(FeatureTemplate{FeatureClosures, "python", `lambda x: x + 1`, "lambda / def"})
	m.Add(FeatureTemplate{FeatureClosures, "javascript", `(x) => x + 1`, "Arrow functions"})

	// --- Modules ---
	m.Add(FeatureTemplate{FeatureModules, "rust", `// Cargo: src/lib.rs declares 'pub mod' re-exports
pub mod submodule;
pub use submodule::*;`, "Cargo crates + modules"})
	m.Add(FeatureTemplate{FeatureModules, "go", `// Go: same-dir files share the same 'package' name
package {{.PackageName}}`, "Packages (one per dir)"})
	m.Add(FeatureTemplate{FeatureModules, "java", `package com.example;`, "Packages + Maven coords"})
	m.Add(FeatureTemplate{FeatureModules, "python", `# __init__.py marks packages; PEP 420 namespace packages
from .submodule import *`, "Packages via __init__.py"})


	return m
}
