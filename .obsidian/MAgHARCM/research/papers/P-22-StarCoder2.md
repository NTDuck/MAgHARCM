---
title: StarCoder 2 and The Stack v2 — The Next Generation
backlink: "[[1.0.0 P-22]]"
bibkey: p22_starcoder2
aliases:
  - "1.0.0 P-22"
  - "P-22"
  - "P-22-StarCoder2"
  - "P-22-StarCoder2"
  - "StarCoder2"
tags: [paper, code-llm, pretraining-data, [[PRIM-23]], [[PRIM-25]], hop-1]
---

# [[1.0.0 P-22]] StarCoder 2 and The Stack v2 — The Next Generation

**Authors**: Anton Lozhkov, Raymond Li, Loubna Ben Allal, Federico Cassano, Joel Lamy-Poirier, Nouamane Tazi, Ao Tang, Leandro von Werra, Harm de Vries (plus 53 co-authors from 32 institutions including Hugging Face, ServiceNow Research, Northeastern, Nvidia)
**Year**: 2024
**Venue**: arXiv preprint (cs.SE / cs.AI, 29 Feb 2024)
**eprint / DOI**: [[arXiv-2402.19173]]
**Cited by**: [[primitives/Primitives-Index]] entry [[PRIM-23]] (Chunked Translation), [[PRIM-25]] (Communicative-De-hallucination Role-Flip Gate)

## Summary

[[Lozhkov-2024-StarCoder2]] introduces the StarCoder2 family of code LLMs (3B / 7B / 15B parameters) trained on The Stack v2 — a 67.5 TB deduplicated, permissively-licensed source corpus built on top of the Software Heritage archive and augmented with GitHub issues, Jupyter notebooks, and code-completion data. The headline scientific claim is that, at fixed compute, a smaller model trained on more rigorously-cleaned data beats a larger model trained on noisier data: StarCoder2-3B matches StarCoder-15B on HumanEval while using ~5x less compute to train. The paper documents the full data pipeline — file-level deduplication, near-deduplication via MinHash, secret scanning, license filtering (only permissive OSS), and quality classification — and publishes the resulting corpus as The Stack v2. The 7B and 15B variants use grouped-query attention and a 16K-token context window.

## Relevance to MAgHARCM

StarCoder2 informs MAgHARCM's coding-model assumptions for the [[PRIM-23]] chunked translator: the smaller (3B) variant is the throughput-friendly fallback when the 8 GB-VRAM profile applies, and the report's empirical scaling laws justify MAgHARCM's empirical policy of preferring the smallest model that meets a fragment's compilation-pass rate. The Stack v2's permissive-license filter is the upstream precedent for MAgHARCM's [[PRIM-10]] (Feature-Mapping & Type-Compatibility Validation) licence audit: any third-party crate the manifest-rewriter ([[PRIM-10 partial]] lives in `internal/agents/manifest_rewriter.go`) introduces must be Apache/BSD/MIT-licensed, mirroring the Stack v2's license filter. StarCoder2 is also one of the candidate backbones for the [[PRIM-25]] role-flip reviewer agent when Qwen2.5-Coder is unavailable; the report's evaluation suite includes code-debugging benchmarks, supporting the same-family capability assumption. Finally, the report's explicit treatment of `near-dedup` and `secret-scan` filters informs the MAgHARCM eval harness's data-leakage audit (the harness must not have seen test inputs from CRUST-Bench during pre-training — The Stack v2's deduplication helps make this empirically defensible).

## Hop-1 References

- [[Qwen2.5-Coder-2024]] — sibling post-training approach (synthetic code-instruction data instead of raw-deduplication); together they define the modern open code-LLM landscape.
- [[DeepSeek-Coder-V2-2024]] — second major open code-LLM released in 2024; cites StarCoder2 as a baseline.
- [[Code-Llama-2023]] — Meta's earlier code-LLM; StarCoder2's data-pipeline advances over Code Llama's training set.
- [[BigCode-2023]] — the BigCode project's charter paper (predecessor to StarCoder2); documents the responsible-AI commitments StarCoder2 inherits.
- [[Mueller-2024-Code-LLM-Survey]] — independent survey of open code-LLMs; situates StarCoder2 in the broader landscape.

## Hop-2 Anchors (software-archaeology lean)

- [[Rajlich-1997]] — concept-locator analysis is the offline analogue of The Stack v2's repository-level file-level deduplication: both treat a project as a graph of concepts / files rather than a flat sequence.
- [[Kazman-2000]] — design-rule hierarchy (L1/L2/L3) predicts the Stack v2 should weight training samples by their layer, not uniformly — a refinement MAgHARCM could experiment with in offline fine-tunes.
- [[Mueller-2000]] — the Müller 5-strategy LIS analysis predicts that a legacy corpus (Stack v1) and a modernised corpus (Stack v2 filtered to permissive licences) co-exist, requiring an Integrate-in-Place strategy for any consumer straddling the two.
