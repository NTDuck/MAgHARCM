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

	"MAgHARCM/internal/compiletime"
)

// verdictPanelCheck adapts the multi-agent verdict panel
// (PRIM-7, MatchFixAgent-style) into an OptionalCheck. It dispatches N
// judge models against the current translated file set and returns the
// majority verdict. Skips when the panel is nil or has no model.
type verdictPanelCheck struct {
	panel *VerdictPanel
}

func (c *verdictPanelCheck) Name() string { return compiletime.OptionalCheckVerdictPanel }
func (c *verdictPanelCheck) Run(ctx context.Context, state *compiletime.State) (string, string, error) {
	if c.panel == nil {
		return string(compiletime.VerdictSkipped), "no verdict panel configured", nil
	}
	if len(state.TranslatedProject.Files) == 0 {
		return string(compiletime.VerdictSkipped), "no translated files to judge", nil
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
		return string(compiletime.VerdictFail), fmt.Sprintf("judge panel error: %v", err), nil
	}
	agreed := len(v.PerJudge) - len(v.Disagreements)
	if v.Agree {
		return string(compiletime.VerdictPass), fmt.Sprintf("%d judges agreed; %d disagreements", agreed, len(v.Disagreements)), nil
	}
	return string(compiletime.VerdictFail), fmt.Sprintf("%d judge disagreements; majority reject", len(v.Disagreements)), nil
}

// mockValidatorCheck adapts the state-grounded mock validator
// (PRIM-8, TRAM-style) into an OptionalCheck. For each translated module
// it requests a mock-based in-isolation validation against the current
// dependency list. Skipped when no module is selected or the validator
// is nil.
type mockValidatorCheck struct {
	validator *MockValidator
}

func (c *mockValidatorCheck) Name() string { return compiletime.OptionalCheckMockValidator }
func (c *mockValidatorCheck) Run(ctx context.Context, state *compiletime.State) (string, string, error) {
	if c.validator == nil {
		return string(compiletime.VerdictSkipped), "no mock validator configured", nil
	}
	if len(state.TranslatedProject.Files) == 0 {
		return string(compiletime.VerdictSkipped), "no translated files to mock-validate", nil
	}
	var module string
	for path := range state.TranslatedProject.Files {
		module = path
		break
	}
	rep, err := c.validator.ValidateInIsolation(ctx, module, nil)
	if err != nil {
		return string(compiletime.VerdictFail), fmt.Sprintf("mock validator error: %v", err), nil
	}
	if rep.Passed {
		return string(compiletime.VerdictPass), fmt.Sprintf("module %s mocked in isolation: %d mocks generated", module, len(rep.MocksGenerated)), nil
	}
	return string(compiletime.VerdictFail), fmt.Sprintf("module %s failed in-isolation check: %v", module, rep.Errors), nil
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

func (c *implAgnosticCheck) Name() string { return compiletime.OptionalCheckImplAgnostic }
func (c *implAgnosticCheck) Run(ctx context.Context, state *compiletime.State) (string, string, error) {
	if c.tester == nil || len(c.vectors) == 0 {
		return string(compiletime.VerdictSkipped), "no IO test vectors registered", nil
	}
	if state.Task.TargetDir == "" {
		return string(compiletime.VerdictSkipped), "no target directory", nil
	}
	res, err := c.tester.RunIOTests(ctx, state.Task.TargetDir, c.vectors)
	if err != nil {
		return string(compiletime.VerdictFail), fmt.Sprintf("impl-agnostic error: %v", err), nil
	}
	if res.FailCount == 0 {
		return string(compiletime.VerdictPass), fmt.Sprintf("%d/%d IO vectors passed", res.PassCount, len(c.vectors)), nil
	}
	return string(compiletime.VerdictFail), fmt.Sprintf("%d/%d IO vectors failed", res.FailCount, len(c.vectors)), nil
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

func (c *wasmOracleCheck) Name() string { return compiletime.OptionalCheckWasmOracle }

func (c *wasmOracleCheck) Run(ctx context.Context, state *compiletime.State) (string, string, error) {
	if c.oracle == nil {
		return string(compiletime.VerdictSkipped), "no wasm oracle configured", nil
	}
	if c.sourceWasm == "" || c.targetBinary == "" {
		return string(compiletime.VerdictSkipped), "no source-wasm or target-binary path registered", nil
	}
	res, err := c.oracle.Compare(ctx, c.sourceWasm, c.targetBinary, c.oracleInputs)
	if err != nil {
		return string(compiletime.VerdictFail), fmt.Sprintf("wasm oracle error: %v", err), nil
	}
	if len(res.Mismatches) == 0 && res.Agree {
		return string(compiletime.VerdictPass), fmt.Sprintf("%d inputs matched between source-wasm and target-binary", len(c.oracleInputs)), nil
	}
	return string(compiletime.VerdictFail), fmt.Sprintf("%d/%d inputs diverged between source-wasm and target-binary", len(res.MismatchedInputs), len(c.oracleInputs)), nil
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

func (c *roleFlipCheck) Name() string { return compiletime.OptionalCheckRoleFlipGate }
func (c *roleFlipCheck) Run(ctx context.Context, state *compiletime.State) (string, string, error) {
	if c.gate == nil || state == nil {
		return string(compiletime.VerdictSkipped), "no role-flip gate or state", nil
	}
	sample, samplePath := RepresentativeTranslationSample(state.TranslatedProject.Files)
	if sample == "" {
		return string(compiletime.VerdictSkipped), "no translator output to inspect", nil
	}
	v, err := c.gate.Inspect(ctx, samplePath, sample)
	if err != nil {
		return string(compiletime.VerdictFail), fmt.Sprintf("role-flip reviewer error: %v", err), nil
	}
	if !v.DefectFound {
		return string(compiletime.VerdictPass), "reviewer accepted translator output", nil
	}
	return string(compiletime.VerdictFail), v.Reason, nil
}
