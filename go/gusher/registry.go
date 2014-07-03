package gusher

import (
	"github.com/deckarep/golang-set"
	"gopkg.in/tomb.v2"
)

//TODO move mapset to clipperhouse.github.io/gen

type registry struct {
	sessionids mapset.Set
	gsessions  map[string]*gsession
	channels   map[string]mapset.Set
	command    chan func()
	t          tomb.Tomb
}

func newRegistry() *registry {
	registry := &registry{
		sessionids: mapset.NewThreadUnsafeSet(),
		gsessions:  make(map[string]*gsession),
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
func (r *registry) subscribe(msg message, gsession *gsession) {
	r.command <- func() {
		channelName := msg.(messageSubscribe).Channel
		sID := gsession.s.ID()
		_, ok := r.channels[channelName]
		if !ok {
			r.channels[channelName] = mapset.NewThreadUnsafeSet()
		}
		r.channels[channelName].Add(sID)
		gsession.subbed.Add(channelName)
	}
}

func (r *registry) unsubscribe(msg message, gsession *gsession) {
	r.command <- func() {
		channelName := msg.(messageSubscribe).Channel
		r.channels[channelName].Remove(gsession.s.ID())
		gsession.subbed.Remove(channelName)
	}
}

func (r *registry) publish(channelName string, payload string) {
	r.command <- func() {
		c := r.channels[channelName].Iter()
		for sID := range c {
			r.gsessions[sID.(string)].toSock <- payload
		}
	}
}

func (r *registry) add(gsession *gsession) {
	r.command <- func() {
		sID := gsession.s.ID()
		r.sessionids.Add(sID)
		r.gsessions[sID] = gsession
	}
}

func (r *registry) remove(gsession *gsession) {
	r.command <- func() {
		sID := gsession.s.ID()
		r.sessionids.Remove(sID)
		delete(r.gsessions, sID)
		for channelName := range gsession.subbed.Iter() {
			r.channels[channelName.(string)].Remove(sID)
		}
	}
}

func (r *registry) send(gsession *gsession, payload string) {
	r.command <- func() {
		gsession.s.Send(payload)
	}
}
