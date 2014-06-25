package gusher

import (
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"sync"
)

type channels struct {
	mutex    sync.Mutex
	channels map[string]*channel
}

func newChannels() *channels {
	return &channels{channels: make(map[string]*channel)}
}

func (chs *channels) get(channelName string) *channel {
	if ch, ok := chs.channels[channelName]; ok {
		return ch
	}
	chs.mutex.Lock()
	defer chs.mutex.Unlock()
	ch := newChannel()
	chs.channels[channelName] = ch
	return ch
}

func (chs *channels) publish(channelName string, payload string) {
	//check if channel is empty
	chs.get(channelName).publish(payload)
}

func (chs *channels) subscribe(channelName string, session sockjs.Session) {
	chs.get(channelName).subscribe(session)
}

func (chs *channels) unsubscribe(channelName string, session sockjs.Session) {
	//assert someone in here
	chs.channels[channelName].unsubscribe(session)
	//check if empty and delete
}
