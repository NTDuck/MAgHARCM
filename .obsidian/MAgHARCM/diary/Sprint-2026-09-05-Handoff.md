---
title: Sprint 2026-09-05 Handoff — P-42..P-45 Synthesis Anchors & Ponytail Closure
backlink: [[2.0.0 Sprint-2026-09-05]]
tags: [sprint, handoff, ponytail, synthesis, lineage, [[2.0.0 MAgHARCM]]]
---

# [[2.0.0 Sprint 2026-09-05 Handoff — P-42..P-45 Synthesis Anchors & Ponytail Closure]]

Continuation of the modernization track from `[[2.0.0 Sprint-2026-09-04-Handoff]]`. Read that file first; current state continues from §5.1 of that handoff.

## 1. Outcome Summary

1. **Ponytail closure applied.** `internal/compiletime/compiletime.go` grew from 172 to 431 lines as the central source of typed enums (ArchitectureStabilityLayer, SpecLifecyclePhase, Verdict), sentinel strings (ErrorUnknown, DefaultProjectDir), magic numbers (BigBangFileMax, PilotFileMin, ParallelCutoverFileMin, IterativeContextBudgetBytes, ConceptTokenMinLength), format strings (CheckpointFilePattern, CheckpointExt), cluster tables (DefaultConceptClusters), and description tables (ArchitectureStabilityDescriptionL1/L2/L3, ComprehensionRecognition*, ComprehensionExplanationDefault). 13 agent files migrated to compile-time references; 11 of those gained `Must*` constructors. The orphaned `CurrentSchemaVersion` re-export in `internal/agents/state.go` was deleted; callers now reference `compiletime.CurrentSchemaVersion` directly.
2. **Imports restored.** `strings`, `unicode`, and `eino/components/model` were restored in `internal/agents/analyzer.go` and `internal/agents/planning.go` after PonytailFix dropped them; `go build ./...` and `go vet ./...` are green; `go test -count=1 ./...` reports 9/9 packages passing.
3. **Four synthesis-anchor papers persisted.** `[[1.0.0 P-42]]` (MacCormack DSM 2006, anchors PRIM-3 + PRIM-19), `[[1.0.0 P-43]]` (Kang FODA 1990, anchors PRIM-10), `[[1.0.0 P-44]]` (Corkill Blackboard 1991, anchors PRIM-17), `[[1.0.0 P-45]]` (Curtis-Kellner-Over Process Modelling 1992, anchors PRIM-24).
4. **Lineage + INDEX rewired.** `Software-Archaeology-Lineage.md` matrix cross-references P-42..P-45 in their respective PRIM rows; new section 5 documents the four synthesis citations. `primitives/INDEX.md` does the same in compact form.
5. **Versioning convention clean.** No stray parentheses in version markers anywhere in the vault. The convention `[[x.y.z ...]]` is the only form used.

## 2. Codebase State

### 2.1. Build & Test
```
$ go build ./...          # no output (clean)
$ go vet ./...            # no output (clean)
$ go test -count=1 ./...  # 9/9 packages pass
```

### 2.2. Compile-Time Centralization
The following modules now consume compile-time constants exclusively (no inline magic literals):

| Module | Old | New |
| :--- | :--- | :--- |
| `strategy.go` | `switch` over StrategyIDs; magic LoC numbers | `compiletime.StrategyRationale*` map; `BigBangFileMax/LoCMax`, `PilotFileMin/LoCMin`, `ParallelCutoverFileMin` |
| `verdict_panel.go` | Stringly-typed `VerdictKind*` constants | `compiletime.VerdictEquivalent/NotEquivalent` typed enum + `VerdictAlias*` slice; `VerdictJudgeIDPrefix` |
| `optional_checks.go` | String verdict tokens | Typed `Verdict` enum (`Pass`/`Fail`/`Skipped`); `OptionalCheck*` names |
| `concept_assignment.go` | Inline keyword clusters | `compiletime.DefaultConceptClusters` table; `ConceptTokenMinLength` |
| `design_rule_hierarchy.go` | `LayerL1/L2/L3` int constants | `compiletime.ArchitectureStabilityLayer` typed enum + `DescriptionL1/L2/L3` |
| `comprehension.go` | 'Math/Numerics Library' literal | `compiletime.ComprehensionRecognition*` + `ComprehensionExplanationDefault` |
| `recruit.go` | Inline tool/agent name strings | `compiletime.Tool*/Agent*` |
| `roleflip.go` | Private constants + `errors.New` | `compiletime.RoleFlip*` typed sentinels + `ErrRoleFlipGateNotConfigured` |
| `iter_retrieval.go` | Hard-coded `4096` | `compiletime.IterativeContextBudgetBytes` |
| `navigator.go` | `'.'` project path | `compiletime.DefaultProjectDir`; `ErrNavigatorNoProvider` |
| `checkpoint.go` | Hard-coded `'default'`/`'iter-%04d.json'`/`0o755`/`0o644` | `compiletime.DefaultRunID`/`CheckpointFilePattern`/`CheckpointDirMode`/`CheckpointFileMode` |
| `state.go` | Re-exported `CurrentSchemaVersion` | Deleted; callers reference `compiletime.CurrentSchemaVersion` directly |

### 2.3. Backbone Integrity
- All 31 primitives (`PRIM-1` .. `PRIM-31`) appear in `primitives/INDEX.md`; all 31 rows appear in `Software-Archaeology-Lineage.md` matrix. No duplicates, no orphans.
- 8 dedicated agent nodes (`archaeologist`, `analyzer`, `planning`, `translator`, `validator`, `roleflip`, `verdict_panel`, `recruiter`) wired through the Eino graph in `internal/graph/graph.go`.

## 3. Research State

### 3.1. New Papers
| ID | Authors | Year | Anchors |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 P-42]]` | MacCormack, Rusnak & Baldwin | 2006 | PRIM-3, PRIM-19 (DSM cycle analysis) |
| `[[1.0.0 P-43]]` | Kang, Cohen, Hess, Novak, Peterson | 1990 | PRIM-10 (FODA feature models) |
| `[[1.0.0 P-44]]` | Corkill | 1991 | PRIM-17 (blackboard architecture) |
| `[[1.0.0 P-45]]` | Curtis, Kellner & Over | 1992 | PRIM-24 (process modelling) |

Each paper follows the standard schema: frontmatter with `title`, `bibkey`, `tags`; frontmatter block with Authors/Year/Venue/eprint/Cited by; Summary; Relevance to MAgHARCM (cross-refs to PRIM entries); Hop-1 References; Hop-2 Anchors (software-archaeology lean); Backlinks.

### 3.2. Lineage Matrix Cross-References
- PRIM-3 row: `Baldwin & Clark (2000) Design Rules; [[P-42]] DSM (MacCormack 2006)`
- PRIM-10 row: `Oxidizer-2023, RustRepoTrans; [[P-43]] FODA (Kang 1990)`
- PRIM-17 row: `CAID-2024; [[P-44]] Corkill (1991)`
- PRIM-19 row: `Baldwin & Clark (2000), [[P-41]] 2024 retrospective, [[P-42]] MacCormack (2006) DSM`
- PRIM-24 row: `MetaGPT-2023 SOP Contracts; [[P-45]] Curtis-Kellner-Over (1992)`

### 3.3. Vault Hygiene
- `[[x.y.z ...]]` is the sole version-marker convention. No stray `(v0.x)`, `(phase X)`, `(method Y)`, `(cycle N)` parentheses survive.
- All version-marker backlinks resolve to existing files in the vault.

## 4. Open Work / Follow-Ups

1. **`internal/agents/state.go` move.** Per the user's directive, `State` and intermediate role artifacts should eventually move out of `internal/agents/state.go` to their producer modules (Locality of Behaviour). Currently `State` holds `AnalyzerOutput`, `PlanningOutput`, `TranslatedProject`, `ValidationReport`, `SpecMinerInvariants`, `ArchaeologyReport`. The cleanest move splits these into per-agent output structs owned by the agent that produces them, with `State` reduced to a slim bag of pointers + iteration counters. This is a structural refactor — touches graph wiring + runner state threading — and was deferred from this sprint.
2. **`ABCoderMCPProvider` rename.** ✅ Completed — `internal/tools/lsp_provider.go` and all callers now use `ABCoderMcpProvider` throughout.
3. **TUI migration to Bubble Tea stack.** The current TUI uses ad-hoc rendering. Migration to Bubble Tea + Bubbles + Lip Gloss + Glamour is a substantial new-feature sprint.
4. **Per-project compilation rate vs pass-rate.** Per user directive: project compilation status is strictly binary (Pass/Fail) — no partial pass rate per file. This is already enforced; needs documentation update only.
5. **`fmt.Print*` audit.** All `fmt.Print*` should be replaced with `logger.*` calls. The graph + runner code is partially migrated; remaining `fmt.Print*` calls in non-agent modules still need a sweep.

## 5. Verification Commands

```
$ go build ./...
$ go vet ./...
$ go test -count=1 ./...     # 9/9 packages pass
$ git status --short        # clean
```

## 6. Pointers

- PonytailFix transcript: `history://PonytailFix` (3h45m, completed)
- Previous sprint handoff: `.obsidian/MAgHARCM/diary/Sprint-2026-09-04-Handoff.md`
- Modernization note: `.obsidian/MAgHARCM/diary/Sprint-2026-09-04-Modernization.md`
- Recon note: `.obsidian/MAgHARCM/diary/Sprint-Recon-2026-09-04.md`
