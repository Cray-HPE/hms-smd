-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Install hmsds schema version 6

BEGIN;

ALTER TABLE components
ADD COLUMN "subrole" VARCHAR(32) NOT NULL DEFAULT '';

ALTER TABLE node_nid_mapping
ADD COLUMN "subrole" VARCHAR(32) NOT NULL DEFAULT '';

-- Bump the schema version
insert into system values(0, 6, '{}'::JSON)
    on conflict(id) do update set schema_version=6;

COMMIT;