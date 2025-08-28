#!/usr/bin/env bash
set -euo pipefail

# Pre-commit hook: Warn for files >500 lines (up to 750),
# block commit if any file exceeds 750 lines.

warn_threshold=500
block_threshold=750

warned=0
blocked=0
declare -a blocked_files
declare -a warned_files

is_regular_file() {
  [[ -f "$1" ]]
}

line_count() {
  # Use wc -l; handle files with no trailing newline
  wc -l < "$1" | awk '{print $1}'
}

for f in "$@"; do
  # Skip if not a regular file (deleted, dir, symlink, etc.)
  if ! is_regular_file "$f"; then
    continue
  fi

  # Skip binary files: if grep treating binary as non-text, skip
  if ! LC_ALL=C grep -Iq . "$f"; then
    continue
  fi

  cnt=$(line_count "$f" || echo 0)

  if (( cnt > block_threshold )); then
    blocked=1
    blocked_files+=("$f ($cnt lines)")
  elif (( cnt > warn_threshold )); then
    warned=1
    warned_files+=("$f ($cnt lines)")
  fi
done

if (( warned == 1 )); then
  echo "[line-count] Warning: some files exceed ${warn_threshold} lines:" >&2
  for w in "${warned_files[@]}"; do
    echo "  - $w" >&2
  done
fi

if (( blocked == 1 )); then
  echo "[line-count] Error: files exceed ${block_threshold} lines; reduce size or split files:" >&2
  for b in "${blocked_files[@]}"; do
    echo "  - $b" >&2
  done
  exit 1
fi

exit 0
