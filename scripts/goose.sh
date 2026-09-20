#!/bin/sh
set -eu

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  . ./.env
  set +a
fi

if [ -z "${DATABASE_URL:-}" ]; then
  echo "DATABASE_URL is not set. Copy .env.example to .env." >&2
  exit 1
fi

exec go run github.com/pressly/goose/v3/cmd/goose@v3.26.0 \
  -dir migrations \
  postgres "$DATABASE_URL" \
  "$@"
