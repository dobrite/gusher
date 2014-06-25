package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"sync"
)

type channel struct {
	mutex    sync.Mutex
	sessions map[string]sockjs.Session
}

func newChannel() *channel {
	return &channel{sessions: make(map[string]sockjs.Session)}
}

func (ch *channel) publish(payload string) {
	//what happens when a session is added while looping?
	for _, session := range ch.sessions {
		session.Send(payload)
	}
}

func (ch *channel) subscribe(session sockjs.Session) {
	ch.mutex.Lock()
	ch.sessions[session.ID()] = session
	ch.mutex.Unlock()
}

func (ch *channel) unsubscribe(session sockjs.Session) {
	id := session.ID()
	_, ok := ch.sessions[id]
	if ok {
		ch.mutex.Lock()
		delete(ch.sessions, id)
		ch.mutex.Unlock()
	}
}
