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
// This file contains structs and function for interfacing with SLS

package slsapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	base "stash.us.cray.com/HMS/hms-base"
	"time"
)

const DefaultSlsUrl string = "http://cray-sls/"

type SLS struct {
	Url    *url.URL
	Client *http.Client
}

type slsReady struct {
	Ready  bool   `json:"Ready"`
	Reason string `json:"Reason,omitempty"`
	Code   int    `json:"Code,omitempty"`
}

type NodeHardware struct {
	Parent          string       `json:"Parent"`
	Children        []string     `json:"Children"`
	Xname           string       `json:"Xname"`
	Type            string       `json:"Type"`
	Class           string       `json:"Class"`
	TypeString      base.HMSType `json:"TypeString"`
	ExtraProperties ComptypeNode `json:"ExtraProperties"`
}

type ComptypeNode struct {
	NID     int    `json:"NID,omitempty"`
	Role    string `json:"Role"`
	SubRole string `json:"SubRole"`
}

// Allocate and initialize new SLS struct.
func NewSLS(slsUrl string, httpClient *http.Client) *SLS {
	var err error
	sls := new(SLS)
	if sls.Url, err = url.Parse(slsUrl); err != nil {
		sls.Url, _ = url.Parse(DefaultSlsUrl)
	} else {
		// Default to using http if not specified
		if len(sls.Url.Scheme) == 0 {
			sls.Url.Scheme = "http"
		}
	}

	// Create an httpClient if one was not given
	if httpClient == nil {
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		sls.Client = &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		}
	} else {
		sls.Client = httpClient
	}
	return sls
}

// Ping SLS to see if it is ready.
func (sls *SLS) IsReady() (bool, error) {
	var ready slsReady
	if sls.Url == nil {
		return false, fmt.Errorf("SLS struct has no URL")
	}
	uri := sls.Url.String() + "/ready"
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return false, err
	}
	body, err := sls.doRequest(req)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(body, &ready)
	if err != nil {
		return false, err
	}

	if !ready.Ready {
		err = fmt.Errorf("%d - %s", ready.Code, ready.Reason)
	}
	return ready.Ready, err
}

// Query SLS for node information. This just picks up the ExtraProperties
// struct for nodes from SLS.
func (sls *SLS) GetNodeInfo(id string) (ComptypeNode, error) {
	var nh NodeHardware
	// Validate inputs
	if sls.Url == nil {
		return ComptypeNode{}, fmt.Errorf("SLS struct has no URL")
	}
	if len(id) == 0 {
		return ComptypeNode{}, fmt.Errorf("Id is missing")
	}

	// Construct a GET to /hardware/<xname> for SLS to get the node info
	uri := sls.Url.String() + "/hardware/" + id
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return ComptypeNode{}, err
	}
	body, err := sls.doRequest(req)
	if err != nil {
		return ComptypeNode{}, err
	}

	// SLS returns null if the component was not found
	if body == nil {
		return ComptypeNode{}, nil
	}

	err = json.Unmarshal(body, &nh)
	if err != nil {
		return ComptypeNode{}, err
	}

	return nh.ExtraProperties, nil
}

// doRequest sends a HTTP request to SLS
func (sls *SLS) doRequest(req *http.Request) ([]byte, error) {
	// Error if there is no client defined
	if sls.Client == nil {
		return nil, fmt.Errorf("SLS struct has no HTTP Client")
	}

	// Send the request
	rsp, err := sls.Client.Do(req)
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
