---
title: Human Intervention Required & Active Project Blockers
date: 2026-09-28
last_updated: 2026-09-28
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
| **BLK-01** | Codebase Build | `internal/memorystore/memorystore.go:375:24: undefined: compiletime.MaTTSDefaultBudget` | `go test ./...` fails in `internal/memorystore` | Add `const MaTTSDefaultBudget = 10` (or appropriate value) in `internal/compiletime/compiletime.go` |
| **BLK-02** | Benchmark Plateau | `Commons-Validator` (Java $	o$ Rust) plateaus at Iteration 12 (18/68 tests passing, 26.5%) | Target project fails functional parity | Provide human architectural guidance on deep Java class inheritance / regex translation |
| **BLK-03** | Literature Verification | Papers P-85, P-86, P-89 are marked `UNVERIFIED placeholder` | Literature citations lack verified venue URLs | Review and replace with confirmed peer-reviewed venue citations (e.g. Li et al. ASE 2024) |
| **BLK-04** | Runtime Environment | Ollama/vLLM local endpoint configuration for 4B-30B models | Offline execution fails without running daemon | Configure GPU VRAM allocation or provide remote fallback API key in `.env` |

---

## Detailed Issue Breakdowns

### BLK-01: MemoryStore MaTTS Compilation Break
- **Location**: `internal/memorystore/memorystore.go:375:24`
- **Error**: `undefined: compiletime.MaTTSDefaultBudget`
- **Context**: The `ApplyMaTTS` function implements Memory-augmented Test-Time Scaling, but references a constant `compiletime.MaTTSDefaultBudget` that was never declared in `internal/compiletime/compiletime.go`.
- **Recommendation**: Human engineer should decide the default budget iteration limit (e.g. 5 or 10) and export `MaTTSDefaultBudget` in `compiletime.go`.

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
