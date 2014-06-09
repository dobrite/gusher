package main

import (
	"github.com/igm/pubsub"
	"github.com/igm/sockjs-go/sockjs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var chat pubsub.Publisher

func main() {
	http.Handle("/echo/", sockjs.NewHandler("/echo", sockjs.DefaultOptions, echoHandler))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/", Index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func echoHandler(session sockjs.Session) {
	log.Println("Client connected")
	closedSession := make(chan struct{})
	chat.Publish("[info] chatter joined")
	defer chat.Publish("[info] chatter left")
	go func() {
		reader, _ := chat.SubChannel(nil)
		for {
			select {
			case <-closedSession:
				return
			case msg := <-reader:
				if err := session.Send(msg.(string)); err != nil {
					return
				}
			}
		}
	}()
	for {
		if msg, err := session.Recv(); err == nil {
			log.Println("Msg rec'd: " + msg)
			chat.Publish(msg)
			continue
		}
		log.Println("Client disconnected")
		break
	}
	close(closedSession)
	log.Println("Session closed")
}

func Index(w http.ResponseWriter, req *http.Request) {
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
