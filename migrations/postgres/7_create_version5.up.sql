-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Install hmsds schema version 5

BEGIN;

-- Replaces the comp_endpoints_info to include the enabled field from rf_endpoints
CREATE OR REPLACE VIEW comp_endpoints_info AS
SELECT 
    comp_endpoints.id              AS  "id",
    comp_endpoints.type            AS  "type",
    comp_endpoints.domain          AS  "domain",
    comp_endpoints.redfish_type    AS  "redfish_type",
    comp_endpoints.redfish_subtype AS  "redfish_subtype",
    comp_endpoints.mac             AS  "mac",
    comp_endpoints.uuid            AS  "uuid",
    comp_endpoints.odata_id        AS  "odata_id",
    comp_endpoints.rf_endpoint_id  AS  "rf_endpoint_id",
    rf_endpoints.fqdn              AS  "rf_endpoint_fqdn",
    comp_endpoints.component_info  AS  "component_info",  -- JSON
    rf_endpoints.user              AS  "rf_endpoint_user",
    rf_endpoints.password          AS  "rf_endpoint_password",
    rf_endpoints.enabled           As  "enabled"
FROM comp_endpoints
LEFT JOIN rf_endpoints on comp_endpoints.rf_endpoint_id = rf_endpoints.id;

-- Bump the schema version
insert into system values(0, 5, '{}'::JSON)
    on conflict(id) do update set schema_version=5;

COMMIT;