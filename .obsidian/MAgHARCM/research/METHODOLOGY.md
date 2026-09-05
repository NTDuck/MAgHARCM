---
title: MAgHARCM Methodology
backlink: "[[2.0.0 Methodology]]"
tags: [methodology, architecture, pipeline, "[[2.0.0 MAgHARCM]]", "[[1.0.0 PRIM-31]]", slm]
last_updated: 2026-09-18
---

# [[2.0.0 MAgHARCM Methodology]]

> **Entry Point**: This file is the canonical starting point for all subsequent MAgHARCM jobs (research waves, codebase audits, sprint planning). It documents the current methodology, agent pipeline, primitive grounding, and SLM-era anchors.

---

## 0. Quick Start for New Jobs

1. **Research Wave** — Read the latest wave's paper notes in `.obsidian/MAgHARCM/research/papers/` (e.g., P-78..P-83 wave 6, P-84..P-89 wave 7). Cross-link new anchors into `Software-Archaeology-Lineage.md` and `primitives/INDEX.md`.
2. **Codebase Compliance** — Run `go build ./...`, `go vet ./...`, `go test ./...`. Verify ponytail directives: Must pattern, try-and-fail strategy, 8-agent graph, Charm TUI, abcoder-mcp default, state.go cohesion, no fmt.Print*, binary compilation.
3. **Paper Sync** — Add new bib entries to `docs/.paper/refs.bib` and `\cite{}` mentions to `docs/.paper/sec_method.tex`.
4. **Handoff** — Write `diary/Sprint-YYYY-MM-DD-Handoff.md` summarizing closed items + commits.

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

- **2026-09-18** — Sprint 2026-09-18: Wave-10 SLM-era anchors (P-102..P-107) persisted; 5 verified new (P-102 SmallCode, P-103 AgentModernize arXiv 2026, P-104 S*, P-105 ChunkKV, P-106 BFCL) + 1 re-anchor slot (P-107 → P-100). Added Wave-10 SLM-era anchors table (10 rows) + Wave-10 anchors list + Wave-9 anchors list. Corrected P-103 venue (arXiv:2605.17535, NOT ICSE 2025). All 7 directive items already verified compliant at 0c1aed5 (no fmt.Print*, abcoder-mcp default, 8-agent graph, try-and-fail strategy registry, state.go centralised, Charm TUI, binary compilation status).
- **2026-09-17** — Sprint 2026-09-17: Ste100 messaging sweep verified clean (zero marketing jargon in user-facing messages; hedge-language only in code comments describing intent). Externalities audit verified comprehensive (yaml.v3, charm stack, abcoder-mcp, container/ring, filepath, env, flag, json). ADR-V-001 sweep reverted initial over-aggressive sed sweep; prose parentheticals `(PRIM-NN)` / `(P-NN)` retained as standard academic-writing convention. No version-slot drift detected. P-06 re-verified absent (closest analog: RepoTransBench = P-09). All gates green.
- **2026-09-16** — Sprint 2026-09-16: Ponytail inline audit + ADR-C-014 locality documentation strengthened via producer-file backlink headers. Dead Charm `errorStyle` removed from `internal/tui/tui.go`. Wave-10 deferred: wave-9 (P-96..P-101) saturated the reasoning-anchors set; next wave launches when new SLM-era mechanisms require anchors. Method entry-point unchanged from Sprint 2026-09-15.
- **2026-09-15** — Sprint 2026-09-15: Added wave 9 SLM-era anchors (P-96..P-101) to §7 SLM-Era Anchors table + cross-cutting list (all 6 verified). Cross-linked into lineage matrix + primitives INDEX.
- **2026-09-14** — Sprint 2026-09-14: Added wave 8 SLM-era anchors (P-90..P-95) to §7 SLM-Era Anchors table + cross-cutting list (all 6 verified).
- **2026-09-13** — Sprint 2026-09-13: Restructured as entry point (Section 0 Quick Start, Section 7 SLM-Era Anchors with verified/unverified status, Section 8 Cross-References, Section 9 Last Updated).
- **2026-09-07** — Sprint 2026-09-07: Centralised state.go into `internal/compiletime/state.go` (ADR-C-014 Locality of Behaviour).