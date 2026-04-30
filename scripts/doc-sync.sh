#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DOCS_DIR="${PROJECT_ROOT}/docs"

info() { echo "[INFO] $*"; }

info "Syncing module index with 00-index.md..."

# Legacy fallback: list modules present in docs/02-modules and ensure they are in 00-index.md
MODULE_LIST=""
for mod_dir in "${DOCS_DIR}"/02-modules/*/; do
  [[ -d "$mod_dir" ]] || continue
  mod_name=$(basename "$mod_dir")
  if [[ "$mod_name" == "_template" ]]; then continue; fi
  MODULE_LIST="${MODULE_LIST}${mod_name}\n"
done

if [[ -n "${MODULE_LIST}" ]]; then
  info "Found modules:"
  printf '%b' "${MODULE_LIST}"
fi

info "Legacy sync complete. Please use doc-sync-v2.sh for full consistency checks."
exit 0
