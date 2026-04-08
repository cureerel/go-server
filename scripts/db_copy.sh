#!/usr/bin/env bash
# Copy PostgreSQL schema and/or data between two databases using connection URIs.
# Usage:
#   ./scripts/db_copy.sh --schema-only "$SOURCE_DATABASE_URL" "$TARGET_DATABASE_URL"
#   ./scripts/db_copy.sh --data-only   "$SOURCE_DATABASE_URL" "$TARGET_DATABASE_URL"
#   ./scripts/db_copy.sh --full        "$SOURCE_DATABASE_URL" "$TARGET_DATABASE_URL"
#
# Requires: pg_dump, psql, pg_restore (for custom format — not used here; uses plain SQL pipe).
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
