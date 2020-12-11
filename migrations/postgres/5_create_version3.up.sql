-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Install hmsds schema version 3

BEGIN;

-- Component power mapping

create table if not exists power_mapping (
    "id"         VARCHAR(63)   PRIMARY KEY NOT NULL, -- RFC 1035 max hostname 63 chars
    "powered_by" VARCHAR(63)[] NOT NULL
);

-- Bump the schema version
insert into system values(0, 3, '{}'::JSON)
    on conflict(id) do update set schema_version=3;

COMMIT;