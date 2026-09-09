---
description: "Execute a research sprint, vault sync, empirical tests, code defense, and command evolution"
---

# MAgHARCM Research Sprint and Modernization Workflow

Commit rules: run commands from the repository root. Make one conventional commit per listed commit line; you may append ` wave-<N>` or a date. Never commit with an empty subject. Skip a commit when the phase changed nothing durable. Data ownership: this file is the procedure. The vault is the data. `.obsidian/MAgHARCM/primitives/Primitives-Index.md` is the single source of truth for the map of primitives to anchors; never hardcode that map here. Section 9 holds standing rules; apply them in every phase and re-check them in phase 7.

## 1. Start the Sprint

1. Get the current date in `YYYY-MM-DD` form from `date -u` or the orchestration layer. Use it for every filename and header this sprint.
2. Read the newest handoff note in `.obsidian/MAgHARCM/diary/`: pick the highest filename date that is not in the future; skip future-dated notes; break ties by the highest `-<N>` suffix. If no past-dated handoff exists, continue and record that in the handoff.
3. Set the handoff path to `Sprint-<today>-Handoff.md`. Continue from one plus the highest existing `-<N>` suffix for today.
4. Read `.obsidian/MAgHARCM/adhoc/Human-Intervention-And-Blockers.md`. Record open blockers in the handoff note. If it is missing, treat blockers as unknown and record that.
5. Delete all files under `.artifacts/local/` (create the directory if missing). Never touch other `.artifacts/` paths.
6. If the git index is dirty on arrival, commit or stash the stray changes before the sprint commit. Commit: `docs(diary): start sprint YYYY-MM-DD`.

## 2. Research and Literature Expansion

1. Load once now: `Primitives-Index.md`, `Research-Database.json`, `Research-Waves-Index.md`, and a glob of `.obsidian/MAgHARCM/research/diary/Wave-*-Candidates.md`. Phases 2, 3, and 6 edit from this loaded state. If `Primitives-Index.md` is unreadable, record a blocker, triage on Q1/Q2 only with Q3 marked PENDING, and cap the search at four triaged candidates.
2. Search top venues only: NeurIPS, ICML, ICLR, FSE, ICSE, ASE, TOSEM, TSE. Take the current topic list from the newest wave's forward plan in `Research-Waves-Index.md`; small language models (4B-30B) and software archaeology are standing scope.
3. Admit a paper only when it passes three gates; empirical or diagnostic papers never qualify.
   - **Q1 (Venue)**: peer-reviewed venue from that list. A workshop paper passes only with a single-paper method-level mechanism no other paper in this wave covers (precedent: P-138 RepairKV).
   - **Q2 (Mechanism)**: a concrete algorithmic mechanism, not a prompt tweak, and not a duplicate of an already-anchored mechanism.
   - **Q3 (Anchor)**: anchors or defends a primitive from `Primitives-Index.md`.
4. Stop the wave at two admitted papers or four triaged candidates, whichever comes first, or when search dries up. Record scanned topic families in the memo so coverage gaps stay visible.
5. `N` is one plus the highest wave number across the index and memos; with no prior wave, use 1.
6. Write a triage memo to `.obsidian/MAgHARCM/research/diary/Wave-<N>-Candidates.md` before you persist any paper note. The memo lists every candidate with a verdict and a one-line rationale. Never fabricate citations or metrics; mark an unverifiable candidate as UNVERIFIED and continue without admitting it.
7. Commit: `feat(research): add wave-<N> research papers`.

## 3. Synchronize the Obsidian Vault

1. Update `.obsidian/MAgHARCM/Research-Database.json`:
   - Add one entry per admitted paper. Mirror the keys of the newest existing entries, including `id`, `title`, `authors`, `year`, `venue`, `bibkey`, `eprint_or_doi`, `hop1_references`, `hop2_references`, `summary`, `anchored_primitives`, `verified`, and `version_marker`.
   - Append rejects and unverified candidates to `reject_registry.wave-<N>` and `watchlist.wave-<N>` per rule 9.2, using the field names of existing entries (`candidate_id`/`short_name`/`reject_class`; wave-slim carry schema for the watchlist).
   - Verify the file parses: `python3 -c "import json; json.load(open('.obsidian/MAgHARCM/Research-Database.json'))"` (or `jq empty` if python3 is absent). If the file is missing, create it with empty `reject_registry` and `watchlist` objects.
2. Save one paper note per admitted paper to `.obsidian/MAgHARCM/research/papers/P-<NN>-<Name>.md`. `NN` is one plus the highest existing P-number; with no prior notes, use 1. Use Camel-Case with hyphens for `<Name>`. Declare YAML aliases, including the version marker (`[[1.0.0 P-<NN>]]`).
3. Update every affected human report in `.obsidian/MAgHARCM/adhoc/` per the vault conventions rule. At minimum: new anchors in `adhoc/Methodology.md`, new metrics in `Project-Progress-And-Milestones.md`, the new wave in `Research-Waves-Index.md`, and updated blockers.
4. Cross-reference new papers in `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` and update `Primitives-Index.md` from the phase-2 loaded state. Keep the primitives index, lineage matrix, and code in full parity.
5. Run `./scripts/lint_vault.sh` from the repository root. Fix the findings, then re-run once. If it still fails after two fix passes, record the findings in the blockers note and commit anyway.
6. Commit: `feat(vault): sync research database and adhoc reports`.

## 4. Run Empirical Experiments

1. Run the translation benchmark suite through the crust pipeline (`benchmarks/crust/scripts/run.sh`), or the K-trial pipeline (`scripts/run-samples-k.sh`). Read the current target repositories from `benchmarks/crust/configs/` and `benchmarks/crust/README.md`. If a target is missing, record it as a blocker; do not silently shrink the suite.
2. Record compilation status as `Pass` only when `cargo check` exits 0 and the entire test suite executes. Record `Fail` when the run executed and did not pass. If the run was stopped before completion by the operator or the environment, record it as deferred with the partial log path. When a run fails, name the failing pipeline node and error in the results note.
3. Record exact test counts and percentages, for example `24/42 (57.1%)`.
4. Record numbers only from machine-readable artifacts: `benchmarks/crust/results/<12-char HEAD>/<proj>/result.yml` when a run completes, `docs/sample-results/k-summary.json` after a K-trial run, or the last `FINALSUM` log line from `internal/agents/validator.go`. Link the artifact path beside every recorded row.
5. Never mock, fake, or skip experiments. If the inference daemon is offline (see the phase-1 blocker record), record the deferred experiments in `.obsidian/MAgHARCM/adhoc/Benchmark-Results-And-Evaluation.md` and continue to phase 5. Do not invent numbers.
6. Write results to `Benchmark-Results-And-Evaluation.md` only; phase 6 owns `sec_eval.tex`.
7. Commit only if result rows changed: `test(benchmarks): record empirical evaluation results`.

## 5. Implement Primitives and Defend the Codebase

1. Keep the codebase aligned with admitted primitives and current research progress.
2. Apply ponytail and eino principles. Prefer the Go standard library. Delete dead code, unused abstractions, and redundant shims. Keep implementations minimal, readable, and executable.
3. Document design intent inline with comments that reference primitive IDs and paper bibkeys.
4. Name constructors with the `Must` prefix only when they panic. Initialize package-level values with `Must*` constructors so a bad value surfaces at init. Panic on missing or unparseable configuration. Read runtime settings from YAML.
5. Artifact structs live in their producer file under `internal/agents/`; `internal/compiletime/state.go` re-exports them as `type X = agents.X` aliases (ADR-C-014 as applied). Never move structs into `compiletime`; the import cycle makes that impossible.
6. Wire the full agent roster (currently eight: Archaeologist, Analyzer, Planner, Translator, Reviewer, VerdictPanel, Validator, Recruiter) in `internal/graph/graph.go`.
7. Keep the five Mueller strategies in canonical order in `internal/agents/strategy.go`: Big Bang, Pilot, Frozen Legacy, Parallel Cutover, Incremental.
8. Log through `internal/logger`. Never call `fmt.Print*` in production Go code.
9. When you implement a substrate named in `Methodology.md`, add its `configs/agents.yml` opt-in key in the same commit.
10. Resolve the blockers recorded in phase 1 where code changes can resolve them.
11. Run `go test ./...` and make sure every test passes.
12. Commit: `feat(primitives): implement and wire PRIM-<XX>`.

## 6. Synchronize the Academic Paper

1. Update `docs/.paper/sec_method.tex` and `main.tex` with newly admitted primitives and mechanisms.
2. Add verified BibTeX entries to `docs/.paper/refs.bib`.
3. Update `docs/.paper/sec_eval.tex` with the verified benchmark numbers from phase 4.
4. Make sure claims match across the paper, the codebase, and the vault.
5. Commit: `docs(paper): update method and references`.

## 7. Evolve This Command File

1. Re-read section 9. Remove directives that describe completed work, expired temporary rules, and per-wave state now tracked in the vault.
2. Promote only constraints the vault cannot express: new research anchors to `Primitives-Index.md`, not here. List every removed directive in the handoff note with the resolving artifact.
3. Keep this file concise and compliant with ASD-STE100. Edit in one pass per sprint; do not iterate on it.
4. Commit: `feat(command): evolve MAgHARCM sprint workflow`.

## 8. Persist Sprint Handoff

1. Write the sprint summary to the handoff path from phase 1 (append `-<N>` again if a session rerun collides). Document closed tasks, commit sequence, audit tables, and active triggers. If a session must pause early, set status `paused` with the next phase; otherwise close it.
2. Commit: `docs(diary): persist sprint YYYY-MM-DD handoff`.

## 9. Standing Governance Rules

1. **Dating**: `date -u` or the orchestration layer is the only authoritative date. Never extrapolate future dates from git history or older notes. Every YAML `date` must equal its filename date. A session may run all phases or a subset; record completed phases and the next phase in the handoff.
2. **Registries**: every triage wave appends rejects to `reject_registry.wave-<N>` and unverified candidates to `watchlist.wave-<N>` in `Research-Database.json`, using the field names of existing entries. Never fabricate citations or metrics. A paper already anchored under another P-number is a non-event: record it in the memo only, never in a registry.
3. **Sandbox Batch Commit**: if git mutations are blocked, accumulate artifacts in the working tree and batch-commit when git works again. Use per-phase commits again as soon as it recovers.

## 10. Durable Watchlist Constraints

The vault owns per-wave watchlist state (`watchlist.wave-<N>` in `Research-Database.json` and the blockers note). This section holds only durable constraints the vault does not.

1. **Comprehension-mechanism slot**: closed by `[[1.0.0 P-153]]` CoReX. Do not reopen broad comprehension queries; partition-aligned summary work lives in `Methodology.md`.
