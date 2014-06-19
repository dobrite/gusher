package gusher

import (
	"github.com/igm/pubsub"
)

type Channel struct {
	pubsub.Publisher
}

var Channels map[string]*Channel
