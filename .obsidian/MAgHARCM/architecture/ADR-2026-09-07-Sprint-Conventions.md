---
title: "ADR-2026-09-07 Sprint Conventions — Authoritative Rule-Set for MAgHARCM Codebase & Vault"
backlink: "[[1.0.0 ADR-2026-09-07-Sprint-Conventions]]"
status: Accepted
date: 2026-09-07
tags: [adr, conventions, [[2.0.0 MAgHARCM]], [[2.0.0 Architecture]]]
---

# [[1.0.0 ADR-2026-09-07 Sprint Conventions]]

## 1. Context

The user-directed refactor set was re-issued in Sprint 2026-09-04, again in 2026-09-05, and again in 2026-09-06. Each issuance forces a fresh audit of the same conventions. Sprint 2026-09-06 closed with three of them already-satisfied (try-and-fail migration, 8-agent graph, Charm TUI adoption) and several others partially-satisfied. Without a single authoritative document, every dispatch re-litigates the same questions.

This ADR fixes the rule-set as of Sprint 2026-09-07 and binds all subsequent worker dispatch to it.

## 2. Decision

The following conventions are authoritative for the MAgHARCM codebase and Obsidian vault. Workers MUST treat this ADR as their starting contract; deviations require a new ADR or an explicit `rejected: <reason>` annotation.

### 2.1. Codebase Conventions

| ID | Rule | Enforcement |
| :--- | :--- | :--- |
| [[1.0.0 ADR-C-001]] | Compile-time inits & configs MUST NOT have fallbacks. Use `MustXxx` patterns (`MustParse`, `MustRead`, `MustCompile`). | `internal/compiletime/` (no fallbacks); all callers use `Must*`. |
| [[1.0.0 ADR-C-002]] | Configurations MUST be provided in external YAML files (`configs/*.yml`); no inline literals in agent modules. | Loader `configs.load()`; agents take typed config struct. |
| [[1.0.0 ADR-C-003]] | Each unit (function/module) MUST communicate through clear interfaces, not internal-struct knowledge. | Agent `Inputs`/`Outputs` structs; no cross-imports into producer state. |
| [[1.0.0 ADR-C-004]] | Migration strategy selection MUST be try-and-fail incremental: `Registry.TryInOrder(strategies, ctx)` returns the first `Matches(ctx).Attempt(ctx)` result; `Miss` advances. | `internal/agents/strategy.go:TryInOrder`. |
| [[1.0.0 ADR-C-005]] | Hard-coded magic values MUST be eliminated. Hard-coded-by-necessity values (e.g. error sentinel `"unknown"`) live in a centralised compile-time config (`internal/consts/`), not in agent modules. | `internal/consts/consts.go`. |
| [[1.0.0 ADR-C-006]] | The agent graph MUST wire ALL agents defined in `internal/agents/`. As of Sprint 2026-09-06, the graph has 8 nodes: archaeologist, analyzer, planning, translator, roleflip, validator, verdict_panel, recruiter. | `internal/graph/graph.go`. |
| [[1.0.0 ADR-C-007]] | All log messages and user-facing strings MUST comply with `asd-ste100` skill (STE-flavored). | AD-HOC review; sweep subagent each sprint. |
| [[1.0.0 ADR-C-008]] | AST analysis MUST default to abcoder's MCP (`ABCoderMcpProvider`), not tree-sitter. Tree-sitter may remain as a documented fallback only when abcoder is unavailable. | `internal/tools/lsp_provider.go` (`ABCoderMcpProvider` is the canonical ident). |
| [[1.0.0 ADR-C-009]] | Per-project compilation status is BINARY: `CompilationPass` or `CompilationFail`. There is no partial / rate. | `internal/compiletime/types.go`. |
| [[1.0.0 ADR-C-010]] | NO `fmt.Print*` in non-agent modules (CLI args / prompt assembly excepted). All runtime logging MUST use the structured logger. | `logger.*` everywhere except CLI entry-point banner. |
| [[1.0.0 ADR-C-011]] | The TUI MUST be built on the Charm stack: `charm.land/bubbletea`, `charm.land/bubbles`, `charm.land/lipgloss`, `charm.land/glamour`. Stack usage MUST be idiomatic (no manual ANSI escapes, no reinvented layout primitives). | `internal/tui/tui.go`. |
| [[1.0.0 ADR-C-012]] | Prefer external libraries over self-implemented code when the library covers ≥80% of the use case idiomatically. | Lazy/Ponytail rule. |
| [[1.0.0 ADR-C-013]] | Identifier naming MUST use Go-conventional camelCase where appropriate. `ABCoderMcpProvider` (NOT `ABCoderMCPProvider`). | Repo-wide grep. |
| [[1.0.0 ADR-C-014]] | `State` schema lives in `internal/compiletime/state.go` (centralised, schema-versioned). Per-agent artifact structs (`AnalyzerOutput`, `PlanningOutput`, etc.) MUST live alongside their producer agent file. | `internal/agents/<agent>.go`; cross-refs from `compiletime` only. |
| [[1.0.0 ADR-C-015]] | Migration to a strategy failing partial steps MUST increment attempt counter; retry policy lives in the registry, not the strategy. | `Strategy.Attempt` returns `Partial`/`Miss`; `Registry.TryInOrder` advances on Miss. |

### 2.2. Vault Conventions

| ID | Rule | Enforcement |
| :--- | :--- | :--- |
| [[1.0.0 ADR-V-001]] | Version markers MUST use `[[x.y.z ...]]` wikilink syntax. Stray parentheses (e.g. `(P-46)`) are PROHIBITED in version slots. | `grep` sweep each sprint. |
| [[1.0.0 ADR-V-002]] | Primitives (`PRIM-*`) are SINGLE-SOURCE-OF-TRUTH. The index `.obsidian/MAgHARCM/primitives/INDEX.md`, the lineage matrix `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`, and the codebase `internal/` MUST agree on the full primitive set. | Parity check `P-31` in INDEX ≡ rows in lineage ≡ files in `internal/`. |
| [[1.0.0 ADR-V-003]] | Hop-1 + hop-2 citation chains MUST be persisted per paper. Hop-1 = the paper's own bibliography; hop-2 = papers-cited-by-hop-1 authors. | Each `P-NN` paper file has a `## Hop-1` and `## Hop-2` section. |
| [[1.0.0 ADR-V-004]] | ADRs live in `.obsidian/MAgHARCM/architecture/`. New ADR = new file `ADR-YYYY-MM-DD-<slug>.md`. | Directory layout. |
| [[1.0.0 ADR-V-005]] | Sprint handoffs live in `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md`. | Directory layout. |
| [[1.0.0 ADR-V-006]] | Methodological / architectural discussion lives in `.obsidian/MAgHARCM/research/`. | Directory layout. |
| [[1.0.0 ADR-V-007]] | A primitive is considered "PRESENT" iff the lineage matrix lists it AND the codebase file under `internal/` declares a function/type matching its name. | Parity check. |

## 3. Consequences

Positive:
- Worker dispatch in subsequent sprints reads this ADR as its contract; no re-litigation of the rule-set.
- Convention drift requires an explicit ADR supersession.
- Audit cost per sprint drops (parity sweep only).

Negative:
- An ADR adds documentation maintenance cost.
- A rule may be superseded by user direction; the ADR must be amended, not silently violated.

## 4. Compliance Map (Sprint 2026-09-06 state)

| Rule | Status | Evidence |
| :--- | :--- | :--- |
| ADR-C-001 | ✅ | `internal/compiletime/*.go` |
| ADR-C-002 | ✅ | `configs/*.yml` |
| ADR-C-003 | ✅ | Agent `Inputs/Outputs` interfaces |
| ADR-C-004 | ✅ | `internal/agents/strategy.go:TryInOrder` |
| ADR-C-005 | ⚠️ | `internal/consts/` exists; sweep ongoing |
| ADR-C-006 | ✅ | `internal/graph/graph.go` (8 nodes) |
| ADR-C-007 | ⚠️ | Sweep pending; flagged residue exists |
| ADR-C-008 | ✅ | `internal/tools/lsp_provider.go` |
| ADR-C-009 | ✅ | `internal/compiletime/types.go` |
| ADR-C-010 | ✅ | Audit clean (Sprint 2026-09-06) |
| ADR-C-011 | ✅ | `internal/tui/tui.go` |
| ADR-C-012 | ✅ | Charm stack adopted |
| ADR-C-013 | ✅ | `ABCoderMcpProvider` |
| ADR-C-014 | ⚠️ | State moved; per-agent split pending |
| ADR-C-015 | ✅ | `Strategy.Attempt` returns Miss/Attempt |
| ADR-V-001..007 | ✅ | Verified Sprint 2026-09-06 |

## 5. Supersession

To supersede a rule, write a new ADR referencing the rule ID with explicit `supersedes: ADR-C-NNN` and update this ADR's compliance map.

## 6. Pointers

- Companion ADR: [[1.0.0 ADR-2026-09-07-Dup-Row-Escape-Recipe]]
- Methodology: `.obsidian/MAgHARCM/research/METHODOLOGY.md`
- Lineage: `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`
- Handoff (current sprint): `.obsidian/MAgHARCM/diary/Sprint-2026-09-06-Handoff.md`
