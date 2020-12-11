-- Copyright 2020 Cray Inc.

-- Completely remove schema 13 if we roll back this far.

BEGIN;

DROP FUNCTION hwinv_by_loc_update_parents();

DROP VIEW IF EXISTS hwinv_by_loc_with_partition;

ALTER TABLE hwinv_by_loc
DROP COLUMN "parent_node";

-- Decrease the schema version
INSERT INTO system VALUES(0, 12, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=12;

COMMIT;