package main

import (
	"github.com/dobrite/gusher/go/gusher"
	"log"
	"net/http"
	"os"
)

func setupLogger() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func getMux() *http.ServeMux {
	gmux := gusher.NewServeMux("/app", "tester")
	gmux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	gmux.HandleFunc("/", IndexHandler)
	return gmux
}

func main() {
	setupLogger()
	gmux := getMux()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	keyFile := "server.key"
	certFile := "server.crt"

	log.Println("server started on port " + port)
	//log.Fatal(http.ListenAndServe(":"+port, gmux))
	log.Fatal(http.ListenAndServeTLS(":"+port, certFile, keyFile, gmux))
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	http.ServeFile(w, req, "./public/index.html")
}
