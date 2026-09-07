---
title: "P-63 — Dong et al. 2025 — A Survey on Code Generation with LLM-based Agents"
backlink: "[[1.0.0 P-63]]"
aliases:
  - "1.0.0 P-63"
  - "P-63"
  - "P-63-Dong-Survey-LLM-Code-Agents-2025"
  - "P-63-Dong-Survey-LLM-Code-Agents-2025"
  - "Dong-Survey-LLM-Code-Agents-2025"
tags: [paper, survey, multi-agent, code-generation, workflow-taxonomy, [[1.0.0 PRIM-23]], [[1.0.0 PRIM-24]], [[1.0.0 PRIM-29]], [[1.0.0 PRIM-17]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-63 — Dong et al. — Survey on Code Generation with LLM-based Agents]]

## Citation

Dong, Y., Jiang, X., Qian, J., Wang, X., Zhang, K., Jin, X., & Li, G. (2025). *A Survey on Code Generation with LLM-based Agents*. arXiv:2508.00083 (July 2025). Peking University. URL: https://arxiv.org/abs/2508.00083.

## Summary

The Dong et al. 2025 survey is the methodologically-oriented sibling to [[1.0.0 P-50]] (Jiang et al. 2024, which surveys single-model code generation). Dong et al. focus on *multi-agent* workflows for code generation and empirically taxonomize them into four families:

1. **Pipeline-based division of labor.** Each agent has a fixed role; data flows forward through the pipeline. Example: ChatDev (Qian et al. 2023), MetaGPT (Hong et al. 2023). MAgHARCM's forward path (archaeologist → analyzer → planning → translator → save_translator_ckpt → reviewer → validator → save_validator_ckpt) is in this family.
2. **Hierarchical.** A manager agent decomposes the task and dispatches subtasks to worker agents; results aggregate back. Example: HyperAgent.
3. **Self-negotiation circular.** Multiple agents iteratively critique each other's output until convergence. Example: Self-Refine (Saunders et al. 2022), Role-Flip (MAgHARCM's reviewer).
4. **Self-evolving.** Agents modify their own workflow based on observed outcomes. Example: AgentVerse (Chen et al. 2023, [[1.0.0 PRIM-29]] anchor).

The survey also maps every canonical role (programmer, reviewer, tester, architect, project manager, QA-checker) onto the workflow families and notes that **multi-agent workflows outperform single-agent workflows** on repository-scale tasks because (a) context budget is partitioned per role, (b) intermediate artifacts are inspectable, and (c) the blackboard-style shared state ([[1.0.0 PRIM-17]] MAgHARCM's `compiletime.State`) provides a writable, readable, scalable global context space.

## Findings Relevant to MAgHARCM

- **MAgHARCM is a hybrid pipeline + self-negotiation circular + self-evolving.** The survey treats these as separate families; MAgHARCM combines them in one graph: pipeline forward (archaeologist → ... → validator), self-negotiation circular for repair (verdict_panel ↔ translator via recruiter), and self-evolving for adaptive replanning (recruiter). This hybridization is the distinctive contribution MAgHARCM makes to the surveyed taxonomy.
- **Checkpoint barriers are the bottleneck of pipeline architectures.** The survey highlights that pipelines without explicit durable-state barriers lose intermediate state on interruption. MAgHARCM's two checkpoint lambdas (`save_translator_ckpt`, `save_validator_ckpt`) are the survey's recommended mitigation.
- **Role specialization matters more than model scale.** Across the surveyed systems, role-specialised 7B-13B SLMs outperform single 70B+ general-purpose models on repository-scale tasks, because the role-specialised prompt fits within the SLM's effective attention window.
- **Blackboard state is the recommended shared-memory pattern.** The survey singles out MetaGPT's SOP contracts (cf. [[1.0.0 PRIM-24]]) and Self-Collaboration (Dong et al. 2023) as the canonical blackboard implementations. MAgHARCM's `compiletime.State` plus `compiletime.DocumentWrapper[T]` typed intermediate artifacts is the same pattern.
- **Repair loops are universal across surveyed systems.** Every successful multi-agent code-generation system in the survey has an explicit repair loop. MAgHARCM's `repairBranch` (END vs verdict_panel) is the canonical instantiation.

## How MAgHARCM Uses It

The 10-node MAgHARCM graph (`internal/graph/graph.go`) is the canonical instantiation of the surveyed multi-agent taxonomy: pipeline forward (8 nodes), self-negotiation circular (verdict_panel ↔ translator via recruiter), self-evolving (recruiter adapts strategy via `SwitchToNextStrategy`), with explicit checkpoint barriers. The survey's bottleneck analysis is the empirical justification for the two checkpoint lambdas. The role-specialization finding is the empirical justification for the SLM-fleet design (cf. [[1.0.0 P-58]] Qwen2.5-Coder 7B-14B as the primary deployment).

## References

### Hop-1
- Qian, C. et al. (2023). *ChatDev: A Sociological Perspective on Code Generation*. arXiv:2307.07924. See [[1.0.0 ChatDev-2023]].
- Hong, S. et al. (2023). *MetaGPT: Meta Programming for a Multi-Agent Collaborative Framework*. arXiv:2308.00352. See [[1.0.0 PRIM-24]].
- Chen, W. et al. (2023). *AgentVerse: Facilitating Multi-Agent Collaboration*. arXiv:2308.10848. See [[1.0.0 P-26]].
- Saunders, W. et al. (2022). *Self-Refine: Iterative Refinement with Self-Feedback*. arXiv:2303.05671.
- Dong, Y. et al. (2023). *Self-Collaboration Code Generation via ChatGPT*. arXiv:2304.07590.
- Jiang, J. et al. (2024). *A Survey on Large Language Models for Code Generation*. arXiv:2406.00515. See [[1.0.0 P-50]].
- Qian, C. et al. (2023). *HyperAgent*. See [[1.0.0 P-13]].

### Hop-2
- Brown, T. B. et al. (2020). *Language Models are Few-Shot Learners*. arXiv:2005.14165.
- Wei, J. et al. (2022). *Chain-of-Thought Prompting*. arXiv:2201.11903.
- Park, J. S. et al. (2023). *Generative Agents: Interactive Simulacra of Human Behavior*. UIST 2023. arXiv:2304.03442.
- Qian, C. et al. (2023). *Experiential Co-Learning of Software-Developing Agents*. arXiv:2306.02325.

## Backlinks

[[1.0.0 PRIM-7]], [[1.0.0 PRIM-17]], [[1.0.0 PRIM-23]], [[1.0.0 PRIM-24]], [[1.0.0 PRIM-25]], [[1.0.0 PRIM-28]], [[1.0.0 PRIM-29]], [[1.0.0 P-13]], [[1.0.0 P-26]], [[1.0.0 P-50]], [[1.0.0 P-58]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-63 is the **methodological anchor** for MAgHARCM's 8-agent (10-node) graph: pipeline + self-negotiation circular + self-evolving, with explicit checkpoint barriers. The survey's bottleneck analysis justifies the checkpoint design.
