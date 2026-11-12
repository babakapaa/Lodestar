#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
out="${1:-$root/../../web/crdt}"
mkdir -p "$out"
out="$(cd "$out" && pwd)"

cd "$root/rust"
cargo build --release --target wasm32-unknown-unknown
wasm-bindgen --target web --no-typescript --out-dir "$out" \
  target/wasm32-unknown-unknown/release/lodestar_crdt.wasm
