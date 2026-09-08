package agents

import (
	"strings"
	"testing"

	"MAgHARCM/internal/compiletime"
)

func normalizerState() *compiletime.State {
	return &compiletime.State{
		Task: compiletime.Task{
			TargetDir:  ".artifacts/GildedRose-Refactoring-Kata/rust",
			TargetLang: "Rust",
			Toolchain:  "cargo",
		},
	}
}

// Regression: the repair SLM repeatedly emitted `mod tests { ... }` with no
// use of the package root, producing error[E0425] on every repair iteration
// (checkpoints iter-0017..0019 of run assets-samples-GildedRose-Refactoring-Kata-C).
func TestNormalizeRustTestImportsInjectsMissingGlob(t *testing.T) {
	in := `#[cfg(test)]
mod tests {
    #[test]
    fn test_aged_brie() {
        let mut item = init_item("Aged Brie", 2, 0);
        update_quality(&mut item);
        assert_eq!(item.quality, 1);
    }
}
`
	out := NormalizeRustTestImports(in, normalizerState())
	if !strings.Contains(out, "use gildedrose_refactoring_kata::*;") {
		t.Fatalf("expected injected glob import, got:\n%s", out)
	}
	// Import must land inside the mod block, before first symbol use.
	importIdx := strings.Index(out, "use gildedrose_refactoring_kata::*;")
	firstUse := strings.Index(out, "init_item(")
	if importIdx == -1 || firstUse == -1 || importIdx > firstUse {
		t.Fatalf("import must precede first symbol use:\n%s", out)
	}
}

func TestNormalizeRustTestImportsKeepsExistingImport(t *testing.T) {
	in := `use gildedrose_refactoring_kata::*;

#[cfg(test)]
mod tests {
    use gildedrose_refactoring_kata::*;
    #[test]
    fn test_aged_brie() {
        let mut item = init_item("Aged Brie", 2, 0);
    }
}
`
	out := NormalizeRustTestImports(in, normalizerState())
	if got := strings.Count(out, "use gildedrose_refactoring_kata::*;"); got != 2 {
		t.Fatalf("expected exactly 2 import lines (file-top + in-mod), got %d:\n%s", got, out)
	}
}

func TestNormalizeRustTestImportsIgnoresNonTestsContent(t *testing.T) {
	in := `// Plain comment mentioning Item but no mod block.
fn helper() -> i32 { 42 }
`
	out := NormalizeRustTestImports(in, normalizerState())
	if strings.Contains(out, "use gildedrose_refactoring_kata::*;") {
		t.Fatalf("must not inject without a mod tests block:\n%s", out)
	}
}

func TestNormalizeRustTestImportsHandlesNestedBraces(t *testing.T) {
	in := `#[cfg(test)]
mod tests {
    #[test]
    fn nested() {
        if true {
            let mut item = init_item("x", 1, 1);
            update_quality(&mut item);
        }
    }
}
`
	out := NormalizeRustTestImports(in, normalizerState())
	importIdx := strings.Index(out, "use gildedrose_refactoring_kata::*;")
	firstUse := strings.Index(out, "init_item(")
	if importIdx == -1 || firstUse == -1 || importIdx > firstUse {
		t.Fatalf("nested-brace case must still inject before first use:\n%s", out)
	}
	// Exactly one injection, not one per nested close.
	if got := strings.Count(out, "use gildedrose_refactoring_kata::*;"); got != 1 {
		t.Fatalf("expected single injection, got %d:\n%s", got, out)
	}
}
