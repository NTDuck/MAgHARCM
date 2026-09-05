---
title: Sprint 2026-09-11 Handoff
backlink: "[[2.0.0 Sprint 2026-09-11]]"
tags: [sprint, handoff, [[2.0.0 MAgHARCM]], [[1.0.0 PRIM-19]], legacy-modernisation, red-teaming, paper-sync, modularity-trap]
---

# [[2.0.0 Sprint 2026-09-11 Handoff]]

## Outcome

Sprint 2026-09-11 closed. Four real anchor papers persisted (P-74..P-77); all 6 remaining `[[NEEDS-LINK]]` stubs resolved; paper sync clean (920121 bytes, 0 undefined citations); `compiletime.ModularityTrapYears = 15` declared + wired into the design-rule partition with a deterministic unit test; sec_method.tex line 307 framing corrected; p77 `howpublished` malformed field fixed. Seven focused commits; clean tree; green `go build`, `go vet`, `go test -count=1 ./...` (9 packages).

## Research Wave 5 (closed)

- `[[1.0.0 P-74]]` Biggerstaff, Mitbander, Webster 1993 — *The Concept Assignment Problem in Program Understanding* (ICSE 1993). Foundational concept-assignment paper; static-first / dynamic-fallback strategy for concept locators.
- `[[1.0.0 P-75]]` Bennett 1995 — *Legacy Systems: Coping with Success* (IEEE Software 12(1)). Canonical 2-page framing for legacy modernisation; informal wrap / reengineer / replace / migrate coping strategies discussed without a formal taxonomy.
- `[[1.0.0 P-76]]` Bennett & Rajlich 2000 — *Software Maintenance and Evolution: A Roadmap* (FOSE 2000 at ICSE 2000). Multi-area research-agenda statement covering concept location, impact analysis, comprehension at scale, software visualisation, reuse-driven maintenance, migration strategies, and empirical validation.
- `[[1.0.0 P-77]]` OpenAI 2024 — *Advancing Red Teaming with People and AI* (November 2024 blog post; `@misc` bib entry). Breadth+depth methodology (manual + automated + mixed); precedent for MAgHARCM's verifier-of-verifier pattern.

**Provenance correction log**: the initial draft of P-74 carried a fabricated "Rajlich & Bennett 1997" co-authored paper and a "Bennett & Rajlich 2000 staged lifecycle" paper; both removed. The real papers (Biggerstaff 1993 for concept-assignment; Bennett & Rajlich 2000 *Roadmap*) are now persisted with correct titles and authors. `Baldwin-Clark-2016-Money` was a conflation of modularity literature with monetary theory — that reference does not exist; the modular-money framing is now anchored in the existing `[[1.0.0 P-70]]` (Baldwin & Clark 2006 design-rules working paper) and `[[1.0.0 P-71]]` (Fleming & Baldwin 2024 retrospective); specific sub-citation years (Baldwin 2008 modular organisations, Colfer & Baldwin 2014 mirroring) deferred pending primary-source verification.

## NEEDS-LINK Resolution (closed)

All 6 remaining stubs resolved:
- `[[NEEDS-LINK Rajlich-1997]]` (in `Software-Archaeology-Lineage.md`, `P-40-foltz-dr-jones-2023.md`) → `[[1.0.0 P-74]]` Biggerstaff 1993 (the actual concept-assignment paper; Rajlich cites this).
- `[[NEEDS-LINK Bennett-2000]]` (in `Software-Archaeology-Lineage.md`) → `[[1.0.0 P-75]]` Bennett 1995 + `[[1.0.0 P-76]]` Bennett & Rajlich 2000 Roadmap.
- `[[NEEDS-LINK Baldwin-Clark-2016-Money]]` (in `P-41-baldwin-clark-design-rules-deep-2024.md`) → conflation dropped; modular-money framing anchored in existing `[[1.0.0 P-70]]` / `[[1.0.0 P-71]]` with sub-citation years deferred.
- `[[NEEDS-LINK Anthropic-2024-Claude35Sonnet]]` (in `P-56-zan-multiswebench-2025.md`) → inline plain-text reference + cross-link to `[[1.0.0 P-38]]`.
- `[[NEEDS-LINK Anthropic-2025]]` (in `P-61-guo-risky-code-execution-2024.md`) → cross-link to `[[1.0.0 P-38]]` + `[[1.0.0 P-69]]`.
- `[[NEEDS-LINK OpenAI-2024-CodexRedTeam]]` (in `P-61-guo-risky-code-execution-2024.md`) → cross-link to `[[1.0.0 P-77]]` (no Codex-specific red-team report exists; OpenAI 2024 *Advancing Red Teaming with People and AI* is the closest analogue).

## Codebase Ponytail (closed)

- `internal/graph/graph.go`: zero changes warranted; already uses `compiletime.MaxGraphRunSteps`; no magic literals.
- `internal/runner/runner.go`: zero changes warranted; already config-driven via `cfg.OllamaBaseURL`, `cfg.ReasoningModel`, `checkpoint.RunIDForSourceDir`.

## Codebase Wiring — Modularity Trap Detection (closed)

Vault docs (Sprint-2026-09-10 handoff line 23 + P-71 paper file line 35) both promised `compiletime.ModularityTrapYears` but the constant was never declared. Resolved this sprint:

- `internal/compiletime/compiletime.go`: new `ModularityTrapYears = 15` constant with provenance comment pointing to `[[1.0.0 P-71]]` (Fleming & Baldwin 2024 retrospective on Design Rules).
- `internal/agents/design_rule_hierarchy.go`: `PartitionedElement` extended with `StableSinceYear *int`, `IsTrap bool`, `TrapReason string` (all `omitempty` JSON). New `PartitionWithStableSince(elements, deps, stableSinceYears map[string]int)` overload flags L1 design rules older than the threshold as modularity traps; original `Partition()` delegates with `nil` map (no-op). `markModularityTrap` helper computes `age = currentYear() - stableSinceYear`; `currentYear` is a package-private seam for test stubbing.
- `internal/agents/design_rule_hierarchy_test.go`: 2 unit tests. `TestPartition_L1FlaggedAsModularityTrap` verifies an ancient root is flagged while a recent root is not. `TestPartition_NilStableSinceYearsLeavesNoTraps` verifies the `Partition()` convenience overload leaves traps unset when no stable-since data is supplied.
- Callers (e.g. the Archaeologist agent's design-rule partition step) can now read git blame or file mtime to populate `stableSinceYears` and surface flagged L1 elements to the user for explicit re-validation, matching the P-71 surface semantics.

## Paper Sync (closed)

- `docs/.paper/sec_method.tex`: 4 new cite mentions. Line 307 framing corrected — dropped "canonical" (Bennett 1995 doesn't lay out a numbered taxonomy) + duplicate "wrap" copy-paste artifact; replaced with "informal wrap/reengineer/replace/migrate coping-strategy framing discussed in Bennett's 1995 editorial".
- `docs/.paper/refs.bib`: 4 new bib entries. `p77_openai_red_teaming_2024` malformed `howpublished` (text-prefix concatenated with `\url{}`) fixed — moved to `note` field where `\url` renders correctly as a hyperlink.
- `pdflatex + bibtex + pdflatex + pdflatex`: 920121 bytes, zero undefined citations.

## Commits this sprint

```
45997d9 feat(archaeology): declare compiletime.ModularityTrapYears=15 + wire into design-rule partition
bc0fe3a docs(paper): correct sec_method.tex line 307 wrap/reengineer/replace/migrate framing + fix p77 howpublished -> note
d5983c7 docs(vault): strip unverified provenance from P-41, P-76, Sprint-2026-09-11 handoff
1ce3d02 docs(diary): Sprint 2026-09-11 handoff - 4 anchor papers; 6 NEEDS-LINK stubs resolved; paper sync; graph+runner ponytail audited clean; provenance correction logged
c1706f8 docs(vault): resolve all 6 remaining NEEDS-LINK stubs
6dd47be docs(paper): sync sec_method.tex + refs.bib with P-74..P-77
b243650 docs(research): persist P-74..P-77 anchor papers
```

(7 commits; oldest→newest.)

## Follow-ups (next sprint)

- **Primary** (the user's actual codebase brief still largely unaddressed): implement the `MigrationStrategy` interface registry as a real try-and-fail cascade. `internal/agents/strategy.go` already has the `MigrationStrategy` interface (`Kind()`, `Matches()`, `Attempt()`) + `Registry.TryInOrder()` + `Registry.NextStrategy()`, but the "try-and-fail incrementally" semantics the user asked for need: (a) explicit state-machine wiring so a failure of one strategy records the failure context and feeds it forward to the next; (b) richer `Profile` fields beyond `FileCount / LoC / HasTests / HasBuild` (e.g. cyclomatic complexity, language distribution, build-tool chain); (c) telemetry around which strategies get tried in real benchmarks.
- **SelectMigrationStrategy vs MigrationStrategy naming**: `SwitchToNextStrategy` already exists in `internal/agents/`; rename + centralise as `compiletime.MigrationStrategy` (currently an enum) and ensure the strategy selection is one canonical path, not duplicated across agents.
- The old `p71_anthropic_sycophancy_2025` bib entry (Raman et al., arXiv 2503.13930) is still in `refs.bib` with no in-text cite. Decide: re-introduce it as the actual P-69 source (replacing Sharma 2025), or remove as dead entry.
- File the research-agenda grand challenges from Bennett & Rajlich 2000 (`[[1.0.0 P-76]]`) into the `[[primitives/INDEX]]` as a research-agenda column.
- Ponytail sweep on `internal/compiletime/compiletime.go` for any remaining inline magic literals; this file has grown to 447 lines.
- Wire the Archaeologist agent to populate `stableSinceYears` from git blame / file mtime so the new modularity-trap detection actually surfaces in user-facing output.
- Verify the two deferred sub-citation claims (Baldwin "Where Do Transactions Come From?" year; Colfer & Baldwin "Mirroring Hypothesis" year) via web_search; fix P-41 line 25 if either year is wrong.
