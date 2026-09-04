// Backlink: [[Methodology]] §1 and [[Primitives]] §NEW-PRIM-25.
// P25 — Communicative-De-hallucination Role-Flip Gate (ChatDev [P46] §2.4).
// The gate asks a critic model to find at least one defect in a translator's
// output. Role-flipping (author → reviewer) is the communicative de-hallucination
// trick: the same model is more likely to spot a hallucination when prompted as
// a fault-finder than when asked to confirm its own work.

package agents

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/logger"
)

// roleFlipSystemPrompt frames the model as a critical reviewer that must
// surface at least one defect — the adversarial framing is the heart of P25.
const (
	roleFlipSystemPrompt   = "you are a critical reviewer who must find at least one bug in the code below"
	reviewerAcceptedReason = "reviewer accepted"
	retryHintForDefect     = "re-check for common bug classes"
	noDefectToken          = "NO_DEFECT"
)

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

// NewRoleFlipGate wires a RoleFlipGate to the supplied chat model.
func NewRoleFlipGate(m model.BaseChatModel) *RoleFlipGate {
	return &RoleFlipGate{Model: m}
}

// Inspect runs the role-flipped review over translatorOutput and returns the
// gate's verdict. The reviewer is instructed to either surface a defect or
// reply with NO_DEFECT; an empty or malformed reply is treated as acceptance
// (DefectFound=false) so a chatty base model cannot deadlock the pipeline.
func (g *RoleFlipGate) Inspect(ctx context.Context, translatorOutput string) (RoleFlipVerdict, error) {
	if g == nil || g.Model == nil {
		return RoleFlipVerdict{}, errRoleFlipGateNotConfigured
	}

	logger.LogAgent("RoleFlipGate", "Invoking role-flipped reviewer on %d bytes of translator output", len(translatorOutput))

	resp, err := g.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage(roleFlipSystemPrompt),
		schema.UserMessage(translatorOutput),
	})
	if err != nil {
		return RoleFlipVerdict{}, err
	}

	reply := strings.TrimSpace(resp.Content)
	if reply == "" || strings.EqualFold(reply, noDefectToken) {
		logger.LogAgent("RoleFlipGate", "Reviewer accepted translator output")
		return RoleFlipVerdict{DefectFound: false, Reason: reviewerAcceptedReason}, nil
	}

	logger.LogValidation("RoleFlipGate surfaced defect: %s", truncate(reply, 160))
	return RoleFlipVerdict{DefectFound: true, Reason: reply, RetryHint: retryHintForDefect}, nil
}

var errRoleFlipGateNotConfigured = errors.New("roleflip: Model is nil")

// truncate caps a reviewer remark for logging so a runaway generation cannot
// blow up the structured log buffer.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
