package gusher

import (
	"github.com/gorilla/websocket"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type websocketTransport struct {
	transport
	ws *websocket.Conn
}

func newWebsocketTransport(ws *websocket.Conn) transport {
	t := &websocketTransport{
		ws: ws,
	}
	t.ws.SetReadLimit(maxMessageSize)
	//t.ws.SetReadDeadline(time.Now().Add(pongWait))
	var zero time.Time
	t.ws.SetReadDeadline(zero)
	t.ws.SetPongHandler(t.pongHandler)
	return t
}

func (t *websocketTransport) close() {
	t.write(websocket.CloseMessage, []byte{})
	t.ws.Close()
}

func (t *websocketTransport) send(msg string) error {
	return t.write(websocket.TextMessage, []byte(msg))
}

func (t *websocketTransport) recv() (string, error) {
	_, message, err := t.ws.ReadMessage()
	return string(message), err
}

func (t *websocketTransport) pongHandler(string) error {
	var zero time.Time
	t.ws.SetReadDeadline(zero)
	//t.ws.SetReadDeadline(time.Now().Add(pongWait))
	return nil
}

// write writes a message with the given message type and payload.
func (t *websocketTransport) write(mt int, payload []byte) error {
	var zero time.Time
	//t.ws.SetWriteDeadline(time.Now().Add(writeWait))
	t.ws.SetWriteDeadline(zero)
	return t.ws.WriteMessage(mt, payload)
}
