package gusher

type message struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Data    string `json:"data"`
}
