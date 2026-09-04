---
title: "CAID — Effective Strategies for Asynchronous Software Engineering Agents"
bibkey: p54_caid
tags: [paper, multi-agent, asynchronous, isolated-workspaces, git-worktree, [PRIM-17], hop-1]
---

# CAID — Effective Strategies for Asynchronous Software Engineering Agents

**Authors**: Jiayi Geng, Graham Neubig (Carnegie Mellon University)
**Year**: 2026 (arXiv v1: 2026-03)
**Venue**: arXiv preprint cs.CL / cs.AI (arXiv:2603.21489); OpenHands companion blog post dated 2026-03
**eprint / DOI**: arXiv:2603.21489 — https://arxiv.org/abs/2603.21489 ; code at https://github.com/JiayiGeng/CAID
**Cited by**: [[primitives/INDEX]] entry [[PRIM-17]] (Asynchronous Software Engineering Agent Architecture); references multiple agents on a shared-memory blackboard, event-driven wake, central scheduler for conflict arbitration

## Summary

[[Geng-2026-CSID]] (arXiv:2603.21489, CMU) introduces **CAID** (Centralized Asynchronous Isolated Delegation), a multi-agent framework for long-horizon software engineering tasks where existing approaches stall. The authors show that scaling a single agent's iteration budget yields only marginal gains and can sometimes *hurt* performance, and that "soft isolation" multi-agent setups where several agents share one workspace produce conflicting edits that are difficult to debug. CAID replaces shared workspaces with three primitives: a **centralized task-delegation manager** that decomposes the overall goal into subtasks, builds the dependency graph, and assigns work; **asynchronous execution** so subtasks run in parallel without blocking; and **isolated workspaces** (per-agent `git worktree`) so coordination happens through branch/commit/merge rather than shared, mutable memory. Empirically, CAID outperforms single-agent baselines by ~26.7 % on PaperBench and ~14.3 % on Commit0; the branch-and-merge workflow is identified as the load-bearing coordination mechanism because it forces integration to be verified by tests, not by trust.

## Relevance to MAgHARCM

[[PRIM-17]] cites CAID as the source for the "shared-memory blackboard + event-driven wake + central scheduler" architecture. MAgHARCM's planned `internal/agents/blackboard.go` (spec-only this sprint) should adopt CAID's *centralised manager + isolated workspaces* rather than the literal "shared mutable memory" interpretation in the appendix text — CAID's empirical result is that shared mutable memory *fails*, and that isolated `git worktree` workspaces with a merge-time integration test are what actually scales. The delegation manager maps cleanly onto `internal/agents/planning.go` ([[PRIM-1]]/[[PRIM-2]]), which already builds a DAG over translation fragments: a CAID-style recruitment step per fragment dispatches one agent per `git worktree`, runs the chunked translator ([[PRIM-23]]) inside, and merges only when the validator ([[PRIM-6]]/[[PRIM-13]]) green-lights the branch. This converts MAgHARCM's current sequential pipeline into a horizontally scalable one without abandoning the multi-stage build/test feedback loop.

## Hop-1 References

- [[Hong-2023-MetaGPT]] — assembly-line SOP encoding; CAID's centralised manager is the "Product Manager" role elevated to a top-level scheduler with isolated workspaces.
- [[Wu-2023-AutoGen]] — conversable-agent substrate and interrupt/resume ([[PRIM-28]]); CAID inherits AutoGen's per-agent autonomy but adds hard workspace isolation.
- [[Chen-2023-AgentVerse]] — recruitment-adaptive planning ([[PRIM-30]]); CAID's per-subtask agent assignment is recruitment at subtask granularity rather than per-iteration granularity.
- [[Qian-2023-ChatDev]] — chat-chain role-flipping ([[PRIM-25]]); CAID's merge-time integration test is an asynchronous analogue of ChatDev's review step.
- [[Yang-2024-SWE-Agent]] — agent-computer interface; CAID's `git worktree` is the SWE-agent sandboxing primitive generalised to long-horizon multi-agent settings.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000-DesignRules]] — CAID's *isolated workspaces + merge integration test* is a load-bearing design rule (the merge contract is the L1 interface; each worktree is a substitutable L3 leaf); the architecture explicitly maximises the option value of "swapping one agent's work for another without breaking the integration rule".
- [[Rajlich-1997-ICSE]] — the central delegation manager needs Rajlich's concept-locator analysis to assign each subtask to the agent whose concept ownership matches; without the concept map, the manager degenerates into blind round-robin.
- [[Kazman-Cai-2024-ArchitecturalRecovery]] — the dependency graph the CAID manager builds is a recovered architectural artefact; for legacy targets it must be reconstructed by [[PRIM-9]]'s tri-representation code graph before delegation can be safe.
