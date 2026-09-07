---
title: MAgHARCM Research Waves Synthesis & Literature Index
date: 2026-09-07
last_updated: 2026-09-07 (iter-6, wave-22)
aliases:
  - "Research-Waves-Index"
  - "Research Waves Index"
  - "Waves"
tags: [adhoc, research-waves, literature, synthesis, papers, "[[2.0.0 MAgHARCM]]"]
---
# [[2.0.0 MAgHARCM Research Waves Index (Waves 1–22)]]

> **Executive Overview**: Complete chronological catalog and thematic synthesis of all 22 research waves, connecting 152 literature papers to the MAgHARCM multi-agent modernization substrate.

> **Wave 22 (this sprint, 2026-09-07 iter-6)** adds 2 ACCEPT (P-151 SmartC2Rust, P-152 Hallu-Eval) for 152 total paper anchors. Vault paper count: 150 -> 152.

---

## Thematic Groupings of Research Waves

| Wave Range | Theme | Anchored Primitives | Key Literature Anchors |
| :--- | :--- | :--- | :--- |
| **Wave 20** | Program-Comprehension Mechanism + SLM Verification Substrate | `PRIM-7, 9, 21, 22, 25, 31` | NESA (`[[1.0.0 P-142]]`), HalluShield (`[[1.0.0 P-143]]`), TraceCoder (`[[1.0.0 P-144]]`), TerraMod (`[[1.0.0 P-145]]`), ContextPRM (`[[1.0.0 P-146]]`) |
| **Wave 21** | KV-Cache × Speculative-Decoding Substrate + Static-Analysis Context Pruning | `PRIM-21, 22, 31` | SpecKV (`[[1.0.0 P-147]]`), LookaheadKV (`[[1.0.0 P-148]]`), Speculative Speculative Decoding (`[[1.0.0 P-149]]`), TestPrune (`[[1.0.0 P-150]]`) |

| **Wave 22** | Feedback-Driven C-to-Rust Translation + Systematic Hallucination Evaluation | `PRIM-22, 23, 25, 29, 31` | SmartC2Rust (`[[1.0.0 P-151]]`), Hallu-Eval (`[[1.0.0 P-152]]`) |

| **Waves 1–5** | Core Systems & Foundations | `PRIM-1` through `PRIM-30` | ReCodeAgent (`[[1.0.0 P-01]]`), AlphaTrans (`[[1.0.0 P-02]]`), CodePlan (`[[1.0.0 P-03]]`), MetaGPT (`[[1.0.0 P-11]]`), ChatDev (`[[1.0.0 P-12]]`), ABCoder (`[[1.0.0 P-14]]`) |
| **Wave 6** | Software Archaeology Foundations | `PRIM-14, 18, 19, 20` | Parnas 1972 (`[[1.0.0 P-31]]`), Lehman 1980 (`[[1.0.0 P-32]]`), Chikofsky & Cross (`[[1.0.0 P-33]]`), Baldwin & Clark (`[[1.0.0 P-34]]`) |
| **Waves 7–10** | SLM-Era Modernization (4B–30B) | `PRIM-3, 6, 7, 21, 23` | Qwen2.5-Coder (`[[1.0.0 P-21]]`, `[[1.0.0 P-58]]`), Phi-3 (`[[1.0.0 P-27]]`, `[[1.0.0 P-54]]`), EAGLE-3 (`[[1.0.0 P-78]]`), Self-Consistency (`[[1.0.0 P-83]]`) |
| **Waves 11–14** | Agentic Benchmarks & Tool Use | `PRIM-5, 9, 22, 25, 29` | SWE-bench Verified (`[[1.0.0 P-109]]`), GraphCoder (`[[1.0.0 P-110]]`), SWE-agent (`[[1.0.0 P-112]]`), AutoCodeRover (`[[1.0.0 P-113]]`), OpenHands (`[[1.0.0 P-115]]`), Aider (`[[1.0.0 P-116]]`) |
| **Waves 15–17** | Test-Time Scaling & Multi-Agent Reasoning | `PRIM-9, 21, 22, 27, 31` | ReasoningBank (`[[1.0.0 P-122]]`), CodeChemist (`[[1.0.0 P-123]]`), Syzygy C-to-Rust Dual Oracle (`[[1.0.0 P-124]]`) |
| **Wave 18** | SLM-Scale Verdict + Software Archaeology | `PRIM-7, 9, 21, 22, 31` | T1 Tool-Integrated Verification (`[[1.0.0 P-125]]`), ARC-Decode (`[[1.0.0 P-126]]`), SLM-as-a-Judge (`[[1.0.0 P-127]]`), KVzip (`[[1.0.0 P-128]]`), LλMDA (`[[1.0.0 P-129]]`), SSAR (`[[1.0.0 P-130]]`), SemArc (`[[1.0.0 P-131]]`), SemRef (`[[1.0.0 P-132]]`), ADI (`[[1.0.0 P-133]]`) |
| **Wave 19** | Speculative-Decoding × KV-Cache Hybrids + SLM Test-Time Scaling | `PRIM-7, 9, 21, 22, 31` | RelayCaching (`[[1.0.0 P-134]]`), SPECS (`[[1.0.0 P-135]]`), CaTS (`[[1.0.0 P-136]]`), SuffixDecoding (`[[1.0.0 P-137]]`), RepairKV (`[[1.0.0 P-138]]`), TypePro (`[[1.0.0 P-139]]`), Panta (`[[1.0.0 P-140]]`), KVFlow (`[[1.0.0 P-141]]`) |

---

## Detailed Chronological Wave Catalog

- **Wave 1–5 (P-01 to P-30)**: Grounded the core modernization loop: bottom-up translation, skeleton emission, dynamic invariants, and multi-tier test cascades.
- **Wave 6 (P-31 to P-49)**: Anchored software archaeology theory: information hiding, software evolution laws, reverse engineering roadmaps, and design rule hierarchy partitioning.
- **Wave 7–8 (P-50 to P-95)**: Established local SLM feasibility: parameter-efficient training, lost-in-the-middle context placement, speculative decoding, and chain-of-thought code generation.
- **Wave 9–10 (P-96 to P-107)**: Formalized prompt contracts: structured output cloze slots, KV cache compression (ChunkKV), and Decomposed Prompting for SLMs.
- **Wave 11–13 (P-108 to P-118)**: Grounded evaluation in canonical benchmarks: SWE-bench original/lite/verified, Graph-RAG retrieval over CPG, and agent-computer interfaces.
- **Wave 14 (P-119 to P-121)**: Evaluated multi-language benchmarks and tool-calling leaderboards (SWE-Rebench, SWE-smith, BFCL).
- **Wave 15 (Deferred)**: Evaluated multimodal SWE-bench extensions; deferred per strict trigger gate criteria.
- **Wave 16 (P-122)**: Distilled self-reflection strategies into ReasoningBank for zero-shot strategy retrieval.
- **Wave 17 (P-123, P-124)**: Test-time scaling for code synthesis (CodeChemist) and dual code-test C-to-Rust runtime oracle (Syzygy).
- **Wave 18 (P-125..P-133)**: SLM-scale verdict validation + LLM-augmented software archaeology. P-125 T1 tool-integrated verification, P-126 ARC-Decode risk-bounded speculative decoding, P-127 SLM-as-a-Judge. P-128 KVzip query-agnostic KV cache compression. P-129..P-132 LLM-aided partial program dependence + architecture recovery (LλMDA, SSAR, SemArc, SemRef). P-133 ADI Frame Lifetime Trace for autonomous debugging (SIGSOFT Distinguished Paper Award at FSE 2026). 4 candidates REJECTED: ReflexiCoder (ACL 2026 Findings, off-list venue), Self-Distillation for Code Generation (Apple arXiv:2604.01193, no venue), SPECS (ICLR 2026 submission unconfirmed at triage), Software-Archaeology strict-mechanism gap (gap partially closed by Wave-18 itself).
- **Wave 19 (P-134..P-141)**: Speculative-decoding × KV-cache compression hybrids + SLM-scale test-time scaling + LLM-augmented static analysis beyond dependence graphs. P-134 RelayCaching cross-agent KV reuse on the 8-agent graph. P-135 SPECS + P-136 CaTS = dual SLM-scale test-time scaling frontier-PRM-replacement anchors. P-137 SuffixDecoding model-free suffix-tree draft. P-138 RepairKV post-compression KV repair (borderline workshop-track ACCEPT per §7 method-level threshold). P-139 TypePro inter-procedural SDG slicing. P-140 Panta iterative hybrid static+dynamic LLM test generation. P-141 KVFlow workflow-aware KV cache eviction for multi-agent retention. 3 candidates REJECTED (Q1): R1 TTA* (workshop redundancy vs P-135/P-136), R2 HELIOS (NDSS LAST-X off-list + off-axis binary decompilation), R3 LongSpec (ACL off-list venue). 2 candidates UNVERIFIED (Watchlist): W1 SliceMate (ISSTA 2026 author-claim unconfirmed), W2 SWE-TRACE (arXiv:2604.14820).
- **Wave 20 (P-142..P-146)**: Program-comprehension-mechanism substrate + SLM-grounded hallucination defence + prompt-trace training data + workflow-aware PRM. P-142 NESA self-evolving graph pre-analysis (partially closes the strict program-comprehension-mechanism slot open since Wave-17). P-143 HalluShield SLM-grounded speculative-decoding hallucination defence. P-144 TraceCoder prompt-trace training-data construction. P-145 TerraMod LLM-driven lexical-based translation strategy selection. P-146 ContextPRM workflow-aware cross-document process reward modelling. 4 candidates REJECTED (Q1+Q3): R1 Nexus ICSE 2026 (ablations only, no mechanism), R2 SWE-Lego ICSE 2026 NIER (engineering pattern, no mechanism), R3 CoPS ICML 2026 (speculative venue), R4 SHIELD-ASR ACL 2026 Findings (off-list venue). 1 candidate UNVERIFIED (Watchlist): U1 NSE ICML 2026 placeholder venue, no DOI / OpenReview / arXiv. 2 watchlist items re-verified: SliceMate (REJECTED — ISSTA 2026 program slot absent on conf.researchr.org); SWE-TRACE (CONFIRMED arXiv-only preprint, REJECTED).
- **Wave 21 (P-147..P-150)**: KV-cache × speculative-decoding substrate reinforcement + static-analysis context pruning for SLM agent comprehension. P-147 SpecKV draft-model-driven KV eviction (Galim et al., ICLR 2026, OpenReview 0vbYakkECY). P-148 LookaheadKV parameter-efficient LoRA-modules on target model (Ahn et al., ICLR 2026, OpenReview RVLMGPXt2i). P-149 SSD/Saguaro asynchronous speculative-decoding pipeline (Kumar, Dao, May, ICLR 2026, OpenReview aL1Wnml9Ef). P-150 TestPrune coverage-driven regression-test minimization (Chen, Ahmed et al., IBM Research, FSE 2026, DOI 10.1145/3808148). 3 candidates REJECTED (Q1+Q3): R1 ABC (Bhardwaj, arXiv:2602.22302 — no peer-reviewed venue; the pre-existing 'Wang et al. ICSE 2026' attribution is a hallucination), R2 NSE Workshop (off-list workshop venue co-located with ICSE 2026), R3 Speculative Actions (ICLR 2026 poster — mechanism overlap with P-137 SuffixDecoding). 1 UNVERIFIED (Watchlist): W1 SWE-TRACE carried Wave-19 W2 → Wave-20 W1 → Wave-21 W1 (arXiv:2604.14820; NeurIPS 2026 notifications scheduled 2026-09-24). Vault paper count: 146 → 150.

