// Copyright 2019-2020 Hewlett Packard Enterprise Development LP

package hmsds

import (
	"testing"
)

// Convert MySQL-style prepared query args to numbered Postgres-style.
func TestToPGQueryArgs(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{{
		in:  "update asdf set a = ?, b = ? where c = ? and d = ?",
		out: "update asdf set a = $1, b = $2 where c = $3 and d = $4",
	}, {
		in:  "update asdf set a = ?, b = ? where c = ? and d = ?;",
		out: "update asdf set a = $1, b = $2 where c = $3 and d = $4;",
	}, {
		in:  "?, b = ? where c = ? and d = ?;",
		out: "$1, b = $2 where c = $3 and d = $4;",
	}, {
		in:  "update aabb set aa = 'hello' where id = 'someid';",
		out: "update aabb set aa = 'hello' where id = 'someid';",
	}, {
		in:  "???",
		out: "$1$2$3",
	}, {
		in:  "?",
		out: "$1",
	}, {
		in:  "'?'",
		out: "'$1'",
	}, {
		in: `
    INSERT INTO components (
        id, type, state, flag, enabled, admin, role, subrole, nid, subtype, nettype, arch, class, reservation_disabled, locked)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    ON DUPLICATE KEY UPDATE 
        state=VALUES(state),
        flag=VALUES(flag),
        subtype=VALUES(subtype),
        nettype=VALUES(nettype),
        arch=VALUES(arch),
        class=VALUES(class);`,

		out: `
    INSERT INTO components (
        id, type, state, flag, enabled, admin, role, subrole, nid, subtype, nettype, arch, class, reservation_disabled, locked)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
    ON DUPLICATE KEY UPDATE 
        state=VALUES(state),
        flag=VALUES(flag),
        subtype=VALUES(subtype),
        nettype=VALUES(nettype),
        arch=VALUES(arch),
        class=VALUES(class);`,
	}}
	for i, pair := range tests {
		gotStr := ToPGQueryArgs(pair.in)
		if gotStr != pair.out {
			t.Errorf("Testcase %d (%s): FAIL: Expected %s, but got %s",
				i, pair.in, pair.out, gotStr)
		} else {
			t.Logf("Testcase %d (%s): PASS: Expected %s and got %s",
				i, pair.in, pair.out, gotStr)
		}
	}
}
