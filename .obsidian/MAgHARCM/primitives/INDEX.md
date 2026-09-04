---
title: Primitives Index
backlink: [[Primitives]]
tags: [primitives, catalog, status, [[PRIM-1]]..[[PRIM-30]]]
---

# Primitives Index

A primitive is a small, named capability the harness reuses across the
pipeline. Each entry links to the file that implements it and notes its
status.

[[Sprint-Recon-2026-09-04]] is the cross-reference of paper appendix
primitives against `internal/agents/*.go`. Status legend: `present`,
`partial`, `absent-implemented` (this sprint), `referenced` (paper cites
it; not a codebase primitive).

## Active — implemented in codebase

| ID | Title | Where | Status |
| -- | ----- | ----- | ------ |
| [[PRIM-1]] | Reverse Topological Translation Ordering | `internal/agents/planning.go::ComputeReverseTopoOrder` | present |
| [[PRIM-2]] | Back-Edge-Conditioned Reverse-Topological Scheduling | `internal/agents/planning.go` | present |
| [[PRIM-3]] | Target Skeleton-First Generation | `internal/agents/planning.go::DefaultProjectSkeleton` | present |
| [[PRIM-4]] | SpecMiner-Style Dynamic Invariant Recovery | `internal/agents/specminer.go` | absent-implemented |
| [[PRIM-5]] | Test Suite Co-Translation & Synthesis | `internal/agents/validator.go::generateAdditionalTests` | present |
| [[PRIM-6]] | Multi-Stage Build/Test Feedback Repair | `internal/agents/validator.go` | present |
| [[PRIM-7]] | Multi-Agent Verdict Validation | `internal/agents/verdict_panel.go` | absent-implemented |
| [[PRIM-8]] | State-Grounded Mock-Based In-Isolation Validation | `internal/agents/mock_validator.go` | absent-implemented |
| [[PRIM-9]] | Tri-Representation Hybrid Code Graph | `internal/agents/cpg.go` | absent-implemented |
| [[PRIM-10]] | Feature-Mapping & Type-Compatibility Validation | `internal/agents/feature_mapping.go` | absent-implemented |
| [[PRIM-11]] | Implementation-Agnostic Testing | `internal/agents/impl_agnostic.go` | absent-implemented |
| [[PRIM-12]] | Wasm-Based Reference Execution Oracle | `internal/agents/wasm_oracle.go` | absent-implemented |
| [[PRIM-13]] | Adversarial Test-Weakening Guard | `internal/agents/validator.go` | present |
| [[PRIM-14]] | Software-Archaeology Stage | `internal/agents/archaeology.go` | absent-implemented |
| [[PRIM-17]] | Asynchronous SE Agent Blackboard (CAID) | `internal/agents/blackboard.go` | absent-implemented (spec) |
| [[PRIM-23]] | Chunked Translation | `internal/agents/chunked_translator.go` | present |
| [[PRIM-25]] | Communicative-De-hallucination Role-Flip Gate | `internal/agents/roleflip.go` | absent-implemented |
| [[PRIM-26]] | Symbol-Aware Navigator | `internal/agents/navigator.go` | present |
| [[PRIM-27]] | Coverage-Guided Plateau Detection | `internal/agents/plateau.go` | present |
| [[PRIM-28]] | Conversable State Checkpoints & Interrupts | `internal/agents/checkpoint.go` | present |
| [[PRIM-29]] | Recruitment-Adaptive Planning | `internal/agents/recruit.go` | absent-implemented |
| [[PRIM-30]] | Iterative Retrieval Refinement | `internal/agents/iter_retrieval.go` | absent-implemented |

## Partial — extend in next sprint

| ID | Title | Where | Status |
| -- | ----- | ----- | ------ |
| [[PRIM-21]] | Migration Strategy Selection | `internal/agents/strategy.go` | present (5 strategies; matches Mueller 5-strategy LIS) |
| [[PRIM-24]] | SOP-Anchored Role-Artifact Schema | `internal/artifacts/versioning.go` | partial (SchemaVersion field added this sprint) |

## Referenced — paper cites but no codebase primitive

These primitives appear in the paper appendix but are derived from cited
literature without a dedicated codebase primitive. They are tracked here
for parity accounting; status `referenced`.

| ID | Title | Source citation | Reason |
| -- | ----- | --------------- | ------ |
| [[PRIM-15]] | Evidence-First Adaptation Pattern | [[Reeper-2024]] | method pattern; not a code primitive |
| [[PRIM-16]] | Spec-Driven Development Lifecycle | [[spec-kit]] | development lifecycle; not a code primitive |
| [[PRIM-18]] | Jaccard-Coupling Architecture Recovery | [[MSR4SA-2017]] | analysis technique; future archaeologist extension |
| [[PRIM-19]] | Design Rule Hierarchy Partitioning | [[Kazman-2000]] | analysis technique; future architect extension |
| [[PRIM-20]] | Concept Assignment and Redocumentation | [[Rajlich-1997]] | analysis technique; future archaeologist extension |
| [[PRIM-22]] | Four Phases of Comprehension | [[Foltz-2023]] | cognitive model; no code primitive |

## Removed

- **[[PRIM-2-legacy]] (Navigator agent)** — was a 5th parallel graph node.
  Removed: symbol-name resolution is now a sub-mechanism of the Planner
  agent invoking the LSP provider directly. The figure
  `magh-pipeline.workflow.json` reflects this. Superseded by
  [[PRIM-26]] Symbol-Aware Navigator (sub-mechanism of planner).

## Package Cohesion

- **`internal/artifacts`** — output data shapes (versioned via
  [[PRIM-24]]); no upstream imports.
- **`internal/types`** — `State` + `TranslationTask`.
- **`internal/consts`** — compile-time constants (error values,
  toolchain names, LSP provider names).

See [[Architecture]] for the dependency DAG and [[METHODOLOGY]] for the
runtime pipeline.

## Versioning convention

All version markers use the `[[x.y.z ...]]` convention. Spontaneous
parentheses like `(CodaMOSA)`, `(Syzygy)`, `(ReCodeAgent)` are NOT
allowed; replace with `[[Author-YEAR]]` or `[[P-NN]]` style backlinks.
