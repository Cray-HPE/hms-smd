-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Install hmsds schema version 2

BEGIN;

create table if not exists component_locks (
    "id"       UUID PRIMARY KEY,
    "created"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "reason"   VARCHAR(255) NOT NULL, 
    "owner"    VARCHAR(255) NOT NULL,
    "lifetime" INT NOT NULL
);

create table if not exists component_lock_members (
    "component_id"    VARCHAR(63) PRIMARY KEY,
    "lock_id"         UUID NOT NULL,
    FOREIGN KEY ("component_id") REFERENCES components ("id") ON DELETE CASCADE,
    FOREIGN KEY ("lock_id") REFERENCES component_locks ("id") ON DELETE CASCADE
);

-- Bump the schema version
insert into system values(0, 2, '{}'::JSON)
    on conflict(id) do update set schema_version=2;

COMMIT;