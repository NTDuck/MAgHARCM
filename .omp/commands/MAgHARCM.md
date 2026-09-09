---
description: "Execute a research sprint, vault sync, empirical tests, code defense, and command evolution"
---

# MAgHARCM Research Sprint and Modernization Workflow

## 1. Start the Sprint

1. Read the system reminder or run `date -u` to get the authoritative current date in `YYYY-MM-DD` format.
2. Never extrapolate future dates from git history or older notes.
3. Make sure note metadata headers (`date:` and `last_updated:`) match the current date (rule BLK-06).
4. Set the sprint handoff path to `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md`.
5. If a handoff file exists for today, increment the suffix (`-1`, `-2`, `-N`).
6. Read the latest handoff note in `.obsidian/MAgHARCM/diary/` for pending tasks.
7. Read `.obsidian/MAgHARCM/adhoc/Human-Intervention-And-Blockers.md` to identify active blockers.
8. Delete stale temporary files from `.artifacts/local/` and test caches.
9. Commit: `docs(diary): start sprint YYYY-MM-DD`.

## 2. Research and Literature Expansion

1. Search for papers in top peer-reviewed venues: NeurIPS, ICML, ICLR, FSE, ICSE, ASE, TOSEM, and TSE.
2. Evaluate each candidate paper against three mandatory gates:
   - **Q1 (Venue)**: Is the paper in a recognized top venue? Workshop papers require single-paper method novelty and a non-redundant theme (P-138 threshold; empirical/diagnostic papers disqualified).
   - **Q2 (Mechanism)**: Does the paper define a concrete algorithmic mechanism rather than a prompt tweak?
   - **Q3 (Anchor)**: Does the mechanism anchor or defend a MAgHARCM primitive (`PRIM-01` to `PRIM-33`)?
3. Enforce the two-hop citation lineage for every admitted paper:
   - Hop 1: Direct mechanism support.
   - Hop 2: Foundational literature anchor.
4. Focus research on small language models (4B-30B) and software archaeology:
   - Program slicing, dynamic invariants, execution traces, concept assignment, dependency graphs.
   - KV cache compression, speculative decoding, SLM verification, test-time compute.
5. Record non-admitted candidates in the candidate memo and registries (rule BLK-08):
   - Save candidate triage to `.obsidian/MAgHARCM/research/diary/Wave-NN-Candidates.md`.
   - Record rejected candidates in `reject_registry.wave-NN` with bibkey, title, venue, verdict (Q1/Q2/Q3), and a four-way classification rationale:
     (a) hallucinated-venue, (b) off-list venue, (c) mechanism-overlap, or (d) off-axis target-language/problem.
   - Record unverified candidates in `watchlist.wave-NN` for re-verification in subsequent waves.
6. Never fabricate citations or metrics. Mark unverified candidates as placeholders in blockers.
7. Commit: `feat(research): add wave research papers`.

## 3. Synchronize the Obsidian Vault

1. Update `.obsidian/MAgHARCM/Research-Database.json` with all new paper and primitive entries:
   - Supply all required fields: `id`, `title`, `authors`, `year`, `venue`, `bibkey`, identifier (DOI, OpenReview ID, or arXiv ID), `citations`, and `summaries`.
   - Append rejected candidates to `reject_registry.wave-NN` and watchlisted candidates to `watchlist.wave-NN`.
   - Make sure JSON parses cleanly: `python3 -c "import json; json.load(open('.obsidian/MAgHARCM/Research-Database.json'))"`.
2. Save paper notes to `.obsidian/MAgHARCM/research/papers/P-NN-<Name>.md`:
   - Declare YAML aliases including the version marker (`[[1.0.0 P-NN]]`).
3. Update human reports in `.obsidian/MAgHARCM/adhoc/`:
   - `Methodology.md`: Update primitive anchors, methodology justifications, closed slots, and changelog.
   - `Project-Progress-And-Milestones.md`: Update metrics, wave tables, and milestone status.
   - `Strategic-Direction-And-Roadmap.md`: Update roadmap and wave plans.
   - `Human-Intervention-And-Blockers.md`: Update active blockers and resolved items.
   - `Benchmark-Results-And-Evaluation.md`: Update empirical evaluation data.
   - `Architecture-And-Dataflow.md`: Update architecture and dataflow diagrams.
   - `Research-Waves-Index.md`: Update the wave index ledger.
4. Update lineage and primitive indices:
   - Cross-reference papers in `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`.
   - Update `.obsidian/MAgHARCM/primitives/Primitives-Index.md`.
   - Maintain full parity across the primitives index, lineage matrix, and codebase.
5. Use Camel-Case with hyphens for all file names in the vault.
6. Run `./scripts/lint_vault.sh` from the repository root and make sure it exits with code 0.
7. Commit: `feat(vault): sync research database and adhoc reports`.

## 4. Run Empirical Experiments

1. Run translation benchmarks on target repositories:
   - GildedRose (C to Rust).
2. Record binary compilation status as `Pass` or `Fail`. Do not record partial compilation.
3. Record exact test pass counts and pass percentages ($N/M$, $P\%$).
4. Never mock, fake, or skip experiments.
5. Record experimental results in:
   - `.obsidian/MAgHARCM/adhoc/Benchmark-Results-And-Evaluation.md`.
   - `docs/.paper/sec_eval.tex`.
6. Commit: `test(benchmarks): record empirical evaluation results`.

## 5. Implement Primitives and Defend the Codebase

### 5.1 Generative Implementation

1. Keep codebase aligned with current research progress and admitted primitives.
2. Apply minimal code principles (ponytail and eino):
   - Prefer Go standard library over external dependencies.
   - Delete dead code, unused abstractions, and redundant shims.
   - Keep implementations minimal, readable, and directly executable.
3. Document all design decisions and intent inline with comments referencing primitive IDs and paper bibkeys for auditability.

### 5.2 Defensive Invariants

1. Enforce compile-time safety and architectural boundaries:
   - Apply the `Must` pattern for compile-time initializations. Reject silent fallbacks.
   - Read runtime configuration from external YAML files in `configs/`.
   - Declare shared state in `internal/compiletime/state.go`. Use type aliases in producer packages to satisfy ADR-C-014 locality of behavior while preventing Go import cycles.
   - Wire all eight agents in `internal/graph/graph.go` to maintain the cyclic repair loop:
     Archaeologist (AST/archaeology), Analyzer (signal profiling), Planner (decomposition), Translator (chunked translation), Reviewer (role-flip gate), VerdictPanel (consensus scoring), Validator (compilation/tests), and Recruiter (strategy iteration).
   - Maintain the five Mueller strategies in canonical order in `internal/agents/strategy.go`:
     Big Bang (tiny project), Pilot (very large), Frozen Legacy (untested), Parallel Cutover (modular with tests), and Incremental (default).
   - Use `internal/logger` for logging. Do not use `fmt.Print*` in production Go code.
2. Defend codebase integrity:
   - Resolve blockers listed in `Human-Intervention-And-Blockers.md`.
   - Run `go test ./...` and make sure all tests pass without errors.
3. Commit: `feat(agents): implement PRIM-XX from P-YY and defend invariants`.

## 6. Synchronize the Academic Paper

1. Update `docs/.paper/sec_method.tex` and `main.tex` with newly admitted research primitives and mechanisms.
2. Add verified BibTeX entries to `docs/.paper/refs.bib`.
3. Update `docs/.paper/sec_eval.tex` with verified benchmark numbers.
4. Make sure claims across paper, codebase, and vault match completely.
5. Commit: `docs(paper): update method and references`.

## 7. Evolve This Command File

1. Read `.omp/commands/MAgHARCM.md`.
2. Remove directives and watchlist entries that describe completed work.
3. Add newly admitted research anchors, active watchlists, and resolved blockers.
4. Keep this command concise, exhaustive, and compliant with ASD-STE100.
5. Commit: `feat(command): evolve MAgHARCM sprint workflow`.

## 8. Persist Sprint Handoff

1. Write the sprint summary to `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md`.
2. Document closed tasks, commit sequence, audit tables, and active triggers.
3. Commit: `docs(diary): persist sprint YYYY-MM-DD handoff`.

## 9. Canonical Substrate and Anchor Matrix

| Primitive | Canonical Substrate Anchors | Mechanism and Operational Role |
| :--- | :--- | :--- |
| **PRIM-07** (SLM Verifier & Judge) | T1 (P-125), ARC-Decode (P-126), SLM-as-a-Judge (P-127), SPECS (P-135), CaTS (P-136), ContextPRM (P-146), HalluShield (P-143) | SLM test-time search, contextual-coherence PRM, value-model guidance for closed-loop verifiers. |
| **PRIM-09** (Hybrid Code Graph) | SSAR (P-130), SemArc (P-131), SemRef (P-132) | Semantic similarity + structural dependency edge weights, canonical pattern base, iterative LLM refinement. |
| **PRIM-12** (Static Analysis Co-Evolution) | LλMDA (P-129), CoReX (P-153), POLA-Tester (P-155) | Partial-PDG context augmentation, context-aware regression slicing, agentic syntactic-dependency mining. |
| **PRIM-21 / PRIM-31** (KV Cache & Speculative Decoding) | KVzip (P-128), RelayCaching (P-134), SuffixDecoding (P-137), SpecKV (P-147), LookaheadKV (P-148), SSD/Saguaro (P-149) | Query-agnostic KV eviction, cross-agent KV cache reuse, LoRA-module eviction, asynchronous draft-verify pipelines. |
| **PRIM-22 / PRIM-25** (Comprehension & De-Hallucination) | ADI (P-133), NESA (P-142), TestPrune (P-150), Hallu-Eval / Hallu-Shield (P-152), CoReX (P-153), ACONITE (P-156) | Function lifetime tracing, relational Datalog policies, coverage-driven test pruning, semantic perturbation benchmark, backward slicing. |
| **PRIM-23 / PRIM-29** (Multi-Language Translation) | Syzygy (P-124), SmartC2Rust (P-151), TransAgent (P-154) | Dual-mode translation: single-LLM 3-signal iterative feedback (C/Go to Rust) + multi-agent execution-aligned critic. |

## 10. Three-Deep Substrate Reinforcement Architecture

The SLM comprehension, de-hallucination, and translation stack is anchored by three reinforced closed-loop substrates:

1. **Static-Graph Refinement (§11.6)**: Relational Datalog policy decomposition (P-142 NESA) combined with context-conditioned regression slicing (P-153 CoReX) anchors deterministic syntactic and semantic decomposition for `PRIM-09` and `PRIM-12`.
2. **Systematic Hallucination-Evaluation Triplet (§11.7)**: 800-pair semantic perturbation benchmark (Hallu-Eval), mutation-aware detection (Hallu-Det), and decode-time value-model re-ranking (Hallu-Shield) anchor `PRIM-22` and `PRIM-25`.
3. **Feedback-Driven Multi-Language Translation (§11.8)**: Dual-mode architecture combining single-LLM three-signal feedback (P-151 SmartC2Rust: compiler errors, AST diffs, residual unsafe blocks) with multi-agent execution-aligned critic (P-154 TransAgent) anchors `PRIM-23` and `PRIM-29`.

## 11. Standing Governance Rules

1. **Dating Rule (BLK-06)**: Use `date -u` or system reminder as authoritative date. Never extrapolate future dates. Ensure YAML metadata date matches filename date.
2. **REJECT and Watchlist Registry (BLK-08)**: Every candidate rejected at Q1/Q2/Q3 must be appended to `reject_registry.wave-NN` in `Research-Database.json` with bibkey, title, venue, verdict, and classification rationale. Unverified candidates must be tracked in `watchlist.wave-NN`.
3. **Temporary Governance Rule (Wave-21+ Sandbox Restriction)**: When environment restrictions block git operations during a sprint, accumulate all artifacts cleanly in the working tree and execute the conventional batch commit when the sandbox recovers. Revert to standard per-phase commits immediately upon recovery.

## 12. Active Watchlist and Strategic Forward Roadmap

1. **Strict Program-Comprehension Mechanism Slot (CLOSED)**: `[[1.0.0 P-153]]` CoReX (context-aware refinement slicing) definitively closed the residual program-comprehension slot carried through Waves 17–22. Autonomous scouts must not reopen broad comprehension search queries.
2. **W23-W1 SWE-TRACE** (`arXiv:2604.14820`): Re-verify peer-reviewed venue status following NeurIPS 2026 author notifications (2026-09-24). Sixth carry triggers automatic retirement to reject archive.
3. **Substrate Configuration Wiring**: Verify opt-in runtime wiring in `configs/agents.yml` for NESA (§11.6), Hallu-Eval triplet (§11.7), and SmartC2Rust feedback (§11.8).
4. **Cross-Pattern Formalization (Wave-24 Roadmap)**:
   - LλMDA (P-129) × TestPrune (P-150): Partial PDG context augmentation with coverage-driven input reduction.
   - TransAgent (P-154) × SmartC2Rust (P-151): Dual-mode translation integrating single-LLM feedback with execution critics.
   - POLA-Tester (P-155) × CoReX (P-153): Syntactic dependency mining paired with refinement slicing.
   - ACONITE (P-156) × Panta (P-140): Backward slicing test generation combined with hybrid static/dynamic analysis.
