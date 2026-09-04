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
// Runtime this sprint: concrete sync.RWMutex, claim/complete/fail semantics,
// subscription table, audit log reaper. The goroutine-driver model is left
// to the caller: callers subscribe via Subscribe and drain the channel they
// receive. Subscribers are expected to be lightweight; heavy lifting stays
// on a dedicated worker pool managed by the caller.

package agents

import (
	"sync"
	"time"

	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/consts"
)

// BlackboardEvent is one entry in the blackboard's append-only event log.
// Kind is a stable event-type identifier. Payload is the opaque, agent-defined
// payload, published as map[string]any so subscribers can access fields by
// name without type assertions. Timestamp records when the event was
// published, not when the underlying work was performed.
type BlackboardEvent struct {
	Kind      string
	Payload   any
	Timestamp time.Time
}

// WorkUnit is one schedulable unit of work claimed and completed by an
// agent. ID is the stable identifier (also used as the map key in
// Blackboard.WorkUnits). ClaimedBy is the agent name currently holding the
// lease, or empty when the unit is unclaimed. Status is one of the lifecycle
// constants in internal/consts/consts.go (WorkUnitPending / WorkUnitClaimed /
// WorkUnitCompleted / WorkUnitFailed). Result is the agent's opaque output
// once Status reaches a terminal state.
type WorkUnit struct {
	ID        string
	ClaimedBy string
	Status    compiletime.WorkUnitStatus
	Result    any
}

// subscriber couples a subscription kind with the buffered channel that
// receives matching events. The scheduler closes the channel to signal
// teardown.
type subscriber struct {
	kind string
	ch   chan BlackboardEvent
}

// Blackboard is the shared memory substrate that agents read from and
// publish to. Events is the append-only event log; WorkUnits is the set of
// in-flight and completed work units keyed by ID. mu guards all mutable
// state. subscribers is the fan-out table for live event delivery.
type Blackboard struct {
	Events     []BlackboardEvent
	WorkUnits  map[string]*WorkUnit
	AuditLog   []BlackboardEvent

	mu          sync.RWMutex
	subscribers []*subscriber
}

// NewBlackboard returns a fresh Blackboard with Events and AuditLog
// initialised to non-nil empty slices and WorkUnits initialised to a non-nil
// empty map, so callers can append events and register work units without
// nil checks.
func NewBlackboard() *Blackboard {
	return &Blackboard{
		Events:    []BlackboardEvent{},
		WorkUnits: map[string]*WorkUnit{},
		AuditLog:  []BlackboardEvent{},
	}
}

// Register adds a pending work unit. No-op when the ID already exists.
func (b *Blackboard) Register(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.WorkUnits[id]; ok {
		return
	}
	b.WorkUnits[id] = &WorkUnit{ID: id, Status: consts.WorkUnitPending}
	b.publishLocked(BlackboardEvent{Kind: "work.registered", Payload: map[string]any{"id": id}, Timestamp: time.Now()})
}

// Claim atomically transitions a PENDING unit to CLAIMED and records the
// holding agent. Returns the unit or false if the unit is missing or not
// pending (already claimed, completed, or failed).
func (b *Blackboard) Claim(id, agent string) (*WorkUnit, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	wu, ok := b.WorkUnits[id]
	if !ok || wu.Status != consts.WorkUnitPending {
		return nil, false
	}
	wu.ClaimedBy = agent
	wu.Status = consts.WorkUnitClaimed
	b.publishLocked(BlackboardEvent{Kind: "work.claimed", Payload: map[string]any{"id": id, "agent": agent}, Timestamp: time.Now()})
	return wu, true
}

// Complete transitions a CLAIMED unit to COMPLETED, records the agent's
// result, and reaps the terminal event into the audit log so future Claim
// calls return false.
func (b *Blackboard) Complete(id string, result any) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	wu, ok := b.WorkUnits[id]
	if !ok || wu.Status != consts.WorkUnitClaimed {
		return false
	}
	wu.Result = result
	wu.Status = consts.WorkUnitCompleted
	b.publishLocked(BlackboardEvent{Kind: "work.completed", Payload: map[string]any{"id": id, "result": result}, Timestamp: time.Now()})
	b.reapLocked(id)
	return true
}

// Fail transitions a CLAIMED unit to FAILED and reaps the terminal event.
// Subscribers can react to "work.failed" before the unit is reaped.
func (b *Blackboard) Fail(id string, reason any) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	wu, ok := b.WorkUnits[id]
	if !ok || wu.Status != consts.WorkUnitClaimed {
		return false
	}
	wu.Result = reason
	wu.Status = consts.WorkUnitFailed
	b.publishLocked(BlackboardEvent{Kind: "work.failed", Payload: map[string]any{"id": id, "reason": reason}, Timestamp: time.Now()})
	b.reapLocked(id)
	return true
}

// Snapshot returns a shallow copy of the current event log and work-unit map
// for read-only inspection (e.g. tests, status endpoints). Safe to call
// concurrently with publishers.
func (b *Blackboard) Snapshot() ([]BlackboardEvent, map[string]*WorkUnit) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	events := make([]BlackboardEvent, len(b.Events))
	copy(events, b.Events)
	units := make(map[string]*WorkUnit, len(b.WorkUnits))
	for k, v := range b.WorkUnits {
		copied := *v
		units[k] = &copied
	}
	return events, units
}

// Subscribe returns a buffered channel that receives events whose Kind
// exactly matches kind. The subscription lives until the returned cancel
// function is called. Buffer size matches the typical burst (16); slow
// consumers will skip live delivery (the audit log still captures the
// event so nothing is lost).
func (b *Blackboard) Subscribe(kind string) (<-chan BlackboardEvent, func()) {
	ch := make(chan BlackboardEvent, 16)
	sub := &subscriber{kind: kind, ch: ch}
	b.mu.Lock()
	b.subscribers = append(b.subscribers, sub)
	b.mu.Unlock()
	return ch, func() { b.unsubscribe(sub) }
}

func (b *Blackboard) unsubscribe(sub *subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, s := range b.subscribers {
		if s == sub {
			b.subscribers = append(b.subscribers[:i], b.subscribers[i+1:]...)
			break
		}
	}
	close(sub.ch)
}

// Publish appends an event under the given kind and fans it out to every
// matching subscriber. Subscribers whose buffer is full are skipped (the
// event is still in Events so the audit log captures it).
func (b *Blackboard) Publish(kind string, payload any) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.publishLocked(BlackboardEvent{Kind: kind, Payload: payload, Timestamp: time.Now()})
}

// publishLocked is the inner publish path; the caller must hold b.mu.
func (b *Blackboard) publishLocked(ev BlackboardEvent) {
	b.Events = append(b.Events, ev)
	for _, sub := range b.subscribers {
		if sub.kind != ev.Kind {
			continue
		}
		select {
		case sub.ch <- ev:
		default:
			// Drop live delivery; the event is still in Events/AuditLog.
		}
	}
}

// reapLocked moves the most recent terminal event for a work unit into the
// AuditLog. Called by Complete and Fail with b.mu held. Payload shapes vary
// by event kind; we accept any map[string]any with an "id" key.
func (b *Blackboard) reapLocked(id string) {
	for i := len(b.Events) - 1; i >= 0; i-- {
		ev := b.Events[i]
		if m, ok := ev.Payload.(map[string]any); ok {
			if s, _ := m["id"].(string); s == id {
				b.AuditLog = append(b.AuditLog, ev)
				return
			}
		}
	}
}
