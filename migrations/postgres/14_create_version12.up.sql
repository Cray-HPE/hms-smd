-- Copyright 2020 Cray Inc. All Rights Reserved.
--
-- Except as permitted by contract or express written permission of Cray Inc.,
-- no part of this work or its content may be modified, used, reproduced or
-- disclosed in any form. Modifications made without express permission of
-- Cray Inc. may damage the system the software is installed within, may
-- disqualify the user from receiving support from Cray Inc. under support or
-- maintenance contracts, or require additional support services outside the
-- scope of those contracts to repair the software or system.

-- Install hmsds schema version 12

BEGIN;

CREATE TABLE IF NOT EXISTS hsn_interfaces (
    nic         VARCHAR(32) PRIMARY KEY,
    macaddr     VARCHAR(32) NOT NULL DEFAULT '',
    hsn         VARCHAR(32) NOT NULL DEFAULT '',
    node        VARCHAR(32) NOT NULL DEFAULT '',
    ipaddr      VARCHAR(64) NOT NULL DEFAULT '',
    last_update TIMESTAMPTZ
);

-- Bump the schema version
insert into system values(0, 12, '{}'::JSON)
    on conflict(id) do update set schema_version=12;

COMMIT;