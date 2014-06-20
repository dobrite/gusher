package gusher

import (
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
	chs.get(channelName).Publish(payload)
}

func (chs *channels) subChannel(channelName string) <-chan interface{} {
	sc, _ := chs.get(channelName).SubChannel(nil)
	return sc
}
