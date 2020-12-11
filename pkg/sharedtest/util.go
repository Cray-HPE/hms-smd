// Copyright 2018-2020 Cray Inc. All Rights Reserved.
//
// Except as permitted by contract or express written permission of Cray Inc.,
// no part of this work or its content may be modified, used, reproduced or
// disclosed in any form. Modifications made without express permission of
// Cray Inc. may damage the system the software is installed within, may
// disqualify the user from receiving support from Cray Inc. under support or
// maintenance contracts, or require additional support services outside the
// scope of those contracts to repair the software or system.

package sharedtest

import (
	base "stash.us.cray.com/HMS/hms-base"
)

func CompareComp(comp1 *base.Component, comp2 *base.Component) bool {
	return CompareComps([]*base.Component{comp1}, []*base.Component{comp2})
}

func CompareComps(comps1 []*base.Component, comps2 []*base.Component) bool {
	if len(comps1) != len(comps2) {
		return false
	}
	// Treat nil entries specially.  There should be the same number,
	// and the remaining non-nil entries should be checked against
	// each other.
	nils := 0
	foundNils := 0
	for _, c1 := range comps1 {
		if c1 == nil {
			nils += 1
		}
	}
	for _, c2 := range comps2 {
		if c2 == nil {
			foundNils += 1
		}
	}
	if nils != foundNils {
		return false
	}
	for _, c1 := range comps1 {
		if c1 == nil {
			continue
		}
		found := false
		for _, c2 := range comps2 {
			if c2 == nil {
				continue
			}
			if c1.ID != c2.ID {
				continue
			}
			if c1.Type != c2.Type {
				continue
			}
			if c1.State != c2.State {
				continue
			}
			if c1.Flag != c2.Flag {
				continue
			}
			if c1.Enabled != nil {
				if c2.Enabled != nil {
					if *c1.Enabled != *c2.Enabled {
						continue
					}
					// Matches
				} else {
					continue
				}
			} else if c2.Enabled != nil {
				continue
			}
			if c1.SwStatus != c2.SwStatus {
				continue
			}
			if c1.Role != c2.Role {
				continue
			}
			if c1.SubRole != c2.SubRole {
				continue
			}
			if c1.NID != c2.NID {
				continue
			}
			if c1.Subtype != c2.Subtype {
				continue
			}
			if c1.NetType != c2.NetType {
				continue
			}
			if c1.Arch != c2.Arch {
				continue
			}
			if c1.Class != c2.Class {
				continue
			}
			found = true
			break
		}
		if !found {
			return false
		}
	}
	return true
}
