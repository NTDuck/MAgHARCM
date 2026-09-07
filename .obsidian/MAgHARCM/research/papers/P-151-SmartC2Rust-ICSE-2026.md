---
title: "P-151 SmartC2Rust: Iterative, Feedback-Driven C-to-Rust Translation via Large Language Models for Safety and Equivalence"
backlink: "[[1.0.0 P-151]]"
aliases:
 - "1.0.0 P-151"
 - "P-151"
 - "P-151-SmartC2Rust-ICSE-2026"
 - "SmartC2Rust-ICSE-2026"
 - "smartc2rust2026icse"
tags: [paper, c-to-rust, iterative-translation, segmentation, compiler-feedback, slm, "[[1.0.0 PRIM-23]]", "[[1.0.0 PRIM-29]]", "[[1.0.0 PRIM-31]]", wave-22]
date: 2026-09-07
last_updated: 2026-09-07 (iter-6, wave-22)
venue: ICSE 2026 (Research Track)
---

# [[1.0.0 P-151]] SmartC2Rust

## TL;DR

SmartC2Rust is an **iterative, feedback-driven C-to-Rust translation mechanism** using **context-aware code segmentation** + a **three-signal feedback loop** (Rust compiler errors, semantic-equivalence diffs, `unsafe`-block counts) to drive LLM refinement until `unsafe` is minimised and semantic equivalence holds. The C-to-Rust analogue of Syzygy P-124 (Go-to-Rust); closes the C-to-Rust slot in the iterative-translation anchor pattern. Anchors `[[1.0.0 PRIM-23]]` Chunked Translation (context-aware C segmentation = chunking substrate), `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter (three-signal feedback loop = recruiter-driven repair cycle for C-to-Rust), and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (compiler-feedback + unsafe-residual feedback = retrieval-refinement substrate).

## Mechanism (Q2)

1. **Context-aware code segmentation** — the C source is divided into smaller, context-bounded units that fit within the LLM's context window while preserving semantic locality (function boundaries + transitive call-context).
2. **Initial LLM translation** — each C unit is translated into a Rust unit by the LLM in a single forward pass.
3. **Three-signal feedback loop** — the LLM receives three orthogonal feedback signals per iteration:
   - **Compiler errors** (lexical/syntactic feedback): the Rust compiler's `cargo check` output is fed back; the LLM fixes type errors, lifetime errors, and borrow-check violations.
   - **Semantic-equivalence diffs** (semantic feedback): a differential semantic check compares the translated Rust against the original C; divergences are fed back as structured feedback (not just "wrong").
   - **`unsafe`-block counts** (security feedback): the number of `unsafe` blocks in the translated Rust is minimised via feedback; residual unsafe is reported per-iteration.
4. **Iterative refinement** — the loop continues until the LLM produces a Rust unit that compiles, is semantically equivalent to the original C, and has minimal residual `unsafe`.
5. **End-to-end safety + equivalence** — the final output is a translated Rust codebase with reduced `unsafe`-block residuals and higher semantic equivalence than prior C-to-Rust translation baselines.

Together these give MAgHARCM an **iterative-translation feedback substrate** for the C-to-Rust target — a multi-signal repair loop that complements P-124 Syzygy's Go-to-Rust dual-oracle anchor and P-145 TerraMod's lexical-based translation strategy selection.

## Anchoring (Q3)

| Primitive | Pre-wave-22 behaviour | SmartC2Rust substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-23]]` Chunked Translation | Chunking anchor on P-124 Syzygy (Go-to-Rust); C-to-Rust slot open | SmartC2Rust context-aware C segmentation = chunking substrate for the C-to-Rust target |
| `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter | Recruiter-driven repair anchor on P-124 Syzygy (Go-to-Rust); C-to-Rust slot open | SmartC2Rust three-signal feedback loop (compiler + semantic-diff + unsafe-residual) = recruiter-driven repair cycle for C-to-Rust |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Retrieval-refinement substrate dominated by KV-cache × speculative-decoding anchors (P-128 KVzip, P-141 KVFlow, P-147..P-149); few repair-loop anchors | SmartC2Rust compiler-feedback + unsafe-residual feedback = retrieval-refinement substrate for C-to-Rust translation loops |

## Hop-1 Citations

- P-124 Syzygy (ICLR 2024) — Go-to-Rust dual code-test oracle; SmartC2Rust is the C-to-Rust analogue (C input, Rust output, multi-signal feedback loop).
- P-145 TerraMod (ICSE 2026 NIER) — lexical-based translation strategy selection; SmartC2Rust's context-aware segmentation complements TerraMod's strategy-selection substrate.
- P-146 ContextPRM (ICLR 2026) — workflow-aware cross-document process reward modelling; SmartC2Rust's semantic-equivalence diff feedback is a single-document analogue of ContextPRM's process-reward signal.
- EAGLE-3 P-78 (2025) — speculative-decoding iterative repair loop; SmartC2Rust is the source-to-source translation analogue.
- LangSec (Erbsen et al., POPL 2025) — foundational source-to-source translation safety substrate; SmartC2Rust inherits the LangSec safety framing for the C-to-Rust reduction.

## Hop-2 Citations

- Differential testing (McKeeman 1998) — foundational differential semantic-check substrate; SmartC2Rust's semantic-equivalence diff is an LLM-era differential-testing analogue.
- Iterative compilation (Burke et al., 1996) — foundational compile-feedback loop; SmartC2Rust is the LLM-era source-to-source analogue.
- Rust unsafe-block minimisation (Jung 2017, RFC 2585) — foundational unsafe-residual substrate; SmartC2Rust uses `unsafe`-count as a feedback signal.
- Code segmentation for translation (Shiraishi 2025, context-aware C segmentation researchgate publication 384075558) — direct precedent for context-aware C segmentation before LLM translation.

## MAgHARCM integration

- **YAML config key**: `agents.translation.feedback_driven: true`; `agents.translation.compiler_feedback: true`; `agents.translation.semantic_diff_feedback: true`; `agents.translation.unsafe_residual_feedback: true`; `agents.translation.max_iterations: <int>`.
- **Implementation file**: `internal/translation/feedback_loop.go::NewSmartC2RustLoop` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-29]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **C-to-Rust only** — SmartC2Rust does not cover Go-to-Rust (already anchored by P-124 Syzygy) or Java-to-Rust (no anchor yet; Commons-Validator BLK-02 plateau).
- **Semantic-diff cost** — the semantic-equivalence differential check requires a runnable test harness; C projects without test coverage cannot use SmartC2Rust end-to-end (the same caveat as P-124 Syzygy).
- **LLM round-trip cost** — the three-signal feedback loop adds 2-3 LLM round-trips per iteration; must be amortised over the compile-fix cost.
- **Feedback-signal weighting** — the three signals (compiler / semantic-diff / unsafe-residual) are not weighted; future work should weight them by signal-strength per iteration.

## Source

- arXiv: 2409.10506.
- Venue: ICSE 2026 (Research Track), verified via u-tokyo.ac.jp publication news 2025-10-16; conference DOI: 10.1145/3744916.3773259; arXiv DOI: 10.48550/arXiv.2409.10506.

## BibTeX

```
@inproceedings{smartc2rust2026icse,
  title     = {SmartC2Rust: Iterative, Feedback-Driven C-to-Rust Translation via Large Language Models for Safety and Equivalence},
  author    = {Shiraishi, Momoko and Cao, Yinzhi and Shinagawa, Takahiro},
  booktitle = {Proceedings of the 48th International Conference on Software Engineering (ICSE)},
  year      = {2026},
  doi       = {10.1145/3744916.3773259},
  eprint    = {arXiv:2409.10506}
}
```
