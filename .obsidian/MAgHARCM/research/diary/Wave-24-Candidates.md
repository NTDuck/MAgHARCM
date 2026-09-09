---
title: Wave-24 Candidate Evaluation (2026-09-09)
backlink: "[[1.0.0 Wave-24 Candidates]]"
tags: [candidates, wave-24, iclr-2026, ase-2026, kv-cache, code-translation, software-archaeology, "[[1.0.0 P-157]]", "[[1.0.0 P-158]]", slm-era, status: FIRED]
status: FIRED
date: 2026-09-09
last_updated: 2026-09-09 (iter-1, wave-24)
aliases:
 - "Wave-24-Candidates"
 - "Wave 24 Candidates"
 - "Wave-24"
---

# [[1.0.0 Wave-24 Candidates]]

## Trigger-gate evaluation

§7 trigger gate (Wave-24 firing criterion): Q1 venue confirmation (peer-reviewed at NeurIPS / ICML / ICLR / ICSE / ASE / TOSEM / TSE / FSE; workshop-track permitted ONLY if method-level single-paper) + Q2 mechanism-vs-benchmark (concrete method, not benchmark / eval / prompt tweak alone) + Q3 anchoring-to-existing-primitive (defends or refutes an existing `[[1.0.0 PRIM-NN]]`).

## Wave-24 focus areas (from Wave-23 forward plan)

1. **Cross-agent KV cache policies** — Wave-23 forward plan asked for policies beyond `[[1.0.0 P-147]]` SpecKV / `[[1.0.0 P-148]]` LookaheadKV / `[[1.0.0 P-149]]` SSD/Saguaro / `[[1.0.0 P-141]]` KVFlow at ICLR 2026 / NeurIPS 2026. Directly relevant to MAgHARCM's 8-agent graph, where RelayCaching (`[[1.0.0 P-134]]`) is the sole cross-agent reuse substrate.
2. **Source-to-source modernization** — Wave-23 forward plan asked for advances beyond `[[1.0.0 P-145]]` TerraMod + `[[1.0.0 P-151]]` SmartC2Rust + `[[1.0.0 P-154]]` TransAgent; ASE 2026 main-track program (October 2026, Munich) is now announced.
3. **Partition-aligned comprehension residuals** — concept-assignment / temporal-coupling / DSM-partition mechanisms remain open at NeurIPS 2026 / ICML 2027 / ICLR 2027 (Wave-23 closure note). NeurIPS 2026 notifications (2026-09-24) post-date this wave; no concept-assignment mechanism surfaced at the concluded venues.
4. **Standing scope re-verification** — W1 SWE-TRACE (fifth carry) re-check; R1 AutoCodeSherpa (ISSTA off-list) companion-paper check. Standing scope: small language models (4B-30B) and software archaeology.

## Scanned topic families

| Family | Venue set | Outcome |
| :--- | :--- | :--- |
| Cross-agent / agentic KV cache (eviction, reuse, region-aware) | ICLR 2026 posters (MixKV, ReST-KV, KVTC, PM-KVQ, MemArt submission), arXiv agent-inference cluster (MemDecay, ReCache) | 2 ACCEPT (ReST-KV, ReCache); 2 REJECT; 1 REJECT (MemArt unconfirmed); MemDecay arXiv-only |
| Source-to-source modernization | ASE 2026 main track (announced program) | 1 REJECT (ReCodeAgent journal-version = duplicate mechanism of P-01; venue upgrade logged, no new mechanism) |
| Concept-assignment / temporal-coupling / DSM partition | ICSE 2026 / FSE 2026 concluded programs + ICLR 2026 listing | 0 mechanisms surfaced; dry family this wave |
| SLM-scale code translation (4B-30B) | NeurIPS 2026 (pre-notification; SLM-Agents workshop is 2026 workshop-track), ICLR 2026 | No admissible method-level single-paper mechanism outside the KV family; SWE-TRACE carried |
| Software archaeology / architecture recovery | ASE 2026, TOSEM 2026 | No new mechanism beyond the Wave-18 trio (SSAR/SemArc/SemRef); dry family this wave |

## Candidates (6 triaged, 1 ACCEPT + 1 REJECT + 4 UNVERIFIED/watchlist; 1 non-event venue note)

### Candidate 1 — ReST-KV: Robust KV Cache Eviction with Layer-wise Output Reconstruction and Spatial-Temporal Smoothing (An, Lu, Zhu, Yu, Zhao, Wu, Tang & Wang, ICLR 2026, arXiv:2605.08840)

**Status: ACCEPT.** ICLR 2026 (poster, iclr.cc/virtual/2026/poster/10009650, OpenReview PhEHuo7oMm) — flagship venue, in-list per §7. Mechanism is concrete: (a) **layer-wise output reconstruction** — KV eviction formulated as an optimization problem that minimizes output discrepancy by directly modeling how each token's removal affects the layer output, capturing attention-redistribution effects that raw-attention scoring misses; (b) **spatial-temporal smoothing** — exponential-moving-average smoothing over temporal token-importance variation plus an adaptive window mechanism for spatial patterns; (c) evaluated at 128k context with 10.61x decoding-latency reduction, +2.58% LongBench, +15.2% RULER. Anchors `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement (eviction quality governs what retained context the retrieval loop sees) and `[[1.0.0 PRIM-21]]` Migration Strategy Selection (output-reconstruction vs attention-scoring = eviction-strategy choice knob alongside P-147 SpecKV / P-148 LookaheadKV).

### Candidate 2 — ReCache: Efficient KV Cache Reuse and Compression for Tool-Augmented LLM Agents (Fang, Wei, Hu & Shen, arXiv:2608.19662)

**Status: UNVERIFIED (watchlist W24-W4; fails Q1 venue gate).** Preprint as of 2026-09-09; no peer-reviewed venue confirmation found on dblp/Semantic Scholar/OpenReview. Mechanism is concrete and directly on the Wave-23 forward-plan family: (a) **resource-wise attention** — removes cross-resource attention interactions and assigns resource-local positions, producing composition-invariant KV blocks reusable across tool-schema combinations and orders; (b) **contribution-guided structural pruning + field-aware semantic pruning** — restricts resource visibility to contribution-selected layer-KV-head-group routes and retains only invocation-critical fields; (c) 92.43% KV-memory reduction, 1.423x attention speedup, 3.655x TTFT speedup at matched accuracy (82.3% vs 82.4% Inv-F1). This is the **cross-agent KV-cache reuse mechanism** the Wave-23 forward plan targeted: the closest analogue to `[[1.0.0 P-134]]` RelayCaching for MAgHARCM's tool-schema-heavy agent roster. Anchors `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement + `[[1.0.0 PRIM-9]]` Tri-Representation Hybrid Code Graph (composition-invariant resource blocks = graph-node-level cache substrate). NOT admitted: Q1 fails (arXiv-only, no venue acceptance as of 2026-09-09). Carried on `watchlist.wave-24` (W24-W4); fires at Q1 confirmation per the P-138 RepairKV precedent.

## Rejected (1)

### Venue note (non-event, not triaged) — ReCodeAgent ASE 2026 main-track version (Ibrahimzada et al.)

**Status: NON-EVENT (already anchored; not in any registry).** `[[1.0.0 P-01]]` ReCodeAgent is already anchored with the identical four-agent mechanism. Search surfaced an ASE 2026 main-track acceptance claim (researchr.org track listing + author page + arXiv:2604.07341) for the same paper; logged here for visibility only. Not a Wave-24 triage candidate — the reject/watchlist registries record papers that failed or await Wave-24 gate decisions, and an already-anchored paper does neither. P-01's venue field gains the ASE 2026 main-track confirmation at the next database touch.

### R1 — MixKV: Mixing Importance with Diversity for KV Cache Compression in LVLMs (Liu et al., ICLR 2026, arXiv:2510.20707)

**Status: REJECT Q3 (off-axis modality).** ICLR 2026 accepted (github.com/xuyang-liu16/MixKV confirms) — Q1 passes; but the mechanism targets large vision-language model KV caches (modality-specific, head-wise semantic redundancy for multi-modal sequences). MAgHARCM's pipeline is text-only code translation; no primitive defends or is defended by an LVLM-specific compression mechanism. Off-axis per Q3.

## Watchlist — UNVERIFIED (3 new + 1 carried)

### W1 — SWE-TRACE: Optimizing Long-Horizon SWE Agents Through Rubric Process Reward Models and Heuristic Test-Time Scaling (Han et al., arXiv:2604.14820)

**Status: UNVERIFIED (sixth carry, Wave-19 W2 → Wave-20 W1 → Wave-21 W1 → Wave-22 W1 → Wave-23 W1 → Wave-24 W1).** Re-verification 2026-09-09: arXiv:2604.14820 remains preprint-only; no venue acceptance found. NeurIPS 2026 author notifications scheduled 2026-09-24 (15 days out, post this wave). Per the Wave-21 carry rule recorded in `watchlist.wave-23` ("sixth carry = retire threshold"), Wave-25 must either confirm a venue or retire the candidate.

### W2 — MemArt: KVCache-Centric Memory for LLM Agents (OpenReview YolJOZOGhI, under review)

**Status: UNVERIFIED (first carry).** In-review ICLR 2026 submission; OpenReview forum behind browser verification, no dblp or conference-program entry found as of 2026-09-09. Mechanism (latent KV-block agent memory replacing plaintext memory re-feeding) is real and adjacent to P-134 RelayCaching. Watchlist for Wave-25 venue re-verification.

### W3 — MemDecay: Region-Aware KV Cache Eviction for Efficient LLM Agent Inference (Matam & Kim, arXiv:2607.10582)

**Status: UNVERIFIED (first carry).** July 2026 preprint; no peer-reviewed venue confirmation. Mechanism concrete (region-specific base priorities + decay rates + pinning; Qwen2.5-1.5B/3B evaluation = SLM-relevant) but fails the Q1 venue gate. Watchlist for Wave-25.

### W4 — ReCache: Efficient KV Cache Reuse and Compression for Tool-Augmented LLM Agents (Fang, Wei, Hu & Shen, arXiv:2608.19662)

**Status: UNVERIFIED (first carry; fails Q1 venue gate).** arXiv:2608.19662 (2026-08-20) preprint only; no venue acceptance as of 2026-09-09. Q2 concrete (resource-wise attention producing composition-invariant tool-schema KV blocks), Q3 on-target (cross-agent/tool-schema KV reuse = Wave-23 forward-plan family; P-134 RelayCaching tool-side analogue). Fires at Q1 confirmation per the P-138 RepairKV precedent.

> **Partition summary.** ACCEPT = 1 (P-157 ReST-KV). REJECT = 1 (R1 MixKV off-axis). UNVERIFIED = 4 (W1 SWE-TRACE sixth carry; W2 MemArt; W3 MemDecay; W4 ReCache). Total triaged this wave = 6 (plus 1 non-event venue note).

## Wave-24 outcome

- **1 new P-NN anchor persisted**: `[[1.0.0 P-157]]` ReST-KV (ICLR 2026).
- **1 rejection** logged with per-candidate rationale (R1 MixKV off-axis LVLM modality). 1 non-event venue note recorded (ReCodeAgent ASE 2026 version of P-01; not triaged, not registered).
- **4 UNVERIFIED items** on `watchlist.wave-24`: W1 SWE-TRACE (sixth carry; retire-or-confirm decision due Wave-25 after NeurIPS 2026 notifications land 2026-09-24), W2 MemArt (in-review ICLR 2026), W3 MemDecay (arXiv-only), W4 ReCache (arXiv-only; Q1-pending method-level mechanism).
- **Vault paper count: 156 → 157.**
- **Coverage gains**:
  - `[[1.0.0 PRIM-31]]` Iterative Retrieval Refinement: P-157 (output-reconstruction-governed eviction determines retained-context fidelity for the retrieval loop). One new retrieval-refinement anchor.
  - `[[1.0.0 PRIM-21]]` Migration Strategy Selection: P-157 (reconstruction-based vs attention-heuristic eviction = strategy knob in the SpecKV/LookaheadKV family). One new strategy-selection anchor.

## Watchlist for Wave-25+

- **W1 SWE-TRACE** (sixth carry) — NeurIPS 2026 notifications land 2026-09-24; Wave-25 applies the retire-or-confirm decision.
- **W2 MemArt + W3 MemDecay + W4 ReCache** — re-verify venue acceptance at NeurIPS 2026 / ICLR 2027 cycles; ReCache fires at Q1 confirmation per the P-138 RepairKV precedent.
- **Partition-aligned comprehension mechanisms** (concept-assignment / temporal-coupling / DSM) — remained dry at all concluded 2026 venues; next look: NeurIPS 2026 (post-09-24) and ICML 2027 / ICLR 2027.
- **AutoCodeSherpa companion** (REJECT Q1 Wave-23, ISSTA) — author pages confirm ISSTA 2026 full-paper acceptance only; watch for an ICSE 2027 / FSE 2027 / TOSEM companion.
- **P-01 ReCodeAgent venue upgrade** — propagate "ASE 2026 main track" into the P-01 `venue` field at the next database touch.

**TOTAL FIRED: 1 ACCEPT + 1 REJECT + 4 UNVERIFIED (3 new + 1 carried).** Headline ACCEPT count: 1. Headline REJECT count: 1 (Q3 off-axis). UNVERIFIED: 4.
