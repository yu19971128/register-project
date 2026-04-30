#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
WORKFLOW_DIR="${PROJECT_ROOT}/workflow"

usage() {
  echo "Usage: $0 <phase>"
  echo "  phase: design | implementation | qa"
  exit 1
}

if [[ $# -lt 1 ]]; then
  usage
fi

PHASE="$1"
ERRORS=0
info() { echo "[INFO] $*"; }
error() { echo "[ERROR] $*"; ((ERRORS++)); }

check_design_gate() {
  info "Checking design phase exit gate..."
  if [[ ! -f "${PROJECT_ROOT}/docs/00-index.md" ]]; then
    error "Missing docs/00-index.md"
  fi
  if [[ ! -f "${PROJECT_ROOT}/workflow/doc-relations.yaml" ]]; then
    error "Missing workflow/doc-relations.yaml"
  fi
}

check_implementation_gate() {
  info "Checking implementation phase exit gate..."
  local task_count=0
  for task in "${PROJECT_ROOT}/docs/07-tasks"/TASK-*.md; do
    [[ -f "$task" ]] || continue
    ((task_count++))
  done
  if (( task_count == 0 )); then
    error "No Task Cards found in docs/07-tasks/"
  fi
}

check_qa_gate() {
  info "Checking QA phase exit gate..."
  # Placeholder: real checks would run test suites
  info "QA gate placeholder: ensure tests pass before calling this gate."
}

case "$PHASE" in
  design) check_design_gate ;;
  implementation) check_implementation_gate ;;
  qa) check_qa_gate ;;
  *) usage ;;
esac

if (( ERRORS > 0 )); then
  echo ""
  error "Phase exit gate '${PHASE}' failed with ${ERRORS} error(s)."
  exit 1
fi

info "Phase exit gate '${PHASE}' passed."
exit 0
