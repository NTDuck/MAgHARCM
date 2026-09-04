---
title: Legacy Code Archaeology Pattern for Autonomous Systems
bibkey: p37_agentpatterns
tags: [esoteric, pattern, software-archaeology, multi-agent, [[1.0.0 PRIM-14]], [[1.0.0 PRIM-15]], hop-1]
---

# [[1.0.0 P-37]] Legacy Code Archaeology Pattern for Autonomous Systems

**Authors**: AgentPatterns.ai Architecture Group  
**Year**: 2024  
**Venue**: AgentPatterns.ai / Architectural Catalog for Agentic Engineering  
**eprint / DOI**: agentpatterns.ai/patterns/legacy-code-archaeology  
**Cited by**: `[[2.0.0 Software-Archaeology-Lineage]]`, `[[1.0.0 PRIM-14]]` (Software Archaeology Stage), `[[1.0.0 PRIM-15]]` (Evidence-First Adaptation Pattern)

## Summary

The AgentPatterns.ai guide formalizes how autonomous LLM multi-agent systems should conduct legacy code excavation before attempting programmatic synthesis or refactoring.
The pattern identifies three fatal failure modes when LLMs attempt naive legacy modernization:
1. *Context Window Exhaustion*: Shoveling sprawling legacy files directly into LLM prompts loses subtle invariant constraints.
2. *Hallucinated Modernity*: LLMs assume modern framework conventions that subtly break legacy behavioral quirks.
3. *Blind Transliteration*: Synthesizing target code without understanding temporal dependencies causes ripple-effect compilation failures.

The pattern outlines the **Archaeologist-to-Planner Pipeline**:
- *Phase 1: Excavation Agent*: Runs deterministic AST, git log, and naming tools to generate an archaeology report.
- *Phase 2: Evidence-First Adaptation*: Synthesizes target-preserving behavioral specs and untrusted source models (`[[1.0.0 PRIM-15]]`).
- *Phase 3: Bounded Synthesis*: Passes only bounded context excerpts to specialized coding agents.

## Relevance to MAgHARCM

Direct inspiration for MAgHARCM's multi-agent decomposition:
- The Archaeologist Agent runs as Node 1 in the Eino Graph (`internal/graph/graph.go`), decoupled from the Planner and Translator.
- `[[1.0.0 PRIM-15]]` (`internal/agents/evidence_adaptation.go`) enforces the Evidence-First pattern before translation begins.
- Keeps local SLMs (30B reasoning + 4B coding) within their effective attention budgets ($\le 4$\,KB bounded retrieval).

## Hop-1 References

- [[pp-besm-Software-Archaeology]] — Upstream source for terminal excavation commands.
- [[Reeper-2024]] — Scientific formulation of the Evidence-First adaptation pattern.
- [[MetaGPT-2023]] — Standard Operating Procedure (SOP) role-artifact contracts.

## Hop-2 Anchors (Software Archaeology Lean)

- [[Foltz-2023]] — Grounds the excavation steps in human cognitive traversal patterns.
- [[Baldwin-Clark-2000]] — Relies on design rule boundaries to scope LLM attention budgets.
