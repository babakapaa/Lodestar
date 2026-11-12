#!/usr/bin/env bash
set -euo pipefail

app="$(cd "$(dirname "$0")/.." && pwd)"
web="$app/web"

locked() {
  awk -v name="$1" '$1 == name":" { found = 1 } found && $1 == "version:" { gsub(/"/, "", $2); print $2; exit }' "$app/pubspec.lock"
}

drift="$(locked drift)"
sqlite3="$(locked sqlite3)"

curl -fsSL -o "$web/drift_worker.js" \
  "https://github.com/simolus3/drift/releases/download/drift-$drift/drift_worker.js"
curl -fsSL -o "$web/sqlite3.wasm" \
  "https://github.com/simolus3/sqlite3.dart/releases/download/sqlite3-$sqlite3/sqlite3.wasm"
"$app/packages/lodestar_crdt/tool/build_web.sh" "$web/crdt"
