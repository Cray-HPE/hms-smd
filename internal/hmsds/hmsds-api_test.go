// Copyright 2019-2020 Hewlett Packard Enterprise Development LP

package hmsds

import (
	"context"
	"log"
	"os"
	"testing"

	"stash.us.cray.com/HMS/hms-smd/pkg/sm"

	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
)

//////////////////////////////////////////////////////////////////////////////
//
// Global initialization for all DB drivers and shared helper functions
//
//////////////////////////////////////////////////////////////////////////////

// Postgres driver
var dPG hmsdbPg
var mockPG sqlmock.Sqlmock

// Compare arrays of xnames to make sure every name exists in both lists (order doesn't matter)
func compareIDs(ids1, ids2 []string) bool {
	if len(ids1) != len(ids2) {
		return false
	}
	for i := 0; i < len(ids1); i++ {
		found := false
		for j := 0; j < len(ids2); j++ {
			if ids1[i] == ids2[j] {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func compareSCNSubs(subs1, subs2 *sm.SCNSubscriptionArray) bool {
	if subs1 == nil && subs2 == nil {
		return true
	} else if subs1 == nil || subs2 == nil {
		return false
	}
	if len(subs1.SubscriptionList) != len(subs2.SubscriptionList) {
		return false
	}
	for i, sub1 := range subs1.SubscriptionList {
		sub2 := subs2.SubscriptionList[i]
		if (sub1.Subscriber == sub2.Subscriber) &&
			(sub1.Url == sub2.Url) &&
			(len(sub1.States) == len(sub2.States)) {
			for j, state1 := range sub1.States {
				state2 := sub2.States[j]
				if state1 != state2 {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

//
// Unit Tests
//

// Set up for both drivers and then run all tests.  The DB-specific parts
// use separate HMSDB instances of their internal type, so we can set
// everything up beforehand and not have to run each types tests separately.
func TestMain(m *testing.M) {

	excode := 1

	InitializeMockDB()

	// Run tests for all drivers
	excode = m.Run()

	// Postgres cleanup
	dPG.Close()
	os.Exit(excode)
}

func InitializeMockDB() {
	var err error

	// Postgres setup.
	dPG.dsn = "user=hmsdsuser dbname=hmsds"
	dPG.db = nil
	dPG.connected = true
	dPG.lgLvl = LOG_DEBUG
	dPG.db, mockPG, err = sqlmock.New()
	if err != nil {
		dPG.LogAlways("Error: Open(): an error '%s' was not expected when opening a stub database connection", err)
		os.Exit(1)
	}
	dPG.lg = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
	dPG.sc = sq.NewStmtCache(dPG.db)
	dPG.ctx = context.TODO()
}

func ResetMockDB() {
	// Close the previous mock DB
	dPG.Close()
	// Create the new mock DB
	InitializeMockDB()
}
