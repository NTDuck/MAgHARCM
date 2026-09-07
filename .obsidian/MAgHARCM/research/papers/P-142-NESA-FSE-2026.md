---
title: "P-142 NESA: Relational Neuro-Symbolic Static Program Analysis"
backlink: "[[1.0.0 P-142]]"
aliases:
  - "1.0.0 P-142"
  - "P-142"
  - "P-142-NESA-FSE-2026"
  - "NESA-FSE-2026"
  - "nesa2026fse"
tags: [paper, neuro-symbolic, static-analysis, datalog, comprehension, "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-9]]", wave-20]
date: 2026-09-07
last_updated: 2026-09-07
venue: FSE 2026 (Research Track)
---

# [[1.0.0 P-142]] NESA

## TL;DR

NESA is a **relational neuro-symbolic** static program-analysis framework that uses a restricted **Datalog** analysis-policy language to decompose complex analysis problems into syntactic sub-problems (handled by parsing-based analysis) and semantic sub-problems (handled by LLM). Lazy and incremental prompting reduces LLM hallucinations. Compilation-free and customizable. Achieves **F1 0.72** on TaintBench (+0.20 over industrial baseline); detects **13 real-world memory-leak bugs** subsequently fixed by developers. Anchors `[[1.0.0 PRIM-22]]` (Four Phases of Comprehension: Reorganization phase via Datalog policy decomposition) and `[[1.0.0 PRIM-9]]` (Tri-Representation Hybrid Code Graph: a *fourth representation* — neuro-symbolic policy — on top of syntactic / lexical / semantic).

## Mechanism (Q2)

1. **Restricted Datalog policy language** — users express an analysis policy as a set of Datalog rules; complex problems decompose into smaller sub-problems.
2. **Sub-problem routing** — syntactic sub-problems are resolved by parsing-based static analysis (avoiding LLM hallucinations); semantic sub-problems are routed to the LLM.
3. **Lazy and incremental prompting** — the analysis policy is evaluated using a strategy that only invokes the LLM when parsing alone is insufficient, and only for the additional predicates needed.
4. **Compilation-free operation** — NESA does not require the target program to compile; it operates on source-code syntax plus LLM reasoning.

Together these give MAgHARCM a *fourth representation* for code understanding beyond syntactic (AST), lexical (token), and semantic (CPG) — the **neuro-symbolic policy** that explicitly separates "what can be computed deterministically" from "what requires an LLM".

## Anchoring (Q3)

| Primitive | Pre-wave-20 behaviour | NESA substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Reorganization / Insight phases) | Reorganization relied on syntactic AST rewriting; Insight required LLM-only reasoning with no formal boundary | NESA's Datalog policy is the *decomposition layer* for Reorganization; LLM only invoked for the semantic sub-problem slice |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Three representations (syntactic AST, lexical tokens, semantic CPG) — all deterministic | NESA adds a *fourth* representation — neuro-symbolic Datalog policy — that captures relational facts spanning all three deterministic representations |

## Hop-1 Citations

- TaintBench (2019, Feng et al., benchmark substrate for the F1=0.72 result).
- P-139 TypePro (FSE 2026, inter-procedural SDG slicing) — sister paper at the same venue, both in the LLM-augmented static-analysis family.
- P-129 LλMDA (ICSE 2026, LLM-aided partial PDG) — Wave-18 anchor that pioneered the partial-program context-augment-then-DG pattern; NESA generalizes this to Datalog policy decomposition.
- Relational program analysis / Datalog lineage (Bravenboer & Smaragdakis 2009, Doop).

## Hop-2 Citations

- Datalog (classic, 1977).
- Abstract interpretation (Cousot & Cousot 1977) — foundations of static analysis.
- Neuro-symbolic AI (Garcez & Lamb 2020 survey).
- P-87 Hou et al. TOSEM 2024 (LLM4SE SLR — confirms the LLM-augmented static-analysis axis as an active research direction).

## MAgHARCM integration

- **YAML config key**: `comprehension.nesa.policy_file: configs/nesa_policy.dl`; `comprehension.nesa.llm_routing.max_calls: <int>`.
- **Implementation file**: `internal/comprehension/nesa_policy.go::NewNESAPolicy` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-9]]`.

## Caveats

- **Datalog expressiveness** — restricted Datalog cannot express all desired policies (e.g., path-sensitive analyses); users must decompose to fit.
- **LLM invocation cost** — even with lazy prompting, each semantic sub-problem incurs an LLM round-trip; large policies may be slow.
- **Compilation-free caveat** — analysis is incomplete by definition; properties that require type information from compilation must be approximated or skipped.

## Source

- arXiv: 2412.14399.
- Venue: FSE 2026 (Research Track), verified via PACMSE Vol 3 Issue FSE Article FSE154, DOI: 10.1145/3808161; conference-publishing.com/authors/FSE26 paper index 103.
