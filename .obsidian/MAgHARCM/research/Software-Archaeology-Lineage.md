---
title: Software Archaeology & Legacy System Modernization Lineage
date: 2026-09-28
backlink: [[2.0.0 Software-Archaeology-Lineage]]
last_updated: 2026-09-28
tags: [research, software-archaeology, modernization, lineage, synthesis, [[2.0.0 MAgHARCM]], "[[1.0.0 P-122]]", "[[1.0.0 P-123]]", "[[1.0.0 P-124]]", wave-16, wave-17]
---

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

- **Baldwin & Clark ([[1.0.0 P-34]]) — *Design Rules: The Power of Modularity*** (`[[1.0.0 P-34]]`):
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
- **Ted J. Biggerstaff, Bharat G. Mitbander & Dallas E. Webster (1993) ([[1.0.0 P-74]]) & Keith H. Bennett (1995) ([[1.0.0 P-75]]) & Bennett & Rajlich (2000) ([[1.0.0 P-76]]) — *Concept Assignment & Legacy Coping***:
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
- **Michael [[1.0.0 P-24 Feathers]] (2004) — *Working Effectively with Legacy Code***:
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
| `[[1.0.0 PRIM-1]]` | Reverse Topological Ordering | ReCodeAgent [[1.0.0 P-01]]; [[1.0.0 P-86]] LLM-empowered modernization survey; [[1.0.0 P-101]] Least-to-Most Prompting | Parnas (1972) Information Hiding; modern LLM-empowered modernization taxonomy; LtM chained-decomposition template | `internal/agents/planning.go` |
| `[[1.0.0 PRIM-2]]` | Back-Edge Cycle Linearization | AlphaTrans [[1.0.0 P-02]] | Tarjan DAG decomposition | `internal/agents/planning.go` |
| `[[1.0.0 PRIM-3]]` | Target Skeleton-First Gen | Skel [[1.0.0 P-03]], ReCodeAgent; [[1.0.0 P-54]] Phi-3 Tech Report; [[1.0.0 P-88]] HiTyper type-annotation migration; [[1.0.0 P-89]] legacy-modernization baseline | Baldwin & Clark ([[1.0.0 P-34]]) Design Rules; [[1.0.0 P-42]] DSM ([[1.0.0 P-42]]); [[1.0.0 P-53]] Lost in the Middle; HiTyper mixed-method type inference as skeleton input | `internal/agents/planning.go` |
| `[[1.0.0 PRIM-4]]` | SpecMiner Dynamic Invariants | Syzygy [[1.0.0 P-14]], Daikon | Dynamic Invariant Detection (Ernst) | `internal/agents/specminer.go` |
| `[[1.0.0 PRIM-5]]` | Test Co-Translation & Synth | Pynguin [[1.0.0 P-23]], ReCodeAgent; [[1.0.0 P-56]] Zan ([[1.0.0 P-56]]) Multi-SWE-bench; [[1.0.0 P-109]] SWE-bench Verified; [[1.0.0 P-111]] SWE-bench original (Jimenez ICLR 2024); [[1.0.0 P-112]] SWE-agent (Yang 2024); [[1.0.0 P-118]] SWE-bench Lite (Jimenez 2024) | [[1.0.0 P-24 Feathers]] (2004) Characterization Tests; multilingual SLM baseline (Java/Go/Rust/C); verified real-GitHub-issue benchmark; original SWE-bench evaluation benchmark; tool-calling agent scaffold for SWE-bench resolution; leaner 300-instance curated subset of SWE-bench for faster SLM-era iteration | `internal/agents/validator.go` |
| `[[1.0.0 PRIM-6]]` | Multi-Stage Build/Test Repair | AlphaTrans [[1.0.0 P-02]], ReCodeAgent; [[1.0.0 P-54]] Phi-3; [[1.0.0 P-109]] SWE-bench Verified; [[1.0.0 P-111]] SWE-bench original (Jimenez ICLR 2024); [[1.0.0 P-112]] SWE-agent (Yang 2024); [[1.0.0 P-118]] SWE-bench Lite (Jimenez 2024) | Automated Program Repair (Le Goues); Phi-3 chat-format alignment for repair prompts; verified real-GitHub-issue repair benchmark; original SWE-bench evaluation benchmark; tool-calling agent scaffold for SWE-bench resolution; leaner 300-instance curated subset of SWE-bench for faster SLM-era iteration | `internal/agents/validator.go` |
| `[[1.0.0 PRIM-7]]` | Multi-Agent Verdict Validation | MatchFixAgent [[1.0.0 P-08]]; [[1.0.0 P-57]] Speculative Decoding; [[1.0.0 P-78]] EAGLE-3; [[1.0.0 P-92]] Lightman process reward model; [[1.0.0 P-95]] Code Llama verifier baseline; [[1.0.0 P-97]] Welleck self-correct; [[1.0.0 P-98]] Large Language Monkeys; [[1.0.0 P-108]] EAGLE-3 training-time-test draft model (SLM speculative decoding path); [[1.0.0 P-114]] Medusa (Cai 2024) | N-Version Programming (Avizienis); [[1.0.0 P-52]] Self-Consistency; [[1.0.0 P-83]] code-specialised self-consistency; speculative draft/target for SLM fleet; step-by-step verifier as PRM critic; trained self-correction alternative to 3-voter panel; EAGLE-3 training-time-test draft model = 4B-target SLM speculative decoding path; Medusa multi-head drafting as parallel-voter-equivalent verdict panel | `internal/agents/verdict_panel.go` |
| `[[1.0.0 PRIM-8]]` | State-Grounded Mock Validation | TRAM [[1.0.0 P-10]] | [[1.0.0 P-24 Feathers]] (2004) Sensing Pins | `internal/agents/mock_validator.go` |
| `[[1.0.0 PRIM-9]]` | Tri-Representation Code Graph | RepoGraph [[1.0.0 P-14]], Yamaguchi; [[1.0.0 P-110]] GraphCoder / CodeGraphRAG; [[1.0.0 P-113]] AutoCodeRover (Zhang 2024); [[1.0.0 P-116]] Aider (Gauthier 2024-2025); [[1.0.0 P-117]] RepoCoder (Zhang ICLR 2023) | Code Property Graphs (Yamaguchi 2014); graph-RAG for code retrieval with LLM-guided subgraph extraction; AutoCodeRover retrieval+synthesis via code-property-graph traversal; Aider tree-sitter-derived repo-map with call-graph + definition + reference edges; RepoCoder iterative retrieval-augmented repository-level completion loop | `internal/agents/cpg.go` |
| `[[1.0.0 PRIM-10]]` | Feature-Mapping Validation | Oxidizer [[1.0.0 P-05]], RustRepoTrans; [[1.0.0 P-43]] FODA (Kang 1990) | Language Idiom Mapping (Czarnecki) | `internal/agents/feature_mapping.go` |
| `[[1.0.0 PRIM-11]]` | Implementation-Agnostic Test | RepoMod-Bench [[1.0.0 P-11]]; [[1.0.0 P-56]] Zan ([[1.0.0 P-56]]) Multi-SWE-bench | Black-Box Specification Testing; multilingual real-issue benchmark | `internal/agents/impl_agnostic.go` |
| `[[1.0.0 PRIM-12]]` | Wasm Reference Oracle | VERT [[1.0.0 P-12]]; [[1.0.0 P-56]] Zan ([[1.0.0 P-56]]) Multi-SWE-bench | Differential Execution Oracles (McKeeman); per-language differential baseline | `internal/agents/wasm_oracle.go` |
| `[[1.0.0 PRIM-13]]` | Adversarial Test Guard | AdvTestGen [[1.0.0 P-25]]; [[1.0.0 P-56]] Zan ([[1.0.0 P-56]]) Multi-SWE-bench | Mutation Testing & Assertion Invariants, [[1.0.0 P-48]] Jia & [[1.0.0 P-48]] Mutation Survey; multilingual adversarial benchmark | `internal/agents/validator.go` |
| `[[1.0.0 PRIM-14]]` | Software-Archaeology Stage | pp-besm, AgentPatterns.ai; [[1.0.0 P-54]] Phi-3; [[1.0.0 P-79]] Wilde-Scully dynamic feature-location; [[1.0.0 P-87]] TOSEM SLR on LLM4SE | Chikofsky & Cross ([[1.0.0 P-33]]) Reverse Eng, [[1.0.0 P-46]] Seacord et al. ([[1.0.0 P-46]]) Modernizing Legacy Systems; Phi-3-mini-128K whole-file archaeology; dynamic-trace reconnaissance as static-fallback; SLR gap analysis for SLM coverage | `internal/agents/archaeology.go` |
| `[[1.0.0 PRIM-15]]` | Evidence-First Adaptation | Reeper [[1.0.0 P-16]] | Cleanroom Software Engineering (Mills) | `internal/agents/evidence_adaptation.go` |
| `[[1.0.0 PRIM-16]]` | Spec-Driven Dev Lifecycle | spec-kit [[1.0.0 P-17]] | Design-by-Contract (Meyer 1988) | `internal/agents/spec_lifecycle.go` |
| `[[1.0.0 PRIM-17]]` | Asynchronous SE Blackboard | CAID [[1.0.0 P-18]]; [[1.0.0 P-44]] Corkill (1991) | Blackboard Architecture (Nii 1986) | `internal/agents/blackboard.go` |
| `[[1.0.0 PRIM-18]]` | Jaccard-Coupling Recovery | MSR4SA [[1.0.0 P-19]] | Mining Software Repositories (Hassan), [[1.0.0 P-47]] Gall Hajek Jazayeri (1998) Logical Coupling | `internal/agents/jaccard_coupling.go` |
| `[[1.0.0 PRIM-19]]` | Design Rule Hierarchy Part | Kazman et al. [[1.0.0 P-20]], Baldwin-Clark | Baldwin & Clark ([[1.0.0 P-34]]) Modularity, [[1.0.0 P-41]] 2024 retrospective, [[1.0.0 P-42]] MacCormack ([[1.0.0 P-42]]) DSM | `internal/agents/design_rule_hierarchy.go` |
| `[[1.0.0 PRIM-20]]` | Concept Assignment & Redoc | Rajlich (1997) [[1.0.0 P-28]]; [[1.0.0 P-79]] Wilde-Scully Software Reconnaissance | Concept Assignment (Biggerstaff 1993 [[1.0.0 P-74]]); static-first + dynamic-trace-set-difference fallback | `internal/agents/concept_assignment.go` |
| `[[1.0.0 PRIM-21]]` | Migration Strategy Selection | Müller et al. [[1.0.0 P-06]]; [[1.0.0 P-57]] Speculative Decoding; [[1.0.0 P-84]] s1 test-time scaling; [[1.0.0 P-91]] Snell test-time compute allocation; [[1.0.0 P-98]] Large Language Monkeys; [[1.0.0 P-108]] EAGLE-3 training-time-test draft model (SLM speculative decoding path); [[1.0.0 P-114]] Medusa (Cai 2024) | Legacy Migration Frameworks (Müller); [[1.0.0 P-52]] Self-Consistency (strategy-level ensemble); speculative decoding strategy registry; wait-token budget per strategy attempt; compute-optimal allocation across strategies; sampling + verifier as cheapest strategy; EAGLE-3 training-time-test draft = 4B-target SLM speculative path; Medusa multi-head drafting as alternative strategy to EAGLE-3 | `internal/agents/strategy.go` |
| `[[1.0.0 PRIM-22]]` | Four Phases Comprehension | Foltz ([[1.0.0 P-40]]) [[1.0.0 P-30]]; canonical hop-1 [[1.0.0 P-40]]; [[1.0.0 P-51]] Cassano ([[1.0.0 P-51]]) Can It Edit; [[1.0.0 P-54]] Abdin ([[1.0.0 P-54]]) Phi-3; [[1.0.0 P-55]] Schick ([[1.0.0 P-55]]) SLM Few-Shot; [[1.0.0 P-56]] Zan ([[1.0.0 P-56]]) Multi-SWE-bench; [[1.0.0 P-80]] StreamingLLM (long-context SLM mitigation); [[1.0.0 P-90]] Wei CoT; [[1.0.0 P-94]] LIMA curation hypothesis; [[1.0.0 P-96]] Zero-Shot CoT; [[1.0.0 P-99]] BIG-Bench Hard; [[1.0.0 P-100]] Decomposed Prompting (SLM multi-agent); [[1.0.0 P-110]] GraphCoder / CodeGraphRAG; [[1.0.0 P-113]] AutoCodeRover (Zhang 2024); [[1.0.0 P-114]] Medusa (Cai 2024); [[1.0.0 P-115]] OpenHands/CodeAct (Wang 2024); [[1.0.0 P-116]] Aider (Gauthier 2024-2025); [[1.0.0 P-117]] RepoCoder (Zhang ICLR 2023); [[1.0.0 P-118]] SWE-bench Lite (Jimenez 2024); [[1.0.0 P-119]] SWE-Rebench (Badertdinov NeurIPS 2025); [[1.0.0 P-120]] SWE-smith (Yang NeurIPS 2025 spotlight); [[1.0.0 P-121]] BFCL (Patil ICML 2025) | DR. JONES cognitive traversal model (Foltz 2023 [[1.0.0 P-40]]) — Decomposition + Recognition + Organisation + Navigation + Explanation + Search; SLM-aware comprehension substrate via StreamingLLM attention sinks + Phi-3 chat-format alignment + CoT trace decomposition; CodeAct action format simplifies ACI for SLM tool-call generation; Aider repo-map as first-class prompt component; RepoCoder retrieve-then-regenerate loop; SWE-bench Lite per-instance comprehension-iteration budget; SWE-Rebench per-instance provenance gate as Search-phase contamination filter; SWE-smith procedural AST mutations as bug-topology-aware Anchoring phase; BFCL AST-based per-turn contract check for tool-call comprehension | `internal/agents/comprehension.go` |
| `[[1.0.0 PRIM-23]]` | Chunked Translation | ChatDev [[1.0.0 P-12]], MetaGPT [[1.0.0 P-11]]; [[1.0.0 P-53]] Lost in the Middle; [[1.0.0 P-101]] Least-to-Most Prompting; [[1.0.0 P-111]] SWE-bench original (Jimenez ICLR 2024); [[1.0.0 P-113]] AutoCodeRover (Zhang 2024); [[1.0.0 P-120]] SWE-smith (Yang NeurIPS 2025 spotlight) | Bounded-Context Translation; LtM chained-decomposition for cross-chunk state propagation; original SWE-bench multi-file patch chunking as evaluation signal; AutoCodeRover patch chunking as synthesis unit; SWE-smith AST-rewrite strategies as adversarial test-synthesis training-corpus analogue | `internal/agents/chunked_translator.go` |
| `[[1.0.0 PRIM-24]]` | SOP-Anchored Role Artifact | MetaGPT [[1.0.0 P-11]]; [[1.0.0 P-45]] Curtis-Kellner-Over ([[1.0.0 P-45]]); [[1.0.0 P-55]] Schick ([[1.0.0 P-55]]) SLM Few-Shot; [[1.0.0 P-93]] DPO alignment signal; [[1.0.0 P-95]] Code Llama instruction tuning; [[1.0.0 P-100]] Decomposed Prompting (SLM multi-agent) | Standard Operating Procedures (SOP); cloze reformulation as artifact schema; DPO as direct alignment signal for SLM artifact schemas; Code Llama instruction-tuning recipe; decomposed prompting = typed I/O specialist modules | `internal/compiletime/compiletime.go`, `internal/agents/state.go` |
| `[[1.0.0 PRIM-25]]` | Role-Flip De-Hallucination | ChatDev [[1.0.0 P-12]]; sycophancy anchor [[1.0.0 P-38]] ([[1.0.0 P-38]] Raman); [[1.0.0 P-51]] Cassano ([[1.0.0 P-51]]) Can It Edit; [[1.0.0 P-55]] Schick ([[1.0.0 P-55]]) SLM Few-Shot; [[1.0.0 P-81]] Gorilla tool-use at SLM scale; [[1.0.0 P-85]] function-calling SLM gap; [[1.0.0 P-112]] SWE-agent (Yang 2024); [[1.0.0 P-115]] OpenHands/CodeAct (Wang 2024) | Adversarial Verification (Sycophancy Gate); SLM few-shot role-flip prompt pattern; function-call accuracy as adversarial probe; SWE-agent agent-environment role separation as communicative de-hallucination scaffold; OpenHands/CodeAct agent-sandbox role separation as communicative de-hallucination scaffold (multi-turn CodeAct host-side sandbox enforces timeouts, resource limits, typed return schema) | `internal/agents/roleflip.go` |
| `[[1.0.0 PRIM-26]]` | Symbol-Aware Navigator | HyperAgent [[1.0.0 P-13]], ABCoder; [[1.0.0 P-50]] Jiang ([[1.0.0 P-50]]) Code-Gen Survey; [[1.0.0 P-110]] GraphCoder / CodeGraphRAG; [[1.0.0 P-113]] AutoCodeRover (Zhang 2024); [[1.0.0 P-116]] Aider (Gauthier 2024-2025); [[1.0.0 P-117]] RepoCoder (Zhang ICLR 2023) | Targeted Context Retrieval; graph-RAG subgraph ranking for symbol-level navigation; AutoCodeRover symbol-navigator as retrieval+synthesis loop driver; Aider tree-sitter-derived repo-map (call-graph + definition + reference edges) as first-class symbol-level navigation prompt component; RepoCoder iterative retrieve-similar-code-regenerate-completion loop as symbol-aware navigation driver | `internal/agents/navigator.go` |
| `[[1.0.0 PRIM-27]]` | Coverage-Guided Plateau Det | CodaMOSA [[1.0.0 P-23]]; [[1.0.0 P-109]] SWE-bench Verified; [[1.0.0 P-111]] SWE-bench original (Jimenez ICLR 2024); [[1.0.0 P-118]] SWE-bench Lite (Jimenez 2024); [[1.0.0 P-119]] SWE-Rebench (Badertdinov NeurIPS 2025); [[1.0.0 P-120]] SWE-smith (Yang NeurIPS 2025 spotlight) | Search-Based Software Testing (Harman); verified real-GitHub-issue plateau benchmark; original SWE-bench coverage plateau evaluation signal; leaner 300-instance curated subset of SWE-bench that enables iteration budgets the full 2,294-instance benchmark cannot; continuously-evolving decontamination substrate distinguishes "model has plateaued" from "model has memorised"; environment-first synthetic-task injection as adaptive difficulty regulator | `internal/agents/plateau.go` |
| `[[1.0.0 PRIM-28]]` | Conversable Checkpoints | AutoGen [[1.0.0 P-26]] | State-Snapshotting & Resumption | `internal/agents/checkpoint.go` |
| `[[1.0.0 PRIM-29]]` | Recruitment-Adaptive Plan | AgentVerse [[1.0.0 P-26]]; controller-expert [[1.0.0 P-39]] (HuggingGPT / Jarvis); [[1.0.0 P-112]] SWE-agent (Yang 2024); [[1.0.0 P-115]] OpenHands/CodeAct (Wang 2024); [[1.0.0 P-121]] BFCL (Patil ICML 2025) | Dynamic Team Organization, [[1.0.0 P-49]] Shehory & Kraus (1998) Coalition Formation; SWE-agent LM-tool controller-expert recruitment as adaptive plan; OpenHands/CodeAct unified multi-turn CodeAct loop with persistent state as dynamic recruitment scaffold; BFCL multi-turn stateful call pattern as tool-return-driven dynamic re-recruitment | `internal/agents/recruit.go` |
| `[[1.0.0 PRIM-30]]` | Source-to-Target Manifest | Syzygy [[1.0.0 P-14]], JavaC2Rust | Dependency Graph Transpilation | `internal/agents/manifest_rewriter.go` |
| `[[1.0.0 PRIM-31]]` | Iterative Retrieval Refine | RepoCoder [[1.0.0 P-13]]; [[1.0.0 P-53]] Lost in the Middle; [[1.0.0 P-56]] Zan ([[1.0.0 P-56]]) Multi-SWE-bench; [[1.0.0 P-57]] Leviathan ([[1.0.0 P-57]]) Speculative Decoding; [[1.0.0 P-78]] EAGLE-3 multi-layer fusion; [[1.0.0 P-80]] StreamingLLM attention sinks; [[1.0.0 P-108]] EAGLE-3 training-time-test draft model (SLM speculative decoding path); [[1.0.0 P-110]] GraphCoder / CodeGraphRAG; [[1.0.0 P-114]] Medusa (Cai 2024); [[1.0.0 P-115]] OpenHands/CodeAct (Wang 2024); [[1.0.0 P-116]] Aider (Gauthier 2024-2025); [[1.0.0 P-117]] RepoCoder (Zhang ICLR 2023); [[1.0.0 P-119]] SWE-Rebench (Badertdinov NeurIPS 2025); [[1.0.0 P-120]] SWE-smith (Yang NeurIPS 2025 spotlight); [[1.0.0 P-121]] BFCL (Patil ICML 2025) | Dynamic Feedback Retrieval-Augmented Gen; speculative decoding for retrieval-decoding; graph-RAG retrieval augmentation; Medusa multi-head draftin…; SWE-Rebench provenance-driven retrieval exclusion complements feedback-driven refinement; SWE-smith PR-mirrored ground-truth fixes as realistic-retrieval training substrate; BFCL call-pattern templates (serial/parallel/multi-turn) for retrieval-refinement training | `internal/agents/iter_retrieval.go` |


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

- **MacCormack, Rusnak & Baldwin ([[1.0.0 P-42]]) — *Exploring the Structure of Complex Software Designs (DSM)* (`[[1.0.0 P-42]]`)**:
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

---

## 6. Wave-15 Deferral (2026-09-26)

3 candidates (P-122 SWE-Rebench V2, P-123 SWE-bench Multimodal, P-124 SWE-Bench Verified Reference Harness) were held against the §7 trigger gate rewritten 2026-09-25 and rejected. Full rationale per candidate in `.obsidian/MAgHARCM/research/diary/Wave-15-Candidates.md`. No `P-NN` anchor added this wave; wave-15 is deferred.

## 7. Wave-16 SLM-Era Anchors (2026-09-27)

Wave-16 fired 2026-09-27 with 1 new SLM-era mechanism paper persisted.

### 7.1 [[1.0.0 P-122]] ReasoningBank (Zhang et al. 2026, ICLR 2026)

Google Research; **strategy-distilled persistent memory** as the canonical substrate for `[[1.0.0 PRIM-31]]` Iterative Retrieval, `[[1.0.0 PRIM-29]]` Recruiter, `[[1.0.0 PRIM-21]]` Migration Strategy Selection.

- **Hop-1 (cited by ReasoningBank)**:
  - `[[1.0.0 P-90]]` Wei et al. 2022 — Chain-of-Thought (self-judge step).
  - `[[1.0.0 P-111]]` Jimenez et al. 2024 — SWE-Bench (experimental setup).
  - Schick et al. 2023 — Toolformer (analog of structured-tool memory).
- **Hop-2**:
  - `[[1.0.0 P-109]]` SWE-bench Verified (OpenAI 2024).
  - `[[1.0.0 P-115]]` OpenHands (Wang 2024).
- **Lineage position**: closes the persistent-memory gap that `PRIM-31` Iterative Retrieval Refinement assumed but did not implement. Adds the **time-axis dimension** to retrieval (strategy persistence across runs vs LM-feedback within a run).



### 7.2 Wave-16 Substrate-Cross-Reference Matrix

| Primitive | ReasoningBank (P-122) |
| :--- | :--- |
| `[[1.0.0 PRIM-7]]` Verdict Validation | — |
| `[[1.0.0 PRIM-21]]` Migration Strategy | informed-switching policy |
| `[[1.0.0 PRIM-29]]` Recruiter Agent | MaTTS compute-memory loop |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | strategy-distilled persistent memory |

### 7.3 Wave-16 Deferral Note

2 candidates REJECTED at the §7 trigger gate:
- **P-124 SWE-Bench Pro** (ICML 2026) — benchmark, not mechanism; re-evaluation fires when a method-level companion lands.
- **P-125 CodeClash** (ICML 2026) — benchmark, not mechanism; re-evaluation fires when a method-level companion lands.

On 2026-09-26 the wave-15 trigger-criterion evaluation was carried out against the `Methodology.md` §7 gate rewritten 2026-09-25 ("wave-N+1 fires when a 2025+ NeurIPS/ICML/ICLR paper introduces an unanchored mechanism that defends or refutes an existing SLM-era primitive's substrate claim"). Three candidates were drafted and the gate rejected all three. Full deferral memo with per-candidate verdicts is preserved on file at `.obsidian/MAgHARCM/research/diary/Wave-15-Candidates.md`.

**Trigger-gate verdicts:**

- **P-122 SWE-Rebench V2** (Badertdinov et al. 2026) — **REJECTED**. arXiv:2602.23866 is a forward-reference id flagged `[INFERENCE]` in the existing `[[1.0.0 P-119]]` SWE-Rebench V1 note; no confirmed NeurIPS/ICML/ICLR venue publication at note-creation time. Not a 2025+ venue paper.
- **P-123 SWE-bench Multimodal** (October 2024 announcement per swebench.com) — **REJECTED**. No peer-reviewed venue; no arXiv id; announcement-only. Not a 2025+ venue paper.
- **P-124 SWE-bench Verified Reference Harness** (OpenAI August 2024) — **REJECTED**. The harness is a component of `[[1.0.0 P-109]]` SWE-bench Verified (OpenAI 2024), not a new mechanism. Already anchored via P-109.

**Why no `P-NN` anchor was added**: the wave-15 trigger criterion was not met this sprint; none of the three candidates survive the 2025+ venue gate, and none introduces an unanchored mechanism that defends or refutes an existing SLM-era primitive's substrate claim. Adding a `P-NN` anchor under these conditions would violate `Methodology.md` §7 and reintroduce the wave-9 saturation problem that motivated the gate rewrite.

**Re-evaluation trigger**: any of the three candidates becomes wave-15-eligible when (a) P-122 lands at a confirmed 2025+ NeurIPS/ICML/ICLR venue, (b) P-123 lands at a peer-reviewed venue with a verifiable arXiv id, or (c) the user explicitly requests the P-124 component-level anchor. Until then, **no `[[1.0.0 P-122]]` / `[[1.0.0 P-123]]` / `[[1.0.0 P-124]]` anchors are added to any lineage row**, and the P-122 prose body produced during the wave-15 foundation research is preserved verbatim in the deferred-candidates memo as research-on-file.

**Why this lineage note does not cite P-122 even as a hop-2 anchor**: the deferred-candidates memo explicitly bans cross-linking `[[1.0.0 P-122]]` into `Software-Archaeology-Lineage.md`, `docs/.paper/refs.bib`, or `docs/.paper/sec_method.tex` until the candidate is promoted to wave-15. The ban is enforced by the new `scripts/lint_vault.sh` ADR-V-001 automation (which would report a stray `(P-NN)` / `(PRIM-NN)` STRAY hit if anyone attempted to slip the citation in via parenthetical form).

**Cross-reference**:

- Deferred-candidates memo: `.obsidian/MAgHARCM/research/diary/Wave-15-Candidates.md`
- Trigger-gate language: `Methodology.md` §7 (rewritten 2026-09-25)
- Lint automation: `.obsidian/MAgHARCM/architecture/ADR-2026-09-26-Vault-Lint-Extension.md` + `scripts/lint_vault.sh`
- Sprint-2026-09-26 audit block: `.obsidian/MAgHARCM/primitives/Primitives-Index.md` §Sprint 2026-09-26

## 8. Wave-17 SLM-Era Anchors (2026-09-28)

Wave-17 fired 2026-09-28 with 2 new SLM-era mechanism papers persisted.

### 8.1 [[1.0.0 P-123]] CodeChemist (Wang et al. 2026, ICML 2026)

University of Illinois Urbana-Champaign; **training-free test-time scaling for low-resource code generation via functional knowledge transfer**. The mechanism is multi-temperature hedged sampling + cross-lingual I/O test oracle (functional knowledge transfer from a high-resource reference language into the low-resource target). Demonstrated on Qwen-1.5B with 60-70% relative gains on Lua. Anchors `[[1.0.0 PRIM-21]]` Migration Strategy Selection, `[[1.0.0 PRIM-23]]` Chunked Translation, `[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection.

**Substrate implications**:
- `[[1.0.0 PRIM-21]]`: confidence-gated switching between in-language majority voting (cheap) and cross-lingual I/O test oracle (expensive but cross-language functional). Replaces blind try-and-fail with informed switching that the oracle can verify.
- `[[1.0.0 PRIM-23]]`: multi-temperature hedged sampling (k candidate drafts at varied temperatures) with cross-language functional verification. The oracle rejects drafts that fail I/O equivalence tests.
- `[[1.0.0 PRIM-27]]`: functional-coverage plateau via I/O oracle — coverage is measured by oracle acceptance rate, replacing frontier-model judges.

**SLM relevance**: demonstrated on Qwen-1.5B. The oracle is the key enabler: it removes the requirement for a frontier model to judge SLM outputs, making 4B-30B SLM deployment feasible for the chunked-translation pipeline.

### 8.2 [[1.0.0 P-124]] Syzygy (Shetty et al. 2025, ICLR 2025 VerifAI Workshop)

CMU + UCSD; **dual code-test C-to-safe-Rust translation via LLMs + dynamic analysis**. Clang/LLVM-instrumented SpecMiner mines type/bounds/nullability/aliasing properties at runtime; LLM generates Rust code AND equivalence test per translation unit; multi-round repair loop. Anchors `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph, `[[1.0.0 PRIM-22]]` Four Phases of Comprehension, `[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter. Complements `[[1.0.0 P-88]]` HiTyper's static Type Dependency Graph with the **dynamic-analysis half** — together they close the type/safety substrate for legacy modernization.

**Substrate implications**:
- `[[1.0.0 PRIM-9]]`: the tri-representation (AST + CFG + DFG) gains a fourth representation — runtime-mined properties (aliasing, bounds, nullability) — that capture information syntactically hidden from static analysis.
- `[[1.0.0 PRIM-22]]`: comprehension now combines static TDG (`[[1.0.0 P-88]]` HiTyper) + dynamic property mining (`[[1.0.0 P-124]]` Syzygy) as complementary dimensions. Pure-static comprehension leaves type/bounds/nullability under-specified; pure-dynamic comprehension is per-incomplete; the union is full type/safety reasoning.
- `[[1.0.0 PRIM-30]]`: manifests become type/bounds/nullability-enriched — necessary for safe-Rust generation, not just functionally-equivalent translation.

**SLM relevance**: bounded per-translation-unit prompt budget (signature + SpecMiner properties + tests). SpecMiner is model-agnostic and grounds SLM type/safety reasoning in real program behavior, mitigating the SLM hallucination analogue that HiTyper documented for Python type inference.

### 8.3 Wave-17 Substrate-Cross-Reference Matrix

| Primitive | CodeChemist (P-123) | Syzygy (P-124) |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | — | runtime-mined 4th rep |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | confidence-gated oracle switching | — |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | — | static TDG + dynamic mining |
| `[[1.0.0 PRIM-23]]` Chunked Translation | multi-temp + functional verification | — |
| `[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection | functional-coverage plateau via oracle | — |
| `[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter | — | enriched manifests (type/bounds/null) |

### 8.4 Wave-17 Deferral Note

3 candidates REJECTED at the §7 trigger gate:
- **MemSearcher** (Yuan et al. 2026, ACL 2026 Findings) — REJECTED Q1. ACL 2026 Findings is off the trigger-gate's approved venue list (2025+ NeurIPS/ICML/ICLR). Withdrawn from ICLR 2026. Watchlist candidate: fires if a method-level companion lands at NeurIPS 2026 / ICML 2027 / ICLR 2027.
- **Verified Tool Calls** (Mansoor et al. 2026, arXiv:2608.02645) — REJECTED Q1. Strongest Q3 anchor to `[[1.0.0 PRIM-7]]` in the candidate class (postcondition verification + idempotency + verify-before-retry), but arXiv-only with no confirmed NeurIPS/ICML/ICLR venue. ToolACE (ICLR 2025) and CoSC (ICLR 2025) reviewed but fail Q2 (data-QC pipeline / execution feedback, not tool-call verdict validation).
- **LLM-IR / program-comprehension candidates** (Jiang et al. ICML 2025, arXiv:2502.06854) — REJECTED Q2/Q3. The ICML 2025 candidate is a benchmark ("Can LLMs Understand Intermediate Representations in Compilers?"), not a mechanism. No program-comprehension / archaeology mechanism paper at 2025+ ICLR/ICML/NeurIPS found.

**Open gap**: the software archaeology / program-comprehension mechanism slot remains empty at the SLM-era anchor level. The TOSEM SLR (`[[1.0.0 P-87]]` Hou 2024) remains the literature anchor until a 2025+ venue paper emerges.

**Re-evaluation trigger**: any of the three candidates becomes wave-18-eligible when (a) MemSearcher lands at a confirmed 2025+ NeurIPS/ICML/ICLR venue, (b) Verified Tool Calls receives venue confirmation, or (c) a program-comprehension mechanism paper emerges at 2025+ venue. Until then, **no `[[1.0.0 P-125]]` / `[[1.0.0 P-126]]` / `[[1.0.0 P-127]]` anchors are added to any lineage row**, and the deferred-candidate prose bodies are preserved verbatim in `.obsidian/MAgHARCM/research/diary/Wave-17-Candidates.md`.

### 8.5 Watchlist for Wave-18

- **Software archaeology mechanism gap** — active research-on-file. Recommend the next wave-18 scout carry an explicit `program-comprehension-mechanism` query against NeurIPS 2026 / ICML 2027 / ICLR 2027 listings.
- **MemSearcher venue confirmation** — fires if accepted at NeurIPS 2026 main track.
- **Verified Tool Calls venue confirmation** — fires if accepted at NeurIPS 2026 / ICML 2027.
- **MigGPT follow-up** (NeurIPS 2025 spotlight) — patch migration, not type-aware translation; out of scope for this slot.

**Cross-reference**:
- Wave-17 candidates memo: `.obsidian/MAgHARCM/research/diary/Wave-17-Candidates.md`
- Wave-17 paper notes: `.obsidian/MAgHARCM/research/papers/P-123-CodeChemist-ICML-2026.md` + `P-124-Syzygy-ICLR2025-Workshop.md`
- Trigger-gate language: `.obsidian/MAgHARCM/research/Methodology.md` §7 (rewritten 2026-09-25)
- Sprint-2026-09-28 audit block: `.obsidian/MAgHARCM/primitives/Primitives-Index.md` §Sprint 2026-09-28

## 9. Wave-18 SLM-Era Anchors (2026-09-07 iter-2)

9 ACCEPT + 4 REJECT = 13 triaged. Wave-18 closes the SLM-scale verdict + LLM-augmented software-archaeology gap.

### 9.1 Accepted anchors

- `[[1.0.0 P-125]]` **T1: Tool-Integrated Verification for Test-time Compute Scaling in SLMs** (Kang, Jeong & Cho, ICLR 2026) — tool-use as verifier at SLM scale. Anchors `[[1.0.0 PRIM-7]]` Verdict Validation + `[[1.0.0 PRIM-22]]` Four Phases of Comprehension.
- `[[1.0.0 P-126]]` **ARC-Decode: Risk-Bounded Acceptance for Speculative Decoding** (Li et al., ICML 2026) — continuous risk-budget knob. Anchors `[[1.0.0 PRIM-7]]` + `[[1.0.0 PRIM-21]]`.
- `[[1.0.0 P-127]]` **SLM-as-a-Judge for Code Generation** (Crupi et al., ICSE 2026) — replaces frontier PRM. Anchors `[[1.0.0 PRIM-7]]`.
- `[[1.0.0 P-128]]` **KVzip: Query-Agnostic KV Cache Compression with Context Reconstruction** (Kim et al., NeurIPS 2025 Oral) — single-query KV compression. Anchors `[[1.0.0 PRIM-22]]` + `[[1.0.0 PRIM-31]]`.
- `[[1.0.0 P-129]]` **LλMDA: LLM-Aided Partial Program Dependence Analysis** (Rong, Yadavally & Nguyen, ICSE 2026) — LLM-aided partial PDG construction. Anchors `[[1.0.0 PRIM-9]]` + `[[1.0.0 PRIM-22]]`.
- `[[1.0.0 P-130]]` **SSAR: Software Architecture Recovery** (Ding, Mo, Wu & Song, ICSE 2026) — structural+semantic hybrid. Anchors `[[1.0.0 PRIM-9]]` + `[[1.0.0 PRIM-22]]`.
- `[[1.0.0 P-131]]` **SemArc: Software Architecture Recovery Augmented with Semantics** (Zhao et al., TSE 2026) — semantic-enhanced architecture recovery. Anchors `[[1.0.0 PRIM-9]]` + `[[1.0.0 PRIM-22]]`.
- `[[1.0.0 P-132]]` **SemRef: Semantic-Enhanced Refinement of Architecture Recovery** (Zhang et al., ICSE 2026) — iterative semantic refinement. Anchors `[[1.0.0 PRIM-9]]` + `[[1.0.0 PRIM-22]]` + `[[1.0.0 PRIM-31]]`.
- `[[1.0.0 P-133]]` **ADI: Empowering Autonomous Debugging Agents with Efficient Dynamic Analysis** (Xiang et al., FSE 2026 — SIGSOFT Distinguished Paper Award) — Frame Lifetime Trace + function-level interactive debugging. Anchors `[[1.0.0 PRIM-22]]` + `[[1.0.0 PRIM-31]]`.

### 9.2 Rejected (4)

- **R1 ReflexiCoder** (Jiang et al., ACL 2026 Findings) — REJECTED Q1. ACL Findings off-list. Mechanism real (RL-internalised self-reflection at 1.5B-14B) but venue fails Q1.
- **R2 Self-Distillation for Code Generation** (Zhang et al., Apple, arXiv:2604.01193) — REJECTED Q1. arXiv-only; no confirmed peer-reviewed acceptance. Strong mechanism (30B 42.4% → 55.3% pass@1 on LiveCodeBench v6).
- **R3 SPECS** (Cemri et al., arXiv:2506.15733) — REJECTED Q1 (Wave-18). ICLR 2026 submission acceptance unconfirmed at Wave-18 triage. **Resolved in Wave-19**: accepted at ICLR 2026 as P-135.
- **R4 Software-Archaeology strict-mechanism gap** — gap partially closed by Wave-18 P-129..P-133; strict program-comprehension-mechanism slot remains open.

**Cross-reference**:
- Wave-18 candidates memo: `.obsidian/MAgHARCM/research/diary/Wave-18-Candidates.md`
- Wave-18 paper notes: `.obsidian/MAgHARCM/research/papers/P-125..P-133-*.md`
- Sprint 2026-09-07 (iter-2) audit block: `.obsidian/MAgHARCM/primitives/Primitives-Index.md` §Sprint 2026-09-07 (Wave-19)

## 10. Wave-19 SLM-Era Anchors (2026-09-07 iter-2)

8 ACCEPT + 3 REJECT + 2 UNVERIFIED = 13 triaged. Wave-19 closes the speculative-decoding × KV-cache hybrid gap and the SLM-scale test-time-scaling gap.

### 10.1 Accepted anchors

- `[[1.0.0 P-134]]` **RelayCaching: Accelerating LLM Collaboration via Decoding KV Cache Reuse** (Geng et al., ICML 2026 Poster #1915) — cross-agent KV reuse via sparse deviation recompute (>80% reuse, 4.7× TTFT). Anchors `[[1.0.0 PRIM-31]]` + `[[1.0.0 PRIM-9]]`.
- `[[1.0.0 P-135]]` **SPECS: Faster Test-Time Scaling through Speculative Drafts** (Cemri et al., ICLR 2026) — speculative drafts + soft verification + dynamic switch. Closes Wave-18 R3. Anchors `[[1.0.0 PRIM-7]]` + `[[1.0.0 PRIM-21]]`.
- `[[1.0.0 P-136]]` **CaTS: Calibrated Test-Time Scaling for Efficient LLM Reasoning** (Huang et al., ICLR 2026 Poster) — Self-Calibration distilled confidence. Anchors `[[1.0.0 PRIM-7]]` as frontier-PRM replacement.
- `[[1.0.0 P-137]]` **SuffixDecoding: Extreme Speculative Decoding for Emerging AI Applications** (Oliaro et al., NeurIPS 2025 Spotlight) — model-free suffix-tree draft. Deployed in Snowflake ArcticInference + vLLM. Anchors `[[1.0.0 PRIM-7]]` + `[[1.0.0 PRIM-21]]`.
- `[[1.0.0 P-138]]` **RepairKV (Cache You Later)** (Rusli et al., ICML 2026 AdaptFM Workshop) — borderline workshop-track ACCEPT per §7 method-level threshold. Post-compression KV repair runtime operator. Anchors `[[1.0.0 PRIM-22]]` + `[[1.0.0 PRIM-31]]`.
- `[[1.0.0 P-139]]` **TypePro: Boosting LLM-Based Type Inference via Inter-Procedural Slicing** (Lin et al., FSE 2026) — SDG + inter-procedural backward/forward slicing. Anchors `[[1.0.0 PRIM-9]]` + `[[1.0.0 PRIM-22]]`.
- `[[1.0.0 P-140]]` **Panta: LLM Test Generation via Iterative Hybrid Program Analysis** (Gu, Nashid & Mesbah, ICSE 2026) — static (cyclomatic) + dynamic (coverage) iterative loop. Anchors `[[1.0.0 PRIM-22]]` + `[[1.0.0 PRIM-21]]`.
- `[[1.0.0 P-141]]` **KVFlow: Workflow-Aware KV Cache Eviction for Multi-Agent LLM Serving** (NeurIPS 2025 Poster) — Agent Step Graph + steps-to-execution metric + prefetching. Anchors `[[1.0.0 PRIM-31]]` + `[[1.0.0 PRIM-21]]`.

### 10.2 Rejected (3)

- **R1 TTA\*** (Braverman, Zhang & Gu, NeurIPS 2025 LAW Workshop) — REJECTED Q1 (workshop redundancy). P-135 SPECS + P-136 CaTS cover same focus area at ICLR 2026 main track. Mechanism (A* search wrapper) concrete but adds little beyond P-84 s1 beam-search lineage.
- **R2 HELIOS** (Achamyeleh, Thomare & Al Faruque, NDSS 2026 LAST-X Workshop) — REJECTED Q1 (off-list venue + off-axis target). NDSS LAST-X workshop is off-list. Mechanism concrete (CFG+FCG → textual prompt; compilability 45.0%→85.2%) but binary-decompilation target is off-axis for MAgHARCM's source→source translation.
- **R3 LongSpec** (ACL 2026 Main) — REJECTED Q1 (off-list venue). ACL is not in Wave-19 whitelist. Same rationale as Wave-18 R1 ReflexiCoder.

### 10.3 Watchlist — UNVERIFIED (2)

- **W1 SliceMate** (Chang et al., arXiv:2507.18957) — UNVERIFIED. Yunbo Lyu's homepage claims ISSTA 2026 acceptance, but no ISSTA 2026 program slot for SliceMate is listed on conf.researchr.org. Mechanism concrete (three LLM agents replace explicit PDG/SDG construction; 22% acc / 28% F1 improvement). Watchlist for Wave-20 venue re-verification.
- **W2 SWE-TRACE** (Han et al., arXiv:2604.14820) — UNVERIFIED. April 2026 arXiv preprint only; no peer-reviewed venue confirmation. Mechanism concrete (cascaded trajectory optimisation + rubric-PRM + heuristic TTS). Watchlist for Wave-20 venue confirmation.

**Cross-reference**:
- Wave-19 candidates memo: `.obsidian/MAgHARCM/research/diary/Wave-19-Candidates.md`
- Wave-19 paper notes: `.obsidian/MAgHARCM/research/papers/P-134..P-141-*.md`
- REJECT registry: `.obsidian/MAgHARCM/Research-Database.json::reject_registry.wave-19` (BLK-08 resolved)
- Watchlist: `.obsidian/MAgHARCM/Research-Database.json::watchlist.wave-19`

## 11. Wave-20 SLM-Era Anchors (2026-09-07 iter-3)

5 ACCEPT + 4 REJECT + 1 UNVERIFIED = 10 triaged. Wave-20 partially closes the strict program-comprehension-mechanism slot (open since Wave-17) via P-142 NESA self-evolving graph pre-analysis; anchors SLM-grounded hallucination defence (P-143 HalluShield), prompt-trace training data (P-144 TraceCoder), LLM-driven translation strategy selection (P-145 TerraMod), and workflow-aware cross-document PRM (P-146 ContextPRM).

### 11.1 Accepted anchors

- `[[1.0.0 P-142]]` **NESA: Relational Neuro-Symbolic Static Program Analysis** (Li et al., FSE 2026) — self-evolving graph pre-analysis; LLM-aided iterative refinement of the static program graph. **Partially closes the strict program-comprehension-mechanism slot.** Anchors `[[1.0.0 PRIM-9]]` + `[[1.0.0 PRIM-22]]`.
- `[[1.0.0 P-143]]` **HalluShield / Hallucination Detection and Mitigation for LLM-based Code Summarization** (Wang et al., FSE 2026) — SLM-grounded speculative-decoding hallucination defence for 4B-30B substrates. Re-anchors the frontier-PRM-as-judge gap. Anchors `[[1.0.0 PRIM-7]]` + `[[1.0.0 PRIM-21]]`.
- `[[1.0.0 P-144]]` **TraceCoder: A Trace-Driven Multi-Agent Framework for Automated Debugging of LLM-Generated Code** (Zhou et al., ICSE 2026) — prompt-trace training-data construction; preserves successful trajectories as agent supervision. Anchors `[[1.0.0 PRIM-29]]` + `[[1.0.0 PRIM-31]]`.
- `[[1.0.0 P-145]]` **TerraMod: LLM-Driven Lexical-Based Translation Strategy Selection** (Zhang et al., ICSE 2026 NIER) — LLM-driven lexical-based translation strategy selection for legacy migration. Anchors `[[1.0.0 PRIM-21]]`.
- `[[1.0.0 P-146]]` **ContextPRM: Workflow-Aware Cross-Document Process Reward Modelling for Agentic SLMs** (Park et al., ICLR 2026) — workflow-aware cross-document process reward modelling; re-anchors P-91..P-92 frontier-PRM-as-judge gap with workflow context. Anchors `[[1.0.0 PRIM-7]]` + `[[1.0.0 PRIM-31]]`.

### 11.2 Rejected (4)

- **R1 Nexus** (ICSE 2026) — REJECTED Q3 (ablations only, no mechanism). Same pattern as Wave-19 REJECTs: ablations do not introduce a new substrate.
- **R2 SWE-Lego** (ICSE 2026 NIER) — REJECTED Q3 (engineering pattern, no mechanism). NIER track is the right venue but the contribution is a scaffolding pattern, not a mechanism.
- **R3 CoPS** (ICML 2026) — REJECTED Q1 (speculative venue, no venue confirmation). Same rationale as Wave-19 R2 HELIOS off-list handling.
- **R4 SHIELD-ASR** (ACL 2026 Findings) — REJECTED Q1 (off-list venue). ACL is not in the §7 trigger-list.

### 11.3 Watchlist — UNVERIFIED (1)

- **U1 NSE** (ICML 2026 placeholder venue, no DOI / OpenReview / arXiv) — Watchlist for Wave-21 venue re-verification.

### 11.4 Watchlist Resolved (2)

- **W1 SliceMate** — REJECTED. ISSTA 2026 program slot still absent on conf.researchr.org; Yunbo Lyu's homepage claim is unsupported. Removed from watchlist.
- **W2 SWE-TRACE** — REJECTED. April 2026 arXiv preprint only (arXiv:2604.14820); no peer-reviewed venue confirmation. Removed from watchlist.

**Cross-reference**:
- Wave-20 candidates memo: `.obsidian/MAgHARCM/research/diary/Wave-20-Candidates.md`
- Wave-20 paper notes: `.obsidian/MAgHARCM/research/papers/P-142..P-146-*.md`
- REJECT registry: `.obsidian/MAgHARCM/Research-Database.json::reject_registry.wave-20` (BLK-08 standing rule applied)
- Watchlist: `.obsidian/MAgHARCM/Research-Database.json::watchlist.wave-20`
- Sprint 2026-09-07 (iter-3) audit block: `.obsidian/MAgHARCM/diary/Sprint-2026-09-07-Handoff-3.md`
- Methodology §11.6 (NEW): Program-Comprehension Mechanism substrate anchored by P-142 + Wave-18 architecture-recovery trio + P-133 ADI
- Sprint 2026-09-07 (iter-2) audit block: `.obsidian/MAgHARCM/primitives/Primitives-Index.md` §Sprint 2026-09-07 (Wave-19)
