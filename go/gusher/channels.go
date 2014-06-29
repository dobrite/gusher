package gusher

import ()

type channels struct {
	channels map[string]*channel
	pub      chan *publish
	unsub    chan *unsubscribe
	sub      chan *subscribe
}

type subscribe struct {
	channelName string
	c           chan string
}

type unsubscribe struct {
	channelName string
	c           chan string
}

type publish struct {
	channelName string
	payload     string
}

func newChannels() *channels {
	channels := &channels{
		channels: make(map[string]*channel),
		sub:      make(chan *subscribe),
		unsub:    make(chan *unsubscribe),
		pub:      make(chan *publish),
	}
	go channels.run()
	return channels
}

func (chs *channels) run() {
	for {
		select {
		case s := <-chs.sub:
			chs.subscribe(s.channelName, s.c)
		case u := <-chs.unsub:
			chs.subscribe(u.channelName, u.c)
		case p := <-chs.pub:
			chs.publish(p.channelName, p.payload)
		}
	}
}

func (chs *channels) _get(channelName string) *channel {
	if ch, ok := chs.channels[channelName]; ok {
		return ch
	}
	return chs.create(channelName)
}

func (chs *channels) create(channelName string) *channel {
	ch := newChannel()
	chs.channels[channelName] = ch
	return ch
}

func (chs *channels) publish(channelName string, payload string) {
	//check if channel is empty
	chs._get(channelName).pub <- payload
}

func (chs *channels) subscribe(channelName string, toSock chan string) {
	chs._get(channelName).sub <- toSock
}

func (chs *channels) unsubscribe(channelName string, toSock chan string) {
	//assert someone in here
	chs.channels[channelName].unsub <- toSock
	//check if empty and delete
}

func (chs *channels) unsubscribeAll(gsession *gsession) {

}
