---
title: Working Effectively with Legacy Code
bibkey: p24_feathers
tags: [paper, book, legacy-code, seams, characterisation-tests, [[PRIM-14]], [[PRIM-5]], hop-1]
---

# [[1.0.0 P-24]] Working Effectively with Legacy Code

**Authors**: Michael C. Feathers
**Year**: 2004
**Venue**: Prentice Hall (Pearson), Object Technology Series, ed. Martin Fowler
**eprint / DOI**: ISBN-10: 0131177052 (print); no DOI
**Cited by**: [[primitives/INDEX]] entry [[PRIM-14]] (Software-Archaeology Stage), [[PRIM-5]] (Test Suite Co-Translation & Synthesis)

## Summary

[[Feathers-2004-WELC]] is the canonical practitioner reference for bringing untested ("legacy") code under test. Feathers defines legacy code operationally: code without tests. The book catalogues ~25 dependency-breaking techniques — collectively called "seams" — that allow a developer to inject tests into a code base whose original design resists testability. Each seam is a place in the code where the behaviour of a module can be altered without editing the source: pre-processor seams, link seams, object-seams (subclassing), and method-object seams are the four families. The book then maps the seams to refactoring moves (Sprout Method, Sprout Class, Wrap Method, Wrap Class, Extract Interface, Parameterise Constructor, Pull-Up Constructor, Push-Down Dependency) and provides decision flowcharts for choosing among them. The third pillar is the "characterisation test" — a test whose expected output is not specified by the developer but recorded from the current behaviour; characterisation tests become the regression safety net for subsequent refactoring. The book also introduces "churn" as a measurable indicator of where the legacy code is being modified and thus most likely to need seams.

## Relevance to MAgHARCM

Feathers is the philosophical ancestor of MAgHARCM's entire pipeline, but especially [[PRIM-14]] (Software-Archaeology Stage). The five sub-steps of the archaeologist (`internal/agents/archaeology.go`, planned) — `ExtractBoundaries`, `BuildTimeCapsule`, `FindChurnHotspots`, `ForensicNaming`, `MapConcepts` — are direct adaptations of Feathers's workflow: extracting boundaries mirrors seam-identification, building the time-capsule mirrors the build-environment-reconstruction step, churn hotspots mirror his churn-driven prioritisation, forensic naming mirrors his "explaining cryptic identifiers" guidance, and concept-mapping mirrors his advice to "sketch what the code is for". Feathers's characterisation-test technique is the offline analogue of [[PRIM-5]] (Test Suite Co-Translation & Synthesis) when source tests are missing: record the current behaviour, freeze it as a test, then proceed. The seam catalogue informs [[PRIM-19]] (Design Rule Hierarchy Partitioning) by teaching the harness to recognise L1/L2 boundaries (the seams) versus L3 leaves (the substitutable modules). Feathers's "rules for editing legacy code" also seed the safety guards in [[PRIM-6]] (Multi-Stage Build/Test Feedback Repair) — never change behaviour and tests in the same commit; prefer small, reversible refactors.

## Hop-1 References

- [[Fowler-1999-Refactoring]] — Fowler's earlier refactoring catalogue that Feathers extends with the legacy-specific "Sprout / Wrap" techniques.
- [[Rajlich-1997]] — concept-locator analysis is the theoretical foundation for Feathers's "sketch what the code is for" guidance; both treat the source as a network of concepts.
- [[Mueller-2000]] — the 5-strategy LIS analysis is the migration-framework complement to Feathers's tactics; Feathers answers "how" and Mueller answers "which migration strategy".
- [[Kazman-Cai-2024]] — architectural recovery extends Feathers's churn hotspots from file-level to module-level using static-dependency analysis.
- [[OpenRewrite-2023]] — modern tool that automates Feathers's "Sprout Method / Wrap Class" refactors at scale; informs [[PRIM-10]] feature-mapping validation.

## Hop-2 Anchors (software-archaeology lean)

- [[pp-besm-Software-Archaeology]] — the dev.to playbook operationalises Feathers's seam catalogue for non-expert developers and gives concrete shell commands for finding seams with grep/ctags.
- [[AgentPatterns-Legacy-Code-Archaeology]] — agent-pattern essay that adapts Feathers's workflow for an LLM-driven archaeologist agent; informs MAgHARCM's [[PRIM-14]] sub-step prompt design.
- [[Baldwin-Clark-2000]] — design rules formalise what Feathers calls "seams": stable inter-module contracts (L1) are the natural seams; substitutable leaves (L3) are where sprouting is safe.
