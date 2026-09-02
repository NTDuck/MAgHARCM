package agents

import (
	"testing"
)

func TestDetectTestWeakening(t *testing.T) {
	v := NewValidatorAgent(nil, "")

	prev := map[string]string{
		"tests/test_basic.rs": `
			#[test]
			fn test_foo() {
				assert_eq!(foo(), 42);
				assert!(is_valid());
				assert_ne!(bar(), 0);
			}
		`,
	}

	// 1. Same or more assertions -> not weakened
	same := map[string]string{
		"tests/test_basic.rs": `
			#[test]
			fn test_foo() {
				assert_eq!(foo(), 42);
				assert!(is_valid());
				assert_ne!(bar(), 0);
			}
		`,
	}
	weakened, reasons := v.DetectTestWeakening(prev, same)
	if weakened {
		t.Errorf("expected not weakened, got reasons: %v", reasons)
	}

	// 2. Reduced assertions -> weakened
	fewer := map[string]string{
		"tests/test_basic.rs": `
			#[test]
			fn test_foo() {
				assert_eq!(foo(), 42);
			}
		`,
	}
	weakened, reasons = v.DetectTestWeakening(prev, fewer)
	if !weakened {
		t.Errorf("expected weakened when assertions reduced from 3 to 1")
	}

	// 3. Removed test file -> weakened
	removed := map[string]string{}
	weakened, reasons = v.DetectTestWeakening(prev, removed)
	if !weakened {
		t.Errorf("expected weakened when test file is removed")
	}
}
