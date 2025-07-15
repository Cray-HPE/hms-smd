/*
 * MIT License
 *
 * (C) Copyright [2025] Hewlett Packard Enterprise Development LP
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */
-- Removes duplicate "Detected" events from the hardware history table

-- For each unique id (xname) and fru_id combination, grab the earliest event
-- for the pair to use as a 'base' event. Then compare the remaining events
-- for the id pair (in time order) to this base event. If the event_type of
-- the next event matches the base event and is a "Detected" event, delete it
-- and move to the event after it, and so on. If they do not match then set
-- the base event to this next event and continue. This will remove all but
-- the first occurrence of each "Detected" event that does not come after a
-- non-"Detected" event for the same id (xname) and fru_id pair.

BEGIN;

CREATE OR REPLACE FUNCTION hwinv_hist_remove_duplicate_detected_events()
RETURNS VOID AS $$
BEGIN
    -- Create a temporary timestamp index to speed up pruning

    CREATE INDEX IF NOT EXISTS hwinvhist_id_ts_idx ON hwinv_hist (id, "timestamp");

    -- Run the pruning logic

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

    -- Drop the temporary index to free up associated resources

    DROP INDEX IF EXISTS hwinvhist_id_ts_idx;
END;
$$ LANGUAGE plpgsql;

-- A full vacuum must be run to reclaim space but cannot run in a migration.
-- The cray-smd-init service will run it manually after the migration is complete.

-- Bump the schema version
insert into system values(0, 21, '{}'::JSON)
    on conflict(id) do update set schema_version=21;

COMMIT;