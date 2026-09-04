---
title: ChatDev — Communicative Agents for Software Development
bibkey: p12_chatdev
tags: [paper, multi-agent, chat-chain, communicative-dehallucination, [[PRIM-25]], hop-1]
---

# ChatDev — Communicative Agents for Software Development

**Authors**: Chen Qian, Xin Cong, Cheng Yang, Weize Chen, Yusheng Su, Juyuan Xu, Zhiyuan Liu, Maosong Sun (Tsinghua University)
**Year**: 2023 (ACL 2024 main)
**Venue**: Annual Meeting of the Association for Computational Linguistics (ACL) 2024
**eprint / DOI**: arXiv:2307.07924
**Cited by**: [[primitives/INDEX]] entry for the communicative-de-hallucination pattern underlying [[PRIM-25]] (Communicative-De-hallucination Role-Flip Gate).

## Summary

[[Qian-2023-ChatDev]] models software development as a waterfall of four stages — Design → Code → Test → Document — and orchestrates the work through a "chatting-chain" of agent dyads (CEO/CTO, Programmer/Reviewer, Programmer/Test Engineer, Programmer/Documenter). Within each phase the chain is broken into atomic subtasks; an "instructor" agent proposes an action and an "assistant" agent critiques or extends it; both sides share a memory stream of prior dialogue so context is preserved across phases. A unifying-bridge insight is that natural-language planning and code artifacts are interleaved within the same conversation, mitigating the inconsistencies that arise when separate systems own each phase. Evaluation reports substantially higher completeness and consistency than prior chat-based multi-agent baselines on small-program generation tasks. The framework is open-source (`openbmb/ChatDev`).

## Relevance to MAgHARCM

ChatDev is the direct citation for [[PRIM-25]] (Communicative-De-hallucination Role-Flip Gate). The pattern — a Translator emits code, a second agent with a role-flipped prompt ("you must find at least one bug") inspects the output — is a thinner, more local version of ChatDev's programmer/reviewer dyad: the reviewer is in-process and bound to the immediate artifact rather than to a long waterfall. ChatDev's memory-stream concept also motivates MAgHARCM's persistent state summary inside the chunked translator ([[PRIM-23]]), which records emitted modules, public symbols, and type aliases between calls. The waterfall-decomposition influence is visible in the multi-stage build / test feedback repair ([[PRIM-6]]); the test-phase dyad is the ancestor of the test-synthesis step in [[PRIM-5]].

## Hop-1 References

- [[Hong-2023-MetaGPT]] — the other 2023 multi-agent framework; ChatDev contrasts with MetaGPT's rigid SOPs by being a chat-chain, and the two are routinely co-cited as the design-space pair.
- [[Wu-2023-AutoGen]] — conversation-as-computation framework; ChatDev's dyads can be implemented as AutoGen agents and many follow-up papers do exactly that.
- [[Park-2023-Generative-Agents]] — foundational work on memory-streamed agent simulation in a sandbox town; the memory-stream mechanism ChatDev uses is borrowed from this line.
- [[Du-2023-Multi-Agent- Debate]] — multi-agent debate with two opposing critics; the role-flipped critic gate in [[PRIM-25]] is a translation-specialised variant of this debate pattern.

## Hop-2 Anchors (software-archaeology lean)

- [[Foltz-2023]] — DR.JONES cognitive model explains why the chat-chain with explicit phases reduces errors: linear, recency-biased traversal of an explicit catalog beats a single-shot inference that must hold many constraints at once.
- [[Baldwin-Clark-2000]] — ChatDev's CEO/CTO dyad is an L1 design-rule decision; the Programmer/Reviewer dyad is an L2 / L3 substitutability check; the four phases correspond to the three layer boundaries.
- [[Muller-2022]] — "COTS" and "Integrate-in-Place" from Mueller's migration-strategy taxonomy are first-class choices in ChatDev's Documenter phase: depending on the strategy, the output is either a generated library or a migration patch.
