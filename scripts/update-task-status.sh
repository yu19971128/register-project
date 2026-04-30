#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DOCS_DIR="${PROJECT_ROOT}/docs"

usage() {
  echo "Usage: $0 <task-id> <new-status>"
  echo "  new-status must be one of: 需求澄清中 | 设计中 | 设计已完成 | 开发中 | 测试中 | 已完成"
  exit 1
}

if [[ $# -lt 2 ]]; then
  usage
fi

TASK_ID="$1"
NEW_STATUS="$2"
TASK_FILE="${DOCS_DIR}/07-tasks/${TASK_ID}.md"

if [[ ! -f "$TASK_FILE" ]]; then
  echo "[ERROR] Task file not found: ${TASK_FILE}"
  exit 1
fi

# Map status to icon
status_icon=""
case "$NEW_STATUS" in
  "需求澄清中") status_icon="📝" ;;
  "设计中") status_icon="🎨" ;;
  "设计已完成") status_icon="✅" ;;
  "开发中") status_icon="🔨" ;;
  "测试中") status_icon="🧪" ;;
  "已完成") status_icon="🏁" ;;
  *) echo "[ERROR] Invalid status: ${NEW_STATUS}"; usage ;;
esac

# Update Task Card status checklist
tmpfile=$(mktemp)
sed -E "s/^- \[x\] (规划中|开发中|审查中|测试中|已完成)$/- [ ] \1/" "$TASK_FILE" > "$tmpfile"
# Mark the closest matching status as checked
sed -i -E "s/^- \[ \] ${NEW_STATUS}$/- [x] ${NEW_STATUS}/" "$tmpfile" 2>/dev/null || true
mv "$tmpfile" "$TASK_FILE"

# Update 00-index.md if entry exists
INDEX_FILE="${DOCS_DIR}/00-index.md"
if [[ -f "$INDEX_FILE" ]] && grep -qF "$TASK_ID" "$INDEX_FILE"; then
  # naive replacement of status column for the row containing TASK_ID
  sed -i -E "s/^([|][^|]*${TASK_ID}[^|]*[|][^|]*[|][^|]*[|][^|]*[|][^|]*[|][^|]*[|][^|]*[|][^|]*[|])[^|]*(.*)$/\1 ${status_icon} ${NEW_STATUS} \2/" "$INDEX_FILE" 2>/dev/null || true
fi

echo "[OK] Updated ${TASK_ID} status to ${status_icon} ${NEW_STATUS}."
exit 0
