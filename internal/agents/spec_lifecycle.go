package agents

import (
	"fmt"
	"time"
)

// Backlink: [[1.0.0 PRIM-16]] Spec-Driven Development Lifecycle (spec-kit-2024).
// Strict sequence: Constitution -> Specify -> Plan -> Tasks -> Implement -> Converge.
// Each phase produces a reviewable artifact and transitions are gated.

// SpecLifecyclePhase identifies the active stage in the spec-driven lifecycle.
type SpecLifecyclePhase string

const (
	PhaseConstitution SpecLifecyclePhase = "CONSTITUTION"
	PhaseSpecify      SpecLifecyclePhase = "SPECIFY"
	PhasePlan         SpecLifecyclePhase = "PLAN"
	PhaseTasks        SpecLifecyclePhase = "TASKS"
	PhaseImplement    SpecLifecyclePhase = "IMPLEMENT"
	PhaseConverge     SpecLifecyclePhase = "CONVERGE"
)

// LifecycleTransition records a gated step between phases.
type LifecycleTransition struct {
	FromPhase SpecLifecyclePhase `json:"from_phase"`
	ToPhase   SpecLifecyclePhase `json:"to_phase"`
	Timestamp time.Time          `json:"timestamp"`
	Artifact  string             `json:"artifact"`
	Approved  bool               `json:"approved"`
}

// SpecLifecycleManager enforces the PRIM-16 sequence and gates transitions.
type SpecLifecycleManager struct {
	CurrentPhase SpecLifecyclePhase    `json:"current_phase"`
	History      []LifecycleTransition `json:"history"`
}

// NewSpecLifecycleManager initializes the manager at the Constitution phase.
func NewSpecLifecycleManager() *SpecLifecycleManager {
	return &SpecLifecycleManager{
		CurrentPhase: PhaseConstitution,
		History:      make([]LifecycleTransition, 0),
	}
}

// NextPhase returns the canonical succeeding phase.
func (s *SpecLifecycleManager) NextPhase() (SpecLifecyclePhase, error) {
	switch s.CurrentPhase {
	case PhaseConstitution:
		return PhaseSpecify, nil
	case PhaseSpecify:
		return PhasePlan, nil
	case PhasePlan:
		return PhaseTasks, nil
	case PhaseTasks:
		return PhaseImplement, nil
	case PhaseImplement:
		return PhaseConverge, nil
	case PhaseConverge:
		return "", fmt.Errorf("lifecycle already in terminal phase CONVERGE")
	default:
		return "", fmt.Errorf("unknown lifecycle phase: %s", s.CurrentPhase)
	}
}

// Advance gates and executes the transition to the next phase.
func (s *SpecLifecycleManager) Advance(artifactName string, artifactContent string) (SpecLifecyclePhase, error) {
	if artifactContent == "" {
		return s.CurrentPhase, fmt.Errorf("transition from %s rejected: artifact content is empty", s.CurrentPhase)
	}

	next, err := s.NextPhase()
	if err != nil {
		return s.CurrentPhase, err
	}

	transition := LifecycleTransition{
		FromPhase: s.CurrentPhase,
		ToPhase:   next,
		Timestamp: time.Now(),
		Artifact:  artifactName,
		Approved:  true,
	}

	s.History = append(s.History, transition)
	s.CurrentPhase = next
	return next, nil
}
