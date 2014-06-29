package gusher

import (
	"sync"
)

type channels struct {
	channels map[string]*channel
	mutex    sync.Mutex
}

func newChannels() *channels {
	channels := &channels{
		channels: make(map[string]*channel),
	}
	go channels.run()
	return channels
}

func (chs *channels) run() {
	for {
		select {}
	}
}

func (chs *channels) _get(channelName string) *channel {
	if ch, ok := chs.channels[channelName]; ok {
		return ch
	}
	return chs.create(channelName)
}

func (chs *channels) create(channelName string) *channel {
	chs.mutex.Lock()
	defer chs.mutex.Unlock()
	ch := newChannel()
	//switch to channel
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
