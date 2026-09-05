package agents

import (
	"testing"

	"MAgHARCM/internal/compiletime"
)

// TestPartition_L1FlaggedAsModularityTrap verifies that an L1 element whose
// stable-since year is older than compiletime.ModularityTrapYears is
// surfaced as IsTrap=true with a populated TrapReason. The shared root
// receives many callers (high fan-in) and calls one leaf (low fan-out) —
// exactly the L1 contract per Baldwin & Clark.
func TestPartition_L1FlaggedAsModularityTrap(t *testing.T) {
	oldYear := currentYear() - compiletime.ModularityTrapYears - 1
	recentYear := currentYear() - 1

	stable := map[string]int{
		"shared_root":  oldYear,
		"shared_root2": recentYear,
	}

	// Two isolated chains: caller -> root -> leaf_a. Each root has fan-in 1
	// (one caller) — that's not enough. Need fan-in >=2 so each caller also
	// depends on the other root.
	elements := []string{"shared_root", "shared_root2", "caller1", "caller2", "leaf_a", "leaf_b"}
	deps := map[string][]string{
		"shared_root":  {"leaf_a"},
		"shared_root2": {"leaf_b"},
		"caller1":      {"shared_root", "shared_root2"},
		"caller2":      {"shared_root", "shared_root2"},
		"leaf_a":       {},
		"leaf_b":       {},
	}

	p := NewDRHierarchyPartitioner()
	h := p.PartitionWithStableSince(elements, deps, stable)

	if len(h.L1Interfaces) != 2 {
		t.Fatalf("expected 2 L1 elements, got %d (L1=%+v)", len(h.L1Interfaces), h.L1Interfaces)
	}

	old := findByName(h.L1Interfaces, "shared_root")
	if old == nil {
		t.Fatalf("shared_root missing from L1Interfaces; got L1=%+v", h.L1Interfaces)
	}
	if !old.IsTrap {
		t.Errorf("shared_root should be flagged as modularity trap; StableSinceYear=%v TrapReason=%q",
			old.StableSinceYear, old.TrapReason)
	}
	if old.StableSinceYear == nil || *old.StableSinceYear != oldYear {
		t.Errorf("shared_root StableSinceYear=%v want %d", old.StableSinceYear, oldYear)
	}
	if old.TrapReason == "" {
		t.Error("shared_root should have a non-empty TrapReason")
	}

	newer := findByName(h.L1Interfaces, "shared_root2")
	if newer == nil {
		t.Fatal("shared_root2 missing from L1Interfaces")
	}
	if newer.IsTrap {
		t.Errorf("shared_root2 should NOT be flagged as modularity trap (age=%d years, threshold=%d)",
			currentYear()-recentYear, compiletime.ModularityTrapYears)
	}
	if newer.TrapReason != "" {
		t.Errorf("shared_root2 should have empty TrapReason, got %q", newer.TrapReason)
	}
}

// TestPartition_NilStableSinceYearsLeavesNoTraps verifies the no-stamp case
// leaves every L1 element IsTrap=false even when the partitioner is called
// via the Partition convenience overload.
func TestPartition_NilStableSinceYearsLeavesNoTraps(t *testing.T) {
	p := NewDRHierarchyPartitioner()
	elements := []string{"shared_root", "caller1", "caller2", "caller3", "leaf_a"}
	deps := map[string][]string{
		"shared_root": {"leaf_a"},
		"caller1":     {"shared_root"},
		"caller2":     {"shared_root"},
		"caller3":     {"shared_root"},
		"leaf_a":      {},
	}
	h := p.Partition(elements, deps)
	if len(h.L1Interfaces) != 1 {
		t.Fatalf("expected 1 L1 element, got %d (L1=%+v)", len(h.L1Interfaces), h.L1Interfaces)
	}
	if h.L1Interfaces[0].IsTrap {
		t.Errorf("Partition() should not flag traps without stableSinceYears input")
	}
	if h.L1Interfaces[0].TrapReason != "" {
		t.Errorf("TrapReason should be empty, got %q", h.L1Interfaces[0].TrapReason)
	}
}

func findByName(elems []PartitionedElement, name string) *PartitionedElement {
	for i := range elems {
		if elems[i].Name == name {
			return &elems[i]
		}
	}
	return nil
}
