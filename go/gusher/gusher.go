package gusher

import (
	"encoding/json"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
)

type handler struct {
	*registry
}

func NewServeMux(prefix string) *http.ServeMux {
	h := &handler{
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
		h.registry.publish(channel, string(payload))
	})
}

func (h *handler) handler(session sockjs.Session) {
	log.Println("client connected")
	toGush := make(chan string)
	toSock := make(chan string)
	gsession := newSession(session, toGush, toSock)
	h.registry.add(gsession)
	h.registry.send(gsession, buildMessageConnectionEstablished(gsession.s.ID()))
	go h.listen(gsession)
}

func (h *handler) teardown(gsession *gsession) {
	log.Println("client disconnected")
	h.registry.remove(gsession)
}

func (h *handler) listen(gsession *gsession) {
	for {
		if raw, ok := <-gsession.toGush; ok {
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
	h.teardown(gsession)
}

func (h *handler) handleMessage(msg message, gsession *gsession) {
	switch msg := msg.(type) {
	case messageSubscribe:
		log.Println("subscribed to " + msg.Channel)
		h.registry.subscribe(msg, gsession)
	case messageUnsubscribe:
		log.Println("unsubscribed to " + msg.Channel)
		h.registry.unsubscribe(msg, gsession)
	default:
		log.Fatal("I give up")
	}
}
