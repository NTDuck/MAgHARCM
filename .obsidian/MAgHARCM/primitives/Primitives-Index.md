---
title: Primitives Index
backlink: "[[2.0.0 Primitives Index]]"
aliases:
  - "2.0.0 Primitives Index"
  - "Primitives Index"
  - "Primitives-Index"
  - "INDEX"
  - "primitives/Primitives-Index"
date: 2026-09-07
last_updated: 2026-09-07
tags: [primitives, catalog, status, [[1.0.0 PRIM-1]]..[[1.0.0 PRIM-31]], "[[1.0.0 P-122]]", "[[1.0.0 P-123]]", "[[1.0.0 P-124]]", "[[1.0.0 P-134]]".."[[1.0.0 P-141]]", wave-16, wave-17, wave-19]
---

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
| `[[1.0.0 PRIM-5]]` | Test Suite Co-Translation & Synthesis | `internal/agents/validator.go::generateAdditionalTests` | [[Pynguin-2021]], [[1.0.0 P-24]] Feathers, [[1.0.0 P-56]] Multi-SWE-bench, [[1.0.0 P-109]] SWE-bench Verified | implemented |
| `[[1.0.0 PRIM-6]]` | Multi-Stage Build/Test Feedback Repair | `internal/agents/validator.go::Run` | [[AlphaTrans-2024]], Le Goues, [[1.0.0 P-54]] Phi-3, [[1.0.0 P-109]] SWE-bench Verified | implemented |
| `[[1.0.0 PRIM-7]]` | Multi-Agent Verdict Validation | `internal/agents/verdict_panel.go` | [[MatchFixAgent-2024]], Avizienis, [[1.0.0 P-57]] Speculative Decoding, [[1.0.0 P-61]] RedCode, [[1.0.0 P-78]] EAGLE-3, [[1.0.0 P-83]] code-specialised self-consistency, [[1.0.0 P-92]] Lightman PRM, [[1.0.0 P-95]] Code Llama verifier baseline, [[1.0.0 P-108]] EAGLE-3 speculative decoding (SLM draft model), [[1.0.0 P-125]] T1 tool-integrated verification (wave-18), [[1.0.0 P-126]] ARC-Decode risk-bounded acceptance (wave-18), [[1.0.0 P-127]] SLM-as-a-Judge (wave-18), [[1.0.0 P-135]] SPECS speculative drafts + soft verification + dynamic switch (wave-19), [[1.0.0 P-136]] CaTS Self-Calibration frontier-PRM replacement (wave-19), [[1.0.0 P-137]] SuffixDecoding model-free suffix-tree draft (wave-19) | implemented |
| `[[1.0.0 PRIM-8]]` | State-Grounded Mock-Based In-Isolation Validation | `internal/agents/mock_validator.go` | [[TRAM-2024]], [[1.0.0 P-24]] Feathers | implemented |
| `[[1.0.0 PRIM-9]]` | Tri-Representation Hybrid Code Graph | `internal/agents/cpg.go` | [[RepoGraph-2024]], Yamaguchi, [[1.0.0 P-110]] GraphCoder / CodeGraphRAG, [[1.0.0 P-124]] Syzygy runtime-mined aliasing/bounds/nullability as fourth representation (wave-17), [[1.0.0 P-129]] LλMDA LLM-aided partial program dependence (wave-18), [[1.0.0 P-130]] SSAR architecture recovery (wave-18), [[1.0.0 P-131]] SemArc semantic-enhanced architecture recovery (wave-18), [[1.0.0 P-132]] SemRef semantic-enhanced refinement (wave-18), [[1.0.0 P-134]] RelayCaching cross-agent prefill reuse on the 8-agent graph (wave-19), [[1.0.0 P-139]] TypePro inter-procedural SDG slicing as fourth representation (wave-19) | implemented |
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
| `[[1.0.0 PRIM-21]]` | Migration Strategy Selection | `internal/agents/strategy.go::Registry.TryInOrder` | [[1.0.0 P-35]] Müller 5 Strategies, [[1.0.0 P-57]] Speculative Decoding, [[1.0.0 P-60]] Bisbal Cascading Fallback, [[1.0.0 P-84]] s1 test-time scaling, [[1.0.0 P-91]] Snell test-time compute allocation, [[1.0.0 P-98]] Large Language Monkeys, [[1.0.0 P-108]] EAGLE-3 speculative decoding (SLM draft model), [[1.0.0 P-122]] ReasoningBank strategy-distilled persistent memory (wave-16), [[1.0.0 P-123]] CodeChemist confidence-gated cross-lingual I/O oracle switching (wave-17), [[1.0.0 P-126]] ARC-Decode risk-bounded acceptance (wave-18), [[1.0.0 P-135]] SPECS continuous draft<->target budget knob (wave-19), [[1.0.0 P-137]] SuffixDecoding model-free suffix-tree draft path (wave-19), [[1.0.0 P-140]] Panta coverage-driven test prioritisation knob (wave-19), [[1.0.0 P-141]] KVFlow agent-activation-time scheduling signal (wave-19) | implemented |
| `[[1.0.0 PRIM-22]]` | Four Phases of Comprehension | `internal/agents/comprehension.go` | [[1.0.0 P-40]] DR. JONES Model, [[1.0.0 P-54]] Phi-3, [[1.0.0 P-55]] SLM Few-Shot, [[1.0.0 P-56]] Multi-SWE-bench, [[1.0.0 P-58]] Qwen2.5-Coder, [[1.0.0 P-62]] Decomposed Prompting, [[1.0.0 P-80]] StreamingLLM, [[1.0.0 P-90]] Wei CoT, [[1.0.0 P-94]] LIMA curation, [[1.0.0 P-96]] Zero-Shot CoT, [[1.0.0 P-99]] BIG-Bench Hard, [[1.0.0 P-100]] Decomposed Prompting (SLM multi-agent), [[1.0.0 P-110]] GraphCoder / CodeGraphRAG, [[1.0.0 P-124]] Syzygy dynamic property mining complements static TDG [[1.0.0 P-88]] HiTyper (wave-17), [[1.0.0 P-128]] KVzip query-agnostic context reconstruction (wave-18), [[1.0.0 P-129]] LλMDA (wave-18), [[1.0.0 P-130]] SSAR (wave-18), [[1.0.0 P-133]] ADI Frame Lifetime Trace (wave-18), [[1.0.0 P-138]] RepairKV post-compression KV repair for long-context agentic retrieval (wave-19), [[1.0.0 P-139]] TypePro inter-procedural type inference for Recognition + Explanation (wave-19), [[1.0.0 P-140]] Panta static+dynamic signal for untested-path search (wave-19) | implemented |
| `[[1.0.0 PRIM-23]]` | Chunked Translation | `internal/agents/chunked_translator.go` | [[ChatDev-2023]], [[MetaGPT-2023]], [[1.0.0 P-58]] Qwen2.5-Coder, [[1.0.0 P-62]] Decomposed Prompting, [[1.0.0 P-65]] PAL, [[1.0.0 P-101]] Least-to-Most Prompting, [[1.0.0 P-123]] CodeChemist multi-temperature hedged sampling + cross-language functional verification (wave-17) | implemented |
| `[[1.0.0 PRIM-24]]` | SOP-Anchored Role-Artifact Schema | `internal/compiletime/compiletime.go`, `internal/agents/state.go` | [[MetaGPT-2023]] SOP Contracts, [[1.0.0 P-45]] Curtis-Kellner-Over, [[1.0.0 P-55]] SLM Few-Shot, [[1.0.0 P-62]] Decomposed Prompting, [[1.0.0 P-63]] Multi-Agent Survey, [[1.0.0 P-93]] DPO alignment, [[1.0.0 P-95]] Code Llama instruction tuning, [[1.0.0 P-100]] Decomposed Prompting (SLM multi-agent) | implemented |
| `[[1.0.0 PRIM-25]]` | Communicative-De-hallucination Role-Flip Gate | `internal/agents/roleflip.go` | [[ChatDev-2023]] Sycophancy Gate, [[1.0.0 P-38]] Anthropic Sycophancy, [[1.0.0 P-55]] SLM Few-Shot, [[1.0.0 P-58]] Qwen2.5-Coder, [[1.0.0 P-61]] RedCode, [[1.0.0 P-81]] Gorilla, [[1.0.0 P-85]] function-calling SLM gap | implemented |
| `[[1.0.0 PRIM-26]]` | Symbol-Aware Navigator | `internal/agents/navigator.go` | [[HyperAgent-2024]], ABCoder, [[1.0.0 P-50]] Code-Gen Survey, [[1.0.0 P-58]] Qwen2.5-Coder, [[1.0.0 P-110]] GraphCoder / CodeGraphRAG | implemented |
| `[[1.0.0 PRIM-27]]` | Coverage-Guided Plateau Detection | `internal/agents/plateau.go` | [[CodaMOSA-2023]], Harman SBST, [[1.0.0 P-109]] SWE-bench Verified, [[1.0.0 P-123]] CodeChemist functional-coverage plateau via cross-language I/O oracle (wave-17) | implemented |
| `[[1.0.0 PRIM-28]]` | Conversable State Checkpoints & Interrupts | `internal/agents/checkpoint.go` | [[AutoGen-2023]], Snapshotting, [[1.0.0 P-59]] Digital Apollo | implemented |
| `[[1.0.0 PRIM-29]]` | Recruitment-Adaptive Planning | `internal/agents/recruit.go` | [[AgentVerse-2023]], Dynamic Org, [[1.0.0 P-49]] Shehory & Kraus Coalition Formation, [[1.0.0 P-63]] Multi-Agent Survey, [[1.0.0 P-122]] ReasoningBank MaTTS compute-memory loop (wave-16) | implemented |
| `[[1.0.0 PRIM-30]]` | Source-to-Target Manifest Rewriter | `internal/agents/manifest_rewriter.go` | [[Syzygy-2024]], JavaC2Rust, [[1.0.0 P-124]] Syzygy type/bounds/nullability-enriched manifests for safe-Rust generation (wave-17) | implemented |
| `[[1.0.0 PRIM-31]]` | Iterative Retrieval Refinement | `internal/agents/iter_retrieval.go` | [[RepoCoder-2024]], Dynamic Context, [[1.0.0 P-56]] Multi-SWE-bench, [[1.0.0 P-57]] Speculative Decoding, [[1.0.0 P-58]] Qwen2.5-Coder, [[1.0.0 P-62]] Decomposed Prompting, [[1.0.0 P-108]] EAGLE-3 speculative decoding (SLM draft model), [[1.0.0 P-110]] GraphCoder / CodeGraphRAG, [[1.0.0 P-122]] ReasoningBank strategy-distilled persistent memory (wave-16), [[1.0.0 P-128]] KVzip (wave-18), [[1.0.0 P-132]] SemRef (wave-18), [[1.0.0 P-133]] ADI (wave-18), [[1.0.0 P-134]] RelayCaching cross-agent prefill reuse (wave-19), [[1.0.0 P-138]] RepairKV second-chance correction mechanism for evicted KV rows (wave-19), [[1.0.0 P-141]] KVFlow workflow-aware eviction for multi-agent retention (wave-19) | implemented |

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
- PR backlog (low-risk follow-ups): re-tag Yamaguchi (`[[1.0.0 PRIM-9]]`) and Nii (`[[1.0.0 PRIM-17]]`) backlinks once their standalone paper notes are materialised.

## Sprint 2026-09-18 Vault Sync Audit

- Indexed 6 wave-10 entries (P-102 SmallCode, P-103 AgentModernize arXiv 2026, P-104 S*, P-105 ChunkKV, P-106 BFCL, P-107 re-anchor of P-100) into relevant primitive rows (PRIM-21, 22, 23, 24, 25, 26, 27, 31).
- P-103 venue corrected: arXiv:2605.17535 (2026), NOT ICSE 2025.

## Sprint 2026-09-19 Vault Sync Audit

- All 31 primitives verified to have implementation files in internal/agents/ (35 files = 31 primitive impl + 1 state (ADR-C-014) + 2 utility (parser, canonical_crates) + 3 test files + 0 orphan).
- Ponytail sweep: zero fmt.Print* / log.Print* / raw panic / os.Stdout in production code; all fmt.* uses are Sprintf/Fprintf to buffers (string construction, not I/O).
- 8-agent graph verified: Archaeologist, Analyzer, Planning, Translator, reviewer (RoleFlipGate), Validator, VerdictPanel, Recruiter. See Architecture.md §4 for the canonical statement.
- Try-and-fail strategy registry confirmed; SelectMigrationStrategy replaced by Registry.TryInOrder.
- Renamed 3 *Default* constants to *Placeholder to align naming with no-fallback rule: TranslatedPackagePlaceholder, ProjectDirPlaceholder, ConceptDescriptionPlaceholder.
- configs/agents.yml now lists `lsp.provider: abcoder-mcp` as canonical example (abcoder MCP default surfaced in docs).
## Sprint 2026-09-20 Vault Sync Audit

- Indexed 3 wave-11 anchor papers: P-108 EAGLE-3 (NeurIPS 2025 — re-anchor of P-78; venue corrected from earlier "2024" attribution), P-109 SWE-bench Verified (OpenAI 2024), P-110 GraphCoder / CodeGraphRAG (graph-RAG for code, 2024).
- All 31 primitives still map to implementation files in `internal/agents/` (35 files = 31 primitive impl + 1 state (ADR-C-014) + 2 utility (parser, canonical_crates) + 3 test files + 0 orphan).
- 8-agent graph still wired; Charm TUI idioms still intact; abcoder-mcp default still in `configs/agents.yml`.
- Wave-11 continued deferral: 5 candidates (EAGLE-3, GraphCoder, MemoryBank-E, TinyRM, SWE-bench Verified 2025) — none introduce a new SLM-era mechanism that requires anchoring.

## Sprint 2026-09-21 Vault Sync Audit

- Indexed 3 wave-11 anchor papers: P-108 EAGLE-3 (NeurIPS 2025 — re-anchor of P-78; venue corrected from earlier "2024" attribution), P-109 SWE-bench Verified (OpenAI 2024), P-110 GraphCoder / CodeGraphRAG (graph-RAG for code, 2024).
- Cross-linked P-108 into PRIM-7, PRIM-21, PRIM-31; P-109 into PRIM-5, PRIM-6, PRIM-27; P-110 into PRIM-9, PRIM-26, PRIM-31.
- Re-verified 31/31 primitives still map to implementation files in `internal/agents/` (35 files = 31 primitive impl + 1 state (ADR-C-014) + 2 utility (parser, canonical_crates) + 3 test files + 0 orphan).
- Wave-11 trigger fired: new SLM-era mechanism (EAGLE-3 training-time-test draft model = 4B-target SLM speculative decoding path).

## Sprint 2026-09-22 Vault Sync Audit

- Indexed 4 wave-12 deep-research papers (hop-1 of wave-11 anchors): P-111 SWE-bench original (Jimenez ICLR 2024), P-112 SWE-agent (Yang 2024), P-113 AutoCodeRover (Zhang 2024), P-114 Medusa (Cai 2024).
- P-111 + P-112 + P-109 form the canonical SWE-bench evaluation lineage (original → Verified → tool-calling agent scaffold).
- P-113 anchors PRIM-9 + PRIM-26 (code-graph + symbol-navigator primitives) with the AutoCodeRover retrieval+synthesis pattern.
- P-114 anchors PRIM-7 + PRIM-21 + PRIM-22 + PRIM-31 (verdict + strategy + comprehension + retrieval) as the alternative drafting strategy to EAGLE-3.
- Cross-linked 16 rows in Software-Archaeology-Lineage.md.
- 31/31 primitives still mapped to implementation files; 8-agent graph still wired; Charm TUI idioms still intact; abcoder-mcp default still in configs/agents.yml.

## Sprint 2026-09-23 Vault Sync Audit

- Indexed 4 wave-13 deep-research papers: P-115 OpenHands/CodeAct (Wang 2024), P-116 Aider (Gauthier 2024-2025), P-117 RepoCoder (Zhang ICLR 2023), P-118 SWE-bench Lite (Jimenez 2024).
- P-115 + P-112 + P-113 + P-116 form the canonical "agent scaffold" cluster for SWE-bench-style issue resolution.
- P-117 + P-110 + P-113 form the canonical "repository-level retrieval-augmented completion" cluster.
- P-118 narrows the SWE-bench evaluation target from 2,294 (`[[1.0.0 P-111]]`) to a 300-instance leaner subset for faster SLM-era iteration.
- Cross-linked 16 rows in Software-Archaeology-Lineage.md.
- 31/31 primitives still mapped to implementation files; 8-agent graph still wired; Charm TUI idioms still intact; abcoder-mcp default still in configs/agents.yml.

## Sprint 2026-09-24 Vault Sync Audit
- ADR-V-001 sweep: 11 stray `(P-NN)` / `(PRIM-NN)` parentheticals rewritten to `[[1.0.0 P-NN]]` / `[[1.0.0 PRIM-NN]]` form across `Methodology.md`, `primitives/Primitives-Index.md`, and `diary/Sprint-2026-09-06-Handoff.md`, `Sprint-2026-09-18-Handoff.md`, `Sprint-2026-09-22-Handoff.md`, `Sprint-2026-09-23-Handoff.md`. ASCII diagram in `Software-Archaeology-Lineage.md:51` retained (diagram context, not version slots); ADR-V-001 rule statement in `architecture/ADR-2026-09-07-Sprint-Conventions.md` retained as the canonical counter-example.
- Primitive-count audit: `grep -oE 'PRIM-[0-9]+' primitives/Primitives-Index.md | sort -u | wc -l` = 31 (no orphan rows, no missing rows).
- Compliance re-verified: 8-agent graph still wired (archaeologist, analyzer, planning, translator, **reviewer** [RoleFlipGate — constructor `agents.NewRoleFlipGate` at `internal/agents/roleflip.go:44`], validator, verdict_panel, recruiter = 8 functional + 2 checkpoint `AddLambdaNode` calls at `graph.go:87,124` = 10 total); `configs/agents.yml` still lists `lsp.provider: abcoder-mcp`; zero `fmt.Print*` I/O in production Go code; **Must pattern** (broader regex `Must[A-Z][a-zA-Z]*\(`) used **40 times repo-root**: archaeology.go:7, iter_retrieval.go:1, parser.go:3, spec_lifecycle.go:1, verdict_panel.go:3, navigator.go:1, validator.go:2, roleflip.go:1, config/yaml.go:2, languages/extractor.go:4, tools/lsp.go:6, tools/exec.go:4, compiletime/state.go:1, compiletime/compiletime.go:1, checkpoint/checkpoint.go:1, cmd/MAgHARCM/main.go:1, tests/internal/tools/tools.go:1. Narrow regex `Must\(|MustNew\(|MustNotNil\(|MustNotEmpty\(|MustTask\(` returns 4 (spec_lifecycle.go:1, compiletime/state.go:2, compiletime/compiletime.go:1); all 31 primitive artifacts still in `internal/agents/*.go` (35 files = 31 primitive impl + 1 state (ADR-C-014) + 2 utility (parser, canonical_crates) + 3 test files + 0 orphan).
- Wave-14 deferred: no new SLM-era mechanism landed this sprint; wave-14 trigger criterion (`new SLM-era primitive OR 2026 venue paper introduces unanchored mechanism`) not satisfied.
- All gates green: `go build ./...`, `go vet ./...`, `go test ./...`.

## Sprint 2026-09-25 Vault Sync Audit
- Indexed 3 wave-14 deep-research papers: P-119 SWE-Rebench (Badertdinov NeurIPS 2025 D&B), P-120 SWE-smith (Yang NeurIPS 2025 D&B spotlight), P-121 BFCL (Patil ICML 2025).
- P-119 anchors PRIM-22 + PRIM-27 + PRIM-31 (decontamination-aware evaluation; complement to existing SWE-bench-family coverage by `[[1.0.0 P-109]]` / `[[1.0.0 P-111]]` / `[[1.0.0 P-118]]` with the **time-axis dimension** the static benchmarks lack).
- P-120 anchors PRIM-22 + PRIM-23 + PRIM-27 + PRIM-31 (environment-first synthetic-task generator; complementary evaluation-data substrate for the open-weights SLM regime).
- P-121 anchors PRIM-22 + PRIM-29 + PRIM-31 (AST-based tool-call evaluation; multi-turn/serial/parallel patterns + cost+latency rubric). P-121 supersedes the earlier wave-10 P-106 BFCL anchor as the canonical tool-calling benchmark reference; P-106 retained alongside P-121 per the §12 stale-directive audit decision (2026-09-26: user did not request retirement).
- P-122 dropped: Qwen3-Coder / SmolLM3 / xLAM-2 are model releases, not mechanism papers; go in hop-2 citations, not as standalone P-NN anchors.
- Cross-links pending for `Software-Archaeology-Lineage.md` PRIM-22, PRIM-23, PRIM-27, PRIM-29, PRIM-31 rows.
- Primitive-count audit: 31/31 primitives still map to implementation files in `internal/agents/` (35 files = 31 primitive impl + 1 state (ADR-C-014) + 2 utility (parser, canonical_crates) + 3 test files + 0 orphan). 8-agent graph still wired; Charm TUI idioms still intact; abcoder-mcp default still in `configs/agents.yml`. Wave-14 trigger criterion met (NeurIPS 2025 + ICML 2025 mechanism papers); fire recorded.
- All gates green: `go build ./...`, `go vet ./...`, `go test ./...`.

## Sprint 2026-09-26 Vault Sync Audit
- Wave-15 status: deferred (3 candidates).
- ADR-C-014 status: locality split applied by Subagent D.
- ADR-C-005 status: magic-string sweep applied by Subagent E.
- ADR-C-011 status: Charm stack audit applied by Subagent G.
- ADR-V-001 status: scripts/lint_vault.sh shipped by this subagent.
- P-106 status: P-106 BFCL retire task deferred to Sprint 2026-09-27 (Subagent F cancelled).
- Parity check: 31/31/31 unchanged.

## Sprint 2026-09-27 Vault Sync Audit (Wave-16)

- **Wave-16 FIRED** 2026-09-27. 1 ACCEPT ([[1.0.0 P-122]] ReasoningBank ICLR 2026 tentative) + 2 REJECT (SWE-Bench Pro ICML 2026 + CodeClash ICML 2026 — both fail Q2 mechanism-vs-benchmark gate).
- **Cross-links propagated** into this INDEX: P-122 into PRIM-21, PRIM-29, PRIM-31. Sources: `.obsidian/MAgHARCM/research/diary/Wave-16-Candidates.md` + paper note `P-122-ReasoningBank-ICLR-2026.md`.
- **Total paper notes**: 122 (121 prior + P-122).
- **Parity**: 31/31/31 unchanged — 31 primitives listed above, 31 implementation files in `internal/agents/*.go` (35 files = 31 primitive impl + 1 state (ADR-C-014) + 2 utility (parser, canonical_crates) + 3 test files + 0 orphan).

## Sprint 2026-09-28 Vault Sync Audit (Wave-17)

- **Wave-17 FIRED** 2026-09-28. 2 ACCEPT (`[[1.0.0 P-123]]` CodeChemist ICML 2026 + `[[1.0.0 P-124]]` Syzygy ICLR 2025 VerifAI Workshop) + 3 REJECT (MemSearcher ACL 2026 Findings — off-list venue; Verified Tool Calls arXiv-only — no venue confirmation; LLM-IR / program-comprehension — benchmark, no mechanism).
- **Cross-links propagated** into this INDEX: P-123 into PRIM-21, PRIM-23, PRIM-27. P-124 into PRIM-9, PRIM-22, PRIM-30. Sources: `.obsidian/MAgHARCM/research/diary/Wave-17-Candidates.md` + paper notes `P-123-CodeChemist-ICML-2026.md` and `P-124-Syzygy-ICLR2025-Workshop.md`.
- **Total paper notes**: 124 (122 prior + P-123 + P-124).
- **Parity**: 31/31/31 unchanged — 31 primitives listed above, 31 implementation files in `internal/agents/*.go` (35 files = 31 primitive impl + 1 state (ADR-C-014) + 2 utility (parser, canonical_crates) + 3 test files + 0 orphan).

## Sprint 2026-09-07 Vault Sync Audit (Wave-19)

- **Wave-19 FIRED** 2026-09-07 (second iteration of 2026-09-07, follows Sprint 2026-09-07-Handoff-1.md). 8 ACCEPT (P-134..P-141) + 3 REJECT (R1 TTA*, R2 HELIOS, R3 LongSpec — all fail Q1: workshop redundancy, off-list venue, off-list venue respectively) + 2 UNVERIFIED (W1 SliceMate ISSTA 2026 claim, W2 SWE-TRACE arXiv-only).
- **Cross-links propagated** into this INDEX: P-134 → PRIM-9 + PRIM-31. P-135 → PRIM-7 + PRIM-21. P-136 → PRIM-7. P-137 → PRIM-7 + PRIM-21. P-138 → PRIM-22 + PRIM-31 (borderline workshop-track; **deliberately excluded from PRIM-7** per §7 method-level threshold). P-139 → PRIM-9 + PRIM-22. P-140 → PRIM-21 + PRIM-22. P-141 → PRIM-21 + PRIM-31.
- **Total paper notes**: 141 (133 prior + P-134..P-141). Wave-19 ACCEPT count: 8.
- **REJECT registry** added to `.obsidian/MAgHARCM/Research-Database.json` under `reject_registry.wave-19` (BLK-08 resolved).
- **Watchlist** added under `watchlist.wave-19` for the 2 UNVERIFIED items.
- **Partition contract**: ACCEPT = 8 / REJECT = 3 / UNVERIFIED = 2 = 13 triaged (matches Wave-19-Candidates.md headline).
- **Sources**: `.obsidian/MAgHARCM/research/diary/Wave-19-Candidates.md` + paper notes P-134..P-141.
- **Parity**: 31/31/31 unchanged — 31 primitives listed above, 31 implementation files in `internal/agents/*.go`. No code changes this sprint (Wave-19 is research-only; integration deferred to future sprints).
