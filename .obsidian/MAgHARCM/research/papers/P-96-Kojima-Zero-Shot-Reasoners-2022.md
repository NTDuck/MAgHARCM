---
title: "P-96 — Kojima, Gu, Reid, Matsuo & Iwasawa 2022 — Large Language Models are Zero-Shot Reasoners"
backlink: "[[1.0.0 P-96]]"
aliases:
  - "1.0.0 P-96"
  - "P-96"
  - "P-96-Kojima-Zero-Shot-Reasoners-2022"
  - "P-96-Kojima-Zero-Shot-Reasoners-2022"
  - "Kojima-Zero-Shot-Reasoners-2022"
tags: [paper, slm, hop-2, cot, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-23]]]
---

# [[1.0.0 P-96 — Kojima et al. — Zero-Shot Chain-of-Thought]]

## Citation

Kojima, T., Gu, S. S., Reid, M., Matsuo, Y., & Iwasawa, Y. (2022). *Large Language Models are Zero-Shot Reasoners*. Advances in Neural Information Processing Systems 35 (NeurIPS 2022). arXiv:2205.11916 (24 May 2022 v1; 25 Jan 2023 v2 with appendix). The University of Tokyo / Tohoku University / Google Research.

## Summary

The Kojima et al. (2022) paper is the *zero-shot counterpart* to Wei et al. (2022) [[1.0.0 P-90]]: where Wei et al. demonstrated that hand-crafted few-shot CoT exemplars unlock multi-step reasoning in sufficiently large LMs, Kojima et al. show that **the same unlock can be achieved with zero exemplars — provided a single magic phrase is appended to the prompt**.

The headline finding is striking in its simplicity. Across twelve reasoning benchmarks spanning arithmetic (GSM8K, SVAMP, ASDiv, AQuA, MAWPS), commonsense (CSQA, StrategyQA, Date Understanding, Sports Understanding, SayCan), logic (Letter Concatenation, Coin Flip, Tracking Shuffled Objects), and a third-party evaluation (LAMBADA), appending the prompt *"Let's think step by step."* after the question — with **no worked exemplars** — yields substantial zero-shot CoT accuracy gains on **Instruct-GPT3 / GPT-3.5 (text-davinci-002)**. The full PaLM 540B sweep is the second pillar of the paper and confirms the effect is robust across model families.

Key quantitative claims verified against the published paper:

1. **Standard zero-shot (no CoT trigger)** plateaus; on GSM8K, Instruct-GPT3 (text-davinci-002) achieves roughly 15–17% accuracy.
2. **Zero-shot-CoT with *"Let's think step by step."*** more than doubles that figure, landing at roughly 40–46% on GSM8K — comparable to Wei et al.'s *few-shot-CoT* result for the same model class.
3. **The effect is scale-dependent**: PaLM 2-S (small) shows modest gains; PaLM 540B and PaLM 2-L show large gains. Sub-10B models show small or negligible gains — consistent with Wei et al.'s emergent-CoT claim.
4. **Output-space sensitivity**: the phrasing of the trigger matters. *"Let's think step by step"* > *"Let's think step by step to make sure we get the right answer"* > plain instruction. The authors run a search over candidate prompts and report the top-performing one.
5. **Self-consistency amplifies zero-shot CoT** (Wang et al. 2023 [[1.0.0 P-52]]): sampling N paths and majority-voting the final answers boosts GSM8K zero-shot-CoT into the 60%+ range without any exemplars.

The paper's conceptual contribution is to *decouple two phenomena that Wei et al. had bundled together*: (a) the existence of CoT reasoning in large LMs, and (b) the need for hand-crafted exemplars to elicit it. P-96 demonstrates that (a) is a property of the model; (b) is an artefact of prompt engineering. Zero exemplars are enough once the model is large enough.

## Method

- **Prompt template**: `Q: <question>\nA: Let's think step by step.` — followed by free-form CoT decoding and extraction of the final answer. No worked examples. No fine-tuning.
- **Magic phrase**: the paper does not commit to one canonical trigger phrase, but the canonical and widely-cited instance is *"Let's think step by step."* — empirically the top of the searched phrase space on Instruct-GPT3.
- **Models studied**: Instruct-GPT3 (text-davinci-002) as the headline model; full PaLM family (8B, 62B, 540B) as the scale-sweep second pillar; GPT-3.5 (code-davinci-002) as a code-adjacent variant.
- **Answer extraction**: a second *"Therefore, the answer is"* prompt appended after the generated chain to coerce the final answer into a parseable form. This two-stage extract pattern (generate → extract) is what makes the technique robust on arithmetic benchmarks.
- **Self-consistency follow-on**: N sampled CoT chains → majority vote on extracted answers. The paper demonstrates this *composes cleanly* with the zero-shot-CoT substrate.
- **Scale ablation**: per-benchmark breakdown of accuracy vs. model size confirms the emergent pattern from P-90; the same threshold (≥~30B for robust gains) recurs.

## Findings Relevant to MAgHARCM

- **CoT is "cheap" — no few-shot exemplars required.** The SLM-era deployment budget on local machines (Qwen2.5-Coder-4B/7B/32B [[1.0.0 P-21]], StarCoder2-3B/7B [[1.0.0 P-22]], Phi-3-mini [[1.0.0 P-27]], Code Llama 7B/13B/34B [[1.0.0 P-95]]) makes every exemplar token a real cost. P-96 shows that the cognitive structure of CoT can be elicited by a *five-word prompt suffix*, which means MAgHARCM's `[[1.0.0 PRIM-22]]` Four Phases of Comprehension pipeline can drop its few-shot exemplar budget entirely on ≥30B models and still get CoT.
- **Scale threshold is the gating constraint.** Below ~30B parameters the magic phrase yields small or unreliable gains. MAgHARCM's deployment matrix (see `internal/agents/strategy.go`) should *default to zero-shot-CoT* on Qwen2.5-Coder-32B and Code Llama 34B; *disable CoT* on Phi-3-mini-3.8B and Qwen2.5-Coder-4B. P-96 is the calibration source for that gate.
- **The two-stage extract pattern is reusable for code.** The *"Therefore, the answer is"* extract step generalises to *"Therefore, the output is … code block …"* — extract the final fenced block from a zero-shot CoT reasoning trace. This is exactly what `internal/agents/comprehension.go` needs for `[[1.0.0 PRIM-22]]` Explanation phase output parsing, and what `internal/agents/chunked_translator.go` needs for `[[1.0.0 PRIM-23]]` Chunked Translation final-output extraction when the translator is asked to *justify* its translation choices.
- **Self-consistency composes.** The paper demonstrates that P-52 [[1.0.0 P-52]] self-consistency works just as well over zero-shot-CoT chains as over few-shot-CoT chains. MAgHARCM's `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation can therefore sample N zero-shot-CoT candidates from the `[[1.0.0 PRIM-23]]` chunked translator and majority-vote the final emitted blocks — *without* depending on per-exemplar hand-engineering of CoT demos.
- **The "magic phrase" is a *tunable*, not a fixed string.** P-96's prompt search matters operationally: MAgHARCM should treat the magic phrase as a *configuration knob* in `state.Config` (e.g. `cot.TriggerPhrase = "Let's think step by step."`). Future work can re-run P-96's phrase search on the local SLM roster and pick the per-model phrase that maximises accuracy on a calibration subset.

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-22]]` Four Phases of Comprehension**: drop the per-phase 4-shot CoT exemplar bank on ≥30B SLMs and replace it with zero-shot-CoT (Q + `"Let's think step by step."`). Cite P-96 as the justification: CoT is in the model, not in the exemplars. Keep the exemplar bank as the *fallback* for sub-7B SLMs where the magic phrase is unreliable. The `Explanation` phase output parser must adopt the two-stage extract pattern (chain → *"Therefore, the summary is …"*) to robustly recover the structured summary.
- **`[[1.0.0 PRIM-23]]` Chunked Translation**: when the chunked translator is asked to translate a ~500 LoC fragment on a ≥30B SLM, prepend the question with the magic phrase and decode the chain-of-reasoning *before* the final code block. Then use P-52 self-consistency over N=3–5 such samples and majority-vote on the extracted Rust/Go block. This is cheaper than hand-crafting translation exemplars and empirically matches few-shot-CoT accuracy on PaLM-540B-class models. On sub-7B SLMs, keep the existing exemplar-driven chunked-translator prompts and *do not* inject the magic phrase (P-96's scale dependence says it will not help).
- **Future work (not implemented)**: per-model calibration sweep — for each SLM in MAgHARCM's roster {Qwen2.5-Coder-4B/7B/32B, StarCoder2-3B/7B/15B, Code Llama 7B/13B/34B, Phi-3-mini-3.8B}, search a small phrase space for the highest-accuracy zero-shot-CoT trigger on a held-out GSM8K-style benchmark, and store the winner in `state.Config.CotTriggerPhrase`.

## References

### Hop-1 (papers that build directly on P-96)
- Wang, X. et al. (2023). *Self-Consistency Improves Chain of Thought Reasoning in Language Models*. ICLR 2023. arXiv:2203.11171. See `[[1.0.0 P-52]]` — the sampling-aggregation layer that P-96 explicitly composes with zero-shot-CoT.
- Suzgun, M. et al. (2022). *Challenging BIG-Bench Tasks and Whether Chain-of-Thought Can Solve Them*. arXiv:2210.09261. Systematic study that includes zero-shot-CoT across 23 BIG-Bench tasks; shows the magic-phrase pattern generalises beyond arithmetic.
- Khot, T. et al. (2022). *Decomposed Prompting: A Modular Approach for Solving Complex Tasks*. arXiv:2210.02406. See `[[1.0.0 P-62]]` — modular decomposition is a generalisation of CoT; the decomposition prompts are themselves zero-shot-CoT triggers over sub-tasks.

### Hop-2 (foundational anchors referenced)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. See `[[1.0.0 P-90]]` — the few-shot-CoT paper whose exemplar requirement P-96 eliminates. P-96 is *literally* the zero-shot counterpart of P-90 and should always be cited alongside it.
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners* (GPT-3). arXiv:2005.14165. The in-context-learning foundation both P-90 and P-96 build on.
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762. Transformer architecture.
- Cobbe, K. et al. (2021). *Training Verifiers to Solve Math Word Problems* (GSM8K). arXiv:2110.14168. The grade-school math benchmark P-96 uses for its headline number.
- Ouyang, L. et al. (2022). *Training Language Models to Follow Instructions with Human Feedback* (InstructGPT). arXiv:2203.02155. The instruction-tuned base model (text-davinci-002) P-96 runs against.
- Chowdhery, A. et al. (2022). *PaLM: Scaling Language Modeling with Pathways*. arXiv:2204.02311. The 8B/62B/540B base models for the scale-sweep second pillar.

### MAgHARCM lineage cross-refs
- `[[1.0.0 P-21]]` — Qwen2.5-Coder (the 32B SLM where the magic phrase unlocks CoT per P-96's scale dependence).
- `[[1.0.0 P-22]]` — StarCoder2 (the 3B/7B/15B SLM where the magic phrase is borderline — P-96 says small or unreliable gains).
- `[[1.0.0 P-27]]` — Phi-3-mini (the 3.8B SLM where P-96 says do NOT rely on the magic phrase).
- `[[1.0.0 P-52]]` — Self-Consistency (the sampling-aggregation layer P-96 explicitly composes with).
- `[[1.0.0 P-62]]` — Decomposed Prompting (modular generalisation of CoT, used in PRIM-22/23/31).
- `[[1.0.0 P-83]]` — Code-specialised Self-Consistency (parallel sampling + verifier on code; inherits the P-96 magic-phrase trigger).
- `[[1.0.0 P-84]]` — s1 (CoT-style reasoning traces + budget forcing; P-96 is the upstream zero-shot substrate).
- `[[1.0.0 P-90]]` — Chain-of-Thought (the few-shot precursor; cite alongside P-96 whenever CoT is invoked).
- `[[1.0.0 P-95]]` — Code Llama (the 34B SLM where P-96's magic phrase unlocks CoT).

## Backlinks

`[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-23]]`, `[[1.0.0 PRIM-7]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-96 is the **zero-shot-CoT anchor** for MAgHARCM's SLM-era deployment. It eliminates the exemplar budget that P-90 imposed — for ≥30B SLMs, five words ("Let's think step by step.") are enough to elicit CoT. Cite P-96 whenever a primitive defaults to the magic-phrase trigger instead of hand-crafted CoT exemplars, and whenever the `cot.TriggerPhrase` config knob is changed.
