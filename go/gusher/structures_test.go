package gusher

import (
	"testing"
)

func TestStructuresSubscribe(t *testing.T) {
	b := []byte(`{
		"event":"gusher:subscribe",
		"data":"{\"channel\":\"tester\",\"auth\":\"auth type stuff\",\"channel_data\":\"channel data type stuff\"}"
	}`)
	msg, err := MessageUnmarshalJSON(b)
	sub := msg.(messageSubscribe)
	if err != nil {
		t.Errorf("MessageUnmarshalJSON error: %s", err)
	}
	if sub.Channel != "tester" {
		t.Errorf("Channel Unmarshal error: %s", sub.Channel)
	}
	if sub.Auth != "auth type stuff" {
		t.Errorf("Auth Unmarshal error: %s", sub.Auth)
	}
	if sub.ChannelData != "channel data type stuff" {
		t.Errorf("ChannelData Unmarshal error: %s", sub.ChannelData)
	}
}

func TestStructuresUnsubscribe(t *testing.T) {
	b := []byte(`{
		"event":"gusher:unsubscribe",
		"data":"{\"channel\":\"tester\"}"
	}`)
	msg, err := MessageUnmarshalJSON(b)
	unsub := msg.(messageUnsubscribe)
	if err != nil {
		t.Errorf("MessageUnmarshalJSON error: %s", err)
	}
	if unsub.Channel != "tester" {
		t.Errorf("Channel Unmarshal error: %s", unsub.Channel)
	}
}
