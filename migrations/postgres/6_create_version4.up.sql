-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Install hmsds schema version 4

BEGIN;

create table if not exists job_sync (
    "id"           UUID PRIMARY KEY,
    "type"         VARCHAR(128),
    "status"       VARCHAR(128),
    "last_update"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "lifetime"     INT NOT NULL
);

create table if not exists job_state_rf_poll (
    "comp_id"      VARCHAR(63) PRIMARY KEY,
    "job_id"       UUID NOT NULL,
    FOREIGN KEY ("job_id") REFERENCES job_sync ("id") ON DELETE CASCADE
);

-- Bump the schema version
insert into system values(0, 4, '{}'::JSON)
    on conflict(id) do update set schema_version=4;

COMMIT;