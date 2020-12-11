// Copyright (c) 2019 Cray Inc. All Rights Reserved.
package main

import (
	"log"
	"os"
	"testing"
)

var ss SmD

func TestStartRFEventMonitor(t *testing.T) {
	ss.msgbusListen = ""
	ss.lg = log.New(os.Stdout, "", log.Lshortfile|log.LstdFlags|log.Lmicroseconds)
	if err := s.SetLogLevel(LogLevel(LOG_DEBUG)); err != nil {
		os.Exit(1)
	}
	// Should exit immediately without trying to connect.
	ss.StartRFEventMonitor()
	t.Logf("Test 0: Pass: StartRFEventMonitor()")
}
