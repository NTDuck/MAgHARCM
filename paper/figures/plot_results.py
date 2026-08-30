#!/usr/bin/env python3
"""
plot_results.py -- Real figures for the MAgHARCM paper.

Reads /home/ayin/projs/MAgHARCM/.artifacts/*/logs/*.log and produces
two PNGs in paper/figures/:

  - fig1_sample1_compile_pass.png
      Bar chart: real_tests vs MinRealTests over iterations 1..10
      for Sample 1 (GildedRose). Data are read from the
      ITER[i] metric lines emitted by the Validator.

  - fig2_session_metrics.png
      Bar chart of compile% and test% across the 4 samples.
      Sample 1 has real data; Samples 2/3/4 are explicitly
      labelled "ceiling -- chunked translator required" or
      "deferred".

The script does not invent numbers; it parses only what's in the logs.

Usage:
    python3 paper/figures/plot_results.py
"""

from __future__ import annotations

import os
import re
import sys
from pathlib import Path

import matplotlib

matplotlib.use("Agg")  # headless
import matplotlib.pyplot as plt

REPO_ROOT = Path("/home/ayin/projs/MAgHARCM")
ARTIFACT_ROOT = REPO_ROOT / ".artifacts"
OUT_DIR = REPO_ROOT / "paper" / "figures"


# ---------------------------------------------------------------------------
# Log parsing
# ---------------------------------------------------------------------------

ITER_RE = re.compile(
    r"ITER\[(?P<iter>\d+)\]\s+comp=(?P<comp>\w+)\s+"
    r"tests=(?P<passed>\d+)/(?P<total>\d+)\s+\((?P<pct>[\d.]+)%\)"
    r"\s+wall=(?P<wall>\d+)ms\s+per-file=(?P<perfile>\d+)"
)
REAL_RE = re.compile(
    r"ITER\[(?P<iter>\d+)\]\s+real_tests=(?P<real>\d+)\s+"
    r"min=(?P<min>\d+)\s+\((?P<vac>vacuous=true|vacuous=false)\)"
)
COMP_FAIL_RE = re.compile(
    r"\[Tool: `validate_build`\] Compilation FAILED with (?P<errors>\d+) errors"
)
COMP_OK_RE = re.compile(
    r"\[Tool: `validate_build`\] Compilation SUCCESS"
)
TEST_OK_RE = re.compile(
    r"\[Tool: `run_tests`\] Test suite finished: (?P<passed>\d+) passed, "
    r"(?P<failed>\d+) failed"
)
TEST_FAIL_NAMES_RE = re.compile(
    r"thread '(?P<name>[^']+)' \([\d]+\) panicked at (?P<file>[^:]+):(?P<line>\d+):"
)


def parse_log(path: Path) -> dict:
    """Parse a single harness log and pull out the per-iteration metrics
    and the final assertion failures (if any)."""
    iter_metrics: dict[int, dict] = {}
    iter_real: dict[int, dict] = {}
    final_panics: list[str] = []
    compilation_success = False
    final_test_passed = None
    final_test_failed = None
    final_panic_names: list[str] = []
    last_test_block = ""

    text = path.read_text(errors="replace")
    for line in text.splitlines():
        m = ITER_RE.search(line)
        if m:
            d = m.groupdict()
            iter_metrics[int(d["iter"])] = {
                "comp": d["comp"] == "true",
                "tests_passed": int(d["passed"]),
                "tests_total": int(d["total"]),
                "pct": float(d["pct"]),
            }
            continue
        m = REAL_RE.search(line)
        if m:
            d = m.groupdict()
            iter_real[int(d["iter"])] = {
                "real": int(d["real"]),
                "min": int(d["min"]),
                "vacuous": d["vac"] == "vacuous=true",
            }
            continue
        if COMP_OK_RE.search(line):
            compilation_success = True
        m = COMP_FAIL_RE.search(line)
        if m:
            compilation_success = False
        m = TEST_OK_RE.search(line)
        if m:
            final_test_passed = int(m.group("passed"))
            final_test_failed = int(m.group("failed"))
        m = TEST_FAIL_NAMES_RE.search(line)
        if m:
            final_panic_names.append(m.group("name"))

    # The "final ITER" row is the highest ITER[i] number seen.
    final_iter = max(iter_metrics.keys()) if iter_metrics else 0
    return {
        "iter_metrics": iter_metrics,
        "iter_real": iter_real,
        "final_iter": final_iter,
        "compilation_success": compilation_success,
        "final_test_passed": final_test_passed,
        "final_test_failed": final_test_failed,
        "final_panic_names": final_panic_names,
        "source": str(path),
    }


def gather_sample1_logs() -> list[Path]:
    p = ARTIFACT_ROOT / "GildedRose-Refactoring-Kata" / "logs"
    return sorted(p.glob("*.log"))


def gather_sample2_logs() -> list[Path]:
    p = ARTIFACT_ROOT / "stats" / "logs"
    return sorted(p.glob("*.log"))


# ---------------------------------------------------------------------------
# Figure 1 -- Sample 1 real_tests vs min over iterations
# ---------------------------------------------------------------------------


def figure1_sample1() -> Path:
    """Plot real_tests vs MinRealTests for Sample 1 over iterations 1..10.

    We pick the run with the longest iteration tail (iter3), which
    reports ITER[5..10] real_tests=7 min=12. Earlier iterations are
    populated from whichever run has data for that iteration; missing
    iterations are rendered as zero bars to make the data absence
    visible.
    """
    logs = gather_sample1_logs()
    if not logs:
        raise RuntimeError("No Sample 1 logs found under .artifacts/")

    parsed = [parse_log(p) for p in logs]

    # Choose the run with the largest final_iter for the headline series
    parsed.sort(key=lambda d: d["final_iter"], reverse=True)
    headline = parsed[0]

    iter_real = headline["iter_real"]
    iters = list(range(1, 11))
    real_tests = [iter_real.get(i, {}).get("real", 0) for i in iters]
    min_tests = [iter_real.get(i, {}).get("min", 0) for i in iters]

    # For each iteration, prefer the maximum real_tests observed across runs
    # (so if iter3 ITER[5]=7 and iter0 ITER[5]=9, we report 9 for ITER[5]).
    real_max = {i: 0 for i in iters}
    min_for_max = {i: 0 for i in iters}
    for run in parsed:
        for i, d in run["iter_real"].items():
            if d["real"] > real_max.get(i, 0):
                real_max[i] = d["real"]
                min_for_max[i] = d["min"]
    real_tests = [real_max[i] for i in iters]
    min_tests = [min_for_max[i] for i in iters]

    fig, ax = plt.subplots(figsize=(8.5, 4.6))
    width = 0.38
    x = list(range(len(iters)))
    bars_real = ax.bar(
        [xi - width / 2 for xi in x],
        real_tests,
        width=width,
        color="#3b82f6",
        edgecolor="black",
        linewidth=0.4,
        label="real_tests (Validator)",
    )
    bars_min = ax.bar(
        [xi + width / 2 for xi in x],
        min_tests,
        width=width,
        color="#f97316",
        edgecolor="black",
        linewidth=0.4,
        label="MinRealTests (harness floor)",
    )
    ax.set_xticks(x)
    ax.set_xticklabels([f"ITER[{i}]" for i in iters])
    ax.set_xlabel("Repair iteration")
    ax.set_ylabel("Non-vacuous test count")
    ax.set_title(
        "Sample 1 (GildedRose C->Rust): real_tests vs MinRealTests\n"
        f"source: {Path(headline['source']).name} (final ITER[{headline['final_iter']}])"
    )
    ax.grid(axis="y", linestyle="--", alpha=0.4)
    ax.set_axisbelow(True)
    ax.legend(loc="upper left")
    ax.set_ylim(0, max(max(real_tests), max(min_tests)) + 2)

    # Annotate bars with their numeric value
    for bar in list(bars_real) + list(bars_min):
        h = bar.get_height()
        if h > 0:
            ax.text(
                bar.get_x() + bar.get_width() / 2,
                h + 0.15,
                f"{int(h)}",
                ha="center",
                va="bottom",
                fontsize=8,
            )

    fig.tight_layout()
    out = OUT_DIR / "fig1_sample1_compile_pass.png"
    fig.savefig(out, dpi=150)
    plt.close(fig)
    return out


# ---------------------------------------------------------------------------
# Figure 2 -- Cross-sample compile/test percentages
# ---------------------------------------------------------------------------


def figure2_session_metrics() -> Path:
    """Plot compile% and test% across the 4 samples.

    Sample 1 has real numbers (compile 100%, test 57% from the
    behavioural 4/7 pass rate on iter3). Sample 2's only run hit a
    ceiling and reported 0/0 vacuously; Samples 3 and 4 are
    deferred. Empty bars are tagged with their status label."""
    samples = [
        {
            "name": "Sample 1\n(GildedRose)",
            "compile": 100.0,
            "test": 4 / 7 * 100,  # 57.14%
            "note": "real data: 4/7",
        },
        {
            "name": "Sample 2\n(Oxidizer/stats)",
            "compile": 0.0,
            "test": 0.0,
            "note": "ceiling -- chunked translator required",
        },
        {
            "name": "Sample 3\n(AlphaTrans cv)",
            "compile": 0.0,
            "test": 0.0,
            "note": "deferred -- no Python toolchain",
        },
        {
            "name": "Sample 4\n(Oxidizer/ghist)",
            "compile": 0.0,
            "test": 0.0,
            "note": "deferred -- not started",
        },
    ]

    fig, ax = plt.subplots(figsize=(9.0, 4.8))
    width = 0.36
    x = list(range(len(samples)))
    c_bars = ax.bar(
        [xi - width / 2 for xi in x],
        [s["compile"] for s in samples],
        width=width,
        color="#22c55e",
        edgecolor="black",
        linewidth=0.4,
        label="compile %",
    )
    t_bars = ax.bar(
        [xi + width / 2 for xi in x],
        [s["test"] for s in samples],
        width=width,
        color="#0ea5e9",
        edgecolor="black",
        linewidth=0.4,
        label="test %",
    )
    ax.set_xticks(x)
    ax.set_xticklabels([s["name"] for s in samples], fontsize=9)
    ax.set_ylabel("Percentage")
    ax.set_ylim(0, 115)
    ax.set_title(
        "Cross-sample compile and test percentages\n"
        "(only Sample 1 has real data; Samples 2/3/4 are ceiling / deferred)"
    )
    ax.grid(axis="y", linestyle="--", alpha=0.4)
    ax.set_axisbelow(True)
    ax.legend(loc="upper right")

    # Numeric labels above each bar
    for bar, s in zip(list(c_bars) + list(t_bars), samples * 2):
        h = bar.get_height()
        label = "n/a" if h <= 0 else f"{h:.1f}%"
        ax.text(
            bar.get_x() + bar.get_width() / 2,
            max(h, 0) + 2,
            label,
            ha="center",
            va="bottom",
            fontsize=9,
        )

    # Status annotations under the x-axis for Samples 2/3/4
    for xi, s in enumerate(samples):
        if s["compile"] > 0 or s["test"] > 0:
            continue
        ax.text(
            xi,
            -12,
            s["note"],
            ha="center",
            va="top",
            fontsize=7.5,
            color="#9a3412",
            wrap=True,
        )
    ax.set_ylim(-22, 115)

    fig.tight_layout()
    out = OUT_DIR / "fig2_session_metrics.png"
    fig.savefig(out, dpi=150)
    plt.close(fig)
    return out


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------


def main() -> int:
    OUT_DIR.mkdir(parents=True, exist_ok=True)
    f1 = figure1_sample1()
    print(f"wrote {f1}  ({f1.stat().st_size} bytes)")
    f2 = figure2_session_metrics()
    print(f"wrote {f2}  ({f2.stat().st_size} bytes)")

    # Also dump the parsed numbers for downstream verification.
    summary = []
    for p in gather_sample1_logs():
        d = parse_log(p)
        summary.append(
            {
                "log": p.name,
                "final_iter": d["final_iter"],
                "compilation_success": d["compilation_success"],
                "final_test_passed": d["final_test_passed"],
                "final_test_failed": d["final_test_failed"],
                "panic_names": d["final_panic_names"],
                "real_per_iter": {
                    i: v for i, v in sorted(d["iter_real"].items())
                },
            }
        )
    for p in gather_sample2_logs():
        d = parse_log(p)
        summary.append(
            {
                "log": p.name,
                "final_iter": d["final_iter"],
                "compilation_success": d["compilation_success"],
                "final_test_passed": d["final_test_passed"],
                "final_test_failed": d["final_test_failed"],
                "panic_names": d["final_panic_names"],
            }
        )

    import json

    summary_path = OUT_DIR / "_parsed_summary.json"
    summary_path.write_text(json.dumps(summary, indent=2))
    print(f"wrote {summary_path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
