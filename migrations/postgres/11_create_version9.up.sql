-- Copyright 2020 Cray Inc. All Rights Reserved.
--
-- Except as permitted by contract or express written permission of Cray Inc.,
-- no part of this work or its content may be modified, used, reproduced or
-- disclosed in any form. Modifications made without express permission of
-- Cray Inc. may damage the system the software is installed within, may
-- disqualify the user from receiving support from Cray Inc. under support or
-- maintenance contracts, or require additional support services outside the
-- scope of those contracts to repair the software or system.

-- Install hmsds schema version 9

BEGIN;

ALTER TABLE rf_endpoints
ADD COLUMN ipAddr VARCHAR(64) NOT NULL DEFAULT '';

-- Bump the schema version
insert into system values(0, 9, '{}'::JSON)
    on conflict(id) do update set schema_version=9;

COMMIT;