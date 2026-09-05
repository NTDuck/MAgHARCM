---
title: "P-104 — S* Test-Time Scaling for Code Generation (arXiv 2025)"
backlink: "[[1.0.0 P-104]]"
tags: [paper, test-time-compute, code-generation, sampling, verifier, [[1.0.0 PRIM-7]], [[1.0.0 PRIM-21]], hop-1]
---

# [[1.0.0 P-104 — S*: Test-Time Scaling for Code Generation]]

## Citation

Anonymous (2025). *S*: Test Time Scaling for Code Generation*. arXiv:2502.14382.
URL: https://arxiv.org/abs/2502.14382

## Summary

S* is the **first hybrid test-time scaling framework** specifically for code generation. Prior test-time scaling work (P-84 s1, P-91 Snell, P-98 LLM Monkeys) focused on math/reasoning; S* extends it to code with three scaling modes:
1. **Sequential refinement**: generate → execute → diagnose → regenerate.
2. **Parallel sampling**: N candidates → pick best via execution signal.
3. **Hybrid**: parallel sample → sequential refine top-K.

Key empirical finding: hybrid mode achieves **2.3× pass@1 improvement** at matched compute budget vs sequential-only or parallel-only. The execution signal (compile + test) is the dominant cost driver, not sampling.

## Relevance to MAgHARCM

- **Direct evolution of P-98 (LLM Monkeys) for code**: validates MAgHARCM's verdict-panel-as-sampling-strategy pattern (PRIM-7) and adds execution-signal guidance for strategy selection (PRIM-21).
- **S* hybrid mode = MAgHARCM verdict panel + repair loop**: verdict_panel samples, recruiter chooses next strategy, translator refines — exactly the hybrid mode.
- **PRIM-27 Coverage-Guided Plateau Detection** can use S*-style execution feedback as a stop signal: when sequential refinement's marginal gain drops below threshold, switch to next strategy.
- **Test-budget allocator**: S* provides a concrete algorithm for allocating test-time compute across strategies (vs static "max_iterations: 3" in `configs/agents.yml`).

## References (hop-1)

- [[1.0.0 P-84]] s1 test-time scaling
- [[1.0.0 P-91]] Snell test-time compute allocation
- [[1.0.0 P-98]] Large Language Monkeys
- [[1.0.0 P-92]] Lightman PRM (process reward model for code)

## References (hop-2)

- P-97 Welleck Self-Correct 2024
- P-83 code-specialised self-consistency
