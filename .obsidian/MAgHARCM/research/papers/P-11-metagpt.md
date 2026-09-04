---
title: MetaGPT — Meta Programming for A Multi-Agent Collaborative Framework
bibkey: p11_metagpt
tags: [paper, multi-agent, sop, role-artifact, [[PRIM-24]], hop-1]
---

# MetaGPT — Meta Programming for A Multi-Agent Collaborative Framework

**Authors**: Sirui Hong, Xiawu Zheng, Jonathan Chen, Yuheng Cheng, Jinlin Wang, Ceyao Zhang, Zili Wang, Steven Ka Shing Yau, Zijuan Lin, Liyang Zhou, Chenyu Ran, Lingfeng Xiao, Chenglin Wu, Jürgen Schmidhuber (DeepWisdom, Xiamen University, CUHK Shenzhen, Nanjing University, University of Pennsylvania, UC Berkeley, KAIST)
**Year**: 2023 (ICLR 2024 Oral)
**Venue**: International Conference on Learning Representations (ICLR) 2024 — Oral track
**eprint / DOI**: arXiv:2308.00352
**Cited by**: [[primitives/INDEX]] entry for the SOP-anchored role-artifact scheme underlying [[PRIM-24]] (SOP-Anchored Role-Artifact Schema).

## Summary

[[Hong-2023-MetaGPT]] frames the engineering of complex software as materialising human Standardised Operating Procedures ("Code = SOP(Team)"). A small team of specialised LLM agents — Product Manager, Architect, Project Manager, Engineer, QA Engineer — are wired together by a prompt-encoded assembly line that flows artifacts (PRD, design doc, API spec, file skeletons, source code) from role to role. Each role produces a strictly-typed intermediate output that the next role consumes; intermediate verification ("review and verify") between adjacent roles reduces the cascading hallucinations observed in pure chat-chain systems. MetaGPT reports substantially higher complexity handling than prior chat-based multi-agent systems on code-generation benchmarks and is the seminal reference for SOP-anchored multi-agent coordination. The framework is open-source (`foundationagents/metagpt`); the design has since been generalised by AFlow (ICLR 2025) and MGX.

## Relevance to MAgHARCM

MetaGPT is the conceptual anchor for [[PRIM-24]] (SOP-Anchored Role-Artifact Schema). MAgHARCM's `internal/artifacts/versioning.go` (planned this sprint) embodies the MetaGPT discipline that every agent emits a versioned, schema-validated intermediate — required fields, cardinality constraints, parse-or-retry gate before the next agent consumes it. The "Code = SOP(Team)" insight also motivates why [[PRIM-6]] (Multi-Stage Build/Test Feedback Repair) sequences validators by artifact rather than by agent: the artifact, not the prompt, is the contract between stages. MetaGPT's role-specialisation roster maps roughly onto MAgHARCM's analyser / planner / translator / navigator / validator decomposition; the difference is that MetaGPT's stages are conversational (chat-chain) while MAgHARCM's are tool-mediated (LSP + tree-sitter + chunked translator).

## Hop-1 References

- [[Qian-2023-ChatDev]] — sibling chat-chain framework; complements MetaGPT's SOP rigidity with a phase-decomposed waterfall (Design → Code → Test → Document) and is the immediate comparator paper.
- [[Wu-2023-AutoGen]] — multi-agent conversation framework; the third leg of the 2023 multi-agent triad (MetaGPT / ChatDev / AutoGen) that establishes the design space.
- [[Schmidhuber-2020-Planning]] — historical thread on hierarchical planning and learning from instructions; Schmidhuber's authorship on MetaGPT links the SOP idea back to his earlier work on learning to think.
- [[Hong-2024-AFlow]] — automated search over MetaGPT-style workflows; meta-evolution of the SOP idea.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997]] — concept assignment is the offline analogue of MetaGPT's product-manager role: the PM's PRD is a reification of which concepts the codebase is supposed to encode.
- [[Baldwin-Clark-2000]] — MetaGPT's role boundaries are an operationalisation of design rules: each role owns a layer (L1 interfaces / L2 subsystems / L3 leaves) and may freely substitute within it.
- [[Kazman-Cai-2024]] — architectural recovery underwrites the "Architect" role's input; without a recovered module / connector view, the architect is hallucinating the structure it claims to design.
