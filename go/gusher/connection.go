package gusher

import (
	"gopkg.in/tomb.v2"
	"log"
	"sync"
)

type connection struct {
	trans   transport
	closed  bool       //makes me sad
	closedm sync.Mutex //makes me sad
	tomb    tomb.Tomb
	toConn  chan string
	toGush  chan string
}

func (c *connection) close() error {
	c.closedm.Lock()
	defer c.closedm.Unlock()
	if !c.closed {
		c.closed = true
		log.Println("connection closed")
		c.trans.close()
		c.tomb.Kill(nil)
		close(c.toGush)
	}
	return c.tomb.Wait()
}

func (c *connection) sender() error {
	//	ticker := time.NewTicker(pingPeriod)
	defer func() {
		log.Println("sender exiting")
		//		ticker.Stop()
		c.close()
	}()

	for {
		select {
		case msg, ok := <-c.toConn:
			if !ok {
				log.Println("not ok")
				return nil // TODO return err?
			}
			if err := c.trans.send(msg); err != nil {
				log.Println("err")
				log.Println(err)
				return err
			}
		//case <-ticker.C:
		//			if err := t.write(websocket.PingMessage, []byte{}); err != nil {
		//				return err
		//			}
		case <-c.tomb.Dying():
			return tomb.ErrDying
		}
	}
}

func (c *connection) receiver() error {
	defer func() {
		log.Println("receiver exiting")
		c.close()
	}()

	for {
		raw, err := c.trans.recv()
		if err != nil {
			log.Println(err)
			return err
		}
		select {
		case c.toGush <- raw:
		case <-c.tomb.Dying():
			return tomb.ErrDying
		}
	}
}
