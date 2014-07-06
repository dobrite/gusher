package gusher

import (
	"github.com/deckarep/golang-set"
)

type session struct {
	conn   *connection
	subbed mapset.Set
	id     string
}

func newSession(transport transport) *session {
	conn := &connection{
		trans: transport,
	}

	sess := &session{
		conn:   conn,
		subbed: mapset.NewThreadUnsafeSet(),
		id:     conn.trans.id_(),
	}

	// prob put this in connection
	// and create newConnection factory function
	sess.conn.trans.go_(sess.conn.trans.sender_)
	sess.conn.trans.go_(sess.conn.trans.receiver_)
	return sess
}
