---
title: "P-139 TypePro: Boosting LLM-Based Type Inference via Inter-Procedural Slicing"
backlink: "[[1.0.0 P-139]]"
aliases:
  - "1.0.0 P-139"
  - "P-139"
  - "P-139-TypePro-FSE-2026"
  - "TypePro-FSE-2026"
  - "lin2026typepro"
tags: [paper, type-inference, system-dependence-graph, inter-procedural-slicing, "[[1.0.0 PRIM-9]]", "[[1.0.0 PRIM-22]]", wave-19]
date: 2026-09-07
last_updated: 2026-09-07
venue: FSE 2026 (Research Papers, PACMSE inaugural issue)
---

# [[1.0.0 P-139]] TypePro

## TL;DR

TypePro combines **System Dependence Graph (SDG)** construction (function-name, parameter, and return-value matching) with **inter-procedural backward/forward code slicing** and a structural-similarity candidate-type selector. Achieves **88.9% Top-1 EM on ManyTypes4Py (+7.1 pp)** and **86.6% on ManyTypes4TypeScript (+10.3 pp)**. LLM slicing substrate with no LLM-only prompt tweaks. Anchors `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph (slicing as a fourth representation) and `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Recognition + Explanation substrate).

## Mechanism (Q2)

1. **System Dependence Graph construction** — extends intra-procedural PDG with inter-procedural edges via function-name, parameter, and return-value matching across call sites.
2. **Inter-procedural backward/forward slicing** — Algorithm 1 of the paper; computes a slice across procedures that captures the user-defined type structure.
3. **Structural-similarity candidate-type selector** — ranks candidate types by structural similarity to the slice; LLM is queried only on the top-k candidates.

Together these give MAgHARCM's migration pipeline a **type-inference substrate** that works across procedures — the existing PDG-based analysis in `internal/agents/cpg.go` is intra-procedural and does not propagate user-defined types across function boundaries.

## Anchoring (Q3)

| Primitive | Pre-wave-19 behaviour | TypePro substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | CPG had syntactic AST edges and dynamic call-graph reality but no inter-procedural type-aware slicing | SDG with type-aware matching adds a fourth representation that bridges syntactic and dynamic edges |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Recognition and Explanation phases operated on local context; cross-procedure type recovery was noisy | Slice-based type inference gives the phases a substrate for inter-procedural type recovery |

## Hop-1 Citations

- Peng 2023 TypeGen (type inference via CoT).
- P-88 HiTyper ICSE 2022 (Type Dependency Graph).
- P-129 LλMDA ICSE 2026 (LLM-aided partial program dependence).

## Hop-2 Citations

- Horwitz-Reps-Binkley 1990 System Dependence Graph (inter-procedural slicing foundation).
- P-28 Weiser 1979 Program Slicing (slicing origin).
- P-72 Yamaguchi 2014 Code Property Graph (CPG substrate).

## MAgHARCM integration

- **YAML config key**: `static_analysis.typepro.enabled: true`; `static_analysis.typepro.slice_max_depth: <int>`.
- **Implementation file**: `internal/staticanalysis/typepro.go::InterProceduralSlice` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-9]]`, `[[1.0.0 PRIM-22]]`.

## Caveats

- **Language coverage** — evaluated on Python and TypeScript; Rust/C/Go adapters require additional matching rules.
- **Slice cost** — inter-procedural slicing is more expensive than intra-procedural; budget for large codebases.
- **Candidate top-k** — the structural-similarity selector requires a knowledge base of candidate types; KB coverage limits recall.

## Source

- arXiv: 2604.02702.
- Venue: FSE 2026 Research Papers (PACMSE inaugural issue; ACM DOI forthcoming; verified via authors' homepage wu-rongxin1987.github.io and FSE '26 PDF copyright header at huaxunhuang.github.io/src/fse2026.pdf).
