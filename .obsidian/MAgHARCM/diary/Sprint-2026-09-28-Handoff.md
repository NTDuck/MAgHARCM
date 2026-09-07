---
title: Sprint 2026-09-28 Handoff
backlink: "[[3.0.0 Sprint 2026-09-28 Handoff]]"
tags: [sprint, handoff, [[2.0.0 MAgHARCM]], wave-17, "[[1.0.0 P-123]]", "[[1.0.0 P-124]]", "[[1.0.0 PRIM-21]]", "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-30]]", "[[1.0.0 P-88]]", cross-lingual-oracle, dynamic-analysis-mining, software-archaeology-gap]
---

# [[3.0.0 Sprint 2026-09-28 Handoff]]

## Outcome

Wave-17 **FIRED** (5 candidates triaged, 2 ACCEPT + 3 REJECT). Two new SLM-era mechanism papers persisted: `[[1.0.0 P-123]]` CodeChemist (ICML 2026) anchors cross-lingual functional-oracle substrate for `[[1.0.0 PRIM-21]]` Migration Strategy Selection, `[[1.0.0 PRIM-23]]` Chunked Translation, `[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection. `[[1.0.0 P-124]]` Syzygy (ICLR 2025 VerifAI Workshop) anchors dynamic-analysis property-mining substrate for `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph, `[[1.0.0 PRIM-22]]` Four Phases of Comprehension, `[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter (complementing `[[1.0.0 P-88]]` HiTyper's static TDG). All gates green (`go build`, `go vet`, `go test ./...`, `bash scripts/lint_vault.sh` exit 0).

## Foundation (closed)

- Re-read Sprint 2026-09-26 handoff + `Methodology.md` (entry point) + ADR-C-014 for grounding.
- Confirmed gate baseline (build/vet/test/lint) before research fan-out.
- Stale-directive audit verified all 7 user directives still satisfied at `a5bba8e` (8-agent graph + 2 checkpoints, abcoder-mcp default, zero `fmt.Print*`, etc.).

## Wave-17 Research

- **5 candidates triaged** in parallel (workpool):
  - `ResearchCodeReasoning` → **ACCEPT** P-123 CodeChemist (ICML 2026).
  - `ResearchASTTranslation` → **ACCEPT** P-124 Syzygy (ICLR 2025 VerifAI Workshop).
  - `ResearchMemSearcher` → **REJECT Q1** (ACL 2026 Findings, off-list venue).
  - `ResearchToolVerify` → **REJECT Q1** (arXiv-only, no venue confirmation).
  - `ResearchSoftwareArch` → **REJECT Q2/Q3** (ICML 2025 candidate is a benchmark; software archaeology SLM-era mechanism gap remains open).
- **2 P-NN notes persisted** (`P-123-CodeChemist-ICML-2026.md`, `P-124-Syzygy-ICLR2025-Workshop.md`).
- **Wave-17 candidates memo** appended at `.obsidian/MAgHARCM/research/diary/Wave-17-Candidates.md`.

## Vault Sync (Wave-17)

- `Methodology.md` — frontmatter bumped to 2026-09-28; §7 wave-17 SLM-era anchors table (6 rows: PRIM-21/23/27/9/22/30) + Wave-17 anchor list appended; §9 entry for 2026-09-28 appended; §11 extended from 2 to 4 ReasoningBank/CodeChemist/Syzygy sub-sections; §11.5 Substrate Application Matrix now 4 rows; §12 stale-directive audit rolled forward to 2026-09-28.
- `Architecture.md` — frontmatter bumped; §7 vault sync audit extended with Wave-17 row; §9 Wave-17 SLM-Era Architectural Implications added (9.1 cross-lingual oracle, 9.2 dynamic-analysis mining, 9.3 backward compatibility, 9.4 parity check).
- `primitives/Primitives-Index.md` — frontmatter bumped; PRIM-9/21/22/23/27/30 rows cross-linked to P-124/P-123; Sprint 2026-09-28 Vault Sync Audit block appended.
- `Software-Archaeology-Lineage.md` — frontmatter bumped; §8 Wave-17 SLM-Era Anchors added (8.1 P-123 CodeChemist, 8.2 P-124 Syzygy, 8.3 substrate cross-reference matrix, 8.4 deferral note, 8.5 watchlist for Wave-18).

## Codebase Refactor

- Stale-directive audit re-verified all 7 user directives still satisfied at `a5bba8e`.
- No Go code touched (Wave-17 is research-only; subsequent sprints implement the two new substrates via `internal/oracle/` and `internal/specminer/` leaf packages).
- All gates green (`go build`, `go vet`, `go test ./...` cached, `bash scripts/lint_vault.sh` exit 0; vault lint: 219 files scanned, 15 refs checked).

## Paper Sync

- `docs/.paper/refs.bib` — 2 new bib entries appended (`p123_wang_codechemist_2026`, `p124_shetty_syzygy_2025`).
- `docs/.paper/sec_method.tex` — 3 new `\cite{}` clusters: line 306 (PRIM-21 cascade) adds `p123_wang_codechemist_2026`; line 335 (PRIM-22/23/31 cluster) adds `p124_shetty_syzygy_2025`; line 471 (CPG context) adds `p124_shetty_syzygy_2025`. Translator section extended with cross-lingual-oracle paragraph + dynamic-analysis-mining sentence in CPG context.

## Open gaps (Wave-18 candidates)

- **Software archaeology mechanism gap** — no SLM-era mechanism paper found at 2025+ ICLR/ICML/NeurIPS. `[[1.0.0 P-87]]` Hou TOSEM SLR remains the literature anchor. Recommend next wave-18 scout carry explicit `program-comprehension-mechanism` query.
- **MemSearcher venue confirmation** — fires if accepted at NeurIPS 2026 main track.
- **Verified Tool Calls venue confirmation** — fires if accepted at NeurIPS 2026 / ICML 2027.
- **Implementation work** for wave-17 anchors: `internal/oracle/` leaf package for CodeChemist-style cross-lingual I/O test oracle; `internal/specminer/` leaf package for Syzygy-style Clang/LLVM dynamic-analysis property mining. Both opt-in via YAML config per Architecture §9.3.

## Files touched (this sprint)

- `.obsidian/MAgHARCM/research/papers/P-123-CodeChemist-ICML-2026.md` (new)
- `.obsidian/MAgHARCM/research/papers/P-124-Syzygy-ICLR2025-Workshop.md` (new)
- `.obsidian/MAgHARCM/research/diary/Wave-17-Candidates.md` (new)
- `.obsidian/MAgHARCM/research/Methodology.md` (frontmatter, §7 wave-17 anchors, §9 entry, §11 extension, §12 refresh)
- `.obsidian/MAgHARCM/research/Architecture.md` (frontmatter, §7 vault sync, §9 wave-17 implications)
- `.obsidian/MAgHARCM/primitives/Primitives-Index.md` (frontmatter, PRIM-9/21/22/23/27/30 rows, Sprint 2026-09-28 audit block)
- `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` (frontmatter, §8 wave-17 anchors)
- `docs/.paper/refs.bib` (2 new entries)
- `docs/.paper/sec_method.tex` (3 cite clusters extended)
- `.obsidian/MAgHARCM/diary/Sprint-2026-09-28-Handoff.md` (this file)
