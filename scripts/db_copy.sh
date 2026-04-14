#!/usr/bin/env bash

set -euo pipefail

MODE="${1:-}"
SRC="${2:-}"
DST="${3:-}"

if [[ -z "$MODE" || -z "$SRC" || -z "$DST" ]]; then
  echo "Usage: $0 --schema-only|--data-only|--full <source_url> <target_url>" >&2
  exit 1
fi

case "$MODE" in
  --schema-only)
    echo "Dumping schema from source, applying to target..."
    pg_dump "$SRC" --schema-only --no-owner --no-acl | psql "$DST" -v ON_ERROR_STOP=1
    ;;
  --data-only)
    echo "Dumping data from source, applying to target (target schema must already match)..."
    pg_dump "$SRC" --data-only --disable-triggers | psql "$DST" -v ON_ERROR_STOP=1
    ;;
  --full)
    echo "Full dump (schema + data) into target..."
    pg_dump "$SRC" --no-owner --no-acl | psql "$DST" -v ON_ERROR_STOP=1
    ;;
  *)
    echo "Unknown mode: $MODE" >&2
    exit 1
    ;;
esac

echo "Done."
