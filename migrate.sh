#!/bin/sh
set -e

echo "Waiting for database to be ready..."
sleep 5

echo "Running migrations..."
for f in ./migrations/*.up.sql; do
  if [ -f "$f" ]; then
    echo "Running migration: $f"
    psql "postgresql://$PGUSER:$PGPASSWORD@$PGHOST:$PGPORT/$POSTGRES_DB?sslmode=disable" -f "$f" || true
  fi
done

echo "Migrations completed!"
