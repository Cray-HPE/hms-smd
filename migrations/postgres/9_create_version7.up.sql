-- Copyright 2020 Cray Inc. All Rights Reserved.
--
-- Except as permitted by contract or express written permission of Cray Inc.,
-- no part of this work or its content may be modified, used, reproduced or
-- disclosed in any form. Modifications made without express permission of
-- Cray Inc. may damage the system the software is installed within, may
-- disqualify the user from receiving support from Cray Inc. under support or
-- maintenance contracts, or require additional support services outside the
-- scope of those contracts to repair the software or system.

-- Install hmsds schema version 7

BEGIN;

create table if not exists hwinv_hist (
    "id"         VARCHAR(63),
    "fru_id"     VARCHAR(128),
    "event_type" VARCHAR(128),
    "timestamp"  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Bump the schema version
insert into system values(0, 7, '{}'::JSON)
    on conflict(id) do update set schema_version=7;

COMMIT;