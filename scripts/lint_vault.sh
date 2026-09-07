#!/usr/bin/env bash
#
# lint_vault.sh — ADR-V-001 enforcement lint
#
# Rule (ADR-V-001): Version markers MUST use `[[x.y.z ...]]` wikilink syntax.
#   - In `internal/`: bare `(P-NN)` outside backticks is FORBIDDEN. `(PRIM-NN)`
#     parentheticals are explicitly retained per the Sprint 2026-09-17 decision
#     (prose parentheticals = standard academic-writing convention).
#   - In `.obsidian/MAgHARCM/`: bare `(PRIM-NN)` outside backticks is scanned
#     and reported; the canonical permitted form is `[[1.0.0 PRIM-NN]]`.
#     Bare `(P-NN)` in vault markdown is exempted (vault content, not code).
#
# Exit codes:
#   0 — clean
#   1 — forbidden reference(s) found
#
set -u

cd "$(dirname "$0")/.." || exit 2
ROOT="$(pwd)"

INTERNAL_DIR="$ROOT/internal"
VAULT_DIR="$ROOT/.obsidian/MAgHARCM"

if [ ! -d "$INTERNAL_DIR" ]; then
  echo "lint_vault: internal/ directory not found at $INTERNAL_DIR" >&2
  exit 2
fi
if [ ! -d "$VAULT_DIR" ]; then
  echo "lint_vault: vault directory not found at $VAULT_DIR" >&2
  exit 2
fi

# Pick ripgrep if available, otherwise grep -E.
if command -v rg >/dev/null 2>&1; then
  SEARCH() {
    rg --no-heading --line-number --color=never "$@"
  }
  SEARCH_FILES_COUNT() {
    rg --files "$@" | wc -l
  }
else
  SEARCH() {
    # Translate rg --no-heading --line-number --color=never PATTERN PATH
    # to grep -rnE.
    local pat="$1"; shift
    local path="$1"; shift
    grep -rnE -- "$pat" "$path"
  }
  SEARCH_FILES_COUNT() {
    # Count regular files (rg --files equivalent for our purposes).
    find "$@" -type f | wc -l
  }
fi

# Helper: print whether a match line has the (P-NN|PRIM-NN) token inside
# backticks on that line. Inline code (``...``) and fenced code blocks
# (``` ... ```) are treated as exempt contexts.
#
# Args:
#   $1 — file path (absolute)
#   $2 — line number
#   $3 — full line content
# Returns:
#   0 — token IS inside a backtick/code-fence context (exempt)
#   1 — token is bare (potentially forbidden)
is_in_code_context() {
  local file="$1"
  local lineno="$2"
  local line="$3"

  # Inline-code check: count backticks on the line. If odd count and the token
  # appears after an opening backtick that hasn't been closed, treat as inline
  # code. Conservative: if the line contains an odd number of backticks AND
  # the token appears after the first backtick, exempt.
  local first_bt="${line%%\`*}"
  local before_token="${line%%\(P*}"
  # If the token position is past any backtick on the line, we treat the
  # substring up to the token as having at least one unmatched backtick.
  if [[ "$before_token" == *\`* ]]; then
    # Count backticks in the prefix before the token.
    local prefix="${line%%\(P*}"
    local count="${prefix//[^\`]/}"
    if (( ${#count} % 2 == 1 )); then
      return 0
    fi
  fi

  # Fenced code-block check: scan preceding lines for ``` or ~~~ fence openers.
  local fence_marker=""
  local i=$((lineno - 1))
  while (( i > 0 )); do
    local prev
    prev="$(sed -n "${i}p" "$file" 2>/dev/null)"
    if [[ -z "$fence_marker" ]]; then
      if [[ "$prev" == *'```'* ]] && [[ "$prev" =~ ^[[:space:]]*\`\`\` ]]; then
        fence_marker='```'
      elif [[ "$prev" == *'~~~'* ]] && [[ "$prev" =~ ^[[:space:]]*\~\~\~ ]]; then
        fence_marker='~~~'
      fi
    else
      # We're inside a fence; close it when we see the matching marker.
      local marker_prefix="^[[:space:]]*${fence_marker}"
      if [[ "$prev" =~ $marker_prefix ]]; then
        fence_marker=""
      fi
    fi
    i=$((i - 1))
  done

  if [[ -n "$fence_marker" ]]; then
    return 0
  fi

  return 1
}

# Counters
files_scanned=0
refs_checked=0
violations=0

# ---------------------------------------------------------------------------
# Pass 1 — internal/: flag bare (P-NN) outside backticks.
#   (PRIM-NN) parentheticals are explicitly retained per Sprint 2026-09-17.
# ---------------------------------------------------------------------------
internal_files="$(SEARCH_FILES_COUNT "$INTERNAL_DIR")"
files_scanned=$((files_scanned + internal_files))

# Pull all (P-NN) matches (NOT PRIM-prefixed) in internal/.
# Pattern: literal "(" then "P-" then digits then ")". Use rg / grep -E.
internal_p_matches="$(SEARCH '\(P-[0-9]+\)' "$INTERNAL_DIR" || true)"

while IFS= read -r line; do
  [ -z "$line" ] && continue
  refs_checked=$((refs_checked + 1))
  # rg format: /path/to/file:lineno:content
  # grep -rn format: /path/to/file:lineno:content
  file="$(echo "$line" | cut -d: -f1)"
  lineno="$(echo "$line" | cut -d: -f2)"
  content="$(echo "$line" | cut -d: -f3-)"
  if is_in_code_context "$file" "$lineno" "$content"; then
    continue
  fi
  echo "FORBIDDEN: $file:$lineno: $content" >&2
  violations=$((violations + 1))
done <<< "$internal_p_matches"

# (PRIM-NN) in internal/ is permitted (Sprint 2026-09-17); count but don't flag.
internal_prim_matches="$(SEARCH '\(PRIM-[0-9]+\)' "$INTERNAL_DIR" || true)"
while IFS= read -r line; do
  [ -z "$line" ] && continue
  refs_checked=$((refs_checked + 1))
done <<< "$internal_prim_matches"

# ---------------------------------------------------------------------------
# Pass 2 — .obsidian/MAgHARCM/: scan (PRIM-NN); permitted form is wikilink.
#   Bare (PRIM-NN) outside backticks is reported. Markdown files in
#   .obsidian/MAgHARCM/ are exempt from the (P-NN) rule (vault content).
# ---------------------------------------------------------------------------
vault_files="$(SEARCH_FILES_COUNT "$VAULT_DIR")"
files_scanned=$((files_scanned + vault_files))

vault_prim_matches="$(SEARCH '\(PRIM-[0-9]+\)' "$VAULT_DIR" || true)"
while IFS= read -r line; do
  [ -z "$line" ] && continue
  refs_checked=$((refs_checked + 1))
  file="$(echo "$line" | cut -d: -f1)"
  lineno="$(echo "$line" | cut -d: -f2)"
  content="$(echo "$line" | cut -d: -f3-)"
  if is_in_code_context "$file" "$lineno" "$content"; then
    continue
  fi
  echo "STRAY: $file:$lineno: $content" >&2
  # STRAY is informational, not a forbidden violation in this sprint.
done <<< "$vault_prim_matches"

# ---------------------------------------------------------------------------
# Result
# ---------------------------------------------------------------------------
if (( violations > 0 )); then
  echo "vault lint FAILED: $violations forbidden reference(s) found" >&2
  exit 1
fi

echo "vault lint clean: $files_scanned files scanned, $refs_checked refs checked"
exit 0
