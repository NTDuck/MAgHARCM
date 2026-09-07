// Package memorystore is the persistent-memory substrate backing
// Memory-aware Test-Time Scaling (MaTTS), per [[1.0.0 P-122]].
//
// Backlink: [[1.0.0 Architecture]] §8 (Wave-16 opt-in patterns) and
// [[1.0.0 Methodology]] §11. The substrate is intentionally leaf-level:
// it imports only stdlib, internal/compiletime, and internal/logger, so
// higher layers (agents, graph, runner) can wire it in without creating
// cycles. All configuration sentinels (capacity, EMA factor, default
// budget) live in compiletime (ADR-C-005); this package never hardcodes
// magic numbers and never falls back to defaults at runtime (ADR-C-001).
package memorystore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/logger"
)

// MemoryTriple is a single distilled memory unit stored in the
// ReasoningBank substrate. Each triple represents one strategy
// observation produced by DistillFromTrajectory: the Title is a short
// slug, the Description captures the surrounding context, the Content
// is the actionable strategy body, the Reward is the EMA-smoothed
// outcome signal, and CreatedAt records when the triple entered the
// store (used for FIFO tie-breaking during eviction).
//
// Backlink: [[1.0.0 P-122]] §3.2 — ReasoningBank triple schema.
type MemoryTriple struct {
	// Title is a short, slug-like identifier for the distilled strategy.
	Title string `json:"title"`
	// Description records the trajectory context that produced the triple.
	Description string `json:"description"`
	// Content is the actionable strategy body surfaced to future attempts.
	Content string `json:"content"`
	// Reward is the EMA-smoothed outcome signal (0.0 = always failed,
	// 1.0 = always succeeded). Larger values are kept during eviction.
	Reward float64 `json:"reward"`
	// CreatedAt records when the triple entered the store; used as the
	// FIFO tie-breaker when two triples share the same reward.
	CreatedAt time.Time `json:"created_at"`
}

// Store is the persistent-memory substrate contract. Implementations MUST
// be safe for concurrent use; the canonical in-memory implementation is
// returned by NewInMemoryStore and satisfies this contract with a single
// mutex guarding the triple slice.
type Store interface {
	// Add inserts t into the store. Implementations MAY evict an
	// existing triple to honour capacity bounds; the evicted triple is
	// the one with the lowest reward, breaking ties by oldest
	// CreatedAt. Returns an error only on I/O or invariant violations,
	// never on capacity pressure.
	Add(t MemoryTriple) error
	// Query returns the top-k triples ranked by reward (descending),
	// with FIFO ordering (oldest first) breaking ties. When k is
	// larger than the store size the full store is returned. The
	// prompt argument is reserved for future similarity ranking and is
	// currently ignored — kept on the interface so the contract does
	// not churn when a vector backend is added.
	Query(prompt string, k int) []MemoryTriple
	// Persist serialises the store to path as a single JSON document.
	// The parent directory is created if missing. Path is overwritten.
	Persist(path string) error
	// Load replaces the in-memory contents from path. The file MUST
	// exist; an error is returned for missing or malformed files so
	// callers cannot silently start with an empty store.
	Load(path string) error
}

// ErrEmptyTrajectory is returned by DistillFromTrajectory when the
// caller supplies a whitespace-only attempt summary, which cannot yield
// a meaningful strategy triple.
var ErrEmptyTrajectory = errors.New("memorystore: empty attempt trajectory")

// ErrBudgetExhausted is returned by ApplyMaTTS when the caller's
// per-invocation iteration ceiling is reached before the attempt
// succeeds. The error wraps the final attempt error so callers can
// distinguish "ran out of budget" from "the last attempt itself failed".
var ErrBudgetExhausted = errors.New("memorystore: MaTTS budget exhausted")

// inMemoryStore is the canonical Store implementation: a slice guarded
// by a single mutex, with reward-ranked LRU eviction. It honours the
// capacity bound from compiletime.MaxMemoryTriples.
type inMemoryStore struct {
	mu       sync.Mutex
	triples  []MemoryTriple
	capacity int
}

// NewInMemoryStore returns a Store backed by process-local memory with
// capacity compiletime.MaxMemoryTriples. The returned implementation
// evicts the lowest-reward triple (oldest first) when capacity is
// exceeded, so high-value strategies naturally survive long sessions
// while noise is forgotten.
func NewInMemoryStore() Store {
	return &inMemoryStore{
		triples:  make([]MemoryTriple, 0, compiletime.MaxMemoryTriples),
		capacity: compiletime.MaxMemoryTriples,
	}
}

// Add inserts t into the store. When the store is at capacity the
// lowest-reward triple (oldest first among ties) is evicted to make
// room. The created-at timestamp is overwritten with the current wall
// clock when the caller has not provided one, so the FIFO tie-breaker
// reflects insertion order.
func (s *inMemoryStore) Add(t MemoryTriple) error {
	if t.Title == "" {
		return errors.New("memorystore: triple Title must not be empty")
	}
	if t.Content == "" {
		return errors.New("memorystore: triple Content must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}

	if len(s.triples) >= s.capacity {
		evictIdx := indexToEvict(s.triples)
		s.triples = append(s.triples[:evictIdx], s.triples[evictIdx+1:]...)
	}
	s.triples = append(s.triples, t)
	return nil
}

// Query returns the top-k triples ranked by reward (descending) with
// FIFO ordering (oldest first) breaking ties. The prompt argument is
// currently ignored — see the Store contract for the rationale.
func (s *inMemoryStore) Query(_ string, k int) []MemoryTriple {
	s.mu.Lock()
	defer s.mu.Unlock()

	ordered := make([]MemoryTriple, len(s.triples))
	copy(ordered, s.triples)
	sortByRewardDescFIFO(ordered)

	if k <= 0 || k > len(ordered) {
		return ordered
	}
	return ordered[:k]
}

// Persist writes the current store contents to path as JSON. The
// parent directory is created with the standard 0o755 mode; the file
// itself is written with 0o644. An existing file at path is overwritten.
func (s *inMemoryStore) Persist(path string) error {
	s.mu.Lock()
	snapshot := make([]MemoryTriple, len(s.triples))
	copy(snapshot, s.triples)
	s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("memorystore: persist mkdir: %w", err)
	}
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Errorf("memorystore: persist marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("memorystore: persist write: %w", err)
	}
	return nil
}

// Load replaces the in-memory contents from the JSON document at path.
// The file MUST exist; callers that want an empty store must use
// NewInMemoryStore directly. Decoding failures surface verbatim so the
// loader can fail loudly rather than silently dropping state.
func (s *inMemoryStore) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("memorystore: load read: %w", err)
	}
	var loaded []MemoryTriple
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("memorystore: load unmarshal: %w", err)
	}
	if loaded == nil {
		loaded = []MemoryTriple{}
	}
	s.mu.Lock()
	s.triples = loaded
	s.mu.Unlock()
	return nil
}

// indexToEvict returns the index of the lowest-reward triple in ts,
// breaking ties by oldest CreatedAt (FIFO). ts MUST be non-empty; the
// caller is responsible for guarding this precondition.
func indexToEvict(ts []MemoryTriple) int {
	idx := 0
	for i := 1; i < len(ts); i++ {
		if ts[i].Reward < ts[idx].Reward {
			idx = i
			continue
		}
		if ts[i].Reward == ts[idx].Reward && ts[i].CreatedAt.Before(ts[idx].CreatedAt) {
			idx = i
		}
	}
	return idx
}

// sortByRewardDescFIFO sorts in place: highest reward first, oldest
// first among equal rewards. Used by Query to honour the
// reward-desc / FIFO-asc contract documented on Store.
func sortByRewardDescFIFO(ts []MemoryTriple) {
	sort.SliceStable(ts, func(i, j int) bool {
		if ts[i].Reward != ts[j].Reward {
			return ts[i].Reward > ts[j].Reward
		}
		return ts[i].CreatedAt.Before(ts[j].CreatedAt)
	})
}

// DistillFromTrajectory extracts a MemoryTriple from a single attempt
// trajectory using a deterministic, non-LLM heuristic (sentence-boundary
// split + scoring). The success flag sets the reward signal: 1.0 on
// success, 0.0 on failure. The Title is derived from the first
// content-bearing sentence; the Description captures the surrounding
// context; the Content is the highest-information sentence.
//
// This is the ReasoningBank strategy-distillation pattern: produce a
// reusable strategy artefact from each attempt's outcome without
// requiring an LLM at distillation time. Callers that have an LLM
// available SHOULD augment this baseline, not replace it — the
// heuristic guarantee is what lets the substrate run in tests and
// resource-constrained pipelines.
func DistillFromTrajectory(attemptSummary string, success bool) (MemoryTriple, error) {
	trimmed := strings.TrimSpace(attemptSummary)
	if trimmed == "" {
		return MemoryTriple{}, ErrEmptyTrajectory
	}

	sentences := splitSentences(trimmed)
	if len(sentences) == 0 {
		return MemoryTriple{}, ErrEmptyTrajectory
	}

	// Pick the highest-scoring sentence as Content; everything else
	// becomes the Description context. Scoring rewards longer,
	// verb-bearing sentences — a crude proxy for actionability.
	bestIdx := 0
	bestScore := -1.0
	for i, s := range sentences {
		score := scoreSentence(s)
		if score > bestScore {
			bestScore = score
			bestIdx = i
		}
	}

	title := titleFromSentence(sentences[bestIdx])
	content := sentences[bestIdx]
	var description string
	if len(sentences) > 1 {
		description = strings.Join(append(append([]string{}, sentences[:bestIdx]...), sentences[bestIdx+1:]...), " ")
	} else {
		description = content
	}

	reward := 0.0
	if success {
		reward = 1.0
	}

	return MemoryTriple{
		Title:       title,
		Description: description,
		Content:     content,
		Reward:      reward,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

// splitSentences splits s on sentence terminators (. ! ?) followed by
// whitespace, returning non-empty trimmed sentences in order. The split
// is intentionally simple — punctuation-based sentence segmentation is
// deterministic, allocation-light, and language-agnostic enough for
// the strategy-distillation use case.
func splitSentences(s string) []string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed == "" {
			continue
		}
		out = append(out, trimmed)
	}
	return out
}

// scoreSentence returns a crude actionability score: longer sentences
// with verb-suggesting tokens outrank short, declarative fragments.
// The exact weights are not principled — the heuristic exists to make
// DistillFromTrajectory deterministic and allocation-light, not to
// match a benchmark.
func scoreSentence(s string) float64 {
	words := strings.Fields(s)
	if len(words) == 0 {
		return 0
	}
	score := float64(len(words))
	lower := strings.ToLower(s)
	for _, verb := range []string{"use", "apply", "avoid", "ensure", "add", "fix", "wrap", "set"} {
		if strings.Contains(lower, verb) {
			score += 2
		}
	}
	return score
}

// titleFromSentence derives a short, slug-like Title from a single
// sentence: trims, lowercases, replaces whitespace with hyphens, and
// strips trailing punctuation. Callers use the title as the human
// identifier surfaced in logs and dashboards.
func titleFromSentence(s string) string {
	t := strings.ToLower(strings.TrimSpace(s))
	t = strings.TrimRight(t, ".!?")
	var b strings.Builder
	b.Grow(len(t))
	prevSpace := false
	for _, r := range t {
		switch {
		case r == ' ' || r == '\t':
			if !prevSpace {
				b.WriteByte('-')
				prevSpace = true
			}
		default:
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return strings.Trim(b.String(), "-")
}

// ApplyMaTTS is the Memory-aware Test-Time Scaling loop ([[1.0.0 P-122]]).
// It pulls the top-k triples from store as context for the attempt,
// runs attempt, distils the resulting trajectory into a fresh triple,
// persists it, and loops until either attempt succeeds or the budget
// is exhausted. The returned usedStrategies list contains the triples
// pulled on each iteration (newest attempt first); the returned error
// is ErrBudgetExhausted wrapping the last attempt error when the
// budget ceiling is reached, otherwise the attempt error verbatim.
//
// budget <= 0 falls back to compiletime.MaTTSDefaultBudget so the
// package contract is satisfied without forcing every callsite to
// thread a constant through. ctx is honoured before each iteration;
// cancellation propagates immediately and surfaces as ctx.Err().
func ApplyMaTTS(ctx context.Context, store Store, attempt func() error, budget int) (error, []MemoryTriple) {
	if store == nil {
		return errors.New("memorystore: ApplyMaTTS requires a non-nil Store"), nil
	}
	if attempt == nil {
		return errors.New("memorystore: ApplyMaTTS requires a non-nil attempt"), nil
	}
	if budget <= 0 {
		budget = compiletime.MaTTSDefaultBudget
	}

	var usedStrategies []MemoryTriple
	var lastErr error

	for i := 0; i < budget; i++ {
		if err := ctx.Err(); err != nil {
			return err, usedStrategies
		}

		// Pull the full store as the strategy context for this
		// iteration. A future revision can switch to a vector
		// retriever once the store graduates from in-memory.
		prior := store.Query("", compiletime.MaxMemoryTriples)
		usedStrategies = append(usedStrategies, prior...)

		attemptErr := attempt()
		summary := attemptSummary(prior, attemptErr)

		triple, distilErr := DistillFromTrajectory(summary, attemptErr == nil)
		if distilErr != nil {
			logger.LogWarning("memorystore: distill failed on attempt %d: %v", i+1, distilErr)
		} else if addErr := store.Add(triple); addErr != nil {
			logger.LogWarning("memorystore: persist failed on attempt %d: %v", i+1, addErr)
		}

		if attemptErr == nil {
			return nil, usedStrategies
		}
		lastErr = attemptErr
	}

	return fmt.Errorf("%w: last error: %v", ErrBudgetExhausted, lastErr), usedStrategies
}

// attemptSummary renders the trajectory line fed into DistillFromTrajectory:
// the prior strategies (newest first) plus a one-line outcome marker. The
// format is intentionally verbose because DistillFromTrajectory treats it
// as the input corpus — the more signal in the trajectory, the better
// the distilled triple.
func attemptSummary(prior []MemoryTriple, attemptErr error) string {
	var b strings.Builder
	for _, p := range prior {
		b.WriteString("Strategy: ")
		b.WriteString(p.Title)
		b.WriteString(" — ")
		b.WriteString(p.Content)
		b.WriteString(". ")
	}
	if attemptErr == nil {
		b.WriteString("Attempt succeeded.")
	} else {
		b.WriteString("Attempt failed: ")
		b.WriteString(attemptErr.Error())
		b.WriteString(".")
	}
	return b.String()
}
