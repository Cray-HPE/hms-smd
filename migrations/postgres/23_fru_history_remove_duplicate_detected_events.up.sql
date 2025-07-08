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

BEGIN;

-- For each unique id (xname) and fru_id combination, grab the earliest event
-- for the pair to use as a 'base' event. Then compare the remaining events
-- for the id pair (in time order) to this base event. If the event_type of
-- the next event matches the base event and is a "Detected" event, delete it
-- and move to the event after it, and so on. If they do not match then set
-- the base event to this next event and continue. This will remove all but
-- the first occurrence of each "Detected" event that does not come after a
-- non-"Detected" event for the same id (xname) and fru_id pair.

CREATE OR REPLACE FUNCTION hwinv_hist_remove_duplicate_detected_events()
RETURNS VOID AS $$
DECLARE
    unique_ids RECORD;
    base_event RECORD;
    next_event RECORD;
BEGIN
    /* Loop through every unique xname + FRUID pair in the event history */
    FOR unique_ids IN SELECT distinct id,fru_id FROM hwinv_hist LOOP

        /* For this unique pair of ids, select the first event in time */
        SELECT * INTO base_event FROM hwinv_hist WHERE id = unique_ids.id AND fru_id = unique_ids.fru_id ORDER BY "timestamp" ASC LIMIT 1;

        /* Starting at the second event for this pair, loop through their remaining events */
        FOR next_event IN SELECT * FROM hwinv_hist WHERE id = unique_ids.id AND fru_id = unique_ids.fru_id AND "timestamp" != base_event.timestamp ORDER BY "timestamp" ASC LOOP

            /* If the event type is 'Detected' and the two events match, delete it */
            IF base_event.event_type = 'Detected' AND base_event.event_type = next_event.event_type THEN

                DELETE FROM hwinv_hist WHERE id = next_event.id AND fru_id = next_event.fru_id AND "timestamp" = next_event.timestamp;

            ELSE

                /* Otherwise, set the base event to the next event and continue */
                base_event := next_event;

            END IF;
        END LOOP;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Bump the schema version
insert into system values(0, 21, '{}'::JSON)
    on conflict(id) do update set schema_version=21;

COMMIT;
