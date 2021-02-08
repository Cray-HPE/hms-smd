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
-- Modifies SMD to have the new CompEthInterface schema.

BEGIN;

-- Add the new IP Address column
ALTER TABLE comp_eth_interfaces
ADD COLUMN "ip_addresses" JSON NOT NULL DEFAULT '[]'::JSON;


-- Migrate IP Addresses to new format
CREATE OR REPLACE FUNCTION comp_ethernet_interfaces_update()
RETURNS VOID AS $$
DECLARE
    comp_ethernet_interface RECORD;
BEGIN
    FOR comp_ethernet_interface IN SELECT id, json_build_array(json_build_object('IPAddress', ipaddr, 'Network', '')) as ip_addresses
                                   FROM comp_eth_interfaces
                                   WHERE ipaddr != ''
        LOOP
            UPDATE comp_eth_interfaces
            SET ip_addresses = comp_ethernet_interface.ip_addresses
            WHERE id = comp_ethernet_interface.id;
        END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT comp_ethernet_interfaces_update();

-- Remove the old IP Address column
ALTER TABLE comp_eth_interfaces
DROP COLUMN "ipaddr";

-- Bump the schema version
insert into system values(0, 15, '{}'::JSON)
    on conflict(id) do update set schema_version=15;

COMMIT;