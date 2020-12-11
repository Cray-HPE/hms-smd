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

-- Remove the new IP Address column
ALTER TABLE comp_eth_interfaces
DROP COLUMN "ip_addresses";

-- Decrease the schema version
INSERT INTO system VALUES(0, 14, '{}'::JSON)
    ON CONFLICT(id) DO UPDATE SET schema_version=14;

COMMIT;