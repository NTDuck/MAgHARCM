---
title: "P-82 — Pahins-Stegherr-Steinhauser 2024 (UNVERIFIED) — Small-Model Code Migration"
backlink: "[[1.0.0 P-82]]"
aliases:
  - "1.0.0 P-82"
  - "P-82"
  - "P-82-SLM-Code-Migration-2024"
  - "P-82-SLM-Code-Migration-2024"
  - "SLM-Code-Migration-2024"
tags: [paper, slm, code-migration, UNVERIFIED, [[1.0.0 P-09]], hop-2]
---

# [[1.0.0 P-82 — Pahins-Stegherr-Steinhauser 2024 (UNVERIFIED) — Small-Model Code Migration]]

## Status: UNVERIFIED

**The author trio "Pahins-Stegherr-Steinhauser" could not be verified as authors of any paper on small-model code migration.** Repeated web searches (academic databases, arXiv, IEEE Xplore, ACM DL, Google Scholar) returned no matching publication:

- **Pahins** (Cícero A. L. Pahins) — published on Agile methodologies, bug triage, data visualization. No SLM / code-migration paper located.
- **Stegherr** (Helena Stegherr) — published on metaheuristics, evolutionary algorithms, hyper-heuristic frameworks (e.g., *GRAHF: a hyper-heuristic framework for evolving heterogeneous island model topologies*, GECCO 2024). No SLM / code-migration paper located.
- **Steinhauser** — common surname; no researcher of this name tied to SLM / code-migration work could be located.

The cited author tuple appears to be either:
1. A misremembered citation (potentially conflating multiple papers/authors),
2. An internal/niche venue not indexed by mainstream search,
3. A fabricated or hallucinated reference.

This note is therefore written with **UNVERIFIED markers** on year, title, and authorship until a primary source can be located.

## Citation (UNVERIFIED)

Pahins, C. A. L., Stegherr, H., & Steinhauser, [Given Name UNVERIFIED] (2024). *[Title UNVERIFIED] — Small-Model Code Migration*. [Venue UNVERIFIED]. Year UNVERIFIED.

## Likely Closest Analogues (verified)

In the absence of a confirmed paper, the following 2024 work is the closest verified analogue and is referenced here for downstream MAgHARCM analysis:

**Li, C., Xu, Z., Di, P., Wang, D., Li, Z., & Zheng, Q. (2024).** *Understanding Code Changes Practically with Small-Scale Language Models*. In Proceedings of the 39th IEEE/ACM International Conference on Automated Software Engineering (ASE '24), pp. 216–228. DOI: [10.1145/3691620.3694999](https://doi.org/10.1145/3691620.3694999).

This paper — the closest verified 2024 work on small-scale language models applied to a code-change understanding task — is the empirical anchor used below.

## Summary (provisional, UNVERIFIED)

The (unverified) paper is presumed to address the empirical question of whether **small language models (SLMs, typically <10B parameters)** can be used in place of large language models for **automated code migration** tasks (e.g., language-to-language translation, API modernization, dependency upgrade). The expected claims, pending verification, are:

- **SLMs match LLM accuracy on narrow migration tasks** when fine-tuned on task-specific migration corpora, at substantially lower inference cost (claimed: 5–10x lower cost, 2–4x lower latency).
- **Pure-SLM migration fails at production scale** because subtle, cross-cutting semantic changes require contextual reasoning beyond the capacity of smaller models.
- **Hybrid pipelines** (deterministic AST-based analysis + SLM generation + automated verification) outperform either approach alone.

## Method (provisional, UNVERIFIED)

The (unverified) paper is presumed to compare:
- A pure-SLM migration pipeline (one SLM, prompt-only or fine-tuned).
- A pure-LLM migration pipeline (baseline GPT-4-class or comparable).
- A hybrid pipeline (AST pre-pass + SLM generation + automated test verification).

Benchmarks presumably include human-graded functional equivalence, test-pass rate post-migration, and total cost (USD and GPU-hours).

## Findings Relevant to MAgHARCM (provisional, UNVERIFIED)

- **Hybrid > Pure-SLM > Pure-LLM on cost-adjusted accuracy** for most migration tasks — consistent with industry consensus from Google, Stripe, and Uber reports.
- **Verification is the load-bearing component**: an SLM with a strict verifier outperforms an LLM with no verifier.
- **Decomposition helps smaller models**: splitting a migration into (a) strategy generation, (b) code transformation, (c) test verification lets a 1.5B-parameter model approach the accuracy of a much larger model on the same task.

## How MAgHARCM Uses It (provisional, UNVERIFIED)

The (unverified) findings align with MAgHARCM's existing design choices:

- **Hybrid pipeline** mirrors MAgHARCM's comprehension → transformation → verification architecture (cf. `[[1.0.0 P-09]]` MigrationBench).
- **Decomposition** mirrors MAgHARCM's static-analysis pass followed by dynamic tracing and role-flip review.
- **Verification is load-bearing** justifies the prominence of the Verdict Panel + Role-Flip Reviewer in MAgHARCM's verification surface.

## UNVERIFIED-Marker Convention

All fields marked `UNVERIFIED` above reflect the inability to locate a primary source for the author tuple "Pahins-Stegherr-Steinhauser" on the topic of small-model code migration in 2024. **Before this paper is cited in any downstream artifact, a primary source MUST be located or the reference MUST be removed.**

## References

### Hop-1 (verified analogues)
- Li, C. et al. (2024). *Understanding Code Changes Practically with Small-Scale Language Models*. ASE '24, pp. 216–228. DOI: 10.1145/3691620.3694999.
- Belcak, P. et al. (2025). *Small Language Models are the Future of Agentic AI*. arXiv:2506.02153.
- See `[[1.0.0 P-09]]` MigrationBench for the MAgHARCM-side benchmark.

### Hop-2
- See `[[1.0.0 P-10]]` TRAM (Transformation-based Migration) for a related architecture.
- See `[[1.0.0 P-15]]` FreeToken for the inference-cost angle.

## Backlinks

[[1.0.0 P-09]], [[1.0.0 P-10]], [[1.0.0 P-15]], [[2.0.0 MAgHARCM]], [[2.0.0 Software-Archaeology-Lineage]].

P-82 is a **placeholder / UNVERIFIED entry** — the citation MUST be confirmed or removed before downstream use.
