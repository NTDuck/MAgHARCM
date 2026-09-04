package agents

import (
	"fmt"
	"time"

	"MAgHARCM/internal/compiletime"
)

// Backlink: [[1.0.0 PRIM-16]] Spec-Driven Development Lifecycle (spec-kit-2024).
// Strict sequence: Constitution -> Specify -> Plan -> Tasks -> Implement -> Converge.
// SpecLifecyclePhase aliases compiletime.SpecLifecyclePhase for centralized enums.
type SpecLifecyclePhase = compiletime.SpecLifecyclePhase

// Phase* are the back-compat aliases for compiletime.Phase*. New code should
// reference compiletime.Phase* directly.
const (
	PhaseConstitution = compiletime.PhaseConstitution
	PhaseSpecify      = compiletime.PhaseSpecify
	PhasePlan         = compiletime.PhasePlan
	PhaseTasks        = compiletime.PhaseTasks
	PhaseImplement    = compiletime.PhaseImplement
	PhaseConverge     = compiletime.PhaseConverged
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

// MustAdvance is the panic-on-error variant of Advance. Use at startup where
// a malformed lifecycle is a fatal configuration error.
func (s *SpecLifecycleManager) MustAdvance(artifactName, artifactContent string) SpecLifecyclePhase {
	phase, err := s.Advance(artifactName, artifactContent)
	compiletime.Must(phase, err)
	return phase
}
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
