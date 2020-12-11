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
// This file contains unit tests for the components POST and PUT structs

package sm

import (
	"fmt"
	"reflect"
	base "stash.us.cray.com/HMS/hms-base"
	"testing"
)

//
// Test component POST and PUT struct constructors
//

func TestNewCompPost(t *testing.T) {
	tests := []struct {
		compsIn     []base.Component
		forceIn     bool
		expectedOut *ComponentsPost
		expectedErr error
	}{{
		compsIn: []base.Component{{
			ID:    "x0c0s0b0n0",
			State: "empty",
		}},
		forceIn: false,
		expectedOut: &ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "x0c0s0b0n0",
						Type:  base.Node.String(),
						State: base.StateEmpty.String(),
						Flag:  base.FlagOK.String(),
					},
				},
			},
			Force: false,
		},
		expectedErr: nil,
	}, {
		compsIn: []base.Component{{
			ID:    "foo",
			State: "empty",
		}},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("xname ID 'foo' is invalid"),
	}, {
		compsIn: []base.Component{{
			ID:    "x0c0s0b0n0",
			State: "Exploded",
		}},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("state 'Exploded' is invalid"),
	}, {
		compsIn: []base.Component{{
			ID:    "x0c0s0b0n0",
			State: "empty",
			Flag:  "None",
		}},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("flag 'None' is invalid"),
	}, {
		compsIn: []base.Component{{
			ID:    "x0c0s0b0n0",
			State: "empty",
			Flag:  "OK",
			Role:  "Master",
		}},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("role 'Master' is invalid"),
	}, {
		compsIn: []base.Component{{
			ID:      "x0c0s0b0n0",
			State:   "empty",
			Flag:    "OK",
			Role:    "Compute",
			NetType: "Fishing",
		}},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("netType 'Fishing' is invalid"),
	}, {
		compsIn: []base.Component{{
			ID:      "x0c0s0b0n0",
			State:   "empty",
			Flag:    "OK",
			Role:    "Compute",
			NetType: "Sling",
			Arch:    "Bridge",
		}},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("arch 'Bridge' is invalid"),
	}, {
		compsIn: []base.Component{
			{
				ID:      "x00c0s0b0n0",
				State:   "empty",
				Flag:    "ok",
				Role:    "compute",
				NetType: "sling",
				Arch:    "x86",
			}, {
				ID:    "x0c0s0b0n1",
				State: "empty",
			},
		},
		forceIn: false,
		expectedOut: &ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:      "x0c0s0b0n0",
						Type:    base.Node.String(),
						State:   base.StateEmpty.String(),
						Flag:    base.FlagOK.String(),
						Role:    base.RoleCompute.String(),
						NetType: base.NetSling.String(),
						Arch:    base.ArchX86.String(),
					},
					&base.Component{
						ID:    "x0c0s0b0n1",
						Type:  base.Node.String(),
						State: base.StateEmpty.String(),
						Flag:  base.FlagOK.String(),
					},
				},
			},
			Force: false,
		},
		expectedErr: nil,
	}}
	for i, test := range tests {
		out, err := NewCompPost(test.compsIn, test.forceIn)
		if err != nil && test.expectedErr == nil {
			t.Errorf("Test %v Failed: Received unexpected error '%v'", i, err)
		} else if test.expectedErr == nil {
			if !reflect.DeepEqual(test.expectedOut, out) {
				t.Errorf("Test %v Failed: Expected compPost struct '%v'; Received compPost struct '%v'", i, test.expectedOut, out)
			}
		}
	}
}

func TestVerifyNormalizeCompPost(t *testing.T) {
	tests := []struct {
		in          ComponentsPost
		expectedOut ComponentsPost
		expectedErr error
	}{{
		in: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "x0c0s0b0n0",
						State: "empty",
					},
				},
			},
			Force: false,
		},
		expectedOut: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "x0c0s0b0n0",
						Type:  base.Node.String(),
						State: base.StateEmpty.String(),
						Flag:  base.FlagOK.String(),
					},
				},
			},
			Force: false,
		},
		expectedErr: nil,
	}, {
		in: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "foo",
						State: "empty",
					},
				},
			},
			Force: false,
		},
		expectedOut: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "foo",
						State: "empty",
					},
				},
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("xname ID 'foo' is invalid"),
	}, {
		in: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "x0c0s0b0n0",
						State: "Exploded",
					},
				},
			},
			Force: false,
		},
		expectedOut: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "x0c0s0b0n0",
						State: "Exploded",
					},
				},
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("state 'Exploded' is invalid"),
	}, {
		in: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "x0c0s0b0n0",
						State: "empty",
						Flag:  "None",
					},
				},
			},
			Force: false,
		},
		expectedOut: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "x0c0s0b0n0",
						State: base.StateEmpty.String(),
						Flag:  "None",
					},
				},
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("flag 'None' is invalid"),
	}, {
		in: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "x0c0s0b0n0",
						State: "empty",
						Flag:  "OK",
						Role:  "Master",
					},
				},
			},
			Force: false,
		},
		expectedOut: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:    "x0c0s0b0n0",
						State: base.StateEmpty.String(),
						Flag:  base.FlagOK.String(),
						Role:  "Master",
					},
				},
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("role 'Master' is invalid"),
	}, {
		in: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:      "x0c0s0b0n0",
						State:   "empty",
						Flag:    "OK",
						Role:    "Compute",
						NetType: "Fishing",
					},
				},
			},
			Force: false,
		},
		expectedOut: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:      "x0c0s0b0n0",
						State:   base.StateEmpty.String(),
						Flag:    base.FlagOK.String(),
						Role:    base.RoleCompute.String(),
						NetType: "Fishing",
					},
				},
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("netType 'Fishing' is invalid"),
	}, {
		in: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:      "x0c0s0b0n0",
						State:   "empty",
						Flag:    "OK",
						Role:    "Compute",
						NetType: "Sling",
						Arch:    "Bridge",
					},
				},
			},
			Force: false,
		},
		expectedOut: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:      "x0c0s0b0n0",
						State:   base.StateEmpty.String(),
						Flag:    base.FlagOK.String(),
						Role:    base.RoleCompute.String(),
						NetType: base.NetSling.String(),
						Arch:    "Bridge",
					},
				},
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("arch 'Bridge' is invalid"),
	}, {
		in: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:      "x00c0s0b0n0",
						State:   "empty",
						Flag:    "ok",
						Role:    "compute",
						NetType: "sling",
						Arch:    "x86",
					},
					&base.Component{
						ID:    "x0c0s0b0n1",
						State: "empty",
					},
				},
			},
			Force: false,
		},
		expectedOut: ComponentsPost{
			ComponentArray: base.ComponentArray{
				Components: []*base.Component{
					&base.Component{
						ID:      "x0c0s0b0n0",
						Type:    base.Node.String(),
						State:   base.StateEmpty.String(),
						Flag:    base.FlagOK.String(),
						Role:    base.RoleCompute.String(),
						NetType: base.NetSling.String(),
						Arch:    base.ArchX86.String(),
					},
					&base.Component{
						ID:    "x0c0s0b0n1",
						Type:  base.Node.String(),
						State: base.StateEmpty.String(),
						Flag:  base.FlagOK.String(),
					},
				},
			},
			Force: false,
		},
		expectedErr: nil,
	}}
	for i, test := range tests {
		err := test.in.VerifyNormalize()
		if err != nil && test.expectedErr == nil {
			t.Errorf("Test %v Failed: Received unexpected error '%v'", i, err)
		} else if test.expectedErr == nil {
			if !reflect.DeepEqual(test.expectedOut, test.in) {
				t.Errorf("Test %v Failed: Expected compPost struct '%v'; Received compPost struct '%v'", i, test.expectedOut, test.in)
			}
		}
	}
}

func TestNewCompPut(t *testing.T) {
	tests := []struct {
		compIn      base.Component
		forceIn     bool
		expectedOut *ComponentPut
		expectedErr error
	}{{
		compIn: base.Component{
			ID:    "x0c0s0b0n0",
			State: "empty",
		},
		forceIn: false,
		expectedOut: &ComponentPut{
			Component: base.Component{
				ID:    "x0c0s0b0n0",
				Type:  base.Node.String(),
				State: base.StateEmpty.String(),
				Flag:  base.FlagOK.String(),
			},
			Force: false,
		},
		expectedErr: nil,
	}, {
		compIn: base.Component{
			ID:    "foo",
			State: "empty",
		},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("xname ID 'foo' is invalid"),
	}, {
		compIn: base.Component{
			ID:    "x0c0s0b0n0",
			State: "Exploded",
		},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("state 'Exploded' is invalid"),
	}, {
		compIn: base.Component{
			ID:    "x0c0s0b0n0",
			State: "empty",
			Flag:  "None",
		},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("flag 'None' is invalid"),
	}, {
		compIn: base.Component{
			ID:    "x0c0s0b0n0",
			State: "empty",
			Flag:  "OK",
			Role:  "Master",
		},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("role 'Master' is invalid"),
	}, {
		compIn: base.Component{
			ID:      "x0c0s0b0n0",
			State:   "empty",
			Flag:    "OK",
			Role:    "Compute",
			NetType: "Fishing",
		},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("netType 'Fishing' is invalid"),
	}, {
		compIn: base.Component{
			ID:      "x0c0s0b0n0",
			State:   "empty",
			Flag:    "OK",
			Role:    "Compute",
			NetType: "Sling",
			Arch:    "Bridge",
		},
		forceIn:     false,
		expectedOut: nil,
		expectedErr: fmt.Errorf("arch 'Bridge' is invalid"),
	}, {
		compIn: base.Component{
			ID:      "x00c0s0b0n0",
			State:   "empty",
			Flag:    "ok",
			Role:    "compute",
			NetType: "sling",
			Arch:    "x86",
		},
		forceIn: false,
		expectedOut: &ComponentPut{
			Component: base.Component{
				ID:      "x0c0s0b0n0",
				Type:    base.Node.String(),
				State:   base.StateEmpty.String(),
				Flag:    base.FlagOK.String(),
				Role:    base.RoleCompute.String(),
				NetType: base.NetSling.String(),
				Arch:    base.ArchX86.String(),
			},
			Force: false,
		},
		expectedErr: nil,
	}}
	for i, test := range tests {
		out, err := NewCompPut(test.compIn, test.forceIn)
		if err != nil && test.expectedErr == nil {
			t.Errorf("Test %v Failed: Received unexpected error '%v'", i, err)
		} else if test.expectedErr == nil {
			if !reflect.DeepEqual(test.expectedOut, out) {
				t.Errorf("Test %v Failed: Expected compPut struct '%v'; Received compPut struct '%v'", i, test.expectedOut, out)
			}
		}
	}
}

func TestVerifyNormalizeCompPut(t *testing.T) {
	tests := []struct {
		in          ComponentPut
		expectedOut ComponentPut
		expectedErr error
	}{{
		in: ComponentPut{
			Component: base.Component{
				ID:    "x0c0s0b0n0",
				State: "empty",
			},
			Force: false,
		},
		expectedOut: ComponentPut{
			Component: base.Component{
				ID:    "x0c0s0b0n0",
				Type:  base.Node.String(),
				State: base.StateEmpty.String(),
				Flag:  base.FlagOK.String(),
			},
			Force: false,
		},
		expectedErr: nil,
	}, {
		in: ComponentPut{
			Component: base.Component{
				ID:    "foo",
				State: "empty",
			},
			Force: false,
		},
		expectedOut: ComponentPut{
			Component: base.Component{
				ID: "foo",
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("xname ID 'foo' is invalid"),
	}, {
		in: ComponentPut{
			Component: base.Component{
				ID:    "x0c0s0b0n0",
				State: "Exploded",
			},
			Force: false,
		},
		expectedOut: ComponentPut{
			Component: base.Component{
				ID:    "x0c0s0b0n0",
				State: "Exploded",
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("state 'Exploded' is invalid"),
	}, {
		in: ComponentPut{
			Component: base.Component{
				ID:    "x0c0s0b0n0",
				State: "empty",
				Flag:  "None",
			},
			Force: false,
		},
		expectedOut: ComponentPut{
			Component: base.Component{
				ID:    "x0c0s0b0n0",
				State: base.StateEmpty.String(),
				Flag:  "None",
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("flag 'None' is invalid"),
	}, {
		in: ComponentPut{
			Component: base.Component{
				ID:    "x0c0s0b0n0",
				State: "empty",
				Flag:  "OK",
				Role:  "Master",
			},
			Force: false,
		},
		expectedOut: ComponentPut{
			Component: base.Component{
				ID:    "x0c0s0b0n0",
				State: base.StateEmpty.String(),
				Flag:  base.FlagOK.String(),
				Role:  "Master",
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("role 'Master' is invalid"),
	}, {
		in: ComponentPut{
			Component: base.Component{
				ID:      "x0c0s0b0n0",
				State:   "empty",
				Flag:    "OK",
				Role:    "Compute",
				NetType: "Fishing",
			},
			Force: false,
		},
		expectedOut: ComponentPut{
			Component: base.Component{
				ID:      "x0c0s0b0n0",
				State:   base.StateEmpty.String(),
				Flag:    base.FlagOK.String(),
				Role:    base.RoleCompute.String(),
				NetType: "Fishing",
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("netType 'Fishing' is invalid"),
	}, {
		in: ComponentPut{
			Component: base.Component{
				ID:      "x0c0s0b0n0",
				State:   "empty",
				Flag:    "OK",
				Role:    "Compute",
				NetType: "Sling",
				Arch:    "Bridge",
			},
			Force: false,
		},
		expectedOut: ComponentPut{
			Component: base.Component{
				ID:      "x0c0s0b0n0",
				State:   base.StateEmpty.String(),
				Flag:    base.FlagOK.String(),
				Role:    base.RoleCompute.String(),
				NetType: base.NetSling.String(),
				Arch:    "Bridge",
			},
			Force: false,
		},
		expectedErr: fmt.Errorf("arch 'Bridge' is invalid"),
	}, {
		in: ComponentPut{
			Component: base.Component{
				ID:      "x00c0s0b0n0",
				State:   "empty",
				Flag:    "ok",
				Role:    "compute",
				NetType: "sling",
				Arch:    "x86",
			},
			Force: false,
		},
		expectedOut: ComponentPut{
			Component: base.Component{
				ID:      "x0c0s0b0n0",
				Type:    base.Node.String(),
				State:   base.StateEmpty.String(),
				Flag:    base.FlagOK.String(),
				Role:    base.RoleCompute.String(),
				NetType: base.NetSling.String(),
				Arch:    base.ArchX86.String(),
			},
			Force: false,
		},
		expectedErr: nil,
	}}
	for i, test := range tests {
		err := test.in.VerifyNormalize()
		if err != nil && test.expectedErr == nil {
			t.Errorf("Test %v Failed: Received unexpected error '%v'", i, err)
		} else if test.expectedErr == nil {
			if !reflect.DeepEqual(test.expectedOut, test.in) {
				t.Errorf("Test %v Failed: Expected compPost struct '%v'; Received compPost struct '%v'", i, test.expectedOut, test.in)
			}
		}
	}
}
