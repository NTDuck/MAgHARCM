---
title: CodePlan — Repository-level Coding using LLMs and Planning
bibkey: p03_codeplan
tags: [paper, repository-coding, planning, [[PRIM-1]], [[PRIM-23]], [[PRIM-29]], hop-1]
---

# CodePlan — Repository-level Coding using LLMs and Planning

**Authors**: Ramakrishna Bairi, Atharv Sonwane, Aditya Kanade, Vageesh D C, Arun Iyer, Suresh Parthasarathy, Sriram Rajamani, B. Ashok, Shashank Shet
**Year**: 2024 (arXiv 2023-09-21; PACMSE FSE 2024)
**Venue**: Proceedings of the ACM on Software Engineering (PACMSE) / FSE (DOI 10.1145/3643757); arXiv:2309.12499
**eprint / DOI**: arXiv:2309.12499
**Cited by**: [[primitives/INDEX]] entry not formally assigned, but the methodology ledger reuses CodePlan's planning backbone for [[PRIM-1]] (Reverse Topological Translation Ordering), [[PRIM-23]] (Chunked Translation), [[PRIM-29]] (Recruitment-Adaptive Planning).

## Summary

[[Bairi-2024-CodePlan]] formulates repository-level coding — package migration, error-repair, type annotation — as a planning problem and presents a task-agnostic neuro-symbolic framework. The planner combines an incremental dependency analysis (what files does the change reach?), a change may-impact analysis (what transitive effects arise?), and an adaptive planning algorithm that emits a multi-step chain of edits; each step is one LLM call on a code location with context derived from the entire repository, prior changes, and task instructions. Evaluated on C# package migration and Python temporal code edits across repositories requiring interdependent edits in 2–97 files, CodePlan passed 5/6 repositories' validity checks where baseline LLM-with-context configurations passed none.

## Relevance to MAgHARCM

CodePlan is the planning lineage that MAgHARCM inherits alongside [[AlphaTrans-2024]]. The incremental-dependency + may-impact pair is what `internal/agents/planning.go::ComputeReverseTopoOrder` ([PRIM-1]) produces when fragmenting the AST before the chunked translator ([PRIM-23]) is invoked. The adaptive-planning idea — re-plan after each iteration's report — reappears in [[PRIM-29]]'s recruiter loop, which picks per-iteration tool/agent configurations based on the validator report. CodePlan also prefigures [[PRIM-30]] (Iterative Retrieval Refinement): the planner re-queries context after each plan step, just as RepoCoder re-indexes after each fragment.

## Hop-1 References

- [[RepoCoder-2023]] — iterative retrieval-augmented code completion; CodePlan cites it for the re-query loop.
- [[AlphaTrans-2024]] — applies a similar incremental-dependency view at translation granularity (per-method) rather than per-edit.
- [[Ling-2024-TransRepo-Bench]] — the skeleton-guided translation benchmark is the Java→C# application of CodePlan-style skeleton planning.
- [[Roziere-2020-Transcoder]] — unsupervised neural baseline CodePlan compares against on translation sub-tasks.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997]] — concept-locator analysis is the offline analogue of CodePlan's dependency-graph "context slices": both group files by what they implement, not by what they import.
- [[Baldwin-Clark-2000]] — design-rule hierarchy explains why CodePlan plans *inter*faces before *implementations*: skeleton (L1) is stable, leaves (L3) are substitutable.
- [[Mueller-LIS-2004]] — the five-strategy migration taxonomy is the catalogue that CodePlan's "task-agnostic" planning implicitly chooses between per repository profile.
