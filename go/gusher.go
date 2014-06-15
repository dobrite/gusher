package main

import (
	"encoding/json"
	"github.com/igm/pubsub"
	"github.com/igm/sockjs-go/sockjs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Channel struct {
	pubsub.Publisher
}

var Channels map[string]*Channel

func main() {
	Channels = make(map[string]*Channel)
	Channels["test-channel"] = new(Channel)
	http.Handle("/gusher/", sockjs.NewHandler("/gusher", sockjs.DefaultOptions, gusherHandler))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/api/", API)
	http.HandleFunc("/", Index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func subscribe(session sockjs.Session, closedSession chan struct{}) {
	reader, _ := Channels["test-channel"].SubChannel(nil)
	for {
		select {
		case <-closedSession:
			log.Println("subscribe closed")
			return
		case message := <-reader:
			msg := message.(string)
			log.Println("Msg sent: " + msg)
			if err := session.Send(msg); err != nil {
				return
			}
		}
	}
}

func gusherHandler(session sockjs.Session) {
	log.Println("Client connected")
	//chat.Publish("[info] chatter joined")
	//defer chat.Publish("[info] chatter left")
	closedSession := make(chan struct{})
	go subscribe(session, closedSession)
	for {
		if raw, err := session.Recv(); err == nil {
			log.Println("Msg rec'd: " + raw)
			Channels["test-channel"].Publish(raw)
			continue
		}
		log.Println("Client disconnected")
		break
	}
	close(closedSession)
	log.Println("Session closed")
}

type Message struct {
	Channel string `json:"channel"`
	Event   string `json:"event"`
	Data    string `json:"data"`
}

func API(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		return
	}
	channel := req.PostFormValue("channel")
	event := req.PostFormValue("event")
	data := req.PostFormValue("data")
	m := Message{channel, event, data}
	payload, _ := json.Marshal(m)
	Channels[channel].Publish(string(payload))
}

func Index(w http.ResponseWriter, req *http.Request) {
	log.Println(req.URL.Path)
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	contents, err := ioutil.ReadFile("public/index.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(contents)
}
