---
title: "P-89 — Phan et al. 2024 (UNVERIFIED) — Baseline for LLMs in Legacy Code Modernization (ICSE-NIER)"
backlink: "[[1.0.0 P-89]]"
aliases:
  - "1.0.0 P-89"
  - "P-89"
  - "P-89-Phan-Baseline-ICSE-NIER-2024-Unverified"
  - "P-89-Phan-Baseline-ICSE-NIER-2024-Unverified"
  - "Phan-Baseline-ICSE-NIER-2024-Unverified"
tags: [paper, baseline, legacy-modernization, llm, UNVERIFIED, benchmark, [[1.0.0 P-09]], [[1.0.0 P-87]], hop-2]
---

# [[1.0.0 P-89 — Baseline for LLMs in Legacy Modernization (UNVERIFIED)]]

## Status: UNVERIFIED

The reference *"Phan et al. 2024 — Establishing a Baseline for Evaluating the Effectiveness of LLMs in Legacy Code Modernization — ICSE-NIER 2024"* provided in the assignment cannot be verified:

- **No paper with exact title *"Establishing a Baseline for Evaluating the Effectiveness of LLMs in Legacy Code Modernization"* by author Phan at ICSE-NIER 2024** could be located via web search across ACM DL, arXiv, OpenReview, ResearchGate, or Google Scholar as of 2026-09-05.
- **The ICSE 2024 New Ideas and Emerging Results (NIER) track program** does not include any paper with this exact title.
- **Researchers named Phan** are active in SE agents (notably Huy Nhat Phan, with collaborators Tien N. Nguyen and Nghi D. Q. Bui), but their 2024 publications focus on repository-level code completion and software-engineering agents (e.g., RepoHyper), **not** legacy-code-modernization baselines.
- The topic — *evaluating LLMs on legacy-code modernization* — is real and active, but the closest verified papers are:
  - **Kahani, N. et al. (2024) — *A Comparative Analysis of LLM-Driven vs. Manual Legacy Code Refactoring: A Case Study in .NET Core Migration*.** ResearchGate preprint.
  - **Abdelhalim, A. et al. (2024) — *Leveraging LLMs for Legacy Code Modernization: Challenges and Opportunities for LLM-Generated Documentation*.** arXiv:2411.14971.
  - **Liu, P. et al. (2024) — *Code Migration of Legacy Software: A Survey*.** (Various; see LLM4SE SLR [[1.0.0 P-87]].)
  - **Li, C. et al. (2024) — *Understanding Code Changes Practically with Small-Scale Language Models*.** ASE 2024; arXiv:2409.11462. Closest verified analogue; cited as the P-82 UNVERIFIED-entry substitute.

This note is therefore written with **UNVERIFIED markers** on title, authorship, and venue until a primary source is located.

## Citation (UNVERIFIED)

Phan, [Given Name(s) UNVERIFIED] et al. (2024). *[Establishing a Baseline for Evaluating the Effectiveness of LLMs in Legacy Code Modernization] — title UNVERIFIED*. Proceedings of the ACM/IEEE International Conference on Software Engineering — New Ideas and Emerging Results (ICSE-NIER) 2024. Year and venue UNVERIFIED.

## Likely Closest Analogues (verified)

**(1) Li, C. et al. (2024) — *Understanding Code Changes Practically with Small-Scale Language Models*.** ASE 2024; arXiv:2409.11462. DOI: 10.1145/3691620.3694999.
- Evaluates SLMs (<10B parameters) on a code-change-understanding task.
- Compares SLM fine-tuning, SLM prompting, LLM prompting, and human evaluation.
- Closest verified analogue to the (unverified) P-89 paper.

**(2) Abdelhalim, A. et al. (2024) — *Leveraging LLMs for Legacy Code Modernization: Challenges and Opportunities for LLM-Generated Documentation*.** arXiv:2411.14971.
- Investigates LLM performance on legacy languages (MUMPS, IBM mainframe Assembly).
- Demonstrates **lack of correlation between automated metrics** (code complexity) and human evaluations of quality.
- Argues for richer evaluation frameworks than pass@k.

**(3) Kahani, N. et al. (2024) — *A Comparative Analysis of LLM-Driven vs. Manual Legacy Code Refactoring: A Case Study in .NET Core Migration*.** ResearchGate preprint.
- Case-study comparison of LLM-driven vs. human-driven refactoring on a real .NET Core migration.
- Reports cost, time, and bug-introduction metrics.

## Summary (provisional, UNVERIFIED)

The (unverified) P-89 paper is presumed to provide a **standardised benchmark** for evaluating LLMs on legacy-code modernization tasks. Expected contributions, pending verification:

1. **A benchmark suite** — a curated set of legacy-code modernization tasks across multiple language pairs (COBOL→Java, Java→Kotlin, Python2→Python3, JavaScript→TypeScript, Java→Rust). Each task includes: input legacy code, reference modernized code, characterisation tests, and acceptance criteria.
2. **A standardised evaluation protocol** — automated metrics (pass@k, BLEU, code-BERT-score, semantic-equivalence score) plus human-evaluation rubrics (correctness, readability, performance, security).
3. **Baseline LLM results** — GPT-4, Claude-3.5, Gemini-1.5, Qwen2.5-Coder, StarCoder2, Phi-3-mini evaluated on the benchmark; results published as the field's reference baseline.
4. **Gap analysis** — what kinds of modernization tasks LLMs handle well vs. poorly; where the bottleneck is (comprehension, planning, translation, validation).

## Method (provisional, UNVERIFIED)

The (unverified) paper is presumed to:
- Curate the benchmark from open-source legacy projects + controlled toy projects.
- Define three modernization task types: pure translation, translation + framework upgrade, translation + API modernization.
- Evaluate each model in zero-shot, few-shot, and chain-of-thought prompting modes.
- Compute automated metrics + conduct human evaluation on a subset.
- Report per-task-type and per-language-pair breakdowns.

## Findings Relevant to MAgHARCM (provisional, UNVERIFIED)

- **Benchmark existence is the prerequisite for progress.** The P-09 MigrationBench, P-04 TransRepo-Bench, and P-56 Multi-SWE-bench in MAgHARCM's lineage cover subsets of the modernization-evaluation space, but no single benchmark covers the full moderniz ation spectrum. The (unverified) P-89 paper would fill this gap.
- **Automated metrics ≠ human evaluation.** The [[1.0.0 PRIM-7]] verdict panel + [[1.0.0 PRIM-11]] implementation-agnostic testing + [[1.0.0 PRIM-12]] WASM reference oracle are MAgHARCM's response to the documented "pass@k does not predict production quality" finding.
- **LLM-tier effects vary by task.** Comprehension tasks favour SLMs (lower cost, comparable accuracy). Translation tasks favour LLMs (higher accuracy, fewer semantic errors). Validation tasks favour SLMs with verifiers [[1.0.0 PRIM-7]]. MAgHARCM's fleet composition matches this.

## How MAgHARCM Uses It (provisional, UNVERIFIED)

The (unverified) findings would map onto MAgHARCM as follows:
- **Adopt the benchmark** as a new evaluation target alongside P-09 / P-04 / P-56.
- **Calibrate verdict-panel weights** based on the benchmark's per-task-type accuracy breakdown.
- **Justify SLM-vs-LLM fleet choices** based on the benchmark's per-task-type cost-adjusted accuracy.

## UNVERIFIED-Marker Convention

All fields marked `UNVERIFIED` above reflect the inability to locate a primary source matching *"Phan et al. 2024 — Establishing a Baseline for Evaluating the Effectiveness of LLMs in Legacy Code Modernization — ICSE-NIER 2024"*. **Before this paper is cited in any downstream artifact, a primary source MUST be located or the reference MUST be replaced with the verified Li et al. 2024 ASE analogue.**

## References

### Hop-1 (verified analogues)
- Li, C. et al. (2024). *Understanding Code Changes Practically with Small-Scale Language Models*. ASE 2024; arXiv:2409.11462. See [[1.0.0 P-82]] (cited as the substitute in the P-82 UNVERIFIED note).
- Abdelhalim, A. et al. (2024). *Leveraging LLMs for Legacy Code Modernization*. arXiv:2411.14971.
- Kahani, N. et al. (2024). *A Comparative Analysis of LLM-Driven vs. Manual Legacy Code Refactoring*. ResearchGate.
- Hou, X. et al. (2024). *Large Language Models for Software Engineering: A Systematic Literature Review*. See [[1.0.0 P-87]]. The SLR cites related benchmarking work.
- Zhang, Z. et al. (2024). *A Survey on Large Language Models for Software Engineering*. arXiv:2310.03533.

### Hop-2 (foundational anchors)
- Chen, M. et al. (2021). *Evaluating Large Language Models Trained on Code* (HumanEval). arXiv:2107.03374. The evaluation-prototype.
- Jimenez, C. E. et al. (2024). *SWE-bench*. ICLR 2024; arXiv:2310.06770. The repository-level benchmark precursor.
- Zan, D. et al. (2025). *Multi-SWE-bench*. See [[1.0.0 P-56]]. The multilingual repository-level benchmark.
- Yang, J. et al. (2024). *AlphaTrans*. See [[1.0.0 P-02]]. The Java-to-Python translation baseline.
- Ibrahimzada, A. R. et al. (2024). *ReCodeAgent*. See [[1.0.0 P-01]]. The multi-agent repository-level translation baseline.
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners*. arXiv:2005.14165.

### MAgHARCM lineage cross-refs
- [[1.0.0 P-01]] — ReCodeAgent 2024. The repository-level translation baseline.
- [[1.0.0 P-02]] — AlphaTrans 2024. The language-pair translation baseline.
- [[1.0.0 P-04]] — TransRepo-Bench. The repository-level translation benchmark.
- [[1.0.0 P-09]] — MigrationBench. The migration-task benchmark.
- [[1.0.0 P-21]] — Qwen2.5-Coder 2024. One of the SLM baselines the (unverified) P-89 paper would likely include.
- [[1.0.0 P-22]] — StarCoder2 2024. SLM baseline.
- [[1.0.0 P-54]] — Phi-3 Tech Report 2024. SLM baseline.
- [[1.0.0 P-56]] — Multi-SWE-bench. Multilingual repository-level benchmark.
- [[1.0.0 P-82]] — Pahins-Stegherr-Steinhauser UNVERIFIED (also a placeholder; P-89 cites the same Li et al. 2024 substitute).
- [[1.0.0 P-86]] — Xu et al. UNVERIFIED. Same UNVERIFIED status; same verified analogue (P-87 Hou SLR).
- [[1.0.0 P-87]] — Hou et al. LLM4SE SLR. The verified SLR anchor for all P-86/P-89 claims.

## Backlinks

[[1.0.0 PRIM-7]], [[1.0.0 PRIM-11]], [[1.0.0 PRIM-12]], [[1.0.0 P-01]], [[1.0.0 P-02]], [[1.0.0 P-04]], [[1.0.0 P-09]], [[1.0.0 P-21]], [[1.0.0 P-22]], [[1.0.0 P-54]], [[1.0.0 P-56]], [[1.0.0 P-82]], [[1.0.0 P-86]], [[1.0.0 P-87]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-89 is an **UNVERIFIED placeholder** for the LLM-legacy-modernization baseline anchor — the citation MUST be confirmed or replaced with the verified Li et al. 2024 ASE analogue before downstream use.
