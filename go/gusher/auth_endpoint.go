package gusher

import (
	"fmt"
	"net/http"
)

func (h *handler) auth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" && req.Method != "GET" {
			w.WriteHeader(405)
			return
		}

		var socketId string
		var callback string
		var channelName string
		var channelData string
		var authJSON string
		var payload string

		if req.Method == "POST" {
			socketId = req.PostFormValue("socket_id")
			channelName = req.PostFormValue("channel_name")
			channelData = req.PostFormValue("channel_data")
		}

		if req.Method == "GET" {
			callback = req.FormValue("callback")
			socketId = req.FormValue("socket_id")
			channelName = req.FormValue("channel_name")
			channelData = req.FormValue("channel_data")
		}

		if channelData == "" {
			authJSON = auth(socketId, channelName)
			payload = fmt.Sprintf("{\"auth\":\"%s\"}", authJSON)
		} else {
			authJSON = auth(socketId, channelName, channelData)
			payload = fmt.Sprintf("{\"auth\":\"%s\", \"channel_data\":\"%s\"}", authJSON, channelData)
		}

		header := "application/json"

		if callback != "" {
			header = "application/javascript"
			payload = fmt.Sprintf("%s(%s)", callback, payload)
		}

		w.Header().Set("Content-Type", header)
		fmt.Fprintf(w, payload)

		// TODO if not authed 403 forbidden plantext
	})
}
