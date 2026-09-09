---
title: "P-158 Agentic Rubrics as Contextual Verifiers for SWE Agents"
backlink: "[[1.0.0 P-158]]"
aliases:
 - "1.0.0 P-158"
 - "P-158"
 - "P-158-Agentic-Rubrics-ACL-2026"
 - "Agentic-Rubrics-ACL-2026"
 - "raghavendra2026agenticrubrics"
tags: [paper, verification, rubric-checklist, execution-free, swe-agents, slm-era, "[[1.0.0 PRIM-7]]", wave-25]
date: 2026-09-09
last_updated: 2026-09-09 (iter-3, wave-25)
venue: ACL 2026 (Long Papers)
---

# [[1.0.0 P-158]] Agentic Rubrics

## TL;DR

Agentic Rubrics is an **execution-free patch-verification** substrate: an expert agent explores the repository, emits a context-grounded rubric checklist, and candidate patches are scored against the checklist without running tests. Anchors `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation — the VerdictPanel gains a verification signal that does not depend on the compile+test cascade when that cascade is too expensive for the SLM repair loop.

## Mechanism (Q2)

1. **Rubric generation** — an expert agent interacts with the repository to produce a structured, context-grounded checklist of criteria for what a correct patch must satisfy (root cause addressed, no regressions, tests still valid, etc.).
2. **Execution-free scoring** — candidate patches are scored against the checklist without requiring test execution, removing environment-setup overhead from the verification loop.
3. **Results** — 54.2% on SWE-Bench Verified with Qwen3-Coder-30B-A3B under parallel test-time scaling (30B-class = SLM-relevant), at least +3.5pp over the strongest baseline; 40.6% on Qwen3-32B. Rubric scores are consistent with ground-truth tests while flagging issues that tests do not capture.
4. **Ablations** — agentic context gathering is essential for producing codebase-specific, unambiguous criteria (the checklist is not generic).

## Anchoring (Q3)

| Primitive | Pre-wave-25 behaviour | Agentic Rubrics substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation | Verdict anchored on P-57 Speculative Decoding, P-78 EAGLE-3, P-92 Lightman PRM, P-125 T1, P-126 ARC-Decode, P-127 SLM-as-a-Judge, P-135 SPECS, P-136 CaTS, P-146 ContextPRM; every signal either executes code or uses a learned reward model | Agentic Rubrics adds an execution-free, self-generated-checklist signal: the VerdictPanel can gate patch acceptance when the test cascade is too expensive at SLM scale, complementing P-146 ContextPRM's learned PRM with a zero-training rubric |

## Hop-1 Citations

- P-146 ContextPRM (ICLR 2026, Wave-20) — workflow-aware process reward modeling. Agentic Rubrics replaces the learned PRM with a self-generated checklist; both score patch/intermediate quality without ground-truth execution.
- P-127 SLM-as-a-Judge (Wave-18) — SLM-scale verdict validation; the rubric gives a small judge model a structured contract to apply.

## Hop-2 Citations

- Self-consistency / majority-vote verification baselines (Wang 2022) — the execution-free family Agentic Rubrics outperforms on SWE-Bench Verified.

## MAgHARCM integration

- **YAML config key**: `agents.verdict_panel.rubric_checklist: true` (hypothetical gate; not wired this session).
- **Implementation file**: `internal/agents/verdict_panel.go` — future substrate; the checklist scorer would sit alongside the existing consensus voters.
- **Affected primitives**: `[[1.0.0 PRIM-7]]`.

## Caveats

- Evaluated on SWE-Bench Verified (Python repair); generalisation to C→Rust translation verdicts is untested — the checklist criteria would need re-grounding for translation-equivalence rather than issue resolution.
- Rubric generation costs one agentic repository exploration per criterion set; cost scales with codebase size.
- ACL 2026 admission follows the Wave-25 venue-list amendment (ACL Long Papers added to the §7 list); if a later sprint reverses that amendment, this anchor reverts to watchlist.

## Source

- Venue: ACL 2026 Long Papers, verified at aclanthology.org/2026.acl-long.697 (DOI 10.18653/v1/2026.acl-long.697, pages 15265–15290, July 2026).
- arXiv: 2601.04171.
- Authors' affiliation: Scale AI.

## BibTeX

```
@inproceedings{raghavendra2026agenticrubrics,
  title     = {Agentic Rubrics as Contextual Verifiers for {SWE} Agents},
  author    = {Raghavendra, Mohit and Gunjal, Anisha and Liu, Bing and He, Yunzhong},
  booktitle = {Proceedings of the 64th Annual Meeting of the Association for Computational Linguistics (Volume 1: Long Papers)},
  year      = {2026},
  pages     = {15265--15290},
  doi       = {10.18653/v1/2026.acl-long.697}
}
```
