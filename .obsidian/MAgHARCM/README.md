---
title: MAgHARCM Vault
backlink: "[[2.0.0 MAgHARCM Vault]]"
aliases:
  - "2.0.0 MAgHARCM Vault"
  - "MAgHARCM Vault"
  - "README"
tags: [vault, index, "[[2.0.0 MAgHARCM]]"]
last_updated: 2026-09-28
---

# [[2.0.0 MAgHARCM Vault]]

Welcome to the canonical knowledgebase for the **MAgHARCM** (Multi-Agent Hardware/Architecture Recovery & Code Modernization) platform.

All 31 primitives (`[[1.0.0 PRIM-1]]` through `[[1.0.0 PRIM-31]]`) and 124 literature papers (`[[1.0.0 P-01]]` through `[[1.0.0 P-124]]`) across 17 research waves are active and documented with 100% codebase parity.

---

## 🧭 Human Executive Reports (`adhoc/`)

For human engineering leads, researchers, and contributors, the **`adhoc/`** directory houses consolidated reports designed for quick consumption:

1. **[[adhoc/Methodology|Methodology & Architecture Hub]]** (`adhoc/Methodology.md`)
   - Complete multi-agent pipeline overview, dynamic strategy registry, validation cascade, and theoretical lineage.
2. **[[adhoc/Project-Progress-And-Milestones|Project Progress & Milestones Tracker]]** (`adhoc/Project-Progress-And-Milestones.md`)
   - High-level metrics, 31/31 primitives parity status, completed architectural milestones, and sprint velocity.
3. **[[adhoc/Strategic-Direction-And-Roadmap|Strategic Direction & Future Roadmap]]** (`adhoc/Strategic-Direction-And-Roadmap.md`)
   - Local 4B–30B SLM modernization vision, test-time reasoning scaling, and phased delivery roadmap.
4. **[[adhoc/Human-Intervention-And-Blockers|Human Intervention & Blockers Triage]]** (`adhoc/Human-Intervention-And-Blockers.md`)
   - Active triage of items needing human decisions: `memorystore` compilation break, unverified literature placeholders, and `Commons-Validator` translation plateau.
5. **[[adhoc/Benchmark-Results-And-Evaluation|Benchmark Results & Empirical Evaluation]]** (`adhoc/Benchmark-Results-And-Evaluation.md`)
   - Quantitative evaluation across 4 benchmark repositories (GildedRose 100%, Gohistogram 62.5%, Stats 57.1%, Commons-Validator 26.5%), research questions, and ablation findings.
6. **[[adhoc/Architecture-And-Dataflow|Architecture & System Dataflow]]** (`adhoc/Architecture-And-Dataflow.md`)
   - Visual execution graph of the 8 agents, shared typed pipeline state (`compiletime.State`), and Locality of Behaviour boundaries.
7. **[[adhoc/Research-Waves-Index|Research Waves Index (Waves 1–17)]]** (`adhoc/Research-Waves-Index.md`)
   - Thematic and chronological synthesis of all 17 research waves connecting 124 papers to primitives.

---

## 🤖 Machine-Readable Autoresearch Database

For autonomous LLM agents conducting research or verification:
- **`Research-Database.json`** — Unified JSON dataset containing structured metadata for all 31 primitives, 124 papers (with authors, venue, year, bibkey, DOI, hop-1/2 citations, and summaries), 8 agent pipeline stages, and research waves.

---

## 📁 Repository & Vault Structure

```
.obsidian/MAgHARCM/
├── README.md                        — Vault overview and navigation guide
├── Research-Database.json           — Machine-readable autoresearch dataset for LLM agents
├── adhoc/                           — Executive reports suite for human readers
│   ├── Methodology.md               — Primary methodology aggregation dashboard
│   ├── Project-Progress-And-Milestones.md — Velocity, milestones, and parity matrix
│   ├── Strategic-Direction-And-Roadmap.md — Strategic vision and future trajectory
│   ├── Human-Intervention-And-Blockers.md — Issues requiring human decision-making
│   ├── Benchmark-Results-And-Evaluation.md — Empirical performance & ablation data
│   ├── Architecture-And-Dataflow.md — Agent graph diagrams and state transitions
│   └── Research-Waves-Index.md      — Thematic synthesis of Waves 1 through 17
├── architecture/
│   ├── ADR-2026-09-07-Sprint-Conventions.md — Governance & conventions (ADR-C-001..015, ADR-V-001..007)
│   ├── ADR-2026-09-07-Dup-Row-Escape-Recipe.md — Lineage multi-anchor escape recipe
│   └── ADR-2026-09-26-Vault-Lint-Extension.md — ADR-V-001 automation
├── diary/
│   ├── Sprint-YYYY-MM-DD-Handoff.md — Sprint handoff records (2026-09-04 through 2026-09-28)
│   └── Sprint-2026-09-26-Charm-Audit.md — Charm stack idiomatic audit
├── primitives/
│   └── Primitives-Index.md          — Canonical catalog of 31 primitives with codebase parity
└── research/
    ├── Architecture.md              — Package DAG, 8-agent decomposition, Locality of Behaviour
    ├── Methodology.md               — Detailed methodology specification
    ├── Software-Archaeology-Lineage.md — 2-hop foundational literature synthesis (Parnas, Lehman, etc.)
    ├── diary/                       — Research wave candidate memos (Wave-15 through Wave-17)
    └── papers/                      — 124 paper notes (P-01-ReCodeAgent through P-124-Syzygy)
```

---

## ⚖️ Governance & Invariants

1. **Camel-Case (with hyphen)**: Every file name across all subdirectories uses Camel-Case with hyphen (e.g. `P-01-ReCodeAgent.md`, `Primitives-Index.md`).
2. **Version Markers (`[[x.y.z ...]]`)**: All version references must use wikilink syntax (`[[1.0.0 PRIM-NN]]`, `[[1.0.0 P-NN]]`).
3. **Quality Gate**: Changes must pass `./scripts/lint_vault.sh` cleanly with zero forbidden references.
