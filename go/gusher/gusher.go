package gusher

import (
	"encoding/json"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
)

type handler struct {
	*channels
}

func NewServeMux(prefix string) *http.ServeMux {
	g := &handler{newChannels()}
	mux := http.NewServeMux()
	mux.Handle(prefix+"/", sockjs.NewHandler(prefix, sockjs.DefaultOptions, g.handler))
	mux.Handle(prefix+"/api/", g.API())
	return mux
}

func (g *handler) API() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		channel := req.PostFormValue("channel")
		event := req.PostFormValue("event")
		data := req.PostFormValue("data")
		m := post{channel, event, data}
		payload, _ := json.Marshal(m)
		g.publish(channel, string(payload))
	})
}

func (g *handler) handleMessage(msg message, session sockjs.Session) {
	switch msg := msg.(type) {
	case messageSubscribe:
		g.get(msg.Channel).subscribe(session)
		log.Println("subscribed " + session.ID() + " to " + msg.Channel)
	case messageUnsubscribe:
		g.get(msg.Channel).unsubscribe(session)
		log.Println("unsubscribed " + session.ID() + " to " + msg.Channel)
	default:
		log.Fatal("I give up")
	}
}

func (g *handler) handler(session sockjs.Session) {
	log.Println("client connected")
	defer g.teardownSession(session)
	for {
		raw, err := session.Recv()
		if err != nil {
			log.Println("client disconnected")
			break
		}

		log.Println("msg rec'd: " + raw)
		msg, err := MessageUnmarshalJSON([]byte(raw))
		if err != nil {
			log.Println("error unmarshaling json: " + err.Error())
			break
		}

		g.handleMessage(msg, session)
	}
}
