---
title: Primitives Index
backlink: "[[2.0.0 Primitives Index]]"
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
| `[[1.0.0 PRIM-1]]` | Reverse Topological Translation Ordering | `internal/agents/planning.go::ComputeReverseTopoOrder` | [[ReCodeAgent-2024]], [[1.0.0 P-31]], [[1.0.0 P-86]] LLM-empowered modernization survey | implemented |
| `[[1.0.0 PRIM-2]]` | Back-Edge-Conditioned Scheduling | `internal/agents/planning.go::linearizeDAG` | [[AlphaTrans-2024]], Tarjan | implemented |
| `[[1.0.0 PRIM-3]]` | Target Skeleton-First Generation | `internal/agents/planning.go::DefaultProjectSkeleton` | [[Skel-2024]], [[1.0.0 P-34]] Baldwin & Clark, [[1.0.0 P-42]] DSM, [[1.0.0 P-54]] Phi-3, [[1.0.0 P-88]] HiTyper type-annotation migration, [[1.0.0 P-89]] legacy-modernization baseline | implemented |
| `[[1.0.0 PRIM-4]]` | SpecMiner Dynamic Invariant Recovery | `internal/agents/specminer.go` | [[Syzygy-2024]], Ernst Daikon | implemented |
| `[[1.0.0 PRIM-5]]` | Test Suite Co-Translation & Synthesis | `internal/agents/validator.go::generateAdditionalTests` | [[Pynguin-2021]], [[1.0.0 P-24]] Feathers, [[1.0.0 P-56]] Multi-SWE-bench | implemented |
| `[[1.0.0 PRIM-6]]` | Multi-Stage Build/Test Feedback Repair | `internal/agents/validator.go::Run` | [[AlphaTrans-2024]], Le Goues, [[1.0.0 P-54]] Phi-3 | implemented |
| `[[1.0.0 PRIM-7]]` | Multi-Agent Verdict Validation | `internal/agents/verdict_panel.go` | [[MatchFixAgent-2024]], Avizienis, [[1.0.0 P-57]] Speculative Decoding, [[1.0.0 P-61]] RedCode, [[1.0.0 P-78]] EAGLE-3, [[1.0.0 P-83]] code-specialised self-consistency, [[1.0.0 P-92]] Lightman PRM, [[1.0.0 P-95]] Code Llama verifier baseline | implemented |
| `[[1.0.0 PRIM-8]]` | State-Grounded Mock-Based In-Isolation Validation | `internal/agents/mock_validator.go` | [[TRAM-2024]], [[1.0.0 P-24]] Feathers | implemented |
| `[[1.0.0 PRIM-9]]` | Tri-Representation Hybrid Code Graph | `internal/agents/cpg.go` | [[RepoGraph-2024]], Yamaguchi | implemented |
| `[[1.0.0 PRIM-10]]` | Feature-Mapping & Type-Compatibility Validation | `internal/agents/feature_mapping.go` | [[Oxidizer-2023]], Czarnecki, [[1.0.0 P-43]] FODA | implemented |
| `[[1.0.0 PRIM-11]]` | Implementation-Agnostic Testing | `internal/agents/impl_agnostic.go` | [[RepoMod-Bench-2024]], Weyuker | implemented |
| `[[1.0.0 PRIM-12]]` | Wasm-Based Reference Execution Oracle | `internal/agents/wasm_oracle.go` | [[VERT-2024]], McKeeman | implemented |
| `[[1.0.0 PRIM-13]]` | Adversarial Test-Weakening Guard | `internal/agents/validator.go::verifyNoTestWeakening` | [[AdvTestGen-2024]], ISSTA, [[1.0.0 P-48]] Jia & Harman Mutation Survey | implemented |
| `[[1.0.0 PRIM-14]]` | Software-Archaeology Stage | `internal/agents/archaeology.go` | pp-besm, [[1.0.0 P-33]] Chikofsky & Cross, [[1.0.0 P-46]] Seacord, [[1.0.0 P-54]] Phi-3, [[1.0.0 P-59]] Digital Apollo, [[1.0.0 P-87]] TOSEM SLR on LLM4SE | implemented |
| `[[1.0.0 PRIM-15]]` | Evidence-First Adaptation Pattern | `internal/agents/evidence_adaptation.go` | [[Reeper-2024]], Cleanroom Mills | implemented |
| `[[1.0.0 PRIM-16]]` | Spec-Driven Development Lifecycle | `internal/agents/spec_lifecycle.go` | [[spec-kit-2024]], Meyer DbC | implemented |
| `[[1.0.0 PRIM-17]]` | Asynchronous SE Agent Blackboard | `internal/agents/blackboard.go` | [[CAID-2024]], Nii Blackboard, [[1.0.0 P-44]] Corkill | implemented |
| `[[1.0.0 PRIM-18]]` | Jaccard-Coupling Architecture Recovery | `internal/agents/jaccard_coupling.go` | [[MSR4SA-2017]], Hassan MSR, [[1.0.0 P-47]] Gall Hajek Jazayeri, [[1.0.0 P-59]] Digital Apollo | implemented |
| `[[1.0.0 PRIM-19]]` | Design Rule Hierarchy Partitioning | `internal/agents/design_rule_hierarchy.go` | [[Kazman-2017]], Baldwin & Clark, [[1.0.0 P-42]] DSM, [[1.0.0 P-59]] Digital Apollo | implemented |
| `[[1.0.0 PRIM-20]]` | Concept Assignment and Redocumentation | `internal/agents/concept_assignment.go` | [[Rajlich-1997]], Biggerstaff, [[1.0.0 P-59]] Digital Apollo | implemented |
| `[[1.0.0 PRIM-21]]` | Migration Strategy Selection | `internal/agents/strategy.go::Registry.TryInOrder` | [[1.0.0 P-35]] Müller 5 Strategies, [[1.0.0 P-57]] Speculative Decoding, [[1.0.0 P-60]] Bisbal Cascading Fallback, [[1.0.0 P-84]] s1 test-time scaling, [[1.0.0 P-91]] Snell test-time compute allocation, [[1.0.0 P-98]] Large Language Monkeys | implemented |
| `[[1.0.0 PRIM-22]]` | Four Phases of Comprehension | `internal/agents/comprehension.go` | [[1.0.0 P-40]] DR. JONES Model, [[1.0.0 P-54]] Phi-3, [[1.0.0 P-55]] SLM Few-Shot, [[1.0.0 P-56]] Multi-SWE-bench, [[1.0.0 P-58]] Qwen2.5-Coder, [[1.0.0 P-62]] Decomposed Prompting, [[1.0.0 P-80]] StreamingLLM, [[1.0.0 P-90]] Wei CoT, [[1.0.0 P-94]] LIMA curation, [[1.0.0 P-96]] Zero-Shot CoT, [[1.0.0 P-99]] BIG-Bench Hard, [[1.0.0 P-100]] Decomposed Prompting (SLM multi-agent) | implemented |
| `[[1.0.0 PRIM-23]]` | Chunked Translation | `internal/agents/chunked_translator.go` | [[ChatDev-2023]], [[MetaGPT-2023]], [[1.0.0 P-58]] Qwen2.5-Coder, [[1.0.0 P-62]] Decomposed Prompting, [[1.0.0 P-65]] PAL, [[1.0.0 P-101]] Least-to-Most Prompting | implemented |
| `[[1.0.0 PRIM-24]]` | SOP-Anchored Role-Artifact Schema | `internal/compiletime/compiletime.go`, `internal/agents/state.go` | [[MetaGPT-2023]] SOP Contracts, [[1.0.0 P-45]] Curtis-Kellner-Over, [[1.0.0 P-55]] SLM Few-Shot, [[1.0.0 P-62]] Decomposed Prompting, [[1.0.0 P-63]] Multi-Agent Survey, [[1.0.0 P-93]] DPO alignment, [[1.0.0 P-95]] Code Llama instruction tuning, [[1.0.0 P-100]] Decomposed Prompting (SLM multi-agent) | implemented |
| `[[1.0.0 PRIM-25]]` | Communicative-De-hallucination Role-Flip Gate | `internal/agents/roleflip.go` | [[ChatDev-2023]] Sycophancy Gate, [[1.0.0 P-38]] Anthropic Sycophancy, [[1.0.0 P-55]] SLM Few-Shot, [[1.0.0 P-58]] Qwen2.5-Coder, [[1.0.0 P-61]] RedCode, [[1.0.0 P-81]] Gorilla, [[1.0.0 P-85]] function-calling SLM gap | implemented |
| `[[1.0.0 PRIM-26]]` | Symbol-Aware Navigator | `internal/agents/navigator.go` | [[HyperAgent-2024]], ABCoder, [[1.0.0 P-50]] Code-Gen Survey, [[1.0.0 P-58]] Qwen2.5-Coder | implemented |
| `[[1.0.0 PRIM-27]]` | Coverage-Guided Plateau Detection | `internal/agents/plateau.go` | [[CodaMOSA-2023]], Harman SBST | implemented |
| `[[1.0.0 PRIM-28]]` | Conversable State Checkpoints & Interrupts | `internal/agents/checkpoint.go` | [[AutoGen-2023]], Snapshotting, [[1.0.0 P-59]] Digital Apollo | implemented |
| `[[1.0.0 PRIM-29]]` | Recruitment-Adaptive Planning | `internal/agents/recruit.go` | [[AgentVerse-2023]], Dynamic Org, [[1.0.0 P-49]] Shehory & Kraus Coalition Formation, [[1.0.0 P-63]] Multi-Agent Survey | implemented |
| `[[1.0.0 PRIM-30]]` | Source-to-Target Manifest Rewriter | `internal/agents/manifest_rewriter.go` | [[Syzygy-2024]], JavaC2Rust | implemented |
| `[[1.0.0 PRIM-31]]` | Iterative Retrieval Refinement | `internal/agents/iter_retrieval.go` | [[RepoCoder-2024]], Dynamic Context, [[1.0.0 P-56]] Multi-SWE-bench, [[1.0.0 P-57]] Speculative Decoding, [[1.0.0 P-58]] Qwen2.5-Coder, [[1.0.0 P-62]] Decomposed Prompting | implemented |

---

## Single Versioning Convention

All version markers use the `[[x.y.z ...]]` convention. Spontaneous parentheses (such as `(CodaMOSA)` or `(ReCodeAgent)`) are forbidden; replace them with `[[Author-Year]]` or `[[x.y.z P-NN]]` style backlinks.
All 31 primitives are compiled into the codebase, preserving architectural integrity.

## Sprint 2026-09-09 Vault Sync Audit

- Indexed 31 primitives; all 31 are implemented in `internal/agents/*.go`.
- Reconciled backlinks to use `[[x.y.z P-NN]]` form (single-versioning convention).
- Added new anchors P-58..P-65 (Qwen2.5-Coder, Digital Apollo, Bisbal, RedCode, Decomposed Prompting, Multi-Agent Survey, EAGLE, PAL) to relevant primitive rows.

## Sprint 2026-09-12 Vault Sync Audit

- Indexed 6 new anchor papers (P-78 EAGLE-3, P-79 Wilde-Scully Software Reconnaissance, P-80 StreamingLLM, P-81 Gorilla, P-82 Pahins-Stegherr-Steinhauser UNVERIFIED, P-83 code-specialised self-consistency) into relevant primitive rows.

## Sprint 2026-09-14 Vault Sync Audit

- Indexed 6 new SLM-era anchor papers (P-90 Wei CoT, P-91 Snell test-time compute, P-92 Lightman PRM, P-93 Rafailov DPO, P-94 LIMA, P-95 Roziere Code Llama) — all verified — into relevant primitive rows.
- PR backlog (low-risk follow-ups): re-tag Yamaguchi (PRIM-9) and Nii (PRIM-17) backlinks once their standalone paper notes are materialised.
