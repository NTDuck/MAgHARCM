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
| `[[1.0.0 PRIM-3]]` | Target Skeleton-First Generation | `internal/agents/planning.go::DefaultProjectSkeleton` | [[Skel-2024]], [[1.0.0 P-34]] Baldwin & Clark, [[P-42]] DSM ([[1.0.0 P-42 maccormack-dsm-2006]]), [[P-54]] Phi-3 ([[1.0.0 P-54 abdin-phi3-techreport-2024]]) | implemented |
| `[[1.0.0 PRIM-4]]` | SpecMiner Dynamic Invariant Recovery | `internal/agents/specminer.go` | [[Syzygy-2024]], Ernst Daikon | implemented |
| `[[1.0.0 PRIM-5]]` | Test Suite Co-Translation & Synthesis | `internal/agents/validator.go::generateAdditionalTests` | [[Pynguin-2021]], [[1.0.0 P-24 Feathers]] (2004), [[P-56]] Multi-SWE-bench ([[1.0.0 P-56 zan-multiswebench-2025]] Zan) | implemented |
| `[[1.0.0 PRIM-6]]` | Multi-Stage Build/Test Feedback Repair | `internal/agents/validator.go::Run` | [[AlphaTrans-2024]], Le Goues, [[P-54]] Phi-3 ([[1.0.0 P-54 abdin-phi3-techreport-2024]] Abdin) | implemented |
| `[[1.0.0 PRIM-7]]` | Multi-Agent Verdict Validation | `internal/agents/verdict_panel.go` | [[MatchFixAgent-2024]], Avizienis, [[P-57]] Speculative Decoding ([[1.0.0 P-57 leviathan-speculative-decoding-2023]] Leviathan) | implemented |
| `[[1.0.0 PRIM-8]]` | State-Grounded Mock-Based In-Isolation Validation | `internal/agents/mock_validator.go` | [[TRAM-2024]], [[1.0.0 P-24 Feathers]] (2004) | implemented |
| `[[1.0.0 PRIM-9]]` | Tri-Representation Hybrid Code Graph | `internal/agents/cpg.go` | [[RepoGraph-2024]], [[NEEDS-LINK Yamaguchi-2014]] | implemented |
| `[[1.0.0 PRIM-10]]` | Feature-Mapping & Type-Compatibility Validation | `internal/agents/feature_mapping.go` | [[Oxidizer-2023]], Czarnecki, [[P-43]] FODA ([[1.0.0 P-43 kang-foda-1990]]) | implemented |
| `[[1.0.0 PRIM-11]]` | Implementation-Agnostic Testing | `internal/agents/impl_agnostic.go` | [[RepoMod-Bench-2024]], Weyuker | implemented |
| `[[1.0.0 PRIM-12]]` | Wasm-Based Reference Execution Oracle | `internal/agents/wasm_oracle.go` | [[VERT-2024]], McKeeman | implemented |
| `[[1.0.0 PRIM-13]]` | Adversarial Test-Weakening Guard | `internal/agents/validator.go::verifyNoTestWeakening` | [[AdvTestGen-2024]], ISSTA, [[P-48]] Jia & [[1.0.0 P-48 jia-harman-mutation-2011]] (2011) Mutation Survey | implemented |
| `[[1.0.0 PRIM-14]]` | Software-Archaeology Stage | `internal/agents/archaeology.go` | pp-besm, Chikofsky & Cross ([[1.0.0 P-33 chikofsky-cross-1990]]), [[P-46]] Seacord et al. ([[1.0.0 P-46 seacord-modernizing-2003]]), [[P-54]] Phi-3 ([[1.0.0 P-54 abdin-phi3-techreport-2024]]) | implemented |
| `[[1.0.0 PRIM-15]]` | Evidence-First Adaptation Pattern | `internal/agents/evidence_adaptation.go` | [[Reeper-2024]], Cleanroom Mills | implemented |
| `[[1.0.0 PRIM-16]]` | Spec-Driven Development Lifecycle | `internal/agents/spec_lifecycle.go` | [[spec-kit-2024]], Meyer DbC | implemented |
| `[[1.0.0 PRIM-17]]` | Asynchronous SE Agent Blackboard | `internal/agents/blackboard.go` | [[CAID-2024]], Nii ([[NEEDS-LINK Nii-1986]]), [[P-44]] [[1.0.0 P-44 corkill-blackboard-1991]] | implemented |
| `[[1.0.0 PRIM-18]]` | Jaccard-Coupling Architecture Recovery | `internal/agents/jaccard_coupling.go` | [[MSR4SA-2017]], Hassan MSR, [[P-47]] Gall Hajek Jazayeri ([[1.0.0 P-47 gall-coupling-1998]]) | implemented |
| `[[1.0.0 PRIM-19]]` | Design Rule Hierarchy Partitioning | `internal/agents/design_rule_hierarchy.go` | [[Kazman-2017]], Baldwin & Clark, [[P-42]] [[1.0.0 P-42 maccormack-dsm-2006]] (2006) DSM | implemented |
| `[[1.0.0 PRIM-20]]` | Concept Assignment and Redocumentation | `internal/agents/concept_assignment.go` | [[Rajlich-1997]], Biggerstaff | implemented |
| `[[1.0.0 PRIM-21]]` | Migration Strategy Selection | `internal/agents/strategy.go::Registry.TryInOrder` | [[Müller-2000]] (5 Strategies), [[P-57]] Speculative Decoding ([[1.0.0 P-57 leviathan-speculative-decoding-2023]]) | implemented |
| `[[1.0.0 PRIM-22]]` | Four Phases of Comprehension | `internal/agents/comprehension.go` | [[Foltz-2023]] DR. JONES Model, [[P-54]] Phi-3 ([[1.0.0 P-54 abdin-phi3-techreport-2024]] Abdin), [[P-55]] SLM Few-Shot ([[1.0.0 P-55 schick-schutze-slm-fewshot-2021]] Schick), [[P-56]] Multi-SWE-bench ([[1.0.0 P-56 zan-multiswebench-2025]] Zan) | implemented |
| `[[1.0.0 PRIM-23]]` | Chunked Translation | `internal/agents/chunked_translator.go` | [[ChatDev-2023]], MetaGPT | implemented |
| `[[1.0.0 PRIM-24]]` | SOP-Anchored Role-Artifact Schema | `internal/compiletime/compiletime.go`, `internal/agents/state.go` | [[MetaGPT-2023]] SOP Contracts, [[P-45]] Curtis-Kellner-Over ([[1.0.0 P-45 curtis-process-modeling-1992]]), [[P-55]] SLM Few-Shot ([[1.0.0 P-55 schick-schutze-slm-fewshot-2021]]) | implemented |
| `[[1.0.0 PRIM-25]]` | Communicative-De-hallucination Role-Flip Gate | `internal/agents/roleflip.go` | [[ChatDev-2023]] Sycophancy Gate, [[P-55]] SLM Few-Shot ([[1.0.0 P-55 schick-schutze-slm-fewshot-2021]]), [[P-38]] Sycophancy Anchor ([[1.0.0 P-38 anthropic-sycophancy-2025]]) | implemented |
| `[[1.0.0 PRIM-26]]` | Symbol-Aware Navigator | `internal/agents/navigator.go` | [[HyperAgent-2024]], ABCoder, [[P-50]] [[1.0.0 P-50 jiang-llm-code-survey-2024]] (2024) Code-Gen Survey | implemented |
| `[[1.0.0 PRIM-27]]` | Coverage-Guided Plateau Detection | `internal/agents/plateau.go` | [[CodaMOSA-2023]], Harman SBST | implemented |
| `[[1.0.0 PRIM-28]]` | Conversable State Checkpoints & Interrupts | `internal/agents/checkpoint.go` | [[AutoGen-2023]], Snapshotting | implemented |
| `[[1.0.0 PRIM-29]]` | Recruitment-Adaptive Planning | `internal/agents/recruit.go` | [[AgentVerse-2023]], Dynamic Org, [[P-49]] Shehory & Kraus ([[1.0.0 P-49 shehory-kraus-coalition-1998]]) | implemented |
| `[[1.0.0 PRIM-30]]` | Source-to-Target Manifest Rewriter | `internal/agents/manifest_rewriter.go` | [[Syzygy-2024]], JavaC2Rust | implemented |
| `[[1.0.0 PRIM-31]]` | Iterative Retrieval Refinement | `internal/agents/iter_retrieval.go` | [[RepoCoder-2024]], Dynamic Context, [[P-56]] Multi-SWE-bench ([[1.0.0 P-56 zan-multiswebench-2025]] Zan), [[P-57]] Speculative Decoding ([[1.0.0 P-57 leviathan-speculative-decoding-2023]] Leviathan) | implemented |

---

## Single Versioning Convention

All version markers use the `[[x.y.z ...]]` convention. Spontaneous parentheses (such as `(CodaMOSA)` or `(ReCodeAgent)`) are forbidden; replace them with `[[Author-YEAR]]` or `[[P-NN]]` style backlinks.
All 31 primitives are compiled into the codebase, preserving architectural integrity.
