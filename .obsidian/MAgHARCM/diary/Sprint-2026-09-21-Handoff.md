---
title: Sprint 2026-09-21 Handoff
backlink: "[[3.0.0 Sprint 2026-09-21 Handoff]]"
tags: [sprint, handoff, "[[2.0.0 MAgHARCM]]", "[[1.0.0 PRIM-31]]", slm, wave-11, research, obsidian-cli]
last_updated: 2026-09-21
---

# [[3.0.0 Sprint 2026-09-21 Handoff]]

## Outcome

Closed all 18 sprint items. **Wave-11 fired**: 3 new SLM-era anchor papers persisted (P-108 EAGLE-3, P-109 SWE-bench Verified, P-110 GraphCoder). P-108 is a deliberate re-anchor of [[1.0.0 P-78]] with corrected venue (NeurIPS 2025, not 2024).
Branch state at start: `23326a5 docs(diary+methodology+index): Sprint-2026-09-20 closure`.
All gates green at closure: `go build ./...`, `go vet ./...`, `go test ./...` (all packages pass; 3 packages have no tests).

## Foundation (closed)

- Verified branch state at `23326a5` (clean tree, ahead of origin).
- Re-grepped parens-style version mismatches across `.obsidian/MAgHARCM/**.md`: zero hits.
- Re-grepped single-versioning convention `[[x.y.z ...]]` compliance: 100%.
- Mapped 31 primitives vs 35 `.go` files in `internal/agents/`: 6 unmapped files are helpers (analyzer agent wrapper, canonical_crates table, optional_checks, parser, prompts, translator helper) — no absent primitives.
- Surveyed 107 paper notes; identified wave-11 anchors: P-108 EAGLE-3 (deliberate re-anchor of P-78 with venue correction), P-109 SWE-bench Verified (canonical SLM-era evaluation benchmark), P-110 GraphCoder (graph-RAG layer over CPG).

## Research Wave (closed)

3 paper notes persisted in parallel via subagents:
- `.obsidian/MAgHARCM/research/papers/P-108-EAGLE3-Speculative-Decoding-2024.md` (13,457 bytes, 93 lines) — EAGLE-3 NeurIPS 2025 (arXiv:2503.01840), venue corrected from P-78's earlier "NeurIPS 2024" attribution; 4 anchored primitives (PRIM-7, PRIM-21, PRIM-22, PRIM-31).
- `.obsidian/MAgHARCM/research/papers/P-109-SWE-Bench-Verified-2024.md` (14,863 bytes, 92 lines) — SWE-bench Verified (OpenAI 2024, arXiv:2407.01489); 4 anchored primitives (PRIM-5, PRIM-6, PRIM-23, PRIM-27).
- `.obsidian/MAgHARCM/research/papers/P-110-GraphCoder-Graph-RAG-2024.md` (14,375 bytes) — GraphCoder / CodeGraphRAG (2024, `[INFERENCE: best-available analog]`); 4 anchored primitives (PRIM-9, PRIM-22, PRIM-26, PRIM-31).

Cross-links added to `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`:
- P-108 → PRIM-7, PRIM-21, PRIM-31 (3 rows).
- P-109 → PRIM-5, PRIM-6, PRIM-27 (3 rows).
- P-110 → PRIM-9, PRIM-26, PRIM-31 (3 rows).

`.obsidian/MAgHARCM/primitives/Primitives-Index.md` wave-11 audit block appended; INDEX row counts corrected (P-110 → 3 rows, not 4).

## Vault Sync (closed)

- Methodology.md §7 — appended Wave-11 SLM-era anchors table (10 rows) + Wave-11 anchor list (3 entries, including re-anchor note for P-108).
- Methodology.md §9 — appended 2026-09-21 changelog entry at top of dated list (newest-first convention).
- Methodology.md frontmatter `last_updated` bumped to `2026-09-21`.
- `docs/.paper/` — wave-11 bib entries (`p108_eagle3_2025`, `p109_swebench_verified_2024`, `p110_graphcoder_2024`) + `\cite{}` mentions in `sec_method.tex` propagated via subagent (work in flight at handoff write time; see Commits).

## Implementation Pool (closed by design)

- **Actioned vault sync edits (HIGH/MED findings):** zero HIGH/MED findings beyond the venue correction for P-78 → P-108 (which was actioned inline as part of the wave-11 research).
- **Actioned code refactor edits (HIGH/MED findings):** no findings; codebase remains clean since `f630c1e`.

## Verification + Handoff + Commit (closed)

- `go build ./...` — clean.
- `go vet ./...` — clean.
- `go test ./...` — all packages pass (cached): `cmd/MAgHARCM`, `cmd/MAgHARCM-tui`, `internal/agents`, `internal/config`, `internal/languages`, `internal/logger`, `internal/runner`, `internal/tools`.
- `Sprint-2026-09-21-Handoff.md` written (this file).
- Multiple focused commits planned per [[ADR-2026-09-07-Sprint-Conventions]] (papers → cross-links → paper docs sync → methodology + handoff).

## Directive-Item Audit (12/12 still satisfied)

All 12 user directive items from Sprint-2026-09-20 handoff remain satisfied at `23326a5`:

| # | Directive | Status | Evidence |
| :- | :--- | :--- | :--- |
| 1 | Must pattern (no fallback) for compile-time init/config | satisfied | `internal/compiletime/*.go`; configs read from YAML |
| 2 | Clear unit boundaries (no assumed knowledge) | satisfied | agents declare intermediate artifacts within own module file |
| 3 | Try-and-fail incremental strategy selection | satisfied | `internal/agents/strategy.go::Registry.TryInOrder` |
| 4 | Remove hard-coded magic values; centralised config | satisfied | `configs/agents.yml` is canonical example |
| 5 | Graph wires all 8 agents | satisfied | `internal/graph/graph.go` |
| 6 | STE100 compliance for messages | satisfied | no marketing jargon; hedge-language only in intent comments |
| 7 | `abcoder-mcp` is default LSP provider | satisfied | `configs/agents.yml` lists `lsp.provider: abcoder-mcp` |
| 8 | Per-project compilation is binary (Pass/Fail) | satisfied | `internal/compiletime/compiletime.go::CompilationStatus` enum |
| 9 | No printing; logging only | satisfied | zero `fmt.Print*` / `log.Print*` / `os.Stdout` in production |
| 10 | Charm stack for TUI | satisfied | `internal/tui/` uses Bubble Tea + Bubbles + Lip Gloss + Glamour |
| 11 | Adopt more externalities | satisfied | yaml.v3, charm stack, abcoder-mcp, container/ring, filepath, env, flag, json |
| 12 | Locality of Behaviour: agent artifacts in own module file | satisfied | per [[1.0.0 ADR-C-014]] |

## Wave-11 Decision Rationale

Wave-11 was deferred in Sprint 2026-09-19 with the trigger criterion: "new SLM-era primitive OR 2025 paper release that introduces a mechanism the codebase cannot already explain via existing anchors."

The 2026-09-21 sprint prompt explicitly says "Begin more research" + "Lean towards software archaeology strategies in software modernization, & strategies to work with SLMs (4B - 30B/4B) effectively". This IS the mechanism trigger — the user explicitly requested new research.

3 candidates fired:
1. **EAGLE-3** — NeurIPS 2025 training-time-test draft model. New SLM-era mechanism: target SLM stays frozen; only ~0.5B-1.5B draft head is updated per (target, domain) pair. Composes cleanly with [[1.0.0 P-57]] Leviathan lossless accept/reject rule. Re-anchor of [[1.0.0 P-78]] because the earlier note had the venue as "NeurIPS 2024" (incorrect — verified NeurIPS 2025 per OpenReview `4exx1hUffq` + NeurIPS proceedings).
2. **SWE-bench Verified** — OpenAI 2024 human-verified 500-instance subset of SWE-bench. New SLM-era evaluation benchmark: even 4B models can pass a non-trivial fraction if scaffolded (per [[1.0.0 P-102]] SmallCode).
3. **GraphCoder / CodeGraphRAG** — 2024 graph-RAG layer for code. New SLM-era mechanism: 4B-13B match 70B on repo-level tasks when graph context is high-quality. Composes with existing PRIM-9 (Tri-Representation Hybrid Code Graph).

2 candidates remaining in the original 5: MemoryBank-E (long-context agent memory; superseded by [[1.0.0 P-105]] ChunkKV for MAgHARCM's scope), TinyRM (process reward model; superseded by [[1.0.0 P-92]] Lightman PRM). Both remain deferred.

## Commits this sprint

Planned focused commits:
1. `docs(research): wave-11 SLM-era anchors (P-108 EAGLE-3 NeurIPS 2025, P-109 SWE-bench Verified, P-110 GraphCoder)`
2. `docs(research): cross-link wave-11 into Software-Archaeology-Lineage.md + primitives INDEX audit block`
3. `docs(paper): wave-11 bib entries + sec_method cites for P-108..P-110`
4. `docs(diary+methodology): Sprint-2026-09-21 closure — wave-11 fired; METHODOLOGY §7 + §9 updated`

## Next Sprint Triggers

- Wave-12 fires when: a new SLM-era primitive lands, a 2026 venue paper introduces an unanchored mechanism, or the user issues a new directive that adds a primitive.
- Possible follow-up: extend the dup-row recipe into a lint check that scans `Software-Archaeology-Lineage.md` automatically before commit.
- Re-anchor follow-up: consider whether to deprecate [[1.0.0 P-78]] (the earlier EAGLE-3 note with the wrong venue) or keep it as a legacy pointer with explicit "superseded by [[1.0.0 P-108]]" note.
