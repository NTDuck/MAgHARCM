---
title: "P-101 — Zhou et al. 2023 — Least-to-Most Prompting Enables Complex Reasoning in Large Language Models"
backlink: "[[1.0.0 P-101]]"
tags: [paper, prompting, least-to-most, decomposition, reasoning, slm, [[1.0.0 PRIM-1]], [[1.0.0 PRIM-23]], hop-1]
---

# [[1.0.0 P-101 — Zhou et al. — Least-to-Most Prompting]]

## Citation

Zhou, D., Schuurmans, D., Li, X., Chi, E. H., Wang, S. Z., & Zhou, D. (2023). *Least-to-Most Prompting Enables Complex Reasoning in Large Language Models*. International Conference on Learning Representations (ICLR 2023). arXiv:2205.10625 (21 May 2022 v1; 29 Jul 2022 v2 with appendices). Google Research / UC Berkeley / UNC.

## Summary

The Zhou et al. (2023) paper extends chain-of-thought prompting with an explicit *sequential decomposition* strategy. The core observation: standard few-shot CoT (cf. `[[1.0.0 P-90]]`) prompts the model with a fixed exemplar pattern — `Q → reasoning → A` — but the model's reasoning is constrained to whatever structure the exemplars suggest. When the target problem is *harder than the exemplars*, the model imitates the exemplar's length and gives up.

Least-to-Most (LtM) prompting instead constructs an *ordered sequence of subproblems* whose answers feed forward:

1. **Decomposition stage**: a few-shot prompt asks the model to break the target problem `Q` into a list of subproblems `Q1, Q2, …, Qn`, where each `Qi` is intended to be *easier* than `Q` and the conjunction of their answers solves `Q`.
2. **Sequential solve stage**: the model then solves `Q1`, appends `(Q1, A1)` to the prompt, and uses this augmented context to solve `Q2`, appends again, and so on. Each subproblem's solution is conditioned on all prior subproblem solutions.

This *chained-answer* pattern is the key insight — answers are not solved in isolation but **conditioned on the prefix of accumulated solutions**, which lets information discovered in early subproblems propagate to later ones. Standard CoT has no analogous mechanism: it commits to one reasoning chain and never branches.

**Headline results on GPT-3 `code-davinci-002` (≈175B) and `text-davinci-002` (≈175B):**
- **SCAN compositional split**: LtM hits ~99.7% accuracy on length-splits that standard prompting solves at ~16%. The compositional generalization gap collapses from ~84pp to ~0.7pp.
- **GSM8K (grade-school math)**: ~62.4% pass@1 with LtM vs ~49.5% for standard CoT prompting — a ~13pp absolute gain at matched model scale.
- **Last-letter concatenation** (symbolic, 4–7 word inputs): standard prompting collapses beyond 4 words (~30% at 7 words); LtM reaches ~94% at 7 words.
- **Color/object reasoning tasks**: LtM solves compositional variants that defeat one-shot prompting.

**Why it matters**: LtM is the *first decomposition-prompting technique* that propagates intermediate answers, distinguishing it from P-90 CoT (single chain), `[[1.0.0 P-62]]` Decomposed Prompting (Khot et al. — a *modular* decomposition with separate sub-task prompts, no sequential conditioning), and `[[1.0.0 P-96]]` Kojima zero-shot CoT ("Let's think step by step" — no decomposition at all). LtM sits between CoT and full modular decomposition.

## Method

- **Two-stage prompt**:
  - *Stage 1 (decomposition)*: `Q → [list of subproblems Q1, Q2, …, Qn]` — typically 4–8-shot exemplars showing how to factor `Q` into easier `Qi`s.
  - *Stage 2 (sequential solve)*: iteratively append `(Qi, Ai)` to the prompt, asking the model to solve each subsequent `Qi` in the context of the growing answer prefix. Output the final answer `A` after solving all `Qi`.
- **No fine-tuning**: pure inference-time prompting. Two-pass inference (decompose → solve).
- **Models studied**: GPT-3 `code-davinci-002` (175B), `text-davinci-002` (175B), `text-davinci-003` (175B). LtM is scale-sensitive — at sub-30B the decomposition quality itself becomes the bottleneck.
- **Critical design choice**: the *ordering* of subproblems matters. The paper's `least-to-most` order is **reverse-topological**: solve prerequisites first, then the target. Each subproblem's solution is **prefix-conditioned** on all earlier answers — not sampled independently.
- **Distinguishing feature from P-62 Khot Decomposed Prompting**: P-62 uses *independent* sub-prompts per sub-task (no chaining). LtM uses *chained* sub-prompts that condition on prior answers — this is what enables compositional generalization on SCAN.
- **Chained-answer vs. parallel-answer ablation**: the paper shows parallel sampling of independent subproblem answers *degrades* on SCAN (~16% accuracy) compared to chained (~99.7%) — prefix-conditioning is not optional.

## Findings Relevant to MAgHARCM

- **Reverse-topological planning maps cleanly onto `[[1.0.0 PRIM-1]]` Reverse-Topological Planner.** LtM's `least-to-most` order is exactly reverse-topological: solve the leaves (callees / imports / inner types) before the root (the API surface / top-level function). MAgHARCM's `internal/agents/planner.go` already builds a topo DAG of translation units; LtM is the *prompt-side* template that asks the SLM to *emit* such a decomposition. This is the missing link between MAgHARCM's structural topo sort and the SLM's output planning.
- **Chained-answer propagation = `[[1.0.0 PRIM-23]]` Chunked Translation.** When a large source file is split into `N` chunks for translation, naive parallel translation of all `N` chunks produces inconsistent types and import bindings (because each chunk is solved in isolation). LtM's chained pattern — solve chunk 1, append `(chunk1, translated1)` to the prompt, solve chunk 2 in that context — gives the SLM the *cross-chunk visibility* it needs to keep type names and module imports coherent. MAgHARCM's chunked-translation routine (cf. `internal/agents/translator.go`) should adopt the LtM chained pattern as its default for SLMs ≥30B.
- **Scale sensitivity gates deployment.** LtM's decomposition step is itself a generation task — at sub-30B SLMs, the decomposition quality is poor and the chained-solve stage amplifies errors. MAgHARCM should gate LtM behind the same scale check as P-90 CoT: use LtM on Qwen2.5-Coder-32B [[1.0.0 P-21]] and s1-32B [[1.0.0 P-84]]; fall back to direct CoT or whole-file translation on smaller SLMs.
- **Two passes are cheap on large models but expensive on small ones.** LtM doubles the inference cost (decompose + solve). On 32B SLMs at ~30 tok/s, this adds ~10–20s to a typical migration step. On 7B SLMs at ~80 tok/s, the same overhead is more acceptable. MAgHARCM's cost model should expose this as a per-strategy overhead line item.
- **LtM is a strict generalization of P-90 CoT.** When the decomposition stage emits a *single* subproblem equal to `Q`, LtM reduces to standard CoT. This means MAgHARCM can implement LtM *and* fall back to CoT within the same code path: if the decomposition stage returns `[Q]`, skip the chained-solve loop and emit one CoT-style response.
- **Chained answers handle cross-chunk state that `[[1.0.0 P-98]]` (Large Language Monkeys parallel sampling) does not.** P-98's repeated-sampling baseline cannot propagate state across chunks; LtM's chained prefix-conditioning is what MAgHARCM needs for chunked translation across language boundaries where type/import visibility matters.
- **`[[1.0.0 P-62]]` Decomposed Prompting vs. LtM: choose by problem structure.** P-62 modular decomposition is stronger when subproblems are *truly independent* (e.g., parallel API endpoint translations with no shared types). LtM is stronger when subproblems have *dependencies* (e.g., a chunk that references a type defined in a prior chunk). MAgHARCM's chunker should classify chunks as `independent` (P-62 path) or `dependent` (LtM path) before routing.

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-1]]` Reverse-Topological Planner**: replace the current static top-down SLM prompt ("translate this file top-to-bottom") with the LtM decomposition prompt — ask the SLM to first emit the dependency order of the translation units, then solve in that order. Cite P-101 as the chained-decomposition template. Combine with MAgHARCM's existing structural topo DAG (computed by `internal/agents/planner.go`) as a *cross-check*: if the SLM-emitted order disagrees with the static DAG, prefer the DAG and ask the SLM to explain the discrepancy.
- **`[[1.0.0 PRIM-23]]` Chunked Translation**: the chunked-translation primitive in `internal/agents/translator.go` should adopt LtM's chained pattern. Instead of translating each chunk independently and concatenating, run a chained pass: solve chunk 1, append `(chunk1, translated1)`, solve chunk 2 in that context, etc. Fall back to independent-parallel chunked translation (the P-100 path) only when chunks are classified as `independent` by the chunker. This is the single largest expected accuracy gain on cross-language migrations with shared types.
- **`[[1.0.0 PRIM-22]]` Four Phases of Comprehension** (extension): add a *decomposition phase* between *search* and *explanation* that asks the SLM to emit the comprehension plan as an LtM-ordered sub-question list, then answer each sub-question sequentially. This makes comprehension reproducible across SLMs and exposes the planning step for audit.
- **Scale gate**: LtM is gated to SLMs ≥30B (Qwen2.5-Coder-32B, s1-32B). On smaller SLMs, MAgHARCM uses P-90-style direct CoT (no decomposition) or P-100-style parallel chunking. The cost of decomposition on sub-30B SLMs exceeds the benefit.
- **Future work (not implemented)**: P-101 + `[[1.0.0 P-91]]` Snell test-time-scaling — at high compute budgets, run LtM with N parallel decomposition candidates and select the most coherent topo order before chaining the solve. This is the "sequential revise" extension of LtM in P-91's framework.

## References

### Hop-1 (papers that build directly on P-101)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. NeurIPS 2022. arXiv:2201.11903. See `[[1.0.0 P-90]]` — the foundational CoT paper LtM extends; LtM reduces to standard CoT when decomposition emits a single subproblem.
- Khot, T. et al. (2022). *Decomposed Prompting: A Modular Approach for Solving Complex Tasks*. arXiv:2210.02406. See `[[1.0.0 P-62]]` — the *modular* decomposition counterpart to LtM's *chained* decomposition; choose by problem structure (P-62 for independent subproblems, P-101 for dependent ones).
- Brown, B. et al. (2024). *Large Language Monkeys: Scaling Inference Compute with Repeated Sampling*. arXiv:2407.21787. See `[[1.0.0 P-100]]` — the *parallel-sampling* baseline P-101 outperforms on compositional tasks; P-100 is the right choice when subproblems are genuinely independent.

### Hop-2 (foundational anchors referenced)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. See `[[1.0.0 P-90]]` — the prompting substrate LtM composes onto; LtM is a *strict generalization* of P-90 in the single-subproblem degenerate case.

### MAgHARCM lineage cross-refs
- `[[1.0.0 P-21]]` — Qwen2.5-Coder-32B (the 32B SLM where LtM's decomposition quality is reliable).
- `[[1.0.0 P-52]]` — Self-Consistency (samples multiple CoT traces; P-101 + P-52 is "LtM with self-consistency over decompositions").
- `[[1.0.0 P-62]]` — Decomposed Prompting (the modular sibling of P-101).
- `[[1.0.0 P-84]]` — s1 (test-time-scaling CoT traces that can host LtM-style decomposition as a planning prelude).
- `[[1.0.0 P-90]]` — Chain-of-Thought (the prompting substrate).
- `[[1.0.0 P-91]]` — Snell test-time-scaling (compute-optimal allocator; future work: LtM × compute-optimal).
- `[[1.0.0 P-96]]` — Kojima zero-shot CoT (the zero-shot alternative; P-101 is the few-shot, sequential-decomposition alternative).
- `[[1.0.0 P-100]]` — Large Language Monkeys (the parallel-sampling alternative for independent-chunk paths).

## Backlinks

`[[1.0.0 PRIM-1]]`, `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-23]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-101 is the **chained-decomposition** anchor for MAgHARCM's planning + chunked-translation primitives. It bridges the gap between MAgHARCM's structural topo DAG (`internal/agents/planner.go`) and the SLM's emitted planning, and gives `[[1.0.0 PRIM-23]]` chunked translation the cross-chunk state propagation it needs. Cite P-101 whenever a primitive decomposes a problem into an ordered subproblem list with prefix-conditioned sequential solve.
