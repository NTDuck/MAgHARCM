---
title: "P-61 — Guo et al. 2024 — RedCode: Risky Code Execution Benchmark"
backlink: "[[1.0.0 P-61]]"
aliases:
  - "1.0.0 P-61"
  - "P-61"
  - "P-61-Guo-Risky-Code-Execution-2024"
  - "P-61-Guo-Risky-Code-Execution-2024"
  - "Guo-Risky-Code-Execution-2024"
tags: [paper, adversarial, red-team, code-generation, risky-execution, verifier-of-verifier, [[1.0.0 PRIM-7]], [[1.0.0 PRIM-25]], [[1.0.0 P-57]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-61 — Guo et al. — RedCode]]

## Citation

Guo, J. et al. (2024). *RedCode: Risky Code Execution Benchmark*. NeurIPS 2024 Datasets and Benchmarks Track. arXiv:2411.07781. URL: https://arxiv.org/abs/2411.07781.

## Summary

RedCode is a benchmark for evaluating the *risky execution* behavior of LLM-generated code: how often does an LLM produce code that, when run, causes system-level damage (file deletion, network exfiltration, privilege escalation, persistence mechanisms)? The benchmark has ~4,000 risky-coding prompts across four risk categories (file-system, network, OS/privilege, persistence) and ten programming languages, scored against the actual outcome when the generated code is sandbox-executed.

The benchmark reveals four empirical findings that directly inform verifier-of-verifier architectures:

1. **Capability-paradox.** Stronger models (GPT-4o, Claude-3.5-Sonnet, Qwen2.5-Coder-32B) produce *more* risky code on absolute count than weaker models, because they generate more code overall. Risk rate per generated-line is similar, but absolute exposure grows with capability.
2. **Rejection-asymmetry.** Subtle-buggy code (e.g., a single-character path-traversal slip) slips past obvious-danger detectors (regex-based "is this a rm -rf?" filters) at much higher rates than obvious-danger code. The rejection accuracy distribution is heavily skewed toward the dangerous end, leaving the subtle end unprotected.
3. **Format-sensitivity.** The same defect exposed in different code-fragment formats (function definition, snippet, library call, configuration) is accepted at different rates by safety classifiers. The classifier's "safe / unsafe" output is unstable across formats.
4. **Framework-brittleness.** LLMs trained with one safety framework (e.g., RLHF on red-teamed refusal data) leak more aggressively under a different framework (e.g., when wrapped in a different code-execution sandbox).

## Findings Relevant to MAgHARCM

- **Capability-paradox ⇒ single-model self-check is insufficient.** A self-judging reviewer (the same model that generated the code also judging it) inherits the model's blind spots. MAgHARCM's [[1.0.0 PRIM-25]] (Role-Flip De-Hallucination Gate) addresses this by priming the reviewer to find defects, but the same model is still the reviewer. RedCode shows that the role-flip helps with surface defects but not with the capability-paradox at the model level.
- **Rejection-asymmetry ⇒ verdict panel must aggregate.** RedCode's subtle-buggy slips past obvious-danger detectors. MAgHARCM's [[1.0.0 PRIM-7]] (Multi-Agent Verdict Panel) aggregates multiple reasoning models: when the panel disagrees, the verdict is "defect" by default, which addresses rejection-asymmetry.
- **Format-sensitivity ⇒ role-flip reviewer is a structured-output harness.** MAgHARCM's role-flip reviewer's response format (defect / no-defect + structured repair hint) is more stable across code-fragment formats than open-ended prompts. RedCode's format-sensitivity finding justifies this design choice.
- **Framework-brittleness ⇒ optional checks.** MAgHARCM's optional-checks pattern (`internal/agents/optional_checks.go`) — each check is independently enableable — is the natural mitigation for framework-brittleness: a check that consistently fails on a new code-fragment format can be disabled and replaced with a new check without affecting the rest of the pipeline.

## How MAgHARCM Uses It

[[1.0.0 PRIM-7]] (Multi-Agent Verdict Validation) is MAgHARCM's defense against capability-paradox: the verdict panel is a different reasoning model than the Translator, so the verdict does not inherit the Translator's blind spots. [[1.0.0 PRIM-25]] (Role-Flip Gate) is the defense against format-sensitivity: the role-flipped reviewer uses structured output. [[1.0.0 P-57]]'s lossless accept-reject discipline is the algebraic structure that keeps draft-model bias (introduced by speculative decoding on the panel judges) from propagating into final verdicts.

## References

### Hop-1
- Leviathan, Y., Kalman, M., Matias, Y. (2023). See [[1.0.0 P-57]] — lossless accept-reject discipline.
- Anthropic (2025) — Sycophancy evaluation. See [[1.0.0 P-38]] for the baseline paper; cf. [[1.0.0 P-69]] for the 2025 update and [[1.0.0 P-77]] for the OpenAI 2024 red-teaming methodology that complements this line.
- OpenAI (2024) — *Advancing Red Teaming with People and AI*. See [[1.0.0 P-77]]. (cited inline as the Codex-red-team analogue; OpenAI did not publish a Codex-specific red-team report.)
- Perry, N., et al. (2023). *Do Users Write More Insecure Code with AI Assistants?* ACM CCS 2023.

### Hop-2
- Chen, M., et al. (2021). *Evaluating Large Language Models Trained on Code* (Codex/HumanEval).
- Hubinger, E., et al. (2024). *Sleeper Agents: Training Deceptive LLMs That Act Through Their Training*. arXiv:2401.05566. Anthropic.
- Carlini, N., et al. (2024). *Are aligned neural networks adversarially aligned?* IEEE S&P 2024.

## Backlinks

[[1.0.0 PRIM-7]], [[1.0.0 PRIM-25]], [[1.0.0 PRIM-13]], [[1.0.0 P-38]], [[1.0.0 P-57]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].
