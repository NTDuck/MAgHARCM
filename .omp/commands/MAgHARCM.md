---
description: "Execute a research sprint, vault sync, empirical tests, code defense, and command evolution"
---

# MAgHARCM Research Sprint and Modernization Workflow

Rules: follow the phases in order. Run commands from the repository root. Make one conventional commit per phase, exactly as written. If a phase has nothing to change, skip its commit. This file is the procedure; the vault is the data. `.obsidian/MAgHARCM/primitives/Primitives-Index.md` is the single source of truth for the primitive-to-anchor mapping and paper identities. Never hardcode that mapping here.

## 1. Start the Sprint

1. Get the current date in `YYYY-MM-DD` form from the system reminder or `date -u`. Use it for every filename and header this sprint.
2. Read the newest handoff note in `.obsidian/MAgHARCM/diary/` by modification time. Set the sprint handoff path to `Sprint-<today>-Handoff.md`. If that file exists, append `-1`, `-2`, and so on.
3. Read `.obsidian/MAgHARCM/adhoc/Human-Intervention-And-Blockers.md`. Record open blockers.
4. Delete stale temporary files from `.artifacts/local/` and test caches.
5. Commit: `docs(diary): start sprint YYYY-MM-DD`.

## 2. Research and Literature Expansion

1. Search top venues only: NeurIPS, ICML, ICLR, FSE, ICSE, ASE, TOSEM, TSE.
2. Admit a paper only when it passes three gates:
   - **Q1 (Venue)**: peer-reviewed main-track venue. A workshop paper passes only with a method-level single-paper mechanism that no other accepted paper in this wave covers (threshold: P-138 RepairKV). Empirical or diagnostic papers never qualify.
   - **Q2 (Mechanism)**: a concrete algorithmic mechanism, not a prompt tweak.
   - **Q3 (Anchor)**: anchors or defends a primitive from `Primitives-Index.md`.
3. Record two citation hops per admitted paper in `Research-Database.json`: `hop1_references` (direct mechanism support) and `hop2_references` (foundational sources). Both must be non-empty.
4. Scout small language models (4B-30B) and software archaeology: program slicing, dynamic invariants, execution traces, concept assignment, dependency graphs, KV cache compression, speculative decoding, SLM verification, test-time compute.
5. Stop the wave when it yields at least two admitted papers or at least four triaged candidates.
6. Write a triage memo to `.obsidian/MAgHARCM/research/diary/Wave-<N>-Candidates.md` before persisting any paper note. `N` is one plus the highest wave in `Research-Waves-Index.md`. The memo lists every candidate with verdict and one-line rationale.
7. Never fabricate citations or metrics. Mark an unverifiable candidate as UNVERIFIED in `watchlist.wave-<N>`.
8. Commit: `feat(research): add wave research papers`.

## 3. Synchronize the Obsidian Vault

1. Update `.obsidian/MAgHARCM/Research-Database.json`:
   - Add one entry per admitted paper with `id`, `title`, `authors`, `year`, `venue`, `bibkey`, one identifier (DOI, OpenReview ID, or arXiv ID), `hop1_references`, `hop2_references`, and `summary`.
   - Append rejects to `reject_registry.wave-<N>` with bibkey, title, venue, verdict (Q1, Q2, or Q3), and rationale.
   - Append unverified candidates to `watchlist.wave-<N>`.
   - Verify the file parses: `python3 -c "import json; json.load(open('.obsidian/MAgHARCM/Research-Database.json'))"`.
2. Save one paper note per admitted paper to `.obsidian/MAgHARCM/research/papers/P-<NN>-<Name>.md`. Declare YAML aliases, including the version marker (`[[1.0.0 P-<NN>]]`).
3. Update all seven human reports in `.obsidian/MAgHARCM/adhoc/` per the vault conventions rule. At minimum: new anchors in `Methodology.md`, new metrics in `Project-Progress-And-Milestones.md`, the new wave in `Research-Waves-Index.md`, and updated blockers.
4. Cross-reference new papers in `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` and update `.obsidian/MAgHARCM/primitives/Primitives-Index.md`. Keep the primitives index, lineage matrix, and code in full parity.
5. Run `./scripts/lint_vault.sh` from the repository root. Fix every finding until it exits with code 0.
6. Commit: `feat(vault): sync research database and adhoc reports`.

## 4. Run Empirical Experiments

1. Run the translation benchmark suite. The current target is GildedRose (C to Rust).
2. Record compilation status as `Pass` only when `cargo check` exits 0 and the entire test suite executes. Record `Fail` for any other outcome, including an interrupted run.
3. Record exact test counts and percentages, for example `24/42 (57.1%)`.
4. Record numbers only from machine-readable artifacts: `benchmarks/crust/results/<commit>/<proj>/result.yml`, `docs/sample-results/k-summary.json`, or the last `FINALSUM` log line. Link the artifact path beside every recorded row.
5. Never mock, fake, or skip experiments. If the inference daemon is offline (blocker BLK-04), record the deferral in `Benchmark-Results-And-Evaluation.md` and move to phase 5. Do not invent numbers.
6. Write results to `.obsidian/MAgHARCM/adhoc/Benchmark-Results-And-Evaluation.md` and `docs/.paper/sec_eval.tex`.
7. Commit: `test(benchmarks): record empirical evaluation results`.

## 5. Implement Primitives and Defend the Codebase

### 5.1 Generative Implementation

1. Keep the codebase aligned with admitted primitives and current research progress.
2. Apply ponytail and eino principles. Prefer the Go standard library. Delete dead code, unused abstractions, and redundant shims. Keep implementations minimal, readable, and executable.
3. Document design intent inline with comments that reference primitive IDs and paper bibkeys.

### 5.2 Defensive Invariants

1. Apply the `Must` pattern for compile-time initialization. Panic on missing or unparseable configuration. Read runtime settings from YAML in `configs/`.
2. Declare shared state in `internal/compiletime/state.go`. Producer packages in `internal/agents/` reference those types through aliases, per ADR-C-014, to avoid import cycles.
3. Wire all eight agents in `internal/graph/graph.go`: Archaeologist, Analyzer, Planner, Translator, Reviewer, VerdictPanel, Validator, Recruiter.
4. Keep the five Mueller strategies in canonical order in `internal/agents/strategy.go`: Big Bang, Pilot, Frozen Legacy, Parallel Cutover, Incremental.
5. Log through `internal/logger`. Never call `fmt.Print*` in production Go code.
6. Resolve blockers recorded in `Human-Intervention-And-Blockers.md` where code changes can resolve them.
7. Run `go test ./...` and make sure every test passes.
8. Commit: `feat(agents): implement PRIM-<XX> from P-<YY> and defend invariants`.

## 6. Synchronize the Academic Paper

1. Update `docs/.paper/sec_method.tex` and `main.tex` with newly admitted primitives and mechanisms.
2. Add verified BibTeX entries to `docs/.paper/refs.bib`.
3. Update `docs/.paper/sec_eval.tex` with verified benchmark numbers.
4. Make sure claims match across paper, codebase, and vault.
5. Commit: `docs(paper): update method and references`.

## 7. Evolve This Command File

1. Remove directives that describe completed work, expired temporary rules, and carried watchlist entries now tracked in the vault.
2. Add new research anchors, active watchlists, and resolved blockers.
3. Keep this file concise, exhaustive, and compliant with ASD-STE100.
4. Commit: `feat(command): evolve MAgHARCM sprint workflow`.

## 8. Persist Sprint Handoff

1. Write the sprint summary to `.obsidian/MAgHARCM/diary/Sprint-<today>-Handoff.md`.
2. Document closed tasks, commit sequence, audit tables, and active triggers.
3. Commit: `docs(diary): persist sprint YYYY-MM-DD handoff`.

## 9. Standing Governance Rules

1. **Dating (BLK-06)**: the system reminder or `date -u` is the only authoritative date. Never extrapolate future dates from git history or older notes. Every YAML `date` must equal its filename date.
2. **Registries (BLK-08)**: every triage wave appends rejects to `reject_registry.wave-<N>` and unverified candidates to `watchlist.wave-<N>` in `Research-Database.json`.
3. **Sandbox Batch Commit**: if git mutations are blocked, accumulate artifacts in the working tree and batch-commit on recovery. Revert to per-phase commits immediately after recovery.

## 10. Active Watchlist

1. **Comprehension-mechanism slot**: CLOSED by `[[1.0.0 P-153]]` CoReX. Do not reopen broad comprehension queries. Residual partition-aligned summary work stays in `Methodology.md`.
2. **SWE-TRACE** (`arXiv:2604.14820`): re-verify the venue after NeurIPS 2026 notifications. After five carries without a peer-reviewed venue, retire it to the reject archive.
3. **Substrate wiring**: verify that `configs/agents.yml` exposes opt-in keys for the comprehension, hallucination-evaluation, and translation substrates named in `Methodology.md`. Add missing keys when implementing those substrates.
4. **Cross-pattern formalization**: skip until the current wave items close. Candidate pairs live in `Strategic-Direction-And-Roadmap.md`.
