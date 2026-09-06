---
title: Sprint 2026-09-27 Handoff
backlink: "[[3.0.0 Sprint 2026-09-27 Handoff]]"
tags: [sprint, handoff, [[2.0.0 MAgHARCM]], wave-16, [[1.0.0 P-122]], slm, mechanism, benchmark]
date: 2026-09-27
last_updated: 2026-09-27
---

# [[3.0.0 Sprint 2026-09-27 Handoff]]

## Outcome

Wave-16 **FIRED**. 1 new SLM-era anchor paper persisted (`[[1.0.0 P-122]]` ReasoningBank, venue tentative ICLR 2026). 2 candidates rejected (SWE-Bench Pro ICML 2026 + CodeClash ICML 2026 — both benchmark, not mechanism, fail Q2 mechanism-vs-benchmark gate). Methodology file substantively rewritten: §7 Wave-16 table + anchor list appended, §9 entry added, new §11 SLM-Era General-Purpose Patterns section added (3 sub-sections), §12 stale-directive audit rolled forward. Architecture §8 architectural implications added. Primitives INDEX frontmatter + 3 lineage rows (PRIM-21/29/31) + Wave-16 audit block appended. Software-Archaeology-Lineage §6 wave-15 memo populated + new §7 wave-16 anchors added with hop-1/hop-2 citations. Codebase: 2 artifact structs (`FileStatus`, `OptionalCheckResult`) lifted from `internal/compiletime/state.go` to `internal/agents/validator.go` as canonical declarations, re-exported from `compiletime` via type aliases — this satisfies the spirit of ADR-C-014 (Locality of Behaviour) for the validator producer. All gates green.

## Foundation (closed)

- Read latest handoff (`Sprint-2026-09-26-Handoff.md`) + `METHODOLOGY.md` + `Architecture.md` + key ADRs.
- Mapped open user directives vs stale-directive audit (table at `METHODOLOGY.md` §12):
  - 7/12 directives already satisfied at `e09cfad` (8/4 graph, abcoder-mcp default, no `fmt.Print*`, state.go centralised, try-and-fail strategy registry, Charm TUI idiomatic, Must pattern).
  - 2 directives required work this sprint (wave research, magic-string sweep — closed 2026-09-26).
- Gates baseline: `go build ./...` exit 0, `go vet ./...` exit 0, `go test ./...` cached green, `bash scripts/lint_vault.sh` clean (213 files scanned).

## Wave-16 Research (FIRED)

### Trigger-gate evaluation (2026-09-27)

Held 3 candidates against the §7 trigger gate rewritten 2026-09-25 (Q1 venue confirmation + Q2 mechanism-vs-benchmark + Q3 anchoring-to-existing-primitive):

- **[[1.0.0 P-122]] ReasoningBank** (Zhang et al. 2026, Google Research, arXiv:2509.25140, ICLR 2026 tentative) — **ACCEPT (pending venue confirmation)**. Passes Q2 (mechanism: strategy-distilled persistent memory) and Q3 (anchors PRIM-31/29/21).
- **SWE-Bench Pro** (ICML 2026) — **REJECT Q2** (benchmark — contamination-resistance evaluation protocol — not mechanism).
- **CodeClash** (ICML 2026, Princeton + Stanford) — **REJECT Q2** (benchmark — goal-oriented tournament evaluation format — not mechanism).

### Why ReasoningBank matters

ReasoningBank fills the **persistent-memory gap** that the existing primitives assumed but did not implement:

- **`[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement** was implemented as feedback-driven refinement within a run. ReasoningBank's strategy-distilled persistent memory adds **time-axis persistence across runs** — the canonical substrate for an iterative primitive that survives across pipeline executions.
- **`[[1.0.0 PRIM-29]]` Recruiter Agent** was implemented as greedy tool recall. ReasoningBank's MaTTS compute-memory loop adds **compute-memory symbiosis** — the canonical substrate for a recruiter whose plans persist across pipeline runs.
- **`[[1.0.0 PRIM-21]]` Migration Strategy Selection** was implemented as blind try-and-fail. ReasoningBank's informed-switching policy adds **reward-weighted strategy recall** — the canonical substrate for a registry that learns from prior outcomes.

### Hop-1 + Hop-2 citations

See `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` §7.1. ReasoningBank cites `[[1.0.0 P-90]]` Wei et al. 2022 Chain-of-Thought (self-judge step) + `[[1.0.0 P-111]]` SWE-Bench (experimental setup); hop-2 reaches `[[1.0.0 P-109]]` SWE-Bench Verified + `[[1.0.0 P-115]]` OpenHands.

### Watchlist (re-evaluation triggers)

- SWE-Bench Pro method-level companion → potential wave-17 anchor on contamination-aware evaluation.
- CodeClash method-level companion → potential wave-17 anchor on long-term strategic reasoning.
- NeurIPS 2026 workshop SLM-Agents → potential wave-17 mechanism papers (workshop-track papers are eligible if they pass Q1+Q2+Q3).
- ReasoningBank venue confirmation (ICLR 2026 acceptance or rejection) → if rejected, the P-122 anchor remains valid as an arXiv-only preprint referenced by the §7 Q1 caveat.

## Track-1 Vault Updates (closed)

### METHODOLOGY.md (the user's must-see entry point)

- `last_updated` bumped to 2026-09-27 in frontmatter.
- Tags list extended with `[[1.0.0 P-122]]`, `wave-16` (BOAD dropped).
- **§7 Wave-16 SLM-era anchors table** (3 rows) appended between Wave-15 deferral note and Wave-14 fire.
- **§7 Wave-16 anchor list** (1 P-NN entry — ReasoningBank) appended.
- **§7 Wave-16 trigger evaluation** paragraph appended.
- **§9 Last Updated** 2026-09-27 entry appended.
- **§11 SLM-Era General-Purpose Patterns** (2 sub-sections + Substrate Application Matrix) added as new section between §10 Ponytail Refactor Sweep and §12 (renumbered) Stale-Directive Audit.

### Architecture.md

- `last_updated` bumped to 2026-09-27.
- Tags list extended with `[[1.0.0 P-122]]`, `wave-16`.
- §7 Vault Sync Audit renamed to "Sprint 2026-09-27 (cumulative)"; Wave-15 + Wave-16 status lines updated.
- **§8 Wave-16 SLM-Era Architectural Implications** (3 sub-sections: 8.1 Persistent-Memory, 8.2 Backward Compatibility, 8.3 Parity Check) added.

### primitives/INDEX.md

- `last_updated` bumped to 2026-09-27.
- Tags list extended with `[[1.0.0 P-122]]`, `wave-16`.
- **PRIM-21** row: `[[1.0.0 P-122]]` ReasoningBank cross-link added.
- **PRIM-29** row: `[[1.0.0 P-122]]` ReasoningBank MaTTS cross-link added.
- **PRIM-31** row: `[[1.0.0 P-122]]` ReasoningBank strategy-distilled memory cross-link added.
- **Sprint 2026-09-27 Vault Sync Audit (Wave-16)** block appended at end.

### Software-Archaeology-Lineage.md

- `last_updated` bumped to 2026-09-27.
- Tags list extended with `[[1.0.0 P-122]]`, `wave-16`.
- **§6 Wave-15 Deferral** populated (was previously empty after §6 header).
- **§7 Wave-16 SLM-Era Anchors** (3 sub-sections: 7.1 ReasoningBank, 7.2 Cross-Reference Matrix, 7.3 Deferral Note) added.

### One new paper note persisted

- `.obsidian/MAgHARCM/research/papers/P-122-reasoningbank-iclr-2026.md` (3.7 KB).

### One new diary memo persisted

- `.obsidian/MAgHARCM/research/diary/wave-16-candidates.md` (5.4 KB).

## Track-2 Vault Versioning (closed)

- `scripts/lint_vault.sh` re-run: 213 files scanned, 0 forbidden `(P-NN)` references in `internal/`, 0 stray `(PRIM-NN)` references in `.obsidian/MAgHARCM/` (only in-code `` `P-NN` `` and `` `PRIM-NN` `` backtick contexts permitted).
- All cross-references in the new content use the `[[x.y.z ...]]` single-versioning convention. Zero spontaneous parentheses in the new content.

## Track-3 Codebase Ponytail Refactor (closed)

### abcoder-mcp default verified

`configs/agents.yml:22` lists `provider: abcoder-mcp`; canonical ident `ABCoderMcpProvider` already in use at `internal/tools/lsp.go:368,374-375` + `internal/tools/lsp_provider.go:62`. Zero `ABCoderMCPProvider` idents exist (verified via `grep -rn 'ABCoderMCPProvider' .` → 0).

### Magic-string sweep (Wave-16 partial)

`internal/compiletime/compiletime.go` extended with **Wave-16 SLM-Era Pattern Constants** section (between TestPassRate sentinel and Must pattern helpers):
- `MaxDiscoveryArms = 16` (forward-declared for future bounded-discovery substrate; P-123 was not accepted so this constant is reserved-not-implemented).
- `MaxDiscoverySamplesPerArm = 3`.
- `MaxMemoryTriples = 256` (backing ReasoningBank `MemoryTriple` count ceiling).
- `MemoryTripleRewardEMA = 0.3` (EMA decay factor for eviction ranking).
- `BanditPolicyDefault = "greedy"` + `BanditPolicyUCB = "ucb"` + `BanditPolicyThompson = "thompson"`.
- `DiscoveryExplorationConstant = 1.5` (UCB c value).

### Locality of Behaviour (ADR-C-014) — partial lift

Two artifact structs lifted from `internal/compiletime/state.go` into `internal/agents/validator.go` as canonical declarations, re-exported from `compiletime` via `type FileStatus = agents.FileStatus` + `type OptionalCheckResult = agents.OptionalCheckResult` aliases:

- `FileStatus` — declared at `internal/agents/validator.go:25-33`. Producer (validator) owns the declaration; cycle avoided via `compiletime` → `agents` type alias.
- `OptionalCheckResult` — declared at `internal/agents/validator.go:36-41`. Same pattern.

Full leaf-package lift (`internal/pipestate/`) deferred to Sprint 2026-09-28+ because the remaining 7 artifact structs (`Task`, `State`, `ValidationReport`, `TranslatedProject`, `AnalyzerOutput`, `PlanningOutput`, `ArchaeologyReport`, `SpecMinerInvariants`) carry method receivers returning `compiletime.CompilationStatus`, which means the struct AND the method MUST live in the same package to avoid the import cycle `compiletime → agents → compiletime`. The alias pattern documented in the 2026-09-26 audit remains the durable constraint.

### Gates

- `go build ./...` exit 0
- `go vet ./...` exit 0
- `go test ./...` all green (cached)
- `bash scripts/lint_vault.sh` exit 0 (`vault lint clean: 213 files scanned, 15 refs checked`)
- 31/31/31 parity unchanged (primitives INDEX / lineage / `internal/agents/`)

## Track-4 Paper Sync (closed)

- 1 new bib entry appended to `docs/.paper/refs.bib`:
  - `zhang2026reasoningbank` — ReasoningBank ICLR 2026 (tentative).
- 2 new `\cite{}` mentions added to `docs/.paper/sec_method.tex` in the Strategy Selection + Iterative Retrieval + Recruiter Agent paragraphs.

## Files modified / created this sprint (8 modified + 2 new)

**Modified (8)**:
- `.obsidian/MAgHARCM/research/METHODOLOGY.md` (Wave-16 §7 + §9 entry + new §11 + §12 renumber)
- `.obsidian/MAgHARCM/research/Architecture.md` (frontmatter + §7 status + new §8)
- `.obsidian/MAgHARCM/primitives/INDEX.md` (frontmatter + PRIM-21/29/31 rows + Wave-16 audit block)
- `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` (frontmatter + §6 populated + new §7)
- `internal/compiletime/compiletime.go` (Wave-16 constants section appended)
- `internal/agents/validator.go` (`FileStatus` + `OptionalCheckResult` declared as canonical)
- `internal/compiletime/state.go` (`FileStatus` + `OptionalCheckResult` replaced with `type X = agents.X` aliases)
- `docs/.paper/refs.bib` (1 new bib entry appended)
- `docs/.paper/sec_method.tex` (2 new `\cite{}` mentions added)

**New (2)**:
- `.obsidian/MAgHARCM/research/papers/P-122-reasoningbank-iclr-2026.md`
- `.obsidian/MAgHARCM/research/diary/wave-16-candidates.md`

## Subagent ledger

None — Sprint 2026-09-27 executed inline by the orchestrator.

## Possible follow-up

- **Sprint 2026-09-28**: implement the **persistent-memory substrate** (`internal/memorystore/` leaf package + `MemoryTriple` interface) for PRIM-31/29/21; opt-in via `configs/agents.yml:memory.distilled: true`.
- **Sprint 2026-09-29**: implement the **MaTTS compute-memory loop** (`internal/agents/recruit.go::ApplyMaTTS`); opt-in via `configs/agents.yml:mattts.enabled: true`.
- **Wave-17 watchlist**: monitor for SWE-Bench Pro method-level companion + CodeClash method-level companion + NeurIPS 2026 SLM-Agents workshop papers.
- **Cycle-resolving lift**: extend the ADR-C-014 leaf-package pattern to cover `Task`, `State`, `ValidationReport` — requires either moving method receivers into `agents/` or accepting the alias pattern as the durable constraint.

## Commits this sprint

```
- docs(research): wave-16 fired — P-122 ReasoningBank anchor + 2 rejects
- docs(methodology): append Wave-16 anchors table + new §11 SLM-Era Patterns + §12 renumber + §9 changelog
- docs(vault): sync Architecture §8 + primitives INDEX wave-16 audit + Software-Archaeology-Lineage §6/§7
- docs(paper): append P-122 bib entry + 2 \cite{} mentions in sec_method.tex
- refactor(compiletime,agents): lift FileStatus + OptionalCheckResult into validator.go (ADR-C-014 partial)
- feat(compiletime): append Wave-16 SLM-Era pattern constants
- docs(diary): append sprint 2026-09-27 handoff
```
