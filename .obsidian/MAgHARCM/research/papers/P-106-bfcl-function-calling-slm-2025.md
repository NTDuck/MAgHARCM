---
title: "P-106 — BFCL Function Calling Leaderboard for SLM (Gorilla 2024-2025)"
backlink: "[[1.0.0 P-106]]"
tags: [paper, function-calling, tool-use, benchmark, slm, [[1.0.0 PRIM-25]], [[1.0.0 PRIM-26]], hop-2]
---

# [[1.0.0 P-106 — BFCL: Berkeley Function Calling Leaderboard]]

## Citation

Patil, S. et al. (2024-2025). *Berkeley Function Calling Leaderboard*. Gorilla LLM team, UC Berkeley.
URL: https://gorilla.cs.berkeley.edu/blogs/8_berkeley_function_calling_leaderboard.html

## Summary

BFCL is the **de facto standard for function-calling evaluation**. It tests models across:
- Simple, parallel, and multiple function calls
- Agentic multi-hop reasoning
- Error recovery
- Context/memory management
- AST evaluation vs executable verification
- Latency and cost tracking

**Key SLM finding (2025)**: 7B-13B fine-tuned models now match GPT-4-class frontier models on simple function-calling tasks, but **trail significantly on "knowing when not to call a tool" (relevance detection)** and multi-turn stateful conversations. This is the dominant gap for SLM-era tool-use agents.

## Relevance to MAgHARCM

- **AST evaluation directly maps to PRIM-7 Verdict Panel**: BFCL's AST-evaluation methodology validates MAgHARCM's verdict-panel structural-equivalence check.
- **Executable verification complements PRIM-12 Wasm Oracle**: BFCL's "executable verification" is the empirical anchor for wasm-based reference execution.
- **Relevance-detection gap → PRIM-25 Role-Flip Gate**: the gate's "do not call the tool" path is the explicit relevance-detection signal. This anchors the gate's design.
- **Performance tracking on BFCL** is a concrete metric for MAgHARCM's tool-calling accuracy. A new configurable `bfcl_eval_path` slot in `configs/agents.yml` would let downstream consumers gate on BFCL score.

## References (hop-1)

- [[1.0.0 P-85]] Yue et al. 2025 function-calling SLM gap (UNVERIFIED — this paper supersedes it)
- [[1.0.0 P-81]] Gorilla
- [[1.0.0 P-102]] SmallCode 4B SLM coding (uses BFCL-compatible tool calling)
- [[1.0.0 P-95]] Code Llama instruction tuning

## References (hop-2)

- ToolBench: https://github.com/openbmb/toolbench
- Tau-Bench: multi-turn customer-service agent benchmark
