#!/usr/bin/env bash
# scripts/run-all-samples.sh
#
# Run every .config/<category>-<project>.yml in sequence, capturing
# per-sample outcome to docs/sample-results/<category>-<project>.md
# and a JSON index to docs/sample-results/index.json.
#
# Behaves like the harness: exits 0 when finished regardless of
# per-sample outcome (each run logs its own [WARNING] Validation
# INCOMPLETE line and a non-zero exit code that we capture but do not
# propagate). Designed for overnight unattended batches.
#
# Usage:
#   bash scripts/run-all-samples.sh                    # all configs
#   bash scripts/run-all-samples.sh crust-*            # glob filter
#   bash scripts/run-all-samples.sh oxidizer-gohistogram   # one specific
#
# Skip-already-run:
#   Touch docs/sample-results/.skip-<name> before running to mark
#   a sample as already-exercised; the script will skip it.
#
# Per-iteration wall-clock budget:
#   Each sample's .config/*.yml has timeout_seconds set. The script
#   uses `timeout` with 1.5x the YAML budget as a hard ceiling.

set -uo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
RESULTS_DIR="$REPO_ROOT/docs/sample-results"
INDEX_FILE="$RESULTS_DIR/index.json"
LOG_DIR="$REPO_ROOT/.artifacts/_batch-logs"
mkdir -p "$RESULTS_DIR" "$LOG_DIR"

cd "$REPO_ROOT"

# Build once; reuse the binary for all samples.
echo "[batch] building magharcm-cli..."
go build -o "$REPO_ROOT/bin/MAgHARCM" ./cmd/MAgHARCM || exit 1
echo "[batch] built $(ls -la "$REPO_ROOT/bin/MAgHARCM" | awk '{print $5}') bytes"

# Resolve config glob. If no args, run all 120.
shopt -s nullglob
configs=("$REPO_ROOT"/.config/"${1:-}"*.yml)
shopt -u nullglob

# If user passed a partial pattern, expand it.
if [[ $# -gt 0 ]]; then
    configs=("$REPO_ROOT"/.config/$1*.yml)
    if [[ ${#configs[@]} -eq 0 ]]; then
        echo "[batch] no configs match pattern: $1*"
        exit 0
    fi
fi

total=${#configs[@]}
echo "[batch] $total configs to run"
echo ""

# Build initial index (empty array) if missing.
[[ -f "$INDEX_FILE" ]] || echo '[]' > "$INDEX_FILE"

ran=0; skipped=0; failed=0
for cfg in "${configs[@]}"; do
    name="$(basename "$cfg" .yml)"

    # Honour skip marker
    if [[ -f "$RESULTS_DIR/.skip-$name" ]]; then
        echo "[skip ] $name (skip marker present)"
        skipped=$((skipped + 1))
        continue
    fi

    echo "[run  ] $name ($((ran + skipped + failed + 1))/$total)"
    log="$LOG_DIR/$name.log"

    # Get the YAML's timeout_seconds; pad by 1.5x for the outer timeout.
    yaml_timeout=$(grep -E '^\s*timeout_seconds:' "$cfg" | awk '{print $2}')
    yaml_timeout=${yaml_timeout:-7200}
    wall=$((yaml_timeout * 3 / 2))

    start=$(date +%s)
    timeout --foreground "$wall" "$REPO_ROOT/bin/MAgHARCM" --config "$cfg" \
        >"$log" 2>&1
    rc=$?
    end=$(date +%s)
    duration=$((end - start))

    # Triage outcome from the log.
    compile_ok="false"; tests_pass="0"; tests_total="0"
    if grep -q 'compilation=true' "$log"; then compile_ok="true"; fi
    if grep -qE 'passed=[0-9]+/[0-9]+' "$log"; then
        tests_pass=$(grep -oE 'passed=[0-9]+' "$log" | tail -1 | cut -d= -f2 | cut -d/ -f1)
        tests_total=$(grep -oE 'passed=[0-9]+/[0-9]+' "$log" | tail -1 | cut -d= -f2 | cut -d/ -f2)
    fi
    status="INCOMPLETE"
    grep -q 'Validation INCOMPLETE' "$log" && status="INCOMPLETE"
    grep -q 'all_success=true' "$log" && status="SUCCESS"

    # Append to JSON index.
    python3 -c "
import json, sys
with open('$INDEX_FILE') as f: idx = json.load(f)
idx.append({
    'name': '$name',
    'config': '.config/$name.yml',
    'log': '.artifacts/_batch-logs/$name.log',
    'wall_seconds': $duration,
    'exit_code': $rc,
    'status': '$status',
    'compile_ok': '$compile_ok',
    'tests_passed': int('$tests_pass'),
    'tests_total': int('$tests_total'),
})
with open('$INDEX_FILE', 'w') as f: json.dump(idx, f, indent=2)
"
    if [[ $rc -eq 0 && "$status" == "SUCCESS" ]]; then
        echo "       ok   ${duration}s  status=$status"
    else
        echo "       warn ${duration}s  status=$status  rc=$rc"
        failed=$((failed + 1))
    fi
    ran=$((ran + 1))
done

echo ""
echo "[batch] done: ran=$ran, skipped=$skipped, failed=$failed"
echo "[batch] index: $INDEX_FILE"
