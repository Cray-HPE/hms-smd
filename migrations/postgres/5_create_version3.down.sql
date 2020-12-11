-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Completely remove schema 3.

BEGIN;

DROP TABLE IF EXISTS power_mapping;

-- Decrease the schema version
INSERT INTO system VALUES(0, 2, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=2;

COMMIT;
