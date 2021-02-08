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