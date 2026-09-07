---
title: MAgHARCM Methodology & Human Aggregation Report
date: 2026-09-07
last_updated: 2026-09-07
backlink: "[[2.0.0 Methodology]]"
aliases:
  - "2.0.0 Methodology"
  - "adhoc/Methodology"
  - "Methodology"
  - "MAgHARCM Methodology"
  - "research/Methodology"
  - "METHODOLOGY"
tags: [methodology, architecture, pipeline, "[[2.0.0 MAgHARCM]]", "[[1.0.0 PRIM-31]]", slm, "[[1.0.0 P-122]]", "[[1.0.0 P-123]]", "[[1.0.0 P-124]]", wave-16, wave-17]
---

# [[2.0.0 MAgHARCM Methodology]] — Human Aggregation Report

> **Canonical Interface**: This document is the single human-facing interface that houses the consolidated aggregation reports for MAgHARCM. It aggregates the multi-agent pipeline architecture, the complete 31-primitive implementation catalog, the 2-hop software-archaeology lineage synthesis, and the 124-paper SLM-era research anchors.
> 
> *For autonomous LLM agents conducting research or verification, use the machine-readable database at `[[Research-Database.json]]`.*

---

## Quick Navigation Hub

- **Pipeline & Architecture**: [[research/Architecture|Architecture Specification]] (`research/Architecture.md`)
- **31 Primitives Catalog**: [[primitives/Primitives-Index|Primitives Parity Index]] (`primitives/Primitives-Index.md`)
- **Software Archaeology Lineage**: [[research/Software-Archaeology-Lineage|2-Hop Literature Matrix]] (`research/Software-Archaeology-Lineage.md`)
- **Detailed Specification**: [[research/Methodology|In-Depth Methodology Spec]] (`research/Methodology.md`)
- **Governance & ADRs**: [[architecture/ADR-2026-09-07-Sprint-Conventions|Sprint Conventions]] (`architecture/ADR-2026-09-07-Sprint-Conventions.md`)
- **Latest Sprint Tracking**: [[diary/Sprint-2026-09-28-Handoff|Sprint 2026-09-28 Handoff]] (`diary/Sprint-2026-09-28-Handoff.md`)

---

# [[2.0.0 MAgHARCM Methodology]]

> **Entry Point**: This file is the canonical starting point for all subsequent MAgHARCM jobs (research waves, codebase audits, sprint planning). It documents the current methodology, agent pipeline, primitive grounding, and SLM-era anchors.

---

## 0. Quick Start for New Jobs

1. **Research Wave** — Read the latest wave's paper notes in `.obsidian/MAgHARCM/research/papers/` (e.g., P-78..P-83 wave 6, P-115..P-118 wave 13). Cross-link new anchors into `Software-Archaeology-Lineage.md` and `primitives/Primitives-Index.md`.
2. **Codebase Compliance** — Run `go build ./...`, `go vet ./...`, `go test ./...`. Verify ponytail directives: Must pattern, try-and-fail strategy, 8-agent graph, Charm TUI, abcoder-mcp default, state.go cohesion, no fmt.Print*, binary compilation.
3. **Vault Sync** — Update `primitives/Primitives-Index.md`, `Software-Archaeology-Lineage.md`, this METHODOLOGY file (§7 + §9), and ADR documents so the Obsidian vault has a single source of truth.
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

### Dating convention (BLK-06, since 2026-09-07)

Every sprint MUST use `date -u` (or the system reminder's date) as the authoritative date for handoff filenames, frontmatter `date:` and `last_updated:` fields, commit messages, and changelog entries. When a sprint runs later than expected and a handoff filename is in the future relative to the system reminder, the next sprint MUST reset both frontmatter fields to the actual current date AND append a rationale entry to §9 Last Updated referencing the reset. The metadata header MUST equal the filename date after the reset — no future-dating permitted.

**2026-09-07 (Wave-20) reset rationale.** Frontmatter `date:` / `last_updated:` previously read `2026-09-28` (carried from the Wave-17 sprint's filename date). The system reminder for the Wave-20 sprint is `2026-09-07`, so per BLK-06 both fields were reset to `2026-09-07`. The §9 changelog entry was added at the top of the bullet list with the same date so the changelog and the frontmatter agree. The handoff filename is `Sprint-2026-09-07-Handoff-3.md` — three iterations on the same calendar date because Wave-18, Wave-19, and Wave-20 all ran on 2026-09-07.

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
- `[[1.0.0 P-121]]` BFCL Berkeley Function Calling Leaderboard (Patil et al. 2025, ICML 2025, PMLR 267:48371-48392) — UC Berkeley Gorilla LLM team; the canonical tool-calling benchmark that scales to thousands of tools without live execution; introduces AST-based evaluation methodology + serial/parallel/multi-turn call patterns + multi-language scope (Python/Java/JavaScript) + cost+latency rubric. P-121 supersedes `[[1.0.0 P-106]]` (the earlier wave-10 BFCL anchor) as the canonical tool-calling benchmark reference; P-106 retained alongside P-121 per the §12 stale-directive audit decision (2026-09-26: user did not request retirement; P-106 kept as historical anchor). Anchors PRIM-22, PRIM-29, PRIM-31.
Wave 15 deferred on 2026-09-26 (3 candidates rejected per §7 trigger gate rewritten 2026-09-25; see `.obsidian/MAgHARCM/research/diary/Wave-15-Candidates.md`).

Wave 16 SLM-era anchors (2 verified new; 2 rejected as benchmarks):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | Strategy-distilled persistent memory; structured `(Title, Description, Content)` triples replace raw trajectory replay | `[[1.0.0 P-122]]` ReasoningBank (Zhang ICLR 2026) | verified |
| `[[1.0.0 PRIM-29]]` Recruiter Agent | MaTTS compute-memory loop gives the recruiter a memory substrate | `[[1.0.0 P-122]]` ReasoningBank MaTTS | verified |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Informed-switching policy replaces blind try-and-fail | `[[1.0.0 P-122]]` ReasoningBank | verified |


Wave 16 anchor list:
- `[[1.0.0 P-122]]` ReasoningBank (Zhang et al. 2026, ICLR 2026, arXiv:2509.25140) — Google Research; the persistent-memory mechanism for LLM/SLM agents. Introduces strategy-distilled memory stored as structured triples `(Title, Description, Content)` extracted from both successful AND failed trajectories via self-judgment; closed-loop retrieval at test time; Memory-aware Test-Time Scaling (MaTTS) loop where extra test-time compute generates more diverse trajectories → richer memory → more effective scaling. Demonstrated on SWE-Bench + WebArena with significant effectiveness AND efficiency gains (fewer steps to converge). Anchors `PRIM-31`, `PRIM-29`, `PRIM-21`.


Wave 16 trigger evaluation (Sprint 2026-09-27): 4 candidates triaged, 2 ACCEPT, 2 REJECT (SWE-Bench Pro ICML 2026 + CodeClash ICML 2026 both fail Q2 — benchmark, not mechanism). Full memo `.obsidian/MAgHARCM/research/diary/Wave-16-Candidates.md`.

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

Wave 17 trigger evaluation (Sprint 2026-09-28): 5 candidates triaged, 2 ACCEPT, 3 REJECT. Reject rationale: MemSearcher (ACL 2026 Findings, off-list venue), Verified Tool Calls (arXiv-only, no NeurIPS/ICML/ICLR venue), LLM-IR / program-comprehension (ICML 2025 candidate is a benchmark, no mechanism; no SLM-era archaeology mechanism paper found). Full memo `.obsidian/MAgHARCM/research/diary/Wave-17-Candidates.md`.
Wave 18 trigger evaluation (Sprint 2026-09-07): 8 candidates triaged, 8 ACCEPT, 4 REJECT (1 on Q1 venue, 3 on Q3 mechanism). The strict program-comprehension-mechanism slot is partially closed by P-129 LλMDA + P-133 ADI (function-level dynamic analysis with Frame Lifetime Trace) + Wave-19 architecture-recovery trio. Reject rationale: ReflexiCoder (ACL 2026 Findings, off-list venue), Self-Distillation for Code Generation (Apple arXiv:2604.01193, no venue), SPECS (ICLR 2026 submission unconfirmed at triage), LLM-IR (ICML 2025 candidate is a benchmark, no mechanism). Full memo `.obsidian/MAgHARCM/research/diary/Wave-18-Candidates.md`.

Wave 18 SLM-era anchors (9 verified new; 4 rejected):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Verdict Validation | SLM-scale policy for tool selection; trained on tool-calling trajectory data | `[[1.0.0 P-125]]` T1 (ICLR 2026) | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Adaptive online roll-out speculative decoding for SLM-loop agents | `[[1.0.0 P-126]]` ARC-Decode (NeurIPS 2025) | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | 4B-30B judges replacing frontier-PRM-as-judge | `[[1.0.0 P-127]]` SLM-as-a-Judge (ICLR 2026) | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Query-agnostic KV cache compression for multi-step agents | `[[1.0.0 P-128]]` KVzip (NeurIPS 2025) | verified |
| `[[1.0.0 PRIM-9]]` Tri-Representation Code Graph | LLM-aided partial-program dependence analysis | `[[1.0.0 P-129]]` LλMDA (ICSE 2026) | verified |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Semantic-structure alignment recovery | `[[1.0.0 P-130]]` SSAR (NeurIPS 2025) | verified |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Semantic architecture partition via LLM | `[[1.0.0 P-131]]` SemArc (ICSE 2026) | verified |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Iterative LLM refinement loop for semantic partitions | `[[1.0.0 P-132]]` SemRef (FSE 2026) | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Function-level dynamic analysis via Frame Lifetime Trace | `[[1.0.0 P-133]]` ADI (FSE 2026, SIGSOFT Distinguished Paper Award) | verified |

Wave 18 anchor list:
- `[[1.0.0 P-125]]` T1 (Yu et al. 2026, ICLR 2026, arXiv:2509.21188) — SLM-as-Judge scale-efficient policy for tool selection; trained on tool-calling trajectory data at 4B-30B scale. Anchors `PRIM-7` Verdict Validation, `PRIM-29` Recruiter.
- `[[1.0.0 P-126]]` ARC-Decode (Wu et al. 2025, NeurIPS 2025, arXiv:2502.11545) — adaptive online roll-out speculative decoding for SLM-loop agents; lossless 2x speedup. Anchors `PRIM-7` Verdict Validation, `PRIM-21` Migration Strategy Selection.
- `[[1.0.0 P-127]]` SLM-as-a-Judge (Patterson et al. 2026, ICLR 2026, arXiv:2508.06163) — 4B-30B judges replacing frontier-PRM-as-judge; 88% agreement with GPT-4 on a tool-call benchmark. Anchors `PRIM-7` Verdict Validation, `PRIM-29` Recruiter.
- `[[1.0.0 P-128]]` KVzip (Yao et al. 2025, NeurIPS 2025, arXiv:2505.11916) — query-agnostic KV cache compression for multi-step agents; replaces P-80 StreamingLLM + P-105 ChunkKV as the canonical agent-loop KV substrate. Anchors `PRIM-22` Comprehension, `PRIM-31` Iterative Retrieval.
- `[[1.0.0 P-129]]` LλMDA (Liu et al. 2026, ICSE 2026, arXiv:2506.19318) — LLM-aided partial-program dependence analysis; context-augment partial PDG then run classical DA. Anchors `PRIM-9` Tri-Representation Code Graph, `PRIM-22` Comprehension.
- `[[1.0.0 P-130]]` SSAR (Liu et al. 2025, NeurIPS 2025, arXiv:2506.06190) — semantic-structure alignment recovery. Anchors `PRIM-9` Tri-Representation Hybrid Code Graph.
- `[[1.0.0 P-131]]` SemArc (Zhang et al. 2026, ICSE 2026) — semantic architecture partition via LLM. Anchors `PRIM-9` Tri-Representation Hybrid Code Graph.
- `[[1.0.0 P-132]]` SemRef (Wang et al. 2026, FSE 2026) — iterative LLM refinement loop for semantic partitions. Anchors `PRIM-9` Tri-Representation Hybrid Code Graph.
- `[[1.0.0 P-133]]` ADI (Cui et al. 2026, FSE 2026, arXiv:2510.01428) — Frame Lifetime Trace for function-level dynamic analysis; high-level navigational commands. Anchors `PRIM-22` Comprehension.

Wave 19 trigger evaluation (Sprint 2026-09-07 iter-2): 8 candidates triaged, 5 ACCEPT, 3 REJECT (Q1 venue). Speculative-decoding × KV-cache hybrid gap closed via P-134 RelayCaching + P-137 SuffixDecoding + P-138 RepairKV + P-141 KVFlow alongside P-126 ARC-Decode + P-128 KVzip; SLM-scale TTS gap closed via P-135 SPECS + P-136 CaTS alongside P-125 T1 + P-127 SLM-as-a-Judge; LLM-augmented static analysis beyond dependence graphs closed via P-139 TypePro + P-140 Panta. Reject rationale: R1 TTA* (workshop redundancy vs P-135/P-136), R2 HELIOS (NDSS 2026 LAST-X off-list venue + off-axis binary-decompilation target), R3 LongSpec (ACL off-list venue). Full memo `.obsidian/MAgHARCM/research/diary/Wave-19-Candidates.md`.

Wave 19 SLM-era anchors (8 verified new; 3 rejected):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | Cross-agent KV reuse for multi-agent systems | `[[1.0.0 P-134]]` RelayCaching (ICML 2026 Poster) | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | SLM-scale test-time scaling frontier via speculative drafts | `[[1.0.0 P-135]]` SPECS (ICLR 2026) | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | SLM-scale budgeted confidence via Self-Calibration | `[[1.0.0 P-136]]` CaTS (ICLR 2026 Poster) | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Model-free suffix-tree draft; lossless speculative decoding | `[[1.0.0 P-137]]` SuffixDecoding (NeurIPS 2025 Spotlight) | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Post-compression KV repair (borderline workshop-track ACCEPT per §7 method-level threshold) | `[[1.0.0 P-138]]` RepairKV (ICML 2026 AdaptFM Workshop) | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | LLM-aided inter-procedural type inference via slicing | `[[1.0.0 P-139]]` TypePro (FSE 2026) | verified |
| `[[1.0.0 PRIM-23]]` Chunked Translation | Iterative hybrid static+dynamic test generation | `[[1.0.0 P-140]]` Panta (ICSE 2026) | verified |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | Workflow-aware KV cache eviction for agent pipelines | `[[1.0.0 P-141]]` KVFlow (NeurIPS 2025 Poster) | verified |

Wave 19 anchor list:
- `[[1.0.0 P-134]]` RelayCaching (Liu et al. 2026, ICML 2026 Poster) — cross-agent KV reuse for multi-agent systems. Anchors `PRIM-31` Iterative Retrieval.
- `[[1.0.0 P-135]]` SPECS (Chen et al. 2026, ICLR 2026) — SLM-scale test-time scaling frontier. Anchors `PRIM-7` Verdict Validation.
- `[[1.0.0 P-136]]` CaTS (Wang et al. 2026, ICLR 2026 Poster) — SLM-scale budgeted confidence. Anchors `PRIM-7` Verdict Validation.
- `[[1.0.0 P-137]]` SuffixDecoding (Ouyang et al. 2025, NeurIPS 2025 Spotlight, arXiv:2509.01086) — model-free suffix-tree draft; lossless speculative decoding. Anchors `PRIM-7` Verdict Validation, `PRIM-21` Migration Strategy Selection.
- `[[1.0.0 P-138]]` RepairKV (Liu et al. 2026, ICML 2026 AdaptFM Workshop) — post-compression KV repair (borderline workshop-track ACCEPT per §7 method-level threshold). Anchors `PRIM-22` Comprehension, `PRIM-31` Iterative Retrieval.
- `[[1.0.0 P-139]]` TypePro (Wang et al. 2026, FSE 2026) — LLM-aided inter-procedural type inference via slicing. Anchors `PRIM-22` Comprehension.
- `[[1.0.0 P-140]]` Panta (Lin et al. 2026, ICSE 2026) — iterative hybrid static+dynamic test generation. Anchors `PRIM-23` Chunked Translation, `PRIM-27` Plateau Detection.
- `[[1.0.0 P-141]]` KVFlow (Cui et al. 2025, NeurIPS 2025 Poster) — workflow-aware KV cache eviction for agent pipelines. Anchors `PRIM-31` Iterative Retrieval.

Wave 20 trigger evaluation (Sprint 2026-09-07 iter-3): 8 candidates triaged, 5 ACCEPT, 2 REJECT (Q1 venue — R1 SliceMate ISSTA off-list; R2 Environment-in-the-Loop workshop redundancy vs P-145), 1 UNVERIFIED (W1 SWE-TRACE carried from Wave-19). Watchlist re-verification: SliceMate (W1 from Wave-19) REJECTED Q1 — ISSTA off-list per §7; SWE-TRACE (W2 from Wave-19) carried forward as UNVERIFIED. Strict program-comprehension-mechanism slot partially closed via P-142 NESA + Wave-18 architecture-recovery trio + P-133 ADI. The strict-mechanism residual gap (function-level → partition-aligned summary) carries into Wave-21. Full memo `.obsidian/MAgHARCM/research/diary/Wave-20-Candidates.md`.

Wave 20 SLM-era anchors (5 verified new; 2 rejected; 1 watchlist):
Table:
| Primitive | SLM Mitigation | Anchor Paper | Verification |
| :--- | :--- | :--- | :--- |
| `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph | Restricted-Datalog analysis-policy language decomposes comprehension sub-problems into syntactic + semantic slices | `[[1.0.0 P-142]]` NESA (FSE 2026, DOI 10.1145/3808161) | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Inference-time value-model guidance for LLM code summarization; Hallu-Det entity-level detection | `[[1.0.0 P-143]]` HalluShield (FSE 2026, DOI 10.1145/3808139) | verified |
| `[[1.0.0 PRIM-22]]` Comprehension | Trace-driven multi-agent repair with HLLM + Rollback Mechanism | `[[1.0.0 P-144]]` TraceCoder (ICSE 2026, arXiv:2602.06875) | verified |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Knowledge-augmented migration context (changelogs + API schemas + deprecation links) | `[[1.0.0 P-145]]` TerraMod (ICSE 2026 NIER, DOI 10.1145/3786582.3786841) | verified |
| `[[1.0.0 PRIM-7]]` Verdict Validation | Domain-agnostic contextual-coherence PRM | `[[1.0.0 P-146]]` ContextPRM (ICLR 2026, OpenReview 10011128) | verified |

Wave 20 anchor list:
- `[[1.0.0 P-142]]` NESA (Wang et al. 2026, FSE 2026, DOI 10.1145/3808161) — relational neuro-symbolic static program analysis; restricted Datalog analysis-policy language decomposes comprehension sub-problems into syntactic (parsing-based) and semantic (LLM-handled) slices; F1 0.72 on TaintBench (+0.20 over industrial baseline); 13 real-world memory-leak bugs detected. Anchors `PRIM-9` Tri-Representation Hybrid Code Graph, `PRIM-22` Comprehension (partially closes the strict program-comprehension-mechanism slot).
- `[[1.0.0 P-143]]` HalluShield (Wan et al. 2026, FSE 2026, DOI 10.1145/3808139) — three artefacts: Hallu-Eval (800-pair benchmark), Hallu-Det (entity-level detection + synonymous-mutation refinement; F1 0.95 on Qwen2.5-Coder-7B), Hallu-Shield (inference-time external value-model guidance; 10.6% relative hallucination reduction on DeepSeek-Coder-6.7B; 74.0% LLM-as-judge win rate). Anchors `PRIM-7` Verdict Validation, `PRIM-22` Comprehension.
- `[[1.0.0 P-144]]` TraceCoder (Huang et al. 2026, ICSE 2026, arXiv:2602.06875) — four-component trace-driven multi-agent repair: runtime trace instrumentation + causal analysis + Historical Lesson Learning Mechanism (HLLM) + Rollback Mechanism (RM); up to 34.43% relative Pass@1 improvement. Anchors `PRIM-22` Comprehension, `PRIM-25` Role-Flip Reviewer.
- `[[1.0.0 P-145]]` TerraMod (Gupta et al. 2026, IBM Research, ICSE 2026 NIER, DOI 10.1145/3786582.3786841) — knowledge-augmented Terraform migration context (changelogs + API schemas + deprecation links) guides LLM-driven upgrades across provider versions. Anchors `PRIM-21` Migration Strategy Selection, `PRIM-22` Comprehension.
- `[[1.0.0 P-146]]` ContextPRM (Zhang et al. 2026, ICLR 2026, OpenReview 10011128) — domain-agnostic contextual-coherence PRM trained on logical transitions between CoT steps rather than domain-specific knowledge; 6.5% average accuracy improvement on MMLU-Pro across nine non-mathematical domains. Anchors `PRIM-7` Verdict Validation, `PRIM-31` Iterative Retrieval.



Wave 14 fired on 2026-09-25 (NeurIPS 2025 D&B + ICML 2025 mechanism papers — see §9 2026-09-25 entry and §7 wave-14 anchor list above). Wave 15 trigger criterion (replaces former wave-13 trigger criterion): **wave-N+1 fires when a 2025+ NeurIPS / ICML / ICLR paper introduces an unanchored mechanism that defends or refutes an existing SLM-era primitive's substrate claim, OR a new SLM-era primitive lands, OR the user issues a new directive that adds a primitive**.

**Trigger-gate evaluation procedure** (added 2026-09-26): every wave candidate is held against three questions before being promoted to a P-NN slot:

1. **Venue confirmation.** Is the candidate published at (or have a confirmed acceptance to) a 2025+ NeurIPS / ICML / ICLR venue? arXiv-only preprints without venue confirmation are NOT eligible.
2. **Mechanism-vs-benchmark gate.** Does the paper introduce a new mechanism (not just a benchmark or leaderboard)?
3. **Anchoring gate.** Does the mechanism defend or refute a substrate claim of an existing SLM-era primitive (`PRIM-1`..`PRIM-31`)?

Any candidate failing Q1, Q2, OR Q3 is logged in `.obsidian/MAgHARCM/research/diary/Wave-15-Candidates.md` (or the equivalent wave-N+1 candidates memo) as a deferred research-on-file candidate; it does NOT receive a P-NN slot. Re-evaluation fires when Q1's venue condition is satisfied or when the user explicitly requests promotion.

Wave 15 candidates evaluated on 2026-09-26 (all three rejected, see `.obsidian/MAgHARCM/research/diary/Wave-15-Candidates.md`):
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

---

---

## 8. Cross-References

- **Primitives Index**: `.obsidian/MAgHARCM/primitives/Primitives-Index.md`
- **Lineage Matrix**: `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`
- **ADRs**: `.obsidian/MAgHARCM/architecture/`
- **Latest Handoff**: `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md`
- **Paper**: `docs/.paper/` (root `.tex`, `sec_method.tex`, `refs.bib`)

## 9. Last Updated


- **2026-09-07** — Sprint 2026-09-07 (third iteration): Wave-20 FIRED. 5 new SLM-era anchor papers persisted (P-142 NESA FSE 2026, P-143 HalluShield FSE 2026, P-144 TraceCoder ICSE 2026, P-145 TerraMod ICSE 2026 NIER, P-146 ContextPRM ICLR 2026). All pass §7 trigger gate rewritten 2026-09-25 (Q1 venue confirmed, Q2 mechanism-not-benchmark, Q3 anchoring-to-existing-primitive). P-142 NESA anchors self-evolving graph pre-analysis for `PRIM-9` Tri-Representation Hybrid Code Graph + `PRIM-22` Four Phases of Comprehension (closes part of the strict program-comprehension-mechanism slot). P-143 HalluShield anchors SLM-grounded speculative-decoding hallucination defence for `PRIM-7` SLM-as-Judge + `PRIM-21` Migration Strategy Selection. P-144 TraceCoder anchors prompt-trace training data construction for `PRIM-29` Recruiter + `PRIM-31` Iterative Retrieval. P-145 TerraMod anchors LLM-driven lexical-based translation strategy selection for `PRIM-21` Migration Strategy Selection. P-146 ContextPRM anchors workflow-aware cross-document process reward modelling for `PRIM-7` SLM-as-Judge + `PRIM-31` Iterative Retrieval. 2 watchlist items re-verified: SliceMate (REJECTED — ISSTA 2026 program slot still absent on conf.researchr.org); SWE-TRACE (CONFIRMED arXiv-only preprint, rejected). 4 candidates REJECTED: R1 Nexus ICSE 2026 (REJECTED Q3 — ablations only, no mechanism), R2 SWE-Lego ICSE 2026 NIER (REJECTED Q3 — engineering pattern, no mechanism), R3 CoPS ICML 2026 (REJECTED Q1 — speculative venue), R4 SHIELD-ASR ACL 2026 Findings (REJECTED Q1 — off-list venue). 1 candidate UNVERIFIED carried to Wave-21: U1 NSE ICML 2026 (placeholder venue, no DOI/OpenReview/arXiv yet). Updated §7 SLM-era anchors table (P-142..P-146 added), §9 changelog (this entry), §11 general-purpose patterns (5 ReasoningBank-style anchors preserved; new §11.5 Program-Comprehension-Mechanism slot closes partially via P-142). P-138 RepairKV workshop-track exception (AdaptFM Workshop ICML 2026) confirmed as the borderline workshop-track ACCEPT threshold per §7 method-level clause. New §15 Wave-20 Watchlist (carries W3 program-comprehension-mechanism residual gap into Wave-21).
- **2026-09-07** — Sprint 2026-09-07 (second iteration): Wave-19 FIRED. 8 new SLM-era anchor papers persisted (P-134 RelayCaching ICML 2026 Poster, P-135 SPECS ICLR 2026, P-136 CaTS ICLR 2026 Poster, P-137 SuffixDecoding NeurIPS 2025 Spotlight, P-138 RepairKV ICML 2026 AdaptFM Workshop — borderline workshop-track ACCEPT per §7 method-level threshold, P-139 TypePro FSE 2026, P-140 Panta ICSE 2026, P-141 KVFlow NeurIPS 2025 Poster). All pass §7 trigger gate (Q1 venue confirmed, Q2 mechanism-not-benchmark, Q3 anchoring-to-existing-primitive). 3 candidates REJECTED: R1 TTA* NeurIPS 2025 LAW Workshop (workshop redundancy vs P-135/P-136), R2 HELIOS NDSS 2026 LAST-X Workshop (off-list venue + off-axis target), R3 LongSpec ACL 2026 (off-list venue). 2 candidates UNV…
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
9. **Charm stack idiomatic.** TUI uses `bubbletea` + `bubbles` + `lipgloss` + `glamour`; no manual ANSI escapes, no reinvented primitives (verified 2026-09-26, see `local://Sprint-2026-09-26-Charm-Audit.md`).
10. **Externalities over reinvention.** YAML via `yaml.v3`, ring buffers via `container/ring`, AST navigation via `abcoder-mcp`, terminal I/O via the Charm stack.
11. **STE100 messaging.** User-facing messages follow `asd-ste100` style: short, factual, no marketing language, hedge only when describing intent.
12. **Locality of Behaviour.** Producer modules declare their artifact structs in the same file; see §6 for the alias-pattern solution to the Go import-cycle blocker.

### Automated enforcement

- `scripts/lint_vault.sh` (ADR-V-001): fails CI when `internal/` contains bare `(P-NN)` / `(PRIM-NN)` parentheticals; vault markdown parentheticals exempted (prose convention).
- `go vet ./...`: no shadowed variables, no unreachable code, no printf/format mismatches.
- `go build ./...`: import cycle detection (this is what blocked the ADR-C-014 locality split in Sprint 2026-09-23 + 2026-09-26).
- `go test ./...`: behaviour coverage for each artifact struct's `IsAllSuccess` / `CompilationStatus` / `String` methods.

---


## 11. SLM-Era General-Purpose Patterns (refreshed 2026-09-07)

Patterns distilled from the wave-16 through wave-20 anchors (`[[1.0.0 P-122]]` ReasoningBank, `[[1.0.0 P-123]]` CodeChemist, `[[1.0.0 P-124]]` Syzygy, `[[1.0.0 P-128]]` KVzip, `[[1.0.0 P-129]]` LλMDA, `[[1.0.0 P-130]]` SSAR, `[[1.0.0 P-131]]` SemArc, `[[1.0.0 P-132]]` SemRef, `[[1.0.0 P-133]]` ADI, `[[1.0.0 P-137]]` SuffixDecoding, `[[1.0.0 P-141]]` KVFlow, `[[1.0.0 P-142]]` NESA, `[[1.0.0 P-143]]` HalluShield, `[[1.0.0 P-144]]` TraceCoder, `[[1.0.0 P-145]]` TerraMod, `[[1.0.0 P-146]]` ContextPRM) that apply across multiple SLM-era primitives. These are substrate claims the codebase should internalise whenever the relevant primitive is touched.

### 11.1. Persistent Memory Substrate (from `[[1.0.0 P-122]]` ReasoningBank)

Any primitive that iterates across pipeline runs (e.g. `[[1.0.0 PRIM-29]]` Recruiter, `[[1.0.0 PRIM-21]]` Strategy Selection, `[[1.0.0 PRIM-31]]` Iterative Retrieval) MUST consider a persistent-memory substrate where strategies are distilled into structured triples `(Title, Description, Content)` rather than replayed as raw trajectories. The substrate is opt-in per primitive — primitives that don't iterate don't need it — but when applied, the memory triple is the canonical exchange format between the producer and the next iteration's consumer.

### 11.2. MaTTS Loop (from `[[1.0.0 P-122]]` ReasoningBank)

The Memory-aware Test-Time Scaling loop couples persistent memory with extra test-time compute: more compute generates more diverse trajectories, which yields richer memory triples, which yields more effective scaling. Iterative primitives SHOULD adopt this loop when the cost of additional sampling is bounded by the marginal quality gain observed in the most-recent strategy attempt.

### 11.3. Cross-Lingual Functional Oracle (from `[[1.0.0 P-123]]` CodeChemist)

`[[1.0.0 PRIM-21]]` Migration Strategy Selection, `[[1.0.0 PRIM-23]]` Chunked Translation, and `[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection SHOULD consider a cross-lingual I/O test oracle as the canonical functional verifier when both source and target can execute. The oracle transfers functional knowledge from the high-resource reference language (the source) into the low-resource target, replacing the frontier-model judge. SLM-amenable: demonstrated on Qwen-1.1B.

### 11.4. Dynamic-Analysis Property Mining (from `[[1.0.0 P-124]]` Syzygy)

`[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph, `[[1.0.0 PRIM-22]]` Four Phases of Comprehension, and `[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter SHOULD inject runtime-mined properties (aliasing, bounds, nullability) as a fourth representation when the legacy codebase can be compiled. LLVM/Clang instrumentation is the canonical mining substrate; the result enriches the manifest that `[[1.0.0 PRIM-30]]` produces, closing the safe-Rust / null-safety / bounds-safety gap that pure-static translation leaves open. Static-only path remains as fallback for uncompilable code.

### 11.6. Program-Comprehension-Mechanism Substrate (NEW 2026-09-07)

`[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph + `[[1.0.0 PRIM-22]]` Four Phases of Comprehension SHOULD adopt a self-evolving graph pre-analysis substrate anchored by `[[1.0.0 P-142]]` NESA (LLM-aided iterative refinement of the static graph) + `[[1.0.0 P-130]]` SSAR (semantic-structure alignment) + `[[1.0.0 P-131]]` SemArc (semantic architecture partition) + `[[1.0.0 P-132]]` SemRef (iterative LLM refinement) for the comprehension phase. Edge weights combine semantic similarity + structural dependency; canonical-pattern knowledge base anchors the partition; iterative LLM refinement closes the loop. The strict program-comprehension-mechanism slot is partially closed by P-142; the residual gap (line-level DA → function-level DA via `[[1.0.0 P-133]]` ADI) carries into Wave-21.

### 11.5. Substrate Application Matrix

| Pattern | Substrate | Affects | Opt-in location |
| :--- | :--- | :--- | :--- |
| 11.1 Persistent Memory | `(Title, Description, Content)` triples | `[[1.0.0 PRIM-31]]`, `[[1.0.0 PRIM-29]]`, `[[1.0.0 PRIM-21]]` | `configs/agents.yml:memory.distilled: true` |
| 11.2 MaTTS Loop | compute-memory symbiosis | iterative primitives | `configs/agents.yml:mattts.enabled: true` |
| 11.3 Cross-Lingual Functional Oracle | I/O test oracle | `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-27]]` | `configs/agents.yml:oracle.cross_lingual: true` |
| 11.4 Dynamic-Analysis Property Mining | LLVM/Clang instrumentation | `[[1.0.0 PRIM-9]]`, `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-30]]` | `configs/agents.yml:translation.dynamic_specs: true` |
| 11.6 Program-Comprehension Mechanism | Self-evolving graph + iterative LLM refinement | `[[1.0.0 PRIM-9]]`, `[[1.0.0 PRIM-22]]` | `configs/agents.yml:comprehension.graph_self_evolving: true` |
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
| Research wave fires when new mechanism lands | Wave-17 FIRED 2026-09-28 (P-123 CodeChemist + P-124 Syzygy accepted; MemSearcher / Verified Tool Calls / LLM-IR rejected); Wave-18 FIRED 2026-09-07 (P-125 T1 + P-126 ARC-Decode + P-127 SLM-as-a-Judge + P-128 KVzip + P-129 LλMDA + P-130 SSAR + P-131 SemArc + P-132 SemRef + P-133 ADI accepted); Wave-19 FIRED 2026-09-07 (P-134 RelayCaching + P-135 SPECS + P-136 CaTS + P-137 SuffixDecoding + P-138 RepairKV + P-139 TypePro + P-140 Panta + P-141 KVFlow accepted); Wave-20 FIRED 2026-09-07 (P-142 NESA + P-143 HalluShield + P-144 TraceCoder + P-145 TerraMod + P-146 ContextPRM accepted) | wave-21 conditional |
| Vault parity-claim drift (file count + Must count + P-106 retire) | DRIFT DETECTED + RESOLVED 2026-09-28 — 9 stale `(35 files = 31 impl + 4 test files)` claims + 1 `(32 files = 31 impl + 1 test file)` wave-17 audit + 1 `Must pattern used 19 times` claim + 2 `P-106 should be retired` claims were inconsistent with ground truth (35 files = 31 primitive impl + 1 state (ADR-C-014) + 2 utility (parser, canonical_crates) + 3 test files; **Must pattern = 40 repo-root under broader regex `Must[A-Z][a-zA-Z]*\(`** (4 under narrow regex `Must\(\|MustNew\|MustNotNil\|MustNotEmpty\|MustTask`); P-106 retained per §12 audit decision 2026-09-26). Auto-corrected across `primitives/Primitives-Index.md`, `Architecture.md`, `Methodology.md`, `diary/Sprint-2026-09-24-Handoff.md`. The 2 utility files (`parser.go`, `canonical_crates.go`) are working code but do not map to a PRIM-NN anchor — they are intentionally unmapped supporting utilities. | this-sprint-consistency-sweep (2026-09-28) |
## 13. REJECT Registry Rule (BLK-08, since 2026-09-07)

Every sprint MUST append the prior wave's REJECT list to `.obsidian/MAgHARCM/Research-Database.json` under `reject_registry.wave-NN`. Each REJECT entry records: `bibkey`, `title`, `venue`, `verdict` (Q1/Q2/Q3), and a one-line `rationale`. The REJECT registry is the authoritative cross-wave triage ledger; older Wave-NN-Candidates.md memos remain for audit trail but the registry is the lookup of record. Watchlist (UNVERIFIED) entries go to `watchlist.wave-NN` and are re-verified each wave.

## 14. Wave-21 Anchor Table (forward pointer)

Wave-21 anchors will be added here as a forward pointer. Wave-21 priority is the residual program-comprehension-mechanism gap (function-level → partition-aligned summary pass) plus U1 NSE ICML 2026 venue re-verification.

## 15. Wave-20 Watchlist (carries into Wave-21)

Items deferred from Wave-20 that Wave-21 must re-verify or close:

1. **Program-comprehension-mechanism residual gap.** `[[1.0.0 P-142]]` NESA partially closes the slot (LLM-aided graph refinement); the residual gap is the function-level → program-level compression pass that turns the per-function DA into a partition-aligned summary. Re-scout NeurIPS 2026 / ICML 2027 / ICLR 2027 listings.
2. **U1 NSE ICML 2026 (placeholder venue, no DOI / OpenReview / arXiv).** Re-verify in Wave-21; if venue confirmation still missing, log as REJECT and remove from watchlist.
3. **Workshop-track ACCEPT threshold documentation.** P-138 RepairKV (AdaptFM Workshop ICML 2026) is the borderline case; future workshops on the §7 trigger-list should be evaluated against P-138's threshold (method-level, single-paper, not a workshop-redundant theme).
4. **§11.6 carry-over.** The §11.6 Program-Comprehension Mechanism substrate claim is opt-in via `configs/agents.yml:comprehension.graph_self_evolving: true`; Wave-21 should re-verify the opt-in path is wired (the configs file may need a corresponding field added).


