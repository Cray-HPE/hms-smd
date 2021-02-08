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

-- Install hmsds schema version 10

BEGIN;

-- Need to drop the view first to alter the tables it references
DROP VIEW IF EXISTS hwinv_by_loc_with_fru;

-- Increase the size of the fru_id in all tables
ALTER TABLE hwinv_by_loc
ALTER COLUMN "fru_id" TYPE VARCHAR(255);

ALTER TABLE hwinv_by_fru
ALTER COLUMN "fru_id" TYPE VARCHAR(255);

ALTER TABLE hwinv_hist
ALTER COLUMN "fru_id" TYPE VARCHAR(255);

-- Replace the view
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

-- Bump the schema version
insert into system values(0, 10, '{}'::JSON)
    on conflict(id) do update set schema_version=10;

COMMIT;