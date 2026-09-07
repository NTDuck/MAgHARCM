package memorystore_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/memorystore"
)

// TestAddQueryOrderingByReward covers the FIFO-by-reward contract on
// Query: triples with a higher reward rank first; among equals, the
// older insertion (FIFO) wins. The test seeds four triples with
// strictly increasing rewards and asserts Query returns them in
// descending reward order.
func TestAddQueryOrderingByReward(t *testing.T) {
	store := memorystore.NewInMemoryStore()

	now := time.Now().UTC()
	triples := []memorystore.MemoryTriple{
		{Title: "low", Content: "low-reward strategy", Reward: 0.1, CreatedAt: now.Add(3 * time.Second)},
		{Title: "high", Content: "high-reward strategy", Reward: 0.9, CreatedAt: now.Add(1 * time.Second)},
		{Title: "mid", Content: "mid-reward strategy", Reward: 0.5, CreatedAt: now.Add(2 * time.Second)},
		{Title: "mid-old", Content: "tie-break older", Reward: 0.5, CreatedAt: now.Add(0 * time.Second)},
	}
	for _, tr := range triples {
		if err := store.Add(tr); err != nil {
			t.Fatalf("Add(%q): %v", tr.Title, err)
		}
	}

	got := store.Query("ignored", 10)
	if len(got) != 4 {
		t.Fatalf("Query returned %d triples, want 4", len(got))
	}
	wantOrder := []string{"high", "mid-old", "mid", "low"}
	for i, want := range wantOrder {
		if got[i].Title != want {
			t.Errorf("Query rank %d: got %q, want %q", i, got[i].Title, want)
		}
	}
}

// TestCapacityOverflowEviction covers the LRU-by-reward eviction
// contract: when the store is at compiletime.MaxMemoryTriples capacity,
// adding a new triple evicts the lowest-reward triple (oldest first
// among ties). We pre-fill the store with
// compiletime.MaxMemoryTriples triples and add one more, asserting the
// evicted triple is the seeded "low" entry and the new triple is
// present.
func TestCapacityOverflowEviction(t *testing.T) {
	store := memorystore.NewInMemoryStore()

	now := time.Now().UTC()
	for i := 0; i < compiletime.MaxMemoryTriples; i++ {
		_ = store.Add(memorystore.MemoryTriple{
			Title:     fmt.Sprintf("t-%d", i),
			Content:   fmt.Sprintf("body %d", i),
			Reward:    0.5,
			CreatedAt: now.Add(time.Duration(i) * time.Millisecond),
		})
	}

	if size := len(store.Query("", compiletime.MaxMemoryTriples)); size != compiletime.MaxMemoryTriples {
		t.Fatalf("expected fill size %d, got %d", compiletime.MaxMemoryTriples, size)
	}

	// The very first seeded triple has the oldest CreatedAt at reward
	// 0.5; a new triple at reward 0.7 should evict it.
	newer := memorystore.MemoryTriple{
		Title:     "t-newer",
		Content:   "new strategy",
		Reward:    0.7,
		CreatedAt: now.Add(time.Hour),
	}
	if err := store.Add(newer); err != nil {
		t.Fatalf("Add(newer): %v", err)
	}

	got := store.Query("", compiletime.MaxMemoryTriples)
	if len(got) != compiletime.MaxMemoryTriples {
		t.Fatalf("post-eviction size: got %d, want %d", len(got), compiletime.MaxMemoryTriples)
	}

	for _, tr := range got {
		if tr.Title == "t-0" {
			t.Fatalf("expected t-0 to be evicted, but it survived")
		}
	}

	// The new triple should now lead the ranking by virtue of its
	// higher reward.
	if got[0].Title != "t-newer" {
		t.Errorf("expected t-newer to rank first after Add, got %q", got[0].Title)
	}
}

// TestPersistLoadRoundTrip covers the JSON serialisation contract:
// Persist writes a document that Load can read back to reconstruct the
// same set of triples. We use t.TempDir to keep the test hermetic.
func TestPersistLoadRoundTrip(t *testing.T) {
	store := memorystore.NewInMemoryStore()

	original := []memorystore.MemoryTriple{
		{Title: "alpha", Description: "first", Content: "do alpha", Reward: 0.8, CreatedAt: time.Date(2026, 9, 7, 12, 0, 0, 0, time.UTC)},
		{Title: "beta", Description: "second", Content: "do beta", Reward: 0.2, CreatedAt: time.Date(2026, 9, 7, 12, 0, 1, 0, time.UTC)},
	}
	for _, tr := range original {
		if err := store.Add(tr); err != nil {
			t.Fatalf("Add(%q): %v", tr.Title, err)
		}
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "memory.json")
	if err := store.Persist(path); err != nil {
		t.Fatalf("Persist: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("Persist did not create file: %v", err)
	}

	reloaded := memorystore.NewInMemoryStore()
	if err := reloaded.Load(path); err != nil {
		t.Fatalf("Load: %v", err)
	}

	got := reloaded.Query("", 10)
	if len(got) != len(original) {
		t.Fatalf("reloaded size: got %d, want %d", len(got), len(original))
	}
	// Reward-desc ordering means alpha (0.8) leads, then beta (0.2).
	if got[0].Title != "alpha" || got[1].Title != "beta" {
		t.Errorf("reloaded order: got [%s, %s], want [alpha, beta]", got[0].Title, got[1].Title)
	}
	if got[0].Reward != 0.8 {
		t.Errorf("alpha reward: got %f, want 0.8", got[0].Reward)
	}
}

// TestDistillFromTrajectoryHappyPath covers the deterministic
// distillation heuristic: a multi-sentence trajectory yields a triple
// whose Title/Description/Content are populated, and whose reward
// reflects the success flag (1.0 success, 0.0 failure).
func TestDistillFromTrajectoryHappyPath(t *testing.T) {
	trajectory := "Avoid global mutable state. Use dependency injection instead. Ensure the struct is initialised at startup. Add tests for the wiring path."

	for _, success := range []bool{true, false} {
		want := 1.0
		if !success {
			want = 0.0
		}

		got, err := memorystore.DistillFromTrajectory(trajectory, success)
		if err != nil {
			t.Fatalf("DistillFromTrajectory(success=%v): %v", success, err)
		}
		if got.Title == "" {
			t.Errorf("DistillFromTrajectory Title empty (success=%v)", success)
		}
		if got.Content == "" {
			t.Errorf("DistillFromTrajectory Content empty (success=%v)", success)
		}
		if got.Reward != want {
			t.Errorf("DistillFromTrajectory Reward: got %f, want %f", got.Reward, want)
		}
		if got.CreatedAt.IsZero() {
			t.Errorf("DistillFromTrajectory CreatedAt not set (success=%v)", success)
		}
		if !strings.Contains(strings.ToLower(got.Content), "use") &&
			!strings.Contains(strings.ToLower(got.Content), "ensure") &&
			!strings.Contains(strings.ToLower(got.Content), "add") {
			t.Errorf("Content lacks verb cue: %q", got.Content)
		}
	}
}

// TestDistillFromTrajectoryEmpty covers the error path: a
// whitespace-only trajectory cannot yield a triple and must surface
// ErrEmptyTrajectory (not a panic, not a zero-valued triple).
func TestDistillFromTrajectoryEmpty(t *testing.T) {
	for _, in := range []string{"", "   ", "\n\t\r"} {
		if _, err := memorystore.DistillFromTrajectory(in, true); !errors.Is(err, memorystore.ErrEmptyTrajectory) {
			t.Errorf("DistillFromTrajectory(%q): got %v, want ErrEmptyTrajectory", in, err)
		}
	}
}

// TestApplyMaTTSBudgetCeiling covers the budget enforcement contract:
// when the attempt always fails, ApplyMaTTS loops exactly `budget`
// times, then returns ErrBudgetExhausted wrapping the last attempt
// error. The usedStrategies slice grows by one entry per iteration.
func TestApplyMaTTSBudgetCeiling(t *testing.T) {
	store := memorystore.NewInMemoryStore()
	_ = store.Add(memorystore.MemoryTriple{
		Title:   "seed",
		Content: "initial strategy",
		Reward:  0.5,
	})

	const budget = 3
	var attempts int
	wantErr := errors.New("simulated failure")
	err, used := memorystore.ApplyMaTTS(context.Background(), store, func() error {
		attempts++
		return wantErr
	}, budget)

	if attempts != budget {
		t.Errorf("attempts: got %d, want %d", attempts, budget)
	}
	if !errors.Is(err, memorystore.ErrBudgetExhausted) {
		t.Errorf("ApplyMaTTS error: got %v, want ErrBudgetExhausted wrap", err)
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("ApplyMaTTS error: got %v, want it to wrap %v", err, wantErr)
	}
	if len(used) != budget {
		t.Errorf("usedStrategies len: got %d, want %d", len(used), budget)
	}

	// Every iteration should have distilled a fresh triple into the
	// store; the seeded triple plus the budget-many distillations.
	if got := len(store.Query("", compiletime.MaxMemoryTriples)); got != budget+1 {
		t.Errorf("store size after MaTTS: got %d, want %d", got, budget+1)
	}
}

// TestApplyMaTTSErrorPath covers the success-then-error path:
// ApplyMaTTS returns nil on the first successful attempt and short-
// circuits the loop. The usedStrategies slice contains only the prior
// store contents pulled for that single iteration.
func TestApplyMaTTSErrorPath(t *testing.T) {
	store := memorystore.NewInMemoryStore()

	var attempts int
	err, used := memorystore.ApplyMaTTS(context.Background(), store, func() error {
		attempts++
		return nil
	}, 5)

	if err != nil {
		t.Fatalf("ApplyMaTTS: got error %v, want nil", err)
	}
	if attempts != 1 {
		t.Errorf("attempts: got %d, want 1", attempts)
	}
	if len(used) != 1 {
		t.Errorf("usedStrategies len: got %d, want 1", len(used))
	}
}

// TestApplyMaTTSContextCancellation covers the cancellation contract:
// a cancelled context aborts the loop before any attempt runs and the
// returned error is ctx.Err().
func TestApplyMaTTSContextCancellation(t *testing.T) {
	store := memorystore.NewInMemoryStore()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err, used := memorystore.ApplyMaTTS(ctx, store, func() error {
		return nil
	}, 5)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("ApplyMaTTS on cancelled ctx: got %v, want context.Canceled", err)
	}
	if len(used) != 0 {
		t.Errorf("usedStrategies on cancelled ctx: got %d, want 0", len(used))
	}
}

// TestApplyMaTTSNilStore covers the precondition guard: passing a nil
// store returns an error without panicking.
func TestApplyMaTTSNilStore(t *testing.T) {
	err, used := memorystore.ApplyMaTTS(context.Background(), nil, func() error { return nil }, 1)
	if err == nil {
		t.Errorf("ApplyMaTTS(nil store): got nil error, want non-nil")
	}
	if used != nil {
		t.Errorf("ApplyMaTTS(nil store): got used strategies, want nil")
	}
}

// TestConcurrentAddSafe covers the concurrency contract: many
// goroutines adding triples simultaneously MUST NOT race or panic.
// Run with -race to assert no data races.
func TestConcurrentAddSafe(t *testing.T) {
	store := memorystore.NewInMemoryStore()
	var wg sync.WaitGroup
	const goroutines = 16
	const perGoroutine = 8
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				_ = store.Add(memorystore.MemoryTriple{
					Title:     fmt.Sprintf("g%d-i%d", g, i),
					Content:   "race body",
					Reward:    float64(i) / float64(perGoroutine),
					CreatedAt: time.Now().UTC(),
				})
			}
		}(g)
	}
	wg.Wait()

	// Capacity ceiling is MaxMemoryTriples; the store MUST NOT exceed
	// it, no matter how many goroutines raced.
	if got := len(store.Query("", compiletime.MaxMemoryTriples)); got > compiletime.MaxMemoryTriples {
		t.Errorf("store size after race: got %d, want <= %d", got, compiletime.MaxMemoryTriples)
	}
}
