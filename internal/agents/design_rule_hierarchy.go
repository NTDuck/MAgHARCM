package agents

import (
	"fmt"
	"slices"
	"cmp"

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
			desc = compiletime.ArchitectureStabilityDescriptionL1
			hierarchy.L1Interfaces = append(hierarchy.L1Interfaces, PartitionedElement{
				Name: el, Layer: layer, InDegree: in, OutDegree: out, Description: desc,
			})
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
	slices.SortFunc(hierarchy.L3Leaves, func(a, b PartitionedElement) int { return cmp.Compare(a.Name, b.Name) })

	return hierarchy
}
