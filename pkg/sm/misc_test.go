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
