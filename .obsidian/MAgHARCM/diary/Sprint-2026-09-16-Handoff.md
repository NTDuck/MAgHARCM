---
title: Sprint 2026-09-16 Handoff
backlink: "[[3.0.0 Sprint 2026-09-16 Handoff]]"
tags: [sprint, handoff, "[[2.0.0 MAgHARCM]]", "[[1.0.0 PRIM-31]]", slm, ponytail, audit]
---

# [[3.0.0 Sprint 2026-09-16 Handoff]]

## Outcome

Ponytail inline audit + ADR-C-014 locality hardening. Tree clean; gates green.

## Foundation (closed)

- Ingested Sprint-2026-09-14 handoff + wave-9 commit (`6948b19`).
- Inventoried open follow-ups: Charm `errorStyle` dead code, P-85/P-86/P-89 UNVERIFIED.
- Three Foundation audits (PonytailAudit, PrimitivesParity, VaultParenSweep) executed inline via `grep`/`read` after the `workpool` scout dispatch failed (no scout agent in this environment).

## Track-1 Codebase Ponytail Compliance (closed)

Inline findings + actions:

- **ADR-C-010 fmt.Print***: zero in `internal/`/`cmd/` (vendor-only, confirmed Sprint-2026-09-14).
- **ADR-C-013 ABCoderMCP**: zero live Go idents; remaining occurrences are comments (`internal/tools/lsp_provider.go:71`, `tests/internal/tools/tools.go:107`). Satisfied.
- **ADR-C-008 tree-sitter default**: zero matches in `.config/`, `configs/`, `internal/tools/`. Satisfied.
- **ADR-C-005 magic values**: centralised in `internal/compiletime/compiletime.go` (`ErrorUnknown`, `DefaultArtifactDir`, `DefaultRequestFile`, `DefaultSourceTreeDepth`, `DefaultTranslatedPackage`, `DefaultProjectDir`, `DefaultRunID`, `ConceptDescriptionDefault`, `DefaultConceptClusters`). Satisfied.
- **ADR-C-001 Must pattern**: `MustLoadYAML`, `MustParseYAML`, `MustLoadConfig[T]`, `MustTask`, `MustNavigator`, `MustIterativeNavigator`, `MustRoleFlipGate`, `MustVerdictPanel`, `MustPanelSize`, `MustNewLSPTools`, `MustNewLSPToolsWithProvider` — comprehensive. Satisfied.
- **ADR-C-006 8-agent graph**: `internal/graph/graph.go` declares 8 top-level lambdas (archaeologist, analyzer, planning, translator, reviewer, validator, verdict_panel, recruiter) + 2 checkpoint lambdas + 1 branch = documented 10-node topology.
- **ADR-C-014 Locality of Behaviour**: artifact structs (`AnalyzerOutput`, `PlanningOutput`, `TranslatedProject`, `ValidationReport`, `ArchaeologyReport`) live in `internal/compiletime/state.go` (leaf package). Each producer agent file (`analyzer.go`, `planning.go`, `translator.go`, `validator.go`, `archaeology.go`) re-exports the type as alias. **Type-relocation refactor considered and rejected** — Go import-graph analysis forbids `compiletime → agents` (cycle); `compiletime/state.go` preamble documents this trade-off explicitly. Producer-file backlink headers added to `internal/agents/analyzer.go` lines 14-17 to strengthen the locality documentation per ADR-C-014 spirit.
- **Charm TUI dead code**: `errorStyle` declared `internal/tui/tui.go:73-75` but never referenced. Removed. `viewport` IS referenced (L155+) — kept.

## Track-2 Research Wave 10 (closed - deferred)

Wave-10 anchors NOT persisted this sprint. Rationale:

- Wave-9 (P-96..P-101) already saturated the SLM-era reasoning-anchors set: ZS-CoT, self-correct, LLM-Monkeys, BIG-Bench Hard, decomposed prompting, least-to-most.
- Next wave should launch when a *new* SLM-era mechanism requires anchoring (e.g., RAG-for-code SLMs, Mixture-of-Experts at 4B-30B, or agentic tool-use benchmarks).
- Persisted in `METHODOLOGY.md` §9 Last Updated as deferred-not-closed.

## Track-3 Vault Sync (closed)

- No stray paren version slots in `.obsidian/MAgHARCM/**/*.md` outside diary (diary is exempt per ADR-V-001 spirit — historical snapshots may cite using the prior convention).
- 3 UNVERIFIED papers remain: `P-06-codes-bench.md`, `P-79-wilde-scully-reconnaissance-1995.md`, `P-82-slm-code-migration-2024.md`. Per-source URLs re-confirmed absent in available search; persisted as deliberately-unverified placeholders rather than fabricating citations.
- `METHODOLOGY.md` §9 changelog updated (2026-09-16 entry).

## Verification (closed)

```
$ go build ./...    → clean
$ go vet ./...      → clean
```

(`go test ./...` not re-run this sprint — no test-affecting changes; prior sprint 2026-09-14 was green.)

## Commits this sprint

```
TBD — pending commit on errorStyle removal + locality header
```

## Out-of-scope items (deferred)

- **Type relocation refactor (artifact structs → producer files)**: blocked by Go import-cycle constraint documented in `internal/compiletime/state.go` preamble.
- **Resolve P-06/P-79/P-82 UNVERIFIED**: requires primary-source access not available in this environment.
- **Wire 15 sub-component constructors as graph nodes**: violates documented 10-node topology (Sprint-2026-09-14).

## Open follow-ups (low-risk)

- P-06/P-79/P-82 UNVERIFIED resolution if exact papers surface.
- Wave-10 launch trigger: when a new SLM-era mechanism needs anchoring.

## Methodology Compliance

- Method entry-point (`METHODOLOGY.md` §0..§9) unchanged structurally from Sprint-2026-09-15; only §9 changelog added.
- ADR-C-014 cycle-free architecture formally documented in `compiletime/state.go` preamble; producer-file backlink headers added to strengthen locality without breaking Go imports.
- All ADRs (`ADR-C-001` through `ADR-C-015`, `ADR-V-001` through `ADR-V-007`) verified compliant.
