package gusher

import (
	"github.com/gorilla/websocket"
	"gopkg.in/tomb.v2"
	"log"
	"sync"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type websocketTransport struct {
	transport
	ws *websocket.Conn
	//send    chan []byte
	closed  bool       //makes me sad
	closedm sync.Mutex //makes me sad
	tomb    tomb.Tomb
	toConn  chan string
	toGush  chan string
}

func newWebsocketTransport(ws *websocket.Conn, toConn chan string, toGush chan string) transport {
	t := &websocketTransport{
		ws:     ws,
		toGush: toGush,
		toConn: toConn,
	}
	t.ws.SetReadLimit(maxMessageSize)
	t.ws.SetReadDeadline(time.Now().Add(pongWait))
	t.ws.SetPongHandler(t.pongHandler)
	return t
}

func (t *websocketTransport) close_() {
	t.write(websocket.CloseMessage, []byte{})
	t.ws.Close()
}

func (t *websocketTransport) send_(msg string) error {
	//if err := t.write(websocket.TextMessage, []byte(message)); err != nil {
	return t.write(websocket.TextMessage, []byte(msg))
}

func (t *websocketTransport) recv_() (string, error) {
	_, message, err := t.ws.ReadMessage()
	// TODO _ above is msg type, prob err out on binary
	return string(message), err
}

func (t *websocketTransport) id_() string {
	return "websocket!" //TODO move this to session maybe?
}

func (t *websocketTransport) go_(f func() error) {
	t.tomb.Go(f)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (t *websocketTransport) pongHandler(string) error {
	t.ws.SetReadDeadline(time.Now().Add(pongWait))
	return nil
}

func (t *websocketTransport) close() error {
	t.closedm.Lock()
	defer t.closedm.Unlock()
	if !t.closed {
		t.closed = true
		log.Println("connection closed")
		t.close_()
		t.tomb.Kill(nil)
		close(t.toGush)
	}
	return t.tomb.Wait()
}

// write writes a message with the given message type and payload.
func (t *websocketTransport) write(mt int, payload []byte) error {
	t.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return t.ws.WriteMessage(mt, payload)
}

func (t *websocketTransport) sender_() error {
	//	ticker := time.NewTicker(pingPeriod)
	defer func() {
		//		ticker.Stop()
		t.close()
	}()

	for {
		select {
		case msg, ok := <-t.toConn:
			if !ok {
				return nil // TODO return err?
			}
			if err := t.send_(msg); err != nil {
				return err
			}
		//case <-ticker.C:
		//			if err := t.write(websocket.PingMessage, []byte{}); err != nil {
		//				return err
		//			}
		case <-t.tomb.Dying():
			return tomb.ErrDying
		}
	}
}

func (t *websocketTransport) receiver_() error {
	defer func() {
		t.close()
	}()

	for {
		raw, err := t.recv_()
		if err != nil {
			return err
		}
		select {
		case t.toGush <- raw:
		case <-t.tomb.Dying():
			return tomb.ErrDying
		}
	}
}
