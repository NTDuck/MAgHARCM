---
title: Software Archaeology & Legacy System Modernization Lineage
backlink: [[2.0.0 Software-Archaeology-Lineage]]
tags: [research, software-archaeology, modernization, lineage, synthesis, [[2.0.0 MAgHARCM]]]
---

# [[2.0.0 Software Archaeology & Legacy System Modernization Lineage]]

## 1. Executive Summary & Epistemological Stance

Software modernization using multi-agent language model architectures does not operate in a historical vacuum. When large language models (LLMs) or small specialized language models (SLMs) encounter legacy enterprise codebases, naive code translation fails because legacy software embodies decades of accumulated invariants, implicit domain assumptions, tacit developer conventions, and emergent coupling.

The MAgHARCM pipeline grounds modern agentic reasoning in 45+ years of rigorous software engineering literature.

---

## 2. Foundational Lineages & 2-Hop Citations

```
[Parnas 1972] Information Hiding & Modular Decomposition
       │
       ▼
[Lehman 1980] Laws of Software Evolution (E-Type Systems)
       │
       ▼
[Chikofsky & Cross 1990] Reverse Engineering & Design Recovery Taxonomy
       │
       ▼
[Rajlich & Bennett 2000] Staged Life Cycle & Concept Assignment (1997)
       │
       ├──────────────────────────────────────────┐
       ▼                                          ▼
[Müller et al. 2000] Roadmap               [Baldwin & Clark 2000] Design Rules
       │                                          │
       ▼                                          ▼
[Feathers 2004] Working with Legacy Code   [Kazman & Cai 2017] DRSpaces
       │                                          │
       └──────────────────┬───────────────────────┘
                          ▼
            [Foltz 2023] DR. JONES Cognition
                          │
                          ▼
        [pp-besm / Esoteric Production Archaeology]
                          │
                          ▼
                          │  ← modern LLM-era hops →
                          │
   ┌──────────────────────┼───────────────────────────┐
   ▼                      ▼                           ▼
[Shen 2023] HuggingGPT   [Raman 2025] Anthropic     [Fleming & Baldwin
  Controller-Expert        Sycophancy (PRIM-25)    2024] Pace of Modular
  (PRIM-29 hop-1)           RoleFlip anchor         Innovation
                                                     (PRIM-19 retrospective)
                          │
                          ▼
      [[2.0.0 MAgHARCM]] Multi-Agent Architecture
```

### 2.1. Decomposition, Information Hiding, & Modularity
- **Parnas (1972) — *On the Criteria To Be Used in Decomposing Systems into Modules*** (`[[1.0.0 P-31]]`):
  - *Core Insight*: Modules should be decomposed based on hidden design decisions (secrets) rather than operational execution steps. The most volatile elements (data structures, hardware quirks) must be enclosed behind unchanging abstract interfaces.
  - *Application in MAgHARCM*: Direct intellectual foundation for `[[1.0.0 PRIM-3]]` (Target Skeleton-First Generation) and `[[1.0.0 PRIM-19]]` (Design Rule Hierarchy Partitioning). Target interfaces (Rust traits) are generated as boundary constraints before method implementations are synthesized.

- **Baldwin & Clark (2000) — *Design Rules: The Power of Modularity*** (`[[1.0.0 P-34]]`):
  - *Core Insight*: Formalized the economic and structural value of modularity through six modular operators: Splitting, Substituting, Augmenting, Excluding, Inverting, and Porting. Defined Design Rules as structural decisions that decouple subsequent module implementations.
  - *Application in MAgHARCM*: Guides `[[1.0.0 PRIM-1]]` (Reverse Topological Ordering) and `[[1.0.0 PRIM-2]]` (Back-Edge Conditioned Scheduling). Identifying design rule interfaces breaks dependency cycles and permits isolated synthesis.

### 2.2. Software Evolution & Archaeological Strata
- **M. M. Lehman (1980, 1996) — *Programs, Life Cycles, and Laws of Software Evolution*** (`[[1.0.0 P-32]]`):
  - *Core Insight*: Formulated the fundamental laws of E-Type software evolution:
    1. *Continuing Change*: A system must continually adapt or become progressively less satisfactory.
    2. *Increasing Complexity*: As an evolving program is modified, its complexity increases unless work is done to maintain or reduce it.
    3. *Conservation of Familiarity*: The incremental growth of systems across releases is statistically invariant.
  - *Application in MAgHARCM*: Explains why legacy source code diverges from original documentation. MAgHARCM's `[[1.0.0 PRIM-14]]` (Software Archaeology Stage) does not trust comments; it analyzes evolutionary churn and executable reality.

- **Chikofsky & Cross (1990) — *Reverse Engineering and Design Recovery: A Taxonomy*** (`[[1.0.0 P-33]]`):
  - *Core Insight*: Formulated the canonical definitions distinguishing Forward Engineering, Reverse Engineering, Redocumentation, Design Recovery, Restructuring, and Reengineering. Defined Design Recovery as recreating software abstractions from a combination of code, external domain knowledge, and developer observation.
  - *Application in MAgHARCM*: Establishes the taxonomy for `[[1.0.0 PRIM-14]]` (Archaeology) and `[[1.0.0 PRIM-20]]` (Concept Assignment and Redocumentation).

### 2.3. Program Comprehension & Concept Assignment
- **Vaclav Rajlich (1997) & Keith Bennett (2000) — *Concept Assignment & Staged Lifecycle***:
  - *Core Insight*: Defined the concept locator methodology. Software maintenance requires mapping human-oriented domain concepts to specific computational locations (AST subtrees, classes, functions).
  - *Application in MAgHARCM*: Realized in `[[1.0.0 PRIM-20]]` (Concept Assignment and Redocumentation), where the archaeologist maps lexical identifiers and call clusters to semantic domain roles.

- **Hausi A. Müller et al. (2000) — *Reverse Engineering: A Roadmap*** (`[[1.0.0 P-35]]`):
  - *Core Insight*: Categorized program comprehension into three interacting elements: cognitive mental models (top-down vs. bottom-up), analysis techniques (static vs. dynamic), and tool interoperability. Formulated the 5 canonical migration strategies: Big Bang, Incremental, Pilot, Frozen Legacy, and Parallel Cutover.
  - *Application in MAgHARCM*: Formalized in `[[1.0.0 PRIM-21]]` (Migration Strategy Selection), implementing dynamic try-and-fail selection across Müller's five strategies.

- **Peter Foltz et al. (1998, 2023) — *DR. JONES Model of Cognitive Traversal*** (`[[1.0.0 P-30]]`):
  - *Core Insight*: Engineers comprehend unfamiliar code through six cognitive phases:
    1. *Decomposition*: Identifying structural units.
    2. *Recognition*: Recognizing known idioms and library patterns.
    3. *Organization*: Constructing hierarchical mental models.
    4. *Navigation*: Linear and associative traversal of call paths.
    5. *Explanation*: Synthesizing operational rationales.
    6. *Search*: Querying localized definitions.
  - *Application in MAgHARCM*: Directly implemented in `[[1.0.0 PRIM-22]]` (Four Phases of Comprehension) and `[[1.0.0 PRIM-26]]` (Symbol-Aware Navigator).

### 2.4. Legacy Modification & Safe Refactoring
- **Michael Feathers (2004) — *Working Effectively with Legacy Code***:
  - *Core Insight*: Defined legacy code rigorously: *"Legacy code is simply code without tests."* Established the Legacy Code Change Algorithm:
    1. Identify change points.
    2. Find test points.
    3. Break dependencies (using sensing and separation pins).
    4. Write automated characterization tests.
    5. Make changes and refactor.
  - *Application in MAgHARCM*: The foundational inspiration for `[[1.0.0 PRIM-5]]` (Test Suite Co-Translation & Synthesis), `[[1.0.0 PRIM-8]]` (State-Grounded Mock-Based In-Isolation Validation), and `[[1.0.0 PRIM-13]]` (Adversarial Test-Weakening Guard).
### 2.5. Esoteric & Industry Archaeological Playbooks
- **pp-besm (dev.to) — *Software Archaeology: Hunting for Lost Knowledge in Production Codebases*** (`[[1.0.0 P-36]]`):
  - *Core Insight*: Industrial legacy systems must be examined across five archaeological strata:
    1. *Commit Strata*: Git churn, author longevity, commit message sentiment.
    2. *Temporal Coupling*: Files that change together without explicit static references (Jaccard co-change metric).
    3. *Bug-Density Hotspots*: Modules with concentrated defect fixes.
    4. *Forensic Naming & Type Invariants*: Hungarian notation, obsolete typedefs, dead configuration flags.
    5. *Executable Time Capsules*: Historical compiler flags, environmental assumptions, and abandoned test harnesses.
  - *Application in MAgHARCM*: Guides `[[1.0.0 PRIM-14]]` (Software Archaeology Stage) and `[[1.0.0 PRIM-18]]` (Jaccard-Coupling Architecture Recovery).

- **AgentPatterns.ai — *Legacy Code Archaeology Pattern for Autonomous Systems*** (`[[1.0.0 P-37]]`):
  - *Core Insight*: Autonomous LLM multi-agent systems must execute an Archaeological excavation pass before synthesis to avoid context window saturation and blind transliteration failure.
  - *Application in MAgHARCM*: Directly motivates the 8-agent decomposition with dedicated Archaeologist and Reviewer execution units.
---

## 3. Comprehensive Mapping: Primitives to Literature Matrix

| Primitive | Title | Primary Source | Theoretical Lineage | Code Implementation |
| :--- | :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-1]]` | Reverse Topological Ordering | ReCodeAgent [[P-01]] | Parnas (1972) Information Hiding | `internal/agents/planning.go` |
| `[[1.0.0 PRIM-2]]` | Back-Edge Cycle Linearization | AlphaTrans [[P-02]] | Tarjan DAG decomposition | `internal/agents/planning.go` |
| `[[1.0.0 PRIM-3]]` | Target Skeleton-First Gen | Skel [[P-03]], ReCodeAgent; [[P-54]] Abdin (2024) Phi-3 Tech Report | Baldwin & Clark (2000) Design Rules; [[P-42]] DSM (MacCormack 2006); [[P-53]] Liu (2023) Lost in the Middle | `internal/agents/planning.go` |
| `[[1.0.0 PRIM-4]]` | SpecMiner Dynamic Invariants | Syzygy [[P-14]], Daikon | Dynamic Invariant Detection (Ernst) | `internal/agents/specminer.go` |
| `[[1.0.0 PRIM-5]]` | Test Co-Translation & Synth | Pynguin [[P-23]], ReCodeAgent; [[P-56]] Zan (2025) Multi-SWE-bench | Feathers (2004) Characterization Tests; multilingual SLM baseline (Java/Go/Rust/C) | `internal/agents/validator.go` |
| `[[1.0.0 PRIM-6]]` | Multi-Stage Build/Test Repair | AlphaTrans [[P-02]], ReCodeAgent; [[P-54]] Abdin (2024) Phi-3 | Automated Program Repair (Le Goues); Phi-3 chat-format alignment for repair prompts | `internal/agents/validator.go` |
| `[[1.0.0 PRIM-7]]` | Multi-Agent Verdict Validation | MatchFixAgent [[P-08]]; [[P-57]] Leviathan (2023) Speculative Decoding | N-Version Programming (Avizienis); [[P-52]] Wang (2023) Self-Consistency; speculative draft/target for SLM fleet | `internal/agents/verdict_panel.go` |
| `[[1.0.0 PRIM-8]]` | State-Grounded Mock Validation | TRAM [[P-10]] | Feathers (2004) Sensing Pins | `internal/agents/mock_validator.go` |
| `[[1.0.0 PRIM-9]]` | Tri-Representation Code Graph | RepoGraph [[P-14]], Yamaguchi | Code Property Graphs (Yamaguchi 2014) | `internal/agents/cpg.go` |
| `[[1.0.0 PRIM-10]]` | Feature-Mapping Validation | Oxidizer [[P-05]], RustRepoTrans; [[P-43]] FODA (Kang 1990) | Language Idiom Mapping (Czarnecki) | `internal/agents/feature_mapping.go` |
| `[[1.0.0 PRIM-11]]` | Implementation-Agnostic Test | RepoMod-Bench [[P-11]]; [[P-56]] Zan (2025) Multi-SWE-bench | Black-Box Specification Testing; multilingual real-issue benchmark | `internal/agents/impl_agnostic.go` |
| `[[1.0.0 PRIM-12]]` | Wasm Reference Oracle | VERT [[P-12]]; [[P-56]] Zan (2025) Multi-SWE-bench | Differential Execution Oracles (McKeeman); per-language differential baseline | `internal/agents/wasm_oracle.go` |
| `[[1.0.0 PRIM-13]]` | Adversarial Test Guard | AdvTestGen [[P-25]]; [[P-56]] Zan (2025) Multi-SWE-bench | Mutation Testing & Assertion Invariants, [[P-48]] Jia & Harman (2011) Mutation Survey; multilingual adversarial benchmark | `internal/agents/validator.go` |
| `[[1.0.0 PRIM-14]]` | Software-Archaeology Stage | pp-besm, AgentPatterns.ai; [[P-54]] Abdin (2024) Phi-3 | Chikofsky & Cross (1990) Reverse Eng, [[P-46]] Seacord et al. (2003) Modernizing Legacy Systems; Phi-3-mini-128K whole-file archaeology | `internal/agents/archaeology.go` |
| `[[1.0.0 PRIM-15]]` | Evidence-First Adaptation | Reeper [[P-16]] | Cleanroom Software Engineering (Mills) | `internal/agents/evidence_adaptation.go` |
| `[[1.0.0 PRIM-16]]` | Spec-Driven Dev Lifecycle | spec-kit [[P-17]] | Design-by-Contract (Meyer 1988) | `internal/agents/spec_lifecycle.go` |
| `[[1.0.0 PRIM-17]]` | Asynchronous SE Blackboard | CAID [[P-18]]; [[P-44]] Corkill (1991) | Blackboard Architecture (Nii 1986) | `internal/agents/blackboard.go` |
| `[[1.0.0 PRIM-18]]` | Jaccard-Coupling Recovery | MSR4SA [[P-19]] | Mining Software Repositories (Hassan), [[P-47]] Gall Hajek Jazayeri (1998) Logical Coupling | `internal/agents/jaccard_coupling.go` |
| `[[1.0.0 PRIM-19]]` | Design Rule Hierarchy Part | Kazman et al. [[P-20]], Baldwin-Clark | Baldwin & Clark (2000) Modularity, [[P-41]] 2024 retrospective, [[P-42]] MacCormack (2006) DSM | `internal/agents/design_rule_hierarchy.go` |
| `[[1.0.0 PRIM-20]]` | Concept Assignment & Redoc | Rajlich (1997) [[P-28]] | Concept Assignment (Biggerstaff 1993) | `internal/agents/concept_assignment.go` |
| `[[1.0.0 PRIM-21]]` | Migration Strategy Selection | Müller et al. [[P-06]]; [[P-57]] Leviathan (2023) Speculative Decoding | Legacy Migration Frameworks (Müller); [[P-52]] Wang (2023) Self-Consistency (strategy-level ensemble); speculative decoding strategy registry | `internal/agents/strategy.go` |
| `[[1.0.0 PRIM-22]]` | Four Phases Comprehension | Foltz (2023) [[P-30]]; canonical hop-1 [[P-40]]; [[P-51]] Cassano (2024) Can It Edit; [[P-54]] Abdin (2024) Phi-3; [[P-55]] Schick (2021) SLM Few-Shot; [[P-56]] Zan (2025) Multi-SWE-bench | Cognitive Program Comprehension (Soloway-Adelson, Pennington, Brooks, Détienne) | `internal/agents/comprehension.go` |
| `[[1.0.0 PRIM-23]]` | Chunked Translation | ChatDev [[P-12]], MetaGPT [[P-11]]; [[P-53]] Liu (2023) Lost in the Middle | Bounded-Context Translation | `internal/agents/chunked_translator.go` |
| `[[1.0.0 PRIM-24]]` | SOP-Anchored Role Artifact | MetaGPT [[P-11]]; [[P-45]] Curtis-Kellner-Over (1992); [[P-55]] Schick (2021) SLM Few-Shot | Standard Operating Procedures (SOP); cloze reformulation as artifact schema | `internal/compiletime/compiletime.go`, `internal/agents/state.go` |
| `[[1.0.0 PRIM-25]]` | Role-Flip De-Hallucination | ChatDev [[P-12]]; sycophancy anchor [[P-38]] (Raman 2025); [[P-51]] Cassano (2024) Can It Edit; [[P-55]] Schick (2021) SLM Few-Shot | Adversarial Verification (Sycophancy Gate); SLM few-shot role-flip prompt pattern | `internal/agents/roleflip.go` |
| `[[1.0.0 PRIM-26]]` | Symbol-Aware Navigator | HyperAgent [[P-13]], ABCoder; [[P-50]] Jiang (2024) Code-Gen Survey | Targeted Context Retrieval | `internal/agents/navigator.go` |
| `[[1.0.0 PRIM-27]]` | Coverage-Guided Plateau Det | CodaMOSA [[P-23]] | Search-Based Software Testing (Harman) | `internal/agents/plateau.go` |
| `[[1.0.0 PRIM-28]]` | Conversable Checkpoints | AutoGen [[P-26]] | State-Snapshotting & Resumption | `internal/agents/checkpoint.go` |
| `[[1.0.0 PRIM-29]]` | Recruitment-Adaptive Plan | AgentVerse [[P-26]]; controller-expert [[P-39]] (HuggingGPT / Jarvis) | Dynamic Team Organization, [[P-49]] Shehory & Kraus (1998) Coalition Formation | `internal/agents/recruit.go` |
| `[[1.0.0 PRIM-30]]` | Source-to-Target Manifest | Syzygy [[P-14]], JavaC2Rust | Dependency Graph Transpilation | `internal/agents/manifest_rewriter.go` |
| `[[1.0.0 PRIM-31]]` | Iterative Retrieval Refine | RepoCoder [[P-13]]; [[P-53]] Liu (2023) Lost in the Middle; [[P-56]] Zan (2025) Multi-SWE-bench; [[P-57]] Leviathan (2023) Speculative Decoding | Dynamic Feedback Retrieval-Augmented Gen; speculative decoding for retrieval-decoding | `internal/agents/iter_retrieval.go` |


---

## 4. Synthesis: The Archaeological Translation Cycle

```
[Legacy Codebase]
       │
       ▼
[Stage 1: Archaeological Excavation]
   ├── PRIM-14: Structural Boundary & Time Capsule Extraction
   ├── PRIM-18: Jaccard-Coupling Temporal Churn Detection
   ├── PRIM-19: Design Rule Hierarchy Layering (L1/L2/L3)
   ├── PRIM-20: Concept Assignment & Redocumentation
   └── PRIM-22: DR. JONES Cognitive Navigation
       │
       ▼
[Stage 2: Specification & Strategy Induction]
   ├── PRIM-15: Evidence-First Adaptation Spec
   ├── PRIM-16: Spec-Driven Lifecycle Gating
   ├── PRIM-21: Incremental Try-and-Fail Strategy Selection
   └── PRIM-4: SpecMiner Dynamic Invariant Recovery
       │
       ▼
[Stage 3: Architectural Planning]
   ├── PRIM-1 & PRIM-2: Reverse-Topological DAG Linearization
   ├── PRIM-3: Target Skeleton Generation
   ├── PRIM-29: Dynamic Agent/Tool Recruitment
   └── PRIM-30: Manifest & Dependency Transpilation
       │
       ▼
[Stage 4: Chunked Translation & Review]
   ├── PRIM-23: Bounded Fragment Dispatch
   ├── PRIM-26 & PRIM-31: Symbol-Aware Iterative Navigator
   └── PRIM-25: Communicative Role-Flip Review Gate
       │
       ▼
[Stage 5: Multi-Tier Validation Cascade]
   ├── PRIM-6: Multi-Stage Build/Test Feedback Repair
   ├── PRIM-5: Test Suite Co-Translation & Synthesis
   ├── PRIM-13: Adversarial Test-Weakening Guard
   ├── PRIM-27: Coverage Plateau Detection
   └── Optional Oracles: PRIM-7 (Verdict), PRIM-8 (Mocks), PRIM-11 (Agnostic), PRIM-12 (Wasm)
       │
       ▼
[Stage 6: Durable Snapshotting]
   └── PRIM-28: Conversable Checkpointing & Idempotent Resumption

## 5. Modern Synthesis Citations (2026 Sprint)

Four academic synthesis-anchor papers added in the 2026-09-05 sprint. Each provides a hop-1 academic reference for a software-archaeology primitive that previously cited only modern implementations.

- **MacCormack, Rusnak & Baldwin (2006) — *Exploring the Structure of Complex Software Designs (DSM)* (`[[1.0.0 P-42]]`)**:
  - *Core Insight*: Static analysis of Linux kernel and Apache HTTP server dependency matrices shows cycle size and cycle count predict post-release defect density; L1 interface stability precedes mass refactoring. DSM-based refactoring strategy (identify cycle → freeze interface → extract submodule) is the empirical operationalisation of skeleton-first synthesis.
  - *Application in MAgHARCM*: Reinforces `[[1.0.0 PRIM-3]]` (Target Skeleton-First Generation) and `[[1.0.0 PRIM-19]]` (Design Rule Hierarchy Partitioning) with quantitative evidence.

- **Kang et al. (1990) — *Feature-Oriented Domain Analysis (FODA) Feasibility Study* (`[[1.0.0 P-43]]`)**:
  - *Core Insight*: Introduces the feature model as a hierarchical tree of mandatory, optional, alternative, and mutually-exclusive end-user-visible characteristics, plus the capability matrix for mapping features to architectural components. Canonical ancestor of every modern feature-model formalism.
  - *Application in MAgHARCM*: Anchors `[[1.0.0 PRIM-10]]` (Feature-Mapping & Type-Compatibility Validation). Capability matrix is the direct ancestor of `feature_mapping.go`'s per-feature source-to-target mapping table.

- **Corkill (1991) — *Blackboard Systems* (`[[1.0.0 P-44]]`)**:
  - *Core Insight*: Practitioner-level survey of blackboard architectures. Five components (knowledge sources, blackboard, control shell, knowledge-base interaction protocol, triggering/event mechanism); four practical lessons (domain-driven KS partitioning, abstraction-level hierarchy, scheduler-is-policy-not-mechanism, granularity trade-off).
  - *Application in MAgHARCM*: Anchors `[[1.0.0 PRIM-17]]` (Asynchronous SE Agent Blackboard) with the practitioner-level companion to the canonical Nii (1986) academic survey.

- **Curtis, Kellner & Over (1992) — *Process Modeling* (`[[1.0.0 P-45]]`)**:
  - *Core Insight*: Survey of software-process modelling organised around four orthogonal dimensions (notation, enactment, content, use). Empirical finding: organisations with multi-aspect, enacted process models exhibit lower defect rates. Introduces Role-Activity Diagram (RAD) where role assignment is a first-class modelling primitive.
  - *Application in MAgHARCM*: Anchors `[[1.0.0 PRIM-24]]` (SOP-Anchored Role-Artifact Schema). SOP definition language is a modern descendant of RAD; each SOP node declares role, artefact type, input/output schema, transition condition.

---
```

Every stage produces strictly typed, schema-versioned artifacts (`[[1.0.0 PRIM-24]]`), ensuring zero information loss and enabling automated backtracking across iterations.
