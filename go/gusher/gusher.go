package gusher

import (
	"encoding/json"
	"github.com/igm/sockjs-go/sockjs"
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

func (g *handler) handler(session sockjs.Session) {
	log.Println("Client connected")
	//chat.Publish("[info] chatter joined")
	//defer chat.Publish("[info] chatter left")
	closedSession := make(chan struct{})
	defer close(closedSession)
	defer session.Close(1, "")
	go g.subscribe(session, closedSession)
	for {
		if raw, err := session.Recv(); err == nil {
			log.Println("Msg rec'd: " + raw)
			if msg, err := MessageUnmarshalJSON([]byte(raw)); err == nil {
				msg.handle(g, &session)
				continue
			}
			log.Println("Error unmarshaling JSON: " + err.Error())
			continue
		}
		log.Println("Client disconnected")
		break
	}
	log.Println("Session closed")
}

func (g *handler) subscribe(session sockjs.Session, closedSession chan struct{}) {
	reader := g.subChannel("test-channel")
	for {
		select {
		case <-closedSession:
			log.Println("subscribe closed")
			return
		case message := <-reader:
			msg := message.(string)
			log.Println("Msg sent: " + msg)
			if err := session.Send(msg); err != nil {
				return
			}
		}
	}
}
