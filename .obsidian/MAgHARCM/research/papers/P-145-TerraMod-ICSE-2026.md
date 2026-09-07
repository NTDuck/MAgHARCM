---
title: "P-145 TerraMod: Automating Terraform Code Migration through Provider Evolution Knowledge"
backlink: "[[1.0.0 P-145]]"
aliases:
  - "1.0.0 P-145"
  - "P-145"
  - "P-145-TerraMod-ICSE-2026"
  - "TerraMod-ICSE-2026"
  - "terramod2026icse"
tags: [paper, code-migration, knowledge-augmentation, infrastructure-as-code, "[[1.0.0 PRIM-21]]", "[[1.0.0 PRIM-22]]", wave-20]
date: 2026-09-07
last_updated: 2026-09-07
venue: ICSE 2026 (NIER Track)
---

# [[1.0.0 P-145]] TerraMod

## TL;DR

TerraMod is a **knowledge-augmented Terraform configuration migration framework** that constructs a structured migration context from **changelogs + API schemas + deprecation links** to guide LLM-driven upgrades across Terraform provider versions. Evaluated on real-world breaking changes from the AWS Terraform Provider; significantly reduces manual effort compared to standard prompting. Anchors `[[1.0.0 PRIM-21]]` (Migration Strategy Selection: knowledge-augmented migration context as a runtime signal for strategy choice) and `[[1.0.0 PRIM-22]]` (Four Phases of Comprehension: Reorganization phase via external knowledge context).

## Mechanism (Q2)

1. **External knowledge harvesting** — TerraMod scrapes Terraform provider **changelogs**, **API schemas**, and **deprecation links** for the source and target provider versions.
2. **Structured migration context construction** — the harvested knowledge is composed into a migration context that captures deprecated resources, attribute changes, and breaking-change semantics.
3. **LLM-guided migration** — the migration context is fed to the LLM at prompt time; the LLM produces the migrated configuration guided by the external knowledge.
4. **Evaluation** — real-world breaking changes from the AWS Terraform Provider; TerraMod outperforms standard prompting and significantly reduces manual effort.

Together these give MAgHARCM's source-to-source modernization loop an **external-knowledge-augmented migration context** as a first-class input alongside the source code and target skeleton — closing the gap identified by BLK-02 (Commons-Validator plateau) for the IaC modernization axis.

## Anchoring (Q3)

| Primitive | Pre-wave-20 behaviour | TerraMod substrate |
| :--- | :--- | :--- |
| `[[1.0.0 PRIM-21]]` Migration Strategy Selection | TryInOrder registry lacked an external-knowledge signal for migration strategy choice; strategies were selected on purely internal signals (compilation, test pass) | TerraMod's migration context is an external-knowledge signal that the registry can use to bias strategy choice toward knowledge-augmented strategies when breaking-change density is high |
| `[[1.0.0 PRIM-22]]` Four Phases of Comprehension (Reorganization phase) | Reorganization relied on syntactic rewriting + LLM reasoning with no formal external-knowledge input | TerraMod's changelog + API-schema + deprecation-link context is the formal external-knowledge input that the Reorganization phase can consume |

## Hop-1 Citations

- AWS Terraform Provider breaking-change history (industrial source).
- IBM enterprise-scale COBOL-to-Java translation (P-145 sister industrial migration research).
- API-derived-rules approach (Forge 2026 Workshop, sister work).

## Hop-2 Citations

- Knowledge-augmented LLM prompting (Lewis et al. RAG 2020, Borgeaud et al. 2022 RETRO).
- Infrastructure-as-code evolution (Munch-Ellingsen et al. 2023).
- Code-migration lineage (P-02 AlphaTrans, P-103 AgentModernize, P-140 Panta).

## MAgHARCM integration

- **YAML config key**: `migration.terramod.knowledge_sources: ["changelog", "api_schema", "deprecation_links"]`; `migration.terramod.provider_versions: {source: "5.x", target: "6.y"}`.
- **Implementation file**: `internal/migration/terramod_context.go::NewTerraModContext` (forthcoming — future sprint).
- **Affected primitives**: `[[1.0.0 PRIM-21]]`, `[[1.0.0 PRIM-22]]`.

## Caveats

- **Knowledge-source freshness** — changelogs and deprecation links must be kept current; stale knowledge produces stale migrations.
- **Provider coverage** — TerraMod is validated on the AWS provider; generalization to GCP / Azure / Oracle providers requires knowledge-base extension.
- **NIER scope** — NIER papers present forward-looking ideas with preliminary data; the full mechanism awaits a Research-Track companion paper.

## Source

- Venue: ICSE 2026 (NIER Track), verified via researchr.org/details/icse-2026/icse-2026-nier/43; IBM Research publication page research.ibm.com/publications/automating-terraform-code-migration-through-provider-evolution-knowledge.
