---
description: "Execute a research sprint, vault sync, empirical tests, code defense, and command evolution"
---

# MAgHARCM Research Sprint and Modernization Workflow

Rules: follow the phases in order. Run commands from the repository root. Make one conventional commit per listed commit line; you may append ` wave-<N>` or a date to the subject. Never commit with an empty subject. Skip a commit when its phase changed nothing. This file is the procedure. The vault is the data. `.obsidian/MAgHARCM/primitives/Primitives-Index.md` is the single source of truth for the map of primitives to anchors. Never hardcode that map here. Standing rules (section 9) apply to every phase; re-check them in phase 7.

## 1. Start the Sprint

1. Get the current date in `YYYY-MM-DD` form from the system reminder or `date -u`. Use it for every filename and header this sprint.
2. Read the newest handoff note in `.obsidian/MAgHARCM/diary/`: highest date in the filename wins, modification time breaks ties. Investigate any future-dated note before trusting it.
3. Set the handoff path to `Sprint-<today>-Handoff.md`. If that file exists, append `-1`, `-2`, and so on.
4. Read `.obsidian/MAgHARCM/adhoc/Human-Intervention-And-Blockers.md`. Record open blockers.
5. Delete stale files under `.artifacts/local/`.
6. Commit: `docs(diary): start sprint YYYY-MM-DD`.

## 2. Research and Literature Expansion

1. Search top venues only: NeurIPS, ICML, ICLR, FSE, ICSE, ASE, TOSEM, TSE.
2. Admit a paper only when it passes three gates:
   - **Q1 (Venue)**: peer-reviewed venue from that list. A workshop paper passes only with a single-paper method-level mechanism no other paper in this wave covers (precedent: P-138 RepairKV). Empirical or diagnostic papers never qualify.
   - **Q2 (Mechanism)**: a concrete algorithmic mechanism, not a prompt tweak, and not a duplicate of an already-anchored mechanism.
   - **Q3 (Anchor)**: anchors or defends a primitive from `Primitives-Index.md`.
3. Give every admitted paper two citation hops in `Research-Database.json`: `hop1_references` (at least one citation that directly supports the mechanism) and `hop2_references` (at least one foundational source). Both must be non-empty.
4. Scout small language models (4B-30B) and software archaeology: program slicing, dynamic invariants, execution traces, concept assignment, dependency graphs, KV cache compression, speculative decoding, SLM verification, test-time compute.
5. Stop the wave when it yields at least two admitted papers or at least four triaged candidates.
6. Write a triage memo to `.obsidian/MAgHARCM/research/diary/Wave-<N>-Candidates.md` before you persist any paper note. `N` is one plus the highest wave number in `Research-Waves-Index.md` and existing `Wave-*-Candidates.md` memos. The memo lists every candidate with a verdict and a one-line rationale.
7. Never fabricate citations or metrics. Mark an unverifiable candidate as UNVERIFIED in `watchlist.wave-<N>`.
8. Commit: `feat(research): add wave-<N> research papers`.

## 3. Synchronize the Obsidian Vault

1. Update `.obsidian/MAgHARCM/Research-Database.json`:
   - Add one entry per admitted paper. Mirror the keys of the newest existing entries, including `id`, `title`, `authors`, `year`, `venue`, `bibkey`, `eprint_or_doi`, `hop1_references`, `hop2_references`, `summary`, `anchored_primitives`, `verified`, and `version_marker`.
   - Append rejects to `reject_registry.wave-<N>` with bibkey, title, venue, verdict (Q1, Q2, or Q3), and rationale.
   - Append unverified candidates to `watchlist.wave-<N>`.
   - Verify the file parses: `python3 -c "import json; json.load(open('.obsidian/MAgHARCM/Research-Database.json'))"`.
2. Save one paper note per admitted paper to `.obsidian/MAgHARCM/research/papers/P-<NN>-<Name>.md`. `NN` is one plus the highest existing P-number. Use Camel-Case with hyphens for `<Name>`. Declare YAML aliases, including the version marker (`[[1.0.0 P-<NN>]]`).
3. Update every affected human report in `.obsidian/MAgHARCM/adhoc/` per the vault conventions rule. At minimum: new anchors in `Methodology.md`, new metrics in `Project-Progress-And-Milestones.md`, the new wave in `Research-Waves-Index.md`, and updated blockers.
4. Cross-reference new papers in `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` and update `Primitives-Index.md`. Keep the primitives index, lineage matrix, and code in full parity.
5. Run `./scripts/lint_vault.sh` from the repository root. Fix every finding until it exits with code 0.
6. Commit: `feat(vault): sync research database and adhoc reports`.

## 4. Run Empirical Experiments

1. Run the translation benchmark suite through the crust pipeline (`benchmarks/crust/scripts/run.sh`), or the K-trial pipeline (`scripts/run-samples-k.sh`). Read the current target repositories from `benchmarks/crust/README.md`.
2. Record compilation status as `Pass` only when `cargo check` exits 0 and the entire test suite executes. Record `Fail` for any other outcome, including an interrupted run.
3. Record exact test counts and percentages, for example `24/42 (57.1%)`.
4. Record numbers only from machine-readable artifacts: `benchmarks/crust/results/<HEAD at run start>/<proj>/result.yml` when a run completes, `docs/sample-results/k-summary.json` after a K-trial run, or the last `FINALSUM` log line from `internal/agents/validator.go`. Link the artifact path beside every recorded row.
5. Never mock, fake, or skip experiments. If the inference daemon is offline (see open blockers), record the deferred experiments in `.obsidian/MAgHARCM/adhoc/Benchmark-Results-And-Evaluation.md` and continue to phase 5. Do not invent numbers.
6. Write results to `Benchmark-Results-And-Evaluation.md` and `docs/.paper/sec_eval.tex`.
7. Commit only if result rows changed: `test(benchmarks): record empirical evaluation results`.

## 5. Implement Primitives and Defend the Codebase

1. Keep the codebase aligned with admitted primitives and current research progress.
2. Apply ponytail and eino principles. Prefer the Go standard library. Delete dead code, unused abstractions, and redundant shims. Keep implementations minimal, readable, and executable.
3. Document design intent inline with comments that reference primitive IDs and paper bibkeys.
4. Initialize values with the `Must` pattern at compile time. Panic on missing or unparseable configuration. Read runtime settings from YAML.
5. Declare shared state in `internal/compiletime/state.go`. Producer packages in `internal/agents/` reference those types through aliases, per ADR-C-014, to avoid import cycles.
6. Wire all eight agents in `internal/graph/graph.go`: Archaeologist, Analyzer, Planner, Translator, Reviewer, VerdictPanel, Validator, Recruiter.
7. Keep the five Mueller strategies in canonical order in `internal/agents/strategy.go`: Big Bang, Pilot, Frozen Legacy, Parallel Cutover, Incremental.
8. Log through `internal/logger`. Never call `fmt.Print*` in production Go code.
9. Resolve blockers from `Human-Intervention-And-Blockers.md` where code changes can resolve them.
10. Run `go test ./...` and make sure every test passes.
11. Commit: `feat(primitives): implement and wire PRIM-<XX>`.

## 6. Synchronize the Academic Paper

1. Update `docs/.paper/sec_method.tex` and `main.tex` with newly admitted primitives and mechanisms.
2. Add verified BibTeX entries to `docs/.paper/refs.bib`.
3. Update `docs/.paper/sec_eval.tex` with verified benchmark numbers.
4. Make sure claims match across the paper, the codebase, and the vault.
5. Commit: `docs(paper): update method and references`.

## 7. Evolve This Command File

1. Remove directives that describe completed work, expired temporary rules, and carried watchlist entries now tracked in the vault.
2. Add new research anchors, active watchlists, and resolved blockers.
3. Keep this file concise and compliant with ASD-STE100. Edit in one pass per sprint; do not iterate on it.
4. Commit: `feat(command): evolve MAgHARCM sprint workflow`.

## 8. Persist Sprint Handoff

1. Write the sprint summary to the handoff path from phase 1. Document closed tasks, commit sequence, audit tables, and active triggers.
2. Commit: `docs(diary): persist sprint YYYY-MM-DD handoff`.

## 9. Standing Governance Rules

1. **Dating (BLK-06)**: the system reminder or `date -u` is the only authoritative date. Never extrapolate future dates from git history or older notes. Every YAML `date` must equal its filename date.
2. **Registries (BLK-08)**: every triage wave appends rejects to `reject_registry.wave-<N>` and unverified candidates to `watchlist.wave-<N>` in `Research-Database.json`.
3. **Sandbox Batch Commit**: if git mutations are blocked, accumulate artifacts in the working tree and batch-commit when the sandbox recovers. Use per-phase commits again as soon as the sandbox recovers.

## 10. Active Watchlist

The vault owns per-wave watchlist state (`watchlist.wave-<N>` in `Research-Database.json` and the blockers note). This section holds only constraints the vault does not.

1. **Comprehension-mechanism slot**: closed by `[[1.0.0 P-153]]` CoReX. Do not reopen broad comprehension queries.
2. **Substrate wiring**: verify that `configs/agents.yml` exposes opt-in keys for the substrates named in `Methodology.md`. Add missing keys when you implement those substrates.
