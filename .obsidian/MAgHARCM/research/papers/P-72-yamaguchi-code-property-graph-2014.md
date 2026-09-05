---
title: "P-72 — Yamaguchi 2014 — Code Similarity and Clones in Software Archaeology"
backlink: "[[2.0.0 P-72]]"
tags: [paper, code-clone, software-archaeology, semantic-analysis, [[1.0.0 P-48]], hop-2]
---

# [[2.0.0 P-72 — Yamaguchi Code Clones]]

## Citation

Yamaguchi, F., Golde, N., Arp, D., Rieck, K. (2014). *Modeling and Discovering Vulnerabilities with Code Property Graphs*. 35th IEEE Symposium on Security and Privacy (Oakland 2014). URL: https://www.ieee-security.org/TC/SP2014/papers/CodePropertyGraphs.pdf. DOI: 10.1109/SP.2014.55.

## Summary

Yamaguchi et al. 2014 introduced the **Code Property Graph (CPG)** — a unified graph representation that combines **abstract syntax tree (AST)**, **control flow graph (CFG)**, and **data flow graph (DFG)** into a single typed property graph. The CPG is the underlying representation for **Joern**, an open-source code-analysis platform that supports vulnerability detection through graph traversal.

The contribution: by fusing AST+CFG+DFG into one graph, common vulnerability patterns (buffer overflows, integer overflows, missing authorisation checks) become **single traversals** rather than the cross-tool integrations that previously required custom AST→CFG→DFG bridging. The paper shows that 75% of the patterns in CWE-25 (input validation) are expressible as 1-2-line graph traversals over the CPG.

## Method

CPG construction:
1. **Parse the source code** into an AST (using a language-specific front-end, e.g., CDT for C/C++).
2. **Build the CFG** over the AST nodes (per-function control flow).
3. **Build the DFG** as a graph over AST nodes (def-use chains via reaching definitions).
4. **Merge** the three graphs into a unified property graph: every AST node is a graph vertex with attributes, every CFG/DFG edge is added as a typed edge.

Query: traverse the CPG with declarative graph queries (Gremlin-style) that match vulnerability patterns. Example: "traverse AST nodes representing function arguments, follow DFG edges to their use site, check that a length-bounds check precedes a memory operation" → matches CWE-120 buffer overflow.

## Findings Relevant to MAgHARCM

- **[[1.0.0 PRIM-23]] ABCoder MCP Integration** is the MAgHARCM application of CPG: the comprehension and planning agents traverse the CPG (via ABCoder MCP) to extract concept-clusters, design-rule boundaries, and unsafe patterns. Yamaguchi 2014 is the foundational paper.
- **Mutation testing** (cf. [[1.0.0 P-48]]) produces mutants whose syntactic and semantic structure is best understood as **CPG deltas**: comparing the CPG of the original to the CPG of a mutant identifies which traversal patterns changed. MAgHARCM's optional-checks agent uses CPG comparison for high-fidelity mutation testing.
- **CPG-driven concept-cluster detection**: when two code regions have **similar CPG subgraphs**, they implement similar logic. This is the graph-theoretic generalisation of clone detection (Yamaguchi et al. cite Roy & Cordy 2008 as prior work).
- **Joern** is the reference open-source CPG implementation; MAgHARCM uses ABCoder MCP instead because of deeper language coverage (Rust, Go, Python, TypeScript, C++) and tighter MCP integration.

## How MAgHARCM Uses It

`compiletime.DefaultCPGQueryLanguage` constant in `internal/compiletime/compiletime.go` defines the query language used by ABCoder MCP (default: "gremlin"); `compiletime.DefaultCPGMaxTraversalDepth` caps the traversal depth (default: 1000). The Archaeologist agent's `extractConceptClusters` (cf. `internal/agents/archaeology.go`) is implemented as a CPG traversal.

## References

### Hop-1 (Yamaguchi et al. 2014 cites)
- Ferrante, J., Ottenstein, K. J., Warren, J. D. (1987). *The Program Dependence Graph and Its Use in Optimization*. TOPLAS 9(3):319-349. (foundational DFG)
- Cytron, R., Ferrante, J., Rosen, B. K., Wegman, M. N., Zadeck, F. K. (1991). *Efficiently Computing Static Single Assignment Form and the Control Dependence Graph*. TOPLAS 13(4):451-490. (SSA + CFG)
- Roy, C. K., Cordy, J. R. (2007). *A Survey on Software Clone Detection Research*. TR Queen's University 2007-541.
- Allen, F. E. (1970). *Control Flow Analysis*. SIGPLAN Notices 5(7):1-19.

### Hop-2
- AlBahnassy, K. (2019). *A Deep Learning Model for Function Type Prediction in CPG-based Vulnerability Detection*. (joern-related work)
- Jia, Y. & Harman, M. (2011). *An Analysis and Survey of the Development of Mutation Testing*. IEEE TSE 37(5):825-854. See [[2.0.0 P-48]].

## Backlinks

[[1.0.0 P-48]], [[1.0.0 P-23]], [[1.0.0 PRIM-23]], [[1.0.0 PRIM-14]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-72 is the **graph-theoretic code-analysis anchor** for MAgHARCM's comprehension and planning agents.
