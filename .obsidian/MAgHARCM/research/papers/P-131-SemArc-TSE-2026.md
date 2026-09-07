---
title: "P-131 SemArc: Software Architecture Recovery Augmented with Semantics"
backlink: "[[1.0.0 P-131]]"
aliases:
  - "1.0.0 P-131"
  - "P-131"
  - "P-131-SemArc-TSE-2026"
  - "SemArc-TSE-2026"
  - "zhao2026semarc"
tags: [paper, architecture-recovery, semantic-comprehension, llm, "[[1.0.0 PRIM-9]]", "[[1.0.0 PRIM-22]]", wave-18]
date: 2026-09-07
last_updated: 2026-09-07
venue: IEEE TSE Vol. 52 Issue 1 (Jan 2026)
---

# [[1.0.0 P-131]] SemArc

## TL;DR

SemArc introduces three sub-mechanisms for semantic-aware architecture recovery: (a) **LLM-powered semantic comprehension** that summarizes implementation-level intent, anchored by a knowledge base of canonical architectural patterns; (b) **integration of implicit + explicit dependencies** to recover semantics missed by structural-only tools; (c) **Component-as-Anchor Guided Clustering** where domain-specific or LLM-generated anchors guide the partition. Achieves **+32pp accuracy** over seven baselines on 15 C/C++/Java/Python systems. Anchors `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph and `[[1.0.0 PRIM-22]]` Four Phases of Comprehension.

## Mechanism (Q2)

Three coupled sub-mechanisms:

1. **LLM-powered semantic comprehension** — for each file/module, the LLM produces an intent summary that is mapped against a knowledge base of canonical architectural patterns (layered, MVC, microservices, etc.). The summary enables downstream stages to reason about the role of each component rather than just its dependencies.
2. **Implicit + explicit dependency fusion** — structural tools capture explicit dependencies (calls, imports). SemArc adds implicit dependencies inferred from naming, co-change history, or semantic co-occurrence. The fused graph captures coupling that structural-only tools miss.
3. **Component-as-Anchor Guided Clustering** — instead of pure bottom-up community detection, SemArc uses anchors (either domain-specific known components or LLM-generated candidate components) to seed the clustering process. Anchors constrain the partition toward architecturally plausible boundaries.

The combination addresses the two failure modes of pure structural recovery: missing implicit coupling (mechanism 2) and missing semantic intent (mechanism 1), while the anchor-guided clustering (mechanism 3) ensures the final partition is anchored in known architectural concepts.

## Anchoring (Q3)

| Primitive | Pre-wave-18 behaviour | SemArc substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Structural-only dependency graphs; pure bottom-up clustering | Intent-augmented tri-representation: structural + implicit + semantic-intent layers; anchor-guided clustering |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Comprehension halts at structural signals; semantic intent is implicit | Explicit semantic-intent phase powered by LLM summaries + architectural pattern knowledge base |

## Hop-1 Citations

- Architectural pattern knowledge bases (Buschmann et al. 1996; Gamma et al. 1994).
- Implicit-dependency inference (Pradel et al. 2012; Pal et al. 2019).
- Anchor-guided clustering lineage.

## Hop-2 Citations

- LLM-for-code lineage (Ma et al. 2024; Wang et al. 2024; Dong et al. 2024).
- Maqbool et al. 2024 / Sahu et al. 2023 — earlier semantic-architecture-recovery work that SemArc generalizes.

## MAgHARCM integration

- **YAML config key**: `architecture_recovery.semarc.enabled: true` (opt-in, default `false`); `architecture_recovery.semarc.anchor_mode: <known|llm|hybrid>`.
- **Implementation file**: `internal/code_graph/semarc.go::RecoverWithSemantics` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-9]]`, `[[1.0.0 PRIM-22]]`.

## Caveats

- **Knowledge-base maintenance**: the canonical pattern KB must be kept current; pattern drift between KB and current practice degrades anchoring quality.
- **LLM intent-summary cost**: every file incurs an LLM call; very large repos need batching and caching to keep wall-clock tractable.
- **Anchor ambiguity**: domain-specific anchors may not exist for novel architectures; the LLM-anchor fallback must be reliable or it will mis-seed the partition.

## Source

- URL: computer.org/csdl/journal/ts/2026/01/11203255/2aOjGlUCArS.
