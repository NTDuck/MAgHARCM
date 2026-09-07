---
title: "Reeper — Evidence-First Adaptation Pattern for External Repository Integration"
backlink: "[[1.0.0 P-16]]"
bibkey: p16_reeper
aliases:
  - "1.0.0 P-16"
  - "P-16"
  - "P-16-Reeper"
  - "P-16-Reeper"
  - "Reeper"
tags: [paper, evidence-first, adaptation, external-integration, [[PRIM-15]], hop-2]
---

# [[1.0.0 P-16]] Reeper — Evidence-First Adaptation Pattern for External Repository Integration

**Authors**: leadgenjay (single-maintainer / community project; no formal paper venue)
**Year**: 2025–2026 (active development on `main` as of late 2025)
**Venue**: GitHub repository + companion skill-card (leadgenjay.com/skills/reeper); **no peer-reviewed paper** — sourced from project README, `docs/architecture.md`, and an OpenHands-style blog write-up
**eprint / DOI**: none — https://github.com/leadgenjay/Reeper
**Cited by**: [[primitives/Primitives-Index]] entry [[PRIM-15]] (Evidence-First Adaptation Pattern); references target-preserving spec, contract-before-code, provenance-aware diff, untrusted-source model

> **MISSING BIBKEY.** `docs/.paper/refs.bib` does not contain a `reeper` or `p16_reeper` entry; the appendix cites this source as `Reeper [S12 \S2]`, an S-prefix marker that has no corresponding bib record. This note uses the synthetic bibkey `p16_reeper` to match the `pNN_<author>` convention.

## Summary

[[Reeper-2025]] targets the failure mode of LLM coding agents that "copy code well but notice badly": a competent agent can transplant a library's modules into a target codebase while silently overwriting the host's authentication, database, billing, deployment, design-system, or operational conventions. Reeper replaces the "copy-then-hope" workflow with an evidence-first protocol that emits four artefacts before any code is moved: (1) a **target-preserving spec** that captures which host conventions must be honoured by every adaptation, (2) a **contract-before-code** map stating the API guarantees the adaptation must keep, (3) a **provenance-aware diff** that records, for every line that moves, where it came from and which host convention it touches, and (4) an **untrusted-source model** that flags every assumption the source library makes about its environment (Python version, glibc, network topology, secret-management, time source). Each artefact is reviewed by the operator, conflict interviews are scheduled when the untrusted-source model contradicts the host spec, and only after the four artefacts are accepted does the integration step begin. The framework is shipped as a Claude Code skill plus a standalone CLI and is documented under `docs/architecture.md` with a five-stage state machine (intake → spec → conflict → adapt → verify).

## Relevance to MAgHARCM

[[PRIM-15]] cites Reeper verbatim for the four-artefact gate that must precede any external-code adaptation. MAgHARCM's commons-validator and Apache HTTP client samples both pull in third-party C/Java helpers that the target Rust crate ecosystem does not 1:1 replicate; without Reeper's evidence-first gate, the `internal/agents/manifest_rewriter.go` ([[PRIM-30]]) cannot tell whether a missing crate is a clean drop-in or a silent convention-violator. The contract-before-code map in particular should be wired into `internal/agents/feature_mapping.go` ([[PRIM-10]]) so that each idiom mapping carries an explicit "preserves target convention?" boolean — this catches the [[Pan-2024-LostInTranslation]] failure class at mapping time rather than at validation time.

## Hop-1 References

- [[Khajeh-Hosseini-2019-AutomatedTransformation]] — same modernisation lineage; Reeper's "target-preserving spec" formalises the legacy-system-preservation goal Khajeh-Hosseini et al. discuss.
- [[Müller-2000-IWPC]] — Reeper's "untrusted-source model" is the LIS-style decomposability assessment applied at the *external* library granularity rather than the legacy-internal granularity.
- [[Feathers-2004-WorkingEffectivelyWithLegacyCode]] — Feathers' "characterisation tests first, change behaviour second" is the conceptual ancestor of Reeper's contract-before-code map.
- [[Rajlich-1997-ICSE]] — concept-locator analysis supplies the "where did this assumption come from?" traceability that Reeper's provenance-aware diff assumes.
- [[Kazman-Cai-2015-ICSE-SEIP]] — design-rule-space decomposition explains *why* an unvetted external import can propagate debt: the imported rule violates the host's L1 interface layer.

## Hop-2 Anchors (software-archaeology lean)

- [[Baldwin-Clark-2000-DesignRules]] — Reeper's evidence-first protocol is a *modular operator* in Baldwin & Clark's sense: each of the four artefacts (spec, contract, diff, model) is an independent module that can be replaced without breaking the others, and the protocol itself adds an option value (reversible adaptation) that closed-loop "just port it" workflows destroy.
- [[Kazman-Cai-2024-ArchitecturalRecovery]] — the untrusted-source model is an architectural-recovery artefact applied to the *imported* codebase; without prior call-graph + dependency recovery on the host, the model cannot flag the cross-boundary propagation paths.
- [[Müller-2000-LegacyIntegration]] — Reeper's conflict-interview stage operationalises Müller's "Integrate-in-Place" strategy at the artefact level rather than the project level, giving operators a per-artefact accept/reject gate.
