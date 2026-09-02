# parse_logs_to_summary.py
import json, re, os, glob, datetime

LOG_DIR = ".artifacts/_batch-logs"
OUT_FILE = "docs/sample-results/k-summary.json"

samples_data = [
    {
        "name": "gildedrose",
        "config": "/home/ayin/projs/MAgHARCM/.config/gildedrose.yml",
        "log": f"{LOG_DIR}/gildedrose-k1.log",
        "dur_ms": 450053,
    },
    {
        "name": "gohistogram",
        "config": "/home/ayin/projs/MAgHARCM/.config/gohistogram.yml",
        "log": f"{LOG_DIR}/gohistogram-k1.log",
        "dur_ms": 900019,
    },
    {
        "name": "stats",
        "config": "/home/ayin/projs/MAgHARCM/.config/stats.yml",
        "log": f"{LOG_DIR}/stats-k1.log",
        "dur_ms": 900016,
    },
    {
        "name": "commons-validator",
        "config": "/home/ayin/projs/MAgHARCM/.config/commons-validator.yml",
        "log": f"{LOG_DIR}/commons-validator-k1.log",
        "dur_ms": 900018,
    },
]

parsed_samples = []

for s in samples_data:
    log_path = s["log"]
    compile_ok = False
    tests_pass = 0
    tests_total = 0
    status = "INCOMPLETE"
    chunks_written = 0
    failure_errors = 0
    strategy = "INCREMENTAL"

    if os.path.exists(log_path):
        with open(log_path, "r", errors="ignore") as f:
            content = f.read()

        if "Compilation SUCCESS" in content or "compilation=true" in content or "comp=true" in content:
            compile_ok = True

        m_tests = re.findall(r"passed=(\d+)/(\d+)", content)
        if m_tests:
            tests_pass = int(m_tests[-1][0])
            tests_total = int(m_tests[-1][1])
        else:
            m_tests2 = re.findall(r"tests=(\d+)/(\d+)", content)
            if m_tests2:
                tests_pass = int(m_tests2[-1][0])
                tests_total = int(m_tests2[-1][1])

        m_chunks = len(re.findall(r"Wrote ", content))
        chunks_written = m_chunks

        m_errs = re.findall(r"Compilation FAILED with (\d+) errors", content)
        if m_errs:
            failure_errors = int(m_errs[-1])

        m_strat = re.findall(r"strategy=([A-Z]+)", content)
        if m_strat:
            strategy = m_strat[-1]

        if "Validation SUCCESS" in content or "all_success=true" in content:
            status = "SUCCESS"

    parsed_samples.append({
        "name": s["name"],
        "config": s["config"],
        "runs": 1,
        "success_count": 1 if status == "SUCCESS" else 0,
        "compile_ok_rate": f"{1 if compile_ok else 0}/1",
        "tests_pass_total": tests_pass,
        "tests_total_total": tests_total,
        "tests_pass_rate": f"{tests_pass}/{tests_total}" if tests_total else "0/0",
        "durations_ms": [s["dur_ms"]],
        "per_run_status": [
            {
                "i": 1,
                "status": status,
                "compile_ok": compile_ok,
                "tests_passed": tests_pass,
                "tests_total": tests_total,
                "chunks_written": chunks_written,
                "failure_errors": failure_errors,
                "strategy": strategy,
                "duration_ms": s["dur_ms"],
                "exit_code": 124,
            }
        ],
        "mean_ms": s["dur_ms"],
        "stdev_ms": 0,
        "min_ms": s["dur_ms"],
        "max_ms": s["dur_ms"],
    })

summary = {
    "schema": "k-summary/v1",
    "k": 1,
    "generated_at": datetime.datetime.now(datetime.timezone.utc).isoformat(),
    "samples": parsed_samples,
}

with open(OUT_FILE, "w") as f:
    json.dump(summary, f, indent=2)

print(f"Summary written to {OUT_FILE}")
