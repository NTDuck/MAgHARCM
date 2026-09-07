---
title: "ADR-2026-09-26 Vault Lint Extension — ADR-V-001 Automation"
backlink: "[[1.0.0 ADR-2026-09-26-Vault-Lint-Extension]]"
status: Accepted
date: 2026-09-26
last_updated: 2026-09-26
tags: [adr, vault-lint, automation, adr-v-001, [[2.0.0 MAgHARCM]], [[2.0.0 Architecture]], [[1.0.0 ADR-V-001]], [[1.0.0 ADR-2026-09-07-Sprint-Conventions]], slm]
---

# [[1.0.0 ADR-2026-09-26 Vault Lint Extension (ADR-V-001 automation)]]

## 1. Context

`[[1.0.0 ADR-V-001]]` (`.obsidian/MAgHARCM/architecture/ADR-2026-09-07-Sprint-Conventions.md`, §2.2) mandates that version markers MUST use `[[x.y.z ...]]` wikilink syntax; stray parentheticals such as `(P-46)` are PROHIBITED in version slots. Enforcement until Sprint 2026-09-26 was entirely human: each sprint's audit step manually grep'd `(P-NN)` / `(PRIM-NN)` matches in `internal/` and `.obsidian/MAgHARCM/`, manually inspected each match for backtick-context exemption, and manually rewrote the violations.

The Sprint 2026-09-24 audit sweep landed 17 violations to 0 across `Methodology.md`, `primitives/Primitives-Index.md`, and four diary handoffs. The sweep took one subagent ~25 minutes and is the model for future drift repair — but the model is not durable. Future drifts between sprints will reintroduce the same violation set unless the rule is automated.

The Sprint 2026-09-17 prose-convention decision ("prose parentheticals `(PRIM-NN)` / `(P-NN)` are retained as standard academic-writing convention in user-facing prose") further complicates the enforcement shape: bare `(PRIM-NN)` is permitted in vault prose but banned in `internal/` Go source. A naive lint script that flags all `(P-NN)` / `(PRIM-NN)` matches without distinguishing these contexts will produce both false positives (vault prose) and false negatives (Go source after backtick-context exemptions).

## 2. Decision

Ship `scripts/lint_vault.sh` as the durable enforcement mechanism for `[[1.0.0 ADR-V-001]]`. The lint script implements **asymmetric enforcement** aligned with the Sprint 2026-09-17 prose-convention decision:

- **`internal/` directory** — bare `(P-NN)` outside backticks is **FORBIDDEN** (exit 1). `(PRIM-NN)` parentheticals are explicitly permitted and counted but not flagged (Sprint 2026-09-17 decision).
- **`.obsidian/MAgHARCM/` directory** — bare `(PRIM-NN)` outside backticks is reported as **STRAY** (informational; exit 0 unless `FORBIDDEN` was also found). Bare `(P-NN)` in vault markdown is exempted (vault content, not code).
- **Backtick / fenced-code context** is exempt on both sides. Inline code (`` `(P-46)` ``) and fenced code blocks (```` ```\n(P-46)\n``` ````) do not trigger the rule.

### 2.1. Asymmetric enforcement rationale

The two domains have different conventions:

| Domain | Bare `(P-NN)` | Bare `(PRIM-NN)` |
| :--- | :--- | :--- |
| `internal/` (Go source) | **FORBIDDEN** (exit 1) | permitted (Sprint 2026-09-17) |
| `.obsidian/MAgHARCM/` (markdown) | exempted (vault content) | **STRAY** (report-only, exit 0) |

A single symmetric rule would either over-flag vault prose or under-flag Go source. Asymmetric enforcement matches the Sprint 2026-09-17 decision: prose parentheticals in vault markdown are normal academic-writing convention; Go source has no such convention.

### 2.2. Exit codes

- `0` — clean (zero FORBIDDEN matches; STRAY matches reported but do not fail)
- `1` — FORBIDDEN match found in `internal/`
- `2` — `internal/` or `.obsidian/MAgHARCM/` directory missing (script error, not a lint failure)

### 2.3. CI integration

The script is wired into CI as the canonical ADR-V-001 gate. Zero FORBIDDEN hits expected on green trees (the Sprint 2026-09-24 sweep landed 17 → 0 and the script preserves that invariant). The STRAY report helps reviewers spot drift that the prose-convention decision tolerates but the single-versioning convention prefers to be rewritten.

## 3. Consequences

Positive:

- **Zero STRAY hits expected on green trees.** The Sprint 2026-09-24 sweep landed 17 → 0 and the script preserves that invariant; any drift between sprints surfaces immediately.
- **Doubles as drift detector.** Future sprints that introduce new `(P-NN)` / `(PRIM-NN)` parentheticals by accident see the lint signal in CI rather than at the next audit.
- **Preserves the Sprint 2026-09-17 prose convention.** The script reports STRAY hits but does not fail on them, so vault prose can keep academic-style parentheticals where they read naturally.
- **Self-documenting.** The script's header comment restates ADR-V-001 + the Sprint 2026-09-17 decision + the exit-code contract; reviewers reading the script learn the rule without a separate rule note.

Negative:

- **Adds a CI dependency.** Pipelines must invoke `scripts/lint_vault.sh` after `go build` / `go vet` / `go test` and gate merges on its exit code. A failed gate points at `scripts/lint_vault.sh:line:col` rather than at the offending markdown file directly; the script's output format is human-readable but not editor-integration-shaped.
- **False positives on legitimate prose.** The STRAY report flags every bare `(PRIM-NN)` in vault prose as informational; reviewers must read each STRAY hit and decide whether the prose-convention exemption applies or whether the citation should be rewritten to `[[1.0.0 PRIM-NN]]` form. This is by design (the prose-convention decision tolerates both forms) but adds reviewer cognitive load.
- **No automatic rewrite.** The script reports and exits; it does not auto-fix. Future sprints may add an `--auto-fix` mode that rewrites `(PRIM-NN)` → `[[1.0.0 PRIM-NN]]` for STRAY hits, but that mode is explicitly out of scope for this ADR.

## 4. Compliance Map

| Rule | Status | Evidence |
| :--- | :--- | :--- |
| `[[1.0.0 ADR-V-001]]` enforcement | ✅ | `scripts/lint_vault.sh` |
| `[[1.0.0 ADR-2026-09-07-Sprint-Conventions]]` §2.2 rule statement | ✅ preserved | unchanged in `architecture/ADR-2026-09-07-Sprint-Conventions.md` (canonical counter-example retained) |
| Sprint 2026-09-17 prose-convention decision | ✅ preserved | asymmetric enforcement + STRAY (not FORBIDDEN) for vault `(PRIM-NN)` |
| `Methodology.md` §10 (vault sync audit trail) | ✅ | lint script emission fed into Sprint 2026-09-26 audit block |
| `scripts/` directory convention | ✅ | `scripts/lint_vault.sh` is the second automation in `scripts/` (alongside `generate-configs.py`, `retry-stable.sh`, etc.) |

## 5. Supersession

To supersede this ADR (e.g. to add an `--auto-fix` mode, to extend the rule to `.obsidian/MAgHARCM/` `(P-NN)`, or to migrate to a YAML linter pre-commit hook as flagged in `Sprint-2026-09-13-Handoff.md` §Phase-1 followups), write a new ADR referencing `supersedes: ADR-2026-09-26-Vault-Lint-Extension` and update the compliance map.

## 6. Pointers

- Lint script: `scripts/lint_vault.sh` (signed-off in Sprint 2026-09-26)
- Companion ADR: `[[1.0.0 ADR-V-001]]` in `.obsidian/MAgHARCM/architecture/ADR-2026-09-07-Sprint-Conventions.md` §2.2
- Sprint-Conventions ADR: `.obsidian/MAgHARCM/architecture/ADR-2026-09-07-Sprint-Conventions.md`
- Methodology: `.obsidian/MAgHARCM/research/Methodology.md` §10 (vault sync audit trail)
- Architecture: `.obsidian/MAgHARCM/research/Architecture.md` §7 (Sprint 2026-09-26 Vault Sync Audit)
- Primitives INDEX audit block: `.obsidian/MAgHARCM/primitives/Primitives-Index.md` §Sprint 2026-09-26
- Lineage note: `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` §6 (Wave-15 Deferral)
- Sprint handoff (current sprint): `.obsidian/MAgHARCM/diary/Sprint-2026-09-26-Handoff.md` (pending)
