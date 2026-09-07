---
title: Sprint 2026-09-07 Handoff (iter-2, Wave-19)
backlink: "[[3.0.0 Sprint 2026-09-07 Handoff-2]]"
date: 2026-09-07
last_updated: 2026-09-07
tags: [sprint, handoff, wave-19, [[2.0.0 MAgHARCM]], slm, kv-cache, speculative-decoding, "[[1.0.0 PRIM-7]]", "[[1.0.0 PRIM-9]]", "[[1.0.0 PRIM-21]]", "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-31]]"]
---

# [[3.0.0 Sprint 2026-09-07 Handoff-2 — Wave-19 Anchors + Standing Rules]]

## Outcome

Wave-19 fired on top of Wave-18. **8 ACCEPT + 3 REJECT + 2 UNVERIFIED = 13 triaged.** 17 paper notes persisted (P-125..P-141). REJECT registry encoded in `Research-Database.json` (BLK-08 resolved). Dating convention reasserted (BLK-06). `go test ./...` green. `scripts/lint_vault.sh` clean. Vault + paper + command all synced. Tree clean.

## Wave-19 (closed)

### Accepted (8)

- **P-134 RelayCaching** (Geng et al., ICML 2026 Poster #1915) — cross-agent KV cache reuse via sparse deviation recompute (>80% reuse, 4.7× TTFT). Anchors PRIM-31 + PRIM-9.
- **P-135 SPECS** (Cemri et al., ICLR 2026) — speculative drafts + soft verification + dynamic switch. Closes Wave-18 R3. Anchors PRIM-7 + PRIM-21.
- **P-136 CaTS** (Huang et al., ICLR 2026 Poster) — Self-Calibration distilled confidence. Anchors PRIM-7 as frontier-PRM replacement.
- **P-137 SuffixDecoding** (Oliaro et al., NeurIPS 2025 Spotlight) — model-free suffix-tree draft, deployed in Snowflake ArcticInference + vLLM. Anchors PRIM-7 + PRIM-21.
- **P-138 RepairKV** (Rusli et al., ICML 2026 AdaptFM Workshop) — borderline workshop-track ACCEPT per §7 method-level threshold. Post-compression KV repair. Anchors PRIM-22 + PRIM-31.
- **P-139 TypePro** (Lin et al., FSE 2026) — SDG + inter-procedural backward/forward slicing. Anchors PRIM-9 + PRIM-22.
- **P-140 Panta** (Gu, Nashid & Mesbah, ICSE 2026) — static (cyclomatic) + dynamic (coverage) iterative loop. Anchors PRIM-22 + PRIM-21.
- **P-141 KVFlow** (NeurIPS 2025 Poster) — Agent Step Graph + steps-to-execution metric + prefetching. Anchors PRIM-31 + PRIM-21.

### Rejected (3)

- **R1 TTA*** (Braverman, Zhang & Gu, NeurIPS 2025 LAW Workshop) — REJECTED Q1 (workshop redundancy vs P-135 + P-136).
- **R2 HELIOS** (Achamyeleh, Thomare & Al Faruque, NDSS 2026 LAST-X Workshop) — REJECTED Q1 (off-list venue + off-axis binary-decompilation target).
- **R3 LongSpec** (ACL 2026 Main) — REJECTED Q1 (off-list venue; same rationale as Wave-18 R1 ReflexiCoder).

### Watchlist — UNVERIFIED (2)

- **W1 SliceMate** (Chang et al., arXiv:2507.18957). Yunbo Lyu's homepage claims ISSTA 2026 acceptance, but no ISSTA 2026 program slot for SliceMate is listed on conf.researchr.org. Wave-20 re-verification.
- **W2 SWE-TRACE** (Han et al., arXiv:2604.14820). April 2026 arXiv preprint only; no peer-reviewed venue confirmation. Wave-20 re-verification.

## Standing Rules Encoded (closed)

- **BLK-08 (REJECT registry)** — encoded in `.omp/commands/MAgHARCM.md §13`. Every sprint MUST append the prior wave's REJECT list to `Research-Database.json` under `reject_registry.wave-NN`. Resolved.
- **BLK-06 (Dating convention)** — encoded in `.omp/commands/MAgHARCM.md §14`. `date -u` is the authoritative date; metadata header MUST equal filename date. Resolved.
- **Wave-20 watchlist** — encoded in `.omp/commands/MAgHARCM.md §15`. SliceMate, SWE-TRACE, and the strict program-comprehension-mechanism slot are the next-wave priorities.

## Commits (in order)

1. `docs(diary): start sprint 2026-09-07 (iter-2, wave-19)`
2. `feat(research): wave-19 ACCEPT papers P-134..P-141 + REJECT registry`
3. `feat(vault): wave-19 entries in Methodology + Research-Database.json + Research-Waves-Index`
4. `feat(vault): wave-19 sync into Architecture/Lineage/Progress/Blockers`
5. `fix(codebase): defend architecture invariants (no fmt.Print* regression, 8-agent graph, go test ./... green)`
6. `docs(paper): add wave-18 + wave-19 BibTeX + cite clusters`
7. `feat(command): encode wave-19 insights, REJECT registry (BLK-08), dating convention (BLK-06), wave-20 watchlist`

## Active Triggers

- **Wave-20 priority**: software-archaeology strict-mechanism slot. Carry an explicit `program-comprehension-mechanism` query against NeurIPS 2026 / ICML 2027 / ICLR 2027 listings. Also re-verify SliceMate + SWE-TRACE venue confirmations.
- **BLK-08 standing rule**: next sprint appends `reject_registry.wave-20` to `Research-Database.json`.
- **BLK-06 standing rule**: next sprint uses `date -u`; metadata header MUST equal filename date.

## Verification

- `go test ./...` — green (cached, no new test files).
- `./scripts/lint_vault.sh` — clean (252 files scanned, 17 refs checked, exit 0).
- `Research-Database.json` — parses cleanly; `reject_registry.wave-19` added with 3 entries.
- Paper notes P-125..P-141 — all in `.obsidian/MAgHARCM/research/papers/`.
- Wave-19 candidates memo — `.obsidian/MAgHARCM/research/diary/Wave-19-Candidates.md`.

## Out-of-Scope / Not Actioned

- **Empirical Experiments (§4)**: GildedRose / Gohistogram / Stats / Commons-Validator benchmarks not re-run this sprint. Harness requires live SLM API keys + bench fixture repositories under `tests/`; neither is provisioned in this environment. sec_eval.tex numbers remain at Sprint 2026-09-28 levels (GildedRose 14/14 100.0%, Gohistogram 5/8 62.5%, Stats 24/42 57.1%, Commons-Validator 18/68 26.5% Fail). No number update was warranted without a fresh run.
