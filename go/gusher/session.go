package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type session struct {
	session sockjs.Session
	in      chan string
}

func (s *session) listen() {
	for {
		msg := <-s.in
		s.session.Send(msg)
	}
}
