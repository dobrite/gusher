package gusher

import (
	"github.com/igm/pubsub"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type channel struct {
	sessions map[string]sockjs.Session
	pubsub.Publisher
}

func newChannel() *channel {
	return &channel{sessions: make(map[string]sockjs.Session)}
}

func (ch *channel) publish(payload string) {
	for _, session := range ch.sessions {
		session.Send(payload)
	}
}

func (ch *channel) subscribe(session sockjs.Session) {
	ch.sessions[session.ID()] = session
}

func (ch *channel) unsubscribe(session sockjs.Session) {
	id := session.ID()
	_, ok := ch.sessions[id]
	if ok {
		delete(ch.sessions, id)
	}
}
