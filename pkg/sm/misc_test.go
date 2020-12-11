// Copyright (c) 2019 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.
//
// This file contains unit tests for the the stucts in misc.go

package sm

import (
	"reflect"
	"testing"
)

//
// Test State Redfish Poll Jobs
//

func TestNewStateRFPollJob(t *testing.T) {
	tests := []struct {
		xnameIn     string
		delayIn     int
		pollIn      int
		lifetimeIn  int
		keepAliveIn int
		expectedOut *Job
		expectedErr bool
	}{{
		xnameIn:     "x0c0s0b0n0",
		delayIn:     30,
		pollIn:      10,
		lifetimeIn:  30,
		keepAliveIn: 20,
		expectedOut: &Job{
			JobData: JobData{
				Type:       JobTypeSRFP,
				Status:     JobNotStarted,
				Lifetime:   30,
				KeepAlive:  20,
			},
			Data: &SrfpJobData{
				CompId: "x0c0s0b0n0",
				Delay:  30,
				Poll:   10,
			},
		},
		expectedErr: false,
	}, {
		xnameIn:     "x0c0s0b0n0",
		delayIn:     0,
		pollIn:      0,
		lifetimeIn:  0,
		keepAliveIn: 0,
		expectedOut: &Job{
			JobData: JobData{
				Type:       JobTypeSRFP,
				Status:     JobNotStarted,
				Lifetime:   10,
				KeepAlive:  5,
			},
			Data: &SrfpJobData{
				CompId: "x0c0s0b0n0",
				Delay:  1,
				Poll:   5,
			},
		},
		expectedErr: false,
	}, {
		xnameIn:     "x0c0s0b0n0",
		delayIn:     30,
		pollIn:      10,
		lifetimeIn:  10,
		keepAliveIn: 10,
		expectedOut: &Job{
			JobData: JobData{
				Type:       JobTypeSRFP,
				Status:     JobNotStarted,
				Lifetime:   15,
				KeepAlive:  10,
			},
			Data: &SrfpJobData{
				CompId: "x0c0s0b0n0",
				Delay:  30,
				Poll:   10,
			},
		},
		expectedErr: false,
	}, {
		xnameIn:     "",
		delayIn:     30,
		pollIn:      10,
		lifetimeIn:  30,
		keepAliveIn: 20,
		expectedOut: nil,
		expectedErr: true,
	}}
	for i, test := range tests {
		out, err := NewStateRFPollJob(test.xnameIn, test.delayIn, test.pollIn, test.lifetimeIn, test.keepAliveIn)
		if err != nil && !test.expectedErr {
			t.Errorf("Test %v Failed: Unexpected error - %v", i, err)
		} else if err == nil && test.expectedErr {
			t.Errorf("Test %v Failed: Expected an error", i)
		} else if !reflect.DeepEqual(test.expectedOut, out) {
			t.Errorf("Test %v Failed: Expected Job struct '%v'; Received Job struct '%v'", i, test.expectedOut, out)
		}
	}
}
