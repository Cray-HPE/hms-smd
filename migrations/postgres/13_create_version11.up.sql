-- Copyright 2020 Cray Inc. All Rights Reserved.
--
-- Except as permitted by contract or express written permission of Cray Inc.,
-- no part of this work or its content may be modified, used, reproduced or
-- disclosed in any form. Modifications made without express permission of
-- Cray Inc. may damage the system the software is installed within, may
-- disqualify the user from receiving support from Cray Inc. under support or
-- maintenance contracts, or require additional support services outside the
-- scope of those contracts to repair the software or system.

-- Install hmsds schema version 11

BEGIN;

CREATE TABLE IF NOT EXISTS comp_eth_interfaces (
    id          VARCHAR(32) PRIMARY KEY,
    description TEXT,
    macaddr     VARCHAR(32) NOT NULL,
    ipaddr      VARCHAR(64) NOT NULL DEFAULT '',
    last_update TIMESTAMPTZ,
    compid      VARCHAR(63) NOT NULL DEFAULT '',
    comptype    VARCHAR(63) NOT NULL DEFAULT ''
);

-- Bump the schema version
insert into system values(0, 11, '{}'::JSON)
    on conflict(id) do update set schema_version=11;

COMMIT;