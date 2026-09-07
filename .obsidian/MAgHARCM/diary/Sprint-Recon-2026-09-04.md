---
title: Sprint Recon — Primitive Parity
backlink: [[Sprint-Recon]]
tags: [recon, primitive-parity, sprint-2026-09-04]
---

# Sprint Recon — Primitive Parity

Cross-reference of paper appendix primitives (`docs/.paper/sec_appendix_primitives.tex`) against
`internal/agents/*.go` to identify gaps between claimed and implemented.

## Status matrix

| PRIM-NN | Paper primitive | Code location | Status |
| -- | -- | -- | -- |
| PRIM-1 | Reverse Topological Translation Ordering | `internal/agents/planning.go::ComputeReverseTopoOrder` | present |
| PRIM-2 | Back-Edge-Conditioned Reverse-Topological Scheduling | `internal/agents/planning.go` | present |
| PRIM-3 | Target Skeleton-First Generation | `internal/agents/planning.go::DefaultProjectSkeleton` | present |
| PRIM-4 | SpecMiner-Style Dynamic Invariant Recovery | — | **absent** |
| PRIM-5 | Test Suite Co-Translation & Synthesis | `internal/agents/validator.go::generateAdditionalTests` | present |
| PRIM-6 | Multi-Stage Build/Test Feedback Repair | `internal/agents/validator.go` | present |
| PRIM-7 | Multi-Agent Verdict Validation | — | **absent** |
| PRIM-8 | State-Grounded Mock-Based In-Isolation Validation | — | **absent** |
| PRIM-9 | Tri-Representation Hybrid Code Graph | — | **absent** |
| PRIM-10 | Feature-Mapping & Type-Compatibility Validation | `internal/agents/manifest_rewriter.go` | partial |
| PRIM-11 | Implementation-Agnostic Testing | — | **absent** |
| PRIM-12 | Wasm-Based Reference Execution Oracle | — | **absent** |
| PRIM-13 | Adversarial Test-Weakening Guard | `internal/agents/validator.go` | present |
| PRIM-14 | Software-Archaeology Stage | — | **absent** |
| PRIM-15 | Evidence-First Adaptation Pattern | — | **absent** |
| PRIM-16 | Spec-Driven Development Lifecycle | — | **absent** |
| PRIM-17 | Asynchronous SE Agent Blackboard (CAID) | — | **absent** |
| PRIM-18 | Jaccard-Coupling Architecture Recovery | — | **absent** |
| PRIM-19 | Design Rule Hierarchy Partitioning | — | **absent** |
| PRIM-20 | Concept Assignment and Redocumentation | — | **absent** |
| PRIM-21 | Migration Strategy Selection | `internal/agents/strategy.go` | present |
| PRIM-22 | Four Phases of Comprehension | — | **absent** |
| PRIM-23 | Chunked Translation | `internal/agents/chunked_translator.go` | present |
| PRIM-24 | SOP-Anchored Role-Artifact Schema | `internal/artifacts/*` | partial (no version field) |
| PRIM-25 | Communicative-De-hallucination Role-Flip Gate | — | **absent** |
| PRIM-26 | Symbol-Aware Navigator | `internal/agents/navigator.go` | present |
| PRIM-27 | Coverage-Guided Plateau Detection | `internal/agents/plateau.go` | present |
| PRIM-28 | Conversable State Checkpoints & Interrupts | `internal/agents/checkpoint.go` | present |
| PRIM-29 | Recruitment-Adaptive Planning | — | **absent** |
| PRIM-30 | Source-to-Target Manifest Rewriter | `internal/agents/manifest_rewriter.go` | present |
| PRIM-31 | Iterative Retrieval Refinement | `internal/agents/iter_retrieval.go` | absent-implemented |

## Decision

Implement the 12 high-value absents in the codebase: PRIM-4, PRIM-7, PRIM-8, PRIM-9,
(add schema version field).

Skip the software-archaeology *analysis* primitives (PRIM-15, PRIM-16, PRIM-18, PRIM-19,
PRIM-20, PRIM-22) as standalone implementations — these are described in the paper appendix
as conceptual primitives drawn from the cited literature; tracking them in
[[primitives/INDEX]] with status **referenced** (not implemented) keeps the parity surface
honest. PRIM-14 (Software-Archaeology Stage) is implemented as a planning-stage hook because
it is invoked by the pipeline; the deeper conceptual primitives stay referenced.

## Vault versioning convention

Use `[[x.y.z ...]]` for every ID/version marker. Replace spontaneous parentheses
(`(CodaMOSA)`, `(ReCodeAgent)`, etc.) with `[[P-NN]]` or `[[Author-YEAR]]` style backlinks.

Examples:

- `(CodaMOSA)` → `[[P27 Coverage-Guided Plateau Detection]]`
- `(Syzygy)` → `[[P58 Syzygy]]`
- `(pp-besm dev.to playbook)` → `[[pp-besm Software Archaeology]]`

See [[primitives/INDEX]], [[research/Architecture]], [[research/METHODOLOGY]].

## Recursive research anchor

Bibliography recursion will pivot off PRIM-1..30 with depth bound = 2 hops. Software-archaeology
sources get priority in hop-2: pp-besm dev.to playbook, AgentPatterns.ai Legacy Code Archaeology,
Rajlich, Müller, Foltz, Baldwin & Clark, Kazman & Cai.

## See also

- `[[1.0.0 ADR-2026-09-26-Vault-Lint-Extension]]` — vault lint script (`scripts/lint_vault.sh`) retroactively validates this recon note's `[[x.y.z ...]]` convention.
