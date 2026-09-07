---
title: MatchFixAgent — Language-Agnostic Autonomous Repository-Level Code Translation Validation and Repair
backlink: "[[1.0.0 P-08]]"
bibkey: p08_matchfixagent
aliases:
  - "1.0.0 P-08"
  - "P-08"
  - "P-08-MatchFixAgent"
  - "P-08-MatchFixAgent"
  - "MatchFixAgent"
tags: [paper, repository-translation, multi-agent, validation, repair, [[PRIM-7]], [[PRIM-13]], hop-1]
---

# [[1.0.0 P-08]] MatchFixAgent — Language-Agnostic Autonomous Repository-Level Code Translation Validation and Repair

**Authors**: Ali Reza Ibrahimzada, Brandon Paulsen, Reyhaneh Jabbarvand, Joey Dodds, Daniel Kroening
**Year**: 2025 (arXiv preprint); 2026 (ICML acceptance announced 2026-04-30)
**Venue**: 43rd International Conference on Machine Learning (ICML 2026); arXiv:2509.16187 (Sep 2025)
**eprint / DOI**: arXiv:2509.16187
**Cited by**: [[primitives/Primitives-Index]] entry [[PRIM-7]] (Multi-Agent Verdict Validation); referenced indirectly by [[PRIM-13]] (Adversarial Test-Weakening Guard)

## Summary

[[Ibrahimzada-2025-MatchFixAgent]] is a language-agnostic multi-agent framework for autonomous validation and repair of repository-level code translations. It addresses the test-suite insufficiency problem that [[Oxidizer-2024]] and [[AlphaTrans-2025]] work around by introducing four specialised agents: (a) Semantic-Analysis agents that compare control/data-flow paths and API-usage patterns between source and target; (b) a Test-Generation agent that synthesises new tests from the source's commit history and existing test patterns; (c) a Repair agent that, given a verdict of inequivalence, localises the fault and proposes a minimal patch; and (d) a Verdict agent that synthesises the other three outputs into a final (in)equivalence decision. On the [[Ou-2024-RustRepoTrans]] benchmark and [[TransRepo-Bench-2025]], MatchFixAgent issues verdicts for 99.2 % of translation pairs (vs. 60–80 % for test-only baselines), matches prior-art verdicts in 72.8 % of cases, and is judged correct in 60.7 % of disagreements. The Repair agent successfully fixes 50.6 % of inequivalent translations (vs. 18.5 % for the strongest baseline). Critically, the framework operates without human-in-the-loop intervention; its adversarial structure (multiple independent judges) is the direct ancestor of MAgHARCM's [[PRIM-7]] verdict panel.

## Relevance to MAgHARCM

MatchFixAgent is the direct citation for [[PRIM-7]] (Multi-Agent Verdict Validation). MAgHARCM's `internal/agents/verdict_panel.go` (planned this sprint per the methodology ledger) implements the multi-judge disagreement pattern from MatchFixAgent §3 — multiple LLM verdict agents score the translation independently, and disagreement triggers a repair re-prompt. MatchFixAgent's Test-Generation agent motivates [[PRIM-13]] (Adversarial Test-Weakening Guard): the test-generation step can silently weaken tests to make them pass, which the Repair agent then exploits, so MAgHARCM's validator must AST-validate that any generated test maintains its assertion strength. The paper's RustRepoTrans evaluation aligns with MAgHARCM's Go→Rust focus, making its empirical claims directly transferable. The 50.6 % repair rate is the benchmark MAgHARCM aims to beat on Sample 4 (commons-validator).

## Hop-1 References

- [[Ibrahimzada-2025-AlphaTrans]] — same group; MatchFixAgent's test-generation agent extends AlphaTrans's test translation to *new* test synthesis.
- [[Ibrahimzada-2026-ReCodeAgent]] — same group, language-agnostic translation; MatchFixAgent adds the verification + repair loop on top.
- [[Ke-2025-TRAM]] — mock-based in-isolation validation; MatchFixAgent is the multi-agent counterpart operating on the *whole* translation.
- [[Oxidizer-2024]] — same neuro-symbolic family; MatchFixAgent's Repair agent borrows its type-driven analysis for localisation.
- [[Ou-2024-RustRepoTrans]] — evaluation benchmark; the RustRepoTrans tasks are the primary test bed for MatchFixAgent's claims.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000]] — the multi-judge disagreement pattern is a load-bearing design rule: when load-bearing rules conflict (multiple judges disagree), the system must revert to a higher-level authority (the repair agent).
- [[Rajlich-1997]] — concept-locator analysis motivates MatchFixAgent's semantic-analysis agent: semantic equivalence is judged on whether the target preserves the source's concept set, not its surface behaviour.
- [[Müller-2000]] — legacy integration strategy: MatchFixAgent's adversarial structure maps to Chicken Little (gradual validation) rather than Cold Turkey (full migration), matching MAgHARCM's strategy selector.
- [[Foltz-2023]] — DR.JONES explains why multi-agent disagreement beats single-agent self-consistency: independent linear traversals empirically outperform a single deep one for software comprehension tasks.

## Hop-2 Deep Archaeology & Upstream Lineage

- **Software Modernization Archaeology (Rajlich & Müller)**: Connects modern LLM decompilation/migration back to early software reverse engineering (program slicing, concept assignment, redocumentation).
- **Cognitive Traversal (Foltz & Landauer / Latent Semantic Analysis)**: How human engineers comprehend legacy systems across semantic hops versus how LLM context windows navigate fragmented symbols.
- **Architectural Coupling & Decomposition (Baldwin & Clark / Kazman)**: Modularity theory and Design Structure Matrices (DSM) underpinning why reverse-topological scheduling ([[PRIM-1]]) avoids cyclic cascade failures.
