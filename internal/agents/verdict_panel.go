package agents

// Backlink: [[Primitives]] §NEW-PRIM-07 (Multi-Agent Verdict Validation, also known
// as the MatchFixAgent described in §P08 of the paper appendix). The panel fans
// the equivalence decision across N independent LLM judges so a single
// hallucinated vote cannot flip a wrong verdict, then aggregates by strict
// majority.

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"

	"MAgHARCM/internal/compiletime"
	"MAgHARCM/internal/logger"
)

// VerdictKindEquivalent / VerdictKindNotEquivalent are the back-compat
// aliases for compiletime.VerdictEquivalent / compiletime.VerdictNotEquivalent.
// New code should reference the compiletime constants directly.
const (
	VerdictKindEquivalent    = compiletime.VerdictEquivalent
	VerdictKindNotEquivalent = compiletime.VerdictNotEquivalent
)

// DefaultPanelSize is the canonical PRIM-7 panel size (5 judges — odd,
// absorbs one hallucinated vote, fits one Ollama session latency budget).
const DefaultPanelSize = 5

// Verdict is the aggregated panel decision.
type Verdict struct {
	Agree         bool
	Disagreements []string
	PerJudge      []JudgeOpinion
}

// JudgeOpinion records a single panel member's vote and reasoning.
type JudgeOpinion struct {
	JudgeID   string
	Verdict   string
	Rationale string
}

// VerdictPanel fans an equivalence question across N independent LLM judges
// and aggregates the votes. The chat model must be a small, deterministic
// reasoning model (typically Models.Reasoning) so judges see roughly the
// same world view and disagreements are meaningful rather than noise.
type VerdictPanel struct {
	Model model.BaseChatModel
}

// NewVerdictPanel constructs a VerdictPanel backed by the given chat model.
// Pass llm.Models.Reasoning so the panel uses the smaller reasoning model
// rather than the coding model.
func NewVerdictPanel(m model.BaseChatModel) *VerdictPanel {
	return &VerdictPanel{Model: m}
}

// MustVerdictPanel constructs a VerdictPanel and panics if the chat model
// is nil. Use this at startup where a missing reasoning model is a fatal
// configuration error.
func MustVerdictPanel(m model.BaseChatModel) *VerdictPanel {
	if m == nil {
		panic("compiletime.MustVerdictPanel: Model must not be nil")
	}
	return NewVerdictPanel(m)
}

// MustPanelSize returns n when positive; otherwise it returns DefaultPanelSize.
// Centralises the "fall back to the canonical size" logic so callers no longer
// have to inline `if n <= 0 { n = DefaultPanelSize }`.
func MustPanelSize(n int) int {
	if n <= 0 {
		return DefaultPanelSize
	}
	return n
}

// verdictEquivalentSet / verdictNotEquivalentSet memoise the alias lists
// from compiletime for O(1) lookup in normalizeVerdictToken.
var (
	verdictEquivalentSet    = stringSet(append([]string{compiletime.VerdictEquivalent}, compiletime.VerdictAliasEquivalent...))
	verdictNotEquivalentSet = stringSet(append([]string{compiletime.VerdictNotEquivalent, "NOT-EQUIVALENT"}, compiletime.VerdictAliasNotEquivalent...))
)

func stringSet(items []string) map[string]struct{} {
	m := make(map[string]struct{}, len(items))
	for _, s := range items {
		m[s] = struct{}{}
	}
	return m
}

// Judge prompts N independent judges to decide whether `source` and `target`
// are functionally equivalent. Each judge sees the same prompt but receives a
// distinct JudgeID ("J1".."Jn") so traces and disagreements are attributable.
//
// Aggregation: if the majority (rounded up) votes NOT_EQUIVALENT, Verdict.Agree
// is false. A tie on an odd panel still resolves to "not equivalent" only when
// the not-equivalent count strictly exceeds the equivalent count — exactly
// half on an odd panel is impossible, so this is equivalent to strict
// majority on the configured panel sizes.
//
// JudgeID and rationale strings are trimmed but otherwise preserved verbatim
// from the model output so external auditors can replay the panel.
func (vp *VerdictPanel) Judge(ctx context.Context, source, target string, n int) (Verdict, error) {
	if vp == nil || vp.Model == nil {
		return Verdict{}, fmt.Errorf("verdict panel: nil model")
	}
	n = MustPanelSize(n)

	opinions := make([]JudgeOpinion, 0, n)
	for i := 1; i <= n; i++ {
		judgeID := fmt.Sprintf("%s%d", compiletime.VerdictJudgeIDPrefix, i)
		opinion, err := vp.singleJudge(ctx, judgeID, source, target)
		if err != nil {
			// Record the failure as a NOT_EQUIVALENT vote with the error
			// message as rationale; the panel can still aggregate and the
			// caller gets visibility into the degraded path.
			logger.LogWarning("verdict panel: judge %s failed: %v", judgeID, err)
			opinions = append(opinions, JudgeOpinion{
				JudgeID:   judgeID,
				Verdict:   VerdictKindNotEquivalent,
				Rationale: fmt.Sprintf("judge error: %v", err),
			})
			continue
		}
		opinions = append(opinions, opinion)
	}

	return aggregateVerdict(opinions), nil
}

// singleJudge issues one equivalence-query prompt and parses the first-line
// verdict token out of the response.
func (vp *VerdictPanel) singleJudge(ctx context.Context, judgeID, source, target string) (JudgeOpinion, error) {
	resp, err := vp.Model.Generate(ctx, []*schema.Message{
		schema.SystemMessage(verdictJudgeSystemPrompt),
		schema.UserMessage(verdictJudgeUserPrompt(judgeID, source, target)),
	})
	if err != nil {
		return JudgeOpinion{}, fmt.Errorf("model generate: %w", err)
	}
	if resp == nil {
		return JudgeOpinion{}, fmt.Errorf("model returned nil response")
	}

	raw := strings.TrimSpace(resp.Content)
	verdict, rationale := parseJudgeVerdict(raw)
	return JudgeOpinion{
		JudgeID:   judgeID,
		Verdict:   verdict,
		Rationale: rationale,
	}, nil
}

// aggregateVerdict reduces the per-judge opinions into a single Verdict by
// strict majority. Equivalence votes are counted; if they reach the majority
// threshold (rounded up) the pair is treated as equivalent. Any judges that
// disagreed with the majority are listed in Disagreements.
func aggregateVerdict(opinions []JudgeOpinion) Verdict {
	if len(opinions) == 0 {
		return Verdict{PerJudge: opinions}
	}

	equivalentCount := 0
	for _, op := range opinions {
		if normalizeVerdictToken(op.Verdict) == VerdictKindEquivalent {
			equivalentCount++
		}
	}

	majority := majorityThreshold(len(opinions))
	agree := equivalentCount >= majority

	var disagreements []string
	if !agree {
		for _, op := range opinions {
			if normalizeVerdictToken(op.Verdict) != VerdictKindEquivalent {
				disagreements = append(disagreements, fmt.Sprintf("%s: %s", op.JudgeID, trimForLog(op.Rationale)))
			}
		}
	}

	logger.LogValidation("verdict panel: %d/%d equivalent votes, majority=%d agree=%t",
		equivalentCount, len(opinions), majority, agree)

	return Verdict{
		Agree:         agree,
		Disagreements: disagreements,
		PerJudge:      opinions,
	}
}

// majorityThreshold returns the smallest integer count of equivalent votes
// that must be reached for the panel to agree. For a panel of size k, the
// threshold is ceil(k/2)+1 — i.e., a strict majority of one over half.
func majorityThreshold(k int) int {
	return k/2 + 1
}

// normalizeVerdictToken upper-cases and trims a verdict label so judges that
// emit "Equivalent", "EQUIVALENT", "equiValent." or "EQUIVALENT\n..." all map
// to the same canonical bucket.
func normalizeVerdictToken(raw string) string {
	cleaned := strings.ToUpper(strings.TrimSpace(raw))
	cleaned = strings.Trim(cleaned, ".:;,-`'\"")
	if _, ok := verdictEquivalentSet[cleaned]; ok {
		return VerdictKindEquivalent
	}
	if _, ok := verdictNotEquivalentSet[cleaned]; ok {
		return VerdictKindNotEquivalent
	}
	return cleaned
}

// parseJudgeVerdict reads the first non-empty line of the model response,
// classifies it as equivalent or not, and returns the full text as the
// rationale so downstream auditors can review the reasoning.
func parseJudgeVerdict(raw string) (string, string) {
	rationale := strings.TrimSpace(raw)
	head := rationale
	if idx := strings.IndexAny(head, "\n\r"); idx >= 0 {
		head = head[:idx]
	}
	return normalizeVerdictToken(head), rationale
}

// trimForLog shortens a rationale to a single line and caps its length so the
// structured log stays readable when judges produce long paragraphs.
func trimForLog(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	const maxLen = 160
	if len(s) > maxLen {
		s = s[:maxLen] + "..."
	}
	return s
}

// verdictJudgeSystemPrompt is the role prompt sent to every judge. It is
// constant across judges so the only thing that varies between votes is the
// JudgeID echoed in the user prompt, keeping the panel genuinely independent
// in spirit (each judge sees the same role + the same data) and easy to
// reason about in audit logs.
const verdictJudgeSystemPrompt = `You are verdict_judge, an independent reviewer in a multi-agent code-translation panel. ` +
	`You decide whether two code fragments are functionally equivalent for the purposes of a translation pipeline. ` +
	`Reply with EXACTLY one of these tokens on the first non-empty line, in upper case, and nothing else on that line:

` +
	`  EQUIVALENT
` +
	`  NOT_EQUIVALENT

` +
	`Follow the verdict line with a short rationale (one or two sentences). ` +
	`Treat any difference in observable behavior — return value, thrown/returned error, side effect on global state, ` +
	`asynchronous ordering, panic vs. graceful error — as NOT_EQUIVALENT. ` +
	`Cosmetic differences (whitespace, comment text, identifier naming that preserves semantics) are EQUIVALENT. ` +
	`If you are uncertain, vote NOT_EQUIVALENT.`

// verdictJudgeUserPrompt formats the per-judge prompt. judgeID is echoed so
// judges can self-identify in their rationale if they wish; this is cosmetic
// and does not influence aggregation.
func verdictJudgeUserPrompt(judgeID, source, target string) string {
	return fmt.Sprintf(`Judge ID: %s

Compare the SOURCE and TARGET fragments below and decide whether they are functionally equivalent.

=== SOURCE ===
%s

=== TARGET ===
%s

Reply with one verdict line (EQUIVALENT or NOT_EQUIVALENT) followed by a short rationale.`,
		judgeID, source, target)
}
