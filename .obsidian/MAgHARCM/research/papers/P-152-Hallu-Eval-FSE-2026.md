---
title: "P-152 Hallucinations in LLM-Based Code Summarization: Unveiling, Detection, and Mitigation (Hallu-Eval + Hallu-Det + Hallu-Shield)"
backlink: "[[1.0.0 P-152]]"
aliases:
 - "1.0.0 P-152"
 - "P-152"
 - "P-152-Hallu-Eval-FSE-2026"
 - "Hallu-Eval-FSE-2026"
 - "hallueval2026fse"
tags: [paper, code-summarization, hallucination, benchmark, detection, mitigation, slm, "[[1.0.0 PRIM-22]]", "[[1.0.0 PRIM-25]]", wave-22]
date: 2026-09-07
last_updated: 2026-09-07 (iter-6, wave-22)
venue: FSE 2026 (Research Track)
---

# [[1.0.0 P-152]] Hallu-Eval + Hallu-Det + Hallu-Shield

## TL;DR

Hallu-Eval is a **systematic hallucination-evaluation + detection + inference-time mitigation triplet** for LLM-based code summarization. (a) **Hallu-Eval** — 800-pair benchmark of original + semantically-perturbed code/summary pairs (natural + induced logical hallucinations); (b) **Hallu-Det** — detection approach for code-summarization hallucinations; (c) **Hallu-Shield** — inference-time mitigation method. Anchors `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Hallu-Eval as the systematic hallucination-evaluation substrate for the Comprehension phase) and `[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination (Hallu-Shield as inference-time mitigation alongside P-55/P-81 cloze-reformulation and P-143 HalluShield SLM-grounded speculative-decoding defence).

## Mechanism (Q2)

1. **Hallu-Eval benchmark construction** — 800 code-summary pairs assembled from (a) original code snippets (capturing naturally-occurring hallucinations) and (b) semantically-perturbed counterparts (systematically inducing logical hallucinations). The perturbation step is the inductive-bias lever: it lets the benchmark measure the LLM's faithfulness under controlled semantic-delta stress.
2. **Hallu-Det (detection)** — a detection approach that takes a (code, summary) pair and outputs a hallucination-likelihood signal; trained on the Hallu-Eval benchmark.
3. **Hallu-Shield (inference-time mitigation)** — an inference-time method that mitigates hallucinations by adjusting the LLM's decoding or re-ranking candidate summaries using the Hallu-Det signal; runs during generation, not post-hoc re-training.
4. **Triplet integration** — Hallu-Eval → Hallu-Det (trained) → Hallu-Shield (inference-time) is a closed loop: benchmark produces the labelled data; detector consumes it; mitigation uses the detector.

Together these give MAgHARCM a **systematic hallucination-evaluation + inference-time mitigation substrate** for the Comprehension phase of code summarization — the first benchmark + detection + mitigation triplet in the SLM-era comprehension substrate.

## Anchoring (Q3)

| Primitive | Pre-wave-22 behaviour | Hallu-Eval / Hallu-Det / Hallu-Shield substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension | Comprehension anchored on P-128 KVzip (query-agnostic context reconstruction), P-129 LλMDA (LLM-aided partial PDG), P-133 ADI (Frame Lifetime Trace), P-139 TypePro (inter-procedural SDG slicing), P-140 Panta (iterative test generation), P-142 NESA (Datalog policy), P-150 TestPrune (coverage-driven context pruning); no systematic hallucination-evaluation substrate | Hallu-Eval as the systematic hallucination-evaluation substrate for the Comprehension phase (closes the evaluation gap) |
| `[[1.0.0 PRIM-25]]` Role-Flip De-Hallucination | De-hallucination anchored on P-55/P-81 cloze-reformulation, P-143 HalluShield SLM-grounded speculative-decoding defence; no inference-time mitigation for code-summarization hallucinations | Hallu-Shield as inference-time mitigation alongside P-143 HalluShield (closes the inference-time mitigation gap) |

## Hop-1 Citations

- P-143 HalluShield (FSE 2026, Wave-20) — SLM-grounded speculative-decoding hallucination defence; Hallu-Shield is the inference-time analogue for code summarization.
- P-127 SLM-as-a-Judge (ICLR 2026, Wave-18) — SLM-scale judgment substrate; Hallu-Det is the LLM-as-detector analogue.
- P-97 Welleck Self-Correct 2024 — trained self-correction as cheaper alternative to multi-voter panel; Hallu-Shield is the inference-time analogue.
- TruthfulQA (Lin 2021) — foundational hallucination benchmark; Hallu-Eval is the code-summarization analogue.
- FactCC (Kryściński 2020) — factual-consistency detection; Hallu-Det is the code-summarization analogue.

## Hop-2 Citations

- Hallucination in NLG (Ji 2023 survey) — foundational taxonomy of NLG hallucinations; Hallu-Eval/Hallu-Det/Hallu-Shield follow the taxonomy's "logical" / "factual" / "faithfulness" distinction.
- Code summarization (Allamanis 2018, Wan 2018, Iyer 2016) — foundational code-summarization substrate; Hallu-Eval's perturbation step is grounded in semantic-equivalence stress-testing from the foundational literature.
- Inference-time intervention (Li 2023, ITI) — foundational inference-time mitigation pattern; Hallu-Shield is the code-summarization analogue.
- Detector-as-mitigation (Manakul 2023, SelfCheckGPT) — detector-driven mitigation pattern; Hallu-Shield follows the SelfCheckGPT pattern.

## MAgHARCM integration

- **YAML config key**: `agents.comprehension.hallucination.detector: "halludet"`; `agents.comprehension.hallucination.mitigator: "hallushield"`; `agents.comprehension.hallucination.benchmark: "hallueval"`; `agents.comprehension.hallucination.mitigation_strength: <float>`.
- **Implementation file**: `internal/comprehension/hallucination.go::NewHalluShield` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-22]]`, `[[1.0.0 PRIM-25]]`.

## Caveats

- **Code-summarization only** — Hallu-Eval is a code-summarization benchmark; generalisation to code generation / code translation is not in scope (different hallucination taxonomy).
- **Single-language** — Hallu-Eval is not language-tagged in the abstract; the 800-pair corpus likely spans one or two languages; multi-language generalisation is future work.
- **Inference-time cost** — Hallu-Shield adds inference-time cost (detector + re-ranking); must be amortised over the summary-quality gain.
- **Detector-label coupling** — Hallu-Det is trained on Hallu-Eval; cross-benchmark detector transfer is not verified.

## Source

- Venue: FSE 2026 (Research Track), verified via conf.researchr.org/fse-2026-research-papers/66; PACMSE Vol. 3, Issue FSE, DOI: 10.1145/3808189.
- arXiv: not surfaced (no arXiv preprint confirmed); FSE 2026 is the canonical venue.

## BibTeX

```
@inproceedings{hallueval2026fse,
  title     = {Hallucinations in LLM-Based Code Summarization: Unveiling, Detection, and Mitigation},
  author    = {Liu, Yufan and others},
  booktitle = {Proceedings of the ACM International Conference on the Foundations of Software Engineering (FSE)},
  year      = {2026},
  doi       = {10.1145/3808189},
  volume    = {3},
  number    = {FSE},
  series    = {Proceedings of the ACM on Software Engineering (PACMSE)}
}
```
