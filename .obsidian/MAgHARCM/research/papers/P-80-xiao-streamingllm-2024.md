---
title: "P-80 — Xiao, Tian, Chen, Han & Lewis 2024 — Efficient Streaming Language Models with Attention Sinks (StreamingLLM)"
backlink: "[[1.0.0 P-80]]"
tags: [paper, long-context, kv-cache, attention-sink, streaming, slm, mitigation, lost-in-the-middle, [[1.0.0 P-53]], [[1.0.0 PRIM-31]], [[1.0.0 PRIM-22]], hop-1]
---

# [[1.0.0 P-80 — Xiao et al. 2024 — StreamingLLM (Attention Sinks)]]

- **Authors**: Guangxuan Xiao (MIT HAN Lab), Yuandong Tian (Meta AI), Beidi Chen (CMU / Meta AI), Song Han (MIT HAN Lab), Mike Lewis (Meta AI).
- **Venue / Year**: ICLR 2024 (arXiv:2309.17453, first posted September 2023, v2 ICLR-ready January 2024).
- **URL**: https://arxiv.org/abs/2309.17453
- **OpenReview**: https://openreview.net/forum?id=NG7sS51zVF
- **Code**: https://github.com/mit-han-lab/streaming-llm
- **Anchors**: [[1.0.0 P-53]] (Liu et al. 2023 — Lost-in-the-Middle, the empirical problem); mitigation lineage for [[1.0.0 PRIM-31]] (IterativeNavigator) and [[1.0.0 PRIM-22]] (Comprehension).

## 1. Core Contribution

Xiao et al. (2024) introduce **StreamingLLM**, a training-free framework that allows decoder-only LLMs trained with a finite attention window to generalize to **infinite / streaming sequence lengths** without performance collapse. The key empirical observation is the **attention sink** phenomenon: in softmax-attention transformers, a large fraction of the attention probability mass lands on the very first few tokens of the sequence — not because those tokens are semantically important, but because softmax must redistribute the "leftover" attention mass somewhere and the initial tokens form a stable, high-norm attractor.

StreamingLLM's cache layout is therefore:
- **Attention sinks**: KV cache entries for the first 4 tokens (typically including the BOS token) are kept permanently.
- **Sliding window**: the most recent $N$ tokens' KV states are kept (default $N \approx 1024$).
- **Eviction**: all tokens between the sinks and the window are discarded.

This yields constant KV-cache memory usage and stable perplexity on sequences up to **4 million tokens** (and beyond) across Llama-2, MPT, Falcon, and Pythia, with **up to 22.2× speedup** over the sliding-window-recomputation baseline.

## 2. Relationship to Liu et al. 2023 — Lost-in-the-Middle

[[1.0.0 P-53]] established the empirical law of long-context LLMs: **U-shaped attention degradation**, where performance drops 20–40% when the relevant information is in the middle of a long prompt. StreamingLLM is a **complementary** finding, not a replacement:

| Dimension | Liu et al. 2023 (Lost-in-the-Middle) | Xiao et al. 2024 (StreamingLLM) |
|---|---|---|
| Phenomenon | Relevant info in the middle of a long prompt is under-attended. | Initial tokens consume outsized attention mass ("sinks"); evicting them causes model collapse. |
| Cause | Query-side positional encoding / attention dilution over long contexts. | Softmax redistribution of attention mass; the first tokens attract the residual. |
| Mitigation strategy | Reorder context (relevant to start/end), or use RAG, or chunk. | Keep first-token KV permanently + sliding window of recent KV; discard middle. |
| Trade-off | Loses middle-context fidelity when the prompt exceeds the model's effective attention span. | Loses middle-context fidelity by design (cannot "remember" the middle), but keeps fluency and streamability. |
| Regime | Fixed finite context (e.g., 16K, 32K, 100K). | Infinite / streaming contexts. |

The two papers point at the **same underlying weakness** of dense softmax attention — the middle of a long sequence is poorly served — but propose different operational mitigations. Liu et al. says "re-order your prompt so the answer-relevant content is at the edges"; Xiao et al. says "the edges (the very first tokens) MUST be retained as anchors, but everything in the middle can be evicted because the model wasn't using it anyway."

**Crucial nuance**: StreamingLLM does NOT claim to recover lost-in-the-middle performance. The model genuinely cannot retrieve information from the middle of a sequence that has been evicted from the cache. What StreamingLLM recovers is **language-model fluency** and **stable perplexity** during streaming generation — i.e., the model continues to produce coherent text rather than collapsing to gibberish once the cache is full. This is the regime Liu et al. did not study: Liu et al. measure closed-book QA accuracy over a full long prompt; Xiao et al. measure perplexity over an infinite token stream.

## 3. Method

### 3.1 Attention-Sink Hypothesis

Let the softmax attention from query $q_t$ to key $k_i$ be
$$
\alpha_{t,i} = \frac{\exp(q_t^\top k_i / \sqrt{d})}{\sum_j \exp(q_t^\top k_j / \sqrt{d})}.
$$
For a decoder-only model, the keys $k_1, \dots, k_t$ correspond to the prefix plus the current token. As $t$ grows, the numerator for the **first** keys remains large because $q_t$ often has a high norm in those directions (early layers learn to keep a "global anchor" slot). The softmax denominator grows without bound, so any key that wants attention mass to sum to 1.0 must either earn it semantically or absorb it as residual. The initial tokens absorb the residual — they are **attention sinks**.

Empirically: even when the initial tokens are uninformative (e.g., a placeholder, repeated punctuation, or a partial word), the model still allocates ~10–30% of total attention to them. They are functionally a "parking lot" for the softmax remainder.

### 3.2 StreamingLLM Cache Layout

At each decoding step $t$:
- **Sink KV cache**: $K_{\text{sink}} = [k_1, k_2, k_3, k_4]$ — fixed size 4 (chosen empirically; 1–4 all work, 4 is the safe default).
- **Window KV cache**: $K_{\text{window}} = [k_{t-N+1}, \dots, k_t]$ — most recent $N$ tokens.
- **Attention**: $\alpha_{t,i}$ is computed only against $K_{\text{sink}} \cup K_{\text{window}}$.
- **Eviction**: tokens $k_5, \dots, k_{t-N}$ are dropped from the cache.

Total cache size is $4 + N$ entries regardless of $t$. Memory is **constant**; throughput is **constant**.

### 3.3 Optional Pre-Training Modification

Adding a dedicated **placeholder sink token** (e.g., a single extra token added to the vocabulary and always prepended to every input during pre-training) further stabilizes streaming deployment. Models trained with an explicit sink token show slightly lower perplexity in the streaming regime than models that rely on the BOS token as the implicit sink.

### 3.4 Reported Results

- Stable perplexity on **streams up to 4 million tokens** for Llama-2-7B/13B, MPT-7B/30B, Falcon-7B, Pythia-6.9B/12B.
- **22.2× speedup** over the sliding-window-recomputation baseline in dense-streaming settings.
- Loss of middle-context retrieval is **expected and accepted** — StreamingLLM is not a recall system; it is a fluency-preserving streaming decoder.
- Sinks-and-window framework has been integrated into Hugging Face Transformers, NVIDIA TensorRT-LLM, and the SwiftInfer engine.

## 4. Findings Relevant to MAgHARCM

### 4.1 Empirical corroboration of [[1.0.0 P-53]]'s U-Shape

StreamingLLM's "middle can be evicted" finding is a **mechanistic explanation** for why Liu et al.'s U-shape exists. If the model never really used the middle of a long prompt anyway (it just didn't have a sink there), then a Liu-style "put relevant info in the middle" strategy was always fighting against the softmax's own allocation bias. MAgHARCM's design — which puts critical content at the **start** (re-indexed symbols) and **end** (fresh navigations) of every prompt — is precisely the Liu-compatible strategy, and StreamingLLM now supplies the **mechanism** for why start/end placement works.

### 4.2 SLM Long-Context Mitigation

[[1.0.0 P-82]] and [[1.0.0 P-78]] establish that MAgHARCM's fleet runs on small language models (3B–7B), which have tighter memory and attention budgets than frontier LLMs. StreamingLLM is directly applicable:

- The comprehension agent's IterativeNavigator ([[1.0.0 PRIM-31]]) accumulates a **streaming history** of symbol resolutions across the run. Without StreamingLLM, this history either (a) overflows the SLM's context window and forces a costly re-summarization, or (b) triggers a context-window reset that loses prior reasoning state. With StreamingLLM, the navigator's KV cache stays bounded and the SLM remains coherent across arbitrarily long sessions.
- The Verdict Panel ([[1.0.0 PRIM-7]]) sees a **streaming stream of candidate proposals**. The same sink-and-window cache lets each panelist agent score many candidates without losing the global signal.

### 4.3 Training-Free Mitigation Composes with MAgHARCM's Verification Stack

Because StreamingLLM is **training-free** — it requires no fine-tuning, no LoRA, no RLHF — it composes cleanly with the rest of MAgHARCM's stack. A model that is already loaded for comprehension or verdict use can adopt StreamingLLM's cache layout by changing only the KV-cache manager; no new model artifact, no new training pipeline, no new evaluation. This matches MAgHARCM's "shipping wins" bias ([[1.0.0 PRIM-25]]): adopt the technique that costs the least engineering time and verify the empirical claim.

### 4.4 Caveat: StreamingLLM is NOT a Recall System

The most important caveat for MAgHARCM: StreamingLLM explicitly does **not** recover retrieval over the middle of an evicted span. If a comprehension agent needs to recall a symbol definition that was resolved 10 navigations ago and has since been evicted from the cache, StreamingLLM cannot help. The mitigation strategy for that case is **structured external state** (the comprehension manifest, the symbol-resolution graph) — not streaming attention. MAgHARCM's design already separates these concerns: structured state lives in the comprehension manifest; the SLM's KV cache holds only the recent working set. StreamingLLM is therefore an **in-context mitigation**, not a replacement for [[1.0.0 PRIM-31]]'s external index.

## 5. How MAgHARCM Uses It

Adoption plan in three layers:

1. **Inference runtime** — wrap MAgHARCM's SLM-fleet loader with the StreamingLLM KV-cache manager (sink count = 4, window size tuned per agent: navigator at 1024, verdict panelist at 512, role-flip reviewer at 768). Verified by perplexity stability over a 50K-token synthetic stream.
2. **Comprehension long-run support** — IterativeNavigator ([[1.0.0 PRIM-31]]) runs now have **unbounded effective session length**: the navigator's "rolling context" is held in the StreamingLLM cache instead of being periodically summarized. The comprehension manifest (external structured state) continues to hold the authoritative symbol resolution; the cache holds the working set.
3. **Evaluation** — add a streaming-coherence benchmark to the test suite: a 50K-token synthetic stream with periodic comprehension probes. Without StreamingLLM, perplexity should spike after the SLM's training-window length (typically 4K–8K); with StreamingLLM, perplexity should stay bounded. This becomes a regression test for any future SLM-fleet change.

Combined with the [[1.0.0 P-53]] re-ordering strategy (relevant symbols at start/end), StreamingLLM closes the long-context robustness gap for MAgHARCM's SLM fleet without requiring a single fine-tuning run.

## 6. References

### Hop-1 (Xiao et al. 2024 cites)
- Liu et al. 2023 — [[1.0.0 P-53]] — *Lost in the Middle* — empirical basis for U-shaped attention; StreamingLLM is a mechanistic complement.
- Vaswani et al. 2017 — *Attention Is All You Need* — softmax-attention backbone whose remainder-distribution property creates sinks.
- Touvron et al. 2023 — *Llama 2* — primary eval base model.
- Touvron et al. 2023 — *MPT* (MosaicML) — eval base model.
- Almazrouei et al. 2023 — *Falcon* — eval base model.
- Biderman et al. 2023 — *Pythia* — eval base model.
- Press et al. 2022 — *ALiBi* — positional encoding baseline; compared against as an alternative long-context enabler.
- Chen et al. 2023 — *Extending Context Window of Large Language Models via Positional Interpolation*.
- Mohtashami & Jaggi 2023 — *Landmark Attention* — alternative sparse-attention long-context strategy.
- Beltagy et al. 2020 — *Longformer* — sparse-attention baseline.

### Hop-2 (papers-cited-by-hop-1)
- Devlin et al. 2019 — *BERT* — transformer encoder baseline.
- Su et al. 2021 — *RoPE* (rotary position embeddings) — used by Llama / MPT.
- Dao et al. 2022 — *FlashAttention* — memory-efficient attention implementation.
- Raffel et al. 2020 — *T5* — encoder-decoder baseline.

## 7. Backlinks

[[1.0.0 P-53]], [[1.0.0 P-78]], [[1.0.0 P-82]], [[1.0.0 PRIM-7]], [[1.0.0 PRIM-22]], [[1.0.0 PRIM-25]], [[1.0.0 PRIM-31]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-80 is the **mechanistic complement** to [[1.0.0 P-53]]'s empirical finding, and the **training-free mitigation** that lets MAgHARCM's SLM fleet ([[1.0.0 P-78]], [[1.0.0 P-82]]) stream across arbitrarily long comprehension sessions without context-window collapse.
