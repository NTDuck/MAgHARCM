---
description: "Execute a research sprint, vault sync, empirical tests, code defense, and command evolution"
---

# MAgHARCM Research Sprint and Modernization Workflow

Run a MAgHARCM research and codebase iteration.
Optional focus argument: `$ARGUMENTS`

## 1. Start the Sprint

1. Get the current date in YYYY-MM-DD format from the system reminder or `date -u`.
2. Do not extrapolate future dates from git history or older notes.
3. Set the sprint handoff path to `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md`.
4. If a handoff file exists for today, append `-1`, `-2`, or the next number.

5. Read the previous sprint handoff note in `.obsidian/MAgHARCM/diary/` for open tasks.
6. Read `.obsidian/MAgHARCM/adhoc/Human-Intervention-And-Blockers.md` to identify active blockers.
7. Delete stale temporary files from `.artifacts/local/` and test caches.
8. Commit: `docs(diary): start sprint YYYY-MM-DD`.

## 2. Research and Literature Expansion

1. Search for papers in major venues like NeurIPS, ICML, ICLR, FSE, ICSE, ASE, and TOSEM.
2. Examine each candidate paper with three questions:
   - Is the paper in a peer-reviewed venue or an established archive?
   - Does it describe a concrete mechanism instead of a prompt tweak?
   - Does it anchor or defend an existing or new MAgHARCM primitive?

3. For each admitted paper, record two citation hops:
   - Hop 1 supports the mechanism.
   - Hop 2 supports the foundational sources.

4. Focus research on software archaeology:
   - Program slicing, dynamic invariants, and concept assignment.
   - Temporal coupling, dependency structure matrices, and design rule hierarchy partitions.

5. Focus research on small language models with 4B to 30B parameters:
   - Structured output cloze slots, KV cache compression, and speculative decoding.
   - Test-time scaling search, distilled self-reflection, and self-consistency.

6. Do not fabricate citations.
7. If you cannot verify a source, label it as an unverified placeholder.
8. Add unverified sources to `.obsidian/MAgHARCM/adhoc/Human-Intervention-And-Blockers.md`.
9. Commit: `feat(research): add wave research papers`.

## 3. Synchronize the Obsidian Vault

1. Update `.obsidian/MAgHARCM/Research-Database.json` with all new paper and primitive entries.
2. Fill all fields: id, title, authors, year, venue, bibkey, DOI, citations, and summaries.
3. Save paper notes in `.obsidian/MAgHARCM/research/papers/P-NN-<Name>.md`.
4. Add YAML aliases in each paper note so wikilinks resolve in Obsidian.

5. Update the human reports in `.obsidian/MAgHARCM/adhoc/`:
   - `adhoc/Methodology.md` updates anchors and the changelog.
   - `adhoc/Project-Progress-And-Milestones.md` updates metrics and milestone status.
   - `adhoc/Strategic-Direction-And-Roadmap.md` updates roadmaps and wave plans.
   - `adhoc/Human-Intervention-And-Blockers.md` updates active blockers.

6. Update remaining human reports in `.obsidian/MAgHARCM/adhoc/`:
   - `adhoc/Benchmark-Results-And-Evaluation.md` updates empirical data.
   - `adhoc/Architecture-And-Dataflow.md` updates state and graph diagrams.
   - `adhoc/Research-Waves-Index.md` updates the wave index.

7. Cross-reference new papers in `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`.
8. Update `.obsidian/MAgHARCM/primitives/Primitives-Index.md` with primitive changes.
9. Maintain full parity across the primitives index, lineage matrix, and code.

10. Use Camel-Case with hyphens for every file name in the vault.
11. Run `./scripts/lint_vault.sh` from the repository root.
12. Make sure the script exits with code 0.
13. Commit: `feat(vault): sync research database and adhoc reports`.

## 4. Run Empirical Experiments

1. Test the system on the four benchmark repositories:
   - GildedRose (C to Rust).
   - Gohistogram (Go to Rust).
   - Stats (Go to Rust).
   - Commons-Validator (Java to Rust).

2. Record the binary compilation status as Pass or Fail. Do not record partial compilation.
3. Record the exact test pass count and test pass percentage.
4. Do not fake, mock, or skip experiments.

5. Record experimental results in `.obsidian/MAgHARCM/adhoc/Benchmark-Results-And-Evaluation.md`.
6. Record experimental results in `docs/.paper/sec_eval.tex`.
7. Commit: `test(benchmarks): record empirical evaluation results`.

## 5. Defend the Codebase and Invariants

1. Make sure compile-time initializations use the Must pattern and panic on error.
2. Do not use fallbacks in compile-time configuration.
3. Read runtime settings from external YAML files in `configs/`.
4. Keep shared state declarations in `internal/compiletime/state.go`.

5. Keep agent artifact types in their producer files. Use type aliases to avoid cycles.
6. Use the try-and-fail strategy registry in `internal/agents/strategy.go`.
7. Keep all eight agents wired in `internal/graph/graph.go`.
8. Use idiomatic Charm libraries in `cmd/MAgHARCM-tui`.

9. Do not use `fmt.Print*` in production Go code. Use `internal/logger` instead.
10. Use plain technical English in messages per ASD-STE100.
11. Resolve blockers listed in `.obsidian/MAgHARCM/adhoc/Human-Intervention-And-Blockers.md`.
12. Run `go test ./...` and fix any build breaks.
13. Commit: `fix(codebase): defend architecture invariants and resolve blockers`.

## 6. Synchronize the Academic Paper

1. Update `docs/.paper/sec_method.tex` and `main.tex` with citations to new papers.
2. Add BibTeX entries to `docs/.paper/refs.bib`.
3. Update `docs/.paper/sec_eval.tex` with new benchmark numbers.
4. Make sure claims in the paper match the codebase and the vault.
5. Commit: `docs(paper): update method and references`.

## 7. Evolve This Command File

1. Read `.omp/commands/MAgHARCM.md`.
2. Remove directives that describe completed work.
3. Add new research insights, newly admitted anchors, or resolved blockers.
4. Keep this command concise, exhaustive, and compliant with ASD-STE100.
5. Commit: `feat(command): evolve MAgHARCM sprint workflow`.

## 8. Persist Sprint Handoff

1. Write the sprint summary to `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md`.
2. Document closed tasks, commits, audit tables, and active triggers.
3. Commit: `docs(diary): persist sprint YYYY-MM-DD handoff`.

## 9. Wave-18 Insights (2026-09-07)

1. **SLM-as-judge is the default verifier.** Frontier-PRM-as-judge (P-91..P-92, P-97..P-98) is superseded for 4B-30B substrates by T1 (P-125), ARC-Decode (P-126), and SLM-as-a-Judge (P-127). Promote these as the canonical PRIM-7 anchor in future waves.
2. **Architecture-recovery trio (SSAR + SemArc + SemRef).** P-130, P-131, P-132 form the canonical pattern for PRIM-9 Tri-Representation Hybrid Code Graph. Edge weights combine semantic similarity + structural dependency; canonical-pattern knowledge base anchors the partition; iterative LLM refinement closes the loop.
3. **Query-agnostic KV cache compression for multi-step agents.** P-128 KVzip replaces P-80 StreamingLLM and P-105 ChunkKV as the right substrate for agent loops because the eviction is independent of the future query.
4. **Function-level dynamic analysis for debugging agents.** P-133 ADI's Frame Lifetime Trace + high-level navigational commands is the right pattern for PRIM-22 observation/structure phases; line-level gdb-style DA is the wrong substrate for SLM-era agents.
5. **Partial-program dependence analysis as LLM-augmented static analysis.** P-129 LλMDA's "context-augment partial program then run classical DA" pattern generalizes beyond dependence graphs to slicing and concept assignment.

## 10. Evidence and Citation Discipline

1. Every accepted paper must have a `bibkey`, a venue, and either a DOI, OpenReview ID, or arXiv ID — never a bare URL.
2. Every accepted paper must record two citation hops (mechanism-support + foundational) in `Research-Database.json` and the paper note.
3. Every paper note must declare YAML aliases including the version marker (`[[1.0.0 P-NN]]`).
4. Every REJECT must be logged in `Wave-NN-Candidates.md` with the Q1/Q2/Q3 verdict and a one-line rationale.
5. `Research-Database.json` MUST parse cleanly with `python3 -c "import json; json.load(...)"` after every edit.

## 11. Wave Memo Discipline

1. Each research wave writes a `.obsidian/MAgHARCM/research/diary/Wave-NN-Candidates.md` memo before any paper notes are persisted.
2. The memo is the triage ledger: every fired candidate appears with verdict + rationale.
3. The memo's bibkeys MUST match the bibkeys written to `refs.bib` and `Research-Database.json`.
3. The memo's bibkeys MUST match the bibkeys written to `refs.bib` and `Research-Database.json`.
4. The memo's ACCEPT list MUST equal the new entries appended to `Research-Database.json`.

## 12. Wave-19 Insights (2026-09-07 iter-2)

1. **Speculative-decoding × KV-cache hybrid is now canonical.** P-134 RelayCaching (cross-agent KV reuse), P-137 SuffixDecoding (model-free suffix-tree draft), P-138 RepairKV (post-compression KV repair), and P-141 KVFlow (workflow-aware KV cache eviction) close the speculative-decoding × KV-cache gap alongside P-126 ARC-Decode + P-128 KVzip. PRIM-21 + PRIM-31 are now both anchored by this family; future waves MUST cite at least one P-134 / P-137 / P-141 paper in any SLM-loop architecture claim.
2. **SLM-scale TTS is convergent.** P-135 SPECS + P-136 CaTS close the frontier-PRM-as-judge gap left by P-91..P-92. Combined with P-125 T1 + P-127 SLM-as-a-Judge, the four-paper cluster (T1, ARC-Decode, SPECS, CaTS) forms the canonical PRIM-7 anchor for SLM substrates.
3. **Software-archaeology strict-mechanism gap is partially closed.** P-129 LλMDA (LLM-aided partial PDG) + P-133 ADI (Frame Lifetime Trace) close the program-comprehension-mechanism gap. LLM-augmented static analysis beyond dependence graphs is now anchored (P-139 TypePro for type inference via inter-procedural slicing; P-140 Panta for iterative hybrid static+dynamic test generation). The strict program-comprehension-mechanism slot remains open and is the Wave-20 priority.

## 13. REJECT Registry Rule (BLK-08)

1. Every sprint MUST append the prior wave's REJECT list to `Research-Database.json` under `reject_registry.wave-NN`.
2. Each REJECT entry records: `bibkey`, `title`, `venue`, `verdict` (Q1/Q2/Q3), and a one-line `rationale`.
3. The REJECT registry is the authoritative cross-wave triage ledger; older Wave-NN-Candidates.md memos remain, but the registry is the lookup of record.
4. Watchlist (UNVERIFIED) entries go to `watchlist.wave-NN` and are re-verified each wave.

## 14. Dating Convention (BLK-06)

1. Every sprint MUST use `date -u` (or the system reminder's date) as the authoritative date.
2. If a sprint is run later than expected and a handoff filename is in the future, the next sprint MUST reset the metadata header (`date:` + `last_updated:`) to the actual current date and document the rationale in `Methodology.md §6`.
3. No future-dating permitted under any circumstance; the metadata header MUST equal the filename date after the reset.

## 15. Wave-20 Watchlist

1. **SliceMate venue re-verification** (W1 from Wave-19). Yunbo Lyu's homepage claims ISSTA 2026 acceptance, but no ISSTA 2026 program slot for SliceMate is listed on conf.researchr.org. Re-verify in Wave-20.
2. **SWE-TRACE venue confirmation** (W2 from Wave-19). arXiv:2604.14820 — April 2026 arXiv preprint only; no peer-reviewed venue confirmation. Re-verify in Wave-20.
3. **Strict program-comprehension-mechanism slot** — Wave-20 scout carries an explicit `program-comprehension-mechanism` query against NeurIPS 2026 / ICML 2027 / ICLR 2027 listings.

## 16. Wave-20 Insights (2026-09-07 iter-3)

1. **Program-comprehension-mechanism slot is partially closed.** `[[1.0.0 P-142]]` NESA (Wang et al., FSE 2026, DOI 10.1145/3808161) anchors relational neuro-symbolic static program analysis: a restricted Datalog analysis-policy language decomposes complex comprehension sub-problems into deterministic syntactic slices (handled by parsing-based analysis) and LLM-handled semantic slices. Combined with the Wave-18 anchors (P-130 SSAR + P-131 SemArc + P-132 SemRef architecture-recovery trio + P-133 ADI function-level DA), the comprehension pipeline is now: static CPG → ADI function-level DA → NESA Datalog-policy decomposition → SSAR alignment → SemArc partition → SemRef iterative LLM refinement. The strict program-comprehension-mechanism slot is partially closed; the residual gap (function-level → partition-aligned summary) carries into Wave-21.
2. **Domain-agnostic PRM is the canonical SLM-loop judge.** `[[1.0.0 P-146]]` ContextPRM (Zhang et al., ICLR 2026, OpenReview 10011128) anchors domain-agnostic contextual-coherence PRM alongside the Wave-19 P-135 SPECS + P-136 CaTS convergence. Any future judge design MUST cite at least one of these four (T1, ARC-Decode, SPECS, CaTS) + ContextPRM.
3. **Inference-time hallucination defence anchors PRIM-7.** `[[1.0.0 P-143]]` HalluShield (Wan et al., FSE 2026, DOI 10.1145/3808139) provides inference-time value-model guidance for LLM-based code summarization, with Hallu-Det (entity-level detection + synonymous-mutation refinement) upstream. The PRIM-7 substrate is now closed-loop: judgement (T1 / SLM-as-Judge) → defence (HalluShield inference-time guidance) → validation (SPECS / CaTS / ContextPRM).
4. **Knowledge-augmented strategy selection anchors PRIM-21.** `[[1.0.0 P-145]]` TerraMod (Gupta et al., IBM Research, ICSE 2026 NIER, DOI 10.1145/3786582.3786841) closes the knowledge-augmented migration-context gap for PRIM-21. The PRIM-21 substrate is now closed-loop: try-and-fail registry (P-122 ReasoningBank) → strategy selection (P-145 TerraMod external-knowledge migration context / P-123 CodeChemist) → execution (P-137 SuffixDecoding) → judgement (P-143 HalluShield).
5. **Trace-driven multi-agent repair anchors PRIM-29 + PRIM-31.** `[[1.0.0 P-144]]` TraceCoder (Huang et al., ICSE 2026, arXiv 2602.06875) introduces trace-driven multi-agent repair with Historical Lesson Learning + Rollback Mechanism for LLM-generated code. Future agent training pipelines MUST cite P-144 alongside P-122 ReasoningBank (trajectory replay) and P-117 RepoCoder (retrieval-augmented).
6. **Workshop-track ACCEPT threshold documented.** P-138 RepairKV (AdaptFM Workshop ICML 2026) is the borderline case; future workshops on the §7 trigger-list are evaluated against P-138's threshold (method-level, single-paper, not workshop-redundant theme). Borderline REJECT examples: R1 TTA* NeurIPS 2025 LAW Workshop (workshop redundancy vs P-135/P-136), R2 HELIOS NDSS 2026 LAST-X Workshop (off-list venue).

## 17. Wave-21 Watchlist

1. **Program-comprehension-mechanism residual gap.** P-142 NESA partially closes the slot; the residual gap is the function-level → partition-aligned summary pass. Re-scout NeurIPS 2026 / ICML 2027 / ICLR 2027 listings.
2. **U1 NSE ICML 2026** (placeholder venue from Wave-20). Re-verify venue confirmation in Wave-21; if still missing, log as REJECT and remove from watchlist.
3. **Workshop-track ACCEPT threshold (P-138 RepairKV precedent).** Future borderline workshops on §7 trigger-list MUST be evaluated against P-138's threshold (method-level, single-paper, not workshop-redundant theme).
4. **§11.6 carry-over.** The §11.6 Program-Comprehension Mechanism substrate claim is opt-in via `configs/agents.yml:comprehension.graph_self_evolving: true`; Wave-21 should re-verify the opt-in path is wired (configs may need a corresponding field added).
5. **R1 Nexus ICSE 2026 + R2 SWE-Lego ICSE 2026 NIER** (Wave-20 REJECTS). Re-scout for follow-up mechanism papers from the same research groups; if no follow-up lands, both remain REJECTED in the registry.

