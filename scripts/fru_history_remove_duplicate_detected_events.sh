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
# Set PGPASSWORD in the env if default credentials don't work
# If default user changed from 'postgres' then update DB_USER below
#
# If any issues occur, uncomment the following line to aid in debug:
#set -x

set -eo pipefail

DB_NAME="hmsds"
DB_USER="postgres"

psql_output=$(
psql "dbname=$DB_NAME user=$DB_USER" 2>&1 <<'SQL'
DO $$
DECLARE
    unique_ids RECORD;
    base_event RECORD;
    next_event RECORD;

    del_total         bigint := 0;           /* Total deletions */
    del_per_id        jsonb := '{}'::jsonb;  /* Map of deletions per xname */
    del_this_id_pair  bigint := 0;           /* Deletions for this id pair */
    pair              RECORD;                /* For looping through the stats */
    tbl_size_before   bigint;                /* hwinv_history size before pruning */
    tbl_size_after    bigint;                /* hwinv_history size after pruning */
    db_size_before    bigint;                /* DB size before pruning */
    db_size_after     bigint;                /* DB size after pruning */
BEGIN
    /* Capture size before pruning */
    SELECT pg_database_size(current_database()) INTO db_size_before;
    SELECT pg_total_relation_size('hwinv_hist') INTO tbl_size_before;
    
    /* Loop through every unique xname + FRUID pair in the event history */
    FOR unique_ids IN SELECT distinct id,fru_id FROM hwinv_hist LOOP

        /* Reset the deletion count for this id pair */
        del_this_id_pair := 0;

        /* For this unique pair of ids, select the first event in time */
        SELECT * INTO base_event FROM hwinv_hist WHERE id = unique_ids.id AND fru_id = unique_ids.fru_id ORDER BY "timestamp" ASC LIMIT 1;

        /* Starting at the second event for this pair, loop through their remaining events */
        FOR next_event IN SELECT * FROM hwinv_hist WHERE id = unique_ids.id AND fru_id = unique_ids.fru_id AND "timestamp" != base_event.timestamp ORDER BY "timestamp" ASC LOOP

            /* If the event type is 'Detected' and the two events match, delete it */
            IF base_event.event_type = 'Detected' AND base_event.event_type = next_event.event_type THEN

                DELETE FROM hwinv_hist WHERE id = next_event.id AND fru_id = next_event.fru_id AND "timestamp" = next_event.timestamp;

                /* Increment the deletion counts */
                del_total := del_total + 1;
                del_this_id_pair := del_this_id_pair + 1;

            ELSE

                /* Otherwise, set the base event to the next event and continue */
                base_event := next_event;

            END IF;

        END LOOP;

        /* Update the deletion count (if non-zero) for this pair into the map for the xname */
        if del_this_id_pair > 0 THEN
            del_per_id := del_per_id || jsonb_build_object(
                unique_ids.id,
                COALESCE((del_per_id ->> unique_ids.id)::int, 0) + del_this_id_pair
            );
        END IF;

    END LOOP;

    /* Capture size after pruning */
    SELECT pg_database_size(current_database()) INTO db_size_after;
    SELECT pg_total_relation_size('hwinv_hist') INTO tbl_size_after;

    /* Output statistics */

    IF del_total = 0 THEN
        RAISE NOTICE 'No duplicate Detected events found.';
    ELSE
        RAISE NOTICE 'Deleted events per xname:';
        RAISE NOTICE '';

        FOR pair IN SELECT * FROM jsonb_each_text(del_per_id) LOOP
            RAISE NOTICE E'\t%:\t%', pair.key, pair.value;
        END LOOP;

        RAISE NOTICE '';
        RAISE NOTICE 'Removed % duplicate Detected events.', del_total;
        RAISE NOTICE '';
        RAISE NOTICE 'hwinv_history table size before pruning: % mb', tbl_size_before / 1024 / 1024;
        RAISE NOTICE 'hwinv_history table size after pruning:  % mb', tbl_size_after / 1024 / 1024;
        RAISE NOTICE '';
        RAISE NOTICE 'Database size before pruning: % mb', db_size_before / 1024 / 1024;
        RAISE NOTICE 'Database size after pruning:  % mb', db_size_after / 1024 / 1024;
        RAISE NOTICE '';
    END IF;
END;
$$ LANGUAGE plpgsql;
SQL
)

if [[ $? -ne 0 ]]; then
  echo "Error executing SQL command" >&2
  exit 1
fi

echo "$psql_output"
