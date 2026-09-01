package agents

import (
	"strings"
	"testing"
	"time"
)

func TestPlateauDetectorEmpty(t *testing.T) {
	p := NewPlateauDetector()
	if p.IsPlateau() {
		t.Fatal("empty detector should not be plateau")
	}
	if got := p.Summary(); got != "plateau detector: no samples" {
		t.Fatalf("unexpected empty summary: %q", got)
	}
}

func TestPlateauDetectorSingleSample(t *testing.T) {
	p := NewPlateauDetector()
	p.Record(CoverageSample{Iteration: 0, PassedTests: 10, UncoveredCount: 5, CapturedAt: time.Now()})
	if p.IsPlateau() {
		t.Fatal("single sample should not be plateau (need Window+1 to compare)")
	}
	if len(p.Samples) != 1 {
		t.Fatalf("expected 1 sample retained, got %d", len(p.Samples))
	}
}

func TestPlateauDetectorBelowThreshold(t *testing.T) {
	p := NewPlateauDetector()
	// Three samples: 0 -> 1 -> 2; deltas: passed +2, uncovered -0 (no progress)
	p.Record(CoverageSample{Iteration: 0, PassedTests: 10, UncoveredCount: 5, CapturedAt: time.Now()})
	p.Record(CoverageSample{Iteration: 1, PassedTests: 12, UncoveredCount: 5, CapturedAt: time.Now()})
	if p.IsPlateau() {
		t.Fatal("Window=2 requires two stagnant pairs; only one pair so far")
	}
	p.Record(CoverageSample{Iteration: 2, PassedTests: 13, UncoveredCount: 5, CapturedAt: time.Now()})
	if !p.IsPlateau() {
		t.Fatal("two consecutive stagnant pairs should trigger plateau")
	}
}

func TestPlateauDetectorAboveThreshold(t *testing.T) {
	p := NewPlateauDetector()
	p.Record(CoverageSample{Iteration: 0, PassedTests: 10, UncoveredCount: 5, CapturedAt: time.Now()})
	// passed +10 (>= MinPassedGain=5): NOT stagnant.
	p.Record(CoverageSample{Iteration: 1, PassedTests: 20, UncoveredCount: 5, CapturedAt: time.Now()})
	// passed +10 again.
	p.Record(CoverageSample{Iteration: 2, PassedTests: 30, UncoveredCount: 5, CapturedAt: time.Now()})
	if p.IsPlateau() {
		t.Fatal("above-threshold progress should not be plateau")
	}
}

func TestPlateauDetectorMixedThenRecover(t *testing.T) {
	p := NewPlateauDetector()
	p.Record(CoverageSample{Iteration: 0, PassedTests: 10, UncoveredCount: 5, CapturedAt: time.Now()})
	p.Record(CoverageSample{Iteration: 1, PassedTests: 11, UncoveredCount: 5, CapturedAt: time.Now()}) // stagnant
	p.Record(CoverageSample{Iteration: 2, PassedTests: 50, UncoveredCount: 5, CapturedAt: time.Now()}) // big jump
	if p.IsPlateau() {
		t.Fatal("recovery should reset plateau detection")
	}
}

func TestPlateauDetectorWindowOne(t *testing.T) {
	p := &PlateauDetector{
		Samples:          []CoverageSample{},
		Window:           1,
		MinPassedGain:    5,
		MinUncoveredGain: 1,
	}
	p.Record(CoverageSample{Iteration: 0, PassedTests: 10, UncoveredCount: 5, CapturedAt: time.Now()})
	p.Record(CoverageSample{Iteration: 1, PassedTests: 11, UncoveredCount: 5, CapturedAt: time.Now()})
	if !p.IsPlateau() {
		t.Fatal("Window=1 with one stagnant pair should trigger plateau")
	}
}

func TestCoverageDelta(t *testing.T) {
	prev := CoverageSample{PassedTests: 10, UncoveredCount: 5}
	curr := CoverageSample{PassedTests: 20, UncoveredCount: 2}
	pd, ud := CoverageDelta(prev, curr)
	if pd != 10 {
		t.Fatalf("expected passedDelta=10, got %d", pd)
	}
	if ud != 3 {
		t.Fatalf("expected uncoveredDelta=3, got %d", ud)
	}

	// Zero / negative cases.
	zero := CoverageSample{PassedTests: 10, UncoveredCount: 5}
	pd, ud = CoverageDelta(prev, zero)
	if pd != 0 || ud != 0 {
		t.Fatalf("expected (0,0), got (%d,%d)", pd, ud)
	}
	neg := CoverageSample{PassedTests: 8, UncoveredCount: 7}
	pd, ud = CoverageDelta(prev, neg)
	if pd != -2 || ud != -2 {
		t.Fatalf("expected (-2,-2), got (%d,%d)", pd, ud)
	}
}

func TestPlateauDetectorSummary(t *testing.T) {
	p := NewPlateauDetector()
	p.Record(CoverageSample{Iteration: 0, PassedTests: 10, UncoveredCount: 5, CapturedAt: time.Now()})
	p.Record(CoverageSample{Iteration: 1, PassedTests: 11, UncoveredCount: 5, CapturedAt: time.Now()})
	p.Record(CoverageSample{Iteration: 2, PassedTests: 12, UncoveredCount: 5, CapturedAt: time.Now()})

	summary := p.Summary()
	if !strings.Contains(summary, "PLATEAU DETECTED") {
		t.Fatalf("expected summary to contain 'PLATEAU DETECTED', got: %s", summary)
	}
	if !strings.Contains(summary, "passed=12") {
		t.Fatalf("expected summary to include last passed count, got: %s", summary)
	}

	// And on a progressing detector the summary should reflect that.
	p2 := NewPlateauDetector()
	p2.Record(CoverageSample{Iteration: 0, PassedTests: 0, UncoveredCount: 10, CapturedAt: time.Now()})
	p2.Record(CoverageSample{Iteration: 1, PassedTests: 50, UncoveredCount: 0, CapturedAt: time.Now()})
	s := p2.Summary()
	if !strings.Contains(s, "progressing") {
		t.Fatalf("expected 'progressing' in summary, got: %s", s)
	}
	if strings.Contains(s, "PLATEAU DETECTED") {
		t.Fatalf("did not expect PLATEAU DETECTED, got: %s", s)
	}
}

func TestPlateauDetectorRecordReturnsPlateauState(t *testing.T) {
	p := NewPlateauDetector()
	plateau1 := p.Record(CoverageSample{Iteration: 0, PassedTests: 0, UncoveredCount: 5, CapturedAt: time.Now()})
	if plateau1 {
		t.Fatal("first Record() should not trigger plateau")
	}
	plateau2 := p.Record(CoverageSample{Iteration: 1, PassedTests: 1, UncoveredCount: 5, CapturedAt: time.Now()})
	if plateau2 {
		t.Fatal("second sample with one stagnant pair should not plateau")
	}
	plateau3 := p.Record(CoverageSample{Iteration: 2, PassedTests: 2, UncoveredCount: 5, CapturedAt: time.Now()})
	if !plateau3 {
		t.Fatal("third sample should trigger plateau (two stagnant pairs)")
	}
}
