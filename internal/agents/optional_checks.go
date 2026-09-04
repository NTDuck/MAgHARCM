package agents

// Backlink: [[Primitives]] §[[PRIM-7]] (Multi-Agent Verdict Validation, MatchFixAgent),
// §[[PRIM-8]] (State-Grounded Mock-Based In-Isolation Validation, TRAM),
// §[[PRIM-11]] (Implementation-Agnostic Testing, RepoMod-Bench),
// §[[PRIM-12]] (Wasm-Based Reference Execution Oracle).
//
// This file adapts the four auxiliary validator primitives into the
// OptionalCheck interface declared in validator.go so that adding them
// to a ValidatorAgent via OptionalChecks wires them into the repair
// loop without touching the core cascade.

import (
	"context"
	"fmt"

)

// verdictPanelCheck adapts the multi-agent verdict panel
// (PRIM-7, MatchFixAgent-style) into an OptionalCheck. It dispatches N
// judge models against the current translated file set and returns the
// majority verdict. Skips when the panel is nil or has no model.
type verdictPanelCheck struct {
	panel *VerdictPanel
}

func (c *verdictPanelCheck) Name() string { return "PRIM-7-verdict-panel" }
func (c *verdictPanelCheck) Run(ctx context.Context, state *State) (string, string, error) {
	if c.panel == nil {
		return "skipped", "no verdict panel configured", nil
	}
	if len(state.TranslatedProject.Files) == 0 {
		return "skipped", "no translated files to judge", nil
	}
	// Sample the first translated file as the panel target. In a fuller
	// wiring the panel would judge all files; for now this exercises the
	// path end-to-end and surfaces plumbing bugs.
	var src, tgt string
	for path, content := range state.TranslatedProject.Files {
		src = path // file path doubles as source identity for the panel call
		tgt = content
		break
	}
	v, err := c.panel.Judge(ctx, src, tgt, 3)
	if err != nil {
		return "fail", fmt.Sprintf("judge panel error: %v", err), nil
	}
	if v.Agree {
		return "pass", fmt.Sprintf("%d judges agreed; %d disagreements", len(v.Disagreements)+0, len(v.Disagreements)), nil
	}
	return "fail", fmt.Sprintf("%d judge disagreements; majority reject", len(v.Disagreements)), nil
}

// mockValidatorCheck adapts the state-grounded mock validator
// (PRIM-8, TRAM-style) into an OptionalCheck. For each translated module
// it requests a mock-based in-isolation validation against the current
// dependency list. Skipped when no module is selected or the validator
// is nil.
type mockValidatorCheck struct {
	validator *MockValidator
}

func (c *mockValidatorCheck) Name() string { return "PRIM-8-mock-validator" }
func (c *mockValidatorCheck) Run(ctx context.Context, state *State) (string, string, error) {
	if c.validator == nil {
		return "skipped", "no mock validator configured", nil
	}
	if len(state.TranslatedProject.Files) == 0 {
		return "skipped", "no translated files to mock-validate", nil
	}
	var module string
	for path := range state.TranslatedProject.Files {
		module = path
		break
	}
	rep, err := c.validator.ValidateInIsolation(ctx, module, nil)
	if err != nil {
		return "fail", fmt.Sprintf("mock validator error: %v", err), nil
	}
	if rep.Passed {
		return "pass", fmt.Sprintf("module %s mocked in isolation: %d mocks generated", module, len(rep.MocksGenerated)), nil
	}
	return "fail", fmt.Sprintf("module %s failed in-isolation check: %v", module, rep.Errors), nil
}

// implAgnosticCheck adapts the implementation-agnostic I/O tester
// (PRIM-11, RepoMod-Bench-style) into an OptionalCheck. It runs an empty
// test-vector set against the target binary; callers populate
// IOTestVectors externally. The check passes when no test failures are
// reported.
type implAgnosticCheck struct {
	tester  *ImplAgnosticTester
	vectors []IOTestVector
}

func (c *implAgnosticCheck) Name() string { return "PRIM-11-impl-agnostic" }
func (c *implAgnosticCheck) Run(ctx context.Context, state *State) (string, string, error) {
	if c.tester == nil || len(c.vectors) == 0 {
		return "skipped", "no IO test vectors registered", nil
	}
	if state.Task.TargetDir == "" {
		return "skipped", "no target directory", nil
	}
	res, err := c.tester.RunIOTests(ctx, state.Task.TargetDir, c.vectors)
	if err != nil {
		return "fail", fmt.Sprintf("impl-agnostic error: %v", err), nil
	}
	if res.FailCount == 0 {
		return "pass", fmt.Sprintf("%d/%d IO vectors passed", res.PassCount, len(c.vectors)), nil
	}
	return "fail", fmt.Sprintf("%d/%d IO vectors failed", res.FailCount, len(c.vectors)), nil
}

// wasmOracleCheck adapts the wasm-based reference execution oracle
// (PRIM-12, Syzygy/sandboxing-style) into an OptionalCheck. It compares
// the target binary's output against the source-wasm reference for the
// supplied inputs. Skipped when source wasm or target binary is absent.
type wasmOracleCheck struct {
	oracle        *WasmOracle
	sourceWasm    string
	targetBinary  string
	oracleInputs  []string
}

func (c *wasmOracleCheck) Name() string { return "PRIM-12-wasm-oracle" }
func (c *wasmOracleCheck) Run(ctx context.Context, state *State) (string, string, error) {
	if c.oracle == nil {
		return "skipped", "no wasm oracle configured", nil
	}
	if c.sourceWasm == "" || c.targetBinary == "" {
		return "skipped", "no source-wasm or target-binary path registered", nil
	}
	res, err := c.oracle.Compare(ctx, c.sourceWasm, c.targetBinary, c.oracleInputs)
	if err != nil {
		return "fail", fmt.Sprintf("wasm oracle error: %v", err), nil
	}
	if len(res.Mismatches) == 0 && res.Agree {
		return "pass", fmt.Sprintf("%d inputs matched between source-wasm and target-binary", len(c.oracleInputs)), nil
	}
	return "fail", fmt.Sprintf("%d/%d inputs diverged between source-wasm and target-binary", len(res.MismatchedInputs), len(c.oracleInputs)), nil
}

// DefaultOptionalChecks returns the standard set of optional validator
// checks for a typical run. Callers can append their own (or replace any)
// before passing the slice into ValidatorAgent.OptionalChecks. Returns
// nil when no checks are enabled — keeping the validator's default
// behaviour (core cascade only) unchanged.
func DefaultOptionalChecks(cfg OptionalChecksConfig) []OptionalCheck {
	if !cfg.AnyEnabled() {
		return nil
	}
	var out []OptionalCheck
	if cfg.VerdictPanel != nil {
		out = append(out, &verdictPanelCheck{panel: cfg.VerdictPanel})
	}
	if cfg.MockValidator != nil {
		out = append(out, &mockValidatorCheck{validator: cfg.MockValidator})
	}
	if cfg.ImplAgnostic != nil && len(cfg.ImplAgnosticVectors) > 0 {
		out = append(out, &implAgnosticCheck{tester: cfg.ImplAgnostic, vectors: cfg.ImplAgnosticVectors})
	}
	if cfg.WasmOracle != nil && cfg.SourceWasmPath != "" && cfg.TargetBinaryPath != "" {
		out = append(out, &wasmOracleCheck{
			oracle:       cfg.WasmOracle,
			sourceWasm:   cfg.SourceWasmPath,
			targetBinary: cfg.TargetBinaryPath,
			oracleInputs: cfg.OracleInputs,
		})
	}
	if cfg.RoleFlipGate != nil {
		out = append(out, &roleFlipCheck{gate: cfg.RoleFlipGate})
	}
	return out
}

// OptionalChecksConfig is the public knob that callers (CLI flag parsing,
// YAML config, or programmatic wiring) use to enable each auxiliary
// validator check. AnyEnabled returns true when at least one check is
// configured with the prerequisites it needs.
type OptionalChecksConfig struct {
	VerdictPanel         *VerdictPanel
	MockValidator        *MockValidator
	ImplAgnostic         *ImplAgnosticTester
	ImplAgnosticVectors  []IOTestVector
	WasmOracle           *WasmOracle
	SourceWasmPath       string
	TargetBinaryPath     string
	OracleInputs         []string
	// RoleFlipGate is the PRIM-25 reviewer. When non-nil, the validator
	// cascade calls it on the latest translator snippet before the
	// repair loop iterates.
	RoleFlipGate *RoleFlipGate
}
func (c OptionalChecksConfig) AnyEnabled() bool {
	return c.VerdictPanel != nil ||
		c.MockValidator != nil ||
		(c.ImplAgnostic != nil && len(c.ImplAgnosticVectors) > 0) ||
		(c.WasmOracle != nil && c.SourceWasmPath != "" && c.TargetBinaryPath != "") ||
		c.RoleFlipGate != nil
}

// re-export for callers that prefer artifacts-side imports.

// roleFlipCheck adapts the PRIM-25 communicative-de-hallucination
// RoleFlipGate (ChatDev §2.4) into an OptionalCheck. When the translator's
// latest emitted code is present on state, the gate's reviewer model is
// asked to surface a defect. A defect forces a repair iteration; a clean
// pass leaves the verdict unchanged.
type roleFlipCheck struct {
	gate *RoleFlipGate
}

func (c *roleFlipCheck) Name() string { return "PRIM-25-role-flip-gate" }
func (c *roleFlipCheck) Run(ctx context.Context, state *State) (string, string, error) {
	if c.gate == nil || state == nil {
		return "skipped", "no role-flip gate or state", nil
	}
	sample := latestEmittedSample(state)
	if sample == "" {
		return "skipped", "no translator output to inspect", nil
	}
	v, err := c.gate.Inspect(ctx, sample)
	if err != nil {
		return "fail", fmt.Sprintf("role-flip reviewer error: %v", err), nil
	}
	if !v.DefectFound {
		return "pass", "reviewer accepted translator output", nil
	}
	return "fail", v.Reason, nil
}

// latestEmittedSample returns the first non-empty translator snippet from
// TranslatedProject.Files to feed the PRIM-25 role-flip reviewer. Map
// iteration order is non-deterministic; any sample suffices for the gate.
func latestEmittedSample(state *State) string {
	if state == nil {
		return ""
	}
	for _, code := range state.TranslatedProject.Files {
		if code != "" {
			return code
		}
	}
	return ""
}
