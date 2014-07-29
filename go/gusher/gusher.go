package gusher

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
	"net/http"
	"strings"
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
	mux.Handle("/pusher/auth", h.auth())
	return mux
}

func (h *handler) API() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(405)
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
	log.Println("sockjs connected")
	h.handle(newSockjsTransport(session))
}

func (h *handler) websocket(w http.ResponseWriter, req *http.Request) {
	log.Println("websocket connected")
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
		if raw, ok := <-session.toGush; ok {
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
		if strings.HasPrefix(msg.Channel, "private-") || strings.HasPrefix(msg.Channel, "presence-") {
			//set authTransport to 'ajax' (default)
			//POST to /pusher/auth w/ socket_id and channel_name
			//set authTransport to 'jsonp', also set authEndpoint (default to /pusher/auth)
			//JSONP to /pusher/auth w/ socket_id, channel_name and callback
			//render :text => params[:callback] + "(" + auth.to_json + ")", :content_type => 'application/javascript'
			//return if authorized application/json
			//{"auth":"278d425bdf160c739803:afaed3695da2ffd16931f457e338e6c9f2921fa133ce7dac49f529792be6304c","channel_data":"{\"user_id\":10,\"user_info\":{\"name\":\"Mr. Pusher\"}}"}
			//otherwise 403 Forbidden plain text
		} else {
			h.registry.subscribe(msg, session)
			h.registry.send(session, buildMessageSubscriptionSucceeded(msg.Channel))
		}
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
