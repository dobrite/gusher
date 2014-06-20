package main

import (
	"github.com/dobrite/gusher/go/gusher"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	gmux := gusher.NewServeMux("/gusher")
	gmux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	gmux.HandleFunc("/", Index)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":"+port, gmux))
}

func Index(w http.ResponseWriter, req *http.Request) {
	//use http.ServeFile https://github.com/fzzy/sockjs-go/blob/master/examples/chat/chat.go
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
