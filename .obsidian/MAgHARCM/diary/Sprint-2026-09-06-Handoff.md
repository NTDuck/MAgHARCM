---
title: Sprint 2026-09-06 Handoff — Stale-Directive Resolution + P-46..P-49 Synthesis Anchors
backlink: [[2.0.0 Sprint-2026-09-06]]
tags: [sprint, handoff, verification, lineage, [[2.0.0 MAgHARCM]]]
---

# [[2.0.0 Sprint 2026-09-06 Handoff — Stale-Directive Resolution + P-46..P-49 Synthesis Anchors]]

Continuation of the modernization track from `[[2.0.0 Sprint-2026-09-05-Handoff]]`. Read that file first; current state continues from §1 of that handoff.

## 1. Outcome Summary

This sprint focused on three concrete deliverables, with several prescriptive directives from prior sprints verified as already-satisfied rather than re-implemented:

1. **Stale directive resolution.** Three user directives from earlier sprints were verified against current code state and found already satisfied:
   - **Try-and-fail migration strategy** — already implemented as `Registry.TryInOrder` in `internal/agents/strategy.go` (proper Matches→Attempt loop, Miss → next strategy).
   - **Graph with all 8 agents** — already wired at `internal/graph/graph.go:32-39` (archaeologist, analyzer, planning, translator, roleflip/reviewer, validator, verdict_panel, recruiter).
   - **Charm Bubble Tea TUI** — already adopted; `internal/tui/tui.go` is a full Bubble Tea Model/Update/View with `textinput.Blink`, `spinner.Tick`, Lip Gloss styles, and an idiomatic Charm `table` in `RenderConfigTable`.
   - **ABCoderMcpProvider rename** — already completed in `internal/tools/lsp_provider.go`; the stale "rename deferred" note in `Sprint-2026-09-05-Handoff.md` was updated to ✅ Completed.

2. **Four new synthesis-anchor papers persisted.** `[[1.0.0 P-46]]` (Seacord et al. 2003 *Modernizing Legacy Systems* — anchors PRIM-14), `[[1.0.0 P-47]]` (Gall Hajek Jazayeri 1998 *Logical Coupling* — anchors PRIM-18), `[[1.0.0 P-48]]` (Jia & Harman 2011 *Mutation Testing Survey* — anchors PRIM-13), `[[1.0.0 P-49]]` (Shehory & Kraus 1998 *Coalition Formation* — anchors PRIM-29).

3. **Vault hygiene verified.** No stray parentheses in version markers. Single-source-of-truth parity confirmed: 31/31 primitives in both `INDEX.md` and `Software-Archaeology-Lineage.md`. All 8 P-42..P-49 cross-refs wired in both files.

## 2. Codebase State

### 2.1. Build & Test
```
$ go build ./...          # clean
$ go vet ./...            # clean
$ go test -count=1 ./...  # 9/9 packages pass
```

### 2.2. fmt.Print* Audit (sweep result)
Zero replacements needed in `internal/` or `cmd/`. All `fmt.Fprintf` matches write into `strings.Builder` / `bytes.Buffer` for prompt assembly (correct usage); no log-line `fmt.Print*` calls in non-agent modules remain.

### 2.3. ASD-STE100 Messaging Residue
Subagent flagged no critical STE100 violations in agent file log messages. Findings (if any) listed in the subagent report; fix-up deferred to a focused STE100 sweep in next sprint.

### 2.4. Files Touched (Codebase)
None. All prescriptive items from prior sprints were already implemented. Only vault files changed.

## 3. Research State

### 3.1. New Papers
| ID | Authors | Year | Anchors |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 P-46]]` | Seacord, Comella-Dorda, Lewis, Place, Plakosh (SEI) | 2003 (book) / 2001 (TR) | PRIM-14 (Software-Archaeology Stage) |
| `[[1.0.0 P-47]]` | Gall, Hajek, Jazayeri (TU Vienna) | 1998 | PRIM-18 (Jaccard-Coupling Architecture Recovery) |
| `[[1.0.0 P-48]]` | Jia, Harman (UCL CREST) | 2011 | PRIM-13 (Adversarial Test-Weakening Guard) |
| `[[1.0.0 P-49]]` | Shehory, Kraus | 1998 | PRIM-29 (Recruitment-Adaptive Planning) |

### 3.2. Lineage + INDEX Cross-References
| File | New cross-refs |
| :--- | :--- |
| `.obsidian/MAgHARCM/primitives/INDEX.md` | PRIM-13 + P-48, PRIM-14 + P-46, PRIM-18 + P-47, PRIM-29 + P-49 |
| `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` | PRIM-13 + P-48, PRIM-14 + P-46, PRIM-18 + P-47, PRIM-29 + P-49 |

All 31/31 primitives present in both matrix tables; all 8 P-42..P-49 cross-refs verified (P-42 appears twice: PRIM-3 + PRIM-19).

### 3.3. Vault Hygiene
- Stray-parentheses sweep: 0 findings.
- Primitive parity: 31/31 in INDEX ↔ 31/31 in lineage ↔ 31/31 in codebase.
- Versioning convention: `[[x.y.z ...]]` only.

## 4. Diary Updates

- `Sprint-2026-09-04-Handoff.md` §2.1: row 5 alias renamed `reviewer` → `roleflip` (matches graph.go:var name `roleflipAgent`).
- `Sprint-2026-09-05-Handoff.md` §4.2: `ABCoderMCPProvider rename` open-work bullet updated from "Deferred" to "✅ Completed — internal/tools/lsp_provider.go and all callers now use ABCoderMcpProvider throughout".

## 5. Open Work / Follow-Ups (Deferred)

1. **`internal/agents/state.go` artifact split (Locality of Behaviour).** Still pending. 74 references across 14 files — high blast-radius. The current 28-line `State` struct (with `SchemaVersioned` interface) is minimal but artifact structs (`AnalyzerOutput`, `PlanningOutput`, etc.) live there instead of in their producer agent files. Should be split in a focused Ponytail pass.
2. **ASD-STE100 focused sweep.** Subagent flagged residue but didn't fix; needs a follow-up subagent to apply STE100-compliant rewrites to flagged log strings.
3. **Paper rewrite (sec_method.tex).** Methodology is stable (no changes this sprint); paper does not need rewriting yet. Will require update when `state.go` split lands, or if Charm stack adoption changes the TUI section narrative.
4. **Per-project compilation status documentation.** Already enforced (Pass/Fail binary); needs a paragraph update in `docs/.paper/sec_method.tex` referencing `compiletime.CompilationPass/CompilationFail`.

## 6. Verification Commands

```
$ go build ./...          # clean
$ go vet ./...            # clean
$ go test -count=1 ./...  # 9/9 packages pass
$ git status --short      # clean
```

## 7. Pointers

- Previous sprint handoff: `.obsidian/MAgHARCM/diary/Sprint-2026-09-05-Handoff.md`
- Previous sprint handoff: `.obsidian/MAgHARCM/diary/Sprint-2026-09-04-Handoff.md`
- Recon note: `.obsidian/MAgHARCM/diary/Sprint-Recon-2026-09-04.md`
- New papers: `.obsidian/MAgHARCM/research/papers/P-46` .. `P-49`
- Lineage: `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`
- Primitives catalog: `.obsidian/MAgHARCM/primitives/INDEX.md`
