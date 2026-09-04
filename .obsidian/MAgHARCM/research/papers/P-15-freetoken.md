---
title: FreeToken — Efficient Edge-Native MoE Serving with Bandwidth-Adaptive Execution
bibkey: p15_freetoken
tags: [paper, moe-serving, edge, bandwidth-adaptive, locality, hop-1]
---

# [[1.0.0 P-15]] FreeToken — Efficient Edge-Native MoE Serving with Bandwidth-Adaptive Execution

**Authors**: Shuo Yang (UC Berkeley, corresponding), Xiaoze Fan, Melissa Pan, Haocheng Xi, Zhe Wang, Shanlin Sun, Kurt Keutzer, Song Han, Matei Zaharia, Chenfeng Xu (UT Austin, co-advise), Ion Stoica (co-advise)
**Year**: 2026
**Venue**: arXiv preprint cs.DC, submitted 17 Aug 2026 (no conference acceptance surfaced as of 2026-09-04)
**eprint / DOI**: arXiv:2608.16157
**Cited by**: [[primitives/INDEX]] entry as the empirical justification for MAgHARCM's "local-first" execution stance (used to support the rationale behind the chunked translator ([[PRIM-23]]) and conversable checkpoints ([[PRIM-28]]) running against local open-weight models rather than hosted APIs).

## Summary

[[Yang-2026-FreeToken]] is an edge-native Mixture-of-Experts serving system that runs large open-weight MoE models on personal hardware by co-designing the prefill and decode stages. The core mechanism is bandwidth-adaptive execution: edge bandwidth, traditionally a fixed bottleneck, is turned into a runtime scheduling signal. During prefill, FreeToken double-buffers expert movement with computation — while the GPU evaluates the current layer, the next expert is staged on-device so that no layer stalls waiting for weight transfer. Across heterogeneous local hardware (mixed CPU / GPU / NPU tiers), FreeToken dynamically maps both computation and model state so that open-weight MoE models with tens of billions of parameters fit and run interactively. The reported claim is that personal machines can now serve frontier-scale MoE models in a "fully local" mode with latency comparable to small-API-hosted dense models, fundamentally shifting the deployment economics of multi-agent pipelines.

## Relevance to MAgHARCM

FreeToken is the methodological justification for the "local-first" constraint that MAgHARCM's design budget inherits. The chunked translator ([[PRIM-23]]) — one LLM call per ~500 LoC fragment, persistent state summary between calls — exists in part because each call must fit in local VRAM / edge bandwidth, exactly the regime FreeToken makes viable. The conversable checkpoint persistence ([[PRIM-28]]) and recruitment-adaptive planning ([[PRIM-29]]) both depend on being able to pause / resume a translation pipeline on a developer machine; that posture is only realistic if the inference runtime is local. FreeToken's bandwidth-adaptive scheduling also informs the design of the multi-stage build/test feedback repair ([[PRIM-6]]): when a long-running feedback loop is in progress, the executor can back off and let the validator consume a different expert without thrashing the host's bandwidth budget.

## Hop-1 References

- [[Frantar-2023-GPTQ]] — post-training quantisation; the dense-weight compression substrate FreeToken inherits for its expert pages.
- [[Shazeer-2017-Outrageously-Large-Networks]] — Sparsely-Gated Mixture-of-Experts; the original MoE router design FreeToken's bandwidth-adaptive execution refines for the edge-memory regime.
- [[Pope-2023-Efficiently-Scaling-Transformer-Inference]] — MoE inference optimisation at data-centre scale; the bandwidth-allocation pattern FreeToken rewrites for edge hardware.
- [[Han-2022-Efficient-Inf-Design-Space]] — survey of efficient on-device inference; positions FreeToken within the wider edge-LLM design space.

## Hop-2 Anchors (software-archaeology lean)

- [[Mueller-2022]] — "Cold Turkey" and "Integrate-in-Place" strategies implicitly assume the legacy system can keep running locally while being incrementally translated; FreeToken makes that local-only assumption viable for the inference side too.
- [[Baldwin-Clark-2000]] — FreeToken's expert-paging is a design-rule decomposition at the hardware layer: each expert is a leaf (L3) that can be swapped freely, while the router is the L1 interface and the layer schedule is the L2 subsystem.
- [[Kazman-Cai-2024]] — architectural recovery of legacy codebases benefits from running the recovered model locally during exploration; FreeToken is what makes a "local-first" comprehension harness interactive enough for the archaeologist to actually use it.
