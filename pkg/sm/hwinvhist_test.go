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

package sm

import (
	"reflect"
	"testing"
)

//
// Test HWInventory History functions
//
func TestNewHWInvHistResp(t *testing.T) {
	testHWInvHist1 := HWInvHist{
		ID: "x5c4s3b2n1p0",
		FruId: "MFR-PARTNUMBER-SERIALNUMBER_1",
		Timestamp: "2020-01-21 11:36:00",
		EventType: "Added",
	}
	testHWInvHist2 := HWInvHist{
		ID: "x5c4s3b2n1p1",
		FruId: "MFR-PARTNUMBER-SERIALNUMBER_2",
		Timestamp: "2020-01-21 11:36:00",
		EventType: "Added",
	}
	testHWInvHist3 := HWInvHist{
		ID: "x5c4s3b2n1p0",
		FruId: "MFR-PARTNUMBER-SERIALNUMBER_1",
		Timestamp: "2020-01-21 12:00:00",
		EventType: "Removed",
	}
	testHWInvHist4 := HWInvHist{
		ID: "x5c4s3b2n1p1",
		FruId: "MFR-PARTNUMBER-SERIALNUMBER_2",
		Timestamp: "2020-01-21 12:00:00",
		EventType: "Removed",
	}
	testHWInvHist5 := HWInvHist{
		ID: "x5c4s3b2n1p0",
		FruId: "MFR-PARTNUMBER-SERIALNUMBER_2",
		Timestamp: "2020-01-21 12:10:00",
		EventType: "Added",
	}
	testHWInvHist6 := HWInvHist{
		ID: "x5c4s3b2n1p1",
		FruId: "MFR-PARTNUMBER-SERIALNUMBER_1",
		Timestamp: "2020-01-21 12:10:00",
		EventType: "Added",
	}
	testHwInvHistArray := []*HWInvHist{
		&testHWInvHist1,
		&testHWInvHist2,
		&testHWInvHist3,
		&testHWInvHist4,
		&testHWInvHist5,
		&testHWInvHist6,
	}

	tests := []struct {
		hwHistIn []*HWInvHist
		formatIn HWInvHistFmt
		expectedOut *HWInvHistResp
		expectedErr bool
	}{{
		hwHistIn:    testHwInvHistArray,
		formatIn:    HWInvHistFmtByLoc,
		expectedOut: &HWInvHistResp{
			Components: []HWInvHistArray{{
				ID: testHWInvHist1.ID,
				History: []*HWInvHist{
					&testHWInvHist1,
					&testHWInvHist3,
					&testHWInvHist5,
				},
			}, {
				ID: testHWInvHist2.ID,
				History: []*HWInvHist{
					&testHWInvHist2,
					&testHWInvHist4,
					&testHWInvHist6,
				},
			}},
		},
		expectedErr: false,
	}, {
		hwHistIn:    testHwInvHistArray,
		formatIn:    HWInvHistFmtByFRU,
		expectedOut: &HWInvHistResp{
			Components: []HWInvHistArray{{
				ID: testHWInvHist1.FruId,
				History: []*HWInvHist{
					&testHWInvHist1,
					&testHWInvHist3,
					&testHWInvHist6,
				},
			}, {
				ID: testHWInvHist2.FruId,
				History: []*HWInvHist{
					&testHWInvHist2,
					&testHWInvHist4,
					&testHWInvHist5,
				},
			}},
		},
		expectedErr: false,
	}}
	for i, test := range tests {
		out, err := NewHWInvHistResp(test.hwHistIn, test.formatIn)
		if err != nil && !test.expectedErr {
			t.Errorf("Test %v Failed: Unexpected error - %v", i, err)
		} else if err == nil && test.expectedErr {
			t.Errorf("Test %v Failed: Expected an error", i)
		} else if !reflect.DeepEqual(test.expectedOut, out) {
			t.Errorf("Test %v Failed: Expected HWInvHistResp struct '%v'; Received HWInvHistResp struct '%v'", i, test.expectedOut, out)
		}
	}
}

func TestVerifyNormalizeHWInvHistEventType(t *testing.T) {
	tests := []struct {
		eventTypeIn string
		expectedOut string
	}{{
		eventTypeIn: "added",
		expectedOut: "Added",
	}, {
		eventTypeIn: "Removed",
		expectedOut: "Removed",
	}, {
		eventTypeIn: "foo",
		expectedOut: "",
	}}
	for i, test := range tests {
		out := VerifyNormalizeHWInvHistEventType(test.eventTypeIn)
		if test.expectedOut != out {
			t.Errorf("Test %v Failed: Expected eventType '%v'; Received eventType '%v'", i, test.expectedOut, out)
		}
	}
}
