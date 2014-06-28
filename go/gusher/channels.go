package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
)

type channelName string
type sessionID string

type channels struct {
	channels map[channelName]*channel
	fins     map[sessionID]chan struct{}
}

func newChannels() *channels {
	channels := &channels{
		channels: make(map[channelName]*channel),
		fins:     make(map[sessionID]chan struct{}),
	}
	go channels.run()
	return channels
}

func (chs *channels) run() {
	for {
		select {}
	}
}

func (chs *channels) get(channelName channelName) *channel {
	if ch, ok := chs.channels[channelName]; ok {
		return ch
	}
	ch := newChannel()
	chs.channels[channelName] = ch
	return ch
}

func (chs *channels) publish(channelName channelName, payload string) {
	//check if channel is empty
	chs.get(channelName).publish(payload)
}

func (chs *channels) subscribe(channelName channelName, session sockjs.Session) {
	chs.get(channelName).subscribe(session)
}

func (chs *channels) unsubscribe(channelName channelName, session sockjs.Session) {
	//assert someone in here
	chs.channels[channelName].unsubscribe(session)
	//check if empty and delete
}
