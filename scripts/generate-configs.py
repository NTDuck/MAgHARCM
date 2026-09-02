#!/usr/bin/env python3
"""Generate .config/<category>-<project>.yml for every assets/samples project.

Detects source language from the named subdir (c/, go/, java/, python/,
javascript/, typescript/) under each project, picks the per-category
canonical source if multiple are present, and emits a translation
request that targets Rust.

Iterates over:
  assets/samples/data/tool_projects/<category>/<project>/<lang>/

Skips projects that already have a manual config in .config/ (e.g.
.gildedrose.yml, .stats.yml) — these use hand-tuned settings.

Output: 116 generated configs + 4 manual = 120 total in .config/.

Usage:
  python3 scripts/generate-configs.py
"""
from __future__ import annotations
import json
import os
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent
SAMPLES = REPO / "assets" / "samples" / "data" / "tool_projects"
CONFIG_DIR = REPO / ".config"

LANG_MAP = {
    "c": "C", "go": "Go", "java": "Java", "python": "Python",
    "javascript": "JavaScript", "typescript": "TypeScript", "rust": "Rust",
}

MODELS = {
    "C":     {"reasoning": "qwen2.5-coder:7b-base-q5_0"},
    "Go":    {"reasoning": "lfm2.5:8b-a1b-q4_K_M"},
    "Java":  {"reasoning": "qwen2.5-coder:7b-base-q5_0"},
    "Python":{"reasoning": "lfm2.5:8b-a1b-q4_K_M"},
    "JavaScript":{"reasoning": "lfm2.5:8b-a1b-q4_K_M"},
    "TypeScript":{"reasoning": "lfm2.5:8b-a1b-q4_K_M"},
}
DEFAULT_CODING = "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL"

CATEGORY_DEFAULT_LANG = {
    "alphatrans": "Java", "oxidizer": "Go",
    "crust": "C",         "skel": "Python",
}

def iterations_for(file_count: int) -> tuple[int, int]:
    if file_count < 10:  return 1,  7200
    if file_count < 50:  return 5,  7200
    return 10, 7200

def main() -> None:
    if not SAMPLES.exists():
        raise SystemExit(f"assets/samples submodule not found: {SAMPLES}\n"
                         "  run: git submodule update --init")
    CONFIG_DIR.mkdir(exist_ok=True)
    manual = {p.stem for p in CONFIG_DIR.glob("*.yml")}
    generated = skipped = 0
    for cat_dir in sorted(SAMPLES.iterdir()):
        if not cat_dir.is_dir(): continue
        cat = cat_dir.name
        for proj in sorted(cat_dir.iterdir()):
            if not proj.is_dir(): continue
            langs = [(LANG_MAP[s.name.lower()], s, s.name.lower())
                     for s in proj.iterdir() if s.is_dir() and s.name.lower() in LANG_MAP]
            sources = [t for t in langs if t[0] != "Rust"]
            if not sources: continue
            pref = CATEGORY_DEFAULT_LANG.get(cat)
            chosen = next((t for t in sources if t[0] == pref), sources[0])
            if proj.name in manual:
                skipped += 1; continue
            src_lang, src_path, src_subdir = chosen
            file_count = sum(1 for _,_,fs in os.walk(src_path) for f in fs if not f.startswith('.'))
            iters, timeout = iterations_for(file_count)
            models = MODELS[src_lang]
            cfg = f"""translation:
  source:
    dir: "assets/samples/data/tool_projects/{cat}/{proj.name}/{src_subdir}"
    language: "{src_lang}"
  target:
    dir: ".artifacts/{cat}/{proj.name}/rust"
    language: "Rust"
    toolchain: "cargo"
  models:
    reasoning: "{models['reasoning']}"
    coding: "{DEFAULT_CODING}"
    ollama_url: "http://localhost:11434"
  execution:
    max_iterations: {iters}
    timeout_seconds: {timeout}
  lsp:
    provider: "native"
"""
            (CONFIG_DIR / f"{cat}-{proj.name}.yml").write_text(cfg)
            generated += 1
    print(f"generated={generated}, skipped(manual)={skipped}, total in .config/={generated + len(manual)}")

if __name__ == "__main__":
    main()
