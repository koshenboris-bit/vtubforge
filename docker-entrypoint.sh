#!/usr/bin/env sh
set -eu

: "${DATABASE_URL:=postgres://postgres:postgres@db:5432/app?sslmode=disable}"
: "${HTTP_ADDR:=:8080}"
: "${JWT_ACCESS_SECRET:=super-secret-access}"
: "${JWT_REFRESH_SECRET:=super-secret-refresh}"
: "${RUN_MIGRATIONS:=true}"

export DATABASE_URL HTTP_ADDR JWT_ACCESS_SECRET JWT_REFRESH_SECRET

if [ "$RUN_MIGRATIONS" = "true" ]; then
  echo "Waiting for PostgreSQL..."
  until pg_isready -d "$DATABASE_URL" >/dev/null 2>&1; do
    sleep 1
  done

  echo "Running migrations..."
  psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f /app/migrations/001_init.sql
fi

exec /usr/local/bin/app
