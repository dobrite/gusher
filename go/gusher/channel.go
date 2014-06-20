package gusher

import (
	"github.com/igm/pubsub"
	"github.com/igm/sockjs-go/sockjs"
)

type channel struct {
	sessions map[string]*sockjs.Session
	pubsub.Publisher
}

func newChannel() *channel {
	return &channel{sessions: make(map[string]*sockjs.Session)}
}
