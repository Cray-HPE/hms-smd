// MIT License
//
// (C) Copyright [2020-2021] Hewlett Packard Enterprise Development LP
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

package hbtdapi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const DefaultHbtdUrl string = "http://cray-hbtd/hmi/v1"

type HBTD struct {
	Url    *url.URL
	Client *http.Client
}

type HBState struct {
	XName        string `json:"XName"`
	Heartbeating bool   `json:"Heartbeating"`
}

type HBStatusRsp struct {
	HBStates []HBState `json:"HBStates"`
}

type HBStatusPayload struct {
	XNames []string `json:"XNames"`
}

// Allocate and initialize new HBTD struct.
func NewHBTD(hbtdUrl string, httpClient *http.Client) *HBTD {
	var err error
	hbtd := new(HBTD)
	if hbtd.Url, err = url.Parse(hbtdUrl); err != nil {
		hbtd.Url, _ = url.Parse(DefaultHbtdUrl)
	} else {
		// Default to using http if not specified
		if len(hbtd.Url.Scheme) == 0 {
			hbtd.Url.Scheme = "http"
		}
	}

	// Create an httpClient if one was not given
	if httpClient == nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		hbtd.Client = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}
	} else {
		hbtd.Client = httpClient
	}
	return hbtd
}

// Query HBTD for node heartbeat status.
func (hbtd *HBTD) GetHeartbeatStatus(ids []string) ([]HBState, error) {
	var rsp HBStatusRsp
	// Validate inputs
	if hbtd.Url == nil {
		return nil, fmt.Errorf("HBTD struct has no URL")
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("No component IDs to query")
	}

	// Construct a GET to /hbstates for HBTD to get the node heartbeat status
	uri := hbtd.Url.String() + "/hbstates"
	hbStatReq := HBStatusPayload{XNames: ids}
	payload, err := json.Marshal(hbStatReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	body, err := hbtd.doRequest(req)
	if err != nil {
		return nil, err
	}

	// HBTD returns null if the component was not found
	if body == nil {
		return nil, nil
	}

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return nil, err
	}

	return rsp.HBStates, nil
}

// doRequest sends a HTTP request to HBTD
func (hbtd *HBTD) doRequest(req *http.Request) ([]byte, error) {
	// Error if there is no client defined
	if hbtd.Client == nil {
		return nil, fmt.Errorf("HBTD struct has no HTTP Client")
	}

	// Send the request
	rsp, err := hbtd.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Read the response
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}
