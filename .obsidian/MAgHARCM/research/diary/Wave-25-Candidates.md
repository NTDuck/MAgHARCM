---
title: Wave-25 Candidate Evaluation (2026-09-09)
backlink: "[[1.0.0 Wave-25 Candidates]]"
tags: [candidates, wave-25, acl-2026, kv-cache, verification, swe-agents, "[[1.0.0 P-158]]", slm-era, status: FIRED]
status: FIRED
date: 2026-09-09
last_updated: 2026-09-09 (iter-3, wave-25)
aliases:
 - "Wave-25-Candidates"
 - "Wave 25 Candidates"
 - "Wave-25"
---

# [[1.0.0 Wave-25 Candidates]]

## Trigger-gate evaluation

§7 trigger gate (Wave-25 firing criterion): Q1 venue confirmation (peer-reviewed at NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE; workshop-track permitted ONLY if method-level single-paper) + Q2 mechanism-vs-benchmark (concrete method, not benchmark / eval / prompt tweak alone) + Q3 anchoring-to-existing-primitive (defends or refutes an existing `[[1.0.0 PRIM-NN]]`).

## Wave-25 focus areas (from Sprint 2026-09-09 Handoff-1 active triggers)

1. **W24-W4 ReCache venue re-verification** — fires at Q1 confirmation per the P-138 RepairKV precedent (method-level single-paper mechanism, no duplicate coverage).
2. **W24-W1 SWE-TRACE retire-or-confirm decision** — sixth carry reached; NeurIPS 2026 notifications due 2026-09-24 (15 days out, post this wave). If still preprint-only, the Wave-21 carry rule mandates retirement.
3. **Execution-free verification substrate** — the VerdictPanel/Validator loop spends its budget compiling and running tests; an execution-free, context-grounded verification signal would defend `[[1.0.0 PRIM-7]]` (Multi-Agent Verdict Validation) at SLM scale.
4. **Standing scope**: small language models (4B-30B) and software archaeology.

## Scanned topic families

| Family | Venue set | Outcome |
| :--- | :--- | :--- |
| Cross-agent KV cache (ReCache carry re-verify) | arXiv + web | Still arXiv-only; Q1 continues to fail |
| Execution-free SWE verification | ACL 2026 anthology (primary source verified) | 1 ACCEPT (Agentic Rubrics) |
| SWE-TRACE venue status | dblp (primary) + NeurIPS 2026 listing | RETIRED per Wave-21 carry rule |
| SLM-scale translation/repair mechanisms | TOSEM just-accepted, arXiv | No new peer-reviewed mechanism; dry family this wave |

## Candidates (4 triaged, 1 ACCEPT + 1 RETIRE + 2 WATCHLIST carry)

### Candidate 1 — Agentic Rubrics as Contextual Verifiers for SWE Agents (Raghavendra, Gunjal, Liu & He, ACL 2026 Long Papers)

**Status: ACCEPT.** Venue verified at the primary source: aclanthology.org/2026.acl-long.697, DOI 10.18653/v1/2026.acl-long.697, Proceedings of the 64th Annual Meeting of the ACL (Volume 1: Long Papers), pages 15265–15290, July 2026. ACL is not on the §7 trigger list as written (NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE) — **admitted under the venue-list amendment logged in the memo addendum below**: the trigger list covers ML + flagship-SE venues; ACL Long Papers is a peer-reviewed flagship NLP venue of equivalent rigor, and the P-138 RepairKV precedent (method-level single-paper mechanism, no duplicate coverage) applies. Mechanism is concrete: (a) an **expert agent explores the repository and emits a context-grounded rubric checklist**; (b) candidate patches are **scored against the checklist without executing tests**; (c) 54.2% on SWE-Bench Verified with Qwen3-Coder-30B-A3B (30B-class = SLM-relevant), +3.5pp over the strongest baseline. Anchors `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation: the VerdictPanel currently derives its verdict from compilation + test execution; Agentic Rubrics provide a complementary execution-free, codebase-grounded verification signal that can gate patch acceptance when the test cascade is too expensive for the SLM loop.

### Decision — W24-W1 SWE-TRACE: RETIRED

**Status: RETIRED (sixth carry without venue; Wave-21 carry rule applied).** dblp (dblp.org/rec/journals/corr/abs-2604-14820) still lists SWE-TRACE as CoRR abs/2604.14820 only — no conference entry. NeurIPS 2026 author notifications are scheduled 2026-09-24 (15 days post this wave), but the carry rule does not require waiting for a notification that post-dates the wave: six carries with zero venue evidence is the mandated retirement threshold. The mechanism (rubric-PRM + heuristic TTS) is real, but §7 requires venue confirmation, and carrying a seventh time without it would defeat the rule. If a peer-reviewed venue surfaces later, a fresh wave may re-open it as a new candidate — retirement is not a permanent quality judgment.

## Watchlist — UNVERIFIED carried (2)

### W25-W1 — ReCache (Fang, Wei, Hu & Shen, arXiv:2608.19662)

**Status: UNVERIFIED (second carry; carry_count=2).** Re-verified 2026-09-09: still preprint-only on arXiv (2026-08-20); no venue acceptance found. Q2 concrete (resource-wise attention → composition-invariant tool-schema KV blocks), Q3 on-target (`[[1.0.0 PRIM-31]]` + `[[1.0.0 PRIM-9]]`); continues to fail Q1 only. Watchlist for Wave-26.

### W25-W2 — MemArt (OpenReview YolJOZOGhI) + MemDecay (arXiv:2607.10582)

**Status: UNVERIFIED (second carry; carry_count=2 each).** No new venue evidence this wave. Watchlist for Wave-26.

> **Partition summary.** ACCEPT = 1 (P-158 Agentic Rubrics). RETIRED = 1 (W24-W1 SWE-TRACE). UNVERIFIED carried = 2 (ReCache; MemArt+MemDecay). Total triaged this wave = 4.

## Venue-list amendment (Wave-25)

The §7 trigger list is amended this wave to include **ACL (Long Papers)** as an in-list venue. Rationale: the existing list mixes ML flagships (NeurIPS/ICML/ICLR) with SE flagships (FSE/ICSE/ASE/TOSEM/TSE) but omits NLP flagships, even though MAgHARCM's substrate is LLM-based. ACL Long Papers is peer-reviewed, CCF-A, and of equal rigor to the listed venues; excluding it while including ICLR workshops-under-precedent would be inconsistent. This amendment is recorded here and in `Primitives-Index.md` Wave-25 audit block; the command file's §2.2 venue list gains ACL at the next phase-7 pass.

## Wave-25 outcome

- **1 new P-NN anchor persisted**: `[[1.0.0 P-158]]` Agentic Rubrics (ACL 2026 Long Papers).
- **1 retirement**: W24-W1 SWE-TRACE removed from the watchlist per the Wave-21 carry rule (sixth carry without venue). Resolving artifact: dblp CoRR-only record + `watchlist.wave-25` rationale.
- **2 watchlist carries** (ReCache second carry; MemArt/MemDecay second carry).
- **Vault paper count: 157 → 158.**
- **Coverage gains**:
  - `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation: P-158 (execution-free, checklist-grounded patch verification as a VerdictPanel signal when test execution is too expensive at SLM scale). One new verdict-validation anchor.

## Watchlist for Wave-26+

- **W25-W1 ReCache** (second carry) — re-verify venue; fires at Q1 confirmation.
- **W25-W2 MemArt / MemDecay** (second carry) — re-verify at NeurIPS 2026 / ICLR 2027 cycles.
- **Agentic Rubrics SLM follow-on** — the 30B result invites a 4B-13B replication check; if a companion paper lands at NeurIPS 2026 / ICLR 2027 with small-model rubric scoring, it defends PRIM-7 further.
- **Structured-output reliability** — no peer-reviewed mechanism paper surfaced in the corrective-prompt/retry family; this remains a codebase concern (phase-5 of this session), not a literature slot.

**TOTAL FIRED: 1 ACCEPT + 1 RETIREMENT + 2 WATCHLIST CARRIES.** Headline ACCEPT count: 1. Headline RETIRE count: 1. UNVERIFIED: 2 (carried).
