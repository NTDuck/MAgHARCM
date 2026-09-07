---
title: "P-129 L*λ*MDA: Large Language Model-Aided Partial Program Dependence Analysis"
backlink: "[[1.0.0 P-129]]"
aliases:
  - "1.0.0 P-129"
  - "P-129"
  - "P-129-LLMDA-Partial-Dependence-ICSE-2026"
  - "LLMDA-Partial-Dependence-ICSE-2026"
  - "rong2026lmda"
tags: [paper, partial-program, dependence-analysis, code-comprehension, "[[1.0.0 PRIM-9]]", "[[1.0.0 PRIM-22]]", wave-18]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICSE 2026 (Research Track)
---

# [[1.0.0 P-129]] L*λ*MDA

## TL;DR

L\*λ\*MDA introduces a **predictive dependence-analysis paradigm**: the LLM acts as a context augmenter that synthesizes an approximately-complete program P_AC = P + C (variable declarations, imports, type info, method signatures) from a partial program P, guided by an iterative compiler-feedback loop. A classical dependence-analysis tool (Joern/CPG) then runs on P_AC and the resulting PDG is pruned back to the original P. Achieves **5-265% F1 over learning-based baselines** and **17.2-71.3% recall improvement over classic DA tools**. Anchors `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph and `[[1.0.0 PRIM-22]]` Four Phases of Comprehension.

## Mechanism (Q2)

The framework proceeds in three iterative steps:

1. **LLM-guided context synthesis** — given partial program P, the LLM generates synthetic context C that completes declarations, imports, type annotations, and method signatures. The synthesis is **iterative**: a compiler-feedback loop rejects context proposals that fail to compile or produce unresolved symbols, prompting the LLM to revise.
2. **Classical DA on P_AC** — the augmented program P_AC is fed into a classical dependence-analysis tool (Joern / Code Property Graphs). Because P_AC is approximately-complete, the tool's parser, type resolver, and call-graph builder all succeed where they would have failed on P alone.
3. **Prune-back to P** — the resulting PDG includes edges that span into the synthetic context C. These edges are filtered, leaving only dependences between statements in the original P. The pruned PDG is the final output.

The iterative compiler-feedback loop is the key reliability mechanism: it bounds the LLM's hallucination surface by requiring every synthesized declaration to be at least compiler-valid. The prune-back step ensures the output is grounded in the original program, not the synthetic context.

## Anchoring (Q3)

| Primitive | Pre-wave-18 behaviour | L\*λ\*MDA substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Classical DA fails on partial programs (unresolved symbols, missing types); learning-based DA hallucinates edges | LLM-augmented context enables classical DA on partial programs; iterative compiler feedback bounds hallucination; prune-back grounds output in P |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Comprehension halts at missing-context phase for partial programs | Augmentation extends the comprehension horizon; DA produces structured dependence output usable by downstream phases |

## Hop-1 Citations

- Joern / Code Property Graphs (Yamaguchi et al. 2014) — classical DA substrate.
- Various learning-based partial-program DA baselines (2023-2024).
- LLM-for-code-comprehension lineage (Ma et al. 2024, Wang et al. 2024).

## Hop-2 Citations

- Compiler-feedback loops for code synthesis (Ellis et al. 2019, Chen et al. 2023).
- PDG pruning / slicing literature (Weiser 1984, Ottenstein et al. 1984).

## MAgHARCM integration

- **YAML config key**: `code_graph.llmda.enabled: true` (opt-in, default `false`); `code_graph.llmda.compiler_feedback_max_iters: <int>`.
- **Implementation file**: `internal/code_graph/llmda.go::SynthesizeContextAndAnalyze` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-9]]`, `[[1.0.0 PRIM-22]]`.

## Caveats

- **Compiler availability**: the feedback loop requires a compilable P_AC; languages without a robust compiler (e.g. dynamic languages with weak type info) cannot fully exploit the loop.
- **Synthetic-context drift**: the prune-back step assumes edges into C can be cleanly identified; spurious edges from synthetic to original code may leak through.
- **Compute cost**: each iteration incurs an LLM call plus a compiler invocation; on long partial programs the loop dominates wall-clock time.

## Source

- DOI: 10.1145/3744916.3773119.
- arXiv: forthcoming.
