package gusher

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
)

type handler struct {
	*registry
}

func NewServeMux(prefix string, appName string) *http.ServeMux {
	h := &handler{
		registry: newRegistry(),
	}
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/"+appName, h.websocket)
	mux.Handle(prefix+"/", sockjs.NewHandler(prefix, sockjs.DefaultOptions, h.sockjs))
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
		m := post{
			Event:   event,
			Channel: channel,
			Data:    data,
		}
		payload, _ := json.Marshal(m)
		h.registry.publish(channel, string(payload))
	})
}

func (h *handler) sockjs(session sockjs.Session) {
	//{"hostname":"sock34.pusher.com","websocket":false,"origins":["*:*"],"cookie_needed":false,"entropy":980616283,"server_heartbeat_interval":25000}
	log.Println("sockjsssssss!")
	h.handle(newSockjsTransport(session))
}

func (h *handler) websocket(w http.ResponseWriter, req *http.Request) {
	log.Println("websocketttttt!")
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	h.handle(newWebsocketTransport(ws))
}

func (h *handler) handle(transport transport) {
	toGush := make(chan string)
	toConn := make(chan string)
	u4, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		return
	}
	id := u4.String()
	log.Printf("client connected: %s", id)
	session := newSession(id, transport, toConn, toGush)
	h.registry.add(session)
	h.registry.send(session, buildMessageConnectionEstablished(session.id))
	go h.listen(session)
}

func (h *handler) teardown(session *session) {
	log.Printf("client disconnected: %s", session.id)
	h.registry.remove(session)
}

func (h *handler) listen(session *session) {
	for {
		if raw, ok := <-session.conn.toGush; ok {
			// first message from sockjs really is {"path": "app/tester..."}
			log.Println("msg rec'd: " + raw)

			msg, err := MessageUnmarshalJSON([]byte(raw))

			if err != nil {
				log.Println("error unmarshaling json: " + err.Error())
				//break
			} else {
				h.handleMessage(msg, session)
			}
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
		h.registry.send(session, buildMessageSubscriptionSucceeded(msg.Channel))
	case messageUnsubscribe:
		log.Println("unsubscribed to " + msg.Channel)
		h.registry.unsubscribe(msg, session)
	case messagePing:
		log.Println("ping")
		h.registry.send(session, buildMessagePong())
	default:
		log.Fatal("I give up")
	}
}
