#!/bin/sh

echo "ENTERED THE START.SH"

set -e

echo "RUNNING DATABASE MIGRATIONS"

#/app/migrate -path /app/migrations -database "$PG_URL" -verbose down
/app/migrate -path /app/migrations -database "$PG_URL" -verbose up

echo "STARTING GOLANG APPLICATION"

exec "$@"
