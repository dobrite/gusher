package gusher

import (
	"github.com/deckarep/golang-set"
)

type session struct {
	*connection
	subbed mapset.Set
	id     string
}

func newSession(id string, transport transport, toConn chan string, toGush chan string) *session {
	conn := &connection{
		trans:  transport,
		toConn: toConn,
		toGush: toGush,
	}

	sess := &session{
		conn,
		mapset.NewThreadUnsafeSet(),
		id,
	}

	// prob put this in connection
	// and create newConnection factory function
	sess.tomb.Go(sess.sender)
	sess.tomb.Go(sess.receiver)
	return sess
}
