// Copyright (c) 2019 Cray Inc. All Rights Reserved.
package main

import (
	"testing"
)

type TypeDecodePair struct {
	utype CompUpdateType
	str   string
}

var goodTests = []TypeDecodePair{
	{StateDataUpdate, "StateData"},
	{FlagOnlyUpdate, "FlagOnly"},
	{EnabledUpdate, "Enabled"},
	{SwStatusUpdate, "SoftwareStatus"},
	{RoleUpdate, "Role"},
	{SingleNIDUpdate, "NID"},
	{StateDataUpdate, "statedata"},
	{FlagOnlyUpdate, "flagonly"},
	{EnabledUpdate, "enabled"},
	{SwStatusUpdate, "softwarestatus"},
	{RoleUpdate, "role"},
	{SingleNIDUpdate, "nid"},
	{StateDataUpdate, "stateData"},
	{FlagOnlyUpdate, "flagOnly"},
	{EnabledUpdate, "enaBled"},
	{SwStatusUpdate, "softwareStatus"},
	{RoleUpdate, "role"},
	{SingleNIDUpdate, "Nid"},
}

// Get the CompUpdateType for a given string, based on its pattern in the recognition
// If the string matches no valid type, CompUpdateInvalid is returned.
func TestGetCompUpdateType(t *testing.T) {
	if gotType := GetCompUpdateType("asdfasdfa"); gotType != CompUpdateInvalid {
		t.Errorf("Testcase a: FAIL: Expected %s, but got %s",
			CompUpdateInvalid.String(), gotType)
	} else {
		t.Logf("Testcase a: PASS: Expected %s and got %s",
			CompUpdateInvalid.String(), gotType)
	}
	if gotType := GetCompUpdateType(""); gotType != CompUpdateInvalid {
		t.Errorf("Testcase b: FAIL: Expected %s, but got %s",
			CompUpdateInvalid.String(), gotType.String())
	} else {
		t.Logf("Testcase b: PASS: Expected %s and got %s",
			CompUpdateInvalid.String(), gotType.String())
	}
	for i, pair := range goodTests {
		gotType := GetCompUpdateType(pair.str)
		if gotType != pair.utype {
			t.Errorf("Testcase %d (%s): FAIL: Expected %s, but got %s",
				i, pair.str, pair.utype, gotType)
		} else {
			t.Logf("Testcase %d (%s): PASS: Expected %s and got %s",
				i, pair.str, pair.utype, gotType)
		}
	}
}

func TestVerifyNormalizeCompUpdateType(t *testing.T) {
	if VerifyNormalizeCompUpdateType("asdfasdfa") != "" {
		t.Errorf("Testcase a: FAIL: Expected '', but got non-empty")
	} else {
		t.Logf("Testcase a: PASS: Expected '' and got ''")
	}
	if VerifyNormalizeCompUpdateType("") != "" {
		t.Errorf("Testcase b: FAIL: Expected '', but got non-empty")
	} else {
		t.Logf("Testcase b: PASS: Expected '' and got ''")
	}
	for i, pair := range goodTests {
		gotType := VerifyNormalizeCompUpdateType(pair.str)
		if gotType != pair.utype.String() {
			t.Errorf("Testcase %d (%s): FAIL: Expected %s, but got %s",
				i, pair.str, pair.utype.String(), gotType)
		} else {
			t.Logf("Testcase %d (%s): PASS: Expected %s and got %s",
				i, pair.str, pair.utype.String(), gotType)
		}
	}
}

func TestDoCompUpdate(t *testing.T) {
	testdataBad1 := &CompUpdate{ComponentIDs: []string{}}
	if err := s.doCompUpdate(testdataBad1, "name"); err != ErrSMDNoIDs {
		t.Errorf("Test 0: Did not get expected error ErrSMDNoIDs")
	}
	testdataBad2 := &CompUpdate{ComponentIDs: []string{"asadfsadf"}}
	if err := s.doCompUpdate(testdataBad2, "name"); err != ErrSMDBadID {
		t.Errorf("Test 1: Did not get expected error ErrSMDBadID")
	}
	testdata := &CompUpdate{ComponentIDs: []string{"x0c0s0b0n0"}}
	testdata.UpdateType = StateDataUpdate.String()
	if err := s.doCompUpdate(testdata, "name"); err != ErrSMDNoState {
		t.Errorf("Test 2: Did not get expected error ErrSMDNoState")
	}
	testdata.UpdateType = FlagOnlyUpdate.String()
	if err := s.doCompUpdate(testdata, "name"); err != ErrSMDNoFlag {
		t.Errorf("Test 3: Did not get expected error ErrSMDNoFlag")
	}
	testdata.UpdateType = EnabledUpdate.String()
	if err := s.doCompUpdate(testdata, "name"); err != ErrSMDNoEnabled {
		t.Errorf("Test 4: Did not get expected error ErrSMDNoEnabled")
	}
	testdata.UpdateType = SwStatusUpdate.String()
	if err := s.doCompUpdate(testdata, "name"); err != ErrSMDNoSwStatus {
		t.Errorf("Test 5: Did not get expected error ErrSMDNoSwStatus")
	}
	testdata.UpdateType = RoleUpdate.String()
	if err := s.doCompUpdate(testdata, "name"); err != ErrSMDNoRole {
		t.Errorf("Test 6: Did not get expected error ErrSMDNoRole")
	}
	testdata.UpdateType = SingleNIDUpdate.String()
	if err := s.doCompUpdate(testdata, "name"); err != ErrSMDNoNID {
		t.Errorf("Test 7: Did not get expected error ErrSMDNoNID")
	}
	var nid int64 = 123
	testdataBad3 := &CompUpdate{
		ComponentIDs: []string{"x0c0s0b0n0", "x0c0s0b0n1"},
		NID:          &nid,
		UpdateType:   SingleNIDUpdate.String(),
	}
	if err := s.doCompUpdate(testdataBad3, "name"); err != ErrSMDTooManyIDs {
		t.Errorf("Test 8: Did not get expected error ErrSMDTooManyIDs")
	}
}
