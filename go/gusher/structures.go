package gusher

import (
	"encoding/json"
	"fmt"
	"github.com/igm/sockjs-go/sockjs"
)

type post struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Data    string `json:"data"`
}

type data interface {
	handle(*handler, *sockjs.Session)
	//prob want to do json.Marshal
}

type dataEvent struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data,string"`
}

// client -> gusher
type dataSubscribe struct {
	Channel     string `json:"channel"`
	Auth        string `json:"auth"`
	ChannelData string `json:"channel_data"`
}

func (d dataSubscribe) handle(g *handler, session *sockjs.Session) {
	g.channels.subscribe(d.Channel, session)
}

// client -> gusher
type dataUnsubscribe struct {
	Channel string `json:"channel"`
}

func (d dataUnsubscribe) handle(g *handler, session *sockjs.Session) {

}

// gusher -> client
type dataError struct {
	Message string `json:"message"`
	Code    uint16 `json:"code"`
}

// gusher -> client
type dataConnectionEstablished struct {
	SocketId        string `json:"socket_id"`
	ActivityTimeout uint8  `json:"activity_timeout"`
}

func MessageUnmarshalJSON(b []byte) (d data, err error) {
	var env dataEvent
	err = json.Unmarshal(b, &env)
	if err != nil {
		return
	}
	switch env.Event {
	case "gusher:subscribe":
		var data dataSubscribe
		err = json.Unmarshal(env.Data, &data)
		d = data
	case "gusher:unsubscribe":
		var data dataUnsubscribe
		err = json.Unmarshal(env.Data, &data)
		d = data
	default:
		err = fmt.Errorf("%s is not a recognized event", env.Event)
	}
	return
}
