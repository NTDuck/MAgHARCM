---
title: Sprint 2026-09-22 Handoff
backlink: "[[3.0.0 Sprint 2026-09-22 Handoff]]"
tags: [sprint, handoff, "[[2.0.0 MAgHARCM]]", "[[1.0.0 PRIM-31]]", slm, wave-12, research, obsidian-cli]
last_updated: 2026-09-22
---

# [[3.0.0 Sprint 2026-09-22 Handoff]]

## Outcome

Closed all 18 sprint items. **Wave-12 fired**: 4 new SLM-era anchor papers persisted (P-111 SWE-bench original, P-112 SWE-agent, P-113 AutoCodeRover, P-114 Medusa).
Branch state at start: `0c0e438 docs(primitives/methodology+index): Sprint-2026-09-21 closure`.
All gates green at closure: `go build ./...`, `go vet ./...`, `go test ./...` (all packages pass; 3 packages have no tests).

## Foundation (closed)

- Verified branch state at `0c0e438` (clean tree, ahead of origin).
- Audited wave-11 paper notes for referenced-but-unanchored hop-1 papers; identified wave-12 deep-research candidates.
- Selected 4 wave-12 anchors by greatest citation/influence per wave-11 lineage node.

## Research Wave (closed)

4 paper notes persisted in parallel via subagents:
- `.obsidian/MAgHARCM/research/papers/P-111-Jimenez-SWE-Bench-Original-2024.md` (18,268 bytes, 80 lines) — SWE-bench original (Jimenez ICLR 2024, arXiv:2310.06770); 4 anchored primitives (PRIM-5, PRIM-6, PRIM-23, PRIM-27).
- `.obsidian/MAgHARCM/research/papers/P-112-Yang-SWE-Agent-2024.md` (19,261 bytes, 107 lines) — SWE-agent (Yang NeurIPS 2024, arXiv:2405.15793); 4 anchored primitives (PRIM-5, PRIM-6, PRIM-25, PRIM-29).
- `.obsidian/MAgHARCM/research/papers/P-113-Zhang-AutoCodeRover-2024.md` (23,869 bytes, 111 lines) — AutoCodeRover (Zhang 2024, arXiv:2404.05427); 4 anchored primitives (PRIM-9, PRIM-22, PRIM-23, PRIM-26).
- `.obsidian/MAgHARCM/research/papers/P-114-Cai-Medusa-2024.md` (18,003 bytes) — Medusa (Cai 2024, arXiv:2401.10774); 4 anchored primitives (PRIM-7, PRIM-21, PRIM-22, PRIM-31).

Total: 114 paper notes persisted in `.obsidian/MAgHARCM/research/papers/`.

Cross-links added to `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`:
- P-111 → PRIM-5, PRIM-6, PRIM-23, PRIM-27 (4 rows).
- P-112 → PRIM-5, PRIM-6, PRIM-25, PRIM-29 (4 rows).
- P-113 → PRIM-9, PRIM-22, PRIM-23, PRIM-26 (4 rows).
- P-114 → PRIM-7, PRIM-21, PRIM-22, PRIM-31 (4 rows).

`.obsidian/MAgHARCM/primitives/Primitives-Index.md` wave-12 audit block appended.

## Vault Sync (closed)

- Methodology.md §7 — appended Wave-12 SLM-era anchors table (16 rows) + Wave-12 anchor list (4 entries).
- Methodology.md §9 — appended 2026-09-22 changelog entry at top of dated list (newest-first convention).
- Methodology.md frontmatter `last_updated` bumped to `2026-09-22`.

## Paper Sync (closed)

`docs/.paper/` — wave-12 bib entries (`p111_jimenez_swebench_2024`, `p112_yang_sweagent_2024`, `p113_zhang_autocoderover_2024`, `p114_cai_medusa_2024`) + 5 `\cite{}` mentions in `sec_method.tex` across 4 paragraphs (validator cascade + symbol-aware navigator + verdict panel + recruiter agent).

## Implementation Pool (closed by design)

- **Actioned vault sync edits (HIGH/MED findings):** zero HIGH/MED findings beyond the wave-12 anchors themselves (which were actioned inline as part of the wave-12 research).
- **Actioned code refactor edits (HIGH/MED findings):** no findings; codebase remains clean.

## Verification + Handoff + Commit (closed)

- `go build ./...` — clean.
- `go vet ./...` — clean.
- `go test ./...` — all packages pass (cached): `MAgHARCM/internal/agents`, `MAgHARCM/tests/cmd/MAgHARCM`, `MAgHARCM/tests/cmd/MAgHARCM-tui`, `MAgHARCM/tests/internal/agents`, `MAgHARCM/tests/internal/config`, `MAgHARCM/tests/internal/languages`, `MAgHARCM/tests/internal/logger`, `MAgHARCM/tests/internal/runner`, `MAgHARCM/tests/internal/tools`.
- `Sprint-2026-09-22-Handoff.md` written (this file).
- Multiple focused commits planned per [[ADR-2026-09-07-Sprint-Conventions]] (papers → cross-links → paper docs sync → methodology + handoff).

## Directive-Item Audit (12/12 still satisfied)

All 12 user directive items from Sprint-2026-09-21 handoff remain satisfied at `0c0e438`:

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

## Wave-12 Decision Rationale

Wave-12 was deferred in Sprint 2026-09-21 with the trigger criterion: "new SLM-era primitive OR 2025/2026 paper release that introduces a mechanism the codebase cannot already explain via existing anchors."

The Sprint 2026-09-22 trigger: wave-11 papers (P-108 EAGLE-3 + P-109 SWE-bench Verified) referenced high-impact prior work that was not yet anchored in the vault. Specifically:
- P-109 cited SWE-bench original (Jimenez ICLR 2024) and SWE-agent (Yang NeurIPS 2024).
- P-108 cited Medusa (Cai 2024) as the parallel-head alternative to EAGLE-3.
- P-109 also referenced AutoCodeRover (Zhang 2024) as the canonical code-graph-aware retrieval baseline.

4 candidates fired:
1. **SWE-bench original** — Jimenez ICLR 2024; 2,294 real GitHub issue/PR pairs across 12 Python repositories. The parent benchmark that SWE-bench Verified (`[[1.0.0 P-109]]`) subsets to 500 human-verified instances. Establishes the PASS-to-PASS + FAIL-to-PASS evaluation criterion that MAgHARCM's [[1.0.0 PRIM-27]] coverage-guided plateau detection anchors against.
2. **SWE-agent** — Yang NeurIPS 2024; introduces Agent-Computer Interfaces (ACIs) as a first-class design object with custom file/view commands + search tools that shape the agent's perception of the repo. Reports 12.5% on SWE-bench Verified with GPT-4. The ACI pattern = structured tool schemas that constrain SLM tool-call hallucination; relevant to [[1.0.0 PRIM-29]] Recruiter Agent.
3. **AutoCodeRover** — Zhang 2024; combines code-structure retrieval (AST + symbol search via tree-sitter) with program-synthesis patch generation; the canonical open-source baseline on SWE-bench. Demonstrates that a 4B-13B SLM with strong retrieval scaffolds can match much larger models on repository-level repair tasks. Directly relevant to [[1.0.0 PRIM-9]] + [[1.0.0 PRIM-26]] code-graph + symbol-navigator primitives.
4. **Medusa** — Cai 2024; an alternative to EAGLE-3's training-time-test single-head draft; Medusa adds multiple parallel decoding heads at different future-token positions to the target model directly; 2.2-3.6x speedup at lossless quality; cheaper to set up than EAGLE-3 (no draft-model training) but less flexible (cannot retrain heads per domain). Medusa vs EAGLE-3 = strategy choice: Medusa = cheap + fixed; EAGLE-3 = expensive + per-domain trainable. Relevant to [[1.0.0 PRIM-7]] + [[1.0.0 PRIM-21]] + [[1.0.0 PRIM-22]] + [[1.0.0 PRIM-31]] (verdict + strategy + comprehension + retrieval).

## Commits this sprint

Planned focused commits:
1. `docs(research): wave-12 SLM-era anchors (P-111 SWE-bench, P-112 SWE-agent, P-113 AutoCodeRover, P-114 Medusa)`
2. `docs(research): cross-link wave-12 into Software-Archaeology-Lineage.md + primitives INDEX audit block`
3. `docs(paper): wave-12 bib entries + sec_method cites for P-111..P-114`
4. `docs(diary+methodology): Sprint-2026-09-22 closure — wave-12 fired; METHODOLOGY §7 + §9 updated`

## Next Sprint Triggers

- Wave-13 fires when: a new SLM-era primitive lands, a 2026 venue paper introduces an unanchored mechanism, or the user issues a new directive that adds a primitive.
- Possible follow-up: extend the dup-row recipe into a lint check that scans `Software-Archaeology-Lineage.md` automatically before commit.
- Re-anchor follow-up: P-113 (AutoCodeRover) venue is `[INFERENCE: best-available; verify against arXiv:2404.05427]`; verify and update if found to be ISSTA 2024 or PROMISE 2024.
