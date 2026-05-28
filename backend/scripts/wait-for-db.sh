#!/bin/bash
# scripts/wait-for-db.sh
# Waits for database to be ready

set -e

host="$1"
port="${2:-5432}"
shift 2
cmd="$@"

until pg_isready -h "$host" -p "$port" > /dev/null 2>&1; do
  >&2 echo "Database is unavailable - sleeping"
  sleep 1
done

>&2 echo "Database is up - executing command"
exec $cmd
