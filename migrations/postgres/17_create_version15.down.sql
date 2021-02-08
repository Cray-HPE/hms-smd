-- Copyright 2020 Hewlett Packard Enterprise LP
-- Reverts SMD to the old CompEthInterface schema.

BEGIN;

-- Restore the ipaddr column
ALTER TABLE comp_eth_interfaces
ADD COLUMN "ipaddr" VARCHAR(64) NOT NULL DEFAULT '';


-- Migrate IP Address information to the old format
-- WARNING: There is potentail for dataloss, as the older Database schema does not
-- support multiple IP Addresses per mac adresses. For the most common case
CREATE OR REPLACE FUNCTION comp_ethernet_interfaces_downgrade()
RETURNS VOID AS $$
DECLARE
    comp_ethernet_interface RECORD;
BEGIN
    FOR comp_ethernet_interface IN  SELECT id, ip_addresses->0->>'IPAddress' AS ipaddr
                                    FROM comp_eth_interfaces
                                    WHERE json_array_length(ip_addresses) != 0
        LOOP
            UPDATE comp_eth_interfaces
            SET ipaddr = comp_ethernet_interface.ipaddr
            WHERE id = comp_ethernet_interface.id;
        END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT comp_ethernet_interfaces_downgrade();

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

-- Remove the new IP Address column
ALTER TABLE comp_eth_interfaces
DROP COLUMN "ip_addresses";

-- Decrease the schema version
INSERT INTO system VALUES(0, 14, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=14;

COMMIT;