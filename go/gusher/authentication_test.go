package gusher

import (
	"testing"
)

var channelName = "private-foobar"
var socketId = "1234.1234"
var expectedHMAC = "58df8b0c36d6982b82c3ecf6b4662e34fe8c25bba48f5369f135bf843651c3a4"
var expectedAuth = "278d425bdf160c739803" + ":" + expectedHMAC

func TestAuth(t *testing.T) {
	auth := auth(socketId, channelName)
	if auth != expectedAuth {
		t.Errorf("auth did not match expected auth")
	}
}

func TestSign(t *testing.T) {
	hexdigest := sign(socketId + ":" + channelName)
	if hexdigest != expectedHMAC {
		t.Errorf("hexdigest did not match expected")
	}
}
