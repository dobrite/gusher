package main

import (
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
	Channels["default"] = new(Channel)
	http.Handle("/gusher/", sockjs.NewHandler("/gusher", sockjs.DefaultOptions, gusherHandler))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/", Index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func subscribe(session sockjs.Session, closedSession chan struct{}) {
	reader, _ := Channels["default"].SubChannel(nil)
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
			Channels["default"].Publish(raw)
			continue
		}
		log.Println("Client disconnected")
		break
	}
	close(closedSession)
	log.Println("Session closed")
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
