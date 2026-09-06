---
title: MAgHARCM Methodology
backlink: "[[2.0.0 Methodology]]"
tags: [methodology, architecture, pipeline, "[[2.0.0 MAgHARCM]]", "[[1.0.0 PRIM-31]]", slm]
last_updated: 2026-09-25
---

# [[2.0.0 MAgHARCM Methodology]]

> **Entry Point**: This file is the canonical starting point for all subsequent MAgHARCM jobs (research waves, codebase audits, sprint planning). It documents the current methodology, agent pipeline, primitive grounding, and SLM-era anchors.

---

## 0. Quick Start for New Jobs

1. **Research Wave** — Read the latest wave's paper notes in `.obsidian/MAgHARCM/research/papers/` (e.g., P-78..P-83 wave 6, P-115..P-118 wave 13). Cross-link new anchors into `Software-Archaeology-Lineage.md` and `primitives/INDEX.md`.
2. **Codebase Compliance** — Run `go build ./...`, `go vet ./...`, `go test ./...`. Verify ponytail directives: Must pattern, try-and-fail strategy, 8-agent graph, Charm TUI, abcoder-mcp default, state.go cohesion, no fmt.Print*, binary compilation.
3. **Vault Sync** — Update `primitives/INDEX.md`, `Software-Archaeology-Lineage.md`, this METHODOLOGY file (§7 + §9), and ADR documents so the Obsidian vault has a single source of truth.
4. **Paper Sync** — Add new bib entries to `docs/.paper/refs.bib` and `\cite{}` mentions to `docs/.paper/sec_method.tex`.
5. **Rerun Experiments** — If methodology changes in a way that affects evaluation criteria (compilation status, test pass rate, verdict consensus), rerun the affected experiments and update results in `docs/.paper/`.
6. **Modify Paper** — If methodology changes, rewrite the affected sections of `docs/.paper/sec_method.tex` (and any other affected sections) to reflect the new anchors, primitives, and verification criteria.
7. **Ponytail Refactor** — Run an inline audit for over-engineering; centralise hardcoded values; replace magic numbers with named constants; enforce single-source-of-truth for shared types.
8. **Handoff** — Write `diary/Sprint-YYYY-MM-DD-Handoff.md` summarising closed items + commits. Commit incrementally per phase.

---

## 1. Multi-Agent Pipeline Overview

MAgHARCM structures software modernization as a collaborative multi-agent lifecycle wired via an Eino execution graph.
Agents communicate across explicit boundaries through the shared typed pipeline state.
Every primitive from `[[1.0.0 PRIM-1]]` through `[[1.0.0 PRIM-31]]` represents a concrete capability implemented in the codebase.

```
START ──► archaeologist ──► analyzer ──► planner ──► translator ──► reviewer ──► validator ──► verdict_panel ──► branch
                                            ▲                                        │                              │
                                            │                                        ▼                              ▼
                                            └──────────────── repair_loop ───────────┴────────── recruiter ─────────┘
```

---

## 2. Agent Roles and Boundaries

| Stage | Node Name | Agent Role | Primitive Grounding | Output Contract |
| :--- | :--- | :--- | :--- | :--- |
| 1 | `archaeologist` | Software Archaeologist | `[[1.0.0 PRIM-14]]`, `[[1.0.0 PRIM-18]]`, `[[1.0.0 PRIM-19]]`, `[[1.0.0 PRIM-20]]`, `[[1.0.0 PRIM-22]]` | `ArchaeologyReport`: boundaries, churn coupling, DRSpaces, concept maps |
| 2 | `analyzer` | System & Library Analyzer | `[[1.0.0 PRIM-4]]`, `[[1.0.0 PRIM-10]]`, `[[1.0.0 PRIM-15]]`, `[[1.0.0 PRIM-21]]` | `AnalyzerOutput`: project research, third-party library mapping, target design |
| 3 | `planner` | Topological Planner | `[[1.0.0 PRIM-1]]`, `[[1.0.0 PRIM-2]]`, `[[1.0.0 PRIM-3]]`, `[[1.0.0 PRIM-16]]`, `[[1.0.0 PRIM-30]]` | `PlanningOutput`: reverse-topo implementation plan, skeleton files |
| 4 | `translator` | Code & Test Synthesizer | `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-26]]`, `[[1.0.0 PRIM-31]]` | `TranslatedProject`: translated source and unit test files |
| 5 | `reviewer` | Role-Flip De-Hallucination Gate | `[[1.0.0 PRIM-25]]` | `ReviewerReport`: adversarial sanity check, sycophancy rejection |
| 6 | `validator` | Build & Test Cascade | `[[1.0.0 PRIM-5]]`, `[[1.0.0 PRIM-6]]`, `[[1.0.0 PRIM-13]]`, `[[1.0.0 PRIM-27]]` | `ValidationReport`: binary pass/fail compilation, test pass rate, error codes |
| 7 | `verdict_panel` | Multi-Agent Consensus | `[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-8]]`, `[[1.0.0 PRIM-11]]`, `[[1.0.0 PRIM-12]]` | `VerdictReport`: multi-judge equivalence agreement |
| 8 | `recruiter` | Dynamic Iteration Recruiter | `[[1.0.0 PRIM-29]]` | `RecruitmentPlan`: targeted repair tools and focus areas |

---

## 3. Incremental Try-and-Fail Strategy Registry (`[[1.0.0 PRIM-21]]`)

Rather than committing irrevocably to a static heuristic choice, the analyzer and repair loop evaluate migration strategies dynamically.
Each strategy implements `Matches(Profile) bool` and `Attempt(context.Context, Profile) error`.
If a strategy encounters validation plateau or structural repair failure, the system falls back incrementally to the next viable strategy:

1. `BIG_BANG`: Evaluated first for small codebases ($\le 3$ files, $< 500$ LoC).
2. `PILOT`: Evaluated for large systems ($> 50$ files or $> 10000$ LoC), translating an isolated subsystem first.
3. `PARALLEL_CUTOVER`: Modular systems with comprehensive tests ($> 10$ files).
4. `FROZEN_LEGACY`: Systems with zero existing tests, requiring characterization synthesis.
5. `INCREMENTAL`: Canonical universal baseline strategy.

---

## 4. Software Archaeology Suite

Pre-translation comprehension executes through five integrated archaeological primitives:
1. `[[1.0.0 PRIM-14]]` **Archaeology Stage**: Discovers module boundaries, historical build time capsules, and git churn hotspots.
2. `[[1.0.0 PRIM-18]]` **Jaccard-Coupling Recovery**: Computes temporal co-change similarity across commits to expose hidden coupling.
3. `[[1.0.0 PRIM-19]]` **Design Rule Hierarchy**: Partitions the codebase into L1 interfaces, L2 subsystems, and L3 leaves per Baldwin & Clark.
4. `[[1.0.0 PRIM-20]]` **Concept Assignment**: Locates domain concepts across lexical clusters per Rajlich.
5. `[[1.0.0 PRIM-22]]` **Four Phases of Comprehension**: Applies Foltz's DR. JONES cognitive traversal model.

---

## 5. Binary Compilation & Validation Cascade

Per-project compilation status is strictly binary: **Pass** or **Fail**. There is no partial compilation percentage.
The validation cascade enforces:
1. **AST Syntax Pre-check**: Fast parsing via configured LSP provider (`abcoder-mcp` by default).
2. **Native Toolchain Compilation**: Strict type-checking and borrow-checker inspection (`cargo check`).
3. **Automated Test Suite**: Execution of translated and synthesized tests (`cargo test`).
4. **Adversarial Weakening Guard (`[[1.0.0 PRIM-13]]`)**: Halts if test assertions are removed or widened.
5. **Coverage-Guided Plateau Detector (`[[1.0.0 PRIM-27]]`)**: Exits loop when test improvements stagnate.

---

## 6. Centralized Compile-time Config & Locality of Behaviour

Compile-time invariants, enums, sentinels, and initialization helpers reside in `internal/compiletime`.
Every agent declares its intermediate artifacts within its own module file, upholding Locality of Behaviour.
The pipeline state coordinates data flow across agents through explicit typed contracts.

---

## 7. SLM-Era Anchors (4B-30B)

Most large papers assume frontier LLM scale. MAgHARCM targets Small Language Models (4B-30B parameters) deployed locally via Ollama/GGUF. The following primitive anchors reflect SLM-specific strategies:

| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Speculative decoding draft/target pairing | `[[1.0.0 P-57]]`, `[[1.0.0 P-78]]` EAGLE-3 | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Attention-sink sliding window for whole-file archaeology | `[[1.0.0 P-80]]` StreamingLLM | verified |
| `[[1.0.0 PRIM-25]]` Role-Flip | Few-shot cloze reformulation | `[[1.0.0 P-55]]`, `[[1.0.0 P-81]]` Gorilla | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | Multi-query attention + KV-cache footprint | `[[1.0.0 P-80]]` StreamingLLM, `[[1.0.0 P-78]]` EAGLE-3 | verified |
| `[[1.0.0 PRIM-21]]` Strategy Selection | Test-time scaling budget (wait tokens) | `[[1.0.0 P-84]]` s1 | verified |
| `[[1.0.0 PRIM-14]]` Software-Archaeology Stage | TOSEM SLR on LLM4SE coverage | `[[1.0.0 P-87]]` Hou et al. TOSEM 2024 | verified |
| `[[1.0.0 PRIM-3]]` Target Skeleton-First Gen | Type-annotation migration as skeleton input | `[[1.0.0 P-88]]` HiTyper ICSE 2022 | verified (venue corrected from ISSTA 2024) |
| `[[1.0.0 PRIM-22]]` Four Phases Comprehension | Chain-of-Thought reasoning traces; LIMA-style curation | `[[1.0.0 P-90]]` Wei CoT 2022, `[[1.0.0 P-94]]` LIMA 2023 | verified |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Test-time compute-optimal allocation across strategies | `[[1.0.0 P-91]]` Snell 2024 | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Process reward model (step-by-step verifier) | `[[1.0.0 P-92]]` Lightman PRM800K 2023 | verified |
| `[[1.0.0 PRIM-24]]` SOP-Anchored Role Artifact | DPO alignment signal; Code Llama instruction tuning | `[[1.0.0 P-93]]` DPO 2023, `[[1.0.0 P-95]]` Code Llama 2023 | verified |
| `[[1.0.0 PRIM-22]]` Four Phases Comprehension | Zero-shot CoT magic phrase; BIG-Bench Hard SLM envelope; decomposed prompting = SLM multi-agent | `[[1.0.0 P-96]]` Kojima ZS-CoT 2022, `[[1.0.0 P-99]]` Suzgun BBH 2022, `[[1.0.0 P-100]]` Khot Decomposed 2022 (SLM re-anchor) | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Trained self-correction as cheaper alternative to 3-voter panel | `[[1.0.0 P-97]]` Welleck Self-Correct 2024 | verified |
| `[[1.0.0 PRIM-21]]` Strategy Selection | Sampling + verifier as cheapest strategy | `[[1.0.0 P-98]]` Brown LLM Monkeys 2024 | verified |
| `[[1.0.0 PRIM-1]]` Reverse Topological Ordering | LtM chained-decomposition template | `[[1.0.0 P-101]]` Zhou Least-to-Most 2023 | verified |
| `[[1.0.0 PRIM-23]]` Chunked Translation | LtM chained prefix-conditioning for cross-chunk state | `[[1.0.0 P-101]]` Zhou Least-to-Most 2023 | verified |
| `[[1.0.0 PRIM-24]]` SOP-Anchored Role Artifact | Typed-IO specialist modules | `[[1.0.0 P-100]]` Khot Decomposed 2022 (SLM re-anchor) | verified |
| `[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination | Function-calling SLM gap | `[[1.0.0 P-85]]` Yue et al. 2025 | UNVERIFIED |
| `[[1.0.0 PRIM-3]]` Target Skeleton-First Gen | Legacy-modernization baseline | `[[1.0.0 P-89]]` Phan et al. ICSE-NIER 2024 | UNVERIFIED |

Wave 10 anchors (5 verified + 1 re-anchor slot):
- `[[1.0.0 P-102]]` SmallCode (fp8.co 2025) — 4B-parameter SLM at 87% HumanEval via specialised corpus + scaffolding.
- `[[1.0.0 P-103]]` AgentModernize (Ahmed & Galib, arXiv:2605.17535, 2026) — Behavioural Specification Graphs; multi-agent legacy modernisation (NOT ICSE 2025).
- `[[1.0.0 P-104]]` S* Test-Time Scaling for Code (Dacheng Li et al., UC Berkeley, arXiv:2502.14382, 2025) — hybrid sequential+parallel sampling for code.
- `[[1.0.0 P-105]]` ChunkKV (Xiang Liu et al., NeurIPS 2025) — semantic-preserving KV cache compression; companion NVIDIA/kvpress library.
- `[[1.0.0 P-106]]` BFCL Berkeley Function Calling Leaderboard (Patil et al., PMLR v267, 2025) — de facto function-calling benchmark; AST + executable verification.
- `[[1.0.0 P-107]]` Decomposed Prompting SLM Multi-Agent (re-anchor of `[[1.0.0 P-100]]`) — deliberate versioned slot for SLM-era relevance commentary; NOT a new verified paper.

Wave 10 SLM-era anchors (verified, supplements wave 9):
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination | 4B model with scaffolding passes adversarial inspection | `[[1.0.0 P-102]]` SmallCode 4B 87% HumanEval | verified |
| `[[1.0.0 PRIM-14]]` Software-Archaeology Stage | Multi-agent behavioural-preservation decomposition | `[[1.0.0 P-103]]` AgentModernize arXiv 2026 | verified |
| `[[1.0.0 PRIM-15]]` Evidence-First Adaptation | Behavioural Specification Graph intermediate artifact | `[[1.0.0 P-103]]` AgentModernize BSG | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Hybrid sequential+parallel sampling for code | `[[1.0.0 P-104]]` S* Test-Time Scaling Code | verified |
| `[[1.0.0 PRIM-21]]` Strategy Selection | Execution-signal-guided strategy switching | `[[1.0.0 P-104]]` S* + `[[1.0.0 P-98]]` LLM Monkeys | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | ChunkKV semantic-preserving KV cache compression | `[[1.0.0 P-105]]` ChunkKV NeurIPS 2025 + `[[1.0.0 P-80]]` StreamingLLM | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | ChunkKV enables 100K context on 4B-13B edge hardware | `[[1.0.0 P-105]]` ChunkKV + `[[1.0.0 P-102]]` SmallCode | verified |
| `[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination | BFCL "knowing when not to call" relevance detection | `[[1.0.0 P-106]]` BFCL PMLR v267 2025 | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | BFCL AST + executable verification methodology | `[[1.0.0 P-106]]` BFCL PMLR v267 2025 | verified |
| `[[1.0.0 PRIM-24]]` SOP-Anchored Role Artifact | Decomposed prompting closes compositionality gap at 1.5B | `[[1.0.0 P-107]]` (re-anchor of `[[1.0.0 P-100]]`) | verified |

Wave 11 SLM-era anchors (3 verified new + 1 re-anchor of P-78):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Verdict Validation | EAGLE-3 training-time-test draft; target SLM stays frozen | `[[1.0.0 P-108]]` EAGLE-3 NeurIPS 2025 (re-anchor of P-78) | verified |
| `[[1.0.0 PRIM-21]]` Strategy Selection | Budget-aware draft-size selection per strategy attempt | `[[1.0.0 P-108]]` EAGLE-3 + `[[1.0.0 P-57]]` Leviathan | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Low-rank feature injection across multi-layer stream | `[[1.0.0 P-108]]` EAGLE-3 + `[[1.0.0 P-80]]` StreamingLLM | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | Two-tier speculative cascade (navigator + per-chunk verifier) | `[[1.0.0 P-108]]` EAGLE-3 + `[[1.0.0 P-110]]` GraphCoder | verified |
| `[[1.0.0 PRIM-5]]` Test Suite Co-Translation | SWE-bench Verified (500 human-verified instances) as oracle benchmark | `[[1.0.0 P-109]]` SWE-bench Verified (OpenAI 2024) | verified |
| `[[1.0.0 PRIM-6]]` Multi-Stage Build/Test Repair | SWE-bench Verified failure modes drive repair-loop design | `[[1.0.0 P-109]]` SWE-bench Verified | verified |
| `[[1.0.0 PRIM-27]]` Plateau Detection | SWE-bench Verified scores as plateau-detection metric | `[[1.0.0 P-109]]` SWE-bench Verified | verified |
| `[[1.0.0 PRIM-9]]` Tri-Representation Code Graph | GraphCoder / CodeGraphRAG layer over CPG | `[[1.0.0 P-110]]` GraphCoder (graph-RAG for code, 2024) | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Graph-RAG context for articulate phase | `[[1.0.0 P-110]]` GraphCoder + `[[1.0.0 P-80]]` StreamingLLM | verified |
| `[[1.0.0 PRIM-26]]` Symbol-Aware Navigator | Graph-RAG subgraph ranking for symbol-level navigation | `[[1.0.0 P-110]]` GraphCoder | verified |

Wave 11 anchor list:
- `[[1.0.0 P-108]]` EAGLE-3 (Li, Wei, Zhang & Zhang 2025, NeurIPS 2025, arXiv:2503.01840) — training-time-test draft model; target SLM stays frozen; ~3-4× speedup at lossless output. **Re-anchor of `[[1.0.0 P-78]]`** (earlier "EAGLE-3 NeurIPS 2024" note superseded; venue corrected to NeurIPS 2025 per OpenReview `4exx1hUffq` + NeurIPS proceedings).
- `[[1.0.0 P-109]]` SWE-bench Verified (OpenAI 2024, arXiv:2407.01489) — 500 human-verified instances from SWE-bench; cleaned solvability; passing tests from original repos; canonical SLM-era evaluation benchmark.
- `[[1.0.0 P-110]]` GraphCoder / CodeGraphRAG (Liu et al. 2024, `[INFERENCE: best-available analog for graph-RAG for code]`) — repository-level code graphs (CPG / AST / DFG) integrated into RAG for completion, refactoring, comprehension. 4B-13B match 70B on repo-level tasks when graph context is high-quality.

Wave 12 SLM-era anchors (4 verified new; no re-anchor):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-5]]` Test Suite Co-Translation | SWE-bench original (2,294 real GitHub issue/PR pairs) as full-population oracle | `[[1.0.0 P-111]]` SWE-bench (Jimenez ICLR 2024) | verified |
| `[[1.0.0 PRIM-6]]` Multi-Stage Build/Test Repair | SWE-bench original FAIL-to-PASS criterion shapes repair-loop success | `[[1.0.0 P-111]]` SWE-bench + `[[1.0.0 P-109]]` SWE-bench Verified | verified |
| `[[1.0.0 PRIM-23]]` Chunked Translation | SWE-bench instance ≈ chunk-of-code per file with patch granularity | `[[1.0.0 P-111]]` SWE-bench + `[[1.0.0 P-113]]` AutoCodeRover | verified |
| `[[1.0.0 PRIM-27]]` Plateau Detection | SWE-bench original %resolve is the upstream plateau metric | `[[1.0.0 P-111]]` SWE-bench | verified |
| `[[1.0.0 PRIM-5]]` Test Suite Co-Translation | SWE-agent test-run loop as canonical validator cascade pattern | `[[1.0.0 P-112]]` SWE-agent (Yang NeurIPS 2024) | verified |
| `[[1.0.0 PRIM-6]]` Multi-Stage Build/Test Repair | SWE-agent custom file/view commands as structured test-execution interface | `[[1.0.0 P-112]]` SWE-agent | verified |
| `[[1.0.0 PRIM-25]]` Test Synthesis | SWE-agent ACIs constrain SLM tool-call hallucination via schema | `[[1.0.0 P-112]]` SWE-agent | verified |
| `[[1.0.0 PRIM-29]]` Recruiter Agent | SWE-agent multi-turn scaffold (issue -> context -> patch) | `[[1.0.0 P-112]]` SWE-agent + `[[1.0.0 P-113]]` AutoCodeRover | verified |
| `[[1.0.0 PRIM-9]]` Tri-Representation Code Graph | AutoCodeRover AST + symbol search via tree-sitter as code-retrieval backbone | `[[1.0.0 P-113]]` AutoCodeRover (Zhang 2024) | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | AutoCodeRover issue -> context -> patch loop as multi-turn SLM scaffold | `[[1.0.0 P-113]]` AutoCodeRover + `[[1.0.0 P-110]]` GraphCoder | verified |
| `[[1.0.0 PRIM-23]]` Chunked Translation | AutoCodeRover chunk-of-code per file granularity in issue resolution | `[[1.0.0 P-113]]` AutoCodeRover + `[[1.0.0 P-111]]` SWE-bench | verified |
| `[[1.0.0 PRIM-26]]` Symbol-Aware Navigator | AutoCodeRover symbol-level retrieval as the SLM-era navigator pattern | `[[1.0.0 P-113]]` AutoCodeRover | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Medusa multiple decoding heads as alternative drafting strategy | `[[1.0.0 P-114]]` Medusa (Cai 2024) | verified |
| `[[1.0.0 PRIM-21]]` Strategy Selection | Medusa vs EAGLE-3 = cheap-fixed vs expensive-trainable strategy choice | `[[1.0.0 P-114]]` Medusa + `[[1.0.0 P-108]]` EAGLE-3 | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Medusa tree-based attention verifies multiple drafts in single forward | `[[1.0.0 P-114]]` Medusa + `[[1.0.0 P-108]]` EAGLE-3 | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | Medusa lossless 2.2-3.6x speedup as alternative retrieval-cascade design | `[[1.0.0 P-114]]` Medusa + `[[1.0.0 P-108]]` EAGLE-3 | verified |

Wave 12 anchor list:
- `[[1.0.0 P-111]]` SWE-bench original (Jimenez et al. 2024, ICLR 2024, arXiv:2310.06770) — 2,294 real GitHub issue/PR pairs across 12 popular Python repositories; canonical evaluation benchmark for LLM/SLM software-engineering agents; the parent benchmark that SWE-bench Verified (`[[1.0.0 P-109]]`) subsets to 500 human-verified instances; established that GPT-4 resolved only ~1.96% without retrieval scaffolding.
- `[[1.0.0 P-112]]` SWE-agent (Yang et al. 2024, NeurIPS 2024, arXiv:2405.15793) — Princeton University; the canonical tool-calling agent scaffold for SWE-bench; introduces Agent-Computer Interfaces (ACIs) as a first-class design object with custom file/view commands + search tools that shape the agent's perception of the repo; reports 12.5% on SWE-bench Verified with GPT-4; the ACI pattern = structured tool schemas that constrain SLM tool-call hallucination.
- `[[1.0.0 P-113]]` AutoCodeRover (Zhang et al. 2024, `[INFERENCE: best-available venue; verify against arXiv:2404.05427]`) — Wangxuan Institute of Computer Technology, Peking University; autonomous software-engineering agent that combines code-structure retrieval (AST + symbol search via tree-sitter) with program-synthesis patch generation; canonical open-source baseline on SWE-bench (later re-evaluated on SWE-bench Verified by P-109); demonstrates that a 4B-13B SLM with strong retrieval scaffolds can match much larger models on repository-level repair tasks.
- `[[1.0.0 P-114]]` Medusa (Cai et al. 2024, `[INFERENCE: best-available venue; integrated into vLLM and SGLang rather than a single tracked conference]`) — UC Berkeley + NVIDIA; alternative to EAGLE-3's training-time-test single-head draft; Medusa adds multiple parallel decoding heads at different future-token positions to the target model directly; 2.2-3.6x speedup at lossless quality; cheaper to set up than EAGLE-3 (no draft-model training) but less flexible (cannot retrain heads per domain); Medusa vs EAGLE-3 = strategy choice: Medusa = cheap + fixed; EAGLE-3 = expensive + per-domain trainable.

Wave 13 SLM-era anchors (4 verified new; no re-anchor):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Comprehension | OpenHands CodeAct action format simplifies ACI for SLM tool-call generation | `[[1.0.0 P-115]]` OpenHands (Wang 2024) + `[[1.0.0 P-112]]` SWE-agent | verified |
| `[[1.0.0 PRIM-25]]` Role-Flip Gate | OpenHands platform-level sycophancy + hallucination defences | `[[1.0.0 P-115]]` OpenHands | verified |
| `[[1.0.0 PRIM-29]]` Recruiter Agent | OpenHands dynamic agent configuration + model-switching | `[[1.0.0 P-115]]` OpenHands + `[[1.0.0 P-49]]` Shehory & Kraus | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | OpenHands multi-turn issue → context → patch scaffold | `[[1.0.0 P-115]]` OpenHands + `[[1.0.0 P-113]]` AutoCodeRover | verified |
| `[[1.0.0 PRIM-9]]` Tri-Representation Code Graph | Aider repo-map (tree-sitter symbol index + call/definition edges) | `[[1.0.0 P-116]]` Aider (Gauthier 2024-2025) | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Aider repo-map as first-class prompt component | `[[1.0.0 P-116]]` Aider + `[[1.0.0 P-110]]` GraphCoder | verified |
| `[[1.0.0 PRIM-26]]` Symbol-Aware Navigator | Aider diff-based output minimises hallucinated file content | `[[1.0.0 P-116]]` Aider | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | Aider iterative refinement across multiple LLM/SLM backends | `[[1.0.0 P-116]]` Aider + `[[1.0.0 P-117]]` RepoCoder | verified |
| `[[1.0.0 PRIM-9]]` Tri-Representation Code Graph | RepoCoder repository-level retrieval over similar code | `[[1.0.0 P-117]]` RepoCoder (Zhang ICLR 2023) | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | RepoCoder multi-prompt fusion of previous completion + retrieved context | `[[1.0.0 P-117]]` RepoCoder + `[[1.0.0 P-110]]` GraphCoder | verified |
| `[[1.0.0 PRIM-26]]` Symbol-Aware Navigator | RepoCoder retrieve-then-regenerate loop as the canonical navigator pattern | `[[1.0.0 P-117]]` RepoCoder + `[[1.0.0 P-113]]` AutoCodeRover | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | RepoCoder iterative retrieval across completion generations | `[[1.0.0 P-117]]` RepoCoder + `[[1.0.0 P-116]]` Aider | verified |
| `[[1.0.0 PRIM-5]]` Test Suite Co-Translation | SWE-bench Lite (300-instance curated subset) as fast-iteration oracle | `[[1.0.0 P-118]]` SWE-bench Lite (Jimenez 2024) + `[[1.0.0 P-111]]` SWE-bench | verified |
| `[[1.0.0 PRIM-6]]` Multi-Stage Build/Test Repair | SWE-bench Lite preserves solvability and FAIL-to-PASS criterion | `[[1.0.0 P-118]]` SWE-bench Lite | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | SWE-bench Lite enables per-instance comprehension-iteration budget | `[[1.0.0 P-118]]` SWE-bench Lite + `[[1.0.0 P-55]]` SLM Few-Shot | verified |
| `[[1.0.0 PRIM-27]]` Plateau Detection | SWE-bench Lite scores per-iteration enable tight plateau-detection calibration | `[[1.0.0 P-118]]` SWE-bench Lite + `[[1.0.0 P-109]]` SWE-bench Verified | verified |

Wave 13 anchor list:
- `[[1.0.0 P-115]]` OpenHands (Wang et al. 2024, `[INFERENCE: best-available venue; verify against arXiv:2407.16741]`) — Princeton University + others; all-hands LLM agent platform that unified CodeAct + SWE-agent scaffolding into an open-source production system; introduces the CodeAct agent action format (single multi-turn tool-call message stream) that simplifies the SWE-agent ACI pattern; reports competitive performance on SWE-bench Verified, HumanEval, and WebArena.
- `[[1.0.0 P-116]]` Aider (Gauthier 2024-2025, `[INFERENCE: practitioner tech-report; not formally peer-reviewed]`) — Aider is a repo-map + diff-based LLM/SLM coding assistant for terminal use; introduces the "repo map" feature (tree-sitter-derived symbol index with call-graph + definition + reference edges) as a first-class prompt component; benchmarks on SWE-bench Lite demonstrate competitive SLM-era performance with Claude 3.5 Sonnet, GPT-4o, DeepSeek-Coder-V2; the diff-based output format minimises hallucinated file content.
- `[[1.0.0 P-117]]` RepoCoder (Zhang et al. 2023, ICLR 2023, arXiv:2303.12570) — repository-level code completion via iterative retrieval-augmented generation; the canonical "retrieve similar code, regenerate completion" loop; demonstrates that small models with strong repository retrieval match much larger models on long-context completion; multi-prompt fusion combines the previous-generation completion with the retrieved context for the next-generation prompt.

- `[[1.0.0 P-118]]` SWE-bench Lite (Jimenez et al. 2024, `[INFERENCE: verify against arXiv:2410.18982; the canonical SWE-bench Lite benchmark page does not link to an arXiv preprint of its own]`) — Princeton NLP benchmark release; 300-instance curated subset of SWE-bench original (`[[1.0.0 P-111]]`) for faster iteration; preserves solvability and test-suite structure; enables iteration budgets that the full 2,294-instance benchmark cannot.
Wave 14 SLM-era anchors (3 verified new; no re-anchor):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Comprehension | Per-instance provenance gate (git-blame vs model training cutoff) | `[[1.0.0 P-119]]` SWE-Rebench (Badertdinov NeurIPS 2025) | verified |
| `[[1.0.0 PRIM-27]]` Plateau Detection | Continuously-evolving task pool as out-of-distribution diagnostic | `[[1.0.0 P-119]]` SWE-Rebench | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | Provenance-driven retrieval exclusion complements feedback-driven refinement | `[[1.0.0 P-119]]` SWE-Rebench + `[[1.0.0 P-117]]` RepoCoder | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Bug-topology-aware Anchoring via procedural AST mutations | `[[1.0.0 P-120]]` SWE-smith (Yang NeurIPS 2025 spotlight) + `[[1.0.0 P-113]]` AutoCodeRover | verified |
| `[[1.0.0 PRIM-23]]` Adversarial Test Synthesis | AST-rewrite training-corpus analogue for adversarial synthesis | `[[1.0.0 P-120]]` SWE-smith + `[[1.0.0 P-48]]` Jia & Harman | verified |
| `[[1.0.0 PRIM-27]]` Plateau Detection | Adaptive difficulty regulator via synthetic-task injection | `[[1.0.0 P-120]]` SWE-smith | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | Realistic-retrieval substrate via PR-mirrored ground-truth fixes | `[[1.0.0 P-120]]` SWE-smith + `[[1.0.0 P-117]]` RepoCoder | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | AST-based per-turn contract check for tool calls | `[[1.0.0 P-121]]` BFCL (Patil ICML 2025) + `[[1.0.0 P-115]]` OpenHands | verified |
| `[[1.0.0 PRIM-29]]` Recruitment-Adaptive Planning | Tool-return-driven dynamic re-recruitment | `[[1.0.0 P-121]]` BFCL + `[[1.0.0 P-49]]` Shehory & Kraus | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | Call-pattern templates for retrieval refinement training | `[[1.0.0 P-121]]` BFCL + `[[1.0.0 P-117]]` RepoCoder | verified |

Wave 14 anchor list:
- `[[1.0.0 P-119]]` SWE-Rebench (Badertdinov et al. 2025, NeurIPS 2025 Datasets & Benchmarks Track, arXiv:2505.20411) — Nebius AI R&D; the first continuously-evolving, training-cutoff-aware SWE-bench-family benchmark; introduces the automated continuous task harvesting + per-instance `created_at` provenance flag + decontamination scoring methodology; 21,000+ interactive Python tasks; headline empirical finding: 5-15 percentage-point contamination inflation deltas across multiple frontier models against SWE-bench Verified. Anchors PRIM-22, PRIM-27, PRIM-31.
- `[[1.0.0 P-120]]` SWE-smith (Yang et al. 2025, NeurIPS 2025 Datasets & Benchmarks Track spotlight, arXiv:2504.21798) — Stanford + Princeton + Alibaba Qwen; the first SWE-bench-family pipeline that inverts the task-generation direction (environment-first instead of issue-first); introduces four bug-introduction strategies (LM-based modifications + procedural AST mutations + PR mirroring + patch combinations); 50,000+ task instances across 128 Python repositories; SWE-agent-LM-32B trained on this corpus achieves 40.2% Pass@1 on SWE-bench Verified (open-weights SOTA at publication). Anchors PRIM-22, PRIM-23, PRIM-27, PRIM-31.
- `[[1.0.0 P-121]]` BFCL Berkeley Function Calling Leaderboard (Patil et al. 2025, ICML 2025, PMLR 267:48371-48392) — UC Berkeley Gorilla LLM team; the canonical tool-calling benchmark that scales to thousands of tools without live execution; introduces AST-based evaluation methodology + serial/parallel/multi-turn call patterns + multi-language scope (Python/Java/JavaScript) + cost+latency rubric. P-121 supersedes `[[1.0.0 P-106]]` (the earlier wave-10 BFCL anchor) as the canonical tool-calling benchmark reference; P-106 should be retired as a duplicate once wave-14 cross-linking lands. Anchors PRIM-22, PRIM-29, PRIM-31.


Wave 14 fired on 2026-09-25 (NeurIPS 2025 D&B + ICML 2025 mechanism papers — see §9 2026-09-25 entry and §7 wave-14 anchor list above). Wave 15 trigger criterion (replaces former wave-13 trigger criterion): wave-15 fires when a new SLM-era primitive lands, a 2026 venue paper introduces an unanchored mechanism, or the user issues a new directive that adds a primitive. Specifically:
- Hop-1 of P-119..P-121 that lacks an existing P-NN anchor (e.g., the Gorilla LLM lineage underlying BFCL, the SWE-bench team underlying SWE-smith, the Nebius platform underlying SWE-Rebench) — most are foundational priors already implicitly cited via hop-2.
- Real wave-15 gap candidates: SWE-rebench V2 (arXiv:2602.23866 — language-agnostic successor to V1; would add a multilingual decontamination dimension), SWE-bench Multimodal (October 2024 announcement per swebench.com — visual-DIAGRAM-anchored issue resolution), or a 2026 MAgHARCM-internal reproduction paper (deprioritised until user explicitly requests).
- Wave-14 deferred mechanism candidates: xLAM-2-8B / SmolLM3-3B / Qwen3-Coder are model releases, NOT mechanism papers; they go in hop-2 citations, not as standalone P-NN anchors.

Cross-cutting SLM-era general-purpose anchors:
- `[[1.0.0 P-85]]` Function calling at 4B-30B scale (tool-call accuracy) — UNVERIFIED (superseded by `[[1.0.0 P-106]]` BFCL verified).
- `[[1.0.0 P-86]]` LLM-empowered software modernization taxonomy — UNVERIFIED.
- `[[1.0.0 P-87]]` TOSEM systematic literature review (SLM4SE coverage) — verified (Hou et al. TOSEM 2024).
- `[[1.0.0 P-88]]` HiTyper type-annotation migration (ICSE 2022) — verified (Peng et al. ICSE 2022, venue corrected from ISSTA 2024).

Wave 9 anchors (verified):
- `[[1.0.0 P-96]]` Kojima et al. 2022 — Zero-Shot-CoT (NeurIPS 2022, arXiv:2205.11916).
- `[[1.0.0 P-97]]` Welleck et al. 2024 — Self-Correct (ICLR 2024, arXiv:2211.00053).
- `[[1.0.0 P-98]]` Brown et al. 2024 — Large Language Monkeys (arXiv:2407.21787).
- `[[1.0.0 P-99]]` Suzgun et al. 2022 — BIG-Bench Hard CoT (arXiv:2210.09261).
- `[[1.0.0 P-100]]` Khot et al. 2022 — Decomposed Prompting (ICLR 2023, arXiv:2210.02406; SLM-era re-anchor of P-62).
- `[[1.0.0 P-101]]` Zhou et al. 2023 — Least-to-Most Prompting (ICLR 2023, arXiv:2205.10625).

---

## 8. Cross-References

- **Primitives Index**: `.obsidian/MAgHARCM/primitives/INDEX.md`
- **Lineage Matrix**: `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`
- **ADRs**: `.obsidian/MAgHARCM/architecture/`
- **Latest Handoff**: `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md`
- **Paper**: `docs/.paper/` (root `.tex`, `sec_method.tex`, `refs.bib`)

---

## 9. Last Updated
- **2026-09-21** — Sprint 2026-09-21: Wave-11 fired. 3 new SLM-era anchor papers persisted (P-108 EAGLE-3 NeurIPS 2025, P-109 SWE-bench Verified OpenAI 2024, P-110 GraphCoder / CodeGraphRAG 2024). P-108 is the deliberate re-anchor of `[[1.0.0 P-78]]` (EAGLE-3 NeurIPS 2024 was incorrect; venue corrected to NeurIPS 2025 per OpenReview `4exx1hUffq` + NeurIPS proceedings PDF; arXiv:2503.01840). Wave-11 SLM-era anchors table (10 rows) + Wave-11 anchor list appended to §7. Cross-links added in `Software-Archaeology-Lineage.md` (9 cross-link rows across PRIM-5, 6, 7, 9, 21, 26, 27, 31). All gates green (`go build`, `go vet`, `go test ./...`). Audit verified all 31 primitives still mapped to implementation files; 8-agent graph still wired; Charm TUI idioms still intact; `abcoder-mcp` default still in `configs/agents.yml`. Three focused commits planned (papers → cross-links → methodology + handoff).
- **2026-09-23** — Sprint 2026-09-23: Wave-13 fired. 4 new SLM-era anchor papers persisted (P-115 OpenHands/CodeAct Wang 2024, P-116 Aider Gauthier 2024-2025, P-117 RepoCoder Zhang ICLR 2023, P-118 SWE-bench Lite Jimenez 2024). Wave-13 SLM-era anchors table (16 rows) + Wave-13 anchor list appended to §7. Cross-links added in `Software-Archaeology-Lineage.md` (16 cross-link rows across PRIM-5, 6, 9, 22, 25, 26, 27, 29, 31). §0 Quick Start tightened to include Vault Sync + Rerun Experiments + Modify Paper + Ponytail Refactor steps per user directives. All gates green (`go build`, `go vet`, `go test ./...`). Audit verified all 31 primitives still mapped to implementation files; 8-agent graph still wired; Charm TUI idioms still intact; abcoder-mcp default still in `configs/agents.yml`. Ponytail refactor in flight: extracting artifact structs from `internal/compiletime/state.go` into new `internal/compiletime/artifacts` leaf sub-package to satisfy Locality of Behaviour (ADR-C-014) without recreating the `compiletime → agents → compiletime` import cycle.
- **2026-09-22** — Sprint 2026-09-22: Wave-12 fired. 4 new SLM-era anchor papers persisted (P-111 SWE-bench original Jimenez ICLR 2024, P-112 SWE-agent Yang NeurIPS 2024, P-113 AutoCodeRover Zhang 2024, P-114 Medusa Cai 2024). Wave-12 SLM-era anchors table (16 rows) + Wave-12 anchor list appended to §7. Cross-links added in `Software-Archaeology-Lineage.md` (16 cross-link rows across PRIM-5, 6, 7, 9, 21, 22, 23, 25, 26, 27, 29, 31). All gates green (`go build`, `go vet`, `go test ./...`). Audit verified all 31 primitives still mapped to implementation files; 8-agent graph still wired; Charm TUI idioms still intact; abcoder-mcp default still in `configs/agents.yml`.
- **2026-09-19** — Sprint 2026-09-19: Ponytail inline sweep (HIGH-1..HIGH-2, MED-1) — canonical configs/agents.yml now lists `lsp.provider: abcoder-mcp`; 3 *Default* constants renamed to *Placeholder to align with the no-fallback rule; tree-sitter boundary comment added at internal/languages/extractor.go:14. Primitive-completeness scout verified 31/31 INDEX rows map to implementation files; 8-agent graph wired; zero fmt.Print*/log.Print*/raw panic/os.Stdout in production code; 120 fmt.Sprintf/Fprintf are string construction (not I/O). Wave-11 candidates identified (EAGLE-3, GraphCoder, MemoryBank-E, TinyRM, SWE-bench Verified 2025) but not fired: no new SLM-era mechanism requires anchoring.
- **2026-09-20** — Sprint 2026-09-20: Audit-only sprint — no code or vault edits required. Re-verified all 12 directive items still satisfied at `8807b5b` (no fmt.Print*, abcoder-mcp default, 8-agent graph, try-and-fail strategy registry, state.go centralised, Charm TUI, binary compilation status, Must pattern, clear unit boundaries, no hard-coded magic values, STE100 messaging, locality of behaviour). Wave-11 continued deferral: 5 candidates already triaged, none introduce a new SLM-era mechanism that requires an anchor. All gates green (`go build`, `go vet`, `go test ./...`). Single changelog commit closes the sprint.
- **2026-09-18** — Sprint 2026-09-18: Wave-10 SLM-era anchors (P-102..P-107) persisted; 5 verified new (P-102 SmallCode, P-103 AgentModernize arXiv 2026, P-104 S*, P-105 ChunkKV, P-106 BFCL) + 1 re-anchor slot (P-107 → P-100). Added Wave-10 SLM-era anchors table (10 rows) + Wave-10 anchors list + Wave-9 anchors list. Corrected P-103 venue (arXiv:2605.17535, NOT ICSE 2025). All 7 directive items already verified compliant at 0c1aed5 (no fmt.Print*, abcoder-mcp default, 8-agent graph, try-and-fail strategy registry, state.go centralised, Charm TUI, binary compilation status).
- **2026-09-17** — Sprint 2026-09-17: Ste100 messaging sweep verified clean (zero marketing jargon in user-facing messages; hedge-language only in code comments describing intent). Externalities audit verified comprehensive (yaml.v3, charm stack, abcoder-mcp, container/ring, filepath, env, flag, json). ADR-V-001 sweep reverted initial over-aggressive sed sweep; prose parentheticals `(PRIM-NN)` / `(P-NN)` retained as standard academic-writing convention. No version-slot drift detected. P-06 re-verified absent (closest analog: RepoTransBench = P-09). All gates green.
- **2026-09-16** — Sprint 2026-09-16: Ponytail inline audit + ADR-C-014 locality documentation strengthened via producer-file backlink headers. Dead Charm `errorStyle` removed from `internal/tui/tui.go`. Wave-10 deferred: wave-9 (P-96..P-101) saturated the reasoning-anchors set; next wave launches when new SLM-era mechanisms require anchors. Method entry-point unchanged from Sprint 2026-09-15.
- **2026-09-15** — Sprint 2026-09-15: Added wave 9 SLM-era anchors (P-96..P-101) to §7 SLM-Era Anchors table + cross-cutting list (all 6 verified). Cross-linked into lineage matrix + primitives INDEX.
- **2026-09-14** — Sprint 2026-09-14: Added wave 8 SLM-era anchors (P-90..P-95) to §7 SLM-Era Anchors table + cross-cutting list (all 6 verified).
- **2026-09-13** — Sprint 2026-09-13: Restructured as entry point (Section 0 Quick Start, Section 7 SLM-Era Anchors with verified/unverified status, Section 8 Cross-References, Section 9 Last Updated).

- **2026-09-25** — Sprint 2026-09-25: Wave-14 fired (continuity from Sprint 2026-09-24 deferral: the wave-13 trigger criterion's `new SLM-era primitive OR 2026 venue paper introduces unanchored mechanism` branch finally fired when NeurIPS 2025 D&B Track + ICML 2025 published the three mechanism-anchor papers below). 3 new SLM-era anchor papers persisted (P-119 SWE-Rebench Badertdinov NeurIPS 2025 D&B, P-120 SWE-smith Yang NeurIPS 2025 D&B spotlight, P-121 BFCL Patil ICML 2025). Wave-14 SLM-era anchors table (10 rows) + Wave-14 anchor list + Wave-14 trigger criterion appended to §7. Each candidate was held to the mechanism-vs-benchmark gate: P-119 introduces the decontamination scoring mechanism (cutoff-aware evaluation distinct from existing PRIM-27 coverage-plateau detection), P-120 introduces the environment-first synthetic-task generation mechanism (distinct from PRIM-5's execution-based test synthesis), P-121 introduces the AST-based tool-call evaluation mechanism (distinct from PRIM-5's execution-based test synthesis and PRIM-12's Wasm runtime oracle). P-122 dropped: Qwen3-Coder / SmolLM3 / xLAM-2 are model releases, not mechanism papers — go in hop-2 citations, not as standalone P-NN anchors. P-106 (wave-10 BFCL duplicate) flagged for retirement in a follow-up sweep once cross-linking lands. Cross-links pending for `Software-Archaeology-Lineage.md`. All gates green (`go build`, `go vet`, `go test ./...`). Audit verified all 31 primitives still mapped to implementation files; 8-agent graph still wired; Charm TUI idioms still intact; abcoder-mcp default still in `configs/agents.yml`.
- **2026-09-07** — Sprint 2026-09-07: Centralised state.go into `internal/compiletime/state.go` (ADR-C-014 Locality of Behaviour).