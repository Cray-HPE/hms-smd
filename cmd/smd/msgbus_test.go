// Copyright (c) 2019 Cray Inc. All Rights Reserved.
package main

import (
	msgbus "stash.us.cray.com/HMS/hms-msgbus"
	"testing"
)

///////////////////////////////////////////////////////////////////////////////
// Unit Tests
///////////////////////////////////////////////////////////////////////////////

func TestMsgBusConfig(t *testing.T) {
	testargs := []struct {
		arg      string
		positive bool
	}{{
		"kafka.sma.svc.cluster.local",
		false,
	}, {
		"10.23.33.22:9092:events",
		true,
	}, {
		"10.23.33.22:9092",
		false,
	}, {
		"10.23.33.22:rf-events",
		false,
	}, {
		"kafka.sma.svc.cluster.local:9092:rf-events",
		true,
	}}

	s.msgbusHandle = &msgbus.MsgBusReader_Kafka{}
	if err := s.MsgBusConfig("10.23.33.22:9092:events"); err == nil {
		t.Errorf("Test 0: FAIL: MsgBusConfig with non-nil connection: no err")
	}
	s.msgbusHandle = nil

	for i, targ := range testargs {
		if err := s.MsgBusConfig(targ.arg); err != nil {
			if targ.positive == true {
				t.Errorf("Test %d: FAIL: MsgBusConfig(%s) unexpected error: %s",
					i, targ.arg, err)
				continue
			}
		} else if targ.positive == false {
			t.Errorf("Test %d: FAIL: MsgBusConfig(%s) expected error return",
				i, targ.arg)
			continue
		}
		t.Logf("Test %d: Pass: MsgBusConfig(%s)", i, targ.arg)
	}
}

func TestMsgBusConnect(t *testing.T) {
	t.Logf("Testing MsgBusConnect()...")
	s.msgbusHandle = &msgbus.MsgBusReader_Kafka{}
	if err := s.MsgBusConnect(); err != error(MsgBusAlreadyConnected) {
		t.Errorf("Test 0a: FAIL: MsgBusConnect with non-nil conn: wrong error")
	} else {
		t.Logf("Test 0a: Pass: MsgBusConnect()")
	}
	s.msgbusHandle = nil
	s.msgbusConfig = msgbus.MsgBusConfig{}
	if err := s.MsgBusConnect(); err != error(MsgBusMissingHostSpec) {
		t.Errorf("Test 0b: FAIL: MsgBusConnect without config: wrong error")
	} else {
		t.Logf("Test 0b: Pass: MsgBusConnect()")
	}
}

func TestMsgBusDisconnect(t *testing.T) {
	s.msgbusHandle = nil
	if err := s.MsgBusDisconnect(); err != nil {
		t.Errorf("Test 0a: FAIL: MsgBusDisconnect(): Got error: %s", err)
	} else {
		t.Logf("Test 0a: Pass: MsgBusDisconnect()")
	}
	if s.msgbusHandle != nil {
		t.Errorf("Test 0b: FAIL: MsgBusDisconnect(): Got handle != nil")
	} else {
		t.Logf("Test 0b: Pass: MsgBusDisconnect()")
	}
}
