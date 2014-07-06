package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type sockjsTransport struct {
	transport
	sock sockjs.Session
}

func (s *sockjsTransport) close() {
	s.sock.Close(1, "arghhhh") //do something better than this
}

func (s *sockjsTransport) send(msg string) error {
	return s.sock.Send(msg)
}

func (s *sockjsTransport) recv() (string, error) {
	return s.sock.Recv()
}

func (s *sockjsTransport) id() string {
	return s.sock.ID()
}
