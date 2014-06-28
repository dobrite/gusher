package gusher

import (
	"log"
)

type channel struct {
	subscribers map[string]chan string
	sub         chan *gsession
	unsub       chan *gsession
	pub         chan string
	fin         chan struct{}
}

func newChannel() *channel {
	channel := &channel{
		subscribers: make(map[string]chan string),
		sub:         make(chan *gsession),
		unsub:       make(chan *gsession),
		pub:         make(chan string),
		fin:         make(chan struct{}),
	}
	go channel.run()
	return channel
}

func (ch *channel) run() {
	defer ch.teardown()
	log.Println("channel running")
	for {
		select {
		case subscriber := <-ch.sub:
			log.Println("subscribing: " + subscriber.ID())
			ch.subscribe(subscriber.ID(), subscriber.in)
		case subscriber := <-ch.unsub:
			log.Println("unsubscribing: " + subscriber.ID())
			ch.unsubscribe(subscriber.ID())
		case payload := <-ch.pub:
			log.Println("payload recd: " + payload)
			ch.broadcast(payload)
		case <-ch.fin:
			break
		}
	}
	log.Println("channel done")
}

func (ch *channel) teardown() {
	log.Println("channel closing shop")
}

func (ch *channel) subscribe(id string, connc chan string) {
	ch.subscribers[id] = connc
}

func (ch *channel) unsubscribe(id string) {
	delete(ch.subscribers, id)
	log.Println("length: ")
	log.Println(len(ch.subscribers))
	if len(ch.subscribers) == 0 {
		close(ch.fin)
	}
}

func (ch *channel) broadcast(payload string) {
	log.Println("broadcasting")
	for _, connc := range ch.subscribers {
		connc <- payload
	}
}
