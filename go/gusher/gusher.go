package gusher

import (
	"encoding/json"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
)

type handler struct {
	*channels
	*registry
}

func NewServeMux(prefix string) *http.ServeMux {
	h := &handler{
		channels: newChannels(),
		registry: newRegistry(),
	}
	mux := http.NewServeMux()
	mux.Handle(prefix+"/", sockjs.NewHandler(prefix, sockjs.DefaultOptions, h.handler))
	mux.Handle(prefix+"/api/", h.API())
	return mux
}

func (h *handler) API() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			return
		}
		channel := req.PostFormValue("channel")
		event := req.PostFormValue("event")
		data := req.PostFormValue("data")
		m := post{channel, event, data}
		payload, _ := json.Marshal(m)
		h.publish(channel, string(payload))
	})
}

func (h *handler) handler(session sockjs.Session) {
	log.Println("client connected")
	in := make(chan string)
	gsession := newSession(session, in)
	h.registry.add <- gsession
	go h.listen(in, gsession)
}

func (h *handler) teardown(gsession *gsession) {
	log.Println("client disconnected")
	h.unsubscribeAll(gsession)
}

func (h *handler) listen(in chan string, gsession *gsession) {
	defer h.teardown(gsession)
	for {
		if raw, ok := <-in; ok {
			log.Println("msg rec'd: " + raw)

			msg, err := MessageUnmarshalJSON([]byte(raw))

			if err != nil {
				log.Println("error unmarshaling json: " + err.Error())
				break
			}

			h.handleMessage(msg, gsession)
		} else {
			break
		}
	}
}

func (h *handler) handleMessage(msg message, gsession *gsession) {
	switch msg := msg.(type) {
	case messageSubscribe:
		log.Println("subscribed " + gsession.ID() + " to " + msg.Channel)
		h.subscribe(msg.Channel, gsession)
	case messageUnsubscribe:
		log.Println("unsubscribed " + gsession.ID() + " to " + msg.Channel)
		h.unsubscribe(msg.Channel, gsession)
	default:
		log.Fatal("I give up")
	}
}
