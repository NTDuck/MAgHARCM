---
title: HuggingGPT / Jarvis — Task Planning with Foundation-Model Controllers
bibkey: p39_hugginggpt_2023
tags: [paper, task-planning, foundation-model-controller, agent-recruitment, [[PRIM-29]], hop-1]
---

# [[1.0.0 P-39]] HuggingGPT / Jarvis — Task Planning with Foundation-Model Controllers

**Authors**: Shen et al. — Yongliang Shen, Kaitao Song, Xu Tan, Dongsheng Li, Weiming Lu, Yu Qiao (Microsoft Research Asia + University of Copenhagen); contemporaneous with the Microsoft Jarvis release team (task-planning portion).
**Year**: 2023
**Venue**: NeurIPS 2023 (HuggingGPT); Microsoft Technical Report 2305.19268 (Jarvis); ICLR 2024 follow-up (HuggingGPT v2 / TaskBench)
**eprint / DOI**: arXiv:2303.17580 (HuggingGPT); arXiv:2305.19268 (Jarvis)
**Cited by**: [[primitives/INDEX]] entry [[PRIM-29]] (Recruitment-Adaptive Planning) — provides the task-graph model the Recruiter uses to plan re-configuration.

## Summary

[[Shen-2023-HuggingGPT]] and [[Microsoft-2023-Jarvis]] propose the same architectural pattern, independently. A large language model acts as a controller; the controller receives a natural-language task, decomposes it into a directed acyclic graph of subtasks, and recruits specialised expert models from a registry (HuggingFace model hub in HuggingGPT; Microsoft internal model zoo in Jarvis). Each subtask is dispatched to the chosen expert, the expert's output is collected, and the controller resolves dependencies before the next subtask begins. The recruitment step is conditioned on a structured description of each available expert (capability tag, expected input/output schema, latency budget); the controller selects the expert whose schema best matches the subtask. The TaskBench evaluation (ICLR 2024) shows that explicit task-graph planning with structured recruitment outperforms free-form prompting by 18–31 percentage points on multi-domain tasks. The pattern is sometimes called "task planning with foundation-model controllers" or "controller–expert orchestration".

## Relevance to MAgHARCM

[[PRIM-29]] (Recruitment-Adaptive Planning) is MAgHARCM's specialised version of HuggingGPT's task-graph planning. Where HuggingGPT recruits static domain experts (image captioning, ASR, OCR), the Recruiter re-configures MAgHARCM's eight-agent roster between validator iterations: it may add a second verdict agent when [[PRIM-7]] reports a semantic-equivalence disagreement, swap the [[PRIM-23]] chunked translator for a trait-mapping specialist when the validator reports a borrow-checker failure, or shrink the team to the core translate/validate pair when the validator reports progress. The structured-expert-description metadata lives in `internal/compiletime/compiletime.go` (the centralised registry) and each agent advertises its capability vector at startup, exactly mirroring HuggingGPT's capability-tag schema. The task graph the Recruiter emits is recorded in the [[PRIM-28]] checkpoint for reproducibility.

## Hop-1 References

- [[Chen-2023-AgentVerse]] — dynamic recruitment in multi-agent LLM systems; both HuggingGPT and AgentVerse address the same "who handles this subtask" question, with AgentVerse focused on dialogue agents and HuggingGPT on model experts.
- [[Shen-2023-HuggingGPTv2]] — TaskBench evaluation suite; second-generation HuggingGPT that exposes the task-planning capability as a separate benchmark.
- [[Lu-2024-Frozen-LM-Planner]] — frozen-LM planners for embodied task graphs; corroborates the controller–expert split when the controller weights are frozen.
- [[Qin-2023-LLM-Planner]] — plan-then-execute decomposition for embodied agents; alternative formulation that motivates the planner–recruiter split in [[PRIM-1]] / [[PRIM-29]].
- [[Wu-2023-AutoGen]] — conversation-as-computation framework whose group-chat manager is a thin version of the recruitment pattern.
- [[Khattab-2023-DSPy]] — programmatic prompt composition with declared module signatures; underlies how HuggingGPT describes expert I/O schemas in a typed way.

## Hop-2 Anchors (software-archaeology lean)

- [[Muller-2000]] — Müller's five migration strategies provide the candidate repertoire the Recruiter draws from; HuggingGPT's task graph is the structural form the Recruiter uses to plan a multi-strategy migration.
- [[Baldwin-Clark-2000]] — the controller / expert split is an L1 / L2 design-rule decomposition: the controller decides which expert (L2) is substitutable for a subtask; the registry of experts is the L3 set.
- [[Rajlich-1997]] — concept-locator and concept-assignment are the offline analogues of structured expert descriptions; the Recruiter's capability vector is a Rajlich-style "what does this module know" description.
- [[Parnas-1972]] — information-hiding module boundaries; the structured expert signature (capability tag + I/O schema) is the published interface of each recruited agent.
- [[Foltz-2023]] — DR.JONES Phase 2 (search for relevant concepts) is exactly what the Recruiter does when scanning the capability registry: it pattern-matches the subtask description against each expert's capability vector using the same cognitive heuristic.

## Backlinks

- [[Methodology]] §PRIM-29 — task-graph planning + capability-tag recruitment is the theoretical backing for the Recruiter primitive.
- [[primitives/INDEX]] — [[1.0.0 PRIM-29]] status entry points at P-39 as the controller-expert hop-1, alongside P-26 (AgentVerse) for dynamic-roster hop-1.
- [[Software-Archaeology-Lineage]] §4 — Recruiter sits at the end of the eight-agent cycle, using HuggingGPT's task-graph model to plan the next iteration's roster.
