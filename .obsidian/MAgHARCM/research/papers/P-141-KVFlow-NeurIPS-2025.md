---
title: "P-141 KVFlow: Efficient Prefix Caching for Multi-Agent LLM Serving via Workflow-Aware Eviction"
backlink: "[[1.0.0 P-141]]"
aliases:
  - "1.0.0 P-141"
  - "P-141"
  - "P-141-KVFlow-NeurIPS-2025"
  - "KVFlow-NeurIPS-2025"
  - "kvflow2025neurips"
tags: [paper, kv-cache, multi-agent, workflow-aware, "[[1.0.0 PRIM-31]]", "[[1.0.0 PRIM-21]]", wave-19]
date: 2026-09-07
last_updated: 2026-09-07
venue: NeurIPS 2025 (Poster)
---

# [[1.0.0 P-141]] KVFlow

## TL;DR

KVFlow introduces **workflow-aware KV cache eviction** using an Agent Step Graph that assigns each agent a "steps-to-execution" metric. Policy evicts high-steps-to-execution agents (not needed soon), preserves low-steps-to-execution agents (needed soon); adds KV prefetching that overlaps CPU→GPU transfer with computation. Achieves **1.83× single-workflow** and **2.19× concurrent-workflow** speedups vs SGLang hierarchical radix cache. Anchors `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement and `[[1.0.0 PRIM-21]]` Migration Strategy Selection.

## Mechanism (Q2)

1. **Agent Step Graph** — a workflow graph where each agent has a known activation schedule; "steps-to-execution" is a forward-looking metric from the current step.
2. **Workflow-aware eviction** — evicts KV blocks owned by agents with high steps-to-execution; preserves blocks owned by agents about to execute.
3. **Prefetching** — prefetches evicted blocks from CPU to GPU during computation, overlapping transfer with useful work.

Together these give MAgHARCM's 8-agent graph a **multi-agent workflow retention policy** that Wave-18 KVzip (P-128, single-query compression) does not address.

## Anchoring (Q3)

| Primitive | Pre-wave-19 behaviour | KVFlow substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement | LRU eviction across the 8-agent graph evicted blocks for soon-to-execute agents, causing unnecessary recompute | Workflow-aware eviction preserves blocks for soon-to-execute agents; iterative retrieval loop sees fewer recomputes |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | TryInOrder registry lacked an agent-activation-time signal for scheduling fast-path (cached) vs slow-path (fresh-prefill) strategies | Agent-activation-time metric drives the registry's scheduling decisions |

## Hop-1 Citations

- P-128 KVzip 2025 (query-agnostic eviction baseline).
- SGLang hierarchical radix cache (2024 system baseline).
- P-80 StreamingLLM 2024 (attention-sink eviction).

## Hop-2 Citations

- LRU cache replacement policy (classic).
- P-68 Shazeer 2019 MQA (KV sharing lineage).
- vLLM PagedAttention (2023 KV memory management).

## MAgHARCM integration

- **YAML config key**: `kv_cache.kvflow.enabled: true`; `kv_cache.kvflow.steps_to_exec_window: <int>`.
- **Implementation file**: `internal/kvflow/workflow_cache.go::NewWorkflowCache` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-31]]`, `[[1.0.0 PRIM-21]]`.

## Caveats

- **Workflow declaration** — the Agent Step Graph must be declared up front; dynamic workflow changes require graph re-declaration.
- **Prefetch bandwidth** — CPU→GPU prefetch bandwidth is bounded; very large KV blocks may not fully overlap.
- **Steps-to-execution accuracy** — the metric depends on accurate activation-time prediction; mis-prediction under- or over-evicts.

## Source

- arXiv: 2507.07400.
- Venue: NeurIPS 2025 Poster (OpenReview 5Iw1nDtYmT; NeurIPS poster 115164).
