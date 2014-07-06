package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"gopkg.in/tomb.v2"
	"log"
	"sync"
)

type sockjsTransport struct {
	transport
	sock    sockjs.Session
	closed  bool       //makes me sad
	closedm sync.Mutex //makes me sad
	tomb    tomb.Tomb
	toConn  chan string
	toGush  chan string
}

func newSockjsTransport(session sockjs.Session, toConn chan string, toGush chan string) transport {
	return &sockjsTransport{
		sock:   session,
		toGush: toGush,
		toConn: toConn,
	}
}

func (t *sockjsTransport) close_() {
	t.sock.Close(1, "arghhhh") //do something better than this
}

func (t *sockjsTransport) send_(msg string) error {
	return t.sock.Send(msg)
}

func (t *sockjsTransport) recv_() (string, error) {
	return t.sock.Recv()
}

func (t *sockjsTransport) id_() string {
	return t.sock.ID()
}

func (t *sockjsTransport) go_(f func() error) {
	t.tomb.Go(f)
}

func (t *sockjsTransport) close() error {
	t.closedm.Lock()
	defer t.closedm.Unlock()
	if !t.closed {
		t.closed = true
		log.Println("connection closed")
		t.close_()
		t.tomb.Kill(nil)
		close(t.toGush)
	}
	return t.tomb.Wait()
}

func (t *sockjsTransport) sender_() error {
	defer func() {
		log.Println("sender exiting")
		t.close()
	}()

	for {
		select {
		case msg, ok := <-t.toConn:
			if !ok {
				return nil // TODO return err?
			}
			if err := t.send_(msg); err != nil {
				return err
			}
		case <-t.tomb.Dying():
			return tomb.ErrDying
		}
	}
}

func (t *sockjsTransport) receiver_() error {
	defer func() {
		log.Println("receiver exiting")
		t.close()
	}()

	for {
		raw, err := t.recv_()
		if err != nil {
			return err
		}
		select {
		case t.toGush <- raw:
		case <-t.tomb.Dying():
			return tomb.ErrDying
		}
	}
}
