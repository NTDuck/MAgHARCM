---
title: AgentVerse — Facilitating Multi-Agent Collaboration and Exploring Emergent Behaviors
backlink: "[[1.0.0 P-26]]"
bibkey: p26_agentverse
aliases:
  - "1.0.0 P-26"
  - "P-26"
  - "P-26-AgentVerse"
  - "P-26-AgentVerse"
  - "AgentVerse"
tags: [paper, multi-agent, recruitment, [[PRIM-29]], hop-1]
---

# [[1.0.0 P-26]] AgentVerse — Facilitating Multi-Agent Collaboration and Exploring Emergent Behaviors

**Authors**: Weize Chen, Yusheng Su, Jingwei Zuo, Cheng Yang, Chenfei Yuan, Chi-Min Chan, Heyang Yu, Yaxi Lu, Yi-Hsin Hung, Chen Qian, Yujia Qin, Xin Cong, Ruobing Xie, Zhiyuan Liu, Maosong Sun, Jie Zhou (Tsinghua University, BUPT, WeChat AI Tencent)
**Year**: 2023 (arXiv preprint), accepted ICLR 2024
**Venue**: International Conference on Learning Representations (ICLR) 2024
**eprint / DOI**: arXiv:2308.10848
**Cited by**: [[primitives/Primitives-Index]] entry [[PRIM-29]] (Recruitment-Adaptive Planning)

## Summary

[[Chen-2023-AgentVerse]] frames multi-agent LLM systems as a dynamic recruitment problem. Rather than freezing the agent roster at the start of a task, AgentVerse introduces a four-stage loop — Recruitment, Decision-Making, Action Execution, Evaluation — where, before each iteration, a separate "recruiter" picks the agent set (specialists, tools, ordering) most likely to make progress on the current state ([P26 §2]). Recruitment is conditioned on the previous iteration's evaluation: agents that failed to contribute are dropped, new specialists with complementary skills are added, and the surviving team is reshuffled. Empirically, the authors show that the dynamic roster consistently outperforms a fixed roster of equal size across reasoning, coding, and consulting benchmarks, and that several emergent capabilities (e.g., autonomous role specialisation, cross-agent verification) only appear when recruitment is in the loop.

## Relevance to MAgHARCM

AgentVerse is the direct citation for [[PRIM-29]] (Recruitment-Adaptive Planning). MAgHARCM's `internal/agents/recruit.go` (planned this sprint) implements the recruitment primitive as a state-conditioned re-configuration step between validator passes: after each iteration of the [[PRIM-6]] feedback loop, the recruiter consults the [[PRIM-26]] navigator's diagnostic output and may swap the [[PRIM-23]] chunked translator for a specialist (e.g., a trait-mapping agent when the validator complains about a type mismatch), introduce an extra [[PRIM-7]] verdict agent when the failure mode is semantic, or downsize the team when progress plateaus. The recruitment decision itself is recorded in the [[PRIM-28]] checkpoint so the run is reproducible. Without AgentVerse's framing, MAgHARCM's roster would be static and would inherit the same brittleness that fixed multi-agent pipelines show on heterogeneous repositories.

## Hop-1 References

- [[MetaGPT-2024]] — fixed SOP-anchored multi-agent framework; the static roster that AgentVerse explicitly contrasts against.
- [[ChatDev-2024]] — communicative multi-agent SE pipeline using a chat-chain; another fixed-roster baseline.
- [[AutoGen-2023]] — conversation-programmable multi-agent framework; AgentVerse cites AutoGen's evaluator pattern as inspiration for the post-iteration evaluation step.
- [[HyperAgent-2026]] — generalist SE agent with a dedicated navigator; the kind of specialist AgentVerse's recruiter might recruit mid-loop.
- [[Chen-2023-LLM-Agent-Survey]] — taxonomic context for recruitment-as-control in LLM agent systems.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997]] — concept-locator analysis is the offline analogue of recruitment: knowing what concepts live in which modules is the same dependency-of-recruitment information the recruiter uses.
- [[Baldwin-Clark-2000]] — design-rule decomposition predicts that recruitment should preserve stable inter-module contracts (L1) and only vary L3 substitutable specialists; this constrains MAgHARCM's recruiter from destabilising the skeleton.
- [[Foltz-2023]] — DR.JONES cognitive model justifies the linear, recency-biased recruitment prompt: the human-like comprehension pattern that explains why dynamic rosters beat fixed ones.
