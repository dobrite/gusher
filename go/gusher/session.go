package gusher

import (
	"github.com/deckarep/golang-set"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"gopkg.in/tomb.v2"
	"log"
	"sync"
)

type gsession struct {
	s       sockjs.Session //don't know how to compose this
	toSock  chan string
	toGush  chan string
	t       tomb.Tomb
	closed  bool       //makes me sad
	closedm sync.Mutex //makes me sad
	subbed  mapset.Set
}

func newSession(sockjsSession sockjs.Session, toGush chan string, toSock chan string) *gsession {
	gs := &gsession{
		s:      sockjsSession,
		toSock: toSock,
		toGush: toGush,
		closed: false,
		subbed: mapset.NewThreadUnsafeSet(),
	}
	gs.t.Go(gs.sender)
	gs.t.Go(gs.receiver)
	return gs
}

func (gs *gsession) close() error {
	gs.closedm.Lock()
	defer gs.closedm.Unlock()
	if !gs.closed {
		gs.closed = true
		log.Println("session closed")
		gs.s.Close(1, "arghhhh") //do something better than this
		gs.t.Kill(nil)
		close(gs.toGush)
	}
	return gs.t.Wait()
}

//writepump
func (gs *gsession) sender() error {
	defer func() {
		log.Println("sender exiting")
		gs.close()
	}()

	for {
		select {
		case msg := <-gs.toSock:
			if err := gs.s.Send(msg); err != nil {
				return err
			}
		case <-gs.t.Dying():
			return tomb.ErrDying
		}
	}
}

//readpump
func (gs *gsession) receiver() error {
	defer func() {
		log.Println("receiver exiting")
		gs.close()
	}()

	for {
		raw, err := gs.s.Recv()
		if err != nil {
			return err
		}
		select {
		case gs.toGush <- raw:
		case <-gs.t.Dying():
			return tomb.ErrDying
		}
	}
}
