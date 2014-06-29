package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"gopkg.in/tomb.v2"
	"log"
)

type gsession struct {
	s      sockjs.Session //don't know how to compose this
	toSock chan string
	toGush chan string
	t      tomb.Tomb
}

func newSession(sockjsSession sockjs.Session, toGush chan string, toSock chan string) *gsession {
	gs := &gsession{
		s:      sockjsSession,
		toSock: toSock,
		toGush: toGush,
	}
	gs.t.Go(gs.sender)
	gs.t.Go(gs.receiver)
	return gs
}

func (gs *gsession) stop() error {
	log.Println("session closed")
	gs.s.Close(1, "arghhhh") //do something better than this
	gs.t.Kill(nil)
	return gs.t.Wait()
	//remove session from all channels
}

//writepump
func (gs *gsession) sender() error {
	defer func() { log.Println("sender exiting") }()
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
	defer func() { log.Println("receiver exiting") }()
	for {
		if raw, err := gs.s.Recv(); err != nil {
			return err
		} else {
			select {
			case gs.toGush <- raw:
			case <-gs.t.Dying():
				return tomb.ErrDying
			}
		}
	}
}

//type msg struct {
//	raw string
//	err error
//}
//
//func (gs *gsession) receiver2() <-chan *msg {
//	defer func() { log.Println("receiver exiting") }()
//	c := make(chan *msg)
//	gs.t.Go(func() error {
//		for {
//			if raw, err := gs.s.Recv(); err != nil {
//				c <- &msg{"", err}
//				return err
//			} else {
//				c <- &msg{raw, nil}
//			}
//		}
//	})
//	return c
//}
