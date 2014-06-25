package main

import (
	"github.com/dobrite/gusher/go/gusher"
	"log"
	"net/http"
	"os"
)

func main() {
	gmux := gusher.NewServeMux("/gusher")
	gmux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	gmux.HandleFunc("/", IndexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":"+port, gmux))
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/index.html")
}
