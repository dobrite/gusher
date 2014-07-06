package gusher

import (
	"github.com/deckarep/golang-set"
)

type session struct {
	conn   *connection
	subbed mapset.Set
	id     string
}

func newSession(transport transport, toGush chan string, toConn chan string) *session {
	conn := &connection{
		trans:  transport,
		toConn: toConn,
		toGush: toGush,
		closed: false,
	}

	sess := &session{
		conn:   conn,
		subbed: mapset.NewThreadUnsafeSet(),
		id:     conn.trans.id(),
	}

	// prob put this in connection
	// and create newConnection factory function
	sess.conn.tomb.Go(sess.conn.sender)
	sess.conn.tomb.Go(sess.conn.receiver)
	return sess
}
