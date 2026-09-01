package agents

import (
	"fmt"
	"strings"
	"time"
)

// CoverageSample captures a single point-in-time measurement of test coverage.
type CoverageSample struct {
	Iteration      int
	TotalTests     int
	PassedTests    int
	UncoveredCount int
	CapturedAt     time.Time
}

// CoverageDelta measures the marginal coverage gain from one sample to the next.
// Positive values = improvement; zero or negative = stagnation or regression.
func CoverageDelta(prev, curr CoverageSample) (passedDelta, uncoveredDelta int) {
	return curr.PassedTests - prev.PassedTests, prev.UncoveredCount - curr.UncoveredCount
}

// PlateauDetector tracks recent coverage samples and decides when coverage
// gain has stagnated. The plateau rule: if the last K samples (default K=2)
// all show passedTestsDelta < MinPassedGain AND uncoveredCountDelta < MinUncoveredGain,
// the coverage loop should terminate.
//
// CodaMOSA reference: P49 ICSE 2023 / DOI 10.1109/ICSE48619.2023.00085 / carolemieux.com
// "the GA only iterates until convergence plateau is detected" — for code-translation
// we adopt the simpler rule: 2 consecutive samples with passedTests delta < 5
// OR uncovered count delta < 1 triggers plateau.
type PlateauDetector struct {
	Samples          []CoverageSample
	Window           int // last K samples to consider (default 2)
	MinPassedGain    int // minimum passed-tests delta to count as progress (default 5)
	MinUncoveredGain int // minimum uncovered-count delta to count as progress (default 1)
}

// NewPlateauDetector returns a PlateauDetector with CodaMOSA-style defaults.
func NewPlateauDetector() *PlateauDetector {
	return &PlateauDetector{
		Samples:          []CoverageSample{},
		Window:           2,
		MinPassedGain:    5,
		MinUncoveredGain: 1,
	}
}

// Record appends a sample to the history and returns true if a plateau is now detected.
// Records ARE retained (not popped) so the report can show the full trajectory.
func (p *PlateauDetector) Record(s CoverageSample) (plateau bool) {
	p.Samples = append(p.Samples, s)
	return p.IsPlateau()
}

// IsPlateau returns true if the last K samples all show insufficient progress.
func (p *PlateauDetector) IsPlateau() bool {
	if len(p.Samples) < p.Window+1 {
		return false // need at least Window+1 samples to compare
	}
	start := len(p.Samples) - p.Window - 1
	var stagnantPairs int
	for i := start; i < len(p.Samples)-1; i++ {
		prev := p.Samples[i]
		curr := p.Samples[i+1]
		pd, ud := CoverageDelta(prev, curr)
		if pd < p.MinPassedGain && ud < p.MinUncoveredGain {
			stagnantPairs++
		}
	}
	return stagnantPairs >= p.Window
}

// Summary returns a one-line summary suitable for log output.
func (p *PlateauDetector) Summary() string {
	if len(p.Samples) == 0 {
		return "plateau detector: no samples"
	}
	last := p.Samples[len(p.Samples)-1]
	var b strings.Builder
	fmt.Fprintf(&b, "plateau detector: %d samples; last passed=%d uncovered=%d; ",
		len(p.Samples), last.PassedTests, last.UncoveredCount)
	if p.IsPlateau() {
		b.WriteString("PLATEAU DETECTED")
	} else {
		b.WriteString("progressing")
	}
	return b.String()
}
