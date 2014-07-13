package gusher

import (
	"testing"
)

var privateChannelName = "private-foobar"
var presenceChannelName = "presence-foobar"
var socketId = "1234.1234"
var userData = "{\"user_id\":10,\"user_info\":{\"name\":\"Mr. Pusher\"}}"

var expectedDigest = "58df8b0c36d6982b82c3ecf6b4662e34fe8c25bba48f5369f135bf843651c3a4"
var expectedAuth = "278d425bdf160c739803" + ":" + expectedDigest
var expectedFullDigest = "afaed3695da2ffd16931f457e338e6c9f2921fa133ce7dac49f529792be6304c"
var expectedFullAuth = "278d425bdf160c739803" + ":" + expectedFullDigest

func TestSign(t *testing.T) {
	hexdigest := sign(socketId + ":" + privateChannelName)
	if hexdigest != expectedDigest {
		t.Errorf("hexdigest did not match expected digest")
	}
}

func TestFullSign(t *testing.T) {
	hexdigest := sign(socketId + ":" + presenceChannelName + ":" + userData)
	if hexdigest != expectedFullDigest {
		t.Errorf("hexdigest did not match expected full digest")
	}
}

func TestAuth(t *testing.T) {
	auth := auth(socketId, privateChannelName)
	if auth != expectedAuth {
		t.Errorf("auth did not match expected auth")
	}
}

func TestFullAuth(t *testing.T) {
	auth := auth(socketId, presenceChannelName, userData)
	if auth != expectedFullAuth {
		t.Errorf("auth did not match expected full auth")
	}
}
