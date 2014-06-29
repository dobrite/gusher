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
	toGush := make(chan string)
	toSock := make(chan string)
	gsession := newSession(session, toGush, toSock)
	h.registry.add <- gsession
	go h.listen(toGush, toSock)
}

func (h *handler) teardown() {
	log.Println("client disconnected")
}

func (h *handler) listen(toGush chan string, toSock chan string) {
	defer h.teardown()
	for {
		if raw, ok := <-toGush; ok {
			log.Println("msg rec'd: " + raw)

			msg, err := MessageUnmarshalJSON([]byte(raw))

			if err != nil {
				log.Println("error unmarshaling json: " + err.Error())
				break
			}

			h.handleMessage(msg, toSock)
		} else {
			break
		}
	}
}

func (h *handler) handleMessage(msg message, toSock chan string) {
	switch msg := msg.(type) {
	case messageSubscribe:
		log.Println("subscribed to " + msg.Channel)
		h.subscribe(msg.Channel, toSock)
	case messageUnsubscribe:
		log.Println("unsubscribed to " + msg.Channel)
		h.unsubscribe(msg.Channel, toSock)
	default:
		log.Fatal("I give up")
	}
}
