#!/bin/sh

set -e

echo "Running migrations..."

/app/migrate -path /app/migration -database "${DB_SOURCE}" -verbose up

echo "Starting server..."

exec "$@"
