#!/usr/bin/env bash
set -euo pipefail

# pre-commit commit-msg hook: check commit message references a GitHub issue
# Usage: scripts/require_issue_ref.sh <commit_msg_file>

msg_file="${1:-}"
if [[ -z "${msg_file}" || ! -f "${msg_file}" ]]; then
  echo "commit-msg hook error: commit message file not found" >&2
  exit 1
fi

# Read first line (subject)
subject="$(head -n1 "$msg_file" | tr -d '\r')"
lower_subject="$(printf '%s' "$subject" | tr '[:upper:]' '[:lower:]')"

# Allow merges and reverts to pass without extra checks
if [[ "$lower_subject" =~ ^merge\  || "$lower_subject" =~ ^revert\  ]]; then
  exit 0
fi

# Allow automated release/version bump commits
if [[ "$lower_subject" =~ ^chore\(release\):\  || "$lower_subject" =~ ^bump: ]]; then
  exit 0
fi

# Require a #<num> anywhere in the message (subject or body)
if rg -N --pcre2 -q "#[0-9]{1,7}\b" "$msg_file" 2>/dev/null; then
  exit 0
fi

# Fallback to grep if ripgrep is unavailable
if command -v grep >/dev/null 2>&1 && grep -Eq "#[0-9]{1,7}(\b|$)" "$msg_file"; then
  exit 0
fi

echo "[issue-ref] Warning: no GitHub issue reference found (e.g., 'refs #123' or 'closes #123')." >&2
echo "Include '#<num>' in the message body or footer when possible." >&2
exit 0
