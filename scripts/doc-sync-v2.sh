#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DOCS_DIR="${PROJECT_ROOT}/docs"
WORKFLOW_DIR="${PROJECT_ROOT}/workflow"

ERRORS=0
info() { echo "[INFO] $*"; }
warn() { echo "[WARN] $*"; }
error() { echo "[ERROR] $*"; ((ERRORS++)); }

validate=false
if [[ "${1:-}" == "--validate" ]]; then
  validate=true
fi

# Ensure 00-index.md exists
check_index() {
  if [[ ! -f "${DOCS_DIR}/00-index.md" ]]; then
    error "Missing docs/00-index.md"
    return
  fi
}

# Check that every module in 02-modules has an entry in 00-index.md
sync_modules_to_index() {
  for mod_dir in "${DOCS_DIR}"/02-modules/*/; do
    [[ -d "$mod_dir" ]] || continue
    mod_name=$(basename "$mod_dir")
    if [[ "$mod_name" == "_template" ]]; then continue; fi
    if ! grep -qF "$mod_name" "${DOCS_DIR}/00-index.md"; then
      if $validate; then
        error "Module '${mod_name}' not referenced in docs/00-index.md"
      else
        warn "Module '${mod_name}' not referenced in docs/00-index.md"
      fi
    fi
  done
}

# Check that every Task Card in 07-tasks is linked from 00-index.md
sync_tasks_to_index() {
  for task_file in "${DOCS_DIR}"/07-tasks/TASK-*.md; do
    [[ -f "$task_file" ]] || continue
    task_name=$(basename "$task_file" .md)
    if ! grep -qF "$task_name" "${DOCS_DIR}/00-index.md"; then
      if $validate; then
        error "Task Card '${task_name}' not referenced in docs/00-index.md"
      else
        warn "Task Card '${task_name}' not referenced in docs/00-index.md"
      fi
    fi
  done
}

check_index
sync_modules_to_index
sync_tasks_to_index

if $validate && (( ERRORS > 0 )); then
  echo ""
  error "Validation failed with ${ERRORS} error(s)."
  exit 1
fi

info "doc-sync-v2 complete."
exit 0
