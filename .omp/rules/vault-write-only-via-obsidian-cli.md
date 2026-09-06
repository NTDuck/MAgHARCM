---
name: vault-write-only-via-obsidian-cli
description: "All research/consumption artefacts MUST be written and read from the MAgHARCM vault (id 41be7241693aebef) via the obsidian-cli skill — no other storage location, no other mechanism."
condition: "write\\s+to\\s+`?\\.obsidian/MAgHARCM|edit\\s+\\.obsidian/MAgHARCM|read\\s+\\.obsidian/MAgHARCM|research/|primitives/|diary/[^.]"
scope: ["tool:edit", "tool:write", "tool:read", "tool:bash"]
---

All MAgHARCM artefacts (research notes, paper notes, diary entries, primitives index, methodology, architecture, lineage) MUST be produced and consumed via the `obsidian-cli` skill targeting vault id `41be7241693aebef` (path `/home/ayin/.obsidian/MAgHARCM/MAgHARCM (old)`). Never write/read markdown through `edit`, `write`, `read`, or bare `bash` against `.obsidian/MAgHARCM/...` or top-level `research/`, `primitives/`, `diary/` paths. If `obsidian-cli` is unreachable, surface the gap and stop — do not silently fall back to direct filesystem tools.