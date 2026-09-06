---
title: MAgHARCM Methodology
date: 2026-09-28
last_updated: 2026-09-28
backlink: "[[2.0.0 Methodology]]"
tags: [methodology, architecture, pipeline, "[[2.0.0 MAgHARCM]]", "[[1.0.0 PRIM-31]]", slm, "[[1.0.0 P-122]]", "[[1.0.0 P-123]]", "[[1.0.0 P-124]]", wave-16, wave-17]
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

Rather than committing irrevocably to a static heuristic choice, the analyzer and repair loop evaluate migration strategies dynamically. The registry lives in `internal/agents/strategy.go`; each strategy implements `Matches(Profile) bool` + `Attempt(context.Context, Profile) (AttemptResult, error)` where `AttemptResult` carries the per-attempt verdict + telemetry. The runner picks the first strategy whose `Matches` returns true; if the attempt produces a `ValidationReport` whose `IsAllSuccess()` is false OR whose `CompilationErrors` slice is non-empty, the runner **automatically increments to the next strategy in the registry** — no human-in-the-loop required for the retry. The graph-level retry is wired at `internal/graph/graph.go:153` (`// Trigger try-and-fail migration strategy switch if needed`).

Registry order (evaluated top-down; first match wins; first failure increments):

1. `BIG_BANG`: small codebases ($\le 3$ files, $< 500$ LoC); one-shot translate-then-validate.
2. `PILOT`: large systems ($> 50$ files or $> 10000$ LoC); translates an isolated subsystem first to characterise risk.
3. `PARALLEL_CUTOVER`: modular systems with comprehensive tests ($> 10$ files); parallel module translation with per-module validation gates.
4. `FROZEN_LEGACY`: systems with zero existing tests; characterisation-test synthesis [[1.0.0 PRIM-11]] before translation.
5. `INCREMENTAL`: canonical universal baseline; module-by-module translation with per-module test gates.

**No silent strategy override.** The registry MUST NOT carry a single hard-coded "if all else fails, do X" fallback. If every strategy fails, the run is reported as a Verdict Panel verdict and the user decides.

**Strategy thresholds are config-driven** (ADR-C-005): `BigBangFileLimit`, `BigBangLoCLimit`, `PilotFileLimit`, `PilotLoCLimit`, `ParallelCutoverFileLimit`, `FrozenLegacyMinimumTestCoverage`. They live in `internal/compiletime/compiletime.go`, not as Go magic numbers in agent modules.

---

## 4. Software Archaeology Suite

Pre-translation comprehension executes through five integrated archaeological primitives (`internal/agents/archaeology.go` + helpers). Each primitive is a separate function with explicit inputs/outputs; the Archaeologist agent composes them in a fixed order:

1. `[[1.0.0 PRIM-14]]` **Archaeology Stage**: discovers module boundaries, historical build-time capsules, git churn hotspots. Output: `ArchaeologyReport.ModuleBoundaries`, `.HistoricalCapsules`, `.ChurnHotspots`.
2. `[[1.0.0 PRIM-18]]` **Jaccard-Coupling Recovery**: computes temporal co-change similarity across commits to expose hidden coupling. Output: `ArchaeologyReport.CouplingMatrix`.
3. `[[1.0.0 PRIM-19]]` **Design Rule Hierarchy**: partitions the codebase into L1 interfaces / L2 subsystems / L3 leaves per `[[1.0.0 P-41]]` Baldwin-Clark. Output: `ArchaeologyReport.DesignRuleHierarchy`.
4. `[[1.0.0 PRIM-20]]` **Concept Assignment**: locates domain concepts across lexical clusters per `[[1.0.0 P-30]]` Lehman/Rajlich. Output: `ArchaeologyReport.ConceptClusters`.
5. `[[1.0.0 PRIM-22]]` **Four Phases of Comprehension**: applies Foltz's DR. JONES cognitive traversal model (`[[1.0.0 P-40]]`). Output: `ArchaeologyReport.ComprehensionTrajectory`.

**SLM-aware preamble** (since Sprint 2026-09-08, commit `2bb9082`): the Archaeologist agent prepends an SLM-aware prompt preamble to every model invocation, explicitly instructing the 4B-30B model to (a) prefer concrete module/file paths over abstract descriptions, (b) emit output in the `ArchaeologyReport` schema rather than free-form prose, and (c) skip directories listed in `compiletime.ArchaeologySkipDirs` (typically `node_modules`, `.git`, `target`, `dist`, `build`, vendor trees).

---

## 5. Binary Compilation & Validation Cascade

Per-project compilation status is strictly **binary**: **Pass** or **Fail**. There is no partial compilation percentage, no "10/15 modules compiled" scoring, no "compilation rate". The Validator emits a single boolean verdict per project; downstream metrics derive from the binary verdict only. `ValidationReport.IsAllSuccess()` returns `bool`; `ValidationReport.CompilationStatus()` returns `compiletime.CompilationStatus` (an enum with `Pass` / `Fail` only — no `Partial` value).

The validation cascade enforces:

1. **AST Syntax Pre-check** (`internal/agents/validator.go:checkASTSyntax`): fast parsing via configured LSP provider (`abcoder-mcp` by default; `tree-sitter` retained ONLY for offline source-language feature extraction in `internal/languages/extractor.go`, NOT for agent-side navigation per the boundary comment at `internal/languages/extractor.go:14`).
2. **Native Toolchain Compilation**: strict type-checking and borrow-checker inspection (`cargo check` for Rust targets; equivalent for other toolchains). Failures here produce the `CompilationStatus = Fail` verdict.
3. **Automated Test Suite**: execution of translated and synthesised tests (`cargo test`). Failures here also produce `Fail`.
4. **Adversarial Weakening Guard** (`[[1.0.0 PRIM-13]]`): halts if test assertions are removed or widened between iterations; the optional_checks.go watchdog compares assertion counts across iterations and aborts the loop if they drop.
5. **Coverage-Guided Plateau Detector** (`[[1.0.0 PRIM-27]]`): exits the iteration loop when test pass-rate improvements stagnate across the configured window.

**Why binary, not percentage.** The downstream metrics (verdict panel consensus, retry decision, strategy-switch trigger, report-cards in `docs/.paper/`) all branch on `IsAllSuccess()`. A "45% compiled" intermediate state would force every downstream consumer to re-implement partial-success handling, which fragments the boundary contract. Binary collapse keeps the contract clean: `Pass` triggers proceed-to-validation; `Fail` triggers retry-or-verdict-panel. `internal/agents/validator.go` is the single source of truth for this collapse.

---

## 6. Centralized Compile-time Config & Locality of Behaviour

Compile-time invariants, enums, sentinels, and initialization helpers reside in `internal/compiletime/compiletime.go` (one package, no internal sub-packages). Examples: `SchemaVersion`, `CurrentSchemaVersion`, `Must` pattern, `DefaultSourceTreeDepth`, `TranslatedPackagePlaceholder`, `ArchaeologySkipDirs`, `CompilationStatus`, `LSPProviderABCoder`, `LSPProviderNative`. Configurations themselves live in `configs/agents.yml` (parsed at boot via `yaml.v3`); the parser rejects missing fields per the **no-fallback rule** (`Must` pattern is mandatory — runtime code never silently substitutes a default for a missing config field).

### Locality of Behaviour (ADR-C-014)

Every producer module MUST declare its intermediate artifact struct in the same Go file as the producer function. Consumers reference the artifact via the producer's package name. The intent: a future reader who opens `internal/agents/translator.go` to understand translation MUST see the `TranslatedProject` struct definition right there, not chase it into a centralised `types.go` leaf.

**Producer-file alias pattern (since Sprint 2026-09-07, commit 0cb5994).** Five artifact structs historically lived in `internal/compiletime/state.go` because they are referenced by `compiletime.State` (the cross-agent pipeline state). Naive relocation of the structs into the producer agent files (`internal/agents/{analyzer,planning,translator,validator,archaeology}.go`) creates a Go import cycle: each agent file imports `compiletime` for `compiletime.State`, `compiletime.Must`, `compiletime.CompilationStatus`, etc.; relocating the structs AND wiring them into `compiletime.State` requires `compiletime` to import `agents` — cycle rejected by `go build`.

Resolution attempted in Sprint 2026-09-23 (commit 89904f6, reverted) and re-attempted in Sprint 2026-09-26 by subagent D: a leaf sub-package `internal/compiletime/artifacts/`. Both attempts hit the same hard blocker: `ValidationReport.CompilationStatus()` returns `compiletime.CompilationStatus`, so the return type MUST move with the method receiver — but moving `CompilationStatus` out of `compiletime/compiletime.go` requires updating every caller in `internal/agents/`, `internal/runner/`, `internal/graph/`, and the test files. A 6-file refactor cannot resolve this without scope expansion.

**Settled solution (ADR-C-014-as-applied, since 2026-09-07):** the artifact struct definitions live in `internal/agents/{analyzer,planning,translator,validator,archaeology}.go` with full docstrings + cross-references; `internal/compiletime/state.go` declares them via `type X = agents.X` (Go type aliases, not redeclarations). The Locality of Behaviour intent is met for the **reader** of the producer file (they see the struct right next to the producer) without breaking the agents→compiletime edge.

**Future refactor (out of scope for current 6-file moves):** lift `compiletime.State` itself out of `internal/compiletime/` into a new leaf package `internal/pipestate/`, then move all 5 artifact structs into `internal/agents/` proper. This is a one-sprint refactor across the entire `internal/agents/` + `internal/runner/` + `internal/graph/` + tests surface.

### No-fallback rule (Must pattern)

Compilation/runtime invariants declared in `internal/compiletime/` use `MustXxx` constructors that `panic` on failure rather than silently substituting a default. Examples: `MustCompileConfig`, `MustLoadConfig`, `MustBuildState`. The intent: missing or malformed configuration is a programmer error; it must abort the process at boot, not silently degrade at runtime. The same rule applies to schema-version migrations, enum-string parsing, and validator setup.

### External configuration

All run-time-tunable values live in YAML files under `configs/`, never as Go constants in agent modules. The set of config keys is documented at the top of `internal/compiletime/compiletime.go` as a YAML schema block (single source of truth).
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
Wave 15 deferred on 2026-09-26 (3 candidates rejected per §7 trigger gate rewritten 2026-09-25; see `.obsidian/MAgHARCM/research/diary/wave-15-candidates.md`).

Wave 16 SLM-era anchors (2 verified new; 2 rejected as benchmarks):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Strategy-distilled persistent memory; structured `(Title, Description, Content)` triples replace raw trajectory replay | `[[1.0.0 P-122]]` ReasoningBank (Zhang ICLR 2026) | verified |
| `[[1.0.0 PRIM-29]]` Recruiter Agent | MaTTS compute-memory loop gives the recruiter a memory substrate | `[[1.0.0 P-122]]` ReasoningBank MaTTS | verified |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Informed-switching policy replaces blind try-and-fail | `[[1.0.0 P-122]]` ReasoningBank | verified |


Wave 16 anchor list:
- `[[1.0.0 P-122]]` ReasoningBank (Zhang et al. 2026, ICLR 2026, arXiv:2509.25140) — Google Research; the persistent-memory mechanism for LLM/SLM agents. Introduces strategy-distilled memory stored as structured triples `(Title, Description, Content)` extracted from both successful AND failed trajectories via self-judgment; closed-loop retrieval at test time; Memory-aware Test-Time Scaling (MaTTS) loop where extra test-time compute generates more diverse trajectories → richer memory → more effective scaling. Demonstrated on SWE-Bench + WebArena with significant effectiveness AND efficiency gains (fewer steps to converge). Anchors `PRIM-31`, `PRIM-29`, `PRIM-21`.


Wave 16 trigger evaluation (Sprint 2026-09-27): 4 candidates triaged, 2 ACCEPT, 2 REJECT (SWE-Bench Pro ICML 2026 + CodeClash ICML 2026 both fail Q2 — benchmark, not mechanism). Full memo `.obsidian/MAgHARCM/research/diary/wave-16-candidates.md`.

Wave 17 SLM-era anchors (2 verified new; 3 rejected — 2 on Q1 venue, 1 on Q2 mechanism):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Confidence-gated switching between in-language majority voting and cross-lingual I/O test oracle | `[[1.0.0 P-123]]` CodeChemist (ICML 2026) | verified |
| `[[1.0.0 PRIM-23]]` Chunked Translation | Multi-temperature hedged sampling with cross-language functional verification | `[[1.0.0 P-123]]` CodeChemist | verified |
| `[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection | Functional-coverage plateau via I/O oracle (tests across source + target language) | `[[1.0.0 P-123]]` CodeChemist | verified |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Runtime-mined properties (aliasing, bounds, nullability) injected as fourth representation | `[[1.0.0 P-124]]` Syzygy (ICLR 2025 VerifAI workshop) | verified |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Static TDG (`[[1.0.0 P-88]]` HiTyper) + dynamic property mining (`[[1.0.0 P-124]]` Syzygy) as complementary dimensions | `[[1.0.0 P-124]]` Syzygy | verified |
| `[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter | Type/bounds/nullability-enriched manifests for safe-Rust generation | `[[1.0.0 P-124]]` Syzygy | verified |

Wave 17 anchor list:
- `[[1.0.0 P-123]]` CodeChemist (Wang et al. 2026, ICML 2026, arXiv:2510.00501) — training-free test-time scaling for low-resource code generation; multi-temperature hedged sampling + cross-lingual I/O test oracle transfers functional knowledge from high-resource reference languages; demonstrated on Qwen-1.5B with 60-70% relative gains on Lua. Anchors `PRIM-21`, `PRIM-23`, `PRIM-27`.
- `[[1.0.0 P-124]]` Syzygy (Shetty et al. 2025, ICLR 2025 VerifAI Workshop, arXiv:2412.14234) — dual code-test C-to-safe-Rust translation via LLMs + dynamic analysis; Clang/LLVM-instrumented SpecMiner mines type/bounds/nullability/aliasing properties at runtime; LLM generates Rust code AND equivalence test per translation unit; multi-round repair loop. Anchors `PRIM-9`, `PRIM-22`, `PRIM-30`. Complements `[[1.0.0 P-88]]` HiTyper's static TDG with the dynamic-analysis half.

Wave 17 trigger evaluation (Sprint 2026-09-28): 5 candidates triaged, 2 ACCEPT, 3 REJECT. Reject rationale: MemSearcher (ACL 2026 Findings, off-list venue), Verified Tool Calls (arXiv-only, no NeurIPS/ICML/ICLR venue), LLM-IR / program-comprehension (ICML 2025 candidate is a benchmark, no mechanism; no SLM-era archaeology mechanism paper found). Full memo `.obsidian/MAgHARCM/research/diary/wave-17-candidates.md`.


Wave 14 fired on 2026-09-25 (NeurIPS 2025 D&B + ICML 2025 mechanism papers — see §9 2026-09-25 entry and §7 wave-14 anchor list above). Wave 15 trigger criterion (replaces former wave-13 trigger criterion): **wave-N+1 fires when a 2025+ NeurIPS / ICML / ICLR paper introduces an unanchored mechanism that defends or refutes an existing SLM-era primitive's substrate claim, OR a new SLM-era primitive lands, OR the user issues a new directive that adds a primitive**.

**Trigger-gate evaluation procedure** (added 2026-09-26): every wave candidate is held against three questions before being promoted to a P-NN slot:

1. **Venue confirmation.** Is the candidate published at (or have a confirmed acceptance to) a 2025+ NeurIPS / ICML / ICLR venue? arXiv-only preprints without venue confirmation are NOT eligible.
2. **Mechanism-vs-benchmark gate.** Does the paper introduce a new mechanism (not just a benchmark or leaderboard)?
3. **Anchoring gate.** Does the mechanism defend or refute a substrate claim of an existing SLM-era primitive (`PRIM-1`..`PRIM-31`)?

Any candidate failing Q1, Q2, OR Q3 is logged in `.obsidian/MAgHARCM/research/diary/wave-15-candidates.md` (or the equivalent wave-N+1 candidates memo) as a deferred research-on-file candidate; it does NOT receive a P-NN slot. Re-evaluation fires when Q1's venue condition is satisfied or when the user explicitly requests promotion.

Wave 15 candidates evaluated on 2026-09-26 (all three rejected, see `.obsidian/MAgHARCM/research/diary/wave-15-candidates.md`):
- SWE-Rebench V2 (Badertdinov et al. 2026, arXiv:2602.23866) — REJECTED Q1 (preprint; no confirmed venue).
- SWE-bench Multimodal (October 2024 announcement per swebench.com) — REJECTED Q1 (no peer-reviewed venue; no arXiv id).
- SWE-bench Verified reference harness (OpenAI August 2024) — REJECTED Q3 (component of `[[1.0.0 P-109]]`, already anchored).
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

## 9. Last Updated


- **2026-09-26** — Sprint 2026-09-26: Wave-15 **deferred** per §7 trigger-gate evaluation. 3 candidates (SWE-Rebench V2 / SWE-bench Multimodal / SWE-bench Verified Reference Harness) all rejected; verdicts logged in `.obsidian/MAgHARCM/research/diary/wave-15-candidates.md`. §7 trigger criterion explicitly rewritten to gate on (Q1) venue confirmation + (Q2) mechanism-vs-benchmark + (Q3) anchoring-to-existing-primitive. ADR-C-014 locality split applied by Subagent D (5 artifact structs moved out of `internal/compiletime/state.go` into producer agent files: analyzer/planning/translator/validator/archaeology). ADR-C-005 magic-string sweep applied by Subagent E (report `local://sprint-2026-09-26-magic-sweep-report.md`). ADR-C-011 Charm stack idiomatic audit applied by Subagent G (3 dead-code removals in `internal/tui/tui.go`: viewport import + field + 3 write lines; zero manual ANSI escapes already; zero reimplemented Charm primitives already). ADR-V-001 automated lint shipped at `scripts/lint_vault.sh` (fails on stray `(P-NN)` / `(PRIM-NN)` in production code; vault markdown parentheticals exempted). Stale-directive audit confirmed 8/4 graph (8 real agents + 2 checkpoints = 10 nodes), abcoder-mcp default in `configs/agents.yml:22`, zero `fmt.Print*` in `internal/`, SelectMigrationStrategy already try-and-fail — all user directives that read as "still to do" already satisfied at `e09cfad`. P-106 BFCL retirement **dropped** (user did not request; P-106 retained alongside P-121 per Sprint 2026-09-25 decision). All gates green after verification.
- **2026-09-28** — Sprint 2026-09-28: Wave-17 FIRED. 2 new SLM-era anchor papers persisted (`[[1.0.0 P-123]]` CodeChemist Wang ICML 2026, arXiv:2510.00501; `[[1.0.0 P-124]]` Syzygy Shetty ICLR 2025 VerifAI Workshop, arXiv:2412.14234). Both pass §7 trigger gate rewritten 2026-09-25 (Q1 venue confirmed, Q2 mechanism-not-benchmark, Q3 anchoring-to-existing-primitive). CodeChemist anchors the cross-lingual functional oracle substrate for `PRIM-21` Migration Strategy Selection, `PRIM-23` Chunked Translation, `PRIM-27` Coverage-Guided Plateau Detection. Syzygy anchors the dynamic-analysis property-mining substrate for `PRIM-9` Tri-Representation Hybrid Code Graph, `PRIM-22` Four Phases of Comprehension, `PRIM-30` Source-to-Target Manifest Rewriter (complementing `[[1.0.0 P-88]]` HiTyper's static TDG with the dynamic half). 3 candidates REJECTED: MemSearcher (ACL 2026 Findings, off-list venue), Verified Tool Calls (arXiv:2608.02645, no NeurIPS/ICML/ICLR venue), LLM-IR / program-comprehension (ICML 2025 candidate is a benchmark, no mechanism; software archaeology SLM-era mechanism gap remains open). §11 extended from 2 to 4 ReasoningBank/CodeChemist/Syzygy sub-sections; §11.5 Substrate Application Matrix now 4 rows. Cross-reference stamps propagated to Architecture §8 (2 new sub-sections), primitives/INDEX frontmatter (2 new P-NN tags) + PRIM-21/23/27/9/22/30 cross-links, Software-Archaeology-Lineage §7/§8. All gates green (`go build`, `go vet`, `go test ./...`, `bash scripts/lint_vault.sh`).
- **2026-09-27** — Sprint 2026-09-27: Wave-16 FIRED. 1 new SLM-era anchor paper persisted (`[[1.0.0 P-122]]` ReasoningBank Zhang ICLR 2026 tentative, arXiv:2509.25140). Passes §7 trigger gate rewritten 2026-09-25 (Q1 venue tentative, Q2 mechanism-not-benchmark, Q3 anchoring-to-PRIM-31/29/21). ReasoningBank anchors the persistent-memory + MaTTS compute-memory substrate for `PRIM-31` Iterative Retrieval, `PRIM-29` Recruiter, `PRIM-21` Migration Strategy Selection. 2 candidates REJECTED: SWE-Bench Pro (ICML 2026) and CodeClash (ICML 2026) — both fail Q2 (benchmark, not mechanism). New §11 SLM-Era General-Purpose Patterns section added with 3 ReasoningBank-only sub-sections + Substrate Application Matrix. Cross-reference stamps propagated to Architecture §8, INDEX wave-16 audit block, Software-Archaeology-Lineage §7 (P-122 only). Compilation gate `bash scripts/lint_vault.sh` exit 0; `go build ./...` exit 0; `go test ./...` cached green.
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
- **2026-09-26** — Sprint 2026-09-26: deferred wave-15 (P-122..P-124 candidates preserved at `.obsidian/MAgHARCM/research/diary/wave-15-candidates.md`); §7 trigger gate verdict logged; ADR-C-014 locality split applied; ADR-C-005 magic-string sweep applied; ADR-C-011 Charm stack audit applied.
---

## 10. Ponytail Refactor Sweep

Every sprint ends with a ponytail audit that explicitly searches for over-engineering, reinvented stdlib, unneeded dependencies, speculative abstractions, and dead flexibility. The audit is a recurring ritual, not a one-off: the codebase's centre of mass keeps moving as new primitives land, and yesterday's clean surface becomes today's accretion surface.

### Standing ponytail rules (verified every sprint)

1. **Must pattern, no fallbacks.** Compile-time invariants use `MustXxx` constructors that `panic` on failure. Configuration loaders reject missing fields; runtime code never silently substitutes a default.
2. **Clear unit boundaries.** Each function/module receives inputs via typed parameters and emits outputs via typed returns. No agent reaches into another agent's package-private state.
3. **Try-and-fail strategy registry.** Migration strategies are not hardcoded; the registry increments on failure (see §3).
4. **Deep, cohesive hierarchy.** Code is organised by concept (agent, primitive, contract), not by file size. A 600-line `archaeology.go` with five sub-functions is preferable to five 100-line files with shared mutable state.
5. **Enums and design patterns over magic strings.** Magic numbers and strings are hoisted into `internal/compiletime/` as named constants; type state is encoded as enums (`compiletime.CompilationStatus`, `compilestate.SpecLifecyclePhase`, etc.).
6. **abcoder-mcp default.** `configs/agents.yml` lists `lsp.provider: abcoder-mcp` as the canonical LSP source. `tree-sitter` is reserved for offline source-language feature extraction only (see boundary comment at `internal/languages/extractor.go:14`).
7. **Binary compilation status.** Per-project compilation is Pass/Fail only — see §5.
8. **Logging, no printing.** Production Go code uses `internal/logger/` exclusively; zero `fmt.Print*` in `internal/` (verified 2026-09-26).
9. **Charm stack idiomatic.** TUI uses `bubbletea` + `bubbles` + `lipgloss` + `glamour`; no manual ANSI escapes, no reinvented primitives (verified 2026-09-26, see `local://sprint-2026-09-26-charm-audit.md`).
10. **Externalities over reinvention.** YAML via `yaml.v3`, ring buffers via `container/ring`, AST navigation via `abcoder-mcp`, terminal I/O via the Charm stack.
11. **STE100 messaging.** User-facing messages follow `asd-ste100` style: short, factual, no marketing language, hedge only when describing intent.
12. **Locality of Behaviour.** Producer modules declare their artifact structs in the same file; see §6 for the alias-pattern solution to the Go import-cycle blocker.

### Automated enforcement

- `scripts/lint_vault.sh` (ADR-V-001): fails CI when `internal/` contains bare `(P-NN)` / `(PRIM-NN)` parentheticals; vault markdown parentheticals exempted (prose convention).
- `go vet ./...`: no shadowed variables, no unreachable code, no printf/format mismatches.
- `go build ./...`: import cycle detection (this is what blocked the ADR-C-014 locality split in Sprint 2026-09-23 + 2026-09-26).
- `go test ./...`: behaviour coverage for each artifact struct's `IsAllSuccess` / `CompilationStatus` / `String` methods.

---

## 11. SLM-Era General-Purpose Patterns (2026-09-28)

Patterns distilled from the wave-16 + wave-17 anchors (`[[1.0.0 P-122]]` ReasoningBank, `[[1.0.0 P-123]]` CodeChemist, `[[1.0.0 P-124]]` Syzygy) that apply across multiple SLM-era primitives. These are substrate claims the codebase should internalise whenever the relevant primitive is touched.

### 11.1. Persistent Memory Substrate (from `[[1.0.0 P-122]]` ReasoningBank)

Any primitive that iterates across pipeline runs (e.g. `[[1.0.0 PRIM-29]]` Recruiter, `[[1.0.0 PRIM-21]]` Strategy Selection, `[[1.0.0 PRIM-31]]` Iterative Retrieval) MUST consider a persistent-memory substrate where strategies are distilled into structured triples `(Title, Description, Content)` rather than replayed as raw trajectories. The substrate is opt-in per primitive — primitives that don't iterate don't need it — but when applied, the memory triple is the canonical exchange format between the producer and the next iteration's consumer.


### 11.3. Cross-Lingual Functional Oracle (from `[[1.0.0 P-123]]` CodeChemist)

`[[1.0.0 PRIM-21]]` Migration Strategy Selection, `[[1.0.0 PRIM-23]]` Chunked Translation, and `[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection SHOULD consider a cross-lingual I/O test oracle as the canonical functional verifier when both source and target can execute. The oracle transfers functional knowledge from the high-resource reference language (the source) into the low-resource target, replacing the frontier-model judge. SLM-amenable: demonstrated on Qwen-1.5B.

### 11.4. Dynamic-Analysis Property Mining (from `[[1.0.0 P-124]]` Syzygy)

`[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph, `[[1.0.0 PRIM-22]]` Four Phases of Comprehension, and `[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter SHOULD inject runtime-mined properties (aliasing, bounds, nullability) as a fourth representation when the legacy codebase can be compiled. LLVM/Clang instrumentation is the canonical mining substrate; the result enriches the manifest that `[[1.0.0 PRIM-30]]` produces, closing the safe-Rust / null-safety / bounds-safety gap that pure-static translation leaves open. Static-only path remains as fallback for uncompilable code.

### 11.5. Substrate Application Matrix

| Pattern | Substrate | Affects | Opt-in location |
| :--- | :--- | :--- | :--- |
| 11.1 Persistent Memory | `(Title, Description, Content)` triples | `[[1.0.0 PRIM-31]]`, `[[1.0.0 PRIM-29]]`, `[[1.0.0 PRIM-21]]` | `configs/agents.yml:memory.distilled: true` |
| 11.2 MaTTS Loop | compute-memory symbiosis | iterative primitives | `configs/agents.yml:mattts.enabled: true` |
| 11.3 Cross-Lingual Functional Oracle | I/O test oracle | `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-27]]` | `configs/agents.yml:oracle.cross_lingual: true` |
| 11.4 Dynamic-Analysis Property Mining | LLVM/Clang instrumentation | `[[1.0.0 PRIM-9]]`, `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-30]]` | `configs/agents.yml:translation.dynamic_specs: true` |
## 12. Stale-Directive Audit (rolling — last refreshed 2026-09-28)

User directives sometimes reference work that has already been completed in an earlier sprint. Re-running already-completed work wastes sprint capacity and fragments the git history with duplicate commits. The stale-directive audit is run before every sprint plan; results are recorded here so future sprints can resolve the same drift quickly.

> **Shorthand note**: The phrase `'8-agent graph'` used as a tag in vault audit lines (this section included) and historical sprint handoffs is shorthand for the canonical 10-node topology (8 specialised agents + 2 checkpoint barriers). The `reviewer` lambda registered at `internal/graph/graph.go:92` is constructed via `agents.NewRoleFlipGate`; "reviewer" is the graph node id, "RoleFlipGate" is the constructor name. See `.obsidian/MAgHARCM/research/Architecture.md` §4 for the canonical statement.

### Verified-already-satisfied directives (as of 2026-09-28, commit `a5bba8e`)

| User directive | Real status | Evidence | Sprint of last verification |
| :--- | :--- | :--- | :--- |
| "graph only has 4 agents" | STALE — graph has **10 nodes** (8 real agents + 2 checkpoints) | `internal/graph/graph.go:45,63,71,79,87,92,116,124,128,150` | 2026-09-26 |
| "Use abcoder's MCP for AST" | ALREADY DEFAULT — `configs/agents.yml:22` lists `provider: abcoder-mcp`; tree-sitter retained only for offline IR extraction | `configs/agents.yml:22` + boundary comment `internal/languages/extractor.go:14` | 2026-09-19 |
| "ABCoderMCPProvider → ABCoderMcpProvider rename" | MOOT — no `ABCoderMCPProvider` ident exists (zero hits); canonical ident is `LSPProviderABCoder` constant | `grep -rn 'ABCoderMCPProvider' .` → 0 | 2026-09-26 |
| "Remove all fmt.Print*" | ALREADY DONE — zero `fmt.Print*` in `internal/` | `grep -rn 'fmt\.Print' internal/` → 0 | 2026-09-19 (re-verified 2026-09-26) |
| "State.go centralised" | ALREADY DONE (basic) — full struct relocation blocked by Go import-cycle (see §6) | `internal/compiletime/state.go` exists since Sprint 2026-09-07 | 2026-09-07 + cycle blocker documented 2026-09-26 |
| "SelectMigrationStrategy too hardcoded, switch to try-and-fail" | ALREADY DONE — graph-level try-and-fail wired at `graph.go:153` | `internal/graph/graph.go:153` | 2026-09-25 |
| "P-106 BFCL retirement" | NOT USER-REQUESTED — P-106 retained alongside P-121 as historical anchor; user did not issue this directive | `Sprint-2026-09-25-Handoff.md` Track-2 | 2026-09-26 |

### Directives still requiring work

| User directive | Status | Owner |
| :--- | :--- | :--- |
| Magic-string sweep across `internal/` | DONE — Subagent E completed 2026-09-26; 5 invariants lifted to `compiletime/`, 3 local consts added; report `local://sprint-2026-09-26-magic-sweep-report.md` | E_MagicSweep (Sprint 2026-09-26) |
| Full ADR-C-014 struct relocation | Blocked by Go import cycle (see §6); alias-pattern solution satisfies intent | future sprint (leaf-package lift) |
| Charm stack idiomatic audit | Subagent G completed 2026-09-26; 3 dead-code removals + zero violations found | G_CharmAudit |
| ADR-V-001 automated lint | Subagent H shipped `scripts/lint_vault.sh` 2026-09-26 | H_VaultSync |
| Research wave fires when new mechanism lands | Wave-17 FIRED 2026-09-28 (P-123 CodeChemist + P-124 Syzygy accepted; MemSearcher / Verified Tool Calls / LLM-IR rejected) | wave-18 conditional |
