---
title: MAgHARCM Architecture
backlink: [[2.0.0 Architecture]]
tags: [architecture, package-graph, [[2.0.0 MAgHARCM]]]
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
- Preconditions and postconditions are strictly typed.
- Error handling uses Go standard errors without hidden defaults.
