---
title: Sprint 2026-09-07 Handoff
backlink: "[[3.0.0 Sprint 2026-09-07 Handoff]]"
tags: [sprint, handoff, [[2.0.0 MAgHARCM]], [[1.0.0 PRIM-31]], slm]
---

# [[3.0.0 Sprint 2026-09-07 Handoff]]

## Outcome

Closed all 22 sprint items. Three phases finished: Foundation, Track-1 Codebase (ADR-C-014 + SLM strategy), Track-2 Research + Vault. Five focused commits; clean tree; green `go build`, `go vet`, `go test ./...`.

## Foundation (closed)

- Audited open items from Sprint-2026-09-06 handoff.
- Persisted two ADRs binding worker rules:
  - [[ADR-2026-09-07-Sprint-Conventions]] — single-versioning convention `[[x.y.z ...]]`, no stray parens.
  - [[ADR-2026-09-07-Dup-Row-Escape-Recipe]] — three questions (Q1 am-I-on-cited-row, Q2 cite-count >=2, Q3 semantically-distinct) that decide when a new citation row is a duplicate.
- Inventoried `state.go` artifact locality and SLM-aware tooling surface.

## Track-1 Codebase (ADR-C-014 + SLM strategy) — closed

- **[[1.0.0 ADR-C-014]]** — `internal/types/state.go` removed. Centralised state + 6 artifact structs now live in `internal/compiletime/state.go`. Producers reference these structs locally; one source of truth. Commit `0cb5994`.
- **SLM-aware patterns (4B-30B)** — injected `SLMPromptContractPreamble` into the agent prompt template pipeline. Centralised `ArchaeologySkipDirs` (vendor / node_modules / .git / build caches) instead of per-call filters. Commit `2bb9082`.
- **Hardcode purge** — swept `internal/` for residual magic values; routed them through `compiletime.*` or per-module `MustCompile*` constructors.
- **ASD-STE100 fix-up** — subagent pass over user-visible strings; one committed tweak to `internal/tui/tui.go` to drop the awkward "(dry-run) would execute with:" phrasing. Commit `375b41e`.

## Track-2 Research + Vault — closed

Four anchor papers persisted with hop-1 + hop-2 sections and cross-referenced into the lineage matrix and primitive INDEX:

| Paper | Title | Anchors |
|---|---|---|
| [[1.0.0 P-50]] | Jiang et al. (2024) Survey on LLMs for Code Generation | [[1.0.0 PRIM-26]] (navigator RAG strategy) |
| [[1.0.0 P-51]] | Cassano et al. (2024) Can It Edit? | [[1.0.0 PRIM-22]], [[1.0.0 PRIM-25]] |
| [[1.0.0 P-52]] | Wang et al. (2023) Self-Consistency | [[1.0.0 PRIM-7]], [[1.0.0 PRIM-21]] |
| [[1.0.0 P-53]] | Liu et al. (2023) Lost in the Middle | [[1.0.0 PRIM-3]], [[1.0.0 PRIM-23]], [[1.0.0 PRIM-31]] |

Lineage matrix updated; INDEX cross-refs added. Dup-row escape recipe (Q1/Q2/Q3) applied — no duplicate rows introduced. Vault swept for stray parentheses in version markers — all clean.

Commit `da424a1` (papers) + `b32dc8d` (INDEX + matrix cross-refs).

## Verification — closed

- `go build ./...` — clean.
- `go vet ./...` — clean.
- `go test -count=1 ./...` — green across `cmd/MAgHARCM-tui`, `internal/agents`, `internal/config`, `internal/languages`, `internal/logger`, `internal/runner`, `internal/tools`.
- `git status` — clean; 5 focused commits this sprint.

## Primitive parity (single source of truth)

- 35 agents/agent-support files in `internal/agents/` ↔ 31 primitives in INDEX.
- Each primitive cites an implementation file in the INDEX and the implementing `.go` file in the matrix.
- Each P-NN citation sits on the row whose PRIM is its *direct* anchor.

## Open items / next sprint

- [[1.0.0 PRIM-32]] — candidate (PRIM-31 already at iterative-retrieval). No new primitive planned this sprint.
- Possible follow-up: extend the dup-row recipe into a lint check that scans `Software-Archaeology-Lineage.md` automatically before commit.
- Continue P-54+ research only when a new SLM-aware mechanism needs an anchor.

## Commits this sprint

```
b32dc8d docs(research): persist P-50..P-53 anchor papers; cross-ref into lineage matrix + INDEX
da424a1 docs(research): persist P-50..P-53 anchor papers; cross-ref into lineage matrix + INDEX
375b41e refactor(ste100): drop awkward '(dry-run) would execute with:' phrasing in TUI
2bb9082 refactor(slm): centralise ArchaeologySkipDirs + inject SLM-aware prompt preamble (4B-30B)
0cb5994 refactor(state): canonicalise State + 6 artifact structs in compiletime/state.go (ADR-C-014)
```
