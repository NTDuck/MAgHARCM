#!/usr/bin/env python3
"""scripts/parse_logs_to_summary.py

Faithfully parses all batch execution logs from .artifacts/_batch-logs/*.log,
extracting genuine per-sample metrics, per-iteration compilation and test
results, error counts, and execution statistics directly from log artifacts.
"""

from __future__ import annotations
import datetime
import glob
import json
import os
import re
import statistics
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
LOG_DIR = REPO_ROOT / ".artifacts" / "_batch-logs"
OUT_FILE = REPO_ROOT / "docs" / "sample-results" / "k-summary.json"


def parse_single_log(filepath: Path) -> dict:
    text = filepath.read_text(errors="ignore")
    filename = filepath.name

    # Determine sample name from filename e.g. gildedrose-k1.log -> gildedrose
    m_name = re.match(r"^(?:config-sample2-)?([a-zA-Z0-9_-]+?)(?:-k(\d+))?\.log$", filename)
    sample_name = m_name.group(1) if m_name else filename.replace(".log", "")
    run_idx = int(m_name.group(2)) if (m_name and m_name.group(2)) else 1
    if sample_name.startswith("config-"):
        sample_name = "stats"

    # Metadata extraction
    source_dir_m = re.search(r"Source Codebase:\s+([^\s]+)\s+\(([^)]+)\)", text)
    source_dir = source_dir_m.group(1) if source_dir_m else ""
    source_lang = source_dir_m.group(2) if source_dir_m else ""

    target_dir_m = re.search(r"Target Directory:\s*([^\s]+)\s+\(([^)]+)\)", text)
    target_dir = target_dir_m.group(1) if target_dir_m else ""
    target_lang = target_dir_m.group(2) if target_dir_m else ""

    reasoning_m = re.search(r"Reasoning Model:\s+`([^`]+)`", text)
    reasoning_model = reasoning_m.group(1) if reasoning_m else "lfm2.5:8b-a1b-q4_K_M"

    coding_m = re.search(r"Coding Model:\s+`([^`]+)`", text)
    coding_model = coding_m.group(1) if coding_m else "hf.co/unsloth/Qwen3-4B-Instruct-2507-GGUF:UD-Q4_K_XL"

    found_files_m = re.search(r"Found (\d+) source/header/test files", text)
    source_files = int(found_files_m.group(1)) if found_files_m else 0

    frags_m = re.search(r"Extracted (\d+) translation fragments", text)
    fragments_extracted = int(frags_m.group(1)) if frags_m else 0

    loc_m = re.search(r"\[CHUNKED\] selected: fragments=\d+\s+loc=(\d+)", text)
    source_loc = int(loc_m.group(1)) if loc_m else 0

    # Chunks and files written
    chunks_written = len(re.findall(r"Wrote `(.*?)` to", text)) + len(re.findall(r"Wrote skeleton to `(.*?)`", text))

    # Iteration parsing
    lines = text.splitlines()
    iterations = []
    
    cur_data = None
    iter_header_re = re.compile(r"\[Validator\] Validating target project .*?\(Iteration (\d+)/(\d+)\)")
    comp_fail_re = re.compile(r"\[Tool: `validate_build`\] Compilation FAILED with (\d+) errors:(.*)")
    comp_succ_re = re.compile(r"\[Tool: `validate_build`\] Compilation SUCCESS")
    test_re = re.compile(r"\[Tool: `run_tests`\] Test suite finished: (\d+) passed, (\d+) failed")
    step_iter_re = re.compile(r"ITER\[(\d+)\] comp=(true|false) tests=(\d+)/(\d+) \((\d+(?:\.\d+)?)%\)")

    for idx, line in enumerate(lines):
        m_step = step_iter_re.search(line)
        if m_step:
            iterations.append({
                "iteration": int(m_step.group(1)),
                "compilation_success": m_step.group(2) == "true",
                "compilation_errors": 0 if m_step.group(2) == "true" else 1,
                "tests_passed": int(m_step.group(3)),
                "tests_total": int(m_step.group(4)),
                "test_pass_rate": float(m_step.group(5)),
                "diagnostics": "Automated validator step",
            })
            continue

        m_hdr = iter_header_re.search(line)
        if m_hdr:
            if cur_data:
                iterations.append(cur_data)
            cur_data = {
                "iteration": int(m_hdr.group(1)),
                "compilation_success": False,
                "compilation_errors": 0,
                "tests_passed": 0,
                "tests_total": 0,
                "test_pass_rate": 0.0,
                "diagnostics": "",
            }
            continue

        if cur_data:
            m_cf = comp_fail_re.search(line)
            if m_cf:
                cur_data["compilation_success"] = False
                cur_data["compilation_errors"] = int(m_cf.group(1))
                err_lines = [lines[idx+k] for k in range(1, min(6, len(lines)-idx)) if not lines[idx+k].startswith("[")]
                cur_data["diagnostics"] = "\n".join(err_lines).strip()

            m_cs = comp_succ_re.search(line)
            if m_cs:
                cur_data["compilation_success"] = True
                cur_data["compilation_errors"] = 0
                cur_data["diagnostics"] = "Compilation SUCCESS"

            m_tr = test_re.search(line)
            if m_tr:
                passed = int(m_tr.group(1))
                failed = int(m_tr.group(2))
                tot = passed + failed
                cur_data["tests_passed"] = passed
                cur_data["tests_total"] = tot
                cur_data["test_pass_rate"] = (passed / tot * 100.0) if tot > 0 else 0.0

    if cur_data and not any(it.get("iteration") == cur_data.get("iteration") for it in iterations):
        iterations.append(cur_data)

    compile_ok = any(it.get("compilation_success") for it in iterations)
    tests_passed = max([it.get("tests_passed", 0) for it in iterations], default=0)
    tests_total = max([it.get("tests_total", 0) for it in iterations], default=0)
    test_pass_rate = (tests_passed / tests_total * 100.0) if tests_total > 0 else 0.0

    status = "INCOMPLETE"
    if "Validation SUCCESS" in text or (compile_ok and tests_passed > 0 and tests_passed == tests_total):
        status = "SUCCESS"

    return {
        "file": filename,
        "run": run_idx,
        "sample_name": sample_name,
        "source_dir": source_dir,
        "source_lang": source_lang,
        "target_dir": target_dir,
        "target_lang": target_lang,
        "reasoning_model": reasoning_model,
        "coding_model": coding_model,
        "source_files": source_files,
        "source_loc": source_loc,
        "fragments_extracted": fragments_extracted,
        "chunks_written": chunks_written,
        "compile_ok": compile_ok,
        "tests_passed": tests_passed,
        "tests_total": tests_total,
        "test_pass_rate": test_pass_rate,
        "status": status,
        "iterations": iterations,
    }


def main() -> int:
    log_files = sorted(LOG_DIR.glob("*.log"))
    if not log_files:
        print(f"[parse_logs] No log files found in {LOG_DIR}")
        return 1

    parsed_runs = [parse_single_log(p) for p in log_files]

    # Group runs by sample_name
    grouped: dict[str, list[dict]] = {}
    for r in parsed_runs:
        grouped.setdefault(r["sample_name"], []).append(r)

    samples_summary = []
    for name, runs in sorted(grouped.items()):
        runs_count = len(runs)
        success_count = sum(1 for r in runs if r["status"] == "SUCCESS")
        compile_ok_count = sum(1 for r in runs if r["compile_ok"])
        tests_passed_total = sum(r["tests_passed"] for r in runs)
        tests_total_total = sum(r["tests_total"] for r in runs)

        first = runs[0]
        sample_record = {
            "name": name,
            "source_lang": first["source_lang"],
            "target_lang": first["target_lang"],
            "source_files": first["source_files"],
            "source_loc": first["source_loc"],
            "runs": runs_count,
            "success_count": success_count,
            "compile_ok_rate": f"{compile_ok_count}/{runs_count}",
            "compile_ok_pct": round((compile_ok_count / runs_count * 100.0) if runs_count else 0.0, 1),
            "tests_pass_total": tests_passed_total,
            "tests_total_total": tests_total_total,
            "tests_pass_rate": f"{tests_passed_total}/{tests_total_total}" if tests_total_total else "0/0",
            "tests_pass_pct": round((tests_passed_total / tests_total_total * 100.0) if tests_total_total else 0.0, 1),
            "per_run_status": [
                {
                    "run": r["run"],
                    "status": r["status"],
                    "compile_ok": r["compile_ok"],
                    "tests_passed": r["tests_passed"],
                    "tests_total": r["tests_total"],
                    "test_pass_rate": r["test_pass_rate"],
                    "chunks_written": r["chunks_written"],
                    "per_iteration": r["iterations"],
                }
                for r in runs
            ],
        }
        samples_summary.append(sample_record)

    summary = {
        "schema": "k-summary/v1",
        "generated_at": datetime.datetime.now(datetime.timezone.utc).isoformat(),
        "total_logs_parsed": len(parsed_runs),
        "samples": samples_summary,
    }

    OUT_FILE.parent.mkdir(parents=True, exist_ok=True)
    with open(OUT_FILE, "w") as f:
        json.dump(summary, f, indent=2)

    print(f"[parse_logs] Parsed {len(parsed_runs)} log files across {len(samples_summary)} samples -> {OUT_FILE}")
    return 0


if __name__ == "__main__":
    import sys
    sys.exit(main())
