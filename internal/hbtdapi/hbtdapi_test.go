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

package hbtdapi

///////////////////////////////////////////////////////////////////////////////
// Pre-Test Setup
///////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

var client *http.Client

// RoundTrip method override
type RTFunc func(req *http.Request) *http.Response

// Implement RoundTrip interface by implementing RoundTrip method
func (f RTFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(f RTFunc) *http.Client {
	return &http.Client{
		Transport: RTFunc(f),
	}
}

// Sets up the http client for testing
func TestMain(m *testing.M) {
	client = NewTestClient(NewRTFuncSLSAPI())

	excode := 1
	excode = m.Run()
	os.Exit(excode)
}

///////////////////////////////////////////////////////////////////////////////
// Unit Tests
///////////////////////////////////////////////////////////////////////////////

func TestNewHBTD(t *testing.T) {
	tests := []struct {
		hbtdUrlIn    string
		expectedUrl string
	}{{
		hbtdUrlIn:   "http://cray-hbtd",
		expectedUrl: "http://cray-hbtd",
	}, {
		hbtdUrlIn:   "cray-hbtd",
		expectedUrl: "http://cray-hbtd",
	}}
	for i, test := range tests {
		out := NewHBTD(test.hbtdUrlIn, nil)
		if test.expectedUrl != out.Url.String() {
			t.Errorf("Test %v Failed: Expected HBTD URL '%v'; Received HBTD URL '%v'", i, test.expectedUrl, out.Url.String())
		}
	}
}

func TestGetHeartbeatStatus(t *testing.T) {
	defaultUrl, _ := url.Parse("http://cray-hbtd/v1")
	goodUrl, _ := url.Parse("http://cray-hbtd.good/v1")
	badUrl, _ := url.Parse("http://cray-hbtd.bad/v1")
	tests := []struct {
		HBTDUrl      *url.URL
		Client       *http.Client
		ids          []string
		expectedInfo []HBState
		expectErr    bool
	}{{
		HBTDUrl:      defaultUrl,
		Client:       client,
		ids:          []string{"x0c0s0b0n0"},
		expectedInfo: []HBState{
			HBState{
				XName: "x0c0s0b0n0",
				Heartbeating: true,
			},
		},
		expectErr:    false,
	}, {
		HBTDUrl:      goodUrl,
		Client:       client,
		ids:          []string{"x0c0s0b0n0","x0c0s1b0n0"},
		expectedInfo: []HBState{
			HBState{
				XName: "x0c0s0b0n0",
				Heartbeating: true,
			},
			HBState{
				XName: "x0c0s1b0n0",
				Heartbeating: false,
			},
		},
		expectErr:    false,
	}, {
		HBTDUrl:      badUrl,
		Client:       client,
		ids:          []string{"x0c0s2b0n0"},
		expectedInfo: nil,
		expectErr:    true,
	}, {
		HBTDUrl:      nil,
		Client:       client,
		ids:          []string{"x0c0s4b0n0"},
		expectedInfo: nil,
		expectErr:    true,
	}, {
		HBTDUrl:      defaultUrl,
		Client:       client,
		ids:          []string{},
		expectedInfo: nil,
		expectErr:    true,
	}, {
		HBTDUrl:      defaultUrl,
		Client:       nil,
		ids:          []string{"x0c0s4b0n0"},
		expectedInfo: nil,
		expectErr:    true,
	}}

	for i, test := range tests {
		hbtd := HBTD{
			Url: test.HBTDUrl,
			Client: test.Client,
		}
		out, err := hbtd.GetHeartbeatStatus(test.ids)
		if err != nil {
			if !test.expectErr {
				t.Errorf("Test %v Failed: Unexpected error - %v", i, err)
			}
		} else {
			if test.expectErr {
				t.Errorf("Test %v Failed: Expected an error", i)
			} else {
				if test.expectedInfo != nil && len(out) == len(test.expectedInfo) {
					for i, stat := range out {
						if stat.XName != test.expectedInfo[i].XName &&
							stat.Heartbeating != test.expectedInfo[i].Heartbeating {
							t.Errorf("Test %v Failed: Expected node info '%v'; Received node info '%v'", i, test.expectedInfo, out)
							break
						}
					}
				}
			}
		}
	}
}

///////////////////////////////////////////////////////////////////////////////
// Mock HBTD Data
///////////////////////////////////////////////////////////////////////////////

const testPayloadHBTDAPI_good = `
{
	"HBStates": [{
		"XName": "x0c0s0b0n0",
		"Heartbeating": true
	}]
}`

const testPayloadHBTDAPI_goodNoHB = `
{
	"HBStates": [{
		"XName": "x0c0s0b0n0",
		"Heartbeating": true
	},{
		"XName": "x0c0s1b0n0",
		"Heartbeating": false
	}]
}`

func NewRTFuncSLSAPI() RTFunc {
	return func(req *http.Request) *http.Response {
		// Test request parameters
		switch req.URL.String() {
		case "http://cray-hbtd/v1/hbstates":
			return &http.Response{
				StatusCode: 200,
				// Send mock response for rpath
				Body:   ioutil.NopCloser(bytes.NewBufferString(testPayloadHBTDAPI_good)),
				Header: make(http.Header),
			}
		case "http://cray-hbtd.bad/v1/hbstates":
			return &http.Response{
				StatusCode: 500,
				// Send mock response for rpath
				Body:   ioutil.NopCloser(bytes.NewBufferString("")),
				Header: make(http.Header),
			}
		case "http://cray-hbtd.good/v1/hbstates":
			return &http.Response{
				StatusCode: 200,
				// Send mock response for rpath
				Body:   ioutil.NopCloser(bytes.NewBufferString(testPayloadHBTDAPI_goodNoHB)),
				Header: make(http.Header),
			}
		default:
			return &http.Response{
				StatusCode: 404,
				// Send mock response for rpath
				Body:   ioutil.NopCloser(bytes.NewBufferString("")),
				Header: make(http.Header),
			}
		}
	}
}