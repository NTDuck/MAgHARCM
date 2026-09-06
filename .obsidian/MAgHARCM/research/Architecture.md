---
title: MAgHARCM Architecture
date: 2026-09-28
backlink: [[2.0.0 Architecture]]
last_updated: 2026-09-28
tags: [architecture, package-graph, [[2.0.0 MAgHARCM]], "[[1.0.0 P-122]]", "[[1.0.0 P-123]]", "[[1.0.0 P-124]]", wave-16, wave-17]
---
# [[2.0.0 MAgHARCM Architecture]]

## 1. Package Decomposition & Dependency DAG

The package hierarchy is structured with deep cohesion, explicit boundaries, and clear communication contracts.
Arrows indicate import dependencies (`A ──► B` means `A` imports `B`).

```
cmd/MAgHARCM, cmd/MAgHARCM-tui
       │
       ▼
internal/runner
       │
       ├───────────────────────────────┐
       ▼                               ▼
internal/graph                  internal/config (YAML runtime loader, Must pattern)
       │                               │
       ▼                               │
internal/agents                        │
  ├── State & Artifacts (Cohesive)     │
  └── 8 Specialized Agent Units        │
       │                               │
       ├───────────────────────────────┘
       ▼
internal/compiletime (centralized compile-time configs, enums, sentinels, Must helpers)
       │
       ▼
internal/tools, internal/languages, internal/logger
```

---

## 2. Core Packages

### 2.1. `internal/compiletime`
Centralizes compile-time configurations, enums, sentinels, and initialization helpers:
- Compile-time sentinels: `ErrorUnknown = "unknown"`
- Migration strategy enums: `StrategyBigBang`, `StrategyIncremental`, etc.
- Toolchain names: `ToolchainCargo = "cargo"`
- Language names: `LangRust = "rust"`, etc.
- LSP provider names: `LSPProviderABCoder = "abcoder-mcp"`, `LSPProviderNative = "native"`
- Binary compilation status enum: `CompilationStatusPass`, `CompilationStatusFail`
- Must pattern helpers: `MustLoadYAML`, `MustParseYAML`, `Must`

### 2.2. `internal/agents` (Locality of Behaviour)
Each agent declares its own role-artifact structures in its own file:
- `analyzer.go`: defines `AnalyzerAgent`, `AnalyzerOutput`, `SourceProjectResearch`, etc.
- `archaeology.go`: defines `ArchaeologistAgent`, `ArchaeologyReport`, etc.
- `planning.go`: defines `PlanningAgent`, `PlanningOutput`, etc.
- `translator.go`: defines `TranslatorAgent`, `TranslatedProject`, etc.
- `roleflip.go`: defines `RoleFlipAgent`, `ReviewerReport`, etc.
- `validator.go`: defines `ValidatorAgent`, `ValidationReport`, etc.
- `verdict_panel.go`: defines `VerdictPanelAgent`, `VerdictReport`, etc.
- `recruit.go`: defines `RecruiterAgent`, `RecruitmentPlan`, etc.
- `state.go`: defines pipeline `State`, coordinating the typed artifacts across the lifecycle.

### 2.3. `internal/graph`
Wires the full multi-agent pipeline into an executable Eino DAG with dedicated nodes:
1. `archaeologist`
2. `analyzer`
3. `planner`
4. `translator`
5. `reviewer`
6. `validator`
7. `verdict_panel`
8. `recruiter`
9. checkpoint anchors

---

## 3. Communication Contracts & Non-Interference

No agent unit assumes knowledge of another agent's internal implementation:
- Agents do not retain pointers to peer agent instances.
- Communication occurs exclusively through `*State`.
- Error handling uses Go standard errors without hidden defaults.

---

## 4. Agent Topology (8-agent graph as of 2026-09-26)

The multi-agent pipeline is wired as an Eino DAG with **10 lambda nodes** (8 specialised agents + 2 checkpoint barriers) and a single cyclic repair edge. All eight agents are registered in `internal/graph/graph.go` via `g.AddLambdaNode(...)` and connected in the following order; every agent's `*LambdaNode` declaration includes the `PRIM-NN` anchor comment documenting which lineage primitives the node implements.

> **Shorthand note**: The phrase `'8-agent graph'` used as a tag elsewhere in vault audit lines (e.g. `primitives/INDEX.md`, `METHODOLOGY.md` §11, historical sprint handoffs) is shorthand for this 10-node topology. The `reviewer` lambda registered at `graph.go:92` is constructed via `agents.NewRoleFlipGate` (`internal/agents/roleflip.go:44`); "reviewer" is the graph node id, "RoleFlipGate" is the constructor name.

| λ-# | Node id | Role | PRIM-NN anchors | `graph.go` line | Cycle role |
| :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | `archaeologist` | Pre-planning software archaeology (commit churn, Jaccard coupling, design-rule layering, concept assignment, DR. JONES cognitive navigation) | `[[1.0.0 PRIM-14]]`, `[[1.0.0 PRIM-18]]`, `[[1.0.0 PRIM-19]]`, `[[1.0.0 PRIM-20]]`, `[[1.0.0 PRIM-22]]` | 45 | Forward head: `START → archaeologist` |
| 2 | `analyzer` | Source-project analysis producing `*AnalyzerOutput` on `*compiletime.State` | (companion to archaeologist — produces the typed artifact the planner consumes) | 63 | Forward |
| 3 | `planning` | Reverse-topological ordering, back-edge-conditioned scheduling, target skeleton-first generation | `[[1.0.0 PRIM-1]]`, `[[1.0.0 PRIM-2]]`, `[[1.0.0 PRIM-3]]` | 71 | Forward |
| 4 | `translator` | Chunked translation + iterative retrieval refinement | `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-31]]` | 79 | Forward + **repair-cycle head** (re-entrant via `recruiter → translator`) |
| 5 | `save_translator_ckpt` | Checkpoint barrier (`checkpointLambda` passthrough) — durable snapshot before role-flip review | `[[1.0.0 PRIM-28]]` (Conversable Checkpoints) | 87 | Forward barrier |
| 6 | `reviewer` | Role-flip de-hallucination gate (PRIM-25 Communicative De-hallucination Role-Flip Gate) | `[[1.0.0 PRIM-25]]` | 92 | Forward |
| 7 | `validator` | Test co-translation + multi-stage build/test repair + adversarial test-weakening guard + coverage plateau detection | `[[1.0.0 PRIM-5]]`, `[[1.0.0 PRIM-6]]`, `[[1.0.0 PRIM-13]]`, `[[1.0.0 PRIM-27]]` | 116 | Forward |
| 8 | `save_validator_ckpt` | Checkpoint barrier — durable snapshot before verdict/recruiter branch | `[[1.0.0 PRIM-28]]` | 124 | Forward barrier; **branch source** |
| 9 | `verdict_panel` | Multi-agent consensus verdict | `[[1.0.0 PRIM-7]]` | 128 | Repair-cycle body (reached only when validation incomplete) |
| 10 | `recruiter` | Dynamic iteration adaptation / try-and-fail strategy switch | `[[1.0.0 PRIM-29]]` | 150 | Repair-cycle tail — `AddEdge("recruiter", "translator")` closes the loop |

**Forward edges** (`internal/graph/graph.go:175-198`): `START → archaeologist → analyzer → planning → translator → save_translator_ckpt → reviewer → validator → save_validator_ckpt`.

**Repair cycle** (`internal/graph/graph.go:200-227`): the `save_validator_ckpt` node carries a `compose.GraphBranch` that routes to `compose.END` on completion / all-success / iteration-ceiling, otherwise to `verdict_panel`. From `verdict_panel` the edges `verdict_panel → recruiter → translator` close the cycle back to node 4. The compile-time safety ceiling is `compose.WithMaxRunSteps(compiletime.MaxGraphRunSteps)` (line 232) so the cycle cannot run unbounded.

**Compliance anchor**: ADR-C-006 binds the 8-agent graph (`internal/graph/graph.go`). The 2026-09-26 sprint re-verified that all 10 `AddLambdaNode` calls are present and that the agent producer files in `internal/agents/` match the constructor calls in `NewMAgHARCMGraph`.

---

## 5. Package Graph

Import topology enforced by `go build` + the ADR-C-001..015 convention set. Arrows indicate import direction (`A ──► B` means `A` imports `B`); the mandatory edges below are the ones a refactor MUST NOT remove without amending the relevant ADR.

```
cmd/MAgHARCM, cmd/MAgHARCM-tui
        │
        ▼
internal/runner ──► internal/graph ──► internal/agents
        │                  │                  │
        │                  │                  │ (mandatory edge:
        │                  │                  │   agents → compiletime
        │                  │                  │   is the only legal
        │                  │                  │   direction — reverse
        │                  │                  │   edge creates a Go
        │                  │                  │   import cycle)
        │                  ▼                  ▼
        │           internal/llm       internal/compiletime
        │                                     │
        ▼                                     │
internal/config                                │
                                              ▼
                       internal/tools, internal/languages, internal/logger
```

**Mandatory edges**:

- **`internal/agents ──► internal/compiletime`** — every producer agent file imports `compiletime` for `*compiletime.State`, `compiletime.Must`, and `compiletime.CompilationStatus`. The reverse edge (`compiletime → agents`) is **forbidden**: the Sprint 2026-09-26 import-cycle blocker (see §6) was the third time a relocation attempt hit the same wall.
- **`internal/graph ──► internal/agents`** — `graph.go:9` (`"MAgHARCM/internal/agents"`) wires the eight `NewXxxAgent()` constructors into lambda nodes. The reverse edge (`agents → graph`) is forbidden; agent files declare their nodes, they don't wire them.
- **`internal/runner ──► internal/graph`** AND **`internal/runner ──► internal/agents`** — the runner constructs the graph and additionally talks to a subset of agent constructors for partial-state recovery (`graph.go:9` + the per-agent constructors invoked through `internal/runner/state_recovery.go`).

**Cycle-avoidance alias pattern**: when a type conceptually lives in `agents/` (Locality-of-Behaviour intent) but is also referenced by `compiletime.State`, the type lives in the producer agent file (full docstring + cross-references) and `compiletime` declares a Go **type alias** (`type X = agents.X`), not a redeclaration. This is the ADR-C-014-as-applied pattern documented in `METHODOLOGY.md` §6 and confirmed in §6 below.

---

## 6. Locality of Behaviour (ADR-C-014) — Durable Constraint

ADR-C-014 requires per-agent artifact structs (`AnalyzerOutput`, `PlanningOutput`, `TranslatedProject`, `ValidationReport`, `ArchaeologyReport`) to live alongside their producer agent file (`internal/agents/{analyzer,planning,translator,validator,archaeology}.go`). The constraint collides with the mandatory `agents → compiletime` edge (because `compiletime.State` references those structs), and that collision has now been verified three times.

**Timeline:**

- **Sprint 2026-09-07, commit `0cb5994`** — the original canonical relocation. Artifact structs moved from `internal/compiletime/state.go` to the producer agent files; `compiletime/state.go` re-exposes them via Go type aliases (`type AnalyzerOutput = agents.AnalyzerOutput`, etc.). This is the **settled solution** and has stood since 2026-09-07. Commit hash: `0cb5994 refactor(state): canonicalise State + 6 artifact structs in compiletime/state.go (ADR-C-014)`.
- **Sprint 2026-09-23, commit `89904f6`** — a banner rewrite in `internal/compiletime/state.go` that clarified the constraint in prose. Reverted/dead-end (no cycle fix found). Commit hash: `89904f6 refactor(compiletime): rewrite misleading ADR-C-014 banner`.
- **Sprint 2026-09-26, subagent D aborted** — a third attempt (leaf sub-package `internal/compiletime/artifacts/`) hit the same hard blocker: `ValidationReport.CompilationStatus()` returns `compiletime.CompilationStatus`, so the return type MUST move with the method receiver — but moving `CompilationStatus` out of `compiletime/compiletime.go` requires updating every caller in `internal/agents/`, `internal/runner/`, `internal/graph/`, which in turn requires a new import edge that recreates the cycle.

**Settled solution (ADR-C-014-as-applied, since 2026-09-07, commit `0cb5994`)**: artifact struct definitions live in `internal/agents/{analyzer,planning,translator,validator,archaeology}.go` with full docstrings + cross-references; `internal/compiletime/state.go` declares them via `type X = agents.X` (Go type aliases, not redeclarations). The Locality-of-Behaviour intent is met for the **reader** of the producer file (they see the struct right next to the producer) without breaking the `agents → compiletime` edge.

**Durable constraint**: future sprints that touch artifact structs MUST consult `internal/compiletime/state.go`'s preamble before proposing any re-homing. The `agents → compiletime` edge is the only legal direction; the reverse edge is a Go import cycle. Producer-file alias pattern is the only known durable resolution.

---

## 7. Vault Sync Audit — Sprint 2026-09-28 (cumulative)
Mirrors `.obsidian/MAgHARCM/primitives/INDEX.md` Sprint 2026-09-26 audit block.

- **Wave-15 status**: deferred on 2026-09-26 (3 candidates rejected per §7 trigger gate). Full rationale in `.obsidian/MAgHARCM/research/diary/wave-15-candidates.md`.
- **Wave-16 status**: FIRED 2026-09-27. 3 candidates triaged; 1 ACCEPT (P-122 ReasoningBank ICLR 2026 tentative) + 2 REJECT (SWE-Bench Pro ICML 2026 + CodeClash ICML 2026 — both fail Q2 mechanism-vs-benchmark gate). Full rationale in `.obsidian/MAgHARCM/research/diary/wave-16-candidates.md`. 1 new SLM-era anchor persisted; 1 new paper note added (`.obsidian/MAgHARCM/research/papers/P-122-reasoningbank-iclr-2026.md`).
- **Wave-17 status**: FIRED 2026-09-28. 5 candidates triaged; 2 ACCEPT (P-123 CodeChemist ICML 2026 + P-124 Syzygy ICLR 2025 VerifAI Workshop) + 3 REJECT (MemSearcher ACL 2026 Findings — off-list venue; Verified Tool Calls arXiv-only — no venue confirmation; LLM-IR / program-comprehension — benchmark, no mechanism). Full rationale in `.obsidian/MAgHARCM/research/diary/wave-17-candidates.md`. 2 new SLM-era anchors persisted; 2 new paper notes added.
- **ADR-C-014 status**: locality split applied by Subagent D. The producer-file alias pattern (`internal/agents/<file>.go` declares the struct, `internal/compiletime/state.go` declares `type X = agents.X`) remains the durable constraint after the third relocation attempt (commit `89904f6` + Sprint 2026-09-26 sub-agent D abort) hit the same hard cycle blocker.
- **ADR-C-005 status**: magic-string sweep applied by Subagent E. `internal/consts/consts.go` is the canonical home for hardcoded-by-necessity values; the Sprint 2026-09-26 sweep swept residual string-literal sentinels from agent files into `compiletime/`.
- **ADR-C-011 status**: Charm stack audit applied by Subagent G. `internal/tui/tui.go` (613 lines after edits) confirmed idiomatic across `bubbletea` / `bubbles/spinner` / `bubbles/table` / `bubbles/textinput` / `lipgloss` / `glamour`. The dead `viewport` import (constructed but never `.View()`-ed) was removed. Full report in `.obsidian/MAgHARCM/diary/sprint-2026-09-26-charm-audit.md`.
- **ADR-V-001 status**: `scripts/lint_vault.sh` shipped by this subagent. Asymmetric enforcement: `FORBIDDEN` (exit 1) for `internal/` `(P-NN)` outside backticks; `STRAY` (report-only) for `.obsidian/MAgHARCM/` `(PRIM-NN)` per Sprint 2026-09-17 prose-convention decision (prose parentheticals retained as standard academic-writing convention). Documented in `.obsidian/MAgHARCM/architecture/ADR-2026-09-26-Vault-Lint-Extension.md`.
- **P-106 status**: `P-106 BFCL` retire task deferred to Sprint 2026-09-27. Sprint 2026-09-26 subagent F was cancelled; the duplicate anchor (P-121 supersedes P-106 as the canonical tool-calling benchmark reference) will be retired in a follow-up sweep.

## 8. Wave-16 SLM-Era Architectural Implications (2026-09-27)

Wave-16 fired 2026-09-27 (`[[1.0.0 P-122]]` ReasoningBank accepted, tentative ICLR 2026). One architectural shift implied; opt-in via YAML config to preserve backward compatibility.

### 8.1 Persistent-Memory Substrate (P-122)

Three iterative primitives acquire a persistent-memory substrate when `configs/agents.yml:memory.distilled: true`:
- **`[[1.0.0 PRIM-31]]` Iterative Retrieval** — feedback-driven refinement is augmented with strategy-distilled `(Title, Description, Content)` triples persisted across iterations.
- **`[[1.0.0 PRIM-29]]` Recruiter Agent** — `RecruitmentPlan` output becomes a structured triple; MaTTS compute-memory loop applies.
- **`[[1.0.0 PRIM-21]]` Migration Strategy Selection** — strategy-by-profile rewards persist; blind try-and-fail becomes informed switching.

The persistent-memory substrate is a leaf package `internal/memorystore/` (forthcoming — sprint 2026-09-28+); the primitives consume it via a typed `MemoryTriple` interface. No graph re-wiring required.



### 8.2 Backward Compatibility

Both shifts are opt-in. Existing users with `configs/agents.yml` default keys continue to get blind try-and-fail + feedback-driven retrieval (the Sprint 2026-09-26 baseline). New YAML keys added; no breaking changes.

### 8.3 Parity Check (post-wave-16, 2026-09-27)

- **Parity**: 31/31/31 unchanged — 31 primitives listed in `.obsidian/MAgHARCM/primitives/INDEX.md`, 31 rows in `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`, 31 implementation files in `internal/agents/*.go` (35 files = 31 impl + 4 test files + 0 orphan).
- **Wave-16 anchors**: 1 added (P-122). Total vault paper notes: 122 (121 prior + P-122).
- **Cross-links**: propagated in `Software-Archaeology-Lineage.md` (sprint 2026-09-27 close).


## 9. Wave-17 SLM-Era Architectural Implications (2026-09-28)

Wave-17 fired 2026-09-28 (`[[1.0.0 P-123]]` CodeChemist ICML 2026 + `[[1.0.0 P-124]]` Syzygy ICLR 2025 VerifAI Workshop, both accepted). Two architectural shifts implied; opt-in via YAML config to preserve backward compatibility.

### 9.1 Cross-Lingual Functional Oracle (P-123 CodeChemist)

Three primitives acquire a cross-lingual I/O test oracle when `configs/agents.yml:oracle.cross_lingual: true`:
- **`[[1.0.0 PRIM-21]]` Migration Strategy Selection** — confidence-gated switching between in-language majority voting (cheap) and cross-lingual I/O test oracle (expensive but cross-language functional) replaces blind try-and-fail.
- **`[[1.0.0 PRIM-23]]` Chunked Translation** — multi-temperature hedged sampling with cross-language functional verification. SLM-amenable: demonstrated on Qwen-1.5B.
- **`[[1.0.0 PRIM-27]]` Coverage-Guided Plateau Detection** — functional-coverage plateau via I/O oracle (tests across source + target language) replaces frontier-model judges.

The oracle is a typed `FunctionalOracle` interface in a leaf package `internal/oracle/` (forthcoming — sprint 2026-09-30+); primitives consume it via typed inputs. No graph re-wiring required.

### 9.2 Dynamic-Analysis Property Mining (P-124 Syzygy)

Three primitives acquire a runtime-mined property enrichment when `configs/agents.yml:translation.dynamic_specs: true`:
- **`[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph** — runtime-mined properties (aliasing, bounds, nullability) injected as fourth representation.
- **`[[1.0.0 PRIM-22]]` Four Phases of Comprehension** — static TDG (`[[1.0.0 P-88]]` HiTyper) + dynamic property mining (`[[1.0.0 P-124]]` Syzygy) as complementary dimensions.
- **`[[1.0.0 PRIM-30]]` Source-to-Target Manifest Rewriter** — type/bounds/nullability-enriched manifests for safe-Rust generation.

The mining substrate is a leaf package `internal/specminer/` (forthcoming — sprint 2026-09-30+); uses Clang/LLVM instrumentation when the legacy codebase compiles. Static-only path remains as fallback for uncompilable code.

### 9.3 Backward Compatibility

Both shifts are opt-in. Existing users with `configs/agents.yml` default keys continue to get blind try-and-fail + feedback-driven retrieval (the Sprint 2026-09-26 baseline) + static-only manifests (the Sprint 2026-09-07 baseline). New YAML keys added; no breaking changes.

### 9.4 Parity Check (post-wave-17, 2026-09-28)

- **Parity**: 31/31/31 unchanged — 31 primitives listed in `.obsidian/MAgHARCM/primitives/INDEX.md`, 31 rows in `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`, 31 implementation files in `internal/agents/*.go` (32 files = 31 impl + 1 test file + 0 orphan).
- **Wave-17 anchors**: 2 added (P-123, P-124). Total vault paper notes: 124 (122 prior + P-123 + P-124).
- **Cross-links**: propagated in `Software-Archaeology-Lineage.md` (sprint 2026-09-28 close) and `primitives/INDEX.md`.
---

