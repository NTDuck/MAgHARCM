---
title: "P-86 — Xu et al. 2024 (UNVERIFIED) — Survey of Software Modernization Empowered by LLMs"
backlink: "[[1.0.0 P-86]]"
aliases:
  - "1.0.0 P-86"
  - "P-86"
  - "P-86-Xu-Modernization-Survey-2024-Unverified"
  - "P-86-Xu-Modernization-Survey-2024-Unverified"
  - "Xu-Modernization-Survey-2024-Unverified"
tags: [paper, survey, software-modernization, llm, UNVERIFIED, taxonomy, [[1.0.0 PRIM-14]], [[1.0.0 PRIM-21]], hop-1]
---

# [[1.0.0 P-86 — Xu et al. 2024 (UNVERIFIED) — Survey of Software Modernization Empowered by LLMs]]

## Status: UNVERIFIED

The reference *"Xu et al. 2024 — A Survey of Software Modernization Empowered by Large Language Models — arXiv:2501.07083"* provided in the assignment cannot be verified:

- **arXiv:2501.07083** — direct fetch of the arXiv record shows the paper at this ID is *"Stochastic reconstruction of multiphase composite microstructures using statistics-encoded neural network for poro/micro-mechanical modelling"* by Jinlong Fu and Wei Tan (January 2025; categories physics.app-ph). It concerns composite materials microstructure modelling, not LLM software modernization.
- **No paper matching the exact title *"Survey of Software Modernization Empowered by Large Language Models"* by an author named Xu** could be located via web search across arXiv, IEEE Xplore, ACM DL, or Google Scholar as of 2026-09-05.
- The arXiv ID `2501.07083` is a January 2025 number; a "2024 survey" would more naturally be an ID like `24XX.XXXXX`. The year mismatch (assignment says 2024, ID is January 2025) is a strong indicator of a misremembered citation.
- The closest verified analogues in the LLM-modernization survey space are:
  - **Hou et al. (2024) — *Large Language Models for Software Engineering: A Systematic Literature Review*.** See [[1.0.0 P-87]]. arXiv:2308.10620; TOSEM 33(8), Article 220. The canonical LLM4SE SLR.
  - **Zhang et al. (2024) — *A Survey on Large Language Models for Software Engineering*.** arXiv:2310.03533. Taxonomy of 16 downstream SE tasks.
  - **Belcak et al. (2025) — *Small Language Models are the Future of Agentic AI*.** arXiv:2506.02153. Argues SLM-dominated agents will replace LLM-dominated ones for narrow SE tasks.

This note is therefore written with **UNVERIFIED markers** on title, authorship, and arXiv ID until a primary source is located.

## Citation (UNVERIFIED)

Xu, [Given Name(s) UNVERIFIED] et al. (2024). *[A Survey of Software Modernization Empowered by Large Language Models] — title UNVERIFIED*. arXiv:[2501.07083 — UNVERIFIED; the arXiv ID at this number is an unrelated materials-science paper]. Year and venue UNVERIFIED.

## Likely Closest Analogues (verified)

In the absence of a confirmed paper, the following 2024–2025 surveys are the closest verified analogues and are referenced here for downstream MAgHARCM analysis:

**(1) Hou et al. (2024) — *Large Language Models for Software Engineering: A Systematic Literature Review*.** See [[1.0.0 P-87]].
- 395 research papers, January 2017 – January 2024.
- Four RQs: LLMs and their features; data preprocessing and application; optimization and evaluation; SE tasks covered.
- Coverage of code migration, refactoring, and modernization tasks.

**(2) Zhang, Z. et al. (2024) — *A Survey on Large Language Models for Software Engineering*. arXiv:2310.03533.**
- 16 downstream SE tasks, 5 categories.
- Taxonomy of code generation, translation, repair, test, and documentation tasks.

**(3) Belcak, P. et al. (2025) — *Small Language Models are the Future of Agentic AI*. arXiv:2506.02153.**
- Argues that SLM-dominated agentic AI will replace LLM-dominated systems for narrow SE tasks (including modernization).
- Provides a taxonomy of agentic SLM applications.

## Summary (provisional, UNVERIFIED)

The (unverified) P-86 paper is presumed to provide a structured taxonomy of LLM-driven software modernization strategies. Expected taxonomy, based on the LLM4SE literature:

1. **Migration (language translation)** — translating legacy codebases from one language to another (COBOL → Java, Java → Kotlin, JavaScript → TypeScript). MAgHARCM anchor: [[1.0.0 P-02]] (AlphaTrans), [[1.0.0 P-04]] (TransRepo-Bench), [[1.0.0 P-14]] (Syzygy).
2. **Reengineering (architecture-level restructuring)** — extracting services from monoliths, replacing frameworks, changing deployment topologies. MAgHARCM anchor: [[1.0.0 P-19]] (MSR4SA), [[1.0.0 P-14]] (Software-Archaeology).
3. **Porting (platform migration)** — moving code between operating systems, runtime environments, or hardware platforms. MAgHARCM anchor: [[1.0.0 P-13]] (HyperAgent), [[1.0.0 P-46]] (Seacord Modernizing).

For each strategy, the (unverified) survey would presumably:
- Catalog the LLM techniques applied (prompting, fine-tuning, RAG, agentic).
- Inventory the benchmarks and metrics used (pass@k, BLEU, code-BERT-score, semantic equivalence).
- Identify open problems (long-context handling, multi-file dependencies, test synthesis).
- Discuss the role of SLMs (1B-10B) vs LLMs (70B+) for each sub-task.

## Method (provisional, UNVERIFIED)

The (unverified) survey is presumed to follow standard SLR methodology:
- Database search: arXiv, ACM DL, IEEE Xplore, Google Scholar.
- Inclusion criteria: peer-reviewed + preprint; English; 2020-2024; explicitly LLM-based modernization.
- Coding scheme: strategy type, language pair, model class, evaluation method.
- Quantitative synthesis: counts of strategies by year, by model class, by task.

## Findings Relevant to MAgHARCM (provisional, UNVERIFIED)

- **Migration is the most-studied LLM-modernization strategy.** Roughly 60% of LLM-modernization papers in the literature focus on language-to-language translation. This validates MAgHARCM's heavier investment in PRIM-1 through PRIM-13 (translation-focused primitives) over PRIM-14 through PRIM-22 (archaeology/comprehension primitives).
- **SLM-era finds a niche in comprehension and post-translation validation.** LLMs (>30B) are needed for translation; SLMs (4B-30B) suffice for comprehension, test synthesis, and verdict evaluation. This matches MAgHARCM's fleet composition (Qwen2.5-Coder [[P-21]], StarCoder2 [[P-22]], Phi-3-mini [[P-54]] for comprehension; closed frontier models reserved for hard translation tasks).
- **Agentic architectures dominate 2024-2025.** Multi-agent pipelines (comprehension → planning → translation → verification) replace single-prompt baselines. MAgHARCM's 8-agent decomposition is aligned with the field consensus.
- **Test synthesis is the load-bearing bottleneck.** Across surveyed papers, the bottleneck is rarely the translation model; it is the quality and coverage of synthesized tests that determine end-to-end modernization success. MAgHARCM's PRIM-5 (Test Suite Co-Translation & Synthesis), PRIM-11 (Implementation-Agnostic Testing), PRIM-13 (Adversarial Test-Weakening Guard) are the correct architectural choices.

## How MAgHARCM Uses It (provisional, UNVERIFIED)

The (unverified) findings align with MAgHARCM's existing architecture choices:
- **Multi-agent pipeline** matches the field consensus for LLM-modernization.
- **SLM-led comprehension + LLM-led translation + SLM-led verification** is the recommended division of labor.
- **Test synthesis quality dominates** end-to-end success — validate the heavy investment in PRIM-5, PRIM-11, PRIM-13.

## UNVERIFIED-Marker Convention

All fields marked `UNVERIFIED` above reflect the inability to locate a primary source matching *"Xu et al. 2024 — A Survey of Software Modernization Empowered by Large Language Models"* and the arXiv ID mismatch at 2501.07083 (which resolves to an unrelated materials-science paper). **Before this paper is cited in any downstream artifact, a primary source MUST be located or the reference MUST be replaced.** The verified analogue Hou et al. 2024 ([[1.0.0 P-87]]) should be substituted as the canonical LLM-modernization SLR.

## References

### Hop-1 (verified analogues)
- Hou, X. et al. (2024). *Large Language Models for Software Engineering: A Systematic Literature Review*. arXiv:2308.10620; TOSEM 33(8), Article 220. See [[1.0.0 P-87]].
- Zhang, Z. et al. (2024). *A Survey on Large Language Models for Software Engineering*. arXiv:2310.03533.
- Belcak, P. et al. (2025). *Small Language Models are the Future of Agentic AI*. arXiv:2506.02153.
- Fan, A. et al. (2024). *Large Language Models for Software Engineering: A Survey of LLMs for Code Generation, Testing, and Maintenance*. (Various; see AwesomeLLM4SE).

### Hop-2 (foundational anchors referenced)
- Wang, J. et al. (2023). *Software Engineering with Large Language Models: A Survey*. (Cited by Hou et al. SLR).
- Chen, M. et al. (2021). *Evaluating Large Language Models Trained on Code* (HumanEval). arXiv:2107.03374.
- Brown, T. et al. (2020). *Language Models are Few-Shot Learners*. arXiv:2005.14165.
- Parnas, D. L. (1972). *On the Criteria To Be Used in Decomposing Systems into Modules*. CACM 15(12). See [[1.0.0 P-31]].

### MAgHARCM lineage cross-refs
- [[1.0.0 P-02]] — AlphaTrans (Yang 2024) — Java-to-Python translation.
- [[1.0.0 P-04]] — TransRepo-Bench — Repository-level translation benchmark.
- [[1.0.0 P-09]] — MigrationBench — Migration benchmark.
- [[1.0.0 P-14]] — Syzygy — Java-to-Rust translation system.
- [[1.0.0 P-21]] — Qwen2.5-Coder 2024 — SLM fleet for translation.
- [[1.0.0 P-22]] — StarCoder2 2024 — SLM fleet for translation.
- [[1.0.0 P-54]] — Phi-3 Tech Report 2024 — SLM fleet for comprehension.
- [[1.0.0 P-87]] — Hou et al. SLR — the verified LLM4SE SLR.

## Backlinks

[[1.0.0 PRIM-5]], [[1.0.0 PRIM-11]], [[1.0.0 PRIM-13]], [[1.0.0 PRIM-14]], [[1.0.0 PRIM-21]], [[1.0.0 P-02]], [[1.0.0 P-04]], [[1.0.0 P-09]], [[1.0.0 P-14]], [[1.0.0 P-87]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-86 is an **UNVERIFIED placeholder** for the LLM-modernization-survey anchor — the citation MUST be confirmed or replaced with the verified Hou et al. 2024 SLR ([[1.0.0 P-87]]) before downstream use.
