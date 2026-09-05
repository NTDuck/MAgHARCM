---
title: "ADR-2026-09-07 Dup-Row Escape Recipe — When the Same [[P-NN]] Belongs in Two Cells"
backlink: "[[1.0.0 ADR-2026-09-07-Dup-Row-Escape-Recipe]]"
status: Accepted
date: 2026-09-07
tags: [adr, lineage, matrix, [[1.0.0 Architecture]]]
---

# [[1.0.0 ADR-2026-09-07 Dup-Row Escape Recipe]]

## 1. Context

The lineage matrix at `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md` is the canonical map: `PRIM-NN → primary citation → secondary citation → producer file`. Several `[[1.0.0 P-NN]]` papers legitimately apply to more than one primitive:

- `[[1.0.0 P-42]]` (MacCormack, Rusnak & Baldwin 2006 DSM) — applies to BOTH `PRIM-3` (Target Skeleton-First Gen, where DSM-driven cycle→freeze→extract is the operationalisation) AND `PRIM-19` (Design Rule Hierarchy Partition, where DSM cycles map to modularisation decisions).
- `[[1.0.0 P-41]]` (2024 retrospective) — referenced twice in PRIM-19 for both Baldwin-Clark legacy lineage and recent empirical re-confirmation.

Two prior sprints (2026-09-05 and 2026-09-06) hit the same issue: a sweep removing duplicate `[P-NN]` entries either over-cleaned (lost the legitimate second anchor) or under-cleaned (left accidental re-listings). The duplication is structural, not a bug.

## 2. Decision

A `[[1.0.0 P-NN]]` citation appearing in **more than one** row of the lineage matrix is **legitimate iff**:

1. The two rows cite the SAME `[[1.0.0 P-NN]]` but the row's `Hop-2 secondary` cell is **DIFFERENT** (genuine multi-anchor application), OR
2. The two rows cite the SAME `[[1.0.0 P-NN]]` in different positions (e.g. once in `Primary citation` and once in `Hop-2 secondary`) — explicitly justified in the row's evidence column.

Otherwise (same row contents, same position) it is **accidental re-listing** and MUST be removed.

## 3. Recipe — Three-Question Test

For each `[[1.0.0 P-NN]]` that appears in ≥2 rows of the lineage matrix, ask:

### Q1. Is the hop-2 secondary citation cell DIFFERENT between the rows?

```text
Row A: PRIM-N | P-XX | Hop-2 = (Smith 2005)
Row B: PRIM-M | P-XX | Hop-2 = (Jones 2008)
```

→ If YES → legitimate multi-anchor. KEEP both. Add a footnote in the row explaining which aspect of `P-XX` applies.

→ If NO → go to Q2.

### Q2. Does the position differ (primary vs hop-2)?

```text
Row A: PRIM-N | Hop-1 = P-XX
Row B: PRIM-M | Hop-2 = P-XX
```

→ If YES → legitimate hierarchical reference (P-XX cites another paper; that other paper cites P-XX via hop-2). KEEP both. The hop-2 slot is a transitively cited primary reference.

→ If NO → accidental re-listing. REMOVE the duplicate. Keep the more-specific row (the one whose hop-2 also references P-XX).

### Q3. Does the row evidence column justify the duplicate?

```text
Row A: PRIM-N | P-XX | Hop-2 = ... | Evidence: "P-XX cycle-prediction metric anchors PRIM-N cycle-breaking"
Row B: PRIM-M | P-XX | Hop-2 = ... | Evidence: "P-XX interface-stability finding underpins PRIM-M's freeze-before-extract"
```

→ If BOTH rows have distinct evidence text → KEEP. The evidence column differentiates the application.

→ If evidence is identical → REMOVE the duplicate.

## 4. Worked Example — `[[1.0.0 P-42]]`

| Row | Cell content | Legitimate? |
| :--- | :--- | :--- |
| `PRIM-3` row 129 | Primary: `P-42`. Hop-2: `Baldwin & Clark ([[1.0.0 P-34]]) Design Rules`. Evidence: "DSM-based refactoring strategy (cycle → freeze interface → extract submodule) operationalises skeleton-first synthesis." | ✅ |
| `PRIM-19` row 145 | Primary: `P-42`. Hop-2: `Baldwin & Clark ([[1.0.0 P-34]]) Modularity, P-41 2024 retrospective`. Evidence: "DSM cycles map directly to design-rule partition boundaries; P-41 confirms in 2024 retrospective." | ✅ |

Q1: Hop-2 cells DIFFER (`Design Rules` vs `Modularity + P-41`). → KEEP both.

Q3: Evidence text DIFFERENT. → KEEP both.

Result: two legitimate rows. Sprint 2026-09-06 retained both. Correct outcome.

## 5. Anti-Patterns

- **Mechanical dedup of any duplicate `P-NN`**. Sweeps that grep for `P-NN` and remove all-but-one BREAK legitimate multi-anchor applications. Always apply the three-question test.
- **Single-link over-extension**. Putting P-NN in many rows with the same hop-2 and no differentiated evidence is NOT multi-anchor — it's sprawl. Apply Q3 strictly.
- **Phantom hop-2 padding**. Adding `P-NN` to a row's hop-2 cell solely to justify a primary citation is not legitimate. The hop-2 cell must contain a paper-cited-by-hop-1.

## 6. Audit Workflow

1. `grep -n '\[\[P-NN\]\]' .obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md | sort` — list all occurrences.
2. Group by `[[1.0.0 P-NN]]` and check the rows.
3. For each group with N≥2 occurrences, apply Q1 → Q2 → Q3.
4. Record findings in `.obsidian/MAgHARCM/diary/Sprint-YYYY-MM-DD-Handoff.md` §"Lineage sweep".
5. Do NOT auto-remove. Auto-removal breaks multi-anchor rows.

## 7. Pointers

- Companion ADR: [[1.0.0 ADR-2026-09-07-Sprint-Conventions]] (full rule-set)
- Lineage matrix: `.obsidian/MAgHARCM/research/Software-Archaeology-Lineage.md`
- Example multi-anchor case: `[[1.0.0 P-42]]` rows 129 + 145
- Sprint audit reports: `.obsidian/MAgHARCM/diary/Sprint-2026-09-05-Handoff.md`, `.obsidian/MAgHARCM/diary/Sprint-2026-09-06-Handoff.md`
