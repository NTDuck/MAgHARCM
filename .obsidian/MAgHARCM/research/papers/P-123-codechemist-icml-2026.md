---
title: "P-123 CodeChemist: Test-Time Scaling for Low-Resource Code Generation"
backlink: "[[1.0.0 P-123]]"
tags: [paper, slm, test-time-scaling, low-resource-code, code-translation, training-free, "[[1.0.0 PRIM-21]]", "[[1.0.0 PRIM-23]]", "[[1.0.0 PRIM-27]]", "[[1.0.0 P-104]]", "[[1.0.0 P-91]]", "[[1.0.0 P-98]]", wave-17]
date: 2026-09-28
last_updated: 2026-09-28
venue: ICML 2026
arxiv: 2510.00501
---

# [[1.0.0 P-123]] CodeChemist

## TL;DR

CodeChemist is a **training-free test-time scaling framework** that improves code generation for low-resource programming languages by transferring functional knowledge from high-resource languages. Demonstrated on Qwen-1.5B (in the 4B-30B SLM band) with 60-70% relative gains on low-resource targets like Lua. Anchors `[[1.0.0 PRIM-21]]` Migration Strategy Selection, `[[1.0.0 PRIM-23]]` Chunked Translation, and `[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection.

## Mechanism (Q2)

Three mechanisms compose the framework:

1. **Multi-temperature hedged sampling** — generate diverse candidate solutions in the target low-resource language using a temperature ladder rather than a single draw.
2. **Uncertainty estimation** — model self-confidence gates the selection strategy: high confidence stays in-language; low confidence triggers cross-lingual fallback.
3. **Cross-lingual I/O test oracle** — when confidence is low, a high-resource reference program is **executed** (not just inspected) to synthesize a behavioural test that the candidate target-language solution must satisfy. This is the **functional knowledge transfer** that makes the training-free property work.

The training-free property is the key SLM-era hook — the technique is **compute-bound** (more samples → better result) rather than parameter-bound, which suits SLMs where post-training is constrained.

## Anchoring (Q3)

| Primitive | Pre-wave-17 behaviour | CodeChemist substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Blind try-and-fail over a fixed strategy registry | Confidence-gated strategy switching (in-language majority voting vs cross-lingual I/O test oracle) |
| `[[1.0.0 PRIM-23]]` Chunked Translation | Single-draw translation per chunk | Multi-temperature hedged sampling per chunk with cross-language functional verification |
| `[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection | Static-coverage plateau only | Functional-coverage plateau via I/O oracle (tests across source + target language) |

## Hop-1 Citations

- `[[1.0.0 P-91]]` Snell et al. 2024 *Scaling LLM Test-Time Compute Optimally* (ICLR 2025) — compute allocation framework.
- `[[1.0.0 P-104]]` Li et al. 2025 *S\*: Test Time Scaling for Code Generation* (EMNLP 2025 Findings) — parallel code-TTC ancestor; CodeChemist adds the low-resource + training-free dimensions.
- `[[1.0.0 P-98]]` Brown et al. 2024 *Large Language Monkeys* (best-of-N sampling).
- Wang et al. 2023 CodeT (execution-based code selection).

## Hop-2 Citations

- Wei et al. 2022 Chain-of-Thought.
- Yao et al. 2023 Tree of Thoughts.

## MAgHARCM integration

- **YAML config key**: `codechemist.enabled: true` (opt-in, default `false`).
- **Implementation file**: `internal/agents/chunked_translator.go::ApplyCodeChemist` (forthcoming — sprint 2026-09-29+).
- **Affected primitives**: `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-27]]`.

## Caveats

- **Cross-lingual I/O test oracle** requires executing a high-resource reference program (Lua ↔ Python, COBOL ↔ Java, etc.); MAgHARCM's existing execution sandbox (`internal/languages/oracle.go`) needs no changes but the orchestrator must preserve the source-language interpreter when the target is low-resource.
- **Multi-temperature hedged sampling** increases per-chunk inference cost linearly with the sample budget; the existing `compiletime.MaxGraphRunSteps = 50` ceiling may need adjustment when CodeChemist is enabled.

## Source

- arXiv:2510.00501 — CodeChemist (ICML 2026).
- OpenReview: `qpZhbxVXY4`.
