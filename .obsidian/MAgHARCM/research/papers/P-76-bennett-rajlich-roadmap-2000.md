---
title: "P-76 — Bennett & Rajlich 2000 — Software Maintenance and Evolution: A Roadmap"
backlink: "[[1.0.0 P-76]]"
tags: [paper, software-maintenance, roadmap, lifecycle, research-agenda, [[1.0.0 P-74]], [[1.0.0 P-75]], hop-2]
---

# [[1.0.0 P-76 — Bennett & Rajlich Roadmap]]

## Citation

Bennett, K. H., & Rajlich, V. T. (2000). *Software Maintenance and Evolution: A Roadmap*. Proceedings of the Conference on the Future of Software Engineering (FOSE 2000) at ICSE 2000, pp. 73-87. ACM. DOI: 10.1145/336512.336534. URL: https://dl.acm.org/doi/10.1145/336512.336534.

## Summary

Bennett & Rajlich's 2000 Roadmap paper (FOSE 2000, a special track at ICSE 2000) is the field-defining **research-agenda statement** for software maintenance and evolution in the 21st century. The paper is structured as:
1. **Definition**: distinguishes *software maintenance* (corrective, adaptive, perfective, preventive) from *software evolution* (the broader category that includes re-engineering, migration, and retirement).
2. **State of the Art**: reviews the 1990s maintenance research (Biggerstaff's concept assignment, Lehman's laws, Chikofsky-Cross reverse-engineering taxonomy, Seacord's modernisation textbook, Ward's maintenance-process models).
3. **Research Challenges**: lists seven grand challenges:
   - Concept location and feature identification.
   - Impact analysis of changes.
   - Program comprehension at scale.
   - Software visualisation for maintenance.
   - Reuse-driven maintenance.
   - Migration strategies.
   - Empirical validation of maintenance techniques.
4. **Roadmap**: argues for empirical validation; predicts the rise of open-source software as a maintenance laboratory; flags the integration of reverse-engineering tools as the next major research direction.

## Method

The paper is a **position / roadmap** paper, not an empirical study. Its contribution is the synthesised research agenda and the empirical-vs-analytical framing that subsequent research has followed.

## Findings Relevant to MAgHARCM

- **Seven grand challenges** map directly onto MAgHARCM primitives: concept location → `[[1.0.0 PRIM-20]]`; impact analysis → `[[1.0.0 PRIM-18]]`; comprehension at scale → `[[1.0.0 PRIM-22]]`; visualisation → TUI; reuse-driven maintenance → domain-adapter training (cf. `[[1.0.0 P-58]]` Qwen2.5-Coder); migration strategies → `[[1.0.0 PRIM-21]]`; empirical validation → the experiment framework.
- **Empirical validation requirement** is the precedent for MAgHARCM's "honest experiments" directive — every primitive must be validated on a real legacy codebase, not just a synthetic example.
- **Open-source as maintenance laboratory** is the precedent for MAgHARCM's choice of case-study codebases (the open-source repositories used in `cmd/magh/bench`).
- **Tool integration** is the precedent for MAgHARCM's ABCoder MCP integration (cf. `[[1.0.0 P-72]]`).

## How MAgHARCM Uses It

The roadmap paper's seven challenges are referenced in `[[2.0.0 Software-Archaeology-Lineage]]` §3 as the **research-agenda table**. Each MAgHARCM primitive's "Application" column points back to one of the seven challenges.

## References

### Hop-1 (Bennett & Rajlich 2000 cites)
- Biggerstaff, T. J., Mitbander, B. G., & Webster, D. E. (1993). *The Concept Assignment Problem in Program Understanding*. ICSE 1993. See [[1.0.0 P-74]].
- Belady, L. A. & Lehman, M. M. (1976). *A Model of Large Program Development*. IBM Systems Journal 15(3):225-252. See [[1.0.0 P-32]].
- Chikofsky, E. J. & Cross, J. H. (1990). *Reverse Engineering and Design Recovery: A Taxonomy*. IEEE Software 7(1):13-17. See [[1.0.0 P-33]].
- Seacord, R. C. et al. (2003). *Modernizing Legacy Systems*. Addison-Wesley. See [[1.0.0 P-46]].
- Bennett, K. H. (1995). *Legacy Systems: Coping with Success*. IEEE Software 12(1):19-23. See [[1.0.0 P-75]].

### Hop-2
- Mens, T. & Tourwé, T. (2004). *A Survey of Software Refactoring*. IEEE TSE 30(2):126-139.
- Zimmermann, T. et al. (2005). *Mining Version Archives to Co-evolve Program and Regression Test Scripts*. MSR 2005.
- Hassan, A. E. & Xie, T. (2010). *Software Innovation: A Systematic Literature Review*. TR.

## Backlinks

[[1.0.0 P-32]], [[1.0.0 P-33]], [[1.0.0 P-46]], [[1.0.0 P-74]], [[1.0.0 P-75]], [[1.0.0 PRIM-18]], [[1.0.0 PRIM-20]], [[1.0.0 PRIM-21]], [[1.0.0 PRIM-22]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-76 is the **research-agenda anchor** for MAgHARCM's primitive roadmap — each primitive maps to one of the seven Bennett-Rajlich challenges.
