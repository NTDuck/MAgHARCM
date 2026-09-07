---
title: "P-91 — Snell, Lee, Xu & Kumar 2024 — Scaling LLM Test-Time Compute Optimally Can be More Effective than Scaling Model Parameters"
backlink: "[[1.0.0 P-91]]"
aliases:
  - "1.0.0 P-91"
  - "P-91"
  - "P-91-Snell-Test-Time-Scaling-2024"
  - "P-91-Snell-Test-Time-Scaling-2024"
  - "Snell-Test-Time-Scaling-2024"
tags: [paper, test-time-compute, scaling, compute-optimal, flops, slm, [[1.0.0 PRIM-21]], [[1.0.0 PRIM-25]], hop-1]
---

# [[1.0.0 P-91 — Snell et al. — Scaling LLM Test-Time Compute]]

## Citation

Snell, C., Lee, J., Xu, K., & Kumar, A. (2024). *Scaling LLM Test-Time Compute Optimally Can be More Effective than Scaling Model Parameters*. arXiv:2408.03314 (6 Aug 2024 v1). Google DeepMind. Not formally peer-reviewed; treated as a foundational preprint widely cited in the s1 line of work (P-84).

## Summary

The Snell et al. paper poses a counterpoint to the prevailing "scale the model" paradigm: *a smaller model with more inference-time compute can outperform a larger model with less.* The headline contribution is a **compute-optimal** scaling strategy that, given a fixed inference-FLOPs budget, dynamically chooses between two strategies: (a) **distribute** the budget across many parallel samples from a small verifier-based proposal model, or (b) **sequentially revise** a single sample from a stronger revision model.

The paper's experimental setup is critical to interpret the results correctly:
- Base models: **PaLM 2-S** (small) and **PaLM 2-L** (large).
- Tasks: MATH (competition math), GSM8K (grade-school math), and code benchmarks.
- Verifier: an Outcome Reward Model (ORM) trained on MATH-style correctness labels.
- Compute budgets: FLOP-matched comparisons across {1×, 4×, 16×, 64×} inference compute.

**Key findings:**

1. **Test-time compute can substitute for model scale when paired with a strong verifier.** On MATH, the compute-optimal allocation of test-time compute to a small proposal model (with revision) outperforms a 14× larger model run with 5× less compute at the same total FLOPs. The crossover is dramatic: at high budgets, *the smaller model wins*.
2. **Difficulty determines strategy.** Easy problems benefit from parallel sampling + best-of-N; hard problems benefit from sequential revision. A single static strategy underperforms a compute-optimal mix.
3. **Verifier quality is the bottleneck.** When the verifier is noisy, additional compute yields diminishing returns. Snell et al. quantify this — the gain from 1× to 16× compute is largest when verifier accuracy is high.
4. **Process reward models (PRMs) outperform outcome reward models (ORMs)** when the same labeling budget is available — consistent with `[[1.0.0 P-92]]` Lightman et al.

The paper also introduces a **practical recipe**: estimate query difficulty using a prior step (e.g., pass-rate of greedy decoding from the small model), then route to {parallel, revision} based on the estimate. This is the compute-optimal scaling heuristic.

## Method

- **Compute budget variable**: total inference FLOPs, varied 1× to 64×.
- **Allocation choices**:
  - **Parallel**: sample N candidates from base model, score with verifier, return best.
  - **Sequential revise**: one candidate iteratively improved by a revision prompt, optionally verifier-gated.
- **Difficulty estimator**: greedy-decoding success rate of the base model on the prompt, or a learned prior.
- **Verifier**: ORM (final-answer correct/incorrect) or PRM (per-step correct/incorrect).
- **Evaluation**: pass@1 (greedy), pass@N (best-of-N with verifier), and FLOP-matched accuracy.

## Findings Relevant to MAgHARCM

- **Compute-optimal is the missing link between s1 and the verdict panel.** P-84 (s1) and P-83 (code-specialised SC) are *strategies*; P-91 is the *allocator* that decides which strategy to spend compute on per query. MAgHARCM's `[[1.0.0 PRIM-21]]` strategy registry should adopt a difficulty-aware router: easy queries → 2-shot SC; hard queries → revision loop with verifier feedback.
- **Sequential revision composes with `[[1.0.0 PRIM-25]]` Role-Flip Gate.** When the role-flip reviewer disagrees with the verdict panel, the paper's "sequential revise" strategy says: spend the extra compute on revising the original candidate *in light of the reviewer's critique* — not on a fresh sample. This is exactly the role-flip-revision pattern MAgHARCM already encodes.
- **Verifier quality gates the entire compute-optimal framework.** MAgHARCM's `[[1.0.0 PRIM-12]]` Wasm-Based Reference Execution Oracle is a perfect-execution oracle — strictly stronger than an ORM/PRM. This means MAgHARCM can apply more aggressive compute-optimal strategies than Snell et al.'s experiments, because the verifier is exact for compiled languages.
- **Difficulty estimation maps to `[[1.0.0 PRIM-21]]` Profile.** The strategy registry already computes a `Profile` (cf. `internal/agents/strategy.go`). Adding a `difficulty` field (easy / medium / hard, estimated from prior pass rate) and routing on it is the Snell-style compute-optimal move.
- **Cost framing is concrete.** At 16× compute, a 7B SLM can match a 14B SLM at 1× compute — about 2× slower wall time at the same accuracy. This is the *exact trade-off* MAgHARCM needs to surface to users: "do you want this migration in 2 min at 92% accuracy, or 1 min at 88%?"

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-21]]` Migration Strategy Selection**: add a `difficulty` field to `Profile`; route to `BIG_BANG` / `INCREMENTAL` / parallel sampling based on it. Cite P-91 as the compute-optimal allocator.
- **`[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation**: the verdict panel can adopt P-91's two-axis routing — easy translations get 3 voters; hard translations get 5 voters with a revision step between rounds.
- **`[[1.0.0 PRIM-25]]` Role-Flip Gate**: when the role-flip reviewer rejects, the repair loop's revision should follow P-91's sequential-revise pattern — feed the reviewer's critique back into the original candidate, not discard it.
- **`[[1.0.0 PRIM-12]]` Wasm-Based Oracle**: the execution oracle provides a *stronger* verifier than ORM/PRM; MAgHARCM should benchmark the P-91 framework against the execution oracle to quantify the gain (planned work, not implemented).
- **Future work (not implemented)**: implement a Snell-style compute-optimal router in `internal/agents/strategy.go`. Track per-query difficulty via rolling pass rate; allocate budget accordingly.

## References

### Hop-1 (papers that build directly on P-91)
- Muennighoff, N. et al. (2025). *s1: Simple Test-Time Scaling*. arXiv:2501.19393. See `[[1.0.0 P-84]]` — s1 builds on P-91's test-time-compute framework, generalising to budget forcing.
- Lightman, H. et al. (2023). *Let's Verify Step by Step*. arXiv:2305.02188. See `[[1.0.0 P-92]]` — the PRM800K verifier dataset P-91 uses for MATH.
- Brown, B. et al. (2024). *Large Language Monkeys: Scaling Inference Compute with Repeated Sampling*. arXiv:2407.21787. Demonstrates that repeated sampling + a verifier is competitive with scaling the base model on code tasks.
- Welleck, S. et al. (2024). *Generating Sequences by Learning to Self-Correct*. ICLR 2024. Self-correction training as an alternative to inference-time verifier search.
- Gehring, J. et al. (2024). *Self-Verification Improves Few-Shot Reasoning*. Cohere For AI. The verifier-as-critic pattern.

### Hop-2 (foundational anchors referenced)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. See `[[1.0.0 P-90]]` — the prompting substrate that test-time-compute scaling operates on.
- Wang, X. et al. (2023). *Self-Consistency Improves Chain of Thought Reasoning in Language Models*. arXiv:2203.11171. See `[[1.0.0 P-52]]` — the parallel-sampling baseline P-91 generalises.
- Cobbe, K. et al. (2021). *Training Verifiers to Solve Math Word Problems*. arXiv:2110.14168. The verifier precursor.
- Anil, R. et al. (2023). *PaLM 2 Technical Report*. arXiv:2305.10403. The base models used.

### MAgHARCM lineage cross-refs
- `[[1.0.0 P-52]]` — Self-Consistency (the parallel-sampling anchor in P-91's framework).
- `[[1.0.0 P-84]]` — s1 (the test-time-scaling descendant of P-91).
- `[[1.0.0 P-92]]` — Lightman PRM (the strongest verifier P-91 evaluates).
- `[[1.0.0 P-83]]` — Code-specialised SC (parallel sampling + verifier on code).
- `[[1.0.0 P-90]]` — Chain-of-Thought (the prompting substrate).
- `[[1.0.0 P-93]]` — DPO (alternative alignment signal for the SLM base of the compute-optimal framework).

## Backlinks

`[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-12]]`, `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-25]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-91 is the **compute-optimal allocator** for MAgHARCM's test-time-compute budget. It provides the principled bridge between strategies (P-52, P-83, P-84) and the SLM deployment budget. Cite P-91 whenever a primitive makes a per-query allocation between parallel sampling and sequential revision.
