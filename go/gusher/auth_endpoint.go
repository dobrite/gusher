package gusher

import (
	"fmt"
	"log"
	"net/http"
)

func (h *handler) auth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			w.WriteHeader(405)
			return
		}

		callback := req.PostFormValue("callback")
		socketId := req.PostFormValue("socket_id")
		channelName := req.PostFormValue("channel_name")
		channelData := req.PostFormValue("channel_data")

		authJSON := ""
		payload := ""

		if channelData != "" {
			authJSON = auth(socketId, channelName)
			payload = fmt.Sprintf("{\"auth\":\"%s\", \"channel_data\":\"%s\"}", authJSON, channelData)
		} else {
			authJSON = auth(socketId, channelName, channelData)
			payload = fmt.Sprintf("{\"auth\":\"%s\"}", authJSON)
		}

		header := "application/json"

		if callback != "" {
			header = "application/javascript"
			payload = fmt.Sprintf("%s(%s)", callback, payload)
		}

		w.Header().Set("Content-Type", header)
		log.Println(payload)
		fmt.Fprintf(w, payload)
		//set authTransport to 'ajax' (default)
		//POST to /pusher/auth w/ socket_id and channel_name
		//set authTransport to 'jsonp', also set authEndpoint (default to /pusher/auth)
		//JSONP to /pusher/auth w/ socket_id, channel_name and callback
		//render :text => params[:callback] + "(" + auth.to_json + ")", :content_type => 'application/javascript'
		//return if authorized application/json
		//{"auth":"278d425bdf160c739803:afaed3695da2ffd16931f457e338e6c9f2921fa133ce7dac49f529792be6304c","channel_data":"{\"user_id\":10,\"user_info\":{\"name\":\"Mr. Pusher\"}}"}
		//otherwise 403 Forbidden plain text
	})
}
