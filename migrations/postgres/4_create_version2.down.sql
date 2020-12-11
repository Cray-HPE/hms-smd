-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Completely remove schema 2, i.e. everything, if we roll back this far.

BEGIN;

DROP TABLE IF EXISTS component_lock_members;
DROP TABLE IF EXISTS component_locks;

-- Decrease the schema version
INSERT INTO system VALUES(0, 1, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=1;

COMMIT;
