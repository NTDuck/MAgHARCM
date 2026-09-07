---
name: incremental-commit-per-prompt
description: "Enforce incremental conventional commits per prompt after each milestone or phase boundary"
scope: ["text", "tool"]
---

# Incremental Commit Enforcement

Commit changes incrementally during every prompt. Do not wait for the end of the prompt to commit all changes together.

## 1. Commit Rules

1. Commit after each completed phase, task, or milestone.
2. Make multiple commits in a single prompt when you touch multiple areas.

3. Use Conventional Commit prefixes:
   - Prefix feat for new capabilities, research papers, or vault artifacts.
   - Prefix fix for bug fixes, compile repairs, or blocker resolutions.
   - Prefix docs for handoffs, documentation updates, or paper edits.

4. Use additional prefixes as needed:
   - Prefix test for benchmark runs or test suite changes.
   - Prefix refactor for structural code improvements.

5. Keep commit messages terse and imperative. Do not use marketing adjectives.
6. Every commit must pass all validation checks before you make it.

## 2. Sprint Workflow Rule

1. When you run the `/MAgHARCM` command, obey its phase commit sequence:
   - Commit after sprint inception.
   - Commit after research paper persistence.
   - Commit after vault synchronization.
   - Commit after empirical benchmark execution.

2. Continue the commit sequence for later phases:
   - Commit after codebase fixes and invariant defense.
   - Commit after academic paper synchronization.
   - Commit after command evolution.
   - Commit after sprint handoff persistence.

3. Every `/MAgHARCM` sprint iteration must edit `.omp/commands/MAgHARCM.md` to keep it updated with research progress.
