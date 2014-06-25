package gusher

import (
	"encoding/json"
	"fmt"
)

type post struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Data    string `json:"data"`
}

type message interface {
	//prob want to do json.Marshal
}

type messageEvent struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data,string"`
}

// client -> gusher
type messageSubscribe struct {
	Channel     string `json:"channel"`
	Auth        string `json:"auth"`
	ChannelData string `json:"channel_data"`
}

// client -> gusher
type messageUnsubscribe struct {
	Channel string `json:"channel"`
}

// gusher -> client
type messageError struct {
	Message string `json:"message"`
	Code    uint16 `json:"code"`
}

// gusher -> client
type messageConnectionEstablished struct {
	SocketId        string `json:"socket_id"`
	ActivityTimeout uint8  `json:"activity_timeout"`
}

func MessageUnmarshalJSON(b []byte) (msg message, err error) {
	var event messageEvent
	err = json.Unmarshal(b, &event)
	if err != nil {
		return
	}
	switch event.Event {
	case "gusher:subscribe":
		var msgSub messageSubscribe
		err = json.Unmarshal(event.Data, &msgSub)
		msg = msgSub
	case "gusher:unsubscribe":
		var msgUnsub messageUnsubscribe
		err = json.Unmarshal(event.Data, &msgUnsub)
		msg = msgUnsub
	default:
		err = fmt.Errorf("%s is not a recognized event", event.Event)
	}
	return
}
