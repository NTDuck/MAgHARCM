---
title: "P-95 — Roziere et al. 2023 — Code Llama: Open Foundation Models for Code"
backlink: "[[1.0.0 P-95]]"
aliases:
  - "1.0.0 P-95"
  - "P-95"
  - "P-95-Roziere-Code-Llama-2023"
  - "P-95-Roziere-Code-Llama-2023"
  - "Roziere-Code-Llama-2023"
tags: [paper, code-llama, code-llm, slm, foundation-model, meta-ai, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-23]], [[1.0.0 PRIM-26]], hop-1]
---

# [[1.0.0 P-95 — Roziere et al. — Code Llama]]

## Citation

Roziere, B., Gehring, J., Gloeckle, F., Sootla, S., Gat, I., Tan, X. E., Adi, Y., Liu, J., Remez, T., Rapin, J., Kozhevnikov, A., Evtimov, I., Bitton, J., Bhatt, M., Canton Ferrer, C., Grattafiori, A., Xiong, W., D{\'e}fossez, A., Copet, J., Azhar, F., Touvron, H., Martin, L., Usunier, N., Scialom, T., & Synnaeve, G. (2023). *Code Llama: Open Foundation Models for Code*. arXiv:2308.12950 (24 Aug 2023 v1; 31 Jan 2024 v3 with extended evaluations). Meta AI (FAIR). Released under a custom **Llama 2 Community License** (commercial use permitted below 700M monthly active users).

## Summary

Code Llama is Meta's open-weights code LLM family, derived from Llama 2 by continued pretraining on a code-heavy corpus. It is the foundational code-SLM release of 2023 and the immediate predecessor (in lineage terms) to MAgHARCM's evaluation fleet. The family covers three sizes — **7B, 13B, 34B** — and three variants per size — **base** (continued pretraining only), **Instruct** (instruction-tuned with self-instruct), **Python** (specialised to Python). Plus **Code Llama 7B / 13B / 34B-Instruct** and the **Code Llama 7B / 13B / 34B-Python** specialisations.

Key technical claims:

1. **Continued pretraining on 500B-1T tokens** of code-heavy data (Infilling objective via *causal masking* on 7B/13B and *next-token prediction* on 34B).
2. **Long context**: 100K-token context window via RoPE scaling. Code Llama-Instruct extends to **16K** by default and supports **100K** with modifications.
3. **Strong humanEval and MBPP numbers at release**: Code Llama-34B-Instruct scored **48.8% on HumanEval** and **57.0% on MBPP**, matching or exceeding GPT-3.5 at the time. The 7B variant achieved **34.8% / 50.2%** — competitive with much larger models of the prior year.
4. **Llama 2 Chat fine-tuning recipe**: the Instruct variants use the same *self-instruct + rejection sampling + RLHF* pipeline as Llama 2 Chat, applied to code prompts.

Operational contribution to MAgHARCM:
- Code Llama-7B and CodeLlama-13B are the smallest viable SLMs in the MAgHARCM evaluation fleet — they fit on consumer GPUs and run locally via Ollama/GGUF.
- Code Llama-34B is the upper end of the SLM range that MAgHARCM targets. Combined with `[[1.0.0 P-90]]` CoT and `[[1.0.0 P-93]]` DPO, a 34B-Instruct variant is a strong baseline for comprehension, translation, and review.

## Method

- **Base model**: Llama 2 (7B, 13B, 34B) — public release from Meta, July 2023.
- **Continued pretraining**: 500B (7B/13B) to 1T (34B) tokens of code-heavy data, longer sequences (16K), with the original Llama 2 pre-training objective.
- **Infilling objective** (7B/13B only): causal masking — predict missing tokens given bidirectional context. This supports code completion at arbitrary positions.
- **Long-context fine-tuning**: RoPE scaling parameters modified to extend the effective context from 4K (Llama 2 default) to 100K.
- **Instruction tuning**: Code Llama-Instruct variants use self-instruct data + Llama 2 Chat recipe (RLHF via PPO); some ablations use rejection sampling without RLHF.
- **Python specialisation**: Code Llama-Python is fine-tuned on 100B Python-specific tokens, achieving the highest HumanEval scores within the family at the time of release.

## Findings Relevant to MAgHARCM

- **Code Llama is the Meta-side baseline for PRIM-22/23/31.** The 7B/13B variants set the **lower edge** of MAgHARCM's SLM fleet; Qwen2.5-Coder `[[1.0.0 P-21]]`, DeepSeek-Coder [[1.0.0 P-18]], StarCoder2 `[[1.0.0 P-22]]` set the middle; Code Llama-34B and Qwen2.5-Coder-32B set the **upper edge**.
- **34B is the upper bound for "local-feasible" SLMs.** A 34B-Int4 model fits on a single 24GB consumer GPU (e.g., RTX 4090, A5000); at Q4_K_M quantisation a 34B model is ~20GB. MAgHARCM's "local 4B-30B+" fleet includes the 34B variant.
- **Long context unlocks repository-scale code tasks.** 100K tokens is the right order of magnitude for repository-scale comprehension (`[[1.0.0 PRIM-22]]`) and chunked translation (`[[1.0.0 PRIM-23]]`); MAgHARCM's bounded 4KB fragment context is a deliberate, much-tighter slice to stay under the SLM attention budget.
- **Self-instruct + RLHF works at SLM scale.** The Code Llama-Instruct variants demonstrate that the Llama 2 Chat recipe (self-instruct + RLHF) is deployable at 7B-34B — but `[[1.0.0 P-93]]` DPO is a cheaper, simpler alternative for SLM alignment, and `[[1.0.0 P-94]]` LIMA shows the data-curation discipline matters more than the RLHF infrastructure.
- **Continued pretraining is necessary.** Llama 2 → Code Llama is a ~10pp HumanEval jump on the base model. Continued pretraining on code-heavy data is the dominant lever; instruction tuning adds ~5-10pp on top. MAgHARCM's SLM-tuning pipeline should respect this ordering.
- **Python specialisation is valuable but narrows coverage.** Code Llama-Python beats the general model on Python tasks but underperforms on multi-language code. MAgHARCM's translation targets include Java, C#, C, C++, Rust — the general Code Llama is the better baseline.

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-22]]` Four Phases of Comprehension**: when the comprehension agent dispatches to a 7B-13B Code Llama variant, apply **compressed** CoT exemplars (per `[[1.0.0 P-90]]`'s scale-dependence finding). When dispatching to Code Llama-34B, full 6-8-shot CoT.
- **`[[1.0.0 PRIM-23]]` Chunked Translation**: Code Llama-Instruct at 13B-34B is the recommended translation-model slot in `internal/agents/chunked_translator.go` when Qwen2.5-Coder is not available. The 100K context window supports large-context chunked translation; MAgHARCM uses the 4KB bounded context for cost reasons.
- **`[[1.0.0 PRIM-26]]` Symbol-Aware Navigator**: Code Llama's 100K context window informs the upper bound on navigator excerpts; current cap is 4 KB per excerpt.
- **`[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination Gate**: Code Llama-Instruct is the recommended reviewer-model on the role-flip gate when a stronger SLM (Qwen2.5-Coder-32B-Instruct) is unavailable.
- **`[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement**: Code Llama's long-context support enables fewer retrieval iterations — entire source files can be in context for the first retrieval pass.
- **Future work (not implemented)**: benchmark Code Llama-34B-Instruct against Qwen2.5-Coder-32B-Instruct on the comprehension/translation/review tasks; report the deltas in the next evaluation sprint.

## References

### Hop-1 (papers that build directly on P-95 or are direct descendants)
- Roziere, B. et al. (2023). *Unsupervised Translation of Programming Languages*. NeurIPS 2020 (cf. `p32_codellama` in `refs.bib` for the **Transcoder** lineage — same authors). The trans-coder paper that birthed the Roziere/Bhatt/Synnaeve research thread.
- Hui, B. et al. (2024). *Qwen2.5-Coder Technical Report*. arXiv:2409.12186. See `[[1.0.0 P-21]]` — the 2024 successor in the same SLM-code-LLM slot.
- Lozhkov, A. et al. (2024). *StarCoder 2 and The Stack v2*. arXiv:2402.19173. See `[[1.0.0 P-22]]` — competing open code-LLM family.
- Guo, D. et al. (2024). *DeepSeek-Coder*. arXiv:2401.14196. See `[[1.0.0 P-18]]` — competing open code-LLM family with stronger HumanEval at the 33B scale.
- Abdin, M. et al. (2024). *Phi-3 Technical Report*. arXiv:2404.14219. See `[[1.0.0 P-27]]` — competing small-model family (Phi-3-mini at 3.8B matches Code Llama-7B).

### Hop-2 (foundational anchors referenced)
- Touvron, H. et al. (2023). *Llama 2: Open Foundation and Fine-Tuned Chat Models*. arXiv:2307.09288. The base model Code Llama continues pre-training on.
- Su, J. et al. (2022). *RoFormer: Enhanced Transformer with Rotary Position Embedding*. The RoPE mechanism Code Llama uses for long context.
- Chen, M. et al. (2021). *Evaluating Large Language Models Trained on Code* (HumanEval). arXiv:2107.03374. The benchmark Code Llama is evaluated on.
- Austin, J. et al. (2021). *Program Synthesis with Large Language Models* (MBPP). arXiv:2108.07732. The benchmark Code Llama is evaluated on.
- Ouyang, L. et al. (2022). *Training Language Models to Follow Instructions with Human Feedback* (InstructGPT / RLHF). arXiv:2203.02155. The instruction-tuning recipe Code Llama-Instruct inherits.

### MAgHARCM lineage cross-refs
- `[[1.0.0 P-18]]` — DeepSeek-Coder (the 33B Code Llama alternative).
- `[[1.0.0 P-21]]` — Qwen2.5-Coder (the 32B Code Llama successor of choice for MAgHARCM).
- `[[1.0.0 P-22]]` — StarCoder2 (the 3B-15B Code Llama alternative).
- `[[1.0.0 P-27]]` — Phi-3 (the 3.8B Code Llama-class small-model alternative).
- `[[1.0.0 P-90]]` — Chain-of-Thought (the prompting substrate Code Llama benefits from at 34B).
- `[[1.0.0 P-93]]` — DPO (the cheaper Code Llama-Instruct alternative).
- `[[1.0.0 P-94]]` — LIMA (the Code Llama data-curation discipline).

## Backlinks

`[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-25]]`, `[[1.0.0 PRIM-26]]`, `[[1.0.0 PRIM-31]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-95 is the **Meta-side code-SLM anchor** for MAgHARCM. Code Llama at 7B/13B/34B establishes the lower, middle, and upper edges of the local SLM fleet and provides the long-context (100K) capabilities that make repository-scale code tasks feasible. Cite P-95 whenever MAgHARCM uses a Code Llama-family SLM or discusses the SLM-fleet baseline.
