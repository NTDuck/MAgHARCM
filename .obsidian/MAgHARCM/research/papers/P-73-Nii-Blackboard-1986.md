---
title: "P-73 — Nii 1986 — Blackboard Systems: The Blackboard Model of Problem Solving"
backlink: "[[1.0.0 P-73]]"
aliases:
  - "1.0.0 P-73"
  - "P-73"
  - "P-73-Nii-Blackboard-1986"
  - "P-73-Nii-Blackboard-1986"
  - "Nii-Blackboard-1986"
tags: [paper, blackboard, multi-agent, problem-solving, knowledge-source, [[1.0.0 P-44]], hop-2]
---

# [[1.0.0 P-73 — Nii Blackboard Systems]]

## Citation

Nii, H. P. (1986). *Blackboard Systems: The Blackboard Model of Problem Solving and the Evolution of Blackboard Architectures*. AI Magazine 7(2):38-53. URL: https://aaai.org/ojs/index.php/aimagazine/article/view/537.

## Summary

The blackboard model is a **problem-solving architecture** where multiple specialised agents (called **knowledge sources, KS**) cooperate by writing to and reading from a shared data structure (the **blackboard**). The blackboard is partitioned into **panels** (sub-spaces), each holding a different kind of intermediate result.

The architecture has three components:
1. **Knowledge sources** — independent agents, each triggered by specific blackboard patterns. Each KS produces modifications to the blackboard; modifications can trigger other KSs.
2. **Blackboard** — the shared data structure, partitioned into panels. Each panel has a schema; modifications must conform.
3. **Control shell** — the scheduler that decides which triggered KS to invoke next.

The control shell is itself a meta-KS that can write to a separate "control panel" of the blackboard, enabling the system to reason about its own problem-solving strategy.

## Method

Nii's contribution is the **evolution taxonomy**: the paper traces blackboard systems from HEARSAY-II (speech recognition, 1971) through HASP (ocean surveillance, 1978) to the modern era, classifying each by:
- **Triggering mode**: data-directed (KS triggered by blackboard patterns) vs. goal-directed (KS triggered by goals).
- **Granularity**: single panel vs. multi-panel.
- **Control**: static (predefined KS scheduling) vs. dynamic (KSes bid on the next action via a control panel).

## Findings Relevant to MAgHARCM

- **[[1.0.0 P-44]] Corkill 1991** is the modern update of Nii's model. MAgHARCM's design draws on both Nii 1986 (foundational architecture) and Corkill 1991 (modern implementation patterns).
- **MAgHARCM's graph-of-agents** is a **typed blackboard with explicit panels**:
  - **Comprehension panel**: concept-clusters, design-rule hierarchy (produced by Archaeologist; consumed by Planner).
  - **Planning panel**: translation plan, dependency DAG (produced by Planner; consumed by Translator).
  - **Translation panel**: translated code fragments (produced by Translator; consumed by Validator).
  - **Verdict panel**: verdict per translated fragment (produced by Verdict Panel; consumed by Runner).
  - **Control panel**: which agent is active, retry counts, current strategy (written by Runner).
- **Each agent is a KS** with a typed trigger pattern: the Archaeologist triggers on the `SourceRepo` panel; the Planner triggers on the `Comprehension` panel; etc.
- **Dynamic control via the control panel** is the empirical justification for the **Registry.TryInOrder** strategy (cf. `internal/agents/strategy.go`): the runner writes to the control panel which strategy is being attempted, agents read it to decide whether to act.

## How MAgHARCM Uses It

`internal/agents/state.go` defines the **panel schemas**: `ConceptAssignment`, `DesignRuleHierarchy`, `TranslationPlan`, `TranslatedFragment`, `Verdict`, `ControlState`. Each panel is a struct in the agents state package, with a schema-versioned serializer.

The graph orchestrator (cf. `internal/graph/graph.go`) implements the **control shell**: it selects which agent to invoke next based on the current control-panel state and the strategy in `Registry.TryInOrder`.

## References

### Hop-1 (Nii 1986 cites)
- Reddy, R. et al. (1973). *HEARSAY-II Speech-Understanding System*. AI Magazine 4(2). (foundational blackboard)
- Erman, L. D. et al. (1980). *The Hearsay-II Speech-Understanding System*. Computing Surveys 12(2):213-253.
- Lesser, V. R. & Corkill, D. D. (1983). *The Distributed Vehicle Monitoring Testbed*. AI Magazine 4(3):63-109.
- Hayes-Roth, B. (1985). *A Blackboard Architecture for Control*. Artificial Intelligence 26(3):251-321.

### Hop-2
- Corkill, D. D. (1991). *Blackboard Systems*. AI Expert 6(9):40-47. See [[1.0.0 P-44]].
- Rao, A. S. & Georgeff, M. P. (1995). *BDI Agents: From Theory to Practice*. ICMAS 1995.
- Wooldridge, M. & Jennings, N. R. (1995). *Intelligent Agents: Theory and Practice*. Knowledge Engineering Review 10(2):115-152.

## Backlinks

[[1.0.0 P-44]], [[1.0.0 P-13]], [[1.0.0 P-11]], [[1.0.0 PRIM-7]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-73 is the **foundational multi-agent blackboard anchor** for MAgHARCM's graph-of-agents architecture.
