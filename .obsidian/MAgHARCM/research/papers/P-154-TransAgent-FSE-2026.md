---
title: "P-154 TransAgent: Enhancing LLM-Based Code Translation via Fine-Grained Execution Alignment"
backlink: "[[1.0.0 P-154]]"
aliases:
 - "1.0.0 P-154"
 - "P-154"
 - "P-154-TransAgent-FSE-2026"
 - "TransAgent-FSE-2026"
 - "transagent2026fse"
tags: [paper, code-translation, multi-agent, execution-alignment, fine-grained-feedback, "[[1.0.0 PRIM-23]]", "[[1.0.0 PRIM-29]]", "[[1.0.0 PRIM-31]]", "[[1.0.0 P-124]]", "[[1.0.0 P-145]]", "[[1.0.0 P-151]]", "[[1.0.0 P-131]]", wave-23]
date: 2026-09-07
last_updated: 2026-09-07 (iter-7, wave-23)
venue: FSE 2026 (Research Track, PACMSE Vol. 3 Issue FSE)
---

# [[1.0.0 P-154]] TransAgent

## TL;DR

TransAgent is a **multi-agent LLM code-translation pipeline** that localises error-prone code blocks by comparing source-code and translated-code **execution behaviour** (fine-grained execution alignment) rather than relying on test-output-only diffs. Multi-agent collaboration (translator + critic + repairer) drives iterative repair. Outperforms UniTrans by up to 33.3% in translation accuracy and improves program repair by 56.7% on average over Agentless. Anchors `[[1.0.0 PRIM-23]]` Chunked Translation (extends P-124 Syzygy / P-145 TerraMod / P-151 SmartC2Rust single-LLM iterative feedback to multi-agent collaboration + execution-aligned critic) and `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter (iterative repair loop as iteration-recruiter anchor) and `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (execution-aligned critic feedback = retrieval refinement substrate).

## Mechanism (Q2)

1. **Multi-agent role decomposition** — three agents collaborate: (a) translator (source → target), (b) critic (compares source/target execution behaviour and localises error-prone blocks), (c) repairer (rewrites the localised blocks). The critic operates on **fine-grained execution alignment** — not test-output diffs but per-block execution-trace comparison.
2. **Fine-grained execution alignment** — the critic compares source-code and translated-code execution behaviour at block granularity rather than at test-output granularity. Localised error-prone blocks are flagged for repair.
3. **Iterative repair loop** — the repairer rewrites the localised blocks; the critic re-aligns; the loop continues until alignment error falls below threshold.
4. **Multi-agent feedback signals** — three orthogonal feedback signals: (a) lexical/syntactic (compiler errors), (b) execution-alignment (critic's localised blocks), (c) semantic-equivalence (test outputs).

Together these give MAgHARCM a **multi-agent execution-aligned translation pipeline** — the first multi-agent + execution-alignment combination in the SLM-era translation substrate.

## Anchoring (Q3)

| Primitive | Pre-wave-23 behaviour | TransAgent substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-23]]` Chunked Translation | Translation pipeline anchored on P-124 Syzygy (Go-to-Rust single-LLM iterative compiler feedback), P-145 TerraMod (source-to-source modernisation via LLM agents), P-151 SmartC2Rust (C-to-Rust single-LLM iterative feedback); no multi-agent translation pipeline | TransAgent as the **multi-agent translation pipeline** with translator + critic + repairer role decomposition; extends the chunked-translation lineage to multi-agent + execution-aligned critic |
| `[[1.0.0 PRIM-29]]` Dynamic Iteration Recruiter | Iteration recruiter anchored on P-124 Syzygy, P-122 ReasoningBank, P-145 TerraMod, P-151 SmartC2Rust; no concrete multi-agent iterative-repair-loop anchor | TransAgent iterative-repair-loop = concrete iteration-recruiter anchor with execution-aligned critic feedback |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Retrieval refinement anchored on P-134 RelayCaching (cross-agent KV reuse), P-137 SuffixDecoding (model-free suffix-tree draft), P-141 KVFlow (workflow-aware KV cache eviction), P-151 SmartC2Rust (compiler + semantic + unsafe feedback); no execution-aligned critic feedback | TransAgent execution-aligned critic feedback = retrieval refinement by execution-behaviour comparison |

## Hop-1 Citations

- P-151 SmartC2Rust (ICSE 2026, Wave-22) — C-to-Rust iterative feedback (compiler + semantic + unsafe). TransAgent extends the iterative-feedback substrate to multi-agent + execution alignment.
- P-145 TerraMod (ICSE 2026, Wave-21) — source-to-source modernisation via LLM agents. TransAgent shares the multi-agent substrate but specialises for translation with execution-alignment critic.
- P-124 Syzygy (ICSE 2025, Wave-17) — Go-to-Rust iterative compilation feedback. TransAgent inherits the compilation-feedback substrate but adds execution alignment.
- P-131 SemArc (Wave-18) — semantic arc recovery for LLM code comprehension. TransAgent shares the semantic-comparison substrate.

## Hop-2 Citations

- UniTrans (Yan et al., ASE 2023) — code-translation baseline that TransAgent outperforms by 33.3% in accuracy.
- Agentless (Xia et al., 2024) — code-repair baseline that TransAgent outperforms by 56.7% on average in repair performance.
- Multi-agent code generation lineage (Wu et al., 2023) — multi-agent collaboration substrate. TransAgent specialises for translation with role-decomposed agents.
- Execution-alignment translation surveys — execution-behaviour-comparison substrate. TransAgent is the first multi-agent + execution-aligned translation instance.

## MAgHARCM integration

- **YAML config key**: `agents.translation.execution_alignment: true`; `agents.translation.role_decomposition: ["translator", "critic", "repairer"]`; `agents.translation.alignment_granularity: "block_level"`.
- **Implementation file**: `internal/translation/execution_aligned.go::NewExecutionAlignedCritic` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-29]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **Multi-agent cost** — three-agent collaboration triples the inference cost vs single-LLM; must be amortised over the translation-accuracy gain.
- **Execution-trace availability** — fine-grained execution alignment requires the source code to be executable; non-executable source (e.g., legacy code with missing dependencies) is out of scope.
- **Block-level granularity** — the alignment granularity is block-level, not line-level; sub-block errors (intra-block) are not localised.
- **Language coverage** — TransAgent is evaluated on translation tasks; generalisation to other multi-agent code tasks (generation, repair) is future work.

## Source

- Venue: FSE 2026 (Research Track), verified via dl.acm.org/doi/abs/10.1145/3797099; PACMSE Vol. 3, Issue FSE, DOI: 10.1145/3797099. Held July 5-9 2026 in Montréal, Canada.
- arXiv: 2409.19894 (preprint).

## BibTeX

```
@inproceedings{transagent2026fse,
  title     = {TransAgent: Enhancing LLM-Based Code Translation via Fine-Grained Execution Alignment},
  author    = {Yuan, Zhiqiang and Chen, Weitong and Wang, Hanlin and Peng, Xin and Chen, Zhenpeng and Lou, Yiling},
  booktitle = {Proceedings of the ACM International Conference on the Foundations of Software Engineering (FSE)},
  year      = {2026},
  volume    = {3},
  number    = {FSE},
  series    = {Proceedings of the ACM on Software Engineering (PACMSE)},
  doi       = {10.1145/3797099}
}
```
