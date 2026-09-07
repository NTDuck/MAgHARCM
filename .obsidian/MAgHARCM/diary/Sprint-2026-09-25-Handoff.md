---
title: Sprint 2026-09-25 Handoff
backlink: "[[3.0.0 Sprint 2026-09-25 Handoff]]"
tags: [sprint, handoff, [[2.0.0 MAgHARCM]], wave-14, swebench, bfcl, slm, vault-sync, paper-sync]
---

# [[3.0.0 Sprint 2026-09-25 Handoff]]

## Outcome

Wave-14 fired. Three new SLM-era paper notes persisted (`[[1.0.0 P-119]]` SWE-Rebench Badertdinov NeurIPS 2025, `[[1.0.0 P-120]]` SWE-smith Yang NeurIPS 2025 spotlight, `[[1.0.0 P-121]]` BFCL Patil ICML 2025). Wave-14 cross-links landed in five `Software-Archaeology-Lineage.md` rows (PRIM-22, PRIM-23, PRIM-27, PRIM-29, PRIM-31). Four `\cite{}` updates in `sec_method.tex` (L336, L471, L532, L590) plus three new bib entries in `refs.bib`. Methodology.md §9 changelog + §7 wave-13 trigger criterion updated. No code logic changes. `go build` / `go vet` / `go test ./...` all green.

## Foundation (closed)

- Read Sprint-2026-09-24 handoff + `Methodology.md` (entry point) for grounding.
- Audited current code state (no drift):
  - Zero `fmt.Print*` I/O in production Go code.
  - Zero hardcoded magic strings outside `internal/compiletime/`.
  - `ABCoderMCPProvider` (wrong casing): **zero** hits in any `.go` file.
  - `SelectMigrationStrategy` references: only the one historical-context comment in `tests/internal/agents/analyzer.go:13` (intentional, retained).
- Inventoried primitives INDEX: 31 unique PRIM-NN anchors (no orphans).
- Wave-14 trigger decision: **FIRE**. Three 2025 NeurIPS/ICML papers introduced unanchored mechanisms the SLM-era primitives (PRIM-22/23/27/29/31) needed to defend their substrate claims.

## Wave-14 Research (fired)

### [[1.0.0 P-119]] — Badertdinov et al. 2025 — SWE-Rebench

- NeurIPS 2025 Datasets and Benchmarks Track (featured). arXiv:2505.20411.
- Mechanism: continuously-evolving task harvest from GitHub with per-instance `created_at` vs. model `training_cutoff` provenance gating.
- Anchors: PRIM-22 (Search-phase contamination filter), PRIM-27 (plateau vs. memorisation substrate), PRIM-31 (provenance-driven retrieval exclusion).
- Headline finding: 5-15pp contamination-inflation deltas vs. SWE-bench Verified for several popular models.

### [[1.0.0 P-120]] — Yang et al. 2025 — SWE-smith

- NeurIPS 2025 Datasets and Benchmarks Track (spotlight). arXiv:2504.21798.
- Mechanism: environment-first procedural AST-rewrite synthesis of new tasks from existing GitHub repositories; PR-mirrored ground-truth fixes as retrieval training substrate.
- Anchors: PRIM-22 (bug-topology-aware Anchoring phase), PRIM-23 (adversarial test-synthesis training-corpus analogue), PRIM-27 (per-iteration training-time sweeps past contamination plateau), PRIM-31 (realistic-retrieval training substrate).

### [[1.0.0 P-121]] — Patil et al. 2025 — BFCL

- ICML 2025 (PMLR 267:48371-48392). arXiv: 2607.05775 [INFERENCE; verify against ICML proceedings PDF on first revisit].
- Mechanism: AST-based per-turn contract evaluation of tool calls; call-pattern templates (serial / parallel / multi-turn) for the function-calling-specific benchmark the open-weights SLM regime needs.
- Anchors: PRIM-22 (AST-based per-turn contract check), PRIM-29 (multi-turn stateful call pattern as tool-return-driven dynamic re-recruitment), PRIM-31 (call-pattern templates for retrieval-refinement training).
- Supersedes wave-10 `[[1.0.0 P-106]]` BFCL anchor as the canonical tool-calling benchmark reference.

### P-122 — deferred

Originally scoped as a fourth wave-14 anchor (a 2026 SLM-era paper I had under consideration during foundation). Decision: drop. Three papers already cover the SWE-bench + BFCL substrate gap; adding a fourth would dilute the trigger criterion without adding substrate-defence evidence. P-122 remains a future-wave candidate, not a wave-14 deliverable.

## Track-1 Vault Sync (closed)

### Methodology.md
- §0 frontmatter `last_updated` bumped 2026-09-23 → 2026-09-25.
- §9 changelog: appended `2026-09-25` entry noting wave-14 fire (P-119/P-120/P-121), three SLM-era substrate-defence anchors, no code logic changes.
- §7 trigger criterion: rewrote the wave-13 → wave-14 successor criterion to point at the 2026-09-25 decisions ("wave-N+1 fires when a 2025+ NeurIPS/ICML/ICLR paper introduces an unanchored mechanism that defends or refutes an existing SLM-era primitive's substrate claim").

### primitives/Primitives-Index.md
- Appended `## Sprint 2026-09-25 Vault Sync Audit` block: 3 papers persisted, 5 lineage rows cross-linked (PRIM-22, PRIM-23, PRIM-27, PRIM-29, PRIM-31), 0 compliance drift.

### Software-Archaeology-Lineage.md
- Line 148 (PRIM-22 row): added `[[1.0.0 P-119]]`, `[[1.0.0 P-120]]`, `[[1.0.0 P-121]]` to Primary Source cell; extended Theoretical Lineage cell with per-instance provenance gate, procedural AST mutations, and AST-based per-turn contract check mechanisms.
- Line 149 (PRIM-23 row): added `[[1.0.0 P-120]]` for SWE-smith AST-rewrite strategies as adversarial test-synthesis training-corpus analogue.
- Line 153 (PRIM-27 row): added `[[1.0.0 P-119]]` for contamination-aware substrate, `[[1.0.0 P-120]]` for environment-first synthetic-task injection as adaptive difficulty regulator.
- Line 155 (PRIM-29 row): added `[[1.0.0 P-121]]` for BFCL multi-turn stateful call pattern as tool-return-driven dynamic re-recruitment.
- Line 157 (PRIM-31 row): added `[[1.0.0 P-119]]` for provenance-driven retrieval exclusion, `[[1.0.0 P-120]]` for PR-mirrored ground-truth fixes as realistic-retrieval training substrate, `[[1.0.0 P-121]]` for call-pattern templates.

## Track-2 Paper Sync (closed)

### docs/.paper/refs.bib
Three new entries appended after `p118_jimenez_swebench_lite_2024`:
- `p119_badertdinov_swerebench_2025` (@misc, NeurIPS 2025 D&B featured, eprint 2505.20411)
- `p120_yang_swesmith_2025` (@misc, NeurIPS 2025 D&B spotlight, eprint 2504.21798)
- `p121_patil_bfcl_2025` (@inproceedings, ICML 2025 PMLR 267:48371-48392, eprint 2607.05775 [INFERENCE])

### docs/.paper/sec_method.tex
Four `\cite{}` cluster updates:
- **L336** (Qwen2.5-Coder SLM-scale cluster): added `p121_patil_bfcl_2025` immediately after `p106_bfcl_2025` (P-121 supersedes P-106 as canonical tool-calling benchmark reference).
- **L471** (CPG/iterative retrieval cluster): added `p120_yang_swesmith_2025` for procedural AST-rewrite training corpus.
- **L532** (solvability taxonomy cluster): added `p119_badertdinov_swerebench_2025` for contamination-aware evaluation lineage.
- **L590** (PRIM-29 recruitment cluster): added `p121_patil_bfcl_2025` for multi-turn stateful call pattern as dynamic re-recruitment.

## Track-3 Verification (closed)

- `go build ./...` — clean (no output, exit 0).
- `go vet ./...` — clean (no output, exit 0).
- `go test ./...` — all 8 test packages green (cmd/MAgHARCM, cmd/MAgHARCM-tui, internal/agents, internal/config, internal/languages, internal/logger, internal/runner, internal/tools).

## Track-4 Commit Plan (closed)

Four focused commits, ordered for safe rollback:

```
1. docs(research): persist wave-14 P-119/P-120/P-121 paper notes
   - 3 new files under .obsidian/MAgHARCM/research/papers/
2. docs(research): cross-link wave-14 anchors into Software-Archaeology-Lineage.md
   - 5 rows (PRIM-22, PRIM-23, PRIM-27, PRIM-29, PRIM-31)
3. docs(methodology): bump last_updated; append 2026-09-25 §9 changelog; rewrite §7 wave-N+1 trigger criterion; append INDEX vault-sync audit block
4. docs(paper): append wave-14 bib entries; update 4 sec_method.tex \cite{} clusters (L336, L471, L532, L590)
5. docs(diary): append sprint 2026-09-25 handoff (this file)
```

## Backlog / open items

- Possible follow-up: verify P-121 arXiv id `2607.05775` against the ICML 2025 proceedings PDF on first revisit; the P-121 paper header notes this is `[INFERENCE]`-marked per the methodology.
- Possible follow-up: SWE-Rebench V2 (arXiv:2602.23866 [INFERENCE]) is a forward reference in `[[1.0.0 P-119]]`; materialise as `[[1.0.0 P-119-v2]]` when the V2 paper note lands.
- Possible follow-up: extend the ADR-V-001 sweep into a pre-commit lint check that also rejects stray `(P-NN)` / `(PRIM-NN)` references in `Software-Archaeology-Lineage.md` figure annotations (currently out of scope per Sprint 2026-09-24 reasoning).
- Possible follow-up: the wave-13 trigger criterion at Methodology.md §7 should be rewritten to point at the wave-N+1 criterion (already done this sprint; verify on next sprint).

## Commits this sprint

(Pending: five commits above.)
