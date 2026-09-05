// Backlink: [[Primitives]] §NEW-PRIM-11 (Implementation-Agnostic Testing).
package agents

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"

	"MAgHARCM/internal/logger"
)

// implAgnosticTimeout caps each per-vector invocation so a runaway target
// binary cannot stall the comparison loop indefinitely.
const implAgnosticTimeout = 30 * time.Second

// IOTestVector is a black-box test case: drive targetBinary with Inputs and
// compare the observed outputs and exit code against the expected ones. No
// source-level inspection happens here — that is what makes the primitive
// implementation-agnostic (RepoMod-Bench [[P11]]).
type IOTestVector struct {
	Name             string
	Inputs           []string
	ExpectedOutputs  []string
	ExpectedExitCode int
}

// IOResult summarises how many vectors passed / failed and lists the names
// of failing vectors so callers can surface diagnostics without rebuilding
// state.
type IOResult struct {
	PassCount int
	FailCount int
	Failures  []string
}

// ImplAgnosticTester drives a target binary against a slice of IOTestVector
// and verifies the observed outputs / exit code via plain string equality.
// It deliberately inspects nothing about the binary's source, AST, or build
// configuration — pass/fail is purely observational.
type ImplAgnosticTester struct{}

// NewImplAgnosticTester returns a zero-value tester ready for use.
func NewImplAgnosticTester() *ImplAgnosticTester {
	return &ImplAgnosticTester{}
}

// RunIOTests executes targetBinary once per vector, passing each element of
// Inputs as a positional argv entry. For every invocation it captures stdout,
// stderr, and exit code, then compares:
//
//   - the joined stdout lines (split on '\n', trailing empty stripped)
//     against ExpectedOutputs element-by-element
//   - the observed exit code against ExpectedExitCode
//
// Comparison is plain string equality; no normalisation, no AST inspection.
// The first disagreement in any vector causes that vector to be recorded as
// a failure but does not abort the remaining vectors. If ctx is cancelled
// before all vectors run, the partially populated IOResult is returned with
// the cancellation error.
func (t *ImplAgnosticTester) RunIOTests(ctx context.Context, targetBinary string, vectors []IOTestVector) (IOResult, error) {
	if targetBinary == "" {
		return IOResult{}, errors.New("impl_agnostic: targetBinary is empty")
	}

	res := IOResult{}
	logger.LogAgent("ImplAgnosticTester", "running %d vectors against %s", len(vectors), targetBinary)

	for _, v := range vectors {
		if err := ctx.Err(); err != nil {
			res.Failures = append(res.Failures, v.Name)
			return res, err
		}

		stdout, _, exitCode, runErr := runVectorOnce(ctx, targetBinary, v.Inputs)
		if runErr != nil {
			// runVectorOnce only returns a non-nil error for context
			// cancellation or true exec failures (binary missing, etc.).
			// An exit-code mismatch is encoded inside exitCode, not here.
			res.FailCount++
			res.Failures = append(res.Failures, v.Name)
			logger.LogWarning("impl_agnostic: vector %q exec error: %v", v.Name, runErr)
			continue
		}

		if !outputsEqual(stdout, v.ExpectedOutputs) || exitCode != v.ExpectedExitCode {
			res.FailCount++
			res.Failures = append(res.Failures, v.Name)
			logger.LogValidation("impl_agnostic: vector %q failed: exit=%d want=%d", v.Name, exitCode, v.ExpectedExitCode)
			continue
		}

		res.PassCount++
	}

	logger.LogAgent("ImplAgnosticTester", "done: pass=%d fail=%d", res.PassCount, res.FailCount)
	return res, nil
}

// runVectorOnce executes targetBinary once with the given argv slice and
// returns captured stdout, stderr, and exit code. An exit-code mismatch
// (non-zero exit) is NOT treated as an error — it is part of the observable
// behaviour the caller compares. Only true exec failures (binary not found,
// spawn error) or context cancellation surface as a non-nil error.
func runVectorOnce(ctx context.Context, targetBinary string, inputs []string) (stdout, stderr []byte, exit int, err error) {
	cctx, cancel := context.WithTimeout(ctx, implAgnosticTimeout)
	defer cancel()

	args := append([]string{}, inputs...)
	cmd := exec.CommandContext(cctx, targetBinary, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	runErr := cmd.Run()
	stdout = outBuf.Bytes()
	stderr = errBuf.Bytes()
	if runErr != nil {
		var ee *exec.ExitError
		if errors.As(runErr, &ee) {
			return stdout, stderr, ee.ExitCode(), nil
		}
		return stdout, stderr, -1, runErr
	}
	return stdout, stderr, 0, nil
}

// outputsEqual splits observed stdout on '\n' and compares the resulting
// line slice against expected via element-wise string equality. A trailing
// empty line (from a final '\n') is stripped so "echo hi" producing
// "hi\n" still matches ExpectedOutputs = ["hi"].
func outputsEqual(observed []byte, expected []string) bool {
	var lines []string
	if len(observed) > 0 {
		raw := string(observed)
		if raw[len(raw)-1] == '\n' {
			raw = raw[:len(raw)-1]
		}
		lines = splitLines(raw)
	} else {
		lines = []string{}
	}
	if len(lines) != len(expected) {
		return false
	}
	for i := range lines {
		if lines[i] != expected[i] {
			return false
		}
	}
	return true
}

// splitLines is bytes-split-on-'\n' without pulling in strings.Split for a
// one-liner. Kept local so the primitive stays free of helpers from sibling
// files.
func splitLines(s string) []string {
	if s == "" {
		return []string{}
	}
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}
