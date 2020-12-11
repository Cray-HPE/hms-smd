-- Copyright 2020 Hewlett Packard Enterprise LP

-- Modifies SMD to revert the new locking schema design.

BEGIN;

DROP FUNCTION IF EXISTS comp_lock_update_reservations;

--drop the new table
DROP TABLE IF EXISTS reservations;


-- remove the added fields in `components`.
ALTER TABLE components
    DROP COLUMN IF EXISTS reservation_disabled,
    DROP COLUMN IF EXISTS locked;

-- Decrease the schema version
INSERT INTO system VALUES(0, 13, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=13;

COMMIT;