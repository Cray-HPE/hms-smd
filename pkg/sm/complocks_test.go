// MIT License
//
// (C) Copyright [2019-2021] Hewlett Packard Enterprise Development LP
//
// Permission is hereby granted, free of charge, to any person obtaining a
// copy of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included
// in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
// THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
// OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

// This file contains unit tests for the CompLock struct

package sm

import (
	"reflect"
	"testing"

	base "github.com/Cray-HPE/hms-base"
)

//
// Test component locks
//

func TestNewCompLock(t *testing.T) {
	tests := []struct {
		reasonIn    string
		ownerIn     string
		lifetimeIn  int
		xnamesIn    []string
		expectedOut *CompLock
		expectedErr error
	}{{
		reasonIn:   "Because I want to",
		ownerIn:    "my_service",
		lifetimeIn: 100,
		xnamesIn:   []string{"x0c0s1b0n0", "x0c0s2b0n0"},
		expectedOut: &CompLock{
			Reason:     "Because I want to",
			Owner:      "my_service",
			Lifetime:   120,
			Xnames:     []string{"x0c0s1b0n0", "x0c0s2b0n0"},
			verified:   true,
			normalized: true,
		},
		expectedErr: nil,
	}, {
		reasonIn:   "Because I want to",
		ownerIn:    "my_service",
		lifetimeIn: 100,
		xnamesIn:   []string{"x0c0s01b0n0", "x00c0s02b0n0"},
		expectedOut: &CompLock{
			Reason:     "Because I want to",
			Owner:      "my_service",
			Lifetime:   120,
			Xnames:     []string{"x0c0s1b0n0", "x0c0s2b0n0"},
			verified:   true,
			normalized: true,
		},
		expectedErr: nil,
	}, {
		reasonIn:    "Because I want to",
		ownerIn:     "my_service",
		lifetimeIn:  0,
		xnamesIn:    []string{"x0c0s1b0n0", "x0c0s2b0n0"},
		expectedOut: &CompLock{},
		expectedErr: ErrCompLockBadLifetime,
	}}
	for i, test := range tests {
		out, err := NewCompLock(test.reasonIn, test.ownerIn, test.lifetimeIn, test.xnamesIn)
		if err != test.expectedErr {
			t.Errorf("Test %v Failed: Expected error '%v'; Received error '%v'", i, test.expectedErr, err)
		} else if test.expectedErr == nil {
			if !reflect.DeepEqual(test.expectedOut, out) {
				t.Errorf("Test %v Failed: Expected partition struct '%v'; Received partition struct '%v'", i, test.expectedOut, out)
			}
		}
	}
}

func TestVerifyCompLock(t *testing.T) {
	tests := []struct {
		in          *CompLock
		expectedOut error
	}{{
		in: &CompLock{
			Reason:   "Because I want to",
			Owner:    "my_service",
			Lifetime: 100,
			Xnames:   []string{"x0c0s1b0n0", "x0c0s2b0n0"},
		},
		expectedOut: nil,
	}, {
		in: &CompLock{
			Reason:   "Because I want to",
			Owner:    "my_service",
			Lifetime: 100,
			Xnames:   []string{"x0c0s1b0n0", "x0c0s2b0n0"},
			verified: true,
		},
		expectedOut: nil,
	}, {
		in: &CompLock{
			Reason:   "Because I want to",
			Owner:    "my_service",
			Lifetime: 0,
			Xnames:   []string{"x0c0s1b0n0", "x0c0s2b0n0"},
		},
		expectedOut: ErrCompLockBadLifetime,
	}, {
		in: &CompLock{
			Reason:   "Because I want to",
			Owner:    "my_service",
			Lifetime: 100,
			Xnames:   []string{"foo", "bar"},
		},
		expectedOut: base.ErrHMSTypeInvalid,
	}}
	for i, test := range tests {
		out := test.in.Verify()
		if test.expectedOut != out {
			t.Errorf("Test %v Failed: Expected error '%v'; Received error '%v'", i, test.expectedOut, out)
		}
	}
}

func TestNormalizeCompLock(t *testing.T) {
	tests := []struct {
		in          *CompLock
		expectedOut *CompLock
	}{{
		in: &CompLock{
			Reason:   "Because I want to",
			Owner:    "my_service",
			Lifetime: 100,
			Xnames:   []string{"x0C0s1B0n0", "x0c0s2b0n0"},
		},
		expectedOut: &CompLock{
			Reason:     "Because I want to",
			Owner:      "my_service",
			Lifetime:   120,
			Xnames:     []string{"x0c0s1b0n0", "x0c0s2b0n0"},
			normalized: true,
		},
	}, {
		in: &CompLock{
			Reason:     "Because I want to",
			Owner:      "my_service",
			Lifetime:   100,
			Xnames:     []string{"x0C0s1B0n0", "x0c0s2b0n0"},
			normalized: true,
		},
		expectedOut: &CompLock{
			Reason:     "Because I want to",
			Owner:      "my_service",
			Lifetime:   100,
			Xnames:     []string{"x0C0s1B0n0", "x0c0s2b0n0"},
			normalized: true,
		},
	}}
	for i, test := range tests {
		test.in.Normalize()
		if !reflect.DeepEqual(test.expectedOut, test.in) {
			t.Errorf("Test %v Failed: Expected CompLock struct '%v'; Received CompLock struct '%v'", i, test.expectedOut, test.in)
		}
	}
}

func TestVerifyCompLockPatch(t *testing.T) {
	plt1 := 100
	plt2 := 0

	tests := []struct {
		in          *CompLockPatch
		expectedOut error
	}{{
		in: &CompLockPatch{
			Lifetime: &plt1,
		},
		expectedOut: nil,
	}, {
		in: &CompLockPatch{
			Lifetime: &plt2,
		},
		expectedOut: ErrCompLockBadLifetime,
	}}
	for i, test := range tests {
		out := test.in.Verify()
		if test.expectedOut != out {
			t.Errorf("Test %v Failed: Expected error '%v'; Received error '%v'", i, test.expectedOut, out)
		}
	}
}

// Locking v2
func TestVerifyNormalizeCompLockV2Filter(t *testing.T) {
	tests := []struct {
		in  *CompLockV2Filter
		out *CompLockV2Filter
		err error
	}{{
		in: &CompLockV2Filter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 1,
		},
		out: &CompLockV2Filter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 1,
		},
		err: nil,
	}, {
		in: &CompLockV2Filter{
			ProcessingModel:     "Rigid",
			ReservationDuration: 1,
		},
		out: &CompLockV2Filter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 1,
		},
		err: nil,
	}, {
		in: &CompLockV2Filter{
			ProcessingModel: CLProcessingModelRigid,
		},
		out: &CompLockV2Filter{
			ProcessingModel: CLProcessingModelRigid,
		},
		err: nil,
	}, {
		in: &CompLockV2Filter{
			ProcessingModel: "foo",
		},
		out: &CompLockV2Filter{
			ProcessingModel: "foo",
		},
		err: ErrCompLockV2BadProcessingModel,
	}, {
		in: &CompLockV2Filter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 16,
		},
		out: &CompLockV2Filter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 16,
		},
		err: ErrCompLockV2BadDuration,
	}}
	for i, test := range tests {
		err := test.in.VerifyNormalize()
		if test.err != err {
			t.Errorf("Test %v Failed: Expected error '%v'; Received error '%v'", i, test.err, err)
		} else if err == nil {
			if !reflect.DeepEqual(test.out, test.in) {
				t.Errorf("Test %v Failed: Expected CompLockV2Filter struct '%v'; Received CompLockV2Filter struct '%v'", i, test.out, test.in)
			}
		}
	}
}

func TestVerifyNormalizeCompLockV2Key(t *testing.T) {
	tests := []struct {
		in  *CompLockV2Key
		out *CompLockV2Key
		err error
	}{{
		in: &CompLockV2Key{
			ID:  "x0c0s0b0n1",
			Key: "x0c0s0b0n1:rk:Some-UUID",
		},
		out: &CompLockV2Key{
			ID:  "x0c0s0b0n1",
			Key: "x0c0s0b0n1:rk:some-uuid",
		},
		err: nil,
	}, {
		in: &CompLockV2Key{
			Key: "x0c0s0b0n1:rk:Some-UUID",
		},
		out: &CompLockV2Key{
			Key: "x0c0s0b0n1:rk:some-uuid",
		},
		err: base.ErrHMSTypeInvalid,
	}}
	for i, test := range tests {
		err := test.in.VerifyNormalize()
		if test.err != err {
			t.Errorf("Test %v Failed: Expected error '%v'; Received error '%v'", i, test.err, err)
		} else if err == nil {
			if !reflect.DeepEqual(test.out, test.in) {
				t.Errorf("Test %v Failed: Expected CompLockV2Key struct '%v'; Received CompLockV2Key struct '%v'", i, test.out, test.in)
			}
		}
	}
}

func TestVerifyNormalizeCompLockV2ReservationFilter(t *testing.T) {
	tests := []struct {
		in  *CompLockV2ReservationFilter
		out *CompLockV2ReservationFilter
		err error
	}{{
		in: &CompLockV2ReservationFilter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 1,
			ReservationKeys: []CompLockV2Key{{
				ID:  "x0c0s0b0n1",
				Key: "x0c0s0b0n1:rk:Some-UUID",
			}},
		},
		out: &CompLockV2ReservationFilter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 1,
			ReservationKeys: []CompLockV2Key{{
				ID:  "x0c0s0b0n1",
				Key: "x0c0s0b0n1:rk:some-uuid",
			}},
		},
		err: nil,
	}, {
		in: &CompLockV2ReservationFilter{
			ProcessingModel:     "Rigid",
			ReservationDuration: 1,
			ReservationKeys: []CompLockV2Key{{
				ID:  "x0c0s0b0n1",
				Key: "x0c0s0b0n1:rk:Some-UUID",
			}},
		},
		out: &CompLockV2ReservationFilter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 1,
			ReservationKeys: []CompLockV2Key{{
				ID:  "x0c0s0b0n1",
				Key: "x0c0s0b0n1:rk:some-uuid",
			}},
		},
		err: nil,
	}, {
		in: &CompLockV2ReservationFilter{
			ProcessingModel: "foo",
		},
		out: &CompLockV2ReservationFilter{
			ProcessingModel: "foo",
		},
		err: ErrCompLockV2BadProcessingModel,
	}, {
		in: &CompLockV2ReservationFilter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 16,
		},
		out: &CompLockV2ReservationFilter{
			ProcessingModel:     CLProcessingModelRigid,
			ReservationDuration: 16,
		},
		err: ErrCompLockV2BadDuration,
	}}
	for i, test := range tests {
		err := test.in.VerifyNormalize()
		if test.err != err {
			t.Errorf("Test %v Failed: Expected error '%v'; Received error '%v'", i, test.err, err)
		} else if err == nil {
			if !reflect.DeepEqual(test.out, test.in) {
				t.Errorf("Test %v Failed: Expected CompLockV2ReservationFilter struct '%v'; Received CompLockV2ReservationFilter struct '%v'", i, test.out, test.in)
			}
		}
	}
}
