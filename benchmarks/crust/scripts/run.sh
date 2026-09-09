#!/usr/bin/env bash
# benchmarks/crust/scripts/run.sh
#
# Long-horizon batch runner for the crust C→Rust benchmark suite.
#
# Walks every benchmarks/crust/configs/*.yml in sequence, executes the
# MAgHARCM pipeline against the matching crust project, and persists per-
# project artefacts under benchmarks/crust/results/<commit>/<proj>/.
#
# Outputs per project:
#   <commit>/<proj>/.log        full console capture
#   <commit>/<proj>/result.yml  scraped metrics (status, wall_seconds,
#                               compile_ok, tests_passed, tests_total,
#                               pass_rate, iterations, final_iter_compile,
#                               final_iter_tests_passed, exit_code)
#   <commit>/<proj>/.interrupted sentinel when SIGINT/SIGTERM killed the run
#                               before the pipeline emitted FINALSUM
#
# Plus an aggregate:
#   <commit>/index.yml          per-project summary across the batch
#
# Resume semantics:
#   - If <proj>/result.yml exists with status: completed (the pipeline
#     emitted FINALSUM before exiting), the project is skipped.
#   - If <proj>/.interrupted exists, or result.yml is missing/partial,
#     the project is re-run; the stale directory is removed first so
#     downstream consumers never see half-finished artefacts.
#
# Usage:
#   bash benchmarks/crust/scripts/run.sh                       # full batch
#   bash benchmarks/crust/scripts/run.sh --only 2dpartint       # single proj
#   bash benchmarks/crust/scripts/run.sh --commit HEAD~1        # pin commit
#   bash benchmarks/crust/scripts/run.sh --dry-run              # list only
#
# Interruption:
#   Ctrl-C / SIGTERM: finishes the current project's log capture,
#   marks it as interrupted, then exits with the next project NOT
#   being started. Re-running the script resumes from after that point.

set -uo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$REPO_ROOT"

CONFIGS_DIR="$REPO_ROOT/benchmarks/crust/configs"
RESULTS_ROOT="$REPO_ROOT/benchmarks/crust/results"
SCRIPT_NAME="$(basename "$0")"

# --- args ----------------------------------------------------------------
ONLY=""
PIN_COMMIT=""
DRY_RUN=0
RESUME=1   # default on; --no-resume disables
while [[ $# -gt 0 ]]; do
    case "$1" in
        --only)        ONLY="$2"; shift 2 ;;
        --commit)      PIN_COMMIT="$2"; shift 2 ;;
        --dry-run)     DRY_RUN=1; shift ;;
        --no-resume)   RESUME=0; shift ;;
        -h|--help)
            sed -n '2,46p' "$0" | sed 's/^# \{0,1\}//'
            exit 0 ;;
        *) echo "[$SCRIPT_NAME] unknown arg: $1" >&2; exit 2 ;;
    esac
done

# --- commit + dirs --------------------------------------------------------
if [[ -n "$PIN_COMMIT" ]]; then
    COMMIT="$PIN_COMMIT"
else
    COMMIT="$(git rev-parse HEAD 2>/dev/null || echo unknown)"
fi
COMMIT_SHORT="${COMMIT:0:12}"
RESULTS_DIR="$RESULTS_ROOT/$COMMIT"
LOG_DIR="$RESULTS_DIR/.logs"
INDEX_FILE="$RESULTS_DIR/index.yml"
if [[ "$DRY_RUN" -eq 0 ]]; then
    mkdir -p "$RESULTS_DIR" "$LOG_DIR"
    touch "$INDEX_FILE"
fi

CURRENT_PROJ=""
INTERRUPTED=0
on_signal() {
    local sig=$1
    INTERRUPTED=1
    echo ""
    echo "[$SCRIPT_NAME] received SIG${sig}; finalising current project ($CURRENT_PROJ)"
    if [[ -n "$CURRENT_PROJ" ]]; then
        touch "$RESULTS_DIR/$CURRENT_PROJ/.interrupted"
    fi
    rebuild_index
    exit 130
}
trap 'on_signal INT'  INT
trap 'on_signal TERM' TERM

# --- helpers --------------------------------------------------------------
build_binary() {
    local bin="$REPO_ROOT/bin/MAgHARCM"
    if [[ -x "$bin" && "$bin" -nt "$REPO_ROOT/cmd/MAgHARCM/main.go" ]]; then
        return 0
    fi
    echo "[$SCRIPT_NAME] building magharcm-cli..."
    go build -o "$bin" ./cmd/MAgHARCM || { echo "[$SCRIPT_NAME] build failed" >&2; exit 1; }
}

# Uses the last FINALSUM line the validator emits per iteration.
scrape_final() {
    local log=$1
    local line
    line="$(grep -E '^\[.*\] FINALSUM ' "$log" | tail -1 || true)"
    if [[ -z "$line" ]]; then
        echo "comp=unknown tests_passed=0 tests_total=0 all_success=false iter=0"
        return 0
    fi
    local comp all_succ pass tot iter
    comp="$(echo "$line" | grep -oE 'comp=[^ ]+' | cut -d= -f2)"
    pass="$(echo "$line" | grep -oE 'tests=[0-9]+' | cut -d= -f2)"
    tot="$(echo "$line" | grep -oE 'tests=[0-9]+/[0-9]+' | cut -d/ -f2)"
    all_succ="$(echo "$line" | grep -oE 'all_success=[^ ]+' | cut -d= -f2)"
    iter="$(echo "$line" | grep -oE 'iter=[0-9]+' | cut -d= -f2)"
    comp=${comp:-unknown}
    pass=${pass:-0}
    tot=${tot:-0}
    all_succ=${all_succ:-false}
    iter=${iter:-0}
    echo "comp=$comp tests_passed=$pass tests_total=$tot all_success=$all_succ iter=$iter"
}

# scrape_iterations <log> -> echoes semicolon-separated ITER entries
# (one per validator iteration). Used so result.yml can record the
# per-iteration trajectory, not just the final snapshot.
scrape_iterations() {
    local log=$1
    grep -E '^\[.*\] ITER\[' "$log" \
        | sed -E 's/^.*ITER\[([0-9]+)\] comp=([a-z]+) tests=([0-9]+)\/([0-9]+).*/\1 comp=\2 tests_passed=\3 tests_total=\4/' \
        | tr '\n' ';' || true
}

write_result_yml() {
    local proj=$1 log=$2 wall=$3 rc=$4 iters=$5 final=$6
    local comp pass tot all_succ iter
    eval "$final"  # sets comp tests_passed tests_total all_success iter
    local rate="0.0"
    if [[ "$tot" -gt 0 ]]; then
        rate="$(awk -v p="$pass" -v t="$tot" 'BEGIN{printf "%.4f", p/t}')"
    fi
    cat > "$RESULTS_DIR/$proj/result.yml" <<EOF
project: $proj
config: benchmarks/crust/configs/${proj}.yml
status: completed
exit_code: $rc
wall_seconds: $wall
compile_ok: $comp
tests_passed: $pass
tests_total: $tot
pass_rate: $rate
iterations: $iter
final_iter_compile: $comp
final_iter_tests_passed: $pass
final_iter_tests_total: $tot
per_iteration:
$(echo "$iters" | tr ';' '\n' | sed '/^$/d' | sed 's/^/  - /')
EOF
}

mark_interrupted() {
    local proj=$1 log=$2 wall=$3 iters=$4
    cat > "$RESULTS_DIR/$proj/result.yml" <<EOF
project: $proj
config: benchmarks/crust/configs/${proj}.yml
status: interrupted
exit_code: -
wall_seconds: $wall
compile_ok: unknown
tests_passed: 0
tests_total: 0
pass_rate: 0.0
iterations: 0
per_iteration:
$(echo "$iters" | tr ';' '\n' | sed '/^$/d' | sed 's/^/  - /')
EOF
    touch "$RESULTS_DIR/$proj/.interrupted"
}

rebuild_index() {
    [[ -f "$INDEX_FILE" ]] && rm -f "$INDEX_FILE"
    {
        echo "commit: $COMMIT"
        echo "generated_at: \"$(date -Iseconds 2>/dev/null || date)\""
        echo "projects:"
        for d in "$RESULTS_DIR"/*/; do
            [[ -f "$d/result.yml" ]] || continue
            local proj
            proj="$(basename "$d")"
            awk -v proj="$proj" '
                /^project:/    {print "  - project: " $2}
                /^status:/     {print "    status: " $2}
                /^wall_seconds:/ {print "    wall_seconds: " $2}
                /^compile_ok:/ {print "    compile_ok: " $2}
                /^tests_passed:/ {print "    tests_passed: " $2}
                /^tests_total:/ {print "    tests_total: " $2}
                /^pass_rate:/ {print "    pass_rate: " $2}
                /^iterations:/ {print "    iterations: " $2}
                /^exit_code:/  {print "    exit_code: " $2}
            ' "$d/result.yml"
        done
    } > "$INDEX_FILE"
}

invalidate_stale() {
    local proj=$1
    local proj_dir="$RESULTS_DIR/$proj"
    # Considered "stale" when:
    #   - .interrupted sentinel exists, OR
    #   - result.yml absent, OR
    #   - result.yml says status != completed
    if [[ -f "$proj_dir/.interrupted" ]]; then return 0; fi
    if [[ ! -f "$proj_dir/result.yml" ]]; then return 0; fi
    if ! grep -qE '^status: completed$' "$proj_dir/result.yml"; then return 0; fi
    return 1
}

run_one() {
    local cfg=$1
    local proj
    proj="$(basename "$cfg" .yml)"
    CURRENT_PROJ="$proj"
    local proj_dir="$RESULTS_DIR/$proj"
    local log="$proj_dir/.log"

    # Resume: skip if completed and resume enabled.
    if [[ "$RESUME" -eq 1 && -f "$proj_dir/result.yml" ]] \
        && grep -qE '^status: completed$' "$proj_dir/result.yml" \
        && [[ ! -f "$proj_dir/.interrupted" ]]; then
        echo "[skip ] $proj (completed)"
        return 0
    fi

    echo "[run  ] $proj"
    mkdir -p "$proj_dir"
    rm -f "$proj_dir/.interrupted" "$proj_dir/result.yml"

    local start end wall rc final iters
    start=$(date +%s)
    # shellcheck disable=SC2086
    "$REPO_ROOT/bin/MAgHARCM" --config "$cfg" >"$log" 2>&1
    end=$(date +%s)
    wall=$((end - start))

    iters="$(scrape_iterations "$log")"
    final="$(scrape_final "$log")"

    if [[ -z "$(grep -E '^\[.*\] FINALSUM ' "$log" 2>/dev/null)" ]]; then
        # pipeline crashed before FINALSUM
        mark_interrupted "$proj" "$log" "$wall" "$iters"
        echo "[err  ] $proj wall=${wall}s (no FINALSUM; rc=$rc)"
    else
        write_result_yml "$proj" "$log" "$wall" "$rc" "$iters" "$final"
        echo "[done ] $proj wall=${wall}s rc=$rc $final"
    fi

    rebuild_index
    return 0
}

# --- main -----------------------------------------------------------------
echo "[$SCRIPT_NAME] commit=$COMMIT_SHORT  results=$RESULTS_DIR"
[[ -d "$CONFIGS_DIR" ]] || { echo "[$SCRIPT_NAME] no configs at $CONFIGS_DIR" >&2; exit 1; }

if [[ "$DRY_RUN" -eq 1 ]]; then
    echo "[$SCRIPT_NAME] dry-run; configs:"
    for cfg in "$CONFIGS_DIR"/*.yml; do
        n="$(basename "$cfg" .yml)"
        if [[ -n "$ONLY" && "$n" != "$ONLY" ]]; then continue; fi
        echo "  $n"
    done
    exit 0
fi
build_binary

ran=0; skipped=0
for cfg in "$CONFIGS_DIR"/*.yml; do
    [[ -n "$ONLY" ]] || true
    n="$(basename "$cfg" .yml)"
    if [[ -n "$ONLY" && "$n" != "$ONLY" ]]; then continue; fi
    if run_one "$cfg"; then
        if grep -qE '^status: completed$' "$RESULTS_DIR/$n/result.yml" 2>/dev/null; then
            skipped=$((skipped + 1))
        else
            ran=$((ran + 1))
        fi
    fi
    [[ "$INTERRUPTED" -eq 1 ]] && break
done

rebuild_index
echo "[$SCRIPT_NAME] done: ran=$ran skipped=$skipped  index=$INDEX_FILE"
