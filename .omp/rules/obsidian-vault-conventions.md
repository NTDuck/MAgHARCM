---
name: obsidian-vault-conventions
description: "Rules and workflow for the MAgHARCM Obsidian vault: adhoc human reports suite, machine-readable autoresearch dataset, Camel-Case (with hyphen) naming, parity, and lint verification"
scope: ["text", "tool"]
---

# MAgHARCM Obsidian Vault Conventions & Workflow

The canonical Obsidian vault lives at `.obsidian/MAgHARCM/`. All future sessions working on research, documentation, primitives, or sprint handoffs MUST comply with these rules and workflows:

## 1. Human-Facing Interface: The `adhoc/` Suite

For human engineers, leads, and researchers, the **`adhoc/`** directory is the primary consumption surface. Future sessions MUST keep these living documents synchronized with codebase and research changes:

1. **`adhoc/Methodology.md`** — Canonical methodology aggregation dashboard: 8-agent pipeline, strategy registry, validation cascade, and theoretical lineage.
2. **`adhoc/Project-Progress-And-Milestones.md`** — Progress tracker: 31/31 primitive parity, completed milestones, sprint velocity, and convergence timeline.
3. **`adhoc/Strategic-Direction-And-Roadmap.md`** — Long-term strategy: local 4B–30B SLM modernization, test-time scaling, and phased delivery roadmap.
4. **`adhoc/Human-Intervention-And-Blockers.md`** — Active triage of issues requiring human decisions: build breaks, unverified literature placeholders, translation plateaus, and hardware/API configurations.
5. **`adhoc/Benchmark-Results-And-Evaluation.md`** — Empirical performance: pass rates across GildedRose, Gohistogram, Stats, and Commons-Validator; ablation studies.
6. **`adhoc/Architecture-And-Dataflow.md`** — Visual execution graph of the 8 agents, shared typed pipeline state (`compiletime.State`), and Locality of Behaviour boundaries.
7. **`adhoc/Research-Waves-Index.md`** — Thematic and chronological synthesis connecting 124 literature papers across Waves 1–17 to primitives.

## 2. LLM Autoresearch Workflow: `Research-Database.json`

- **Machine-First Retrieval**: Autonomous agents conducting research waves or compliance audits MUST query `.obsidian/MAgHARCM/Research-Database.json` instead of scanning 124 individual markdown notes.
- **Bi-Directional Sync**: When a new research wave or primitive lands:
  1. Add/update the structured entry in `Research-Database.json`.
  2. Persist the detailed paper note in `research/papers/P-NN-<SystemOrAuthor>.md` with full YAML aliases.
  3. Update `primitives/Primitives-Index.md` and `research/Software-Archaeology-Lineage.md`.
  4. Update the relevant `adhoc/` reports (`Project-Progress-And-Milestones.md`, `Research-Waves-Index.md`, etc.).

## 3. Vault Structure & Roles

- `README.md` — Vault welcome entrance and navigation guide (routes humans to `adhoc/` and agents to `Research-Database.json`).
- `adhoc/` — Living human executive reports suite.
- `architecture/` — Architectural Decision Records (`ADR-YYYY-MM-DD-<Slug>.md`).
- `diary/` — Sprint tracking handoffs (`Sprint-YYYY-MM-DD-Handoff.md`).
- `primitives/` — `Primitives-Index.md` (single source of truth for the 31 primitives).
- `research/` — In-depth architectural/methodological specifications (`Architecture.md`, `Methodology.md`, `Software-Archaeology-Lineage.md`, `papers/`, `diary/`).

## 4. File Naming Rules

- Root `README.md` uses uppercase per standard convention.
- EVERY other file across `adhoc/`, `architecture/`, `diary/`, `primitives/`, `research/` MUST use **Camel-Case (with hyphen)** (e.g. `Word-Word-Word.md`).
- Proper nouns, author surnames, benchmark names, and system acronyms MUST preserve canonical capitalization (e.g. `P-01-ReCodeAgent.md`, `P-21-Qwen2.5-Coder.md`, `P-109-SWE-Bench-Verified-2024.md`, `Primitives-Index.md`, `Sprint-2026-09-07-Handoff.md`).
- NO lowercase-only filenames, NO spaces, and NO underscores in filenames.

## 5. Parity & Conflict Prevention

- All primitives in `primitives/Primitives-Index.md` MUST strictly agree with `research/Software-Archaeology-Lineage.md`, `adhoc/Methodology.md`, and the Go implementations in `internal/`.
- Every paper note MUST declare `aliases:` in its YAML frontmatter (including `1.0.0 P-NN`, `P-NN`, and canonical name) so wikilinks resolve cleanly in Obsidian.

## 6. Verification Quality Gate

- After any edits to the vault or code references, MUST execute `./scripts/lint_vault.sh` from the repo root.
- The lint check MUST exit with status 0 (zero forbidden references in `internal/`).
