// Backlink: [[Methodology]] §1 and [[Primitives]] §NEW-PRIM-25.
// P25 — Communicative-De-hallucination Role-Flip Gate (ChatDev [P46] §2.4).
// The gate asks a critic model to find at least one defect in a translator's
// output. Role-flipping (author → reviewer) is the communicative de-hallucination
// trick: the same model is more likely to spot a hallucination when prompted as
// a fault-finder than when asked to confirm its own work.

package agents

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/logger"
)

// roleflipLogTruncateLen is the maximum length (in bytes) the reviewer
// remark is allowed to occupy in the structured log buffer before
// truncate() caps it. Local constant; centralised values live in
// internal/compiletime.
const roleflipLogTruncateLen = 160

// RoleFlipVerdict is the gate's structured outcome. DefectFound drives whether the
// pipeline re-invokes the translator; Reason and RetryHint feed the next pass.
type RoleFlipVerdict struct {
	DefectFound bool
	Reason      string
	RetryHint   string
}

// RoleFlipGate runs the communicative de-hallucination check on a translator's
// output. Construct it with the chat model that should play the reviewer role
// (typically the reasoning model — see internal/llm/llm.go).
type RoleFlipGate struct {
	Model model.BaseChatModel
}

// NewRoleFlipGate wires a RoleFlipGate to the supplied chat model. Pass nil
// to obtain a disabled gate (Inspect returns ErrRoleFlipGateNotConfigured).
func NewRoleFlipGate(m model.BaseChatModel) *RoleFlipGate {
	return &RoleFlipGate{Model: m}
}

// MustRoleFlipGate wires a RoleFlipGate to the supplied chat model and panics
// if the model is nil. Use this at startup where a missing reviewer model is
// a fatal configuration error rather than a runtime fallback.
func MustRoleFlipGate(m model.BaseChatModel) *RoleFlipGate {
	if m == nil {
		panic("compiletime.MustRoleFlipGate: Model must not be nil")
	}
	return NewRoleFlipGate(m)
}
// Inspect runs the role-flipped review over translatorOutput and returns the
// gate's verdict. The reviewer is instructed to either surface a defect or
// reply with NO_DEFECT; an empty or malformed reply is treated as acceptance
// (DefectFound=false) so a chatty base model cannot deadlock the pipeline.
func (g *RoleFlipGate) Inspect(ctx context.Context, translatorOutput string) (RoleFlipVerdict, error) {
	if g == nil || g.Model == nil {
		return RoleFlipVerdict{}, compiletime.ErrRoleFlipGateNotConfigured
	}

	logger.LogAgent(compiletime.LogScopeRoleFlip, "Invoking role-flipped reviewer on %d bytes of translator output", len(translatorOutput))

	resp, err := g.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage(compiletime.RoleFlipSystemPrompt),
		schema.UserMessage(translatorOutput),
	})
	if err != nil {
		return RoleFlipVerdict{}, err
	}

	reply := strings.TrimSpace(resp.Content)
	if reply == "" || strings.EqualFold(reply, compiletime.RoleFlipNoDefectToken) {
		return RoleFlipVerdict{DefectFound: false, Reason: compiletime.RoleFlipAcceptedReason}, nil
	}
	logger.LogValidation("RoleFlipGate surfaced defect: %s", truncate(reply, roleflipLogTruncateLen))
	return RoleFlipVerdict{DefectFound: true, Reason: reply, RetryHint: compiletime.RoleFlipRetryHint}, nil
}

// truncate caps a reviewer remark for logging so a runaway generation cannot
// blow up the structured log buffer.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
