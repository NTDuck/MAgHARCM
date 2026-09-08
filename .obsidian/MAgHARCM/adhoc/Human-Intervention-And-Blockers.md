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

> **Executive Overview**: Central triage dashboard tracking active technical blockers, compilation plateaus, and resolved engineering issues requiring or having received human intervention. For definitions of system components and acronyms, consult the [[Glossary|Domain Acronyms & Terminology Glossary]].

---

## 1. Active & Informational Blockers

These items represent current project constraints or decision points requiring human guidance or infrastructure setup.

| Issue ID | Category | Status | Summary Description | Required Action |
| :--- | :--- | :--- | :--- | :--- |
| **BLK-02** | Benchmark Plateau | **ACTIVE** | `Commons-Validator` (Java to Rust) plateaus at Iteration 12 (18/68 tests passing, 26.5%). | Provide architectural guidance on Java inheritance hierarchies and regex crate mapping in `configs/commons-validator.yml`. |
| **BLK-03** | Literature Verification | **INFORMATIONAL** | Papers P-85, P-86, P-89 carry `UNVERIFIED placeholder` tags in older research waves. | Confirm peer-reviewed venue publication or replace with verified counterparts (e.g. Li et al. ASE 2024). |
| **BLK-04** | Runtime Environment | **INFORMATIONAL** | Local GPU inference endpoint (Ollama/vLLM daemon) is offline in standard CI sandbox. | Configure workstation GPU allocation or provide fallback remote API credentials in `.env` to enable stochastic re-runs. |
| **W23-W1** | Research Watchlist | **WATCHLIST** | SWE-TRACE (`arXiv:2604.14820`) is on its fifth carry awaiting NeurIPS 2026 notifications (2026-09-24). | Re-verify on 2026-09-24. If still lacking a peer-reviewed venue, apply the Wave-21 carry rule to retire the paper. |

---

### Detailed Active Issue Breakdowns

#### BLK-02: Commons-Validator Translation Plateau
- **Location**: `docs/.paper/sec_eval.tex`, `Table 2`
- **Issue**: Translating Apache Commons-Validator (121 files, 28,110 LoC) hits a plateau where 18/68 tests pass and compilation fails due to complex Java inheritance hierarchies and reflection patterns.
- **Root Cause**: Small language models struggle with inter-file inheritance cycles and dynamic regex validation mappings without explicit trait bounds.
- **Action Plan**: Define a custom feature mapping profile in `configs/commons-validator.yml` mapping Java regex primitives to `regex` crate equivalents.

#### BLK-03: Unverified Literature Placeholders
- **Files**:
  - `research/papers/P-85-Yue-Function-Calling-2025-Unverified.md`
  - `research/papers/P-86-Xu-Modernization-Survey-2024-Unverified.md`
  - `research/papers/P-89-Phan-Baseline-ICSE-NIER-2024-Unverified.md`
- **Issue**: These notes represent provisional anchors identified in earlier research waves whose publication venues could not be confirmed in available open databases.
- **Action Plan**: Review and verify peer-reviewed venue publication or replace with confirmed counterparts.

#### BLK-04: Runtime Environment (Ollama/vLLM Endpoint)
- **Issue**: Local SLM inference (4B–30B parameters) requires an active Ollama or vLLM daemon. In sandboxed CI/test environments, the daemon is unreachable.
- **Impact**: Empirical benchmark evaluation remains pegged to verified baseline trials ($K=3$); new substrate gates remain disabled by default.
- **Action Plan**: Launch the Ollama daemon on an available GPU workstation (RTX 4090 / Apple Silicon) or configure remote API fallback keys.

---

## 2. Resolved Technical Blockers & Triage History

The following technical blockers and environmental issues have been fully diagnosed and resolved by engineering actions.

| Issue ID | Resolution Date | Category | Root Cause & Resolution |
| :--- | :--- | :--- | :--- |
| **BLK-01** | 2026-09-07 | Build Break | `undefined: compiletime.MaTTSDefaultBudget` in `memorystore.go`. Resolved by exporting `const MaTTSDefaultBudget = 10` in `internal/compiletime/compiletime.go`. |
| **BLK-05** | 2026-09-07 | Test Regression | Latent test failure in `ApplyMaTTS` loop accumulation and error wrapping (`%v` vs `%w`). Fixed loop triple appending and wrapped both errors via `%w`. |
| **BLK-06** | 2026-09-07 (iter-3) | Metadata Reset | Reset frontmatter dates to authoritative current date (`2026-09-07` / `2026-09-08`) per the standing dating convention. |
| **BLK-08** | 2026-09-07 / 08 | Triage Registry | Established canonical `reject_registry` and `watchlist` blocks in `Research-Database.json` across Waves 19–23 to prevent literature amnesia. |
| **P-153 Slot** | 2026-09-08 (iter-1) | Research Slot | Closed the strict program-comprehension-mechanism residual slot (open since Wave-17) via `[[1.0.0 P-153]]` CoReX context-aware slicing. |
| **Sandbox Git** | 2026-09-08 (iter-1) | Environment | Stale command syntax diagnostics resolved during Wave-23 close-out; verified clean commit `792ed5f` landing all 19 vault/paper sync files. |

---

### Detailed Resolved Issue Diagnostics

#### BLK-01: MemoryStore MaTTS Compilation Break
- **Location**: `internal/memorystore/memorystore.go:375:24`
- **Error**: `undefined: compiletime.MaTTSDefaultBudget`
- **Context**: The `ApplyMaTTS` function implements Memory-augmented Test-Time Scaling, but referenced a constant `compiletime.MaTTSDefaultBudget` that was never declared in `internal/compiletime/compiletime.go`.
- **Resolution**: Exported `const MaTTSDefaultBudget = 10` in `internal/compiletime/compiletime.go` adjacent to the memory-substrate constants (`MaxMemoryTriples`, `MemoryTripleRewardEMA`). A human engineer may override by passing an explicit positive budget at the callsite.
- **Status**: **RESOLVED 2026-09-07**.

#### BLK-05: ApplyMaTTS Test Regression (Latent Bug Resolved)
- **Location**: `internal/memorystore/memorystore.go:386-410`; tests at `internal/memorystore/memorystore_test.go:201-258`
- **Symptom**: `TestApplyMaTTSBudgetCeiling` reported `len(used)=6, want 3`; `TestApplyMaTTSErrorPath` reported `len(used)=0, want 1`. `errors.Is(err, wantErr)` failed on the budget-exhausted path.
- **Root Cause**:
  1. The loop body performed `usedStrategies = append(usedStrategies, prior...)` with full store contents instead of the single distilled triple.
  2. The budget-exhausted return formatted `lastErr` with `%v` rather than wrapping it with `%w`.
- **Resolution**:
  1. Appended the distilled `triple` once per successful iteration.
  2. Wrapped both errors: `fmt.Errorf("%w: last error: %w", ErrBudgetExhausted, lastErr)`.
  3. All `internal/memorystore` unit tests now pass cleanly (`go test -count=1 ./internal/memorystore/...` exit 0).
- **Status**: **RESOLVED 2026-09-07**.

#### Sandbox Blocker Resolution (2026-09-08 iter-1)
- **Context**: Previous sprint handoffs noted persistent sandbox git execution errors (`pi-natives:command: syntax error`).
- **Resolution**: Verified during Wave-23 close-out that `git add -A` and `git commit` succeeded, landing commit `792ed5f` (19 vault/paper sync files). Subsequent commits continue to succeed normally.
- **Status**: **RESOLVED 2026-09-08**.
