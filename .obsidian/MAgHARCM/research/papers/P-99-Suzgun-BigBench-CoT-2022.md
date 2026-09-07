---
title: "P-99 — Suzgun et al. 2022 — Challenging BIG-Bench Tasks and Whether Chain-of-Thought Can Solve Them"
backlink: "[[1.0.0 P-99]]"
aliases:
  - "1.0.0 P-99"
  - "P-99"
  - "P-99-Suzgun-BigBench-CoT-2022"
  - "P-99-Suzgun-BigBench-CoT-2022"
  - "Suzgun-BigBench-CoT-2022"
tags: [paper, big-bench, chain-of-thought, reasoning, slm, scale, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-23]], hop-1]
---

# [[1.0.0 P-99 — Suzgun et al. — BIG-Bench Hard / CoT Can Solve Them]]

## Citation

Suzgun, M., Scales, N., Schärli, N., Gehrmann, S., Tay, Y., Chung, H. W., Chen, Q., Gao, J., Han, B., Hart, J., Hashimoto, T., Johnson, R., Madaan, A., Mishra, P., Raghunathan, A., Rocktäschel, T., van den Driessche, G., Welbl, J., Wieting, C., Wei, J., Wu, Y., Yasunaga, M., & Zhou, D. (2022). *Challenging BIG-Bench Tasks and Whether Chain-of-Thought Can Solve Them*. Proceedings of the 61st Annual Meeting of the Association for Computational Linguistics (ACL 2023, Volume 1: Long Papers). arXiv:2210.09261 (17 Oct 2022 v1; 29 Jun 2023 v3 with camera-ready revisions). Google Research / Brown University / Stanford / Tübingen / Cohere AI.

## Summary

Suzgun, Scales, Schärli, Gehrmann, Tay et al. (2022) introduce **BIG-Bench Hard (BBH)**: a curated subset of **23 tasks** drawn from the 204-task BIG-Bench benchmark that were *solved by the vast majority of humans* but *consistently defeated language models* (including PaLM 540B and code-davinci-002) under standard few-shot prompting. The paper is the systematic sibling of Wei et al. 2022 [[1.0.0 P-90]]: where Wei et al. demonstrated that CoT unlocks a handful of reasoning benchmarks, Suzgun et al. ask how far that unlock extends across a broader reasoning portfolio.

The headline empirical claim is that **chain-of-thought prompting, applied with carefully written per-task exemplars, recovers a substantial fraction of the human-vs-model gap on these 23 "hard" tasks** for the instruction-tuned PaLM 540B (text-bison) and code-davinci-002 models. Concretely:

1. **Without CoT**: PaLM 540B under standard few-shot prompting either plateaus or fails on the 23 BBH tasks — average accuracy roughly matches or trails the human-rater baseline on a minority of tasks.
2. **With CoT**: the same PaLM 540B, prompted with task-specific CoT exemplars, **exceeds the average human-rater baseline on 17 of 23 BBH tasks**. On the remaining 6 tasks, CoT still substantially narrows the gap.
3. **Code-davinci-002 + CoT** also improves on every BBH task vs. its standard-prompting baseline, often dramatically (e.g. **logical-args** and **penguins-in-a-table** jump from near-zero to ~50%+ with CoT).

The paper's conceptual contribution is to define and characterise a specific *reasoning-envelope* in which CoT works: BBH is precisely the task set where the model's "latent competence" exceeds its "prompted expression" — CoT closes the gap by giving the model permission to *think out loud*. The remainder of the paper dissects the envelope: which tasks CoT unlocks cleanly, which tasks it leaves stuck, and what prompt-engineering choices matter.

A complementary finding is the **CoT self-consistency** ablation: sampling N CoT chains (temperature > 0) and majority-voting the final answers further raises accuracy on most BBH tasks — empirically confirming that Wang et al. 2023 [[1.0.0 P-52]] self-consistency composes cleanly with the per-task CoT templates BBH ships.

The paper also catalogues a **sub-class of BBH tasks that are "LLM-resistant" even with CoT**: tasks involving multi-step symbolic manipulation over long sequences (e.g. tracking shuffled objects across 7+ steps, formal-language deduction, counting-with-conditions). For these, CoT *helps* but does not close the gap to humans. This is the empirical anchor for MAgHARCM's SLM-envelope reasoning limits.

## Method

- **BBH task suite**: 23 tasks grouped into six reasoning families — (1) linguistic, (2) symbolic, (3) logical, (4) commonsense, (5) mathematical, (6) multi-modal-and-domain-specific. Each task comes with a "hard" filter: a task enters BBH only if state-of-the-art language models at the time of submission scored well below the human-rater baseline.
- **CoT exemplars**: 3-shot or 6-shot per task, hand-written. The paper ships every exemplar string in the appendix so the experiment is fully reproducible.
- **Models studied**: PaLM 540B (text-bison on the BIG-Bench API) and code-davinci-002 as the headline configurations; ablations on smaller PaLM variants.
- **Decoding**: greedy for the headline number; temperature > 0 with majority voting for the self-consistency ablation (N = 21 sampled chains).
- **Evaluation**: exact-match accuracy on the BBH test split; per-task breakdown; comparison to the published average-human-rater baseline per task.
- **Self-consistency** ablation: per task, run N = 21 CoT samples per query at temperature 0.7, extract the final answer, majority-vote, compare to greedy-CoT accuracy.
- **Failure-mode analysis**: the paper enumerates a "still-hard" sub-class of tasks (formal-language deduction, tracking shuffled objects, counting-with-conditions, Dyck languages, word-sorting, multistep-arithmetic) where CoT helps but the model still trails humans.

## Findings Relevant to MAgHARCM

- **The CoT envelope is bounded.** Wei et al. [[1.0.0 P-90]] established that CoT emerges with scale; Suzgun et al. **characterises the upper boundary of that envelope**. Even PaLM 540B, with task-specific CoT, does *not* close the gap to humans on tasks involving: (a) tracking 7+ permuted objects, (b) Dyck-language / formal-language parsing, (c) counting under combinatorial constraints, (d) formal logical deduction over multi-step rules. MAgHARCM's `[[1.0.0 PRIM-22]]` Four Phases of Comprehension must therefore *route these task families away from CoT-only reasoning* and toward explicit verifier-based or symbolic mechanisms. Cite P-99 whenever a comprehension task looks BBH-resilient.
- **The 23 BBH tasks are the canonical CoT-reasoning benchmark.** Suzgun et al. effectively defines "what reasoning means for LLMs in 2022" by enumeration. MAgHARCM's SLM evaluation fleet (Qwen2.5-Coder-32B [[1.0.0 P-21]], Code Llama 34B [[1.0.0 P-95]], Phi-3-mini [[1.0.0 P-27]], StarCoder2-15B [[1.0.0 P-22]]) should be **calibrated against BBH** — not against MMLU, not against GSM8K alone — when the deployment scenario asks "can this SLM reason about X?". BBH is the most informative per-task benchmark because it is curated to be the *failure modes* of standard prompting, which is exactly the regime MAgHARCM's translation and comprehension primitives operate in.
- **CoT + self-consistency amplifies BBH gains.** The paper's N = 21 majority-vote ablation confirms that P-52 [[1.0.0 P-52]] self-consistency composes with per-task CoT templates on BBH. MAgHARCM's `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation should be wired so that verdict candidates are sampled **with CoT exemplars enabled** — not just direct-decoding samples — when the underlying task family overlaps with BBH (symbolic, logical, multi-step).
- **Per-task CoT exemplars ship with the paper.** P-99's appendix is a *ready-made CoT exemplar bank* for the 23 BBH tasks. MAgHARCM can reuse these exemplars verbatim when building `[[1.0.0 PRIM-23]]` Chunked Translation prompts for tasks whose input distribution resembles a BBH task family (e.g. multi-step-arithmetic reasoning embedded in a code comment, formal-language type-system questions in a comprehension summary).
- **SLM-era envelope implication.** P-99's headline numbers are PaLM 540B / code-davinci-002. The paper does *not* sweep sub-30B models in detail, but the implicit reading — combined with the Wei [[1.0.0 P-90]] scale-dependence claim and Kojima [[1.0.0 P-96]] zero-shot-CoT scale-dependence claim — is that **on sub-7B SLMs the BBH envelope shrinks sharply**. MAgHARCM's `[[1.0.0 PRIM-22]]` Explanation phase should treat CoT as *unavailable* on Phi-3-mini-3.8B and Qwen2.5-Coder-4B for tasks that resemble BBH's hard sub-class. For ≥30B SLMs (Qwen2.5-Coder-32B, Code Llama 34B), BBH-derived CoT exemplars are usable.
- **The 17/23 threshold is a calibration anchor.** On the original 23 tasks, CoT exceeds the average human rater on 17. This is the *operational success rate* of CoT on a "hard reasoning" portfolio. MAgHARCM can use this as the upper-bound expectation when claiming "CoT unlocks X" for any new reasoning primitive: if the underlying task resembles a BBH-hard task and the SLM is ≥30B, expect ~70-75% chance CoT unlocks the task; otherwise expect substantial failure rate and budget a verifier / oracle accordingly.

## How MAgHARCM Uses It

- **`[[1.0.0 PRIM-22]]` Four Phases of Comprehension (SLM Envelope Calibration)**: maintain a routing table mapping comprehension-task descriptors → BBH task-family membership. For descriptors that map to BBH's "still-hard" sub-class (Dyck languages, tracking shuffled objects, multi-step arithmetic, formal deduction), route the comprehension question through an explicit verifier / oracle pathway (cf. `[[1.0.0 PRIM-12]]`) instead of CoT-only reasoning. Cite P-99 as the envelope-defining anchor. On ≥30B SLMs, BBH's 3-shot CoT exemplars from the paper appendix become a fallback bank when `[[1.0.0 PRIM-23]]` exemplar store has no closer match.
- **`[[1.0.0 PRIM-23]]` Chunked Translation (CoT self-consistency sampling)**: when the chunked translator is asked to translate a code chunk that contains embedded reasoning artifacts (e.g. arithmetic in a comment, a table-of-symbols table, a multi-step algorithm), sample N = 5–21 CoT chains with BBH-style per-task exemplars and majority-vote the extracted Rust/Go blocks. The 21-sample ablation in P-99 is the upper bound; MAgHARCM should pick the smallest N that retains ≥95% of the gain (typically N = 5). Wire this into `[[1.0.0 PRIM-7]]` Multi-Agent Verdict Validation as the default sampling strategy for BBH-overlapping inputs.
- **`[[1.0.0 PRIM-7]]` Verdict Panel (BBH-anchored test set)**: add a **BBH-mini subset** to the verdict panel's regression test corpus — at minimum the six BBH tasks Suzgun et al. flag as "still hard with CoT". These tasks should *always fail* in the regression suite when the verifier or panel regresses, because they are the empirical anchor of the SLM reasoning envelope. Cite P-99 when explaining why this subset is in the test corpus.
- **`state.Config.CotBbHardExemplarsPath`**: a new config knob that points to the P-99 appendix-derived per-task exemplar bank. Default ON for ≥30B SLMs, OFF for sub-7B SLMs. When ON, the comprehension / translation prompts may load 3-shot exemplars from this bank before falling back to the per-primitive exemplar store.
- **Future work (not implemented)**: a per-model BBH-mini sweep across MAgHARCM's SLM roster {Qwen2.5-Coder-4B/7B/32B, StarCoder2-3B/7B/15B, Code Llama 7B/13B/34B, Phi-3-mini-3.8B} to empirically measure where the SLM envelope boundary lies on each model. The P-99 finding that even 540B does not close the gap on BBH-hard sub-class is the *ceiling*; the SLM sweep is the *floor*.

## References

### Hop-1 (papers that build directly on P-99)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. See `[[1.0.0 P-90]]` — the foundational CoT paper P-99 systematically scales across 23 BBH tasks. P-99's per-task exemplars are the operational descendants of Wei et al.'s CoT format.
- Wang, X. et al. (2023). *Self-Consistency Improves Chain of Thought Reasoning in Language Models*. ICLR 2023. arXiv:2203.11171. See `[[1.0.0 P-52]]` — the sampling-aggregation layer that P-99's N = 21 majority-vote ablation explicitly composes with on BBH. P-52 + P-99 is the standard "CoT + ensemble" stack for hard reasoning.
- Suzgun, M. et al. (2023) follow-on work (T0 / FLAN-T5 instruction-tuned model variants) — explicitly tracks BBH per-task accuracy as a reasoning benchmark.

### Hop-2 (foundational anchors referenced)
- Wei, J. et al. (2022). *Chain-of-Thought Prompting Elicits Reasoning in Large Language Models*. arXiv:2201.11903. See `[[1.0.0 P-90]]` — directly cited as the substrate P-99 builds on (per-task exemplars extend the CoT format to 23 reasoning families).
- Wang, X. et al. (2023). *Self-Consistency Improves Chain of Thought Reasoning in Language Models*. ICLR 2023. arXiv:2203.11171. See `[[1.0.0 P-52]]` — directly cited as the sampling-aggregation layer P-99's majority-vote ablation uses.
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners* (GPT-3). arXiv:2005.14165. The in-context-learning foundation P-99's few-shot CoT templates build on.
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762. Transformer architecture.
- Srivastava, A. et al. (2022). *Beyond the Imitation Game: Quantifying and extrapolating the capabilities of language models* (BIG-Bench). arXiv:2206.04615. The 204-task parent benchmark from which BBH's 23 tasks are filtered.
- Chowdhery, A. et al. (2022). *PaLM: Scaling Language Modeling with Pathways*. arXiv:2204.02311. The base model P-99's headline numbers are reported on.
- Cobbe, K. et al. (2021). *Training Verifiers to Solve Math Word Problems* (GSM8K). arXiv:2110.14168. The grade-school-math benchmark whose multi-step-arithmetic task overlaps BBH.
- Hendrycks, D. et al. (2021). *Measuring Mathematical Problem Solving With the MATH Dataset*. arXiv:2103.03874. The competition-math benchmark whose word-problems overlap BBH's "logical-deduction" sub-class.

### MAgHARCM lineage cross-refs
- `[[1.0.0 P-21]]` — Qwen2.5-Coder (the 32B SLM where BBH-derived CoT exemplars become usable per P-99's scale envelope).
- `[[1.0.0 P-22]]` — StarCoder2 (the 3B/7B/15B SLM where the BBH envelope shrinks; cite P-99 when limiting CoT routing on these models).
- `[[1.0.0 P-27]]` — Phi-3-mini (the 3.8B SLM where P-99's envelope says CoT-on-BBH-substrate is unreliable).
- `[[1.0.0 P-52]]` — Self-Consistency (the sampling-aggregation layer P-99 explicitly composes with; MAgHARCM's verdict panel routing inherits this).
- `[[1.0.0 P-62]]` — Decomposed Prompting (modular generalisation of CoT; relevant when BBH-overlap is detected and the task can be decomposed into sub-tasks).
- `[[1.0.0 P-83]]` — Code-specialised Self-Consistency (parallel sampling + verifier on code; inherits BBH-derived CoT exemplars).
- `[[1.0.0 P-84]]` — s1 (CoT-style reasoning traces + budget forcing; P-99's per-task CoT templates are an upstream substrate).
- `[[1.0.0 P-90]]` — Chain-of-Thought (the foundational CoT paper P-99 systematically scales; cite alongside P-99 whenever CoT is invoked on hard reasoning).
- `[[1.0.0 P-95]]` — Code Llama (the 34B SLM where BBH-derived CoT exemplars are usable per P-99's scale envelope).
- `[[1.0.0 P-96]]` — Zero-Shot CoT (the magic-phrase substrate; P-99's per-task exemplars are an alternative for hard-reasoning tasks where the magic phrase alone is insufficient).

## Backlinks

`[[1.0.0 PRIM-7]]`, `[[1.0.0 PRIM-12]]`, `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-23]]`, `[[2.0.0 MAgHARCM]]`, `[[2.0.0 Software-Archaeology-Lineage]]`.

P-99 is the **SLM-era reasoning-envelope anchor** for MAgHARCM. It defines, by enumeration, the 23 reasoning tasks that constitute the upper boundary of what CoT-with-scale can solve — and the sub-class where even CoT-with-540B still fails. MAgHARCM's `[[1.0.0 PRIM-22]]` Four Phases of Comprehension routes comprehension questions through BBH-derived exemplars on ≥30B SLMs and routes the BBH-hard sub-class to verifier/oracle pathways instead. Cite P-99 whenever a primitive calibrates its CoT routing against a reasoning-task family, whenever BBH-mini is added to the verdict panel regression corpus, or whenever the SLM envelope boundary is invoked.
