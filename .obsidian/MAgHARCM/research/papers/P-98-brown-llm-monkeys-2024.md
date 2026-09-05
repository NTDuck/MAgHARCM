---
title: "P-98 — Brown et al. 2024 — Large Language Monkeys: Scaling Inference Compute with Repeated Sampling"
backlink: "[[1.0.0 P-98]]"
tags: [paper, repeated-sampling, verifier, slm, best-of-n, inference-compute, code-generation, role-flip, [[1.0.0 PRIM-7]], [[1.0.0 PRIM-21]], [[1.0.0 PRIM-25]], hop-1]
---

# [[1.0.0 P-98 — Brown et al. — Large Language Monkeys]]

## Citation

Brown, B., Juravsky, J., Ehrlich, R., Clark, R., Rezende LeQuesne, Q., Mohan, A., & Pappu, A. (2024). *Large Language Monkeys: Scaling Inference Compute with Repeated Sampling*. arXiv:2407.21787 (31 Jul 2024 v1; 3 Oct 2024 v2). Stanford University. Not formally peer-reviewed; treated as a foundational preprint widely cited in the SLM-era repeated-sampling / verifier-scaling line (P-91, P-92).

## Summary

The Brown et al. paper poses the question: **can you buy accuracy by simply sampling more times, instead of scaling the model?** The answer is a qualified but striking *yes* for code and coding-adjacent tasks. By drawing $N$ repeated samples from a single LLM and scoring them with a lightweight verifier, the paper shows that "**repeated sampling + verifier**" competes with, and often beats, trading the same compute budget for a larger base model.

The core experimental setup is deliberately controlled:
- **Models**: public foundation models sampled at multiple sizes (a 70B-class model, a 34B-class model, and a 7B-class **small language model (SLM)**).
- **Tasks**: code generation (HumanEval code-completion and code-repair subsets), math (MATH500), plus a set of "coding" tasks where the verifier is programmatic.
- **Verifier**: on code, an **execution-based unit-test verifier** (pass/fail against reference tests); on math, a learned correctness model.
- **Compute budget**: total inference compute held fixed while $N$ (the number of samples) grows and the base model shrinks (or vice versa).

Key findings:

1. **Sampling beats scaling on code up to a regime change.** On the code-repair and code-replacement tasks, sampling $N$ candidates from a **smaller** model and selecting with the verifier outperforms a **larger** model sampled once — at matched inference FLOPs — for $N$ up to ~10⁴. This is the exact FLOP-for-FLOP argument that later papers (P-91) generalise into compute-optimal scaling.
2. **The verifier, not the sampler, is where the gains come from.** The paper shows that *with* a verifier, accuracy rises roughly logarithmically in $N$ and continues rising well past $N$ where raw pass-rate plateaus; *without* a verifier, repeated sampling adds almost nothing. The verifier converts a flat sampling distribution into a monotonically improving selector.
3. **Embedding-distance verifiers rival exact unit tests on code.** For tasks where reference execution tests are unavailable, the paper demonstrates that an **embedding-based similarity verifier** (comparing candidate output embeddings to reference-answer embeddings) approaches the accuracy of the exact unit-test verifier. This matters for MAgHARCM because it de-risks sampling-based verification when a full execution oracle is absent.
4. **Small-model repetition is a distinctive SLM strategy.** The 7B SLM benefit most from repeated sampling precisely because their single-sample pass rate is lower — the sampling headroom is larger. The paper's cost framing makes this concrete: a fixed token budget buys many 7B samples but only one or two 70B samples, and the 7B-many-samples allocation wins on code.

## Method

- **Sampling regime**: draw $N$ i.i.d. samples $y_1..y_N$ from a single LLM given one prompt $x$; return the sample that maximises the verifier score. $N$ ranges from ~10⁰ to ~10⁵.
- **Verifier options**:
  - *Programmatic unit-test* (exact): candidate passes iff it satisfies the reference tests. Perfect accuracy on the labelled eval set.
  - *Embedding-distance* (approximate): run candidate and reference through a shared embedding model; score by cosine/euclidean similarity. No execution required.
  - *Learned reward model* (math): a correctness classifier over final answers, in the family of P-92 ORMs.
- **Compute-holding**: for each model pair (small/large) and each $N$, the total FLOPs (samples × per-sample forward cost) is matched so comparisons isolate the sampling-vs-scale trade-off.
- **Difficulty stratification**: results are reported by task-difficulty quartile; the sampling benefit is concentrated on middle-difficulty tasks where the verifier can actually rank candidates (too-easy tasks saturate, too-hard tasks have no correct sample in the bucket).

## Findings Relevant to MAgHARCM

- **`[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation is the MAgHARCM instantiation of verifier-best-of-N.** P-98's repeated-sampling + verifier is exactly what the verdict panel's vote-with-execution-oracle does at the repository-translation layer: each SLM proposal is a sample, the `[[1.0.0 PRIM-12]]` Wasm-based reference-execution oracle is the perfect verifier. P-98 justifies *why* fanning out more SLM samples to the oracle is a better accuracy buy than upgrading the base SLM alone.
- **The embedding-verifier result matters for pre-execution routing.** MAgHARCM's repository translations must often be scored *before* a full item can be executed (module-level in-isolation). P-98's embedding-similarity verifier (Section "Method") gives a cheap, execution-free pre-filter: rank candidate chunks by embedding distance to a reference skeleton, keep the top-$k$, run the oracle only on those. This complements, not replaces, the exact oracle.
- **SLM sampling is the cheapest accuracy lever on the fleet.** For the 7B-class SLMs MAgHARCM deploys (Qwen2.5-Coder-7B, cf. `[[1.0.0 P-21]]`), P-98 shows the largest relative sampling headroom — a fixed token budget buys many 7B samples and wins over one 70B sample on code. This is the quantitative basis for MAgHARCM's SLM-fleet thesis.
- **Difficulty stratification maps to `[[1.0.0 PRIM-21]]` strategy selection.** The sampling benefit concentrates on middle-difficulty tasks; too-easy tasks need one good sample, too-hard tasks need a different strategy (revision, cf. P-91's sequential revise). MAgHARCM's `[[1.0.0 PRIM-21]]` router should not uniform-N sample; it should vary $N$ with estimated difficulty.
- **The role-flip reviewer is the cheapest practical verifier.** P-98 assumes a verifier exists (unit tests or embeddings). When neither is available, MAgHARCM's `[[1.0.0 PRIM-25]]` Communicative De-Hallucination Role-Flip Gate *is* the verifier: a reviewer model inverted into a critical role scores each candidate for defects. P-98 makes the case that even a *noisy* role-flip verifier, repeated over more samples, beats a single unverified sample — so the role-flip gate should be run over the full sample bucket, not just the greedily-decoded one.

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation**: bind $N$ to the fan-out count of the verdict panel. Instead of a fixed 3-voter panel, adopt the P-98 pattern: draw $N$ SLM translation proposals, score each with the `[[1.0.0 PRIM-12]]` Wasm oracle, return the best-scoring. Cite P-98 as the repeated-sampling + verifier justification for the panel's design. When the oracle is unavailable for a chunk, fall back to the embedding-similarity pre-filter.
- **`[[1.0.0 PRIM-21]]` Migration Strategy Selection**: route $N$ (samples per chunk) by the strategy registry's `Profile`. Add a `sampling_budget` field (low/medium/high) driven by repetition-demand: independent easy chunks get low $N$ (P-62-style parallel), dependent hard chunks get high $N$ with oracle selection. P-98 is the anchor defining how many samples a given SLM size earns within a fixed token budget.
- **`[[1.0.0 PRIM-25]]` Role-Flip Gate**: run the gate as a *batch verifier* over the full sample bucket, not a per-candidate binary gate. Collect the role-flip reviewer's defect scores across all $N$ candidates and pick the lowest-defect one — converting the gate from a pass/fail checkpoint into a P-98-style selector. This is the cheapest verifier when no execution oracle and no embedding models are present.
- **Future work (not implemented)**: (a) benchmark the oracle-selected best-of-$N$ against the current 3-voter panel on `internal/agents/verdict_panel.go`; (b) train/select an embedding-similarity scorer for execution-free pre-filtering of in-isolation chunks; (c) add per-chunk `sampling_budget` to the `Profile` in `internal/agents/strategy.go`.

## References

### Hop-1 (papers that build directly on P-98)
- Snell, C. et al. (2024). *Scaling LLM Test-Time Compute Optimally Can be More Effective than Scaling Model Parameters*. arXiv:2408.03314. See `[[1.0.0 P-91]]` — generalises P-98's repeated-sampling + verifier into a compute-optimal allocator that chooses between parallel sampling (the P-98 path) and sequential revision per query. P-98 is the parallel-sampling special case of P-91's framework.
- Lightman, H. et al. (2023). *Let's Verify Step by Step*. arXiv:2305.02188. See `[[1.0.0 P-92]]` — the process/outcome reward models that act as the verifier component in P-98's repeated-sampling loop; P-98 shows the verifier (not the sampler) is where accuracy comes from.
- Li, Y. et al. (2025). *EAGLE-3: Scaling up Inference Acceleration of LLMs via Training-Time Test*. arXiv:2503.01840. See `[[1.0.0 P-78]]` — makes the many-sample regimes of P-98 affordable: at ~6.5× lossless speedup over vanilla decoding, the per-sample cost of the $N$-bucket drops enough that best-of-$N$ selection runs on consumer SLM hardware.
- Welleck, S. et al. (2024). *Generating Sequences by Learning to Self-Correct*. ICLR 2024. See `[[1.0.0 P-97]]` — the trained-self-correction alternative: instead of sampling many candidates and verifying (P-98), train the base model to correct its own single output, trading inference compute for amortised training compute.
- Muennighoff, N. et al. (2025). *s1: Simple Test-Time Scaling*. arXiv:2501.19393. See `[[1.0.0 P-84]]` — budget-forcing allocates test-time compute within a *single* generation; P-98 allocates it *across* samples. Both are test-time-compute strategies that the P-91 allocator routes between.

### Hop-2 (foundational anchors referenced)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. See `[[1.0.0 P-90]]` — the per-sample reasoning substrate whose accuracy P-98 amplifies by sampling; CoT raises the per-sample floor that repeated sampling then exploits.
- Wang, X. et al. (2023). *Self-Consistency Improves Chain of Thought Reasoning in Language Models*. ICLR 2023. arXiv:2203.11171. See `[[1.0.0 P-52]]` — the parallel-sampling-with-aggregation ancestor of P-98. Where P-52 majority-votes the *plurality* answer, P-98 uses a verifier to select the *best* single answer; P-98 generalises P-52 by replacing the voting rule with a verifier score.

### MAgHARCM lineage cross-refs
- `[[1.0.0 P-91]]` — Test-Time Compute Scaling (the allocator whose parallel branch is P-98's repeated sampling).
- `[[1.0.0 P-92]]` — Lightman PRM (the verifier component P-98's loop consumes).
- `[[1.0.0 P-78]]` — EAGLE-3 (the inference-speedup that makes P-98's high-$N$ regimes affordable on SLMs).
- `[[1.0.0 P-52]]` — Self-Consistency (the sampling-aggregation predecessor P-98's verifier-selection generalises).
- `[[1.0.0 P-90]]` — Chain-of-Thought (the per-sample reasoning format P-98 amplifies).
- `[[1.0.0 P-21]]` — Qwen2.5-Coder (the 7B-class SLM in MAgHARCM's fleet where P-98's sampling headroom is largest).

## Backlinks

`[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-12]]`, `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-25]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-98 is the **repeated-sampling + verifier anchor** for MAgHARCM. It provides the quantified argument that sampling a small SLM many times and selecting with a verifier beats scaling the base model on code, and that the role-flip reviewer is a legitimate (noisy but effective) verifier when no execution oracle or embedding scorer is available. Cite P-98 whenever a primitive fans out multiple proposals and selects with a verifier, or whenever the verdict panel / role-flip gate budget how many samples to draw per query.
