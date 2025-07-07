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

set -eu

DB_NAME="hmsds"
DB_USER="postgres"

# PGPASSWORD must be set, or exit
#if [[ -z "${PGPASSWORD:-}" ]]; then
#  echo "Environment variable PGPASSWORD is required but not set" >&2
#  exit 1
#fi

output=$(
psql "dbname=$DB_NAME user=$DB_USER" 2>&1 <<'SQL'
DO $$
DECLARE
    unique_ids RECORD;
    fru_event1 RECORD;
    fru_event2 RECORD;
    deleted    bigint := 0;
BEGIN
    FOR unique_ids IN SELECT distinct id,fru_id FROM hwinv_hist LOOP
        SELECT * INTO fru_event1 FROM hwinv_hist WHERE id = unique_ids.id AND fru_id = unique_ids.fru_id ORDER BY "timestamp" ASC LIMIT 1;
        FOR fru_event2 IN SELECT * FROM hwinv_hist WHERE id = unique_ids.id AND fru_id = unique_ids.fru_id AND "timestamp" != fru_event1.timestamp ORDER BY "timestamp" ASC LOOP
            IF fru_event1.event_type = 'Detected' AND fru_event1.event_type = fru_event2.event_type THEN
                DELETE FROM hwinv_hist WHERE id = fru_event2.id AND fru_id = fru_event2.fru_id AND "timestamp" = fru_event2.timestamp;
                deleted := deleted + 1;
            ELSE
                fru_event1 := fru_event2;
            END IF;
        END LOOP;
    END LOOP;
    RAISE NOTICE 'Removed % duplicate Detected events.', deleted;
END;
$$ LANGUAGE plpgsql;
SQL
)

if [[ $? -ne 0 ]]; then
  echo "Error executing SQL command" >&2
  exit 1
fi

echo "psql output:"
echo "$output"
