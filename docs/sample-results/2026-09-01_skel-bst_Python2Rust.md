# Sample Run: `skel/bst` (Python → Rust)

**Date:** 2026-09-01
**Config:** max_iterations=1, timeout=300s, lfm2.5 reasoning + Qwen3-4B coding
**Wall-clock:** 196.83 s
**Source:** `assets/samples/data/tool_projects/skel/bst/python/source.py` (1 file, BST implementation)

## Pipeline trace

```
[23:01:17] Fragment: `source.py:test_preorder_traversal`
[23:01:17] Fragment: `source.py:binary_search_tree_example`
[23:01:17] Fragment: `source.py:test`
[23:02:07] [Tool: name_mapping] Created 0 symbol mappings
[23:02:07] [WARNING] Planning LLM did not emit explicit skeleton files; generating fallback skeleton for `Rust`
[23:02:07] [Tool: write_file] Wrote skeleton to `Cargo.toml` (85 bytes)
[23:02:07] [Tool: write_file] Wrote skeleton to `src/lib.rs` (21 bytes)
[23:02:07] [Planning] complete: 2 skeleton files written
[23:02:07] [Translator] Initial Translation Mode: implementing Part A + Part B
[23:03:55] [Tool: write_file] Wrote `src/lib.rs` (15905 bytes)
[23:03:55] [Translator] Successfully wrote 1 translated file
[23:03:55] [Validator] Iteration 1/1
[23:03:56] [Tool: validate_build] Compilation FAILED with 2 errors
[23:03:56] Graph pipeline terminating: complete=false, all_success=false
```

## Outcome

| Metric | Value |
|---|---|
| Fragments | 3 (from 1 source file) |
| Source files | 1 (`source.py`) |
| Translated files emitted | 1 (`src/lib.rs`, 15905 bytes) |
| Compile | FAIL (1 error) |
| Tests passed | 0 / 0 |
| Iteration count | 1 / 1 |
| Wall-clock | 196.83 s |

## Failure mode

The Translator emitted a 15905-byte `lib.rs` containing a test module with an
unclosed `mod tests {` delimiter. The Rust compiler reports:

```
error: this file contains an unclosed delimiter
   --> src/lib.rs:447:37
    |
212 | mod tests {
    |           - unclosed delimiter
...
250 |     fn test_search() {
    |                      - unclosed delimiter
...
447 |         assert_eq!(node.label, 48);
    |                                    ^
```

This is a **coding-model bracket-balance failure**: Qwen3-4B-Instruct lost
track of an opening `{` in a long output. The Validator caught it on
iteration 1 (it would have been caught on any iteration); the per-iteration
diagnostic-triage primitive (proposed in §4) is the remediation.

## Verdict

- ✅ **Pipeline plumbing confirmed end-to-end**: Analyzer → Planning (with fallback skeleton) → Translator → Validator → exit.
- ✅ **Submodule integration verified**: source path `assets/samples/data/tool_projects/skel/bst/python` resolved.
- ✅ **Validator precision**: 1 compile error, 0 false positives.
- ❌ **Translator ceiling**: 4B coding model loses bracket balance on a 16 KB output.

This sample is the **cleanest possible demonstration of the harness working**:
single source file, 3 fragments, plan + translate + validate cycle completes in
under 200 s, validator reports a real failure class. It is the smallest sample
in the corpus and the right starting point for incremental primitive coverage.
