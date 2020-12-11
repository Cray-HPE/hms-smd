// Copyright (c) 2019 Cray Inc. All Rights Reserved.
package sm_test

import (
	stest "stash.us.cray.com/HMS/hms-smd/pkg/sharedtest"
	"stash.us.cray.com/HMS/hms-smd/pkg/sm"
	"testing"
)

///////////////////////////////////////////////////////////////////////////////
// Tests.  NOTE: INTENDED TO BE RUN IN THIS ORDER.  Some tests modify global
// data.
///////////////////////////////////////////////////////////////////////////////

func TestNewSystemHWInventory(t *testing.T) {
	hwinv, err := sm.NewSystemHWInventory(
		stest.HWInvByLocArray1,
		"s0",
		sm.HWInvFormatFullyFlat)
	if err != nil {
		t.Errorf("Test 1 Failed: Got error '%s'", err)
	} else if hwinv == nil {
		t.Errorf("Test 1 Failed: Got nil hwinv")
	}
	t.Log("Test 1 PASS")
	hwinv, err = sm.NewSystemHWInventory(
		stest.HWInvByLocArray1,
		"s0",
		sm.HWInvFormatNestNodesOnly)
	if err != nil {
		t.Errorf("Test 2 Failed: Got error '%s'", err)
	} else if hwinv == nil {
		t.Errorf("Test 2 Failed: Got nil hwinv")
	}
	t.Log("Test 2 PASS")
	hwinv, err = sm.NewSystemHWInventory(
		stest.HWInvByLocArray1,
		"s0",
		sm.HWInvFormatHierarchical)
	if err == nil {
		t.Errorf("Test 3 Failed: Got nil err")
	} else if err != sm.ErrHWInvFmtNI {
		t.Errorf("Test 3 Failed: Got wrong error '%s'", err)
	}
	t.Log("Test 3 PASS")
	hwinv, err = sm.NewSystemHWInventory(
		stest.HWInvByLocArray1,
		"s0",
		"SomeInvalidFormatType11111")
	if err == nil || hwinv != nil {
		t.Errorf("Test 4 Failed: Got nil err or hwinv != nil")
	} else if err != sm.ErrHWInvFmtInvalid {
		t.Errorf("Test 4 Failed: Got wrong error '%s'", err)
	}
	t.Log("Test 4 PASS")
}

func TestEncodeLocationInfo(t *testing.T) {
	for i, hwloc := range stest.HWInvByLocArray1 {
		bytes, err := hwloc.EncodeLocationInfo()
		if err != nil {
			t.Errorf("Test %d Failed: Got error '%s'", i, err)
		} else if len(bytes) == 0 {
			t.Errorf("Test %d Failed: Zero bytes returned", i)
		} else {
			t.Logf("Test %d PASS", i)
		}
	}
}

func TestDecodeLocationInfo(t *testing.T) {
	for i, hwloc := range stest.HWInvByLocArray1 {
		bytes, err := hwloc.EncodeLocationInfo()
		if err != nil {
			t.Errorf("Test %d Failed (pre-test): Got error '%s'", i, err)
			continue
		}
		err = hwloc.DecodeLocationInfo(bytes)
		if err != nil {
			t.Errorf("Test %d Failed: Got error '%s'", i, err)
			continue
		}
		t.Logf("Test %d PASS", i)
	}
}

func TestEncodeFRUInfo(t *testing.T) {
	for i, hwfru := range stest.HWInvByFRUArray1 {
		bytes, err := hwfru.EncodeFRUInfo()
		if err != nil {
			t.Errorf("Test %d Failed: Got error '%s'", i, err)
		} else if len(bytes) == 0 {
			t.Errorf("Test %d Failed: Zero bytes returned", i)
		} else {
			t.Logf("Test %d PASS", i)
		}
	}
}

func TestDecodeFRUInfo(t *testing.T) {
	for i, hwfru := range stest.HWInvByFRUArray1 {
		bytes, err := hwfru.EncodeFRUInfo()
		if err != nil {
			t.Errorf("Test %d Failed (pre-test): Got error '%s'", i, err)
			continue
		}
		err = hwfru.DecodeFRUInfo(bytes)
		if err != nil {
			t.Errorf("Test %d Failed: Got error '%s'", i, err)
			continue
		}
		t.Logf("Test %d PASS", i)
	}
}
