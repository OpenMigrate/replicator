#!/usr/bin/env bash
set -euo pipefail

# pre-commit commit-msg hook: ensure commit message references a GitHub issue
# Usage: scripts/require_issue_ref.sh <commit_msg_file>

msg_file="${1:-}"
if [[ -z "${msg_file}" || ! -f "${msg_file}" ]]; then
  echo "commit-msg hook error: commit message file not found" >&2
  exit 1
fi

# Read first line (subject)
subject="$(head -n1 "$msg_file" | tr -d '\r')"

# Allow merges and reverts to pass without extra checks
if [[ "$subject" =~ ^Merge\  || "$subject" =~ ^Revert\  ]]; then
  exit 0
fi

# Allow automated release/version bump commits
if [[ "$subject" =~ ^chore\(release\):\  || "$subject" =~ ^bump: ]]; then
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

echo "Commit message must reference a GitHub issue (e.g., 'Refs #123' or 'Closes #123')." >&2
echo "Branch already encodes the issue ID; please include '#<num>' in the message body or footer." >&2
exit 1

