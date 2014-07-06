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

func (h *handler) handler(transport transport) {
	log.Println("client connected")
	toGush := make(chan string)
	toSock := make(chan string)
	session := newSession(transport, toGush, toSock)
	h.registry.add(session)
	h.registry.send(session, buildMessageConnectionEstablished(session.conn.trans.id()))
	go h.listen(session)
}

func (h *handler) teardown(session *session) {
	log.Println("client disconnected")
	h.registry.remove(session)
}

func (h *handler) listen(session *session) {
	for {
		if raw, ok := <-session.conn.toGush; ok {
			log.Println("msg rec'd: " + raw)

			msg, err := MessageUnmarshalJSON([]byte(raw))

			if err != nil {
				log.Println("error unmarshaling json: " + err.Error())
				break
			}

			h.handleMessage(msg, session)
		} else {
			break
		}
	}
	h.teardown(session)
}

func (h *handler) handleMessage(msg message, session *session) {
	switch msg := msg.(type) {
	case messageSubscribe:
		log.Println("subscribed to " + msg.Channel)
		h.registry.subscribe(msg, session)
	case messageUnsubscribe:
		log.Println("unsubscribed to " + msg.Channel)
		h.registry.unsubscribe(msg, session)
	default:
		log.Fatal("I give up")
	}
}
