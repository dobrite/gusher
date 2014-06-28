package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
)

type session struct {
	sockjs.Session
	in  chan string
	out chan string
}

func newSession(sockjsSession sockjs.Session, out chan string) *session {
	session := &session{sockjsSession, make(chan string), out}
	session.setup()
	return session
}

func (s *session) setup() {
	log.Println("client connected")
	go s.sendPump()
	go s.recvPump()
	log.Println("client disconnected")
}

func (s *session) teardown() {
	log.Println("session closed")
	s.Close(1, "") //do something better than this
	//remove session from all channels
}

//writepump
func (s *session) sendPump() {
	for {
		msg := <-s.in
		s.Send(msg)
	}
}

//readpump
func (s *session) recvPump() {
	defer s.teardown()
	for {
		raw, err := s.Recv()
		if err != nil {
			break
		}
		//should be like handler.handle <- raw
	}
}
