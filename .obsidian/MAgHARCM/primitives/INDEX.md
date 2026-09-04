---
title: Primitives Index
backlink: [[2.0.0 Primitives]]
tags: [primitives, catalog, status, [[1.0.0 PRIM-1]]..[[1.0.0 PRIM-31]]]
---

# [[2.0.0 Primitives Index]]

A primitive is a small, named, reusable capability within the MAgHARCM pipeline.
Each entry links to its specification, theoretical lineage, and implementing code file.
All 31 primitives are active implementations in the codebase, establishing a single source of truth.

Status legend:
- `implemented`: fully active and wired in the codebase.
- `integrated`: active sub-mechanism or cascade component.

---

## Complete Catalog of 31 Primitives

| ID | Title | Implementation Location | Lineage | Status |
| :--- | :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-1]]` | Reverse Topological Translation Ordering | `internal/agents/planning.go::ComputeReverseTopoOrder` | [[ReCodeAgent-2024]], Parnas (1972) | implemented |
| `[[1.0.0 PRIM-2]]` | Back-Edge-Conditioned Scheduling | `internal/agents/planning.go::linearizeDAG` | [[AlphaTrans-2024]], Tarjan | implemented |
| `[[1.0.0 PRIM-3]]` | Target Skeleton-First Generation | `internal/agents/planning.go::DefaultProjectSkeleton` | [[Skel-2024]], Baldwin & Clark (2000) | implemented |
| `[[1.0.0 PRIM-4]]` | SpecMiner Dynamic Invariant Recovery | `internal/agents/specminer.go` | [[Syzygy-2024]], Ernst Daikon | implemented |
| `[[1.0.0 PRIM-5]]` | Test Suite Co-Translation & Synthesis | `internal/agents/validator.go::generateAdditionalTests` | [[Pynguin-2021]], Feathers (2004) | implemented |
| `[[1.0.0 PRIM-6]]` | Multi-Stage Build/Test Feedback Repair | `internal/agents/validator.go::Run` | [[AlphaTrans-2024]], Le Goues | implemented |
| `[[1.0.0 PRIM-7]]` | Multi-Agent Verdict Validation | `internal/agents/verdict_panel.go` | [[MatchFixAgent-2024]], Avizienis | implemented |
| `[[1.0.0 PRIM-8]]` | State-Grounded Mock-Based In-Isolation Validation | `internal/agents/mock_validator.go` | [[TRAM-2024]], Feathers (2004) | implemented |
| `[[1.0.0 PRIM-9]]` | Tri-Representation Hybrid Code Graph | `internal/agents/cpg.go` | [[RepoGraph-2024]], Yamaguchi (2014) | implemented |
| `[[1.0.0 PRIM-10]]` | Feature-Mapping & Type-Compatibility Validation | `internal/agents/feature_mapping.go` | [[Oxidizer-2023]], Czarnecki | implemented |
| `[[1.0.0 PRIM-11]]` | Implementation-Agnostic Testing | `internal/agents/impl_agnostic.go` | [[RepoMod-Bench-2024]], Weyuker | implemented |
| `[[1.0.0 PRIM-12]]` | Wasm-Based Reference Execution Oracle | `internal/agents/wasm_oracle.go` | [[VERT-2024]], McKeeman | implemented |
| `[[1.0.0 PRIM-13]]` | Adversarial Test-Weakening Guard | `internal/agents/validator.go::verifyNoTestWeakening` | [[AdvTestGen-2024]], ISSTA | implemented |
| `[[1.0.0 PRIM-14]]` | Software-Archaeology Stage | `internal/agents/archaeology.go` | pp-besm, Chikofsky & Cross (1990) | implemented |
| `[[1.0.0 PRIM-15]]` | Evidence-First Adaptation Pattern | `internal/agents/evidence_adaptation.go` | [[Reeper-2024]], Cleanroom Mills | implemented |
| `[[1.0.0 PRIM-16]]` | Spec-Driven Development Lifecycle | `internal/agents/spec_lifecycle.go` | [[spec-kit-2024]], Meyer DbC | implemented |
| `[[1.0.0 PRIM-17]]` | Asynchronous SE Agent Blackboard | `internal/agents/blackboard.go` | [[CAID-2024]], Nii (1986) | implemented |
| `[[1.0.0 PRIM-18]]` | Jaccard-Coupling Architecture Recovery | `internal/agents/jaccard_coupling.go` | [[MSR4SA-2017]], Hassan MSR | implemented |
| `[[1.0.0 PRIM-19]]` | Design Rule Hierarchy Partitioning | `internal/agents/design_rule_hierarchy.go` | [[Kazman-2017]], Baldwin & Clark | implemented |
| `[[1.0.0 PRIM-20]]` | Concept Assignment and Redocumentation | `internal/agents/concept_assignment.go` | [[Rajlich-1997]], Biggerstaff | implemented |
| `[[1.0.0 PRIM-21]]` | Migration Strategy Selection | `internal/agents/strategy.go::Registry.TryInOrder` | [[Müller-2000]] (5 Strategies) | implemented |
| `[[1.0.0 PRIM-22]]` | Four Phases of Comprehension | `internal/agents/comprehension.go` | [[Foltz-2023]] DR. JONES Model | implemented |
| `[[1.0.0 PRIM-23]]` | Chunked Translation | `internal/agents/chunked_translator.go` | [[ChatDev-2023]], MetaGPT | implemented |
| `[[1.0.0 PRIM-24]]` | SOP-Anchored Role-Artifact Schema | `internal/compiletime/compiletime.go`, `internal/agents/state.go` | [[MetaGPT-2023]] SOP Contracts | implemented |
| `[[1.0.0 PRIM-25]]` | Communicative-De-hallucination Role-Flip Gate | `internal/agents/roleflip.go` | [[ChatDev-2023]] Sycophancy Gate | implemented |
| `[[1.0.0 PRIM-26]]` | Symbol-Aware Navigator | `internal/agents/navigator.go` | [[HyperAgent-2024]], ABCoder | implemented |
| `[[1.0.0 PRIM-27]]` | Coverage-Guided Plateau Detection | `internal/agents/plateau.go` | [[CodaMOSA-2023]], Harman SBST | implemented |
| `[[1.0.0 PRIM-28]]` | Conversable State Checkpoints & Interrupts | `internal/agents/checkpoint.go` | [[AutoGen-2023]], Snapshotting | implemented |
| `[[1.0.0 PRIM-29]]` | Recruitment-Adaptive Planning | `internal/agents/recruit.go` | [[AgentVerse-2023]], Dynamic Org | implemented |
| `[[1.0.0 PRIM-30]]` | Source-to-Target Manifest Rewriter | `internal/agents/manifest_rewriter.go` | [[Syzygy-2024]], JavaC2Rust | implemented |
| `[[1.0.0 PRIM-31]]` | Iterative Retrieval Refinement | `internal/agents/iter_retrieval.go` | [[RepoCoder-2024]], Dynamic Context | implemented |

---

## Single Versioning Convention

All version markers use the `[[x.y.z ...]]` convention. Spontaneous parentheses (such as `(CodaMOSA)` or `(ReCodeAgent)`) are forbidden; replace them with `[[Author-YEAR]]` or `[[P-NN]]` style backlinks.
All 31 primitives are compiled into the codebase, preserving architectural integrity.
