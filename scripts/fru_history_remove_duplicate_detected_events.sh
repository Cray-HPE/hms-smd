#!/usr/bin/env bash
# MIT License
#
# (C) Copyright [2025] Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.

# This script manually applies the pruning changes in:
#
#     migrations/postgres/23_fru_history_remove_duplicate_detected_events.up.sql
#
# This is provided outside of a migration so that databases can be pruned
# outside the scope of an SMD upgrade.
#
# Set PGPASSWORD in the env if default credentials don't work.
#
# If default user changed from 'postgres' then update DB_USER below
#
# If any issues occur, uncomment the following line to aid in debug:
#set -x

set -eo pipefail

DB_NAME="hmsds"
DB_USER="postgres"

# Capture and print sizes before pruning and vacuuming

starting_stats_output=$(
psql "dbname=$DB_NAME user=$DB_USER" 2>&1 <<'SQL'
DO $$
DECLARE
    tbl_size_before    bigint := 0;   /* hwinv_history size before pruning */
    db_size_before     bigint := 0;   /* DB size before pruning */
BEGIN
    SELECT pg_total_relation_size('hwinv_hist') INTO tbl_size_before;
    SELECT pg_database_size(current_database()) INTO db_size_before;

    RAISE NOTICE 'hwinv_history table size before pruning:  % mb', tbl_size_before / 1024 / 1024;
    RAISE NOTICE 'Database size before pruning:             % mb', db_size_before / 1024 / 1024;
END;
$$ LANGUAGE plpgsql;
SQL
)

echo "$starting_stats_output"

# Creating a timestamp index speeds up execution several orders of magnitude.
# We will drop it after pruning is complete.

echo 'Creating temporary "timestamp" index on hwinv_hist table...'

psql "dbname=$DB_NAME user=$DB_USER" -c \
  "CREATE INDEX IF NOT EXISTS hwinvhist_id_ts_idx ON hwinv_hist (id, "timestamp");"

if [[ $? -ne 0 ]]; then
  echo "Error creating temporary timestamp index" >&2
  exit 1
fi

# Now run main pruning logic

echo 'Pruning hwinv_hist table... This will take a while - please do not interrupt'

pruning_output=$(
psql "dbname=$DB_NAME user=$DB_USER" 2>&1 <<'SQL'
    WITH dups AS (
        SELECT id, "timestamp"
        FROM (
            SELECT id, "timestamp", event_type,
                   LAG(event_type) OVER (PARTITION BY id ORDER BY "timestamp") AS prev_type
            FROM hwinv_hist
            WHERE id IN (
                SELECT hist.id
                FROM hwinv_hist hist
                JOIN hwinv_by_loc loc ON loc.id = hist.id
                WHERE loc.type IN ('Processor', 'NodeAccel')
            )
        ) sub
        WHERE event_type = 'Detected' AND prev_type = 'Detected'
    )
    DELETE FROM hwinv_hist h
        USING dups
        WHERE h.id = dups.id AND h."timestamp" = dups."timestamp";
SQL
)

if [[ $? -ne 0 ]]; then
  echo "Error executing SQL command" >&2
  exit 1
fi

echo "$pruning_output"

# The pruning logic above removed a large part of the hwinv_hist table. This
# did not however change the teble and database sizes.  In order to free
# space, a vacuum must be run.  Here are the options:
#
#     1. Do nothing
#            * Standard vacuum will eventually run but could be days or weeks
#     2. Standard vacuum
#            * Non-blocking
#            * Frees internal space
#            * Does not return disk space to the OS
#     2. Full vacuum
#            * Blocking - no updates allowed to table until complete
#            * Frees internal space
#            * Returns disk space to the OS
#
# So, run a full vacuum

echo "Running VACUUM FULL on hwinv_hist table to reclaim disk space..."

psql "dbname=$DB_NAME user=$DB_USER" -c "VACUUM FULL hwinv_hist;" || {
  echo "Error running VACUUM FULL" >&2
  exit 1
}

# Now drop the index to free up associated resources

echo 'Dropping temporary "timestamp" index on hwinv_hist table...'

psql "dbname=$DB_NAME user=$DB_USER" -c "DROP INDEX IF EXISTS hwinvhist_id_ts_idx;"

if [[ $? -ne 0 ]]; then
  echo "Error dropping temporary timestamp index" >&2
  exit 1
fi

# Capture and print sizes after pruning and vacuuming

ending_stats_output=$(
psql "dbname=$DB_NAME user=$DB_USER" 2>&1 <<'SQL'
DO $$
DECLARE
    tbl_size_after    bigint := 0;   /* hwinv_history size after pruning */
    db_size_after     bigint := 0;   /* DB size after pruning */
BEGIN
    SELECT pg_total_relation_size('hwinv_hist') INTO tbl_size_after;
    SELECT pg_database_size(current_database()) INTO db_size_after;

    RAISE NOTICE 'hwinv_history table size after pruning:  % mb', tbl_size_after / 1024 / 1024;
    RAISE NOTICE 'Database size after pruning:             % mb', db_size_after / 1024 / 1024;
END;
$$ LANGUAGE plpgsql;
SQL
)

echo "$ending_stats_output"
