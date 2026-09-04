package agents

import (
	"fmt"
	"sort"
)

// Backlink: [[1.0.0 PRIM-19]] Design Rule Hierarchy Partitioning (Kazman et al. DRSpaces-2017).
// Classifies architectural elements by stability layer:
// L1: Design Rules / Interfaces (highest stability, abstract contracts)
// L2: Subsystems / Intermediaries (moderate stability)
// L3: Leaves / Concrete implementations (most volatile)
// Applies load-bearing concentration checks to detect modularity violations.

// ArchitectureStabilityLayer defines the architectural strata.
type ArchitectureStabilityLayer string

const (
	LayerL1Interface ArchitectureStabilityLayer = "L1_INTERFACE"
	LayerL2Subsystem ArchitectureStabilityLayer = "L2_SUBSYSTEM"
	LayerL3Leaf      ArchitectureStabilityLayer = "L3_LEAF"
)

// PartitionedElement represents one module or file mapped to a stability layer.
type PartitionedElement struct {
	Name        string                     `json:"name"`
	Layer       ArchitectureStabilityLayer `json:"layer"`
	InDegree    int                        `json:"in_degree"`
	OutDegree   int                        `json:"out_degree"`
	Description string                     `json:"description"`
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
			layer = LayerL1Interface
			desc = "High-stability design rule contract"
			hierarchy.L1Interfaces = append(hierarchy.L1Interfaces, PartitionedElement{
				Name: el, Layer: layer, InDegree: in, OutDegree: out, Description: desc,
			})
		} else if out >= 2 && in <= 1 {
			layer = LayerL3Leaf
			desc = "Volatile concrete leaf implementation"
			hierarchy.L3Leaves = append(hierarchy.L3Leaves, PartitionedElement{
				Name: el, Layer: layer, InDegree: in, OutDegree: out, Description: desc,
			})
		} else {
			layer = LayerL2Subsystem
			desc = "Intermediate subsystem coordinator"
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
	sort.Slice(hierarchy.L1Interfaces, func(i, j int) bool { return hierarchy.L1Interfaces[i].Name < hierarchy.L1Interfaces[j].Name })
	sort.Slice(hierarchy.L2Subsystems, func(i, j int) bool { return hierarchy.L2Subsystems[i].Name < hierarchy.L2Subsystems[j].Name })
	sort.Slice(hierarchy.L3Leaves, func(i, j int) bool { return hierarchy.L3Leaves[i].Name < hierarchy.L3Leaves[j].Name })

	return hierarchy
}
