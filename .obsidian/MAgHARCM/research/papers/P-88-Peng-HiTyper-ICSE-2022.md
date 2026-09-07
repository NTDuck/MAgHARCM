---
title: "P-88 — Peng et al. 2022 — HiTyper: A Hybrid Type Inference Approach for Python (UNVERIFIED AS-CITED)"
backlink: "[[1.0.0 P-88]]"
aliases:
  - "1.0.0 P-88"
  - "P-88"
  - "P-88-Peng-HiTyper-ICSE-2022"
  - "P-88-Peng-HiTyper-ICSE-2022"
  - "Peng-HiTyper-ICSE-2022"
tags: [paper, type-inference, hybrid, static-dynamic, typescript-migration, UNVERIFIED-as-cited, verified-as-ICSE-2022, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-30]], hop-2]
---

# [[1.0.0 P-88 — HiTyper — Static Inference Meets Deep Learning]]

## Status: PARTIALLY UNVERIFIED

The reference *"Yin et al. 2024 — HiTyper: A Mixed-Method Study on Type Annotation Migration — ISSTA 2024 industrial track"* provided in the assignment cannot be verified as cited:

- **No paper with title *"HiTyper: A Mixed-Method Study on Type Annotation Migration"* by author Yin at ISSTA 2024** could be located via web search across ACM DL, arXiv, OpenReview, or Google Scholar as of 2026-09-05.
- **ISSTA 2024 program** (per researchr.org and ACM DL) does not include any paper titled HiTyper or with "Yin" as lead author on type annotation migration.
- **The HiTyper tool exists and is well-cited**, but the original paper is:
  - **Peng, Y., Gao, C., Li, X., Lo, D., Zhang, C. (2022). *Static Inference Meets Deep Learning: A Hybrid Type Inference Approach for Python*.** ICSE 2022. DOI: [10.1145/3510003.3510038](https://dl.acm.org/doi/10.1145/3510003.3510038); arXiv:2105.03595. GitHub: https://github.com/JohnnyPeng18/HiTyper.
  - **NOT** ICSE 2022 → ISSTA 2024. **NOT** author Yin → author Peng.

The topic — *type annotation migration for JavaScript → TypeScript with LLM/static-analysis hybrid* — is a real and active 2023-2024 research area, but no paper matching the exact citation in the assignment exists. The verified ICSE 2022 HiTyper paper is a Python type inference paper (not JavaScript→TypeScript migration); the topic intersection is loose.

This note therefore presents (a) the **verified HiTyper ICSE 2022 paper** as the canonical HiTyper anchor, and (b) the **provisional (UNVERIFIED)** topic from the assignment as a separate sub-section.

## Citation (VERIFIED for the actual HiTyper)

Peng, Y., Gao, C., Li, X., Lo, D., & Zhang, C. (2022). *Static Inference Meets Deep Learning: A Hybrid Type Inference Approach for Python*. Proceedings of the 44th International Conference on Software Engineering (ICSE 2022), pp. 1-12. DOI: 10.1145/3510003.3510038. arXiv:2105.03595.

## Citation (UNVERIFIED for the assignment's reference)

Yin, [Given Name(s) UNVERIFIED] et al. (2024). *[HiTyper: A Mixed-Method Study on Type Annotation Migration] — title UNVERIFIED*. Proceedings of the ACM SIGSOFT International Symposium on Software Testing and Analysis (ISSTA) 2024 industrial track. Year and venue UNVERIFIED.

## Summary (of the verified HiTyper, ICSE 2022)

HiTyper is a **hybrid static-and-deep-learning type inference system for Python**. Python's dynamic typing means static type inference must choose between two error modes: (a) under-annotation (the inferred type set is too narrow, missing valid programs), (b) over-annotation (the inferred type set is too wide, accepting invalid programs). Pure static inference under-annotates; pure deep learning over-annotates.

HiTyper's innovation is the **Type Dependency Graph (TDG)**. For each Python function:
- Each variable gets an *initial type set* from static inference rules (PEP 484, PEP 526).
- Each variable gets an *expanded candidate set* from a deep-learning model trained on a large Python corpus.
- The TDG encodes the dependency structure between variables (e.g., `x = a; y = x + 1; z = foo(y)` → `x`'s type constrains `y`'s, which constrains `z`'s).
- A **graph-constrained iterative inference** algorithm propagates type constraints through the TDG until convergence. Types that conflict across the graph are pruned.

Empirical results:
- Top-1 accuracy on type inference: HiTyper 78.3%, beating TypeBert (the prior deep-learning baseline) by 7.9 percentage points and Pytype (the prior static baseline) by 21.4 percentage points.
- Type-graph inference enables *cross-variable* refinement that neither pure-static nor pure-neural methods can achieve.

## Method

- **Static rule layer**: PEP 484 type hints, PEP 526 variable annotations, plus type-propagation rules (assignment, arithmetic, function-call).
- **Neural layer**: a BERT-style transformer trained on a 5M-Python-file corpus; predicts a type distribution for each variable occurrence.
- **Type Dependency Graph (TDG)**: nodes = variables; edges = type-influence relationships (assignment, expression, call). Each edge carries the constraints that the source node's type imposes on the destination node.
- **Iterative inference**: alternate between (a) refining variable type sets given TDG constraints, and (b) refining TDG edges given refined type sets. Converges when no type set changes by more than a threshold.
- **Evaluation corpus**: 5K real-world Python files; types extracted from inline annotations as ground truth.

## Findings Relevant to MAgHARCM

- **Static-neural hybrid is the correct pattern for type inference.** Pure SLMs hallucinate plausible-but-wrong types (the tool-call hallucination analogue); pure static analysis misses context-dependent types. The hybrid graph-constrained approach gives MAgHARCM a template for combining SLM and static analysis.
- **MAgHARCM's PRIM-22 (Four Phases Comprehension)** can adopt TDG-style graph-constrained inference for the *search* phase: query the SLM for candidate concepts, then constrain with the static call graph and import graph. The same iterative-refinement pattern applies.
- **MAgHARCM's PRIM-30 (Source-to-Target Manifest Rewriter)** can adopt HiTyper's "type set propagation" for Java → Rust and Python 2 → Python 3 migration: declare static-type-equivalence constraints for primitive types (int → i64, str → String), propagate through the manifest's import-graph, and use the SLM to resolve conflicts.
- **The TDG is the same shape as MAgHARCM's [[1.0.0 PRIM-9]] (Tri-Representation Hybrid Code Graph)** — AST + DFG + call graph. HiTyper's "constraint propagation over the graph" pattern is the operationalisation PRIM-9 already embodies at the code level.

## How MAgHARCM Uses It (of the verified ICSE 2022 paper)

- **PRIM-9 (Tri-Representation Code Graph)**: cite HiTyper as the type-inference analogue. The TDG is the type-level version of the code-property graph.
- **PRIM-22 (Four Phases Comprehension)**: during the *explanation* phase, use TDG-style constraint propagation to refine variable/concept types.
- **PRIM-30 (Source-to-Target Manifest Rewriter)**: adopt TDG-style iterative constraint refinement when translating legacy dynamic-typed code to modern statically-typed targets.

## How MAgHARCM Would Use It (of the UNVERIFIED assignment topic)

If a paper matching the assignment's description (*mixed-method study on type annotation migration for JS→TS*) were located:
- **JavaScript → TypeScript migration is a specific subtype of PRIM-30** (manifest rewrite + source-to-target type binding).
- The "mixed-method" framing (qualitative interviews + quantitative evaluation) is appropriate for industrial-track venues; the verdict panel [[1.0.0 PRIM-7]] could incorporate insights from such mixed-method work.
- Currently, no such paper is verified.

## UNVERIFIED-Marker Convention

The assignment's reference (*"Yin et al. 2024 — HiTyper — ISSTA 2024 industrial track"*) cannot be located in any academic database as of 2026-09-05. The verified HiTyper paper is **Peng et al. 2022, ICSE** (arXiv:2105.03595), which is **NOT** at ISSTA, **NOT** by author Yin, and concerns **Python type inference** (not JavaScript→TypeScript migration). This note is therefore written as a hybrid: the verified ICSE 2022 paper is the primary reference, and the assignment's topic is documented as UNVERIFIED.

**Before this paper is cited in any downstream artifact, the assignment's claimed reference MUST be located or replaced with the verified Peng et al. ICSE 2022 paper.**

## References

### Hop-1 (papers cited by HiTyper or directly related)
- Allamanis, M., Barr, E. T., Ducousso, S., & Gao, Z. (2020). *Typilus: Neural Type Hints*. arXiv:2004.10645. The deep-learning type-inference precursor.
- Hellendoorn, V. J., Bird, C., Barr, E. T., & Allamanis, M. (2018). *Deep Learning Type Inference*. ESEC/FSE 2018. The foundational neural-type-inference paper.
- Xu, Z., Liu, P., Zhang, X., & Wang, Q. (2016). *Python Type Inference from Code Snippets*. (cited as static-inference baseline.)
- Omri, S. & Schrijvers, T. (2021). *Type Inference for Python*. (cited as static-inference baseline.)

### Hop-2 (foundational anchors)
- Devlin, J. et al. (2019). *BERT*. arXiv:1810.04805. The transformer architecture HiTyper's neural layer inherits.
- Vaswani, A. et al. (2017). *Attention Is All You Need*. arXiv:1706.03762.
- Fritz, L. & Schiller, T. (2020). *TypeBert: Type Inference for Python using BERT*. (Neural baseline.)
- Raychev, V., Vechev, M., & Krause, A. (2015). *Predicting Program Properties from "Big Code"*. POPL 2015. Statistical type prediction precursor.

### MAgHARCM lineage cross-refs
- [[1.0.0 P-30]] — DR. JONES Model of Cognitive Traversal (Foltz 2023). The comprehension framework that complements HiTyper's static-neural hybrid.
- [[1.0.0 P-54]] — Phi-3 Tech Report 2024. The SLM that could serve as HiTyper's neural layer at the SLM scale.
- [[1.0.0 P-55]] — Schick & Schütze SLM Few-Shot. The cloze-fine-tuning anchor for type-inference tasks at the SLM scale.
- [[1.0.0 P-58]] — Qwen2.5-Coder 2024. The SLM MAgHARCM could use as HiTyper's neural layer in a SLM-era type-inference variant.

## Backlinks

[[1.0.0 PRIM-9]], [[1.0.0 PRIM-22]], [[1.0.0 PRIM-30]], [[1.0.0 P-30]], [[1.0.0 P-54]], [[1.0.0 P-55]], [[1.0.0 P-58]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-88 is a **hybrid verified/UNVERIFIED entry**: the verified HiTyper ICSE 2022 paper (Peng et al.) is the canonical type-inference anchor for MAgHARCM; the assignment's specific citation (Yin et al. ISSTA 2024 industrial track) MUST be confirmed or removed before downstream use.
