#!/usr/bin/env bash
set -euo pipefail

if [ $# -ne 2 ]; then
  echo "usage: $0 <module_path> <web_module_name>" >&2
  exit 1
fi

MOD="$1"
WEB="$2"

SCRIPT_PATH="scripts/rename.sh"

# replace placeholders
git ls-files | while read -r f; do
  # skip this script itself so its placeholders stay intact
  if [ "$f" = "$SCRIPT_PATH" ]; then
    continue
  fi

  sed -i \
    -e "s|{{MODULE_PATH}}|$MOD|g" \
    -e "s|{{webmod}}|$WEB|g" \
    "$f"
done

# move folders
git mv "cmd/{{webmod}}" "cmd/$WEB"
git mv "pkg/{{webmod}}" "pkg/$WEB"
