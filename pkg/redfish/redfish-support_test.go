// Copyright (c) 2019 Cray Inc. All Rights Reserved.
package rf

import (
	"testing"
)

type checkMsgId struct {
	input    *EventRecord
	Registry string
	Version  string
	MsgId    string
}

// Others tested implicitly.

func TestEventRecordMsgId(t *testing.T) {
	var tests = []checkMsgId{{
		input:    &EventRecord{MessageId: " Alert.1.0.0.ShutDown  "},
		Registry: "Alert",
		Version:  "1.0.0",
		MsgId:    "ShutDown",
	}, {
		input:    &EventRecord{MessageId: "Alert.1.0.ShutDown"},
		Registry: "Alert",
		Version:  "1.0",
		MsgId:    "ShutDown",
	}, {
		input:    &EventRecord{MessageId: "Alert.v1.ShutDown"},
		Registry: "Alert",
		Version:  "v1",
		MsgId:    "ShutDown",
	}, {
		input:    &EventRecord{MessageId: "Alert.ShutDown"},
		Registry: "Alert",
		Version:  "",
		MsgId:    "ShutDown",
	}, {
		input:    &EventRecord{MessageId: "ShutDown"},
		Registry: "",
		Version:  "",
		MsgId:    "ShutDown",
	},
	}
	for i, test := range tests {
		reg, ver, id := EventRecordMsgId(test.input)
		if reg != test.Registry || ver != test.Version || id != test.MsgId {
			t.Errorf("Testcase %da: FAIL Got '%s/%s/%s', expected '%s/%s/%s'",
				i, reg, ver, id, test.Registry, test.Version, test.MsgId)
		} else {
			t.Logf("Testcase %da: PASS Got '%s/%s/%s'",
				i, reg, ver, id)
		}
	}
}

type checkVersionFields struct {
	input    string
	delim    string
	num      int
	output   string
	included int
}

func TestVersionFields(t *testing.T) {
	var tests = []checkVersionFields{{
		input:    "1.0.0",
		delim:    ".",
		num:      1,
		output:   "1",
		included: 1,
	}, {
		input:    "1.0.0",
		delim:    ".",
		num:      1,
		output:   "1",
		included: 1,
	}, {
		input:    "1.0.0",
		delim:    ".",
		num:      2,
		output:   "1.0",
		included: 2,
	}, {
		input:    "1.0.0",
		delim:    ".",
		num:      3,
		output:   "1.0.0",
		included: 3,
	}, {
		input:    "1.0.0",
		delim:    ".",
		num:      4,
		output:   "1.0.0",
		included: 3,
	}, {
		input:    "100",
		delim:    ".",
		num:      2,
		output:   "100",
		included: 1,
	},
	}
	for i, test := range tests {
		output, included := VersionFields(test.input, test.delim, test.num)
		if output != test.output || included != test.included {
			t.Errorf("Testcase %da: FAIL Got '%s': %d, expected '%s': %d",
				i, output, included, test.output, test.included)
		} else {
			t.Logf("Testcase %da: PASS Got '%s': %d", i, output, included)
		}
	}
}
