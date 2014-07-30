package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type sockjsTransport struct {
	transport
	sock sockjs.Session
}

func newSockjsTransport(session sockjs.Session) transport {
	return &sockjsTransport{
		sock: session,
	}
}

func (t *sockjsTransport) close() {
	t.sock.Close(1, "")
}

func (t *sockjsTransport) send(msg string) error {
	return t.sock.Send(msg)
}

func (t *sockjsTransport) recv() (string, error) {
	return t.sock.Recv()
}
