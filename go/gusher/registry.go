package gusher

import (
	"github.com/deckarep/golang-set"
	"gopkg.in/tomb.v2"
)

//TODO move mapset to clipperhouse.github.io/gen

type registry struct {
	sessionids mapset.Set
	sessions   map[string]*session
	channels   map[string]mapset.Set
	command    chan func()
	t          tomb.Tomb
}

func newRegistry() *registry {
	registry := &registry{
		sessionids: mapset.NewThreadUnsafeSet(),
		sessions:   make(map[string]*session),
		channels:   make(map[string]mapset.Set),
		command:    make(chan func()),
	}
	registry.t.Go(registry.run)
	return registry
}

func (r *registry) run() error {
	for {
		select {
		case command := <-r.command:
			command()
		case <-r.t.Dying():
			return tomb.ErrDying
		}
	}
}

// see http://blog.igormihalik.com/2012/12/sockjs-for-go.html
// for returning values
func (r *registry) subscribe(msg message, session *session) {
	r.command <- func() {
		channelName := msg.(messageSubscribe).Channel
		sID := session.id
		_, ok := r.channels[channelName]
		if !ok {
			r.channels[channelName] = mapset.NewThreadUnsafeSet()
		}
		r.channels[channelName].Add(sID)
		session.subbed.Add(channelName)
	}
}

func (r *registry) unsubscribe(msg message, session *session) {
	r.command <- func() {
		channelName := msg.(messageSubscribe).Channel
		r.channels[channelName].Remove(session.id)
		session.subbed.Remove(channelName)
	}
}

func (r *registry) publish(channelName string, payload string) {
	r.command <- func() {
		c := r.channels[channelName].Iter()
		for sID := range c {
			// whole server can be blocked by a slow client
			r.sessions[sID.(string)].toConn <- payload
		}
	}
}

func (r *registry) add(session *session) {
	r.command <- func() {
		sID := session.id
		r.sessionids.Add(sID)
		r.sessions[sID] = session
	}
}

func (r *registry) remove(session *session) {
	r.command <- func() {
		sID := session.id
		r.sessionids.Remove(sID)
		delete(r.sessions, sID)
		for channelName := range session.subbed.Iter() {
			r.channels[channelName.(string)].Remove(sID)
		}
	}
}

func (r *registry) send(session *session, payload string) {
	r.command <- func() {
		session.toConn <- payload
	}
}
