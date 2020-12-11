// Copyright 2020 Hewlett Packard Enterprise Development LP

package sm

import (
	"reflect"
	"testing"

	base "stash.us.cray.com/HMS/hms-base"
)

//
// Test Component Ethernet Interfaces V1
//

func TestNewCompEthInterface(t *testing.T) {
	tests := []struct {
		desc        string
		macAddr     string
		ipAddr      string
		compID      string
		expectedOut *CompEthInterface
		expectedErr error
	}{{ // Test 0 - Normalize MAC address and comp ID
		desc:    "My description",
		macAddr: "A4:BF:01:38:EE:65",
		ipAddr:  "10.254.2.14",
		compID:  "x0003000c0s26b0",
		expectedOut: &CompEthInterface{
			ID:      "a4bf0138ee65",
			Desc:    "My description",
			MACAddr: "a4:bf:01:38:ee:65",
			IPAddr:  "10.254.2.14",
			CompID:  "x3000c0s26b0",
			Type:    base.NodeBMC.String(),
		},
		expectedErr: nil,
	}, { // Test 1 - Minimal info
		macAddr: "a4:bf:01:38:ee:65",
		expectedOut: &CompEthInterface{
			ID:      "a4bf0138ee65",
			MACAddr: "a4:bf:01:38:ee:65",
		},
		expectedErr: nil,
	}, { // Test 2 - Bad MAC address
		macAddr:     "",
		expectedOut: nil,
		expectedErr: ErrCompEthInterfaceBadMAC,
	}, { // Test 3 - Bad comp ID
		macAddr:     "a4:bf:01:38:ee:65",
		compID:      "y3000",
		expectedOut: nil,
		expectedErr: ErrCompEthInterfaceBadCompID,
	}}
	for i, test := range tests {
		out, err := NewCompEthInterface(test.desc, test.macAddr, test.ipAddr, test.compID)
		if err != test.expectedErr {
			t.Errorf("Test %v Failed: Expected error '%v'; Received error '%v'", i, test.expectedErr, err)
		} else if test.expectedErr == nil {
			if !reflect.DeepEqual(test.expectedOut, out) {
				t.Errorf("Test %v Failed: Expected CompEthInterface struct '%v'; Received CompEthInterface struct '%v'", i, test.expectedOut, out)
			}
		}
	}
}

//
// Test Component Ethernet Interfaces V2
//

func TestNewCompEthInterfaceV2(t *testing.T) {
	tests := []struct {
		desc        string
		macAddr     string
		ipAddrs     []IPAddressMapping
		compID      string
		expectedOut *CompEthInterfaceV2
		expectedErr error
	}{{ // Test 0 - Normalize MAC address and comp ID
		desc:    "My description",
		macAddr: "A4:BF:01:38:EE:65",
		ipAddrs: []IPAddressMapping{{IPAddr: "10.254.2.14"}},
		compID:  "x0003000c0s26b0",
		expectedOut: &CompEthInterfaceV2{
			ID:      "a4bf0138ee65",
			Desc:    "My description",
			MACAddr: "a4:bf:01:38:ee:65",
			IPAddrs: []IPAddressMapping{
				{IPAddr: "10.254.2.14"},
			},
			CompID: "x3000c0s26b0",
			Type:   base.NodeBMC.String(),
		},
		expectedErr: nil,
	}, { // Test 1 - Minimal info
		macAddr: "a4:bf:01:38:ee:65",
		expectedOut: &CompEthInterfaceV2{
			ID:      "a4bf0138ee65",
			MACAddr: "a4:bf:01:38:ee:65",
		},
		expectedErr: nil,
	}, { // Test 2 - Bad MAC address
		macAddr:     "",
		expectedOut: nil,
		expectedErr: ErrCompEthInterfaceBadMAC,
	}, { // Test 3 - Bad comp ID
		macAddr:     "a4:bf:01:38:ee:65",
		compID:      "y3000",
		expectedOut: nil,
		expectedErr: ErrCompEthInterfaceBadCompID,
	}, { // Test 4 - Bad (empty) IP Address
		desc:        "My description",
		macAddr:     "A4:BF:01:38:EE:65",
		ipAddrs:     []IPAddressMapping{{IPAddr: ""}},
		compID:      "x0003000c0s26b0",
		expectedOut: nil,
		expectedErr: ErrCompEthInterfaceBadIPAddress,
	}}
	for i, test := range tests {
		out, err := NewCompEthInterfaceV2(test.desc, test.macAddr, test.compID, test.ipAddrs)
		if err != test.expectedErr {
			t.Errorf("Test %v Failed: Expected error '%v'; Received error '%v'", i, test.expectedErr, err)
		} else if test.expectedErr == nil {
			if !reflect.DeepEqual(test.expectedOut, out) {
				t.Errorf("Test %v Failed: Expected CompEthInterface struct '%v'; Received CompEthInterface struct '%v'", i, test.expectedOut, out)
			}
		}
	}
}

func TestCompEthInterfaceV2ToV1(t *testing.T) {
	tests := []struct {
		desc        string
		macAddr     string
		ipAddrs     []IPAddressMapping
		compID      string
		expectedOut *CompEthInterface
	}{{ // Test 0 - Normalize MAC address and comp ID
		desc:    "My description",
		macAddr: "A4:BF:01:38:EE:65",
		ipAddrs: []IPAddressMapping{{IPAddr: "10.254.2.14"}},
		compID:  "x0003000c0s26b0",
		expectedOut: &CompEthInterface{
			ID:      "a4bf0138ee65",
			Desc:    "My description",
			MACAddr: "a4:bf:01:38:ee:65",
			IPAddr:  "10.254.2.14",
			CompID:  "x3000c0s26b0",
			Type:    base.NodeBMC.String(),
		},
	}, { // Test 1 - Minimal info
		macAddr: "a4:bf:01:38:ee:65",
		expectedOut: &CompEthInterface{
			ID:      "a4bf0138ee65",
			MACAddr: "a4:bf:01:38:ee:65",
		},
	}, { // Test 2 - Multiple IPs
		desc:    "My description",
		macAddr: "A4:BF:01:38:EE:65",
		ipAddrs: []IPAddressMapping{
			{IPAddr: "10.254.2.14"},
			{IPAddr: "10.252.2.14"},
		},
		compID: "x0003000c0s26b0",
		expectedOut: &CompEthInterface{
			ID:      "a4bf0138ee65",
			Desc:    "My description",
			MACAddr: "a4:bf:01:38:ee:65",
			IPAddr:  "10.254.2.14",
			CompID:  "x3000c0s26b0",
			Type:    base.NodeBMC.String(),
		},
	}}
	for i, test := range tests {
		ceiV2, err := NewCompEthInterfaceV2(test.desc, test.macAddr, test.compID, test.ipAddrs)
		if err != nil {
			t.Errorf("Test %v Failed: Failed to Create CompEthInterfaceV2; error '%v'", i, err)
			continue
		}

		out := ceiV2.ToV1()

		if !reflect.DeepEqual(test.expectedOut, out) {
			t.Errorf("Test %v Failed: Expected CompEthInterface struct '%v'; Received CompEthInterface struct '%v'", i, test.expectedOut, out)
		}
	}
}

func TestNewIPAddressMapping(t *testing.T) {
	tests := []struct {
		ipAddr      string
		network     string
		expectedOut *IPAddressMapping
		expectedErr error
	}{{ // Test 0 - IP Address Only
		ipAddr: "10.254.2.14",
		expectedOut: &IPAddressMapping{
			IPAddr: "10.254.2.14",
		},
		expectedErr: nil,
	}, { // Test 1 - IP Address and Network
		ipAddr:  "10.254.2.14",
		network: "HMN",
		expectedOut: &IPAddressMapping{
			IPAddr:  "10.254.2.14",
			Network: "HMN",
		},
		expectedErr: nil,
	}, { // Test 2 - No IP Address
		expectedOut: nil,
		expectedErr: ErrCompEthInterfaceBadIPAddress,
	}, { // Test 3 - No IP Address, only network
		network:     "HMN",
		expectedOut: nil,
		expectedErr: ErrCompEthInterfaceBadIPAddress,
	}}
	for i, test := range tests {
		out, err := NewIPAddressMapping(test.ipAddr, test.network)
		if err != test.expectedErr {
			t.Errorf("Test %v Failed: Expected error '%v'; Received error '%v'", i, test.expectedErr, err)
		} else if test.expectedErr == nil {
			if !reflect.DeepEqual(test.expectedOut, out) {
				t.Errorf("Test %v Failed: Expected CompEthInterface struct '%v'; Received CompEthInterface struct '%v'", i, test.expectedOut, out)
			}
		}
	}
}
