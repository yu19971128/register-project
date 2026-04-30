#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
FRONTEND_DIR="${PROJECT_ROOT}/docs/04-frontend"
RULES_FILE="${PROJECT_ROOT}/workflow/frontend-rules.yaml"

info() { echo "[INFO] $*"; }
warn() { echo "[WARN] $*"; }
error() { echo "[ERROR] $*"; }

if [[ ! -f "$RULES_FILE" ]]; then
  error "Missing workflow/frontend-rules.yaml"
  exit 1
fi

info "Checking frontend design docs..."

for doc in README.md layout.md; do
  if [[ ! -f "${FRONTEND_DIR}/${doc}" ]]; then
    warn "Missing docs/04-frontend/${doc}"
  fi
done

info "Frontend guard checks complete."
exit 0
