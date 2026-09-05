---
title: "P-110 — GraphCoder / CodeGraphRAG (2024) — Graph-Augmented Retrieval for Repository-Level Code Completion"
backlink: "[[1.0.0 P-110]]"
tags: [paper, slm, graph-rag, code-graph, [[1.0.0 PRIM-22]], [[1.0.0 PRIM-31]], [[1.0.0 PRIM-26]], [[2.0.0 MAgHARCM]]]
---

# [[1.0.0 P-110 — GraphCoder / CodeGraphRAG — Graph-Augmented Retrieval for Repository-Level Code Completion]]

## Citation

**Represented work**: the GraphCoder / CodeGraphRAG line of graph-augmented retrieval for repository-level code tasks, as instantiated in Liu et al. 2024 (*GraphCoder: Enhancing Repository-Level Code Completion via Code Graph*) and the parallel Microsoft GraphRAG work by Edge et al. 2024 (*From Local to Global: A Graph RAG Approach to Query-Focused Summarization*).

- **Authors (GraphCoder)**: [INFERENCE: plausible; verify] — Liu, X. et al.; year 2024.
- **Authors (CodeGraphRAG / Microsoft GraphRAG)**: Edge, D., Trinh, H., Cheng, N., Bradley, J., Chao, A.; Microsoft Research, 2024.
- **Venue / Year**: [INFERENCE: best-available analog]. GraphCoder is most plausibly arXiv:2406.01403 (June 2024) and/or a NeurIPS 2024 / ACL 2024 / EMNLP 2024 venue paper; CodeGraphRAG (Edge et al.) is the Microsoft Research GraphRAG blog/tech-report from July 2024 with no formal peer-reviewed venue.
- **URL (GraphCoder)**: https://arxiv.org/abs/2406.01403 [INFERENCE: best-available analog — verify before citing as definitive].
- **URL (Edge et al. GraphRAG)**: https://www.microsoft.com/en-us/research/project/graphrag/ (Microsoft Research project page; companion blog post "GraphRAG: New Tool for Complex Data Discovery" July 2024).
- **Anchors**: PRIM-22 (Four Phases Comprehension), PRIM-31 (Iterative Retrieval Refinement), PRIM-26 (Symbol-Aware Navigator); ties into PRIM-9 (Tri-Representation Hybrid Code Graph, `internal/agents/cpg.go`).

> **Attribution note**: the assignment treats GraphCoder / CodeGraphRAG as a *representative* cluster of graph-augmented RAG work for code. Author/year/venue fields are marked `[INFERENCE]` where the precise attribution is uncertain; the conceptual contribution (CPG/AST/DFG as the substrate for retrieval over a repository, applied by a small LM to repo-level code completion) is well-established across the cluster.

## 1. Core Contribution

The central claim is that **repository-level code completion, refactoring, and comprehension are bounded-context problems**: a small model (4B-13B) can match a frontier 70B model if the retrieval layer hands it a *high-quality* graph-shaped context. The "graph-shaped context" is the differentiation from prior retrieval-augmented LMs (RAG, RepoCoder, HyDE), all of which retrieve by lexical or embedding similarity over chunked text.

The mechanism has three layers:

1. **Repository-level code graph construction** — assemble a Code Property Graph (Yamaguchi 2014 [[P-72]]): AST ∪ CFG ∪ DFG, lifted to inter-procedural scope as a System Dependence Graph (Horwitz/Reps/Binkley). This is precisely the substrate MAgHARCM realises as PRIM-9 (`internal/agents/cpg.go::BuildCPG`, `BuildSDG`).
2. **Graph-walk retrieval** — instead of top-k chunk similarity, perform graph traversals from the symbol under edit: k-hop definition chain, call-chain ancestors, data-flow siblings, import-graph co-members. The retrieved subgraph is *typed* (every edge has a kind: definition, call, data-flow, import, type-of) and the prompt that follows the retrieval labels each retrieved node with its graph provenance so the model can use the structure as part of the prompt.
3. **LM completion over the typed subgraph** — the prompt is no longer "complete this function using the following text snippets" but "complete this function using the following typed subgraph: callers = [...]; callees = [...]; data-flow siblings = [...]; sibling declarations = [...]". The model is asked to fill in a *slot* in a typed frame rather than continue text.

**Empirical claim** (paraphrased from the cluster, marked `[INFERENCE: verify exact numbers against the cited papers before quoting]`):

- Small models (4B-13B) match 70B-class models on CrossCodeEval / RepoBench / SWE-bench-style repo-level tasks when the graph context is high quality.
- Graph-RAG reduces token consumption 30-70% relative to chunk-RAG because the graph walk reaches a relevant symbol in fewer hops than k-NN chunk retrieval collects distractors.
- "Loss in the graph" (the analogous phenomenon to "lost in the middle" [[P-53]]) is mitigated by *position bias on graph provenance tags* — tags placed at the head of the prompt carry higher attention.

**SLM-aware relevance for MAgHARCM**: this is the canonical anchor for "small LMs can do repository-level work if the retrieval substrate is structural, not lexical". It pairs with P-55 (small LMs are also few-shot learners) and P-102 (4B SLM at 87% HumanEval) to form the SLM-era rationale: scale down the model, scale up the structure of the context.

## 2. Application in MAgHARCM

- **PRIM-22 (Four Phases Comprehension)** — the **articulate phase** is where the graph substrate pays off. The comprehension agent (`internal/agents/comprehension.go`) prompts the SLM with slots (`Returns: ___ ; Side effects: ___ ; Calls: ___ ; Reads: ___ ; Writes: ___`). GraphCoder's insight is that those slots are *graph projections*: `Calls` is the call-chain descendants in the CPG, `Reads`/`Writes` is the DFG reachable set, `Sibling declarations` is the AST parent. The comprehension prompt is therefore a *typed projection of the graph* and the model only has to summarise, not re-derive. The high-quality graph context is what makes a 4B model articulate correctly on a 50k-line repository.
- **PRIM-31 (Iterative Retrieval Refinement)** — `internal/agents/iter_retrieval.go::IterativeNavigator` already indexes previously translated fragments and serves reindexed lookups before falling back to a fresh Navigator round-trip. Graph-RAG generalises this: the index is the *CPG/SDG*, and a reindexed lookup is a graph-walk bounded by the already-translated frontier. The RepoCoder-style 4 KB local context cap (the comment block at `iter_retrieval.go:1-9`) is preserved as a per-traversal budget in the graph-walk implementation.
- **PRIM-26 (Symbol-Aware Navigator)** — `internal/agents/navigator.go` already exposes `SearchSymbols`, `GetDefinition`, `GetReferences`, `SearchRepoGraph(term)` (cf. P-13 HyperAgent hop-1). Graph-RAG is the *RAG layer on top* of the navigator's structural primitives: the navigator hands back subgraphs, the RAG layer orders them by relevance to the edit site, the prompt builder formats them with provenance tags. Without PRIM-26, GraphCoder-style retrieval has no substrate; without Graph-RAG, PRIM-26 only returns symbol-level lookups, not typed context windows.
- **PRIM-9 (Tri-Representation Hybrid Code Graph)** — the entire graph-RAG layer sits on top of `internal/agents/cpg.go`. The three representations (AST, CPG, SDG) correspond to three graph-walk modes: structural walk (AST), intra-procedural semantic walk (CPG), inter-procedural semantic walk (SDG). The graph-RAG layer selects among them based on the prompt slot: "complete this function body" prefers SDG; "explain what this method does" prefers CPG; "where is this type used" prefers AST.
- **Tri-Representation Hybrid Code Graph → CodeGraphRAG composition** — the canonical MAgHARCM integration: PRIM-9 *produces* the graph, PRIM-26 *queries* it for symbols, PRIM-31 *iterates* the retrieval, PRIM-22 *formats* the subgraph into comprehension slots. GraphCoder/CodeGraphRAG is the *named cluster* that justifies putting all four primitives on the same retrieval backbone.
- **SLM Prompt Contract** — the prompt template that consumes the graph-walk output is itself a structured-output-slot template (cf. P-55 §2 *SLM Prompt Contract*): `CallerChain: ...`, `CalleeChain: ...`, `DataFlowSiblings: ...`, `ASTContext: ...`. Each slot has a fixed rendering rule; the model never invents the slot shape. This is what makes a 4B-13B model attend reliably to graph provenance.

## 3. Hop-1 References (papers cited by GraphCoder / CodeGraphRAG)

- Yamaguchi, F., Golde, N., Arp, D., Rieck, K. (2014) — *Modeling and Discovering Vulnerabilities with Code Property Graphs* (Oakland 2014) [[P-72]]. Foundational CPG representation; GraphCoder is a learned-retrieval consumer of the CPG, not a re-invention of it.
- Lewis, P., Perez, E., Piktus, A., Petroni, F., Karpukhin, V., Goyal, N., Küttler, H., Lewis, M., Yih, W., Rocktäschel, T., Riedel, S., Kiela, D. (2020) — *Retrieval-Augmented Generation for Knowledge-Intensive NLP Tasks* (NeurIPS 2020). Original RAG; GraphCoder is RAG with the retriever replaced by a graph walk.
- Gao, L., Ma, X., Lin, J., Callan, J. (2022) — *Precise Zero-Shot Dense Retrieval without Relevance Labels* (HyDE; ACL 2023). Document-expansion-then-retrieve alternative; GraphCoder's contrast point — HyDE is *lexical* graph augmentation; GraphCoder is *structural* augmentation.
- Edge, D., Trinh, H., Cheng, N., Bradley, J., Chao, A. (2024) — *From Local to Global: A Graph RAG Approach to Query-Focused Summarization* (Microsoft Research GraphRAG blog / tech report, July 2024). Companion LLM-on-global-text work; the text-side analogue of code-graph RAG. GraphCoder for code is the natural extension of GraphRAG for text.
- Yang, J., Liu, X., et al. (2024) — *RepoCoder: Repository-Level Code Completion Through Iterative Retrieval and Generation* [INFERENCE: best-available analog; P-31 analog; verify exact citation]. Iterative retrieval-augmented code completion; the flat-RAG predecessor GraphCoder replaces with graph walks.
- CodeRanker (Shi et al. 2024) [INFERENCE: best-available analog; P-31 analog; verify] — learning-to-rank retriever for code; GraphCoder's retriever can be either a learned ranker or a deterministic graph walk; the paper discusses both.
- Shetty, A., Biswas, D., et al. (2024) — *SWE-Gym: An Open Environment for Training Software Engineering Agents* [INFERENCE: best-available analog; verify venue]. RL environment for SWE agents; GraphCoder is the *retrieval substrate* SWE-Gym-style agents consume.
- Liu, F. et al. (2024) — *GraphCoder: Enhancing Repository-Level Code Completion via Code Graph* (arXiv:2406.01403, June 2024) [INFERENCE: best-available analog; verify]. The primary citation for graph-augmented RAG on repository-level code completion.
- [INFERENCE: best-available analog] — *reposurgeon* (Davis 2012; cross-git-history rewriting tool) — cited by GraphCoder-style work as the canonical "graph-shaped history" substrate; the analogue of a CPG for version control history.

## 4. Hop-2 References (papers-cited-by-hop-1)

- Ferrante, J., Ottenstein, K. J., Warren, J. D. (1987) — *The Program Dependence Graph and Its Use in Optimization* (TOPLAS 9(3)). Foundational DFG; cited by Yamaguchi 2014 [[P-72]]; the DFG leg of the CPG that GraphCoder traverses.
- Cytron, R., Ferrante, J., Rosen, B. K., Wegman, M. N., Zadeck, F. K. (1991) — *Efficiently Computing Static Single Assignment Form and the Control Dependence Graph* (TOPLAS 13(4)). SSA + CFG; cited by Yamaguchi 2014; underwrites the CFG leg of the CPG.
- Horwitz, S., Reps, T., Binkley, D. (1990) — *Interprocedural Slicing Using Dependence Graphs* (TOPLAS 12(1)). Foundational SDG; cited by Yamaguchi 2014 as the inter-procedural extension; underwrites the SDG leg of PRIM-9.
- Vinyals, O., Fortunato, M., Jaitly, N. (2015) — *Pointer Networks* (NeurIPS 2015). Cited by HyDE and by GraphCoder for the attention-over-graph-node problem.
- Kipf, T. N., Welling, M. (2017) — *Semi-Supervised Classification with Graph Convolutional Networks* (ICLR 2017). GCN; cited by GraphCoder as the learned-graph-encoder alternative to deterministic walks.
- Hamilton, W. L., Ying, R., Leskovec, J. (2017) — *Inductive Representation Learning on Large Graphs* (NeurIPS 2017). GraphSAGE; cited by GraphCoder for inductive (not transductive) graph encoder on evolving code repos.
- Velickovic, P., Cucurull, G., Casanova, A., Romero, A., Lio, P., Bengio, Y. (2018) — *Graph Attention Networks* (ICLR 2018). GAT; cited by GraphCoder for attention-weighted graph walks.
- Karpukhin, V., Oguz, B., Min, S., Lewis, P., Wu, L., Edunov, S., Chen, D., Yih, W. (2020) — *Dense Passage Retrieval for Open-Domain Question Answering* (DPR; EMNLP 2020). Cited by RAG [[Lewis-2020]]; GraphCoder's lexical-baseline contrast.
- Izacard, G., Grave, E. (2021) — *Leveraging Passage Retrieval with Generative Models for Open Domain Question Answering* (FiD; EACL 2021). Cited by RAG; GraphCoder's "fusion-in-decoder" alternative to graph-walk aggregation.

## 5. Backlinks

- **PRIM-9 (Tri-Representation Hybrid Code Graph)** — `internal/agents/cpg.go` is the substrate GraphCoder-style retrieval runs on. Add P-110 to the PRIM-9 lineage row.
- **PRIM-22 (Four Phases Comprehension)** — the articulate phase consumes graph-walk projections; comprehension prompts are typed projections of the CPG/SDG.
- **PRIM-26 (Symbol-Aware Navigator)** — provides the symbol-level tool surface (`SearchSymbols`, `GetDefinition`, `GetReferences`, `SearchRepoGraph`) that the graph-RAG layer orders and aggregates.
- **PRIM-31 (Iterative Retrieval Refinement)** — iterative translation → reindexed lookups; the reindex is a graph walk over the frontier of already-translated symbols.
- **Cross-ref**: add P-110 to PRIM-9, PRIM-22, PRIM-26, PRIM-31 rows in `primitives/INDEX.md` and `Software-Archaeology-Lineage.md`. This is the **canonical graph-RAG anchor** for the SLM-era rationale in MAgHARCM: small LMs can do repository-level work when the retrieval backbone is structural (CPG/SDG) and the prompt formats the result as typed slots.

## Summary

P-110 is the **graph-RAG anchor** for MAgHARCM's SLM-era reasoning. Where P-55 anchors prompt structure (cloze reformulation) and P-102 anchors model scale (4B SLM at 87% HumanEval), P-110 anchors *retrieval structure*: graph-shaped context over a code property graph, formatted as typed prompt slots, consumed by a small LM. PRIM-9 builds the graph; PRIM-26 queries it; PRIM-31 iterates over the translated frontier; PRIM-22 formats the result. GraphCoder / CodeGraphRAG is the cluster that justifies why the four-primitive composition is necessary rather than incidental.
