-- Copyright 2020 Hewlett Packard Enterprise LP
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