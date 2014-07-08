package gusher

import (
	"github.com/deckarep/golang-set"
)

type session struct {
	conn   *connection
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
		conn:   conn,
		subbed: mapset.NewThreadUnsafeSet(),
		id:     id,
	}

	// prob put this in connection
	// and create newConnection factory function
	sess.conn.tomb.Go(sess.conn.sender)
	sess.conn.tomb.Go(sess.conn.receiver)
	return sess
}
