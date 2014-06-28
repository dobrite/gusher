package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
)

type connections struct {
	connections map[string]*session
	reg         <-chan *session
	unreg       <-chan *session
	fin         chan interface{}
}

func (conns *connections) run() {
	defer conns.teardown()
	for {
		select {
		case conn := <-conns.reg:
			conns.register(conn)
		case conn := <-conns.unreg:
			conns.unregister(conn.ID())
		case <-conns.fin:
			break
		}
	}
}

func (conns *connections) teardown() {
	log.Println("connections teardown")
}

func (conns *connections) register(session sockjs.Session) {
	conns.connections[session.ID()] = newSession(session)
}

func (conns *connections) unregister(id string) {
	delete(conns.connections, id)
}
