/*
 * MIT License
 *
 * (C) Copyright [2018-2021] Hewlett Packard Enterprise Development LP
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included
 * in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 * OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 * ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 * OTHER DEALINGS IN THE SOFTWARE.
 */

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