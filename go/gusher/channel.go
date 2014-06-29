package gusher

import (
	"log"
)

type channel struct {
	subscribers map[string]chan string
	sub         chan (chan string)
	unsub       chan (chan string)
	pub         chan string
	fin         chan struct{}
}

func newChannel() *channel {
	channel := &channel{
		subscribers: make(map[string]chan string),
		sub:         make(chan (chan string)),
		unsub:       make(chan (chan string)),
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
		case c := <-ch.sub:
			ch.subscribe("", c)
		case <-ch.unsub:
			ch.unsubscribe("")
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
