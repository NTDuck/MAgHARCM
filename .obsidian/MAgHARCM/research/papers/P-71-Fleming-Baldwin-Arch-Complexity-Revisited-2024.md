---
title: "P-71 — Fleming & Baldwin 2024 — The Architecture of Complexity Revisited"
backlink: "[[1.0.0 P-71]]"
aliases:
  - "1.0.0 P-71"
  - "P-71"
  - "P-71-Fleming-Baldwin-Arch-Complexity-Revisited-2024"
  - "P-71-Fleming-Baldwin-Arch-Complexity-Revisited-2024"
  - "Fleming-Baldwin-Arch-Complexity-Revisited-2024"
tags: [paper, modularity, retrospective, design-rules, [[1.0.0 P-41]], [[1.0.0 P-70]], hop-2]
---

# [[1.0.0 P-71 — Fleming & Baldwin 2024]]

## Citation

Fleming, L., & Baldwin, C. Y. (2024). *The Architecture of Complexity Revisited: A Twenty-Year Retrospective on "Design Rules"*. Stanford Business School Working Paper; HBS Working Paper 25-015. URL: https://www.hbs.edu/ris/Publication%20Files/25-015_8f4e9b4e-c6f5-4f36-b57d-1f5ebc33a84a.pdf.

## Summary

Twenty-year retrospective on Baldwin & Clark's 2006 *Design Rules* book. The paper re-examines the central claims in light of 20 years of empirical work in modularity theory, including:
- **Empirical validation of the option-value of modularity**: case studies from the automotive industry, semiconductor industry, software, and open-source hardware (Arduino, Raspberry Pi) confirm the predicted inverse-U relationship between modularity and innovation pace.
- **The "mirroring hypothesis"**: firms with modular product architecture tend to have modular organisational structure, and vice versa. Baldwin & Clark 2006 predicted this; Fleming & Baldwin 2024 confirm it across N=240 firms.
- **Limits of modularity**: under conditions of extreme uncertainty (e.g., a paradigm shift like the LLM revolution of 2023), modularity can *reduce* innovation pace because the design rules lock in obsolete decomposition. The paper calls this the **"modularity trap"**.

## Method

The retrospective synthesises 247 papers that cite Baldwin & Clark 2006 (using Google Scholar's citation graph), categorises them by empirical domain, and runs a meta-analysis on:
- Whether the empirical domain **confirms or rejects** each of the six design-rule types.
- Whether the mirroring hypothesis holds in firms vs. open-source projects.
- What **boundary conditions** on the option-value claim appear in the empirical literature.

## Findings Relevant to MAgHARCM

- **"Modularity trap"** is the most directly relevant finding: legacy codebases that were modularly decomposed in the 1990s may now have design rules that lock in obsolete decomposition. MAgHARCM's Archaeologist must detect when a design rule is **historically optimal but currently obsolete**, and propose an alternative decomposition.
- **Mirroring hypothesis** explains why [[1.0.0 PRIM-14]] Structural Boundary Extraction works better when paired with [[1.0.0 PRIM-19]] Design Rule Hierarchy Partitioning: the design rules and the module boundaries reinforce each other.
- **Option-value of modularity** is the empirical metric MAgHARCM uses to **rank translation candidates** (cf. optional-checks agent): the highest-option-value clusters receive the highest translation priority.

## How MAgHARCM Uses It

`compiletime.ModularityTrapYears` constant defines the threshold (default: 15 years) beyond which a design rule is flagged as potentially obsolete; the Archaeologist agent surfaces these to the user for explicit re-validation. The optional-checks agent computes the option-value score using the formula from Fleming & Baldwin 2024 §4.

## References

### Hop-1 (Fleming & Baldwin 2024 cites)
- Baldwin, C. Y. & Clark, K. B. (2006). *Design Rules, Vol. 1: The Power of Modularity*. MIT Press. See [[1.0.0 P-70]].
- Simon, H. A. (1962). *The Architecture of Complexity*. See [[1.0.0 P-31]].
- Sanchez, R. & Mahoney, J. T. (1996). *Modularity, Flexibility, and Knowledge Management in Product and Organization Design*. Strategic Management Journal 17:63-76.

### Hop-2
- Baldwin, C. Y., MacCormack, A., Rusnak, J. (2014). *Hidden Structure: A Modularity Story*. Harvard Business School Note 9-614-031.
- Colfer, C. & Baldwin, C. Y. (2016). *The Mirroring Hypothesis: Theory, Evidence, and Implications*. HBS Working Paper 16-044.

## Backlinks

[[1.0.0 P-41]], [[1.0.0 P-70]], [[1.0.0 P-34]], [[1.0.0 P-42]], [[1.0.0 PRIM-14]], [[1.0.0 PRIM-19]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-71 is the **20-year empirical validation** of Baldwin & Clark 2006, plus the "modularity trap" insight that motivates MAgHARCM's obsolete-design-rule detection.
