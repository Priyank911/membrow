#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

(
  cd "$ROOT_DIR/apps/api"
  go run ./cmd/api
) &
API_PID=$!

(
  cd "$ROOT_DIR/apps/web"
  npm run dev
) &
WEB_PID=$!

cleanup() {
  kill "$API_PID" "$WEB_PID" 2>/dev/null || true
}

trap cleanup EXIT INT TERM
wait
