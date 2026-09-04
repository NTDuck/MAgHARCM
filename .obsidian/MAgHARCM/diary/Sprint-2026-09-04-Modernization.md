---
title: Sprint Modernization & Architectural Convergence
backlink: [[2.0.0 Sprint-Modernization]]
tags: [sprint, modernization, software-archaeology, refactor, [[2.0.0 MAgHARCM]]]
---

# [[2.0.0 Sprint Modernization & Architectural Convergence]]

## Context & Handoff

This sprint resumes from [[1.0.0 Sprint-Recon-2026-09-04]].
The prior iteration established the basic multi-agent pipeline and cataloged 31 primitives.
However, discrepancies remained:
1. Primitives parity was split: several primitives ([[PRIM-15]], [[PRIM-16]], [[PRIM-18]], [[PRIM-19]], [[PRIM-20]], [[PRIM-22]]) were marked "referenced" rather than implemented, creating multiple sources of truth.
2. The Eino graph only wired 4 nodes while the system contains more specialized agents.
3. Migration strategy was static instead of incremental try-and-fail.
4. Struct definitions were scattered across `internal/artifacts` rather than co-located with their owning agents (Locality of Behaviour).
5. Compile-time configs and initializations retained fallback patterns instead of strict `Must` validation.

## Sprint Goals

1. **Research Deepening**: Recursively research cited literature and 2-hop foundational works (Michael Feathers, Lehman's laws, Parnas modularity, Chikofsky & Cross reverse engineering, Rajlich & Müller concept assignment, Baldwin & Clark design rules, Kazman DRSpaces, Foltz comprehension).
2. **Vault Synchronization**: Normalize all version markers to `[[x.y.z ...]]` format. Eliminate unanchored parentheses and enforce single source of truth across all 31 primitives.
3. **Full Primitive Parity**: Implement all remaining primitives directly in Go under `internal/agents/`:
   - `[[1.0.0 PRIM-15]]` Evidence-First Adaptation Pattern
   - `[[1.0.0 PRIM-16]]` Spec-Driven Development Lifecycle
   - `[[1.0.0 PRIM-18]]` Jaccard-Coupling Architecture Recovery
   - `[[1.0.0 PRIM-19]]` Design Rule Hierarchy Partitioning
   - `[[1.0.0 PRIM-20]]` Concept Assignment and Redocumentation
   - `[[1.0.0 PRIM-22]]` Four Phases of Comprehension
   - Complete `[[1.0.0 PRIM-24]]` (Role-gated emission enforcement) and `[[1.0.0 PRIM-31]]` (Iterative retrieval refinement loop).
4. **Architectural Refactoring**:
   - Centralize compile-time configuration and constants in `internal/compiletime` using strict `Must` patterns without fallbacks.
   - Extract agent artifacts into agent files to enforce Locality of Behaviour and cohesion.
   - Refactor `SelectMigrationStrategy` to an incremental try-and-fail loop with clean enum types.
   - Decouple agent units: eliminate hidden direct references between agents; communicate via explicit contracts.
   - Expand the Eino graph to represent all active agents as dedicated nodes.
   - Default AST provider to `abcoder-mcp` across all configurations.
   - Standardize project compilation to binary Pass/Fail (no partial compilation rate).
   - Enforce ASD-STE100 across all user-facing copy and logs.
   - Re-architect interactive TUI with idiomatic Charm stack (Bubble Tea, Bubbles viewport/spinner/progress/table, Lip Gloss, Glamour).
5. **Verification & Paper Alignment**:
   - Full test suite verification with zero regressions.
   - Synchronize paper text, tables, and architectural figures with updated methodology.

## Closure

The sprint closed with all listed goals delivered or explicitly deferred with reason:

1. **Research Deepening** — research lineage expanded to 2-hop software-archaeology anchors; the deferred sub-goal "complete `magh-pipeline.workflow.json` geometry" remains open and is tracked under §Open Items below.
2. **Vault Synchronization** — all version markers normalized to `[[x.y.z ...]]`; PRIM-24 mapped to `internal/compiletime/compiletime.go` and `internal/agents/state.go`; single source of truth restored.
3. **Full Primitive Parity** — every referenced primitive was promoted to a real implementation; the 31-entry catalog is now one-to-one with `internal/agents/*.go`.
4. **Architectural Refactoring** — compile-time configs centralised under strict `Must` patterns; artifact structs moved next to their owning agents (Locality of Behaviour); `SelectMigrationStrategy` replaced by `Registry.TryInOrder`; agent units decoupled; Eino graph widened to 8 nodes; default AST provider is `abcoder-mcp`; per-project compilation is binary Pass/Fail; ASD-STE100 phrasing applied across user-facing copy; TUI is idiomatic Charm (Bubble Tea + Bubbles + Lip Gloss + Glamour).
5. **Verification & Paper Alignment** — `go test -count=1 ./...` green; `pdflatex main.tex` clean (17 pages, zero errors); `fig1_workflow.png` regenerated to reflect the 8-agent pipeline; abstract, intro, method, related work, eval, and conclusion narrate the eight-agent set; `sec_related.tex` preserves ReCodeAgent's published 4-agent description as the contrast baseline.

## Open Items

- `docs/.paper/figures/magh-pipeline.workflow.json` (workflow schema-v2 source) still carries unresolved edge-crossing diagnostics; the paper currently embeds the architecture PNG (`fig1_workflow.png`) so submission is unaffected, but a future pass must close the workflow geometry to keep the source-of-truth in sync.
