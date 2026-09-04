---
title: "GitHub Spec Kit — Spec-Driven Development Lifecycle Toolkit"
bibkey: p17_spec_kit
tags: [paper, spec-driven-development, lifecycle, scaffolding, [[PRIM-16]], hop-2]
---

# [[1.0.0 P-17]] GitHub Spec Kit — Spec-Driven Development Lifecycle Toolkit

**Authors**: Den Delimarsky, Aaron Blakely, GitHub engineering + community contributors
**Year**: 2025 (initial open-source release; reached 1.0.0 maturity in 2025)
**Venue**: Open-source toolkit (github/spec-kit on GitHub); documented at github.github.com/spec-kit and announced on github.blog (Generative-AI category) and developer.microsoft.com
**eprint / DOI**: none — https://github.com/github/spec-kit
**Cited by**: [[primitives/INDEX]] entry [[PRIM-16]] (Spec-Driven Development Lifecycle); referenced sequence is constitution → specify → plan → tasks → implement → converge

> **MISSING BIBKEY.** `docs/.paper/refs.bib` does not contain a `spec-kit` or `p17_spec_kit` entry; the appendix cites this source as `spec-kit [S13]`, an S-prefix marker with no corresponding bib record. This note uses the synthetic bibkey `p17_spec_kit` to match the `pNN_<author>` convention. Treat as project/tooling documentation rather than peer-reviewed literature.

## Summary

[[Spec-Kit-2025]] is GitHub's open-source toolkit for *Spec-Driven Development* (SDD). Its thesis: "vibe coding" — letting an LLM agent draft code from a vague prompt — produces output that *looks right* but quietly violates architecture, security, or edge-case requirements because the agent never had a precise spec to grade itself against. Spec Kit installs a structured lifecycle (`constitution → specify → plan → tasks → implement → converge`) in which every transition produces a reviewable Markdown artefact (`spec.md`, `plan.md`, `tasks.md`) and every gate requires explicit human/agent approval. The Specify CLI bootstraps the directory structure and templates; the toolkit is model-agnostic and ships adapters for 30+ coding agents (GitHub Copilot, Claude Code, Gemini CLI, Cursor, Zed, etc.). Spec Kit reached 1.0.0 in 2025, signalling its maturity as a *harness* (not a replacement) for existing agents. Each phase artefact is version-controlled and lives next to the code it specifies, making the spec a "living" artefact rather than a discarded PRD.

## Relevance to MAgHARCM

[[PRIM-16]] cites spec-kit as the canonical reference for the constitution → specify → plan → tasks → implement → converge sequence. MAgHARCM's `internal/agents/` directory already follows a phased structure (planning.go → manifest_rewriter.go → chunked_translator.go → validator.go) but the transitions are not gated by reviewable artefacts: a planner's output is consumed by a translator without a persisted intermediate that a human can sign off on. Spec-kit's template discipline should be back-ported into `internal/artifacts/` so each phase's output is a typed, persisted Markdown/JSON envelope that the next phase must parse-or-retry against ([[PRIM-24]]). The constitution stage in particular needs a one-page document capturing MAgHARCM's non-negotiables (no `unsafe` Rust without justification; preserve Go error-wrapping semantics; never drop test cases during repair) that every subsequent phase can be diffed against.

## Hop-1 References

- [[Hong-2023-MetaGPT]] — SOP-anchored role-artifact schema ([[PRIM-24]]) is the agent-side cousin of spec-kit's phase artefact; MetaGPT encodes SOPs in code, spec-kit encodes them in templates, both reject "vibe coding".
- [[Qian-2023-ChatDev]] — chat-chain role-flipping ([[PRIM-25]]) is the agent-level quality gate that spec-kit's per-phase human/agent review replicates at the workflow level.
- [[Wu-2023-AutoGen]] — AutoGen's `reply_func` + interrupt/resume primitive ([[PRIM-28]]) supplies the runtime machinery spec-kit needs to *enforce* a phase gate (pause the agent until the artefact is signed off).
- [[Chen-2023-AgentVerse]] — recruitment-adaptive planning ([[PRIM-29]]) reuses spec-kit's phase artefacts as the input that the recruitment step reads when deciding which agent configuration to spin up next.
- [[Müller-2000-IWPC]] — spec-kit's "constitution" stage is a modern repackaging of Müller's migration-strategy selection gate; both force the human to commit to non-negotiables *before* the technical plan is drafted.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997-ICSE]] — spec-kit's `spec.md` template is a *redocumentation* artefact in Rajlich's sense; for a legacy modernisation target like commons-validator, the "what is this code for?" answer is the concept map spec-kit's specify stage demands.
- [[Kazman-Cai-2024-ArchitecturalRecovery]] — spec-kit's `plan.md` template should be populated from a Kazman-style design-rule-space decomposition so the plan respects the load-bearing L1/L2 boundaries.
- [[Baldwin-Clark-2000-DesignRules]] — spec-kit's per-phase gating is a *visible-in-the-architecture* search rule; each phase artefact is a stable interface whose replacement (new template, new gate rule) is a localised, reversible change — the option value Baldwin & Clark prescribe.
