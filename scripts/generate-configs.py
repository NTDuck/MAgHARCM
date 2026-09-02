#!/usr/bin/env python3
"""Generate .config/<category>-<project>.yml for every assets/samples project.

Detects source language from the named subdir (c/, go/, java/, python/,
javascript/, typescript/) under each project, picks the per-category
canonical source if multiple are present, and emits a translation
request that targets Rust.

Iterates over:
  assets/samples/<category>/<project>/<lang>/

Skips projects whose name already has a manual config in .config/
(e.g. gildedrose.yml, stats.yml) — these use hand-tuned settings.

Output: N generated configs + manual = total in .config/.

Usage:
  python3 scripts/generate-configs.py [--dry-run]
"""
from __future__ import annotations
import argparse
import os
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent
SAMPLES = REPO / "assets" / "samples"
CONFIG_DIR = REPO / ".config"

LANG_MAP = {
    "c": "C", "go": "Go", "java": "Java", "python": "Python",
    "javascript": "JavaScript", "typescript": "TypeScript", "rust": "Rust",
}

MODELS = {
    "C":          {"reasoning": "qwen2.5-coder:7b-base-q5_0"},
    "Go":         {"reasoning": "lfm2.5:8b-a1b-q4_K_M"},
    "Java":       {"reasoning": "qwen2.5-coder:7b-base-q5_0"},
    "Python":     {"reasoning": "lfm2.5:8b-a1b-q4_K_M"},
    "JavaScript": {"reasoning": "lfm2.5:8b-a1b-q4_K_M"},
    "TypeScript": {"reasoning": "lfm2.5:8b-a1b-q4_K_M"},
}
DEFAULT_CODING = "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL"

# Pick the canonical source language per category when the corpus contains
# multiple language subdirs. Categories not in this map fall back to the
# first non-Rust language encountered.
CATEGORY_DEFAULT_LANG = {
    "alphatrans": "Java",
    "oxidizer":   "Go",
}

# GildedRose-Refactoring-Kata is hand-curated; its config lives at
# .config/gildedrose.yml. Skip its regeneration so we don't overwrite
# the hand-tuned version.
MANUAL_PROJECT_NAMES = {"GildedRose-Refactoring-Kata"}


def iterations_for(file_count: int) -> tuple[int, int]:
    if file_count < 10: return 1, 7200
    if file_count < 50: return 5, 7200
    return 10, 7200


def build_cfg(cat: str, proj: str, src_lang: str, src_subdir: str, file_count: int) -> str:
    iters, timeout = iterations_for(file_count)
    models = MODELS[src_lang]
    return (
        "translation:\n"
        "  source:\n"
        f'    dir: "assets/samples/{cat}/{proj}/{src_subdir}"\n'
        f'    language: "{src_lang}"\n'
        "  target:\n"
        f'    dir: ".artifacts/{cat}/{proj}/rust"\n'
        '    language: "Rust"\n'
        '    toolchain: "cargo"\n'
        "  models:\n"
        f'    reasoning: "{models["reasoning"]}"\n'
        f'    coding: "{DEFAULT_CODING}"\n'
        '    ollama_url: "http://localhost:11434"\n'
        "  execution:\n"
        f"    max_iterations: {iters}\n"
        f"    timeout_seconds: {timeout}\n"
        "  lsp:\n"
        '    provider: "native"\n'
    )


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("--dry-run", action="store_true", help="print what would be written, write nothing")
    args = ap.parse_args()

    if not SAMPLES.exists():
        raise SystemExit(
            f"assets/samples/ corpus not found: {SAMPLES}\n"
            "  run: ls assets/samples/ to verify the corpus is checked out"
        )

    CONFIG_DIR.mkdir(exist_ok=True)

    # Treat any existing .config/<name>.yml as a manual config that should
    # suppress regeneration. Match by project name (case-insensitive) so
    # `gildedrose.yml` blocks regeneration of `GildedRose-Refactoring-Kata`.
    manual_stems = {p.stem.lower() for p in CONFIG_DIR.glob("*.yml")}
    manual_stems |= {name.lower() for name in MANUAL_PROJECT_NAMES}

    generated = skipped = 0
    for cat_dir in sorted(SAMPLES.iterdir()):
        if not cat_dir.is_dir():
            continue
        cat = cat_dir.name
        # Skip the rust/ reference-translation subdir if it sits at the
        # corpus top level (e.g. assets/samples/GildedRose-Refactoring-Kata/rust/).
        if cat.endswith(".rust") or cat == "rust":
            continue
        for proj_dir in sorted(cat_dir.iterdir()):
            if not proj_dir.is_dir():
                continue
            proj = proj_dir.name

            # GildedRose is special: its C source lives under
            # assets/manual-samples/ (split out of the corpus). The generated
            # config would point into an empty path, so we skip it entirely.
            if proj.lower() in manual_stems:
                skipped += 1
                continue

            langs = [
                (LANG_MAP[s.name.lower()], s, s.name.lower())
                for s in proj_dir.iterdir()
                if s.is_dir() and s.name.lower() in LANG_MAP
            ]
            sources = [t for t in langs if t[0] != "Rust"]
            if not sources:
                continue
            pref = CATEGORY_DEFAULT_LANG.get(cat)
            chosen = next((t for t in sources if t[0] == pref), sources[0])
            src_lang, src_path, src_subdir = chosen

            file_count = sum(
                1 for _, _, fs in os.walk(src_path) for f in fs if not f.startswith(".")
            )

            cfg = build_cfg(cat, proj, src_lang, src_subdir, file_count)
            out = CONFIG_DIR / f"{cat}-{proj}.yml"
            if args.dry_run:
                print(f"[dry-run] would write {out} (lang={src_lang}, files={file_count})")
            else:
                out.write_text(cfg)
            generated += 1

    action = "would generate" if args.dry_run else "generated"
    print(
        f"{action}={generated}, skipped(manual)={skipped}, "
        f"total in .config/={generated + len(manual_stems)}"
    )


if __name__ == "__main__":
    main()
