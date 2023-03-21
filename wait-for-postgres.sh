#!/bin/sh
# wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "postgres" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

#psql -U "postgres" -d "postgres" -f /docker-entrypoint-initdb.d/init.sql

>&2 echo "Postgres is up - executing command"
exec $cmd