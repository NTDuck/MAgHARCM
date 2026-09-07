---
name: sprint-name-must-match-today
description: "Sprint/iteration names MUST use today's actual date (YYYY-MM-DD per system reminder), never extrapolate from forward-dated artifacts"
condition: "Sprint[-_\\s](2026-0[89]|2026-1[0-2]|202[7-9])"
scope: ["text", "tool"]
---

Every sprint or iteration name MUST equal the current timestamp as YYYY-MM-DD taken from the system reminder or `date -u`. Never extrapolate a next-sprint date from a previous handoff's filename, the vault's date stamp, or git log dates — those are project conventions, not the calendar. If a handoff for today already exists, append `-{i}` (e.g. `2026-09-07-1`). Verify the actual current date before writing any `Sprint-YYYY-MM-DD-...` filename, diary entry, todo phase name, or prose reference. Re-read the system reminder at the start of every turn if a forward-dated handoff is in context.