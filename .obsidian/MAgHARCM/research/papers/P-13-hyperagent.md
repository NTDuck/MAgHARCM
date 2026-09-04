---
title: HyperAgent — Generalist Software Engineering Agents to Solve Coding Tasks at Scale
bibkey: p13_hyperagent
tags: [paper, software-engineering-agent, generalist, navigator, [[PRIM-26]], hop-1]
---

# HyperAgent — Generalist Software Engineering Agents to Solve Coding Tasks at Scale

**Authors**: Huy Nhat Phan, Phong X. Nguyen, Nghi D. Q. Bui, Tien N. Nguyen (FPT Software AI Center; University of Texas at Dallas)
**Year**: 2024
**Venue**: arXiv preprint (Sept 2024; no conference acceptance surfaced as of 2026-08-30)
**eprint / DOI**: arXiv:2409.16299
**Cited by**: [[primitives/INDEX]] entry for [[PRIM-26]] (Symbol-Aware Navigator) and the role-specialisation model behind the planner/navigator/editor/executor decomposition referenced in the appendix.

## Summary

[[Phan-2024-HyperAgent]] is a multi-agent system that emulates the end-to-end workflow of a human developer over a repository: Plan → Navigate → Edit → Execute. Four specialised agent roles cooperate: a Planner coordinates strategy and assigns sub-tasks, a Navigator explores the repository to locate relevant code, a Code Editor emits patches, and an Executor runs tests to verify. HyperAgent is explicitly positioned as "generalist" rather than issue-resolver-only: it tackles SWE-bench, RepoExec (repository-level code generation), Defects4J (fault localisation and repair) under a single architecture. Reported numbers are state of the art at publication — 25.01% on SWE-bench-Lite and 31.40% on SWE-bench-Verified. The Navigator exposes code-aware tools (SearchSymbols, GetDefinition, GetReferences, SearchRepoGraph) that return tight, ranked context windows — the immediate ancestor of MAgHARCM's [[PRIM-26]]. Code: `FSoft-AI4Code/HyperAgent`.

## Relevance to MAgHARCM

HyperAgent is the direct citation for [[PRIM-26]] (Symbol-Aware Navigator). MAgHARCM's `internal/agents/navigator.go` realises the same tool surface (`SearchSymbols`, `GetDefinition`, `GetReferences`, `SearchRepoGraph`) and the same ≤4 KB tight-context-window discipline before each fragment translation. The Planner / Navigator / Editor / Executor role quartet matches MAgHARCM's pipeline: Planner corresponds to `internal/agents/planning.go`, Navigator to `internal/agents/navigator.go`, Editor to the chunked translator ([[PRIM-23]]), and Executor to the validator ([[PRIM-6]]) running rustc + translated tests. HyperAgent's generalist framing — single architecture across issue-resolution, code-generation, fault-localisation, repair — is the methodological template MAgHARCM adopts for covering both Go→Rust translation and legacy-code archaeology.

## Hop-1 References

- [[Yang-2024-SWE-Agent]] — earlier generalist SWE-agent; introduces the agent-computer-interface (ACI) concept that HyperAgent extends with a multi-agent decomposition.
- [[Jimenez-2024-SWE-Bench]] — the benchmark HyperAgent is evaluated on; defines the issue-resolution task that motivates generalist SE agents.
- [[Yang-2024-RepoCoder]] — iterative retrieval-augmented repo-level code completion; precursor of MAgHARCM's [[PRIM-31]] (Iterative Retrieval Refinement) and HyperAgent's Navigator re-indexing loop.
- [[CloudWeGo-2025-ABCoder]] — AST-based code-RAG via MCP; an alternative to HyperAgent's navigator with stronger structural grounding but no generalist agent layer.
- [[Phan-2024-HyperAgent-Appendix]] — supplementary material with the exact tool-call protocols and ranker weights the public implementation uses.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997]] — concept-locator analysis is what the Navigator does at code level: locate the symbols, definitions, and references that implement a given concept; HyperAgent's tools are operationalisations of the concept→location map.
- [[Kazman-Cai-2024]] — architectural recovery motivates why HyperAgent needs a SearchRepoGraph term (not just textual search): symbol-overloading and cross-module call chains are invisible to lexical retrieval, exactly the gap architectural recovery fills.
- [[AgentPatterns-2024-Legacy-Code-Archaeology]] — the "symbol-aware navigator" pattern is the live-system analogue of the archaeological walkthrough; same discipline (named entry points, defined context windows), different substrate (running repo vs. legacy codebase).
