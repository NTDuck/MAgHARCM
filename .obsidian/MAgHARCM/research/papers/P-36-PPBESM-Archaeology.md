---
title: Software Archaeology - Hunting for Lost Knowledge in Production Codebases
backlink: "[[1.0.0 P-36]]"
bibkey: p36_ppbesm
aliases:
  - "1.0.0 P-36"
  - "P-36"
  - "P-36-PPBESM-Archaeology"
  - "P-36-PPBESM-Archaeology"
  - "PPBESM-Archaeology"
tags: [blog, playbook, software-archaeology, forensics, [[1.0.0 PRIM-14]], [[1.0.0 PRIM-18]], hop-1]
---

# [[1.0.0 P-36]] Software Archaeology: Hunting for Lost Knowledge in Production Codebases

**Authors**: pp-besm (Petar P.)  
**Year**: 2023  
**Venue**: dev.to Engineering Playbooks & Production Archaeology  
**eprint / DOI**: dev.to/ppbesm/software-archaeology-hunting-for-lost-knowledge  
**Cited by**: `[[2.0.0 Software-Archaeology-Lineage]]`, `[[1.0.0 PRIM-14]]` (Software Archaeology Stage), `[[1.0.0 PRIM-18]]` (Jaccard-Coupling Architecture Recovery)

## Summary

The pp-besm playbook is an influential practitioner treatise on executing software archaeology in legacy production environments where original authors have departed and documentation is lost or misleading. It formalizes code excavation across **Five Geological Strata**:
1. **Commit Strata**: Mining git log churn to identify which files change most frequently and which author communities created them. High-churn files represent volatile load-bearing hubs.
2. **Temporal Coupling (Co-Change Strata)**: Uncovering files that commit together in the same changeset without having direct syntactic dependencies (measured via Jaccard similarity).
3. **Bug-Density Hotspots**: Bisection of defect remediation commits to map legacy risk zones.
4. **Forensic Naming & Type Strata**: Identifying dead configuration switches, Hungarian naming residue, historical typedef wrappers, and implicit nullability assumptions.
5. **Executable Time Capsules**: Capturing historical build flags, compiler environments, and execution scripts needed to deterministically reproduce the original runtime.

## Relevance to MAgHARCM

Direct blueprint for `[[1.0.0 PRIM-14]]` (Software Archaeology Stage) and `[[1.0.0 PRIM-18]]` (Jaccard-Coupling Architecture Recovery).
In MAgHARCM's `internal/agents/archaeology.go` and `internal/agents/jaccard_coupling.go`:
- `FindChurnHotspots` analyzes git churn strata.
- `JaccardCouplingAnalyzer` calculates temporal co-change similarity.
- `BuildTimeCapsule` preserves executable environmental reproduction steps.
- `ForensicNaming` extracts deprecated naming patterns before planning begins.

## Hop-1 References

- [[Feathers-2004-WELC]] — Connects archaeological strata to seam identification and test points.
- [[MSR4SA-2017]] — Scientific validation of pp-besm's co-change temporal coupling metrics.
- [[Chikofsky-Cross-1990]] — Conceptual categorization of pp-besm's strata under Design Recovery.

## Hop-2 Anchors (Software Archaeology Lean)

- [[AgentPatterns-Legacy-Code-Archaeology]] — Implements pp-besm's terminal commands as structured MCP tool calls for autonomous agent loops.
- [[Kazman-Cai-2017]] — Uses temporal coupling to confirm static DRSpaces architecture debt patterns.
