// Backlink: [[Primitives]] §NEW-PRIM-12 (Wasm-Based Reference Execution Oracle).
package agents

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"MAgHARCM/internal/logger"
)

// wasmOracleTimeout caps each per-input invocation so a runaway target or
// Wasm runtime cannot stall the comparison loop indefinitely.
const wasmOracleTimeout = 30 * time.Second

// OracleResult summarises the per-input behavioural agreement between a
// reference Wasm binary and the migrated native binary.
//
// Agree is true only when every observed channel (stdout, stderr, exit code)
// matched for every input. Mismatches enumerates human-readable descriptions
// of every disagreement (input + channel); MismatchedInputs lists the inputs
// that produced any disagreement.
type OracleResult struct {
	Mismatches      []string
	MismatchedInputs []string
	Agree           bool
}

// WasmOracle compares observable behaviour of a compiled WebAssembly module
// (the reference implementation) against the translated native target binary
// across a list of inputs. Reference: paper appendix primitive VERT [[P12]].
//
// The oracle deliberately treats the Wasm artifact as already compiled.
// Compiling source → Wasm lives in a future hook (the spec notes this as the
// preceding step in the pipeline; not implemented here).
type WasmOracle struct{}

// NewWasmOracle returns a ready-to-use oracle.
func NewWasmOracle() *WasmOracle { return &WasmOracle{} }

// Compare executes both binaries for each input and reports whether their
// stdout, stderr, and exit codes agree. An empty input slice is an error:
// the caller is expected to drive the oracle with at least one concrete
// stimulus. If sourceWasmPath does not exist, Compare returns an error
// without invoking the target binary, per the [[P12]] contract that the
// reference Wasm must already be present.
//
// Future hook: a compiler step that materialises sourceWasmPath from the
// reference source tree will sit upstream of this method.
func (o *WasmOracle) Compare(ctx context.Context, sourceWasmPath, targetBinaryPath string, inputs []string) (OracleResult, error) {
	if len(inputs) == 0 {
		return OracleResult{}, errors.New("wasm oracle: inputs must contain at least one entry")
	}
	if _, err := os.Stat(sourceWasmPath); err != nil {
		return OracleResult{}, fmt.Errorf("wasm oracle: reference wasm missing at %s: %w", sourceWasmPath, err)
	}
	if _, err := os.Stat(targetBinaryPath); err != nil {
		return OracleResult{}, fmt.Errorf("wasm oracle: target binary missing at %s: %w", targetBinaryPath, err)
	}

	result := OracleResult{}
	for _, in := range inputs {
		wOut, wErr, wExit, wErrRun := runOnce(ctx, sourceWasmPath, in)
		tOut, tErr, tExit, tErrRun := runOnce(ctx, targetBinaryPath, in)
		if wErrRun != nil || tErrRun != nil {
			result.record(in, fmt.Sprintf("input %q: invocation failed (wasm=%v target=%v)", in, wErrRun, tErrRun))
			continue
		}
		if !bytes.Equal(wOut, tOut) {
			result.record(in, fmt.Sprintf("input %q: stdout differs", in))
		}
		if !bytes.Equal(wErr, tErr) {
			result.record(in, fmt.Sprintf("input %q: stderr differs", in))
		}
		if wExit != tExit {
			result.record(in, fmt.Sprintf("input %q: exit code differs (wasm=%d target=%d)", in, wExit, tExit))
		}
	}
	result.Agree = len(result.Mismatches) == 0
	logger.LogValidation("wasm oracle: %d input(s), %d mismatch(es), agree=%t", len(inputs), len(result.Mismatches), result.Agree)
	return result, nil
}

// record appends a mismatch description and the offending input exactly once.
func (r *OracleResult) record(input, msg string) {
	r.Mismatches = append(r.Mismatches, msg)
	r.MismatchedInputs = append(r.MismatchedInputs, input)
}

// runOnce executes path with the given input as argv[1] (positional CLI arg
// form, matching the [[P12]] test-harness convention) and captures stdout,
// stderr, and exit code. The supplied context bounds the run.
func runOnce(ctx context.Context, path, input string) (stdout, stderr []byte, exit int, err error) {
	cctx, cancel := context.WithTimeout(ctx, wasmOracleTimeout)
	defer cancel()
	cmd := exec.CommandContext(cctx, path, input)
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
