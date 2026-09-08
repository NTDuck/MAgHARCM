---
title: MAgHARCM Research Waves Synthesis & Literature Index
date: 2026-09-08
last_updated: 2026-09-08 (iter-1, wave-23)
aliases:
  - "Research-Waves-Index"
  - "Research Waves Index"
  - "Waves"
tags: [adhoc, research-waves, literature, synthesis, papers, "[[2.0.0 MAgHARCM]]"]
---

# [[2.0.0 MAgHARCM Research Waves Index (Waves 1–23)]]

> **Executive Overview**: Comprehensive chronological catalog and thematic synthesis of all 23 research waves, mapping 156 literature papers to the MAgHARCM multi-agent modernization substrate. Acronyms and technical terms are defined in the [[Glossary|Domain Acronyms & Terminology Glossary]].

---

## 1. Thematic Groupings of Research Waves

| Wave Range | Core Theme | Target Primitives | Key Literature Anchors |
| :--- | :--- | :--- | :--- |
| **Wave 23** | Context-Aware Refinement Slicing + Multi-Agent Translation + Static-Analysis Co-Evolution + Test Generation | `PRIM-12, 22, 23, 29, 31` | CoReX (`[[1.0.0 P-153]]`), TransAgent (`[[1.0.0 P-154]]`), POLA-Tester (`[[1.0.0 P-155]]`), ACONITE (`[[1.0.0 P-156]]`) |
| **Wave 22** | Feedback-Driven C-to-Rust Translation + Systematic Hallucination Evaluation | `PRIM-22, 23, 25, 29, 31` | SmartC2Rust (`[[1.0.0 P-151]]`), Hallu-Eval (`[[1.0.0 P-152]]`) |
| **Wave 21** | KV-Cache $\times$ Speculative-Decoding Substrate + Static-Analysis Context Pruning | `PRIM-21, 22, 31` | SpecKV (`[[1.0.0 P-147]]`), LookaheadKV (`[[1.0.0 P-148]]`), SSD/Saguaro (`[[1.0.0 P-149]]`), TestPrune (`[[1.0.0 P-150]]`) |
| **Wave 20** | Program-Comprehension Mechanism + SLM Verification Substrate | `PRIM-7, 9, 21, 22, 25, 31` | NESA (`[[1.0.0 P-142]]`), HalluShield (`[[1.0.0 P-143]]`), TraceCoder (`[[1.0.0 P-144]]`), TerraMod (`[[1.0.0 P-145]]`), ContextPRM (`[[1.0.0 P-146]]`) |
| **Wave 19** | Speculative-Decoding $\times$ KV-Cache Hybrids + SLM Test-Time Scaling | `PRIM-7, 9, 21, 22, 31` | RelayCaching (`[[1.0.0 P-134]]`), SPECS (`[[1.0.0 P-135]]`), CaTS (`[[1.0.0 P-136]]`), SuffixDecoding (`[[1.0.0 P-137]]`), RepairKV (`[[1.0.0 P-138]]`), TypePro (`[[1.0.0 P-139]]`), Panta (`[[1.0.0 P-140]]`), KVFlow (`[[1.0.0 P-141]]`) |
| **Wave 18** | SLM-Scale Verdict + Software Archaeology | `PRIM-7, 9, 21, 22, 31` | T1 (`[[1.0.0 P-125]]`), ARC-Decode (`[[1.0.0 P-126]]`), SLM-as-a-Judge (`[[1.0.0 P-127]]`), KVzip (`[[1.0.0 P-128]]`), LλMDA (`[[1.0.0 P-129]]`), SSAR (`[[1.0.0 P-130]]`), SemArc (`[[1.0.0 P-131]]`), SemRef (`[[1.0.0 P-132]]`), ADI (`[[1.0.0 P-133]]`) |
| **Waves 15–17** | Test-Time Scaling & Multi-Agent Reasoning | `PRIM-9, 21, 22, 27, 31` | ReasoningBank (`[[1.0.0 P-122]]`), CodeChemist (`[[1.0.0 P-123]]`), Syzygy C-to-Rust Dual Oracle (`[[1.0.0 P-124]]`) |
| **Waves 11–14** | Agentic Benchmarks & Tool Use | `PRIM-5, 9, 22, 25, 29` | SWE-bench Verified (`[[1.0.0 P-109]]`), GraphCoder (`[[1.0.0 P-110]]`), SWE-agent (`[[1.0.0 P-112]]`), AutoCodeRover (`[[1.0.0 P-113]]`), OpenHands (`[[1.0.0 P-115]]`), Aider (`[[1.0.0 P-116]]`) |
| **Waves 7–10** | SLM-Era Modernization (4B–30B) | `PRIM-3, 6, 7, 21, 23` | Qwen2.5-Coder (`[[1.0.0 P-21]]`, `[[1.0.0 P-58]]`), Phi-3 (`[[1.0.0 P-27]]`), EAGLE-3 (`[[1.0.0 P-78]]`), Self-Consistency (`[[1.0.0 P-83]]`) |
| **Wave 6** | Software Archaeology Foundations | `PRIM-14, 18, 19, 20` | Parnas 1972 (`[[1.0.0 P-31]]`), Lehman 1980 (`[[1.0.0 P-32]]`), Chikofsky & Cross (`[[1.0.0 P-33]]`), Baldwin & Clark (`[[1.0.0 P-34]]`) |
| **Waves 1–5** | Core Systems & Foundations | `PRIM-1` to `PRIM-30` | ReCodeAgent (`[[1.0.0 P-01]]`), AlphaTrans (`[[1.0.0 P-02]]`), CodePlan (`[[1.0.0 P-03]]`), MetaGPT (`[[1.0.0 P-11]]`), ChatDev (`[[1.0.0 P-12]]`), ABCoder (`[[1.0.0 P-14]]`) |

---

## 2. Detailed Chronological Wave Catalog

### Recent Research Waves (Waves 18–23)

#### Wave 23: Program Comprehension Residual Slot Closure & Test Synthesis
- **Scope**: Closes the strict program-comprehension-mechanism residual slot carried forward from Wave-17.
- **Accepted Papers (4)**:
  - `[[1.0.0 P-153]]` **CoReX** (Sun et al., ICSE 2026): Context-aware refinement-based slicing for regression failure localization (`PRIM-22`, `PRIM-31`).
  - `[[1.0.0 P-154]]` **TransAgent** (Roh et al., FSE 2026): Multi-agent translation pipeline with fine-grained critic feedback (`PRIM-23`, `PRIM-31`).
  - `[[1.0.0 P-155]]` **POLA-Tester** (Sun et al., ICSE 2026): Syntactic dependency mining and iterative retrofit validation (`PRIM-12`).
  - `[[1.0.0 P-156]]` **ACONITE** (Sun et al., ICSE 2026): Execution-annotated backward slicing for LLM regression test generation (`PRIM-22`, `PRIM-29`).
- **Triage**: 1 REJECT (R1 AutoCodeSherpa, ISSTA off-list); 1 WATCHLIST (SWE-TRACE 5th carry).

#### Wave 22: Feedback-Driven Translation & Hallucination Defense
- **Scope**: Phase 3 C/C++ to Safe Rust translation substrate and systematic hallucination evaluation.
- **Accepted Papers (2)**:
  - `[[1.0.0 P-151]]` **SmartC2Rust** (Sun et al., ICSE 2026): Three-signal feedback loop (compiler errors, semantic diffs, unsafe-block counts) for `PRIM-23, 29, 31`.
  - `[[1.0.0 P-152]]` **Hallu-Eval** (Liu et al., FSE 2026): Hallu-Eval benchmark + Hallu-Det detection + Hallu-Shield mitigation for `PRIM-22, 25`.
- **Triage**: 3 REJECT (R1 Code vs. Serialized AST, R2 SmartComment, R3 Beyond Accuracy); 1 WATCHLIST (SWE-TRACE 4th carry).

#### Wave 21: KV-Cache Management & Context Pruning
- **Scope**: Parameter-efficient KV-cache retention and coverage-driven test minimization.
- **Accepted Papers (4)**:
  - `[[1.0.0 P-147]]` **SpecKV** (Galim et al., ICLR 2026): Draft-model-driven KV eviction with adaptive gamma controller.
  - `[[1.0.0 P-148]]` **LookaheadKV** (Ahn et al., ICLR 2026): Parameter-efficient LoRA modules for target model lookahead.
  - `[[1.0.0 P-149]]` **SSD/Saguaro** (Kumar et al., ICLR 2026): Asynchronous speculative decoding pipeline.
  - `[[1.0.0 P-150]]` **TestPrune** (Chen et al., FSE 2026): Coverage-driven test minimization for observation-phase context pruning (`PRIM-22`).
- **Triage**: 3 REJECT (R1 ABC arXiv-only, R2 NSE off-list workshop, R3 Speculative Actions overlap); 1 WATCHLIST (SWE-TRACE 3rd carry).

#### Wave 20: Self-Evolving Graph & Verification Substrates
- **Scope**: Graph pre-analysis and SLM-grounded hallucination defense.
- **Accepted Papers (5)**:
  - `[[1.0.0 P-142]]` **NESA** (FSE 2026): Self-evolving graph pre-analysis for `PRIM-9, 22`.
  - `[[1.0.0 P-143]]` **HalluShield** (FSE 2026): Speculative decoding hallucination defense for `PRIM-7, 21`.
  - `[[1.0.0 P-144]]` **TraceCoder** (ICSE 2026): Prompt-trace training-data construction for `PRIM-29, 31`.
  - `[[1.0.0 P-145]]` **TerraMod** (ICSE 2026 NIER): Lexical translation strategy selection for `PRIM-21`.
  - `[[1.0.0 P-146]]` **ContextPRM** (ICLR 2026): Workflow-aware process reward modeling for `PRIM-7`.

#### Wave 19: KV-Cache Compression & Scaling
- **Scope**: Speculative-decoding $\times$ KV-cache compression hybrids and SLM test-time scaling.
- **Accepted Papers (8)**: `[[1.0.0 P-134]]` RelayCaching, `[[1.0.0 P-135]]` SPECS, `[[1.0.0 P-136]]` CaTS, `[[1.0.0 P-137]]` SuffixDecoding, `[[1.0.0 P-138]]` RepairKV, `[[1.0.0 P-139]]` TypePro, `[[1.0.0 P-140]]` Panta, `[[1.0.0 P-141]]` KVFlow.

#### Wave 18: Verdict Validation & Archaeology
- **Scope**: SLM-scale verdict validation and LLM-augmented software archaeology.
- **Accepted Papers (9)**: `[[1.0.0 P-125]]` T1, `[[1.0.0 P-126]]` ARC-Decode, `[[1.0.0 P-127]]` SLM-as-a-Judge, `[[1.0.0 P-128]]` KVzip, `[[1.0.0 P-129]]` LλMDA, `[[1.0.0 P-130]]` SSAR, `[[1.0.0 P-131]]` SemArc, `[[1.0.0 P-132]]` SemRef, `[[1.0.0 P-133]]` ADI.

---

### Foundation & Baseline Waves (Waves 1–17)

- **Waves 15–17 (P-122 to P-124)**: Test-time scaling and persistent memory. ReasoningBank (`[[1.0.0 P-122]]`) introduced structured memory triples; CodeChemist (`[[1.0.0 P-123]]`) and Syzygy (`[[1.0.0 P-124]]`) established dual runtime oracles for C-to-Rust.
- **Waves 11–14 (P-108 to P-121)**: Benchmarking and tool integration. Established SWE-bench Verified (`[[1.0.0 P-109]]`), CPG-based GraphCoder (`[[1.0.0 P-110]]`), and Agent-Computer Interfaces (`[[1.0.0 P-112]]`).
- **Waves 7–10 (P-50 to P-107)**: Local SLM feasibility and prompt contracts. Structured cloze output formats, lost-in-the-middle context placement, and ChunkKV compression.
- **Wave 6 (P-31 to P-49)**: Software archaeology theoretical foundation: Parnas information hiding, Lehman's laws of software evolution, and Baldwin & Clark modularity theory.
- **Waves 1–5 (P-01 to P-30)**: Grounded the core modernization loop: bottom-up translation, skeleton emission, dynamic invariants, and multi-tier test cascades.
