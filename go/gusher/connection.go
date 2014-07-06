package gusher

import (
	"gopkg.in/tomb.v2"
	"log"
	"sync"
)

type connection struct {
	trans   transport
	tomb    tomb.Tomb
	toConn  chan string
	toGush  chan string
	closed  bool       //makes me sad
	closedm sync.Mutex //makes me sad
}

func (c *connection) close() error {
	c.closedm.Lock()
	defer c.closedm.Unlock()
	if !c.closed {
		c.closed = true
		log.Println("connection closed")
		c.trans.close() //do something better than this
		c.tomb.Kill(nil)
		close(c.toGush)
	}
	return c.tomb.Wait()
}

func (c *connection) sender() error {
	defer func() {
		log.Println("sender exiting")
		c.trans.close()
	}()

	for {
		select {
		case msg := <-c.toConn:
			if err := c.trans.send(msg); err != nil {
				return err
			}
		case <-c.tomb.Dying():
			return tomb.ErrDying
		}
	}
}

func (c *connection) receiver() error {
	defer func() {
		log.Println("receiver exiting")
		c.trans.close()
	}()

	for {
		raw, err := c.trans.recv()
		if err != nil {
			return err
		}
		select {
		case c.toGush <- raw:
		case <-c.tomb.Dying():
			return tomb.ErrDying
		}
	}
}
