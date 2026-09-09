---
description: "Execute a research sprint, vault sync, empirical tests, code defense, and command evolution"
---

# MAgHARCM Research Sprint and Modernization Workflow

## 1. Start the Sprint

1. Get the authoritative date in `YYYY-MM-DD` format from the system reminder or `date -u`; never extrapolate future dates (BLK-06).
2. Make sure note metadata headers (`date:`, `last_updated:`) match the current date.
3. Set the sprint handoff path to `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md`; increment suffix (`-1`, `-2`, `-N`) if it exists.
4. Read the latest handoff note for pending tasks.
5. Read `.obsidian/MAgHARCM/adhoc/Human-Intervention-And-Blockers.md` for active blockers.
6. Delete stale temporary files from `.artifacts/local/` and test caches.
7. Commit: `docs(diary): start sprint YYYY-MM-DD`.

## 2. Research and Literature Expansion

1. Search top venues: NeurIPS, ICML, ICLR, FSE, ICSE, ASE, TOSEM, TSE.
2. Admit a paper only when it passes three gates:
   - **Q1 (Venue)**: peer-reviewed top venue; workshop papers need single-paper method novelty and non-redundant theme (P-138 threshold; empirical/diagnostic work disqualified).
   - **Q2 (Mechanism)**: concrete algorithmic mechanism, not a prompt tweak.
   - **Q3 (Anchor)**: anchors or defends a MAgHARCM primitive (`PRIM-01`–`PRIM-33`).
3. Record two citation hops per admitted paper: Hop 1 direct mechanism support, Hop 2 foundational anchor.
4. Scout SLM (4B–30B) and software-archaeology candidates; the active substrate focus is §6, not a free-form list.
5. Record non-admitted candidates (BLK-08):
   - Triage memo: `.obsidian/MAgHARCM/research/diary/Wave-NN-Candidates.md`.
   - `reject_registry.wave-NN` entries: bibkey, title, venue, verdict (Q1/Q2/Q3), and a four-way rationale: (a) hallucinated-venue, (b) off-list venue, (c) mechanism-overlap, (d) off-axis target-language/problem.
   - Unverified candidates: `watchlist.wave-NN` for re-verification in later waves.
6. Never fabricate citations or metrics; label unverified candidates as placeholders in the blockers note.
7. Commit: `feat(research): add wave research papers`.

## 3. Synchronize the Obsidian Vault

1. Update `.obsidian/MAgHARCM/Research-Database.json`:
   - Required fields: `id`, `title`, `authors`, `year`, `venue`, `bibkey`, identifier (DOI, OpenReview ID, or arXiv ID), `citations`, `summaries`.
   - Append rejects to `reject_registry.wave-NN`; unverified candidates to `watchlist.wave-NN`.
   - Verify JSON parses: `python3 -c "import json; json.load(open('.obsidian/MAgHARCM/Research-Database.json'))"`.
2. Save paper notes to `.obsidian/MAgHARCM/research/papers/P-NN-<Name>.md`; declare YAML aliases including the version marker (`[[1.0.0 P-NN]]`).
3. Update human reports in `.obsidian/MAgHARCM/adhoc/`:
   - `Methodology.md`: anchors, justifications, closed slots, changelog.
   - `Project-Progress-And-Milestones.md`: metrics, wave tables, milestones.
   - `Strategic-Direction-And-Roadmap.md`: roadmap, wave plans.
   - `Human-Intervention-And-Blockers.md`: blockers, resolutions.
   - `Benchmark-Results-And-Evaluation.md`: empirical data.
   - `Architecture-And-Dataflow.md`: state and graph diagrams.
   - `Research-Waves-Index.md`: wave index ledger.
4. Cross-reference new papers in `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`; update `.obsidian/MAgHARCM/primitives/Primitives-Index.md`.
5. Maintain full parity across primitives index, lineage matrix, and code.
6. Use Camel-Case with hyphens for all vault file names.
7. Run `./scripts/lint_vault.sh` from the repo root; make sure it exits with code 0.
8. Commit: `feat(vault): sync research database and adhoc reports`.

## 4. Run Empirical Experiments

1. Run translation benchmarks on GildedRose (C to Rust).
2. Record binary compilation status as `Pass` or `Fail`; no partial-compile results.
3. Record exact test pass count and percentage ($N/M$, $P\%$).
4. Never mock, fake, or skip experiments.
5. Write results to `.obsidian/MAgHARCM/adhoc/Benchmark-Results-And-Evaluation.md` and `docs/.paper/sec_eval.tex`.
6. Commit: `test(benchmarks): record empirical evaluation results`.

## 5. Implement Primitives and Defend the Codebase

### 5.1 Generative Implementation

1. Keep code aligned with admitted primitives and current research progress.
2. Apply ponytail and eino principles: prefer the Go standard library; delete dead code, unused abstractions, and redundant shims; keep implementations minimal, readable, and executable.
3. Document design intent inline with comments referencing primitive IDs and paper bibkeys.

### 5.2 Defensive Invariants

1. Initialization and configuration: apply the `Must` pattern for compile-time initialization (no silent fallbacks); read runtime settings from YAML in `configs/`.
2. Architecture:
   - Shared state lives in `internal/compiletime/state.go`; producer packages use type aliases to satisfy ADR-C-014 locality of behavior without Go import cycles.
   - Wire all eight agents in `internal/graph/graph.go` (cyclic repair loop): Archaeologist, Analyzer, Planner, Translator, Reviewer, VerdictPanel, Validator, Recruiter.
   - Keep the five Mueller strategies in canonical order in `internal/agents/strategy.go`: Big Bang, Pilot, Frozen Legacy, Parallel Cutover, Incremental.
   - Log through `internal/logger`; never use `fmt.Print*` in production Go code.
3. Resolve blockers from `Human-Intervention-And-Blockers.md`; run `go test ./...` and make sure every test passes.
4. Commit: `feat(agents): implement PRIM-XX from P-YY and defend invariants`.

## 6. Synchronize the Academic Paper

1. Update `docs/.paper/sec_method.tex` and `main.tex` with newly admitted primitives and mechanisms.
2. Add verified BibTeX entries to `docs/.paper/refs.bib`.
3. Update `docs/.paper/sec_eval.tex` with verified benchmark numbers.
4. Make sure claims match across paper, codebase, and vault.
5. Commit: `docs(paper): update method and references`.

## 7. Evolve This Command File

1. Read `.omp/commands/MAgHARCM.md`.
2. Remove directives that describe completed work.
3. Add new research anchors, active watchlists, and resolved blockers.
4. Keep this command concise, exhaustive, and compliant with ASD-STE100.
5. Commit: `feat(command): evolve MAgHARCM sprint workflow`.

## 8. Persist Sprint Handoff

1. Write the sprint summary to `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md`.
2. Document closed tasks, commit sequence, audit tables, and active triggers.
3. Commit: `docs(diary): persist sprint YYYY-MM-DD handoff`.

## 9. Canonical Substrate and Anchor Matrix

Three closed-loop substrate stacks anchor the SLM comprehension, de-hallucination, and translation primitives: static-graph refinement (relational Datalog policies + context-conditioned slicing), systematic hallucination evaluation (perturbation benchmark + detection + decode-time re-ranking), and feedback-driven translation (single-LLM three-signal loop + multi-agent execution critic).

| Primitive | Canonical Anchors | Operational Mechanism |
| :--- | :--- | :--- |
| **PRIM-07** (SLM Verifier & Judge) | T1 (P-125), ARC-Decode (P-126), SLM-as-a-Judge (P-127), SPECS (P-135), CaTS (P-136), ContextPRM (P-146), HalluShield (P-143) | SLM test-time search, contextual-coherence PRM, value-model guidance. |
| **PRIM-09** (Hybrid Code Graph) | SSAR (P-130), SemArc (P-131), SemRef (P-132) | Semantic + structural edge weights, iterative LLM refinement. |
| **PRIM-12** (Static Analysis Co-Evolution) | LλMDA (P-129), CoReX (P-153), POLA-Tester (P-155) | Partial-PDG augmentation, regression slicing, dependency mining. |
| **PRIM-21 / PRIM-31** (KV Cache & Speculative Decoding) | KVzip (P-128), RelayCaching (P-134), SuffixDecoding (P-137), SpecKV (P-147), LookaheadKV (P-148), SSD/Saguaro (P-149) | Query-agnostic eviction, KV reuse, asynchronous draft-verify. |
| **PRIM-22 / PRIM-25** (Comprehension & De-Hallucination) | ADI (P-133), NESA (P-142), TestPrune (P-150), Hallu-Eval (P-152), CoReX (P-153), ACONITE (P-156) | Lifetime tracing, Datalog decomposition, coverage pruning, perturbation benchmark. |
| **PRIM-23 / PRIM-29** (Multi-Language Translation) | Syzygy (P-124), SmartC2Rust (P-151), TransAgent (P-154) | Three-signal feedback loop, execution-aligned critic pipeline. |

## 10. Standing Governance Rules

1. **Dating Rule (BLK-06)**: `date -u` or the system reminder is the only authoritative date; never extrapolate; YAML date must equal filename date.
2. **REJECT and Watchlist Registry (BLK-08)**: append every Q1/Q2/Q3 reject to `reject_registry.wave-NN` with a four-way rationale; track unverified candidates in `watchlist.wave-NN`.
3. **Sandbox Batch Commit (temporary)**: when git mutations are blocked, accumulate artifacts in the working tree and batch-commit on recovery; revert to per-phase commits immediately after.

## 11. Active Watchlist and Strategic Roadmap

1. **Program-Comprehension Mechanism Slot: CLOSED** by `[[1.0.0 P-153]]` CoReX (context-aware refinement slicing, Waves 17–22); scouts must not reopen broad comprehension queries.
2. **W23-W1 SWE-TRACE** (`arXiv:2604.14820`): re-verify venue after NeurIPS 2026 notifications (2026-09-24); sixth carry retires it to the reject archive.
3. **Substrate Wiring Check**: verify opt-in runtime wiring in `configs/agents.yml` for NESA (§11.6), Hallu-Eval triplet (§11.7), and SmartC2Rust feedback (§11.8).
4. **Wave-24 Cross-Pattern Formalization (speculative)**: LλMDA × TestPrune, TransAgent × SmartC2Rust, POLA-Tester × CoReX, ACONITE × Panta.
