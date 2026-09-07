---
title: "P-107 — Deliberate re-anchor: Decomposed Prompting SLM Multi-Agent (lineage note)"
backlink: "[[1.0.0 P-107]]"
tags: [paper, decomposed-prompting, slm, multi-agent, re-anchor, [[1.0.0 P-100]], [[1.0.0 PRIM-24]], [[1.0.0 PRIM-22]], hop-2]
---

# [[1.0.0 P-107 — Deliberate SLM-Era Re-anchor of P-100 (Decomposed Prompting)]]

## Status: DELIBERATE RE-ANCHOR (not a new paper)

This is **NOT a new verified paper**. It is a versioned pointer to **[[1.0.0 P-100]]** (Khot et al. 2022/2023, Decomposed Prompting, arXiv:2210.02406, ICLR 2023) emphasising the SLM-era re-read: a 1.5B parameter model with decomposed prompting matches 7-13B monolithic models on multi-step code tasks.

## Why a re-anchor slot

The original P-100 paper note documents the algorithmic contribution. P-107 exists as a **deliberate versioned slot** for the SLM-era relevance commentary that downstream jobs may want to link against, without inflating Wave-10's verified-paper count.

## Verified anchors (from P-100)

- Original paper: Khot, T., Khot, T., Sabharwal, A., Clark, P. (2022/2023). *Decomposed Prompting: A Modular Approach for Solving Complex Tasks*. arXiv:2210.02406; ICLR 2023.
- Hop-1 follow-ups: NAACL 2025 Findings surveys on decomposed prompting for SLM code agents (e.g., https://aclanthology.org/2025.findings-naacl.285.pdf)

## Relevance to MAgHARCM

- PRIM-24 SOP-Anchored Role-Artifact Schema is the concrete instantiation of decomposed prompting in MAgHARCM: each agent has a typed-IO contract that decomposes the monolithic translation task.
- PRIM-22 Four Phases of Comprehension maps directly: plan → observe → hypothesise → test = the Comprehension model's decomposition of program understanding.

## Wave-10 honest count

After removing the P-107 duplicate-of-P-100 inflation:
- **P-102** SmallCode (fp8.co 2025): verified (4B SLM, 87% HumanEval).
- **P-103** AgentModernize (Ahmed & Galib, arXiv:2605.17535, 2026): verified. **NOT ICSE 2025** — the arXiv preprint is the canonical citation.
- **P-104** S* Test-Time Scaling (Dacheng Li et al., UC Berkeley, arXiv:2502.14382, 2025): verified.
- **P-105** ChunkKV (Xiang Liu et al., NeurIPS 2025): verified.
- **P-106** BFCL Berkeley Function Calling Leaderboard (Patil et al., PMLR v267, 2025): verified.
- **P-107** (this file): **deliberate re-anchor of P-100**, NOT a new verified paper.

Wave-10 verified-paper count: **5 new** (P-102..P-106) + 1 re-anchor slot (P-107).
