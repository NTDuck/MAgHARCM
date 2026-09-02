#!/usr/bin/env bash
# scripts/run-samples-k.sh
#
# Run every sample .config/<category>-<project>.yml in sequence, K
# times each, capturing per-run outcome and aggregating per-sample
# metrics into a single JSON file (docs/sample-results/k-summary.json).
#
# Purpose: stability measurement. One-shot runs of a stochastic LLM
# pipeline tell you nothing about variance. Re-running each sample K
# times and aggregating (success rate, compile rate, test-pass rate,
# wall-clock mean + stdev) gives a stable signal for paper claims and
# drift detection.
#
# Usage:
#   bash scripts/run-samples-k.sh                       # all samples, K=3
#   K=5 bash scripts/run-samples-k.sh                   # all samples, K=5
#   bash scripts/run-samples-k.sh gohistogram            # substring filter
#   K=2 bash scripts/run-samples-k.sh stats gohistogram
#
# Skip-already-run:
#   Touch docs/sample-results/.skip-<name> before running to mark
#   a sample as already-exercised; the script will skip every run of it.
#
# Per-iteration wall-clock budget:
#   Each sample's .config/*.yml has timeout_seconds set. The script
#   uses `timeout` with 1.5x the YAML budget as a hard ceiling per
#   individual run.
#
# Output JSON schema (k-summary/v1, written atomically to
# docs/sample-results/k-summary.json):
#   {
#     "schema": "k-summary/v1",
#     "k": <int>,
#     "generated_at": "<RFC3339>",
#     "samples": [
#       {
#         "name": "<category>-<project>",
#         "config": "/abs/path/to/.config/<name>.yml",
#         "runs": <int>,
#         "success_count": <int>,         // status == SUCCESS
#         "compile_ok_rate": "<n>/<runs>",
#         "tests_pass_total": <int>,
#         "tests_total_total": <int>,
#         "tests_pass_rate": "<n>/<denom>",
#         "durations_ms": [<int>, ...],   // wall_ms per run
#         "mean_ms": <int>,
#         "stdev_ms": <int>,
#         "min_ms": <int>,
#         "max_ms": <int>,
#         "per_run_status": [
#           {"i": 1, "status": "...", "compile_ok": bool,
#            "tests_passed": int, "tests_total": int,
#            "duration_ms": int, "exit_code": int}
#         ]
#       }
#     ]
#   }
#
# Behaves like the harness: exits 0 when finished regardless of per-run
# outcome. Designed for unattended batches.

set -uo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
RESULTS_DIR="$REPO_ROOT/docs/sample-results"
SUMMARY_FILE="$RESULTS_DIR/k-summary.json"
LOG_DIR="$REPO_ROOT/.artifacts/_batch-logs"
mkdir -p "$RESULTS_DIR" "$LOG_DIR"

cd "$REPO_ROOT"

K="${K:-3}"
if ! [[ "$K" =~ ^[0-9]+$ ]] || [[ "$K" -lt 1 ]]; then
    echo "[k-run] K must be a positive integer (got: $K)" >&2
    exit 2
fi

# CONFIG_DIR overrides the directory the script scans for *.yml files
# (default: .config/). Used by smoke tests that need a low-timeout
# config without mutating the production .config/ tree.
CONFIG_DIR="${CONFIG_DIR:-$REPO_ROOT/.config}"

# Build once; reuse the binary for all samples.
echo "[k-run] building magharcm-cli..."
if ! go build -o "$REPO_ROOT/bin/MAgHARCM" ./cmd/MAgHARCM; then
    echo "[k-run] build failed" >&2
    exit 1
fi
echo "[k-run] built $(stat -c %s "$REPO_ROOT/bin/MAgHARCM") bytes"

shopt -s nullglob
all_configs=("$CONFIG_DIR"/*.yml)
shopt -u nullglob

# Drop meta-configs (default.yml is the flag default for ad-hoc runs,
# not a sample to benchmark). Sample configs are everything else; the
# exclusion list is explicit so single-word sample names like
# `stats.yml` or `gildedrose.yml` are NOT skipped.
META_CONFIGS=(default)
configs=()
for cfg in "${all_configs[@]}"; do
    name="$(basename "$cfg" .yml)"
    skip=0
    for meta in "${META_CONFIGS[@]}"; do
        if [[ "$name" == "$meta" ]]; then skip=1; break; fi
    done
    if [[ $skip -eq 1 ]]; then
        echo "[skip ] $name (meta-config, not a sample)"
    else
        configs+=("$cfg")
    fi
done

# If user passed positional args, treat each as a substring filter
# against the config basename (sans .yml).
if [[ $# -gt 0 ]]; then
    filtered=()
    for cfg in "${configs[@]}"; do
        name="$(basename "$cfg" .yml)"
        for pat in "$@"; do
            if [[ "$name" == *"$pat"* ]]; then
                filtered+=("$cfg")
                break
            fi
        done
    done
    if [[ ${#filtered[@]} -eq 0 ]]; then
        echo "[k-run] no configs match any of: $*"
        exit 0
    fi
    configs=("${filtered[@]}")
fi

total=${#configs[@]}
echo "[k-run] $total configs × K=$K runs = $((total * K)) total runs"
echo ""

# triage_log extracts compile_ok / tests_pass / status / exit_code
# from a harness log. Mirrors the parsing in run-all-samples.sh but
# returns tab-separated values for reliable python-side aggregation.
triage_log() {
    local log="$1"
    local compile_ok="false" tests_pass="0" tests_total="0" status="INCOMPLETE"
    if grep -q 'compilation=true' "$log" || grep -q 'Compilation SUCCESS' "$log"; then compile_ok="true"; fi
    if grep -qE 'passed=[0-9]+/[0-9]+' "$log"; then
        tests_pass=$(grep -oE 'passed=[0-9]+' "$log" | tail -1 | cut -d= -f2 | cut -d/ -f1)
        tests_total=$(grep -oE 'passed=[0-9]+/[0-9]+' "$log" | tail -1 | cut -d= -f2 | cut -d/ -f2)
    elif grep -qE 'tests=[0-9]+/[0-9]+' "$log"; then
        tests_pass=$(grep -oE 'tests=[0-9]+' "$log" | tail -1 | cut -d= -f2 | cut -d/ -f1)
        tests_total=$(grep -oE 'tests=[0-9]+/[0-9]+' "$log" | tail -1 | cut -d= -f2 | cut -d/ -f2)
    fi
    if grep -q 'Validation INCOMPLETE' "$log"; then status="INCOMPLETE"; fi
    if grep -q 'all_success=true' "$log" || grep -q 'Validation SUCCESS' "$log"; then status="SUCCESS"; fi
    printf '%s\t%s\t%s\t%s\n' "$compile_ok" "$tests_pass" "$tests_total" "$status"
}

# Build the samples array as a python-friendly JSON file by streaming
# each run's result through a small embedded python helper that owns
# the per-sample aggregation (mean, stdev, rates).
SAMPLES_JSON="$RESULTS_DIR/.k-samples-$$.json"
echo '[]' > "$SAMPLES_JSON"

# trim_summary is invoked once at the end to finalise the JSON file
# (add schema, k, generated_at) atomically via a tempfile rename.

ran=0; skipped=0
total_runs=$((total * K))
for cfg in "${configs[@]}"; do
    name="$(basename "$cfg" .yml)"

    # Honour skip marker (skips every K runs of the sample).
    if [[ -f "$RESULTS_DIR/.skip-$name" ]]; then
        echo "[skip ] $name (skip marker present; skipping all $K runs)"
        skipped=$((skipped + 1))
        continue
    fi

    echo "[run  ] $name × $K"

    # Per-sample run loop.
    per_run_status_py="[]"
    durations_py="[]"
    success_count=0
    compile_ok_count=0
    tests_pass_total=0
    tests_total_total=0

    yaml_timeout=$(grep -E '^\s*timeout_seconds:' "$cfg" | awk '{print $2}')
    yaml_timeout=${yaml_timeout:-7200}
    wall=$((yaml_timeout * 3 / 2))

    for ((i=1; i<=K; i++)); do
        log="$LOG_DIR/${name}-k${i}.log"
        echo "       iter $i/$K (timeout=${wall}s) log=$log"
        start=$(date +%s%N)
        timeout --foreground "$wall" "$REPO_ROOT/bin/MAgHARCM" --config "$cfg" \
            >"$log" 2>&1
        rc=$?
        end=$(date +%s%N)
        duration_ms=$(( (end - start) / 1000000 ))

        read -r compile_ok tests_pass tests_total status < <(triage_log "$log")

        if [[ "$status" == "SUCCESS" ]]; then
            success_count=$((success_count + 1))
        fi
        if [[ "$compile_ok" == "true" ]]; then
            compile_ok_count=$((compile_ok_count + 1))
        fi
        tests_pass_total=$((tests_pass_total + tests_pass))
        tests_total_total=$((tests_total_total + tests_total))
        ran=$((ran + 1))

        echo "       -> status=$status compile=$compile_ok tests=$tests_pass/$tests_total rc=$rc dur=${duration_ms}ms"

        # Append this run to the per-sample python lists.
        per_run_status_py=$(python3 -c "
import json
arr = json.loads('''$per_run_status_py''')
arr.append({
    'i': $i,
    'status': '$status',
    'compile_ok': '$compile_ok' == 'true',
    'tests_passed': int('$tests_pass'),
    'tests_total': int('$tests_total'),
    'duration_ms': $duration_ms,
    'exit_code': $rc,
})
print(json.dumps(arr))
")
        durations_py=$(python3 -c "
import json
arr = json.loads('''$durations_py''')
arr.append($duration_ms)
print(json.dumps(arr))
")
    done

    # Append this sample's aggregate record to the samples file.
    python3 -c "
import json
with open('$SAMPLES_JSON') as f:
    samples = json.load(f)
samples.append({
    'name': '$name',
    'config': '$cfg',
    'runs': $K,
    'success_count': $success_count,
    'compile_ok_rate': '$compile_ok_count/$K',
    'tests_pass_total': $tests_pass_total,
    'tests_total_total': $tests_total_total,
    'tests_pass_rate': '$tests_pass_total/$tests_total_total' if $tests_total_total else '0/0',
    'durations_ms': json.loads('''$durations_py'''),
    'per_run_status': json.loads('''$per_run_status_py'''),
})
with open('$SAMPLES_JSON', 'w') as f:
    json.dump(samples, f, indent=2)
"
done

# Finalise: compute mean/stdev/min/max and wrap with schema metadata.
python3 -c "
import json, statistics, datetime, os, sys

with open('$SAMPLES_JSON') as f:
    samples = json.load(f)
os.remove('$SAMPLES_JSON')

for s in samples:
    d = s['durations_ms']
    s['mean_ms'] = int(statistics.fmean(d)) if d else 0
    s['stdev_ms'] = int(statistics.pstdev(d)) if len(d) > 1 else 0
    s['min_ms'] = min(d) if d else 0
    s['max_ms'] = max(d) if d else 0

summary = {
    'schema': 'k-summary/v1',
    'k': $K,
    'generated_at': datetime.datetime.now(datetime.timezone.utc).isoformat(),
    'samples': samples,
}

# Atomic write: tempfile + rename, so a kill mid-write leaves the
# previous summary intact rather than corrupting it.
tmp = '$SUMMARY_FILE' + '.tmp'
with open(tmp, 'w') as f:
    json.dump(summary, f, indent=2)
os.replace(tmp, '$SUMMARY_FILE')
"

echo ""
echo "[k-run] done: ran=$ran runs (skipped=$skipped configs)"
echo "[k-run] summary: $SUMMARY_FILE"
