// Asynchronous Software Engineering Agent Blackboard ([[P17]], cited in CAID [P54]).
//
// Multiple agents operate concurrently over a shared memory blackboard: each
// agent reads the current state of the board, claims a work unit, performs
// specialised work, and posts results back as events. Agents wake on the
// arrival of new events matching their subscription and re-evaluate the board
// to decide their next action. A central scheduler arbitrates conflicts by
// serialising claims against the same work unit, broadcasting fresh events to
// interested subscribers, and reaping terminal events into the audit log.
//
// Paper appendix primitive P17 (CAID [P54]) — Asynchronous SE Agent Blackboard.
// Spec-only this sprint: types and constructor only. No scheduler, claim, or
// subscription logic yet; runtime semantics land in a follow-up sprint once the
// central scheduler primitive is specified.
//
// Field notes (to be honoured by the future scheduler implementation):
//   - Mu is declared as `any` so the file is compile-safe across Go versions
//     and so the lock policy can be swapped (sync.RWMutex vs. channel-based
//     serialiser) without churning callers. Concrete locking will land in
//     the scheduler iteration that adopts this type.
//   - Events is an append-only log; consumers must not mutate entries in
//     place. The constructor seeds it as a non-nil empty slice so callers
//     can append directly.
//   - WorkUnits is keyed by stable string ID; ClaimedBy is the agent name
//     holding the lease, Status is one of the lifecycle constants defined in
//     internal/consts/consts.go (added when the scheduler lands).
//   - Result carries the agent's opaque output; the scheduler never
//     inspects it, only re-publishes it as an event payload.

package agents

import "time"

// BlackboardEvent is one entry in the blackboard's append-only event log.
// Kind is a stable event-type identifier (added to internal/consts/consts.go
// when the scheduler is wired up). Payload is the opaque, agent-defined
// payload. Timestamp records when the event was published, not when the
// underlying work was performed.
type BlackboardEvent struct {
	Kind      string
	Payload   any
	Timestamp time.Time
}

// WorkUnit is one schedulable unit of work claimed and completed by an
// agent. ID is the stable identifier (also used as the map key in
// Blackboard.WorkUnits). ClaimedBy is the agent name currently holding the
// lease, or empty when the unit is unclaimed. Status is the lifecycle state
// (sentinel constants land in internal/consts/consts.go). Result is the
// agent's opaque output once Status reaches a terminal state.
type WorkUnit struct {
	ID        string
	ClaimedBy string
	Status    string
	Result    any
}

// Blackboard is the shared memory substrate that agents read from and
// publish to. Events is the append-only event log; WorkUnits is the set of
// in-flight and completed work units keyed by ID. Mu is reserved for the
// concurrency primitive the central scheduler will use to serialise claims
// and event publishes — kept as `any` here to keep this sprint compile-safe
// and to defer the concrete lock-policy decision.
type Blackboard struct {
	Events    []BlackboardEvent
	WorkUnits map[string]*WorkUnit
	// Mu will be a sync.RWMutex (or equivalent) in the scheduler iteration
	// that adopts this type. Declared as `any` for now so this spec-only
	// file carries no runtime locking and no sync import.
	Mu any
}

// NewBlackboard returns a fresh Blackboard with Events initialised to a
// non-nil empty slice and WorkUnits initialised to a non-nil empty map, so
// callers can append events and register work units without nil checks.
// Mu is left as its zero value; the scheduler iteration will populate it.
func NewBlackboard() *Blackboard {
	return &Blackboard{
		Events:    []BlackboardEvent{},
		WorkUnits: map[string]*WorkUnit{},
	}
}
