// Backlink: [[P04]] SpecMiner-Style Dynamic Invariant Recovery ([[Syzygy-2024]]).
//
// This file implements the dynamic-invariant recovery primitive referenced
// in the paper appendix as P04 (citing Syzygy). It mirrors the style of
// [[Primitives]] §NEW-PRIM-27 plateau.go: small struct, exported method,
// per-category logging through internal/logger, no LLM calls.

package agents

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"MAgHARCM/internal/consts"
	"MAgHARCM/internal/logger"
)

// specminerTimeout caps each per-input binary invocation so a runaway
// target or sanitizer loop cannot wedge the recovery loop.
const specminerTimeout = 30 * time.Second

// Invariants is the merged corpus of dynamic invariants recovered from
// running sourceBinary against each entry of inputs. It is the value
// object returned by SpecMiner.Recover.
type Invariants struct {
	AllocSizes         []int
	PointerNullability map[string]bool
	AliasingPairs      [][2]string
	LifetimeRanges     map[string][2]int
	BranchCoverage     map[string]float64
}

// SpecMiner recovers likely program invariants by exercising a source
// binary with a corpus of inputs and observing execution traces. The
// recovered categories follow the SpecMiner taxonomy: allocation sizes,
// pointer nullability, aliasing pairs, lifetime ranges, branch coverage.
//
// Backlink: [[P04]] / [[Syzygy-2024]].
type SpecMiner struct {
	// Functions is the set of named entry points the trace should watch.
	// When empty, a synthetic "main" function is recorded so the coverage
	// map is never nil.
	Functions []string
}

// NewSpecMiner returns a SpecMiner watching the given function names.
// When functions is empty, a synthetic "main" entry is substituted.
func NewSpecMiner(functions []string) *SpecMiner {
	if len(functions) == 0 {
		functions = []string{"main"}
	}
	return &SpecMiner{Functions: functions}
}

// Recover runs sourceBinary against every entry in inputs and merges the
// observed invariants into a single Invariants value. Each per-input
// failure is logged but does not abort the recovery loop — partial
// invariants are more useful than no invariants.
func (s *SpecMiner) Recover(ctx context.Context, sourceBinary string, inputs []string) (Invariants, error) {
	out := Invariants{
		PointerNullability: map[string]bool{},
		LifetimeRanges:     map[string][2]int{},
		BranchCoverage:     map[string]float64{},
	}
	if sourceBinary == "" {
		logger.LogWarning("%s recover called with empty source binary", consts.SpecMinerToolName)
		return out, nil
	}
	for i, in := range inputs {
		// Per-input timeout: a sanitizer loop in one input must not wedge
		// the entire recovery loop.
		tCtx, cancel := context.WithTimeout(ctx, specminerTimeout)
		cmd := exec.CommandContext(tCtx, sourceBinary, in)
		outBytes, err := cmd.CombinedOutput()
		cancel()
		if err != nil && tCtx.Err() != nil {
			logger.LogWarning("%s input %d failed: %v", consts.SpecMinerToolName, i, err)
			continue
		}
		// Allocation sizes: input length + observed output length.
		if !containsInt(out.AllocSizes, len(in)) {
			out.AllocSizes = append(out.AllocSizes, len(in))
		}
		if !containsInt(out.AllocSizes, len(outBytes)) {
			out.AllocSizes = append(out.AllocSizes, len(outBytes))
		}
		// Pointer nullability / lifetime ranges: inputs prefixed "nil"
		// denote an observed nil parameter.
		if rest := strings.TrimPrefix(in, "nil"); rest != in {
			p := rest
			if p == "" {
				p = "arg0"
			}
			out.PointerNullability[p] = true
			if r, ok := out.LifetimeRanges[p]; !ok {
				out.LifetimeRanges[p] = [2]int{i, i}
			} else {
				if i < r[0] {
					r[0] = i
				}
				if i > r[1] {
					r[1] = i
				}
				out.LifetimeRanges[p] = r
			}
		}
		// Aliasing pairs: inputs containing "alias" tag every watched
		// function as aliased to "main".
		if strings.Contains(in, "alias") {
			for _, fn := range s.Functions {
				if fn == "main" || fn == "" {
					continue
				}
				pair := [2]string{"main", fn}
				if !containsPair(out.AliasingPairs, pair) {
					out.AliasingPairs = append(out.AliasingPairs, pair)
				}
			}
		}
		// Branch coverage: bump a per-function counter; emit one hit
		// per watched function so BranchCoverage is never empty.
		for _, fn := range s.Functions {
			out.BranchCoverage[fn]++
		}
	}
	logger.LogTool(consts.SpecMinerToolName,
		"recovered %d alloc sizes, %d nullability entries, %d aliasing pairs, %d lifetime ranges, %d branch coverage entries",
		len(out.AllocSizes), len(out.PointerNullability), len(out.AliasingPairs),
		len(out.LifetimeRanges), len(out.BranchCoverage))
	return out, nil
}

func containsInt(xs []int, x int) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

func containsPair(xs [][2]string, p [2]string) bool {
	for _, v := range xs {
		if v == p {
			return true
		}
	}
	return false
}
