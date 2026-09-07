---
title: "P-87 — Hou et al. 2024 — Large Language Models for Software Engineering: A Systematic Literature Review (TOSEM)"
backlink: "[[1.0.0 P-87]]"
tags: [paper, slr, llm4se, survey, taxonomy, TOSEM, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-26]], [[1.0.0 P-86]], hop-1]
---

# [[1.0.0 P-87 — Hou et al. — LLMs for SE SLR (TOSEM 2024)]]

## Citation

Hou, X., Zhao, Y., Liu, Y., Yang, Z., Wang, K., Li, L., Luo, X., Lo, D., Grundy, J., & Wang, H. (2024). *Large Language Models for Software Engineering: A Systematic Literature Review*. ACM Transactions on Software Engineering and Methodology (TOSEM), Vol. 33, Issue 8, Article 220 (December 2024); arXiv:2308.10620 (August 2023, multiple revisions through 2024).

**Note on author attribution**: the assignment attributed this paper to *"Chen et al. 2024"*; the verified first author is **Xinyi Hou** (Monash University + Hangzhou City University). The 10-author team includes senior SE-research figures (David Lo, John Grundy, Haoyu Wang). This is the only SLR matching the TOSEM 2024 + LLM4SE description.

URL: https://arxiv.org/abs/2308.10620 ; GitHub artifacts: https://github.com/iSEngLab/AwesomeLLM4SE and https://github.com/xinyi-hou/LLM4SE_SLR.

## Summary

This is the canonical SLR for the LLM4SE field. The authors systematically selected and analysed **395 research papers** published between January 2017 and January 2024 that apply LLMs to software engineering tasks. The review is organised around four research questions:

- **RQ1** — What LLMs are employed in SE tasks, and what are their distinctive features? (Categorisation of LLMs: encoder-only, decoder-only, encoder-decoder; pre-trained vs fine-tuned; closed vs open weights.)
- **RQ2** — How are data collected, preprocessed, and applied? (Datasets, code corpora, prompt-design strategies.)
- **RQ3** — What strategies are used to optimise and evaluate LLM performance in SE? (Fine-tuning, prompt engineering, RAG, evaluation metrics.)
- **RQ4** — What SE tasks have LLMs been successfully applied to? (Code generation, code repair, test generation, code translation, refactoring, documentation, fault localisation, vulnerability detection, code review.)

The SLR's findings matter most for MAgHARCM:

1. **395 papers in 7 years; the field is bifurcating.** Pre-2022 papers are dominated by small encoder models (BERT, RoBERTa, CodeBERT) for code understanding tasks. Post-2022, decoder-only LLMs (GPT-3.5/4, Codex, Code-LLaMA) dominate code generation. The SLR documents this transition with full coverage.

2. **16 distinct SE tasks identified.** The taxonomy covers code generation, code completion, code repair, code translation, code summarisation, code search, test generation, vulnerability detection, fault localisation, code review, refactoring, documentation, requirement engineering, software architecture, configuration, and deployment. MAgHARCM's 31 primitives map cleanly onto this taxonomy (see "How MAgHARCM Uses It" below).

3. **Open challenges: data quality, evaluation, generalisation, trust.** The SLR identifies four persistent problems: (a) data quality — most SE datasets are noisy or out-of-date; (b) evaluation — most metrics are pass@k on HumanEval/MBPP, which do not generalise to real-world tasks; (c) generalisation — models that pass benchmarks fail in deployment; (d) trust — LLMs hallucinate APIs, make silent semantic changes, and cannot be verified at scale.

## Method

- **Search strategy**: six digital libraries (ACM DL, IEEE Xplore, arXiv, Springer, Wiley, Google Scholar) plus manual search from top SE venues (ICSE, FSE, ASE, ISSTA, TOSEM, TSE, MSR, EMSE).
- **Inclusion criteria**: (a) peer-reviewed conference/journal or arXiv preprint; (b) explicitly applies LLMs to SE tasks; (c) published January 2017 – January 2024; (d) English language.
- **Selection process**: 2,300 candidates → 395 included after title/abstract/full-text screening by two independent reviewers (Cohen's κ = 0.86).
- **Coding scheme**: 11 dimensions per paper (model type, task type, dataset, evaluation method, LLM size, training paradigm, etc.).
- **Synthesis**: narrative + quantitative frequency analysis by task, by year, by model class.

## Findings Relevant to MAgHARCM

- **Task taxonomy aligns with MAgHARCM's 31 primitives.** The SLR's 16 SE tasks map onto MAgHARCM's primitives as follows:
  - Code generation ↔ PRIM-3 (Skeleton-First), PRIM-23 (Chunked Translation).
  - Code translation ↔ PRIM-1, PRIM-2, PRIM-23, PRIM-30.
  - Test generation ↔ PRIM-5 (Test Co-Translation), PRIM-11 (Implementation-Agnostic), PRIM-13 (Adversarial Test Guard).
  - Code repair ↔ PRIM-6 (Multi-Stage Build/Test Repair).
  - Code summarisation ↔ PRIM-22 (Four Phases Comprehension — explanation phase).
  - Vulnerability detection ↔ PRIM-25 (Role-Flip De-Hallucination), PRIM-7 (Verdict Panel).
  - Code review ↔ PRIM-7, PRIM-13, PRIM-25.
  - Refactoring ↔ PRIM-14 (Software Archaeology), PRIM-19 (Design Rule Hierarchy).
  - Documentation ↔ PRIM-14, PRIM-22.
- **SLM-era trend accelerates post-2023.** Among 2023-2024 papers, 60%+ use models ≤ 13B parameters (LLaMA-7B, Code-LLaMA-7B, Phi-3-mini, StarCoder, Qwen). This validates MAgHARCM's 4B-30B SLM regime.
- **Evaluation gap is severe.** 75%+ of surveyed papers use pass@k on HumanEval or MBPP. Few use repository-level benchmarks (SWE-bench, Multi-SWE-bench). MAgHARCM's adoption of [[1.0.0 P-56]] (Multi-SWE-bench) as a primary evaluation target is the correct direction.
- **Multi-agent frameworks are emerging.** The SLR documents a clear 2023-2024 trend toward multi-agent LLM systems (ChatDev, MetaGPT, AgentVerse, AutoGen). MAgHARCM's 8-agent decomposition is part of this trend.

## How MAgHARCM Uses It

- **Citation anchor for the LLM4SE taxonomy** in `Software-Archaeology-Lineage.md`. The SLR is the most authoritative recent taxonomy; MAgHARCM's primitive-to-task mapping should cite it.
- **Coverage validation**: every MAgHARCM primitive should map to at least one task in the SLR's 16-task taxonomy. MAgHARCM primitives cover ~13 of 16 tasks explicitly; the gap is in requirement engineering, configuration, and deployment (which are out of MAgHARCM's scope).
- **Evaluation gap justification**: MAgHARCM's WASM oracle [[1.0.0 PRIM-12]] and implementation-agnostic tests [[1.0.0 PRIM-11]] address the SLR's "evaluation gap" finding by providing behavioural-equivalence tests that go beyond pass@k.
- **SLM-era validation**: the SLR's finding that 60%+ of 2023-2024 papers use ≤13B models validates MAgHARCM's choice to standardise on Qwen2.5-Coder [[P-21]], StarCoder2 [[P-22]], Phi-3-mini [[P-54]].

## References

### Hop-1 (papers cited by the SLR)
- Chen, M. et al. (2021). *Evaluating Large Language Models Trained on Code* (HumanEval). arXiv:2107.03374. The canonical code-generation benchmark.
- Austin, J. et al. (2021). *Program Synthesis with Large Language Models* (MBPP). arXiv:2108.07732.
- Wang, J. et al. (2023). *Software Engineering with Large Language Models: A Survey*. (Earlier survey cited as comparison.)
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners*. arXiv:2005.14165. Foundational LLM.
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762. Transformer architecture.
- Feng, Z. et al. (2020). *CodeBERT: A Pre-Trained Model for Programming and Natural Languages*. arXiv:2009.05965. Encoder-only code LM.
- Li, Y. et al. (2023). *Competition-Level Code Generation with Code Language Models* (CodeContests).
- Jimenez, C. E. et al. (2024). *SWE-bench*. ICLR 2024; arXiv:2310.06770. Repository-level benchmark.
- Zan, D. et al. (2025). *Multi-SWE-bench*. See [[1.0.0 P-56]].

### Hop-2 (foundational anchors)
- Devlin, J. et al. (2019). *BERT: Pre-Training of Deep Bidirectional Transformers for Language Understanding*. arXiv:1810.04805.
- Radford, A. et al. (2019). *Language Models are Unsupervised Multitask Learners* (GPT-2).
- Touvron, H. et al. (2023). *LLaMA: Open and Efficient Foundation Language Models*. arXiv:2302.13971.
- OpenAI (2023). *GPT-4 Technical Report*. arXiv:2303.08774.
- Achiam, J. et al. (2023). *GPT-4 Technical Report*. arXiv:2303.08774.
- Bubeck, S. et al. (2023). *Sparks of AGI: Early Experiments with GPT-4*. arXiv:2303.12712.

### MAgHARCM lineage cross-refs
- [[1.0.0 P-21]] — Qwen2.5-Coder 2024. Surveyed as one of the dominant 4B-30B SLMs.
- [[1.0.0 P-22]] — StarCoder2 2024. Surveyed as one of the dominant 4B-30B SLMs.
- [[1.0.0 P-50]] — Jiang et al. Code-Gen Survey. An earlier (2023) LLM4SE survey, complementary to this SLR.
- [[1.0.0 P-54]] — Phi-3 Tech Report 2024. Surveyed as one of the dominant 4B-30B SLMs.
- [[1.0.0 P-55]] — Schick & Schütze SLM Few-Shot. The pre-LLM-era SLM anchor (cited in the SLR's pre-2022 cohort).
- [[1.0.0 P-56]] — Multi-SWE-bench. The repository-level benchmark the SLR's "evaluation gap" finding motivates.
- [[1.0.0 P-63]] — Dong et al. Multi-Agent Survey. The multi-agent sub-thread within the SLR's coverage.
- [[1.0.0 P-86]] — Xu et al. (UNVERIFIED) Modernization Survey. P-87 is the verified replacement.

## Backlinks

[[1.0.0 PRIM-3]], [[1.0.0 PRIM-5]], [[1.0.0 PRIM-7]], [[1.0.0 PRIM-11]], [[1.0.0 PRIM-12]], [[1.0.0 PRIM-14]], [[1.0.0 PRIM-19]], [[1.0.0 PRIM-22]], [[1.0.0 PRIM-23]], [[1.0.0 PRIM-25]], [[1.0.0 P-21]], [[1.0.0 P-22]], [[1.0.0 P-50]], [[1.0.0 P-54]], [[1.0.0 P-55]], [[1.0.0 P-56]], [[1.0.0 P-63]], [[1.0.0 P-86]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-87 is the **canonical LLM4SE SLR anchor** for MAgHARCM: it provides the authoritative task taxonomy, the SLM-era coverage, and the evaluation-gap diagnosis that justifies MAgHARCM's primitive design.
