package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"log"
)

type gsession struct {
	sockjs.Session
	in  chan string
	out chan string
}

func newSession(sockjsSession sockjs.Session, out chan string) *gsession {
	gsession := &gsession{
		sockjsSession,
		make(chan string),
		out,
	}
	go gsession.sender()
	go gsession.receiver()
	return gsession
}

func (g *gsession) teardown() {
	log.Println("session closed")
	g.Close(1, "") //do something better than this
	close(g.in)
	close(g.out)
	//remove session from all channels
}

//writepump
func (g *gsession) sender() {
	for {
		if msg, ok := <-g.in; ok {
			g.Send(msg)
		}
		log.Println("sender closed")
		break
	}
}

//readpump
func (g *gsession) receiver() {
	defer g.teardown()
	for {
		raw, err := g.Recv()
		if err != nil {
			break
		}
		g.out <- raw
	}
	log.Println("receiver closed")
}
