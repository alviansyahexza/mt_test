#!/bin/bash
set -e

# Wait for PostgreSQL to be ready
until pg_isready -U "$POSTGRES_USER" > /dev/null 2>&1; do
  echo "Waiting for PostgreSQL..."
  sleep 2
done

# Create mt_test database if it doesn't exist
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d postgres -tc \
"SELECT 1 FROM pg_database WHERE datname = 'mt_test'" | grep -q 1 || \
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" -d postgres -c "CREATE DATABASE mt_test;"

echo "✅ Database mt_test created (or already exists)"
