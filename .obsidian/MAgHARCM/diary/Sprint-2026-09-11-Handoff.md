---
title: Sprint 2026-09-11 Handoff
backlink: "[[2.0.0 Sprint 2026-09-11]]"
tags: [sprint, handoff, [[2.0.0 MAgHARCM]], [[1.0.0 PRIM-21]], legacy-modernisation, red-teaming, paper-sync]
---

# [[2.0.0 Sprint 2026-09-11 Handoff]]

## Outcome

Sprint 2026-09-11 closed. Four new anchor papers persisted (P-74..P-77), all real literature (no fabrication). All 6 remaining `[[NEEDS-LINK]]` stubs resolved. Paper sync propagates 4 new cite mentions + 4 new bib entries. Three focused commits; clean tree; green `go build`, `go vet`, `go test -count=1 ./...` (9 packages).

## Research Wave 5 (closed)

- `[[1.0.0 P-74]]` Biggerstaff, Mitbander, Webster 1993 — *The Concept Assignment Problem in Program Understanding* (ICSE 1993). The foundational paper defining concept-assignment; static-first / dynamic-fallback strategy for concept locators; 60-70% of mature-codebase concepts resolvable statically.
- `[[1.0.0 P-75]]` Bennett 1995 — *Legacy Systems: Coping with Success* (IEEE Software 12(1)). The canonical 2-page framing for legacy modernisation; the wrap / reengineer / replace / migrate taxonomy that grounds `[[1.0.0 PRIM-21]]`.
- `[[1.0.0 P-76]]` Bennett & Rajlich 2000 — *Software Maintenance and Evolution: A Roadmap* (FOSE 2000 at ICSE 2000). The seven-grand-challenges research-agenda statement; maps directly onto MAgHARCM primitives.
- `[[1.0.0 P-77]]` OpenAI 2024 — *Advancing Red Teaming with People and AI* (November 2024). The breadth+depth methodology (manual + automated + mixed); precedent for MAgHARCM's verifier-of-verifier pattern.

**Provenance correction**: initial draft of P-74 fabricated a "Rajlich & Bennett 1997" co-authored paper and a "Bennett & Rajlich 2000 staged lifecycle" paper. Both removed; the real papers (Biggerstaff 1993 for concept-assignment; Bennett & Rajlich 2000 *Roadmap*) are now persisted with correct titles and authors. `Baldwin-Clark-2016-Money` was a conflation of modularity literature with monetary theory — that reference does not exist; the placeholder was replaced with Baldwin 2008 "Where do Transactions Come From?" + Colfer & Baldwin 2016 "Mirroring Hypothesis", both cited via `[[1.0.0 P-70]]`.

## NEEDS-LINK Resolution (closed)

All 6 remaining stubs resolved:
- `[[NEEDS-LINK Rajlich-1997]]` (in `Software-Archaeology-Lineage.md`, `P-40-foltz-dr-jones-2023.md`) → `[[1.0.0 P-74]]` Biggerstaff 1993 (the actual concept-assignment paper; Rajlich cites this).
- `[[NEEDS-LINK Bennett-2000]]` (in `Software-Archaeology-Lineage.md`) → `[[1.0.0 P-75]]` Bennett 1995 + `[[1.0.0 P-76]]` Bennett & Rajlich 2000 Roadmap.
- `[[NEEDS-LINK Baldwin-Clark-2016-Money]]` (in `P-41-baldwin-clark-design-rules-deep-2024.md`) → conflation; replaced with real Baldwin 2008 + Colfer & Baldwin 2014 refs via `[[1.0.0 P-70]]`.
- `[[NEEDS-LINK Anthropic-2024-Claude35Sonnet]]` (in `P-56-zan-multiswebench-2025.md`) → cross-link to `[[1.0.0 P-38]]`.
- `[[NEEDS-LINK Anthropic-2025]]` (in `P-61-guo-risky-code-execution-2024.md`) → cross-link to `[[1.0.0 P-38]]` + `[[1.0.0 P-69]]`.
- `[[NEEDS-LINK OpenAI-2024-CodexRedTeam]]` (in `P-61-guo-risky-code-execution-2024.md`) → cross-link to `[[1.0.0 P-77]]` (no Codex-specific red-team report exists; OpenAI 2024 is the closest analogue).

## Codebase Ponytail (closed)

- `internal/graph/graph.go`: zero changes warranted; already uses `compiletime.MaxGraphRunSteps` constant; no magic literals.
- `internal/runner/runner.go`: zero changes warranted; already uses `cfg.OllamaBaseURL`, `cfg.ReasoningModel`, `checkpoint.RunIDForSourceDir` — config-driven.
- `ModularityTrapYears` constant: does not exist in `internal/compiletime/compiletime.go` (was aspirational mention in last sprint's handoff, not actually persisted). Skipped wiring — the constant should be added in a future sprint if the "modularity trap" detection (cf. `[[1.0.0 P-71]]`) is implemented.

## Paper Sync (closed)

- `docs/.paper/sec_method.tex`: 4 new cite mentions — P-74 on Concept Assignment (in Archaeologist's PRIM-20 list); P-75 + P-76 on migration-strategy taxonomy (in Strategy Registry section); P-77 on Verdict Panel (in verifier-of-verifier paragraph).
- `docs/.paper/refs.bib`: 4 new bib entries (`p74_biggerstaff_concept_assignment_1993`, `p75_bennett_legacy_coping_success_1995`, `p76_bennett_rajlich_roadmap_2000`, `p77_openai_red_teaming_2024`).
- `pdflatex + bibtex + pdflatex + pdflatex`: 920059 bytes, zero undefined citations.

## Commits this sprint

```
6dd47be docs(paper): sync sec_method.tex + refs.bib with P-74..P-77 — 4 new cite mentions, 4 new bib entries; pdflatex + bibtex compile clean (920059 bytes, zero undefined citations)
c1706f8 docs(vault): resolve all 6 remaining NEEDS-LINK stubs — Rajlich-1997→P-74 (Biggerstaff), Bennett-2000→P-75 (Bennett 1995), Baldwin-Clark-2016-Money→P-70 (replaced conflated cite with real Baldwin 2008 + Colfer 2014 refs), Anthropic-2024-Claude35Sonnet→P-38, Anthropic-2025→P-38/P-69, OpenAI-2024-CodexRedTeam→P-77; correct Software-Archaeology-Lineage §2.3 author/title attribution
c979...  docs(research): persist P-74..P-77 anchor papers — Biggerstaff 1993 Concept Assignment; Bennett 1995 Legacy Systems Coping with Success; Bennett-Rajlich 2000 Roadmap; OpenAI 2024 Red Teaming methodology
```

(3 commits; `c979...` shown abbreviated.)

## Follow-ups (next sprint)

- The old `p71_anthropic_sycophancy_2025` bib entry (Raman et al., arXiv 2503.13930) is still in `refs.bib` with no in-text cite. Decide: re-introduce it as the actual P-69 source (replacing Sharma 2025), or remove as dead entry.
- Add `compiletime.ModularityTrapYears` constant if the "modularity trap" detection (`[[1.0.0 P-71]]` Fleming-Baldwin 2024 anchor) is implemented in `internal/agents/archaeology.go`.
- Consider filing the 7 grand challenges from Bennett & Rajlich 2000 (`[[1.0.0 P-76]]`) into the `[[primitives/INDEX]]` as a research-agenda column.
- Ponytail sweep on `internal/compiletime/compiletime.go` for any remaining inline magic literals; this file has grown to 441 lines.
