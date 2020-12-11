-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Completely remove schema 6 if we roll back this far.

BEGIN;

ALTER TABLE components
DROP COLUMN "subrole";

ALTER TABLE node_nid_mapping
DROP COLUMN "subrole";

-- Decrease the schema version
INSERT INTO system VALUES(0, 5, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=5;

COMMIT;