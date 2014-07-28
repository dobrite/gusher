package gusher

import (
	"encoding/json"
	"fmt"
)

type post struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

type message interface {
	//prob want to do json.Marshal
}

type messageEvent struct {
	Event   string          `json:"event"`
	Channel string          `json:"channel,omitempty"`
	Data    json.RawMessage `json:"data"`
}

// client -> gusher
type messageSubscribe struct {
	Channel     string `json:"channel"`
	Auth        string `json:"auth"`
	ChannelData string `json:"channel_data"`
}

// client -> gusher
type messageUnsubscribe struct {
	Channel string `json:"channel"`
}

// gusher -> client
type messageError struct {
	Message string `json:"message"`
	Code    uint16 `json:"code"`
}

type messagePing struct {
}

func buildMessageConnectionEstablished(id string) string {
	return fmt.Sprintf(`{"event": "pusher:connection_established", "data": "{\"socket_id\":\"%s\", \"activity_timeout\": 120}"}`, id)
}

func buildMessageSubscriptionSucceeded(channelName string) string {
	return fmt.Sprintf(`{"event": "pusher_internal:subscription_succeeded", "channel": "%s", "data": {}}`, channelName)
}

func buildMessagePresenceSubscriptionSucceeded(channelName string, ids []string, hash string, count int) string {
	// TODO hash is likely a JSON object
	return fmt.Sprintf(`{"event": "pusher_internal:subscription_succeeded", "channel": "%s", "data": {"presence": {"ids": %s, "hash": %s, "count": , %s}}}`, channelName, ids, hash, count)
}

func buildMessagePresenceMemberAdded(channelName string, userId string, userInfo string) string {
	// TODO userInfo likely a JSON object
	return fmt.Sprintf(`{"event": "pusher_internal:member_added", "channel": "%s", "data": {"user_id": "%s", "user_info": %s, }}`, channelName, userId, userInfo)
}

func buildMessagePresenceMemberRemoved(channelName string, userId string) string {
	return fmt.Sprintf(`{"event": "pusher_internal:member_added", "channel": "%s", "data": {"user_id": "%s"}}`, channelName, userId)
}

func buildMessagePong() string {
	return `{"event": "pusher:pong", "data":{}}`
}

func MessageUnmarshalJSON(b []byte) (msg message, err error) {
	var event messageEvent
	err = json.Unmarshal(b, &event)
	if err != nil {
		return
	}
	switch event.Event {
	case "pusher:subscribe":
		var msgSub messageSubscribe
		err = json.Unmarshal(event.Data, &msgSub)
		msg = msgSub
	case "pusher:unsubscribe":
		var msgUnsub messageUnsubscribe
		err = json.Unmarshal(event.Data, &msgUnsub)
		msg = msgUnsub
	case "pusher:ping":
		var msgPing messagePing
		err = json.Unmarshal(event.Data, &msgPing)
		msg = msgPing
	default:
		err = fmt.Errorf("%s is not a recognized event", event.Event)
	}
	return
}
