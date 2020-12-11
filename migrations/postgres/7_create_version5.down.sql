-- Copyright (c) 2019 Cray Inc. All Rights Reserved.

-- Completely remove schema 5 if we roll back this far.

BEGIN;

-- Revert this view to the previous version
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
    rf_endpoints.password          AS  "rf_endpoint_password"
FROM comp_endpoints
LEFT JOIN rf_endpoints on comp_endpoints.rf_endpoint_id = rf_endpoints.id;

-- Decrease the schema version
INSERT INTO system VALUES(0, 4, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=4;

COMMIT;