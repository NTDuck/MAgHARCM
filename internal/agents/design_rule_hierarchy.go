package agents
import (
	"fmt"
	"slices"
	"cmp"
	"time"

	"MAgHARCM/internal/compiletime"
)

type ArchitectureStabilityLayer = compiletime.ArchitectureStabilityLayer

// Layer* are the back-compat aliases for compiletime.ArchitectureStabilityLayer*.
// New code should reference the compiletime constants directly.
const (
	LayerL1Interface ArchitectureStabilityLayer = compiletime.ArchitectureStabilityLayerL1
	LayerL2Subsystem ArchitectureStabilityLayer = compiletime.ArchitectureStabilityLayerL2
	LayerL3Leaf      ArchitectureStabilityLayer = compiletime.ArchitectureStabilityLayerL3
)

type PartitionedElement struct {
	Name        string                     `json:"name"`
	Layer       ArchitectureStabilityLayer `json:"layer"`
	InDegree    int                        `json:"in_degree"`
	OutDegree   int                        `json:"out_degree"`
	Description string                     `json:"description"`
	// StableSinceYear is the year the design rule became visible / stable;
	// nil when unknown. Set by callers who can read git blame or file mtime;
	// consumed by IsModularityTrap to flag rules older than
	// [[1.0.0 PRIM-71]]'s threshold.
	StableSinceYear *int    `json:"stable_since_year,omitempty"`
	IsTrap          bool    `json:"is_modularity_trap,omitempty"`
	TrapReason      string  `json:"trap_reason,omitempty"`
}

// ModularityViolation records an architectural smell where an L1 element depends on an L3 leaf.
type ModularityViolation struct {
	SourceL1 string `json:"source_l1"`
	TargetL3 string `json:"target_l3"`
	Details  string `json:"details"`
}

// DRHierarchy holds the partitioned architecture layers and detected violations.
type DRHierarchy struct {
	L1Interfaces []PartitionedElement  `json:"l1_interfaces"`
	L2Subsystems []PartitionedElement  `json:"l2_subsystems"`
	L3Leaves     []PartitionedElement  `json:"l3_leaves"`
	Violations   []ModularityViolation `json:"violations"`
}

// DRHierarchyPartitioner implements PRIM-19.
type DRHierarchyPartitioner struct{}

// NewDRHierarchyPartitioner constructs a new partitioner.
func NewDRHierarchyPartitioner() *DRHierarchyPartitioner {
	return &DRHierarchyPartitioner{}
}
// Partition classifies elements based on their fan-in (in-degree) and fan-out (out-degree).
// Elements with high fan-in and low fan-out are L1 design rules.
// Elements with zero fan-in are L3 leaf implementations.
func (d *DRHierarchyPartitioner) Partition(elements []string, dependencies map[string][]string) *DRHierarchy {
	return d.PartitionWithStableSince(elements, dependencies, nil)
}

// PartitionWithStableSince extends Partition with optional per-element
// stable-since-year input. When provided, each L1 element whose stable-since
// year is older than [[1.0.0 PRIM-71]]'s ModularityTrapYears threshold is
// flagged as a modularity trap; the caller surfaces these to the user for
// explicit re-validation. Years reference [[1.0.0 P-71]] (Fleming &amp;
// Baldwin 2024 retrospective on Design Rules).
func (d *DRHierarchyPartitioner) PartitionWithStableSince(elements []string, dependencies map[string][]string, stableSinceYears map[string]int) *DRHierarchy {
	inDegree := make(map[string]int)
	outDegree := make(map[string]int)

	for _, el := range elements {
		inDegree[el] = 0
		outDegree[el] = len(dependencies[el])
	}

	for _, deps := range dependencies {
		for _, dep := range deps {
			inDegree[dep]++
		}
	}

	hierarchy := &DRHierarchy{
		L1Interfaces: make([]PartitionedElement, 0),
		L2Subsystems: make([]PartitionedElement, 0),
		L3Leaves:     make([]PartitionedElement, 0),
		Violations:   make([]ModularityViolation, 0),
	}

	layerMap := make(map[string]ArchitectureStabilityLayer)

	for _, el := range elements {
		in := inDegree[el]
		out := outDegree[el]

		var layer ArchitectureStabilityLayer
		var desc string

		if in >= 2 && out <= 1 {
			desc = compiletime.ArchitectureStabilityDescriptionL1
			elem := PartitionedElement{Name: el, Layer: layer, InDegree: in, OutDegree: out, Description: desc}
			if trapEl, isTrap := markModularityTrap(elem, el, stableSinceYears); isTrap {
				elem = trapEl
			}
			hierarchy.L1Interfaces = append(hierarchy.L1Interfaces, elem)
		} else if out >= 2 && in <= 1 {
			layer = LayerL3Leaf
			desc = compiletime.ArchitectureStabilityDescriptionL3
			hierarchy.L3Leaves = append(hierarchy.L3Leaves, PartitionedElement{
				Name: el, Layer: layer, InDegree: in, OutDegree: out, Description: desc,
			})
		} else {
			layer = LayerL2Subsystem
			desc = compiletime.ArchitectureStabilityDescriptionL2
			hierarchy.L2Subsystems = append(hierarchy.L2Subsystems, PartitionedElement{
				Name: el, Layer: layer, InDegree: in, OutDegree: out, Description: desc,
			})
		}
		layerMap[el] = layer
	}

	// Load-bearing concentration check: L1 cannot depend on L3
	for _, l1 := range hierarchy.L1Interfaces {
		for _, dep := range dependencies[l1.Name] {
			if layerMap[dep] == LayerL3Leaf {
				hierarchy.Violations = append(hierarchy.Violations, ModularityViolation{
					SourceL1: l1.Name,
					TargetL3: dep,
					Details:  fmt.Sprintf("Design rule interface %q illegally depends on concrete leaf %q", l1.Name, dep),
				})
			}
		}
	}

	// Stable sort for deterministic outputs
	slices.SortFunc(hierarchy.L1Interfaces, func(a, b PartitionedElement) int { return cmp.Compare(a.Name, b.Name) })
	slices.SortFunc(hierarchy.L2Subsystems, func(a, b PartitionedElement) int { return cmp.Compare(a.Name, b.Name) })

	return hierarchy
}

// markModularityTrap flags an L1 design rule as a modularity trap when its
// stable-since year is older than [[1.0.0 PRIM-71]]'s ModularityTrapYears
// threshold. Returns the original element unchanged when stableSinceYears
// is nil or does not contain the element. Caller surfaces flagged elements
// to the user for explicit re-validation per [[1.0.0 P-71]] §4.
func markModularityTrap(elem PartitionedElement, name string, stableSinceYears map[string]int) (PartitionedElement, bool) {
	if stableSinceYears == nil {
		return elem, false
	}
	year, ok := stableSinceYears[name]
	if !ok {
		return elem, false
	}
	age := currentYear() - year
	if age < compiletime.ModularityTrapYears {
		return elem, false
	}
	yearCopy := year
	elem.StableSinceYear = &yearCopy
	elem.IsTrap = true
	elem.TrapReason = fmt.Sprintf("Design rule stable since %d (%d years ago, exceeds %d-year threshold)",
		year, age, compiletime.ModularityTrapYears)
	return elem, true
}

// currentYear is a package-private seam so tests can stub the calendar
// without monkey-patching time.Now everywhere.
var currentYear = func() int { return time.Now().Year() }

