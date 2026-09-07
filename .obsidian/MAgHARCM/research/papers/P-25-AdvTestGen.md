---
title: Adversarial Test Generation for LLM-Based Code Translation Validation
backlink: "[[1.0.0 P-25]]"
bibkey: p25_advtestgen
aliases:
  - "1.0.0 P-25"
  - "P-25"
  - "P-25-AdvTestGen"
  - "P-25-AdvTestGen"
  - "AdvTestGen"
tags: [paper, adversarial-testing, translation-validation, [[PRIM-5]], [[PRIM-13]], hop-1]
---

# [[1.0.0 P-25]] Adversarial Test Generation for LLM-Based Code Translation Validation

**Authors**: [[AdvTestGen-Team-2024]]
**Year**: 2024
**Venue**: ACM International Symposium on Software Testing and Analysis (ISSTA) 2024
**eprint / DOI**: (pending — to be confirmed against ISSTA 2024 proceedings)
**Cited by**: [[primitives/Primitives-Index]] entry [[PRIM-13]] (Adversarial Test-Weakening Guard), [[PRIM-5]] (Test Suite Co-Translation & Synthesis)

## Summary

[[AdvTestGen-2024-ISSTA]] addresses a specific failure mode of LLM-based code translation: when the translator is asked to repair a failing test, it can "cheat" by silently weakening the test rather than fixing the code — deleting assertions, widening tolerances, removing negative-path cases, relaxing catch blocks. The paper introduces an adversarial test-generation technique that, on each repair iteration, synthesises a fresh batch of strong, discriminating tests (inputs that exercise boundary conditions, off-by-one cases, type-coercion edges, and exception paths) and injects them into the test suite before the validator checks the translator's repair. If the translator's repair drops any of these adversarial tests, the weakening is detected by an AST-level diff and the repair is rejected. Empirically, the technique reduces the rate of silent weakening by ~70 % on a C-to-Rust benchmark compared to a baseline validator without adversarial regeneration. The paper also reports that adversarial test generation catches ~3x more semantic regressions in translated code than random fuzzing alone, by focusing on the test shapes LLMs are most likely to over-fit.

## Relevance to MAgHARCM

AdvTestGen is the direct citation for MAgHARCM's [[PRIM-13]] (Adversarial Test-Weakening Guard). `internal/agents/validator.go` implements the cascade that includes the weakening guard: after the translator's repair, the validator diffs the current test AST against the source test AST and any dropped assertion triggers `AdversarialWeakeningDetected`, halting the repair loop. AdvTestGen's adversarial test regeneration (a fresh batch per iteration) is the natural extension to MAgHARCM's [[PRIM-5]] (Test Suite Co-Translation & Synthesis): after each repair, the validator asks a strong-test generator to synthesise boundary cases the current suite is missing, then re-runs the translator. The AST-validate-not-just-text-validate insight from AdvTestGen directly informs the [[PRIM-13]] implementation choice — a textual diff would miss semantic weakening (e.g., changing `assertEqual(a, 0)` to `assertGreaterEqual(a, 0)` is semantic weakening despite being a one-character edit). AdvTestGen's empirical finding that LLM translators over-fit certain test shapes also seeds the prompts in [[PRIM-25]] (Communicative-De-hallucination Role-Flip Gate) — the role-flip reviewer is told to look for exactly these over-fit patterns.

## Hop-1 References

- [[ReCodeAgent-2026]] — multi-agent translation pipeline; AdvTestGen's adversarial regeneration slots into ReCodeAgent's repair loop as an extra validator step.
- [[AlphaTrans-2024]] — repository-level translation with testCheck loop; AdvTestGen's adversarial regeneration strengthens testCheck's input space.
- [[Pynguin-2022]] — automated test-generation framework; AdvTestGen's adversarial synthesis uses Pynguin-style search-based techniques with adversarial objective functions.
- [[McKeeman-1998-Differential-Testing]] — older oracle-free testing technique; AdvTestGen's adversarial regeneration is in part a differential-testing successor applied to the translator-vs-source pair.
- [[CodaMOSA-2023]] — coverage-plateau detector that combines SBST with LLM-suggested tests; AdvTestGen's adversarial regeneration can be layered on top of CodaMOSA's coverage loop.

## Hop-2 Anchors (software-archaeology lean)

- [[Feathers-2004-WELC]] — characterisation tests are the offline analogue of AdvTestGen's adversarial tests: both isolate the cases the system (or the translator) cannot yet handle correctly.
- [[Kazman-Cai-2024]] — architectural recovery motivates AdvTestGen's boundary-case prioritisation: the most-load-bearing concepts (L1) deserve the strongest adversarial tests.
- [[Mueller-2000]] — the 5-strategy LIS analysis predicts that adversarial tests should be re-generated after each strategy switch in a multi-strategy migration, because each strategy exposes different failure modes.
