---
title: MAgHARCM Methodology Changelog & Sprint History
date: 2026-09-08
last_updated: 2026-09-08 (iter-1, wave-23)
aliases:
  - "Methodology-Changelog"
  - "Methodology Changelog"
  - "CHANGELOG"
  - "Changelog"
tags: [adhoc, changelog, history, methodology, "[[2.0.0 MAgHARCM]]"]
---

# [[2.0.0 MAgHARCM Methodology Changelog & Sprint History]]

> **Executive Overview**: Detailed historical log of architectural updates, literature wave triggers, triage decisions, and methodology enhancements across MAgHARCM research sprints. For core domain concepts, refer to the [[Glossary|Domain Acronyms & Terminology Glossary]]. Canonical methodology is documented in [[Methodology|Methodology Specification]].

---

## Sprint 2026-09-08 (Iteration 1: Wave-23 Close-Out)

- **Research Wave Trigger**: Wave-23 fired with 6 candidates evaluated against the §7 trigger gate.
- **Accepted Papers (4)**:
  - `[[1.0.0 P-153]]` **CoReX** (Sun et al., ICSE 2026): Context-aware refinement-based slicing for regression-failure localisation. Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Structure phase, extending P-133 ADI) + `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement. **Closes the strict program-comprehension-mechanism residual slot** open since Wave-17.
  - `[[1.0.0 P-154]]` **TransAgent** (Roh et al., FSE 2026): Multi-agent translation pipeline with fine-grained execution-aligned critic feedback. Anchors `[[1.0.0 PRIM-23]]` Chunked Translation + `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (complementing P-151 SmartC2Rust single-LLM loop).
  - `[[1.0.0 P-155]]` **POLA-Tester** (Sun et al., ICSE 2026): Agentic wait + syntactic dependency mining + iterative retrofit validation for LLM-augmented static analysis. Anchors `[[1.0.0 PRIM-12]]` Static Analysis Co-Evolution.
  - `[[1.0.0 P-156]]` **ACONITE** (Sun et al., ICSE 2026): Backward slicing + close-test retrieval + in-line execution annotations for coverage-plateau regression test generation. Anchors `[[1.0.0 PRIM-22]]` + `[[1.0.0 PRIM-29]]`.
- **Triage Decisions**:
  - **REJECT (Q1)**: R1 AutoCodeSherpa (Yunbo Lyu et al., ISSTA 2026) — off-list venue; ISSTA is not on the §7 trigger list.
  - **WATCHLIST**: W23-W1 SWE-TRACE (`arXiv:2604.14820`) carried 5th time awaiting NeurIPS 2026 author notifications (2026-09-24).
- **Substrate Wiring**: Opt-in gates enabled via `configs/agents.yml:comprehension.graph_self_evolving` and `configs/agents.yml:translation.feedback_driven`.
- **Vault & Paper Sync**: Persisted 4 paper notes, updated `Research-Database.json` (`reject_registry.wave-23`, `watchlist.wave-23`), updated LaTeX paper sections (`sec_method.tex`, `sec_eval.tex`).

---

## Sprint 2026-09-07 (Iteration 6: Wave-22 C-to-Rust & Hallucination)

- **Research Wave Trigger**: Wave-22 fired with 5 candidates evaluated.
- **Accepted Papers (2)**:
  - `[[1.0.0 P-151]]` **SmartC2Rust** (Sun et al., ICSE 2026): Feedback-driven iterative C-to-Rust translation with context-aware segmentation + three-signal feedback loop (target compiler errors, semantic diffs, unsafe-block counts). Anchors `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-29]]`, `[[1.0.0 PRIM-31]]`.
  - `[[1.0.0 P-152]]` **Hallu-Eval** (Liu et al., FSE 2026): Systematic hallucination-evaluation triplet (Hallu-Eval benchmark + Hallu-Det detection + Hallu-Shield mitigation). Anchors `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-25]]`.
- **Triage Decisions**:
  - **REJECT**: R1 Code vs. Serialized AST (LLM4Code workshop, off-list); R2 SmartComment (off-axis Solidity target); R3 Beyond Accuracy (DeepTest workshop, off-list).
  - **WATCHLIST**: W1 SWE-TRACE carried 4th time.

---

## Sprint 2026-09-07 (Iteration 4: Wave-21 KV-Cache Eviction)

- **Research Wave Trigger**: Wave-21 fired with 8 candidates evaluated.
- **Accepted Papers (4)**:
  - `[[1.0.0 P-147]]` **SpecKV** (Galim et al., ICLR 2026): Draft-model-driven KV eviction with adaptive gamma controller.
  - `[[1.0.0 P-148]]` **LookaheadKV** (Ahn et al., ICLR 2026): Parameter-efficient LoRA-modules on target model for speculative KV retention.
  - `[[1.0.0 P-149]]` **SSD/Saguaro** (Kumar et al., ICLR 2026): Asynchronous draft-verify speculative decoding pipeline.
  - `[[1.0.0 P-150]]` **TestPrune** (Chen et al., FSE 2026): Coverage-driven regression test minimization for Observation-phase context pruning (`[[1.0.0 PRIM-22]]`).
- **Triage Decisions**:
  - **REJECT**: R1 ABC (`arXiv:2602.22302`, unconfirmed venue); R2 NSE Workshop (off-list); R3 Speculative Actions (overlap with P-137).
  - **WATCHLIST**: W1 SWE-TRACE carried 3rd time.

---

## Sprint 2026-09-07 (Iteration 3: Wave-20 Pre-Analysis & Verification)

- **Research Wave Trigger**: Wave-20 fired with 10 candidates evaluated.
- **Accepted Papers (5)**:
  - `[[1.0.0 P-142]]` **NESA** (FSE 2026): Self-evolving graph pre-analysis for `[[1.0.0 PRIM-9]]` and `[[1.0.0 PRIM-22]]`.
  - `[[1.0.0 P-143]]` **HalluShield** (FSE 2026): Speculative decoding hallucination defense for `[[1.0.0 PRIM-7]]` and `[[1.0.0 PRIM-21]]`.
  - `[[1.0.0 P-144]]` **TraceCoder** (ICSE 2026): Prompt-trace training data construction for `[[1.0.0 PRIM-29]]` and `[[1.0.0 PRIM-31]]`.
  - `[[1.0.0 P-145]]` **TerraMod** (ICSE 2026 NIER): Lexical-based translation strategy selection for `[[1.0.0 PRIM-21]]`.
  - `[[1.0.0 P-146]]` **ContextPRM** (ICLR 2026): Process reward model across multi-turn migration workflows.
- **Administrative Action**: Applied BLK-06 dating convention reset (frontmatter date synchronized with actual system date `2026-09-07`).

---

## Sprint 2026-09-07 (Iteration 2: Wave-19 Speculative-Decoding & Test-Time Scaling)

- **Research Wave Trigger**: Wave-19 fired with 13 candidates evaluated.
- **Accepted Papers (8)**:
  - `[[1.0.0 P-134]]` **RelayCaching** (ICML 2026 Poster)
  - `[[1.0.0 P-135]]` **SPECS** (ICLR 2026)
  - `[[1.0.0 P-136]]` **CaTS** (ICLR 2026 Poster)
  - `[[1.0.0 P-137]]` **SuffixDecoding** (NeurIPS 2025 Spotlight)
  - `[[1.0.0 P-138]]` **RepairKV** (ICML 2026 AdaptFM Workshop)
  - `[[1.0.0 P-139]]` **TypePro** (FSE 2026)
  - `[[1.0.0 P-140]]` **Panta** (ICSE 2026)
  - `[[1.0.0 P-141]]` **KVFlow** (NeurIPS 2025 Poster)
- **Triage Decisions**: Added structured `reject_registry.wave-19` to `Research-Database.json` (BLK-08 resolution).

---

## Foundation & Parity Sprints (Waves 1–18)

- **Sprint 2026-09-25 to 09-28 (Waves 14–17)**: Achieved full 31/31 primitive parity; automated `scripts/lint_vault.sh` CI enforcement; anchored ReasoningBank (`[[1.0.0 P-122]]`), CodeChemist (`[[1.0.0 P-123]]`), Syzygy (`[[1.0.0 P-124]]`).
- **Sprint 2026-09-21 to 09-24 (Waves 11–13)**: Grounded agentic evaluation benchmarks: SWE-bench Verified (`[[1.0.0 P-109]]`), SWE-agent (`[[1.0.0 P-112]]`), OpenHands (`[[1.0.0 P-115]]`), Aider (`[[1.0.0 P-116]]`).
- **Sprint 2026-09-15 to 09-20 (Waves 9–10)**: Formalized prompt contracts, structured cloze slots, and ChunkKV compression.
- **Sprint 2026-09-09 to 09-14 (Waves 6–8)**: Enforced strict versioning conventions (`[[x.y.z ...]]`) and single source of truth in `primitives/Primitives-Index.md`.
- **Sprint 2026-09-04 to 09-08 (Inception)**: Scaffolding of 8-agent CloudWeGo Eino execution graph, baseline Go primitives, and Charm TUI.
