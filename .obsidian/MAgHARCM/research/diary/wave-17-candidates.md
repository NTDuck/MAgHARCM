---
title: Wave-17 Candidate Evaluation (2026-09-28)
backlink: "[[1.0.0 Wave-17 Candidates]]"
tags: [candidates, wave-17, mechanism-vs-benchmark, "[[1.0.0 P-123]]", "[[1.0.0 P-124]]", "[[1.0.0 P-104]]", "[[1.0.0 P-88]]", "[[1.0.0 P-91]]", "[[1.0.0 P-98]]", "[[1.0.0 P-122]]", test-time-scaling, dynamic-analysis, type-aware-translation, persistent-memory, tool-call-verification]
status: FIRED
date: 2026-09-28
last_updated: 2026-09-28
---

# [[1.0.0 Wave-17 Candidates]]

## Trigger-gate evaluation

§7 trigger gate rewritten 2026-09-25 (Q1 venue + Q2 mechanism-vs-benchmark + Q3 anchoring-to-existing-primitive).

## Candidates (5 triaged, 2 ACCEPT + 3 REJECT)

### Candidate 1 — CodeChemist (Wang et al. 2026, ICML 2026)

**Status: ACCEPT.** (Per scout `ResearchCodeReasoning`.)

- **Title**: CodeChemist: Test-Time Scaling for Low-Resource Code Generation via Functional Knowledge Transfer.
- **arXiv ID**: 2510.00501.
- **Venue**: ICML 2026 (peer-reviewed).
- **Q2 Mechanism?** YES — training-free test-time scaling via multi-temperature hedged sampling + cross-lingual I/O test oracle (functional knowledge transfer).
- **Q3 Anchoring?** YES — closes the low-resource + training-free gaps for `[[1.0.0 PRIM-21]]` Migration Strategy Selection, `[[1.0.0 PRIM-23]]` Chunked Translation, `[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection.
- **SLM relevance**: demonstrated on Qwen-1.5B (1.5B-parameter, training-free, compute-bound).
- **Hop-1**: `[[1.0.0 P-91]]` Snell 2024 (ICLR 2025), `[[1.0.0 P-104]]` S\* 2025 (EMNLP Findings), `[[1.0.0 P-98]]` Large Language Monkeys 2024.
- **Hop-2**: Wei 2022 CoT; Yao 2023 Tree of Thoughts.

### Candidate 2 — Syzygy (Shetty et al. 2025, ICLR 2025 VerifAI Workshop)

**Status: ACCEPT.** (Per scout `ResearchASTTranslation`.)

- **Title**: Syzygy: Dual Code-Test C to (safe) Rust Translation using LLMs and Dynamic Analysis.
- **arXiv ID**: 2412.14234.
- **Venue**: ICLR 2025 VerifAI Workshop (peer-reviewed, method-level; satisfies §7 Q1 workshop-track exemption).
- **Q2 Mechanism?** YES — Clang/LLVM-instrumented dynamic-analysis specification mining + dual code+test LLM generation + I/O-equivalence validation; multi-round repair loop.
- **Q3 Anchoring?** YES — closes the **dynamic** type/safety substrate for `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph, `[[1.0.0 PRIM-22]]` Four Phases of Comprehension, `[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter. Complements `[[1.0.0 P-88]]` HiTyper's static TDG.
- **SLM relevance**: bounded per-translation-unit prompt budget, model-agnostic; 4B-30B SLMs feasible.
- **Hop-1**: `[[1.0.0 P-88]]` HiTyper (Peng ICSE 2022), C2Rust, VERT.
- **Hop-2**: Le Goues 2019 automated program repair; LLVM instrumentation framework.

### Candidate 3 — MemSearcher (Yuan et al. 2026, ACL 2026 Findings)

**Status: REJECT Q1.** (Per scout `ResearchMemSearcher`.)

- ACL 2026 Findings is off the trigger-gate's approved venue list (2025+ NeurIPS/ICML/ICLR). Withdrawn from ICLR 2026 submission. No method-level companion paper at an approved venue.

### Candidate 4 — Verified Tool Calls (Mansoor et al. 2026, arXiv:2608.02645)

**Status: REJECT Q1.** (Per scout `ResearchToolVerify`.)

- Strongest Q3 anchor to `[[1.0.0 PRIM-7]]` in the candidate class (postcondition verification + idempotency + verify-before-retry map directly to PRIM-7 verdict validation), but arXiv-only with no confirmed NeurIPS/ICML/ICLR acceptance. ToolACE (ICLR 2025) and CoSC (ICLR 2025) reviewed but fail Q2 (data-QC pipeline / execution feedback, not tool-call verdict validation).

### Candidate 5 — Software Archaeology / Program Comprehension Mechanism Candidates

**Status: REJECT Q2/Q3.** (Per scout `ResearchSoftwareArch`.)

- Closest candidate is the ICML 2025 "Can LLMs Understand Intermediate Representations in Compilers?" (Jiang et al., arXiv:2502.06854) — an empirical benchmark, not a mechanism. No program-comprehension / archaeology mechanism paper at 2025+ ICLR/ICML/NeurIPS found. **Vault position**: no SLM-era mechanism anchor available for software archaeology as of 2026-09-28; the existing `[[1.0.0 P-87]]` TOSEM SLR remains the literature anchor.

## Verdict per candidate

| ID | Title | Q1 | Q2 | Q3 | Verdict |
| :--- | :--- | :--- | :--- | :--- | :--- |
| CodeChemist | Training-Free Test-Time Scaling for Low-Resource Code Generation | PASS | PASS | PASS | **ACCEPT** |
| Syzygy | Dual Code-Test C-to-Safe-Rust Translation via LLMs + Dynamic Analysis | PASS (workshop) | PASS | PASS | **ACCEPT** |
| MemSearcher | RL-trained Search Agent with Iterative Compact Memory | FAIL (ACL 2026 Findings) | n/a | n/a | **REJECT Q1** |
| Verified Tool Calls | Postcondition Verification for Tool Calls | FAIL (arXiv-only) | PASS | PASS | **REJECT Q1** |
| LLM-IR (Jiang ICML 2025) | Empirical benchmark, no mechanism | PASS | FAIL | FAIL | **REJECT Q2/Q3** |

## Wave-17 outcome

- **2 new P-NN anchors added**: `[[1.0.0 P-123]]` CodeChemist (ICML 2026) + `[[1.0.0 P-124]]` Syzygy (ICLR 2025 workshop).
- **3 rejections** logged with per-candidate rationale.
- **Vault coverage gains**: PRIM-21/23/27 (CodeChemist), PRIM-9/22/30 (Syzygy). All six primitives previously had only static-analysis anchors; CodeChemist adds compute-bound TTC and Syzygy adds runtime-mined type/safety properties.

## Watchlist for Wave-18+

- MemSearcher if a method-level companion lands at NeurIPS 2026 / ICML 2027 / ICLR 2027.
- Verified Tool Calls if venue confirmation lands.
- Software archaeology mechanism papers at NeurIPS 2026 / ICML 2027 / ICLR 2027 (active gap; recommend research-on-file in `.obsidian/MAgHARCM/research/diary/wave-18-watchlist.md` next sprint).
- MigGPT follow-up (NeurIPS 2025 spotlight) — patch migration, not type-aware translation; out of scope for this slot.

## Vault file stamp propagation (Sprint 2026-09-28)

- `.obsidian/MAgHARCM/research/METHODOLOGY.md` — §7 anchor list + §9 changelog entry + §11 SLM-Era General-Purpose Patterns section extended (2 new sub-sections for CodeChemist + Syzygy).
- `.obsidian/MAgHARCM/research/Architecture.md` — §8 extended with 2 new sub-sections (8.4 CodeChemist TTC Substrate, 8.5 Syzygy Dynamic-Specs Substrate); §7 Vault Sync Audit updated.
- `.obsidian/MAgHARCM/primitives/INDEX.md` — frontmatter extended; PRIM-21/23/27 rows + CodeChemist cross-links; PRIM-9/22/30 rows + Syzygy cross-links; Sprint 2026-09-28 audit block appended.
- `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` — §7 wave-16 + §8 wave-17 added; cross-reference matrix extended.
- `.obsidian/MAgHARCM/research/papers/P-123-codechemist-icml-2026.md` — new.
- `.obsidian/MAgHARCM/research/papers/P-124-syzygy-iclr2025-workshop.md` — new.
- `docs/.paper/refs.bib` — 2 new bib entries appended.
- `docs/.paper/sec_method.tex` — 2 new `\cite{}` clusters added.
- `.obsidian/MAgHARCM/diary/Sprint-2026-09-28-Handoff.md` — appended.
