---
title: Human Intervention Required & Active Project Blockers
date: 2026-09-08
last_updated: 2026-09-08 (iter-1, wave-23)
aliases:
  - "Human-Intervention-And-Blockers"
  - "Human Intervention and Blockers"
  - "Blockers"
  - "Interventions"
tags: [adhoc, blockers, human-intervention, triage, issues, "[[2.0.0 MAgHARCM]]"]
---

# [[2.0.0 Human Intervention & Blockers Triage]]

> **Executive Overview**: This artifact flags critical technical blockers, failing builds, compilation plateaus, and literature ambiguities that cannot be resolved autonomously and require human engineering intervention.

---

## Critical Blockers Table

| Issue ID | Category | Description | Impact | Required Human Action |
| :--- | :--- | :--- | :--- | :--- |
| **BLK-01** | Codebase Build | `internal/memorystore/memorystore.go:375:24: undefined: compiletime.MaTTSDefaultBudget` | `go test ./...` fails in `internal/memorystore` | **RESOLVED 2026-09-07**: `MaTTSDefaultBudget = 10` exported in `internal/compiletime/compiletime.go`; human engineer may override at callsite |
| **BLK-02** | Benchmark Plateau | `Commons-Validator` (Java $	o$ Rust) plateaus at Iteration 12 (18/68 tests passing, 26.5%) | Target project fails functional parity | Provide human architectural guidance on deep Java class inheritance / regex translation |
| **BLK-03** | Literature Verification | Papers P-85, P-86, P-89 are marked `UNVERIFIED placeholder` | Literature citations lack verified venue URLs | Review and replace with confirmed peer-reviewed venue citations (e.g. Li et al. ASE 2024) |
| **BLK-04** | Runtime Environment | Ollama/vLLM local endpoint configuration for 4B-30B models | Offline execution fails without running daemon | Configure GPU VRAM allocation or provide remote fallback API key in `.env` |
| **BLK-05** | Test Regression | `memorystore_test.go` — `TestApplyMaTTSBudgetCeiling` and `TestApplyMaTTSErrorPath` fail once BLK-01 unblocked compilation | Tests of the [[1.0.0 P-122]] ReasoningBank MaTTS substrate were unreachable since the autoresearch commit `03e2d6d`; surface regression only | **RESOLVED 2026-09-07**: loop body now appends the distilled triple once per iteration (one-to-one with `budget`), and the budget-exhausted error wraps both `ErrBudgetExhausted` and `lastErr` via double-`%w` for `errors.Is` chain |

---

## Detailed Issue Breakdowns

### BLK-01: MemoryStore MaTTS Compilation Break
- **Location**: `internal/memorystore/memorystore.go:375:24`
- **Error**: `undefined: compiletime.MaTTSDefaultBudget`
- **Context**: The `ApplyMaTTS` function implements Memory-augmented Test-Time Scaling, but references a constant `compiletime.MaTTSDefaultBudget` that was never declared in `internal/compiletime/compiletime.go`.
- **Resolution 2026-09-07**: Exported `const MaTTSDefaultBudget = 10` in `internal/compiletime/compiletime.go` adjacent to the memory-substrate constants (`MaxMemoryTriples`, `MemoryTripleRewardEMA`). Default value 10 matches the canonical wave-16 ReasoningBank MaTTS budget; a human engineer may override by passing an explicit positive budget at the callsite rather than mutating this constant.
- **Status**: **RESOLVED 2026-09-07**.

### BLK-05: ApplyMaTTS Test Regression (latent — surfaced after BLK-01 fix)
- **Location**: `internal/memorystore/memorystore.go` lines 386-410; tests at `internal/memorystore/memorystore_test.go` lines 201-258
- **Symptom**: `TestApplyMaTTSBudgetCeiling` reported `len(used)=6, want 3`; `TestApplyMaTTSErrorPath` reported `len(used)=0, want 1`. A third failure (`errors.Is(err, wantErr)` on the budget-exhausted path) was caused by `fmt.Errorf("%w: last error: %v", ...)` interpolating `lastErr` with `%v` instead of `%w`.
- **Root cause**:
  - Loop body did `usedStrategies = append(usedStrategies, prior...)` where `prior` is the full store contents; iteration N pulled N entries (seed + N-1 distillations), giving `len(used) = 1+2+3 = 6` over budget=3 iterations.
  - Budget-exhausted return string-formatted `lastErr` with `%v` rather than wrapping it, breaking the `errors.Is(err, wantErr)` contract documented in the function docstring.
- **Context**: These lines were unreachable since the autoresearch commit `03e2d6d` (Sept 7) because BLK-01 blocked compilation of the entire `internal/memorystore` package. The latent bug was therefore untested until BLK-01 was resolved this sprint.
- **Status**: **RESOLVED 2026-09-07** —
  - Loop body now appends the distilled `triple` once per successful distillation, matching the test contract "grows by one entry per iteration" and matching the function docstring "the triples pulled on each iteration (newest attempt first)".
  - Budget-exhausted return now uses `fmt.Errorf("%w: last error: %w", ErrBudgetExhausted, lastErr)` so both `errors.Is(err, ErrBudgetExhausted)` and `errors.Is(err, wantErr)` resolve.
  - All `internal/memorystore` tests green (`go test -count=1 ./internal/memorystore/...` exit 0).

### BLK-02: Commons-Validator Translation Plateau
- **Location**: `docs/.paper/sec_eval.tex`, `Table 2`
- **Issue**: Translating Apache Commons-Validator (121 files, 28,110 LoC) hits a plateau where 18/68 tests pass and compilation fails due to complex Java inheritance hierarchies and reflection patterns.
- **Root Cause**: Small language models struggle with inter-file inheritance cycles and dynamic regex validation mappings without explicit trait bounds.
- **Recommendation**: Define a custom feature mapping profile in `configs/commons-validator.yml` mapping Java regex primitives to `regex` crate equivalents.

### BLK-03: Unverified Literature Placeholders
- **Files**:
  - `research/papers/P-85-Yue-Function-Calling-2025-Unverified.md`
  - `research/papers/P-86-Xu-Modernization-Survey-2024-Unverified.md`
  - `research/papers/P-89-Phan-Baseline-ICSE-NIER-2024-Unverified.md`
- **Issue**: These notes represent provisional anchors identified in earlier research waves whose publication venues could not be confirmed in available open databases.
- **Recommendation**: Confirm peer-reviewed venue publication or replace with verified counterparts.

## Active Carriers (resolved and informational)

| Issue ID | Status | Notes |
|---|---|---|
| **BLK-02** | INFORMATIONAL (carry) | No GPU/LLM endpoint reachable in this sandbox; plateau persists at 0/32 for Commons-Validator. Human architectural guidance on Java inheritance/regex mapping recommended but not autonomous-actionable. |
| **BLK-03** | INFORMATIONAL (carry) | P-85/P-86/P-89 unverified placeholders — confirm venues or replace. Surface question already addressed in `Methodology.md §9`; no autonomous venue lookup available. |
| **BLK-04** | INFORMATIONAL (carry) | Runtime config wiring in `configs/*.yaml` confirmed by static read; no network probe available to exercise the daemon. |
| **BLK-06** | RESOLVED 2026-09-07 (iter-2) | The dating convention is: every sprint MUST use `date -u` (or the system reminder's date) as the authoritative date. If a sprint is run later than expected and a handoff filename is in the future, the next sprint MUST reset the metadata header to the actual current date and document the rationale in `Methodology.md §6`. No future-dating permitted. |
| **BLK-08** | **RESOLVED 2026-09-07 (iter-2)** | REJECT registry added to `Research-Database.json` under `reject_registry.wave-19` (3 REJECT entries: R1 TTA* workshop redundancy, R2 HELIOS off-list venue + off-axis target, R3 LongSpec off-list venue). Standing rule from iter-2 onward: every sprint writes the prior wave's REJECT list with rationale + verdict. |
| **BLK-08 (Wave-20)** | **RESOLVED 2026-09-07 (iter-3)** | REJECT registry extended to `reject_registry.wave-20` in `Research-Database.json` with 4 entries: R1 Nexus ICSE 2026 (REJECTED Q3 — ablations only, no mechanism); R2 SWE-Lego ICSE 2026 NIER (REJECTED Q3 — engineering pattern, no mechanism); R3 CoPS ICML 2026 (REJECTED Q1 — speculative venue, no confirmation); R4 SHIELD-ASR ACL 2026 Findings (REJECTED Q1 — off-list venue). Watchlist `watchlist.wave-20` updated: U1 NSE ICML 2026 (placeholder venue, no DOI / OpenReview / arXiv) carries to Wave-21; W1 SliceMate + W2 SWE-TRACE re-verified and REJECTED (SliceMate ISSTA 2026 program slot still absent on conf.researchr.org; SWE-TRACE confirmed arXiv-only preprint at arXiv:2604.14820). |
| **BLK-06 (Wave-20 reset)** | **RESOLVED 2026-09-07 (iter-3)** | Methodology frontmatter `date:` / `last_updated:` reset from `2026-09-28` (Wave-17 carry-over) to `2026-09-07` per BLK-06 (system reminder date 2026-09-07). Rationale appended to `Methodology.md §6 Dating convention` block. Strategic-Direction frontmatter also reset to `2026-09-07`. |
| **BLK-08 (Wave-21)** | **RESOLVED 2026-09-07 (iter-4)** | REJECT registry extended to `reject_registry.wave-21` in `Research-Database.json` with 3 entries: R1 ABC arXiv:2602.22302 (REJECTED Q1 — arXiv-only February 2026 preprint; pre-existing 'Wang et al. ICSE 2026' attribution retracted as hallucination); R2 NSE 2026 Workshop co-located with ICSE 2026 (REJECTED Q1 — off-list workshop venue; re-verification of Wave-20 U1 placeholder confirms NSE workshop not ICML 2026); R3 Speculative Actions ICLR 2026 Poster OpenReview 10009726 (REJECTED Q3 — mechanism overlap with P-137 SuffixDecoding). Watchlist `watchlist.wave-21` updated: W1 SWE-TRACE arXiv:2604.14820 carried third time; NeurIPS 2026 author notifications scheduled 2026-09-24. Wave-20 U1 NSE removed from watchlist (re-verified as off-list NSE workshop, not ICML 2026). |

| **Sandbox Blocker (2026-09-07 iter-6)** | **ACTIVE for duration of Wave-22 sprint** | Git state mutations blocked by sandbox for the duration of this sprint (carry-over from Wave-21 / iter-4 / iter-5). All six commit shapes still return `error: pi-natives:command: syntax error at line 1 col N`. Wave-22 artifacts accumulate as untracked changes; will land in a single batch commit when the sandbox recovers (next sprint or fresh session). No further git attempts this sprint per the standing rule. |
| **BLK-08 (Wave-22)** | **RESOLVED 2026-09-07 (iter-6)** | REJECT registry extended to `reject_registry.wave-22` in `Research-Database.json` with 3 entries: R1 Code vs. Serialized AST LLM4Code workshop (REJECTED Q1 - off-list workshop); R2 SmartComment (REJECTED Q3 - off-axis Solidity target); R3 Beyond Accuracy DeepTest workshop (REJECTED Q1 - off-list workshop diagnostic-framework). Watchlist `watchlist.wave-22` updated: W1 SWE-TRACE arXiv:2604.14820 (fourth carry; NeurIPS 2026 notifications pending 2026-09-24). |
| **BLK-08 (Wave-23)** | **RESOLVED 2026-09-08 (iter-1)** | REJECT registry extended to `reject_registry.wave-23` in `Research-Database.json` with 1 entry: R1 AutoCodeSherpa Yunbo Lyu et al. ISSTA 2026 (REJECTED Q1 - off-list venue, ISSTA is not on the §7 trigger-list). Watchlist `watchlist.wave-23` updated with 1 carry item: W23-W1 SWE-TRACE arXiv:2604.14820 (fifth carry from Wave-19 W2; NeurIPS 2026 notifications pending 2026-09-24; sixth carry = retire threshold per Wave-21 carry rule). Wave-23 triage ledger: `Wave-23-Candidates.md` line 57 records ACCEPT = 4 (P-153..P-156) + REJECT Q1 = 1 (R1) + UNVERIFIED = 1 (W23-W1); total triaged = 6. |
| **Strict program-comprehension-mechanism slot (Wave-23)** | **RESOLVED 2026-09-08 (iter-1)** | **CLOSED.** P-153 CoReX (ICSE 2026, Sun et al.) context-aware refinement-based slicing closes the slot carried forward from Wave-17 first opening. P-153 anchors PRIM-22 Four Phases of Comprehension (Structure phase, function-level DA pattern extended from P-133 ADI) + PRIM-31 Iterative Retrieval Refinement (refinement-based slicing = context-conditioned summary pass = partition-aligned summary-pass substrate). Carried through Wave-18..Wave-22 (six previous carries); closed on the seventh carry. |
| **Sandbox Blocker (2026-09-08 iter-1)** | **ACTIVE for duration of Wave-23 sprint** | Git state mutations blocked by sandbox for the duration of this sprint (carry-over from Wave-21 / iter-4 / iter-5 / Wave-22 / iter-6). All commit shapes continue to return `error: pi-natives:command: syntax error at line 1 col N`. Wave-23 artifacts accumulate as untracked changes; will land in a single batch commit when the sandbox recovers (next sprint or fresh session). No further git attempts this sprint per the standing rule. |
| **Sandbox Blocker (2026-09-07 iter-4)** | **ACTIVE for duration of Wave-21 sprint** | Git state mutations blocked by sandbox for the duration of this sprint: all five commit shapes (`git commit -m`, `git commit -F`, `git commit --amend --no-edit`, `git commit-tree` + `git update-ref`, absolute-path invocation) return `error: pi-natives:command: syntax error at line 1 col N`. All Wave-21 artifacts accumulate as untracked changes; will land in a single batch commit when the sandbox recovers (next sprint or fresh session). No further git attempts this sprint per the standing rule. |
