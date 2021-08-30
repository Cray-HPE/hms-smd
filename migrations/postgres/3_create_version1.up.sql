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
-- Install hmsds schema version 1

SET client_min_messages = ERROR;

BEGIN;

--
-- System table - Global parameters
--
create table if not exists system (
    "id"             INT PRIMARY KEY NOT NULL,
    "schema_version" INT UNIQUE NOT NULL,
    "system_info"    JSON     -- JSON blob
);

--
-- Components - State Management and other higher-level activities
--

create table if not exists components (
    "id"          VARCHAR(63) PRIMARY KEY NOT NULL, -- RFC 1035 max hostname 63 chars
    "type"        VARCHAR(63)     NOT NULL, 
    "state"       VARCHAR(32)  NOT NULL,
    "admin"       VARCHAR(32)  NOT NULL DEFAULT '',
    "enabled"     BOOL         NOT NULL DEFAULT '1',
    "flag"        VARCHAR(32)  NOT NULL,
    "role"        VARCHAR(32)  NOT NULL,
    "nid"         BIGINT       NOT NULL,  -- -1 if no NID
    "subtype"     VARCHAR(64)  NOT NULL,
    "nettype"     VARCHAR(64)  NOT NULL,
    "arch"        VARCHAR(64)  NOT NULL,
    "disposition" VARCHAR(64)  NOT NULL  DEFAULT ''
);


-- Component xname->NID map

create table if not exists node_nid_mapping (
    "id"        VARCHAR(63) PRIMARY KEY NOT NULL, -- RFC 1035 max hostname 63 chars
    "nid"       BIGINT       UNIQUE,           -- 
    "role"      VARCHAR(32)  NOT NULL,
    "name"      VARCHAR(32)  NOT NULL DEFAULT '',
    "node_info" JSON                           -- JSON blob
);

--
-- Redfish Endpoint Management
--

create table if not exists rf_endpoints (
     "id"                 VARCHAR(63) PRIMARY KEY, -- RFC 1035 max hostname 63 chars
     "type"               VARCHAR(63) NOT NULL,
     "name"               TEXT,              -- Descriptive user-supplied name
     "hostname"           VARCHAR(63),          -- RFC 1035 max hostname 63 chars
     "domain"             VARCHAR(192),         -- RFC 1035 max FQDN 255 chars
     "fqdn"               VARCHAR(255) UNIQUE,
     "ip_info"            JSON DEFAULT '{}'::JSON,   -- JSON blob
     "enabled"            BOOL,
     "uuid"               VARCHAR(64),
     "user"               VARCHAR(128),
     "password"           VARCHAR(128),
     useSSDP              BOOL,
     macRequired          BOOL,
     macAddr              VARCHAR(32),
     rediscoverOnUpdate   BOOL,
     templateID           VARCHAR(128),
     discovery_info       JSON               -- JSON blob
);

--
--  Redfish Endpoints - component discovery info, for interfacing with Redfish
--

create table if not exists comp_endpoints (
    "id"              VARCHAR(63) PRIMARY KEY,  -- RFC 1035 max hostname 63 chars
    "type"            VARCHAR(63)  NOT NULL,
    "domain"          VARCHAR(192) NOT NULL,    -- RFC 1035 max FQDN 255 chars
    "redfish_type"    VARCHAR(63) NOT NULL,
    "redfish_subtype" VARCHAR(63) NOT NULL,
    "rf_endpoint_id"  VARCHAR(63) NOT NULL,
    "mac"             VARCHAR(32),
    "uuid"            VARCHAR(64),
    "odata_id"        VARCHAR(512) NOT NULL,
    "component_info"  JSON,                  -- JSON blob
    FOREIGN KEY("rf_endpoint_id") REFERENCES rf_endpoints("id") ON DELETE CASCADE
);


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

create table if not exists service_endpoints (
    "rf_endpoint_id"  VARCHAR(63),
    "redfish_type"    VARCHAR(63),
    "redfish_subtype" VARCHAR(63) NOT NULL,
    "uuid"            VARCHAR(64),
    "odata_id"        VARCHAR(512) NOT NULL,
    "service_info"    JSON,         -- JSON blob
    PRIMARY KEY(rf_endpoint_id, redfish_type),
    FOREIGN KEY("rf_endpoint_id") REFERENCES rf_endpoints("id") ON DELETE CASCADE
);


CREATE OR REPLACE VIEW service_endpoints_info AS
SELECT
    service_endpoints.rf_endpoint_id  AS  "rf_endpoint_id",
    service_endpoints.redfish_type    AS  "redfish_type",
    service_endpoints.redfish_subtype AS  "redfish_subtype",
    service_endpoints.uuid            AS  "uuid",
    service_endpoints.odata_id        AS  "odata_id",
    rf_endpoints.fqdn                 AS  "rf_endpoint_fqdn",
    service_endpoints.service_info    AS  "service_info"  -- JSON
FROM service_endpoints
LEFT JOIN rf_endpoints on service_endpoints.rf_endpoint_id = rf_endpoints.id;

create table if not exists discovery_status (
    "id"          INT PRIMARY KEY NOT NULL,
    "status"      VARCHAR(128),
    "last_update" TIMESTAMPTZ, 
    "details"     JSON   -- JSON blob 
);

insert into discovery_status values(0, 'NotStarted', NOW(), '{}'::JSON)
    on conflict(id) do update set id = 0;

--
--  Hardware Inventory Tables
--

create table if not exists hwinv_by_fru (
    "fru_id"        VARCHAR(128) PRIMARY KEY NOT NULL,  --
    "type"          VARCHAR(63) NOT NULL,  -- HMS type
    "subtype"       VARCHAR(63) NOT NULL,
    "serial_number" VARCHAR(255) NOT NULL DEFAULT '',
    "part_number"   VARCHAR(255) NOT NULL DEFAULT '',
    "manufacturer"  VARCHAR(255) NOT NULL DEFAULT '',
    "fru_info"      JSON     NOT NULL   -- JSON blob
);

create table if not exists hwinv_by_loc (
    "id"            VARCHAR(63) PRIMARY KEY NOT NULL, -- RFC 1035 max hostname 63 chars
    "type"          VARCHAR(63) NOT NULL,  -- HMS type
    "ordinal"       INT      NOT NULL,  -- Normalized Nth child component, from 0
    "status"        VARCHAR(63) NOT NULL,
    "parent"        VARCHAR(63) NOT NULL DEFAULT '',
    "location_info" JSON,               -- JSON blob 
    "fru_id"        VARCHAR(63),
     FOREIGN KEY("fru_id") REFERENCES hwinv_by_fru("fru_id")
);


CREATE OR REPLACE VIEW hwinv_by_loc_with_fru AS
SELECT
    hwinv_by_loc.id             AS  "id",
    hwinv_by_loc.type           AS  "type",
    hwinv_by_loc.ordinal        AS  "ordinal",
    hwinv_by_loc.status         AS  "status",
    hwinv_by_loc.location_info  AS  "location_info", -- JSON blob
    hwinv_by_loc.fru_id         AS  "fru_id",
    hwinv_by_fru.type           AS  "fru_type",
    hwinv_by_fru.subtype        AS  "fru_subtype",
    hwinv_by_fru.fru_info       AS  "fru_info"       -- JSON blob
FROM hwinv_by_loc
LEFT JOIN hwinv_by_fru ON hwinv_by_loc.fru_id = hwinv_by_fru.fru_id;

create table if not exists scn_subscriptions (
    "id"           SERIAL PRIMARY KEY NOT NULL,
    "sub_url"      VARCHAR(255) UNIQUE NOT NULL,  -- Combined subscriber name an URL for unique sub id
    "subscription" JSON DEFAULT '{}'::JSON
);

-- Used in table 'component_groups' below.
-- Note: No create type if exists option.  So this will generate an error,
--       on a second attempt, though it should/need be run only once.
--CREATE TYPE group_type AS ENUM ('partition', 'exclusive', 'shared');
--CREATE TYPE group_namespace AS ENUM ('partition', 'group');

CREATE TABLE IF NOT EXISTS component_groups (
    "id"               UUID PRIMARY KEY, -- Stored internally as 8 byte binary
    "name"             VARCHAR(255) NOT NULL ,  -- name or label, unique per type
    "description"      VARCHAR(255) NOT NULL,
    "tags"             VARCHAR(255)[],
    "annotations"      JSON DEFAULT '{}'::JSON,  -- reserved for future use
    "type"             group_type, -- ENUM type - Note: partition is a special case of exclusive
    "namespace"        group_namespace,  -- ENUM type:  Group or partition namespace.
    "exclusive_group_identifier" VARCHAR(253),  -- identifier for exclusive type groups
    UNIQUE (name, namespace)
);

CREATE TABLE IF NOT EXISTS component_group_members (
    "component_id"    VARCHAR(63) NOT NULL,
    "group_id"        UUID NOT NULL,
    "group_namespace" VARCHAR(255) NOT NULL, -- Most general that exists of (group name, %exclusive_group%) or special %partition% string
    "joined_at"       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("component_id") REFERENCES components ("id") ON DELETE CASCADE,
    FOREIGN KEY ("group_id") REFERENCES component_groups ("id") ON DELETE CASCADE,
    PRIMARY KEY ("component_id", "group_id"),
    UNIQUE ("component_id", "group_namespace")
);


INSERT INTO system VALUES(0, 1, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=1;

COMMIT;
