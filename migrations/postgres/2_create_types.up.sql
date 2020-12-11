-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

BEGIN;

-- Used in table 'component_groups' table in version 1 schema.
-- Note: No create type if exists option.  So this will generate an error,
--       on a second attempt, though it should/need be run only once.
CREATE TYPE group_type AS ENUM ('partition', 'exclusive', 'shared');
CREATE TYPE group_namespace AS ENUM ('partition', 'group');

COMMIT;