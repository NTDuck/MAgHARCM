---
title: "P-122 ReasoningBank: Distilled Strategies for Persistent LLM Reasoning"
backlink: "[[1.0.0 P-122]]"
tags: [paper, slm, persistent-memory, strategy-distillation, bandit-policy, "[[1.0.0 PRIM-21]]", "[[1.0.0 PRIM-29]]", "[[1.0.0 PRIM-31]]", "[[1.0.0 P-115]]", "[[1.0.0 P-109]]", "[[1.0.0 P-111]]", "[[1.0.0 P-90]]", wave-16]
date: 2026-09-27
last_updated: 2026-09-27
venue: ICLR 2026 (tentative)
arxiv: 2509.25140
---

# [[1.0.0 P-122]] ReasoningBank

## TL;DR

Google Research introduces a **persistent-memory substrate** for LLM agents that distills successful strategies from prior runs into reusable `MemoryTriple` records. The substrate closes the cross-run learning gap that self-evolving agents assumed but did not implement, and is the canonical substrate for `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement.

## Mechanism (Q2)

The mechanism is the **`MemoryTriple` lifecycle**:

1. **Extract** — after each agent run, distill the trajectory into a `(problem-type, strategy, reward)` triple via a self-judge step (analogous to `[[1.0.0 P-90]]` Wei et al. 2022 Chain-of-Thought self-judge).
2. **Store** — persist the triple in a bounded memory store (default 256 triples; eviction by lowest-reward).
3. **Retrieve** — at the start of a new run, query the store for triples matching the new problem type and prepend them to the prompt.
4. **Update** — re-distill after each run; triples are versioned by `(problem-type, iteration)`.

The substrate is **time-axis persistent**: retrieval operates across runs, not just within a run. This is the missing dimension that `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement assumed but did not implement.

## Anchoring (Q3)

| Primitive | Pre-wave-16 behaviour | Wave-16 substrate (ReasoningBank) |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval | LM-feedback-driven refinement within a run | Strategy-distilled persistent memory across runs |
| `[[1.0.0 PRIM-29]]` Recruiter Agent | Greedy tool recall (no cross-run learning) | MaTTS compute-memory loop (ReasoningBank's MaTTS budget-aware query) |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | Blind try-and-fail registry | Informed-switching policy (Reward-weighted strategy recall) |

## Hop-1 Citations

- `[[1.0.0 P-90]]` Wei et al. 2022 — Chain-of-Thought Prompting Elicits Reasoning in Large Language Models (self-judge step).
- `[[1.0.0 P-111]]` Jimenez et al. 2024 — SWE-bench: Can Language Models Resolve Real-World GitHub Issues? (experimental setup).
- Schick et al. 2023 — Toolformer (analog of structured-tool memory; not currently in vault).

## Hop-2 Citations

- `[[1.0.0 P-109]]` OpenAI 2024 — SWE-bench Verified (held-out test set).
- `[[1.0.0 P-115]]` Wang et al. 2024 — OpenHands: An Open Platform for AI Software Developers as Generalist Agents (hierarchical multi-agent baseline).

## MAgHARCM integration

- **YAML config key:** `memory.distilled: true` (opt-in, default `false`).
- **Implementation file:** `internal/agents/memorystore/` (Sprint 2026-09-28 target; compile-time constants in `internal/compiletime/compiletime.go` already present: `MaxMemoryTriples = 256`, `MemoryTripleRewardEMA = 0.3`).
- **Affected primitives:** `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-29]]`, `[[1.0.0 PRIM-31]]`.

## Caveats

- **Venue confirmation pending**: arXiv 2509.25140 posted 2025-09 with `UnderReview: ICLR 2026` banner; if ICLR 2026 rejects, this P-NN is re-evaluated under Q1.
- **Self-judge step** inherits the well-known biases of CoT self-evaluation; ReasoningBank partially mitigates via the reward-weighted eviction policy but does not eliminate them.

## Source

- arXiv 2509.25140 — ReasoningBank: Distilled Strategies for Persistent LLM Reasoning (tentative ICLR 2026).
