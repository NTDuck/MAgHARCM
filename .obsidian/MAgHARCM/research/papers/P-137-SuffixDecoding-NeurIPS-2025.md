---
title: "P-137 SuffixDecoding: Extreme Speculative Decoding for Emerging AI Applications"
backlink: "[[1.0.0 P-137]]"
aliases:
  - "1.0.0 P-137"
  - "P-137"
  - "P-137-SuffixDecoding-NeurIPS-2025"
  - "SuffixDecoding-NeurIPS-2025"
  - "oliaro2025suffixdecoding"
tags: [paper, speculative-decoding, suffix-tree, agentic-workload, "[[1.0.0 PRIM-7]]", "[[1.0.0 PRIM-21]]", wave-19]
date: 2026-09-07
last_updated: 2026-09-07
venue: NeurIPS 2025 (Spotlight)
---

# [[1.0.0 P-137]] SuffixDecoding

## TL;DR

SuffixDecoding is a **model-free speculative decoding** mechanism using suffix trees built over prompt + previous output tokens; speculation length adapts dynamically to acceptance likelihood. Deployed in **Snowflake ArcticInference** and **vLLM**. Achieves **up to 5.3× over vanilla** and **2.8× over EAGLE-2/3** on AgenticSQL; 2.5× on SWE-Bench. Anchors `[[1.0.0 PRIM-7]]` and `[[1.0.0 PRIM-21]]` by exploiting template repetitiveness in agent loops.

## Mechanism (Q2)

1. **Suffix-tree draft** — a suffix tree built over the prompt + previous output tokens proposes continuations whose leaves match recent outputs. No draft-model training required.
2. **Adaptive speculation length** — speculation length is dynamically tuned based on acceptance likelihood (recent acceptance ratio drives the next speculation window).
3. **Production-grade integration** — deployed at scale in Snowflake ArcticInference and vLLM; no GPU memory pressure beyond the suffix-tree overhead.

Together these give MAgHARCM's verdict-panel-as-sampling-strategy pattern a draft mechanism that exploits **structural repetitiveness** (tool-call boilerplate, JSON templates, repeated error phrasing) without training a draft model.

## Anchoring (Q3)

| Primitive | Pre-wave-19 behaviour | SuffixDecoding substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation | Required a draft model (EAGLE-3 / Medusa) or a hand-crafted template; both cost training time or engineering effort | Suffix-tree draft is training-free; built from prompt + output stream at runtime |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | TryInOrder registry lacked a code-template-aware draft path | Suffix-tree naturally captures template repetitiveness; gives the registry a draft path that exploits translator-pipeline templates |

## Hop-1 Citations

- P-57 Leviathan 2023 Speculative Decoding.
- P-114 Cai 2024 Medusa (heads-based alternative).
- P-78 / P-108 EAGLE-3 2025 (model-based draft baseline).

## Hop-2 Citations

- REST NAACL 2024 (retrieval-based speculation).
- Prompt Lookup / Token Recycling (2024 model-free baselines).
- Suffix tree data structure (classic, Weiner 1973).

## MAgHARCM integration

- **YAML config key**: `speculative.suffixdecoding.enabled: true`; `speculative.suffixdecoding.tree_max_depth: <int>`.
- **Implementation file**: `internal/speculative/suffix_decoder.go::NewSuffixDecoder` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-21]]`.

## Caveats

- **Tree construction cost** — building the suffix tree over a long prompt is O(N) time and O(N) memory; cost is amortised across speculation rounds but front-loaded.
- **Template dependence** — speedup assumes high template repetitiveness; novel-agent outputs (no repeated structure) see marginal gains.
- **Substring alignment** — token-level alignment between suffix-tree leaves and target tokenizer must be maintained; mismatch degrades acceptance rate.

## Source

- arXiv: 2411.04975.
- Venue: NeurIPS 2025 Spotlight (verified via NeurIPS proceedings.com index 085713-4211, OpenReview uwL0vbeEVn, CMU blog csd-phd-blog/2025/suffix-decoding).
