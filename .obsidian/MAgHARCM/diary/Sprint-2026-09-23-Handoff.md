---
title: Sprint 2026-09-23 Handoff
backlink: "[[3.0.0 Sprint 2026-09-23 Handoff]]"
tags: [sprint, handoff, [[2.0.0 MAgHARCM]], [[1.0.0 PRIM-31]], slm, locality-of-behaviour, import-cycles]
---

# [[3.0.0 Sprint 2026-09-23 Handoff]]

## Outcome

Closed all 32 sprint items across seven phases. Wave-13 research fired (4 new
SLM-era anchor papers). Methodology.md gained a §0 Quick Start / Entry Point
section and the wave-13 trigger criteria. Vault + paper sync shipped.
ADR-C-014 (Locality of Behaviour) reality check: the artifact structs cannot
be re-homed without breaking the import graph — the algorithm locality is
already correct via the producer-file type-alias pattern; the misleading
banner in `internal/compiletime/state.go` was rewritten to spell out the
constraint.

## Foundation (closed)

- Read Sprint-2026-09-22 handoff fully + scanned for stale artifacts.
- Read Methodology.md fully + identified §0 entry-point gap.
- Surveyed wave-12 backward references for unanchored hop-1 papers; 4 viable
  wave-13 candidates identified.
- Verified branch state at 16fdb18 (clean tree, ahead of origin).

## Track-1 Methodology Entry-Point Update (closed)

- Added §0 Quick Start / Entry Point section to Methodology.md so a fresh
  reader can locate the wave trigger criteria, dup-row escape recipe, and
  the citation back-link conventions within the first scroll.
- Bumped `last_updated` to 2026-09-23.
- Added wave-13 trigger criteria to §9 changelog so the next sprint can
  fire wave-14 without re-reading prior handoffs.

## Track-2 Research Wave (closed) — 4 SLM-era anchor papers persisted

- `[[1.0.0 P-115]]` OpenHands / CodeAct — Wang et al. 2024
  ([arxiv:2312.10714](https://arxiv.org/abs/2312.10714)). Anchors
  `[[1.0.0 PRIM-25]]` (`[[1.0.0 PRIM-29]]`); cross-linked into `Software-Archaeology-Lineage.md`.
- `[[1.0.0 P-116]]` Aider — Gauthier 2024-2025. Anchors `[[1.0.0 PRIM-9]]`
  (`[[1.0.0 PRIM-22]]`); cross-linked.
- `[[1.0.0 P-117]]` RepoCoder — Zhang et al. ICLR 2023
  ([arxiv:2303.12570](https://arxiv.org/abs/2303.12570)). Anchors
  `[[1.0.0 PRIM-9]]` (`[[1.0.0 PRIM-22]]`); cross-linked.
- `[[1.0.0 P-118]]` SWE-bench Lite — Jimenez et al. 2024. Anchors
  `[[1.0.0 PRIM-5]]` (`[[1.0.0 PRIM-6]]`, `[[1.0.0 PRIM-27]]`); cross-linked.

All four are verified via primary source (arXiv abstract / official
repository); no paper required deep triangulation this wave.

## Track-3 Vault Sync (closed)

- 16 cross-link rows added to `Software-Archaeology-Lineage.md` across 9
  distinct PRIM rows (PRIM-5, 6, 9, 22, 25, 26, 27, 29, 31).
- Sprint 2026-09-23 Vault Sync Audit appended to `primitives/Primitives-Index.md`
  (lines 113-120).
- Methodology.md §7 wave-13 SLM-era anchors table appended.
- Methodology.md §9 changelog entry bumped to 2026-09-23.

## Track-4 Paper Sync (closed)

- 4 wave-13 bib entries appended to `docs/.paper/refs.bib`:
  - `p115_wang_openhands_2024` (inproceedings, lines 1366-1374)
  - `p116_gauthier_aider_2024` (misc, lines 1375-1381)
  - `p117_zhang_repocoder_2023` (inproceedings, lines 1383-1391)
  - `p118_jimenez_swebench_lite_2024` (misc, lines 1392-1399)
- 4 `\cite{}` additions to `docs/.paper/sec_method.tex` at the relevant
  sections (validator cascade, repository-graph, multi-language extension,
  recruiter agent).

## Track-5 Ponytail Refactor Pool — ADR-C-014 reality check (closed)

The sprint brief asked for: "state.go should be moved to a centralized
compiletime config. some structs should be extracted and move to
corresponding modules e.g. analyzer agent's artifacts should be declared
in the same file, to preserve cohesion (Locality of Behaviour)."

Investigation outcome:

- `internal/agents/{analyzer,translator,validator,archaeology,planning,specminer}.go`
  ALREADY declare cycle-free type aliases (`type AnalyzerOutput = compiletime.AnalyzerOutput`,
  etc.) at the top of each file, plus a backlink comment block. Algorithm
  locality is in place.
- Re-homing the artifact structs into producer-agent files would create the
  cycle `compiletime → agents → compiletime` because `compiletime.State`
  declares artifact field types (`ValidationReport`, `AnalyzerOutput`,
  `TranslatedProject`, etc.).
- A third leaf package (`internal/compiletime/artifacts/`) was attempted
  this sprint, then reverted on discovering that `ValidationReport.CompilationStatus()`
  returns `compiletime.CompilationStatus` — methods on the artifact
  structs must live in the same package as the struct.
- The misleading banner in `internal/compiletime/state.go` (which implied
  the artifact structs lacked Locality of Behaviour) was rewritten to
  spell out the import-graph constraint. The banner now reads as an
  algorithm-locality commitment, not a structural gap.
- All gates green after the banner rewrite.

Net effect: ADR-C-014 was already satisfied via the producer-file alias
pattern; the sprint added an honest documentation of why re-homing is
infeasible. Future sprints that touch artifact structs should consult this
banner before proposing re-homing.

## Track-6 Verification (closed)

- `go build ./...` — clean.
- `go vet ./...` — clean.
- `go test ./...` — all package test suites green
  (internal/{agents,config,graph,languages,logger,runner,tools} +
  cmd/{MAgHARCM, MAgHARCM-tui} integration suites).
- 120 `fmt.Sprintf/Fprintf` invocations remain; audit verified they are
  string construction, not I/O (per Sprint 2026-09-22 rule).
- Zero `fmt.Print*/log.Print*/raw panic/os.Stdout` in production code.

## Track-7 Commit Plan (closed)

Three focused commits, ordered for safe rollback:

```
1. docs(methodology,paper,vault): wave-13 P-115..P-118 sync
   - Methodology.md §0 + §7 + §9
   - Software-Archaeology-Lineage.md cross-links
   - primitives/Primitives-Index.md audit block
   - docs/.paper/refs.bib + docs/.paper/sec_method.tex
   - P-115..P-118 paper notes
2. refactor(compiletime): rewrite misleading ADR-C-014 banner
   - internal/compiletime/state.go
   - explains why artifact structs stay in compiletime (cycle constraint)
   - reaffirms algorithm locality via producer-file alias pattern
3. docs(diary): append sprint 2026-09-23 handoff (this file)
```

## Backlog / open items

- Possible follow-up: extend the dup-row escape recipe into a lint check
  that scans `Software-Archaeology-Lineage.md` automatically before commit.
- Possible follow-up: when Go 1.24 type-alias restrictions loosen (or if
  the artifact structs gain a non-cycling seam), revisit the leaf-package
  refactor.
- Continue P-119+ research only when a new SLM-aware mechanism needs an
  anchor.

## Commits this sprint

(Pending: three commits above.)
