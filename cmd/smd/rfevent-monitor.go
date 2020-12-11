// Copyright (c) 2019 Cray Inc. All Rights Reserved.
package main

import (
	"time"
)

// This is the main thread that monitors the event bus
func (s *SmD) StartRFEventMonitor() {
	if s.msgbusListen != "" {
		if err := s.MsgBusConfig(s.msgbusListen); err != nil {
			s.LogAlways("WARNING: Cannot parse message bus host: %s", err)
		} else {
			// Loop forever trying to connect or reconnect if an error
			// occurs.
			for {
				if err := s.MsgBusReconnect(); err != nil {
					s.LogAlways("ERROR: Cannot connect to  message bus host: %s",
						err)
					s.LogAlways("Retrying msg bus connection in 5 seconds")
					time.Sleep(5 * time.Second)
					continue
				}
				s.LogAlways("Connected to message bus: %s:%d:%s",
					s.msgbusConfig.Host, s.msgbusConfig.Port, s.msgbusConfig.Topic)

				// Consume events from topic on message bus.  If an
				// error occurs, report it and reconnect.
				err := s.rfEventMonitor()
				s.LogAlways("ERROR: Event monitor returned '%s', reconnecting",
					err)
			}
		}
	} else {
		s.LogAlways("No message bus host given.")
	}
	s.LogAlways("Not listening for events on the message bus.")
	return
}

// Dequeue a single event and hand it off for processing.
func (s *SmD) rfEventMonitor() error {
	// Read the event
	s.LogAlways("Starting event monitor...")
	var err error
	for {
		// Read the event
		payload, err := s.MsgBusReadNext()
		if err != nil {
			s.Log(LOG_INFO, "Got error reading event: %s", err)
			return err
		}
		rfEventJob := NewJobRFEvent(payload, s)
		s.wpRFEvent.Queue(rfEventJob)
	}
	// Shouldn't happen
	return err
}
