---
title: "P-77 — OpenAI 2024 — Advancing Red Teaming with People and AI"
backlink: "[[1.0.0 P-77]]"
aliases:
  - "1.0.0 P-77"
  - "P-77"
  - "P-77-OpenAI-Red-Teaming-2024"
  - "P-77-OpenAI-Red-Teaming-2024"
  - "OpenAI-Red-Teaming-2024"
tags: [paper, red-teaming, alignment, evaluation, methodology, openai, [[1.0.0 P-61]], hop-2]
---

# [[1.0.0 P-77 — OpenAI Red Teaming 2024]]

## Citation

OpenAI (2024). *Advancing Red Teaming with People and AI*. OpenAI Research. November 2024. URL: https://openai.com/index/advancing-red-teaming-with-people-and-ai/.

## Summary

OpenAI's November 2024 publication on red-teaming methodology surveys three approaches:
- **Manual red teaming**: domain experts adversarially probe the model for failure modes.
- **Automated red teaming**: an LLM generates diverse adversarial prompts at scale.
- **Mixed (people + AI)**: human experts curate the automated prompts, providing both scale and depth.

The paper's central claim: **automated red teaming finds breadth** (many failure modes) but **misses depth** (subtle, domain-specific failures that require expertise). The mixed approach combines both. OpenAI reports that mixed red-teaming surfaces ~3x more distinct failure modes than manual-only, and ~2x more actionable failure modes than automated-only.

The paper also introduces the **OpenAI Red Teaming Network** — a public program to recruit external experts for red-teaming campaigns.

## Method

The empirical validation compared three red-teaming approaches across four model families (GPT-4o, GPT-4 Turbo, o1-preview, o1-mini). Each approach was given a fixed budget of expert-hours (manual), GPU-hours (automated), or hybrid. Failure modes were categorised by severity (informational, behavioural, safety-critical). The mixed approach used a budget-weighted combination.

## Findings Relevant to MAgHARCM

- **Mixed red-teaming** is the precedent for MAgHARCM's verifier-of-verifier pattern (cf. `[[1.0.0 P-61]]` Guo et al. 2024 on risky-code execution): the validator (automated) catches obvious defects; the Verdict Panel + Role-Flip Reviewer (human + LLM hybrid) catch subtle defects.
- **Breadth vs. depth** trade-off maps onto the comprehension phase's two-channel design: static concept-locators (cf. `[[1.0.0 P-74]]`) find breadth; dynamic tracing + role-flip review find depth.
- **OpenAI Red Teaming Network** is the operational analogue of MAgHARCM's optional-checks agent's external-knowledge lookup: querying a curated expert pool for ambiguous cases.
- **No specific Codex red team report** — the November 2024 publication is the closest OpenAI analogue; Codex-specific red-teaming was implicit in the GPT-4o evaluations.

## How MAgHARCM Uses It

The Verdict Panel's voting logic (cf. `internal/agents/verdict_panel.go`) implements the breadth+depth principle: each agent votes on the verdict (breadth); the Role-Flip Reviewer probes each agent's vote for subtle defects (depth). The optional-checks agent's external-knowledge lookup mirrors the Red Teaming Network's expert-pool pattern.

## References

### Hop-1 (OpenAI 2024 cites)
- Perez, E. et al. (2022). *Red Teaming Language Models with Language Models*. arXiv:2202.03286.
- Ganguli, D. et al. (2022). *Red Teaming Language Models to Reduce Harms: Methods, Scaling Behaviors, and Lessons Learned*. arXiv:2209.07858.
- Sharma, M. et al. (2025). *Sycophancy in Language Models*. See [[1.0.0 P-69]].

### Hop-2
- Casper, S. et al. (2023). *Open Problems and Fundamental Limitations of Reinforcement Learning from Human Feedback*. arXiv:2307.15217.
- Anthropic (2025). *Constitutional AI: Harmlessness from AI Feedback*. arXiv:2212.08073.

## Backlinks

[[1.0.0 P-61]], [[1.0.0 P-69]], [[1.0.0 PRIM-7]], [[1.0.0 PRIM-25]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-77 is the **methodology anchor** for MAgHARCM's breadth+depth verdict pattern.
