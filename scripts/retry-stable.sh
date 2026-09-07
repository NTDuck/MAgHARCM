#!/usr/bin/env bash
# scripts/retry-stable.sh
#
# Run samples whose first-attempt result was "too weird" (per user spec:
# e.g. compile=0% / tests=0% on a sample that should pass). Retry up to
# 3 times with config perturbations; declare stable when two
# consecutive attempts agree on (compile_ok, tests_passed).
#
# Stability criterion: max(3 attempts), success when (a) two consecutive
# attempts agree on (compile_ok, tests_passed), OR (b) the third attempt
# matches the first (rule of three, sample-of-three).
#
# Config perturbations applied between attempts (cheap, no model swap):
#   1. Initial attempt as-is.
#   2. Increase max_iterations to 10 (cap) and timeout to 10800s.
#   3. Add lsp_provider override (native -> abcoder if available).
#
# Usage:
#   bash scripts/retry-stable.sh .config/oxidizer-gohistogram.yml
#   bash scripts/retry-stable.sh .config/crust-csyncmers.yml

set -uo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
RESULTS_DIR="$REPO_ROOT/docs/sample-results"

cfg="$1"
[[ -f "$cfg" ]] || { echo "usage: $0 <path-to-config.yml>"; exit 1; }

name="$(basename "$cfg" .yml)"
echo "[retry] target: $name"

# Read base fields.
base_yaml="$(cat "$cfg")"
max_iter=$(echo "$base_yaml" | grep -E 'max_iterations:' | awk '{print $2}')
timeout=$(echo "$base_yaml" | grep -E 'timeout_seconds:' | awk '{print $2}')
[[ -z "$max_iter" ]] && max_iter=1
[[ -z "$timeout" ]] && timeout=7200

declare -a history
for attempt in 1 2 3; do
    case $attempt in
        1) iter=$max_iter; tout=$timeout; lsp='native';;
        2) iter=10;        tout=10800;    lsp='native';;
        3) iter=10;        tout=10800;    lsp='abcoder';;
    esac

    # Build per-attempt config (rewrite 3 fields; keep everything else).
    attempt_cfg="$RESULTS_DIR/.retry-${name}-${attempt}.yml"
    echo "$base_yaml" \
        | sed -E "s/(max_iterations:[[:space:]]*)[0-9]+/\1$iter/" \
        | sed -E "s/(timeout_seconds:[[:space:]]*)[0-9]+/\1$tout/" \
        | sed -E "s/(lsp:[[:space:]]*$)/lsp:/" \
        > "$attempt_cfg"
    # patch lsp provider explicitly
    python3 -c "
import sys
text = open('$attempt_cfg').read()
text = text.replace('provider: native', 'provider: $lsp', 1)
open('$attempt_cfg','w').write(text)
"

    log="$RESULTS_DIR/.retry-${name}-${attempt}.log"
    echo "[retry] attempt $attempt: iter=$iter timeout=$tout lsp=$lsp"
    timeout --foreground "$((tout + 600))" "$REPO_ROOT/bin/MAgHARCM" --config "$attempt_cfg" \
        >"$log" 2>&1
    rc=$?

    compile_ok=false; tests_pass=0
    grep -q 'compilation=true' "$log" && compile_ok=true
    tests_pass=$(grep -oE 'passed=[0-9]+' "$log" | tail -1 | cut -d= -f2 | cut -d/ -f1)
    tests_pass=${tests_pass:-0}

    history+=("compile_ok=$compile_ok tests_pass=$tests_pass rc=$rc")
    echo "         → $compile_ok tests=$tests_pass rc=$rc"

    # Stability check: 2 consecutive matches OR 3-of-3 identical.
    if [[ $attempt -ge 2 ]]; then
        prev=$(echo "${history[$((attempt-2))]}" | grep -oE 'compile_ok=\w+ tests_pass=\w+')
        curr=$(echo "${history[$((attempt-1))]}" | grep -oE 'compile_ok=\w+ tests_pass=\w+')
        if [[ "$prev" == "$curr" ]]; then
            echo "[retry] STABLE on attempt $attempt (consecutive match)"
            echo "[retry] final state: $curr"
            exit 0
        fi
    fi
    if [[ $attempt -eq 3 ]]; then
        # Compare first vs third.
        first=$(echo "${history[0]}" | grep -oE 'compile_ok=\w+ tests_pass=\w+')
        third=$(echo "${history[2]}" | grep -oE 'compile_ok=\w+ tests_pass=\w+')
        if [[ "$first" == "$third" ]]; then
            echo "[retry] STABLE on attempt 3 (rule-of-three)"
            echo "[retry] final state: $third"
            exit 0
        fi
        echo "[retry] UNSTABLE after 3 attempts: history=${history[*]}"
        exit 2
    fi
done
