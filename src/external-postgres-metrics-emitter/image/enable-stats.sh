#!/bin/bash

set -e

echo "shared_preload_libraries = 'pg_stat_statements'" >> $PGDATA/postgresql.conf
echo "pg_stat_statements.max = 10000" >> $PGDATA/postgresql.conf
echo "pg_stat_statements.track = all" >> $PGDATA/postgresql.conf

"${psql[@]}" --dbname="$POSTGRES_DB" <<-'EOSQL'
	CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

