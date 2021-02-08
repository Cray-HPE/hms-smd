-- Copyright 2020 Hewlett Packard Enterprise LP
-- Modifies SMD to have the new locking schema design.

BEGIN;

-- Add the new fields we need in `components`.

ALTER TABLE components
    ADD COLUMN IF NOT EXISTS reservation_disabled bool DEFAULT FALSE NOT NULL,
    ADD COLUMN IF NOT EXISTS locked bool DEFAULT FALSE NOT NULL;

-- Add the new `reservations` table.
CREATE TABLE IF NOT EXISTS reservations
(
    component_id          VARCHAR(63)              NOT NULL
        CONSTRAINT locks_component_id_pk
            PRIMARY KEY
        CONSTRAINT locks_hardware_component_id_fk
            REFERENCES components ("id")
            ON DELETE CASCADE,
    create_timestamp      TIMESTAMP WITH TIME ZONE NOT NULL,
    expiration_timestamp  TIMESTAMP WITH TIME ZONE,
    deputy_key            VARCHAR,
    reservation_key       VARCHAR,
    v1_lock_id            UUID
);

CREATE INDEX locks_create_timestamp_index
    ON reservations (create_timestamp);

CREATE INDEX locks_expiration_timestamp_index
    ON reservations (expiration_timestamp);

CREATE INDEX locks_reservation_key_index
    ON reservations (reservation_key);

CREATE INDEX locks_deputy_key_index
    ON reservations (deputy_key);

-- Migrate any existing v1/locks to v2/locks
CREATE OR REPLACE FUNCTION comp_lock_update_reservations()
RETURNS VOID AS $$
DECLARE
    lock_member RECORD;
BEGIN
    FOR lock_member IN SELECT
        component_lock_members.component_id AS "comp_id",
        component_lock_members.lock_id AS "lock_id",
        component_locks.created AS "created",
        component_locks.lifetime AS "lifetime"
    FROM component_lock_members LEFT JOIN component_locks ON component_lock_members.lock_id = component_locks.id LOOP
        INSERT INTO reservations (
            component_id, create_timestamp, expiration_timestamp, deputy_key, reservation_key, v1_lock_id)
        VALUES (
            lock_member.comp_id,
            lock_member.created,
            lock_member.created + (lock_member.lifetime || ' seconds')::interval,
            lock_member.comp_id || ':dk:' || lock_member.lock_id::text,
            lock_member.comp_id || ':rk:' || lock_member.lock_id::text,
            lock_member.lock_id);
    END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT comp_lock_update_reservations();

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

-- Bump the schema version
INSERT INTO system
VALUES (0, 14, '{}'::JSON)
ON CONFLICT(id) DO UPDATE SET schema_version=14;

COMMIT;