// Copyright 2019,2020 Hewlett Packard Enterprise Development LP
package rf

import (
	"testing"
	"time"
)

func TestSetHTTPClientTimeout(t *testing.T) {
	SetHTTPClientTimeout(50)
	if GetHTTPClientTimeout() != 50 {
		t.Errorf("Test 1: FAIL: bad value: %d", GetHTTPClientTimeout())
	}
	SetHTTPClientTimeout(-1)
	if GetHTTPClientTimeout() != 50 {
		t.Errorf("Test 2: FAIL: bad value: %d", GetHTTPClientTimeout())
	}
}

/*
func TestSetHTTPClientProxyURL(t *testing.T) {
	SetHTTPClientProxyURL("socks5://127.0.0.1:9999")
	if GetHTTPClientProxyURL() != "socks5://127.0.0.1:9999" {
		t.Errorf("Test 1: FAIL: bad value: %s", GetHTTPClientProxyURL())
	}
}

func TestSetHTTPClientInsecureSkipVerify(t *testing.T) {
	SetHTTPClientInsecureSkipVerify(false)
	if GetHTTPClientInsecureSkipVerify() != false {
		t.Errorf("Test 1: FAIL: wrong value (true)")
	}
}
*/

func TestRfDefaultClient(t *testing.T) {
	SetHTTPClientTimeout(35)
	client := RfDefaultClient()
	if client.InsecureClient.Timeout != time.Duration(35)*time.Second {
		t.Errorf("Test 1: FAIL: Got unexpected/no timeout")
	}
}

/*
func TestRfProxyClient(t *testing.T) {
	SetHTTPClientTimeout(66)
	client := RfProxyClient("asdfasdfasdf")
	if client.Timeout != time.Duration(66)*time.Second {
		t.Errorf("Test 1: FAIL: Got unexpected/no timeout")
	}
	// Don't know a way to verify other fields, we get a generic
	// interface type back and can't seem to cast it.
	SetHTTPClientTimeout(77)
	client = RfProxyClient("socks5://127.0.0.1:9999")
	if client.Timeout != time.Duration(77)*time.Second {
		t.Errorf("Test 2: FAIL: Got unexpected/no timeout")
	}
}
*/
