-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Completely remove schema 1, i.e. everything, if we roll back this far.
SET client_min_messages = ERROR;

BEGIN;

DROP TABLE IF EXISTS component_group_members;
DROP TABLE IF EXISTS component_groups;

DROP table if exists scn_subscriptions;

DROP VIEW IF EXISTS hwinv_by_loc_with_fru;
DROP table if exists hwinv_by_loc;
DROP table if exists hwinv_by_fru;

DROP table if exists discovery_status;

DROP VIEW IF EXISTS service_endpoints_info;
DROP table if exists service_endpoints;

DROP VIEW IF EXISTS comp_endpoints_info;
DROP TABLE IF EXISTS comp_endpoints;

DROP TABLE IF EXISTS rf_endpoints;

DROP TABLE IF EXISTS node_nid_mapping;

DROP TABLE IF EXISTS components;
DROP table if exists system;


COMMIT;