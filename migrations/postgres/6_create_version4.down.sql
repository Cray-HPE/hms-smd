-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Completely remove schema 4 if we roll back this far.

BEGIN;

DROP TABLE IF EXISTS job_sync;
DROP TABLE IF EXISTS job_state_rf_poll;

-- Decrease the schema version
INSERT INTO system VALUES(0, 3, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=3;

COMMIT;