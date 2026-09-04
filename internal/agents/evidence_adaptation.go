package agents

import (
	"context"
	"fmt"
	"strings"
)

// Backlink: [[1.0.0 PRIM-15]] Evidence-First Adaptation Pattern (Reeper-2024).
// Before adapting legacy code into a target system, produce four artifacts:
// (1) Target-preserving spec
// (2) Contract-before-code
// (3) Provenance-aware diff
// (4) Untrusted-source model

// TargetPreservingSpec defines invariants that target code must preserve.
type TargetPreservingSpec struct {
	SourceUnit       string   `json:"source_unit"`
	RequiredBehaviors []string `json:"required_behaviors"`
	DisallowedSideEffects []string `json:"disallowed_side_effects"`
}

// ContractBeforeCode defines input/output types and error bounds prior to synthesis.
type ContractBeforeCode struct {
	OperationName string            `json:"operation_name"`
	InputTypes    map[string]string `json:"input_types"`
	ReturnType    string            `json:"return_type"`
	ErrorBounds   []string          `json:"error_bounds"`
}

// ProvenanceDiffEntry maps a source construct to its target construct with origin tracking.
type ProvenanceDiffEntry struct {
	SourceSymbol string `json:"source_symbol"`
	TargetSymbol string `json:"target_symbol"`
	OriginFile   string `json:"origin_file"`
	Rationale    string `json:"rationale"`
}

// UntrustedSourceModel encapsulates isolation and safety assumptions for foreign code.
type UntrustedSourceModel struct {
	ContainsUnsafeBlocks bool     `json:"contains_unsafe_blocks"`
	PotentialMemoryLeaks bool     `json:"potential_memory_leaks"`
	SanitizationGates    []string `json:"sanitization_gates"`
}

// EvidenceAdaptationReport combines the four required Reeper artifacts.
type EvidenceAdaptationReport struct {
	Spec        TargetPreservingSpec  `json:"spec"`
	Contract    ContractBeforeCode    `json:"contract"`
	Provenance  []ProvenanceDiffEntry `json:"provenance"`
	SafetyModel UntrustedSourceModel  `json:"safety_model"`
}

// EvidenceFirstAdaptor implements PRIM-15.
type EvidenceFirstAdaptor struct{}

// NewEvidenceFirstAdaptor constructs a new PRIM-15 adaptor.
func NewEvidenceFirstAdaptor() *EvidenceFirstAdaptor {
	return &EvidenceFirstAdaptor{}
}

// Adapt evaluates source code and produces an EvidenceAdaptationReport.
func (e *EvidenceFirstAdaptor) Adapt(ctx context.Context, unitName, sourceContent, targetLang string) (*EvidenceAdaptationReport, error) {
	if unitName == "" {
		return nil, fmt.Errorf("unitName is required")
	}

	report := &EvidenceAdaptationReport{
		Spec: TargetPreservingSpec{
			SourceUnit: unitName,
			RequiredBehaviors: []string{
				fmt.Sprintf("Preserve deterministic I/O of %s", unitName),
				"Ensure nullability safety and proper resource disposal",
			},
			DisallowedSideEffects: []string{
				"Do not panic on invalid user input",
				"Do not mutate shared state without synchronization",
			},
		},
		Contract: ContractBeforeCode{
			OperationName: unitName,
			InputTypes:    map[string]string{"input": "standard_context"},
			ReturnType:    "Result<TargetOutput, TargetError>",
			ErrorBounds:   []string{"InvalidInputError", "InternalExecutionError"},
		},
		Provenance: []ProvenanceDiffEntry{
			{
				SourceSymbol: unitName,
				TargetSymbol: strings.ToLower(unitName),
				OriginFile:   unitName,
				Rationale:    fmt.Sprintf("Direct functional mapping to %s idiom", targetLang),
			},
		},
		SafetyModel: UntrustedSourceModel{
			ContainsUnsafeBlocks: strings.Contains(sourceContent, "unsafe") || strings.Contains(sourceContent, "pointer"),
			PotentialMemoryLeaks: strings.Contains(sourceContent, "malloc") || strings.Contains(sourceContent, "free"),
			SanitizationGates: []string{
				"Wrap raw pointer dereferences in safe bounds-checked abstractions",
				"Enforce static borrow-checking invariants",
			},
		},
	}

	return report, nil
}
