// Package agents: recruit.go implements Recruitment-Adaptive Planning
// (paper appendix primitive P29, cited as AgentVerse [P26] §2).
//
// At the end of each iteration the orchestrator receives a ValidationSummary
// describing how the last migration pass fared. The Recruiter inspects those
// signals — compilation success, test pass-rate, plateau detection, and
// adversarial weakening — and produces a RecruitmentPlan listing the tools
// and downstream agents that should run next. This lets the pipeline adapt
// to regressions without re-running the full pipeline blindly.
package agents

import (
	"context"

	"MAgHARCM/internal/logger"
)

// ValidationSummary is the per-iteration verdict consumed by the Recruiter.
// It is computed by the validator+verdict-panel stage and passed verbatim
// into Recruit so the selection policy stays a pure function of inputs.
type ValidationSummary struct {
	// CompilationSuccess is true when the migrated tree builds without errors.
	CompilationSuccess bool
	// PassRate is the fraction of tests passing in [0,1].
	PassRate float64
	// PlateauDetected is true when consecutive iterations show negligible
	// improvement in PassRate (see internal/agents/plateau.go PRIM-27).
	PlateauDetected bool
	// AdversarialWeakeningDetected is true when the adversarial suite shows
	// the migration regressing under pressure (see internal/agents/roleflip.go).
	AdversarialWeakeningDetected bool
}

// RecruitmentPlan is the Recruiter's output: the next iteration's tool and
// agent roster plus a human-readable rationale suitable for logging.
type RecruitmentPlan struct {
	Tools     []string
	Agents    []string
	Rationale string
}

// Recruiter selects the next iteration's roster from the prior ValidationSummary
// and the static Profile. It holds no per-call state; all decisions live in
// the selection policy documented below.
type Recruiter struct{}

// Default policy:
//  1. Compilation failed → re-run the validator (the previous translator
//     produced an unbuildable tree; re-validation surfaces the exact errors).
//  2. Plateau detected → re-run the chunked translator with the
//     `recruit_translator_v2` flavour to break out of the local optimum.
//  3. Adversarial weakening → re-run the role-flip gate to collect
//     compensating evidence.
//  4. Otherwise → standard pipeline (translator → validator → verifier).
//
// NewRecruiter returns a ready-to-use Recruiter.
func NewRecruiter() *Recruiter { return &Recruiter{} }

// Recruit applies the selection policy and returns the next iteration's plan.
// It never returns an error: every input maps to a defined plan.
func (r *Recruiter) Recruit(ctx context.Context, profile Profile, lastReport ValidationSummary) (RecruitmentPlan, error) {
	_ = ctx
	_ = profile

	switch {
	case !lastReport.CompilationSuccess:
		logger.LogStep("recruiter: compilation failed → re-run validator")
		return RecruitmentPlan{
			Tools:     []string{"validator", "diagnostics"},
			Agents:    []string{"validator"},
			Rationale: "Compilation failed on the prior pass; re-run the validator to surface concrete errors before any further translation.",
		}, nil

	case lastReport.PlateauDetected:
		logger.LogStep("recruiter: plateau detected → recruit_translator_v2")
		return RecruitmentPlan{
			Tools:     []string{"chunked_translator", "recruit_translator_v2", "plateau_breaker"},
			Agents:    []string{"chunked_translator"},
			Rationale: "Plateau detected; re-run the chunked translator with the recruit_translator_v2 flavour to escape the local optimum.",
		}, nil

	case lastReport.AdversarialWeakeningDetected:
		logger.LogStep("recruiter: adversarial weakening → re-run role-flip gate")
		return RecruitmentPlan{
			Tools:     []string{"roleflip_gate", "adversarial_suite"},
			Agents:    []string{"roleflip"},
			Rationale: "Adversarial weakening detected; re-run the role-flip gate to collect compensating evidence before resuming translation.",
		}, nil

	default:
		logger.LogStep("recruiter: standard pipeline")
		return RecruitmentPlan{
			Tools:     []string{"chunked_translator", "validator", "verifier"},
			Agents:    []string{"chunked_translator", "validator", "verdict_panel"},
			Rationale: "Standard pipeline: chunked translator followed by validator and verdict panel.",
		}, nil
	}
}
