---
title: Shehory & Kraus — Methods for Task Allocation via Agent Coalition Formation (Artificial Intelligence 1998)
backlink: "[[1.0.0 P-49]]"
bibkey: p49_shehory_kraus_coalition_1998
aliases:
  - "1.0.0 P-49"
  - "P-49"
  - "P-49-Shehory-Kraus-Coalition-1998"
  - "P-49-Shehory-Kraus-Coalition-1998"
  - "Shehory-Kraus-Coalition-1998"
tags: [paper, multi-agent, coalition-formation, task-allocation, [[1.0.0 PRIM-29]], hop-1]
---

# [[1.0.0 P-49]] Shehory & Kraus — Methods for Task Allocation via Agent Coalition Formation

**Authors**: Onn Shehory (Carnegie-Mellon University, Robotics Institute), Sarit Kraus (Bar-Ilan University / University of Maryland)
**Year**: 1998 (May)
**Venue**: *Artificial Intelligence*, Vol. 101, Issues 1-2, pp. 165-200 (Elsevier)
**eprint / DOI**: DOI 10.1016/S0004-3702(98)00045-9
**Cited by**: [[primitives/Primitives-Index]] entry [[PRIM-29]] (Recruitment-Adaptive Planning) — provides the canonical hop-1 academic reference for dynamic multi-agent coalition formation and task allocation.

## Summary

[[Shehory-1998-Coalition-Formation]] establishes the foundational algorithmic framework for *coalition formation* and *task allocation* in multi-agent systems. A coalition is a subset of the agents that jointly commit to executing a task; each task has a *capability requirement* (the set of skills it needs) and a *utility* (the reward for successful completion); each agent has a *capability endowment* (the set of skills it can provide). Coalition formation is the process of partitioning the agent population into disjoint coalitions, one per task, such that every task's capability requirement is covered by the union of its coalition members' capabilities and the total utility is maximized.

The paper makes four foundational contributions:

1. *Distributed coalition-formation algorithms*. Shehory & Kraus give polynomial-time algorithms that allow a population of autonomous agents to form coalitions by exchanging capability-bid messages, without central coordination. The algorithms are *anytime*: they produce a feasible coalition structure at any time during execution, and refine it as more messages are exchanged.
2. *Coalition value calculation*. The paper defines the *coalition value* as the marginal utility gain from adding an agent to a coalition, and shows that the utility must be *super-additive* (the joint utility of a coalition exceeds the sum of its members' individual utilities) for coalition formation to be worthwhile.
3. *Rationality and stability*. The paper proves that the coalition-formation algorithms converge to a *Nash-stable* partition (no agent can improve its utility by unilaterally switching coalitions) under super-additive utility, and that this stable partition is *individually rational* (every agent earns at least as much as it would by working alone).
4. *Application domains*. The paper surveys three canonical application domains — distributed vehicle routing, manufacturing scheduling, and software-agent team formation — and shows that the same coalition-formation framework applies across all three, with the domain-specific structure entering through the utility function and the capability graph.

The paper is the canonical academic anchor for *dynamic team formation* in multi-agent systems. It is the direct ancestor of every subsequent work on multi-agent task allocation, including the modern LLM-era framework's recruitment step. The paper's emphasis on *bounded rationality* and *computational tractability* (the algorithms are polynomial, not exponential, in the number of agents and tasks) is precisely the constraint that modern recruitment algorithms must respect.

## Relevance to MAgHARCM

1. [[PRIM-29]] (Recruitment-Adaptive Planning, `internal/agents/recruit.go`): the recruitment step that selects which agent (or tool) handles each translation chunk is a modern realization of coalition formation. Each chunk is a *task* with a capability requirement (the chunk's source-language constructs, target-language constructs, and contextual dependencies); each agent is a *worker* with a capability endowment (the agent's specializations, tools, and prompts); the recruitment algorithm partitions chunks into per-agent queues. The paper's anytime algorithm is the theoretical justification for the bounded-time recruitment step.
2. [[PRIM-17]] (Asynchronous SE Agent Blackboard): the blackboard trigger mechanism decides when a recruitment step is necessary; the recruitment step in turn forms a coalition for the newly-triggered task. Shehory & Kraus's stability result is the bridge between blackboard events and stable per-agent task assignments.
3. [[PRIM-24]] (SOP-Anchored Role-Artifact Schema): SOP nodes declare the role, the artefact, and the transition condition; the recruitment step maps each SOP-node transition to a coalition-formation instance where the role defines the capability endowment and the artefact defines the utility function.
4. [[P-26]] (AgentVerse-2023, dynamic recruitment): AgentVerse's post-iteration evaluation is a coalition-formation heuristic; Shehory & Kraus's framework is the academic ancestor that justifies the heuristic.
5. [[P-39]] (HuggingGPT-2023, controller-expert decomposition): HuggingGPT's controller routes tasks to expert models; the routing is a single-task coalition-formation instance (the controller is the sole coalition member).

## Hop-1 References

- [[Sandholm-1999-Distributed-Reasoning-Coalitions]] — distributed reasoning about coalitions; the algorithmic completeness results for coalition formation in uncertain environments.
- [[Dang-2006-Protocol-Based-Coalition-Formation]] — protocol-based coalition formation; the operational specification of the negotiation messages Shehory & Kraus describe informally.
- [[Chalkiadakis-2011-Coalition-Formation-Overview]] — the modern survey of coalition formation; positions Shehory & Kraus (1998) as the foundational paper.
- [[Horling-2001-Survey-Organizational-Structures-MAS]] — survey of organizational structures in multi-agent systems; positions coalition formation among the alternative organizational paradigms (hierarchies, markets, congregations).
- [[Sycara-1998-Perspectives-Agency-Cooperation]] — multi-agent cooperation perspectives; the broader theoretical context for coalition formation.
- [[Kraus-1997-Negotiation-Automated-Group]] — automated negotiation for coalition formation; the negotiation-protocol ancestor of the capability-bid exchange.

## Hop-2 Anchors (software-archaeology lean)

- [[Curtis-1992-Process-Modeling]] (`[[1.0.0 P-45]]`) — Curtis's process-model framework's role-assignment axis is the *static* counterpart of coalition formation: in a process model, roles are assigned by the SOP definition; in coalition formation, agents self-organize into coalitions dynamically. Both are role-allocation mechanisms.
- [[Corkill-1991-Blackboard-Systems]] (`[[1.0.0 P-44]]`) — blackboard control shell's scheduler selects which knowledge source to invoke next; coalition formation is a generalization where the scheduler also decides which *subset* of agents should jointly handle a task.
- [[MacCormack-2006-DSM]] (`[[1.0.0 P-42]]`) — DSM cycles predict which files will need to be translated together; coalition formation uses these co-translation clusters as capability requirements (a single agent cannot independently translate a DSM cycle).
- [[Muller-2000-Strategies]] (`[[1.0.0 P-35]]`) — Müller's migration strategies are *coarse-grained* tasks that may be allocated to separate coalitions (one coalition per strategy); recruitment dispatches the chosen strategy to its coalition.
- [[Baldwin-Clark-2000]] (`[[1.0.0 P-34]]`) — design rules can be reinterpreted as coalition-fairness constraints: an L1 design rule constrains which agents may collaborate without violating the interface contract.

## Backlinks

- [[METHODOLOGY]] §PRIM-29 — Shehory & Kraus (1998) is the canonical academic anchor for the Recruitment-Adaptive Planning primitive; complements [[P-26]] (AgentVerse) and [[P-39]] (HuggingGPT) in the multi-agent dynamic-recruitment lineage.
- [[primitives/Primitives-Index]] — [[1.0.0 PRIM-29]] status entry now references [[P-49]] alongside [[P-26]] (AgentVerse) and [[P-39]] (HuggingGPT).
- [[Software-Archaeology-Lineage]] §4 — coalition formation is the team-organization theory that PRIM-29's recruitment step instantiates and that Stage 3 of the modernization cycle consumes.
