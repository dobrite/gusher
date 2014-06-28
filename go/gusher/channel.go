package gusher

import (
	"log"
)

type channel struct {
	subscribers map[string]*session
	sub         <-chan *session
	unsub       <-chan *session
	pub         <-chan string
	fin         chan struct{}
}

func newChannel() *channel {
	channel := &channel{subscribers: make(map[string]*session)}
	go channel.run()
	return channel
}

func (ch *channel) run() {
	defer ch.teardown()
	for {
		select {
		case subscriber := <-ch.sub:
			ch.add(subscriber.ID(), subscriber)
		case subscriber := <-ch.unsub:
			ch.remove(subscriber.ID())
		case payload := <-ch.pub:
			ch.broadcast(payload)
		case <-ch.fin:
			break
		}
	}
}

func (ch *channel) teardown() {
	log.Println("channel closing shop")
}

func (ch *channel) add(id string, subscriber *session) {
	ch.subscribers[id] = subscriber
}

func (ch *channel) remove(id string) {
	delete(ch.subscribers, id)
	log.Println("length: ")
	log.Println(len(ch.subscribers))
	if len(ch.subscribers) == 0 {
		close(ch.fin)
	}
}

func (ch *channel) broadcast(payload string) {
	for id, subscriber := range ch.subscribers {
		select {
		case subscriber.in <- payload:
		default:
			ch.remove(id)
		}
	}
}
