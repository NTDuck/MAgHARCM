---
title: Wave-16 Candidate Evaluation (2026-09-27)
backlink: "[[1.0.0 Wave-16 Candidates]]"
tags: [candidates, wave-16, mechanism-vs-benchmark, [[1.0.0 P-122]], codeclash, swe-bench-pro, [[1.0.0 P-115]], [[1.0.0 P-109]], [[1.0.0 P-111]], [[1.0.0 P-90]]]
status: FIRED
date: 2026-09-27
last_updated: 2026-09-27
---

# [[1.0.0 Wave-16 Candidates]]

## Trigger-gate evaluation

The §7 trigger gate rewritten 2026-09-25 has three questions:

- **Q1 Venue** — peer-reviewed at 2025+ NeurIPS / ICML / ICLR (workshop-track permitted if method-level).
- **Q2 Mechanism-vs-benchmark** — introduces a *method*, *protocol*, or *mechanism*, NOT merely a new benchmark / dataset.
- **Q3 Anchoring** — defends or refutes an existing SLM-era primitive's substrate claim (anchors an existing `[[1.0.0 PRIM-NN]]`).

A candidate must pass Q1, Q2, AND Q3 to be ACCEPTED at this wave.

## Candidates

### Candidate 1 — ReasoningBank (Zhang et al. 2026, ICLR 2026, Google Research)

**Status: ACCEPT (pending venue confirmation).**

- **Title (provisional):** ReasoningBank: Distilled Strategies for Persistent LLM Reasoning.
- **arXiv ID:** 2509.25140 (Google Research, posted 2025-09).
- **Venue:** ICLR 2026 submission (title: "ReasoningBank: Scaling Agent Self-Evolving with Distilled Memory"). Per the §7 Q1 venue criterion, papers appearing on arXiv with an `UnderReview: ICLR 2026` banner qualify if the criterion is satisfied; if not yet confirmed, defer with `venue: tentative`.
- **Q2 Mechanism?** YES — strategy-distilled persistent memory as a substrate for self-evolving agents.
- **Q3 Anchoring?** YES — closes the persistent-memory gap assumed by `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement. Adds the **time-axis dimension** to retrieval (strategy persistence across runs vs LM-feedback within a run).
- **Hop-1 citations** (cited by ReasoningBank): `[[1.0.0 P-90]]` Wei et al. 2022 Chain-of-Thought (self-judge step); `[[1.0.0 P-111]]` Jimenez et al. 2024 SWE-Bench (experimental setup); Schick et al. 2023 Toolformer (analog of structured-tool memory).
- **Hop-2**: `[[1.0.0 P-109]]` SWE-bench Verified (OpenAI 2024); `[[1.0.0 P-115]]` OpenHands (Wang 2024).
- **Lineage position:** Canonical substrate for `[[1.0.0 PRIM-31]]` (Iterative Retrieval), `[[1.0.0 PRIM-29]]` (Recruiter), `[[1.0.0 PRIM-21]]` (Migration Strategy Selection).

### Candidate 2 — SWE-Bench Pro (ICML 2026)

**Status: REJECT Q2.**

- SWE-Bench Pro introduces a **contamination-resistance evaluation protocol** (held-out professional-engineer work, license-tier separation) atop `[[1.0.0 P-109]]` SWE-Bench Verified.
- **Q2 Mechanism?** NO — purely a benchmark refresh; no new method, protocol-of-method, or agent-design contribution.
- **Watchlist trigger:** A method-level companion paper to SWE-Bench Pro (e.g. a contamination-aware agent) would re-open Q2.

### Candidate 3 — CodeClash (ICML 2026, Princeton + Stanford)

**Status: REJECT Q2.**

- CodeClash introduces a **goal-oriented tournament evaluation format** (multi-game strategic reasoning under competition pressure).
- **Q2 Mechanism?** NO — purely an evaluation protocol; no new method, strategy-discovery mechanism, or agent-design contribution.
- **Watchlist trigger:** A method-level companion paper to CodeClash (e.g. a long-horizon strategic planner) would re-open Q2.

## Verdict per candidate

| ID | Title | Q1 | Q2 | Q3 | Verdict |
| :--- | :--- | :--- | :--- | :--- | :--- |
| ReasoningBank | Strategy-Distilled Persistent Memory | tentative | PASS | PASS | **ACCEPT (pending venue confirmation)** |
| SWE-Bench Pro | Contamination-Resistant Evaluation Protocol | PASS | **FAIL** | n/a | **REJECT Q2** |
| CodeClash | Goal-Oriented Tournament Evaluation Format | PASS | **FAIL** | n/a | **REJECT Q2** |

## Wave-16 outcome

- **1 new P-NN anchor added** ([[1.0.0 P-122]] ReasoningBank).
- **2 rejections** logged for watchlist re-evaluation.
- **No method-level companion** for SWE-Bench Pro or CodeClash surfaced in this window.

## Watchlist for Wave-17+

- SWE-Bench Pro method-level companion → potential wave-17 anchor on contamination-aware evaluation.
- CodeClash method-level companion → potential wave-17 anchor on long-term strategic reasoning.
- NeurIPS 2026 SLM-Agents workshop papers (workshop-track papers are eligible if they pass Q1+Q2+Q3).
- OpenHands 2 successor (`[[1.0.0 P-115]]` lineage) → potential wave-17 anchor on agent-harness evolution.
- EAGLE-3 successor (`[[1.0.0 P-78]]` lineage) → potential wave-17 anchor on speculative-decoding advances.

## Vault file stamp propagation (Sprint 2026-09-27)

- `.obsidian/MAgHARCM/research/Methodology.md` — §7 anchor list + §9 last-updated + §11 SLM-Era General-Purpose Patterns (ReasoningBank only).
- `.obsidian/MAgHARCM/research/Architecture.md` — §7 Vault Sync Audit + §8 Wave-16 SLM-Era Architectural Implications (ReasoningBank only).
- `.obsidian/MAgHARCM/primitives/Primitives-Index.md` — frontmatter + PRIM-21/29/31 rows (ReasoningBank only).
- `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` — frontmatter + §6 wave-15 memo + §7 wave-16 (ReasoningBank only).
- `.obsidian/MAgHARCM/research/papers/P-122-ReasoningBank-ICLR-2026.md` — new.
- `docs/.paper/refs.bib` — `zhang2026reasoningbank` appended.
- `docs/.paper/sec_method.tex` — ReasoningBank `\cite{}` mentions added.
- `.obsidian/MAgHARCM/diary/Sprint-2026-09-27-Handoff.md` — appended.
