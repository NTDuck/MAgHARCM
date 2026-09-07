---
description: "Execute an autonomous research sprint, literature expansion, vault synchronization, empirical evaluation, codebase refactor, and academic paper update"
---

# MAgHARCM Autonomous Research Sprint & Modernization Workflow

Execute an iteration/sprint of the MAgHARCM research and codebase lifecycle.
Optional focus argument: `$ARGUMENTS`

---

## Phase 0: Sprint Inception & Context Ingestion

1. **Sprint Naming & Date Alignment**:
   - Determine today's actual date as `YYYY-MM-DD` from the system reminder or `date -u`.
   - Never extrapolate forward dates from git history or older handoffs.
   - The sprint handoff artifact MUST be named `Sprint-YYYY-MM-DD-Handoff.md` (append `-{i}` if multiple sprints occur today, e.g. `Sprint-YYYY-MM-DD-1-Handoff.md`).
2. **Context Recovery**:
   - Read the latest sprint handoff note in `.obsidian/MAgHARCM/diary/` to ingest open watchlist items, deferred triggers, and in-flight tasks.
   - Read `.obsidian/MAgHARCM/adhoc/Human-Intervention-And-Blockers.md` to identify active blockers requiring triage.
3. **Artifact Hygiene**:
   - Gracefully sweep stale temporary artifacts from `.artifacts/local/` and test output caches. Do not build upon obsolete assumptions.

---

## Phase 1: Autonomous Research & Literature Expansion

1. **Trigger-Gate Evaluation (Wave-N+1 Criterion)**:
   - Evaluate candidate literature against the 3-question gate:
     - **Q1 (Venue Verification)**: Is the paper verified in a major peer-reviewed venue (NeurIPS, ICML, ICLR, FSE, ICSE, ASE, TOSEM) or authoritative archival repository?
     - **Q2 (Mechanism Distinction)**: Does it contribute an actionable mechanism or algorithm (not merely a benchmark leaderboard rank or surface prompt)?
     - **Q3 (Primitive Anchoring)**: Does it defend, refute, or supply the substrate for a MAgHARCM primitive?
2. **2-Hop Citation Graph Tracking**:
   - For every admitted paper, systematically track:
     - **Hop-1**: Direct references that ground the paper's mechanism.
     - **Hop-2**: Foundational antecedents cited by Hop-1 authors.
3. **Domain Focus & Biases**:
   - **Software Archaeology**: Prioritize temporal coupling, program slicing, dynamic invariants, concept assignment, dependency structure matrices (DSMs), and design rule hierarchy partitioning.
   - **Small Language Models (4B–30B)**: Focus on techniques enabling local SLMs (Qwen2.5-Coder, StarCoder2, Phi-3) to perform whole-repo modernization without frontier cloud APIs (structured cloze slots, KV cache compression, speculative decoding, test-time scaling, self-consistency).
4. **Citation Integrity**:
   - Never fabricate or guess citations. If an important source lacks confirmed venue verification, tag it as `UNVERIFIED placeholder` and log it in `adhoc/Human-Intervention-And-Blockers.md`.

---

## Phase 2: Dual-Interface Vault Synchronization

1. **Machine-Readable Autoresearch Database**:
   - Update `.obsidian/MAgHARCM/Research-Database.json` whenever new papers, primitives, or waves land. Ensure all fields (`id`, `title`, `authors`, `year`, `venue`, `eprint_or_doi`, `bibkey`, `anchored_primitives`, `hop1_references`, `hop2_references`, `summary`, `relevance`, `verified`) are populated.
2. **Human Executive Reports Suite (`adhoc/`)**:
   - Synchronize living documents under `.obsidian/MAgHARCM/adhoc/`:
     - `adhoc/Methodology.md` — Section 7 SLM anchors, Section 9 changelog, Section 11 patterns.
     - `adhoc/Project-Progress-And-Milestones.md` — metrics dashboard, velocity, completed milestones.
     - `adhoc/Strategic-Direction-And-Roadmap.md` — strategic vision, wave planning.
     - `adhoc/Human-Intervention-And-Blockers.md` — active blockers, compilation breaks, unverified sources.
     - `adhoc/Benchmark-Results-And-Evaluation.md` — empirical trial results and ablation data.
     - `adhoc/Architecture-And-Dataflow.md` — multi-agent graph and state schema changes.
     - `adhoc/Research-Waves-Index.md` — thematic and chronological synthesis of waves.
3. **Markdown Notes & Parity**:
   - Persist paper notes in `research/papers/P-NN-<CanonicalName>.md` with full YAML `aliases:` so wikilinks resolve in Obsidian.
   - Cross-link into `research/Software-Archaeology-Lineage.md` and `primitives/Primitives-Index.md`.
   - Maintain 100% parity across `Primitives-Index.md`, `Software-Archaeology-Lineage.md`, `Methodology.md`, and `internal/`.
4. **Naming & Lint Invariants**:
   - Every file name across all subdirectories MUST use **Camel-Case (with hyphen)** (`Word-Word-Word.md`).
   - Root `README.md` uses standard uppercase.
   - Version slots MUST use `[[x.y.z ...]]` wikilink syntax.
   - Execute `./scripts/lint_vault.sh` from repo root and ensure exit code 0.

---

## Phase 3: Honest Empirical Evaluation

1. **Benchmark Suite**:
   - Evaluate against the 4 canonical repositories:
     - `GildedRose` (C $\to$ Rust, 199 LoC)
     - `Gohistogram` (Go $\to$ Rust, 15,470 LoC)
     - `Stats` (Go $\to$ Rust, 5,625 LoC)
     - `Commons-Validator` (Java $\to$ Rust, 28,110 LoC)
2. **Measurement Protocol**:
   - Measure real outcomes: Binary compilation status (`Pass`/`Fail`, no partial rates) and test pass rate ($X/Y$, percentage).
   - Never fake, mock, or skip experiments.
   - Record findings in `adhoc/Benchmark-Results-And-Evaluation.md` and `docs/.paper/sec_eval.tex`.

---

## Phase 4: Codebase Defense & Ponytail Engineering

1. **Defend Standing Invariants**:
   - **Must Pattern**: Compile-time inits and static configs MUST NOT have silent fallbacks; panic on invariant violation. Runtime configs come from external files (`configs/*.yml`).
   - **Cohesive State & Locality of Behaviour (ADR-C-014)**: Shared state is declared in `internal/compiletime/state.go`. Agent-local artifact types reside within their producer files, using type aliases where import cycles would occur.
   - **Dynamic Strategy Registry (`PRIM-21`)**: Try-and-fail strategy cascade (`internal/agents/strategy.go`). Increments to the next strategy upon validation failure.
   - **8-Agent Eino Execution Graph**: All 8 agents (`archaeologist`, `analyzer`, `planner`, `translator`, `reviewer`, `validator`, `verdict_panel`, `recruiter`) wired in `internal/graph/graph.go`.
   - **Idiomatic Charm TUI**: Proper Bubble Tea model/update/view, Lip Gloss styling, Glamour markdown rendering in `cmd/MAgHARCM-tui`.
   - **Zero Raw Printing**: Zero `fmt.Print*` in production Go code. Use structured logging (`internal/logger`).
   - **Externalities First**: Prefer standard/battle-tested external libraries (`yaml.v3`, `bubbletea`, `eino`) over reinvented custom helpers.
   - **Terse Technical Messaging (`asd-ste100`)**: Imperative, direct prose in errors, logs, comments, and commit messages.
2. **Active Blocker Triage**:
   - Resolve any open blockers listed in `adhoc/Human-Intervention-And-Blockers.md`.

---

## Phase 5: Academic Paper Synchronization (`docs/.paper/`)

1. Update `docs/.paper/sec_method.tex` and `main.tex` with newly anchored research papers (`\cite{...}`).
2. Append valid BibTeX entries to `docs/.paper/refs.bib`.
3. Synchronize `docs/.paper/sec_eval.tex` with updated empirical metrics.
4. Ensure paper claims strictly mirror codebase and vault ground truth.

---

## Phase 6: Handoff Persistence & Incremental Commits

1. **Sprint Handoff Note**:
   - Persist summary at `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md` detailing closed items, commits, audit tables, and watchlist triggers.
2. **Conventional Commits**:
   - Commit incrementally per phase boundary:
     - `feat(vault): ...` for research waves and vault sync.
     - `fix(...)` / `feat(...)` for codebase implementation.
     - `docs(paper): ...` for LaTeX paper updates.
     - `docs(diary): ...` for sprint handoff.
