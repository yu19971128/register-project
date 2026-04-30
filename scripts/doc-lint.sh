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

check_file_exists() {
  local file="$1"
  if [[ ! -f "${PROJECT_ROOT}/${file}" ]]; then
    error "Missing required file: ${file}"
  fi
}

check_frontmatter() {
  local file="$1"
  local req_keys=("$@")
  req_keys=("${req_keys[@]:1}")
  for key in "${req_keys[@]}"; do
    if ! grep -qE "^${key}:" "$file"; then
      error "Missing frontmatter key '${key}' in ${file#${PROJECT_ROOT}/}"
    fi
  done
}

info "Checking root GSD files..."
check_file_exists "PROJECT.md"
check_file_exists "AGENTS.md"
check_file_exists "DECISIONS.md"
check_file_exists "KNOWLEDGE.md"
check_file_exists "M001-ROADMAP.md"

info "Checking docs tree..."
for dir in 01-background 02-modules 03-architecture 04-frontend 05-database 06-api 07-tasks; do
  if [[ ! -d "${DOCS_DIR}/${dir}" ]]; then
    error "Missing docs directory: docs/${dir}"
  fi
done

info "Checking 00-index.md columns..."
if [[ -f "${DOCS_DIR}/00-index.md" ]]; then
  for col in "功能点/需求" "feature_id" "优先级" "需求文档" "模块设计" "前端设计" "数据库设计" "API 文档" "当前状态" "任务拆解"; do
    if ! grep -qF "${col}" "${DOCS_DIR}/00-index.md"; then
      error "Missing column '${col}' in docs/00-index.md"
    fi
  done
fi

info "Checking module README frontmatter..."
for readme in "${DOCS_DIR}"/02-modules/*/README.md; do
  [[ -f "$readme" ]] || continue
  # Skip non-module dirs like _template
  if [[ "$(basename "$(dirname "$readme")")" == "_template" ]]; then continue; fi
  if [[ "$(basename "$(dirname "$readme")")" == "orchestrator" ]]; then continue; fi
  check_frontmatter "$readme" "module_id" "module_name" "feature_ids" "status"
done

info "Checking workflow files..."
for wf in doc-relations.yaml doc-statuses.yaml frontend-rules.yaml task-rules.yaml runtime-policies.yaml; do
  check_file_exists "workflow/${wf}"
done

if (( ERRORS > 0 )); then
  echo ""
  error "Lint failed with ${ERRORS} error(s)."
  exit 1
fi

info "All checks passed."
exit 0
